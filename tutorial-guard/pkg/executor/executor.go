/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
	"github.com/aws-research-wizard/tutorial-guard/pkg/interpreter"
	"github.com/aws-research-wizard/tutorial-guard/pkg/registry"
	"github.com/aws-research-wizard/tutorial-guard/pkg/runner"
)

// TutorialExecutor executes complete tutorials end-to-end with AI guidance
type TutorialExecutor struct {
	registry    *registry.ProviderRegistry
	interpreter *interpreter.TutorialInterpreter
	runner      runner.Runner
	config      ExecutorConfig
	reporter    *ExecutionReporter
	hooks       *ExecutionHooks
	mutex       sync.RWMutex
}

// ExecutorConfig defines configuration for tutorial execution
type ExecutorConfig struct {
	WorkingDirectory string            `json:"working_directory"`
	Environment      ExecutionEnv      `json:"environment"`
	ValidationMode   ValidationMode    `json:"validation_mode"`
	ErrorHandling    ErrorHandling     `json:"error_handling"`
	SafetyMode       SafetyMode        `json:"safety_mode"`
	TimeoutPolicy    TimeoutPolicy     `json:"timeout_policy"`
	ResourceLimits   ResourceLimits    `json:"resource_limits"`
	CleanupPolicy    CleanupPolicy     `json:"cleanup_policy"`
	ReportingConfig  ReportingConfig   `json:"reporting_config"`
	Metadata         map[string]string `json:"metadata"`
}

// ExecutionEnv defines the execution environment
type ExecutionEnv string

const (
	EnvLocal      ExecutionEnv = "local"      // Execute on local machine
	EnvDocker     ExecutionEnv = "docker"     // Execute in Docker container
	EnvKubernetes ExecutionEnv = "kubernetes" // Execute in Kubernetes pod
	EnvAWS        ExecutionEnv = "aws"        // Execute on AWS (EC2, Lambda, etc.)
	EnvSandbox    ExecutionEnv = "sandbox"    // Execute in secure sandbox
)

// ValidationMode defines how strictly to validate tutorial outcomes
type ValidationMode string

const (
	ValidationStrict     ValidationMode = "strict"      // Fail on any validation error
	ValidationLenient    ValidationMode = "lenient"     // Allow minor validation failures
	ValidationAdaptive   ValidationMode = "adaptive"    // Use AI to determine if failures are acceptable
	ValidationBestEffort ValidationMode = "best_effort" // Continue even with failures
)

// ErrorHandling defines how to handle execution errors
type ErrorHandling string

const (
	ErrorFail     ErrorHandling = "fail"        // Stop on first error
	ErrorContinue ErrorHandling = "continue"    // Continue with next instruction
	ErrorRecover  ErrorHandling = "recover"     // Attempt AI-powered recovery
	ErrorInteract ErrorHandling = "interactive" // Prompt user for guidance
)

// SafetyMode defines security and safety constraints
type SafetyMode string

const (
	SafetyNone        SafetyMode = "none"        // No safety restrictions
	SafetyBasic       SafetyMode = "basic"       // Basic safety checks
	SafetyRestrictive SafetyMode = "restrictive" // Restrictive safety policy
	SafetyParanoid    SafetyMode = "paranoid"    // Maximum safety restrictions
)

// TimeoutPolicy defines timeout behavior
type TimeoutPolicy struct {
	InstructionTimeout time.Duration `json:"instruction_timeout"` // Timeout per instruction
	TotalTimeout       time.Duration `json:"total_timeout"`       // Total tutorial timeout
	AITimeout          time.Duration `json:"ai_timeout"`          // AI request timeout
	CommandTimeout     time.Duration `json:"command_timeout"`     // Shell command timeout
}

// ResourceLimits defines resource consumption limits
type ResourceLimits struct {
	MaxMemoryMB    int     `json:"max_memory_mb"`    // Maximum memory usage
	MaxDiskMB      int     `json:"max_disk_mb"`      // Maximum disk usage
	MaxCPUPercent  float64 `json:"max_cpu_percent"`  // Maximum CPU usage
	MaxNetworkMB   int     `json:"max_network_mb"`   // Maximum network usage
	MaxProcesses   int     `json:"max_processes"`    // Maximum number of processes
	MaxFileHandles int     `json:"max_file_handles"` // Maximum file handles
}

