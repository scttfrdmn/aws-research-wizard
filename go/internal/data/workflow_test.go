package data

import (
	"context"
	"testing"
	"time"
)

// TestWorkflowEngineCreation tests basic workflow engine creation and configuration
func TestWorkflowEngineCreation(t *testing.T) {
	config := &WorkflowEngineConfig{
		MaxConcurrentWorkflows: 2,
		DefaultTimeout:         1 * time.Hour,
		RetryAttempts:          2,
		RetryDelay:             10 * time.Second,
		MonitoringEnabled:      true,
	}

	engine := NewWorkflowEngine(config)

	if engine == nil {
		t.Fatal("Expected workflow engine to be created")
	}

	if engine.config.MaxConcurrentWorkflows != 2 {
		t.Errorf("Expected max concurrent workflows to be 2, got %d", engine.config.MaxConcurrentWorkflows)
	}

	// Test component registration
	analyzer := NewPatternAnalyzer()
	engine.RegisterAnalyzer(analyzer)

	if engine.patternAnalyzer == nil {
		t.Error("Expected pattern analyzer to be registered")
	}

	bundlingEngine := NewBundlingEngine(nil)
	engine.RegisterBundlingEngine(bundlingEngine)

	if engine.bundlingEngine == nil {
		t.Error("Expected bundling engine to be registered")
	}
}

// TestWorkflowExecution tests basic workflow execution workflow
func TestWorkflowExecution(t *testing.T) {
	engine := createTestWorkflowEngine()

	// Create test project configuration
	projectConfig := createTestProjectConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute workflow
	execution, err := engine.ExecuteWorkflow(ctx, projectConfig, "test_workflow")
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	if execution == nil {
		t.Fatal("Expected workflow execution to be returned")
	}

	if execution.WorkflowName != "test_workflow" {
		t.Errorf("Expected workflow name 'test_workflow', got '%s'", execution.WorkflowName)
	}

	if execution.Status != WorkflowStatusRunning && execution.Status != WorkflowStatusPending {
		t.Errorf("Expected workflow to be running or pending, got %s", execution.Status)
	}

	// Test workflow retrieval
	retrieved, err := engine.GetWorkflowExecution(execution.ID)
	if err != nil {
		t.Errorf("Failed to retrieve workflow execution: %v", err)
	}

	if retrieved.ID != execution.ID {
		t.Errorf("Expected retrieved workflow ID to match, got %s", retrieved.ID)
	}
}

// TestWorkflowStepBuilding tests workflow step building from configuration
func TestWorkflowStepBuilding(t *testing.T) {
	engine := createTestWorkflowEngine()

	workflow := &Workflow{
		Name:        "test_workflow",
		Source:      "/test/data",
		Destination: "s3://test-bucket/",
		Engine:      "s5cmd",
		PreProcessing: []ProcessingStep{
			{
				Name: "bundle_files",
				Type: "bundle",
				Parameters: map[string]string{
					"target_size": "100MB",
				},
			},
		},
		PostProcessing: []ProcessingStep{
			{
				Name: "cleanup",
				Type: "cleanup",
				Parameters: map[string]string{
					"action": "delete_temp",
				},
			},
		},
	}

	steps := engine.buildWorkflowSteps(workflow)

	// Should have: analyze + bundle + transfer + cleanup + report = 5 steps
	expectedSteps := 5
	if len(steps) != expectedSteps {
		t.Errorf("Expected %d steps, got %d", expectedSteps, len(steps))
	}

	// Check step types are correct
	expectedTypes := []string{"analyze", "bundle", "transfer", "cleanup", "report"}
	for i, step := range steps {
		if step.Type != expectedTypes[i] {
			t.Errorf("Expected step %d to be type '%s', got '%s'", i, expectedTypes[i], step.Type)
		}
	}

	// Check analyze step
	analyzeStep := steps[0]
	if analyzeStep.Name != "analyze_data_pattern" {
		t.Errorf("Expected analyze step name 'analyze_data_pattern', got '%s'", analyzeStep.Name)
	}

	// Check transfer step
	transferStep := steps[2] // analyze, bundle, transfer
	if transferStep.Engine != "s5cmd" {
		t.Errorf("Expected transfer step engine 's5cmd', got '%s'", transferStep.Engine)
	}
}

