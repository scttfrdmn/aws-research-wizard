package data

import (
	"time"
)

// createGenomicsProfile creates the genomics research domain profile
func (dpm *ResearchDomainProfileManager) createGenomicsProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "genomics",
		Description: "Optimizations for genomics and bioinformatics research data",

		FileTypeHints: map[string]FileTypeHint{
			".fastq": {
				Description:           "Raw sequencing reads",
				TypicalSize:           "100MB-2GB",
				CompressionRatio:      2.5,
				AccessPattern:         "write_once_read_frequently",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
				SpecialHandling: map[string]string{
					"compression": "gzip_recommended",
					"bundling":    "pair_aware",
				},
			},
			".bam": {
				Description:           "Aligned sequencing data",
				TypicalSize:           "1GB-50GB",
				CompressionRatio:      1.0, // Already compressed
				AccessPattern:         "frequent_random_access",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
				SpecialHandling: map[string]string{
					"index_required": ".bai",
					"checksum":       "required",
				},
			},
			".vcf": {
				Description:           "Variant call data",
				TypicalSize:           "10MB-1GB",
				CompressionRatio:      3.0,
				AccessPattern:         "write_once_read_many",
				PreferredStorageClass: "STANDARD_IA",
				PreferredEngine:       "s5cmd",
				SpecialHandling: map[string]string{
					"compression": "bgzip_recommended",
					"indexing":    "tabix_compatible",
				},
			},
		},

		TransferOptimization: TransferOptimization{
			PreferredEngines:   []string{"s5cmd", "aws_cli"},
			OptimalConcurrency: 20,
			OptimalPartSize:    "64MB",
			BandwidthSettings: BandwidthSettings{
				MaxUploadSpeed:   "1Gbps",
				MaxDownloadSpeed: "1Gbps",
				TimeBasedLimits: map[string]string{
					"09:00-17:00": "500MB/s", // Business hours
					"17:00-09:00": "1GB/s",   // Off hours
				},
			},
			ProgressMonitoring: true,
			VerificationLevel:  "checksum",
			RetryStrategy: RetryStrategy{
				MaxAttempts:  5,
				BackoffType:  "exponential",
				InitialDelay: 5 * time.Second,
				MaxDelay:     5 * time.Minute,
				RetryableErrors: []string{
					"network_timeout",
					"throttling",
					"temporary_failure",
				},
			},
		},

		StorageOptimization: StorageOptimization{
			DefaultStorageClass: "STANDARD",
			LifecycleRules: []LifecycleRule{
				{
					Name:                "raw_sequencing_data",
					TransitionToIA:      "30d",
					TransitionToGlacier: "90d",
					TransitionToDA:      "365d",
					ExpirationDays:      "2555d", // 7 years
					FilePatterns:        []string{"*.fastq.gz", "*.fq.gz"},
				},
				{
					Name:                "processed_data",
					TransitionToIA:      "60d",
					TransitionToGlacier: "180d",
					TransitionToDA:      "730d",  // 2 years
					ExpirationDays:      "3650d", // 10 years
					FilePatterns:        []string{"*.bam", "*.vcf.gz"},
				},
			},
			IntelligentTiering: true,
			CrossRegionSettings: CrossRegionSettings{
				EnableReplication:   true,
				ReplicationRegions:  []string{"us-west-2", "eu-west-1"},
				CrossRegionAccess:   true,
				LatencyOptimization: true,
			},
			ArchiveStrategy: ArchiveStrategy{
				AutoArchive:        true,
				ArchiveAfterDays:   365,
				ArchiveDestination: "DEEP_ARCHIVE",
				CompressionLevel:   "high",
				IndexingStrategy:   "metadata_catalog",
				RetrievalStrategy:  "bulk",
				ArchiveMetadata: map[string]string{
					"retention_period": "7_years",
					"data_type":        "genomics",
					"compliance":       "NIH_guidelines",
				},
			},
			MetadataRequirements: map[string]string{
				"sample_id":           "required",
				"sequencing_platform": "recommended",
				"reference_genome":    "recommended",
				"processing_pipeline": "recommended",
			},
		},

		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:     true,
			MinFilesForBundle:  100,
			MaxFilesPerBundle:  1000,
			TargetBundleSize:   "500MB",
			BundlingCriteria:   []string{"directory", "type", "sample_id"},
			PreservePath:       true,
			CompressionEnabled: true,
			ToolPreference:     "suitcase",
			GroupingStrategy: DomainGroupingStrategy{
				Strategy:          "sample_id",
				MaxGroupSize:      "1GB",
				GroupingPattern:   ".*_([^_]+)_.*", // Extract sample ID
				PreserveStructure: true,
				CustomRules: map[string]string{
					"pair_files": "group_together", // R1/R2 pairs
					"same_lane":  "bundle_together",
				},
			},
		},

		QualityChecks: []QualityCheck{
			{
				Name:        "fastq_validation",
				Type:        "format_validation",
				Command:     "fastqc",
				Parameters:  map[string]string{"output_format": "summary"},
				Severity:    "warning",
				AutoFix:     false,
				Description: "Validate FASTQ file format and quality metrics",
			},
			{
				Name:        "bam_integrity",
				Type:        "file_integrity",
				Command:     "samtools quickcheck",
				Parameters:  map[string]string{"verbose": "true"},
				Severity:    "error",
				AutoFix:     false,
				Description: "Check BAM file integrity and index consistency",
			},
		},

		CostOptimization: DomainCostOptimization{
			BudgetLimits: BudgetLimits{
				MonthlyLimit:  5000.0,
				DailyLimit:    200.0,
				PerJobLimit:   100.0,
				StorageLimit:  2000.0,
				TransferLimit: 500.0,
			},
			CostMonitoring: CostMonitoring{
				Enabled:              true,
				GranularityLevel:     "daily",
				CostBreakdownByType:  true,
				PredictiveBudgeting:  true,
				CostOptimizationTips: true,
				ReportingFrequency:   "weekly",
				ReportRecipients:     []string{"genomics-team@research.org"},
			},
			OptimizationTargets: []string{
				"storage_cost_reduction",
				"transfer_optimization",
				"compute_efficiency",
			},
			CostAlerts: []CostAlert{
				{
					Name:          "monthly_budget_80_percent",
					Threshold:     80.0,
					ThresholdType: "percentage",
					Frequency:     "daily",
					Channels:      []string{"email", "slack"},
				},
			},
			UsagePatterns: UsagePatterns{
				TypicalDataVolume:     "10TB-100TB",
				SeasonalPatterns:      map[string]string{"Q1": "high", "Q2": "medium"},
				AccessFrequency:       "daily_analysis",
				RetentionRequirements: "7_years_minimum",
				GrowthProjections:     map[string]string{"annual": "50%"},
			},
		},

		SecurityRequirements: SecurityRequirements{
			EncryptionRequired:  true,
			EncryptionStandards: []string{"AES-256", "SSE-S3"},
			AccessControls: AccessControls{
				RequiresMFA:         true,
				AllowedPrincipals:   []string{"genomics-researchers", "data-analysts"},
				RequiredPermissions: []string{"s3:GetObject", "s3:PutObject"},
				SessionTimeout:      "8h",
				IPRestrictions:      []string{"university_network"},
			},
			AuditRequirements: AuditRequirements{
				AuditLogging:         true,
				LogRetentionPeriod:   "7_years",
				ComplianceFrameworks: []string{"NIH", "institutional_policy"},
				AuditReporting:       true,
				AccessLogging:        true,
				ChangeTracking:       true,
			},
			DataClassification:     "sensitive_research",
			GeographicRestrictions: []string{"US", "EU"},
		},

		ComplianceSettings: ComplianceSettings{
			Framework:           "NIH",
			DataResidency:       []string{"US"},
			RetentionPolicies:   map[string]string{"raw_data": "7_years", "processed_data": "10_years"},
			DeletionPolicies:    map[string]string{"temp_files": "90_days", "logs": "1_year"},
			ConsentManagement:   false, // Typically not human subjects
			DataLineage:         true,
			ComplianceReporting: true,
		},
	}
}