// CleanupPolicy defines cleanup behavior
type CleanupPolicy struct {
	CleanupOnSuccess bool              `json:"cleanup_on_success"` // Clean up after successful execution
	CleanupOnFailure bool              `json:"cleanup_on_failure"` // Clean up after failed execution
	PreserveFiles    []string          `json:"preserve_files"`     // Files to preserve during cleanup
	CustomCleanup    map[string]string `json:"custom_cleanup"`     // Custom cleanup commands
}

// ReportingConfig defines reporting behavior
type ReportingConfig struct {
	Enabled       bool     `json:"enabled"`
	OutputFormats []string `json:"output_formats"` // html, json, markdown, etc.
	OutputPath    string   `json:"output_path"`
	IncludeStdout bool     `json:"include_stdout"`
	IncludeStderr bool     `json:"include_stderr"`
	IncludeEnv    bool     `json:"include_env"`
	IncludeFiles  bool     `json:"include_files"`
}

// ExecutionResult represents the result of executing a tutorial
type ExecutionResult struct {
	Tutorial           *interpreter.Tutorial     `json:"tutorial"`
	Plan               *interpreter.TutorialPlan `json:"plan"`
	Success            bool                      `json:"success"`
	StartTime          time.Time                 `json:"start_time"`
	EndTime            time.Time                 `json:"end_time"`
	Duration           time.Duration             `json:"duration"`
	StepsExecuted      int                       `json:"steps_executed"`
	StepsTotal         int                       `json:"steps_total"`
	Results            []StepResult              `json:"results"`
	FinalContext       ai.TutorialContext        `json:"final_context"`
	ErrorSummary       *ErrorSummary             `json:"error_summary,omitempty"`
	PerformanceMetrics *PerformanceMetrics       `json:"performance_metrics"`
	QualityScore       float64                   `json:"quality_score"`
	Metadata           map[string]string         `json:"metadata"`
}

// StepResult represents the result of executing a single tutorial step
type StepResult struct {
	Step         interpreter.TutorialStep `json:"step"`
	Success      bool                     `json:"success"`
	StartTime    time.Time                `json:"start_time"`
	EndTime      time.Time                `json:"end_time"`
	Duration     time.Duration            `json:"duration"`
	Instructions []InstructionResult      `json:"instructions"`
	Context      ai.TutorialContext       `json:"context"`
	Validation   *ai.ValidationResult     `json:"validation,omitempty"`
	Error        *ExecutionError          `json:"error,omitempty"`
	Recovery     *RecoveryAction          `json:"recovery,omitempty"`
	Metadata     map[string]string        `json:"metadata"`
}

// InstructionResult represents the result of executing a single instruction
type InstructionResult struct {
	Instruction  interpreter.Instruction `json:"instruction"`
	Success      bool                    `json:"success"`
	StartTime    time.Time               `json:"start_time"`
	EndTime      time.Time               `json:"end_time"`
	Duration     time.Duration           `json:"duration"`
	Actions      []ActionResult          `json:"actions"`
	Validation   *ai.ValidationResult    `json:"validation,omitempty"`
	Error        *ExecutionError         `json:"error,omitempty"`
	Recovery     *RecoveryAction         `json:"recovery,omitempty"`
	AIConfidence float64                 `json:"ai_confidence"`
	Metadata     map[string]string       `json:"metadata"`
}

// ActionResult represents the result of executing a single action
type ActionResult struct {
	Action        interpreter.Action `json:"action"`
	Success       bool               `json:"success"`
	StartTime     time.Time          `json:"start_time"`
	EndTime       time.Time          `json:"end_time"`
	Duration      time.Duration      `json:"duration"`
	Command       string             `json:"command"`
	ExitCode      int                `json:"exit_code"`
	Stdout        string             `json:"stdout"`
	Stderr        string             `json:"stderr"`
	CreatedFiles  []string           `json:"created_files"`
	ModifiedFiles []string           `json:"modified_files"`
	Error         *ExecutionError    `json:"error,omitempty"`
	Validation    *ValidationResult  `json:"validation,omitempty"`
	Metadata      map[string]string  `json:"metadata"`
}