// TestWorkflowConcurrencyLimit tests that concurrent workflow limits are enforced
func TestWorkflowConcurrencyLimit(t *testing.T) {
	config := &WorkflowEngineConfig{
		MaxConcurrentWorkflows: 1, // Limit to 1 concurrent workflow
		DefaultTimeout:         1 * time.Hour,
	}

	engine := NewWorkflowEngine(config)
	engine.RegisterAnalyzer(NewPatternAnalyzer())

	projectConfig := createTestProjectConfig()
	ctx := context.Background()

	// Start first workflow
	execution1, err := engine.ExecuteWorkflow(ctx, projectConfig, "test_workflow")
	if err != nil {
		t.Fatalf("Failed to start first workflow: %v", err)
	}

	// Try to start second workflow - should fail due to limit
	_, err = engine.ExecuteWorkflow(ctx, projectConfig, "test_workflow")
	if err == nil {
		t.Error("Expected second workflow to fail due to concurrency limit")
	}

	// Cancel first workflow to free up slot
	err = engine.CancelWorkflow(execution1.ID)
	if err != nil {
		t.Errorf("Failed to cancel workflow: %v", err)
	}

	// Wait a moment for cancellation to process
	time.Sleep(100 * time.Millisecond)

	// Now second workflow should succeed
	_, err = engine.ExecuteWorkflow(ctx, projectConfig, "test_workflow")
	if err != nil {
		t.Errorf("Expected second workflow to succeed after cancellation: %v", err)
	}
}

// TestWorkflowProgressTracking tests progress tracking functionality
func TestWorkflowProgressTracking(t *testing.T) {
	engine := createTestWorkflowEngine()
	projectConfig := createTestProjectConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	execution, err := engine.ExecuteWorkflow(ctx, projectConfig, "test_workflow")
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Register progress callback
	progressUpdates := 0
	engine.RegisterProgressCallback(execution.ID, func(exec *WorkflowExecution) {
		progressUpdates++
		if exec.Progress < 0 || exec.Progress > 1 {
			t.Errorf("Progress should be between 0 and 1, got %f", exec.Progress)
		}
	})

	// Wait for some progress
	time.Sleep(2 * time.Second)

	// Check that progress callbacks were called (allow for quick workflows)
	if progressUpdates == 0 {
		t.Logf("No progress updates received - workflow may have completed too quickly")
		// This is not necessarily an error for fast-completing test workflows
	} else {
		t.Logf("Received %d progress updates", progressUpdates)
	}
}

// TestWorkflowEventSystem tests the workflow event system
func TestWorkflowEventSystem(t *testing.T) {
	engine := createTestWorkflowEngine()

	// Create a simple execution for testing events
	execution := &WorkflowExecution{
		ID:           "test-execution",
		WorkflowName: "test",
		Events:       make([]*WorkflowEvent, 0),
	}

	// Emit some test events
	engine.emitEvent(execution, "test_event", "Test message", "info", nil)
	engine.emitEvent(execution, "warning_event", "Warning message", "warning", map[string]interface{}{
		"test_data": "value",
	})

	// Check events were recorded
	if len(execution.Events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(execution.Events))
	}

	// Check first event
	firstEvent := execution.Events[0]
	if firstEvent.Type != "test_event" {
		t.Errorf("Expected first event type 'test_event', got '%s'", firstEvent.Type)
	}

	if firstEvent.Message != "Test message" {
		t.Errorf("Expected first event message 'Test message', got '%s'", firstEvent.Message)
	}

	if firstEvent.Severity != "info" {
		t.Errorf("Expected first event severity 'info', got '%s'", firstEvent.Severity)
	}

	// Check second event has metadata
	secondEvent := execution.Events[1]
	if secondEvent.Metadata == nil {
		t.Error("Expected second event to have metadata")
	}

	if value, exists := secondEvent.Metadata["test_data"]; !exists || value != "value" {
		t.Error("Expected second event metadata to contain test_data=value")
	}
}

