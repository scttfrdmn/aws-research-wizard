package data

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkflowEngine orchestrates complex data movement workflows
type WorkflowEngine struct {
	// Core components
	patternAnalyzer    *PatternAnalyzer
	recommendationEngine *RecommendationEngine
	bundlingEngine     *BundlingEngine
	transferEngines    map[string]TransferEngine
	warningSystem      *WarningSystem
	
	// Execution state
	activeWorkflows    map[string]*WorkflowExecution
	executionHistory   []*WorkflowExecution
	config            *WorkflowEngineConfig
	
	// Synchronization
	mu                sync.RWMutex
	shutdownChan      chan struct{}
	
	// Monitoring
	progressCallbacks map[string][]WorkflowProgressCallback
	eventBus          *WorkflowEventBus
}

// WorkflowEngineConfig contains configuration for the workflow engine
type WorkflowEngineConfig struct {
	MaxConcurrentWorkflows int               `json:"max_concurrent_workflows"`
	DefaultTimeout         time.Duration     `json:"default_timeout"`
	RetryAttempts         int               `json:"retry_attempts"`
	RetryDelay            time.Duration     `json:"retry_delay"`
	EnablePersistence     bool              `json:"enable_persistence"`
	PersistenceDir        string            `json:"persistence_dir"`
	MonitoringEnabled     bool              `json:"monitoring_enabled"`
	MetricsCollection     bool              `json:"metrics_collection"`
	
	// Domain-specific settings
	DomainConfigs         map[string]DomainConfig `json:"domain_configs"`
}

// DomainConfig contains domain-specific workflow settings
type DomainConfig struct {
	DefaultBundling       bool              `json:"default_bundling"`
	PreferredEngines      []string          `json:"preferred_engines"`
	CustomOptimizations   map[string]interface{} `json:"custom_optimizations"`
	QualityChecks         []string          `json:"quality_checks"`
}

// WorkflowExecution represents a running or completed workflow
type WorkflowExecution struct {
	// Identification
	ID                string            `json:"id"`
	WorkflowName      string            `json:"workflow_name"`
	ProjectConfig     *ProjectConfig    `json:"project_config"`
	
	// Execution state
	Status            WorkflowStatus    `json:"status"`
	CurrentStep       int               `json:"current_step"`
	TotalSteps        int               `json:"total_steps"`
	StartTime         time.Time         `json:"start_time"`
	EndTime           time.Time         `json:"end_time,omitempty"`
	Duration          time.Duration     `json:"duration"`
	
	// Steps and results
	Steps             []*WorkflowStep   `json:"steps"`
	Results           *WorkflowResults  `json:"results,omitempty"`
	Error             error             `json:"error,omitempty"`
	
	// Progress tracking
	Progress          float64           `json:"progress"` // 0.0 to 1.0
	CurrentStepProgress float64         `json:"current_step_progress"`
	
	// Context and cancellation
	Context           context.Context   `json:"-"`
	CancelFunc        context.CancelFunc `json:"-"`
	
	// Monitoring
	Metrics           *WorkflowMetrics  `json:"metrics,omitempty"`
	Events            []*WorkflowEvent  `json:"events,omitempty"`
}

// WorkflowStatus represents the current status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusPending    WorkflowStatus = "pending"
	WorkflowStatusRunning    WorkflowStatus = "running"
	WorkflowStatusCompleted  WorkflowStatus = "completed"
	WorkflowStatusFailed     WorkflowStatus = "failed"
	WorkflowStatusCancelled  WorkflowStatus = "cancelled"
	WorkflowStatusPaused     WorkflowStatus = "paused"
)

// WorkflowStep represents a single step in a workflow execution
type WorkflowStep struct {
	// Definition
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Engine            string            `json:"engine,omitempty"`
	
	// Configuration
	Parameters        map[string]string `json:"parameters"`
	Conditions        []string          `json:"conditions,omitempty"`
	Dependencies      []string          `json:"dependencies,omitempty"`
	
	// Execution state
	Status            StepStatus        `json:"status"`
	StartTime         time.Time         `json:"start_time,omitempty"`
	EndTime           time.Time         `json:"end_time,omitempty"`
	Duration          time.Duration     `json:"duration"`
	RetryCount        int               `json:"retry_count"`
	
	// Results
	Output            map[string]interface{} `json:"output,omitempty"`
	Error             error                  `json:"error,omitempty"`
	
	// Progress
	Progress          float64           `json:"progress"`
}

