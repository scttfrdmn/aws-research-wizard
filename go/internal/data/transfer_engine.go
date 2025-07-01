package data

import (
	"context"
	"fmt"
	"time"
)

// TransferEngine defines the interface for all data transfer implementations
type TransferEngine interface {
	// GetName returns the name of the transfer engine
	GetName() string

	// GetType returns the type/category of the engine (e.g., "s3", "multi-cloud", "local")
	GetType() string

	// IsAvailable checks if the engine is available and properly configured
	IsAvailable(ctx context.Context) error

	// GetCapabilities returns the capabilities supported by this engine
	GetCapabilities() EngineCapabilities

	// Upload transfers data from local to remote location
	Upload(ctx context.Context, req *TransferRequest) (*TransferResult, error)

	// Download transfers data from remote to local location
	Download(ctx context.Context, req *TransferRequest) (*TransferResult, error)

	// Sync synchronizes data between two locations
	Sync(ctx context.Context, req *SyncRequest) (*TransferResult, error)

	// GetProgress returns current progress for an active transfer
	GetProgress(ctx context.Context, transferID string) (*TransferProgress, error)

	// Cancel cancels an active transfer
	Cancel(ctx context.Context, transferID string) error

	// Validate validates the engine configuration
	Validate() error
}

// EngineCapabilities describes what features an engine supports
type EngineCapabilities struct {
	// Protocols supported (s3, gcs, azure, ftp, sftp, local, etc.)
	Protocols []string

	// Features supported
	SupportsResume         bool
	SupportsProgress       bool
	SupportsParallel       bool
	SupportsCompression    bool
	SupportsEncryption     bool
	SupportsValidation     bool
	SupportsBandwidthLimit bool
	SupportsRetry          bool

	// Performance characteristics
	OptimalFileSizeMin int64 // Bytes - below this size, engine may not be optimal
	OptimalFileSizeMax int64 // Bytes - above this size, engine may not be optimal
	MaxConcurrency     int   // Maximum recommended concurrent transfers

	// Cloud provider optimizations
	CloudOptimized []string // Cloud providers this engine is optimized for
}

// TransferRequest represents a transfer operation request
type TransferRequest struct {
	// Unique identifier for this transfer
	ID string

	// Source and destination
	Source      string // URI or path
	Destination string // URI or path

	// Transfer options
	Options TransferOptions

	// Progress callback
	ProgressCallback func(progress TransferProgress)

	// Context for cancellation
	Context context.Context
}

// SyncRequest represents a synchronization operation request
type SyncRequest struct {
	// Unique identifier for this sync
	ID string

	// Source and destination
	Source      string // URI or path
	Destination string // URI or path

	// Sync options
	Options SyncOptions

	// Progress callback
	ProgressCallback func(progress TransferProgress)

	// Context for cancellation
	Context context.Context
}

// TransferOptions contains options for transfer operations
type TransferOptions struct {
	// Transfer behavior
	Overwrite           bool
	Resume              bool
	Verify              bool
	DeleteAfterTransfer bool

	// Performance options
	Concurrency    int
	PartSize       int64
	BandwidthLimit int64 // bytes per second

	// Retry options
	MaxRetries int
	RetryDelay time.Duration

	// Compression and encryption
	Compress bool
	Encrypt  bool

	// Tool-specific options
	ToolSpecific map[string]interface{}
}

// SyncOptions contains options for sync operations
type SyncOptions struct {
	// Sync behavior
	DeleteExtraneous bool
	SkipNewer        bool
	DryRun           bool

	// Filters
	Include []string
	Exclude []string

	// Performance options
	Concurrency    int
	BandwidthLimit int64

	// Tool-specific options
	ToolSpecific map[string]interface{}
}

// TransferResult represents the result of a transfer operation
type TransferResult struct {
	// Transfer identification
	TransferID string
	Engine     string

	// Transfer details
	Source      string
	Destination string

	// Results
	Success bool
	Error   error

	// Statistics
	BytesTransferred int64
	FilesTransferred int
	StartTime        time.Time
	EndTime          time.Time
	Duration         time.Duration
	AverageSpeed     int64 // bytes per second

	// Additional metadata
	Metadata map[string]interface{}
}

