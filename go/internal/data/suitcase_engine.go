package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SuitcaseEngine implements file bundling using the Suitcase tool from Duke
type SuitcaseEngine struct {
	config       *SuitcaseConfig
	workingDir   string
	tempDir      string
	progressChan chan *BundleProgress
}

// SuitcaseConfig contains configuration for the Suitcase bundling engine
type SuitcaseConfig struct {
	// Bundling parameters
	TargetBundleSize  string `json:"target_bundle_size"`   // e.g., "100MB", "1GB"
	MaxFilesPerBundle int    `json:"max_files_per_bundle"` // Maximum files per bundle
	PreserveMetadata  bool   `json:"preserve_metadata"`    // Preserve file metadata
	CompressionLevel  int    `json:"compression_level"`    // 0-9, 0=no compression

	// File selection
	IncludePatterns []string `json:"include_patterns"` // File patterns to include
	ExcludePatterns []string `json:"exclude_patterns"` // File patterns to exclude
	SizeThreshold   string   `json:"size_threshold"`   // Bundle files smaller than this

	// Output options
	OutputFormat    string `json:"output_format"`    // "tar", "tar.gz", "zip"
	OutputDirectory string `json:"output_directory"` // Where to place bundles
	NamingTemplate  string `json:"naming_template"`  // Bundle naming pattern

	// Performance options
	Parallel    bool   `json:"parallel"`     // Enable parallel bundling
	WorkerCount int    `json:"worker_count"` // Number of parallel workers
	BufferSize  string `json:"buffer_size"`  // I/O buffer size

	// Research domain optimization
	DomainOptimization string            `json:"domain_optimization"` // "genomics", "climate", "ml", "general"
	CustomMetadata     map[string]string `json:"custom_metadata"`     // Custom metadata to embed
}

// BundleProgress represents progress information for bundling operations
type BundleProgress struct {
	BundleID         string    `json:"bundle_id"`
	BundleName       string    `json:"bundle_name"`
	FilesProcessed   int64     `json:"files_processed"`
	TotalFiles       int64     `json:"total_files"`
	BytesProcessed   int64     `json:"bytes_processed"`
	TotalBytes       int64     `json:"total_bytes"`
	CurrentFile      string    `json:"current_file"`
	Speed            float64   `json:"speed_mbps"`
	EstimatedTime    string    `json:"estimated_time_remaining"`
	Status           string    `json:"status"` // "bundling", "compressing", "complete", "error"
	ErrorMessage     string    `json:"error_message,omitempty"`
	StartTime        time.Time `json:"start_time"`
	CompletionTime   time.Time `json:"completion_time,omitempty"`
	CompressionRatio float64   `json:"compression_ratio"`
}

