package data

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DataPattern represents the analyzed characteristics of a dataset
type DataPattern struct {
	// Basic statistics
	TotalFiles     int64  `json:"total_files"`
	TotalSize      int64  `json:"total_size"`
	TotalSizeHuman string `json:"total_size_human"`

	// File size analysis
	FileSizes FileSizeAnalysis `json:"file_sizes"`

	// File type analysis
	FileTypes map[string]FileTypeInfo `json:"file_types"`

	// Directory analysis
	DirectoryDepth DirectoryAnalysis `json:"directory_analysis"`

	// Access pattern hints
	AccessPatterns AccessPatternAnalysis `json:"access_patterns"`

	// Efficiency metrics
	Efficiency EfficiencyMetrics `json:"efficiency"`

	// Research domain hints
	DomainHints DomainAnalysis `json:"domain_hints"`

	// Analysis metadata
	AnalyzedPath string    `json:"analyzed_path"`
	AnalysisTime time.Time `json:"analysis_time"`
	SampleSize   int64     `json:"sample_size"` // If sampling was used
}

// FileSizeAnalysis provides detailed file size statistics
type FileSizeAnalysis struct {
	MinSize     int64   `json:"min_size"`
	MaxSize     int64   `json:"max_size"`
	MeanSize    int64   `json:"mean_size"`
	MedianSize  int64   `json:"median_size"`
	P95Size     int64   `json:"p95_size"`
	P99Size     int64   `json:"p99_size"`
	StandardDev float64 `json:"standard_deviation"`

	// Size distribution buckets
	SizeBuckets map[string]int64 `json:"size_buckets"`

	// Small file analysis (important for S3 efficiency)
	SmallFiles SmallFileAnalysis `json:"small_files"`
}

// SmallFileAnalysis focuses on small file patterns that affect S3 efficiency
type SmallFileAnalysis struct {
	CountUnder1KB   int64 `json:"count_under_1kb"`
	CountUnder10KB  int64 `json:"count_under_10kb"`
	CountUnder100KB int64 `json:"count_under_100kb"`
	CountUnder1MB   int64 `json:"count_under_1mb"`

	SizeUnder1KB   int64 `json:"size_under_1kb"`
	SizeUnder10KB  int64 `json:"size_under_10kb"`
	SizeUnder100KB int64 `json:"size_under_100kb"`
	SizeUnder1MB   int64 `json:"size_under_1mb"`

	PercentageSmall  float64 `json:"percentage_small_files"`
	PotentialSavings float64 `json:"potential_bundling_savings"` // Estimated cost savings from bundling
}

// FileTypeInfo represents information about a specific file type
type FileTypeInfo struct {
	Extension      string  `json:"extension"`
	Count          int64   `json:"count"`
	TotalSize      int64   `json:"total_size"`
	AverageSize    int64   `json:"average_size"`
	Percentage     float64 `json:"percentage_of_total"`
	Compressible   bool    `json:"compressible"`
	CompressionEst float64 `json:"estimated_compression_ratio"`
}

// DirectoryAnalysis analyzes directory structure patterns
type DirectoryAnalysis struct {
	MaxDepth       int     `json:"max_depth"`
	AverageDepth   float64 `json:"average_depth"`
	DirectoryCount int64   `json:"directory_count"`
	FilesPerDir    float64 `json:"average_files_per_directory"`

	// Patterns that might indicate structure
	HasDateDirs bool `json:"has_date_directories"`
	HasTypeDirs bool `json:"has_type_directories"`
	IsFlat      bool `json:"is_flat_structure"`
	IsDeep      bool `json:"is_deep_structure"`
}

