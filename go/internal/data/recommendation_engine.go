package data

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// RecommendationEngine provides intelligent optimization suggestions
type RecommendationEngine struct {
	patternAnalyzer  *PatternAnalyzer
	costCalculator   *S3CostCalculator
	engineRegistry   *EngineRegistry
	configManager    *ConfigManager
}

// RecommendationResult contains all recommendations for a dataset
type RecommendationResult struct {
	AnalysisID       string                    `json:"analysis_id"`
	DataPath         string                    `json:"data_path"`
	Timestamp        time.Time                 `json:"timestamp"`
	
	// Analysis results
	DataPattern      *DataPattern             `json:"data_pattern"`
	CostAnalysis     *CostAnalysis            `json:"cost_analysis"`
	
	// Recommendations
	ToolRecommendations    []ToolRecommendation    `json:"tool_recommendations"`
	OptimizationSuggestions []OptimizationSuggestion `json:"optimization_suggestions"`
	WarningAlerts          []WarningAlert          `json:"warning_alerts"`
	
	// Summary
	TopRecommendation      *OptimizationSuggestion `json:"top_recommendation"`
	EstimatedSavings       float64                 `json:"estimated_total_savings"`
	ImplementationPriority []string                `json:"implementation_priority"`
	
	// Configuration generation
	GeneratedConfig        *ProjectConfig          `json:"generated_config,omitempty"`
}

// ToolRecommendation suggests optimal tools for specific tasks
type ToolRecommendation struct {
	Task            string             `json:"task"`
	RecommendedTool string             `json:"recommended_tool"`
	AlternativeTools []string          `json:"alternative_tools"`
	Reasoning       string             `json:"reasoning"`
	Confidence      float64            `json:"confidence"`
	Configuration   map[string]interface{} `json:"configuration"`
	ExpectedPerformance PerformanceEstimate `json:"expected_performance"`
}

// OptimizationSuggestion represents a specific optimization opportunity
type OptimizationSuggestion struct {
	ID              string             `json:"id"`
	Type            string             `json:"type"` // "bundling", "compression", "storage_class", "tool_chain"
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Impact          ImpactAssessment   `json:"impact"`
	Implementation  Implementation     `json:"implementation"`
	Prerequisites   []string           `json:"prerequisites"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// WarningAlert represents potential problems or anti-patterns
type WarningAlert struct {
	Severity    string    `json:"severity"` // "critical", "warning", "info"
	Category    string    `json:"category"` // "cost", "performance", "reliability"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"`
	Solution    string    `json:"solution"`
	LearnMoreURL string   `json:"learn_more_url,omitempty"`
}

// ImpactAssessment quantifies the impact of an optimization
type ImpactAssessment struct {
	CostSavingsMonthly    float64   `json:"cost_savings_monthly"`
	CostSavingsPercent    float64   `json:"cost_savings_percent"`
	PerformanceImprovement float64  `json:"performance_improvement_percent"`
	TimeToTransfer        string    `json:"time_to_transfer"`
	RiskLevel             string    `json:"risk_level"` // "low", "medium", "high"
	Confidence            float64   `json:"confidence"`
}

// Implementation describes how to implement a recommendation
type Implementation struct {
	Complexity      string           `json:"complexity"` // "simple", "moderate", "complex"
	EstimatedTime   string           `json:"estimated_time"`
	Steps           []string         `json:"steps"`
	Commands        []string         `json:"commands,omitempty"`
	ConfigChanges   map[string]interface{} `json:"config_changes,omitempty"`
	ToolsRequired   []string         `json:"tools_required"`
	Automation      AutomationInfo   `json:"automation"`
}

// AutomationInfo describes automation possibilities
type AutomationInfo struct {
	Automatable     bool     `json:"automatable"`
	ScriptGenerated bool     `json:"script_generated"`
	ConfigGenerated bool     `json:"config_generated"`
	MonitoringSetup bool     `json:"monitoring_setup"`
}