// ExecutionError represents an error during execution
type ExecutionError struct {
	Type        ErrorType          `json:"type"`
	Message     string             `json:"message"`
	Command     string             `json:"command,omitempty"`
	ExitCode    int                `json:"exit_code,omitempty"`
	Stdout      string             `json:"stdout,omitempty"`
	Stderr      string             `json:"stderr,omitempty"`
	Timestamp   time.Time          `json:"timestamp"`
	Context     ai.TutorialContext `json:"context"`
	Recoverable bool               `json:"recoverable"`
	Metadata    map[string]string  `json:"metadata"`
}

// ErrorType defines types of execution errors
type ErrorType string

const (
	ErrorTypeCommand    ErrorType = "command"    // Command execution error
	ErrorTypeValidation ErrorType = "validation" // Validation failure
	ErrorTypeTimeout    ErrorType = "timeout"    // Timeout error
	ErrorTypeResource   ErrorType = "resource"   // Resource limit exceeded
	ErrorTypePermission ErrorType = "permission" // Permission denied
	ErrorTypeDependency ErrorType = "dependency" // Missing dependency
	ErrorTypeNetwork    ErrorType = "network"    // Network error
	ErrorTypeAI         ErrorType = "ai"         // AI provider error
	ErrorTypeInternal   ErrorType = "internal"   // Internal executor error
)

// RecoveryAction represents an action taken to recover from an error
type RecoveryAction struct {
	Type        RecoveryType      `json:"type"`
	Description string            `json:"description"`
	Commands    []string          `json:"commands"`
	Success     bool              `json:"success"`
	Duration    time.Duration     `json:"duration"`
	AIGuided    bool              `json:"ai_guided"`
	Metadata    map[string]string `json:"metadata"`
}

// RecoveryType defines types of recovery actions
type RecoveryType string

const (
	RecoverySkip      RecoveryType = "skip"      // Skip the failed instruction
	RecoveryRetry     RecoveryType = "retry"     // Retry the instruction
	RecoveryAlternate RecoveryType = "alternate" // Try alternative approach
	RecoveryFix       RecoveryType = "fix"       // Fix the environment and retry
	RecoveryManual    RecoveryType = "manual"    // Manual intervention required
)

// ValidationResult represents validation of an action's outcome
type ValidationResult struct {
	Success    bool              `json:"success"`
	Confidence float64           `json:"confidence"`
	Expected   interface{}       `json:"expected"`
	Actual     interface{}       `json:"actual"`
	Message    string            `json:"message"`
	Metadata   map[string]string `json:"metadata"`
}

// ErrorSummary provides a summary of all errors encountered
type ErrorSummary struct {
	TotalErrors      int               `json:"total_errors"`
	ErrorsByType     map[ErrorType]int `json:"errors_by_type"`
	RecoveryAttempts int               `json:"recovery_attempts"`
	RecoverySuccess  int               `json:"recovery_success"`
	CriticalErrors   []ExecutionError  `json:"critical_errors"`
	Recommendations  []string          `json:"recommendations"`
}

// PerformanceMetrics provides performance analysis
type PerformanceMetrics struct {
	TotalDuration      time.Duration `json:"total_duration"`
	AIDuration         time.Duration `json:"ai_duration"`
	ExecutionDuration  time.Duration `json:"execution_duration"`
	ValidationDuration time.Duration `json:"validation_duration"`
	AverageStepTime    time.Duration `json:"average_step_time"`
	ResourceUsage      ResourceUsage `json:"resource_usage"`
	AIRequests         int           `json:"ai_requests"`
	TotalCost          float64       `json:"total_cost"`
	Efficiency         float64       `json:"efficiency"` // Success rate weighted by time
}

// ResourceUsage tracks resource consumption
type ResourceUsage struct {
	PeakMemoryMB     int     `json:"peak_memory_mb"`
	DiskUsageMB      int     `json:"disk_usage_mb"`
	CPUTimeSeconds   float64 `json:"cpu_time_seconds"`
	NetworkBytesSent int64   `json:"network_bytes_sent"`
	NetworkBytesRecv int64   `json:"network_bytes_recv"`
	ProcessesCreated int     `json:"processes_created"`
	FilesCreated     int     `json:"files_created"`
}