// createClimateProfile creates the climate science domain profile
func (dpm *ResearchDomainProfileManager) createClimateProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "climate",
		Description: "Optimizations for climate science and atmospheric research data",

		FileTypeHints: map[string]FileTypeHint{
			".nc": {
				Description:           "NetCDF climate data files",
				TypicalSize:           "100MB-10GB",
				CompressionRatio:      1.5,
				AccessPattern:         "frequent_analysis",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "rclone",
				SpecialHandling: map[string]string{
					"compression": "internal_compression",
					"chunking":    "time_series_optimal",
				},
			},
			".grib": {
				Description:           "Meteorological data in GRIB format",
				TypicalSize:           "50MB-5GB",
				CompressionRatio:      1.2,
				AccessPattern:         "time_series_access",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
			},
			".hdf5": {
				Description:           "Hierarchical scientific data",
				TypicalSize:           "500MB-50GB",
				CompressionRatio:      2.0,
				AccessPattern:         "selective_access",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
				SpecialHandling: map[string]string{
					"chunking":    "dataset_aware",
					"compression": "lz4_recommended",
				},
			},
		},

		TransferOptimization: TransferOptimization{
			PreferredEngines:   []string{"rclone", "s5cmd"},
			OptimalConcurrency: 15,
			OptimalPartSize:    "128MB",
			BandwidthSettings: BandwidthSettings{
				MaxUploadSpeed:   "800MB/s",
				MaxDownloadSpeed: "1GB/s",
			},
			ProgressMonitoring: true,
			VerificationLevel:  "checksum",
			RetryStrategy: RetryStrategy{
				MaxAttempts:  3,
				BackoffType:  "linear",
				InitialDelay: 10 * time.Second,
				MaxDelay:     2 * time.Minute,
			},
		},

		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:     true,
			MinFilesForBundle:  50,
			MaxFilesPerBundle:  200,
			TargetBundleSize:   "1GB",
			BundlingCriteria:   []string{"time_period", "variable_type"},
			PreservePath:       true,
			CompressionEnabled: false, // Climate data often pre-compressed
			ToolPreference:     "tar",
			GroupingStrategy: DomainGroupingStrategy{
				Strategy:          "time_based",
				MaxGroupSize:      "2GB",
				GroupingPattern:   ".*_(\\d{4})(\\d{2}).*", // YYYYMM pattern
				PreserveStructure: true,
			},
		},

		CostOptimization: DomainCostOptimization{
			BudgetLimits: BudgetLimits{
				MonthlyLimit: 8000.0,
				DailyLimit:   300.0,
				StorageLimit: 5000.0,
			},
			UsagePatterns: UsagePatterns{
				TypicalDataVolume:     "50TB-500TB",
				SeasonalPatterns:      map[string]string{"summer": "high", "winter": "medium"},
				AccessFrequency:       "research_cycles",
				RetentionRequirements: "indefinite",
				GrowthProjections:     map[string]string{"annual": "30%"},
			},
		},

		SecurityRequirements: SecurityRequirements{
			EncryptionRequired: false, // Often public data
			DataClassification: "public_research",
		},

		ComplianceSettings: ComplianceSettings{
			Framework:         "open_science",
			DataResidency:     []string{"global"},
			RetentionPolicies: map[string]string{"all_data": "indefinite"},
		},
	}
}

