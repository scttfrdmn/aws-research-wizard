package data

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RcloneEngine implements TransferEngine for rclone
type RcloneEngine struct {
	*BaseTransferEngine
	executablePath  string
	configPath      string
	activeTransfers map[string]*rcloneTransfer
	mu              sync.RWMutex
}

// rcloneTransfer tracks an active rclone transfer
type rcloneTransfer struct {
	id       string
	cmd      *exec.Cmd
	progress *TransferProgress
	callback ProgressCallback
	stats    *RcloneStats
	done     chan bool
	mu       sync.RWMutex
}

// RcloneStats represents rclone's JSON statistics output
type RcloneStats struct {
	Bytes        int64                    `json:"bytes"`
	Checks       int64                    `json:"checks"`
	Deletes      int64                    `json:"deletes"`
	Elapsed      float64                  `json:"elapsedTime"`
	Errors       int64                    `json:"errors"`
	FatalError   bool                     `json:"fatalError"`
	Renames      int64                    `json:"renames"`
	RetryError   bool                     `json:"retryError"`
	Speed        float64                  `json:"speed"`
	TotalBytes   int64                    `json:"totalBytes"`
	TotalChecks  int64                    `json:"totalChecks"`
	Transfers    int64                    `json:"transfers"`
	TransferTime float64                  `json:"transferTime"`
	Transferring []RcloneTransferringFile `json:"transferring"`
}

// RcloneTransferringFile represents a file being transferred
type RcloneTransferringFile struct {
	Bytes      int64   `json:"bytes"`
	Name       string  `json:"name"`
	Percentage int     `json:"percentage"`
	Size       int64   `json:"size"`
	Speed      float64 `json:"speed"`
	SpeedAvg   float64 `json:"speedAvg"`
}

// NewRcloneEngine creates a new rclone transfer engine
func NewRcloneEngine(executablePath, configPath string) *RcloneEngine {
	capabilities := EngineCapabilities{
		Protocols: []string{
			"s3", "gcs", "azure", "dropbox", "gdrive", "onedrive",
			"ftp", "sftp", "webdav", "http", "local",
		},
		SupportsResume:         true,
		SupportsProgress:       true,
		SupportsParallel:       true,
		SupportsCompression:    true,
		SupportsEncryption:     true,
		SupportsValidation:     true,
		SupportsBandwidthLimit: true,
		SupportsRetry:          true,
		OptimalFileSizeMin:     1024,                      // 1KB
		OptimalFileSizeMax:     1024 * 1024 * 1024 * 1024, // 1TB
		MaxConcurrency:         32,
		CloudOptimized:         []string{"aws", "gcp", "azure", "multi-cloud"},
	}

	base := NewBaseTransferEngine("rclone", "multi-cloud", capabilities)

	engine := &RcloneEngine{
		BaseTransferEngine: base,
		executablePath:     executablePath,
		configPath:         configPath,
		activeTransfers:    make(map[string]*rcloneTransfer),
	}

	return engine
}

// IsAvailable checks if rclone is available and working
func (e *RcloneEngine) IsAvailable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, e.executablePath, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("rclone not available: %w", err)
	}

	// Check if output contains version information
	if !strings.Contains(string(output), "rclone") {
		return fmt.Errorf("unexpected rclone version output: %s", string(output))
	}

	return nil
}

// Upload transfers data from local to remote
func (e *RcloneEngine) Upload(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// For upload, source should be local and destination should be remote
	if strings.Contains(req.Source, ":") && !strings.Contains(req.Destination, ":") {
		return nil, fmt.Errorf("rclone upload requires local source and remote destination")
	}

	// Build rclone copy command
	args := e.buildCopyArgs(req, "copy")

	// Execute transfer
	return e.executeTransfer(ctx, req, args)
}

// Download transfers data from remote to local
func (e *RcloneEngine) Download(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// For download, source should be remote and destination should be local
	if !strings.Contains(req.Source, ":") && strings.Contains(req.Destination, ":") {
		return nil, fmt.Errorf("rclone download requires remote source and local destination")
	}

	// Build rclone copy command
	args := e.buildCopyArgs(req, "copy")

	// Execute transfer
	return e.executeTransfer(ctx, req, args)
}

// Sync synchronizes data between two locations
func (e *RcloneEngine) Sync(ctx context.Context, req *SyncRequest) (*TransferResult, error) {
	// Convert SyncRequest to TransferRequest for consistency
	transferReq := &TransferRequest{
		ID:               req.ID,
		Source:           req.Source,
		Destination:      req.Destination,
		ProgressCallback: req.ProgressCallback,
		Context:          req.Context,
		Options: TransferOptions{
			Overwrite:      true, // sync typically overwrites
			Concurrency:    req.Options.Concurrency,
			BandwidthLimit: req.Options.BandwidthLimit,
			ToolSpecific: map[string]interface{}{
				"sync_options": req.Options,
			},
		},
	}

	// Build rclone sync command
	args := e.buildSyncArgs(transferReq, &req.Options)

	// Execute transfer
	return e.executeTransfer(ctx, transferReq, args)
}

