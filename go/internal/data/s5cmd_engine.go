package data

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// S5cmdEngine implements TransferEngine for s5cmd
type S5cmdEngine struct {
	*BaseTransferEngine
	executablePath  string
	activeTransfers map[string]*s5cmdTransfer
	mu              sync.RWMutex
}

// s5cmdTransfer tracks an active s5cmd transfer
type s5cmdTransfer struct {
	id       string
	cmd      *exec.Cmd
	progress *TransferProgress
	callback ProgressCallback
	done     chan bool
	mu       sync.RWMutex
}

// NewS5cmdEngine creates a new s5cmd transfer engine
func NewS5cmdEngine(executablePath string) *S5cmdEngine {
	capabilities := EngineCapabilities{
		Protocols:              []string{"s3"},
		SupportsResume:         false, // s5cmd doesn't natively support resume
		SupportsProgress:       true,
		SupportsParallel:       true,
		SupportsCompression:    false,
		SupportsEncryption:     true, // via S3 server-side encryption
		SupportsValidation:     true,
		SupportsBandwidthLimit: false,
		SupportsRetry:          true,
		OptimalFileSizeMin:     1024 * 1024,              // 1MB
		OptimalFileSizeMax:     1024 * 1024 * 1024 * 100, // 100GB
		MaxConcurrency:         50,
		CloudOptimized:         []string{"aws"},
	}

	base := NewBaseTransferEngine("s5cmd", "s3", capabilities)

	engine := &S5cmdEngine{
		BaseTransferEngine: base,
		executablePath:     executablePath,
		activeTransfers:    make(map[string]*s5cmdTransfer),
	}

	return engine
}

// IsAvailable checks if s5cmd is available and working
func (e *S5cmdEngine) IsAvailable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, e.executablePath, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("s5cmd not available: %w", err)
	}

	// Check if output contains version information
	if !strings.Contains(string(output), "s5cmd") {
		return fmt.Errorf("unexpected s5cmd version output: %s", string(output))
	}

	return nil
}

// Upload transfers data from local to S3
func (e *S5cmdEngine) Upload(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// Validate request
	if !strings.HasPrefix(req.Destination, "s3://") {
		return nil, fmt.Errorf("s5cmd upload requires S3 destination, got: %s", req.Destination)
	}

	// Build s5cmd command
	args := e.buildUploadArgs(req)

	// Execute transfer
	return e.executeTransfer(ctx, req, args)
}

// Download transfers data from S3 to local
func (e *S5cmdEngine) Download(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// Validate request
	if !strings.HasPrefix(req.Source, "s3://") {
		return nil, fmt.Errorf("s5cmd download requires S3 source, got: %s", req.Source)
	}

	// Build s5cmd command
	args := e.buildDownloadArgs(req)

	// Execute transfer
	return e.executeTransfer(ctx, req, args)
}

// Sync synchronizes data between local and S3
func (e *S5cmdEngine) Sync(ctx context.Context, req *SyncRequest) (*TransferResult, error) {
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

	// Build s5cmd sync command
	args := e.buildSyncArgs(transferReq, &req.Options)

	// Execute transfer
	return e.executeTransfer(ctx, transferReq, args)
}

// buildUploadArgs builds command arguments for upload
func (e *S5cmdEngine) buildUploadArgs(req *TransferRequest) []string {
	args := []string{"cp"}

	// Add performance options
	if req.Options.Concurrency > 0 {
		args = append(args, "--numworkers", strconv.Itoa(req.Options.Concurrency))
	}

	// Add part size if specified
	if req.Options.PartSize > 0 {
		partSizeMB := req.Options.PartSize / (1024 * 1024)
		args = append(args, "--part-size", fmt.Sprintf("%dMB", partSizeMB))
	}

	// Add dry run if specified
	if dryRun, ok := req.Options.ToolSpecific["dry_run"].(bool); ok && dryRun {
		args = append(args, "--dry-run")
	}

	// Add source and destination
	args = append(args, req.Source, req.Destination)

	return args
}