// EngineRegistry manages available transfer engines
type EngineRegistry struct {
	engines map[string]TransferEngine
	config  *TransferConfig
}

// NewEngineRegistry creates a new engine registry
func NewEngineRegistry(config *TransferConfig) *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[string]TransferEngine),
		config:  config,
	}
}

// RegisterEngine registers a transfer engine
func (r *EngineRegistry) RegisterEngine(engine TransferEngine) error {
	if engine == nil {
		return fmt.Errorf("engine cannot be nil")
	}

	name := engine.GetName()
	if name == "" {
		return fmt.Errorf("engine name cannot be empty")
	}

	if err := engine.Validate(); err != nil {
		return fmt.Errorf("engine validation failed: %w", err)
	}

	r.engines[name] = engine
	return nil
}

// GetEngine retrieves a transfer engine by name
func (r *EngineRegistry) GetEngine(name string) (TransferEngine, error) {
	engine, exists := r.engines[name]
	if !exists {
		return nil, fmt.Errorf("engine not found: %s", name)
	}

	return engine, nil
}

// ListEngines returns all registered engines
func (r *EngineRegistry) ListEngines() []TransferEngine {
	engines := make([]TransferEngine, 0, len(r.engines))
	for _, engine := range r.engines {
		engines = append(engines, engine)
	}
	return engines
}

// ListAvailableEngines returns engines that are currently available
func (r *EngineRegistry) ListAvailableEngines(ctx context.Context) []TransferEngine {
	available := make([]TransferEngine, 0)
	for _, engine := range r.engines {
		if err := engine.IsAvailable(ctx); err == nil {
			available = append(available, engine)
		}
	}
	return available
}

// SelectOptimalEngine selects the best engine for a transfer based on request parameters
func (r *EngineRegistry) SelectOptimalEngine(ctx context.Context, req *TransferRequest) (TransferEngine, error) {
	available := r.ListAvailableEngines(ctx)
	if len(available) == 0 {
		return nil, fmt.Errorf("no transfer engines available")
	}

	// Apply selection logic based on:
	// 1. User preferences from config
	// 2. Protocol compatibility
	// 3. File size optimization
	// 4. Performance characteristics

	return r.selectEngineByStrategy(available, req)
}

// selectEngineByStrategy implements engine selection logic
func (r *EngineRegistry) selectEngineByStrategy(engines []TransferEngine, req *TransferRequest) (TransferEngine, error) {
	if len(engines) == 0 {
		return nil, fmt.Errorf("no engines available")
	}

	// For now, implement simple priority-based selection
	// TODO: Implement sophisticated selection based on request characteristics

	// Check if user has a preference
	if r.config != nil && r.config.PreferredEngine != "" {
		for _, engine := range engines {
			if engine.GetName() == r.config.PreferredEngine {
				return engine, nil
			}
		}
	}

	// Check protocol compatibility
	protocol := extractProtocol(req.Source)
	for _, engine := range engines {
		caps := engine.GetCapabilities()
		for _, supportedProtocol := range caps.Protocols {
			if supportedProtocol == protocol {
				return engine, nil
			}
		}
	}

	// Fallback to first available engine
	return engines[0], nil
}

// EnginePerformanceHint provides hints for engine selection
type EnginePerformanceHint struct {
	FileSize        int64
	FileCount       int
	SourceProtocol  string
	DestProtocol    string
	NetworkDistance string // "local", "region", "cross-region", "cross-cloud"
	Priority        string // "speed", "reliability", "cost"
}

// GetEngineRecommendation gets engine recommendation based on performance hints
func (r *EngineRegistry) GetEngineRecommendation(ctx context.Context, hint EnginePerformanceHint) ([]TransferEngine, error) {
	available := r.ListAvailableEngines(ctx)
	if len(available) == 0 {
		return nil, fmt.Errorf("no engines available")
	}

	scored := make([]engineScore, 0, len(available))

	for _, engine := range available {
		score := r.scoreEngine(engine, hint)
		scored = append(scored, engineScore{
			engine: engine,
			score:  score,
		})
	}

	// Sort by score (highest first)
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Return top recommendations
	recommendations := make([]TransferEngine, 0, len(scored))
	for _, s := range scored {
		recommendations = append(recommendations, s.engine)
	}

	return recommendations, nil
}