// PerformanceEstimate provides performance predictions
type PerformanceEstimate struct {
	TransferSpeed       string    `json:"transfer_speed"`
	TransferTime        string    `json:"transfer_time"`
	NetworkEfficiency   float64   `json:"network_efficiency"`
	ConcurrencyOptimal  int       `json:"concurrency_optimal"`
	PartSizeOptimal     string    `json:"part_size_optimal"`
}

// NewRecommendationEngine creates a new recommendation engine
func NewRecommendationEngine(patternAnalyzer *PatternAnalyzer, costCalculator *S3CostCalculator, 
							 engineRegistry *EngineRegistry, configManager *ConfigManager) *RecommendationEngine {
	return &RecommendationEngine{
		patternAnalyzer: patternAnalyzer,
		costCalculator:  costCalculator,
		engineRegistry:  engineRegistry,
		configManager:   configManager,
	}
}

// GenerateRecommendations analyzes data and generates comprehensive recommendations
func (re *RecommendationEngine) GenerateRecommendations(ctx context.Context, dataPath string) (*RecommendationResult, error) {
	result := &RecommendationResult{
		AnalysisID:    fmt.Sprintf("rec-%d", time.Now().Unix()),
		DataPath:      dataPath,
		Timestamp:     time.Now(),
	}

	// Step 1: Analyze data patterns
	pattern, err := re.patternAnalyzer.AnalyzePattern(ctx, dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze data pattern: %w", err)
	}
	result.DataPattern = pattern

	// Step 2: Analyze costs
	costAnalysis, err := re.costCalculator.AnalyzeCosts(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze costs: %w", err)
	}
	result.CostAnalysis = costAnalysis

	// Step 3: Generate tool recommendations
	result.ToolRecommendations = re.generateToolRecommendations(pattern)

	// Step 4: Generate optimization suggestions
	result.OptimizationSuggestions = re.generateOptimizationSuggestions(pattern, costAnalysis)

	// Step 5: Generate warning alerts
	result.WarningAlerts = re.generateWarningAlerts(pattern, costAnalysis)

	// Step 6: Identify top recommendation and calculate total savings
	result.TopRecommendation = re.identifyTopRecommendation(result.OptimizationSuggestions)
	result.EstimatedSavings = re.calculateTotalSavings(result.OptimizationSuggestions)

	// Step 7: Generate implementation priority
	result.ImplementationPriority = re.generateImplementationPriority(result.OptimizationSuggestions)

	// Step 8: Generate configuration if requested
	result.GeneratedConfig = re.generateProjectConfig(pattern, result)

	return result, nil
}

// generateToolRecommendations recommends optimal tools for the dataset
func (re *RecommendationEngine) generateToolRecommendations(pattern *DataPattern) []ToolRecommendation {
	var recommendations []ToolRecommendation

	// Primary upload/download tool recommendation
	uploadRec := re.recommendUploadTool(pattern)
	recommendations = append(recommendations, uploadRec)

	// Bundling tool recommendation for small files
	if pattern.FileSizes.SmallFiles.CountUnder1MB > 100 {
		bundlingRec := re.recommendBundlingTool(pattern)
		recommendations = append(recommendations, bundlingRec)
	}

	// Compression tool recommendation
	if re.shouldRecommendCompression(pattern) {
		compressionRec := re.recommendCompressionTool(pattern)
		recommendations = append(recommendations, compressionRec)
	}

	// Monitoring tool recommendation
	monitoringRec := re.recommendMonitoringTool(pattern)
	recommendations = append(recommendations, monitoringRec)

	return recommendations
}