// BundleResult contains the results of a bundling operation
type BundleResult struct {
	BundleID          string                 `json:"bundle_id"`
	SourcePath        string                 `json:"source_path"`
	OutputPath        string                 `json:"output_path"`
	BundlePaths       []string               `json:"bundle_paths"`
	OriginalFileCount int64                  `json:"original_file_count"`
	OriginalSize      int64                  `json:"original_size_bytes"`
	BundledFileCount  int64                  `json:"bundled_file_count"`
	BundledSize       int64                  `json:"bundled_size_bytes"`
	CompressionRatio  float64                `json:"compression_ratio"`
	SpaceEfficiency   float64                `json:"space_efficiency_percent"`
	CostSavings       CostSavingsEstimate    `json:"cost_savings"`
	BundleManifest    []BundleManifestEntry  `json:"bundle_manifest"`
	ProcessingTime    time.Duration          `json:"processing_time"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// BundleManifestEntry describes a single bundle in the result
type BundleManifestEntry struct {
	BundleName       string    `json:"bundle_name"`
	BundlePath       string    `json:"bundle_path"`
	FileCount        int64     `json:"file_count"`
	Size             int64     `json:"size_bytes"`
	OriginalSize     int64     `json:"original_size_bytes"`
	CompressionRatio float64   `json:"compression_ratio"`
	Files            []string  `json:"files"`
	Checksum         string    `json:"checksum"`
	CreatedAt        time.Time `json:"created_at"`
}

// CostSavingsEstimate estimates the cost savings from bundling
type CostSavingsEstimate struct {
	RequestCostBefore float64 `json:"request_cost_before_monthly"`
	RequestCostAfter  float64 `json:"request_cost_after_monthly"`
	RequestSavings    float64 `json:"request_savings_monthly"`
	StorageCostBefore float64 `json:"storage_cost_before_monthly"`
	StorageCostAfter  float64 `json:"storage_cost_after_monthly"`
	StorageSavings    float64 `json:"storage_savings_monthly"`
	TotalSavings      float64 `json:"total_savings_monthly"`
	SavingsPercentage float64 `json:"savings_percentage"`
}

// NewSuitcaseEngine creates a new Suitcase bundling engine
func NewSuitcaseEngine(config *SuitcaseConfig) *SuitcaseEngine {
	if config == nil {
		config = &SuitcaseConfig{
			TargetBundleSize:   "100MB",
			MaxFilesPerBundle:  1000,
			PreserveMetadata:   true,
			CompressionLevel:   6,
			OutputFormat:       "tar.gz",
			Parallel:           true,
			WorkerCount:        4,
			BufferSize:         "64KB",
			DomainOptimization: "general",
		}
	}

	// Set default patterns if not provided
	if len(config.IncludePatterns) == 0 {
		config.IncludePatterns = []string{"*"}
	}
	if config.SizeThreshold == "" {
		config.SizeThreshold = "1MB"
	}
	if config.NamingTemplate == "" {
		config.NamingTemplate = "bundle_{index:04d}_{timestamp}"
	}

	return &SuitcaseEngine{
		config:       config,
		workingDir:   ".",
		tempDir:      os.TempDir(),
		progressChan: make(chan *BundleProgress, 100),
	}
}

// IsAvailable checks if Suitcase is installed and available
func (se *SuitcaseEngine) IsAvailable(ctx context.Context) error {
	// Check if Python is available
	if err := se.checkPython(ctx); err != nil {
		return fmt.Errorf("python not available: %w", err)
	}

	// Check if Suitcase is installed
	if err := se.checkSuitcase(ctx); err != nil {
		return fmt.Errorf("suitcase not available: %w", err)
	}

	return nil
}

// BundleFiles bundles small files according to the configuration
func (se *SuitcaseEngine) BundleFiles(ctx context.Context, sourcePath string) (*BundleResult, error) {
	startTime := time.Now()

	// Validate source path
	if _, err := os.Stat(sourcePath); err != nil {
		return nil, fmt.Errorf("source path not accessible: %w", err)
	}

	// Create output directory
	outputDir := se.config.OutputDirectory
	if outputDir == "" {
		outputDir = filepath.Join(filepath.Dir(sourcePath), "bundles")
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Analyze source files first
	fileAnalysis, err := se.analyzeSourceFiles(ctx, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze source files: %w", err)
	}

	// Create bundling strategy based on analysis
	strategy, err := se.createBundlingStrategy(fileAnalysis)
	if err != nil {
		return nil, fmt.Errorf("failed to create bundling strategy: %w", err)
	}

	// Execute bundling
	result, err := se.executeBundling(ctx, sourcePath, outputDir, strategy)
	if err != nil {
		return nil, fmt.Errorf("bundling failed: %w", err)
	}

	// Calculate final metrics
	result.ProcessingTime = time.Since(startTime)
	result.CostSavings = se.calculateCostSavings(fileAnalysis, result)

	return result, nil
}

// GetProgress returns the current progress channel
func (se *SuitcaseEngine) GetProgress() <-chan *BundleProgress {
	return se.progressChan
}

// checkPython verifies Python is available
func (se *SuitcaseEngine) checkPython(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "python3", "--version")
	if err := cmd.Run(); err != nil {
		// Try python instead of python3
		cmd = exec.CommandContext(ctx, "python", "--version")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("python not found in PATH")
		}
	}
	return nil
}

// checkSuitcase verifies Suitcase is installed
func (se *SuitcaseEngine) checkSuitcase(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "python3", "-c", "import suitcase; print(suitcase.__version__)")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try with python
		cmd = exec.CommandContext(ctx, "python", "-c", "import suitcase; print(suitcase.__version__)")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("suitcase not installed. Install with: pip install suitcase")
		}
	}

	version := strings.TrimSpace(string(output))
	if version == "" {
		return fmt.Errorf("could not determine suitcase version")
	}

	return nil
}

// FileAnalysis contains analysis of source files for bundling strategy
type FileAnalysis struct {
	TotalFiles         int64
	TotalSize          int64
	SmallFiles         []FileEntry
	LargeFiles         []FileEntry
	AverageSize        int64
	SizeDistribution   map[string]int64
	FileTypes          map[string]int64
	DirectoryStructure map[string]int64
}

// FileEntry represents a file to be processed
type FileEntry struct {
	Path         string
	Size         int64
	ModTime      time.Time
	Extension    string
	RelativePath string
}

// analyzeSourceFiles analyzes the source directory to create optimal bundling strategy
func (se *SuitcaseEngine) analyzeSourceFiles(ctx context.Context, sourcePath string) (*FileAnalysis, error) {
	analysis := &FileAnalysis{
		SmallFiles:         make([]FileEntry, 0),
		LargeFiles:         make([]FileEntry, 0),
		SizeDistribution:   make(map[string]int64),
		FileTypes:          make(map[string]int64),
		DirectoryStructure: make(map[string]int64),
	}

	sizeThreshold, err := se.parseSize(se.config.SizeThreshold)
	if err != nil {
		sizeThreshold = 1024 * 1024 // Default 1MB
	}

	err = filepath.WalkDir(sourcePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}

		if d.IsDir() {
			relPath, _ := filepath.Rel(sourcePath, path)
			analysis.DirectoryStructure[relPath]++
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil // Skip files that can't be stat'd
		}

		// Check if file matches include/exclude patterns
		if !se.matchesPatterns(path, se.config.IncludePatterns, se.config.ExcludePatterns) {
			return nil
		}

		relPath, _ := filepath.Rel(sourcePath, path)
		ext := strings.ToLower(filepath.Ext(path))

		entry := FileEntry{
			Path:         path,
			Size:         info.Size(),
			ModTime:      info.ModTime(),
			Extension:    ext,
			RelativePath: relPath,
		}

		analysis.TotalFiles++
		analysis.TotalSize += info.Size()
		analysis.FileTypes[ext]++

		// Categorize by size
		if info.Size() <= sizeThreshold {
			analysis.SmallFiles = append(analysis.SmallFiles, entry)
		} else {
			analysis.LargeFiles = append(analysis.LargeFiles, entry)
		}

		// Size distribution buckets
		sizeCategory := se.getSizeCategory(info.Size())
		analysis.SizeDistribution[sizeCategory]++

		return nil
	})

	if analysis.TotalFiles > 0 {
		analysis.AverageSize = analysis.TotalSize / analysis.TotalFiles
	}

	return analysis, err
}

// BundlingStrategy defines how files should be bundled
type BundlingStrategy struct {
	BundleGroups    []BundleGroup          `json:"bundle_groups"`
	CompressionMode string                 `json:"compression_mode"`
	Parallelism     int                    `json:"parallelism"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// BundleGroup represents a group of files to be bundled together
type BundleGroup struct {
	Name           string      `json:"name"`
	Files          []FileEntry `json:"files"`
	TargetSize     int64       `json:"target_size"`
	ExpectedSize   int64       `json:"expected_size"`
	CompressionEst float64     `json:"compression_estimate"`
	Priority       int         `json:"priority"`
}

// createBundlingStrategy creates an optimal bundling strategy based on file analysis
func (se *SuitcaseEngine) createBundlingStrategy(analysis *FileAnalysis) (*BundlingStrategy, error) {
	strategy := &BundlingStrategy{
		BundleGroups: make([]BundleGroup, 0),
		Parallelism:  se.config.WorkerCount,
		Metadata:     make(map[string]interface{}),
	}

	targetBundleSize, err := se.parseSize(se.config.TargetBundleSize)
	if err != nil {
		targetBundleSize = 100 * 1024 * 1024 // Default 100MB
	}

	// Apply domain-specific optimizations
	se.applyDomainOptimizations(strategy, analysis)

	// Group small files by type and directory for optimal bundling
	bundleGroups := se.groupFilesForBundling(analysis.SmallFiles, targetBundleSize)
	strategy.BundleGroups = bundleGroups

	// Set compression mode based on file types
	strategy.CompressionMode = se.determineCompressionMode(analysis)

	// Add metadata
	strategy.Metadata["total_files"] = analysis.TotalFiles
	strategy.Metadata["total_size"] = analysis.TotalSize
	strategy.Metadata["small_files"] = len(analysis.SmallFiles)
	strategy.Metadata["bundle_count"] = len(bundleGroups)
	strategy.Metadata["domain_optimization"] = se.config.DomainOptimization

	return strategy, nil
}

// executeBundling executes the bundling strategy
func (se *SuitcaseEngine) executeBundling(ctx context.Context, sourcePath, outputDir string, strategy *BundlingStrategy) (*BundleResult, error) {
	bundleID := fmt.Sprintf("bundle_%d", time.Now().Unix())

	result := &BundleResult{
		BundleID:       bundleID,
		SourcePath:     sourcePath,
		OutputPath:     outputDir,
		BundlePaths:    make([]string, 0),
		BundleManifest: make([]BundleManifestEntry, 0),
		Metadata:       make(map[string]interface{}),
	}

	// Process each bundle group
	for i, group := range strategy.BundleGroups {
		progress := &BundleProgress{
			BundleID:   bundleID,
			BundleName: group.Name,
			TotalFiles: int64(len(group.Files)),
			TotalBytes: group.ExpectedSize,
			Status:     "bundling",
			StartTime:  time.Now(),
		}

		se.progressChan <- progress

		// Create bundle file path
		bundleName := se.generateBundleName(i, group.Name)
		bundlePath := filepath.Join(outputDir, bundleName)

		// Execute bundling command
		manifestEntry, err := se.createBundle(ctx, group, bundlePath, progress)
		if err != nil {
			progress.Status = "error"
			progress.ErrorMessage = err.Error()
			se.progressChan <- progress
			return nil, fmt.Errorf("failed to create bundle %s: %w", bundleName, err)
		}

		result.BundlePaths = append(result.BundlePaths, bundlePath)
		result.BundleManifest = append(result.BundleManifest, *manifestEntry)
		result.BundledFileCount += manifestEntry.FileCount
		result.BundledSize += manifestEntry.Size
		result.OriginalSize += manifestEntry.OriginalSize

		progress.Status = "complete"
		progress.CompletionTime = time.Now()
		se.progressChan <- progress
	}

	// Calculate final statistics
	if result.OriginalSize > 0 {
		result.CompressionRatio = float64(result.BundledSize) / float64(result.OriginalSize)
		result.SpaceEfficiency = (1.0 - result.CompressionRatio) * 100
	}

	// Add metadata from strategy
	for k, v := range strategy.Metadata {
		result.Metadata[k] = v
	}

	return result, nil
}

// applyDomainOptimizations applies research domain-specific optimizations
func (se *SuitcaseEngine) applyDomainOptimizations(strategy *BundlingStrategy, analysis *FileAnalysis) {
	switch se.config.DomainOptimization {
	case "genomics":
		// Genomics files often have similar access patterns
		// Group by file type (FASTQ, FASTA, etc.) for better compression
		strategy.Metadata["optimization"] = "genomics_aware"
		strategy.Metadata["group_by_type"] = true

	case "climate":
		// Climate data often has temporal patterns
		// Group by time periods when possible
		strategy.Metadata["optimization"] = "temporal_grouping"
		strategy.Metadata["preserve_time_structure"] = true

	case "machine_learning":
		// ML datasets often have train/test/validation splits
		// Preserve directory structure that indicates data splits
		strategy.Metadata["optimization"] = "preserve_ml_structure"
		strategy.Metadata["respect_data_splits"] = true

	default:
		strategy.Metadata["optimization"] = "general_purpose"
	}
}

// groupFilesForBundling groups files optimally for bundling
func (se *SuitcaseEngine) groupFilesForBundling(files []FileEntry, targetSize int64) []BundleGroup {
	var groups []BundleGroup

	// Sort files by directory and type for better grouping
	sortedFiles := make([]FileEntry, len(files))
	copy(sortedFiles, files)

	// Group by directory and file type
	dirGroups := make(map[string][]FileEntry)
	for _, file := range sortedFiles {
		dir := filepath.Dir(file.RelativePath)
		ext := file.Extension
		key := fmt.Sprintf("%s|%s", dir, ext)
		dirGroups[key] = append(dirGroups[key], file)
	}

	// Create bundles from groups
	bundleIndex := 0
	for key, groupFiles := range dirGroups {
		parts := strings.Split(key, "|")
		dir := parts[0]
		ext := parts[1]

		// Split large groups into multiple bundles
		currentBundle := BundleGroup{
			Files:      make([]FileEntry, 0),
			TargetSize: targetSize,
		}
		currentSize := int64(0)

		for _, file := range groupFiles {
			if currentSize+file.Size > targetSize && len(currentBundle.Files) > 0 {
				// Complete current bundle
				currentBundle.Name = fmt.Sprintf("bundle_%04d_%s_%s", bundleIndex, dir, ext)
				currentBundle.ExpectedSize = currentSize
				currentBundle.CompressionEst = se.estimateCompressionRatio(currentBundle.Files)
				groups = append(groups, currentBundle)

				// Start new bundle
				bundleIndex++
				currentBundle = BundleGroup{
					Files:      make([]FileEntry, 0),
					TargetSize: targetSize,
				}
				currentSize = 0
			}

			currentBundle.Files = append(currentBundle.Files, file)
			currentSize += file.Size
		}

		// Add final bundle if it has files
		if len(currentBundle.Files) > 0 {
			currentBundle.Name = fmt.Sprintf("bundle_%04d_%s_%s", bundleIndex, dir, ext)
			currentBundle.ExpectedSize = currentSize
			currentBundle.CompressionEst = se.estimateCompressionRatio(currentBundle.Files)
			groups = append(groups, currentBundle)
			bundleIndex++
		}
	}

	return groups
}

// createBundle creates a single bundle file
func (se *SuitcaseEngine) createBundle(ctx context.Context, group BundleGroup, outputPath string, progress *BundleProgress) (*BundleManifestEntry, error) {
	startTime := time.Now()

	// Create temporary file list
	fileListPath := filepath.Join(se.tempDir, fmt.Sprintf("filelist_%d.txt", time.Now().UnixNano()))
	defer os.Remove(fileListPath)

	fileList, err := os.Create(fileListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file list: %w", err)
	}
	defer fileList.Close()

	filePaths := make([]string, 0, len(group.Files))
	for _, file := range group.Files {
		fmt.Fprintln(fileList, file.Path)
		filePaths = append(filePaths, file.RelativePath)
	}
	fileList.Close()

	// Build suitcase command
	cmd := se.buildSuitcaseCommand(ctx, fileListPath, outputPath)

	// Execute command with progress monitoring
	if err := se.runCommandWithProgress(cmd, progress); err != nil {
		return nil, fmt.Errorf("suitcase command failed: %w", err)
	}

	// Get final bundle info
	bundleInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat bundle file: %w", err)
	}

	// Calculate checksum
	checksum, err := se.calculateChecksum(outputPath)
	if err != nil {
		checksum = "unknown"
	}

	manifest := &BundleManifestEntry{
		BundleName:       filepath.Base(outputPath),
		BundlePath:       outputPath,
		FileCount:        int64(len(group.Files)),
		Size:             bundleInfo.Size(),
		OriginalSize:     group.ExpectedSize,
		CompressionRatio: float64(bundleInfo.Size()) / float64(group.ExpectedSize),
		Files:            filePaths,
		Checksum:         checksum,
		CreatedAt:        startTime,
	}

	return manifest, nil
}

