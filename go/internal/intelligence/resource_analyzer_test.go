package intelligence

import (
	"strings"
	"testing"

	"github.com/aws-research-wizard/go/internal/data"
)

func TestResourceAnalyzer_AnalyzeResourceRequirements(t *testing.T) {
	analyzer := NewResourceAnalyzer()
	
	tests := []struct {
		name           string
		domain         string
		dataPattern    *data.DataPattern
		expectedMinCPU int
	}{
		{
			name:   "genomics_small",
			domain: "genomics",
			dataPattern: &data.DataPattern{
				TotalSize:  10 * 1024 * 1024 * 1024, // 10 GB
				TotalFiles: 1000,
			},
			expectedMinCPU: 4,
		},
		{
			name:   "genomics_large",
			domain: "genomics",
			dataPattern: &data.DataPattern{
				TotalSize:  1000 * 1024 * 1024 * 1024, // 1 TB
				TotalFiles: 100000,
			},
			expectedMinCPU: 8,
		},
		{
			name:   "machine_learning",
			domain: "machine_learning",
			dataPattern: &data.DataPattern{
				TotalSize:  100 * 1024 * 1024 * 1024, // 100 GB
				TotalFiles: 10000,
			},
			expectedMinCPU: 4,
		},
		{
			name:   "climate_modeling",
			domain: "climate",
			dataPattern: &data.DataPattern{
				TotalSize:  500 * 1024 * 1024 * 1024, // 500 GB
				TotalFiles: 50000,
			},
			expectedMinCPU: 4,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &data.ResearchDomainProfile{Name: tt.domain}
			hints := DomainHints{}
			
			requirements := analyzer.AnalyzeResourceRequirements(tt.domain, tt.dataPattern, profile, hints)
			
			if requirements == nil {
				t.Fatal("AnalyzeResourceRequirements() returned nil")
			}
			
			if requirements.MinCPUs < tt.expectedMinCPU {
				t.Errorf("AnalyzeResourceRequirements() MinCPUs = %v, want at least %v", requirements.MinCPUs, tt.expectedMinCPU)
			}
			
			// Verify required fields are set
			if requirements.MinMemoryGB <= 0 {
				t.Error("Expected MinMemoryGB to be positive")
			}
			
			if requirements.MinStorageGB <= 0 {
				t.Error("Expected MinStorageGB to be positive")
			}
			
			if requirements.NetworkRequirements.MinBandwidthMbps <= 0 {
				t.Error("Expected MinBandwidthMbps to be positive")
			}
		})
	}
}

func TestResourceAnalyzer_FindOptimalInstance(t *testing.T) {
	analyzer := NewResourceAnalyzer()
	
	tests := []struct {
		name             string
		requirements     *ResourceRequirement
		expectedCount    int
		expectContains   string
	}{
		{
			name: "small_requirements",
			requirements: &ResourceRequirement{
				MinCPUs:     4,
				MinMemoryGB: 8,
			},
			expectedCount:  1,
			expectContains: "c6i",
		},
		{
			name: "medium_requirements",
			requirements: &ResourceRequirement{
				MinCPUs:     16,
				MinMemoryGB: 32,
			},
			expectedCount:  1,
			expectContains: "c6i",
		},
		{
			name: "large_requirements",
			requirements: &ResourceRequirement{
				MinCPUs:     32,
				MinMemoryGB: 64,
			},
			expectedCount:  1,
			expectContains: "c6i",
		},
		{
			name: "memory_intensive",
			requirements: &ResourceRequirement{
				MinCPUs:     8,
				MinMemoryGB: 64,
			},
			expectedCount:  1,
			expectContains: "r6i",
		},
		{
			name: "gpu_requirements",
			requirements: &ResourceRequirement{
				MinCPUs:     16,
				MinMemoryGB: 32,
				GPURequired: true,
				MinGPUs:     1,
			},
			expectedCount:  1,
			expectContains: "g5",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendations := analyzer.FindOptimalInstance(tt.requirements)
			
			if len(recommendations) < tt.expectedCount {
				t.Errorf("FindOptimalInstance() returned %d recommendations, want at least %d", len(recommendations), tt.expectedCount)
			}
			
			// Check that expected instance type is present
			if tt.expectContains != "" && len(recommendations) > 0 {
				found := false
				for _, instanceType := range recommendations {
					if strings.Contains(instanceType, tt.expectContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected instance type containing %s not found in recommendations %v", tt.expectContains, recommendations)
				}
			}
		})
	}
}