// AccessPatternAnalysis analyzes file access patterns and provides hints
type AccessPatternAnalysis struct {
	// File age analysis
	NewestFile time.Time `json:"newest_file"`
	OldestFile time.Time `json:"oldest_file"`
	AverageAge float64   `json:"average_age_days"`

	// Modification patterns
	RecentlyModified int64 `json:"recently_modified_count"`
	StaleFiles       int64 `json:"stale_files_count"`

	// Hints about access patterns
	LikelyWriteOnce  bool `json:"likely_write_once"`
	LikelyFreqAccess bool `json:"likely_frequent_access"`
	LikelyArchival   bool `json:"likely_archival"`

	// Seasonal patterns (if detectable)
	HasSeasonality bool `json:"has_seasonal_pattern"`
}

// EfficiencyMetrics calculates various efficiency and cost-related metrics
type EfficiencyMetrics struct {
	// S3 efficiency metrics
	EstimatedPutRequests  int64   `json:"estimated_put_requests"`
	EstimatedGetRequests  int64   `json:"estimated_get_requests"`
	EstimatedRequestCosts float64 `json:"estimated_request_costs_monthly"`
	EstimatedStorageCosts float64 `json:"estimated_storage_costs_monthly"`

	// Transfer efficiency
	NetworkEfficiency float64 `json:"network_efficiency_score"`
	SmallFilePenalty  float64 `json:"small_file_penalty_score"`

	// Bundling potential
	BundlingRecommended bool    `json:"bundling_recommended"`
	EstimatedBundles    int64   `json:"estimated_bundles_after_bundling"`
	BundlingCostSavings float64 `json:"bundling_cost_savings_monthly"`

	// Storage class recommendations
	RecommendedStorageClass string  `json:"recommended_storage_class"`
	StorageClassSavings     float64 `json:"storage_class_savings_monthly"`
}

// DomainAnalysis provides hints about the research domain based on file patterns
type DomainAnalysis struct {
	DetectedDomains          []string               `json:"detected_domains"`
	Confidence               map[string]float64     `json:"domain_confidence"`
	DomainSpecificHints      map[string]interface{} `json:"domain_specific_hints"`
	RecommendedOptimizations []string               `json:"recommended_optimizations"`
}

// PatternAnalyzer analyzes data patterns and provides optimization recommendations
type PatternAnalyzer struct {
	sampleThreshold int64 // If more than this many files, use sampling
	maxSampleSize   int64 // Maximum number of files to sample
}

// NewPatternAnalyzer creates a new pattern analyzer
func NewPatternAnalyzer() *PatternAnalyzer {
	return &PatternAnalyzer{
		sampleThreshold: 10000, // Sample if more than 10k files
		maxSampleSize:   5000,  // Sample at most 5k files
	}
}

// AnalyzePattern analyzes a directory or file pattern and returns detailed analysis
func (pa *PatternAnalyzer) AnalyzePattern(ctx context.Context, path string) (*DataPattern, error) {
	startTime := time.Now()

	// Initialize the pattern analysis
	pattern := &DataPattern{
		AnalyzedPath: path,
		AnalysisTime: startTime,
		FileTypes:    make(map[string]FileTypeInfo),
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access path %s: %w", path, err)
	}

	// Collect file information
	var files []FileInfo
	if info.IsDir() {
		files, err = pa.scanDirectory(ctx, path)
	} else {
		files = []FileInfo{pa.getFileInfo(path, info)}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan files: %w", err)
	}

	// If too many files, sample them
	if int64(len(files)) > pa.sampleThreshold {
		files = pa.sampleFiles(files)
		pattern.SampleSize = int64(len(files))
	}

	// Perform various analyses
	pa.analyzeBasicStats(pattern, files)
	pa.analyzeFileSizes(pattern, files)
	pa.analyzeFileTypes(pattern, files)
	pa.analyzeDirectoryStructure(pattern, path)
	pa.analyzeAccessPatterns(pattern, files)
	pa.calculateEfficiencyMetrics(pattern, files)
	pa.analyzeDomainHints(pattern, files)

	return pattern, nil
}

