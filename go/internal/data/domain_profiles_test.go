package data

import (
	"testing"
)

// TestDomainProfileManager tests the domain profile manager functionality
func TestDomainProfileManager(t *testing.T) {
	dpm := NewResearchDomainProfileManager()

	// Test that built-in profiles are loaded
	expectedDomains := []string{
		"genomics", "climate", "machine_learning", "astronomy",
		"geospatial", "chemistry", "physics", "materials_science",
	}

	for _, domain := range expectedDomains {
		profile, exists := dpm.GetProfile(domain)
		if !exists {
			t.Errorf("Expected domain profile '%s' to exist", domain)
			continue
		}

		if profile.Name != domain {
			t.Errorf("Expected profile name '%s', got '%s'", domain, profile.Name)
		}

		if profile.Description == "" {
			t.Errorf("Expected non-empty description for domain '%s'", domain)
		}
	}

	// Test getting all profiles
	allProfiles := dpm.GetAllProfiles()
	if len(allProfiles) < len(expectedDomains) {
		t.Errorf("Expected at least %d profiles, got %d", len(expectedDomains), len(allProfiles))
	}
}

// TestGenomicsProfile tests genomics-specific optimizations
func TestGenomicsProfile(t *testing.T) {
	dpm := NewResearchDomainProfileManager()
	profile, exists := dpm.GetProfile("genomics")

	if !exists {
		t.Fatal("Genomics profile should exist")
	}

	// Test file type hints
	fastqHint, exists := profile.FileTypeHints[".fastq"]
	if !exists {
		t.Error("Expected .fastq file type hint")
	} else {
		if fastqHint.CompressionRatio <= 1.0 {
			t.Error("FASTQ files should have compression ratio > 1.0")
		}

		if fastqHint.PreferredEngine == "" {
			t.Error("FASTQ files should have preferred engine")
		}
	}

	// Test transfer optimization
	if profile.TransferOptimization.OptimalConcurrency <= 0 {
		t.Error("Expected positive optimal concurrency")
	}

	if profile.TransferOptimization.OptimalPartSize == "" {
		t.Error("Expected non-empty optimal part size")
	}

	// Test bundling strategy
	if !profile.BundlingStrategy.EnableBundling {
		t.Error("Genomics should enable bundling for small files")
	}

	if profile.BundlingStrategy.MinFilesForBundle <= 0 {
		t.Error("Expected positive minimum files for bundle")
	}

	// Test security requirements
	if !profile.SecurityRequirements.EncryptionRequired {
		t.Error("Genomics data should require encryption")
	}

	// Test quality checks
	if len(profile.QualityChecks) == 0 {
		t.Error("Expected quality checks for genomics data")
	}
}

// TestClimateProfile tests climate science optimizations
func TestClimateProfile(t *testing.T) {
	dpm := NewResearchDomainProfileManager()
	profile, exists := dpm.GetProfile("climate")

	if !exists {
		t.Fatal("Climate profile should exist")
	}

	// Test NetCDF file hints
	ncHint, exists := profile.FileTypeHints[".nc"]
	if !exists {
		t.Error("Expected .nc file type hint for climate data")
	} else {
		if ncHint.AccessPattern == "" {
			t.Error("NetCDF files should have defined access pattern")
		}
	}

	// Climate data often doesn't require encryption (public data)
	if profile.SecurityRequirements.EncryptionRequired {
		// This might be acceptable, but let's check data classification
		if profile.SecurityRequirements.DataClassification != "public_research" {
			t.Error("Expected public_research classification for climate data")
		}
	}

	// Test bundling for time-series data
	if profile.BundlingStrategy.GroupingStrategy.Strategy != "time_based" {
		t.Error("Expected time-based grouping strategy for climate data")
	}
}

// TestMachineLearningProfile tests ML-specific optimizations
func TestMachineLearningProfile(t *testing.T) {
	dpm := NewResearchDomainProfileManager()
	profile, exists := dpm.GetProfile("machine_learning")

	if !exists {
		t.Fatal("Machine learning profile should exist")
	}

	// Test ML file types
	expectedTypes := []string{".pkl", ".h5", ".model"}
	for _, fileType := range expectedTypes {
		hint, exists := profile.FileTypeHints[fileType]
		if !exists {
			t.Errorf("Expected %s file type hint for ML", fileType)
		} else {
			if hint.TypicalSize == "" {
				t.Errorf("Expected typical size for %s files", fileType)
			}
		}
	}

	// ML should have high concurrency for large dataset transfers
	if profile.TransferOptimization.OptimalConcurrency < 20 {
		t.Error("Expected high concurrency for ML workloads")
	}

	// ML models should be encrypted (proprietary)
	if !profile.SecurityRequirements.EncryptionRequired {
		t.Error("ML models should require encryption")
	}
}

