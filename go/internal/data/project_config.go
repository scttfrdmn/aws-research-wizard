package data

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// ProjectConfig represents a complete declarative configuration for a research data project
type ProjectConfig struct {
	// Project metadata
	Project      ProjectInfo              `yaml:"project"`
	
	// Data definitions
	DataProfiles map[string]DataProfile   `yaml:"data_profiles"`
	Destinations map[string]Destination   `yaml:"destinations"`
	
	// Workflows
	Workflows    []Workflow               `yaml:"workflows"`
	
	// Configuration
	Settings     ProjectSettings          `yaml:"settings"`
	
	// Optimization
	Optimization OptimizationSettings    `yaml:"optimization"`
	
	// Monitoring
	Monitoring   MonitoringSettings       `yaml:"monitoring"`
}

// ProjectInfo contains basic project information
type ProjectInfo struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Domain      string            `yaml:"domain"`  // e.g., "genomics", "climate", "ml"
	Owner       string            `yaml:"owner"`
	Budget      string            `yaml:"budget"`  // e.g., "$500/month"
	Tags        map[string]string `yaml:"tags"`
	CreatedAt   time.Time         `yaml:"created_at"`
	UpdatedAt   time.Time         `yaml:"updated_at"`
}

// DataProfile describes a data source or collection
type DataProfile struct {
	Name          string            `yaml:"name"`
	Description   string            `yaml:"description"`
	Path          string            `yaml:"path"`           // Local path or pattern
	FileCount     int64             `yaml:"file_count,omitempty"`
	TotalSize     string            `yaml:"total_size,omitempty"`     // e.g., "1.5TB"
	AvgFileSize   string            `yaml:"avg_file_size,omitempty"`  // e.g., "500MB"
	FileTypes     []string          `yaml:"file_types,omitempty"`     // e.g., [".fastq", ".bam"]
	AccessPattern string            `yaml:"access_pattern"`           // "write_once_read_many", "frequent_access", "archival"
	Priority      string            `yaml:"priority"`                 // "high", "medium", "low"
	Retention     string            `yaml:"retention,omitempty"`      // e.g., "7_years"
	Metadata      map[string]string `yaml:"metadata,omitempty"`
}

// Destination describes where data should be stored
type Destination struct {
	Name         string            `yaml:"name"`
	URI          string            `yaml:"uri"`              // e.g., "s3://bucket/prefix/"
	StorageClass string            `yaml:"storage_class"`    // "STANDARD", "IA", "GLACIER", etc.
	Goals        []string          `yaml:"goals"`            // ["cost_optimize", "fast_access", "long_term_storage"]
	Region       string            `yaml:"region,omitempty"`
	Encryption   EncryptionConfig  `yaml:"encryption,omitempty"`
	Lifecycle    LifecycleConfig   `yaml:"lifecycle,omitempty"`
	Metadata     map[string]string `yaml:"metadata,omitempty"`
}

// EncryptionConfig defines encryption settings
type EncryptionConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Type      string `yaml:"type"`      // "SSE-S3", "SSE-KMS", "SSE-C"
	KeyID     string `yaml:"key_id,omitempty"`
	Algorithm string `yaml:"algorithm,omitempty"`
}

// LifecycleConfig defines lifecycle management rules
type LifecycleConfig struct {
	TransitionIA     string `yaml:"transition_ia,omitempty"`     // e.g., "30_days"
	TransitionGlacier string `yaml:"transition_glacier,omitempty"` // e.g., "90_days"
	TransitionDA     string `yaml:"transition_da,omitempty"`     // e.g., "365_days"
	Expiration       string `yaml:"expiration,omitempty"`        // e.g., "7_years"
	IncompleteUploads string `yaml:"incomplete_uploads,omitempty"` // e.g., "7_days"
}

// Workflow defines a data movement or processing workflow
type Workflow struct {
	Name          string                 `yaml:"name"`
	Description   string                 `yaml:"description"`
	Source        string                 `yaml:"source"`       // Reference to DataProfile
	Destination   string                 `yaml:"destination"`  // Reference to Destination
	Engine        string                 `yaml:"engine,omitempty"`    // "s5cmd", "rclone", "auto"
	Triggers      []string               `yaml:"triggers"`     // ["manual", "schedule", "file_watcher"]
	Schedule      string                 `yaml:"schedule,omitempty"`  // Cron expression
	PreProcessing []ProcessingStep       `yaml:"preprocessing,omitempty"`
	PostProcessing []ProcessingStep      `yaml:"postprocessing,omitempty"`
	Configuration WorkflowConfiguration  `yaml:"configuration,omitempty"`
	Enabled       bool                   `yaml:"enabled"`
}