// recommendUploadTool recommends the best tool for uploading data
func (re *RecommendationEngine) recommendUploadTool(pattern *DataPattern) ToolRecommendation {
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
	avgFileSizeMB := float64(pattern.FileSizes.MeanSize) / (1024 * 1024)
	fileCount := pattern.TotalFiles

	var tool string
	var alternatives []string
	var reasoning string
	var confidence float64
	var config map[string]interface{}

	// Decision logic based on data characteristics
	if totalSizeGB > 100 && avgFileSizeMB > 50 {
		// Large files, large dataset - s5cmd optimal
		tool = "s5cmd"
		alternatives = []string{"rclone", "aws-cli"}
		reasoning = "Large files and dataset size make s5cmd optimal for maximum throughput"
		confidence = 0.95
		config = map[string]interface{}{
			"concurrency": 20,
			"part_size":   "64MB",
		}
	} else if fileCount > 10000 {
		// Many files - s5cmd for parallel efficiency
		tool = "s5cmd"
		alternatives = []string{"rclone"}
		reasoning = "High file count benefits from s5cmd's parallel processing capabilities"
		confidence = 0.90
		config = map[string]interface{}{
			"concurrency": 15,
			"part_size":   "32MB",
		}
	} else if avgFileSizeMB < 1 {
		// Small files - bundle first, then s5cmd
		tool = "suitcase + s5cmd"
		alternatives = []string{"tar + s5cmd", "rclone"}
		reasoning = "Small files should be bundled first to reduce S3 request costs"
		confidence = 0.85
		config = map[string]interface{}{
			"bundle_size":   "100MB",
			"bundle_count":  100,
			"concurrency":   10,
		}
	} else {
		// General purpose - rclone for reliability
		tool = "rclone"
		alternatives = []string{"s5cmd", "aws-cli"}
		reasoning = "Balanced file sizes make rclone a reliable general-purpose choice"
		confidence = 0.80
		config = map[string]interface{}{
			"concurrency": 8,
			"transfers":   8,
		}
	}

	// Calculate performance estimate
	performance := re.estimatePerformance(pattern, tool, config)

	return ToolRecommendation{
		Task:                "primary_upload",
		RecommendedTool:     tool,
		AlternativeTools:    alternatives,
		Reasoning:           reasoning,
		Confidence:          confidence,
		Configuration:       config,
		ExpectedPerformance: performance,
	}
}

// recommendBundlingTool recommends tools for bundling small files
func (re *RecommendationEngine) recommendBundlingTool(pattern *DataPattern) ToolRecommendation {
	smallFileCount := pattern.FileSizes.SmallFiles.CountUnder1MB
	
	var tool string
	var reasoning string

	if smallFileCount > 10000 {
		tool = "suitcase"
		reasoning = "Suitcase excels at handling large numbers of small files with metadata preservation"
	} else if smallFileCount > 1000 {
		tool = "tar"
		reasoning = "Standard tar archiving suitable for moderate numbers of small files"
	} else {
		tool = "zip"
		reasoning = "ZIP compression good for smaller numbers of files with individual access needs"
	}

	return ToolRecommendation{
		Task:             "file_bundling",
		RecommendedTool:  tool,
		AlternativeTools: []string{"tar", "zip", "7zip"},
		Reasoning:        reasoning,
		Confidence:       0.85,
		Configuration: map[string]interface{}{
			"target_bundle_size": "100MB",
			"preserve_metadata":  true,
		},
	}
}

// recommendCompressionTool recommends compression strategies
func (re *RecommendationEngine) recommendCompressionTool(pattern *DataPattern) ToolRecommendation {
	// Analyze file types for compression potential
	hasTextFiles := false
	hasCompressedFiles := false
	
	for ext, typeInfo := range pattern.FileTypes {
		if typeInfo.Compressible {
			hasTextFiles = true
		}
		if strings.Contains(ext, "gz") || strings.Contains(ext, "zip") {
			hasCompressedFiles = true
		}
	}

	tool := "gzip"
	reasoning := "Standard gzip compression for text-based files"
	
	if hasTextFiles && !hasCompressedFiles {
		tool = "gzip"
		reasoning = "Gzip compression recommended for text-based research data"
	} else if hasCompressedFiles {
		tool = "none"
		reasoning = "Files already compressed - additional compression not beneficial"
	}

	return ToolRecommendation{
		Task:             "compression",
		RecommendedTool:  tool,
		AlternativeTools: []string{"bzip2", "xz", "lz4"},
		Reasoning:        reasoning,
		Confidence:       0.80,
	}
}

