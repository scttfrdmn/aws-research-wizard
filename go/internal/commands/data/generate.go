package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws-research-wizard/go/internal/data"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command for creating project configurations
var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate optimized project configuration from data analysis",
	Long: `Generate an optimized project configuration by analyzing your data and 
applying research domain best practices. This command creates a complete 
project.yaml file that can be used immediately with the workflow engine.

The command performs:
- Comprehensive data pattern analysis
- Research domain detection and optimization
- Cost optimization recommendations
- Transfer engine selection
- Workflow template generation

Examples:
  # Generate config for current directory
  aws-research-wizard data generate .

  # Generate config with specific output file
  aws-research-wizard data generate /data/genomics --output genomics-project.yaml

  # Generate config with domain hint and custom project settings
  aws-research-wizard data generate /data --domain genomics --project-name "My Research" --owner "researcher@example.com"

  # Generate minimal config for quick testing
  aws-research-wizard data generate . --template minimal`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGenerate,
}

var (
	generateOutputFile   string
	generateProjectName  string
	generateOwner        string
	generateBudget       string
	generateDomain       string
	generateTemplate     string
	generateDestinations []string
	generateVerbose      bool
	generateOverwrite    bool
)

func init() {
	// Add generate command to data command
	DataCmd.AddCommand(generateCmd)
	
	// Output and project settings
	generateCmd.Flags().StringVarP(&generateOutputFile, "output", "o", "project.yaml", "Output file for generated configuration")
	generateCmd.Flags().StringVar(&generateProjectName, "project-name", "", "Project name (auto-detected if not provided)")
	generateCmd.Flags().StringVar(&generateOwner, "owner", "", "Project owner email")
	generateCmd.Flags().StringVar(&generateBudget, "budget", "1000", "Monthly budget in USD")
	
	// Analysis hints
	generateCmd.Flags().StringVar(&generateDomain, "domain", "", "Research domain hint (genomics, climate, ml, etc.)")
	generateCmd.Flags().StringVar(&generateTemplate, "template", "optimized", "Configuration template (minimal, optimized, comprehensive)")
	generateCmd.Flags().StringSliceVar(&generateDestinations, "destination", []string{}, "S3 destinations (format: name=s3://bucket/prefix)")
	
	// Behavior
	generateCmd.Flags().BoolVarP(&generateVerbose, "verbose", "v", false, "Show detailed generation process")
	generateCmd.Flags().BoolVar(&generateOverwrite, "overwrite", false, "Overwrite existing configuration file")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Determine path to analyze
	analyzePath := "."
	if len(args) > 0 {
		analyzePath = args[0]
	}
	
	// Convert to absolute path
	absPath, err := filepath.Abs(analyzePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	
	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", absPath)
	}
	
	// Check output file
	outputPath, err := filepath.Abs(generateOutputFile)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}
	
	if !generateOverwrite {
		if _, err := os.Stat(outputPath); err == nil {
			return fmt.Errorf("output file already exists: %s (use --overwrite to replace)", outputPath)
		}
	}
	
	fmt.Printf("🔧 AWS Research Wizard - Configuration Generator\n")
	fmt.Printf("================================================\n\n")
	fmt.Printf("📂 Analyzing data in: %s\n", absPath)
	fmt.Printf("📝 Output file: %s\n", outputPath)
	fmt.Printf("🎯 Template: %s\n", generateTemplate)
	if generateDomain != "" {
		fmt.Printf("🔬 Domain hint: %s\n", generateDomain)
	}
	fmt.Println()
	
	// Step 1: Analyze data patterns
	if generateVerbose {
		fmt.Println("🔍 Step 1: Analyzing data patterns...")
	}
	
	analyzer := data.NewPatternAnalyzer()
	ctx := context.Background()
	
	pattern, err := analyzer.AnalyzePattern(ctx, absPath)
	if err != nil {
		return fmt.Errorf("pattern analysis failed: %w", err)
	}
	
	// Apply domain hint if provided
	if generateDomain != "" {
		found := false
		for _, domain := range pattern.DomainHints.DetectedDomains {
			if domain == generateDomain {
				found = true
				break
			}
		}
		if !found {
			pattern.DomainHints.DetectedDomains = append(pattern.DomainHints.DetectedDomains, generateDomain)
			pattern.DomainHints.Confidence[generateDomain] = 0.9 // High confidence for user hint
		}
	}
	
	if generateVerbose {
		fmt.Printf("   • Total files: %d\n", pattern.TotalFiles)
		fmt.Printf("   • Total size: %s\n", pattern.TotalSizeHuman)
		fmt.Printf("   • Detected domains: %v\n", pattern.DomainHints.DetectedDomains)
		fmt.Printf("   • Small files under 1MB: %d\n", pattern.FileSizes.SmallFiles.CountUnder1MB)
	}
	
	// Step 2: Generate recommendations
	if generateVerbose {
		fmt.Println("🚀 Step 2: Generating optimization recommendations...")
	}
	
	costCalculator := data.NewS3CostCalculator("us-east-1")
	recommendationEngine := data.NewRecommendationEngine(analyzer, costCalculator, nil, nil)
	
	recommendations, err := recommendationEngine.GenerateRecommendations(ctx, absPath)
	if err != nil {
		if generateVerbose {
			fmt.Printf("   ⚠️  Warning: Could not generate recommendations: %v\n", err)
		}
		// Continue without recommendations
	}
	
	if generateVerbose && recommendations != nil {
		fmt.Printf("   • Tool recommendations: %d\n", len(recommendations.ToolRecommendations))
		fmt.Printf("   • Optimization suggestions: %d\n", len(recommendations.OptimizationSuggestions))
		if recommendations.CostAnalysis != nil {
			fmt.Printf("   • Potential savings: $%.2f/month\n", recommendations.CostAnalysis.PotentialSavings)
		}
	}
	
	// Step 3: Generate project configuration
	if generateVerbose {
		fmt.Println("📋 Step 3: Generating project configuration...")
	}
	
	pcm := data.NewProjectConfigManager(filepath.Dir(outputPath))
	
	// Generate base configuration
	projectConfig, err := pcm.GenerateConfig(pattern, recommendations)
	if err != nil {
		return fmt.Errorf("config generation failed: %w", err)
	}
	
	// Apply user customizations
	err = applyUserCustomizations(projectConfig, absPath)
	if err != nil {
		return fmt.Errorf("failed to apply customizations: %w", err)
	}
	
	// Apply template modifications
	err = applyTemplate(projectConfig, generateTemplate)
	if err != nil {
		return fmt.Errorf("failed to apply template: %w", err)
	}
	
	if generateVerbose {
		fmt.Printf("   • Project name: %s\n", projectConfig.Project.Name)
		fmt.Printf("   • Data profiles: %d\n", len(projectConfig.DataProfiles))
		fmt.Printf("   • Destinations: %d\n", len(projectConfig.Destinations))
		fmt.Printf("   • Workflows: %d\n", len(projectConfig.Workflows))
	}
	
	// Step 4: Save configuration
	if generateVerbose {
		fmt.Println("💾 Step 4: Saving configuration...")
	}
	
	err = pcm.SaveConfig(projectConfig, outputPath)
	if err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	
	// Success output
	fmt.Printf("✅ Configuration generated successfully!\n\n")
	
	// Show summary
	fmt.Println("📊 Configuration Summary")
	fmt.Println("========================")
	fmt.Printf("Project: %s (%s)\n", projectConfig.Project.Name, projectConfig.Project.Domain)
	fmt.Printf("Data profiles: %d\n", len(projectConfig.DataProfiles))
	fmt.Printf("Destinations: %d\n", len(projectConfig.Destinations))
	fmt.Printf("Workflows: %d\n", len(projectConfig.Workflows))
	
	// Show key optimizations
	if len(pattern.DomainHints.DetectedDomains) > 0 {
		fmt.Printf("\n🎯 Applied Optimizations:\n")
		for _, domain := range pattern.DomainHints.DetectedDomains {
			confidence := pattern.DomainHints.Confidence[domain] * 100
			fmt.Printf("• %s domain optimizations (%.1f%% confidence)\n", domain, confidence)
		}
	}
	
	if projectConfig.Optimization.CostOptimization.Enabled {
		fmt.Printf("• Cost optimization enabled\n")
		if projectConfig.Optimization.CostOptimization.AutoBundleSmallFiles {
			fmt.Printf("• Small file bundling enabled\n")
		}
	}
	
	// Show next steps
	fmt.Printf("\n🚀 Next Steps:\n")
	fmt.Printf("1. Review and customize: %s\n", outputPath)
	fmt.Printf("2. Test with dry-run: aws-research-wizard data workflow run --dry-run\n")
	fmt.Printf("3. Execute workflows: aws-research-wizard data workflow run\n")
	fmt.Printf("4. Monitor progress: aws-research-wizard data workflow status\n")
	
	return nil
}

