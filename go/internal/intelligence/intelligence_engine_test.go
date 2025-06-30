package intelligence

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws-research-wizard/go/internal/data"
)

func TestIntelligenceEngine_detectDomain(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	tests := []struct {
		name        string
		dataPath    string
		hints       DomainHints
		expectedDomain string
		minConfidence float64
	}{
		{
			name:         "genomics_with_fastq_files",
			dataPath:     "/data/samples/sample1.fastq",
			hints:        DomainHints{},
			expectedDomain: "genomics",
			minConfidence: 0.2, // Lower expectation as only 1 file type matches
		},
		{
			name:         "climate_with_netcdf_files",
			dataPath:     "/data/weather/temperature.nc",
			hints:        DomainHints{},
			expectedDomain: "climate",
			minConfidence: 0.2, // Lower expectation as only 1 file type matches
		},
		{
			name:         "explicit_domain_hint",
			dataPath:     "/data/unknown/file.dat",
			hints:        DomainHints{ExplicitDomain: "machine_learning"},
			expectedDomain: "machine_learning",
			minConfidence: 0.8,
		},
		{
			name:         "workflow_hints",
			dataPath:     "/data/analysis/",
			hints:        DomainHints{WorkflowHints: []string{"variant_calling", "genomics"}},
			expectedDomain: "genomics",
			minConfidence: 0.3,
		},
		{
			name:         "tool_hints",
			dataPath:     "/data/ml/",
			hints:        DomainHints{ToolHints: []string{"pytorch", "tensorflow"}},
			expectedDomain: "machine_learning",
			minConfidence: 0.2,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, confidence := ie.detectDomain(tt.dataPath, tt.hints)
			
			if domain != tt.expectedDomain {
				t.Errorf("detectDomain() domain = %v, want %v", domain, tt.expectedDomain)
			}
			
			if confidence < tt.minConfidence {
				t.Errorf("detectDomain() confidence = %v, want >= %v", confidence, tt.minConfidence)
			}
		})
	}
}

func TestIntelligenceEngine_analyzeFileExtensions(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	tests := []struct {
		name        string
		dataPath    string
		expectedExts map[string]bool
	}{
		{
			name:     "single_fastq_file",
			dataPath: "/data/sample.fastq",
			expectedExts: map[string]bool{".fastq": true},
		},
		{
			name:     "netcdf_file",
			dataPath: "/data/climate_model.nc",
			expectedExts: map[string]bool{".nc": true},
		},
		{
			name:     "genomics_path_with_multiple_hints",
			dataPath: "/data/genomics/sample.bam",
			expectedExts: map[string]bool{".bam": true},
		},
		{
			name:     "no_extension",
			dataPath: "/data/unknown",
			expectedExts: map[string]bool{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exts := ie.analyzeFileExtensions(tt.dataPath)
			
			for expectedExt := range tt.expectedExts {
				if !exts[expectedExt] {
					t.Errorf("analyzeFileExtensions() missing expected extension %v", expectedExt)
				}
			}
			
			for foundExt := range exts {
				if !tt.expectedExts[foundExt] {
					t.Errorf("analyzeFileExtensions() found unexpected extension %v", foundExt)
				}
			}
		})
	}
}

func TestIntelligenceEngine_detectFromCommonPatterns(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	tests := []struct {
		name        string
		dataPath    string
		expectedDomain string
	}{
		{
			name:        "genomics_pattern",
			dataPath:    "/data/genome-sequencing/samples/",
			expectedDomain: "genomics",
		},
		{
			name:        "climate_pattern",
			dataPath:    "/data/weather-forecast/models/",
			expectedDomain: "climate",
		},
		{
			name:        "ml_pattern",
			dataPath:    "/data/neural-network-training/",
			expectedDomain: "machine_learning",
		},
		{
			name:        "no_pattern_match",
			dataPath:    "/data/unknown/files/",
			expectedDomain: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := ie.detectFromCommonPatterns(tt.dataPath)
			
			if domain != tt.expectedDomain {
				t.Errorf("detectFromCommonPatterns() = %v, want %v", domain, tt.expectedDomain)
			}
		})
	}
}

func TestIntelligenceEngine_assessWorkloadSize(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	tests := []struct {
		name        string
		totalSizeGB float64
		expectedSize string
	}{
		{
			name:        "small_workload",
			totalSizeGB: 5.0,
			expectedSize: "small",
		},
		{
			name:        "medium_workload",
			totalSizeGB: 100.0,
			expectedSize: "medium",
		},
		{
			name:        "large_workload",
			totalSizeGB: 1000.0,
			expectedSize: "large",
		},
		{
			name:        "massive_workload",
			totalSizeGB: 10000.0,
			expectedSize: "massive",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := &data.DataPattern{
				TotalSize: int64(tt.totalSizeGB * 1024 * 1024 * 1024),
			}
			
			size := ie.assessWorkloadSize(pattern)
			
			if size != tt.expectedSize {
				t.Errorf("assessWorkloadSize() = %v, want %v", size, tt.expectedSize)
			}
		})
	}
}