// buildDownloadArgs builds command arguments for download
func (e *S5cmdEngine) buildDownloadArgs(req *TransferRequest) []string {
	args := []string{"cp"}

	// Add performance options
	if req.Options.Concurrency > 0 {
		args = append(args, "--numworkers", strconv.Itoa(req.Options.Concurrency))
	}

	// Add part size if specified
	if req.Options.PartSize > 0 {
		partSizeMB := req.Options.PartSize / (1024 * 1024)
		args = append(args, "--part-size", fmt.Sprintf("%dMB", partSizeMB))
	}

	// Add source and destination
	args = append(args, req.Source, req.Destination)

	return args
}

// buildSyncArgs builds command arguments for sync
func (e *S5cmdEngine) buildSyncArgs(req *TransferRequest, syncOpts *SyncOptions) []string {
	args := []string{"sync"}

	// Add performance options
	if syncOpts.Concurrency > 0 {
		args = append(args, "--numworkers", strconv.Itoa(syncOpts.Concurrency))
	}

	// Add delete option
	if syncOpts.DeleteExtraneous {
		args = append(args, "--delete")
	}

	// Add dry run option
	if syncOpts.DryRun {
		args = append(args, "--dry-run")
	}

	// Add exclude patterns
	for _, pattern := range syncOpts.Exclude {
		args = append(args, "--exclude", pattern)
	}

	// Add include patterns (s5cmd uses exclude, so we need to be careful here)
	// For now, we'll handle include patterns in a future enhancement

	// Add source and destination
	args = append(args, req.Source, req.Destination)

	return args
}

// executeTransfer executes the s5cmd command and monitors progress
func (e *S5cmdEngine) executeTransfer(ctx context.Context, req *TransferRequest, args []string) (*TransferResult, error) {
	startTime := time.Now()

	// Create command with progress monitoring
	cmd := exec.CommandContext(ctx, e.executablePath, args...)

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
	transfer := &s5cmdTransfer{
		id:  req.ID,
		cmd: cmd,
		progress: &TransferProgress{
			StartTime: startTime,
		},
		callback: req.ProgressCallback,
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
		return nil, fmt.Errorf("failed to start s5cmd: %w", err)
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
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			output.WriteString(scanner.Text() + "\n")
		}
	}()

	// Monitor stderr for errors
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			errorOutput.WriteString(line + "\n")
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
		},
	}

	if err != nil {
		result.Error = fmt.Errorf("s5cmd failed: %w, stderr: %s", err, errorOutput.String())
		transferError = result.Error
	}

	// Update final progress
	transfer.mu.Lock()
	if transferError == nil {
		transfer.progress.Percentage = 100.0
		transfer.progress.BytesTransferred = transfer.progress.TotalBytes
	}
	if transfer.callback != nil {
		transfer.callback(*transfer.progress)
	}
	transfer.mu.Unlock()

	// Calculate statistics
	if result.Success {
		result.BytesTransferred = transfer.progress.BytesTransferred
		if duration.Seconds() > 0 {
			result.AverageSpeed = int64(float64(result.BytesTransferred) / duration.Seconds())
		}
	}

	return result, transferError
}

