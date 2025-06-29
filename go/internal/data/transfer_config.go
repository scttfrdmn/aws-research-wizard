package data

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// TransferConfig represents the main configuration for the transfer system
type TransferConfig struct {
	// Global settings
	PreferredEngine    string                   `yaml:"preferred_engine"`
	FallbackEngines    []string                 `yaml:"fallback_engines"`
	MaxConcurrentJobs  int                      `yaml:"max_concurrent_jobs"`
	DefaultTimeout     time.Duration            `yaml:"default_timeout"`
	
	// Engine configurations
	Engines            map[string]EngineConfig  `yaml:"engines"`
	
	// Domain-specific profiles
	DomainProfiles     map[string]DomainProfile `yaml:"domain_profiles"`
	
	// Transfer optimization rules
	OptimizationRules  []OptimizationRule       `yaml:"optimization_rules"`
	
	// Monitoring settings
	Monitoring         MonitoringConfig         `yaml:"monitoring"`
	
	// Storage settings
	Storage            StorageConfig            `yaml:"storage"`
}

// EngineConfig represents configuration for a specific transfer engine
type EngineConfig struct {
	// Engine identification
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Enabled     bool   `yaml:"enabled"`
	
	// Tool location and validation
	ExecutablePath string            `yaml:"executable_path"`
	Version        string            `yaml:"version"`
	CheckCommand   string            `yaml:"check_command"`
	
	// Performance settings
	MaxConcurrency     int               `yaml:"max_concurrency"`
	DefaultPartSize    string            `yaml:"default_part_size"`    // e.g., "16MB", "64MB"
	BandwidthLimit     string            `yaml:"bandwidth_limit"`      // e.g., "100MB/s"
	
	// Retry settings
	MaxRetries         int               `yaml:"max_retries"`
	RetryDelay         time.Duration     `yaml:"retry_delay"`
	
	// Feature toggles
	EnableResume       bool              `yaml:"enable_resume"`
	EnableCompression  bool              `yaml:"enable_compression"`
	EnableEncryption   bool              `yaml:"enable_encryption"`
	EnableValidation   bool              `yaml:"enable_validation"`
	
	// Tool-specific settings
	ToolSettings       map[string]interface{} `yaml:"tool_settings"`
	
	// Environment variables
	Environment        map[string]string `yaml:"environment"`
	
	// Priority and selection
	Priority           int               `yaml:"priority"`
	UseFor             []string          `yaml:"use_for"`   // ["upload", "download", "sync"]
	OptimalFor         []string          `yaml:"optimal_for"` // ["small_files", "large_files", "many_files"]
}

// DomainProfile represents transfer optimization settings for specific research domains
type DomainProfile struct {
	Name               string            `yaml:"name"`
	Description        string            `yaml:"description"`
	
	// Preferred engines for this domain
	PreferredEngines   []string          `yaml:"preferred_engines"`
	
	// Domain-specific optimizations
	FileTypes          []string          `yaml:"file_types"`
	TypicalFileSize    string            `yaml:"typical_file_size"`
	TypicalDataVolume  string            `yaml:"typical_data_volume"`
	
	// Transfer patterns
	TransferPatterns   TransferPatterns  `yaml:"transfer_patterns"`
	
	// Monitoring and alerts
	MonitoringProfile  string            `yaml:"monitoring_profile"`
	
	// Custom settings
	CustomSettings     map[string]interface{} `yaml:"custom_settings"`
}

// TransferPatterns describes common transfer patterns for a domain
type TransferPatterns struct {
	// Common operations
	PrimaryOperation   string   `yaml:"primary_operation"`   // "upload", "download", "sync"
	SecondaryOperation string   `yaml:"secondary_operation"`
	
	// Data characteristics
	AccessPattern      string   `yaml:"access_pattern"`      // "sequential", "random", "batch"
	CompressionRatio   float64  `yaml:"compression_ratio"`   // Expected compression ratio
	
	// Scheduling preferences
	PreferredHours     []int    `yaml:"preferred_hours"`     // Hours of day (0-23)
	AvoidPeakHours     bool     `yaml:"avoid_peak_hours"`
	
	// Quality of service
	PriorityLevel      string   `yaml:"priority_level"`      // "low", "normal", "high", "critical"
	MaxLatency         time.Duration `yaml:"max_latency"`
}