func TestIntelligenceEngine_selectOptimalInstance(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	// Mock domain pack loader to return test data
	ie.domainPackLoader = &mockDomainPackLoader{
		domainPacks: map[string]*DomainPackInfo{
			"genomics": {
				Name: "genomics",
				InstanceTypes: map[string]string{
					"small":  "r6i.2xlarge",
					"medium": "r6i.4xlarge",
					"large":  "r6i.8xlarge",
				},
			},
			"machine_learning": {
				Name: "machine_learning",
				InstanceTypes: map[string]string{
					"small":  "g5.xlarge",
					"medium": "g5.4xlarge",
					"large":  "p4d.24xlarge",
				},
			},
		},
	}
	
	tests := []struct {
		name         string
		domain       string
		workloadSize string
		expectedInstance string
	}{
		{
			name:         "genomics_medium",
			domain:       "genomics",
			workloadSize: "medium",
			expectedInstance: "r6i.4xlarge",
		},
		{
			name:         "ml_large",
			domain:       "machine_learning", 
			workloadSize: "large",
			expectedInstance: "p4d.24xlarge",
		},
		{
			name:         "unknown_domain_fallback",
			domain:       "unknown",
			workloadSize: "medium",
			expectedInstance: "c6i.4xlarge",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &data.ResearchDomainProfile{Name: tt.domain}
			hints := DomainHints{}
			
			instance := ie.selectOptimalInstance(profile, tt.workloadSize, hints)
			
			if instance != tt.expectedInstance {
				t.Errorf("selectOptimalInstance() = %v, want %v", instance, tt.expectedInstance)
			}
		})
	}
}

func TestIntelligenceEngine_GenerateIntelligentRecommendations(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	// Setup test data
	ctx := context.Background()
	dataPath := "/data/genomics/samples.fastq"
	hints := DomainHints{
		ExplicitDomain: "genomics",
		DataSizeHint:   "medium",
	}
	
	// Mock the recommendation engine
	ie.recommendationEngine = &mockRecommendationEngine{}
	
	recommendation, err := ie.GenerateIntelligentRecommendations(ctx, dataPath, hints)
	
	if err != nil {
		t.Fatalf("GenerateIntelligentRecommendations() error = %v", err)
	}
	
	if recommendation == nil {
		t.Fatal("GenerateIntelligentRecommendations() returned nil recommendation")
	}
	
	// Verify basic fields
	if recommendation.Domain != "genomics" {
		t.Errorf("Expected domain = genomics, got %v", recommendation.Domain)
	}
	
	if recommendation.Confidence < 0.5 {
		t.Errorf("Expected confidence >= 0.5, got %v", recommendation.Confidence)
	}
	
	if recommendation.ResourcePlan == nil {
		t.Error("Expected ResourcePlan to be populated")
	}
	
	if recommendation.CostOptimization == nil {
		t.Error("Expected CostOptimization to be populated")
	}
	
	if recommendation.Implementation == nil {
		t.Error("Expected Implementation to be populated")
	}
	
	if recommendation.Impact == nil {
		t.Error("Expected Impact to be populated")
	}
}

func TestIntelligenceEngine_generateResourcePlan(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	// Create test data
	dataPattern := &data.DataPattern{
		TotalSize:  1024 * 1024 * 1024 * 100, // 100 GB
		TotalFiles: 1000,
		FileSizes: data.FileSizeAnalysis{
			MeanSize: 1024 * 1024 * 100, // 100 MB average
		},
	}
	
	dataRec := &data.RecommendationResult{
		DataPattern: dataPattern,
	}
	
	profile := ie.domainProfileManager.GetAllProfiles()["genomics"]
	if profile == nil {
		t.Fatal("Failed to get genomics profile")
	}
	
	hints := DomainHints{}
	
	plan := ie.generateResourcePlan("genomics", dataRec, hints)
	
	if plan == nil {
		t.Fatal("generateResourcePlan() returned nil")
	}
	
	if plan.RecommendedInstance == "" {
		t.Error("Expected RecommendedInstance to be set")
	}
	
	if len(plan.AlternativeInstances) == 0 {
		t.Error("Expected AlternativeInstances to be populated")
	}
	
	if plan.StorageConfiguration.PrimaryStorage.SizeGB == 0 {
		t.Error("Expected storage configuration to be set")
	}
	
	if plan.Reasoning == "" {
		t.Error("Expected reasoning to be provided")
	}
}

// Mock implementations for testing

type mockDomainPackLoader struct {
	domainPacks map[string]*DomainPackInfo
}

func (m *mockDomainPackLoader) LoadDomainPack(domainName string) (*DomainPackInfo, error) {
	if pack, exists := m.domainPacks[domainName]; exists {
		return pack, nil
	}
	return nil, fmt.Errorf("domain pack not found: %s", domainName)
}

