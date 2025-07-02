/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package ai

import (
	"context"
	"time"
)

// Provider defines the interface that all AI providers must implement
type Provider interface {
	// Core AI capabilities
	ParseInstruction(ctx context.Context, instruction string, context TutorialContext) (*ParsedInstruction, error)
	ValidateExpectation(ctx context.Context, expected, actual string, context TutorialContext) (*ValidationResult, error)
	CompressContext(ctx context.Context, fullContext TutorialContext) (*CompressedContext, error)
	InterpretError(ctx context.Context, errorMsg string, context TutorialContext) (*ErrorInterpretation, error)

	// Provider metadata
	GetCapabilities() ProviderCapabilities
	GetCostEstimate(request AIRequest) CostEstimate
	GetPerformanceMetrics() PerformanceMetrics

	// Health and monitoring
	HealthCheck(ctx context.Context) error
	GetUsageStats() UsageStats
}

// TutorialContext represents the current state of tutorial execution
type TutorialContext struct {
	WorkingDirectory string            `json:"working_directory"`
	EnvironmentVars  map[string]string `json:"environment_vars"`
	CreatedFiles     []string          `json:"created_files"`
	ExecutedCommands []string          `json:"executed_commands"`
	PreviousOutputs  []string          `json:"previous_outputs"`
	CurrentStep      int               `json:"current_step"`
	TotalSteps       int               `json:"total_steps"`
	Metadata         map[string]string `json:"metadata"`
}

// ParsedInstruction represents an AI's understanding of a tutorial instruction
type ParsedInstruction struct {
	OriginalText     string            `json:"original_text"`
	Intent           string            `json:"intent"`            // What the instruction wants to accomplish
	Actions          []Action          `json:"actions"`           // Specific actions to take
	Prerequisites    []string          `json:"prerequisites"`     // What must be done first
	ExpectedOutcomes []string          `json:"expected_outcomes"` // What should happen
	Confidence       float64           `json:"confidence"`        // AI confidence (0-1)
	Reasoning        string            `json:"reasoning"`         // Why the AI interpreted it this way
	Metadata         map[string]string `json:"metadata"`
}

// Action represents a specific action to be executed
type Action struct {
	Type        ActionType        `json:"type"`
	Command     string            `json:"command"`     // Shell command to execute
	Description string            `json:"description"` // Human-readable description
	Validation  ValidationRule    `json:"validation"`  // How to verify success
	Timeout     time.Duration     `json:"timeout"`     // Maximum execution time
	Metadata    map[string]string `json:"metadata"`
}

// ActionType defines the types of actions Tutorial Guard can perform
type ActionType string

const (
	ActionCommand     ActionType = "command"     // Execute shell command
	ActionValidate    ActionType = "validate"    // Validate a condition
	ActionWait        ActionType = "wait"        // Wait for a condition
	ActionDownload    ActionType = "download"    // Download a file
	ActionExtract     ActionType = "extract"     // Extract an archive
	ActionNavigate    ActionType = "navigate"    // Change directory
	ActionCheck       ActionType = "check"       // Check if something exists
	ActionConditional ActionType = "conditional" // Conditional logic
)

// ValidationRule defines how to validate an action's success
type ValidationRule struct {
	Type      ValidationType    `json:"type"`
	Condition string            `json:"condition"` // What to check
	Expected  interface{}       `json:"expected"`  // Expected result (can be string, bool, number)
	Operator  string            `json:"operator"`  // Comparison operator
	Metadata  map[string]string `json:"metadata"`
}

// ValidationType defines types of validation
type ValidationType string

const (
	ValidationExitCode   ValidationType = "exit_code"   // Check command exit code
	ValidationFileExists ValidationType = "file_exists" // Check if file exists
	ValidationOutput     ValidationType = "output"      // Check command output
	ValidationContains   ValidationType = "contains"    // Check if output contains text
	ValidationRegex      ValidationType = "regex"       // Check against regex pattern
	ValidationCustom     ValidationType = "custom"      // Custom validation logic
)

// ValidationResult represents the result of expectation validation
type ValidationResult struct {
	Success     bool              `json:"success"`
	Confidence  float64           `json:"confidence"`  // How confident the AI is (0-1)
	Reasoning   string            `json:"reasoning"`   // Why it succeeded/failed
	Differences []string          `json:"differences"` // What was different
	Suggestions []string          `json:"suggestions"` // Suggestions for improvement
	Metadata    map[string]string `json:"metadata"`
}