// OptimizationRule defines rules for automatic engine selection
type OptimizationRule struct {
	Name               string            `yaml:"name"`
	Description        string            `yaml:"description"`
	
	// Conditions
	Conditions         RuleConditions    `yaml:"conditions"`
	
	// Actions
	Actions            RuleActions       `yaml:"actions"`
	
	// Priority
	Priority           int               `yaml:"priority"`
	Enabled            bool              `yaml:"enabled"`
}

// RuleConditions defines conditions for optimization rules
type RuleConditions struct {
	// File characteristics
	FileSizeMin        string   `yaml:"file_size_min"`       // e.g., "100MB"
	FileSizeMax        string   `yaml:"file_size_max"`
	FileCount          int      `yaml:"file_count"`
	FileTypes          []string `yaml:"file_types"`
	
	// Transfer characteristics
	SourceProtocol     string   `yaml:"source_protocol"`
	DestProtocol       string   `yaml:"dest_protocol"`
	TransferType       string   `yaml:"transfer_type"`       // "upload", "download", "sync"
	
	// Network characteristics
	NetworkDistance    string   `yaml:"network_distance"`    // "local", "region", "cross-region", "cross-cloud"
	BandwidthAvailable string   `yaml:"bandwidth_available"`
	
	// Time constraints
	TimeOfDay          []int    `yaml:"time_of_day"`         // Hours
	DayOfWeek          []string `yaml:"day_of_week"`
	
	// Domain context
	Domain             string   `yaml:"domain"`
	
	// Custom conditions
	CustomConditions   map[string]interface{} `yaml:"custom_conditions"`
}

// RuleActions defines actions to take when conditions are met
type RuleActions struct {
	// Engine selection
	PreferEngine       string   `yaml:"prefer_engine"`
	AvoidEngines       []string `yaml:"avoid_engines"`
	
	// Performance tuning
	SetConcurrency     int      `yaml:"set_concurrency"`
	SetPartSize        string   `yaml:"set_part_size"`
	SetBandwidthLimit  string   `yaml:"set_bandwidth_limit"`
	
	// Feature controls
	EnableCompression  *bool    `yaml:"enable_compression"`
	EnableEncryption   *bool    `yaml:"enable_encryption"`
	EnableValidation   *bool    `yaml:"enable_validation"`
	
	// Scheduling
	DelayTransfer      time.Duration `yaml:"delay_transfer"`
	ScheduleFor        string        `yaml:"schedule_for"`   // Time specification
	
	// Notifications
	NotifyOnStart      bool     `yaml:"notify_on_start"`
	NotifyOnComplete   bool     `yaml:"notify_on_complete"`
	NotifyOnError      bool     `yaml:"notify_on_error"`
	
	// Custom actions
	CustomActions      map[string]interface{} `yaml:"custom_actions"`
}

// MonitoringConfig defines monitoring and alerting settings
type MonitoringConfig struct {
	// Progress reporting
	ProgressInterval   time.Duration `yaml:"progress_interval"`
	DetailedProgress   bool          `yaml:"detailed_progress"`
	
	// Performance monitoring
	CollectMetrics     bool          `yaml:"collect_metrics"`
	MetricsInterval    time.Duration `yaml:"metrics_interval"`
	
	// Alerting
	AlertThresholds    AlertThresholds `yaml:"alert_thresholds"`
	
	// Logging
	LogLevel           string        `yaml:"log_level"`
	LogFormat          string        `yaml:"log_format"`
	LogFile            string        `yaml:"log_file"`
	
	// Dashboard
	EnableDashboard    bool          `yaml:"enable_dashboard"`
	DashboardPort      int           `yaml:"dashboard_port"`
}