func TestResourceAnalyzer_GetInstanceSpec(t *testing.T) {
	analyzer := NewResourceAnalyzer()
	
	tests := []struct {
		name         string
		instanceType string
		expectExists bool
	}{
		{
			name:         "valid_c6i_instance",
			instanceType: "c6i.xlarge",
			expectExists: true,
		},
		{
			name:         "valid_r6i_instance",
			instanceType: "r6i.2xlarge",
			expectExists: true,
		},
		{
			name:         "valid_g5_instance",
			instanceType: "g5.xlarge",
			expectExists: true,
		},
		{
			name:         "invalid_instance",
			instanceType: "invalid.type",
			expectExists: false,
		},
		{
			name:         "empty_instance",
			instanceType: "",
			expectExists: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, exists := analyzer.GetInstanceSpec(tt.instanceType)
			
			if exists != tt.expectExists {
				t.Errorf("GetInstanceSpec() exists = %v, want %v", exists, tt.expectExists)
			}
			
			if exists {
				if spec.vCPUs <= 0 {
					t.Error("Expected positive vCPUs for valid instance")
				}
				if spec.MemoryGB <= 0 {
					t.Error("Expected positive memory for valid instance")
				}
				if spec.InstanceFamily == "" {
					t.Error("Expected non-empty instance family")
				}
			}
		})
	}
}

func TestResourceAnalyzer_EdgeCases(t *testing.T) {
	analyzer := NewResourceAnalyzer()
	
	t.Run("nil_data_pattern", func(t *testing.T) {
		profile := &data.ResearchDomainProfile{Name: "genomics"}
		hints := DomainHints{}
		
		requirements := analyzer.AnalyzeResourceRequirements("genomics", nil, profile, hints)
		if requirements == nil {
			t.Error("Expected non-nil requirements for nil data pattern")
		}
	})
	
	t.Run("nil_profile", func(t *testing.T) {
		dataPattern := &data.DataPattern{
			TotalSize:  100 * 1024 * 1024 * 1024, // 100 GB
			TotalFiles: 10000,
		}
		hints := DomainHints{}
		
		requirements := analyzer.AnalyzeResourceRequirements("genomics", dataPattern, nil, hints)
		if requirements == nil {
			t.Error("Expected non-nil requirements for nil profile")
		}
	})
	
	t.Run("empty_requirements", func(t *testing.T) {
		emptyReq := &ResourceRequirement{}
		instances := analyzer.FindOptimalInstance(emptyReq)
		
		// Should still return some default instances
		if len(instances) == 0 {
			t.Error("Expected at least one instance recommendation for empty requirements")
		}
	})
	
	t.Run("extreme_requirements", func(t *testing.T) {
		extremeReq := &ResourceRequirement{
			MinCPUs:     1000,
			MinMemoryGB: 10000,
		}
		instances := analyzer.FindOptimalInstance(extremeReq)
		
		// Should handle gracefully, even if no instances meet requirements
		if instances == nil {
			t.Error("FindOptimalInstance should not return nil for extreme requirements")
		}
	})
}

func BenchmarkResourceAnalyzer_AnalyzeResourceRequirements(b *testing.B) {
	analyzer := NewResourceAnalyzer()
	profile := &data.ResearchDomainProfile{Name: "genomics"}
	dataPattern := &data.DataPattern{
		TotalSize:  100 * 1024 * 1024 * 1024, // 100 GB
		TotalFiles: 10000,
	}
	hints := DomainHints{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.AnalyzeResourceRequirements("genomics", dataPattern, profile, hints)
	}
}

func BenchmarkResourceAnalyzer_FindOptimalInstance(b *testing.B) {
	analyzer := NewResourceAnalyzer()
	requirements := &ResourceRequirement{
		MinCPUs:     16,
		MinMemoryGB: 32,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.FindOptimalInstance(requirements)
	}
}