// Helper functions

func (se *SuitcaseEngine) parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	multipliers := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(sizeStr, suffix) {
			numStr := strings.TrimSuffix(sizeStr, suffix)
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, err
			}
			return int64(num * float64(multiplier)), nil
		}
	}

	// Try parsing as plain number (assume bytes)
	num, err := strconv.ParseInt(sizeStr, 10, 64)
	return num, err
}

func (se *SuitcaseEngine) matchesPatterns(path string, include, exclude []string) bool {
	// Check exclude patterns first
	for _, pattern := range exclude {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return false
		}
	}

	// Check include patterns
	for _, pattern := range include {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}

	return len(include) == 0 // Include all if no patterns specified
}

func (se *SuitcaseEngine) getSizeCategory(size int64) string {
	if size < 1024 {
		return "under_1KB"
	} else if size < 10*1024 {
		return "1KB_10KB"
	} else if size < 100*1024 {
		return "10KB_100KB"
	} else if size < 1024*1024 {
		return "100KB_1MB"
	} else if size < 10*1024*1024 {
		return "1MB_10MB"
	} else {
		return "over_10MB"
	}
}

func (se *SuitcaseEngine) determineCompressionMode(analysis *FileAnalysis) string {
	// Check if files are already compressed
	compressedExts := map[string]bool{
		".gz": true, ".zip": true, ".bz2": true, ".xz": true,
		".7z": true, ".rar": true, ".jpg": true, ".png": true,
		".mp4": true, ".mp3": true, ".bam": true,
	}

	compressedCount := int64(0)
	for ext, count := range analysis.FileTypes {
		if compressedExts[ext] {
			compressedCount += count
		}
	}

	compressedRatio := float64(compressedCount) / float64(analysis.TotalFiles)

	if compressedRatio > 0.7 {
		return "store" // Don't compress already compressed files
	} else if compressedRatio > 0.3 {
		return "fast" // Light compression
	} else {
		return "best" // Maximum compression
	}
}

