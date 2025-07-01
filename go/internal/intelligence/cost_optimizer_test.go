package intelligence

import (
	"strings"
	"testing"

	"github.com/aws-research-wizard/go/internal/data"
)

func TestCostOptimizer_calculateBaseMonthlyCost(t *testing.T) {
	co := NewCostOptimizer()

	tests := []struct {
		name         string
		resourcePlan *ResourcePlan
		expectedMin  float64
		expectedMax  float64
	}{
		{
			name: "small_instance_basic_storage",
			resourcePlan: &ResourcePlan{
				RecommendedInstance: "c6i.2xlarge",
				StorageConfiguration: StorageConfiguration{
					PrimaryStorage: StorageType{
						Type:   "gp3",
						SizeGB: 100,
					},
				},
			},
			expectedMin: 200.0, // Roughly $0.34 * 24 * 30 + storage
			expectedMax: 300.0,
		},
		{
			name: "large_gpu_instance",
			resourcePlan: &ResourcePlan{
				RecommendedInstance: "p4d.24xlarge",
				StorageConfiguration: StorageConfiguration{
					PrimaryStorage: StorageType{
						Type:   "gp3",
						SizeGB: 1000,
					},
				},
			},
			expectedMin: 20000.0, // Roughly $32.77 * 24 * 30 + storage
			expectedMax: 30000.0,
		},
		{
			name: "unknown_instance_fallback",
			resourcePlan: &ResourcePlan{
				RecommendedInstance: "unknown.instance",
				StorageConfiguration: StorageConfiguration{
					PrimaryStorage: StorageType{
						Type:   "gp3",
						SizeGB: 100,
					},
				},
			},
			expectedMin: 700.0, // Fallback rate * 24 * 30 + storage
			expectedMax: 800.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := co.calculateBaseMonthlyCost(tt.resourcePlan)

			if cost < tt.expectedMin || cost > tt.expectedMax {
				t.Errorf("calculateBaseMonthlyCost() = %v, want between %v and %v",
					cost, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

func TestCostOptimizer_calculateSpotInstanceSavings(t *testing.T) {
	co := NewCostOptimizer()

	tests := []struct {
		name              string
		instanceType      string
		expectedNotNil    bool
		minSavingsPercent float64
	}{
		{
			name:              "compute_optimized_instance",
			instanceType:      "c6i.4xlarge",
			expectedNotNil:    true,
			minSavingsPercent: 60.0,
		},
		{
			name:              "gpu_instance",
			instanceType:      "p4d.24xlarge",
			expectedNotNil:    true,
			minSavingsPercent: 50.0,
		},
		{
			name:              "memory_optimized_instance",
			instanceType:      "r6i.8xlarge",
			expectedNotNil:    true,
			minSavingsPercent: 60.0,
		},
		{
			name:              "unknown_instance",
			instanceType:      "unknown.instance",
			expectedNotNil:    false,
			minSavingsPercent: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savings := co.calculateSpotInstanceSavings(tt.instanceType)

			if tt.expectedNotNil && savings == nil {
				t.Error("calculateSpotInstanceSavings() returned nil, expected savings")
				return
			}

			if !tt.expectedNotNil && savings != nil {
				t.Error("calculateSpotInstanceSavings() returned savings, expected nil")
				return
			}

			if savings != nil {
				if savings.PotentialSavingsPercent < tt.minSavingsPercent {
					t.Errorf("PotentialSavingsPercent = %v, want >= %v",
						savings.PotentialSavingsPercent, tt.minSavingsPercent)
				}

				if savings.EstimatedMonthlySavings <= 0 {
					t.Error("EstimatedMonthlySavings should be positive")
				}

				if savings.RiskAssessment == "" {
					t.Error("RiskAssessment should not be empty")
				}

				if savings.RecommendedStrategy == "" {
					t.Error("RecommendedStrategy should not be empty")
				}
			}
		})
	}
}

func TestCostOptimizer_calculateReservedInstanceSavings(t *testing.T) {
	co := NewCostOptimizer()

	tests := []struct {
		name                string
		instanceType        string
		monthlyOnDemandCost float64
		expectedMinSavings  float64
	}{
		{
			name:                "medium_cost_instance",
			instanceType:        "c6i.4xlarge",
			monthlyOnDemandCost: 500.0,
			expectedMinSavings:  1000.0, // Should save at least this much over 1 year
		},
		{
			name:                "high_cost_instance",
			instanceType:        "p4d.24xlarge",
			monthlyOnDemandCost: 24000.0,
			expectedMinSavings:  50000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savings := co.calculateReservedInstanceSavings(tt.instanceType, tt.monthlyOnDemandCost)

			if savings == nil {
				t.Fatal("calculateReservedInstanceSavings() returned nil")
			}

			if savings.OneYearSavings < tt.expectedMinSavings {
				t.Errorf("OneYearSavings = %v, want >= %v",
					savings.OneYearSavings, tt.expectedMinSavings)
			}

			if savings.ThreeYearSavings <= savings.OneYearSavings {
				t.Error("ThreeYearSavings should be greater than OneYearSavings")
			}

			if savings.RecommendedTerm == "" {
				t.Error("RecommendedTerm should not be empty")
			}

			if savings.PaymentOption == "" {
				t.Error("PaymentOption should not be empty")
			}

			if savings.BreakevenPoint == "" {
				t.Error("BreakevenPoint should not be empty")
			}
		})
	}
}

func TestCostOptimizer_generateStorageOptimizations(t *testing.T) {
	co := NewCostOptimizer()

	// Create test data
	resourcePlan := &ResourcePlan{
		StorageConfiguration: StorageConfiguration{
			PrimaryStorage: StorageType{
				Type:   "gp2",
				SizeGB: 1000,
			},
		},
	}

	dataPattern := &data.DataPattern{
		TotalSize: 1024 * 1024 * 1024 * 200, // 200 GB
		AccessPatterns: data.AccessPatternAnalysis{
			LikelyArchival:  true,
			LikelyWriteOnce: true,
		},
	}

	dataRec := &data.RecommendationResult{
		DataPattern: dataPattern,
	}

	optimizations := co.generateStorageOptimizations(resourcePlan, dataRec)

	if len(optimizations) == 0 {
		t.Error("generateStorageOptimizations() returned no optimizations")
	}

	// Should have intelligent tiering for large datasets
	hasIntelligentTiering := false
	hasGlacierArchival := false
	hasEBSOptimization := false

	for _, opt := range optimizations {
		switch opt.Type {
		case "intelligent_tiering":
			hasIntelligentTiering = true
		case "glacier_archival":
			hasGlacierArchival = true
		case "ebs_gp3_migration":
			hasEBSOptimization = true
		}

		// Validate optimization structure
		if opt.Description == "" {
			t.Error("Optimization description should not be empty")
		}

		if opt.SavingsPercent <= 0 {
			t.Error("SavingsPercent should be positive")
		}

		if opt.MonthlySavings < 0 {
			t.Error("MonthlySavings should not be negative")
		}

		if opt.Implementation == "" {
			t.Error("Implementation should not be empty")
		}
	}

	if dataPattern.TotalSize > 1024*1024*1024*100 && !hasIntelligentTiering {
		t.Error("Should recommend intelligent tiering for large datasets")
	}

	if dataPattern.AccessPatterns.LikelyArchival && !hasGlacierArchival {
		t.Error("Should recommend Glacier archival for archival patterns")
	}

	if resourcePlan.StorageConfiguration.PrimaryStorage.Type == "gp2" && !hasEBSOptimization {
		t.Error("Should recommend gp3 migration for gp2 volumes")
	}
}

func TestCostOptimizer_GenerateCostOptimizationPlan(t *testing.T) {
	co := NewCostOptimizer()

	resourcePlan := &ResourcePlan{
		RecommendedInstance: "c6i.4xlarge",
		StorageConfiguration: StorageConfiguration{
			PrimaryStorage: StorageType{
				Type:   "gp3",
				SizeGB: 500,
			},
		},
	}

	dataPattern := &data.DataPattern{
		TotalSize:  1024 * 1024 * 1024 * 100, // 100 GB
		TotalFiles: 1000,
	}

	dataRecommendations := &data.RecommendationResult{
		DataPattern: dataPattern,
	}

	plan := co.GenerateCostOptimizationPlan("genomics", resourcePlan, dataRecommendations)

	if plan == nil {
		t.Fatal("GenerateCostOptimizationPlan() returned nil")
	}

	// Validate plan structure
	if plan.EstimatedMonthlyCost <= 0 {
		t.Error("EstimatedMonthlyCost should be positive")
	}

	if plan.OptimizedMonthlyCost < 0 {
		t.Error("OptimizedMonthlyCost should not be negative")
	}

	if plan.OptimizedMonthlyCost > plan.EstimatedMonthlyCost {
		t.Error("OptimizedMonthlyCost should be less than or equal to EstimatedMonthlyCost")
	}

	if plan.PotentialSavings < 0 {
		t.Error("PotentialSavings should not be negative")
	}

	if plan.SavingsPercentage < 0 || plan.SavingsPercentage > 100 {
		t.Errorf("SavingsPercentage should be between 0 and 100, got %v", plan.SavingsPercentage)
	}

	if len(plan.Recommendations) == 0 {
		t.Error("Should have at least one recommendation")
	}

	// Validate spot instance savings
	if plan.SpotInstanceSavings != nil {
		if plan.SpotInstanceSavings.PotentialSavingsPercent <= 0 {
			t.Error("SpotInstanceSavings PotentialSavingsPercent should be positive")
		}
		if plan.SpotInstanceSavings.EstimatedMonthlySavings <= 0 {
			t.Error("SpotInstanceSavings EstimatedMonthlySavings should be positive")
		}
	}

	// Validate reserved instance savings
	if plan.ReservedInstanceSavings != nil {
		if plan.ReservedInstanceSavings.OneYearSavings <= 0 {
			t.Error("ReservedInstanceSavings OneYearSavings should be positive")
		}
		if plan.ReservedInstanceSavings.ThreeYearSavings <= plan.ReservedInstanceSavings.OneYearSavings {
			t.Error("ThreeYearSavings should be greater than OneYearSavings")
		}
	}
}

func TestCostOptimizer_generateCostRecommendations(t *testing.T) {
	co := NewCostOptimizer()

	resourcePlan := &ResourcePlan{
		RecommendedInstance: "c6i.4xlarge",
	}

	spotSavings := &SpotInstanceSavings{
		PotentialSavingsPercent: 70.0,
		RiskAssessment:          "medium",
	}

	reservedSavings := &ReservedInstanceSavings{
		OneYearSavings:  1000.0,
		RecommendedTerm: "1-year",
	}

	storageOpts := []StorageOptimization{
		{
			Description:    "Enable S3 Intelligent Tiering",
			MonthlySavings: 50.0,
			SavingsPercent: 30.0,
		},
	}

	recommendations := co.generateCostRecommendations(
		"genomics", resourcePlan, spotSavings, reservedSavings, storageOpts)

	if len(recommendations) == 0 {
		t.Error("generateCostRecommendations() should return recommendations")
	}

	// Check for spot instance recommendation
	hasSpotRecommendation := false
	hasReservedRecommendation := false
	hasStorageRecommendation := false
	hasDomainSpecificRecommendation := false

	for _, rec := range recommendations {
		if rec == "" {
			t.Error("Recommendation should not be empty")
		}

		recLower := strings.ToLower(rec)
		if strings.Contains(recLower, "spot") {
			hasSpotRecommendation = true
		}
		if strings.Contains(recLower, "reserved") {
			hasReservedRecommendation = true
		}
		if strings.Contains(recLower, "storage") || strings.Contains(recLower, "tiering") {
			hasStorageRecommendation = true
		}
		if strings.Contains(recLower, "genomics") {
			hasDomainSpecificRecommendation = true
		}
	}

	if !hasSpotRecommendation {
		t.Error("Should include spot instance recommendation")
	}

	if !hasReservedRecommendation {
		t.Error("Should include reserved instance recommendation")
	}

	if !hasStorageRecommendation {
		t.Error("Should include storage optimization recommendation")
	}

	if !hasDomainSpecificRecommendation {
		t.Error("Should include domain-specific recommendation")
	}
}

func TestCostOptimizer_EstimateMonthlyCost(t *testing.T) {
	co := NewCostOptimizer()

	tests := []struct {
		name          string
		instanceType  string
		hoursPerMonth float64
		expectedMin   float64
		expectedMax   float64
	}{
		{
			name:          "c6i_2xlarge_full_month",
			instanceType:  "c6i.2xlarge",
			hoursPerMonth: 720.0, // 24 * 30
			expectedMin:   200.0,
			expectedMax:   300.0,
		},
		{
			name:          "p4d_24xlarge_half_month",
			instanceType:  "p4d.24xlarge",
			hoursPerMonth: 360.0, // 12 hours/day * 30 days
			expectedMin:   10000.0,
			expectedMax:   15000.0,
		},
		{
			name:          "unknown_instance",
			instanceType:  "unknown.instance",
			hoursPerMonth: 720.0,
			expectedMin:   0.0,
			expectedMax:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := co.EstimateMonthlyCost(tt.instanceType, tt.hoursPerMonth)

			if tt.expectedMin == 0.0 && tt.expectedMax == 0.0 {
				if cost != 0.0 {
					t.Errorf("EstimateMonthlyCost() = %v, want 0 for unknown instance", cost)
				}
			} else {
				if cost < tt.expectedMin || cost > tt.expectedMax {
					t.Errorf("EstimateMonthlyCost() = %v, want between %v and %v",
						cost, tt.expectedMin, tt.expectedMax)
				}
			}
		})
	}
}

func TestCostOptimizer_CompareInstanceCosts(t *testing.T) {
	co := NewCostOptimizer()

	instanceTypes := []string{"c6i.2xlarge", "c6i.4xlarge", "r6i.4xlarge", "unknown.instance"}

	costs := co.CompareInstanceCosts(instanceTypes)

	// Should have costs for known instances
	knownInstances := []string{"c6i.2xlarge", "c6i.4xlarge", "r6i.4xlarge"}
	for _, instance := range knownInstances {
		if cost, exists := costs[instance]; !exists || cost <= 0 {
			t.Errorf("Should have positive cost for %s, got %v", instance, cost)
		}
	}

	// Should not have cost for unknown instance
	if cost, exists := costs["unknown.instance"]; exists {
		t.Errorf("Should not have cost for unknown instance, got %v", cost)
	}

	// Verify cost ordering (larger instances should cost more)
	if costs["c6i.4xlarge"] <= costs["c6i.2xlarge"] {
		t.Error("c6i.4xlarge should cost more than c6i.2xlarge")
	}
}

// Test helper functions

func TestCostOptimizer_HelperFunctions(t *testing.T) {
	t.Run("isComputeOptimized", func(t *testing.T) {
		tests := []struct {
			instanceType string
			expected     bool
		}{
			{"c6i.2xlarge", true},
			{"c6a.4xlarge", true},
			{"c5.xlarge", true},
			{"r6i.2xlarge", false},
			{"g5.xlarge", false},
			{"p4d.24xlarge", false},
		}

		for _, tt := range tests {
			result := isComputeOptimized(tt.instanceType)
			if result != tt.expected {
				t.Errorf("isComputeOptimized(%s) = %v, want %v",
					tt.instanceType, result, tt.expected)
			}
		}
	})

	t.Run("isGPUInstance", func(t *testing.T) {
		tests := []struct {
			instanceType string
			expected     bool
		}{
			{"g5.xlarge", true},
			{"g4dn.2xlarge", true},
			{"p4d.24xlarge", true},
			{"p3.8xlarge", true},
			{"c6i.2xlarge", false},
			{"r6i.4xlarge", false},
		}

		for _, tt := range tests {
			result := isGPUInstance(tt.instanceType)
			if result != tt.expected {
				t.Errorf("isGPUInstance(%s) = %v, want %v",
					tt.instanceType, result, tt.expected)
			}
		}
	})

	t.Run("isMemoryOptimized", func(t *testing.T) {
		tests := []struct {
			instanceType string
			expected     bool
		}{
			{"r6i.2xlarge", true},
			{"r5.4xlarge", true},
			{"x1e.xlarge", true},
			{"z1d.large", true},
			{"c6i.2xlarge", false},
			{"g5.xlarge", false},
		}

		for _, tt := range tests {
			result := isMemoryOptimized(tt.instanceType)
			if result != tt.expected {
				t.Errorf("isMemoryOptimized(%s) = %v, want %v",
					tt.instanceType, result, tt.expected)
			}
		}
	})
}

// Benchmark tests

func BenchmarkCostOptimizer_calculateBaseMonthlyCost(b *testing.B) {
	co := NewCostOptimizer()
	resourcePlan := &ResourcePlan{
		RecommendedInstance: "c6i.4xlarge",
		StorageConfiguration: StorageConfiguration{
			PrimaryStorage: StorageType{
				Type:   "gp3",
				SizeGB: 500,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		co.calculateBaseMonthlyCost(resourcePlan)
	}
}

func BenchmarkCostOptimizer_GenerateCostOptimizationPlan(b *testing.B) {
	co := NewCostOptimizer()
	resourcePlan := &ResourcePlan{
		RecommendedInstance: "c6i.4xlarge",
		StorageConfiguration: StorageConfiguration{
			PrimaryStorage: StorageType{
				Type:   "gp3",
				SizeGB: 500,
			},
		},
	}

	dataRecommendations := &data.RecommendationResult{
		DataPattern: &data.DataPattern{
			TotalSize:  1024 * 1024 * 1024 * 100,
			TotalFiles: 1000,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		co.GenerateCostOptimizationPlan("genomics", resourcePlan, dataRecommendations)
	}
}
