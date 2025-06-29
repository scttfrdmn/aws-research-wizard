package data

import (
	"context"
	"fmt"
	"math"
	"time"
)

// S3CostCalculator provides detailed cost analysis for S3 operations and storage
type S3CostCalculator struct {
	region       string
	pricingModel *S3PricingModel
}

// S3PricingModel contains pricing information for different S3 services and regions
type S3PricingModel struct {
	Region           string                    `json:"region"`
	StorageClasses   map[string]StoragePricing `json:"storage_classes"`
	RequestPricing   RequestPricing            `json:"request_pricing"`
	TransferPricing  TransferPricing           `json:"transfer_pricing"`
	LifecyclePricing LifecyclePricing          `json:"lifecycle_pricing"`
	LastUpdated      time.Time                 `json:"last_updated"`
}

// StoragePricing represents pricing for a specific storage class
type StoragePricing struct {
	Name                string  `json:"name"`
	PricePerGBMonth     float64 `json:"price_per_gb_month"`
	MinimumStorageDays  int     `json:"minimum_storage_days"`
	MinimumObjectSize   int64   `json:"minimum_object_size"`
	RetrievalFeePerGB   float64 `json:"retrieval_fee_per_gb"`
	FirstTierGB         int64   `json:"first_tier_gb"`         // First pricing tier threshold
	SecondTierGB        int64   `json:"second_tier_gb"`        // Second pricing tier threshold
	FirstTierPrice      float64 `json:"first_tier_price"`      // Price for first tier
	SecondTierPrice     float64 `json:"second_tier_price"`     // Price for second tier
	ThirdTierPrice      float64 `json:"third_tier_price"`      // Price for third tier (beyond second tier)
}

// RequestPricing represents pricing for different types of requests
type RequestPricing struct {
	PutCopyPostList    float64 `json:"put_copy_post_list_per_1000"`    // PUT, COPY, POST, LIST requests
	Get                float64 `json:"get_per_1000"`                   // GET and all other requests
	LifecycleTransition float64 `json:"lifecycle_transition_per_1000"`  // Lifecycle transition requests
	Select             float64 `json:"select_per_million"`             // S3 Select requests
	SelectDataScanned  float64 `json:"select_data_scanned_per_gb"`     // S3 Select data scanned
	SelectDataReturned float64 `json:"select_data_returned_per_gb"`    // S3 Select data returned
}

// TransferPricing represents data transfer pricing
type TransferPricing struct {
	InboundFree        bool    `json:"inbound_free"`
	OutboundFirstGB    float64 `json:"outbound_first_gb_price"`    // First 1 GB per month
	OutboundNext9GB    float64 `json:"outbound_next_9gb_price"`    // Next 9.999 TB per month
	OutboundNext40GB   float64 `json:"outbound_next_40gb_price"`   // Next 40 TB per month
	OutboundNext100GB  float64 `json:"outbound_next_100gb_price"`  // Next 100 TB per month
	OutboundBeyond150GB float64 `json:"outbound_beyond_150gb_price"` // Beyond 150 TB per month
	CrossRegionPer     float64 `json:"cross_region_per_gb"`        // Cross-region transfer
	CloudFrontPer      float64 `json:"cloudfront_per_gb"`          // Transfer to CloudFront
}

// LifecyclePricing represents lifecycle transition costs
type LifecyclePricing struct {
	ToIA            float64 `json:"to_ia_per_1000_objects"`
	ToGlacier       float64 `json:"to_glacier_per_1000_objects"`
	ToGlacierIR     float64 `json:"to_glacier_ir_per_1000_objects"`
	ToGlacierDA     float64 `json:"to_glacier_da_per_1000_objects"`
	ToIntelligent   float64 `json:"to_intelligent_per_1000_objects"`
}

// CostAnalysis represents a complete cost analysis for a dataset
type CostAnalysis struct {
	AnalysisID       string                 `json:"analysis_id"`
	DataPattern      *DataPattern          `json:"data_pattern"`
	Scenarios        []CostScenario        `json:"scenarios"`
	Recommendations  []CostRecommendation  `json:"recommendations"`
	Optimizations    []OptimizationOption  `json:"optimizations"`
	TotalCostRange   CostRange            `json:"total_cost_range"`
	PotentialSavings float64              `json:"potential_savings_monthly"`
	AnalysisTime     time.Time            `json:"analysis_time"`
}