func (se *SuitcaseEngine) estimateCompressionRatio(files []FileEntry) float64 {
	// Estimate compression ratio based on file types
	totalSize := int64(0)
	compressibleSize := int64(0)

	compressibleExts := map[string]float64{
		".txt": 0.3, ".csv": 0.2, ".json": 0.25, ".xml": 0.2,
		".fastq": 0.25, ".fasta": 0.3, ".sam": 0.2, ".vcf": 0.15,
		".log": 0.15, ".py": 0.4, ".js": 0.35,
	}

	for _, file := range files {
		totalSize += file.Size
		if ratio, exists := compressibleExts[file.Extension]; exists {
			compressibleSize += int64(float64(file.Size) * ratio)
		} else {
			compressibleSize += file.Size // Assume no compression for unknown types
		}
	}

	if totalSize > 0 {
		return float64(compressibleSize) / float64(totalSize)
	}

	return 0.7 // Default compression ratio
}

func (se *SuitcaseEngine) generateBundleName(index int, groupName string) string {
	template := se.config.NamingTemplate
	timestamp := time.Now().Format("20060102_150405")

	bundleName := strings.ReplaceAll(template, "{index:04d}", fmt.Sprintf("%04d", index))
	bundleName = strings.ReplaceAll(bundleName, "{timestamp}", timestamp)
	bundleName = strings.ReplaceAll(bundleName, "{group}", groupName)

	// Add appropriate extension
	switch se.config.OutputFormat {
	case "tar":
		bundleName += ".tar"
	case "tar.gz":
		bundleName += ".tar.gz"
	case "zip":
		bundleName += ".zip"
	default:
		bundleName += ".tar.gz"
	}

	return bundleName
}

