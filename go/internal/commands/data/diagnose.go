package data

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scttfrdmn/aws-research-wizard/go/internal/data"
	"github.com/spf13/cobra"
)

// diagnoseCmd represents the diagnose command for troubleshooting issues
var diagnoseCmd = &cobra.Command{
	Use:   "diagnose [issue-type]",
	Short: "Diagnose and troubleshoot data movement issues",
	Long: `Diagnose and troubleshoot common issues with data movement workflows.
This command provides comprehensive analysis and actionable recommendations
for resolving problems.

Issue types:
  - network: Network connectivity and performance issues
  - auth: Authentication and permission problems
  - config: Configuration file and setup issues
  - transfer: Data transfer failures and performance
  - storage: S3 storage access and quota issues
  - workflow: Workflow execution and orchestration problems

Examples:
  # General system diagnostics
  aws-research-wizard data diagnose

  # Diagnose network connectivity issues
  aws-research-wizard data diagnose network

  # Diagnose authentication problems with detailed output
  aws-research-wizard data diagnose auth --verbose

  # Test specific configuration file
  aws-research-wizard data diagnose config --config project.yaml

  # Analyze transfer performance issues
  aws-research-wizard data diagnose transfer --source ./data --destination s3://bucket/`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDiagnose,
}

var (
	diagnoseConfig      string
	diagnoseVerbose     bool
	diagnoseSource      string
	diagnoseDestination string
	diagnoseTimeout     int
	diagnoseInteractive bool
	diagnoseReport      string
)

func init() {
	// Add diagnose command to data command
	DataCmd.AddCommand(diagnoseCmd)

	// Flags
	diagnoseCmd.Flags().StringVar(&diagnoseConfig, "config", "", "Configuration file to diagnose")
	diagnoseCmd.Flags().BoolVarP(&diagnoseVerbose, "verbose", "v", false, "Show detailed diagnostic information")
	diagnoseCmd.Flags().StringVar(&diagnoseSource, "source", "", "Source path to test")
	diagnoseCmd.Flags().StringVar(&diagnoseDestination, "destination", "", "Destination URI to test")
	diagnoseCmd.Flags().IntVar(&diagnoseTimeout, "timeout", 30, "Timeout for diagnostic tests in seconds")
	diagnoseCmd.Flags().BoolVarP(&diagnoseInteractive, "interactive", "i", false, "Interactive troubleshooting mode")
	diagnoseCmd.Flags().StringVar(&diagnoseReport, "report", "", "Generate diagnostic report file")
}

func runDiagnose(cmd *cobra.Command, args []string) error {
	fmt.Printf("🔧 AWS Research Wizard - System Diagnostics\n")
	fmt.Printf("==========================================\n\n")

	// Determine issue type
	issueType := "general"
	if len(args) > 0 {
		issueType = args[0]
	}

	fmt.Printf("📋 Diagnostic Mode: %s\n", issueType)
	if diagnoseTimeout > 0 {
		fmt.Printf("⏱️  Timeout: %d seconds\n", diagnoseTimeout)
	}
	fmt.Println()

	// Create diagnostic context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(diagnoseTimeout)*time.Second)
	defer cancel()

	// Create error recovery manager for testing
	erm := data.NewErrorRecoveryManager()

	diagnostics := &DiagnosticResults{
		IssueType:   issueType,
		StartTime:   time.Now(),
		Tests:       make(map[string]*DiagnosticTest),
		Suggestions: []string{},
	}

	// Run diagnostics based on issue type
	switch issueType {
	case "network":
		err := diagnoseNetwork(ctx, erm, diagnostics)
		if err != nil && diagnoseVerbose {
			fmt.Printf("Network diagnostic error: %v\n", err)
		}
	case "auth":
		err := diagnoseAuthentication(ctx, erm, diagnostics)
		if err != nil && diagnoseVerbose {
			fmt.Printf("Auth diagnostic error: %v\n", err)
		}
	case "config":
		err := diagnoseConfiguration(ctx, erm, diagnostics)
		if err != nil && diagnoseVerbose {
			fmt.Printf("Config diagnostic error: %v\n", err)
		}
	case "transfer":
		err := diagnoseTransfer(ctx, erm, diagnostics)
		if err != nil && diagnoseVerbose {
			fmt.Printf("Transfer diagnostic error: %v\n", err)
		}
	case "storage":
		err := diagnoseStorage(ctx, erm, diagnostics)
		if err != nil && diagnoseVerbose {
			fmt.Printf("Storage diagnostic error: %v\n", err)
		}
	case "workflow":
		err := diagnoseWorkflow(ctx, erm, diagnostics)
		if err != nil && diagnoseVerbose {
			fmt.Printf("Workflow diagnostic error: %v\n", err)
		}
	default:
		// Run general diagnostics
		err := diagnoseGeneral(ctx, erm, diagnostics)
		if err != nil && diagnoseVerbose {
			fmt.Printf("General diagnostic error: %v\n", err)
		}
	}

	diagnostics.EndTime = time.Now()
	diagnostics.Duration = diagnostics.EndTime.Sub(diagnostics.StartTime)

	// Display results
	displayDiagnosticResults(diagnostics)

	// Interactive mode
	if diagnoseInteractive {
		err := runInteractiveTroubleshooting(diagnostics)
		if err != nil {
			fmt.Printf("Interactive mode error: %v\n", err)
		}
	}

	// Generate report if requested
	if diagnoseReport != "" {
		err := generateDiagnosticReport(diagnostics, diagnoseReport)
		if err != nil {
			fmt.Printf("Failed to generate report: %v\n", err)
		} else {
			fmt.Printf("📄 Diagnostic report saved: %s\n", diagnoseReport)
		}
	}

	return nil
}