// FileInfo holds information about a single file
type FileInfo struct {
	Path         string
	Size         int64
	ModTime      time.Time
	IsDir        bool
	Extension    string
	RelativePath string
	Depth        int
}

// scanDirectory recursively scans a directory and collects file information
func (pa *PatternAnalyzer) scanDirectory(ctx context.Context, rootPath string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return nil // Skip files that can't be accessed
		}

		info, err := d.Info()
		if err != nil {
			return nil // Skip files that can't be stat'd
		}

		relPath, _ := filepath.Rel(rootPath, path)
		depth := strings.Count(relPath, string(filepath.Separator))

		fileInfo := FileInfo{
			Path:         path,
			Size:         info.Size(),
			ModTime:      info.ModTime(),
			IsDir:        info.IsDir(),
			Extension:    strings.ToLower(filepath.Ext(path)),
			RelativePath: relPath,
			Depth:        depth,
		}

		files = append(files, fileInfo)

		// Stop early if we have too many files (we'll sample later)
		if int64(len(files)) > pa.sampleThreshold*2 {
			return filepath.SkipAll
		}

		return nil
	})

	return files, err
}

// getFileInfo creates FileInfo for a single file
func (pa *PatternAnalyzer) getFileInfo(path string, info os.FileInfo) FileInfo {
	return FileInfo{
		Path:         path,
		Size:         info.Size(),
		ModTime:      info.ModTime(),
		IsDir:        info.IsDir(),
		Extension:    strings.ToLower(filepath.Ext(path)),
		RelativePath: filepath.Base(path),
		Depth:        0,
	}
}

// sampleFiles selects a representative sample of files for analysis
func (pa *PatternAnalyzer) sampleFiles(files []FileInfo) []FileInfo {
	if int64(len(files)) <= pa.maxSampleSize {
		return files
	}

	// Sort files by size to ensure we get a good distribution
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size < files[j].Size
	})

	// Take every nth file to get a representative sample
	step := len(files) / int(pa.maxSampleSize)
	if step < 1 {
		step = 1
	}

	var sample []FileInfo
	for i := 0; i < len(files) && len(sample) < int(pa.maxSampleSize); i += step {
		sample = append(sample, files[i])
	}

	return sample
}

// analyzeBasicStats calculates basic statistics about the dataset
func (pa *PatternAnalyzer) analyzeBasicStats(pattern *DataPattern, files []FileInfo) {
	var totalSize int64
	var fileCount int64

	for _, file := range files {
		if !file.IsDir {
			totalSize += file.Size
			fileCount++
		}
	}

	pattern.TotalFiles = fileCount
	pattern.TotalSize = totalSize
	pattern.TotalSizeHuman = formatBytes(totalSize)
}

// analyzeFileSizes performs detailed file size analysis
func (pa *PatternAnalyzer) analyzeFileSizes(pattern *DataPattern, files []FileInfo) {
	var sizes []int64
	var totalSize int64

	// Collect all file sizes
	for _, file := range files {
		if !file.IsDir {
			sizes = append(sizes, file.Size)
			totalSize += file.Size
		}
	}

	if len(sizes) == 0 {
		return
	}

	// Sort sizes for percentile calculations
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] < sizes[j]
	})

	analysis := FileSizeAnalysis{
		MinSize:     sizes[0],
		MaxSize:     sizes[len(sizes)-1],
		MeanSize:    totalSize / int64(len(sizes)),
		MedianSize:  sizes[len(sizes)/2],
		P95Size:     sizes[int(float64(len(sizes))*0.95)],
		P99Size:     sizes[int(float64(len(sizes))*0.99)],
		SizeBuckets: make(map[string]int64),
	}

	// Calculate standard deviation
	var variance float64
	for _, size := range sizes {
		diff := float64(size - analysis.MeanSize)
		variance += diff * diff
	}
	analysis.StandardDev = math.Sqrt(variance / float64(len(sizes)))

	// Create size buckets
	buckets := map[string]struct{ min, max int64 }{
		"under_1KB":  {0, 1024},
		"1KB_10KB":   {1024, 10 * 1024},
		"10KB_100KB": {10 * 1024, 100 * 1024},
		"100KB_1MB":  {100 * 1024, 1024 * 1024},
		"1MB_10MB":   {1024 * 1024, 10 * 1024 * 1024},
		"10MB_100MB": {10 * 1024 * 1024, 100 * 1024 * 1024},
		"100MB_1GB":  {100 * 1024 * 1024, 1024 * 1024 * 1024},
		"over_1GB":   {1024 * 1024 * 1024, math.MaxInt64},
	}

	for _, size := range sizes {
		for bucketName, bucket := range buckets {
			if size >= bucket.min && size < bucket.max {
				analysis.SizeBuckets[bucketName]++
				break
			}
		}
	}

	// Analyze small files (important for S3 efficiency)
	pa.analyzeSmallFiles(&analysis, sizes, totalSize)

	pattern.FileSizes = analysis
}

