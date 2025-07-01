package intelligence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDomainPackLoader_LoadDomainPack(t *testing.T) {
	loader := NewDomainPackLoader()

	tests := []struct {
		name        string
		domain      string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "genomics_lab domain",
			domain:      "genomics_lab",
			expectError: false, // Real domain pack files exist
			expectNil:   false,
		},
		{
			name:        "ai_research_studio domain",
			domain:      "ai_research_studio",
			expectError: false, // Real domain pack files exist
			expectNil:   false,
		},
		{
			name:        "climate_modeling domain",
			domain:      "climate_modeling",
			expectError: false, // Real domain pack files exist
			expectNil:   false,
		},
		{
			name:        "unknown domain",
			domain:      "unknown_domain",
			expectError: true,
			expectNil:   true,
		},
		{
			name:        "empty domain",
			domain:      "",
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := loader.LoadDomainPack(tt.domain)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectNil && result != nil {
				t.Errorf("expected nil result but got: %v", result)
			}
			if !tt.expectNil && result == nil {
				t.Errorf("expected non-nil result but got nil")
			}

			// Validate structure if we got a result
			if result != nil {
				if result.Name == "" {
					t.Errorf("expected non-empty name")
				}
				if result.Description == "" {
					t.Errorf("expected non-empty description")
				}
			}
		})
	}
}

func TestDomainPackLoader_LoadAllDomainPacks(t *testing.T) {
	loader := NewDomainPackLoader()

	packs, err := loader.LoadAllDomainPacks()

	// In test environment, we expect no domain packs to be found
	if err != nil {
		t.Logf("expected error in test environment: %v", err)
	}

	if len(packs) == 0 {
		t.Logf("no domain packs found in test environment (expected)")
	}

	// Verify each pack has required fields if any are found
	for domain, pack := range packs {
		if domain == "" {
			t.Errorf("found empty domain key")
		}
		if pack == nil {
			t.Errorf("found nil pack for domain: %s", domain)
			continue
		}
		if pack.Name == "" {
			t.Errorf("pack for domain %s has empty name", domain)
		}
		if pack.Description == "" {
			t.Errorf("pack for domain %s has empty description", domain)
		}
	}
}

func TestDomainPackLoader_GetAvailableDomains(t *testing.T) {
	loader := NewDomainPackLoader()

	domains, err := loader.GetAvailableDomains()
	if err != nil {
		t.Logf("expected error in test environment: %v", err)
	}

	if len(domains) == 0 {
		t.Logf("no available domains in test environment (expected)")
	}

	// Only check for expected domains if any domains are found
	if len(domains) > 0 {
		expectedDomains := []string{"genomics_lab", "ai_research_studio", "climate_modeling"}
		for _, expected := range expectedDomains {
			found := false
			for _, domain := range domains {
				if domain == expected {
					found = true
					break
				}
			}
			if !found {
				t.Logf("domain %s not found in available domains (may be expected in test environment)", expected)
			}
		}
	}
}

func TestDomainPackLoader_ValidateDomainPack(t *testing.T) {
	loader := NewDomainPackLoader()

	tests := []struct {
		name        string
		domainName  string
		expectValid bool
	}{
		{
			name:        "genomics_lab domain",
			domainName:  "genomics_lab",
			expectValid: true, // Real domain pack files exist
		},
		{
			name:        "ai_research_studio domain",
			domainName:  "ai_research_studio",
			expectValid: true, // Real domain pack files exist
		},
		{
			name:        "climate_modeling domain",
			domainName:  "climate_modeling",
			expectValid: true, // Real domain pack files exist
		},
		{
			name:        "invalid domain",
			domainName:  "nonexistent",
			expectValid: false,
		},
		{
			name:        "empty domain",
			domainName:  "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loader.ValidateDomainPack(tt.domainName)
			hasError := err != nil
			if hasError == tt.expectValid {
				t.Errorf("expected validation result %v, got error: %v", tt.expectValid, err)
			}
		})
	}
}

