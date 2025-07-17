package data

import "time"

// Cargoship integration types (interface definitions for future integration)

// ArchiveJob represents a data archiving job
type ArchiveJob struct {
	ProjectID    string
	Sources      []string
	Destination  string
	StorageClass string
	Metadata     map[string]string
	Tags         map[string]string
}

// CargoshipJob represents the status and details of a Cargoship job
type CargoshipJob struct {
	ID               string
	Status           string
	EstimatedCost    float64
	EstimatedTime    time.Duration
	TransferredBytes int64
	CreatedAt        time.Time
}

// RestoreRequest represents a data restoration request
type RestoreRequest struct {
	JobID       string
	Destination string
	RestoreTier string
}

// CostEstimateRequest represents a cost estimation request
type CostEstimateRequest struct {
	DataSources  []string
	StorageClass string
}

// CostEstimate represents cost estimation results
type CostEstimate struct {
	TransferCost float64
	StorageCost  float64
	TotalCost    float64
}

// StorageOptimizationRequest represents a storage optimization request
type StorageOptimizationRequest struct {
	DataSources   []string
	AccessPattern string
}

// StorageRecommendation represents storage class recommendations
type StorageRecommendation struct {
	RecommendedClass  string
	PotentialSavings  float64
	SavingsPercentage float64
}

// ArchiveSchedule represents automatic archiving schedule
type ArchiveSchedule struct {
	CronExpression string
	ProjectID      string
	StorageClass   string
}

// TransferUpdate represents real-time transfer updates
type TransferUpdate struct {
	JobID    string
	Progress float64
	Speed    int64
	ETA      time.Duration
	Status   string
}

// TransferMetrics represents detailed transfer metrics
type TransferMetrics struct {
	JobID            string
	TotalBytes       int64
	TransferredBytes int64
	AverageSpeed     int64
	Duration         time.Duration
	Efficiency       float64
}
