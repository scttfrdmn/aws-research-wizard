package validator

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/runner"
)

// Config holds configuration for validation
type Config struct {
	Format string // Documentation format to validate
	Strict bool   // Enable strict validation rules
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Success  bool                   `json:"success"`
	Errors   []string               `json:"errors"`
	Warnings []string               `json:"warnings"`
	Summary  ValidationSummary      `json:"summary"`
	Details  map[string]interface{} `json:"details"`
}

// ValidationSummary provides aggregate validation statistics
type ValidationSummary struct {
	TotalTests   int           `json:"total_tests"`
	PassedTests  int           `json:"passed_tests"`
	FailedTests  int           `json:"failed_tests"`
	SkippedTests int           `json:"skipped_tests"`
	WarningCount int           `json:"warning_count"`
	AvgDuration  time.Duration `json:"avg_duration"`
	MaxDuration  time.Duration `json:"max_duration"`
	MinDuration  time.Duration `json:"min_duration"`
}

// Validator handles validation of test results and documentation
type Validator struct {
	config Config
	rules  []ValidationRule
}

// ValidationRule defines a validation rule
type ValidationRule struct {
	Name        string
	Description string
	Severity    string // error, warning, info
	Check       func(*runner.TestResult) (bool, string)
}

// New creates a new validator with default rules
func New() *Validator {
	v := &Validator{
		rules: []ValidationRule{},
	}

	v.addDefaultRules()
	return v
}

// ValidateResults validates a collection of test results
func (v *Validator) ValidateResults(results []*runner.TestResult) error {
	validationResult := &ValidationResult{
		Success: true,
		Details: make(map[string]interface{}),
	}

	var totalDuration time.Duration
	var maxDuration time.Duration
	minDuration := time.Duration(^uint64(0) >> 1) // Max duration

	// Validate each test result
	for _, result := range results {
		v.validateSingleResult(result, validationResult)

		// Collect duration statistics
		totalDuration += result.Duration
		if result.Duration > maxDuration {
			maxDuration = result.Duration
		}
		if result.Duration < minDuration {
			minDuration = result.Duration
		}
	}

	// Calculate summary
	summary := ValidationSummary{
		TotalTests:   len(results),
		WarningCount: len(validationResult.Warnings),
		MaxDuration:  maxDuration,
		MinDuration:  minDuration,
	}

	if len(results) > 0 {
		summary.AvgDuration = totalDuration / time.Duration(len(results))
	}

	// Count passed/failed/skipped
	for _, result := range results {
		if result.Success {
			summary.PassedTests++
		} else {
			if skipped, exists := result.Metadata["status"]; exists && skipped == "skipped" {
				summary.SkippedTests++
			} else {
				summary.FailedTests++
			}
		}
	}

	validationResult.Summary = summary

	// Overall success depends on whether we have failures
	if summary.FailedTests > 0 {
		validationResult.Success = false
	}

	// Report results
	return v.reportValidationResults(validationResult)
}

// ValidatePath validates documentation at a given path
func (v *Validator) ValidatePath(path string, config Config) error {
	// TODO: Implement path validation
	// This would validate:
	// - Documentation format consistency
	// - Code block syntax
	// - Link validity
	// - Example completeness
	// - Best practices compliance

	fmt.Printf("Validating documentation at: %s\n", path)
	fmt.Printf("Format: %s, Strict: %v\n", config.Format, config.Strict)

	return nil
}

// validateSingleResult validates a single test result against all rules
func (v *Validator) validateSingleResult(result *runner.TestResult, validationResult *ValidationResult) {
	for _, rule := range v.rules {
		passed, message := rule.Check(result)

		if !passed {
			switch rule.Severity {
			case "error":
				validationResult.Errors = append(validationResult.Errors,
					fmt.Sprintf("%s: %s", rule.Name, message))
			case "warning":
				validationResult.Warnings = append(validationResult.Warnings,
					fmt.Sprintf("%s: %s", rule.Name, message))
			}
		}
	}
}