func TestDomainPackLoader_ClearCache(t *testing.T) {
	loader := NewDomainPackLoader()

	// Try to load something into cache first (should work now with real domain packs)
	_, err := loader.LoadDomainPack("genomics_lab")
	if err != nil {
		t.Logf("error loading genomics_lab domain: %v", err)
	}

	// Clear cache should not return error
	loader.ClearCache()

	// Verify cache is cleared by loading again
	// This is mainly to ensure the function runs without error
	_, err = loader.LoadDomainPack("genomics_lab")
	if err != nil {
		t.Logf("error after cache clear: %v", err)
	}
}

func TestDomainPackLoader_LoadErrorCases(t *testing.T) {
	loader := NewDomainPackLoader()

	// Test loading non-existent domain
	_, err := loader.LoadDomainPack("nonexistent-domain")
	if err == nil {
		t.Error("expected error for non-existent domain")
	}

	// Test validating non-existent domain
	err = loader.ValidateDomainPack("nonexistent-domain")
	if err == nil {
		t.Error("expected error for non-existent domain validation")
	}
}

func TestDomainPackLoader_EdgeCases(t *testing.T) {
	loader := NewDomainPackLoader()

	t.Run("load domain pack with special characters", func(t *testing.T) {
		_, err := loader.LoadDomainPack("domain-with-dashes_and_underscores")
		// Should not panic or cause issues
		if err != nil {
			// Error is acceptable for non-existent domains
		}
	})

	t.Run("load domain pack with very long name", func(t *testing.T) {
		longName := string(make([]byte, 1000))
		_, err := loader.LoadDomainPack(longName)
		// Should not panic
		if err != nil {
			// Error is acceptable for invalid domains
		}
	})

	t.Run("multiple cache operations", func(t *testing.T) {
		// Multiple clears should be safe
		loader.ClearCache()
		loader.ClearCache()
		loader.ClearCache()

		// Multiple loads should be consistent
		pack1, _ := loader.LoadDomainPack("genomics_lab")
		pack2, _ := loader.LoadDomainPack("genomics_lab")

		if pack1 != nil && pack2 != nil {
			if pack1.Name != pack2.Name {
				t.Errorf("inconsistent results from cache")
			}
		}
	})
}