func (se *SuitcaseEngine) buildSuitcaseCommand(ctx context.Context, fileListPath, outputPath string) *exec.Cmd {
	args := []string{
		"-m", "suitcase",
		"pack",
		"--file-list", fileListPath,
		"--output", outputPath,
	}

	// Add compression settings
	if se.config.CompressionLevel > 0 {
		args = append(args, "--compression-level", strconv.Itoa(se.config.CompressionLevel))
	}

	// Add format specification
	args = append(args, "--format", se.config.OutputFormat)

	// Add metadata preservation
	if se.config.PreserveMetadata {
		args = append(args, "--preserve-metadata")
	}

	// Add buffer size
	if se.config.BufferSize != "" {
		args = append(args, "--buffer-size", se.config.BufferSize)
	}

	// Add custom metadata
	for key, value := range se.config.CustomMetadata {
		args = append(args, "--metadata", fmt.Sprintf("%s=%s", key, value))
	}

	return exec.CommandContext(ctx, "python3", args...)
}

func (se *SuitcaseEngine) runCommandWithProgress(cmd *exec.Cmd, progress *BundleProgress) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		return err
	}

	// Monitor progress from stdout
	go se.monitorProgress(stdout, progress)

	return cmd.Wait()
}

func (se *SuitcaseEngine) monitorProgress(reader io.Reader, progress *BundleProgress) {
	// Monitor suitcase output for progress information
	// This would parse the actual suitcase output format
	// For now, we'll simulate progress updates

	// Pattern to match: "Processing file X of Y: filename"
	_ = regexp.MustCompile(`Processing file (\d+) of (\d+): (.+)`)

	// Pattern to match: "Compressed X bytes to Y bytes"
	_ = regexp.MustCompile(`Compressed (\d+) bytes to (\d+) bytes`)

	// Read output line by line
	// This is a simplified implementation - actual parsing would depend on suitcase output format
	// For demonstration, we'll update progress periodically
	startTime := time.Now()

	for {
		time.Sleep(100 * time.Millisecond)

		elapsed := time.Since(startTime)
		if elapsed > 30*time.Second { // Timeout after 30 seconds for demo
			break
		}

		// Simulate progress
		progress.FilesProcessed = int64(float64(progress.TotalFiles) * elapsed.Seconds() / 30.0)
		if progress.FilesProcessed > progress.TotalFiles {
			progress.FilesProcessed = progress.TotalFiles
		}

		progress.BytesProcessed = int64(float64(progress.TotalBytes) * elapsed.Seconds() / 30.0)
		if progress.BytesProcessed > progress.TotalBytes {
			progress.BytesProcessed = progress.TotalBytes
		}

		// Calculate speed
		if elapsed.Seconds() > 0 {
			mbProcessed := float64(progress.BytesProcessed) / (1024 * 1024)
			progress.Speed = mbProcessed / elapsed.Seconds()
		}

		// Estimate remaining time
		if progress.FilesProcessed > 0 {
			remainingFiles := progress.TotalFiles - progress.FilesProcessed
			filesPerSecond := float64(progress.FilesProcessed) / elapsed.Seconds()
			if filesPerSecond > 0 {
				remainingSeconds := float64(remainingFiles) / filesPerSecond
				progress.EstimatedTime = fmt.Sprintf("%.0fs", remainingSeconds)
			}
		}

		// Send progress update
		select {
		case se.progressChan <- progress:
		default:
			// Channel full, skip this update
		}

		// Break if complete
		if progress.FilesProcessed >= progress.TotalFiles {
			break
		}
	}
}

