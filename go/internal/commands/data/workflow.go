package data

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws-research-wizard/go/internal/data"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// workflowCmd represents the workflow command
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Execute and manage intelligent data movement workflows",
	Long: `Execute and manage intelligent data movement workflows that combine analysis,
bundling, transfer, and monitoring into automated pipelines.

Examples:
  # Execute a workflow from a project configuration
  aws-research-wizard data workflow run --config project.yaml --workflow upload_genomics_data

  # List active workflows
  aws-research-wizard data workflow list

  # Get status of a specific workflow
  aws-research-wizard data workflow status --id wf_1234567890

  # Cancel a running workflow
  aws-research-wizard data workflow cancel --id wf_1234567890`,
}

var runWorkflowCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a data movement workflow",
	Long: `Execute a data movement workflow from a project configuration file.
The workflow will analyze the data, apply optimizations like bundling,
transfer data using the optimal engine, and provide detailed reporting.`,
	RunE: runWorkflow,
}

var listWorkflowsCmd = &cobra.Command{
	Use:   "list",
	Short: "List active workflows",
	Long:  `List all currently active workflow executions with their status and progress.`,
	RunE:  listWorkflows,
}

var statusWorkflowCmd = &cobra.Command{
	Use:   "status",
	Short: "Get workflow execution status",
	Long:  `Get detailed status information for a specific workflow execution.`,
	RunE:  statusWorkflow,
}

var cancelWorkflowCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a running workflow",
	Long:  `Cancel a currently running workflow execution.`,
	RunE:  cancelWorkflow,
}

// Command flags
var (
	configFile           string
	workflowName         string
	executionID          string
	workflowOutputFormat string
	followProgress       bool
	dryRun               bool
)

func init() {
	// Add subcommands
	workflowCmd.AddCommand(runWorkflowCmd)
	workflowCmd.AddCommand(listWorkflowsCmd)
	workflowCmd.AddCommand(statusWorkflowCmd)
	workflowCmd.AddCommand(cancelWorkflowCmd)

	// Run command flags
	runWorkflowCmd.Flags().StringVarP(&configFile, "config", "c", "", "Project configuration file (required)")
	runWorkflowCmd.Flags().StringVarP(&workflowName, "workflow", "w", "", "Workflow name to execute (required)")
	runWorkflowCmd.Flags().BoolVarP(&followProgress, "follow", "f", false, "Follow workflow progress in real-time")
	runWorkflowCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be executed without running")
	runWorkflowCmd.MarkFlagRequired("config")
	runWorkflowCmd.MarkFlagRequired("workflow")

	// Status command flags
	statusWorkflowCmd.Flags().StringVar(&executionID, "id", "", "Workflow execution ID (required)")
	statusWorkflowCmd.MarkFlagRequired("id")

	// Cancel command flags
	cancelWorkflowCmd.Flags().StringVar(&executionID, "id", "", "Workflow execution ID (required)")
	cancelWorkflowCmd.MarkFlagRequired("id")

	// Global flags
	workflowCmd.PersistentFlags().StringVarP(&workflowOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load project configuration
	projectConfig, err := loadProjectConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load project configuration: %w", err)
	}

	// Validate workflow exists
	var workflow *data.Workflow
	for _, w := range projectConfig.Workflows {
		if w.Name == workflowName {
			workflow = &w
			break
		}
	}

	if workflow == nil {
		return fmt.Errorf("workflow '%s' not found in configuration", workflowName)
	}

	if dryRun {
		return showWorkflowPlan(projectConfig, workflow)
	}

	// Create and configure workflow engine
	engine := createWorkflowEngine()

	// Execute workflow
	fmt.Printf("🚀 Starting workflow '%s'...\n", workflowName)
	execution, err := engine.ExecuteWorkflow(ctx, projectConfig, workflowName)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}

	fmt.Printf("✅ Workflow started with ID: %s\n", execution.ID)

	if followProgress {
		return followWorkflowProgress(engine, execution.ID)
	}

	return nil
}