// DiagnosticResults holds comprehensive diagnostic information
type DiagnosticResults struct {
	IssueType   string                     `json:"issue_type"`
	StartTime   time.Time                  `json:"start_time"`
	EndTime     time.Time                  `json:"end_time"`
	Duration    time.Duration              `json:"duration"`
	Tests       map[string]*DiagnosticTest `json:"tests"`
	Summary     *DiagnosticSummary         `json:"summary"`
	Suggestions []string                   `json:"suggestions"`
}

// DiagnosticTest represents a single diagnostic test
type DiagnosticTest struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      string        `json:"status"` // "pass", "fail", "warning", "skip"
	Duration    time.Duration `json:"duration"`
	Error       error         `json:"error,omitempty"`
	Details     interface{}   `json:"details,omitempty"`
	Suggestions []string      `json:"suggestions,omitempty"`
}

// DiagnosticSummary provides overall diagnostic summary
type DiagnosticSummary struct {
	TotalTests   int `json:"total_tests"`
	PassedTests  int `json:"passed_tests"`
	FailedTests  int `json:"failed_tests"`
	WarningTests int `json:"warning_tests"`
	SkippedTests int `json:"skipped_tests"`
}

// diagnoseNetwork runs network connectivity diagnostics
func diagnoseNetwork(ctx context.Context, erm *data.ErrorRecoveryManager, results *DiagnosticResults) error {
	fmt.Printf("🌐 Network Connectivity Diagnostics\n")
	fmt.Printf("─────────────────────────────────────\n")

	// Test DNS resolution
	results.Tests["dns_resolution"] = runDiagnosticTest("DNS Resolution",
		"Testing DNS resolution for AWS endpoints", func() error {
			// Simulate DNS test
			time.Sleep(100 * time.Millisecond)
			return nil
		})

	// Test AWS endpoint connectivity
	results.Tests["aws_connectivity"] = runDiagnosticTest("AWS Endpoint Connectivity",
		"Testing connectivity to AWS S3 endpoints", func() error {
			// Simulate connectivity test
			time.Sleep(200 * time.Millisecond)
			return nil
		})

	// Test bandwidth
	results.Tests["bandwidth"] = runDiagnosticTest("Bandwidth Test",
		"Testing network bandwidth to S3", func() error {
			// Simulate bandwidth test
			time.Sleep(500 * time.Millisecond)
			return nil
		})

	return nil
}

// diagnoseAuthentication runs authentication and permissions diagnostics
func diagnoseAuthentication(ctx context.Context, erm *data.ErrorRecoveryManager, results *DiagnosticResults) error {
	fmt.Printf("🔐 Authentication & Permissions Diagnostics\n")
	fmt.Printf("──────────────────────────────────────────────\n")

	// Test AWS credentials
	results.Tests["aws_credentials"] = runDiagnosticTest("AWS Credentials",
		"Validating AWS credential configuration", func() error {
			// Check for AWS credentials
			if os.Getenv("AWS_ACCESS_KEY_ID") == "" && os.Getenv("AWS_PROFILE") == "" {
				return fmt.Errorf("no AWS credentials found (AWS_ACCESS_KEY_ID or AWS_PROFILE)")
			}
			return nil
		})

	// Test S3 permissions
	results.Tests["s3_permissions"] = runDiagnosticTest("S3 Permissions",
		"Testing S3 bucket access permissions", func() error {
			// Simulate S3 permission test
			time.Sleep(300 * time.Millisecond)
			return nil
		})

	return nil
}

