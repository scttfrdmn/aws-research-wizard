package data

import (
	"time"
)

// ResearchDomainProfile contains comprehensive optimization settings for a research domain
type ResearchDomainProfile struct {
	Name                 string                  `json:"name"`
	Description          string                  `json:"description"`
	FileTypeHints        map[string]FileTypeHint `json:"file_type_hints"`
	TransferOptimization TransferOptimization    `json:"transfer_optimization"`
	StorageOptimization  StorageOptimization     `json:"storage_optimization"`
	BundlingStrategy     DomainBundlingStrategy  `json:"bundling_strategy"`
	QualityChecks        []QualityCheck          `json:"quality_checks"`
	CostOptimization     DomainCostOptimization  `json:"cost_optimization"`
	SecurityRequirements SecurityRequirements    `json:"security_requirements"`
	ComplianceSettings   ComplianceSettings      `json:"compliance_settings"`
}

// FileTypeHint provides domain-specific guidance for file types
type FileTypeHint struct {
	Description           string            `json:"description"`
	TypicalSize           string            `json:"typical_size"`
	CompressionRatio      float64           `json:"compression_ratio"`
	AccessPattern         string            `json:"access_pattern"`
	PreferredStorageClass string            `json:"preferred_storage_class"`
	PreferredEngine       string            `json:"preferred_engine"`
	SpecialHandling       map[string]string `json:"special_handling"`
}

// TransferOptimization contains transfer-specific optimizations
type TransferOptimization struct {
	PreferredEngines   []string          `json:"preferred_engines"`
	OptimalConcurrency int               `json:"optimal_concurrency"`
	OptimalPartSize    string            `json:"optimal_part_size"`
	BandwidthSettings  BandwidthSettings `json:"bandwidth_settings"`
	ProgressMonitoring bool              `json:"progress_monitoring"`
	VerificationLevel  string            `json:"verification_level"`
	RetryStrategy      RetryStrategy     `json:"retry_strategy"`
}

// BandwidthSettings for domain-specific network optimization
type BandwidthSettings struct {
	MaxUploadSpeed   string            `json:"max_upload_speed"`
	MaxDownloadSpeed string            `json:"max_download_speed"`
	TimeBasedLimits  map[string]string `json:"time_based_limits"` // e.g., "09:00-17:00": "100MB/s"
}

// RetryStrategy defines how retries should be handled
type RetryStrategy struct {
	MaxAttempts     int           `json:"max_attempts"`
	BackoffType     string        `json:"backoff_type"` // "exponential", "linear", "fixed"
	InitialDelay    time.Duration `json:"initial_delay"`
	MaxDelay        time.Duration `json:"max_delay"`
	RetryableErrors []string      `json:"retryable_errors"`
}

// StorageOptimization contains S3 storage optimizations
type StorageOptimization struct {
	DefaultStorageClass  string              `json:"default_storage_class"`
	LifecycleRules       []LifecycleRule     `json:"lifecycle_rules"`
	IntelligentTiering   bool                `json:"intelligent_tiering"`
	CrossRegionSettings  CrossRegionSettings `json:"cross_region_settings"`
	ArchiveStrategy      ArchiveStrategy     `json:"archive_strategy"`
	MetadataRequirements map[string]string   `json:"metadata_requirements"`
}

// LifecycleRule defines automated storage transitions
type LifecycleRule struct {
	Name                string   `json:"name"`
	TransitionToIA      string   `json:"transition_to_ia"`      // e.g., "30d"
	TransitionToGlacier string   `json:"transition_to_glacier"` // e.g., "90d"
	TransitionToDA      string   `json:"transition_to_da"`      // e.g., "365d"
	ExpirationDays      string   `json:"expiration_days"`       // e.g., "2555d" (7 years)
	FilePatterns        []string `json:"file_patterns"`         // e.g., ["*.fastq.gz", "*.bam"]
}

// CrossRegionSettings for multi-region optimization
type CrossRegionSettings struct {
	EnableReplication   bool     `json:"enable_replication"`
	ReplicationRegions  []string `json:"replication_regions"`
	CrossRegionAccess   bool     `json:"cross_region_access"`
	LatencyOptimization bool     `json:"latency_optimization"`
}