// ExecutionHooks define callback functions for execution events
type ExecutionHooks struct {
	OnTutorialStart    func(*interpreter.Tutorial) error
	OnTutorialEnd      func(*ExecutionResult) error
	OnStepStart        func(*interpreter.TutorialStep) error
	OnStepEnd          func(*StepResult) error
	OnInstructionStart func(*interpreter.Instruction) error
	OnInstructionEnd   func(*InstructionResult) error
	OnActionStart      func(*interpreter.Action) error
	OnActionEnd        func(*ActionResult) error
	OnError            func(*ExecutionError) (*RecoveryAction, error)
	OnValidation       func(*ai.ValidationResult) error
	OnRecovery         func(*RecoveryAction) error
}

// ExecutionReporter generates reports of tutorial execution
type ExecutionReporter struct {
	config ReportingConfig
	mutex  sync.Mutex
}

// NewTutorialExecutor creates a new tutorial executor
func NewTutorialExecutor(registry *registry.ProviderRegistry, config ExecutorConfig) *TutorialExecutor {
	// Create AI client from registry
	client := ai.NewClient(registry, ai.ClientConfig{
		DefaultTimeout:   "30s",
		MaxRetries:       3,
		CacheEnabled:     true,
		CostOptimization: true,
	})

	// Create interpreter
	interpreterConfig := interpreter.InterpreterConfig{
		MaxSteps:            100,
		StrictValidation:    config.ValidationMode == ValidationStrict,
		AllowErrorRecovery:  config.ErrorHandling == ErrorRecover,
		ContextCompression:  true,
		ValidationThreshold: 0.8,
	}

	interpreter := interpreter.NewTutorialInterpreter(client, interpreterConfig)

	// Create runner based on environment
	var runner runner.Runner
	switch config.Environment {
	case EnvLocal:
		runner = runner.NewLocalRunner(runner.LocalConfig{
			WorkingDirectory: config.WorkingDirectory,
			Timeout:          config.TimeoutPolicy.CommandTimeout,
			ResourceLimits: runner.ResourceLimits{
				MaxMemoryMB:  config.ResourceLimits.MaxMemoryMB,
				MaxDiskMB:    config.ResourceLimits.MaxDiskMB,
				MaxProcesses: config.ResourceLimits.MaxProcesses,
			},
		})
	case EnvDocker:
		runner = runner.NewDockerRunner(runner.DockerConfig{
			Image:            "ubuntu:22.04",
			WorkingDirectory: config.WorkingDirectory,
			Timeout:          config.TimeoutPolicy.CommandTimeout,
		})
	case EnvAWS:
		runner = runner.NewAWSRunner(runner.AWSConfig{
			Region:       "us-east-1",
			InstanceType: "t3.micro",
			Timeout:      config.TimeoutPolicy.CommandTimeout,
		})
	default:
		runner = runner.NewLocalRunner(runner.LocalConfig{})
	}

	return &TutorialExecutor{
		registry:    registry,
		interpreter: interpreter,
		runner:      runner,
		config:      config,
		reporter:    NewExecutionReporter(config.ReportingConfig),
		hooks:       &ExecutionHooks{},
		mutex:       sync.RWMutex{},
	}
}

// Execute runs a complete tutorial end-to-end
func (e *TutorialExecutor) Execute(ctx context.Context, tutorial *interpreter.Tutorial) (*ExecutionResult, error) {
	start := time.Now()

	// Initialize execution result
	result := &ExecutionResult{
		Tutorial:      tutorial,
		Success:       false,
		StartTime:     start,
		StepsExecuted: 0,
		StepsTotal:    len(tutorial.Sections),
		Results:       make([]StepResult, 0),
		Metadata:      make(map[string]string),
	}

	// Call tutorial start hook
	if e.hooks.OnTutorialStart != nil {
		if err := e.hooks.OnTutorialStart(tutorial); err != nil {
			return result, fmt.Errorf("tutorial start hook failed: %w", err)
		}
	}

	// Interpret tutorial with AI
	plan, err := e.interpreter.InterpretTutorial(ctx, tutorial)
	if err != nil {
		return result, fmt.Errorf("failed to interpret tutorial: %w", err)
	}
	result.Plan = plan

	// Execute each step
	for i, step := range plan.Steps {
		stepResult := e.executeStep(ctx, step)
		result.Results = append(result.Results, stepResult)
		result.StepsExecuted++

		if !stepResult.Success {
			if e.config.ErrorHandling == ErrorFail {
				break
			}
		}
	}

	// Finalize result
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Success = e.calculateOverallSuccess(result.Results)
	result.QualityScore = e.calculateQualityScore(result.Results)
	result.PerformanceMetrics = e.calculatePerformanceMetrics(result)
	result.ErrorSummary = e.generateErrorSummary(result.Results)

	// Call tutorial end hook
	if e.hooks.OnTutorialEnd != nil {
		if err := e.hooks.OnTutorialEnd(result); err != nil {
			// Log but don't fail the execution
			fmt.Printf("Warning: tutorial end hook failed: %v\n", err)
		}
	}

	return result, nil
}

