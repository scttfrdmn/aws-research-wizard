package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// New creates a new test runner with the given configuration
func New(config Config) Runner {
	// Set defaults
	if config.Timeout == "" {
		config.Timeout = "5m"
	}
	if config.Parallel == 0 {
		config.Parallel = 1
	}
	if config.WorkDir == "" {
		config.WorkDir = "/tmp/tutorial-guard"
	}
	if config.RetryCount == 0 {
		config.RetryCount = 1
	}

	var env Environment
	switch config.Environment {
	case "docker":
		env = NewDockerEnvironment(config)
	case "local":
		env = NewLocalEnvironment(config)
	case "aws":
		env = NewAWSEnvironment(config)
	default:
		env = NewLocalEnvironment(config) // Default to local
	}

	return &TestRunner{
		config:      config,
		environment: env,
	}
}

// RunExample executes a single code example
func (r *TestRunner) RunExample(ctx context.Context, example extractor.Example) (*TestResult, error) {
	startTime := time.Now()

	result := &TestResult{
		ExampleID:   example.ID,
		StartTime:   startTime,
		Environment: r.config.Environment,
		Metadata:    make(map[string]string),
	}

	// Parse timeout
	timeout, err := time.ParseDuration(r.config.Timeout)
	if err != nil {
		timeout = 5 * time.Minute
	}

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Setup environment if needed
	if err := r.environment.Setup(execCtx); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Environment setup failed: %v", err))
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Execute the example with retries
	var lastErr error
	for attempt := 0; attempt < r.config.RetryCount; attempt++ {
		if attempt > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Retry attempt %d", attempt))
		}

		testResult, err := r.environment.Execute(execCtx, example)
		if err != nil {
			lastErr = err
			continue
		}

		// Copy results
		*result = *testResult
		break
	}

	if lastErr != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("All retry attempts failed: %v", lastErr))
	}

	// Collect resources
	result.Resources = r.environment.GetResources()

	// Set timing
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Add metadata
	result.Metadata["example_source"] = example.Source
	result.Metadata["example_language"] = example.Language
	result.Metadata["retry_count"] = fmt.Sprintf("%d", r.config.RetryCount)

	return result, nil
}

// RunExamples executes multiple code examples
func (r *TestRunner) RunExamples(ctx context.Context, examples []extractor.Example) (*TestSuite, error) {
	startTime := time.Now()

	suite := &TestSuite{
		Name:        "Tutorial Examples",
		Description: fmt.Sprintf("Execution of %d extracted examples", len(examples)),
		Environment: r.config.Environment,
		StartTime:   startTime,
		Metadata:    make(map[string]string),
	}

	// Filter examples by language support
	var validExamples []extractor.Example
	for _, example := range examples {
		if r.isLanguageSupported(example.Language) {
			validExamples = append(validExamples, example)
		} else {
			// Add skipped result
			skippedResult := &TestResult{
				ExampleID:   example.ID,
				Success:     false,
				StartTime:   time.Now(),
				EndTime:     time.Now(),
				Environment: r.config.Environment,
				Errors:      []string{fmt.Sprintf("Unsupported language: %s", example.Language)},
				Metadata:    map[string]string{"status": "skipped"},
			}
			suite.Results = append(suite.Results, *skippedResult)
		}
	}

	// Execute examples in parallel
	results := make(chan *TestResult, len(validExamples))
	var wg sync.WaitGroup

	// Create worker pool
	semaphore := make(chan struct{}, r.config.Parallel)

	for _, example := range validExamples {
		wg.Add(1)
		go func(ex extractor.Example) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result, err := r.RunExample(ctx, ex)
			if err != nil && result == nil {
				// Create error result if RunExample failed completely
				result = &TestResult{
					ExampleID:   ex.ID,
					Success:     false,
					StartTime:   time.Now(),
					EndTime:     time.Now(),
					Environment: r.config.Environment,
					Errors:      []string{err.Error()},
					Metadata:    make(map[string]string),
				}
			}

			results <- result

			// Stop on first failure if FailFast is enabled
			if r.config.FailFast && !result.Success {
				// Cancel context to stop other tests
				// Note: This is a simple implementation - could be improved
			}
		}(example)
	}

	// Wait for all tests to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		suite.Results = append(suite.Results, *result)
	}

	// Calculate summary
	suite.Summary = r.calculateSummary(suite.Results)

	// Set timing
	suite.EndTime = time.Now()
	suite.Duration = suite.EndTime.Sub(suite.StartTime)

	// Add metadata
	suite.Metadata["total_examples"] = fmt.Sprintf("%d", len(examples))
	suite.Metadata["valid_examples"] = fmt.Sprintf("%d", len(validExamples))
	suite.Metadata["parallel_workers"] = fmt.Sprintf("%d", r.config.Parallel)

	return suite, nil
}

// RunExamplesFromFile loads examples from a file and executes them
func (r *TestRunner) RunExamplesFromFile(filename string) (*TestSuite, error) {
	// Read examples file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read examples file: %w", err)
	}

	var exampleSet extractor.ExampleSet
	if err := json.Unmarshal(data, &exampleSet); err != nil {
		return nil, fmt.Errorf("failed to parse examples file: %w", err)
	}

	ctx := context.Background()
	return r.RunExamples(ctx, exampleSet.Examples)
}

// Cleanup performs cleanup of resources
func (r *TestRunner) Cleanup() error {
	if r.config.Cleanup {
		return r.environment.Cleanup(context.Background())
	}
	return nil
}

// calculateSummary calculates aggregate statistics from test results
func (r *TestRunner) calculateSummary(results []TestResult) TestSummary {
	summary := TestSummary{}

	for _, result := range results {
		summary.Total++

		if result.Success {
			summary.Passed++
		} else {
			// Check if it was skipped
			if skipped, exists := result.Metadata["status"]; exists && skipped == "skipped" {
				summary.Skipped++
			} else {
				summary.Failed++
			}
		}

		if len(result.Warnings) > 0 {
			summary.Warnings++
		}
	}

	return summary
}

// isLanguageSupported checks if a language is supported for execution
func (r *TestRunner) isLanguageSupported(language string) bool {
	supportedLanguages := map[string]bool{
		"bash":       true,
		"sh":         true,
		"shell":      true,
		"go":         true,
		"python":     true,
		"javascript": false, // Not implemented yet
		"yaml":       false, // Not executable
		"json":       false, // Not executable
		"dockerfile": false, // Special handling needed
	}

	supported, exists := supportedLanguages[language]
	return exists && supported
}

// Helper function to ensure working directory exists
func ensureWorkDir(workDir string) error {
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		return os.MkdirAll(workDir, 0755)
	}
	return nil
}