// createMachineLearningProfile creates the machine learning domain profile
func (dpm *ResearchDomainProfileManager) createMachineLearningProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "machine_learning",
		Description: "Optimizations for machine learning and AI research data",

		FileTypeHints: map[string]FileTypeHint{
			".pkl": {
				Description:           "Pickled Python objects (models, datasets)",
				TypicalSize:           "1MB-10GB",
				CompressionRatio:      2.0,
				AccessPattern:         "frequent_loading",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
			},
			".h5": {
				Description:           "HDF5 datasets and model weights",
				TypicalSize:           "100MB-100GB",
				CompressionRatio:      1.5,
				AccessPattern:         "sequential_access",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
				SpecialHandling: map[string]string{
					"chunking": "layer_aware",
					"caching":  "recommended",
				},
			},
			".model": {
				Description:           "Trained model artifacts",
				TypicalSize:           "10MB-50GB",
				CompressionRatio:      1.8,
				AccessPattern:         "versioned_access",
				PreferredStorageClass: "STANDARD_IA",
				PreferredEngine:       "s5cmd",
			},
		},

		TransferOptimization: TransferOptimization{
			PreferredEngines:   []string{"s5cmd", "aws_cli"},
			OptimalConcurrency: 25,
			OptimalPartSize:    "32MB",
			BandwidthSettings: BandwidthSettings{
				MaxUploadSpeed:   "1.5GB/s",
				MaxDownloadSpeed: "2GB/s",
			},
			ProgressMonitoring: true,
			VerificationLevel:  "basic",
		},

		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:     true,
			MinFilesForBundle:  200,
			MaxFilesPerBundle:  2000,
			TargetBundleSize:   "1GB",
			BundlingCriteria:   []string{"experiment_id", "dataset_version"},
			PreservePath:       true,
			CompressionEnabled: true,
			ToolPreference:     "suitcase",
			GroupingStrategy: DomainGroupingStrategy{
				Strategy:          "directory",
				MaxGroupSize:      "5GB",
				PreserveStructure: true,
				CustomRules: map[string]string{
					"training_data":     "separate_bundle",
					"model_checkpoints": "version_aware",
				},
			},
		},

		CostOptimization: DomainCostOptimization{
			BudgetLimits: BudgetLimits{
				MonthlyLimit: 10000.0,
				DailyLimit:   400.0,
				PerJobLimit:  200.0,
			},
			UsagePatterns: UsagePatterns{
				TypicalDataVolume:     "1TB-1PB",
				AccessFrequency:       "experiment_driven",
				RetentionRequirements: "model_lifecycle",
				GrowthProjections:     map[string]string{"annual": "100%"},
			},
		},

		SecurityRequirements: SecurityRequirements{
			EncryptionRequired:  true,
			EncryptionStandards: []string{"AES-256"},
			DataClassification:  "proprietary_research",
		},
	}
}