// analyzeSmallFiles focuses on small file analysis for S3 optimization
func (pa *PatternAnalyzer) analyzeSmallFiles(analysis *FileSizeAnalysis, sizes []int64, totalSize int64) {
	smallFiles := SmallFileAnalysis{}

	for _, size := range sizes {
		if size < 1024 { // 1KB
			smallFiles.CountUnder1KB++
			smallFiles.SizeUnder1KB += size
		}
		if size < 10*1024 { // 10KB
			smallFiles.CountUnder10KB++
			smallFiles.SizeUnder10KB += size
		}
		if size < 100*1024 { // 100KB
			smallFiles.CountUnder100KB++
			smallFiles.SizeUnder100KB += size
		}
		if size < 1024*1024 { // 1MB
			smallFiles.CountUnder1MB++
			smallFiles.SizeUnder1MB += size
		}
	}

	// Calculate percentage of small files
	if len(sizes) > 0 {
		smallFiles.PercentageSmall = float64(smallFiles.CountUnder1MB) / float64(len(sizes)) * 100
	}

	// Estimate potential savings from bundling small files
	// S3 PUT requests cost $0.0005 per 1000 requests (us-east-1)
	putCostPer1000 := 0.0005
	currentPutCost := float64(smallFiles.CountUnder1MB) * putCostPer1000 / 1000

	// Assume bundling would reduce to 1/100th the number of files
	bundledPutCost := float64(smallFiles.CountUnder1MB/100) * putCostPer1000 / 1000
	smallFiles.PotentialSavings = currentPutCost - bundledPutCost

	analysis.SmallFiles = smallFiles
}

// analyzeFileTypes analyzes file type distribution and characteristics
func (pa *PatternAnalyzer) analyzeFileTypes(pattern *DataPattern, files []FileInfo) {
	typeStats := make(map[string]*FileTypeInfo)

	for _, file := range files {
		if file.IsDir {
			continue
		}

		ext := file.Extension
		if ext == "" {
			ext = "(no extension)"
		}

		if typeStats[ext] == nil {
			typeStats[ext] = &FileTypeInfo{
				Extension: ext,
			}
		}

		typeStats[ext].Count++
		typeStats[ext].TotalSize += file.Size
	}

	// Calculate derived statistics and add compression estimates
	for ext, info := range typeStats {
		if info.Count > 0 {
			info.AverageSize = info.TotalSize / info.Count
			info.Percentage = float64(info.TotalSize) / float64(pattern.TotalSize) * 100

			// Add compression and domain-specific information
			pa.addFileTypeHints(info, ext)
		}

		pattern.FileTypes[ext] = *info
	}
}