// ArchiveStrategy defines long-term archival approach
type ArchiveStrategy struct {
	AutoArchive        bool              `json:"auto_archive"`
	ArchiveAfterDays   int               `json:"archive_after_days"`
	ArchiveDestination string            `json:"archive_destination"`
	CompressionLevel   string            `json:"compression_level"`
	IndexingStrategy   string            `json:"indexing_strategy"`
	RetrievalStrategy  string            `json:"retrieval_strategy"`
	ArchiveMetadata    map[string]string `json:"archive_metadata"`
}

// DomainBundlingStrategy defines how files should be bundled for research domains
type DomainBundlingStrategy struct {
	EnableBundling     bool                   `json:"enable_bundling"`
	MinFilesForBundle  int                    `json:"min_files_for_bundle"`
	MaxFilesPerBundle  int                    `json:"max_files_per_bundle"`
	TargetBundleSize   string                 `json:"target_bundle_size"`
	BundlingCriteria   []string               `json:"bundling_criteria"` // "size", "type", "directory", "timestamp"
	PreservePath       bool                   `json:"preserve_path"`
	CompressionEnabled bool                   `json:"compression_enabled"`
	ToolPreference     string                 `json:"tool_preference"`
	GroupingStrategy   DomainGroupingStrategy `json:"grouping_strategy"`
}

// DomainGroupingStrategy defines how files are grouped for bundling
type DomainGroupingStrategy struct {
	Strategy          string            `json:"strategy"` // "directory", "file_type", "sample_id", "time_based"
	MaxGroupSize      string            `json:"max_group_size"`
	GroupingPattern   string            `json:"grouping_pattern"` // regex or glob pattern
	PreserveStructure bool              `json:"preserve_structure"`
	CustomRules       map[string]string `json:"custom_rules"`
}

// QualityCheck defines domain-specific validation
type QualityCheck struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // "file_integrity", "format_validation", "metadata_check"
	Command     string            `json:"command"`
	Parameters  map[string]string `json:"parameters"`
	Severity    string            `json:"severity"` // "error", "warning", "info"
	AutoFix     bool              `json:"auto_fix"`
	Description string            `json:"description"`
}

// DomainCostOptimization contains cost-specific optimizations
type DomainCostOptimization struct {
	BudgetLimits        BudgetLimits   `json:"budget_limits"`
	CostMonitoring      CostMonitoring `json:"cost_monitoring"`
	OptimizationTargets []string       `json:"optimization_targets"`
	CostAlerts          []CostAlert    `json:"cost_alerts"`
	UsagePatterns       UsagePatterns  `json:"usage_patterns"`
}

// BudgetLimits define spending constraints
type BudgetLimits struct {
	MonthlyLimit  float64 `json:"monthly_limit"`
	DailyLimit    float64 `json:"daily_limit"`
	PerJobLimit   float64 `json:"per_job_limit"`
	StorageLimit  float64 `json:"storage_limit"`
	TransferLimit float64 `json:"transfer_limit"`
}

// CostMonitoring defines cost tracking preferences
type CostMonitoring struct {
	Enabled              bool     `json:"enabled"`
	GranularityLevel     string   `json:"granularity_level"` // "daily", "hourly", "per_operation"
	CostBreakdownByType  bool     `json:"cost_breakdown_by_type"`
	PredictiveBudgeting  bool     `json:"predictive_budgeting"`
	CostOptimizationTips bool     `json:"cost_optimization_tips"`
	ReportingFrequency   string   `json:"reporting_frequency"`
	ReportRecipients     []string `json:"report_recipients"`
}

// CostAlert defines cost-based alerting
type CostAlert struct {
	Name          string   `json:"name"`
	Threshold     float64  `json:"threshold"`
	ThresholdType string   `json:"threshold_type"` // "percentage", "absolute"
	Frequency     string   `json:"frequency"`      // "daily", "weekly", "real_time"
	Channels      []string `json:"channels"`       // "email", "slack", "webhook"
}

// UsagePatterns help predict and optimize costs
type UsagePatterns struct {
	TypicalDataVolume     string            `json:"typical_data_volume"`
	SeasonalPatterns      map[string]string `json:"seasonal_patterns"`
	AccessFrequency       string            `json:"access_frequency"`
	RetentionRequirements string            `json:"retention_requirements"`
	GrowthProjections     map[string]string `json:"growth_projections"`
}

