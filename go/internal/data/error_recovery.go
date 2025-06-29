package data

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ErrorRecoveryManager handles comprehensive error recovery strategies
type ErrorRecoveryManager struct {
	strategies       map[string]RecoveryStrategy
	circuitBreakers  map[string]*CircuitBreaker
	errorClassifier  *ErrorClassifier
	mu               sync.RWMutex
}

// RecoveryStrategy defines how to handle specific types of errors
type RecoveryStrategy struct {
	Name                 string        `json:"name"`
	MaxRetryAttempts     int           `json:"max_retry_attempts"`
	BaseRetryDelay       time.Duration `json:"base_retry_delay"`
	MaxRetryDelay        time.Duration `json:"max_retry_delay"`
	ExponentialBackoff   bool          `json:"exponential_backoff"`
	BackoffMultiplier    float64       `json:"backoff_multiplier"`
	RetryableErrors      []string      `json:"retryable_errors"`
	RecoveryActions      []string      `json:"recovery_actions"`
	CircuitBreakerConfig *CircuitBreakerConfig `json:"circuit_breaker_config"`
}

// CircuitBreakerConfig defines circuit breaker parameters
type CircuitBreakerConfig struct {
	ErrorThreshold     int           `json:"error_threshold"`
	TimeWindow         time.Duration `json:"time_window"`
	RecoveryTimeout    time.Duration `json:"recovery_timeout"`
	HalfOpenMaxCalls   int           `json:"half_open_max_calls"`
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config       CircuitBreakerConfig
	state        CircuitBreakerState
	errorCount   int
	lastFailTime time.Time
	halfOpenCalls int
	mu           sync.RWMutex
}

// CircuitBreakerState represents the current state of a circuit breaker
type CircuitBreakerState int

const (
	CircuitBreakerClosed CircuitBreakerState = iota
	CircuitBreakerOpen
	CircuitBreakerHalfOpen
)

// ErrorClassifier categorizes errors for appropriate recovery strategies
type ErrorClassifier struct {
	networkErrors    []*regexp.Regexp
	rateLimitErrors  []*regexp.Regexp
	authErrors       []*regexp.Regexp
	storageErrors    []*regexp.Regexp
	configErrors     []*regexp.Regexp
	temporaryErrors  []*regexp.Regexp
	permanentErrors  []*regexp.Regexp
}

// ErrorCategory represents different types of errors
type ErrorCategory string

const (
	ErrorCategoryNetwork    ErrorCategory = "network"
	ErrorCategoryRateLimit  ErrorCategory = "rate_limit"
	ErrorCategoryAuth       ErrorCategory = "authentication"
	ErrorCategoryStorage    ErrorCategory = "storage"
	ErrorCategoryConfig     ErrorCategory = "configuration"
	ErrorCategoryTemporary  ErrorCategory = "temporary"
	ErrorCategoryPermanent  ErrorCategory = "permanent"
	ErrorCategoryUnknown    ErrorCategory = "unknown"
)

// RecoveryResult holds the outcome of a recovery attempt
type RecoveryResult struct {
	Success         bool          `json:"success"`
	AttemptCount    int           `json:"attempt_count"`
	TotalDuration   time.Duration `json:"total_duration"`
	LastError       error         `json:"last_error,omitempty"`
	RecoveryActions []string      `json:"recovery_actions"`
	Suggestions     []string      `json:"suggestions"`
}

// NewErrorRecoveryManager creates a new error recovery manager with default strategies
func NewErrorRecoveryManager() *ErrorRecoveryManager {
	erm := &ErrorRecoveryManager{
		strategies:      make(map[string]RecoveryStrategy),
		circuitBreakers: make(map[string]*CircuitBreaker),
		errorClassifier: NewErrorClassifier(),
	}
	
	// Initialize default recovery strategies
	erm.initializeDefaultStrategies()
	
	return erm
}