// addFileTypeHints adds compression estimates and domain hints for file types
func (pa *PatternAnalyzer) addFileTypeHints(info *FileTypeInfo, ext string) {
	// Define known file type characteristics
	fileTypeHints := map[string]struct {
		compressible bool
		compression  float64
		domain       string
	}{
		// Already compressed formats
		".gz":  {false, 1.0, ""},
		".zip": {false, 1.0, ""},
		".bz2": {false, 1.0, ""},
		".xz":  {false, 1.0, ""},
		".7z":  {false, 1.0, ""},

		// Images (already compressed)
		".jpg":  {false, 1.0, ""},
		".jpeg": {false, 1.0, ""},
		".png":  {false, 1.0, ""},
		".gif":  {false, 1.0, ""},

		// Genomics files
		".fastq": {true, 0.25, "genomics"},
		".fasta": {true, 0.3, "genomics"},
		".sam":   {true, 0.2, "genomics"},
		".bam":   {false, 1.0, "genomics"}, // Already compressed
		".vcf":   {true, 0.15, "genomics"},
		".bed":   {true, 0.3, "genomics"},

		// Climate/Scientific data
		".nc":     {true, 0.5, "climate"},
		".hdf5":   {false, 1.0, "climate"}, // Can be compressed internally
		".grib":   {false, 1.0, "climate"}, // Already compressed
		".netcdf": {true, 0.5, "climate"},

		// Text and data files
		".txt":  {true, 0.3, ""},
		".csv":  {true, 0.2, ""},
		".json": {true, 0.25, ""},
		".xml":  {true, 0.2, ""},
		".log":  {true, 0.15, ""},

		// Code and config
		".py":   {true, 0.4, ""},
		".js":   {true, 0.35, ""},
		".yaml": {true, 0.3, ""},
		".yml":  {true, 0.3, ""},

		// Documents
		".pdf":  {false, 1.0, ""},
		".doc":  {true, 0.6, ""},
		".docx": {true, 0.6, ""},
	}

	if hints, exists := fileTypeHints[ext]; exists {
		info.Compressible = hints.compressible
		info.CompressionEst = hints.compression
	} else {
		// Default: assume text-like files are compressible
		info.Compressible = true
		info.CompressionEst = 0.4
	}
}

// analyzeDirectoryStructure analyzes the directory structure patterns
func (pa *PatternAnalyzer) analyzeDirectoryStructure(pattern *DataPattern, rootPath string) {
	analysis := DirectoryAnalysis{}

	var totalDepth, dirCount, fileCount int
	var maxDepth int

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(rootPath, path)
		depth := strings.Count(relPath, string(filepath.Separator))

		if d.IsDir() {
			dirCount++
			totalDepth += depth
			if depth > maxDepth {
				maxDepth = depth
			}

			// Check for common directory patterns
			dirName := strings.ToLower(d.Name())
			if pa.isDateDirectory(dirName) {
				analysis.HasDateDirs = true
			}
			if pa.isTypeDirectory(dirName) {
				analysis.HasTypeDirs = true
			}
		} else {
			fileCount++
		}

		return nil
	})

	analysis.MaxDepth = maxDepth
	analysis.DirectoryCount = int64(dirCount)

	if dirCount > 0 {
		analysis.AverageDepth = float64(totalDepth) / float64(dirCount)
		analysis.FilesPerDir = float64(fileCount) / float64(dirCount)
	}

	// Determine structure characteristics
	analysis.IsFlat = maxDepth <= 2
	analysis.IsDeep = maxDepth > 5

	pattern.DirectoryDepth = analysis
}

// isDateDirectory checks if a directory name suggests date-based organization
func (pa *PatternAnalyzer) isDateDirectory(name string) bool {
	datePatterns := []string{
		"2023", "2024", "2025", // Years
		"jan", "feb", "mar", "apr", "may", "jun",
		"jul", "aug", "sep", "oct", "nov", "dec", // Months
		"01", "02", "03", "04", "05", "06",
		"07", "08", "09", "10", "11", "12", // Numeric months
	}

	for _, pattern := range datePatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}

	return false
}