// ProcessingStep defines a processing step in a workflow
type ProcessingStep struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"`        // "bundle", "compress", "validate", "transform"
	Parameters  map[string]string `yaml:"parameters,omitempty"`
	Condition   string            `yaml:"condition,omitempty"`   // Conditional execution
	OnFailure   string            `yaml:"on_failure,omitempty"`  // "stop", "continue", "retry"
}

// WorkflowConfiguration defines workflow-specific settings
type WorkflowConfiguration struct {
	Concurrency      int               `yaml:"concurrency,omitempty"`
	PartSize         string            `yaml:"part_size,omitempty"`
	RetryAttempts    int               `yaml:"retry_attempts,omitempty"`
	Timeout          string            `yaml:"timeout,omitempty"`
	BandwidthLimit   string            `yaml:"bandwidth_limit,omitempty"`
	Checksum         bool              `yaml:"checksum"`
	OverwritePolicy  string            `yaml:"overwrite_policy"`  // "always", "never", "if_newer"
	FailurePolicy    string            `yaml:"failure_policy"`    // "stop", "continue", "retry"
	NotificationURL  string            `yaml:"notification_url,omitempty"`
	CustomParameters map[string]string `yaml:"custom_parameters,omitempty"`
}

// ProjectSettings contains global project settings
type ProjectSettings struct {
	DefaultRegion     string            `yaml:"default_region"`
	DefaultEngine     string            `yaml:"default_engine"`
	WorkingDirectory  string            `yaml:"working_directory"`
	LogLevel          string            `yaml:"log_level"`
	ConfigDirectory   string            `yaml:"config_directory"`
	CacheDirectory    string            `yaml:"cache_directory"`
	TempDirectory     string            `yaml:"temp_directory"`
	MaxConcurrent     int               `yaml:"max_concurrent_workflows"`
	GlobalTags        map[string]string `yaml:"global_tags,omitempty"`
}

// OptimizationSettings defines optimization preferences
type OptimizationSettings struct {
	EnableAutoOptimization bool                   `yaml:"enable_auto_optimization"`
	CostOptimization       CostOptimizationConfig `yaml:"cost_optimization"`
	PerformanceOptimization PerformanceOptimizationConfig `yaml:"performance_optimization"`
	ReliabilityOptimization ReliabilityOptimizationConfig `yaml:"reliability_optimization"`
}

// CostOptimizationConfig defines cost optimization settings
type CostOptimizationConfig struct {
	Enabled               bool     `yaml:"enabled"`
	BudgetLimit           string   `yaml:"budget_limit,omitempty"`          // e.g., "$1000/month"
	AutoBundleSmallFiles  bool     `yaml:"auto_bundle_small_files"`
	AutoCompression       bool     `yaml:"auto_compression"`
	AutoStorageClass      bool     `yaml:"auto_storage_class_optimization"`
	AutoLifecyclePolicies bool     `yaml:"auto_lifecycle_policies"`
	CostAlerts            []string `yaml:"cost_alerts,omitempty"`           // ["50%", "80%", "100%"]
}

// PerformanceOptimizationConfig defines performance optimization settings
type PerformanceOptimizationConfig struct {
	Enabled              bool   `yaml:"enabled"`
	AutoConcurrency      bool   `yaml:"auto_concurrency_tuning"`
	AutoPartSize         bool   `yaml:"auto_part_size_tuning"`
	AutoEngineSelection  bool   `yaml:"auto_engine_selection"`
	PrefetchData         bool   `yaml:"prefetch_data"`
	NetworkOptimization  bool   `yaml:"network_optimization"`
	MaxTransferSpeed     string `yaml:"max_transfer_speed,omitempty"`  // e.g., "1Gbps"
}

// ReliabilityOptimizationConfig defines reliability optimization settings
type ReliabilityOptimizationConfig struct {
	Enabled              bool   `yaml:"enabled"`
	AutoRetry            bool   `yaml:"auto_retry"`
	AutoVerification     bool   `yaml:"auto_verification"`
	AutoBackup           bool   `yaml:"auto_backup"`
	AutoVersioning       bool   `yaml:"auto_versioning"`
	AutoReplication      bool   `yaml:"auto_replication"`
	MaxRetryAttempts     int    `yaml:"max_retry_attempts"`
	RetryDelay           string `yaml:"retry_delay"`               // e.g., "5s", "1m"
	HealthCheckInterval  string `yaml:"health_check_interval"`    // e.g., "5m"
}

