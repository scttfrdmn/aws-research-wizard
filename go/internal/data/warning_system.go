package data

import (
	"context"
	"fmt"
	"time"
)

// WarningSystem provides proactive alerts for data management anti-patterns
type WarningSystem struct {
	thresholds WarningThresholds
	rules      []WarningRule
}

// WarningThresholds defines thresholds for various warning conditions
type WarningThresholds struct {
	// File count thresholds
	SmallFileCountWarning  int64 `json:"small_file_count_warning"`  // 1000
	SmallFileCountCritical int64 `json:"small_file_count_critical"` // 10000

	// Size thresholds
	LargeDatasetGB float64 `json:"large_dataset_gb"` // 500 GB
	HugeDatasetGB  float64 `json:"huge_dataset_gb"`  // 5000 GB

	// Cost thresholds
	MonthlyCostWarning  float64 `json:"monthly_cost_warning"`  // $100
	MonthlyCostCritical float64 `json:"monthly_cost_critical"` // $1000

	// Performance thresholds
	SlowTransferMBps     float64 `json:"slow_transfer_mbps"`      // 10 MB/s
	HighRequestCostRatio float64 `json:"high_request_cost_ratio"` // 0.5 (requests > 50% of total cost)

	// Efficiency thresholds
	LowNetworkEfficiency float64 `json:"low_network_efficiency"` // 30%
	HighSmallFileRatio   float64 `json:"high_small_file_ratio"`  // 70%
}

// WarningRule defines a rule for detecting anti-patterns
type WarningRule struct {
	ID          string                                   `json:"id"`
	Name        string                                   `json:"name"`
	Category    string                                   `json:"category"`
	Severity    string                                   `json:"severity"`
	Description string                                   `json:"description"`
	Condition   func(*DataPattern, *CostAnalysis) bool   `json:"-"`
	Message     func(*DataPattern, *CostAnalysis) string `json:"-"`
	Solution    string                                   `json:"solution"`
	LearnMore   string                                   `json:"learn_more"`
	Enabled     bool                                     `json:"enabled"`
}