// monitorProgress monitors s5cmd output for progress information
func (e *S5cmdEngine) monitorProgress(reader io.Reader, transfer *s5cmdTransfer) {
	scanner := bufio.NewScanner(reader)

	// Regular expressions for parsing s5cmd output
	// s5cmd outputs progress in various formats, we'll try to parse what we can
	progressRegex := regexp.MustCompile(`(\d+)\s*\/\s*(\d+)`)
	speedRegex := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(B|KB|MB|GB)\/s`)

	for scanner.Scan() {
		line := scanner.Text()

		transfer.mu.Lock()

		// Try to parse progress information
		if matches := progressRegex.FindStringSubmatch(line); len(matches) == 3 {
			current, err1 := strconv.ParseInt(matches[1], 10, 64)
			total, err2 := strconv.ParseInt(matches[2], 10, 64)

			if err1 == nil && err2 == nil && total > 0 {
				transfer.progress.BytesTransferred = current
				transfer.progress.TotalBytes = total
				transfer.progress.Percentage = float64(current) / float64(total) * 100.0
			}
		}

		// Try to parse speed information
		if matches := speedRegex.FindStringSubmatch(line); len(matches) == 3 {
			speed, err := strconv.ParseFloat(matches[1], 64)
			unit := matches[2]

			if err == nil {
				var bytesPerSec int64
				switch unit {
				case "B":
					bytesPerSec = int64(speed)
				case "KB":
					bytesPerSec = int64(speed * 1024)
				case "MB":
					bytesPerSec = int64(speed * 1024 * 1024)
				case "GB":
					bytesPerSec = int64(speed * 1024 * 1024 * 1024)
				}

				transfer.progress.Speed = bytesPerSec

				// Calculate ETA
				if bytesPerSec > 0 && transfer.progress.TotalBytes > transfer.progress.BytesTransferred {
					remaining := transfer.progress.TotalBytes - transfer.progress.BytesTransferred
					transfer.progress.ETA = time.Duration(remaining/bytesPerSec) * time.Second
				}
			}
		}

		// Update last update time
		transfer.progress.LastUpdate = time.Now()

		// Call progress callback
		if transfer.callback != nil {
			transfer.callback(*transfer.progress)
		}

		transfer.mu.Unlock()
	}
}

// GetProgress returns current progress for an active transfer
func (e *S5cmdEngine) GetProgress(ctx context.Context, transferID string) (*TransferProgress, error) {
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
func (e *S5cmdEngine) Cancel(ctx context.Context, transferID string) error {
	e.mu.RLock()
	transfer, exists := e.activeTransfers[transferID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("transfer not found: %s", transferID)
	}

	// Kill the process
	if transfer.cmd != nil && transfer.cmd.Process != nil {
		if err := transfer.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill s5cmd process: %w", err)
		}
	}

	return nil
}

// Validate validates the s5cmd engine configuration
func (e *S5cmdEngine) Validate() error {
	if e.executablePath == "" {
		return fmt.Errorf("executable path not set")
	}

	// Check if executable exists and is executable
	if _, err := exec.LookPath(e.executablePath); err != nil {
		return fmt.Errorf("s5cmd executable not found: %w", err)
	}

	return nil
}

// GetActiveTransfers returns all active transfers for this engine
func (e *S5cmdEngine) GetActiveTransfers() map[string]*TransferProgress {
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

// SetExecutablePath sets the path to the s5cmd executable
func (e *S5cmdEngine) SetExecutablePath(path string) {
	e.executablePath = path
}

// GetExecutablePath returns the path to the s5cmd executable
func (e *S5cmdEngine) GetExecutablePath() string {
	return e.executablePath
}

// GetVersion returns the version of s5cmd
func (e *S5cmdEngine) GetVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, e.executablePath, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get s5cmd version: %w", err)
	}

	// Parse version from output
	version := strings.TrimSpace(string(output))
	return version, nil
}

// SupportsFeature checks if s5cmd supports a specific feature
func (e *S5cmdEngine) SupportsFeature(feature string) bool {
	caps := e.GetCapabilities()

	switch feature {
	case "resume":
		return caps.SupportsResume
	case "progress":
		return caps.SupportsProgress
	case "parallel":
		return caps.SupportsParallel
	case "compression":
		return caps.SupportsCompression
	case "encryption":
		return caps.SupportsEncryption
	case "validation":
		return caps.SupportsValidation
	case "bandwidth_limit":
		return caps.SupportsBandwidthLimit
	case "retry":
		return caps.SupportsRetry
	default:
		return false
	}
}