// recommendMonitoringTool recommends monitoring solutions
func (re *RecommendationEngine) recommendMonitoringTool(pattern *DataPattern) ToolRecommendation {
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
	
	if totalSizeGB > 1000 {
		return ToolRecommendation{
			Task:             "monitoring",
			RecommendedTool:  "aws-research-wizard monitor + cloudwatch",
			AlternativeTools: []string{"aws-research-wizard monitor"},
			Reasoning:        "Large datasets benefit from comprehensive monitoring with CloudWatch integration",
			Confidence:       0.90,
		}
	}

	return ToolRecommendation{
		Task:             "monitoring",
		RecommendedTool:  "aws-research-wizard monitor",
		AlternativeTools: []string{"manual monitoring"},
		Reasoning:        "Built-in monitoring sufficient for dataset size",
		Confidence:       0.85,
	}
}

// generateOptimizationSuggestions creates optimization suggestions
func (re *RecommendationEngine) generateOptimizationSuggestions(pattern *DataPattern, costAnalysis *CostAnalysis) []OptimizationSuggestion {
	var suggestions []OptimizationSuggestion

	// Small file bundling suggestion
	if pattern.FileSizes.SmallFiles.CountUnder1MB > 100 {
		bundlingSuggestion := re.createBundlingSuggestion(pattern, costAnalysis)
		suggestions = append(suggestions, bundlingSuggestion)
	}

	// Storage class optimization
	if pattern.AccessPatterns.LikelyArchival || pattern.AccessPatterns.LikelyWriteOnce {
		storageClassSuggestion := re.createStorageClassSuggestion(pattern, costAnalysis)
		suggestions = append(suggestions, storageClassSuggestion)
	}

	// Compression suggestion
	if re.shouldRecommendCompression(pattern) {
		compressionSuggestion := re.createCompressionSuggestion(pattern, costAnalysis)
		suggestions = append(suggestions, compressionSuggestion)
	}

	// Tool chain optimization
	toolChainSuggestion := re.createToolChainSuggestion(pattern, costAnalysis)
	suggestions = append(suggestions, toolChainSuggestion)

	// Sort by impact (cost savings)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Impact.CostSavingsMonthly > suggestions[j].Impact.CostSavingsMonthly
	})

	return suggestions
}

// createBundlingSuggestion creates a suggestion for small file bundling
func (re *RecommendationEngine) createBundlingSuggestion(pattern *DataPattern, costAnalysis *CostAnalysis) OptimizationSuggestion {
	smallFileCount := pattern.FileSizes.SmallFiles.CountUnder1MB
	currentCost := costAnalysis.Scenarios[0].MonthlyCosts.Total
	
	// Find bundled scenario cost
	var bundledCost float64
	for _, scenario := range costAnalysis.Scenarios {
		if scenario.Name == "Bundled Small Files" {
			bundledCost = scenario.MonthlyCosts.Total
			break
		}
	}
	
	savings := currentCost - bundledCost
	savingsPercent := (savings / currentCost) * 100

	return OptimizationSuggestion{
		ID:          "bundling-small-files",
		Type:        "bundling",
		Title:       "Bundle Small Files",
		Description: fmt.Sprintf("Bundle %d small files to reduce S3 request costs by %.0f%%", smallFileCount, savingsPercent),
		Impact: ImpactAssessment{
			CostSavingsMonthly:     savings,
			CostSavingsPercent:     savingsPercent,
			PerformanceImprovement: 25.0, // Faster uploads with fewer requests
			TimeToTransfer:         "Reduced by 40%",
			RiskLevel:              "low",
			Confidence:             0.90,
		},
		Implementation: Implementation{
			Complexity:    "moderate",
			EstimatedTime: "2-4 hours",
			Steps: []string{
				"Install Suitcase bundling tool",
				"Configure bundling parameters (100 files per bundle)",
				"Test bundling with sample data",
				"Update upload workflow to include bundling step",
				"Verify bundle integrity and metadata preservation",
			},
			Commands: []string{
				"# Install Suitcase",
				"pip install suitcase",
				"# Bundle small files",
				"suitcase pack --target-size 100MB --preserve-metadata ./small-files/ ./bundles/",
				"# Upload bundles",
				"aws-research-wizard data upload ./bundles/ s3://bucket/bundles/",
			},
			ToolsRequired: []string{"suitcase", "aws-research-wizard"},
			Automation: AutomationInfo{
				Automatable:     true,
				ScriptGenerated: true,
				ConfigGenerated: true,
				MonitoringSetup: true,
			},
		},
		Prerequisites: []string{
			"Python environment with pip",
			"Sufficient local storage for temporary bundles",
		},
	}
}