// AntiPattern represents a detected anti-pattern with specific details
type AntiPattern struct {
	RuleID      string                 `json:"rule_id"`
	Severity    string                 `json:"severity"`
	Category    string                 `json:"category"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Impact      AntiPatternImpact      `json:"impact"`
	Evidence    map[string]interface{} `json:"evidence"`
	Remediation Remediation            `json:"remediation"`
	DetectedAt  time.Time              `json:"detected_at"`
}

// AntiPatternImpact quantifies the impact of an anti-pattern
type AntiPatternImpact struct {
	CostImpact        float64 `json:"cost_impact_monthly"`
	PerformanceImpact string  `json:"performance_impact"`
	ReliabilityImpact string  `json:"reliability_impact"`
	MaintenanceImpact string  `json:"maintenance_impact"`
	SeverityScore     int     `json:"severity_score"` // 1-10
}

// Remediation provides specific steps to fix an anti-pattern
type Remediation struct {
	PrimaryAction   string   `json:"primary_action"`
	Steps           []string `json:"steps"`
	Commands        []string `json:"commands,omitempty"`
	EstimatedTime   string   `json:"estimated_time"`
	Complexity      string   `json:"complexity"`
	Prerequisites   []string `json:"prerequisites"`
	AutomationLevel string   `json:"automation_level"` // "manual", "semi-automated", "fully-automated"
}

// WarningReport contains all detected warnings and anti-patterns
type WarningReport struct {
	ReportID    string    `json:"report_id"`
	GeneratedAt time.Time `json:"generated_at"`
	DataPath    string    `json:"data_path"`

	// Summary
	Summary WarningReportSummary `json:"summary"`

	// Detected issues
	CriticalIssues []AntiPattern `json:"critical_issues"`
	Warnings       []AntiPattern `json:"warnings"`
	InfoAlerts     []AntiPattern `json:"info_alerts"`

	// Recommendations
	QuickFixes      []AntiPattern `json:"quick_fixes"`
	LongTermActions []AntiPattern `json:"long_term_actions"`

	// Scores
	OverallScore   int            `json:"overall_score"` // 1-100
	CategoryScores map[string]int `json:"category_scores"`
}

// WarningReportSummary provides high-level summary
type WarningReportSummary struct {
	TotalIssues          int     `json:"total_issues"`
	CriticalCount        int     `json:"critical_count"`
	WarningCount         int     `json:"warning_count"`
	InfoCount            int     `json:"info_count"`
	EstimatedWastedCost  float64 `json:"estimated_wasted_cost_monthly"`
	MaxPotentialSavings  float64 `json:"max_potential_savings_monthly"`
	HighestSeverityIssue string  `json:"highest_severity_issue"`
	RecommendedFirstStep string  `json:"recommended_first_step"`
}

// NewWarningSystem creates a new warning system with default thresholds
func NewWarningSystem() *WarningSystem {
	ws := &WarningSystem{
		thresholds: WarningThresholds{
			SmallFileCountWarning:  1000,
			SmallFileCountCritical: 10000,
			LargeDatasetGB:         500,
			HugeDatasetGB:          5000,
			MonthlyCostWarning:     100,
			MonthlyCostCritical:    1000,
			SlowTransferMBps:       10,
			HighRequestCostRatio:   0.5,
			LowNetworkEfficiency:   30,
			HighSmallFileRatio:     70,
		},
		rules: make([]WarningRule, 0),
	}

	ws.initializeDefaultRules()
	return ws
}

// initializeDefaultRules sets up the default warning rules
func (ws *WarningSystem) initializeDefaultRules() {
	ws.rules = []WarningRule{
		{
			ID:          "small-file-explosion",
			Name:        "Small File Explosion",
			Category:    "cost",
			Severity:    "critical",
			Description: "Large number of small files will cause excessive S3 request costs",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				return pattern.FileSizes.SmallFiles.CountUnder1MB > ws.thresholds.SmallFileCountCritical
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				return fmt.Sprintf("Detected %d files under 1MB, estimated extra cost: $%.2f/month",
					pattern.FileSizes.SmallFiles.CountUnder1MB,
					pattern.FileSizes.SmallFiles.PotentialSavings)
			},
			Solution:  "Bundle small files using Suitcase or tar before uploading to S3",
			LearnMore: "https://docs.aws.amazon.com/s3/latest/userguide/optimizing-performance.html",
			Enabled:   true,
		},
		{
			ID:          "inefficient-tool-usage",
			Name:        "Inefficient Transfer Tool",
			Category:    "performance",
			Severity:    "warning",
			Description: "Using suboptimal tools for data characteristics",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				// Detect when large files would benefit from s5cmd
				avgFileSizeMB := float64(pattern.FileSizes.MeanSize) / (1024 * 1024)
				return avgFileSizeMB > 100 && pattern.TotalFiles > 100
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				avgFileSizeMB := float64(pattern.FileSizes.MeanSize) / (1024 * 1024)
				return fmt.Sprintf("Large files (avg %.0f MB) would benefit from s5cmd instead of aws CLI", avgFileSizeMB)
			},
			Solution:  "Use s5cmd with high concurrency for large file uploads",
			LearnMore: "https://github.com/peak/s5cmd",
			Enabled:   true,
		},
		{
			ID:          "compression-opportunity",
			Name:        "Missed Compression Opportunity",
			Category:    "cost",
			Severity:    "info",
			Description: "Compressible files not being compressed before upload",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				// Check if majority of files are compressible
				compressibleSize := int64(0)
				for _, typeInfo := range pattern.FileTypes {
					if typeInfo.Compressible && typeInfo.CompressionEst < 0.8 {
						compressibleSize += typeInfo.TotalSize
					}
				}
				return float64(compressibleSize)/float64(pattern.TotalSize) > 0.6
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				return "60%+ of data is compressible - could save ~30% storage costs"
			},
			Solution:  "Enable compression in upload pipeline for text-based research data",
			LearnMore: "https://docs.aws.amazon.com/s3/latest/userguide/compression.html",
			Enabled:   true,
		},
		{
			ID:          "suboptimal-storage-class",
			Name:        "Suboptimal Storage Class",
			Category:    "cost",
			Severity:    "warning",
			Description: "Data access patterns suggest different storage class would be more cost-effective",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				return pattern.AccessPatterns.LikelyArchival || pattern.AccessPatterns.LikelyWriteOnce
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				if pattern.AccessPatterns.LikelyArchival {
					return "Access patterns suggest Glacier storage class for 75% cost savings"
				}
				return "Infrequent access patterns suggest Standard-IA for 45% savings"
			},
			Solution:  "Set up lifecycle policies to automatically transition to appropriate storage class",
			LearnMore: "https://docs.aws.amazon.com/s3/latest/userguide/lifecycle-transition-general-considerations.html",
			Enabled:   true,
		},
		{
			ID:          "high-cost-without-monitoring",
			Name:        "High Cost Without Monitoring",
			Category:    "reliability",
			Severity:    "warning",
			Description: "High-cost dataset lacks proper monitoring and alerting",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				return cost.Scenarios[0].MonthlyCosts.Total > ws.thresholds.MonthlyCostWarning
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				return fmt.Sprintf("Monthly cost of $%.2f requires monitoring setup", cost.Scenarios[0].MonthlyCosts.Total)
			},
			Solution:  "Set up CloudWatch monitoring and cost alerts",
			LearnMore: "https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/monitor_estimated_charges_with_cloudwatch.html",
			Enabled:   true,
		},
		{
			ID:          "genomics-specific-inefficiency",
			Name:        "Genomics Data Inefficiency",
			Category:    "domain-specific",
			Severity:    "info",
			Description: "Genomics data could be optimized with domain-specific practices",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				// Check if genomics domain detected
				for _, domain := range pattern.DomainHints.DetectedDomains {
					if domain == "genomics" {
						// Check for uncompressed FASTQ files
						for ext, typeInfo := range pattern.FileTypes {
							if ext == ".fastq" && typeInfo.Percentage > 20 {
								return true
							}
						}
					}
				}
				return false
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				return "Uncompressed FASTQ files detected - genomics data typically compresses to 25% of original size"
			},
			Solution:  "Compress FASTQ files with gzip before upload, consider bgzip for indexed access",
			LearnMore: "https://www.ncbi.nlm.nih.gov/sra/docs/submitformats/#fastq-files",
			Enabled:   true,
		},
		{
			ID:          "climate-data-chunking",
			Name:        "Climate Data Not Optimally Chunked",
			Category:    "domain-specific",
			Severity:    "info",
			Description: "Climate data files could benefit from optimal chunking strategies",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				// Check for climate domain with NetCDF files
				for _, domain := range pattern.DomainHints.DetectedDomains {
					if domain == "climate" {
						for ext := range pattern.FileTypes {
							if ext == ".nc" || ext == ".netcdf" {
								return true
							}
						}
					}
				}
				return false
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				return "NetCDF files detected - consider rechunking for optimal access patterns"
			},
			Solution:  "Use tools like NCO or xarray to rechunk NetCDF files for your access patterns",
			LearnMore: "https://nco.sourceforge.net/nco.html#Chunking",
			Enabled:   true,
		},
		{
			ID:          "ml-model-storage-inefficiency",
			Name:        "ML Model Storage Inefficiency",
			Category:    "domain-specific",
			Severity:    "info",
			Description: "Machine learning models and checkpoints could use different storage strategies",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				// Check for ML domain with model files
				for _, domain := range pattern.DomainHints.DetectedDomains {
					if domain == "machine_learning" {
						for ext := range pattern.FileTypes {
							if ext == ".model" || ext == ".ckpt" || ext == ".pkl" || ext == ".h5" {
								return true
							}
						}
					}
				}
				return false
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				return "ML model files detected - distinguish between active models and archived checkpoints"
			},
			Solution:  "Use Standard storage for active models, Standard-IA for checkpoints, Glacier for archived experiments",
			LearnMore: "https://aws.amazon.com/blogs/machine-learning/",
			Enabled:   true,
		},
		{
			ID:          "request-cost-dominance",
			Name:        "Request Costs Dominate Storage Costs",
			Category:    "cost",
			Severity:    "critical",
			Description: "Request costs are unusually high compared to storage costs",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				total := cost.Scenarios[0].MonthlyCosts.Total
				requests := cost.Scenarios[0].MonthlyCosts.Requests
				return requests/total > ws.thresholds.HighRequestCostRatio
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				requests := cost.Scenarios[0].MonthlyCosts.Requests
				total := cost.Scenarios[0].MonthlyCosts.Total
				percentage := (requests / total) * 100
				return fmt.Sprintf("Request costs are %.0f%% of total costs - indicates small file problem", percentage)
			},
			Solution:  "Urgently implement file bundling to reduce request counts by 90%+",
			LearnMore: "https://docs.aws.amazon.com/s3/latest/userguide/optimizing-performance.html",
			Enabled:   true,
		},
		{
			ID:          "large-dataset-no-versioning",
			Name:        "Large Dataset Without Versioning Strategy",
			Category:    "reliability",
			Severity:    "warning",
			Description: "Large valuable dataset lacks versioning or backup strategy",
			Condition: func(pattern *DataPattern, cost *CostAnalysis) bool {
				totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
				return totalSizeGB > ws.thresholds.LargeDatasetGB
			},
			Message: func(pattern *DataPattern, cost *CostAnalysis) string {
				totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
				return fmt.Sprintf("%.1f GB dataset should have versioning and backup strategy", totalSizeGB)
			},
			Solution:  "Enable S3 versioning and set up cross-region replication for critical data",
			LearnMore: "https://docs.aws.amazon.com/s3/latest/userguide/Versioning.html",
			Enabled:   true,
		},
	}
}

// AnalyzePattern analyzes a data pattern and generates a warning report
func (ws *WarningSystem) AnalyzePattern(ctx context.Context, pattern *DataPattern, costAnalysis *CostAnalysis) (*WarningReport, error) {
	report := &WarningReport{
		ReportID:        fmt.Sprintf("warn-%d", time.Now().Unix()),
		GeneratedAt:     time.Now(),
		DataPath:        pattern.AnalyzedPath,
		CriticalIssues:  make([]AntiPattern, 0),
		Warnings:        make([]AntiPattern, 0),
		InfoAlerts:      make([]AntiPattern, 0),
		QuickFixes:      make([]AntiPattern, 0),
		LongTermActions: make([]AntiPattern, 0),
		CategoryScores:  make(map[string]int),
	}

	// Run all enabled warning rules
	for _, rule := range ws.rules {
		if !rule.Enabled {
			continue
		}

		if rule.Condition(pattern, costAnalysis) {
			antiPattern := ws.createAntiPattern(rule, pattern, costAnalysis)

			// Categorize by severity
			switch antiPattern.Severity {
			case "critical":
				report.CriticalIssues = append(report.CriticalIssues, antiPattern)
			case "warning":
				report.Warnings = append(report.Warnings, antiPattern)
			case "info":
				report.InfoAlerts = append(report.InfoAlerts, antiPattern)
			}

			// Categorize by remediation complexity
			if antiPattern.Remediation.Complexity == "simple" {
				report.QuickFixes = append(report.QuickFixes, antiPattern)
			} else {
				report.LongTermActions = append(report.LongTermActions, antiPattern)
			}
		}
	}

	// Generate summary
	report.Summary = ws.generateSummary(report, pattern, costAnalysis)

	// Calculate scores
	report.OverallScore = ws.calculateOverallScore(report)
	report.CategoryScores = ws.calculateCategoryScores(report)

	return report, nil
}

// createAntiPattern creates an AntiPattern from a triggered rule
func (ws *WarningSystem) createAntiPattern(rule WarningRule, pattern *DataPattern, costAnalysis *CostAnalysis) AntiPattern {
	antiPattern := AntiPattern{
		RuleID:      rule.ID,
		Severity:    rule.Severity,
		Category:    rule.Category,
		Title:       rule.Name,
		Description: rule.Message(pattern, costAnalysis),
		DetectedAt:  time.Now(),
		Evidence:    make(map[string]interface{}),
	}

	// Calculate impact based on rule type
	impact := ws.calculateImpact(rule, pattern, costAnalysis)
	antiPattern.Impact = impact

	// Generate remediation
	remediation := ws.generateRemediation(rule, pattern, costAnalysis)
	antiPattern.Remediation = remediation

	// Add evidence based on rule type
	antiPattern.Evidence = ws.gatherEvidence(rule, pattern, costAnalysis)

	return antiPattern
}

// calculateImpact calculates the impact of an anti-pattern
func (ws *WarningSystem) calculateImpact(rule WarningRule, pattern *DataPattern, costAnalysis *CostAnalysis) AntiPatternImpact {
	impact := AntiPatternImpact{}

	switch rule.ID {
	case "small-file-explosion":
		impact.CostImpact = pattern.FileSizes.SmallFiles.PotentialSavings
		impact.PerformanceImpact = "Very High - 10x slower uploads"
		impact.ReliabilityImpact = "Medium - more failure points"
		impact.MaintenanceImpact = "High - complex file management"
		impact.SeverityScore = 9

	case "inefficient-tool-usage":
		impact.PerformanceImpact = "High - 3x slower transfers"
		impact.CostImpact = 0 // No direct cost impact
		impact.ReliabilityImpact = "Low"
		impact.MaintenanceImpact = "Medium"
		impact.SeverityScore = 6

	case "compression-opportunity":
		totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
		impact.CostImpact = totalSizeGB * 0.3 * 0.023 // 30% compression savings
		impact.PerformanceImpact = "Medium - faster uploads"
		impact.ReliabilityImpact = "Low"
		impact.MaintenanceImpact = "Low"
		impact.SeverityScore = 4

	case "request-cost-dominance":
		requests := costAnalysis.Scenarios[0].MonthlyCosts.Requests
		impact.CostImpact = requests * 0.9 // 90% savings possible
		impact.PerformanceImpact = "Very High - 10x slower"
		impact.ReliabilityImpact = "High - many failure points"
		impact.MaintenanceImpact = "Very High"
		impact.SeverityScore = 10

	default:
		impact.SeverityScore = 5
		impact.PerformanceImpact = "Medium"
		impact.ReliabilityImpact = "Medium"
		impact.MaintenanceImpact = "Medium"
	}

	return impact
}

// generateRemediation generates specific remediation steps
func (ws *WarningSystem) generateRemediation(rule WarningRule, pattern *DataPattern, costAnalysis *CostAnalysis) Remediation {
	remediation := Remediation{
		PrimaryAction: rule.Solution,
	}

	switch rule.ID {
	case "small-file-explosion":
		remediation.Steps = []string{
			"Install Suitcase bundling tool",
			"Configure bundling parameters (target 100MB bundles)",
			"Test bundling process with sample data",
			"Update upload workflow to include bundling",
			"Monitor bundle sizes and adjust parameters",
		}
		remediation.Commands = []string{
			"pip install suitcase",
			"suitcase pack --target-size 100MB ./small-files/ ./bundles/",
			"aws-research-wizard data upload ./bundles/ s3://bucket/",
		}
		remediation.EstimatedTime = "2-4 hours"
		remediation.Complexity = "moderate"
		remediation.AutomationLevel = "semi-automated"
		remediation.Prerequisites = []string{"Python environment", "Sufficient disk space"}

	case "inefficient-tool-usage":
		remediation.Steps = []string{
			"Install s5cmd",
			"Configure optimal concurrency settings",
			"Test transfer performance",
			"Update automation scripts",
		}
		remediation.Commands = []string{
			"curl -L https://github.com/peak/s5cmd/releases/latest/download/s5cmd_*_linux_amd64.tar.gz | tar -xz",
			"s5cmd --numworkers 20 sync ./data/ s3://bucket/data/",
		}
		remediation.EstimatedTime = "1 hour"
		remediation.Complexity = "simple"
		remediation.AutomationLevel = "fully-automated"

	case "compression-opportunity":
		remediation.Steps = []string{
			"Identify compressible file types",
			"Test compression ratios",
			"Add compression to upload pipeline",
			"Monitor storage savings",
		}
		remediation.Commands = []string{
			"gzip -r ./data/",
			"aws-research-wizard data upload ./data/ s3://bucket/compressed/",
		}
		remediation.EstimatedTime = "30 minutes"
		remediation.Complexity = "simple"
		remediation.AutomationLevel = "fully-automated"

	default:
		remediation.Steps = []string{"Follow solution guidance"}
		remediation.EstimatedTime = "1-2 hours"
		remediation.Complexity = "moderate"
		remediation.AutomationLevel = "manual"
	}

	return remediation
}

// gatherEvidence collects specific evidence for an anti-pattern
func (ws *WarningSystem) gatherEvidence(rule WarningRule, pattern *DataPattern, costAnalysis *CostAnalysis) map[string]interface{} {
	evidence := make(map[string]interface{})

	switch rule.ID {
	case "small-file-explosion":
		evidence["small_file_count"] = pattern.FileSizes.SmallFiles.CountUnder1MB
		evidence["small_file_percentage"] = pattern.FileSizes.SmallFiles.PercentageSmall
		evidence["potential_savings"] = pattern.FileSizes.SmallFiles.PotentialSavings
		evidence["request_cost_monthly"] = costAnalysis.Scenarios[0].MonthlyCosts.Requests

	case "inefficient-tool-usage":
		evidence["avg_file_size_mb"] = float64(pattern.FileSizes.MeanSize) / (1024 * 1024)
		evidence["total_file_count"] = pattern.TotalFiles
		evidence["estimated_speedup"] = "3x faster with s5cmd"

	case "compression-opportunity":
		compressibleSize := int64(0)
		for _, typeInfo := range pattern.FileTypes {
			if typeInfo.Compressible {
				compressibleSize += typeInfo.TotalSize
			}
		}
		evidence["compressible_size_gb"] = float64(compressibleSize) / (1024 * 1024 * 1024)
		evidence["compressible_percentage"] = float64(compressibleSize) / float64(pattern.TotalSize) * 100
		evidence["estimated_compression_ratio"] = 0.7
	}

	return evidence
}

// generateSummary creates a summary of the warning report
func (ws *WarningSystem) generateSummary(report *WarningReport, pattern *DataPattern, costAnalysis *CostAnalysis) WarningReportSummary {
	summary := WarningReportSummary{
		CriticalCount: len(report.CriticalIssues),
		WarningCount:  len(report.Warnings),
		InfoCount:     len(report.InfoAlerts),
	}

	summary.TotalIssues = summary.CriticalCount + summary.WarningCount + summary.InfoCount

	// Calculate potential savings
	maxSavings := 0.0
	for _, issue := range append(report.CriticalIssues, report.Warnings...) {
		if issue.Impact.CostImpact > maxSavings {
			maxSavings = issue.Impact.CostImpact
		}
		summary.EstimatedWastedCost += issue.Impact.CostImpact
	}
	summary.MaxPotentialSavings = maxSavings

	// Identify highest severity issue
	if len(report.CriticalIssues) > 0 {
		summary.HighestSeverityIssue = report.CriticalIssues[0].Title
		summary.RecommendedFirstStep = report.CriticalIssues[0].Remediation.PrimaryAction
	} else if len(report.Warnings) > 0 {
		summary.HighestSeverityIssue = report.Warnings[0].Title
		summary.RecommendedFirstStep = report.Warnings[0].Remediation.PrimaryAction
	}

	return summary
}

// calculateOverallScore calculates an overall efficiency score (1-100)
func (ws *WarningSystem) calculateOverallScore(report *WarningReport) int {
	score := 100

	// Deduct points for issues
	for _, issue := range report.CriticalIssues {
		score -= issue.Impact.SeverityScore * 2 // Critical issues count double
	}
	for _, issue := range report.Warnings {
		score -= issue.Impact.SeverityScore
	}
	for _, issue := range report.InfoAlerts {
		score -= issue.Impact.SeverityScore / 2 // Info issues count half
	}

	if score < 0 {
		score = 0
	}

	return score
}

// calculateCategoryScores calculates scores by category
func (ws *WarningSystem) calculateCategoryScores(report *WarningReport) map[string]int {
	categoryScores := map[string]int{
		"cost":            100,
		"performance":     100,
		"reliability":     100,
		"domain-specific": 100,
	}

	// Deduct points by category
	allIssues := append(append(report.CriticalIssues, report.Warnings...), report.InfoAlerts...)
	for _, issue := range allIssues {
		multiplier := 1
		if issue.Severity == "critical" {
			multiplier = 2
		} else if issue.Severity == "info" {
			multiplier = 1
		}

		categoryScores[issue.Category] -= issue.Impact.SeverityScore * multiplier
		if categoryScores[issue.Category] < 0 {
			categoryScores[issue.Category] = 0
		}
	}

	return categoryScores
}

// GetEnabledRules returns all enabled warning rules
func (ws *WarningSystem) GetEnabledRules() []WarningRule {
	enabled := make([]WarningRule, 0)
	for _, rule := range ws.rules {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	return enabled
}

// UpdateThresholds updates warning thresholds
func (ws *WarningSystem) UpdateThresholds(thresholds WarningThresholds) {
	ws.thresholds = thresholds
}

// EnableRule enables a specific warning rule
func (ws *WarningSystem) EnableRule(ruleID string) {
	for i := range ws.rules {
		if ws.rules[i].ID == ruleID {
			ws.rules[i].Enabled = true
			break
		}
	}
}

// DisableRule disables a specific warning rule
func (ws *WarningSystem) DisableRule(ruleID string) {
	for i := range ws.rules {
		if ws.rules[i].ID == ruleID {
			ws.rules[i].Enabled = false
			break
		}
	}
}