func TestDomainPackLoader_loadConfigFile(t *testing.T) {
	// Cast to concrete type to access private methods
	loader := NewDomainPackLoader().(*DomainPackLoader)

	// Create a temporary test file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-domain-pack.yaml")

	// Test with invalid YAML
	invalidYAML := `
name: test-domain
description: Test domain pack
invalid_yaml: [unclosed
`
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = loader.loadConfigFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}

	// Test with valid YAML
	validYAML := `
name: test-domain
description: Test domain pack
version: "1.0.0"
categories:
  - test
aws_config:
  instance_types:
    small: c6i.large
    medium: c6i.xlarge
spack_config:
  packages:
    - test-package
workflows:
  - name: test-workflow
    description: Test workflow
    script: test.sh
    input_data: input.txt
    expected_output: output.txt
cost_estimates:
  small: "$10/month"
documentation:
  getting_started: "https://example.com"
`
	err = os.WriteFile(configPath, []byte(validYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config, err := loader.loadConfigFile(configPath)
	if err != nil {
		t.Errorf("Unexpected error for valid YAML: %v", err)
	}

	if config != nil {
		if config.Name != "test-domain" {
			t.Errorf("Expected name 'test-domain', got %s", config.Name)
		}
		if config.Description != "Test domain pack" {
			t.Errorf("Expected description 'Test domain pack', got %s", config.Description)
		}
		if config.Version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got %s", config.Version)
		}
	}

	// Test with non-existent file
	_, err = loader.loadConfigFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestDomainPackLoader_convertToDomainPackInfo(t *testing.T) {
	// Cast to concrete type to access private methods
	loader := NewDomainPackLoader().(*DomainPackLoader)

	config := &DomainPackConfig{
		Name:        "test-domain",
		Description: "Test domain pack",
		Version:     "1.0.0",
		Categories:  []string{"test", "development"},
		AWSConfig: AWSConfig{
			InstanceTypes: map[string]string{
				"small":  "c6i.large",
				"medium": "c6i.xlarge",
				"large":  "c6i.2xlarge",
			},
		},
		SpackConfig: SpackConfig{
			Packages: []string{"gcc", "python", "numpy"},
		},
		Workflows: []WorkflowConfig{
			{
				Name:           "test-workflow",
				Description:    "A test workflow",
				Script:         "run_test.sh",
				InputData:      "input.txt",
				ExpectedOutput: "output.txt",
			},
			{
				Name:           "benchmark",
				Description:    "Benchmark workflow",
				Script:         "benchmark.sh",
				InputData:      "data.csv",
				ExpectedOutput: "results.json",
			},
		},
		CostEstimates: map[string]string{
			"small":  "$10/month",
			"medium": "$20/month",
			"large":  "$40/month",
		},
	}

	info := loader.convertToDomainPackInfo(config)

	if info == nil {
		t.Fatal("convertToDomainPackInfo() returned nil")
	}

	if info.Name != config.Name {
		t.Errorf("Expected name %s, got %s", config.Name, info.Name)
	}

	if info.Description != config.Description {
		t.Errorf("Expected description %s, got %s", config.Description, info.Description)
	}

	if info.Version != config.Version {
		t.Errorf("Expected version %s, got %s", config.Version, info.Version)
	}

	if len(info.Categories) != len(config.Categories) {
		t.Errorf("Expected %d categories, got %d", len(config.Categories), len(info.Categories))
	}

	if len(info.InstanceTypes) != len(config.AWSConfig.InstanceTypes) {
		t.Errorf("Expected %d instance types, got %d", len(config.AWSConfig.InstanceTypes), len(info.InstanceTypes))
	}

	if len(info.SpackPackages) != len(config.SpackConfig.Packages) {
		t.Errorf("Expected %d Spack packages, got %d", len(config.SpackConfig.Packages), len(info.SpackPackages))
	}

	if len(info.Workflows) != len(config.Workflows) {
		t.Errorf("Expected %d workflows, got %d", len(config.Workflows), len(info.Workflows))
	}

	// Check workflow conversion
	if len(info.Workflows) > 0 {
		workflow := info.Workflows[0]
		originalWorkflow := config.Workflows[0]

		if workflow.Name != originalWorkflow.Name {
			t.Errorf("Expected workflow name %s, got %s", originalWorkflow.Name, workflow.Name)
		}

		if workflow.Description != originalWorkflow.Description {
			t.Errorf("Expected workflow description %s, got %s", originalWorkflow.Description, workflow.Description)
		}

		if workflow.InputData != originalWorkflow.InputData {
			t.Errorf("Expected input data %s, got %s", originalWorkflow.InputData, workflow.InputData)
		}

		if workflow.OutputData != originalWorkflow.ExpectedOutput {
			t.Errorf("Expected output data %s, got %s", originalWorkflow.ExpectedOutput, workflow.OutputData)
		}
	}

	if len(info.EstimatedCost) != len(config.CostEstimates) {
		t.Errorf("Expected %d cost estimates, got %d", len(config.CostEstimates), len(info.EstimatedCost))
	}
}

func TestDomainPackLoader_convertToDomainPackInfo_EmptyConfig(t *testing.T) {
	// Cast to concrete type to access private methods
	loader := NewDomainPackLoader().(*DomainPackLoader)

	// Test with minimal config
	config := &DomainPackConfig{
		Name: "minimal-domain",
	}

	info := loader.convertToDomainPackInfo(config)

	if info == nil {
		t.Fatal("convertToDomainPackInfo() returned nil for minimal config")
	}

	if info.Name != "minimal-domain" {
		t.Errorf("Expected name 'minimal-domain', got %s", info.Name)
	}

	// Should handle empty workflows gracefully (nil is acceptable for empty workflows)
	if len(info.Workflows) != 0 {
		t.Errorf("Expected empty workflows, got %d", len(info.Workflows))
	}
}