// AlertThresholds defines thresholds for various alerts
type AlertThresholds struct {
	// Performance alerts
	SlowTransferSpeed  string        `yaml:"slow_transfer_speed"`  // e.g., "1MB/s"
	HighErrorRate      float64       `yaml:"high_error_rate"`      // Percentage
	LongDuration       time.Duration `yaml:"long_duration"`
	
	// Resource alerts
	HighCPUUsage       float64       `yaml:"high_cpu_usage"`       // Percentage
	HighMemoryUsage    float64       `yaml:"high_memory_usage"`    // Percentage
	HighDiskUsage      float64       `yaml:"high_disk_usage"`      // Percentage
	
	// Queue alerts
	LongQueue          int           `yaml:"long_queue"`           // Number of jobs
	StuckJobs          time.Duration `yaml:"stuck_jobs"`
}

// StorageConfig defines storage settings for the transfer system
type StorageConfig struct {
	// Temporary storage
	TempDirectory      string        `yaml:"temp_directory"`
	CleanupInterval    time.Duration `yaml:"cleanup_interval"`
	MaxTempSize        string        `yaml:"max_temp_size"`       // e.g., "100GB"
	
	// Metadata storage
	MetadataDirectory  string        `yaml:"metadata_directory"`
	BackupMetadata     bool          `yaml:"backup_metadata"`
	
	// Resume data
	ResumeDirectory    string        `yaml:"resume_directory"`
	ResumeRetention    time.Duration `yaml:"resume_retention"`
	
	// Caching
	EnableCache        bool          `yaml:"enable_cache"`
	CacheDirectory     string        `yaml:"cache_directory"`
	CacheSize          string        `yaml:"cache_size"`
	CacheTTL           time.Duration `yaml:"cache_ttl"`
}

// ConfigManager manages transfer configuration
type ConfigManager struct {
	configPath string
	config     *TransferConfig
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		configPath: configPath,
	}
}

// LoadConfig loads configuration from file
func (cm *ConfigManager) LoadConfig() (*TransferConfig, error) {
	if cm.config != nil {
		return cm.config, nil
	}
	
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Create default config
		cm.config = cm.createDefaultConfig()
		if err := cm.SaveConfig(); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		return cm.config, nil
	}
	
	// Load existing config
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config TransferConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Validate and apply defaults
	if err := cm.validateAndApplyDefaults(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	cm.config = &config
	return cm.config, nil
}

