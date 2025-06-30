package intelligence

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws-research-wizard/go/internal/data"
)

// IntelligenceEngine provides domain-aware intelligent recommendations and optimization
type IntelligenceEngine struct {
	domainProfileManager *data.ResearchDomainProfileManager
	recommendationEngine RecommendationEngineInterface
	domainPackLoader     DomainPackLoaderInterface
	costOptimizer        *CostOptimizer
	resourceAnalyzer     *ResourceAnalyzer
}

// IntelligentRecommendation represents a comprehensive recommendation with domain context
type IntelligentRecommendation struct {
	ID               string                    `json:"id"`
	Timestamp        time.Time                 `json:"timestamp"`
	Domain           string                    `json:"domain"`
	DomainPack       *DomainPackInfo           `json:"domain_pack,omitempty"`
	DataAnalysis     *data.RecommendationResult `json:"data_analysis"`
	ResourcePlan     *ResourcePlan             `json:"resource_plan"`
	CostOptimization *CostOptimizationPlan     `json:"cost_optimization"`
	Implementation   *ImplementationPlan       `json:"implementation"`
	Confidence       float64                   `json:"confidence"`
	Impact           *ImpactAssessment         `json:"impact"`
}

// DomainPackInfo contains information about the recommended domain pack
type DomainPackInfo struct {
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Description    string            `json:"description"`
	Categories     []string          `json:"categories"`
	InstanceTypes  map[string]string `json:"instance_types"`
	Workflows      []WorkflowInfo    `json:"workflows"`
	SpackPackages  []string          `json:"spack_packages"`
	EstimatedCost  map[string]string `json:"estimated_cost"`
}

// WorkflowInfo describes available workflows in a domain pack
type WorkflowInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputData   string `json:"input_data"`
	OutputData  string `json:"output_data"`
}

// ResourcePlan contains infrastructure recommendations
type ResourcePlan struct {
	RecommendedInstance  string                 `json:"recommended_instance"`
	AlternativeInstances []string               `json:"alternative_instances"`
	StorageConfiguration StorageConfiguration   `json:"storage_configuration"`
	NetworkConfiguration NetworkConfiguration   `json:"network_configuration"`
	SecurityConfiguration SecurityConfiguration `json:"security_configuration"`
	Reasoning            string                 `json:"reasoning"`
}

// StorageConfiguration defines storage recommendations
type StorageConfiguration struct {
	PrimaryStorage   StorageType `json:"primary_storage"`
	BackupStorage    StorageType `json:"backup_storage"`
	ArchiveStorage   StorageType `json:"archive_storage"`
	LifecyclePolicies []string   `json:"lifecycle_policies"`
}

// StorageType defines storage type details
type StorageType struct {
	Type       string `json:"type"`
	SizeGB     int    `json:"size_gb"`
	IOPS       int    `json:"iops"`
	Throughput int    `json:"throughput_mbps"`
}

// NetworkConfiguration defines network recommendations
type NetworkConfiguration struct {
	PlacementGroup      bool     `json:"placement_group"`
	EnhancedNetworking  bool     `json:"enhanced_networking"`
	EFAEnabled          bool     `json:"efa_enabled"`
	SecurityGroups      []string `json:"security_groups"`
	VPCConfiguration    string   `json:"vpc_configuration"`
}

// SecurityConfiguration defines security recommendations
type SecurityConfiguration struct {
	EncryptionAtRest    bool     `json:"encryption_at_rest"`
	EncryptionInTransit bool     `json:"encryption_in_transit"`
	IAMRoles            []string `json:"iam_roles"`
	ComplianceFramework string   `json:"compliance_framework"`
	AuditingEnabled     bool     `json:"auditing_enabled"`
}

// CostOptimizationPlan contains cost optimization strategies
type CostOptimizationPlan struct {
	EstimatedMonthlyCost    float64                `json:"estimated_monthly_cost"`
	OptimizedMonthlyCost    float64                `json:"optimized_monthly_cost"`
	PotentialSavings        float64                `json:"potential_savings"`
	SavingsPercentage       float64                `json:"savings_percentage"`
	SpotInstanceSavings     *SpotInstanceSavings   `json:"spot_instance_savings,omitempty"`
	ReservedInstanceSavings *ReservedInstanceSavings `json:"reserved_instance_savings,omitempty"`
	StorageOptimizations    []StorageOptimization  `json:"storage_optimizations"`
	Recommendations         []string               `json:"recommendations"`
}