// MonitoringSettings defines monitoring and alerting configuration
type MonitoringSettings struct {
	Enabled           bool              `yaml:"enabled"`
	DashboardEnabled  bool              `yaml:"dashboard_enabled"`
	MetricsCollection bool              `yaml:"metrics_collection"`
	LogCollection     bool              `yaml:"log_collection"`
	AlertsEnabled     bool              `yaml:"alerts_enabled"`
	NotificationChannels []NotificationChannel `yaml:"notification_channels,omitempty"`
	Metrics           MetricsConfig     `yaml:"metrics"`
	Alerts            AlertsConfig      `yaml:"alerts"`
}

// NotificationChannel defines a notification endpoint
type NotificationChannel struct {
	Name     string            `yaml:"name"`
	Type     string            `yaml:"type"`     // "email", "slack", "webhook", "sns"
	Endpoint string            `yaml:"endpoint"`
	Enabled  bool              `yaml:"enabled"`
	Settings map[string]string `yaml:"settings,omitempty"`
}

// MetricsConfig defines metrics collection settings
type MetricsConfig struct {
	CollectionInterval string   `yaml:"collection_interval"`  // e.g., "1m"
	RetentionPeriod    string   `yaml:"retention_period"`     // e.g., "90d"
	CustomMetrics      []string `yaml:"custom_metrics,omitempty"`
	ExportToCloudWatch bool     `yaml:"export_to_cloudwatch"`
	ExportToPrometheus bool     `yaml:"export_to_prometheus"`
}

// AlertsConfig defines alerting rules
type AlertsConfig struct {
	TransferFailure      bool   `yaml:"transfer_failure"`
	HighCost             bool   `yaml:"high_cost"`
	SlowPerformance      bool   `yaml:"slow_performance"`
	LargeQueueSize       bool   `yaml:"large_queue_size"`
	DiskSpaceLow         bool   `yaml:"disk_space_low"`
	CostThreshold        string `yaml:"cost_threshold,omitempty"`        // e.g., "$100/day"
	PerformanceThreshold string `yaml:"performance_threshold,omitempty"` // e.g., "10MB/s"
	QueueSizeThreshold   int    `yaml:"queue_size_threshold,omitempty"`
}

// ProjectConfigManager manages project configurations
type ProjectConfigManager struct {
	configPath string
}

// NewProjectConfigManager creates a new project config manager
func NewProjectConfigManager(configPath string) *ProjectConfigManager {
	return &ProjectConfigManager{
		configPath: configPath,
	}
}