// createStorageClassSuggestion creates storage class optimization suggestion
func (re *RecommendationEngine) createStorageClassSuggestion(pattern *DataPattern, costAnalysis *CostAnalysis) OptimizationSuggestion {
	recommendedClass := pattern.Efficiency.RecommendedStorageClass
	savings := pattern.Efficiency.StorageClassSavings

	var description string
	var steps []string

	switch recommendedClass {
	case "GLACIER":
		description = "Archive old data to Glacier for 82% storage cost reduction"
		steps = []string{
			"Create lifecycle policy for automatic archival",
			"Set up 90-day transition rule",
			"Configure monitoring for archived data",
			"Test retrieval process",
		}
	case "STANDARD_IA":
		description = "Move infrequently accessed data to Standard-IA for 46% savings"
		steps = []string{
			"Analyze access patterns",
			"Create lifecycle policy for IA transition",
			"Set up 30-day minimum storage monitoring",
			"Configure alerts for retrieval costs",
		}
	default:
		description = "Optimize storage class based on access patterns"
		steps = []string{
			"Analyze current access patterns",
			"Choose appropriate storage class",
			"Implement lifecycle policies",
		}
	}

	return OptimizationSuggestion{
		ID:          "storage-class-optimization",
		Type:        "storage_class",
		Title:       fmt.Sprintf("Use %s Storage Class", recommendedClass),
		Description: description,
		Impact: ImpactAssessment{
			CostSavingsMonthly: savings,
			CostSavingsPercent: (savings / costAnalysis.Scenarios[0].MonthlyCosts.Total) * 100,
			RiskLevel:          "low",
			Confidence:         0.85,
		},
		Implementation: Implementation{
			Complexity:    "simple",
			EstimatedTime: "1-2 hours",
			Steps:         steps,
			ToolsRequired: []string{"aws-cli", "aws console"},
			Automation: AutomationInfo{
				Automatable:     true,
				ConfigGenerated: true,
			},
		},
	}
}

// createCompressionSuggestion creates compression optimization suggestion
func (re *RecommendationEngine) createCompressionSuggestion(pattern *DataPattern, costAnalysis *CostAnalysis) OptimizationSuggestion {
	// Estimate compression savings
	totalSize := float64(pattern.TotalSize)
	compressedSize := totalSize * 0.7 // Assume 30% compression
	storageSavings := (totalSize - compressedSize) / (1024 * 1024 * 1024) * 0.023 // $0.023 per GB
	
	return OptimizationSuggestion{
		ID:          "enable-compression",
		Type:        "compression",
		Title:       "Enable Data Compression",
		Description: "Compress data before upload to reduce storage costs by 30%",
		Impact: ImpactAssessment{
			CostSavingsMonthly: storageSavings,
			CostSavingsPercent: 30.0,
			RiskLevel:          "low",
			Confidence:         0.80,
		},
		Implementation: Implementation{
			Complexity:    "simple",
			EstimatedTime: "30 minutes",
			Steps: []string{
				"Identify compressible file types",
				"Add compression step to upload workflow",
				"Test compression ratios",
				"Update documentation",
			},
			Commands: []string{
				"# Compress files before upload",
				"gzip -r ./data/",
				"# Upload compressed files",
				"aws-research-wizard data upload ./data/ s3://bucket/compressed-data/",
			},
			ToolsRequired: []string{"gzip"},
			Automation: AutomationInfo{
				Automatable:     true,
				ScriptGenerated: true,
			},
		},
	}
}