// buildCopyArgs builds command arguments for copy operations
func (e *RcloneEngine) buildCopyArgs(req *TransferRequest, operation string) []string {
	args := []string{operation}

	// Add configuration file if specified
	if e.configPath != "" {
		args = append(args, "--config", e.configPath)
	}

	// Add progress reporting
	args = append(args, "--stats", "1s", "--stats-one-line")

	// Add performance options
	if req.Options.Concurrency > 0 {
		args = append(args, "--transfers", strconv.Itoa(req.Options.Concurrency))
	}

	// Add bandwidth limit
	if req.Options.BandwidthLimit > 0 {
		bandwidthMBps := req.Options.BandwidthLimit / (1024 * 1024)
		args = append(args, "--bwlimit", fmt.Sprintf("%dM", bandwidthMBps))
	}

	// Add verification
	if req.Options.Verify {
		args = append(args, "--checksum")
	}

	// Add compression if supported and requested
	if req.Options.Compress {
		// Note: rclone doesn't have a universal compression flag,
		// but some backends support it via specific flags
	}

	// Add overwrite behavior
	if req.Options.Overwrite {
		// Default behavior for copy is to overwrite
	} else {
		args = append(args, "--ignore-existing")
	}

	// Add tool-specific options
	if toolOpts, ok := req.Options.ToolSpecific["rclone"]; ok {
		if opts, ok := toolOpts.(map[string]interface{}); ok {
			for key, value := range opts {
				args = append(args, fmt.Sprintf("--%s", key))
				if value != nil && value != true {
					args = append(args, fmt.Sprintf("%v", value))
				}
			}
		}
	}

	// Add source and destination
	args = append(args, req.Source, req.Destination)

	return args
}

// buildSyncArgs builds command arguments for sync operations
func (e *RcloneEngine) buildSyncArgs(req *TransferRequest, syncOpts *SyncOptions) []string {
	args := []string{"sync"}

	// Add configuration file if specified
	if e.configPath != "" {
		args = append(args, "--config", e.configPath)
	}

	// Add progress reporting
	args = append(args, "--stats", "1s", "--stats-one-line")

	// Add performance options
	if syncOpts.Concurrency > 0 {
		args = append(args, "--transfers", strconv.Itoa(syncOpts.Concurrency))
	}

	// Add bandwidth limit
	if syncOpts.BandwidthLimit > 0 {
		bandwidthMBps := syncOpts.BandwidthLimit / (1024 * 1024)
		args = append(args, "--bwlimit", fmt.Sprintf("%dM", bandwidthMBps))
	}

	// Add delete behavior
	if !syncOpts.DeleteExtraneous {
		// Default sync behavior deletes, so we don't need to add anything
	}

	// Add dry run option
	if syncOpts.DryRun {
		args = append(args, "--dry-run")
	}

	// Add skip newer files
	if syncOpts.SkipNewer {
		args = append(args, "--update")
	}

	// Add exclude patterns
	for _, pattern := range syncOpts.Exclude {
		args = append(args, "--exclude", pattern)
	}

	// Add include patterns
	for _, pattern := range syncOpts.Include {
		args = append(args, "--include", pattern)
	}

	// Add tool-specific options
	if toolOpts, ok := syncOpts.ToolSpecific["rclone"]; ok {
		if opts, ok := toolOpts.(map[string]interface{}); ok {
			for key, value := range opts {
				args = append(args, fmt.Sprintf("--%s", key))
				if value != nil && value != true {
					args = append(args, fmt.Sprintf("%v", value))
				}
			}
		}
	}

	// Add source and destination
	args = append(args, req.Source, req.Destination)

	return args
}

