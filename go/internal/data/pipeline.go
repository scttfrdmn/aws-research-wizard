package data

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Pipeline represents a data processing pipeline
type Pipeline struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Jobs        []Job          `json:"jobs"`
	Status      PipelineStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Config      PipelineConfig `json:"config"`
}

// Job represents a single job within a pipeline
type Job struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         JobType           `json:"type"`
	Status       JobStatus         `json:"status"`
	Dependencies []string          `json:"dependencies"` // Job IDs this job depends on
	Config       map[string]string `json:"config"`
	StartTime    *time.Time        `json:"start_time,omitempty"`
	EndTime      *time.Time        `json:"end_time,omitempty"`
	Error        string            `json:"error,omitempty"`
	Progress     float64           `json:"progress"` // 0-100
}

// PipelineStatus represents the overall status of a pipeline
type PipelineStatus string

const (
	PipelineStatusPending   PipelineStatus = "pending"
	PipelineStatusRunning   PipelineStatus = "running"
	PipelineStatusCompleted PipelineStatus = "completed"
	PipelineStatusFailed    PipelineStatus = "failed"
	PipelineStatusCancelled PipelineStatus = "cancelled"
)

// JobStatus represents the status of a job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusSkipped   JobStatus = "skipped"
)

// JobType represents the type of job
type JobType string

const (
	JobTypeDownload  JobType = "download"
	JobTypeUpload    JobType = "upload"
	JobTypeTransform JobType = "transform"
	JobTypeValidate  JobType = "validate"
	JobTypeCleanup   JobType = "cleanup"
)

// PipelineConfig holds configuration for pipeline execution
type PipelineConfig struct {
	MaxConcurrentJobs int           `json:"max_concurrent_jobs"`
	RetryAttempts     int           `json:"retry_attempts"`
	RetryDelay        time.Duration `json:"retry_delay"`
	Timeout           time.Duration `json:"timeout"`
	WorkingDirectory  string        `json:"working_directory"`
}

// PipelineManager manages data processing pipelines
type PipelineManager struct {
	s3Manager       *S3Manager
	openData        *OpenDataRegistry
	pipelines       map[string]*Pipeline
	activePipelines map[string]*pipelineExecution
	mu              sync.RWMutex
	configPath      string
}

// pipelineExecution tracks the execution state of a pipeline
type pipelineExecution struct {
	pipeline    *Pipeline
	ctx         context.Context
	cancel      context.CancelFunc
	jobStatuses map[string]JobStatus
	errors      []error
	startTime   time.Time
}

// NewPipelineManager creates a new pipeline manager
func NewPipelineManager(s3Manager *S3Manager, openData *OpenDataRegistry, configPath string) *PipelineManager {
	return &PipelineManager{
		s3Manager:       s3Manager,
		openData:        openData,
		pipelines:       make(map[string]*Pipeline),
		activePipelines: make(map[string]*pipelineExecution),
		configPath:      configPath,
	}
}

// CreatePipeline creates a new data processing pipeline
func (pm *PipelineManager) CreatePipeline(name, description string, config PipelineConfig) *Pipeline {
	pipeline := &Pipeline{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Jobs:        make([]Job, 0),
		Status:      PipelineStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Config:      config,
	}

	pm.mu.Lock()
	pm.pipelines[pipeline.ID] = pipeline
	pm.mu.Unlock()

	return pipeline
}

// AddJob adds a job to a pipeline
func (pm *PipelineManager) AddJob(pipelineID string, job Job) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pipeline, exists := pm.pipelines[pipelineID]
	if !exists {
		return fmt.Errorf("pipeline not found: %s", pipelineID)
	}

	job.ID = generateID()
	job.Status = JobStatusPending
	pipeline.Jobs = append(pipeline.Jobs, job)
	pipeline.UpdatedAt = time.Now()

	return nil
}

// AddDownloadJob adds a download job to a pipeline
func (pm *PipelineManager) AddDownloadJob(pipelineID, name, bucket, key, localPath string, dependencies []string) error {
	job := Job{
		Name:         name,
		Type:         JobTypeDownload,
		Dependencies: dependencies,
		Config: map[string]string{
			"bucket":     bucket,
			"key":        key,
			"local_path": localPath,
		},
	}

	return pm.AddJob(pipelineID, job)
}

// AddUploadJob adds an upload job to a pipeline
func (pm *PipelineManager) AddUploadJob(pipelineID, name, localPath, bucket, key string, dependencies []string) error {
	job := Job{
		Name:         name,
		Type:         JobTypeUpload,
		Dependencies: dependencies,
		Config: map[string]string{
			"local_path": localPath,
			"bucket":     bucket,
			"key":        key,
		},
	}

	return pm.AddJob(pipelineID, job)
}