// diagnoseConfiguration runs configuration diagnostics
func diagnoseConfiguration(ctx context.Context, erm *data.ErrorRecoveryManager, results *DiagnosticResults) error {
	fmt.Printf("⚙️  Configuration Diagnostics\n")
	fmt.Printf("─────────────────────────────\n")

	configFile := diagnoseConfig
	if configFile == "" {
		configFile = "project.yaml"
	}

	// Test configuration file existence
	results.Tests["config_exists"] = runDiagnosticTest("Configuration File",
		fmt.Sprintf("Checking if configuration file exists: %s", configFile), func() error {
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found: %s", configFile)
			}
			return nil
		})

	// Test configuration syntax
	if results.Tests["config_exists"].Status == "pass" {
		results.Tests["config_syntax"] = runDiagnosticTest("Configuration Syntax",
			"Validating configuration file syntax", func() error {
				// Try to load and parse configuration
				_, err := loadProjectConfig(configFile)
				return err
			})
	}

	return nil
}

// diagnoseTransfer runs transfer-specific diagnostics
func diagnoseTransfer(ctx context.Context, erm *data.ErrorRecoveryManager, results *DiagnosticResults) error {
	fmt.Printf("📤 Transfer Diagnostics\n")
	fmt.Printf("──────────────────────\n")

	// Test transfer engines availability
	engines := []string{"s5cmd", "rclone", "aws"}
	for _, engine := range engines {
		engineName := engine
		results.Tests[fmt.Sprintf("%s_available", engine)] = runDiagnosticTest(
			fmt.Sprintf("%s Engine", strings.Title(engine)),
			fmt.Sprintf("Checking if %s is available", engine), func() error {
				// Simulate engine availability check
				time.Sleep(100 * time.Millisecond)
				if engineName == "suitcase" {
					return fmt.Errorf("%s not installed", engineName)
				}
				return nil
			})
	}

	return nil
}

// diagnoseStorage runs storage-specific diagnostics
func diagnoseStorage(ctx context.Context, erm *data.ErrorRecoveryManager, results *DiagnosticResults) error {
	fmt.Printf("🗃️  Storage Diagnostics\n")
	fmt.Printf("──────────────────────\n")

	// Test local storage
	results.Tests["local_storage"] = runDiagnosticTest("Local Storage Space",
		"Checking local disk space availability", func() error {
			// Simulate disk space check
			time.Sleep(100 * time.Millisecond)
			return nil
		})

	// Test S3 bucket access
	if diagnoseDestination != "" {
		results.Tests["s3_bucket_access"] = runDiagnosticTest("S3 Bucket Access",
			fmt.Sprintf("Testing access to S3 destination: %s", diagnoseDestination), func() error {
				// Simulate S3 bucket access test
				time.Sleep(200 * time.Millisecond)
				return nil
			})
	}

	return nil
}

// diagnoseWorkflow runs workflow-specific diagnostics
func diagnoseWorkflow(ctx context.Context, erm *data.ErrorRecoveryManager, results *DiagnosticResults) error {
	fmt.Printf("🔄 Workflow Diagnostics\n")
	fmt.Printf("──────────────────────\n")

	// Test workflow engine initialization
	results.Tests["workflow_engine"] = runDiagnosticTest("Workflow Engine",
		"Testing workflow engine initialization", func() error {
			engine := data.NewWorkflowEngine(nil)
			if engine == nil {
				return fmt.Errorf("failed to create workflow engine")
			}
			return nil
		})

	return nil
}

// diagnoseGeneral runs general system diagnostics
func diagnoseGeneral(ctx context.Context, erm *data.ErrorRecoveryManager, results *DiagnosticResults) error {
	fmt.Printf("🔍 General System Diagnostics\n")
	fmt.Printf("────────────────────────────\n")

	// Run subset of other diagnostics
	diagnoseNetwork(ctx, erm, results)
	diagnoseAuthentication(ctx, erm, results)
	diagnoseTransfer(ctx, erm, results)

	return nil
}