// createToolChainSuggestion creates tool chain optimization suggestion
func (re *RecommendationEngine) createToolChainSuggestion(pattern *DataPattern, costAnalysis *CostAnalysis) OptimizationSuggestion {
	return OptimizationSuggestion{
		ID:          "optimize-tool-chain",
		Type:        "tool_chain",
		Title:       "Optimize Transfer Tool Chain",
		Description: "Use s5cmd + bundling for optimal performance and cost",
		Impact: ImpactAssessment{
			PerformanceImprovement: 300.0, // 3x faster than aws CLI
			TimeToTransfer:         "Reduced by 67%",
			RiskLevel:              "low",
			Confidence:             0.85,
		},
		Implementation: Implementation{
			Complexity:    "moderate",
			EstimatedTime: "1-2 hours",
			Steps: []string{
				"Install s5cmd",
				"Configure optimal concurrency settings",
				"Set up monitoring",
				"Create automated workflow",
			},
			Commands: []string{
				"# Install s5cmd",
				"curl -L https://github.com/peak/s5cmd/releases/latest/download/s5cmd_*_linux_amd64.tar.gz | tar -xz",
				"# Configure and use",
				"aws-research-wizard data upload --engine s5cmd --concurrency 20 ./data/ s3://bucket/",
			},
			ToolsRequired: []string{"s5cmd", "aws-research-wizard"},
			Automation: AutomationInfo{
				Automatable:     true,
				ScriptGenerated: true,
				ConfigGenerated: true,
			},
		},
	}
}

// generateWarningAlerts creates warning alerts for potential issues
func (re *RecommendationEngine) generateWarningAlerts(pattern *DataPattern, costAnalysis *CostAnalysis) []WarningAlert {
	var alerts []WarningAlert

	// Small file warning
	if pattern.FileSizes.SmallFiles.CountUnder1MB > 1000 {
		severity := "warning"
		if pattern.FileSizes.SmallFiles.CountUnder1MB > 10000 {
			severity = "critical"
		}

		alerts = append(alerts, WarningAlert{
			Severity:    severity,
			Category:    "cost",
			Title:       "High Small File Count Detected",
			Description: fmt.Sprintf("Found %d files under 1MB which will result in high S3 request costs", pattern.FileSizes.SmallFiles.CountUnder1MB),
			Impact:      fmt.Sprintf("Estimated extra cost: $%.2f/month", pattern.FileSizes.SmallFiles.PotentialSavings),
			Solution:    "Bundle small files using Suitcase or tar before uploading",
			LearnMoreURL: "https://docs.aws.amazon.com/s3/latest/userguide/optimizing-performance.html",
		})
	}

	// Large dataset without monitoring warning
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
	if totalSizeGB > 500 {
		alerts = append(alerts, WarningAlert{
			Severity:    "warning",
			Category:    "performance",
			Title:       "Large Dataset Without Monitoring",
			Description: fmt.Sprintf("Dataset size of %.1f GB should have monitoring enabled", totalSizeGB),
			Impact:      "Risk of unnoticed failures or performance issues",
			Solution:    "Enable CloudWatch monitoring and set up alerts",
		})
	}

	// Cost threshold warning
	if costAnalysis.Scenarios[0].MonthlyCosts.Total > 500 {
		alerts = append(alerts, WarningAlert{
			Severity:    "warning",
			Category:    "cost",
			Title:       "High Monthly Cost Detected",
			Description: fmt.Sprintf("Estimated monthly cost of $%.2f is above typical research budgets", costAnalysis.Scenarios[0].MonthlyCosts.Total),
			Impact:      "May exceed research budget limits",
			Solution:    "Consider storage class optimization and lifecycle policies",
		})
	}

	return alerts
}