// createAstronomyProfile creates the astronomy domain profile
func (dpm *ResearchDomainProfileManager) createAstronomyProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "astronomy",
		Description: "Optimizations for astronomy and astrophysics research data",

		FileTypeHints: map[string]FileTypeHint{
			".fits": {
				Description:           "Flexible Image Transport System files",
				TypicalSize:           "10MB-10GB",
				CompressionRatio:      1.3,
				AccessPattern:         "archive_and_analyze",
				PreferredStorageClass: "STANDARD_IA",
				PreferredEngine:       "rclone",
				SpecialHandling: map[string]string{
					"compression": "rice_compression",
					"headers":     "preserve_metadata",
				},
			},
		},

		TransferOptimization: TransferOptimization{
			PreferredEngines:   []string{"rclone", "globus"},
			OptimalConcurrency: 10,
			OptimalPartSize:    "256MB",
		},

		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:    true,
			MinFilesForBundle: 20,
			MaxFilesPerBundle: 100,
			TargetBundleSize:  "2GB",
			BundlingCriteria:  []string{"observation_date", "instrument"},
			GroupingStrategy: DomainGroupingStrategy{
				Strategy:     "time_based",
				MaxGroupSize: "5GB",
			},
		},

		CostOptimization: DomainCostOptimization{
			UsagePatterns: UsagePatterns{
				TypicalDataVolume:     "10TB-10PB",
				AccessFrequency:       "periodic_analysis",
				RetentionRequirements: "indefinite",
			},
		},

		SecurityRequirements: SecurityRequirements{
			EncryptionRequired: false, // Often public survey data
			DataClassification: "public_research",
		},
	}
}

// createGeospatialProfile creates the geospatial domain profile
func (dpm *ResearchDomainProfileManager) createGeospatialProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "geospatial",
		Description: "Optimizations for geospatial and remote sensing data",

		FileTypeHints: map[string]FileTypeHint{
			".tif": {
				Description:           "GeoTIFF raster images",
				TypicalSize:           "50MB-5GB",
				CompressionRatio:      2.5,
				AccessPattern:         "spatial_queries",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
				SpecialHandling: map[string]string{
					"compression": "LZW_or_JPEG",
					"tiling":      "recommended",
					"pyramids":    "generate_overviews",
				},
			},
			".shp": {
				Description:           "Shapefile vector data",
				TypicalSize:           "1MB-500MB",
				CompressionRatio:      3.0,
				AccessPattern:         "frequent_queries",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
				SpecialHandling: map[string]string{
					"bundle_components": "shp_shx_dbf_prj",
				},
			},
		},

		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:    true,
			MinFilesForBundle: 100,
			MaxFilesPerBundle: 500,
			TargetBundleSize:  "800MB",
			BundlingCriteria:  []string{"geographic_region", "data_source"},
			GroupingStrategy: DomainGroupingStrategy{
				Strategy:     "directory",
				MaxGroupSize: "2GB",
				CustomRules: map[string]string{
					"shapefile_sets": "keep_together",
					"tile_pyramids":  "bundle_by_zoom",
				},
			},
		},

		CostOptimization: DomainCostOptimization{
			UsagePatterns: UsagePatterns{
				TypicalDataVolume: "5TB-500TB",
				AccessFrequency:   "project_based",
			},
		},
	}
}

