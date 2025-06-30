package intelligence

import (
	"fmt"
	"strings"

	"github.com/aws-research-wizard/go/internal/data"
)

// CostOptimizer provides intelligent cost optimization recommendations
type CostOptimizer struct {
	// AWS pricing data would typically be loaded from AWS APIs
	// For now, we'll use static pricing data
	instancePricing map[string]float64
	storagePricing  map[string]float64
}

// AWSInstancePricing represents instance pricing information
type AWSInstancePricing struct {
	InstanceType    string  `json:"instance_type"`
	HourlyPrice     float64 `json:"hourly_price"`
	MonthlyPrice    float64 `json:"monthly_price"`
	vCPUs           int     `json:"vcpus"`
	MemoryGB        float64 `json:"memory_gb"`
	NetworkSpeed    string  `json:"network_speed"`
	StorageType     string  `json:"storage_type"`
}

// NewCostOptimizer creates a new cost optimizer
func NewCostOptimizer() *CostOptimizer {
	co := &CostOptimizer{
		instancePricing: make(map[string]float64),
		storagePricing:  make(map[string]float64),
	}
	
	co.initializePricingData()
	return co
}

// initializePricingData sets up static pricing data
// In production, this would fetch real-time pricing from AWS APIs
func (co *CostOptimizer) initializePricingData() {
	// Instance pricing (hourly rates in USD)
	co.instancePricing = map[string]float64{
		// General purpose
		"t3.micro":     0.0104,
		"t3.small":     0.0208,
		"t3.medium":    0.0416,
		"t3.large":     0.0832,
		"t3.xlarge":    0.1664,
		"t3.2xlarge":   0.3328,
		
		// Compute optimized
		"c6i.large":    0.0850,
		"c6i.xlarge":   0.1700,
		"c6i.2xlarge":  0.3400,
		"c6i.4xlarge":  0.6800,
		"c6i.8xlarge":  1.3600,
		"c6i.12xlarge": 2.0400,
		"c6i.16xlarge": 2.7200,
		"c6i.24xlarge": 4.0800,
		"c6i.32xlarge": 5.4400,
		
		// Memory optimized
		"r6i.large":    0.1260,
		"r6i.xlarge":   0.2520,
		"r6i.2xlarge":  0.5040,
		"r6i.4xlarge":  1.0080,
		"r6i.8xlarge":  2.0160,
		"r6i.12xlarge": 3.0240,
		"r6i.16xlarge": 4.0320,
		"r6i.24xlarge": 6.0480,
		"r6i.32xlarge": 8.0640,
		
		// GPU instances
		"g5.xlarge":    1.0060,
		"g5.2xlarge":   1.2120,
		"g5.4xlarge":   1.6240,
		"g5.8xlarge":   2.4480,
		"g5.12xlarge":  4.8960,
		"g5.16xlarge":  4.8960,
		"g5.24xlarge":  9.7920,
		"g5.48xlarge":  19.5840,
		"p4d.24xlarge": 32.7726,
		
		// AMD instances
		"c6a.large":    0.0765,
		"c6a.xlarge":   0.1530,
		"c6a.2xlarge":  0.3060,
		"c6a.4xlarge":  0.6120,
		"c6a.8xlarge":  1.2240,
		"c6a.12xlarge": 1.8360,
		"c6a.16xlarge": 2.4480,
		"c6a.24xlarge": 3.6720,
		"c6a.32xlarge": 4.8960,
		"c6a.48xlarge": 7.3440,
	}
	
	// Storage pricing (per GB per month)
	co.storagePricing = map[string]float64{
		"gp3":              0.080,
		"gp2":              0.100,
		"io1":              0.125,
		"io2":              0.125,
		"st1":              0.045,
		"sc1":              0.025,
		"s3_standard":      0.023,
		"s3_standard_ia":   0.0125,
		"s3_glacier":       0.004,
		"s3_deep_archive":  0.00099,
	}
}