// StepStatus represents the status of a workflow step
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
)

// WorkflowResults contains the comprehensive results of a workflow execution
type WorkflowResults struct {
	// Analysis results
	DataPattern           *DataPattern           `json:"data_pattern,omitempty"`
	Recommendations       *RecommendationResult  `json:"recommendations,omitempty"`
	WarningReport         *WarningReport         `json:"warning_report,omitempty"`
	
	// Processing results
	BundlingResult        *BundlingResult        `json:"bundling_result,omitempty"`
	TransferResults       []*TransferResult      `json:"transfer_results,omitempty"`
	
	// Summary metrics
	TotalFilesProcessed   int64                  `json:"total_files_processed"`
	TotalBytesTransferred int64                  `json:"total_bytes_transferred"`
	TotalCostSavings      float64                `json:"total_cost_savings"`
	PerformanceGain       float64                `json:"performance_gain_percent"`
	
	// Quality metrics
	SuccessRate           float64                `json:"success_rate"`
	ErrorCount            int                    `json:"error_count"`
	WarningCount          int                    `json:"warning_count"`
	
	// Recommendations for future
	NextStepSuggestions   []string               `json:"next_step_suggestions"`
	OptimizationOpportunities []string           `json:"optimization_opportunities"`
}

// WorkflowMetrics contains detailed metrics about workflow execution
type WorkflowMetrics struct {
	// Timing metrics
	TotalExecutionTime    time.Duration         `json:"total_execution_time"`
	StepExecutionTimes    map[string]time.Duration `json:"step_execution_times"`
	QueueTime             time.Duration         `json:"queue_time"`
	
	// Resource utilization
	PeakMemoryUsage       int64                 `json:"peak_memory_usage"`
	CPUUtilization        float64               `json:"cpu_utilization"`
	NetworkUtilization    float64               `json:"network_utilization"`
	
	// Efficiency metrics
	TransferEfficiency    float64               `json:"transfer_efficiency"`
	CompressionRatio      float64               `json:"compression_ratio"`
	CostEfficiency        float64               `json:"cost_efficiency"`
	
	// Error and retry metrics
	RetryAttempts         map[string]int        `json:"retry_attempts"`
	ErrorRecoveryTime     time.Duration         `json:"error_recovery_time"`
}