// isTypeDirectory checks if a directory name suggests type-based organization
func (pa *PatternAnalyzer) isTypeDirectory(name string) bool {
	typePatterns := []string{
		"data", "raw", "processed", "results", "output",
		"images", "docs", "logs", "temp", "backup",
		"src", "source", "bin", "lib", "config",
	}

	for _, pattern := range typePatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}

	return false
}

// analyzeAccessPatterns analyzes file timestamps and access patterns
func (pa *PatternAnalyzer) analyzeAccessPatterns(pattern *DataPattern, files []FileInfo) {
	if len(files) == 0 {
		return
	}

	analysis := AccessPatternAnalysis{}

	var totalAge float64
	var fileCount int
	now := time.Now()

	analysis.NewestFile = files[0].ModTime
	analysis.OldestFile = files[0].ModTime

	recentThreshold := now.AddDate(0, 0, -30) // 30 days ago
	staleThreshold := now.AddDate(0, 0, -365) // 1 year ago

	for _, file := range files {
		if file.IsDir {
			continue
		}

		fileCount++
		age := now.Sub(file.ModTime).Hours() / 24 // Age in days
		totalAge += age

		if file.ModTime.After(analysis.NewestFile) {
			analysis.NewestFile = file.ModTime
		}
		if file.ModTime.Before(analysis.OldestFile) {
			analysis.OldestFile = file.ModTime
		}

		if file.ModTime.After(recentThreshold) {
			analysis.RecentlyModified++
		}
		if file.ModTime.Before(staleThreshold) {
			analysis.StaleFiles++
		}
	}

	if fileCount > 0 {
		analysis.AverageAge = totalAge / float64(fileCount)
	}

	// Determine access pattern hints
	staleFraction := float64(analysis.StaleFiles) / float64(fileCount)
	recentFraction := float64(analysis.RecentlyModified) / float64(fileCount)

	analysis.LikelyArchival = staleFraction > 0.8
	analysis.LikelyWriteOnce = recentFraction < 0.1 && staleFraction > 0.5
	analysis.LikelyFreqAccess = recentFraction > 0.3

	pattern.AccessPatterns = analysis
}

// calculateEfficiencyMetrics calculates S3 and transfer efficiency metrics
func (pa *PatternAnalyzer) calculateEfficiencyMetrics(pattern *DataPattern, files []FileInfo) {
	metrics := EfficiencyMetrics{}

	// S3 cost calculations (US East 1 pricing)
	putCostPer1000 := 0.0005  // $0.0005 per 1,000 PUT requests
	getCostPer1000 := 0.0004  // $0.0004 per 1,000 GET requests
	storageCostPerGB := 0.023 // $0.023 per GB per month (Standard)

	fileCount := pattern.TotalFiles
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)

	// Estimate requests and costs
	metrics.EstimatedPutRequests = fileCount
	metrics.EstimatedGetRequests = fileCount // Assume each file is accessed once per month

	metrics.EstimatedRequestCosts = (float64(fileCount) * putCostPer1000 / 1000) +
		(float64(fileCount) * getCostPer1000 / 1000)
	metrics.EstimatedStorageCosts = totalSizeGB * storageCostPerGB

	// Calculate efficiency scores
	smallFileRatio := float64(pattern.FileSizes.SmallFiles.CountUnder1MB) / float64(fileCount)
	metrics.SmallFilePenalty = smallFileRatio * 100 // Higher = worse

	// Network efficiency (larger files = more efficient)
	avgFileSizeMB := float64(pattern.FileSizes.MeanSize) / (1024 * 1024)
	if avgFileSizeMB > 100 {
		metrics.NetworkEfficiency = 100
	} else if avgFileSizeMB > 10 {
		metrics.NetworkEfficiency = 80
	} else if avgFileSizeMB > 1 {
		metrics.NetworkEfficiency = 60
	} else {
		metrics.NetworkEfficiency = 30
	}

	// Bundling recommendations
	if pattern.FileSizes.SmallFiles.CountUnder1MB > 100 && smallFileRatio > 0.5 {
		metrics.BundlingRecommended = true
		metrics.EstimatedBundles = pattern.FileSizes.SmallFiles.CountUnder1MB / 100 // Assume 100 files per bundle
		metrics.BundlingCostSavings = pattern.FileSizes.SmallFiles.PotentialSavings
	}

	// Storage class recommendations
	if pattern.AccessPatterns.LikelyArchival {
		metrics.RecommendedStorageClass = "GLACIER"
		glacierCostPerGB := 0.004 // $0.004 per GB per month
		metrics.StorageClassSavings = totalSizeGB * (storageCostPerGB - glacierCostPerGB)
	} else if pattern.AccessPatterns.LikelyWriteOnce {
		metrics.RecommendedStorageClass = "STANDARD_IA"
		iaCostPerGB := 0.0125 // $0.0125 per GB per month
		metrics.StorageClassSavings = totalSizeGB * (storageCostPerGB - iaCostPerGB)
	} else {
		metrics.RecommendedStorageClass = "STANDARD"
	}

	pattern.Efficiency = metrics
}