// GenerateCostOptimizationPlan creates a comprehensive cost optimization plan
func (co *CostOptimizer) GenerateCostOptimizationPlan(
	domain string,
	resourcePlan *ResourcePlan,
	dataRecommendations *data.RecommendationResult,
) *CostOptimizationPlan {
	
	// Calculate base costs
	baseMonthlyCost := co.calculateBaseMonthlyCost(resourcePlan)
	
	// Generate optimization strategies
	spotSavings := co.calculateSpotInstanceSavings(resourcePlan.RecommendedInstance)
	reservedSavings := co.calculateReservedInstanceSavings(resourcePlan.RecommendedInstance, baseMonthlyCost)
	storageOptimizations := co.generateStorageOptimizations(resourcePlan, dataRecommendations)
	
	// Calculate optimized cost
	optimizedCost := baseMonthlyCost
	totalSavings := 0.0
	
	// Apply spot instance savings (conservative estimate)
	if spotSavings != nil {
		potentialSpotSavings := baseMonthlyCost * (spotSavings.PotentialSavingsPercent / 100.0) * 0.7 // 70% uptime assumption
		optimizedCost -= potentialSpotSavings
		totalSavings += potentialSpotSavings
	}
	
	// Apply storage optimizations
	for _, opt := range storageOptimizations {
		optimizedCost -= opt.MonthlySavings
		totalSavings += opt.MonthlySavings
	}
	
	savingsPercentage := (totalSavings / baseMonthlyCost) * 100
	
	recommendations := co.generateCostRecommendations(domain, resourcePlan, spotSavings, reservedSavings, storageOptimizations)
	
	return &CostOptimizationPlan{
		EstimatedMonthlyCost:    baseMonthlyCost,
		OptimizedMonthlyCost:    optimizedCost,
		PotentialSavings:        totalSavings,
		SavingsPercentage:       savingsPercentage,
		SpotInstanceSavings:     spotSavings,
		ReservedInstanceSavings: reservedSavings,
		StorageOptimizations:    storageOptimizations,
		Recommendations:         recommendations,
	}
}

// calculateBaseMonthlyCost calculates the base monthly cost for the resource plan
func (co *CostOptimizer) calculateBaseMonthlyCost(resourcePlan *ResourcePlan) float64 {
	// Instance cost
	instanceHourlyRate, exists := co.instancePricing[resourcePlan.RecommendedInstance]
	if !exists {
		instanceHourlyRate = 1.0 // Default fallback
	}
	
	instanceMonthlyCost := instanceHourlyRate * 24 * 30 // Assume 30 days, 24/7 usage
	
	// Storage cost
	storageCost := 0.0
	if resourcePlan.StorageConfiguration.PrimaryStorage.SizeGB > 0 {
		storageType := resourcePlan.StorageConfiguration.PrimaryStorage.Type
		storageRate, exists := co.storagePricing[storageType]
		if exists {
			storageCost += float64(resourcePlan.StorageConfiguration.PrimaryStorage.SizeGB) * storageRate
		}
	}
	
	// Additional storage (backup, archive)
	if resourcePlan.StorageConfiguration.BackupStorage.SizeGB > 0 {
		backupRate, exists := co.storagePricing[resourcePlan.StorageConfiguration.BackupStorage.Type]
		if exists {
			storageCost += float64(resourcePlan.StorageConfiguration.BackupStorage.SizeGB) * backupRate
		}
	}
	
	return instanceMonthlyCost + storageCost
}

// calculateSpotInstanceSavings calculates potential spot instance savings
func (co *CostOptimizer) calculateSpotInstanceSavings(instanceType string) *SpotInstanceSavings {
	onDemandPrice, exists := co.instancePricing[instanceType]
	if !exists {
		return nil
	}
	
	// Spot instances typically provide 60-90% savings, but with interruption risk
	spotSavingsPercent := 70.0 // Conservative estimate
	monthlyOnDemandCost := onDemandPrice * 24 * 30
	monthlySavings := monthlyOnDemandCost * (spotSavingsPercent / 100.0)
	
	// Risk assessment based on instance type
	riskLevel := "medium"
	strategy := "mixed"
	
	if isComputeOptimized(instanceType) {
		riskLevel = "low"
		strategy = "spot_primary"
		spotSavingsPercent = 75.0
	} else if isGPUInstance(instanceType) {
		riskLevel = "high"
		strategy = "on_demand_primary"
		spotSavingsPercent = 60.0
	}
	
	return &SpotInstanceSavings{
		PotentialSavingsPercent: spotSavingsPercent,
		EstimatedMonthlySavings: monthlySavings,
		RiskAssessment:          riskLevel,
		RecommendedStrategy:     strategy,
	}
}

