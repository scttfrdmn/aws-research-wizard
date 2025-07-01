package data

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
)

// BundlingEngine integrates with the transfer engine framework to provide file bundling capabilities
type BundlingEngine struct {
	suitcase *SuitcaseEngine
	config   *BundlingConfig
}

// BundlingTransferRequest represents an internal transfer request for bundling
type BundlingTransferRequest struct {
	SourcePath      string                 `json:"source_path"`
	DestinationPath string                 `json:"destination_path"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// BundlingConfig contains configuration for the bundling engine
type BundlingConfig struct {
	// Integration settings
	Enabled             bool   `json:"enabled"`
	AutoBundle          bool   `json:"auto_bundle"`            // Automatically bundle small files
	BundleThreshold     string `json:"bundle_threshold"`       // Bundle files smaller than this
	MinFilesForBundling int    `json:"min_files_for_bundling"` // Minimum files to trigger bundling

	// Suitcase configuration
	SuitcaseConfig *SuitcaseConfig `json:"suitcase_config"`

	// Integration with transfer engines
	ChainWithUpload     bool   `json:"chain_with_upload"`     // Automatically upload after bundling
	PreferredUploadTool string `json:"preferred_upload_tool"` // "s5cmd", "rclone", "aws-cli"
	CleanupOriginals    bool   `json:"cleanup_originals"`     // Remove original files after bundling

	// Research domain settings
	DomainOptimizations map[string]interface{} `json:"domain_optimizations"`
}

// NewBundlingEngine creates a new bundling engine with default configuration
func NewBundlingEngine(config *BundlingConfig) *BundlingEngine {
	if config == nil {
		config = &BundlingConfig{
			Enabled:             true,
			AutoBundle:          true,
			BundleThreshold:     "1MB",
			MinFilesForBundling: 100,
			ChainWithUpload:     true,
			PreferredUploadTool: "s5cmd",
			CleanupOriginals:    false, // Conservative default
		}
	}

	// Create default Suitcase config if not provided
	if config.SuitcaseConfig == nil {
		config.SuitcaseConfig = &SuitcaseConfig{
			TargetBundleSize:   "100MB",
			MaxFilesPerBundle:  1000,
			PreserveMetadata:   true,
			CompressionLevel:   6,
			OutputFormat:       "tar.gz",
			Parallel:           true,
			WorkerCount:        4,
			DomainOptimization: "general",
		}
	}

	suitcase := NewSuitcaseEngine(config.SuitcaseConfig)

	return &BundlingEngine{
		suitcase: suitcase,
		config:   config,
	}
}

// GetName returns the name of this transfer engine
func (be *BundlingEngine) GetName() string {
	return "bundling"
}

// GetType returns the type/category of the engine
func (be *BundlingEngine) GetType() string {
	return "preprocessing"
}

// IsAvailable checks if the bundling engine is available
func (be *BundlingEngine) IsAvailable(ctx context.Context) error {
	if !be.config.Enabled {
		return fmt.Errorf("bundling engine is disabled")
	}

	return be.suitcase.IsAvailable(ctx)
}

// GetCapabilities returns the capabilities of the bundling engine
func (be *BundlingEngine) GetCapabilities() EngineCapabilities {
	return EngineCapabilities{
		Protocols:              []string{"local"},
		SupportsParallel:       true,
		SupportsProgress:       true,
		SupportsResume:         false, // Bundling doesn't support resume
		SupportsCompression:    true,
		SupportsEncryption:     false,
		SupportsValidation:     true,
		SupportsBandwidthLimit: false,
		SupportsRetry:          true,
		OptimalFileSizeMin:     0,
		OptimalFileSizeMax:     1024 * 1024, // 1MB - optimal for files smaller than this
		MaxConcurrency:         be.config.SuitcaseConfig.WorkerCount,
		CloudOptimized:         []string{}, // Not cloud-specific
	}
}

// Upload performs bundling operation (adapts bundling to transfer interface)
func (be *BundlingEngine) Upload(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// Convert TransferRequest to our internal format
	bundlingReq := &BundlingTransferRequest{
		SourcePath:      req.Source,
		DestinationPath: req.Destination,
		Metadata:        make(map[string]interface{}),
	}

	// Copy relevant metadata
	if req.Options.ToolSpecific != nil {
		for k, v := range req.Options.ToolSpecific {
			bundlingReq.Metadata[k] = v
		}
	}

	startTime := time.Now()

	// Perform bundling
	bundleResult, err := be.ProcessForBundling(ctx, bundlingReq)
	if err != nil {
		return &TransferResult{
			TransferID:  req.ID,
			Engine:      be.GetName(),
			Source:      req.Source,
			Destination: req.Destination,
			Success:     false,
			Error:       err,
			StartTime:   startTime,
			EndTime:     time.Now(),
			Duration:    time.Since(startTime),
		}, err
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// Convert bundling result to transfer result
	return &TransferResult{
		TransferID:       req.ID,
		Engine:           be.GetName(),
		Source:           req.Source,
		Destination:      bundleResult.OutputPath,
		Success:          true,
		BytesTransferred: bundleResult.BundledSize,
		FilesTransferred: int(bundleResult.BundledFileCount),
		StartTime:        startTime,
		EndTime:          endTime,
		Duration:         duration,
		AverageSpeed:     int64(float64(bundleResult.BundledSize) / duration.Seconds()),
		Metadata: map[string]interface{}{
			"original_file_count": bundleResult.OriginalFileCount,
			"bundled_file_count":  bundleResult.BundledFileCount,
			"compression_ratio":   bundleResult.CompressionRatio,
			"space_efficiency":    bundleResult.SpaceEfficiency,
			"cost_savings":        bundleResult.CostSavings,
			"bundle_paths":        bundleResult.BundlePaths,
		},
	}, nil
}

// Download is not applicable for bundling engine
func (be *BundlingEngine) Download(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	return nil, fmt.Errorf("download operation not supported by bundling engine")
}

// Sync is not applicable for bundling engine
func (be *BundlingEngine) Sync(ctx context.Context, req *SyncRequest) (*TransferResult, error) {
	return nil, fmt.Errorf("sync operation not supported by bundling engine")
}

// GetProgress returns current progress for bundling operations
func (be *BundlingEngine) GetProgress(ctx context.Context, transferID string) (*TransferProgress, error) {
	// This would need to be implemented with transfer ID tracking
	// For now, return an error
	return nil, fmt.Errorf("progress tracking by transfer ID not yet implemented")
}

// Cancel cancels an active bundling operation
func (be *BundlingEngine) Cancel(ctx context.Context, transferID string) error {
	// This would need to be implemented with transfer ID tracking and cancellation
	return fmt.Errorf("cancellation by transfer ID not yet implemented")
}

// Validate validates the bundling engine configuration
func (be *BundlingEngine) Validate() error {
	if be.config == nil {
		return fmt.Errorf("bundling configuration is required")
	}

	if be.config.SuitcaseConfig == nil {
		return fmt.Errorf("suitcase configuration is required")
	}

	// Validate suitcase availability
	return be.suitcase.IsAvailable(context.Background())
}

// ShouldBundle analyzes a dataset and determines if bundling is recommended
func (be *BundlingEngine) ShouldBundle(ctx context.Context, pattern *DataPattern) (*BundlingRecommendation, error) {
	recommendation := &BundlingRecommendation{
		Recommended:      false,
		Confidence:       0.0,
		Reasoning:        []string{},
		EstimatedSavings: 0.0,
		EstimatedTime:    "unknown",
		Complexity:       "moderate",
		Prerequisites:    []string{"suitcase", "python"},
	}

	// Check if bundling is enabled
	if !be.config.Enabled {
		recommendation.Reasoning = append(recommendation.Reasoning, "Bundling is disabled in configuration")
		return recommendation, nil
	}

	// Analyze small file patterns
	smallFiles := pattern.FileSizes.SmallFiles

	// Primary recommendation criteria: lots of small files
	if smallFiles.CountUnder1MB >= int64(be.config.MinFilesForBundling) {
		recommendation.Recommended = true
		recommendation.Confidence = 0.9
		recommendation.Reasoning = append(recommendation.Reasoning,
			fmt.Sprintf("Found %d files under 1MB (%.1f%% of total)",
				smallFiles.CountUnder1MB, smallFiles.PercentageSmall))

		// Calculate estimated savings
		recommendation.EstimatedSavings = smallFiles.PotentialSavings

		// Adjust confidence based on small file percentage
		if smallFiles.PercentageSmall > 70 {
			recommendation.Confidence = 0.95
			recommendation.Reasoning = append(recommendation.Reasoning, "Very high percentage of small files")
		} else if smallFiles.PercentageSmall > 50 {
			recommendation.Confidence = 0.90
		} else if smallFiles.PercentageSmall > 30 {
			recommendation.Confidence = 0.80
		} else {
			recommendation.Confidence = 0.70
		}

		// Estimate bundling time based on file count and sizes
		recommendation.EstimatedTime = be.estimateBundlingTime(pattern)

		// Adjust complexity based on domain and file patterns
		recommendation.Complexity = be.assessComplexity(pattern)

		// Add domain-specific recommendations
		be.addDomainSpecificRecommendations(recommendation, pattern)

	} else {
		recommendation.Reasoning = append(recommendation.Reasoning,
			fmt.Sprintf("Only %d small files found, below threshold of %d",
				smallFiles.CountUnder1MB, be.config.MinFilesForBundling))
	}

	// Additional factors that increase recommendation confidence
	if pattern.Efficiency.EstimatedRequestCosts > pattern.Efficiency.EstimatedStorageCosts {
		recommendation.Confidence += 0.05
		recommendation.Reasoning = append(recommendation.Reasoning, "Request costs dominate storage costs")
	}

	// Research domain factors
	for _, domain := range pattern.DomainHints.DetectedDomains {
		switch domain {
		case "genomics":
			recommendation.Confidence += 0.03
			recommendation.Reasoning = append(recommendation.Reasoning, "Genomics data often benefits from bundling")
		case "climate":
			recommendation.Confidence += 0.02
			recommendation.Reasoning = append(recommendation.Reasoning, "Climate data with many files benefits from bundling")
		}
	}

	// Cap confidence at 1.0
	if recommendation.Confidence > 1.0 {
		recommendation.Confidence = 1.0
	}

	return recommendation, nil
}

// BundlingRecommendation contains the recommendation for bundling a dataset
type BundlingRecommendation struct {
	Recommended        bool                   `json:"recommended"`
	Confidence         float64                `json:"confidence"` // 0.0 to 1.0
	Reasoning          []string               `json:"reasoning"`
	EstimatedSavings   float64                `json:"estimated_savings_monthly"`
	EstimatedTime      string                 `json:"estimated_time"`
	Complexity         string                 `json:"complexity"` // "simple", "moderate", "complex"
	Prerequisites      []string               `json:"prerequisites"`
	DomainHints        map[string]interface{} `json:"domain_hints,omitempty"`
	AlternativeOptions []string               `json:"alternative_options,omitempty"`
}

// ProcessForBundling processes a dataset for bundling with the transfer engine framework
func (be *BundlingEngine) ProcessForBundling(ctx context.Context, req *BundlingTransferRequest) (*BundlingResult, error) {
	// Validate request
	if req.SourcePath == "" {
		return nil, fmt.Errorf("source path is required")
	}

	// Set up bundling configuration based on request
	be.configureBundlingForRequest(req)

	// Execute bundling
	bundleResult, err := be.suitcase.BundleFiles(ctx, req.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("bundling failed: %w", err)
	}

	// Convert to our result format
	result := &BundlingResult{
		BundleResult: bundleResult,
		TransferMetadata: TransferMetadata{
			Engine:        "bundling",
			StartTime:     bundleResult.ProcessingTime,
			TransferSpeed: "N/A", // Bundling doesn't have traditional transfer speed
			Success:       true,
		},
		NextSteps: make([]NextStep, 0),
	}

	// Add next steps if chaining is enabled
	if be.config.ChainWithUpload {
		result.NextSteps = append(result.NextSteps, NextStep{
			Action:      "upload",
			Engine:      be.config.PreferredUploadTool,
			SourcePath:  bundleResult.OutputPath,
			Description: fmt.Sprintf("Upload %d bundles using %s", len(bundleResult.BundlePaths), be.config.PreferredUploadTool),
			Priority:    "high",
		})
	}

	// Add cleanup step if configured
	if be.config.CleanupOriginals {
		result.NextSteps = append(result.NextSteps, NextStep{
			Action:      "cleanup",
			Engine:      "filesystem",
			SourcePath:  req.SourcePath,
			Description: "Remove original files after successful bundling and upload",
			Priority:    "low",
			Conditions:  []string{"upload_successful"},
		})
	}

	return result, nil
}

// BundlingResult extends the standard bundle result with transfer engine integration
type BundlingResult struct {
	*BundleResult
	TransferMetadata TransferMetadata `json:"transfer_metadata"`
	NextSteps        []NextStep       `json:"next_steps"`
	Recommendations  []string         `json:"recommendations,omitempty"`
}

// NextStep describes a recommended next action after bundling
type NextStep struct {
	Action      string   `json:"action"`               // "upload", "cleanup", "verify"
	Engine      string   `json:"engine"`               // Which engine to use
	SourcePath  string   `json:"source_path"`          // Path for the action
	Description string   `json:"description"`          // Human-readable description
	Priority    string   `json:"priority"`             // "high", "medium", "low"
	Conditions  []string `json:"conditions,omitempty"` // Conditions that must be met
}

// TransferMetadata provides metadata compatible with the transfer engine framework
type TransferMetadata struct {
	Engine        string        `json:"engine"`
	StartTime     time.Duration `json:"start_time"`
	TransferSpeed string        `json:"transfer_speed"`
	Success       bool          `json:"success"`
	ErrorMessage  string        `json:"error_message,omitempty"`
}

// CreateWorkflowFromBundling creates a complete workflow configuration from bundling analysis
func (be *BundlingEngine) CreateWorkflowFromBundling(pattern *DataPattern, recommendation *BundlingRecommendation) (*Workflow, error) {
	if !recommendation.Recommended {
		return nil, fmt.Errorf("bundling not recommended for this dataset")
	}

	workflow := &Workflow{
		Name:        "auto_bundle_upload",
		Description: "Automatically bundle small files and upload to optimized storage",
		Source:      "small_files_dataset",
		Destination: "optimized_storage",
		Engine:      "bundling",
		Triggers:    []string{"manual"},
		Enabled:     true,
	}

	// Add preprocessing steps
	if recommendation.Complexity == "complex" {
		workflow.PreProcessing = append(workflow.PreProcessing, ProcessingStep{
			Name: "analyze_patterns",
			Type: "analyze",
			Parameters: map[string]string{
				"deep_analysis": "true",
				"domain_hints":  fmt.Sprintf("%v", pattern.DomainHints.DetectedDomains),
			},
		})
	}

	// Add bundling step
	workflow.PreProcessing = append(workflow.PreProcessing, ProcessingStep{
		Name: "bundle_small_files",
		Type: "bundle",
		Parameters: map[string]string{
			"tool":                be.suitcase.config.OutputFormat,
			"target_size":         be.suitcase.config.TargetBundleSize,
			"compression_level":   fmt.Sprintf("%d", be.suitcase.config.CompressionLevel),
			"preserve_metadata":   fmt.Sprintf("%t", be.suitcase.config.PreserveMetadata),
			"domain_optimization": be.suitcase.config.DomainOptimization,
		},
	})

	// Configure optimal settings based on bundling results
	workflow.Configuration = WorkflowConfiguration{
		Concurrency:     be.suitcase.config.WorkerCount,
		RetryAttempts:   3,
		Timeout:         "2h", // Bundling can take time
		Checksum:        true,
		OverwritePolicy: "if_newer",
		FailurePolicy:   "stop",
	}

	// Add post-processing if cleanup is enabled
	if be.config.CleanupOriginals {
		workflow.PostProcessing = append(workflow.PostProcessing, ProcessingStep{
			Name: "cleanup_originals",
			Type: "cleanup",
			Parameters: map[string]string{
				"condition": "upload_successful",
				"backup":    "false",
			},
			OnFailure: "continue", // Don't fail workflow if cleanup fails
		})
	}

	return workflow, nil
}

// Helper methods

func (be *BundlingEngine) configureBundlingForRequest(req *BundlingTransferRequest) {
	// Configure output directory based on request destination
	if req.DestinationPath != "" {
		bundleDir := filepath.Join(filepath.Dir(req.SourcePath), "bundled_for_upload")
		be.suitcase.config.OutputDirectory = bundleDir
	}

	// Apply domain-specific optimizations if available
	if req.Metadata != nil {
		if domain, exists := req.Metadata["domain"]; exists {
			if domainStr, ok := domain.(string); ok {
				be.suitcase.config.DomainOptimization = domainStr
			}
		}

		// Apply custom metadata
		if customMeta, exists := req.Metadata["custom_metadata"]; exists {
			if metaMap, ok := customMeta.(map[string]string); ok {
				be.suitcase.config.CustomMetadata = metaMap
			}
		}
	}
}

func (be *BundlingEngine) estimateBundlingTime(pattern *DataPattern) string {
	// Estimate based on file count and total size
	totalFiles := pattern.TotalFiles
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)

	// Base time estimate: 1 second per 100 files + 1 second per GB
	estimatedSeconds := float64(totalFiles)/100 + totalSizeGB

	// Add compression overhead
	estimatedSeconds *= 1.5

	// Adjust for parallelism
	estimatedSeconds /= float64(be.suitcase.config.WorkerCount)

	if estimatedSeconds < 60 {
		return fmt.Sprintf("%.0f seconds", estimatedSeconds)
	} else if estimatedSeconds < 3600 {
		return fmt.Sprintf("%.0f minutes", estimatedSeconds/60)
	} else {
		return fmt.Sprintf("%.1f hours", estimatedSeconds/3600)
	}
}

func (be *BundlingEngine) assessComplexity(pattern *DataPattern) string {
	complexityScore := 0

	// File count factor
	if pattern.TotalFiles > 100000 {
		complexityScore += 2
	} else if pattern.TotalFiles > 10000 {
		complexityScore += 1
	}

	// Directory depth factor
	if pattern.DirectoryDepth.MaxDepth > 10 {
		complexityScore += 2
	} else if pattern.DirectoryDepth.MaxDepth > 5 {
		complexityScore += 1
	}

	// File type diversity factor
	if len(pattern.FileTypes) > 20 {
		complexityScore += 1
	}

	// Size variance factor
	if pattern.FileSizes.StandardDev > float64(pattern.FileSizes.MeanSize) {
		complexityScore += 1
	}

	switch {
	case complexityScore >= 4:
		return "complex"
	case complexityScore >= 2:
		return "moderate"
	default:
		return "simple"
	}
}

func (be *BundlingEngine) addDomainSpecificRecommendations(recommendation *BundlingRecommendation, pattern *DataPattern) {
	recommendation.DomainHints = make(map[string]interface{})

	for _, domain := range pattern.DomainHints.DetectedDomains {
		switch domain {
		case "genomics":
			recommendation.DomainHints["genomics"] = map[string]interface{}{
				"group_by_type":        true,
				"preserve_fastq_pairs": true,
				"compress_text_files":  true,
				"bundle_size":          "500MB", // Larger bundles for genomics
			}
			recommendation.AlternativeOptions = append(recommendation.AlternativeOptions,
				"Consider bgzip for FASTQ files instead of bundling")

		case "climate":
			recommendation.DomainHints["climate"] = map[string]interface{}{
				"preserve_time_structure": true,
				"group_by_date":           true,
				"netcdf_chunking":         true,
			}
			recommendation.AlternativeOptions = append(recommendation.AlternativeOptions,
				"Consider NetCDF rechunking for better access patterns")

		case "machine_learning":
			recommendation.DomainHints["ml"] = map[string]interface{}{
				"preserve_data_splits": true,
				"separate_checkpoints": true,
				"compress_logs":        true,
			}
			recommendation.AlternativeOptions = append(recommendation.AlternativeOptions,
				"Keep training/validation splits in separate bundles")
		}
	}
}

// GetProgressChannel returns the progress channel for monitoring bundling operations
func (be *BundlingEngine) GetProgressChannel() <-chan *BundleProgress {
	return be.suitcase.GetProgress()
}

// ExtractBundle provides access to bundle extraction functionality
func (be *BundlingEngine) ExtractBundle(ctx context.Context, bundlePath, outputDir string) error {
	return be.suitcase.ExtractBundle(ctx, bundlePath, outputDir)
}

// ListBundleContents provides access to bundle content listing
func (be *BundlingEngine) ListBundleContents(ctx context.Context, bundlePath string) ([]string, error) {
	return be.suitcase.ListBundleContents(ctx, bundlePath)
}