// initializeDefaultStrategies sets up default recovery strategies for common error types
func (erm *ErrorRecoveryManager) initializeDefaultStrategies() {
	// Network error strategy - aggressive retry with exponential backoff
	erm.strategies["network"] = RecoveryStrategy{
		Name:               "Network Error Recovery",
		MaxRetryAttempts:   5,
		BaseRetryDelay:     2 * time.Second,
		MaxRetryDelay:      60 * time.Second,
		ExponentialBackoff: true,
		BackoffMultiplier:  2.0,
		RetryableErrors:    []string{"connection refused", "timeout", "network unreachable", "temporary failure"},
		RecoveryActions:    []string{"check_network_connectivity", "verify_dns_resolution", "test_endpoint_availability"},
		CircuitBreakerConfig: &CircuitBreakerConfig{
			ErrorThreshold:   3,
			TimeWindow:       5 * time.Minute,
			RecoveryTimeout:  2 * time.Minute,
			HalfOpenMaxCalls: 2,
		},
	}
	
	// Rate limiting strategy - longer delays, fewer retries
	erm.strategies["rate_limit"] = RecoveryStrategy{
		Name:               "Rate Limit Recovery",
		MaxRetryAttempts:   3,
		BaseRetryDelay:     10 * time.Second,
		MaxRetryDelay:      300 * time.Second,
		ExponentialBackoff: true,
		BackoffMultiplier:  3.0,
		RetryableErrors:    []string{"rate limit", "throttling", "too many requests", "429"},
		RecoveryActions:    []string{"reduce_concurrency", "implement_backoff", "check_rate_limits"},
		CircuitBreakerConfig: &CircuitBreakerConfig{
			ErrorThreshold:   2,
			TimeWindow:       10 * time.Minute,
			RecoveryTimeout:  5 * time.Minute,
			HalfOpenMaxCalls: 1,
		},
	}
	
	// Authentication strategy - limited retries with immediate escalation
	erm.strategies["authentication"] = RecoveryStrategy{
		Name:               "Authentication Error Recovery",
		MaxRetryAttempts:   2,
		BaseRetryDelay:     1 * time.Second,
		MaxRetryDelay:      5 * time.Second,
		ExponentialBackoff: false,
		BackoffMultiplier:  1.0,
		RetryableErrors:    []string{"expired token", "invalid credentials", "unauthorized", "403"},
		RecoveryActions:    []string{"refresh_credentials", "check_permissions", "verify_configuration"},
		CircuitBreakerConfig: &CircuitBreakerConfig{
			ErrorThreshold:   1,
			TimeWindow:       15 * time.Minute,
			RecoveryTimeout:  10 * time.Minute,
			HalfOpenMaxCalls: 1,
		},
	}
	
	// Storage/S3 strategy - moderate retry with storage-specific recovery
	erm.strategies["storage"] = RecoveryStrategy{
		Name:               "Storage Error Recovery",
		MaxRetryAttempts:   4,
		BaseRetryDelay:     3 * time.Second,
		MaxRetryDelay:      120 * time.Second,
		ExponentialBackoff: true,
		BackoffMultiplier:  2.5,
		RetryableErrors:    []string{"internal error", "service unavailable", "slow down", "500", "502", "503"},
		RecoveryActions:    []string{"verify_bucket_access", "check_storage_quotas", "validate_object_permissions"},
		CircuitBreakerConfig: &CircuitBreakerConfig{
			ErrorThreshold:   4,
			TimeWindow:       10 * time.Minute,
			RecoveryTimeout:  3 * time.Minute,
			HalfOpenMaxCalls: 2,
		},
	}
	
	// Configuration strategy - no retries, immediate escalation
	erm.strategies["configuration"] = RecoveryStrategy{
		Name:               "Configuration Error Recovery",
		MaxRetryAttempts:   0,
		BaseRetryDelay:     0,
		MaxRetryDelay:      0,
		ExponentialBackoff: false,
		BackoffMultiplier:  1.0,
		RetryableErrors:    []string{}, // Configuration errors are not retryable
		RecoveryActions:    []string{"validate_configuration", "check_file_paths", "verify_settings"},
	}
	
	// Default/fallback strategy
	erm.strategies["default"] = RecoveryStrategy{
		Name:               "Default Error Recovery",
		MaxRetryAttempts:   3,
		BaseRetryDelay:     5 * time.Second,
		MaxRetryDelay:      30 * time.Second,
		ExponentialBackoff: true,
		BackoffMultiplier:  2.0,
		RetryableErrors:    []string{"error", "failed", "timeout"},
		RecoveryActions:    []string{"general_recovery", "check_system_status"},
		CircuitBreakerConfig: &CircuitBreakerConfig{
			ErrorThreshold:   3,
			TimeWindow:       5 * time.Minute,
			RecoveryTimeout:  2 * time.Minute,
			HalfOpenMaxCalls: 2,
		},
	}
}