// LoadConfig loads a project configuration from file
func (pcm *ProjectConfigManager) LoadConfig(configFile string) (*ProjectConfig, error) {
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(pcm.configPath, configFile)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults and validate
	if err := pcm.validateAndApplyDefaults(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// SaveConfig saves a project configuration to file
func (pcm *ProjectConfigManager) SaveConfig(config *ProjectConfig, configFile string) error {
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(pcm.configPath, configFile)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Update timestamps
	config.Project.UpdatedAt = time.Now()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GenerateConfig creates a configuration based on data analysis
func (pcm *ProjectConfigManager) GenerateConfig(pattern *DataPattern, recommendations *RecommendationResult) (*ProjectConfig, error) {
	config := &ProjectConfig{
		Project: ProjectInfo{
			Name:        "auto-generated-project",
			Description: "Auto-generated configuration based on data analysis",
			Domain:      inferDomain(pattern),
			Budget:      estimateBudget(recommendations),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Tags: map[string]string{
				"generated": "true",
				"tool":      "aws-research-wizard",
			},
		},
		DataProfiles: make(map[string]DataProfile),
		Destinations: make(map[string]Destination),
		Workflows:    make([]Workflow, 0),
	}

	// Generate data profile
	config.DataProfiles["main_dataset"] = DataProfile{
		Name:          "Main Dataset",
		Description:   "Primary research dataset",
		Path:          pattern.AnalyzedPath,
		FileCount:     pattern.TotalFiles,
		TotalSize:     pattern.TotalSizeHuman,
		AvgFileSize:   formatBytes(pattern.FileSizes.MeanSize),
		AccessPattern: inferAccessPattern(pattern),
		Priority:      "high",
		FileTypes:     extractFileTypes(pattern),
	}

	// Generate destinations based on recommendations
	config.Destinations = pcm.generateDestinations(pattern, recommendations)

	// Generate workflows based on recommendations
	config.Workflows = pcm.generateWorkflows(pattern, recommendations)

	// Generate settings
	config.Settings = pcm.generateSettings(pattern, recommendations)

	// Generate optimization settings
	config.Optimization = pcm.generateOptimizationSettings(pattern, recommendations)

	// Generate monitoring settings
	config.Monitoring = pcm.generateMonitoringSettings(pattern, recommendations)

	return config, nil
}

// generateDestinations creates destination configurations
func (pcm *ProjectConfigManager) generateDestinations(pattern *DataPattern, recommendations *RecommendationResult) map[string]Destination {
	destinations := make(map[string]Destination)

	// Primary destination
	storageClass := "STANDARD"
	if pattern.AccessPatterns.LikelyArchival {
		storageClass = "GLACIER"
	} else if pattern.AccessPatterns.LikelyWriteOnce {
		storageClass = "STANDARD_IA"
	}

	destinations["primary"] = Destination{
		Name:         "Primary Storage",
		URI:          "s3://your-bucket/primary/",
		StorageClass: storageClass,
		Goals:        []string{"cost_optimize", "reliable_access"},
		Encryption: EncryptionConfig{
			Enabled: true,
			Type:    "SSE-S3",
		},
		Lifecycle: LifecycleConfig{
			TransitionIA:     "30_days",
			TransitionGlacier: "90_days",
			IncompleteUploads: "7_days",
		},
	}

	// Archive destination if needed
	if pattern.AccessPatterns.LikelyArchival {
		destinations["archive"] = Destination{
			Name:         "Long-term Archive",
			URI:          "s3://your-bucket/archive/",
			StorageClass: "DEEP_ARCHIVE",
			Goals:        []string{"minimum_cost", "long_term_storage"},
			Encryption: EncryptionConfig{
				Enabled: true,
				Type:    "SSE-S3",
			},
		}
	}

	return destinations
}

// generateWorkflows creates workflow configurations
func (pcm *ProjectConfigManager) generateWorkflows(pattern *DataPattern, recommendations *RecommendationResult) []Workflow {
	var workflows []Workflow

	// Find the best engine recommendation
	engine := "auto"
	for _, toolRec := range recommendations.ToolRecommendations {
		if toolRec.Task == "primary_upload" {
			engine = toolRec.RecommendedTool
			break
		}
	}

	// Primary upload workflow
	workflow := Workflow{
		Name:        "primary_upload",
		Description: "Upload main dataset to primary storage",
		Source:      "main_dataset",
		Destination: "primary",
		Engine:      engine,
		Triggers:    []string{"manual"},
		Enabled:     true,
		Configuration: WorkflowConfiguration{
			Concurrency:     10,
			PartSize:        "32MB",
			RetryAttempts:   3,
			Timeout:         "24h",
			Checksum:        true,
			OverwritePolicy: "if_newer",
			FailurePolicy:   "stop",
		},
	}

	// Add preprocessing steps based on recommendations
	for _, suggestion := range recommendations.OptimizationSuggestions {
		switch suggestion.Type {
		case "bundling":
			workflow.PreProcessing = append(workflow.PreProcessing, ProcessingStep{
				Name: "bundle_small_files",
				Type: "bundle",
				Parameters: map[string]string{
					"target_size": "100MB",
					"tool":        "suitcase",
				},
			})
		case "compression":
			workflow.PreProcessing = append(workflow.PreProcessing, ProcessingStep{
				Name: "compress_data",
				Type: "compress",
				Parameters: map[string]string{
					"algorithm": "gzip",
				},
			})
		}
	}

	workflows = append(workflows, workflow)

	return workflows
}

// generateSettings creates project settings
func (pcm *ProjectConfigManager) generateSettings(pattern *DataPattern, recommendations *RecommendationResult) ProjectSettings {
	return ProjectSettings{
		DefaultRegion:        "us-east-1",
		DefaultEngine:        "auto",
		WorkingDirectory:     "/tmp/aws-research-wizard",
		LogLevel:            "info",
		ConfigDirectory:     "~/.aws-research-wizard/config",
		CacheDirectory:      "~/.aws-research-wizard/cache",
		TempDirectory:       "/tmp",
		MaxConcurrent:       5,
		GlobalTags: map[string]string{
			"project":     "research",
			"managed-by":  "aws-research-wizard",
		},
	}
}

// generateOptimizationSettings creates optimization settings
func (pcm *ProjectConfigManager) generateOptimizationSettings(pattern *DataPattern, recommendations *RecommendationResult) OptimizationSettings {
	// Enable optimizations based on detected patterns
	enableBundling := pattern.FileSizes.SmallFiles.CountUnder1MB > 100
	enableCompression := pcm.shouldEnableCompression(pattern)
	
	return OptimizationSettings{
		EnableAutoOptimization: true,
		CostOptimization: CostOptimizationConfig{
			Enabled:              true,
			AutoBundleSmallFiles: enableBundling,
			AutoCompression:      enableCompression,
			AutoStorageClass:     true,
			AutoLifecyclePolicies: true,
			CostAlerts:           []string{"80%", "100%"},
		},
		PerformanceOptimization: PerformanceOptimizationConfig{
			Enabled:             true,
			AutoConcurrency:     true,
			AutoPartSize:        true,
			AutoEngineSelection: true,
			NetworkOptimization: true,
		},
		ReliabilityOptimization: ReliabilityOptimizationConfig{
			Enabled:         true,
			AutoRetry:       true,
			AutoVerification: true,
			MaxRetryAttempts: 3,
			RetryDelay:      "5s",
		},
	}
}

// generateMonitoringSettings creates monitoring settings
func (pcm *ProjectConfigManager) generateMonitoringSettings(pattern *DataPattern, recommendations *RecommendationResult) MonitoringSettings {
	totalSizeGB := float64(pattern.TotalSize) / (1024 * 1024 * 1024)
	
	// Enable more monitoring for larger datasets
	enableDashboard := totalSizeGB > 100
	
	return MonitoringSettings{
		Enabled:          true,
		DashboardEnabled: enableDashboard,
		MetricsCollection: true,
		LogCollection:    true,
		AlertsEnabled:    true,
		Metrics: MetricsConfig{
			CollectionInterval: "1m",
			RetentionPeriod:   "90d",
			ExportToCloudWatch: enableDashboard,
		},
		Alerts: AlertsConfig{
			TransferFailure:   true,
			HighCost:         true,
			SlowPerformance:  true,
			CostThreshold:    "$100/day",
			PerformanceThreshold: "10MB/s",
		},
	}
}

// validateAndApplyDefaults validates configuration and applies defaults
func (pcm *ProjectConfigManager) validateAndApplyDefaults(config *ProjectConfig) error {
	// Validate required fields
	if config.Project.Name == "" {
		return fmt.Errorf("project name is required")
	}

	// Apply defaults
	if config.Settings.DefaultRegion == "" {
		config.Settings.DefaultRegion = "us-east-1"
	}
	if config.Settings.LogLevel == "" {
		config.Settings.LogLevel = "info"
	}
	if config.Settings.MaxConcurrent <= 0 {
		config.Settings.MaxConcurrent = 5
	}

	// Validate workflows
	for i, workflow := range config.Workflows {
		if workflow.Name == "" {
			return fmt.Errorf("workflow %d: name is required", i)
		}
		if workflow.Source == "" {
			return fmt.Errorf("workflow %s: source is required", workflow.Name)
		}
		if workflow.Destination == "" {
			return fmt.Errorf("workflow %s: destination is required", workflow.Name)
		}

		// Check if referenced data profile and destination exist
		if _, exists := config.DataProfiles[workflow.Source]; !exists {
			return fmt.Errorf("workflow %s: data profile '%s' not found", workflow.Name, workflow.Source)
		}
		if _, exists := config.Destinations[workflow.Destination]; !exists {
			return fmt.Errorf("workflow %s: destination '%s' not found", workflow.Name, workflow.Destination)
		}
	}

	return nil
}

// Helper functions

func inferDomain(pattern *DataPattern) string {
	if len(pattern.DomainHints.DetectedDomains) > 0 {
		return pattern.DomainHints.DetectedDomains[0]
	}
	return "general"
}

func inferAccessPattern(pattern *DataPattern) string {
	if pattern.AccessPatterns.LikelyArchival {
		return "archival"
	}
	if pattern.AccessPatterns.LikelyFreqAccess {
		return "frequent_access"
	}
	return "write_once_read_many"
}

func extractFileTypes(pattern *DataPattern) []string {
	var types []string
	for ext := range pattern.FileTypes {
		if ext != "(no extension)" {
			types = append(types, ext)
		}
	}
	return types
}

func estimateBudget(recommendations *RecommendationResult) string {
	if recommendations.CostAnalysis != nil && len(recommendations.CostAnalysis.Scenarios) > 0 {
		monthlyCost := recommendations.CostAnalysis.Scenarios[0].MonthlyCosts.Total
		budget := int(monthlyCost * 1.2) // 20% buffer
		return fmt.Sprintf("$%d/month", budget)
	}
	return "$500/month"
}

func (pcm *ProjectConfigManager) shouldEnableCompression(pattern *DataPattern) bool {
	compressibleSize := int64(0)
	for _, typeInfo := range pattern.FileTypes {
		if typeInfo.Compressible {
			compressibleSize += typeInfo.TotalSize
		}
	}
	return float64(compressibleSize)/float64(pattern.TotalSize) > 0.5
}