// addDefaultRules adds default validation rules
func (v *Validator) addDefaultRules() {
	// Rule: Test should complete within reasonable time
	v.rules = append(v.rules, ValidationRule{
		Name:        "Execution Time",
		Description: "Test should complete within reasonable time",
		Severity:    "warning",
		Check: func(result *runner.TestResult) (bool, string) {
			if result.Duration > 5*time.Minute {
				return false, fmt.Sprintf("Test took %v, which is longer than expected (5m)", result.Duration)
			}
			return true, ""
		},
	})

	// Rule: Successful tests should have zero exit code
	v.rules = append(v.rules, ValidationRule{
		Name:        "Exit Code Consistency",
		Description: "Successful tests should have exit code 0",
		Severity:    "error",
		Check: func(result *runner.TestResult) (bool, string) {
			if result.Success && result.ExitCode != 0 {
				return false, fmt.Sprintf("Test marked as successful but has non-zero exit code: %d", result.ExitCode)
			}
			return true, ""
		},
	})

	// Rule: Failed tests should have error information
	v.rules = append(v.rules, ValidationRule{
		Name:        "Error Information",
		Description: "Failed tests should provide error details",
		Severity:    "warning",
		Check: func(result *runner.TestResult) (bool, string) {
			if !result.Success && len(result.Errors) == 0 && result.ErrorOutput == "" {
				return false, "Failed test has no error information"
			}
			return true, ""
		},
	})

	// Rule: Check for security issues in code
	v.rules = append(v.rules, ValidationRule{
		Name:        "Security Check",
		Description: "Code should not contain obvious security issues",
		Severity:    "warning",
		Check: func(result *runner.TestResult) (bool, string) {
			// Check output for signs of exposed credentials or unsafe operations
			output := strings.ToLower(result.Output)

			securityIssues := []string{
				"password=",
				"secret=",
				"token=",
				"api_key=",
				"private_key",
				"rm -rf /",
				"chmod 777",
			}

			for _, issue := range securityIssues {
				if strings.Contains(output, issue) {
					return false, fmt.Sprintf("Potential security issue detected: %s", issue)
				}
			}
			return true, ""
		},
	})

	// Rule: Resource cleanup verification
	v.rules = append(v.rules, ValidationRule{
		Name:        "Resource Cleanup",
		Description: "Tests should clean up created resources",
		Severity:    "warning",
		Check: func(result *runner.TestResult) (bool, string) {
			if len(result.Resources) > 5 {
				return false, fmt.Sprintf("Test created %d resources - consider cleanup", len(result.Resources))
			}
			return true, ""
		},
	})

	// Rule: Output should not be empty for successful tests
	v.rules = append(v.rules, ValidationRule{
		Name:        "Output Presence",
		Description: "Successful tests should produce some output",
		Severity:    "info",
		Check: func(result *runner.TestResult) (bool, string) {
			if result.Success && strings.TrimSpace(result.Output) == "" {
				return false, "Successful test produced no output"
			}
			return true, ""
		},
	})

	// Rule: Check for common error patterns
	v.rules = append(v.rules, ValidationRule{
		Name:        "Common Error Patterns",
		Description: "Check for common error patterns in output",
		Severity:    "error",
		Check: func(result *runner.TestResult) (bool, string) {
			output := strings.ToLower(result.Output + " " + result.ErrorOutput)

			errorPatterns := []string{
				"command not found",
				"permission denied",
				"no such file or directory",
				"connection refused",
				"timeout",
			}

			for _, pattern := range errorPatterns {
				if strings.Contains(output, pattern) {
					return false, fmt.Sprintf("Common error pattern detected: %s", pattern)
				}
			}
			return true, ""
		},
	})
}

// reportValidationResults reports the validation results
func (v *Validator) reportValidationResults(result *ValidationResult) error {
	fmt.Printf("\n=== Validation Results ===\n")
	fmt.Printf("Overall Success: %v\n", result.Success)
	fmt.Printf("Total Tests: %d\n", result.Summary.TotalTests)
	fmt.Printf("Passed: %d\n", result.Summary.PassedTests)
	fmt.Printf("Failed: %d\n", result.Summary.FailedTests)
	fmt.Printf("Skipped: %d\n", result.Summary.SkippedTests)
	fmt.Printf("Warnings: %d\n", result.Summary.WarningCount)

	if result.Summary.TotalTests > 0 {
		fmt.Printf("Average Duration: %v\n", result.Summary.AvgDuration)
		fmt.Printf("Max Duration: %v\n", result.Summary.MaxDuration)
		fmt.Printf("Min Duration: %v\n", result.Summary.MinDuration)
	}

	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range result.Errors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range result.Warnings {
			fmt.Printf("  ⚠️  %s\n", warning)
		}
	}

	fmt.Printf("\n")

	if !result.Success {
		return fmt.Errorf("validation failed with %d errors", len(result.Errors))
	}

	return nil
}

// AddRule adds a custom validation rule
func (v *Validator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// SetConfig sets the validator configuration
func (v *Validator) SetConfig(config Config) {
	v.config = config
}
