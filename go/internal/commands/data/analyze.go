package data

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws-research-wizard/go/internal/data"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze data patterns and generate intelligent recommendations",
	Long: `Analyze data patterns in a directory and generate intelligent recommendations
for transfer optimization, cost savings, and domain-specific configurations.

This command performs comprehensive analysis including:
- File pattern detection and classification
- Research domain identification  
- Cost optimization analysis
- Transfer engine recommendations
- Bundling optimization suggestions
- Security and compliance recommendations

Examples:
  # Analyze current directory
  aws-research-wizard data analyze .

  # Analyze specific path with detailed output
  aws-research-wizard data analyze /data/genomics --output json --verbose

  # Generate project configuration from analysis
  aws-research-wizard data analyze /data/genomics --generate-config project.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAnalyze,
}

var (
	outputFormat    string
	verbose         bool
	generateConfig  bool
	configOutput    string
	includeEstimates bool
	domainHint      string
)

func init() {
	// Add analyze command to data command
	DataCmd.AddCommand(analyzeCmd)
	
	// Flags
	analyzeCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	analyzeCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed analysis")
	analyzeCmd.Flags().BoolVar(&generateConfig, "generate-config", false, "Generate project configuration")
	analyzeCmd.Flags().StringVar(&configOutput, "config-output", "project.yaml", "Output file for generated config")
	analyzeCmd.Flags().BoolVar(&includeEstimates, "include-estimates", true, "Include cost estimates")
	analyzeCmd.Flags().StringVar(&domainHint, "domain", "", "Hint for research domain (genomics, climate, ml, etc.)")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
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
	
	fmt.Printf("🔍 Analyzing data patterns in: %s\n\n", absPath)
	
	// Create analyzer
	analyzer := data.NewPatternAnalyzer()
	
	// Analyze patterns
	ctx := context.Background()
	pattern, err := analyzer.AnalyzePattern(ctx, absPath)
	if err != nil {
		return fmt.Errorf("pattern analysis failed: %w", err)
	}
	
	// Apply domain hint if provided
	if domainHint != "" {
		// Add domain hint to detected domains if not already present
		found := false
		for _, domain := range pattern.DomainHints.DetectedDomains {
			if domain == domainHint {
				found = true
				break
			}
		}
		if !found {
			pattern.DomainHints.DetectedDomains = append(pattern.DomainHints.DetectedDomains, domainHint)
			pattern.DomainHints.Confidence[domainHint] = 0.8 // User-provided hint
		}
	}
	
	// Generate recommendations
	recommendations, err := generateRecommendations(ctx, pattern, absPath)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not generate recommendations: %v\n", err)
	}
	
	// Display results
	switch outputFormat {
	case "json":
		err = outputJSON(pattern, recommendations)
	case "yaml":
		err = outputYAML(pattern, recommendations)
	default:
		err = outputTable(pattern, recommendations)
	}
	
	if err != nil {
		return fmt.Errorf("output failed: %w", err)
	}
	
	// Generate configuration if requested
	if generateConfig {
		err = generateProjectConfig(pattern, recommendations, configOutput)
		if err != nil {
			return fmt.Errorf("config generation failed: %w", err)
		}
		fmt.Printf("\n✅ Project configuration generated: %s\n", configOutput)
	}
	
	return nil
}

func generateRecommendations(ctx context.Context, pattern *data.DataPattern, path string) (*data.RecommendationResult, error) {
	// Create cost calculator
	costCalculator := data.NewS3CostCalculator("us-east-1")
	
	// Create recommendation engine
	analyzer := data.NewPatternAnalyzer()
	recommendationEngine := data.NewRecommendationEngine(analyzer, costCalculator, nil, nil)
	
	// Generate recommendations
	return recommendationEngine.GenerateRecommendations(ctx, path)
}

func outputTable(pattern *data.DataPattern, recommendations *data.RecommendationResult) error {
	fmt.Println("📊 Data Pattern Analysis Results")
	fmt.Println("================================")
	
	// Basic statistics
	fmt.Printf("Total Files:     %d\n", pattern.TotalFiles)
	fmt.Printf("Total Size:      %s\n", pattern.TotalSizeHuman)
	fmt.Printf("Analysis Time:   %s\n", pattern.AnalysisTime.Format("2006-01-02 15:04:05"))
	
	// Domain detection
	if len(pattern.DomainHints.DetectedDomains) > 0 {
		fmt.Printf("\n🔬 Detected Research Domains:\n")
		for _, domain := range pattern.DomainHints.DetectedDomains {
			confidence := pattern.DomainHints.Confidence[domain] * 100
			fmt.Printf("  • %s (%.1f%% confidence)\n", domain, confidence)
		}
	}
	
	// File type breakdown
	fmt.Printf("\n📁 File Type Analysis:\n")
	for ext, info := range pattern.FileTypes {
		if info.Count > 0 {
			percentage := float64(info.TotalSize) / float64(pattern.TotalSize) * 100
			fmt.Printf("  • %-10s: %6d files, %8s (%.1f%%)\n", 
				ext, info.Count, formatBytes(info.TotalSize), percentage)
		}
	}
	
	// Small file analysis
	smallFiles := pattern.FileSizes.SmallFiles
	if smallFiles.CountUnder1MB > 0 {
		fmt.Printf("\n⚠️  Small File Analysis:\n")
		fmt.Printf("  Files < 1KB:   %d\n", smallFiles.CountUnder1KB)
		fmt.Printf("  Files < 10KB:  %d\n", smallFiles.CountUnder10KB)
		fmt.Printf("  Files < 100KB: %d\n", smallFiles.CountUnder100KB)
		fmt.Printf("  Files < 1MB:   %d\n", smallFiles.CountUnder1MB)
		fmt.Printf("  Small file ratio: %.1f%%\n", smallFiles.PercentageSmall)
		
		if smallFiles.PercentageSmall > 50 {
			fmt.Printf("  💡 High small file ratio detected - bundling recommended\n")
		}
	}
	
	// Efficiency metrics
	if pattern.Efficiency.EstimatedBundles > 0 {
		fmt.Printf("\n💰 Efficiency Analysis:\n")
		fmt.Printf("  Current S3 PUTs:    %d requests\n", pattern.Efficiency.EstimatedPutRequests)
		fmt.Printf("  Estimated bundles:  %d\n", pattern.Efficiency.EstimatedBundles)
		if pattern.Efficiency.BundlingCostSavings > 0 {
			fmt.Printf("  Bundling savings:   $%.2f/month\n", pattern.Efficiency.BundlingCostSavings)
		}
		if pattern.Efficiency.StorageClassSavings > 0 {
			fmt.Printf("  Storage savings:    $%.2f/month\n", pattern.Efficiency.StorageClassSavings)
		}
	}
	
	// Recommendations
	if recommendations != nil {
		fmt.Printf("\n🚀 Optimization Recommendations:\n")
		
		// Tool recommendations
		for _, toolRec := range recommendations.ToolRecommendations {
			fmt.Printf("  • %s: Use %s (%.1f%% confidence)\n", 
				toolRec.Task, toolRec.RecommendedTool, toolRec.Confidence*100)
			if toolRec.Reasoning != "" {
				fmt.Printf("    Reason: %s\n", toolRec.Reasoning)
			}
		}
		
		// Optimization suggestions
		for _, opt := range recommendations.OptimizationSuggestions {
			fmt.Printf("  • %s: %s\n", opt.Type, opt.Description)
			if opt.Impact.CostSavingsMonthly > 0 {
				fmt.Printf("    Potential savings: $%.2f/month\n", opt.Impact.CostSavingsMonthly)
			}
		}
		
		// Cost analysis summary
		if recommendations.CostAnalysis != nil {
			fmt.Printf("\n💵 Cost Analysis:\n")
			if len(recommendations.CostAnalysis.Scenarios) > 0 {
				current := recommendations.CostAnalysis.Scenarios[0]
				fmt.Printf("  Current monthly cost: $%.2f\n", current.MonthlyCosts.Total)
			}
			if recommendations.CostAnalysis.PotentialSavings > 0 {
				fmt.Printf("  Potential savings:    $%.2f/month\n", recommendations.CostAnalysis.PotentialSavings)
			}
		}
	}
	
	// Domain-specific recommendations
	if len(pattern.DomainHints.DetectedDomains) > 0 {
		fmt.Printf("\n🎯 Domain-Specific Recommendations:\n")
		for _, domain := range pattern.DomainHints.DetectedDomains {
			showDomainRecommendations(domain)
		}
	}
	
	return nil
}

func showDomainRecommendations(domain string) {
	dpm := data.NewResearchDomainProfileManager()
	profile, exists := dpm.GetProfile(domain)
	if !exists {
		return
	}
	
	fmt.Printf("  %s:\n", domain)
	fmt.Printf("    • Preferred engines: %v\n", profile.TransferOptimization.PreferredEngines)
	fmt.Printf("    • Optimal concurrency: %d\n", profile.TransferOptimization.OptimalConcurrency)
	fmt.Printf("    • Bundling enabled: %t\n", profile.BundlingStrategy.EnableBundling)
	if profile.SecurityRequirements.EncryptionRequired {
		fmt.Printf("    • ⚠️  Encryption required for this domain\n")
	}
}

func outputJSON(pattern *data.DataPattern, recommendations *data.RecommendationResult) error {
	output := map[string]interface{}{
		"pattern":         pattern,
		"recommendations": recommendations,
		"timestamp":       pattern.AnalysisTime,
	}
	
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

func outputYAML(pattern *data.DataPattern, recommendations *data.RecommendationResult) error {
	output := map[string]interface{}{
		"pattern":         pattern,
		"recommendations": recommendations,
		"timestamp":       pattern.AnalysisTime,
	}
	
	return yaml.NewEncoder(os.Stdout).Encode(output)
}

func generateProjectConfig(pattern *data.DataPattern, recommendations *data.RecommendationResult, outputFile string) error {
	// Create project config manager
	pcm := data.NewProjectConfigManager(filepath.Dir(outputFile))
	
	// Generate configuration
	projectConfig, err := pcm.GenerateConfig(pattern, recommendations)
	if err != nil {
		return err
	}
	
	// Save to file
	return pcm.SaveConfig(projectConfig, outputFile)
}