// TestDomainOptimizationApplication tests applying domain optimizations to workflows
func TestDomainOptimizationApplication(t *testing.T) {
	// Create workflow engine with domain profiles
	engine := NewWorkflowEngine(nil)

	// Create test project config with genomics domain
	projectConfig := &ProjectConfig{
		Project: ProjectInfo{
			Domain: "genomics",
		},
		DataProfiles: map[string]DataProfile{
			"test_data": {
				Name: "Test Genomics Data",
				Path: "/test/genomics",
			},
		},
		Destinations: map[string]Destination{
			"test_dest": {
				Name: "Test Destination",
				URI:  "s3://test-bucket/",
			},
		},
	}

	// Create test workflow with minimal configuration
	workflow := &Workflow{
		Name:          "test_workflow",
		Source:        "test_data",
		Destination:   "test_dest",
		Engine:        "auto", // Should be optimized
		Enabled:       true,
		Configuration: WorkflowConfiguration{
			// Leave empty to test auto-optimization
		},
	}

	// Apply domain optimizations
	err := engine.ApplyDomainOptimizations(workflow, projectConfig)
	if err != nil {
		t.Fatalf("Failed to apply domain optimizations: %v", err)
	}

	// Check that optimizations were applied
	if workflow.Configuration.Concurrency == 0 {
		t.Error("Expected concurrency to be set from domain profile")
	}

	if workflow.Configuration.PartSize == "" {
		t.Error("Expected part size to be set from domain profile")
	}

	if workflow.Engine == "auto" {
		t.Error("Expected engine to be selected from domain profile")
	}

	// Check that bundling step was added
	if len(workflow.PreProcessing) == 0 {
		t.Error("Expected bundling preprocessing step to be added")
	}

	// Check that quality checks were added
	if len(workflow.PostProcessing) == 0 {
		t.Error("Expected quality check postprocessing steps to be added")
	}
}

// TestCustomDomainProfile tests adding custom domain profiles
func TestCustomDomainProfile(t *testing.T) {
	dpm := NewResearchDomainProfileManager()

	// Create custom domain profile
	customProfile := &ResearchDomainProfile{
		Name:        "custom_domain",
		Description: "Custom research domain for testing",
		TransferOptimization: TransferOptimization{
			OptimalConcurrency: 42,
			OptimalPartSize:    "256MB",
		},
		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:   false,
			TargetBundleSize: "1GB",
		},
	}

	// Add custom profile
	dpm.AddProfile("custom_domain", customProfile)

	// Retrieve and verify
	retrieved, exists := dpm.GetProfile("custom_domain")
	if !exists {
		t.Fatal("Custom profile should exist after adding")
	}

	if retrieved.TransferOptimization.OptimalConcurrency != 42 {
		t.Error("Custom profile settings should be preserved")
	}

	if retrieved.BundlingStrategy.EnableBundling {
		t.Error("Custom bundling setting should be preserved")
	}
}

// TestDomainProfileIntegration tests integration with existing systems
func TestDomainProfileIntegration(t *testing.T) {
	// This test verifies that domain profiles integrate properly with
	// the existing pattern analyzer and recommendation engine

	dpm := NewResearchDomainProfileManager()

	// Test that we can get profiles for detected domains
	detectedDomains := []string{"genomics", "climate", "machine_learning"}

	for _, domain := range detectedDomains {
		profile, exists := dpm.GetProfile(domain)
		if !exists {
			t.Errorf("Should have profile for detected domain: %s", domain)
			continue
		}

		// Verify profile has essential components
		if len(profile.FileTypeHints) == 0 {
			t.Errorf("Domain %s should have file type hints", domain)
		}

		if profile.TransferOptimization.OptimalConcurrency <= 0 {
			t.Errorf("Domain %s should have valid transfer optimization", domain)
		}
	}
}