// calculateReservedInstanceSavings calculates reserved instance savings
func (co *CostOptimizer) calculateReservedInstanceSavings(instanceType string, monthlyOnDemandCost float64) *ReservedInstanceSavings {
	// Reserved instances typically provide 30-60% savings depending on term
	oneYearSavings := monthlyOnDemandCost * 12 * 0.35  // 35% savings for 1-year
	threeYearSavings := monthlyOnDemandCost * 36 * 0.55 // 55% savings for 3-year
	
	// Breakeven analysis
	breakevenPoint := "3-4 months"
	if monthlyOnDemandCost > 1000 {
		breakevenPoint = "2-3 months"
	}
	
	return &ReservedInstanceSavings{
		OneYearSavings:   oneYearSavings,
		ThreeYearSavings: threeYearSavings,
		RecommendedTerm:  "1-year",
		PaymentOption:    "partial_upfront",
		BreakevenPoint:   breakevenPoint,
	}
}

// generateStorageOptimizations creates storage optimization recommendations
func (co *CostOptimizer) generateStorageOptimizations(
	resourcePlan *ResourcePlan,
	dataRecommendations *data.RecommendationResult,
) []StorageOptimization {
	
	var optimizations []StorageOptimization
	
	// S3 Intelligent Tiering optimization
	if dataRecommendations.DataPattern.TotalSize > 1024*1024*1024*100 { // > 100GB
		standardCost := float64(dataRecommendations.DataPattern.TotalSize) / (1024*1024*1024) * co.storagePricing["s3_standard"]
		intelligentTieringCost := standardCost * 0.68 // Typically 32% savings
		savings := standardCost - intelligentTieringCost
		
		optimizations = append(optimizations, StorageOptimization{
			Type:           "intelligent_tiering",
			Description:    "Enable S3 Intelligent Tiering for automatic cost optimization",
			SavingsPercent: 32.0,
			MonthlySavings: savings,
			Implementation: "Configure lifecycle policy with intelligent tiering",
		})
	}
	
	// Lifecycle policy optimization
	if dataRecommendations.DataPattern.AccessPatterns.LikelyArchival {
		standardCost := float64(dataRecommendations.DataPattern.TotalSize) / (1024*1024*1024) * co.storagePricing["s3_standard"]
		glacierCost := float64(dataRecommendations.DataPattern.TotalSize) / (1024*1024*1024) * co.storagePricing["s3_glacier"]
		savings := standardCost - glacierCost
		
		optimizations = append(optimizations, StorageOptimization{
			Type:           "glacier_archival",
			Description:    "Archive infrequently accessed data to Glacier",
			SavingsPercent: 82.6,
			MonthlySavings: savings,
			Implementation: "Create lifecycle policy to transition to Glacier after 90 days",
		})
	}
	
	// EBS optimization
	if resourcePlan.StorageConfiguration.PrimaryStorage.Type == "gp2" {
		gp2Cost := float64(resourcePlan.StorageConfiguration.PrimaryStorage.SizeGB) * co.storagePricing["gp2"]
		gp3Cost := float64(resourcePlan.StorageConfiguration.PrimaryStorage.SizeGB) * co.storagePricing["gp3"]
		savings := gp2Cost - gp3Cost
		
		if savings > 0 {
			optimizations = append(optimizations, StorageOptimization{
				Type:           "ebs_gp3_migration",
				Description:    "Migrate from gp2 to gp3 for better performance and cost",
				SavingsPercent: 20.0,
				MonthlySavings: savings,
				Implementation: "Convert EBS volumes to gp3 type",
			})
		}
	}
	
	return optimizations
}