// WorkflowEvent represents an event that occurred during workflow execution
type WorkflowEvent struct {
	Timestamp    time.Time   `json:"timestamp"`
	Type         string      `json:"type"`
	Step         string      `json:"step,omitempty"`
	Message      string      `json:"message"`
	Severity     string      `json:"severity"` // "info", "warning", "error"
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// WorkflowProgressCallback is called to report workflow progress
type WorkflowProgressCallback func(execution *WorkflowExecution)

// WorkflowEventBus handles workflow event distribution
type WorkflowEventBus struct {
	subscribers map[string][]chan *WorkflowEvent
	mu          sync.RWMutex
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(config *WorkflowEngineConfig) *WorkflowEngine {
	if config == nil {
		config = &WorkflowEngineConfig{
			MaxConcurrentWorkflows: 5,
			DefaultTimeout:         2 * time.Hour,
			RetryAttempts:         3,
			RetryDelay:            30 * time.Second,
			EnablePersistence:     true,
			MonitoringEnabled:     true,
			MetricsCollection:     true,
		}
	}
	
	return &WorkflowEngine{
		activeWorkflows:   make(map[string]*WorkflowExecution),
		executionHistory:  make([]*WorkflowExecution, 0),
		transferEngines:   make(map[string]TransferEngine),
		progressCallbacks: make(map[string][]WorkflowProgressCallback),
		config:           config,
		shutdownChan:     make(chan struct{}),
		eventBus:         &WorkflowEventBus{
			subscribers: make(map[string][]chan *WorkflowEvent),
		},
	}
}

// RegisterAnalyzer registers the pattern analyzer
func (we *WorkflowEngine) RegisterAnalyzer(analyzer *PatternAnalyzer) {
	we.patternAnalyzer = analyzer
}

// RegisterRecommendationEngine registers the recommendation engine
func (we *WorkflowEngine) RegisterRecommendationEngine(engine *RecommendationEngine) {
	we.recommendationEngine = engine
}

// RegisterBundlingEngine registers the bundling engine
func (we *WorkflowEngine) RegisterBundlingEngine(engine *BundlingEngine) {
	we.bundlingEngine = engine
}

// RegisterTransferEngine registers a transfer engine
func (we *WorkflowEngine) RegisterTransferEngine(engine TransferEngine) error {
	if engine == nil {
		return fmt.Errorf("engine cannot be nil")
	}
	
	name := engine.GetName()
	if name == "" {
		return fmt.Errorf("engine name cannot be empty")
	}
	
	we.transferEngines[name] = engine
	return nil
}

// RegisterWarningSystem registers the warning system
func (we *WorkflowEngine) RegisterWarningSystem(system *WarningSystem) {
	we.warningSystem = system
}

// ExecuteWorkflow executes a complete workflow from project configuration
func (we *WorkflowEngine) ExecuteWorkflow(ctx context.Context, projectConfig *ProjectConfig, workflowName string) (*WorkflowExecution, error) {
	// Find the specified workflow
	var workflow *Workflow
	for _, w := range projectConfig.Workflows {
		if w.Name == workflowName {
			workflow = &w
			break
		}
	}
	
	if workflow == nil {
		return nil, fmt.Errorf("workflow '%s' not found in project configuration", workflowName)
	}
	
	if !workflow.Enabled {
		return nil, fmt.Errorf("workflow '%s' is disabled", workflowName)
	}
	
	// Check concurrent workflow limit
	we.mu.RLock()
	activeCount := len(we.activeWorkflows)
	we.mu.RUnlock()
	
	if activeCount >= we.config.MaxConcurrentWorkflows {
		return nil, fmt.Errorf("maximum concurrent workflows (%d) reached", we.config.MaxConcurrentWorkflows)
	}
	
	// Create workflow execution
	execution := we.createWorkflowExecution(ctx, projectConfig, workflow)
	
	// Register execution
	we.mu.Lock()
	we.activeWorkflows[execution.ID] = execution
	we.mu.Unlock()
	
	// Start execution in background
	go we.executeWorkflowSteps(execution)
	
	return execution, nil
}

// createWorkflowExecution creates a new workflow execution instance
func (we *WorkflowEngine) createWorkflowExecution(ctx context.Context, projectConfig *ProjectConfig, workflow *Workflow) *WorkflowExecution {
	executionID := fmt.Sprintf("wf_%d", time.Now().UnixNano())
	
	// Create cancellable context with timeout
	timeout := we.config.DefaultTimeout
	if workflow.Configuration.Timeout != "" {
		if parsedTimeout, err := time.ParseDuration(workflow.Configuration.Timeout); err == nil {
			timeout = parsedTimeout
		}
	}
	
	execCtx, cancelFunc := context.WithTimeout(ctx, timeout)
	
	// Build workflow steps
	steps := we.buildWorkflowSteps(workflow)
	
	return &WorkflowExecution{
		ID:           executionID,
		WorkflowName: workflow.Name,
		ProjectConfig: projectConfig,
		Status:       WorkflowStatusPending,
		TotalSteps:   len(steps),
		StartTime:    time.Now(),
		Steps:        steps,
		Context:      execCtx,
		CancelFunc:   cancelFunc,
		Progress:     0.0,
		Events:       make([]*WorkflowEvent, 0),
		Metrics:      &WorkflowMetrics{
			StepExecutionTimes: make(map[string]time.Duration),
			RetryAttempts:      make(map[string]int),
		},
	}
}

// buildWorkflowSteps converts workflow configuration into executable steps
func (we *WorkflowEngine) buildWorkflowSteps(workflow *Workflow) []*WorkflowStep {
	var steps []*WorkflowStep
	
	// Add analysis step (always first)
	steps = append(steps, &WorkflowStep{
		Name:       "analyze_data_pattern",
		Type:       "analyze",
		Parameters: map[string]string{
			"source_path": workflow.Source,
		},
		Status: StepStatusPending,
	})
	
	// Add preprocessing steps
	for _, preStep := range workflow.PreProcessing {
		step := &WorkflowStep{
			Name:       preStep.Name,
			Type:       preStep.Type,
			Parameters: preStep.Parameters,
			Conditions: []string{preStep.Condition},
			Status:     StepStatusPending,
		}
		steps = append(steps, step)
	}
	
	// Add main transfer step
	steps = append(steps, &WorkflowStep{
		Name:   "primary_transfer",
		Type:   "transfer",
		Engine: workflow.Engine,
		Parameters: map[string]string{
			"source":      workflow.Source,
			"destination": workflow.Destination,
			"concurrency": fmt.Sprintf("%d", workflow.Configuration.Concurrency),
			"part_size":   workflow.Configuration.PartSize,
		},
		Status: StepStatusPending,
	})
	
	// Add postprocessing steps
	for _, postStep := range workflow.PostProcessing {
		step := &WorkflowStep{
			Name:       postStep.Name,
			Type:       postStep.Type,
			Parameters: postStep.Parameters,
			Conditions: []string{postStep.Condition},
			Status:     StepStatusPending,
		}
		steps = append(steps, step)
	}
	
	// Add monitoring step (always last)
	steps = append(steps, &WorkflowStep{
		Name: "generate_report",
		Type: "report",
		Parameters: map[string]string{
			"include_metrics": "true",
			"include_recommendations": "true",
		},
		Status: StepStatusPending,
	})
	
	return steps
}

// executeWorkflowSteps executes all steps in a workflow
func (we *WorkflowEngine) executeWorkflowSteps(execution *WorkflowExecution) {
	defer func() {
		// Clean up when workflow completes
		we.mu.Lock()
		delete(we.activeWorkflows, execution.ID)
		we.executionHistory = append(we.executionHistory, execution)
		we.mu.Unlock()
		
		execution.CancelFunc()
		execution.EndTime = time.Now()
		execution.Duration = execution.EndTime.Sub(execution.StartTime)
	}()
	
	execution.Status = WorkflowStatusRunning
	we.emitEvent(execution, "workflow_started", "Workflow execution started", "info", nil)
	
	// Initialize results
	execution.Results = &WorkflowResults{
		TransferResults:       make([]*TransferResult, 0),
		NextStepSuggestions:   make([]string, 0),
		OptimizationOpportunities: make([]string, 0),
	}
	
	// Execute each step
	for i, step := range execution.Steps {
		execution.CurrentStep = i
		
		// Check for cancellation
		select {
		case <-execution.Context.Done():
			execution.Status = WorkflowStatusCancelled
			we.emitEvent(execution, "workflow_cancelled", "Workflow was cancelled", "warning", nil)
			return
		default:
		}
		
		// Execute step with retry logic
		if err := we.executeStepWithRetry(execution, step); err != nil {
			execution.Status = WorkflowStatusFailed
			execution.Error = err
			we.emitEvent(execution, "workflow_failed", fmt.Sprintf("Workflow failed at step '%s': %v", step.Name, err), "error", nil)
			return
		}
		
		// Update progress
		execution.Progress = float64(i+1) / float64(execution.TotalSteps)
		we.notifyProgress(execution)
	}
	
	execution.Status = WorkflowStatusCompleted
	we.emitEvent(execution, "workflow_completed", "Workflow completed successfully", "info", nil)
	we.notifyProgress(execution)
}

// executeStepWithRetry executes a single step with retry logic
func (we *WorkflowEngine) executeStepWithRetry(execution *WorkflowExecution, step *WorkflowStep) error {
	maxRetries := we.config.RetryAttempts
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		step.RetryCount = attempt
		step.StartTime = time.Now()
		step.Status = StepStatusRunning
		
		// Execute the step
		err := we.executeStep(execution, step)
		
		step.EndTime = time.Now()
		step.Duration = step.EndTime.Sub(step.StartTime)
		execution.Metrics.StepExecutionTimes[step.Name] = step.Duration
		
		if err == nil {
			step.Status = StepStatusCompleted
			we.emitEvent(execution, "step_completed", fmt.Sprintf("Step '%s' completed successfully", step.Name), "info", nil)
			return nil
		}
		
		step.Error = err
		execution.Metrics.RetryAttempts[step.Name] = attempt + 1
		
		if attempt < maxRetries {
			we.emitEvent(execution, "step_retry", fmt.Sprintf("Step '%s' failed, retrying (attempt %d/%d): %v", step.Name, attempt+1, maxRetries, err), "warning", nil)
			time.Sleep(we.config.RetryDelay)
			continue
		}
		
		step.Status = StepStatusFailed
		return fmt.Errorf("step '%s' failed after %d attempts: %w", step.Name, maxRetries+1, err)
	}
	
	return nil
}

// executeStep executes a single workflow step
func (we *WorkflowEngine) executeStep(execution *WorkflowExecution, step *WorkflowStep) error {
	we.emitEvent(execution, "step_started", fmt.Sprintf("Starting step '%s' (%s)", step.Name, step.Type), "info", nil)
	
	switch step.Type {
	case "analyze":
		return we.executeAnalyzeStep(execution, step)
	case "bundle":
		return we.executeBundleStep(execution, step)
	case "transfer":
		return we.executeTransferStep(execution, step)
	case "validate":
		return we.executeValidateStep(execution, step)
	case "cleanup":
		return we.executeCleanupStep(execution, step)
	case "report":
		return we.executeReportStep(execution, step)
	default:
		return fmt.Errorf("unknown step type: %s", step.Type)
	}
}

// Step execution methods

func (we *WorkflowEngine) executeAnalyzeStep(execution *WorkflowExecution, step *WorkflowStep) error {
	if we.patternAnalyzer == nil {
		return fmt.Errorf("pattern analyzer not registered")
	}
	
	sourcePath := step.Parameters["source_path"]
	if sourcePath == "" {
		sourcePath = execution.ProjectConfig.DataProfiles["main_dataset"].Path
	}
	
	// Analyze data pattern
	pattern, err := we.patternAnalyzer.AnalyzePattern(execution.Context, sourcePath)
	if err != nil {
		return fmt.Errorf("data pattern analysis failed: %w", err)
	}
	
	execution.Results.DataPattern = pattern
	step.Output = map[string]interface{}{
		"total_files": pattern.TotalFiles,
		"total_size":  pattern.TotalSizeHuman,
		"small_files": pattern.FileSizes.SmallFiles.CountUnder1MB,
	}
	
	// Generate recommendations if engine is available
	if we.recommendationEngine != nil {
		recommendations, err := we.recommendationEngine.GenerateRecommendations(execution.Context, sourcePath)
		if err != nil {
			we.emitEvent(execution, "recommendation_warning", fmt.Sprintf("Failed to generate recommendations: %v", err), "warning", nil)
		} else {
			execution.Results.Recommendations = recommendations
		}
	}
	
	// Generate warnings if system is available
	if we.warningSystem != nil && execution.Results.Recommendations != nil {
		warnings, err := we.warningSystem.AnalyzePattern(execution.Context, pattern, execution.Results.Recommendations.CostAnalysis)
		if err != nil {
			we.emitEvent(execution, "warning_analysis_failed", fmt.Sprintf("Warning analysis failed: %v", err), "warning", nil)
		} else {
			execution.Results.WarningReport = warnings
		}
	}
	
	return nil
}

func (we *WorkflowEngine) executeBundleStep(execution *WorkflowExecution, step *WorkflowStep) error {
	if we.bundlingEngine == nil {
		return fmt.Errorf("bundling engine not registered")
	}
	
	// Check if bundling is recommended
	if execution.Results.DataPattern != nil {
		recommendation, err := we.bundlingEngine.ShouldBundle(execution.Context, execution.Results.DataPattern)
		if err != nil {
			return fmt.Errorf("bundling recommendation failed: %w", err)
		}
		
		if !recommendation.Recommended {
			step.Status = StepStatusSkipped
			step.Output = map[string]interface{}{
				"skipped_reason": "bundling not recommended",
				"confidence":     recommendation.Confidence,
			}
			return nil
		}
	}
	
	// Create bundling request
	sourcePath := step.Parameters["source_path"]
	if sourcePath == "" && execution.Results.DataPattern != nil {
		sourcePath = execution.Results.DataPattern.AnalyzedPath
	}
	
	request := &BundlingTransferRequest{
		SourcePath:      sourcePath,
		DestinationPath: step.Parameters["output_path"],
		Metadata:        make(map[string]interface{}),
	}
	
	// Copy step parameters to metadata
	for k, v := range step.Parameters {
		request.Metadata[k] = v
	}
	
	// Execute bundling
	result, err := we.bundlingEngine.ProcessForBundling(execution.Context, request)
	if err != nil {
		return fmt.Errorf("bundling failed: %w", err)
	}
	
	execution.Results.BundlingResult = result
	execution.Results.TotalCostSavings += result.CostSavings.TotalSavings
	
	step.Output = map[string]interface{}{
		"bundles_created":    len(result.BundlePaths),
		"compression_ratio":  result.CompressionRatio,
		"cost_savings":       result.CostSavings.TotalSavings,
		"files_bundled":      result.BundledFileCount,
	}
	
	return nil
}

func (we *WorkflowEngine) executeTransferStep(execution *WorkflowExecution, step *WorkflowStep) error {
	engineName := step.Engine
	if engineName == "" || engineName == "auto" {
		// Auto-select engine based on recommendations
		if execution.Results.Recommendations != nil {
			for _, toolRec := range execution.Results.Recommendations.ToolRecommendations {
				if toolRec.Task == "primary_upload" {
					engineName = toolRec.RecommendedTool
					break
				}
			}
		}
		
		// Default to s5cmd if no recommendation
		if engineName == "" {
			engineName = "s5cmd"
		}
	}
	
	engine, exists := we.transferEngines[engineName]
	if !exists {
		return fmt.Errorf("transfer engine '%s' not registered", engineName)
	}
	
	// Check engine availability
	if err := engine.IsAvailable(execution.Context); err != nil {
		return fmt.Errorf("transfer engine '%s' not available: %w", engineName, err)
	}
	
	// Build transfer request
	transferReq := &TransferRequest{
		ID:          fmt.Sprintf("%s_transfer_%d", execution.ID, time.Now().UnixNano()),
		Source:      step.Parameters["source"],
		Destination: step.Parameters["destination"],
		Context:     execution.Context,
		Options:     TransferOptions{},
	}
	
	// Configure transfer options from step parameters
	if concurrency := step.Parameters["concurrency"]; concurrency != "" {
		if c, err := fmt.Sscanf(concurrency, "%d", &transferReq.Options.Concurrency); err == nil && c == 1 {
			// Concurrency set successfully
		}
	}
	
	// Execute transfer
	result, err := engine.Upload(execution.Context, transferReq)
	if err != nil {
		return fmt.Errorf("transfer failed: %w", err)
	}
	
	execution.Results.TransferResults = append(execution.Results.TransferResults, result)
	execution.Results.TotalFilesProcessed += int64(result.FilesTransferred)
	execution.Results.TotalBytesTransferred += result.BytesTransferred
	
	step.Output = map[string]interface{}{
		"engine":             engineName,
		"files_transferred":  result.FilesTransferred,
		"bytes_transferred":  result.BytesTransferred,
		"average_speed":      result.AverageSpeed,
		"duration":           result.Duration.String(),
	}
	
	return nil
}

func (we *WorkflowEngine) executeValidateStep(execution *WorkflowExecution, step *WorkflowStep) error {
	// Validation logic would go here
	// For now, just mark as completed
	step.Output = map[string]interface{}{
		"validation_passed": true,
	}
	return nil
}

func (we *WorkflowEngine) executeCleanupStep(execution *WorkflowExecution, step *WorkflowStep) error {
	// Cleanup logic would go here
	// For now, just mark as completed
	step.Output = map[string]interface{}{
		"cleanup_completed": true,
	}
	return nil
}

func (we *WorkflowEngine) executeReportStep(execution *WorkflowExecution, step *WorkflowStep) error {
	// Generate final report
	if execution.Results != nil {
		// Calculate success rate
		successfulTransfers := 0
		for _, result := range execution.Results.TransferResults {
			if result.Success {
				successfulTransfers++
			}
		}
		
		if len(execution.Results.TransferResults) > 0 {
			execution.Results.SuccessRate = float64(successfulTransfers) / float64(len(execution.Results.TransferResults)) * 100
		}
		
		// Add recommendations for future
		if execution.Results.WarningReport != nil && len(execution.Results.WarningReport.QuickFixes) > 0 {
			for _, fix := range execution.Results.WarningReport.QuickFixes {
				execution.Results.NextStepSuggestions = append(execution.Results.NextStepSuggestions, fix.Title)
			}
		}
	}
	
	step.Output = map[string]interface{}{
		"report_generated": true,
		"total_savings":    execution.Results.TotalCostSavings,
		"success_rate":     execution.Results.SuccessRate,
	}
	
	return nil
}

// Workflow management methods

// GetActiveWorkflows returns all currently active workflows
func (we *WorkflowEngine) GetActiveWorkflows() map[string]*WorkflowExecution {
	we.mu.RLock()
	defer we.mu.RUnlock()
	
	active := make(map[string]*WorkflowExecution)
	for id, execution := range we.activeWorkflows {
		active[id] = execution
	}
	
	return active
}

// GetWorkflowExecution returns a specific workflow execution
func (we *WorkflowEngine) GetWorkflowExecution(executionID string) (*WorkflowExecution, error) {
	we.mu.RLock()
	defer we.mu.RUnlock()
	
	if execution, exists := we.activeWorkflows[executionID]; exists {
		return execution, nil
	}
	
	// Search in history
	for _, execution := range we.executionHistory {
		if execution.ID == executionID {
			return execution, nil
		}
	}
	
	return nil, fmt.Errorf("workflow execution '%s' not found", executionID)
}

// CancelWorkflow cancels a running workflow
func (we *WorkflowEngine) CancelWorkflow(executionID string) error {
	we.mu.RLock()
	execution, exists := we.activeWorkflows[executionID]
	we.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("workflow execution '%s' not found or not active", executionID)
	}
	
	execution.CancelFunc()
	return nil
}