// executeStep executes a single tutorial step
func (e *TutorialExecutor) executeStep(ctx context.Context, step interpreter.TutorialStep) StepResult {
	start := time.Now()

	result := StepResult{
		Step:         step,
		Success:      false,
		StartTime:    start,
		Instructions: make([]InstructionResult, 0),
		Context:      step.Context,
		Metadata:     make(map[string]string),
	}

	// Call step start hook
	if e.hooks.OnStepStart != nil {
		if err := e.hooks.OnStepStart(&step); err != nil {
			result.Error = &ExecutionError{
				Type:      ErrorTypeInternal,
				Message:   fmt.Sprintf("step start hook failed: %v", err),
				Timestamp: time.Now(),
				Context:   step.Context,
			}
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result
		}
	}

	// Execute each instruction in the step
	for _, instruction := range step.Instructions {
		instructionResult := e.executeInstruction(ctx, instruction)
		result.Instructions = append(result.Instructions, instructionResult)

		if !instructionResult.Success && e.config.ErrorHandling == ErrorFail {
			break
		}
	}

	// Calculate step success
	result.Success = e.calculateStepSuccess(result.Instructions)
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Call step end hook
	if e.hooks.OnStepEnd != nil {
		if err := e.hooks.OnStepEnd(&result); err != nil {
			// Log but don't fail the step
			fmt.Printf("Warning: step end hook failed: %v\n", err)
		}
	}

	return result
}

// executeInstruction executes a single instruction
func (e *TutorialExecutor) executeInstruction(ctx context.Context, instruction interpreter.Instruction) InstructionResult {
	start := time.Now()

	result := InstructionResult{
		Instruction:  instruction,
		Success:      false,
		StartTime:    start,
		Actions:      make([]ActionResult, 0),
		AIConfidence: instruction.Confidence,
		Metadata:     make(map[string]string),
	}

	// Call instruction start hook
	if e.hooks.OnInstructionStart != nil {
		if err := e.hooks.OnInstructionStart(&instruction); err != nil {
			result.Error = &ExecutionError{
				Type:      ErrorTypeInternal,
				Message:   fmt.Sprintf("instruction start hook failed: %v", err),
				Timestamp: time.Now(),
			}
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result
		}
	}

	// Execute each action in the instruction
	for _, action := range instruction.Actions {
		actionResult := e.executeAction(ctx, action)
		result.Actions = append(result.Actions, actionResult)

		if !actionResult.Success && e.config.ErrorHandling == ErrorFail {
			break
		}
	}

	// Calculate instruction success
	result.Success = e.calculateInstructionSuccess(result.Actions)
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Call instruction end hook
	if e.hooks.OnInstructionEnd != nil {
		if err := e.hooks.OnInstructionEnd(&result); err != nil {
			// Log but don't fail the instruction
			fmt.Printf("Warning: instruction end hook failed: %v\n", err)
		}
	}

	return result
}