// ExecuteWithRecovery executes a function with comprehensive error recovery
func (erm *ErrorRecoveryManager) ExecuteWithRecovery(
	ctx context.Context,
	operationName string,
	operation func() error,
) *RecoveryResult {
	startTime := time.Now()
	result := &RecoveryResult{
		RecoveryActions: []string{},
		Suggestions:     []string{},
	}
	
	for attempt := 0; attempt < 10; attempt++ { // Max 10 attempts across all strategies
		// Check circuit breaker
		if cb := erm.getCircuitBreaker(operationName); cb != nil {
			if !cb.AllowRequest() {
				result.LastError = fmt.Errorf("circuit breaker is open for operation: %s", operationName)
				result.Suggestions = append(result.Suggestions, "Circuit breaker is open. Wait for recovery timeout or investigate underlying issues.")
				break
			}
		}
		
		// Execute operation
		err := operation()
		
		if err == nil {
			// Success - record success in circuit breaker
			if cb := erm.getCircuitBreaker(operationName); cb != nil {
				cb.RecordSuccess()
			}
			result.Success = true
			result.AttemptCount = attempt + 1
			result.TotalDuration = time.Since(startTime)
			return result
		}
		
		// Record error in circuit breaker
		if cb := erm.getCircuitBreaker(operationName); cb != nil {
			cb.RecordError()
		}
		
		result.LastError = err
		result.AttemptCount = attempt + 1
		
		// Classify error and get recovery strategy
		category := erm.errorClassifier.ClassifyError(err)
		strategy := erm.getRecoveryStrategy(string(category))
		
		// Check if we should retry
		if attempt >= strategy.MaxRetryAttempts {
			break
		}
		
		// Check if error is retryable
		if !erm.isRetryableError(err, strategy) {
			result.Suggestions = append(result.Suggestions, 
				fmt.Sprintf("Error is not retryable: %v", err))
			break
		}
		
		// Calculate retry delay
		delay := erm.calculateRetryDelay(attempt, strategy)
		
		result.RecoveryActions = append(result.RecoveryActions, 
			fmt.Sprintf("Attempt %d: Retrying after %v due to %s error", 
				attempt+1, delay, category))
		
		// Add recovery suggestions
		result.Suggestions = append(result.Suggestions, strategy.RecoveryActions...)
		
		// Wait before retry (with context cancellation)
		select {
		case <-ctx.Done():
			result.LastError = ctx.Err()
			result.Suggestions = append(result.Suggestions, "Operation cancelled by context")
			break
		case <-time.After(delay):
			// Continue to next attempt
		}
	}
	
	result.TotalDuration = time.Since(startTime)
	
	// Add final suggestions if all retries failed
	if !result.Success {
		result.Suggestions = append(result.Suggestions, erm.generateFailureSuggestions(result.LastError)...)
	}
	
	return result
}