// createChemistryProfile creates the chemistry domain profile
func (dpm *ResearchDomainProfileManager) createChemistryProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "chemistry",
		Description: "Optimizations for computational chemistry and molecular data",

		FileTypeHints: map[string]FileTypeHint{
			".xyz": {
				Description:           "Molecular coordinate files",
				TypicalSize:           "1KB-10MB",
				CompressionRatio:      3.0,
				AccessPattern:         "analysis_workflows",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
			},
			".pdb": {
				Description:           "Protein structure files",
				TypicalSize:           "100KB-50MB",
				CompressionRatio:      2.5,
				AccessPattern:         "reference_access",
				PreferredStorageClass: "STANDARD_IA",
				PreferredEngine:       "s5cmd",
			},
			".log": {
				Description:           "Quantum chemistry calculation logs",
				TypicalSize:           "1MB-1GB",
				CompressionRatio:      4.0,
				AccessPattern:         "post_processing",
				PreferredStorageClass: "STANDARD_IA",
				PreferredEngine:       "s5cmd",
			},
		},

		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:    true,
			MinFilesForBundle: 500,
			MaxFilesPerBundle: 5000,
			TargetBundleSize:  "500MB",
			BundlingCriteria:  []string{"calculation_type", "molecule_family"},
		},
	}
}

// createPhysicsProfile creates the physics domain profile
func (dpm *ResearchDomainProfileManager) createPhysicsProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "physics",
		Description: "Optimizations for physics simulation and experimental data",

		FileTypeHints: map[string]FileTypeHint{
			".root": {
				Description:           "ROOT physics analysis files",
				TypicalSize:           "100MB-100GB",
				CompressionRatio:      1.5,
				AccessPattern:         "analysis_chains",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
			},
		},

		TransferOptimization: TransferOptimization{
			PreferredEngines:   []string{"s5cmd", "xrootd"},
			OptimalConcurrency: 30,
			OptimalPartSize:    "128MB",
		},

		CostOptimization: DomainCostOptimization{
			UsagePatterns: UsagePatterns{
				TypicalDataVolume: "100TB-10PB",
				AccessFrequency:   "computation_intensive",
			},
		},
	}
}

// createMaterialsScienceProfile creates the materials science domain profile
func (dpm *ResearchDomainProfileManager) createMaterialsScienceProfile() *ResearchDomainProfile {
	return &ResearchDomainProfile{
		Name:        "materials_science",
		Description: "Optimizations for materials science and engineering data",

		FileTypeHints: map[string]FileTypeHint{
			".cif": {
				Description:           "Crystallographic information files",
				TypicalSize:           "1KB-100KB",
				CompressionRatio:      2.0,
				AccessPattern:         "database_queries",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
			},
			".xrd": {
				Description:           "X-ray diffraction data",
				TypicalSize:           "10KB-10MB",
				CompressionRatio:      2.5,
				AccessPattern:         "pattern_matching",
				PreferredStorageClass: "STANDARD",
				PreferredEngine:       "s5cmd",
			},
		},

		BundlingStrategy: DomainBundlingStrategy{
			EnableBundling:    true,
			MinFilesForBundle: 1000,
			MaxFilesPerBundle: 10000,
			TargetBundleSize:  "200MB",
			BundlingCriteria:  []string{"material_type", "experiment_batch"},
			GroupingStrategy: DomainGroupingStrategy{
				Strategy:     "directory",
				MaxGroupSize: "1GB",
			},
		},
	}
}