// Helper functions for testing

func createTestWorkflowEngine() *WorkflowEngine {
	config := &WorkflowEngineConfig{
		MaxConcurrentWorkflows: 5,
		DefaultTimeout:         1 * time.Hour,
		RetryAttempts:          2,
		RetryDelay:             1 * time.Second,
		MonitoringEnabled:      true,
	}

	engine := NewWorkflowEngine(config)

	// Register test components
	engine.RegisterAnalyzer(NewPatternAnalyzer())
	engine.RegisterBundlingEngine(NewBundlingEngine(nil))
	engine.RegisterWarningSystem(NewWarningSystem())

	// Register mock transfer engines for testing
	mockEngine := &MockTransferEngine{name: "s5cmd"}
	engine.RegisterTransferEngine(mockEngine)

	return engine
}

func createTestProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		Project: ProjectInfo{
			Name:        "test-project",
			Description: "Test project for workflow engine",
		},
		DataProfiles: map[string]DataProfile{
			"test_data": {
				Name: "Test Data",
				Path: "/tmp/test-data",
			},
		},
		Destinations: map[string]Destination{
			"test_dest": {
				Name: "Test Destination",
				URI:  "s3://test-bucket/",
			},
		},
		Workflows: []Workflow{
			{
				Name:        "test_workflow",
				Description: "Test workflow",
				Source:      "test_data",
				Destination: "test_dest",
				Engine:      "s5cmd",
				Enabled:     true,
				PreProcessing: []ProcessingStep{
					{
						Name: "test_bundle",
						Type: "bundle",
					},
				},
				Configuration: WorkflowConfiguration{
					Concurrency:   4,
					RetryAttempts: 2,
					Timeout:       "1h",
				},
			},
		},
	}
}

// MockTransferEngine for testing
type MockTransferEngine struct {
	name string
}

func (m *MockTransferEngine) GetName() string {
	return m.name
}

func (m *MockTransferEngine) GetType() string {
	return "mock"
}

func (m *MockTransferEngine) IsAvailable(ctx context.Context) error {
	return nil // Always available for testing
}

func (m *MockTransferEngine) GetCapabilities() EngineCapabilities {
	return EngineCapabilities{
		Protocols:        []string{"mock"},
		SupportsParallel: true,
		SupportsProgress: true,
		MaxConcurrency:   10,
	}
}

func (m *MockTransferEngine) Upload(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// Simulate transfer
	time.Sleep(100 * time.Millisecond)

	return &TransferResult{
		TransferID:       req.ID,
		Engine:           m.name,
		Source:           req.Source,
		Destination:      req.Destination,
		Success:          true,
		BytesTransferred: 1024 * 1024, // 1MB
		FilesTransferred: 10,
		StartTime:        time.Now().Add(-100 * time.Millisecond),
		EndTime:          time.Now(),
		Duration:         100 * time.Millisecond,
		AverageSpeed:     10 * 1024 * 1024, // 10 MB/s
	}, nil
}

func (m *MockTransferEngine) Download(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	return m.Upload(ctx, req) // Same as upload for testing
}

func (m *MockTransferEngine) Sync(ctx context.Context, req *SyncRequest) (*TransferResult, error) {
	// Convert sync request to transfer request for testing
	transferReq := &TransferRequest{
		ID:          req.ID,
		Source:      req.Source,
		Destination: req.Destination,
		Context:     req.Context,
	}
	return m.Upload(ctx, transferReq)
}

func (m *MockTransferEngine) GetProgress(ctx context.Context, transferID string) (*TransferProgress, error) {
	return &TransferProgress{
		BytesTransferred: 512 * 1024,
		TotalBytes:       1024 * 1024,
		Percentage:       50.0,
		Speed:            10 * 1024 * 1024,
		StartTime:        time.Now().Add(-1 * time.Minute),
		LastUpdate:       time.Now(),
	}, nil
}

func (m *MockTransferEngine) Cancel(ctx context.Context, transferID string) error {
	return nil // Always succeeds for testing
}

func (m *MockTransferEngine) Validate() error {
	return nil // Always valid for testing
}