// applyUserCustomizations applies user-provided settings to the project config
func applyUserCustomizations(config *data.ProjectConfig, dataPath string) error {
	// Project name
	if generateProjectName != "" {
		config.Project.Name = generateProjectName
	} else if config.Project.Name == "" {
		// Generate name from directory
		config.Project.Name = filepath.Base(dataPath) + "-project"
	}
	
	// Owner
	if generateOwner != "" {
		config.Project.Owner = generateOwner
	}
	
	// Budget
	if generateBudget != "" {
		config.Project.Budget = generateBudget
	}
	
	// Add custom destinations
	for _, dest := range generateDestinations {
		err := addCustomDestination(config, dest)
		if err != nil {
			return fmt.Errorf("invalid destination format '%s': %w", dest, err)
		}
	}
	
	return nil
}

// addCustomDestination parses and adds a custom destination
func addCustomDestination(config *data.ProjectConfig, destSpec string) error {
	// Parse format: name=s3://bucket/prefix
	parts := splitOnFirst(destSpec, "=")
	if len(parts) != 2 {
		return fmt.Errorf("expected format 'name=s3://bucket/prefix'")
	}
	
	name, uri := parts[0], parts[1]
	if name == "" || uri == "" {
		return fmt.Errorf("name and URI cannot be empty")
	}
	
	// Add to destinations
	config.Destinations[name] = data.Destination{
		Name: name,
		URI:  uri,
	}
	
	return nil
}