// CostScenario represents a specific cost scenario (e.g., current state, optimized state)
type CostScenario struct {
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	StorageClass    string             `json:"storage_class"`
	Configuration   ScenarioConfig     `json:"configuration"`
	MonthlyCosts    DetailedCosts      `json:"monthly_costs"`
	YearlyCosts     DetailedCosts      `json:"yearly_costs"`
	CostBreakdown   map[string]float64 `json:"cost_breakdown"`
	Assumptions     []string           `json:"assumptions"`
}

// ScenarioConfig represents the configuration for a cost scenario
type ScenarioConfig struct {
	FileCount           int64   `json:"file_count"`
	TotalSizeGB         float64 `json:"total_size_gb"`
	StorageClass        string  `json:"storage_class"`
	BundleSmallFiles    bool    `json:"bundle_small_files"`
	CompressionEnabled  bool    `json:"compression_enabled"`
	CompressionRatio    float64 `json:"compression_ratio"`
	AccessFrequency     string  `json:"access_frequency"` // "daily", "weekly", "monthly", "yearly", "rarely"
	DownloadPercentage  float64 `json:"download_percentage"`
	LifecyclePolicyDays int     `json:"lifecycle_policy_days"`
}

// DetailedCosts represents detailed cost breakdown
type DetailedCosts struct {
	Storage         float64 `json:"storage"`
	Requests        float64 `json:"requests"`
	DataTransfer    float64 `json:"data_transfer"`
	Lifecycle       float64 `json:"lifecycle"`
	Retrieval       float64 `json:"retrieval"`
	Monitoring      float64 `json:"monitoring"`
	Total           float64 `json:"total"`
}

// CostRange represents a range of costs (min/max estimates)
type CostRange struct {
	MinMonthly float64 `json:"min_monthly"`
	MaxMonthly float64 `json:"max_monthly"`
	MinYearly  float64 `json:"min_yearly"`
	MaxYearly  float64 `json:"max_yearly"`
}