func (m *mockDomainPackLoader) LoadAllDomainPacks() (map[string]*DomainPackInfo, error) {
	return m.domainPacks, nil
}

func (m *mockDomainPackLoader) GetAvailableDomains() ([]string, error) {
	var domains []string
	for name := range m.domainPacks {
		domains = append(domains, name)
	}
	return domains, nil
}

func (m *mockDomainPackLoader) ValidateDomainPack(domainName string) error {
	_, exists := m.domainPacks[domainName]
	if !exists {
		return fmt.Errorf("domain pack not found: %s", domainName)
	}
	return nil
}

func (m *mockDomainPackLoader) ClearCache() {
	// No-op for mock
}

type mockRecommendationEngine struct{}

func (m *mockRecommendationEngine) GenerateRecommendations(ctx context.Context, dataPath string) (*data.RecommendationResult, error) {
	return &data.RecommendationResult{
		AnalysisID: "test-analysis",
		DataPath:   dataPath,
		DataPattern: &data.DataPattern{
			TotalSize:  1024 * 1024 * 1024 * 50, // 50 GB
			TotalFiles: 500,
			FileSizes: data.FileSizeAnalysis{
				MeanSize: 1024 * 1024 * 100, // 100 MB
			},
			AccessPatterns: data.AccessPatternAnalysis{
				LikelyArchival:   false,
				LikelyWriteOnce:  true,
			},
		},
		CostAnalysis: &data.CostAnalysis{
			Scenarios: []data.CostScenario{
				{
					Name: "Standard",
					MonthlyCosts: data.DetailedCosts{
						Total: 100.0,
					},
				},
			},
		},
	}, nil
}

// Helper function to create a test intelligence engine
func createTestIntelligenceEngine() *IntelligenceEngine {
	domainProfileManager := data.NewResearchDomainProfileManager()
	
	// Create mock recommendation engine
	mockRecEngine := &mockRecommendationEngine{}
	
	ie := NewIntelligenceEngine(domainProfileManager, mockRecEngine)
	
	// Set up mock domain pack loader
	ie.domainPackLoader = &mockDomainPackLoader{
		domainPacks: map[string]*DomainPackInfo{
			"genomics": {
				Name:        "genomics",
				Version:     "1.0.0",
				Description: "Genomics research tools",
				InstanceTypes: map[string]string{
					"small":  "r6i.2xlarge",
					"medium": "r6i.4xlarge",
					"large":  "r6i.8xlarge",
				},
			},
		},
	}
	
	return ie
}

// Benchmark tests for performance validation

func BenchmarkIntelligenceEngine_detectDomain(b *testing.B) {
	ie := createTestIntelligenceEngine()
	dataPath := "/data/genomics/sample.fastq"
	hints := DomainHints{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ie.detectDomain(dataPath, hints)
	}
}

func BenchmarkIntelligenceEngine_GenerateIntelligentRecommendations(b *testing.B) {
	ie := createTestIntelligenceEngine()
	ctx := context.Background()
	dataPath := "/data/genomics/samples.fastq"
	hints := DomainHints{ExplicitDomain: "genomics"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ie.GenerateIntelligentRecommendations(ctx, dataPath, hints)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test edge cases and error conditions

func TestIntelligenceEngine_EdgeCases(t *testing.T) {
	ie := createTestIntelligenceEngine()
	
	t.Run("nil_data_pattern", func(t *testing.T) {
		size := ie.assessWorkloadSize(nil)
		if size == "" {
			t.Error("assessWorkloadSize() should handle nil gracefully")
		}
	})
	
	t.Run("empty_domain_hints", func(t *testing.T) {
		domain, confidence := ie.detectDomain("", DomainHints{})
		if domain == "" {
			t.Error("detectDomain() should return default domain for empty input")
		}
		if confidence < 0 || confidence > 1 {
			t.Errorf("detectDomain() confidence should be between 0 and 1, got %v", confidence)
		}
	})
	
	t.Run("invalid_data_path", func(t *testing.T) {
		ctx := context.Background()
		_, err := ie.GenerateIntelligentRecommendations(ctx, "", DomainHints{})
		// Should handle gracefully, may return error or default recommendations
		if err != nil {
			t.Logf("Expected error for invalid path: %v", err)
		}
	})
}

// Test concurrent access safety

func TestIntelligenceEngine_ConcurrentAccess(t *testing.T) {
	ie := createTestIntelligenceEngine()
	ctx := context.Background()
	
	// Run multiple goroutines simultaneously
	const numGoroutines = 10
	errors := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			dataPath := fmt.Sprintf("/data/test%d.fastq", id)
			hints := DomainHints{ExplicitDomain: "genomics"}
			
			_, err := ie.GenerateIntelligentRecommendations(ctx, dataPath, hints)
			errors <- err
		}(i)
	}
	
	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		if err != nil {
			t.Errorf("Concurrent access error: %v", err)
		}
	}
}