// splitOnFirst splits a string on the first occurrence of a delimiter
func splitOnFirst(s, delimiter string) []string {
	idx := len(s)
	for i := 0; i < len(s)-len(delimiter)+1; i++ {
		if s[i:i+len(delimiter)] == delimiter {
			idx = i
			break
		}
	}
	
	if idx == len(s) {
		return []string{s}
	}
	
	return []string{s[:idx], s[idx+len(delimiter):]}
}

// applyTemplate applies template-specific modifications
func applyTemplate(config *data.ProjectConfig, template string) error {
	switch template {
	case "minimal":
		return applyMinimalTemplate(config)
	case "optimized":
		return applyOptimizedTemplate(config)
	case "comprehensive":
		return applyComprehensiveTemplate(config)
	default:
		return fmt.Errorf("unknown template: %s (supported: minimal, optimized, comprehensive)", template)
	}
}

// applyMinimalTemplate creates a minimal configuration for quick testing
func applyMinimalTemplate(config *data.ProjectConfig) error {
	// Disable most optimizations for simplicity
	config.Optimization.CostOptimization.Enabled = false
	config.Optimization.PerformanceOptimization.Enabled = false
	
	// Use simple engines
	for i := range config.Workflows {
		workflow := &config.Workflows[i]
		if workflow.Engine == "auto" {
			workflow.Engine = "s5cmd" // Simple default
		}
		// Clear complex preprocessing/postprocessing
		workflow.PreProcessing = []data.ProcessingStep{}
		workflow.PostProcessing = []data.ProcessingStep{}
		
		// Simple configuration
		workflow.Configuration = data.WorkflowConfiguration{
			Concurrency: 4,
			PartSize:    "64MB",
		}
	}
	
	return nil
}

// applyOptimizedTemplate applies balanced optimizations (default)
func applyOptimizedTemplate(config *data.ProjectConfig) error {
	// This is the default behavior from the generator
	// Enable key optimizations without complexity
	config.Optimization.CostOptimization.Enabled = true
	config.Optimization.PerformanceOptimization.Enabled = true
	
	return nil
}

// applyComprehensiveTemplate enables all available optimizations
func applyComprehensiveTemplate(config *data.ProjectConfig) error {
	// Enable all optimizations
	config.Optimization.CostOptimization.Enabled = true
	config.Optimization.CostOptimization.AutoBundleSmallFiles = true
	config.Optimization.CostOptimization.AutoCompression = true
	config.Optimization.CostOptimization.AutoStorageClass = true
	
	config.Optimization.PerformanceOptimization.Enabled = true
	config.Optimization.PerformanceOptimization.AutoConcurrency = true
	config.Optimization.PerformanceOptimization.NetworkOptimization = true
	
	// Add comprehensive monitoring
	for i := range config.Workflows {
		workflow := &config.Workflows[i]
		
		// Add comprehensive validation
		workflow.PreProcessing = append(workflow.PreProcessing, data.ProcessingStep{
			Name: "comprehensive_validation",
			Type: "validation",
			Parameters: map[string]string{
				"check_integrity": "true",
				"check_permissions": "true",
				"estimate_costs": "true",
			},
		})
		
		// Add comprehensive verification
		workflow.PostProcessing = append(workflow.PostProcessing, data.ProcessingStep{
			Name: "comprehensive_verification",
			Type: "verification",
			Parameters: map[string]string{
				"verify_transfer": "true",
				"check_checksums": "true",
				"generate_report": "true",
			},
		})
	}
	
	return nil
}