// executeTransfer executes the rclone command and monitors progress
func (e *RcloneEngine) executeTransfer(ctx context.Context, req *TransferRequest, args []string) (*TransferResult, error) {
	startTime := time.Now()

	// Create command with JSON stats output for better progress monitoring
	statsArgs := append([]string{}, args...)
	statsArgs = append(statsArgs, "--stats-file-name-length", "0", "--use-json-log")

	cmd := exec.CommandContext(ctx, e.executablePath, statsArgs...)

	// Set up pipes for output monitoring
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Create transfer tracking
	transfer := &rcloneTransfer{
		id:  req.ID,
		cmd: cmd,
		progress: &TransferProgress{
			StartTime: startTime,
		},
		callback: req.ProgressCallback,
		stats:    &RcloneStats{},
		done:     make(chan bool, 1),
	}

	// Register active transfer
	e.mu.Lock()
	e.activeTransfers[req.ID] = transfer
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		delete(e.activeTransfers, req.ID)
		e.mu.Unlock()
	}()

	// Start command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start rclone: %w", err)
	}

	// Start output monitoring goroutines
	var wg sync.WaitGroup
	var transferError error
	var output strings.Builder
	var errorOutput strings.Builder

	// Monitor stdout for progress
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.monitorProgress(stdout, transfer)
	}()

	// Monitor stderr for errors and additional output
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			errorOutput.WriteString(line + "\n")

			// Some rclone output goes to stderr, try to parse it too
			e.parseProgressLine(line, transfer)
		}
	}()

	// Wait for command completion
	err = cmd.Wait()
	wg.Wait()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// Create result
	result := &TransferResult{
		TransferID:  req.ID,
		Engine:      e.GetName(),
		Source:      req.Source,
		Destination: req.Destination,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    duration,
		Success:     err == nil,
		Metadata: map[string]interface{}{
			"command_output": output.String(),
			"command_args":   args,
			"rclone_stats":   transfer.stats,
		},
	}

	if err != nil {
		result.Error = fmt.Errorf("rclone failed: %w, stderr: %s", err, errorOutput.String())
		transferError = result.Error
	}

	// Update final progress
	transfer.mu.Lock()
	if transferError == nil {
		transfer.progress.Percentage = 100.0
	}
	if transfer.callback != nil {
		transfer.callback(*transfer.progress)
	}

	// Get final statistics from rclone stats
	if transfer.stats.TotalBytes > 0 {
		result.BytesTransferred = transfer.stats.Bytes
		result.FilesTransferred = int(transfer.stats.Transfers)
	} else {
		result.BytesTransferred = transfer.progress.BytesTransferred
	}
	transfer.mu.Unlock()

	// Calculate statistics
	if result.Success && duration.Seconds() > 0 {
		result.AverageSpeed = int64(float64(result.BytesTransferred) / duration.Seconds())
	}

	return result, transferError
}

// monitorProgress monitors rclone output for progress information
func (e *RcloneEngine) monitorProgress(reader io.Reader, transfer *rcloneTransfer) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		e.parseProgressLine(line, transfer)
	}
}

// parseProgressLine parses a line of rclone output for progress information
func (e *RcloneEngine) parseProgressLine(line string, transfer *rcloneTransfer) {
	transfer.mu.Lock()
	defer transfer.mu.Unlock()

	// Try to parse JSON stats output
	if strings.HasPrefix(line, "{") {
		var stats RcloneStats
		if err := json.Unmarshal([]byte(line), &stats); err == nil {
			transfer.stats = &stats

			// Update progress from stats
			if stats.TotalBytes > 0 {
				transfer.progress.TotalBytes = stats.TotalBytes
				transfer.progress.BytesTransferred = stats.Bytes
				transfer.progress.Percentage = float64(stats.Bytes) / float64(stats.TotalBytes) * 100.0
			}

			if stats.Speed > 0 {
				transfer.progress.Speed = int64(stats.Speed)

				// Calculate ETA
				if stats.TotalBytes > stats.Bytes && stats.Speed > 0 {
					remaining := stats.TotalBytes - stats.Bytes
					transfer.progress.ETA = time.Duration(float64(remaining)/stats.Speed) * time.Second
				}
			}

			transfer.progress.LastUpdate = time.Now()

			// Call progress callback
			if transfer.callback != nil {
				transfer.callback(*transfer.progress)
			}

			return
		}
	}

	// Fallback: try to parse traditional rclone progress output
	// Example: "Transferred:   	  234.567 MiB / 1.234 GiB, 19%, 12.345 MiB/s, ETA 1m23s"
	transferredRegex := regexp.MustCompile(`Transferred:\s+([0-9.]+)\s+(\w+)\s+/\s+([0-9.]+)\s+(\w+),\s+(\d+)%,\s+([0-9.]+)\s+(\w+)/s`)

	if matches := transferredRegex.FindStringSubmatch(line); len(matches) == 8 {
		transferred, _ := strconv.ParseFloat(matches[1], 64)
		transferredUnit := matches[2]
		total, _ := strconv.ParseFloat(matches[3], 64)
		totalUnit := matches[4]
		percentage, _ := strconv.ParseFloat(matches[5], 64)
		speed, _ := strconv.ParseFloat(matches[6], 64)
		speedUnit := matches[7]

		// Convert to bytes
		transferredBytes := convertToBytes(transferred, transferredUnit)
		totalBytes := convertToBytes(total, totalUnit)
		speedBytes := convertToBytes(speed, speedUnit)

		transfer.progress.BytesTransferred = transferredBytes
		transfer.progress.TotalBytes = totalBytes
		transfer.progress.Percentage = percentage
		transfer.progress.Speed = speedBytes
		transfer.progress.LastUpdate = time.Now()

		// Calculate ETA
		if speedBytes > 0 && totalBytes > transferredBytes {
			remaining := totalBytes - transferredBytes
			transfer.progress.ETA = time.Duration(float64(remaining)/float64(speedBytes)) * time.Second
		}

		// Call progress callback
		if transfer.callback != nil {
			transfer.callback(*transfer.progress)
		}
	}
}