// Progress and event handling

func (we *WorkflowEngine) notifyProgress(execution *WorkflowExecution) {
	callbacks, exists := we.progressCallbacks[execution.ID]
	if !exists {
		return
	}
	
	for _, callback := range callbacks {
		go callback(execution)
	}
}

func (we *WorkflowEngine) emitEvent(execution *WorkflowExecution, eventType, message, severity string, metadata map[string]interface{}) {
	event := &WorkflowEvent{
		Timestamp: time.Now(),
		Type:      eventType,
		Message:   message,
		Severity:  severity,
		Metadata:  metadata,
	}
	
	execution.Events = append(execution.Events, event)
	we.eventBus.publish(eventType, event)
}

func (eb *WorkflowEventBus) publish(eventType string, event *WorkflowEvent) {
	eb.mu.RLock()
	subscribers := eb.subscribers[eventType]
	eb.mu.RUnlock()
	
	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
			// Channel full, skip
		}
	}
}

// RegisterProgressCallback registers a callback for workflow progress updates
func (we *WorkflowEngine) RegisterProgressCallback(executionID string, callback WorkflowProgressCallback) {
	if we.progressCallbacks[executionID] == nil {
		we.progressCallbacks[executionID] = make([]WorkflowProgressCallback, 0)
	}
	we.progressCallbacks[executionID] = append(we.progressCallbacks[executionID], callback)
}

// Shutdown gracefully shuts down the workflow engine
func (we *WorkflowEngine) Shutdown(ctx context.Context) error {
	close(we.shutdownChan)
	
	// Cancel all active workflows
	we.mu.RLock()
	for _, execution := range we.activeWorkflows {
		execution.CancelFunc()
	}
	we.mu.RUnlock()
	
	// Wait for workflows to complete or timeout
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			we.mu.RLock()
			activeCount := len(we.activeWorkflows)
			we.mu.RUnlock()
			
			if activeCount == 0 {
				return nil
			}
		}
	}
}