// ExecutePipeline executes a pipeline asynchronously
func (pm *PipelineManager) ExecutePipeline(ctx context.Context, pipelineID string) error {
	pm.mu.RLock()
	pipeline, exists := pm.pipelines[pipelineID]
	pm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("pipeline not found: %s", pipelineID)
	}

	// Check if pipeline is already running
	pm.mu.Lock()
	if _, running := pm.activePipelines[pipelineID]; running {
		pm.mu.Unlock()
		return fmt.Errorf("pipeline is already running: %s", pipelineID)
	}

	// Create execution context
	execCtx, cancel := context.WithCancel(ctx)
	execution := &pipelineExecution{
		pipeline:    pipeline,
		ctx:         execCtx,
		cancel:      cancel,
		jobStatuses: make(map[string]JobStatus),
		startTime:   time.Now(),
	}

	pm.activePipelines[pipelineID] = execution
	pm.mu.Unlock()

	// Update pipeline status
	pipeline.Status = PipelineStatusRunning
	pipeline.UpdatedAt = time.Now()

	// Execute pipeline in goroutine
	go func() {
		defer func() {
			pm.mu.Lock()
			delete(pm.activePipelines, pipelineID)
			pm.mu.Unlock()
		}()

		err := pm.executePipelineJobs(execution)
		if err != nil {
			pipeline.Status = PipelineStatusFailed
		} else {
			pipeline.Status = PipelineStatusCompleted
		}
		pipeline.UpdatedAt = time.Now()
	}()

	return nil
}