// generateCostRecommendations creates actionable cost recommendations
func (co *CostOptimizer) generateCostRecommendations(
	domain string,
	resourcePlan *ResourcePlan,
	spotSavings *SpotInstanceSavings,
	reservedSavings *ReservedInstanceSavings,
	storageOpts []StorageOptimization,
) []string {
	
	var recommendations []string
	
	// Instance recommendations
	if spotSavings != nil && spotSavings.PotentialSavingsPercent > 60 {
		recommendations = append(recommendations, 
			fmt.Sprintf("Consider using Spot Instances for %.0f%% cost savings (risk: %s)", 
				spotSavings.PotentialSavingsPercent, spotSavings.RiskAssessment))
	}
	
	if reservedSavings != nil {
		recommendations = append(recommendations,
			fmt.Sprintf("Reserved Instances can save $%.0f annually with %s commitment",
				reservedSavings.OneYearSavings, reservedSavings.RecommendedTerm))
	}
	
	// Storage recommendations
	for _, opt := range storageOpts {
		if opt.MonthlySavings > 10 { // Only recommend if savings > $10/month
			recommendations = append(recommendations,
				fmt.Sprintf("%s: Save $%.0f/month (%.0f%% reduction)",
					opt.Description, opt.MonthlySavings, opt.SavingsPercent))
		}
	}
	
	// Domain-specific recommendations
	switch domain {
	case "genomics":
		recommendations = append(recommendations,
			"Use S3 Intelligent Tiering for genomics data - typically 40-60% storage cost reduction",
			"Consider Glacier for long-term storage of raw sequencing data after analysis")
	case "machine_learning":
		recommendations = append(recommendations,
			"Use Spot Instances for training workloads - ML jobs are typically fault-tolerant",
			"Consider S3 Standard-IA for model artifacts and training datasets")
	case "climate":
		recommendations = append(recommendations,
			"Archive processed climate model outputs to Glacier after 90 days",
			"Use HPC instances with Spot for batch climate simulations")
	}
	
	// Budget-based recommendations
	baseCost := co.calculateBaseMonthlyCost(resourcePlan)
	if baseCost > 1000 {
		recommendations = append(recommendations,
			"Set up AWS Budgets with alerts for cost control",
			"Enable Cost and Usage Reports for detailed cost analysis")
	}
	
	return recommendations
}

// Helper functions

func isComputeOptimized(instanceType string) bool {
	return strings.HasPrefix(instanceType, "c5") || 
		   strings.HasPrefix(instanceType, "c6i") || 
		   strings.HasPrefix(instanceType, "c6a")
}

func isGPUInstance(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p3") || 
		   strings.HasPrefix(instanceType, "p4d") || 
		   strings.HasPrefix(instanceType, "g4") || 
		   strings.HasPrefix(instanceType, "g5")
}

func isMemoryOptimized(instanceType string) bool {
	return strings.HasPrefix(instanceType, "r5") || 
		   strings.HasPrefix(instanceType, "r6i") || 
		   strings.HasPrefix(instanceType, "x1") ||
		   strings.HasPrefix(instanceType, "z1d")
}

// EstimateMonthlyCost provides a quick cost estimation for an instance type
func (co *CostOptimizer) EstimateMonthlyCost(instanceType string, hoursPerMonth float64) float64 {
	hourlyRate, exists := co.instancePricing[instanceType]
	if !exists {
		return 0.0
	}
	
	return hourlyRate * hoursPerMonth
}

// CompareInstanceCosts compares costs between different instance types
func (co *CostOptimizer) CompareInstanceCosts(instanceTypes []string) map[string]float64 {
	costs := make(map[string]float64)
	
	for _, instanceType := range instanceTypes {
		if hourlyRate, exists := co.instancePricing[instanceType]; exists {
			costs[instanceType] = hourlyRate * 24 * 30 // Monthly cost
		}
	}
	
	return costs
}