// SaveConfig saves configuration to file
func (cm *ConfigManager) SaveConfig() error {
	if cm.config == nil {
		return fmt.Errorf("no config to save")
	}
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *TransferConfig {
	return cm.config
}

// UpdateConfig updates the configuration
func (cm *ConfigManager) UpdateConfig(config *TransferConfig) error {
	if err := cm.validateAndApplyDefaults(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	
	cm.config = config
	return nil
}

// createDefaultConfig creates a default configuration
func (cm *ConfigManager) createDefaultConfig() *TransferConfig {
	return &TransferConfig{
		PreferredEngine:   "s5cmd",
		FallbackEngines:   []string{"aws-cli", "built-in"},
		MaxConcurrentJobs: 5,
		DefaultTimeout:    30 * time.Minute,
		
		Engines: map[string]EngineConfig{
			"s5cmd": {
				Name:             "s5cmd",
				Type:             "s3",
				Enabled:          true,
				ExecutablePath:   "s5cmd",
				CheckCommand:     "s5cmd version",
				MaxConcurrency:   10,
				DefaultPartSize:  "16MB",
				MaxRetries:       3,
				RetryDelay:       5 * time.Second,
				EnableResume:     true,
				EnableValidation: true,
				Priority:         90,
				UseFor:          []string{"upload", "download", "sync"},
				OptimalFor:      []string{"large_files", "s3"},
				ToolSettings: map[string]interface{}{
					"dry_run": false,
					"no_sign_request": false,
				},
			},
			"rclone": {
				Name:             "rclone",
				Type:             "multi-cloud",
				Enabled:          true,
				ExecutablePath:   "rclone",
				CheckCommand:     "rclone version",
				MaxConcurrency:   8,
				DefaultPartSize:  "8MB",
				MaxRetries:       3,
				RetryDelay:       5 * time.Second,
				EnableResume:     true,
				EnableCompression: false,
				EnableValidation: true,
				Priority:         80,
				UseFor:          []string{"upload", "download", "sync"},
				OptimalFor:      []string{"multi_cloud", "many_files"},
				ToolSettings: map[string]interface{}{
					"config": "~/.config/rclone/rclone.conf",
				},
			},
			"aws-cli": {
				Name:             "aws-cli",
				Type:             "s3",
				Enabled:          true,
				ExecutablePath:   "aws",
				CheckCommand:     "aws --version",
				MaxConcurrency:   5,
				MaxRetries:       3,
				RetryDelay:       5 * time.Second,
				Priority:         70,
				UseFor:          []string{"upload", "download"},
				OptimalFor:      []string{"small_files", "fallback"},
			},
			"built-in": {
				Name:             "built-in",
				Type:             "s3",
				Enabled:          true,
				MaxConcurrency:   10,
				DefaultPartSize:  "16MB",
				MaxRetries:       3,
				RetryDelay:       5 * time.Second,
				EnableResume:     true,
				EnableValidation: true,
				Priority:         60,
				UseFor:          []string{"upload", "download"},
				OptimalFor:      []string{"reliable", "fallback"},
			},
		},
		
		DomainProfiles: map[string]DomainProfile{
			"genomics": {
				Name:             "Genomics Research",
				Description:      "Optimized for genomics data transfers",
				PreferredEngines: []string{"s5cmd", "rclone"},
				FileTypes:        []string{".fastq", ".bam", ".vcf", ".fasta"},
				TypicalFileSize:  "1GB",
				TypicalDataVolume: "100GB",
				TransferPatterns: TransferPatterns{
					PrimaryOperation:   "upload",
					SecondaryOperation: "download",
					AccessPattern:      "sequential",
					CompressionRatio:   0.3,
					PriorityLevel:      "high",
				},
			},
			"climate": {
				Name:             "Climate Science",
				Description:      "Optimized for climate data transfers",
				PreferredEngines: []string{"rclone", "s5cmd"},
				FileTypes:        []string{".nc", ".hdf5", ".grib"},
				TypicalFileSize:  "500MB",
				TypicalDataVolume: "1TB",
				TransferPatterns: TransferPatterns{
					PrimaryOperation: "sync",
					AccessPattern:    "batch",
					CompressionRatio: 0.5,
					PriorityLevel:    "normal",
				},
			},
		},
		
		OptimizationRules: []OptimizationRule{
			{
				Name:        "Large File S5cmd",
				Description: "Use s5cmd for files larger than 100MB",
				Conditions: RuleConditions{
					FileSizeMin:     "100MB",
					SourceProtocol:  "s3",
				},
				Actions: RuleActions{
					PreferEngine:    "s5cmd",
					SetConcurrency:  10,
					SetPartSize:     "32MB",
				},
				Priority: 90,
				Enabled:  true,
			},
			{
				Name:        "Multi-cloud Rclone",
				Description: "Use rclone for cross-cloud transfers",
				Conditions: RuleConditions{
					NetworkDistance: "cross-cloud",
				},
				Actions: RuleActions{
					PreferEngine:   "rclone",
					SetConcurrency: 5,
				},
				Priority: 85,
				Enabled:  true,
			},
		},
		
		Monitoring: MonitoringConfig{
			ProgressInterval: 5 * time.Second,
			DetailedProgress: true,
			CollectMetrics:   true,
			MetricsInterval:  30 * time.Second,
			LogLevel:         "info",
			LogFormat:        "json",
			AlertThresholds: AlertThresholds{
				SlowTransferSpeed: "1MB/s",
				HighErrorRate:     10.0,
				LongDuration:      2 * time.Hour,
				LongQueue:         20,
				StuckJobs:         30 * time.Minute,
			},
		},
		
		Storage: StorageConfig{
			TempDirectory:     "/tmp/aws-research-wizard",
			CleanupInterval:   1 * time.Hour,
			MaxTempSize:       "10GB",
			MetadataDirectory: "~/.aws-research-wizard/metadata",
			BackupMetadata:    true,
			ResumeDirectory:   "~/.aws-research-wizard/resume",
			ResumeRetention:   7 * 24 * time.Hour, // 7 days
			EnableCache:       true,
			CacheDirectory:    "~/.aws-research-wizard/cache",
			CacheSize:         "1GB",
			CacheTTL:          24 * time.Hour,
		},
	}
}

// validateAndApplyDefaults validates configuration and applies defaults
func (cm *ConfigManager) validateAndApplyDefaults(config *TransferConfig) error {
	// Apply defaults for missing values
	if config.MaxConcurrentJobs <= 0 {
		config.MaxConcurrentJobs = 5
	}
	
	if config.DefaultTimeout <= 0 {
		config.DefaultTimeout = 30 * time.Minute
	}
	
	// Validate engines
	if config.Engines == nil {
		config.Engines = make(map[string]EngineConfig)
	}
	
	// Validate monitoring config
	if config.Monitoring.ProgressInterval <= 0 {
		config.Monitoring.ProgressInterval = 5 * time.Second
	}
	
	// Expand tilde in paths
	if err := cm.expandPaths(config); err != nil {
		return fmt.Errorf("failed to expand paths: %w", err)
	}
	
	return nil
}

// expandPaths expands tilde (~) in configuration paths
func (cm *ConfigManager) expandPaths(config *TransferConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	// Expand storage paths
	if config.Storage.MetadataDirectory != "" {
		config.Storage.MetadataDirectory = expandPath(config.Storage.MetadataDirectory, homeDir)
	}
	if config.Storage.ResumeDirectory != "" {
		config.Storage.ResumeDirectory = expandPath(config.Storage.ResumeDirectory, homeDir)
	}
	if config.Storage.CacheDirectory != "" {
		config.Storage.CacheDirectory = expandPath(config.Storage.CacheDirectory, homeDir)
	}
	
	return nil
}

// expandPath expands ~ to home directory in a path
func expandPath(path, homeDir string) string {
	if len(path) > 0 && path[0] == '~' {
		if len(path) == 1 {
			return homeDir
		}
		if path[1] == '/' {
			return filepath.Join(homeDir, path[2:])
		}
	}
	return path
}

// GetEngineConfig gets configuration for a specific engine
func (cm *ConfigManager) GetEngineConfig(engineName string) (*EngineConfig, error) {
	if cm.config == nil {
		return nil, fmt.Errorf("config not loaded")
	}
	
	engineConfig, exists := cm.config.Engines[engineName]
	if !exists {
		return nil, fmt.Errorf("engine config not found: %s", engineName)
	}
	
	return &engineConfig, nil
}

// GetDomainProfile gets configuration for a specific domain
func (cm *ConfigManager) GetDomainProfile(domain string) (*DomainProfile, error) {
	if cm.config == nil {
		return nil, fmt.Errorf("config not loaded")
	}
	
	profile, exists := cm.config.DomainProfiles[domain]
	if !exists {
		return nil, fmt.Errorf("domain profile not found: %s", domain)
	}
	
	return &profile, nil
}