// executePipelineJobs executes all jobs in a pipeline respecting dependencies
func (pm *PipelineManager) executePipelineJobs(execution *pipelineExecution) error {
	pipeline := execution.pipeline
	jobMap := make(map[string]*Job)
	for i := range pipeline.Jobs {
		jobMap[pipeline.Jobs[i].ID] = &pipeline.Jobs[i]
	}

	// Create dependency graph
	dependsOn := make(map[string][]string)
	dependents := make(map[string][]string)

	for _, job := range pipeline.Jobs {
		dependsOn[job.ID] = job.Dependencies
		for _, dep := range job.Dependencies {
			dependents[dep] = append(dependents[dep], job.ID)
		}
	}

	// Find jobs with no dependencies to start with
	ready := make([]string, 0)
	for _, job := range pipeline.Jobs {
		if len(job.Dependencies) == 0 {
			ready = append(ready, job.ID)
		}
	}

	// Track completed jobs
	completed := make(map[string]bool)
	semaphore := make(chan struct{}, pipeline.Config.MaxConcurrentJobs)

	var wg sync.WaitGroup

	for len(completed) < len(pipeline.Jobs) {
		if len(ready) == 0 {
			// Wait for some job to complete
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Execute ready jobs
		for len(ready) > 0 && len(completed) < len(pipeline.Jobs) {
			jobID := ready[0]
			ready = ready[1:]

			if completed[jobID] {
				continue
			}

			wg.Add(1)
			go func(jID string) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				job := jobMap[jID]
				err := pm.executeJob(execution.ctx, job)

				pm.mu.Lock()
				if err != nil {
					job.Status = JobStatusFailed
					job.Error = err.Error()
					execution.errors = append(execution.errors, err)
				} else {
					job.Status = JobStatusCompleted
					completed[jID] = true

					// Add dependent jobs to ready queue
					for _, dependent := range dependents[jID] {
						if pm.areJobDependenciesCompleted(dependsOn[dependent], completed) {
							ready = append(ready, dependent)
						}
					}
				}
				pm.mu.Unlock()
			}(jobID)
		}

		// Brief pause to avoid busy waiting
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()

	if len(execution.errors) > 0 {
		return fmt.Errorf("pipeline execution failed with %d errors", len(execution.errors))
	}

	return nil
}

// executeJob executes a single job
func (pm *PipelineManager) executeJob(ctx context.Context, job *Job) error {
	now := time.Now()
	job.StartTime = &now
	job.Status = JobStatusRunning

	defer func() {
		endTime := time.Now()
		job.EndTime = &endTime
	}()

	switch job.Type {
	case JobTypeDownload:
		return pm.executeDownloadJob(ctx, job)
	case JobTypeUpload:
		return pm.executeUploadJob(ctx, job)
	case JobTypeValidate:
		return pm.executeValidateJob(ctx, job)
	case JobTypeCleanup:
		return pm.executeCleanupJob(ctx, job)
	default:
		return fmt.Errorf("unsupported job type: %s", job.Type)
	}
}

// executeDownloadJob executes a download job
func (pm *PipelineManager) executeDownloadJob(ctx context.Context, job *Job) error {
	bucket := job.Config["bucket"]
	key := job.Config["key"]
	localPath := job.Config["local_path"]

	if bucket == "" || key == "" || localPath == "" {
		return fmt.Errorf("missing required config for download job")
	}

	// Create progress callback
	progressCallback := func(progress TransferProgress) {
		job.Progress = progress.Percentage
	}

	return pm.s3Manager.DownloadFile(ctx, bucket, key, localPath, progressCallback)
}

// executeUploadJob executes an upload job
func (pm *PipelineManager) executeUploadJob(ctx context.Context, job *Job) error {
	localPath := job.Config["local_path"]
	bucket := job.Config["bucket"]
	key := job.Config["key"]

	if localPath == "" || bucket == "" || key == "" {
		return fmt.Errorf("missing required config for upload job")
	}

	// Create progress callback
	progressCallback := func(progress TransferProgress) {
		job.Progress = progress.Percentage
	}

	return pm.s3Manager.UploadFile(ctx, bucket, key, localPath, progressCallback)
}

// executeValidateJob executes a validation job
func (pm *PipelineManager) executeValidateJob(ctx context.Context, job *Job) error {
	localPath := job.Config["local_path"]
	if localPath == "" {
		return fmt.Errorf("missing local_path for validation job")
	}

	// Check if file exists and is readable
	if _, err := os.Stat(localPath); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	job.Progress = 100
	return nil
}

// executeCleanupJob executes a cleanup job
func (pm *PipelineManager) executeCleanupJob(ctx context.Context, job *Job) error {
	localPath := job.Config["local_path"]
	if localPath == "" {
		return fmt.Errorf("missing local_path for cleanup job")
	}

	// Remove file or directory
	err := os.RemoveAll(localPath)
	if err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	job.Progress = 100
	return nil
}

// areJobDependenciesCompleted checks if all dependencies of a job are completed
func (pm *PipelineManager) areJobDependenciesCompleted(dependencies []string, completed map[string]bool) bool {
	for _, dep := range dependencies {
		if !completed[dep] {
			return false
		}
	}
	return true
}

// GetPipeline retrieves a pipeline by ID
func (pm *PipelineManager) GetPipeline(pipelineID string) (*Pipeline, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	pipeline, exists := pm.pipelines[pipelineID]
	if !exists {
		return nil, fmt.Errorf("pipeline not found: %s", pipelineID)
	}

	return pipeline, nil
}

// ListPipelines returns all pipelines
func (pm *PipelineManager) ListPipelines() []*Pipeline {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	pipelines := make([]*Pipeline, 0, len(pm.pipelines))
	for _, pipeline := range pm.pipelines {
		pipelines = append(pipelines, pipeline)
	}

	return pipelines
}

// CancelPipeline cancels a running pipeline
func (pm *PipelineManager) CancelPipeline(pipelineID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	execution, exists := pm.activePipelines[pipelineID]
	if !exists {
		return fmt.Errorf("pipeline is not running: %s", pipelineID)
	}

	execution.cancel()
	execution.pipeline.Status = PipelineStatusCancelled
	execution.pipeline.UpdatedAt = time.Now()

	return nil
}

// SavePipeline saves a pipeline to disk
func (pm *PipelineManager) SavePipeline(pipelineID string) error {
	pm.mu.RLock()
	pipeline, exists := pm.pipelines[pipelineID]
	pm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("pipeline not found: %s", pipelineID)
	}

	if pm.configPath == "" {
		return fmt.Errorf("config path not set")
	}

	// Ensure directory exists
	pipelineDir := filepath.Join(pm.configPath, "pipelines")
	if err := os.MkdirAll(pipelineDir, 0755); err != nil {
		return fmt.Errorf("failed to create pipeline directory: %w", err)
	}

	// Save pipeline to file
	filePath := filepath.Join(pipelineDir, fmt.Sprintf("%s.json", pipelineID))
	data, err := json.MarshalIndent(pipeline, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pipeline: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write pipeline file: %w", err)
	}

	return nil
}

// LoadPipeline loads a pipeline from disk
func (pm *PipelineManager) LoadPipeline(pipelineID string) error {
	if pm.configPath == "" {
		return fmt.Errorf("config path not set")
	}

	filePath := filepath.Join(pm.configPath, "pipelines", fmt.Sprintf("%s.json", pipelineID))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read pipeline file: %w", err)
	}

	var pipeline Pipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return fmt.Errorf("failed to unmarshal pipeline: %w", err)
	}

	pm.mu.Lock()
	pm.pipelines[pipelineID] = &pipeline
	pm.mu.Unlock()

	return nil
}

// generateID generates a unique ID for pipelines and jobs
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