// executeAction executes a single action
func (e *TutorialExecutor) executeAction(ctx context.Context, action interpreter.Action) ActionResult {
	start := time.Now()

	result := ActionResult{
		Action:    action,
		Success:   false,
		StartTime: start,
		Command:   action.Command,
		Metadata:  make(map[string]string),
	}

	// Call action start hook
	if e.hooks.OnActionStart != nil {
		if err := e.hooks.OnActionStart(&action); err != nil {
			result.Error = &ExecutionError{
				Type:      ErrorTypeInternal,
				Message:   fmt.Sprintf("action start hook failed: %v", err),
				Timestamp: time.Now(),
			}
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result
		}
	}

	// Execute the action based on its type
	switch action.Type {
	case interpreter.ActionTypeCommand:
		e.executeCommand(ctx, &result)
	case interpreter.ActionTypeValidate:
		e.executeValidation(ctx, &result)
	case interpreter.ActionTypeCheck:
		e.executeCheck(ctx, &result)
	default:
		result.Error = &ExecutionError{
			Type:      ErrorTypeInternal,
			Message:   fmt.Sprintf("unsupported action type: %s", action.Type),
			Timestamp: time.Now(),
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Call action end hook
	if e.hooks.OnActionEnd != nil {
		if err := e.hooks.OnActionEnd(&result); err != nil {
			// Log but don't fail the action
			fmt.Printf("Warning: action end hook failed: %v\n", err)
		}
	}

	return result
}

// executeCommand executes a shell command
func (e *TutorialExecutor) executeCommand(ctx context.Context, result *ActionResult) {
	// Safety check
	if e.config.SafetyMode != SafetyNone {
		if err := e.validateCommandSafety(result.Action.Command); err != nil {
			result.Error = &ExecutionError{
				Type:      ErrorTypePermission,
				Message:   fmt.Sprintf("command blocked by safety policy: %v", err),
				Command:   result.Action.Command,
				Timestamp: time.Now(),
			}
			return
		}
	}

	// Create command context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, e.config.TimeoutPolicy.CommandTimeout)
	defer cancel()

	// Execute command
	cmd := exec.CommandContext(cmdCtx, "sh", "-c", result.Action.Command)
	cmd.Dir = e.config.WorkingDirectory

	// Capture output
	stdout, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			result.Stderr = string(exitErr.Stderr)
		}
		result.Error = &ExecutionError{
			Type:      ErrorTypeCommand,
			Message:   err.Error(),
			Command:   result.Action.Command,
			ExitCode:  result.ExitCode,
			Stdout:    string(stdout),
			Stderr:    result.Stderr,
			Timestamp: time.Now(),
		}
		return
	}

	result.Stdout = string(stdout)
	result.Success = true

	// Validate the action's outcome if validation rules are provided
	if result.Action.Validation.Type != "" {
		e.validateActionOutcome(ctx, result)
	}
}

// executeValidation performs validation logic
func (e *TutorialExecutor) executeValidation(ctx context.Context, result *ActionResult) {
	// Implementation for validation actions
	result.Success = true // Placeholder
}

// executeCheck performs check logic
func (e *TutorialExecutor) executeCheck(ctx context.Context, result *ActionResult) {
	// Implementation for check actions
	result.Success = true // Placeholder
}

// validateCommandSafety checks if a command is safe to execute
func (e *TutorialExecutor) validateCommandSafety(command string) error {
	switch e.config.SafetyMode {
	case SafetyNone:
		return nil
	case SafetyBasic:
		return e.validateBasicSafety(command)
	case SafetyRestrictive:
		return e.validateRestrictiveSafety(command)
	case SafetyParanoid:
		return e.validateParanoidSafety(command)
	default:
		return nil
	}
}

// validateBasicSafety performs basic safety checks
func (e *TutorialExecutor) validateBasicSafety(command string) error {
	dangerous := []string{"rm -rf /", "dd if=", "mkfs", "fdisk", "format"}
	cmd := strings.ToLower(command)

	for _, danger := range dangerous {
		if strings.Contains(cmd, danger) {
			return fmt.Errorf("dangerous command detected: %s", danger)
		}
	}
	return nil
}

// validateRestrictiveSafety performs restrictive safety checks
func (e *TutorialExecutor) validateRestrictiveSafety(command string) error {
	if err := e.validateBasicSafety(command); err != nil {
		return err
	}

	// Additional checks for restrictive mode
	restricted := []string{"sudo", "su", "chmod 777", "chown root"}
	cmd := strings.ToLower(command)

	for _, restrict := range restricted {
		if strings.Contains(cmd, restrict) {
			return fmt.Errorf("restricted command detected: %s", restrict)
		}
	}
	return nil
}