// Helper functions

func (re *RecommendationEngine) shouldRecommendCompression(pattern *DataPattern) bool {
	compressibleSize := int64(0)
	for _, typeInfo := range pattern.FileTypes {
		if typeInfo.Compressible {
			compressibleSize += typeInfo.TotalSize
		}
	}
	
	// Recommend compression if >50% of data is compressible
	return float64(compressibleSize)/float64(pattern.TotalSize) > 0.5
}

func (re *RecommendationEngine) estimatePerformance(pattern *DataPattern, tool string, config map[string]interface{}) PerformanceEstimate {
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
	
	// Base transfer speeds (MB/s) by tool
	baseSpeeds := map[string]float64{
		"s5cmd":     80.0,
		"rclone":    60.0,
		"aws-cli":   25.0,
		"suitcase":  40.0,
	}
	
	speed := baseSpeeds[tool]
	if speed == 0 {
		speed = 50.0 // Default
	}
	
	// Adjust for concurrency
	if concurrency, ok := config["concurrency"].(int); ok && concurrency > 1 {
		speed *= math.Min(float64(concurrency), 20) / 10 // Max 2x improvement
	}
	
	transferTimeHours := totalSizeGB * 1024 / speed / 3600
	
	return PerformanceEstimate{
		TransferSpeed:      fmt.Sprintf("%.0f MB/s", speed),
		TransferTime:       fmt.Sprintf("%.1f hours", transferTimeHours),
		NetworkEfficiency:  85.0,
		ConcurrencyOptimal: int(math.Min(float64(pattern.TotalFiles/100), 20)),
		PartSizeOptimal:    "32MB",
	}
}

func (re *RecommendationEngine) identifyTopRecommendation(suggestions []OptimizationSuggestion) *OptimizationSuggestion {
	if len(suggestions) == 0 {
		return nil
	}
	
	// Return the one with highest cost savings
	return &suggestions[0]
}

func (re *RecommendationEngine) calculateTotalSavings(suggestions []OptimizationSuggestion) float64 {
	total := 0.0
	for _, suggestion := range suggestions {
		total += suggestion.Impact.CostSavingsMonthly
	}
	return total
}

func (re *RecommendationEngine) generateImplementationPriority(suggestions []OptimizationSuggestion) []string {
	// Sort by impact vs complexity
	priorities := make([]string, len(suggestions))
	for i, suggestion := range suggestions {
		priorities[i] = suggestion.Title
	}
	return priorities
}

func (re *RecommendationEngine) generateProjectConfig(pattern *DataPattern, result *RecommendationResult) *ProjectConfig {
	// Generate a project configuration based on recommendations
	config := &ProjectConfig{
		Project: ProjectInfo{
			Name:        "generated-project",
			Domain:      strings.Join(pattern.DomainHints.DetectedDomains, ","),
			Description: "Auto-generated configuration based on data analysis",
		},
		DataProfiles: make(map[string]DataProfile),
		Workflows:    make([]Workflow, 0),
	}
	
	// Add data profile
	config.DataProfiles["main_dataset"] = DataProfile{
		Path:        result.DataPath,
		FileCount:   pattern.TotalFiles,
		AvgFileSize: formatBytes(pattern.FileSizes.MeanSize),
		AccessPattern: "write_once_read_many", // Default assumption
	}
	
	// Add workflow based on top recommendation
	if result.TopRecommendation != nil {
		workflow := Workflow{
			Name:        "optimized_upload",
			Source:      "main_dataset",
			Destination: "s3://bucket/optimized-data/",
			Triggers:    []string{"manual"},
		}
		config.Workflows = append(config.Workflows, workflow)
	}
	
	return config
}