// analyzeDomainHints analyzes file patterns to detect research domains
func (pa *PatternAnalyzer) analyzeDomainHints(pattern *DataPattern, files []FileInfo) {
	analysis := DomainAnalysis{
		DetectedDomains:          make([]string, 0),
		Confidence:               make(map[string]float64),
		DomainSpecificHints:      make(map[string]interface{}),
		RecommendedOptimizations: make([]string, 0),
	}

	// Domain detection based on file extensions and patterns
	domainScores := make(map[string]float64)

	for ext, typeInfo := range pattern.FileTypes {
		switch ext {
		case ".fastq", ".fasta", ".sam", ".bam", ".vcf", ".bed", ".gff":
			domainScores["genomics"] += typeInfo.Percentage
		case ".nc", ".hdf5", ".grib", ".netcdf":
			domainScores["climate"] += typeInfo.Percentage
		case ".pkl", ".h5", ".model", ".ckpt":
			domainScores["machine_learning"] += typeInfo.Percentage
		case ".tif", ".tiff", ".geotiff":
			domainScores["geospatial"] += typeInfo.Percentage
		case ".fits":
			domainScores["astronomy"] += typeInfo.Percentage
		}
	}

	// Determine detected domains (threshold of 10% file content)
	for domain, score := range domainScores {
		if score > 10.0 {
			analysis.DetectedDomains = append(analysis.DetectedDomains, domain)
			analysis.Confidence[domain] = score / 100.0

			// Add domain-specific recommendations
			switch domain {
			case "genomics":
				analysis.RecommendedOptimizations = append(analysis.RecommendedOptimizations,
					"Consider compression for FASTQ files",
					"Use multipart upload for large BAM files")
				analysis.DomainSpecificHints["genomics"] = map[string]interface{}{
					"compression_recommended": true,
					"typical_access_pattern":  "write_once_read_many",
					"bundling_threshold":      1000,
				}
			case "climate":
				analysis.RecommendedOptimizations = append(analysis.RecommendedOptimizations,
					"NetCDF files often benefit from chunking",
					"Consider time-based directory organization")
			case "machine_learning":
				analysis.RecommendedOptimizations = append(analysis.RecommendedOptimizations,
					"Model checkpoints suitable for infrequent access storage",
					"Training data benefits from high-throughput access")
			}
		}
	}

	pattern.DomainHints = analysis
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GenerateAnalysisID creates a unique ID for an analysis based on path and timestamp
func (pa *PatternAnalyzer) GenerateAnalysisID(path string) string {
	hasher := sha256.New()
	hasher.Write([]byte(path + time.Now().Format(time.RFC3339)))
	return fmt.Sprintf("%x", hasher.Sum(nil))[:16]
}