// SecurityRequirements define domain-specific security needs
type SecurityRequirements struct {
	EncryptionRequired     bool              `json:"encryption_required"`
	EncryptionStandards    []string          `json:"encryption_standards"`
	AccessControls         AccessControls    `json:"access_controls"`
	AuditRequirements      AuditRequirements `json:"audit_requirements"`
	DataClassification     string            `json:"data_classification"`
	GeographicRestrictions []string          `json:"geographic_restrictions"`
}

// AccessControls define who can access the data
type AccessControls struct {
	RequiresMFA         bool     `json:"requires_mfa"`
	AllowedPrincipals   []string `json:"allowed_principals"`
	RequiredPermissions []string `json:"required_permissions"`
	SessionTimeout      string   `json:"session_timeout"`
	IPRestrictions      []string `json:"ip_restrictions"`
}

// AuditRequirements define compliance auditing needs
type AuditRequirements struct {
	AuditLogging         bool     `json:"audit_logging"`
	LogRetentionPeriod   string   `json:"log_retention_period"`
	ComplianceFrameworks []string `json:"compliance_frameworks"`
	AuditReporting       bool     `json:"audit_reporting"`
	AccessLogging        bool     `json:"access_logging"`
	ChangeTracking       bool     `json:"change_tracking"`
}

// ComplianceSettings for regulatory requirements
type ComplianceSettings struct {
	Framework           string            `json:"framework"` // "HIPAA", "GDPR", "SOX", "FedRAMP"
	DataResidency       []string          `json:"data_residency"`
	RetentionPolicies   map[string]string `json:"retention_policies"`
	DeletionPolicies    map[string]string `json:"deletion_policies"`
	ConsentManagement   bool              `json:"consent_management"`
	DataLineage         bool              `json:"data_lineage"`
	ComplianceReporting bool              `json:"compliance_reporting"`
}

// ResearchDomainProfileManager manages research domain optimization profiles
type ResearchDomainProfileManager struct {
	profiles map[string]*ResearchDomainProfile
}

// NewResearchDomainProfileManager creates a new domain profile manager with built-in profiles
func NewResearchDomainProfileManager() *ResearchDomainProfileManager {
	dpm := &ResearchDomainProfileManager{
		profiles: make(map[string]*ResearchDomainProfile),
	}

	// Initialize built-in domain profiles
	dpm.initializeBuiltInProfiles()

	return dpm
}

// GetProfile returns a domain profile by name
func (dpm *ResearchDomainProfileManager) GetProfile(domain string) (*ResearchDomainProfile, bool) {
	profile, exists := dpm.profiles[domain]
	return profile, exists
}

// GetAllProfiles returns all available domain profiles
func (dpm *ResearchDomainProfileManager) GetAllProfiles() map[string]*ResearchDomainProfile {
	return dpm.profiles
}

// AddProfile adds a custom domain profile
func (dpm *ResearchDomainProfileManager) AddProfile(name string, profile *ResearchDomainProfile) {
	dpm.profiles[name] = profile
}

// initializeBuiltInProfiles creates the standard research domain profiles
func (dpm *ResearchDomainProfileManager) initializeBuiltInProfiles() {
	// Genomics domain profile
	dpm.profiles["genomics"] = dpm.createGenomicsProfile()

	// Climate science domain profile
	dpm.profiles["climate"] = dpm.createClimateProfile()

	// Machine learning domain profile
	dpm.profiles["machine_learning"] = dpm.createMachineLearningProfile()

	// Astronomy domain profile
	dpm.profiles["astronomy"] = dpm.createAstronomyProfile()

	// Geospatial domain profile
	dpm.profiles["geospatial"] = dpm.createGeospatialProfile()

	// Chemistry domain profile
	dpm.profiles["chemistry"] = dpm.createChemistryProfile()

	// Physics domain profile
	dpm.profiles["physics"] = dpm.createPhysicsProfile()

	// Materials science domain profile
	dpm.profiles["materials_science"] = dpm.createMaterialsScienceProfile()
}
