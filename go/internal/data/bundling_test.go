package data

import (
	"context"
	"testing"
)

// TestBundlingEngineIntegration tests the basic bundling engine functionality
func TestBundlingEngineIntegration(t *testing.T) {
	// Create a bundling engine with test configuration
	config := &BundlingConfig{
		Enabled:             true,
		AutoBundle:          true,
		BundleThreshold:     "1MB",
		MinFilesForBundling: 10,
		SuitcaseConfig: &SuitcaseConfig{
			TargetBundleSize:   "50MB",
			MaxFilesPerBundle:  500,
			PreserveMetadata:   true,
			CompressionLevel:   6,
			OutputFormat:       "tar.gz",
			DomainOptimization: "general",
		},
	}

	engine := NewBundlingEngine(config)

	// Test engine interface compliance
	if engine.GetName() != "bundling" {
		t.Errorf("Expected engine name 'bundling', got '%s'", engine.GetName())
	}

	if engine.GetType() != "preprocessing" {
		t.Errorf("Expected engine type 'preprocessing', got '%s'", engine.GetType())
	}

	// Test capabilities
	capabilities := engine.GetCapabilities()
	if !capabilities.SupportsCompression {
		t.Error("Expected bundling engine to support compression")
	}

	if !capabilities.SupportsParallel {
		t.Error("Expected bundling engine to support parallel processing")
	}

	// Test validation
	if err := engine.Validate(); err != nil {
		// This is expected to fail in CI since Suitcase won't be installed
		t.Logf("Validation failed as expected (Suitcase not installed): %v", err)
	}
}

// TestBundlingRecommendation tests the bundling recommendation logic
func TestBundlingRecommendation(t *testing.T) {
	config := &BundlingConfig{
		Enabled:             true,
		MinFilesForBundling: 100,
	}

	engine := NewBundlingEngine(config)

	// Create test data pattern with many small files
	pattern := &DataPattern{
		TotalFiles: 5000,
		TotalSize:  500 * 1024 * 1024, // 500MB
		FileSizes: FileSizeAnalysis{
			SmallFiles: SmallFileAnalysis{
				CountUnder1MB:    4500, // 90% small files
				PercentageSmall:  90.0,
				PotentialSavings: 150.0, // $150/month potential savings
			},
		},
		DomainHints: DomainAnalysis{
			DetectedDomains: []string{"genomics"},
		},
	}

	recommendation, err := engine.ShouldBundle(context.Background(), pattern)
	if err != nil {
		t.Fatalf("Failed to get bundling recommendation: %v", err)
	}

	if !recommendation.Recommended {
		t.Error("Expected bundling to be recommended for dataset with many small files")
	}

	if recommendation.Confidence < 0.8 {
		t.Errorf("Expected high confidence (>0.8), got %f", recommendation.Confidence)
	}

	if recommendation.EstimatedSavings != 150.0 {
		t.Errorf("Expected estimated savings of $150, got $%f", recommendation.EstimatedSavings)
	}

	// Test that genomics domain is detected
	if genomicsHints, exists := recommendation.DomainHints["genomics"]; exists {
		hints := genomicsHints.(map[string]interface{})
		if !hints["group_by_type"].(bool) {
			t.Error("Expected genomics optimization to recommend grouping by type")
		}
	}
}

// TestWorkflowGeneration tests workflow generation from bundling analysis
func TestWorkflowGeneration(t *testing.T) {
	engine := NewBundlingEngine(nil) // Use default config

	pattern := &DataPattern{
		TotalFiles: 2000,
		FileSizes: FileSizeAnalysis{
			SmallFiles: SmallFileAnalysis{
				CountUnder1MB:   1800,
				PercentageSmall: 90.0,
			},
		},
		DomainHints: DomainAnalysis{
			DetectedDomains: []string{"genomics"},
		},
	}

	recommendation := &BundlingRecommendation{
		Recommended:   true,
		Confidence:    0.9,
		Complexity:    "moderate",
		EstimatedTime: "2 hours",
	}

	workflow, err := engine.CreateWorkflowFromBundling(pattern, recommendation)
	if err != nil {
		t.Fatalf("Failed to create workflow: %v", err)
	}

	if workflow.Name != "auto_bundle_upload" {
		t.Errorf("Expected workflow name 'auto_bundle_upload', got '%s'", workflow.Name)
	}

	if workflow.Engine != "bundling" {
		t.Errorf("Expected workflow engine 'bundling', got '%s'", workflow.Engine)
	}

	// Check for bundling preprocessing step
	found := false
	for _, step := range workflow.PreProcessing {
		if step.Type == "bundle" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected bundling preprocessing step in workflow")
	}
}