// SpotInstanceSavings calculates spot instance cost benefits
type SpotInstanceSavings struct {
	PotentialSavingsPercent float64 `json:"potential_savings_percent"`
	EstimatedMonthlySavings float64 `json:"estimated_monthly_savings"`
	RiskAssessment          string  `json:"risk_assessment"`
	RecommendedStrategy     string  `json:"recommended_strategy"`
}

// ReservedInstanceSavings calculates reserved instance benefits
type ReservedInstanceSavings struct {
	OneYearSavings       float64 `json:"one_year_savings"`
	ThreeYearSavings     float64 `json:"three_year_savings"`
	RecommendedTerm      string  `json:"recommended_term"`
	PaymentOption        string  `json:"payment_option"`
	BreakevenPoint       string  `json:"breakeven_point"`
}

// StorageOptimization describes storage cost optimizations
type StorageOptimization struct {
	Type            string  `json:"type"`
	Description     string  `json:"description"`
	SavingsPercent  float64 `json:"savings_percent"`
	MonthlySavings  float64 `json:"monthly_savings"`
	Implementation  string  `json:"implementation"`
}

// ImplementationPlan provides step-by-step implementation guidance
type ImplementationPlan struct {
	EstimatedDuration  string                 `json:"estimated_duration"`
	Complexity         string                 `json:"complexity"`
	Prerequisites      []string               `json:"prerequisites"`
	Steps              []ImplementationStep   `json:"steps"`
	DeploymentScript   string                 `json:"deployment_script"`
	ValidationChecks   []string               `json:"validation_checks"`
	RollbackProcedure  []string               `json:"rollback_procedure"`
}