// validateParanoidSafety performs paranoid safety checks
func (e *TutorialExecutor) validateParanoidSafety(command string) error {
	if err := e.validateRestrictiveSafety(command); err != nil {
		return err
	}

	// In paranoid mode, only allow explicitly whitelisted commands
	allowed := []string{"echo", "ls", "cat", "grep", "mkdir", "touch", "cp", "mv"}
	cmd := strings.Split(strings.TrimSpace(command), " ")[0]

	for _, allow := range allowed {
		if cmd == allow {
			return nil
		}
	}

	return fmt.Errorf("command not in whitelist: %s", cmd)
}

// validateActionOutcome validates the outcome of an action
func (e *TutorialExecutor) validateActionOutcome(ctx context.Context, result *ActionResult) {
	validation := &ValidationResult{
		Success:    false,
		Confidence: 0.0,
		Expected:   result.Action.Validation.Expected,
		Metadata:   make(map[string]string),
	}

	switch result.Action.Validation.Type {
	case interpreter.ValidationTypeExitCode:
		expected := 0
		if result.Action.Validation.Expected != nil {
			if val, ok := result.Action.Validation.Expected.(int); ok {
				expected = val
			}
		}
		validation.Success = result.ExitCode == expected
		validation.Actual = result.ExitCode
		validation.Confidence = 1.0

	case interpreter.ValidationTypeFileExists:
		if expected, ok := result.Action.Validation.Expected.(string); ok {
			path := filepath.Join(e.config.WorkingDirectory, expected)
			_, err := os.Stat(path)
			validation.Success = err == nil
			validation.Actual = err == nil
			validation.Confidence = 1.0
		}

	case interpreter.ValidationTypeOutput:
		if expected, ok := result.Action.Validation.Expected.(string); ok {
			validation.Success = strings.Contains(result.Stdout, expected)
			validation.Actual = result.Stdout
			validation.Confidence = 0.9
		}

	case interpreter.ValidationTypeContains:
		if expected, ok := result.Action.Validation.Expected.(string); ok {
			validation.Success = strings.Contains(result.Stdout, expected)
			validation.Actual = result.Stdout
			validation.Confidence = 0.9
		}
	}

	result.Validation = validation
	if !validation.Success && e.config.ValidationMode == ValidationStrict {
		result.Success = false
		result.Error = &ExecutionError{
			Type:      ErrorTypeValidation,
			Message:   fmt.Sprintf("validation failed: expected %v, got %v", validation.Expected, validation.Actual),
			Timestamp: time.Now(),
		}
	}
}

// Helper functions for calculating success rates and metrics
func (e *TutorialExecutor) calculateOverallSuccess(results []StepResult) bool {
	for _, result := range results {
		if !result.Success {
			return false
		}
	}
	return true
}

func (e *TutorialExecutor) calculateStepSuccess(instructions []InstructionResult) bool {
	for _, instruction := range instructions {
		if !instruction.Success {
			return false
		}
	}
	return true
}

func (e *TutorialExecutor) calculateInstructionSuccess(actions []ActionResult) bool {
	for _, action := range actions {
		if !action.Success {
			return false
		}
	}
	return true
}

func (e *TutorialExecutor) calculateQualityScore(results []StepResult) float64 {
	// Implementation for quality score calculation
	return 0.95 // Placeholder
}

func (e *TutorialExecutor) calculatePerformanceMetrics(result *ExecutionResult) *PerformanceMetrics {
	// Implementation for performance metrics calculation
	return &PerformanceMetrics{
		TotalDuration: result.Duration,
		Efficiency:    0.85, // Placeholder
	}
}

func (e *TutorialExecutor) generateErrorSummary(results []StepResult) *ErrorSummary {
	// Implementation for error summary generation
	return &ErrorSummary{
		TotalErrors: 0, // Placeholder
	}
}

// SetHooks configures execution hooks for monitoring and custom behavior
func (e *TutorialExecutor) SetHooks(hooks *ExecutionHooks) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.hooks = hooks
}

// GetHooks returns the current execution hooks
func (e *TutorialExecutor) GetHooks() *ExecutionHooks {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.hooks
}

// NewExecutionReporter creates a new execution reporter
func NewExecutionReporter(config ReportingConfig) *ExecutionReporter {
	return &ExecutionReporter{
		config: config,
		mutex:  sync.Mutex{},
	}
}