// CompressedContext represents a compressed version of tutorial context
type CompressedContext struct {
	Summary      string            `json:"summary"`       // Brief summary of what's happened
	KeyFiles     []string          `json:"key_files"`     // Important files created
	CurrentState string            `json:"current_state"` // Current working state
	Metadata     map[string]string `json:"metadata"`
}

// ErrorInterpretation represents AI's understanding of an error
type ErrorInterpretation struct {
	ErrorType   string            `json:"error_type"`  // Category of error
	Explanation string            `json:"explanation"` // What went wrong
	Solutions   []Solution        `json:"solutions"`   // Possible solutions
	Confidence  float64           `json:"confidence"`  // AI confidence
	Metadata    map[string]string `json:"metadata"`
}

// Solution represents a potential solution to an error
type Solution struct {
	Description string            `json:"description"` // What to do
	Commands    []string          `json:"commands"`    // Commands to run
	Probability float64           `json:"probability"` // Likelihood this will work
	Metadata    map[string]string `json:"metadata"`
}

// ProviderCapabilities describes what an AI provider can do
type ProviderCapabilities struct {
	Name               string             `json:"name"`
	Version            string             `json:"version"`
	MaxContextLength   int                `json:"max_context_length"`
	SupportedLanguages []string           `json:"supported_languages"`
	Features           []string           `json:"features"`
	CostModel          CostModel          `json:"cost_model"`
	LatencyProfile     LatencyProfile     `json:"latency_profile"`
	QualityScore       float64            `json:"quality_score"`
	CertificationLevel CertificationLevel `json:"certification_level"`
}

// CostModel represents pricing information for a provider
type CostModel struct {
	InputTokenPrice  float64 `json:"input_token_price"`  // Price per input token
	OutputTokenPrice float64 `json:"output_token_price"` // Price per output token
	Currency         string  `json:"currency"`           // Currency (e.g., "USD")
}

// LatencyProfile represents typical response times
type LatencyProfile struct {
	P50 string `json:"p50"` // 50th percentile latency
	P95 string `json:"p95"` // 95th percentile latency
	P99 string `json:"p99"` // 99th percentile latency
}

// CertificationLevel represents the quality certification level
type CertificationLevel string

const (
	CertifiedGold   CertificationLevel = "gold"
	CertifiedSilver CertificationLevel = "silver"
	CertifiedBronze CertificationLevel = "bronze"
	Community       CertificationLevel = "community"
	Enterprise      CertificationLevel = "enterprise"
)

// AIRequest represents a request to an AI provider
type AIRequest struct {
	Type            RequestType       `json:"type"`
	Content         string            `json:"content"`
	Context         TutorialContext   `json:"context"`
	Priority        Priority          `json:"priority"`
	MaxTokens       int               `json:"max_tokens"`
	Temperature     float64           `json:"temperature"`
	EstimatedTokens int               `json:"estimated_tokens"`
	Metadata        map[string]string `json:"metadata"`
}

// RequestType defines types of AI requests
type RequestType string

const (
	RequestParseInstruction    RequestType = "parse_instruction"
	RequestValidateExpectation RequestType = "validate_expectation"
	RequestCompressContext     RequestType = "compress_context"
	RequestInterpretError      RequestType = "interpret_error"
)

// Priority defines request priority levels
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// CostEstimate represents estimated cost for a request
type CostEstimate struct {
	EstimatedTokens int     `json:"estimated_tokens"`
	EstimatedCost   float64 `json:"estimated_cost"`
	Currency        string  `json:"currency"`
}

// PerformanceMetrics represents provider performance statistics
type PerformanceMetrics struct {
	AverageLatency      time.Duration `json:"average_latency"`
	AverageCost         float64       `json:"average_cost"`
	AccuracyRate        float64       `json:"accuracy_rate"`
	ErrorRate           float64       `json:"error_rate"`
	Availability        float64       `json:"availability"`
	CostEfficiencyScore float64       `json:"cost_efficiency_score"`
}

// UsageStats represents current usage statistics
type UsageStats struct {
	RequestsToday   int     `json:"requests_today"`
	TokensUsedToday int     `json:"tokens_used_today"`
	CostToday       float64 `json:"cost_today"`
	RequestsTotal   int     `json:"requests_total"`
	TokensUsedTotal int     `json:"tokens_used_total"`
	CostTotal       float64 `json:"cost_total"`
}

// CacheEntry represents a cached AI response
type CacheEntry struct {
	Key       string            `json:"key"`
	Request   AIRequest         `json:"request"`
	Response  interface{}       `json:"response"`
	CreatedAt time.Time         `json:"created_at"`
	ExpiresAt time.Time         `json:"expires_at"`
	HitCount  int               `json:"hit_count"`
	Metadata  map[string]string `json:"metadata"`
}