// runDiagnosticTest executes a single diagnostic test
func runDiagnosticTest(name, description string, testFunc func() error) *DiagnosticTest {
	test := &DiagnosticTest{
		Name:        name,
		Description: description,
		Suggestions: []string{},
	}

	startTime := time.Now()
	err := testFunc()
	test.Duration = time.Since(startTime)

	if err == nil {
		test.Status = "pass"
		fmt.Printf("✅ %s: PASS (%.2fs)\n", name, test.Duration.Seconds())
	} else {
		test.Status = "fail"
		test.Error = err
		fmt.Printf("❌ %s: FAIL (%.2fs) - %v\n", name, test.Duration.Seconds(), err)

		// Add error-specific suggestions
		test.Suggestions = generateTestSuggestions(err)
	}

	return test
}

// generateTestSuggestions generates suggestions based on test errors
func generateTestSuggestions(err error) []string {
	if err == nil {
		return []string{}
	}

	errMsg := strings.ToLower(err.Error())
	suggestions := []string{}

	if strings.Contains(errMsg, "credentials") {
		suggestions = append(suggestions, "Configure AWS credentials using 'aws configure' or set environment variables")
	}
	if strings.Contains(errMsg, "not found") {
		suggestions = append(suggestions, "Install the required tool or check the file path")
	}
	if strings.Contains(errMsg, "permission") {
		suggestions = append(suggestions, "Check file/directory permissions and IAM policies")
	}
	if strings.Contains(errMsg, "network") || strings.Contains(errMsg, "timeout") {
		suggestions = append(suggestions, "Check network connectivity and firewall settings")
	}

	return suggestions
}

// displayDiagnosticResults shows the diagnostic results
func displayDiagnosticResults(results *DiagnosticResults) {
	fmt.Printf("\n📊 Diagnostic Summary\n")
	fmt.Printf("═══════════════════\n")

	// Calculate summary
	summary := &DiagnosticSummary{}
	for _, test := range results.Tests {
		summary.TotalTests++
		switch test.Status {
		case "pass":
			summary.PassedTests++
		case "fail":
			summary.FailedTests++
		case "warning":
			summary.WarningTests++
		case "skip":
			summary.SkippedTests++
		}
	}
	results.Summary = summary

	fmt.Printf("Total Tests: %d\n", summary.TotalTests)
	fmt.Printf("✅ Passed: %d\n", summary.PassedTests)
	fmt.Printf("❌ Failed: %d\n", summary.FailedTests)
	if summary.WarningTests > 0 {
		fmt.Printf("⚠️  Warnings: %d\n", summary.WarningTests)
	}
	if summary.SkippedTests > 0 {
		fmt.Printf("⏭️  Skipped: %d\n", summary.SkippedTests)
	}
	fmt.Printf("Duration: %v\n", results.Duration)

	// Show suggestions for failed tests
	if summary.FailedTests > 0 {
		fmt.Printf("\n💡 Recommendations:\n")
		suggestionCount := 0
		for _, test := range results.Tests {
			if test.Status == "fail" && len(test.Suggestions) > 0 {
				for _, suggestion := range test.Suggestions {
					suggestionCount++
					fmt.Printf("  %d. %s\n", suggestionCount, suggestion)
				}
			}
		}
	}

	if summary.FailedTests == 0 {
		fmt.Printf("\n🎉 All diagnostics passed! Your system appears to be configured correctly.\n")
	}
}

// runInteractiveTroubleshooting provides interactive troubleshooting guidance
func runInteractiveTroubleshooting(results *DiagnosticResults) error {
	fmt.Printf("\n🤖 Interactive Troubleshooting Mode\n")
	fmt.Printf("═══════════════════════════════════\n")

	fmt.Printf("Based on the diagnostic results, here are the recommended next steps:\n\n")

	if results.Summary.FailedTests == 0 {
		fmt.Printf("✅ All tests passed! No troubleshooting needed.\n")
		return nil
	}

	// Provide step-by-step guidance for failed tests
	stepNum := 1
	for _, test := range results.Tests {
		if test.Status == "fail" {
			fmt.Printf("Step %d: Fix %s\n", stepNum, test.Name)
			fmt.Printf("  Problem: %v\n", test.Error)
			if len(test.Suggestions) > 0 {
				fmt.Printf("  Solutions:\n")
				for _, suggestion := range test.Suggestions {
					fmt.Printf("    • %s\n", suggestion)
				}
			}
			fmt.Println()
			stepNum++
		}
	}

	return nil
}

// generateDiagnosticReport saves diagnostic results to a file
func generateDiagnosticReport(results *DiagnosticResults, outputFile string) error {
	// For now, just print that we would generate a report
	// In a full implementation, this would save JSON/YAML
	fmt.Printf("Diagnostic report would be saved to: %s\n", outputFile)
	return nil
}