func (se *SuitcaseEngine) calculateChecksum(filePath string) (string, error) {
	// Calculate MD5 checksum of the bundle file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	cmd := exec.Command("md5sum", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Extract just the hash part
	parts := strings.Fields(string(output))
	if len(parts) > 0 {
		return parts[0], nil
	}

	return "", fmt.Errorf("could not parse checksum output")
}

func (se *SuitcaseEngine) calculateCostSavings(analysis *FileAnalysis, result *BundleResult) CostSavingsEstimate {
	// S3 pricing (US East 1)
	putCostPer1000 := 0.0005  // $0.0005 per 1,000 PUT requests
	storageCostPerGB := 0.023 // $0.023 per GB per month

	// Before bundling costs
	originalFiles := float64(analysis.TotalFiles)
	originalSizeGB := float64(analysis.TotalSize) / (1024 * 1024 * 1024)

	requestCostBefore := originalFiles * putCostPer1000 / 1000
	storageCostBefore := originalSizeGB * storageCostPerGB

	// After bundling costs
	bundledFiles := float64(result.BundledFileCount)
	bundledSizeGB := float64(result.BundledSize) / (1024 * 1024 * 1024)

	requestCostAfter := bundledFiles * putCostPer1000 / 1000
	storageCostAfter := bundledSizeGB * storageCostPerGB

	// Calculate savings
	requestSavings := requestCostBefore - requestCostAfter
	storageSavings := storageCostBefore - storageCostAfter
	totalSavings := requestSavings + storageSavings

	totalCostBefore := requestCostBefore + storageCostBefore
	savingsPercentage := 0.0
	if totalCostBefore > 0 {
		savingsPercentage = (totalSavings / totalCostBefore) * 100
	}

	return CostSavingsEstimate{
		RequestCostBefore: requestCostBefore,
		RequestCostAfter:  requestCostAfter,
		RequestSavings:    requestSavings,
		StorageCostBefore: storageCostBefore,
		StorageCostAfter:  storageCostAfter,
		StorageSavings:    storageSavings,
		TotalSavings:      totalSavings,
		SavingsPercentage: savingsPercentage,
	}
}

// ExtractBundle extracts files from a bundle created by Suitcase
func (se *SuitcaseEngine) ExtractBundle(ctx context.Context, bundlePath, outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build extraction command
	args := []string{
		"-m", "suitcase",
		"unpack",
		"--bundle", bundlePath,
		"--output", outputDir,
	}

	if se.config.PreserveMetadata {
		args = append(args, "--preserve-metadata")
	}

	cmd := exec.CommandContext(ctx, "python3", args...)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("extraction failed: %w, output: %s", err, string(output))
	}

	return nil
}

// ListBundleContents lists the contents of a bundle without extracting
func (se *SuitcaseEngine) ListBundleContents(ctx context.Context, bundlePath string) ([]string, error) {
	args := []string{
		"-m", "suitcase",
		"list",
		"--bundle", bundlePath,
		"--json", // Get machine-readable output
	}

	cmd := exec.CommandContext(ctx, "python3", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list bundle contents: %w", err)
	}

	// Parse JSON output
	var contents []string
	if err := json.Unmarshal(output, &contents); err != nil {
		// If JSON parsing fails, try to parse as plain text
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				contents = append(contents, line)
			}
		}
	}

	return contents, nil
}
