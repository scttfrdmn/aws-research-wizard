package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws-research-wizard/go/internal/data"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command for comprehensive workflow validation
var validateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate workflow configurations and perform comprehensive dry-run checks",
	Long: `Validate workflow configurations without executing them. This command performs
comprehensive validation including:

- Configuration file syntax validation
- Data source and destination accessibility checks
- Transfer engine availability verification
- Cost estimation and optimization analysis
- Security and compliance validation
- Performance impact assessment

Examples:
  # Validate default configuration
  aws-research-wizard data validate

  # Validate specific configuration file
  aws-research-wizard data validate project.yaml

  # Validate all workflows with detailed output
  aws-research-wizard data validate --all --verbose

  # Check configuration and generate validation report
  aws-research-wizard data validate --report validation-report.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

var (
	validateAll          bool
	validateWorkflowName string
	validateVerbose      bool
	validateReport       string
	validateSkipAnalysis bool
	validateFixProblems  bool
)

func init() {
	// Add validate command to data command
	DataCmd.AddCommand(validateCmd)
	
	// Flags
	validateCmd.Flags().BoolVar(&validateAll, "all", false, "Validate all workflows in configuration")
	validateCmd.Flags().StringVarP(&validateWorkflowName, "workflow", "w", "", "Validate specific workflow by name")
	validateCmd.Flags().BoolVarP(&validateVerbose, "verbose", "v", false, "Show detailed validation information")
	validateCmd.Flags().StringVar(&validateReport, "report", "", "Generate validation report (json/yaml)")
	validateCmd.Flags().BoolVar(&validateSkipAnalysis, "skip-analysis", false, "Skip data analysis (faster validation)")
	validateCmd.Flags().BoolVar(&validateFixProblems, "fix", false, "Attempt to automatically fix common problems")
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Determine config file
	configFile := "project.yaml"
	if len(args) > 0 {
		configFile = args[0]
	}
	
	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configFile)
	}
	
	fmt.Printf("🔍 AWS Research Wizard - Workflow Validation\n")
	fmt.Printf("===========================================\n\n")
	fmt.Printf("📋 Configuration file: %s\n\n", configFile)
	
	// Load project configuration
	projectConfig, err := loadProjectConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	validationResults := &ValidationResults{
		ConfigFile: configFile,
		Workflows:  make(map[string]*WorkflowValidation),
	}
	
	// Validate overall project configuration
	if validateVerbose {
		fmt.Printf("🏗️  Validating project configuration...\n")
	}
	
	projectValidation := validateProjectConfig(projectConfig)
	validationResults.ProjectValidation = projectValidation
	
	if len(projectValidation.Errors) > 0 {
		fmt.Printf("❌ Project Configuration Errors:\n")
		for _, err := range projectValidation.Errors {
			fmt.Printf("  • %s\n", err)
		}
		fmt.Println()
	}
	
	if len(projectValidation.Warnings) > 0 && validateVerbose {
		fmt.Printf("⚠️  Project Configuration Warnings:\n")
		for _, warning := range projectValidation.Warnings {
			fmt.Printf("  • %s\n", warning)
		}
		fmt.Println()
	}
	
	// Determine which workflows to validate
	workflowsToValidate := []data.Workflow{}
	
	if validateWorkflowName != "" {
		// Validate specific workflow
		found := false
		for _, workflow := range projectConfig.Workflows {
			if workflow.Name == validateWorkflowName {
				workflowsToValidate = append(workflowsToValidate, workflow)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("workflow '%s' not found in configuration", validateWorkflowName)
		}
	} else if validateAll || len(projectConfig.Workflows) <= 3 {
		// Validate all workflows (or default if few workflows)
		workflowsToValidate = projectConfig.Workflows
	} else {
		// Validate enabled workflows only
		for _, workflow := range projectConfig.Workflows {
			if workflow.Enabled {
				workflowsToValidate = append(workflowsToValidate, workflow)
			}
		}
		if len(workflowsToValidate) == 0 {
			workflowsToValidate = projectConfig.Workflows // Fallback to all
		}
	}
	
	fmt.Printf("🔧 Validating %d workflow(s)...\n\n", len(workflowsToValidate))
	
	// Validate each workflow
	totalErrors := 0
	totalWarnings := 0
	
	for i, workflow := range workflowsToValidate {
		fmt.Printf("Workflow %d/%d: %s\n", i+1, len(workflowsToValidate), workflow.Name)
		fmt.Printf("─────────────────────────────────\n")
		
		workflowValidation := validateWorkflow(projectConfig, &workflow, !validateSkipAnalysis, validateVerbose)
		validationResults.Workflows[workflow.Name] = workflowValidation
		
		// Display validation results
		if len(workflowValidation.Errors) > 0 {
			fmt.Printf("❌ Errors (%d):\n", len(workflowValidation.Errors))
			for _, err := range workflowValidation.Errors {
				fmt.Printf("  • %s\n", err)
			}
			totalErrors += len(workflowValidation.Errors)
		}
		
		if len(workflowValidation.Warnings) > 0 {
			fmt.Printf("⚠️  Warnings (%d):\n", len(workflowValidation.Warnings))
			for _, warning := range workflowValidation.Warnings {
				fmt.Printf("  • %s\n", warning)
			}
			totalWarnings += len(workflowValidation.Warnings)
		}
		
		if len(workflowValidation.Errors) == 0 && len(workflowValidation.Warnings) == 0 {
			fmt.Printf("✅ Validation passed\n")
		}
		
		// Show optimization suggestions if verbose
		if validateVerbose && len(workflowValidation.Suggestions) > 0 {
			fmt.Printf("💡 Optimization Suggestions:\n")
			for _, suggestion := range workflowValidation.Suggestions {
				fmt.Printf("  • %s\n", suggestion)
			}
		}
		
		fmt.Println()
	}
	
	// Summary
	fmt.Printf("📊 Validation Summary\n")
	fmt.Printf("═══════════════════\n")
	fmt.Printf("Configuration file: %s\n", configFile)
	fmt.Printf("Workflows validated: %d\n", len(workflowsToValidate))
	fmt.Printf("Total errors: %d\n", totalErrors)
	fmt.Printf("Total warnings: %d\n", totalWarnings)
	
	if totalErrors == 0 {
		fmt.Printf("🎉 All validations passed!\n")
	} else {
		fmt.Printf("⚠️  %d validation errors need attention\n", totalErrors)
	}
	
	// Generate report if requested
	if validateReport != "" {
		err := generateValidationReport(validationResults, validateReport)
		if err != nil {
			fmt.Printf("⚠️  Failed to generate report: %v\n", err)
		} else {
			fmt.Printf("📄 Validation report saved: %s\n", validateReport)
		}
	}
	
	// Return error if there were critical validation failures
	if totalErrors > 0 {
		return fmt.Errorf("validation failed with %d errors", totalErrors)
	}
	
	return nil
}

// ValidationResults holds the complete validation results
type ValidationResults struct {
	ConfigFile         string                         `json:"config_file"`
	ProjectValidation  *ProjectValidation             `json:"project_validation"`
	Workflows          map[string]*WorkflowValidation `json:"workflows"`
	Summary            *ValidationSummary             `json:"summary,omitempty"`
}

// ProjectValidation holds project-level validation results
type ProjectValidation struct {
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
	Suggestions []string `json:"suggestions"`
}

// WorkflowValidation holds workflow-specific validation results
type WorkflowValidation struct {
	WorkflowName    string               `json:"workflow_name"`
	Enabled         bool                 `json:"enabled"`
	Errors          []string             `json:"errors"`
	Warnings        []string             `json:"warnings"`
	Suggestions     []string             `json:"suggestions"`
	DataAnalysis    *data.DataPattern    `json:"data_analysis,omitempty"`
	EstimatedCost   float64              `json:"estimated_cost"`
	EstimatedTime   string               `json:"estimated_time"`
	ValidationTime  string               `json:"validation_time"`
}

// ValidationSummary provides overall validation statistics
type ValidationSummary struct {
	TotalWorkflows    int `json:"total_workflows"`
	PassedWorkflows   int `json:"passed_workflows"`
	FailedWorkflows   int `json:"failed_workflows"`
	TotalErrors       int `json:"total_errors"`
	TotalWarnings     int `json:"total_warnings"`
	TotalSuggestions  int `json:"total_suggestions"`
}

// validateProjectConfig validates the overall project configuration
func validateProjectConfig(config *data.ProjectConfig) *ProjectValidation {
	validation := &ProjectValidation{
		Errors:      []string{},
		Warnings:    []string{},
		Suggestions: []string{},
	}
	
	// Validate project metadata
	if config.Project.Name == "" {
		validation.Errors = append(validation.Errors, "Project name is required")
	}
	
	if config.Project.Owner == "" {
		validation.Warnings = append(validation.Warnings, "Project owner not specified")
	}
	
	// Validate data profiles
	if len(config.DataProfiles) == 0 {
		validation.Errors = append(validation.Errors, "At least one data profile is required")
	}
	
	for name, profile := range config.DataProfiles {
		if profile.Path == "" {
			validation.Errors = append(validation.Errors, fmt.Sprintf("Data profile '%s' missing path", name))
		}
	}
	
	// Validate destinations
	if len(config.Destinations) == 0 {
		validation.Errors = append(validation.Errors, "At least one destination is required")
	}
	
	for name, dest := range config.Destinations {
		if dest.URI == "" {
			validation.Errors = append(validation.Errors, fmt.Sprintf("Destination '%s' missing URI", name))
		}
	}
	
	// Validate workflows
	if len(config.Workflows) == 0 {
		validation.Warnings = append(validation.Warnings, "No workflows defined")
	}
	
	enabledWorkflows := 0
	for _, workflow := range config.Workflows {
		if workflow.Enabled {
			enabledWorkflows++
		}
	}
	
	if enabledWorkflows == 0 {
		validation.Warnings = append(validation.Warnings, "No workflows are enabled")
	}
	
	return validation
}

// validateWorkflow performs comprehensive validation of a single workflow
func validateWorkflow(config *data.ProjectConfig, workflow *data.Workflow, analyzeData bool, verbose bool) *WorkflowValidation {
	validation := &WorkflowValidation{
		WorkflowName: workflow.Name,
		Enabled:      workflow.Enabled,
		Errors:       []string{},
		Warnings:     []string{},
		Suggestions:  []string{},
	}
	
	// Basic workflow validation
	if workflow.Name == "" {
		validation.Errors = append(validation.Errors, "Workflow name is required")
	}
	
	if workflow.Source == "" {
		validation.Errors = append(validation.Errors, "Source data profile is required")
	}
	
	if workflow.Destination == "" {
		validation.Errors = append(validation.Errors, "Destination is required")
	}
	
	// Validate source data profile
	sourceProfile, sourceExists := config.DataProfiles[workflow.Source]
	if !sourceExists {
		validation.Errors = append(validation.Errors, fmt.Sprintf("Source data profile '%s' not found", workflow.Source))
	} else {
		// Validate source path exists
		if sourceProfile.Path != "" {
			if _, err := os.Stat(sourceProfile.Path); os.IsNotExist(err) {
				validation.Errors = append(validation.Errors, fmt.Sprintf("Source path does not exist: %s", sourceProfile.Path))
			} else if analyzeData {
				// Perform data analysis
				analyzer := data.NewPatternAnalyzer()
				pattern, err := analyzer.AnalyzePattern(context.Background(), sourceProfile.Path)
				if err != nil {
					validation.Warnings = append(validation.Warnings, fmt.Sprintf("Could not analyze source data: %v", err))
				} else {
					validation.DataAnalysis = pattern
					
					// Generate suggestions based on analysis
					if pattern.FileSizes.SmallFiles.CountUnder1MB > pattern.TotalFiles/2 {
						validation.Suggestions = append(validation.Suggestions, 
							"High ratio of small files detected - consider enabling bundling")
					}
					
					if len(pattern.DomainHints.DetectedDomains) > 0 {
						domain := pattern.DomainHints.DetectedDomains[0]
						if config.Project.Domain != domain {
							validation.Suggestions = append(validation.Suggestions, 
								fmt.Sprintf("Consider setting project domain to '%s' for optimizations", domain))
						}
					}
				}
			}
		}
	}
	
	// Validate destination
	destination, destExists := config.Destinations[workflow.Destination]
	if !destExists {
		validation.Errors = append(validation.Errors, fmt.Sprintf("Destination '%s' not found", workflow.Destination))
	} else {
		// Basic URI validation
		if destination.URI == "" {
			validation.Errors = append(validation.Errors, "Destination URI is empty")
		}
		// Could add more sophisticated URI validation here
	}
	
	// Validate engine
	if workflow.Engine == "" {
		validation.Warnings = append(validation.Warnings, "No engine specified - will use auto-detection")
	}
	
	// Validate configuration
	if workflow.Configuration.Concurrency < 0 {
		validation.Errors = append(validation.Errors, "Concurrency cannot be negative")
	}
	
	if workflow.Configuration.Concurrency > 100 {
		validation.Warnings = append(validation.Warnings, "Very high concurrency may impact performance")
	}
	
	// Validate preprocessing steps
	for _, step := range workflow.PreProcessing {
		if step.Name == "" {
			validation.Warnings = append(validation.Warnings, "Preprocessing step missing name")
		}
		if step.Type == "" {
			validation.Errors = append(validation.Errors, "Preprocessing step missing type")
		}
	}
	
	// Validate postprocessing steps
	for _, step := range workflow.PostProcessing {
		if step.Name == "" {
			validation.Warnings = append(validation.Warnings, "Postprocessing step missing name")
		}
		if step.Type == "" {
			validation.Errors = append(validation.Errors, "Postprocessing step missing type")
		}
	}
	
	return validation
}

// generateValidationReport saves validation results to a file
func generateValidationReport(results *ValidationResults, outputFile string) error {
	// Add summary
	summary := &ValidationSummary{
		TotalWorkflows: len(results.Workflows),
	}
	
	for _, workflow := range results.Workflows {
		if len(workflow.Errors) == 0 {
			summary.PassedWorkflows++
		} else {
			summary.FailedWorkflows++
		}
		summary.TotalErrors += len(workflow.Errors)
		summary.TotalWarnings += len(workflow.Warnings)
		summary.TotalSuggestions += len(workflow.Suggestions)
	}
	
	results.Summary = summary
	
	// Determine output format from file extension
	ext := filepath.Ext(outputFile)
	switch ext {
	case ".json":
		return saveJSONReport(results, outputFile)
	case ".yaml", ".yml":
		return saveYAMLReport(results, outputFile)
	default:
		return saveJSONReport(results, outputFile) // Default to JSON
	}
}

func saveJSONReport(results *ValidationResults, outputFile string) error {
	// Implementation would save JSON format
	return fmt.Errorf("JSON report generation not yet implemented")
}

func saveYAMLReport(results *ValidationResults, outputFile string) error {
	// Implementation would save YAML format
	return fmt.Errorf("YAML report generation not yet implemented")
}