// CostRecommendation represents a specific cost optimization recommendation
type CostRecommendation struct {
	Type            string             `json:"type"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	EstimatedSavings float64           `json:"estimated_savings_monthly"`
	Confidence      float64           `json:"confidence"`      // 0-1 scale
	Complexity      string            `json:"complexity"`      // "low", "medium", "high"
	Implementation  string            `json:"implementation"`  // How to implement
	Tradeoffs       []string          `json:"tradeoffs"`       // What the user gives up
	Metadata        map[string]interface{} `json:"metadata"`
}

// OptimizationOption represents a specific optimization strategy
type OptimizationOption struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CurrentCost      float64   `json:"current_cost_monthly"`
	OptimizedCost    float64   `json:"optimized_cost_monthly"`
	Savings          float64   `json:"savings_monthly"`
	SavingsPercent   float64   `json:"savings_percent"`
	ImplementationSteps []string `json:"implementation_steps"`
	TimeToImplement  string    `json:"time_to_implement"`
	RiskLevel        string    `json:"risk_level"` // "low", "medium", "high"
}

// NewS3CostCalculator creates a new S3 cost calculator for the specified region
func NewS3CostCalculator(region string) *S3CostCalculator {
	calculator := &S3CostCalculator{
		region: region,
	}
	
	// Load pricing model for the region
	calculator.pricingModel = calculator.loadPricingModel(region)
	
	return calculator
}

// loadPricingModel loads pricing information for the specified region
func (c *S3CostCalculator) loadPricingModel(region string) *S3PricingModel {
	// This would typically load from a pricing API or database
	// For now, we'll use hardcoded US East 1 pricing as a baseline
	
	baseModel := &S3PricingModel{
		Region: region,
		StorageClasses: map[string]StoragePricing{
			"STANDARD": {
				Name:              "S3 Standard",
				PricePerGBMonth:   0.023,
				MinimumStorageDays: 0,
				FirstTierGB:       50 * 1024,    // 50 TB
				SecondTierGB:      450 * 1024,   // 450 TB  
				FirstTierPrice:    0.023,        // First 50 TB
				SecondTierPrice:   0.022,        // Next 450 TB
				ThirdTierPrice:    0.021,        // Beyond 500 TB
			},
			"STANDARD_IA": {
				Name:                "S3 Standard-IA",
				PricePerGBMonth:     0.0125,
				MinimumStorageDays:  30,
				MinimumObjectSize:   128 * 1024, // 128 KB
				RetrievalFeePerGB:   0.01,
				FirstTierPrice:      0.0125,
			},
			"ONEZONE_IA": {
				Name:                "S3 One Zone-IA", 
				PricePerGBMonth:     0.01,
				MinimumStorageDays:  30,
				MinimumObjectSize:   128 * 1024,
				RetrievalFeePerGB:   0.01,
				FirstTierPrice:      0.01,
			},
			"GLACIER": {
				Name:                "S3 Glacier Flexible Retrieval",
				PricePerGBMonth:     0.004,
				MinimumStorageDays:  90,
				RetrievalFeePerGB:   0.03, // Standard retrieval
				FirstTierPrice:      0.004,
			},
			"GLACIER_IR": {
				Name:                "S3 Glacier Instant Retrieval",
				PricePerGBMonth:     0.004,
				MinimumStorageDays:  90,
				MinimumObjectSize:   128 * 1024,
				RetrievalFeePerGB:   0.03,
				FirstTierPrice:      0.004,
			},
			"DEEP_ARCHIVE": {
				Name:                "S3 Glacier Deep Archive",
				PricePerGBMonth:     0.00099,
				MinimumStorageDays:  180,
				RetrievalFeePerGB:   0.02,
				FirstTierPrice:      0.00099,
			},
			"INTELLIGENT_TIERING": {
				Name:              "S3 Intelligent-Tiering",
				PricePerGBMonth:   0.023, // Frequent access tier pricing
				FirstTierPrice:    0.023,
				// Additional monitoring fee of $0.0025 per 1,000 objects
			},
		},
		RequestPricing: RequestPricing{
			PutCopyPostList:     0.0005, // $0.50 per 1,000,000 requests
			Get:                 0.0004, // $0.40 per 1,000,000 requests  
			LifecycleTransition: 0.01,   // $0.01 per 1,000 transitions
			Select:              0.002,  // $0.002 per 1,000,000 requests
			SelectDataScanned:   0.002,  // $0.002 per GB scanned
			SelectDataReturned:  0.0007, // $0.0007 per GB returned
		},
		TransferPricing: TransferPricing{
			InboundFree:         true,
			OutboundFirstGB:     0.00,   // First 1 GB free
			OutboundNext9GB:     0.09,   // $0.09 per GB (next 9.999 TB)
			OutboundNext40GB:    0.085,  // $0.085 per GB (next 40 TB)
			OutboundNext100GB:   0.07,   // $0.07 per GB (next 100 TB) 
			OutboundBeyond150GB: 0.05,   // $0.05 per GB (beyond 150 TB)
			CrossRegionPer:      0.02,   // $0.02 per GB
			CloudFrontPer:       0.085,  // $0.085 per GB
		},
		LifecyclePricing: LifecyclePricing{
			ToIA:          0.01,  // $0.01 per 1,000 objects
			ToGlacier:     0.05,  // $0.05 per 1,000 objects
			ToGlacierIR:   0.02,  // $0.02 per 1,000 objects 
			ToGlacierDA:   0.05,  // $0.05 per 1,000 objects
			ToIntelligent: 0.01,  // $0.01 per 1,000 objects
		},
		LastUpdated: time.Now(),
	}
	
	// Apply regional pricing adjustments
	return c.adjustPricingForRegion(baseModel, region)
}

// adjustPricingForRegion applies regional pricing adjustments
func (c *S3CostCalculator) adjustPricingForRegion(baseModel *S3PricingModel, region string) *S3PricingModel {
	// Regional pricing multipliers (approximations)
	regionalMultipliers := map[string]float64{
		"us-east-1":      1.0,   // Baseline
		"us-east-2":      1.0,   // Same as us-east-1
		"us-west-1":      1.02,  // Slightly higher
		"us-west-2":      1.0,   // Same as us-east-1
		"eu-west-1":      1.05,  // European pricing
		"eu-central-1":   1.08,  // Higher European pricing
		"ap-southeast-1": 1.12,  // Asia Pacific pricing
		"ap-northeast-1": 1.15,  // Japan pricing (typically highest)
	}
	
	multiplier, exists := regionalMultipliers[region]
	if !exists {
		multiplier = 1.1 // Default 10% higher for unknown regions
	}
	
	// Apply multiplier to storage pricing
	for className, pricing := range baseModel.StorageClasses {
		pricing.PricePerGBMonth *= multiplier
		pricing.FirstTierPrice *= multiplier
		pricing.SecondTierPrice *= multiplier
		pricing.ThirdTierPrice *= multiplier
		pricing.RetrievalFeePerGB *= multiplier
		baseModel.StorageClasses[className] = pricing
	}
	
	// Apply multiplier to request pricing
	baseModel.RequestPricing.PutCopyPostList *= multiplier
	baseModel.RequestPricing.Get *= multiplier
	
	// Transfer pricing typically doesn't vary by region for outbound
	
	return baseModel
}

// AnalyzeCosts performs comprehensive cost analysis for a dataset
func (c *S3CostCalculator) AnalyzeCosts(ctx context.Context, pattern *DataPattern) (*CostAnalysis, error) {
	analysis := &CostAnalysis{
		AnalysisID:   fmt.Sprintf("cost-%d", time.Now().Unix()),
		DataPattern:  pattern,
		AnalysisTime: time.Now(),
	}
	
	// Create different cost scenarios
	scenarios := []CostScenario{
		c.createCurrentStateScenario(pattern),
		c.createOptimizedScenario(pattern),
		c.createBundledScenario(pattern),
		c.createArchivalScenario(pattern),
	}
	
	analysis.Scenarios = scenarios
	
	// Generate recommendations
	analysis.Recommendations = c.generateRecommendations(pattern, scenarios)
	
	// Generate optimization options
	analysis.Optimizations = c.generateOptimizations(pattern, scenarios)
	
	// Calculate cost ranges and potential savings
	analysis.TotalCostRange = c.calculateCostRange(scenarios)
	analysis.PotentialSavings = c.calculatePotentialSavings(scenarios)
	
	return analysis, nil
}

// createCurrentStateScenario creates a scenario representing the current state
func (c *S3CostCalculator) createCurrentStateScenario(pattern *DataPattern) CostScenario {
	config := ScenarioConfig{
		FileCount:          pattern.TotalFiles,
		TotalSizeGB:        float64(pattern.TotalSize) / (1024 * 1024 * 1024),
		StorageClass:       "STANDARD",
		BundleSmallFiles:   false,
		CompressionEnabled: false,
		CompressionRatio:   1.0,
		AccessFrequency:    "monthly",
		DownloadPercentage: 10.0, // Assume 10% of data is downloaded monthly
	}
	
	return CostScenario{
		Name:          "Current State",
		Description:   "Uploading files as-is to S3 Standard storage",
		StorageClass:  "STANDARD",
		Configuration: config,
		MonthlyCosts:  c.calculateScenarioCosts(config),
		YearlyCosts:   c.calculateYearlyCosts(c.calculateScenarioCosts(config)),
		Assumptions: []string{
			"All files uploaded to S3 Standard",
			"No compression or bundling",
			"10% of data downloaded monthly",
			"Standard access patterns",
		},
	}
}

// createOptimizedScenario creates an optimized scenario with compression and appropriate storage class
func (c *S3CostCalculator) createOptimizedScenario(pattern *DataPattern) CostScenario {
	// Determine optimal storage class based on access patterns
	storageClass := "STANDARD"
	if pattern.AccessPatterns.LikelyArchival {
		storageClass = "GLACIER"
	} else if pattern.AccessPatterns.LikelyWriteOnce {
		storageClass = "STANDARD_IA"
	}
	
	// Estimate compression ratio based on file types
	compressionRatio := c.estimateCompressionRatio(pattern)
	
	config := ScenarioConfig{
		FileCount:          pattern.TotalFiles,
		TotalSizeGB:        float64(pattern.TotalSize) / (1024 * 1024 * 1024),
		StorageClass:       storageClass,
		BundleSmallFiles:   false,
		CompressionEnabled: compressionRatio < 0.9,
		CompressionRatio:   compressionRatio,
		AccessFrequency:    "monthly",
		DownloadPercentage: 5.0, // Optimized access
	}
	
	return CostScenario{
		Name:          "Optimized Storage",
		Description:   fmt.Sprintf("Optimized with %s storage class and compression", storageClass),
		StorageClass:  storageClass,
		Configuration: config,
		MonthlyCosts:  c.calculateScenarioCosts(config),
		YearlyCosts:   c.calculateYearlyCosts(c.calculateScenarioCosts(config)),
		Assumptions: []string{
			fmt.Sprintf("Files stored in %s", storageClass),
			"Compression applied where beneficial",
			"Reduced download frequency",
			"Optimized access patterns",
		},
	}
}

// createBundledScenario creates a scenario with small file bundling
func (c *S3CostCalculator) createBundledScenario(pattern *DataPattern) CostScenario {
	// Calculate bundling effects
	smallFileCount := pattern.FileSizes.SmallFiles.CountUnder1MB
	bundleReduction := float64(0.01) // Assume 100:1 bundling ratio
	if smallFileCount > 1000 {
		bundleReduction = 0.01 // 1% of original file count after bundling
	}
	
	newFileCount := pattern.TotalFiles - smallFileCount + int64(float64(smallFileCount)*bundleReduction)
	
	config := ScenarioConfig{
		FileCount:          newFileCount,
		TotalSizeGB:        float64(pattern.TotalSize) / (1024 * 1024 * 1024),
		StorageClass:       "STANDARD",
		BundleSmallFiles:   true,
		CompressionEnabled: true,
		CompressionRatio:   0.7, // Additional compression from bundling
		AccessFrequency:    "monthly",
		DownloadPercentage: 10.0,
	}
	
	return CostScenario{
		Name:          "Bundled Small Files",
		Description:   fmt.Sprintf("Small files bundled (reduced from %d to %d objects)", pattern.TotalFiles, newFileCount),
		StorageClass:  "STANDARD",
		Configuration: config,
		MonthlyCosts:  c.calculateScenarioCosts(config),
		YearlyCosts:   c.calculateYearlyCosts(c.calculateScenarioCosts(config)),
		Assumptions: []string{
			"Small files bundled using tools like Suitcase",
			fmt.Sprintf("File count reduced by %d%%", int((1-bundleReduction)*100)),
			"Additional compression from bundling",
			"Metadata preserved for extraction",
		},
	}
}

// createArchivalScenario creates a long-term archival scenario
func (c *S3CostCalculator) createArchivalScenario(pattern *DataPattern) CostScenario {
	config := ScenarioConfig{
		FileCount:           pattern.TotalFiles,
		TotalSizeGB:         float64(pattern.TotalSize) / (1024 * 1024 * 1024),
		StorageClass:        "DEEP_ARCHIVE",
		BundleSmallFiles:    true,
		CompressionEnabled:  true,
		CompressionRatio:    c.estimateCompressionRatio(pattern),
		AccessFrequency:     "yearly",
		DownloadPercentage:  1.0, // Very low access
		LifecyclePolicyDays: 90,  // Transition after 90 days
	}
	
	return CostScenario{
		Name:          "Long-term Archive",
		Description:   "Optimized for long-term storage with minimal access",
		StorageClass:  "DEEP_ARCHIVE",
		Configuration: config,
		MonthlyCosts:  c.calculateScenarioCosts(config),
		YearlyCosts:   c.calculateYearlyCosts(c.calculateScenarioCosts(config)),
		Assumptions: []string{
			"Data transitioned to Glacier Deep Archive",
			"Very infrequent access (yearly)",
			"Bundling and compression applied",
			"Lifecycle policy for automatic transition",
		},
	}
}

// calculateScenarioCosts calculates detailed costs for a scenario
func (c *S3CostCalculator) calculateScenarioCosts(config ScenarioConfig) DetailedCosts {
	costs := DetailedCosts{}
	
	// Adjust size for compression
	effectiveSizeGB := config.TotalSizeGB * config.CompressionRatio
	
	// Storage costs
	storagePricing := c.pricingModel.StorageClasses[config.StorageClass]
	costs.Storage = c.calculateTieredStorageCost(effectiveSizeGB, storagePricing)
	
	// Request costs (initial upload)
	putRequests := float64(config.FileCount)
	costs.Requests = putRequests * c.pricingModel.RequestPricing.PutCopyPostList / 1000
	
	// Add monthly access request costs
	accessFrequencyMultiplier := map[string]float64{
		"daily":   30.0,
		"weekly":  4.0,
		"monthly": 1.0,
		"yearly":  1.0 / 12.0,
		"rarely":  1.0 / 24.0,
	}
	
	frequency := accessFrequencyMultiplier[config.AccessFrequency]
	if frequency == 0 {
		frequency = 1.0 // Default to monthly
	}
	
	getRequests := putRequests * frequency * (config.DownloadPercentage / 100)
	costs.Requests += getRequests * c.pricingModel.RequestPricing.Get / 1000
	
	// Data transfer costs (for downloads)
	downloadSizeGB := effectiveSizeGB * (config.DownloadPercentage / 100) * frequency
	costs.DataTransfer = c.calculateTransferCosts(downloadSizeGB)
	
	// Retrieval costs (for cold storage)
	if storagePricing.RetrievalFeePerGB > 0 {
		costs.Retrieval = downloadSizeGB * storagePricing.RetrievalFeePerGB
	}
	
	// Lifecycle transition costs (if applicable)
	if config.LifecyclePolicyDays > 0 {
		transitionRequests := float64(config.FileCount)
		costs.Lifecycle = transitionRequests * c.pricingModel.LifecyclePricing.ToGlacier / 1000
	}
	
	// Monitoring costs (for Intelligent Tiering)
	if config.StorageClass == "INTELLIGENT_TIERING" {
		costs.Monitoring = float64(config.FileCount) * 0.0025 / 1000 // $0.0025 per 1,000 objects
	}
	
	// Calculate total
	costs.Total = costs.Storage + costs.Requests + costs.DataTransfer + 
				  costs.Retrieval + costs.Lifecycle + costs.Monitoring
	
	return costs
}

// calculateTieredStorageCost calculates storage cost with tiered pricing
func (c *S3CostCalculator) calculateTieredStorageCost(sizeGB float64, pricing StoragePricing) float64 {
	if pricing.FirstTierGB == 0 {
		// Simple pricing
		return sizeGB * pricing.FirstTierPrice
	}
	
	// Tiered pricing
	cost := 0.0
	remaining := sizeGB
	
	// First tier
	if remaining > 0 {
		firstTierSize := math.Min(remaining, float64(pricing.FirstTierGB))
		cost += firstTierSize * pricing.FirstTierPrice
		remaining -= firstTierSize
	}
	
	// Second tier
	if remaining > 0 && pricing.SecondTierGB > 0 {
		secondTierSize := math.Min(remaining, float64(pricing.SecondTierGB))
		cost += secondTierSize * pricing.SecondTierPrice
		remaining -= secondTierSize
	}
	
	// Third tier (everything beyond)
	if remaining > 0 {
		cost += remaining * pricing.ThirdTierPrice
	}
	
	return cost
}

// calculateTransferCosts calculates data transfer costs
func (c *S3CostCalculator) calculateTransferCosts(sizeGB float64) float64 {
	transfer := c.pricingModel.TransferPricing
	cost := 0.0
	remaining := sizeGB
	
	// First 1 GB free
	if remaining > 1 {
		remaining -= 1
	} else {
		return 0 // Everything is free
	}
	
	// Next 9.999 TB
	nextTier := math.Min(remaining, 9999)
	cost += nextTier * transfer.OutboundNext9GB
	remaining -= nextTier
	
	// Next 40 TB
	if remaining > 0 {
		nextTier = math.Min(remaining, 40*1024)
		cost += nextTier * transfer.OutboundNext40GB
		remaining -= nextTier
	}
	
	// Next 100 TB
	if remaining > 0 {
		nextTier = math.Min(remaining, 100*1024)
		cost += nextTier * transfer.OutboundNext100GB
		remaining -= nextTier
	}
	
	// Beyond 150 TB
	if remaining > 0 {
		cost += remaining * transfer.OutboundBeyond150GB
	}
	
	return cost
}

// calculateYearlyCosts calculates yearly costs from monthly costs
func (c *S3CostCalculator) calculateYearlyCosts(monthly DetailedCosts) DetailedCosts {
	return DetailedCosts{
		Storage:      monthly.Storage * 12,
		Requests:     monthly.Requests * 12,
		DataTransfer: monthly.DataTransfer * 12,
		Retrieval:    monthly.Retrieval * 12,
		Lifecycle:    monthly.Lifecycle, // One-time cost
		Monitoring:   monthly.Monitoring * 12,
		Total:        (monthly.Total * 12) - (monthly.Lifecycle * 11), // Adjust for one-time lifecycle cost
	}
}

// estimateCompressionRatio estimates compression ratio based on file types
func (c *S3CostCalculator) estimateCompressionRatio(pattern *DataPattern) float64 {
	totalSize := float64(pattern.TotalSize)
	compressedSize := 0.0
	
	for _, typeInfo := range pattern.FileTypes {
		fileTypeSize := float64(typeInfo.TotalSize)
		if typeInfo.Compressible {
			compressedSize += fileTypeSize * typeInfo.CompressionEst
		} else {
			compressedSize += fileTypeSize
		}
	}
	
	if totalSize > 0 {
		return compressedSize / totalSize
	}
	
	return 0.7 // Default 30% compression
}

// generateRecommendations generates cost optimization recommendations
func (c *S3CostCalculator) generateRecommendations(pattern *DataPattern, scenarios []CostScenario) []CostRecommendation {
	var recommendations []CostRecommendation
	
	currentCost := scenarios[0].MonthlyCosts.Total
	
	// Small file bundling recommendation
	if pattern.FileSizes.SmallFiles.CountUnder1MB > 100 {
		bundledCost := 0.0
		for _, scenario := range scenarios {
			if scenario.Name == "Bundled Small Files" {
				bundledCost = scenario.MonthlyCosts.Total
				break
			}
		}
		
		savings := currentCost - bundledCost
		if savings > 0 {
			recommendations = append(recommendations, CostRecommendation{
				Type:            "bundling",
				Title:           "Bundle Small Files",
				Description:     fmt.Sprintf("Bundle %d small files to reduce request costs", pattern.FileSizes.SmallFiles.CountUnder1MB),
				EstimatedSavings: savings,
				Confidence:      0.9,
				Complexity:      "medium",
				Implementation:  "Use tools like Suitcase to bundle small files before upload",
				Tradeoffs:       []string{"Requires extraction step to access individual files", "Slightly more complex workflow"},
			})
		}
	}
	
	// Storage class optimization
	if pattern.AccessPatterns.LikelyArchival {
		recommendations = append(recommendations, CostRecommendation{
			Type:            "storage_class",
			Title:           "Use Glacier Deep Archive",
			Description:     "Data appears to be archival - consider Glacier Deep Archive for 75% storage cost savings",
			EstimatedSavings: currentCost * 0.75,
			Confidence:      0.8,
			Complexity:      "low",
			Implementation:  "Set up lifecycle policy to transition to Glacier Deep Archive after 90 days",
			Tradeoffs:       []string{"12+ hour retrieval time", "Minimum 180-day storage commitment"},
		})
	}
	
	// Compression recommendation
	compressionRatio := c.estimateCompressionRatio(pattern)
	if compressionRatio < 0.8 {
		recommendations = append(recommendations, CostRecommendation{
			Type:            "compression",
			Title:           "Enable Compression",
			Description:     fmt.Sprintf("Compress data before upload to save %.0f%% storage costs", (1-compressionRatio)*100),
			EstimatedSavings: currentCost * (1 - compressionRatio) * 0.5, // 50% of storage is typically storage costs
			Confidence:      0.85,
			Complexity:      "low",
			Implementation:  "Compress files before upload using gzip or similar",
			Tradeoffs:       []string{"Additional CPU time for compression/decompression", "Slightly more complex pipeline"},
		})
	}
	
	return recommendations
}

// generateOptimizations generates specific optimization options
func (c *S3CostCalculator) generateOptimizations(pattern *DataPattern, scenarios []CostScenario) []OptimizationOption {
	var optimizations []OptimizationOption
	
	currentCost := scenarios[0].MonthlyCosts.Total
	
	for _, scenario := range scenarios[1:] { // Skip current state scenario
		savings := currentCost - scenario.MonthlyCosts.Total
		if savings > 0 {
			optimizations = append(optimizations, OptimizationOption{
				Name:             scenario.Name,
				Description:      scenario.Description,
				CurrentCost:      currentCost,
				OptimizedCost:    scenario.MonthlyCosts.Total,
				Savings:          savings,
				SavingsPercent:   (savings / currentCost) * 100,
				ImplementationSteps: c.getImplementationSteps(scenario),
				TimeToImplement:  c.getImplementationTime(scenario),
				RiskLevel:        c.getRiskLevel(scenario),
			})
		}
	}
	
	return optimizations
}

// getImplementationSteps returns implementation steps for a scenario
func (c *S3CostCalculator) getImplementationSteps(scenario CostScenario) []string {
	switch scenario.Name {
	case "Optimized Storage":
		return []string{
			"Analyze access patterns for your data",
			"Set up lifecycle policies for automatic transitions", 
			"Enable compression in your upload pipeline",
			"Monitor cost savings and adjust policies",
		}
	case "Bundled Small Files":
		return []string{
			"Install and configure Suitcase or similar bundling tool",
			"Create bundling workflow for small files",
			"Test bundle creation and extraction",
			"Deploy bundling in production pipeline",
			"Monitor bundle sizes and optimization",
		}
	case "Long-term Archive":
		return []string{
			"Set up lifecycle policies for Deep Archive transition",
			"Configure monitoring for archived data",
			"Test retrieval process and timing",
			"Document retrieval procedures for users",
			"Monitor archival costs and policies",
		}
	default:
		return []string{"Configure optimized settings", "Test configuration", "Deploy to production"}
	}
}

// getImplementationTime estimates implementation time
func (c *S3CostCalculator) getImplementationTime(scenario CostScenario) string {
	switch scenario.Name {
	case "Optimized Storage":
		return "1-2 days"
	case "Bundled Small Files":
		return "3-5 days"
	case "Long-term Archive":
		return "1-2 days"
	default:
		return "1-3 days"
	}
}

// getRiskLevel assesses implementation risk
func (c *S3CostCalculator) getRiskLevel(scenario CostScenario) string {
	switch scenario.Name {
	case "Optimized Storage":
		return "low"
	case "Bundled Small Files":
		return "medium"
	case "Long-term Archive":
		return "low"
	default:
		return "medium"
	}
}

// calculateCostRange calculates the total cost range across scenarios
func (c *S3CostCalculator) calculateCostRange(scenarios []CostScenario) CostRange {
	minMonthly := math.MaxFloat64
	maxMonthly := 0.0
	
	for _, scenario := range scenarios {
		if scenario.MonthlyCosts.Total < minMonthly {
			minMonthly = scenario.MonthlyCosts.Total
		}
		if scenario.MonthlyCosts.Total > maxMonthly {
			maxMonthly = scenario.MonthlyCosts.Total
		}
	}
	
	return CostRange{
		MinMonthly: minMonthly,
		MaxMonthly: maxMonthly,
		MinYearly:  minMonthly * 12,
		MaxYearly:  maxMonthly * 12,
	}
}

// calculatePotentialSavings calculates maximum potential savings
func (c *S3CostCalculator) calculatePotentialSavings(scenarios []CostScenario) float64 {
	if len(scenarios) < 2 {
		return 0
	}
	
	currentCost := scenarios[0].MonthlyCosts.Total
	minCost := currentCost
	
	for _, scenario := range scenarios[1:] {
		if scenario.MonthlyCosts.Total < minCost {
			minCost = scenario.MonthlyCosts.Total
		}
	}
	
	return currentCost - minCost
}