// convertToBytes converts a value with unit to bytes
func convertToBytes(value float64, unit string) int64 {
	switch strings.ToUpper(unit) {
	case "B", "BYTES":
		return int64(value)
	case "K", "KB", "KIB":
		return int64(value * 1024)
	case "M", "MB", "MIB":
		return int64(value * 1024 * 1024)
	case "G", "GB", "GIB":
		return int64(value * 1024 * 1024 * 1024)
	case "T", "TB", "TIB":
		return int64(value * 1024 * 1024 * 1024 * 1024)
	default:
		return int64(value)
	}
}

// GetProgress returns current progress for an active transfer
func (e *RcloneEngine) GetProgress(ctx context.Context, transferID string) (*TransferProgress, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	transfer, exists := e.activeTransfers[transferID]
	if !exists {
		return nil, fmt.Errorf("transfer not found: %s", transferID)
	}

	transfer.mu.RLock()
	defer transfer.mu.RUnlock()

	// Return a copy of the progress
	progress := *transfer.progress
	return &progress, nil
}

// Cancel cancels an active transfer
func (e *RcloneEngine) Cancel(ctx context.Context, transferID string) error {
	e.mu.RLock()
	transfer, exists := e.activeTransfers[transferID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("transfer not found: %s", transferID)
	}

	// Kill the process
	if transfer.cmd != nil && transfer.cmd.Process != nil {
		if err := transfer.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill rclone process: %w", err)
		}
	}

	return nil
}

// Validate validates the rclone engine configuration
func (e *RcloneEngine) Validate() error {
	if e.executablePath == "" {
		return fmt.Errorf("executable path not set")
	}

	// Check if executable exists and is executable
	if _, err := exec.LookPath(e.executablePath); err != nil {
		return fmt.Errorf("rclone executable not found: %w", err)
	}

	return nil
}

// GetActiveTransfers returns all active transfers for this engine
func (e *RcloneEngine) GetActiveTransfers() map[string]*TransferProgress {
	e.mu.RLock()
	defer e.mu.RUnlock()

	transfers := make(map[string]*TransferProgress)
	for id, transfer := range e.activeTransfers {
		transfer.mu.RLock()
		progress := *transfer.progress
		transfer.mu.RUnlock()
		transfers[id] = &progress
	}

	return transfers
}

// SetExecutablePath sets the path to the rclone executable
func (e *RcloneEngine) SetExecutablePath(path string) {
	e.executablePath = path
}

// SetConfigPath sets the path to the rclone configuration file
func (e *RcloneEngine) SetConfigPath(path string) {
	e.configPath = path
}

// GetExecutablePath returns the path to the rclone executable
func (e *RcloneEngine) GetExecutablePath() string {
	return e.executablePath
}

// GetConfigPath returns the path to the rclone configuration file
func (e *RcloneEngine) GetConfigPath() string {
	return e.configPath
}

// GetVersion returns the version of rclone
func (e *RcloneEngine) GetVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, e.executablePath, "version", "--check=false")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get rclone version: %w", err)
	}

	// Parse version from output
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return strings.TrimSpace(string(output)), nil
}

// ListRemotes returns a list of configured rclone remotes
func (e *RcloneEngine) ListRemotes(ctx context.Context) ([]string, error) {
	args := []string{"listremotes"}
	if e.configPath != "" {
		args = append([]string{"--config", e.configPath}, args...)
	}

	cmd := exec.CommandContext(ctx, e.executablePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list rclone remotes: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	remotes := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Remove trailing colon if present
			if strings.HasSuffix(line, ":") {
				line = line[:len(line)-1]
			}
			remotes = append(remotes, line)
		}
	}

	return remotes, nil
}

// TestRemote tests connectivity to a specific remote
func (e *RcloneEngine) TestRemote(ctx context.Context, remote string) error {
	args := []string{"lsd", remote}
	if e.configPath != "" {
		args = append([]string{"--config", e.configPath}, args...)
	}

	cmd := exec.CommandContext(ctx, e.executablePath, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to test remote %s: %w", remote, err)
	}

	return nil
}