// NewErrorClassifier creates a new error classifier with predefined patterns
func NewErrorClassifier() *ErrorClassifier {
	return &ErrorClassifier{
		networkErrors: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(connection.*refused|network.*unreachable|timeout|no route to host)`),
			regexp.MustCompile(`(?i)(dial tcp.*connect.*refused|i/o timeout)`),
			regexp.MustCompile(`(?i)(temporary failure|name resolution|dns)`),
		},
		rateLimitErrors: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(rate limit|throttling|too many requests)`),
			regexp.MustCompile(`(?i)(429|slow down|quota exceeded)`),
		},
		authErrors: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(unauthorized|forbidden|access denied)`),
			regexp.MustCompile(`(?i)(401|403|invalid.*credentials|expired.*token)`),
		},
		storageErrors: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(internal error|service unavailable|bad gateway)`),
			regexp.MustCompile(`(?i)(500|502|503|504|bucket.*not.*found)`),
		},
		configErrors: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(file not found|no such file|invalid.*configuration)`),
			regexp.MustCompile(`(?i)(bad.*format|parse.*error|invalid.*syntax)`),
		},
		temporaryErrors: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(temporary|try again|retry)`),
			regexp.MustCompile(`(?i)(transient|recoverable)`),
		},
		permanentErrors: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(not found|does not exist|permission denied)`),
			regexp.MustCompile(`(?i)(404|invalid.*request|malformed)`),
		},
	}
}

// ClassifyError determines the category of an error
func (ec *ErrorClassifier) ClassifyError(err error) ErrorCategory {
	if err == nil {
		return ErrorCategoryUnknown
	}
	
	errMsg := strings.ToLower(err.Error())
	
	// Check each category in order of specificity
	if ec.matchesPatterns(errMsg, ec.configErrors) {
		return ErrorCategoryConfig
	}
	if ec.matchesPatterns(errMsg, ec.authErrors) {
		return ErrorCategoryAuth
	}
	if ec.matchesPatterns(errMsg, ec.rateLimitErrors) {
		return ErrorCategoryRateLimit
	}
	if ec.matchesPatterns(errMsg, ec.networkErrors) {
		return ErrorCategoryNetwork
	}
	if ec.matchesPatterns(errMsg, ec.storageErrors) {
		return ErrorCategoryStorage
	}
	if ec.matchesPatterns(errMsg, ec.permanentErrors) {
		return ErrorCategoryPermanent
	}
	if ec.matchesPatterns(errMsg, ec.temporaryErrors) {
		return ErrorCategoryTemporary
	}
	
	return ErrorCategoryUnknown
}

// matchesPatterns checks if error message matches any of the given patterns
func (ec *ErrorClassifier) matchesPatterns(errMsg string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(errMsg) {
			return true
		}
	}
	return false
}

// getRecoveryStrategy returns the appropriate recovery strategy for an error category
func (erm *ErrorRecoveryManager) getRecoveryStrategy(category string) RecoveryStrategy {
	erm.mu.RLock()
	defer erm.mu.RUnlock()
	
	if strategy, exists := erm.strategies[category]; exists {
		return strategy
	}
	return erm.strategies["default"]
}

// getCircuitBreaker returns or creates a circuit breaker for an operation
func (erm *ErrorRecoveryManager) getCircuitBreaker(operationName string) *CircuitBreaker {
	erm.mu.Lock()
	defer erm.mu.Unlock()
	
	if cb, exists := erm.circuitBreakers[operationName]; exists {
		return cb
	}
	
	// Create new circuit breaker with default config
	strategy := erm.strategies["default"]
	if strategy.CircuitBreakerConfig != nil {
		cb := NewCircuitBreaker(*strategy.CircuitBreakerConfig)
		erm.circuitBreakers[operationName] = cb
		return cb
	}
	
	return nil
}

// isRetryableError checks if an error should be retried based on the strategy
func (erm *ErrorRecoveryManager) isRetryableError(err error, strategy RecoveryStrategy) bool {
	if len(strategy.RetryableErrors) == 0 {
		return false
	}
	
	errMsg := strings.ToLower(err.Error())
	for _, retryablePattern := range strategy.RetryableErrors {
		if strings.Contains(errMsg, strings.ToLower(retryablePattern)) {
			return true
		}
	}
	
	return false
}

// calculateRetryDelay calculates the delay before the next retry attempt
func (erm *ErrorRecoveryManager) calculateRetryDelay(attempt int, strategy RecoveryStrategy) time.Duration {
	if strategy.BaseRetryDelay == 0 {
		return 0
	}
	
	delay := strategy.BaseRetryDelay
	
	if strategy.ExponentialBackoff && attempt > 0 {
		multiplier := math.Pow(strategy.BackoffMultiplier, float64(attempt))
		delay = time.Duration(float64(delay) * multiplier)
	}
	
	if strategy.MaxRetryDelay > 0 && delay > strategy.MaxRetryDelay {
		delay = strategy.MaxRetryDelay
	}
	
	return delay
}

// generateFailureSuggestions provides actionable suggestions for persistent failures
func (erm *ErrorRecoveryManager) generateFailureSuggestions(err error) []string {
	if err == nil {
		return []string{}
	}
	
	category := erm.errorClassifier.ClassifyError(err)
	suggestions := []string{}
	
	switch category {
	case ErrorCategoryNetwork:
		suggestions = []string{
			"Check network connectivity and DNS resolution",
			"Verify firewall settings and proxy configuration",
			"Consider using a different network or VPN",
			"Check if the target service is experiencing outages",
		}
	case ErrorCategoryRateLimit:
		suggestions = []string{
			"Reduce transfer concurrency to stay within rate limits",
			"Implement exponential backoff with jitter",
			"Consider upgrading to a higher tier service plan",
			"Distribute load across multiple time periods",
		}
	case ErrorCategoryAuth:
		suggestions = []string{
			"Verify AWS credentials are correct and up-to-date",
			"Check IAM permissions for the required S3 operations",
			"Ensure credentials have not expired",
			"Validate AWS region configuration",
		}
	case ErrorCategoryStorage:
		suggestions = []string{
			"Verify S3 bucket exists and is accessible",
			"Check bucket permissions and policies",
			"Ensure sufficient storage quota is available",
			"Validate object key names and paths",
		}
	case ErrorCategoryConfig:
		suggestions = []string{
			"Review configuration file syntax and format",
			"Verify all required fields are present",
			"Check file paths and directory permissions",
			"Validate YAML/JSON structure",
		}
	default:
		suggestions = []string{
			"Check system logs for additional error details",
			"Verify all prerequisites are met",
			"Consider contacting support if the issue persists",
			"Try running the operation with verbose logging enabled",
		}
	}
	
	return suggestions
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  CircuitBreakerClosed,
	}
}

// AllowRequest determines if a request should be allowed through the circuit breaker
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	switch cb.state {
	case CircuitBreakerClosed:
		return true
	case CircuitBreakerOpen:
		if time.Since(cb.lastFailTime) >= cb.config.RecoveryTimeout {
			cb.state = CircuitBreakerHalfOpen
			cb.halfOpenCalls = 0
			return true
		}
		return false
	case CircuitBreakerHalfOpen:
		return cb.halfOpenCalls < cb.config.HalfOpenMaxCalls
	default:
		return false
	}
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	switch cb.state {
	case CircuitBreakerClosed:
		cb.errorCount = 0
	case CircuitBreakerHalfOpen:
		cb.halfOpenCalls++
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			cb.state = CircuitBreakerClosed
			cb.errorCount = 0
		}
	}
}

// RecordError records a failed operation
func (cb *CircuitBreaker) RecordError() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.errorCount++
	cb.lastFailTime = time.Now()
	
	switch cb.state {
	case CircuitBreakerClosed:
		if cb.errorCount >= cb.config.ErrorThreshold {
			cb.state = CircuitBreakerOpen
		}
	case CircuitBreakerHalfOpen:
		cb.state = CircuitBreakerOpen
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}