// ImplementationStep describes a single implementation step
type ImplementationStep struct {
	Order         int      `json:"order"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Commands      []string `json:"commands"`
	Duration      string   `json:"duration"`
	Dependencies  []string `json:"dependencies"`
	SuccessCriteria string `json:"success_criteria"`
}

// ImpactAssessment evaluates the impact of implementing recommendations
type ImpactAssessment struct {
	PerformanceImpact  PerformanceImpact  `json:"performance_impact"`
	CostImpact         CostImpact         `json:"cost_impact"`
	SecurityImpact     SecurityImpact     `json:"security_impact"`
	OperationalImpact  OperationalImpact  `json:"operational_impact"`
	RiskAssessment     RiskAssessment     `json:"risk_assessment"`
}

// PerformanceImpact measures performance improvements
type PerformanceImpact struct {
	ComputePerformance  float64 `json:"compute_performance_improvement_percent"`
	IOPerformance       float64 `json:"io_performance_improvement_percent"`
	NetworkPerformance  float64 `json:"network_performance_improvement_percent"`
	TransferSpeed       string  `json:"estimated_transfer_speed"`
	ProcessingTime      string  `json:"estimated_processing_time"`
}

// CostImpact measures cost implications
type CostImpact struct {
	InitialSetupCost    float64 `json:"initial_setup_cost"`
	OngoingMonthlyCost  float64 `json:"ongoing_monthly_cost"`
	CostSavings         float64 `json:"cost_savings"`
	PaybackPeriod       string  `json:"payback_period"`
	ROI                 float64 `json:"roi_percent"`
}

// SecurityImpact measures security improvements
type SecurityImpact struct {
	SecurityPosture     string   `json:"security_posture_improvement"`
	ComplianceGaps      []string `json:"compliance_gaps_addressed"`
	VulnerabilitiesFixed []string `json:"vulnerabilities_fixed"`
	AuditImprovements   []string `json:"audit_improvements"`
}

// OperationalImpact measures operational effects
type OperationalImpact struct {
	MaintenanceReduction string   `json:"maintenance_reduction"`
	AutomationLevel      string   `json:"automation_level"`
	MonitoringImprovement string  `json:"monitoring_improvement"`
	StaffEfficiency      string   `json:"staff_efficiency_gain"`
	RequiredSkills       []string `json:"required_skills"`
}

// RiskAssessment evaluates implementation risks
type RiskAssessment struct {
	OverallRisk         string   `json:"overall_risk_level"`
	TechnicalRisks      []string `json:"technical_risks"`
	OperationalRisks    []string `json:"operational_risks"`
	SecurityRisks       []string `json:"security_risks"`
	MitigationStrategies []string `json:"mitigation_strategies"`
	ContingencyPlans    []string `json:"contingency_plans"`
}

// NewIntelligenceEngine creates a new intelligence engine
func NewIntelligenceEngine(
	domainProfileManager *data.ResearchDomainProfileManager,
	recommendationEngine RecommendationEngineInterface,
) *IntelligenceEngine {
	return &IntelligenceEngine{
		domainProfileManager: domainProfileManager,
		recommendationEngine: recommendationEngine,
		domainPackLoader:     NewDomainPackLoader(),
		costOptimizer:        NewCostOptimizer(),
		resourceAnalyzer:     NewResourceAnalyzer(),
	}
}

// GenerateIntelligentRecommendations creates comprehensive domain-aware recommendations
func (ie *IntelligenceEngine) GenerateIntelligentRecommendations(
	ctx context.Context,
	dataPath string,
	hints DomainHints,
) (*IntelligentRecommendation, error) {
	
	// Step 1: Detect or validate domain
	detectedDomain, confidence := ie.detectDomain(dataPath, hints)
	
	// Step 2: Load domain pack information
	domainPack, err := ie.domainPackLoader.LoadDomainPack(detectedDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to load domain pack for %s: %w", detectedDomain, err)
	}
	
	// Step 3: Generate data analysis recommendations
	dataRecommendations, err := ie.recommendationEngine.GenerateRecommendations(ctx, dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate data recommendations: %w", err)
	}
	
	// Step 4: Generate resource plan
	resourcePlan := ie.generateResourcePlan(detectedDomain, dataRecommendations, hints)
	
	// Step 5: Generate cost optimization plan
	costPlan := ie.costOptimizer.GenerateCostOptimizationPlan(
		detectedDomain, 
		resourcePlan, 
		dataRecommendations,
	)
	
	// Step 6: Generate implementation plan
	implPlan := ie.generateImplementationPlan(detectedDomain, domainPack, resourcePlan)
	
	// Step 7: Assess overall impact
	impact := ie.assessImpact(resourcePlan, costPlan, dataRecommendations)
	
	recommendation := &IntelligentRecommendation{
		ID:               fmt.Sprintf("intel-%d", time.Now().Unix()),
		Timestamp:        time.Now(),
		Domain:           detectedDomain,
		DomainPack:       domainPack,
		DataAnalysis:     dataRecommendations,
		ResourcePlan:     resourcePlan,
		CostOptimization: costPlan,
		Implementation:   implPlan,
		Confidence:       confidence,
		Impact:           impact,
	}
	
	return recommendation, nil
}

// detectDomain identifies the research domain based on data and hints
func (ie *IntelligenceEngine) detectDomain(dataPath string, hints DomainHints) (string, float64) {
	scores := make(map[string]float64)
	
	// Score based on file extensions
	if extensions := ie.analyzeFileExtensions(dataPath); len(extensions) > 0 {
		for domain, profile := range ie.domainProfileManager.GetAllProfiles() {
			domainScore := 0.0
			totalHints := len(profile.FileTypeHints)
			
			for ext := range extensions {
				if _, exists := profile.FileTypeHints[ext]; exists {
					domainScore += 1.0
				}
			}
			
			if totalHints > 0 {
				scores[domain] = domainScore / float64(totalHints)
			}
		}
	}
	
	// Score based on explicit domain hints
	if hints.ExplicitDomain != "" {
		scores[hints.ExplicitDomain] += 0.8
	}
	
	// Score based on workflow hints
	for _, workflow := range hints.WorkflowHints {
		for domain := range ie.domainProfileManager.GetAllProfiles() {
			if strings.Contains(domain, workflow) || strings.Contains(workflow, domain) {
				scores[domain] += 0.3
			}
		}
	}
	
	// Score based on tool hints
	for _, tool := range hints.ToolHints {
		for domain, profile := range ie.domainProfileManager.GetAllProfiles() {
			// Check transfer engines
			for _, preferredEngine := range profile.TransferOptimization.PreferredEngines {
				if strings.Contains(tool, preferredEngine) || strings.Contains(preferredEngine, tool) {
					scores[domain] += 0.2
				}
			}
			
			// Check domain-specific tools
			domainSpecificTools := ie.getDomainSpecificTools(domain)
			for _, domainTool := range domainSpecificTools {
				if ie.toolsMatch(tool, domainTool) {
					scores[domain] += 0.4 // Higher score for domain-specific tools
				}
			}
		}
	}
	
	// Find the highest scoring domain
	bestDomain := "general"
	bestScore := 0.0
	
	for domain, score := range scores {
		if score > bestScore {
			bestDomain = domain
			bestScore = score
		}
	}
	
	// If no domain scored well, try to detect from common patterns
	if bestScore < 0.3 {
		if detected := ie.detectFromCommonPatterns(dataPath); detected != "" {
			return detected, 0.6
		}
	}
	
	return bestDomain, bestScore
}

// analyzeFileExtensions extracts file extensions from the data path
func (ie *IntelligenceEngine) analyzeFileExtensions(dataPath string) map[string]bool {
	extensions := make(map[string]bool)
	
	// This would typically scan the directory or analyze file listings
	// For now, we'll extract from the path itself if it contains file names
	if strings.Contains(dataPath, ".") {
		ext := filepath.Ext(dataPath)
		if ext != "" {
			extensions[ext] = true
		}
	}
	
	// Common genomics extensions
	genomicsExts := []string{".fastq", ".fq", ".bam", ".sam", ".vcf", ".gff", ".gtf"}
	for _, ext := range genomicsExts {
		if strings.Contains(strings.ToLower(dataPath), ext) {
			extensions[ext] = true
		}
	}
	
	// Common climate/weather extensions
	climateExts := []string{".nc", ".grib", ".hdf", ".netcdf"}
	for _, ext := range climateExts {
		if strings.Contains(strings.ToLower(dataPath), ext) {
			extensions[ext] = true
		}
	}
	
	return extensions
}

// detectFromCommonPatterns detects domain from common naming patterns
func (ie *IntelligenceEngine) detectFromCommonPatterns(dataPath string) string {
	lowerPath := strings.ToLower(dataPath)
	
	// Genomics patterns
	genomicsPatterns := []string{"genome", "sequencing", "fastq", "bam", "variant", "rna-seq", "dna"}
	for _, pattern := range genomicsPatterns {
		if strings.Contains(lowerPath, pattern) {
			return "genomics"
		}
	}
	
	// Machine learning patterns (check specific patterns first)
	mlPatterns := []string{"ml", "machine_learning", "neural", "tensorflow", "pytorch", "training", "dataset"}
	for _, pattern := range mlPatterns {
		if strings.Contains(lowerPath, pattern) {
			return "machine_learning"
		}
	}
	
	// Climate patterns
	climatePatterns := []string{"climate", "weather", "temperature", "precipitation", "forecast"}
	for _, pattern := range climatePatterns {
		if strings.Contains(lowerPath, pattern) {
			return "climate"
		}
	}
	
	return ""
}

// generateResourcePlan creates infrastructure recommendations
func (ie *IntelligenceEngine) generateResourcePlan(
	domain string, 
	dataRec *data.RecommendationResult, 
	hints DomainHints,
) *ResourcePlan {
	
	profile, exists := ie.domainProfileManager.GetProfile(domain)
	if !exists {
		// Use general recommendations
		return ie.generateGeneralResourcePlan(dataRec, hints)
	}
	
	// Generate domain-specific resource plan
	plan := &ResourcePlan{}
	
	// Determine instance type based on workload size
	workloadSize := ie.assessWorkloadSize(dataRec.DataPattern)
	plan.RecommendedInstance = ie.selectOptimalInstance(profile, workloadSize, hints)
	plan.AlternativeInstances = ie.generateAlternativeInstances(profile, workloadSize)
	
	// Configure storage
	plan.StorageConfiguration = ie.generateStorageConfiguration(profile, dataRec.DataPattern)
	
	// Configure networking
	plan.NetworkConfiguration = ie.generateNetworkConfiguration(profile, workloadSize)
	
	// Configure security
	plan.SecurityConfiguration = ie.generateSecurityConfiguration(profile)
	
	// Generate reasoning
	plan.Reasoning = ie.generateResourceReasoning(domain, profile, workloadSize, dataRec)
	
	return plan
}

// DomainHints provides additional context for domain detection
type DomainHints struct {
	ExplicitDomain   string   `json:"explicit_domain,omitempty"`
	WorkflowHints    []string `json:"workflow_hints,omitempty"`
	ToolHints        []string `json:"tool_hints,omitempty"`
	DataSizeHint     string   `json:"data_size_hint,omitempty"`
	PerformanceHints []string `json:"performance_hints,omitempty"`
	BudgetConstraint float64  `json:"budget_constraint,omitempty"`
}

// Additional helper methods would continue here...
// For brevity, I'll implement the core structure and a few key methods

// assessWorkloadSize determines workload characteristics
func (ie *IntelligenceEngine) assessWorkloadSize(pattern *data.DataPattern) string {
	if pattern == nil {
		return "small" // Default fallback for nil pattern
	}
	
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
	
	if totalSizeGB < 10 {
		return "small"
	} else if totalSizeGB < 500 {
		return "medium"
	} else if totalSizeGB < 5000 {
		return "large"
	} else {
		return "massive"
	}
}

// selectOptimalInstance chooses the best instance type for the workload
func (ie *IntelligenceEngine) selectOptimalInstance(
	profile *data.ResearchDomainProfile, 
	workloadSize string, 
	hints DomainHints,
) string {
	
	// Start with domain pack instance types if we have them loaded
	if domainPack, err := ie.domainPackLoader.LoadDomainPack(profile.Name); err == nil && domainPack != nil {
		if instanceType, exists := domainPack.InstanceTypes[workloadSize]; exists {
			return instanceType
		}
	}
	
	// Fallback to profile-based selection
	baseInstance := "c6i.2xlarge" // Default
	
	switch workloadSize {
	case "small":
		baseInstance = "c6i.xlarge"
	case "medium":
		baseInstance = "c6i.4xlarge"
	case "large":
		baseInstance = "c6i.12xlarge"
	case "massive":
		baseInstance = "c6i.24xlarge"
	}
	
	// Adjust based on domain characteristics
	if profile.Name == "machine_learning" {
		// Prefer GPU instances for ML workloads
		switch workloadSize {
		case "small":
			return "g5.xlarge"
		case "medium":
			return "g5.4xlarge"
		case "large":
			return "p4d.24xlarge"
		}
	}
	
	if profile.Name == "genomics" {
		// Prefer memory-optimized instances for genomics
		switch workloadSize {
		case "small":
			return "r6i.2xlarge"
		case "medium":
			return "r6i.4xlarge"
		case "large":
			return "r6i.8xlarge"
		case "massive":
			return "r6i.16xlarge"
		}
	}
	
	return baseInstance
}

// generateImplementationPlan creates a step-by-step implementation plan
func (ie *IntelligenceEngine) generateImplementationPlan(
	domain string,
	domainPack *DomainPackInfo,
	resourcePlan *ResourcePlan,
) *ImplementationPlan {
	
	steps := []ImplementationStep{
		{
			Order:       1,
			Title:       "Configure AWS Environment",
			Description: "Set up AWS credentials and verify access",
			Commands: []string{
				"aws configure",
				"aws sts get-caller-identity",
			},
			Duration:        "15 minutes",
			Dependencies:    []string{},
			SuccessCriteria: "AWS CLI configured and working",
		},
		{
			Order:       2,
			Title:       "Deploy Domain Pack Infrastructure",
			Description: "Deploy the recommended AWS infrastructure",
			Commands: []string{
				fmt.Sprintf("aws-research-wizard deploy --domain %s --instance %s", 
					domain, resourcePlan.RecommendedInstance),
			},
			Duration:        "20-30 minutes",
			Dependencies:    []string{"AWS Environment"},
			SuccessCriteria: "Infrastructure deployed and accessible",
		},
		{
			Order:       3,
			Title:       "Install Domain Pack Software",
			Description: "Install and configure domain-specific software stack",
			Commands: []string{
				"aws-research-wizard config install-domain-pack",
				"spack env activate research-env",
			},
			Duration:        "30-60 minutes",
			Dependencies:    []string{"Infrastructure"},
			SuccessCriteria: "All software packages installed and configured",
		},
		{
			Order:       4,
			Title:       "Validate Setup",
			Description: "Run validation tests to ensure everything is working",
			Commands: []string{
				"aws-research-wizard validate --domain " + domain,
			},
			Duration:        "10 minutes",
			Dependencies:    []string{"Software Installation"},
			SuccessCriteria: "All validation tests pass",
		},
	}
	
	complexity := "moderate"
	if domain == "machine_learning" || domain == "climate" {
		complexity = "complex"
	}
	
	return &ImplementationPlan{
		EstimatedDuration: "1.5-2.5 hours",
		Complexity:        complexity,
		Prerequisites: []string{
			"AWS account with appropriate permissions",
			"AWS CLI installed",
			"Basic knowledge of " + domain + " workflows",
		},
		Steps:            steps,
		DeploymentScript: ie.generateDeploymentScript(domain, resourcePlan),
		ValidationChecks: []string{
			"Instance accessibility",
			"Software installation status",
			"Domain pack configuration",
			"Network connectivity",
		},
		RollbackProcedure: []string{
			"Stop all running instances",
			"Delete CloudFormation stack",
			"Clean up S3 resources",
		},
	}
}

// assessImpact evaluates the overall impact of implementing recommendations
func (ie *IntelligenceEngine) assessImpact(
	resourcePlan *ResourcePlan,
	costPlan *CostOptimizationPlan,
	dataRecommendations *data.RecommendationResult,
) *ImpactAssessment {
	
	return &ImpactAssessment{
		PerformanceImpact: PerformanceImpact{
			ComputePerformance: 150.0, // 50% improvement over baseline
			IOPerformance:      200.0, // 100% improvement with optimized storage
			NetworkPerformance: 120.0, // 20% improvement with enhanced networking
			TransferSpeed:      "5-8 GB/s",
			ProcessingTime:     "Reduced by 40-60%",
		},
		CostImpact: CostImpact{
			InitialSetupCost:   100.0, // Setup costs
			OngoingMonthlyCost: costPlan.OptimizedMonthlyCost,
			CostSavings:        costPlan.PotentialSavings,
			PaybackPeriod:      "2-3 months",
			ROI:                (costPlan.PotentialSavings * 12 / 100.0) * 100, // Annual ROI
		},
		SecurityImpact: SecurityImpact{
			SecurityPosture:     "Significantly improved",
			ComplianceGaps:      []string{"Automated encryption", "Access logging"},
			VulnerabilitiesFixed: []string{"Unencrypted data transfer", "Default security groups"},
			AuditImprovements:   []string{"CloudTrail logging", "Resource tagging"},
		},
		OperationalImpact: OperationalImpact{
			MaintenanceReduction: "60% reduction in manual tasks",
			AutomationLevel:      "High - automated provisioning and scaling",
			MonitoringImprovement: "Real-time dashboards and alerting",
			StaffEfficiency:      "40% productivity improvement",
			RequiredSkills:       []string{"AWS basics", "Domain expertise", "Command line"},
		},
		RiskAssessment: RiskAssessment{
			OverallRisk:     "Low",
			TechnicalRisks:  []string{"Learning curve", "Configuration complexity"},
			OperationalRisks: []string{"Staff training", "Process changes"},
			SecurityRisks:   []string{"Minimal - enhanced security"},
			MitigationStrategies: []string{
				"Comprehensive documentation",
				"Gradual rollout",
				"Training programs",
			},
			ContingencyPlans: []string{
				"Rollback procedures",
				"Manual fallback options",
				"Expert support available",
			},
		},
	}
}

// generateGeneralResourcePlan creates a general resource plan when no domain profile exists
func (ie *IntelligenceEngine) generateGeneralResourcePlan(
	dataRec *data.RecommendationResult,
	hints DomainHints,
) *ResourcePlan {
	
	workloadSize := ie.assessWorkloadSize(dataRec.DataPattern)
	
	// General instance selection
	instanceType := "c6i.2xlarge"
	switch workloadSize {
	case "small":
		instanceType = "c6i.xlarge"
	case "medium":
		instanceType = "c6i.4xlarge"
	case "large":
		instanceType = "c6i.8xlarge"
	case "massive":
		instanceType = "c6i.16xlarge"
	}
	
	return &ResourcePlan{
		RecommendedInstance:  instanceType,
		AlternativeInstances: []string{"c6a.4xlarge", "r6i.4xlarge"},
		StorageConfiguration: StorageConfiguration{
			PrimaryStorage: StorageType{
				Type:       "gp3",
				SizeGB:     500,
				IOPS:       3000,
				Throughput: 125,
			},
		},
		NetworkConfiguration: NetworkConfiguration{
			PlacementGroup:     false,
			EnhancedNetworking: true,
			EFAEnabled:         false,
		},
		SecurityConfiguration: SecurityConfiguration{
			EncryptionAtRest:    true,
			EncryptionInTransit: true,
			IAMRoles:           []string{"ResearchRole"},
		},
		Reasoning: "General-purpose configuration suitable for most research workloads",
	}
}

// generateAlternativeInstances suggests alternative instance types
func (ie *IntelligenceEngine) generateAlternativeInstances(
	profile *data.ResearchDomainProfile,
	workloadSize string,
) []string {
	
	var alternatives []string
	
	// Add cost-optimized alternatives
	switch workloadSize {
	case "small":
		alternatives = append(alternatives, "c6a.xlarge", "t3.2xlarge")
	case "medium":
		alternatives = append(alternatives, "c6a.4xlarge", "r6i.2xlarge")
	case "large":
		alternatives = append(alternatives, "c6a.8xlarge", "r6i.8xlarge")
	case "massive":
		alternatives = append(alternatives, "c6a.24xlarge", "r6i.16xlarge")
	}
	
	// Add domain-specific alternatives
	if profile.Name == "machine_learning" {
		alternatives = append(alternatives, "g5.4xlarge", "p4d.24xlarge")
	}
	
	return alternatives
}

// generateStorageConfiguration creates storage configuration recommendations
func (ie *IntelligenceEngine) generateStorageConfiguration(
	profile *data.ResearchDomainProfile,
	dataPattern *data.DataPattern,
) StorageConfiguration {
	
	baseSize := 500
	if dataPattern != nil {
		dataSizeGB := int(dataPattern.TotalSize / (1024 * 1024 * 1024))
		if dataSizeGB > baseSize {
			baseSize = dataSizeGB + 200 // Add 200GB buffer
		}
	}
	
	storageType := "gp3"
	iops := 3000
	throughput := 125
	
	// Adjust based on domain requirements
	if profile != nil {
		switch profile.Name {
		case "genomics":
			iops = 8000
			throughput = 500
		case "astronomy":
			throughput = 1000
			baseSize = baseSize * 2 // Larger storage for image files
		case "machine_learning":
			iops = 16000
			throughput = 1000
		}
	}
	
	return StorageConfiguration{
		PrimaryStorage: StorageType{
			Type:       storageType,
			SizeGB:     baseSize,
			IOPS:       iops,
			Throughput: throughput,
		},
		BackupStorage: StorageType{
			Type:   "s3_standard_ia",
			SizeGB: baseSize / 2,
		},
		ArchiveStorage: StorageType{
			Type:   "s3_glacier",
			SizeGB: baseSize,
		},
		LifecyclePolicies: []string{
			"Transition to IA after 30 days",
			"Archive to Glacier after 90 days",
		},
	}
}

// generateNetworkConfiguration creates network configuration recommendations
func (ie *IntelligenceEngine) generateNetworkConfiguration(
	profile *data.ResearchDomainProfile,
	workloadSize string,
) NetworkConfiguration {
	
	config := NetworkConfiguration{
		PlacementGroup:     false,
		EnhancedNetworking: true,
		EFAEnabled:         false,
		SecurityGroups:     []string{"research-sg"},
		VPCConfiguration:   "default",
	}
	
	// Adjust based on workload size and domain
	if workloadSize == "large" || workloadSize == "massive" {
		config.PlacementGroup = true
	}
	
	if profile != nil {
		switch profile.Name {
		case "climate", "machine_learning":
			config.EFAEnabled = true
			config.PlacementGroup = true
		}
	}
	
	return config
}

// generateSecurityConfiguration creates security configuration recommendations
func (ie *IntelligenceEngine) generateSecurityConfiguration(
	profile *data.ResearchDomainProfile,
) SecurityConfiguration {
	
	config := SecurityConfiguration{
		EncryptionAtRest:    true,
		EncryptionInTransit: true,
		IAMRoles:           []string{"ResearchRole"},
		AuditingEnabled:    true,
	}
	
	// Adjust based on domain requirements
	if profile != nil {
		if profile.SecurityRequirements.EncryptionRequired {
			config.ComplianceFramework = "Enhanced"
		}
		
		if len(profile.ComplianceSettings.Framework) > 0 {
			config.ComplianceFramework = profile.ComplianceSettings.Framework
		}
	}
	
	return config
}

// generateResourceReasoning creates reasoning explanation for resource recommendations
func (ie *IntelligenceEngine) generateResourceReasoning(
	domain string,
	profile *data.ResearchDomainProfile,
	workloadSize string,
	dataRec *data.RecommendationResult,
) string {
	
	var reasons []string
	
	reasons = append(reasons, fmt.Sprintf("Domain: %s workloads require specific optimizations", domain))
	reasons = append(reasons, fmt.Sprintf("Workload size: %s requires appropriate resource allocation", workloadSize))
	
	if dataRec.DataPattern != nil {
		totalSizeGB := float64(dataRec.DataPattern.TotalSize) / (1024 * 1024 * 1024)
		reasons = append(reasons, fmt.Sprintf("Dataset size: %.1f GB influences storage and memory requirements", totalSizeGB))
	}
	
	if profile != nil {
		reasons = append(reasons, "Domain profile optimizations applied for performance and cost efficiency")
	}
	
	return strings.Join(reasons, "; ")
}

// toolsMatch determines if two tool names match using intelligent comparison
func (ie *IntelligenceEngine) toolsMatch(tool1, tool2 string) bool {
	t1 := strings.ToLower(strings.TrimSpace(tool1))
	t2 := strings.ToLower(strings.TrimSpace(tool2))
	
	// Exact match
	if t1 == t2 {
		return true
	}
	
	// Prevent single character matches
	if len(t1) <= 2 || len(t2) <= 2 {
		return false
	}
	
	// One contains the other (but both must be substantial length)
	if len(t1) >= 3 && len(t2) >= 3 {
		if strings.Contains(t1, t2) || strings.Contains(t2, t1) {
			return true
		}
	}
	
	return false
}

// getDomainSpecificTools returns domain-specific tools and software
func (ie *IntelligenceEngine) getDomainSpecificTools(domain string) []string {
	switch domain {
	case "machine_learning":
		return []string{"pytorch", "tensorflow", "keras", "scikit-learn", "pandas", "numpy", "jupyter", "conda", "pip"}
	case "genomics":
		return []string{"blast", "bwa", "samtools", "bcftools", "gatk", "star", "hisat2", "bowtie2", "trimmomatic"}
	case "climate":
		return []string{"ncl", "cdo", "nco", "ferret", "grads", "matlab", "r", "python", "netcdf", "hdf5"}
	case "astronomy":
		return []string{"casa", "miriad", "aips", "ds9", "saoimage", "astropy", "aplpy", "montage"}
	case "chemistry":
		return []string{"gaussian", "orca", "gamess", "molpro", "amber", "gromacs", "namd", "lammps"}
	case "physics":
		return []string{"mathematica", "matlab", "labview", "root", "geant4", "openfoam", "ansys"}
	default:
		return []string{}
	}
}

// generateDeploymentScript creates a deployment script for the recommendations
func (ie *IntelligenceEngine) generateDeploymentScript(
	domain string,
	resourcePlan *ResourcePlan,
) string {
	
	script := fmt.Sprintf(`#!/bin/bash
# AWS Research Wizard Deployment Script
# Domain: %s
# Generated: %s

set -e

echo "Deploying AWS Research Wizard for %s domain..."

# Step 1: Verify AWS configuration
echo "Verifying AWS configuration..."
aws sts get-caller-identity || { echo "AWS CLI not configured"; exit 1; }

# Step 2: Deploy infrastructure
echo "Deploying infrastructure..."
aws-research-wizard deploy \
  --domain %s \
  --instance-type %s \
  --storage-size %d \
  --storage-type %s

# Step 3: Wait for deployment
echo "Waiting for deployment to complete..."
aws-research-wizard wait deployment-complete

# Step 4: Validate deployment
echo "Validating deployment..."
aws-research-wizard validate --domain %s

echo "Deployment complete! Access your research environment:"
aws-research-wizard info --domain %s
`,
		domain,
		time.Now().Format("2006-01-02 15:04:05"),
		domain,
		domain,
		resourcePlan.RecommendedInstance,
		resourcePlan.StorageConfiguration.PrimaryStorage.SizeGB,
		resourcePlan.StorageConfiguration.PrimaryStorage.Type,
		domain,
		domain,
	)
	
	return script
}