func listWorkflows(cmd *cobra.Command, args []string) error {
	engine := createWorkflowEngine()
	activeWorkflows := engine.GetActiveWorkflows()

	if len(activeWorkflows) == 0 {
		fmt.Println("No active workflows")
		return nil
	}

	switch workflowOutputFormat {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(activeWorkflows)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(activeWorkflows)
	default:
		printWorkflowTable(activeWorkflows)
		return nil
	}
}

func statusWorkflow(cmd *cobra.Command, args []string) error {
	engine := createWorkflowEngine()
	execution, err := engine.GetWorkflowExecution(executionID)
	if err != nil {
		return fmt.Errorf("failed to get workflow status: %w", err)
	}

	switch workflowOutputFormat {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(execution)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(execution)
	default:
		printWorkflowStatus(execution)
		return nil
	}
}

func cancelWorkflow(cmd *cobra.Command, args []string) error {
	engine := createWorkflowEngine()

	err := engine.CancelWorkflow(executionID)
	if err != nil {
		return fmt.Errorf("failed to cancel workflow: %w", err)
	}

	fmt.Printf("✅ Workflow %s has been cancelled\n", executionID)
	return nil
}

// Helper functions

func loadProjectConfig(configPath string) (*data.ProjectConfig, error) {
	if !filepath.IsAbs(configPath) {
		wd, _ := os.Getwd()
		configPath = filepath.Join(wd, configPath)
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config data.ProjectConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func createWorkflowEngine() *data.WorkflowEngine {
	// Create workflow engine with default configuration
	config := &data.WorkflowEngineConfig{
		MaxConcurrentWorkflows: 3,
		DefaultTimeout:         2 * time.Hour,
		RetryAttempts:          3,
		RetryDelay:             30 * time.Second,
		MonitoringEnabled:      true,
	}

	engine := data.NewWorkflowEngine(config)

	// Register components
	engine.RegisterAnalyzer(data.NewPatternAnalyzer())
	engine.RegisterBundlingEngine(data.NewBundlingEngine(nil))
	engine.RegisterWarningSystem(data.NewWarningSystem())

	// Register transfer engines
	s5cmdEngine := data.NewS5cmdEngine("s5cmd") // Use default executable name
	engine.RegisterTransferEngine(s5cmdEngine)

	rcloneEngine := data.NewRcloneEngine("rclone", "") // Use default executable name and empty config path
	engine.RegisterTransferEngine(rcloneEngine)

	// Register recommendation engine
	costCalculator := data.NewS3CostCalculator("us-east-1")
	recommendationEngine := data.NewRecommendationEngine(
		data.NewPatternAnalyzer(),
		costCalculator,
		nil, // engine registry would be set up here
		nil, // config manager would be set up here
	)
	engine.RegisterRecommendationEngine(recommendationEngine)

	return engine
}

func showWorkflowPlan(projectConfig *data.ProjectConfig, workflow *data.Workflow) error {
	fmt.Printf("🔍 Workflow Dry-Run Validation: %s\n", workflow.Name)
	fmt.Printf("==========================================\n\n")

	// Basic workflow information
	fmt.Printf("📋 Workflow Configuration:\n")
	fmt.Printf("  Name: %s\n", workflow.Name)
	fmt.Printf("  Description: %s\n", workflow.Description)
	fmt.Printf("  Source: %s\n", workflow.Source)
	fmt.Printf("  Destination: %s\n", workflow.Destination)
	fmt.Printf("  Engine: %s\n", workflow.Engine)
	fmt.Printf("  Enabled: %t\n", workflow.Enabled)
	fmt.Println()

	// Validation results
	validationErrors := []string{}
	validationWarnings := []string{}

	// Validate source data profile
	sourceProfile, sourceExists := projectConfig.DataProfiles[workflow.Source]
	if !sourceExists {
		validationErrors = append(validationErrors, fmt.Sprintf("Source data profile '%s' not found", workflow.Source))
	}

	// Validate destination
	destination, destExists := projectConfig.Destinations[workflow.Destination]
	if !destExists {
		validationErrors = append(validationErrors, fmt.Sprintf("Destination '%s' not found", workflow.Destination))
	}

	// Validate engine
	if workflow.Engine == "" || workflow.Engine == "auto" {
		validationWarnings = append(validationWarnings, "Engine set to 'auto' - will be selected based on data pattern analysis")
	}

	// Validate workflow is enabled
	if !workflow.Enabled {
		validationWarnings = append(validationWarnings, "Workflow is disabled and will not execute")
	}

	// Validate source path exists (if available)
	if sourceExists {
		if sourceProfile.Path != "" {
			if _, err := os.Stat(sourceProfile.Path); os.IsNotExist(err) {
				validationErrors = append(validationErrors, fmt.Sprintf("Source path does not exist: %s", sourceProfile.Path))
			}
		}
	}

	// Display validation results
	if len(validationErrors) > 0 {
		fmt.Printf("❌ Validation Errors:\n")
		for _, err := range validationErrors {
			fmt.Printf("  • %s\n", err)
		}
		fmt.Println()
	}

	if len(validationWarnings) > 0 {
		fmt.Printf("⚠️  Validation Warnings:\n")
		for _, warning := range validationWarnings {
			fmt.Printf("  • %s\n", warning)
		}
		fmt.Println()
	}

	if len(validationErrors) == 0 && len(validationWarnings) == 0 {
		fmt.Printf("✅ Validation passed - no issues found\n\n")
	}

	// Only continue with execution plan if validation passes
	if len(validationErrors) > 0 {
		return fmt.Errorf("workflow validation failed with %d errors", len(validationErrors))
	}

	// Perform data analysis to get realistic estimates
	if sourceExists && sourceProfile.Path != "" {
		fmt.Printf("🔍 Data Analysis Preview:\n")
		analyzer := data.NewPatternAnalyzer()
		pattern, err := analyzer.AnalyzePattern(context.Background(), sourceProfile.Path)
		if err != nil {
			fmt.Printf("  ⚠️  Could not analyze data: %v\n", err)
		} else {
			fmt.Printf("  Files to transfer: %d\n", pattern.TotalFiles)
			fmt.Printf("  Total data size: %s\n", pattern.TotalSizeHuman)
			if pattern.FileSizes.SmallFiles.CountUnder1MB > 0 {
				fmt.Printf("  Small files (<1MB): %d (%.1f%%)\n",
					pattern.FileSizes.SmallFiles.CountUnder1MB,
					pattern.FileSizes.SmallFiles.PercentageSmall)
			}

			// Show domain detection
			if len(pattern.DomainHints.DetectedDomains) > 0 {
				fmt.Printf("  Detected domains: ")
				for i, domain := range pattern.DomainHints.DetectedDomains {
					if i > 0 {
						fmt.Printf(", ")
					}
					confidence := pattern.DomainHints.Confidence[domain] * 100
					fmt.Printf("%s (%.1f%%)", domain, confidence)
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}

	// Show estimated execution plan
	fmt.Printf("📋 Estimated Execution Plan:\n")
	stepNum := 1

	// Analysis step
	fmt.Printf("  %d. 🔍 Data Pattern Analysis\n", stepNum)
	fmt.Printf("     • Scan source directory structure\n")
	fmt.Printf("     • Analyze file types and sizes\n")
	fmt.Printf("     • Detect research domain patterns\n")
	fmt.Printf("     • Generate optimization recommendations\n")
	stepNum++

	// Preprocessing steps
	if len(workflow.PreProcessing) > 0 {
		for _, step := range workflow.PreProcessing {
			fmt.Printf("  %d. 🔧 %s (%s)\n", stepNum, step.Name, step.Type)
			switch step.Type {
			case "bundle":
				fmt.Printf("     • Identify small files for bundling\n")
				fmt.Printf("     • Create optimized bundles\n")
				fmt.Printf("     • Estimate cost savings\n")
			case "compress":
				fmt.Printf("     • Compress files for transfer\n")
				fmt.Printf("     • Estimate space savings\n")
			case "validate":
				fmt.Printf("     • Validate file integrity\n")
				fmt.Printf("     • Check permissions\n")
			default:
				fmt.Printf("     • Execute %s processing\n", step.Type)
			}
			stepNum++
		}
	}

	// Main transfer step
	fmt.Printf("  %d. 🚀 Primary Data Transfer\n", stepNum)
	fmt.Printf("     • Engine: %s\n", workflow.Engine)
	if destExists {
		fmt.Printf("     • Destination: %s\n", destination.URI)
		fmt.Printf("     • Storage class: %s\n", destination.StorageClass)
	}
	if workflow.Configuration.Concurrency > 0 {
		fmt.Printf("     • Concurrency: %d parallel transfers\n", workflow.Configuration.Concurrency)
	}
	if workflow.Configuration.PartSize != "" {
		fmt.Printf("     • Part size: %s\n", workflow.Configuration.PartSize)
	}
	if workflow.Configuration.Checksum {
		fmt.Printf("     • Checksums: enabled\n")
	}
	stepNum++

	// Postprocessing steps
	if len(workflow.PostProcessing) > 0 {
		for _, step := range workflow.PostProcessing {
			fmt.Printf("  %d. ✅ %s (%s)\n", stepNum, step.Name, step.Type)
			switch step.Type {
			case "verify":
				fmt.Printf("     • Verify transfer completeness\n")
				fmt.Printf("     • Validate checksums\n")
			case "cleanup":
				fmt.Printf("     • Clean up temporary files\n")
				fmt.Printf("     • Remove intermediate bundles\n")
			case "report":
				fmt.Printf("     • Generate transfer report\n")
				fmt.Printf("     • Calculate metrics\n")
			default:
				fmt.Printf("     • Execute %s processing\n", step.Type)
			}
			stepNum++
		}
	}

	// Final reporting
	fmt.Printf("  %d. 📊 Final Report Generation\n", stepNum)
	fmt.Printf("     • Cost analysis and savings summary\n")
	fmt.Printf("     • Transfer performance metrics\n")
	fmt.Printf("     • Recommendations for optimization\n")
	fmt.Printf("     • Workflow execution summary\n")

	fmt.Printf("\n🎯 Dry-run completed successfully!\n")
	fmt.Printf("💡 To execute this workflow: aws-research-wizard data workflow run --workflow=%s\n", workflow.Name)

	return nil
}

func followWorkflowProgress(engine *data.WorkflowEngine, executionID string) error {
	fmt.Printf("📊 Following progress for workflow %s...\n\n", executionID)

	// Register progress callback
	engine.RegisterProgressCallback(executionID, func(execution *data.WorkflowExecution) {
		printProgressUpdate(execution)
	})

	// Poll for completion
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		execution, err := engine.GetWorkflowExecution(executionID)
		if err != nil {
			return fmt.Errorf("failed to get workflow status: %w", err)
		}

		if execution.Status == data.WorkflowStatusCompleted ||
			execution.Status == data.WorkflowStatusFailed ||
			execution.Status == data.WorkflowStatusCancelled {
			fmt.Printf("\n🏁 Workflow %s: %s\n", execution.ID, execution.Status)

			if execution.Results != nil {
				printFinalResults(execution.Results)
			}

			break
		}

		<-ticker.C
	}

	return nil
}

func printProgressUpdate(execution *data.WorkflowExecution) {
	progress := execution.Progress * 100
	currentStep := "unknown"

	if execution.CurrentStep < len(execution.Steps) {
		currentStep = execution.Steps[execution.CurrentStep].Name
	}

	fmt.Printf("\r🔄 Progress: %.1f%% | Step: %s", progress, currentStep)
}

func printWorkflowTable(workflows map[string]*data.WorkflowExecution) {
	fmt.Printf("%-20s %-25s %-12s %-10s %-15s\n", "ID", "Workflow", "Status", "Progress", "Duration")
	fmt.Println(string(make([]byte, 85, 85)))

	for _, execution := range workflows {
		duration := time.Since(execution.StartTime).Truncate(time.Second)
		progress := fmt.Sprintf("%.1f%%", execution.Progress*100)

		fmt.Printf("%-20s %-25s %-12s %-10s %-15s\n",
			execution.ID[:20],
			execution.WorkflowName,
			execution.Status,
			progress,
			duration.String())
	}
}

func printWorkflowStatus(execution *data.WorkflowExecution) {
	fmt.Printf("🔧 Workflow Status: %s\n", execution.ID)
	fmt.Printf("Name: %s\n", execution.WorkflowName)
	fmt.Printf("Status: %s\n", execution.Status)
	fmt.Printf("Progress: %.1f%%\n", execution.Progress*100)
	fmt.Printf("Duration: %s\n", time.Since(execution.StartTime).Truncate(time.Second))
	fmt.Printf("Steps: %d/%d\n", execution.CurrentStep+1, execution.TotalSteps)
	fmt.Println()

	// Show step details
	fmt.Println("📋 Steps:")
	for i, step := range execution.Steps {
		status := "⏸️"
		switch step.Status {
		case data.StepStatusRunning:
			status = "🔄"
		case data.StepStatusCompleted:
			status = "✅"
		case data.StepStatusFailed:
			status = "❌"
		case data.StepStatusSkipped:
			status = "⏭️"
		}

		current := ""
		if i == execution.CurrentStep {
			current = " (current)"
		}

		fmt.Printf("  %s %s%s\n", status, step.Name, current)

		if step.Duration > 0 {
			fmt.Printf("     Duration: %s\n", step.Duration.Truncate(time.Second))
		}

		if step.Error != nil {
			fmt.Printf("     Error: %v\n", step.Error)
		}
	}

	// Show recent events
	if len(execution.Events) > 0 {
		fmt.Println("\n📝 Recent Events:")
		eventCount := len(execution.Events)
		start := eventCount - 5
		if start < 0 {
			start = 0
		}

		for i := start; i < eventCount; i++ {
			event := execution.Events[i]
			timestamp := event.Timestamp.Format("15:04:05")
			fmt.Printf("  [%s] %s: %s\n", timestamp, event.Type, event.Message)
		}
	}
}

func printFinalResults(results *data.WorkflowResults) {
	fmt.Println("\n📊 Final Results:")

	if results.TotalFilesProcessed > 0 {
		fmt.Printf("Files Processed: %d\n", results.TotalFilesProcessed)
	}

	if results.TotalBytesTransferred > 0 {
		fmt.Printf("Bytes Transferred: %s\n", formatBytes(results.TotalBytesTransferred))
	}

	if results.TotalCostSavings > 0 {
		fmt.Printf("💰 Cost Savings: $%.2f/month\n", results.TotalCostSavings)
	}

	if results.SuccessRate > 0 {
		fmt.Printf("Success Rate: %.1f%%\n", results.SuccessRate)
	}

	if len(results.NextStepSuggestions) > 0 {
		fmt.Println("\n💡 Recommendations for Next Steps:")
		for _, suggestion := range results.NextStepSuggestions {
			fmt.Printf("  • %s\n", suggestion)
		}
	}
}

// formatBytes function is imported from the data package

// GetWorkflowCmd returns the workflow command for registration
func GetWorkflowCmd() *cobra.Command {
	return workflowCmd
}