type engineScore struct {
	engine TransferEngine
	score  float64
}

// scoreEngine calculates a score for an engine based on performance hints
func (r *EngineRegistry) scoreEngine(engine TransferEngine, hint EnginePerformanceHint) float64 {
	caps := engine.GetCapabilities()
	score := 0.0

	// Protocol compatibility
	protocolSupported := false
	for _, protocol := range caps.Protocols {
		if protocol == hint.SourceProtocol {
			protocolSupported = true
			score += 10.0
			break
		}
	}
	if !protocolSupported {
		return 0.0 // Engine can't handle this protocol
	}

	// File size optimization
	if hint.FileSize > 0 {
		if caps.OptimalFileSizeMin <= hint.FileSize && hint.FileSize <= caps.OptimalFileSizeMax {
			score += 5.0
		} else if hint.FileSize < caps.OptimalFileSizeMin {
			score += 2.0 // Less optimal for small files
		} else {
			score += 3.0 // Less optimal for large files
		}
	}

	// Feature bonuses based on priority
	switch hint.Priority {
	case "speed":
		if caps.SupportsParallel {
			score += 3.0
		}
		if caps.MaxConcurrency > 5 {
			score += 2.0
		}
	case "reliability":
		if caps.SupportsRetry {
			score += 3.0
		}
		if caps.SupportsResume {
			score += 2.0
		}
		if caps.SupportsValidation {
			score += 2.0
		}
	}

	// Engine-specific bonuses
	switch engine.GetName() {
	case "s5cmd":
		if hint.SourceProtocol == "s3" || hint.DestProtocol == "s3" {
			score += 3.0 // s5cmd is optimized for S3
		}
	case "rclone":
		if hint.NetworkDistance == "cross-cloud" {
			score += 3.0 // rclone excels at cross-cloud transfers
		}
	}

	return score
}

// Utility functions

// extractProtocol extracts the protocol from a URI
func extractProtocol(uri string) string {
	// Simple protocol extraction
	// TODO: Implement proper URI parsing
	if len(uri) > 3 && uri[0:3] == "s3:" {
		return "s3"
	}
	if len(uri) > 4 && uri[0:4] == "gcs:" {
		return "gcs"
	}
	if len(uri) > 6 && uri[0:6] == "azure:" {
		return "azure"
	}
	if len(uri) > 4 && (uri[0:4] == "ftp:" || uri[0:4] == "sftp") {
		return "ftp"
	}
	if len(uri) > 7 && (uri[0:7] == "http://" || uri[0:8] == "https://") {
		return "http"
	}
	return "local"
}

// Common transfer engine base implementation
type BaseTransferEngine struct {
	name         string
	engineType   string
	capabilities EngineCapabilities
	config       map[string]interface{}
}

// NewBaseTransferEngine creates a new base transfer engine
func NewBaseTransferEngine(name, engineType string, capabilities EngineCapabilities) *BaseTransferEngine {
	return &BaseTransferEngine{
		name:         name,
		engineType:   engineType,
		capabilities: capabilities,
		config:       make(map[string]interface{}),
	}
}

// GetName returns the engine name
func (b *BaseTransferEngine) GetName() string {
	return b.name
}

// GetType returns the engine type
func (b *BaseTransferEngine) GetType() string {
	return b.engineType
}

// GetCapabilities returns the engine capabilities
func (b *BaseTransferEngine) GetCapabilities() EngineCapabilities {
	return b.capabilities
}

// SetConfig sets configuration for the engine
func (b *BaseTransferEngine) SetConfig(config map[string]interface{}) {
	b.config = config
}

// GetConfig gets configuration value
func (b *BaseTransferEngine) GetConfig(key string) (interface{}, bool) {
	value, exists := b.config[key]
	return value, exists
}
