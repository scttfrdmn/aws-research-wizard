package intelligence

import (
	"fmt"
	"strings"

	"github.com/aws-research-wizard/go/internal/data"
)

// ResourceAnalyzer provides intelligent resource requirement analysis
type ResourceAnalyzer struct {
	instanceSpecs map[string]InstanceSpec
}

// InstanceSpec represents the specifications of an AWS instance type
type InstanceSpec struct {
	vCPUs            int     `json:"vcpus"`
	MemoryGB         float64 `json:"memory_gb"`
	NetworkSpeed     string  `json:"network_speed"`
	StorageType      string  `json:"storage_type"`
	GPUs             int     `json:"gpus"`
	GPUMemoryGB      float64 `json:"gpu_memory_gb"`
	Architecture     string  `json:"architecture"`
	InstanceFamily   string  `json:"instance_family"`
	GenerationNumber int     `json:"generation_number"`
	BaselinePerformance float64 `json:"baseline_performance"`
}

// ResourceRequirement represents computed resource requirements
type ResourceRequirement struct {
	MinCPUs          int     `json:"min_cpus"`
	RecommendedCPUs  int     `json:"recommended_cpus"`
	MinMemoryGB      float64 `json:"min_memory_gb"`
	RecommendedMemoryGB float64 `json:"recommended_memory_gb"`
	MinStorageGB     int     `json:"min_storage_gb"`
	RecommendedStorageGB int `json:"recommended_storage_gb"`
	GPURequired      bool    `json:"gpu_required"`
	MinGPUs          int     `json:"min_gpus"`
	NetworkRequirements NetworkRequirements `json:"network_requirements"`
	IORequirements   IORequirements      `json:"io_requirements"`
	Reasoning        []string            `json:"reasoning"`
}

// NetworkRequirements specifies network performance needs
type NetworkRequirements struct {
	MinBandwidthMbps    int    `json:"min_bandwidth_mbps"`
	LatencyRequirement  string `json:"latency_requirement"`
	PlacementGroupNeeded bool  `json:"placement_group_needed"`
	EFARequired         bool   `json:"efa_required"`
	MultiAZRequired     bool   `json:"multi_az_required"`
}

// IORequirements specifies storage I/O needs
type IORequirements struct {
	MinIOPS           int    `json:"min_iops"`
	ThroughputMBps    int    `json:"throughput_mbps"`
	LatencyRequirement string `json:"latency_requirement"`
	RandomAccessNeeded bool   `json:"random_access_needed"`
	SequentialPattern  bool   `json:"sequential_pattern"`
}

// NewResourceAnalyzer creates a new resource analyzer
func NewResourceAnalyzer() *ResourceAnalyzer {
	ra := &ResourceAnalyzer{
		instanceSpecs: make(map[string]InstanceSpec),
	}
	
	ra.initializeInstanceSpecs()
	return ra
}

// initializeInstanceSpecs populates the instance specifications database
func (ra *ResourceAnalyzer) initializeInstanceSpecs() {
	// General purpose instances
	ra.instanceSpecs["t3.micro"] = InstanceSpec{
		vCPUs: 2, MemoryGB: 1, NetworkSpeed: "low", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "t3", GenerationNumber: 3,
		BaselinePerformance: 0.2,
	}
	ra.instanceSpecs["t3.small"] = InstanceSpec{
		vCPUs: 2, MemoryGB: 2, NetworkSpeed: "low", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "t3", GenerationNumber: 3,
		BaselinePerformance: 0.4,
	}
	ra.instanceSpecs["t3.medium"] = InstanceSpec{
		vCPUs: 2, MemoryGB: 4, NetworkSpeed: "low", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "t3", GenerationNumber: 3,
		BaselinePerformance: 0.4,
	}
	
	// Compute optimized instances
	ra.instanceSpecs["c6i.large"] = InstanceSpec{
		vCPUs: 2, MemoryGB: 4, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6i.xlarge"] = InstanceSpec{
		vCPUs: 4, MemoryGB: 8, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6i.2xlarge"] = InstanceSpec{
		vCPUs: 8, MemoryGB: 16, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6i.4xlarge"] = InstanceSpec{
		vCPUs: 16, MemoryGB: 32, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6i.8xlarge"] = InstanceSpec{
		vCPUs: 32, MemoryGB: 64, NetworkSpeed: "10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6i.12xlarge"] = InstanceSpec{
		vCPUs: 48, MemoryGB: 96, NetworkSpeed: "12g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6i.16xlarge"] = InstanceSpec{
		vCPUs: 64, MemoryGB: 128, NetworkSpeed: "25g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6i.24xlarge"] = InstanceSpec{
		vCPUs: 96, MemoryGB: 192, NetworkSpeed: "25g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	
	// Memory optimized instances
	ra.instanceSpecs["r6i.large"] = InstanceSpec{
		vCPUs: 2, MemoryGB: 16, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "r6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["r6i.xlarge"] = InstanceSpec{
		vCPUs: 4, MemoryGB: 32, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "r6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["r6i.2xlarge"] = InstanceSpec{
		vCPUs: 8, MemoryGB: 64, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "r6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["r6i.4xlarge"] = InstanceSpec{
		vCPUs: 16, MemoryGB: 128, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "r6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["r6i.8xlarge"] = InstanceSpec{
		vCPUs: 32, MemoryGB: 256, NetworkSpeed: "10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "r6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["r6i.16xlarge"] = InstanceSpec{
		vCPUs: 64, MemoryGB: 512, NetworkSpeed: "25g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "r6i", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	
	// GPU instances
	ra.instanceSpecs["g5.xlarge"] = InstanceSpec{
		vCPUs: 4, MemoryGB: 16, NetworkSpeed: "up_to_10g", StorageType: "nvme_ssd",
		GPUs: 1, GPUMemoryGB: 24, Architecture: "x86_64", InstanceFamily: "g5", GenerationNumber: 5,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["g5.2xlarge"] = InstanceSpec{
		vCPUs: 8, MemoryGB: 32, NetworkSpeed: "up_to_10g", StorageType: "nvme_ssd",
		GPUs: 1, GPUMemoryGB: 24, Architecture: "x86_64", InstanceFamily: "g5", GenerationNumber: 5,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["g5.4xlarge"] = InstanceSpec{
		vCPUs: 16, MemoryGB: 64, NetworkSpeed: "up_to_10g", StorageType: "nvme_ssd",
		GPUs: 1, GPUMemoryGB: 24, Architecture: "x86_64", InstanceFamily: "g5", GenerationNumber: 5,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["p4d.24xlarge"] = InstanceSpec{
		vCPUs: 96, MemoryGB: 1152, NetworkSpeed: "400g", StorageType: "nvme_ssd",
		GPUs: 8, GPUMemoryGB: 320, Architecture: "x86_64", InstanceFamily: "p4d", GenerationNumber: 4,
		BaselinePerformance: 1.0,
	}
	
	// AMD instances
	ra.instanceSpecs["c6a.large"] = InstanceSpec{
		vCPUs: 2, MemoryGB: 4, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6a", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6a.xlarge"] = InstanceSpec{
		vCPUs: 4, MemoryGB: 8, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6a", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6a.2xlarge"] = InstanceSpec{
		vCPUs: 8, MemoryGB: 16, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6a", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6a.4xlarge"] = InstanceSpec{
		vCPUs: 16, MemoryGB: 32, NetworkSpeed: "up_to_10g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6a", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
	ra.instanceSpecs["c6a.48xlarge"] = InstanceSpec{
		vCPUs: 192, MemoryGB: 384, NetworkSpeed: "50g", StorageType: "ebs",
		Architecture: "x86_64", InstanceFamily: "c6a", GenerationNumber: 6,
		BaselinePerformance: 1.0,
	}
}

// AnalyzeResourceRequirements determines optimal resource requirements for a workload
func (ra *ResourceAnalyzer) AnalyzeResourceRequirements(
	domain string,
	dataPattern *data.DataPattern,
	profile *data.ResearchDomainProfile,
	hints DomainHints,
) *ResourceRequirement {
	
	req := &ResourceRequirement{
		Reasoning: make([]string, 0),
	}
	
	// Analyze based on domain-specific requirements
	ra.analyzeDomainSpecificRequirements(domain, req)
	
	// Analyze based on data characteristics
	ra.analyzeDataBasedRequirements(dataPattern, req)
	
	// Apply domain profile optimizations
	if profile != nil {
		ra.applyDomainProfileOptimizations(profile, req)
	}
	
	// Apply hints and constraints
	ra.applyHintsAndConstraints(hints, req)
	
	// Generate network and I/O requirements
	ra.generateNetworkRequirements(domain, dataPattern, req)
	ra.generateIORequirements(domain, dataPattern, req)
	
	return req
}

// analyzeDomainSpecificRequirements sets domain-specific base requirements
func (ra *ResourceAnalyzer) analyzeDomainSpecificRequirements(domain string, req *ResourceRequirement) {
	switch domain {
	case "genomics":
		// Genomics workloads are typically memory-intensive
		req.MinCPUs = 8
		req.RecommendedCPUs = 16
		req.MinMemoryGB = 32
		req.RecommendedMemoryGB = 64
		req.MinStorageGB = 500
		req.RecommendedStorageGB = 1000
		req.Reasoning = append(req.Reasoning, 
			"Genomics workflows require substantial memory for sequence alignment and variant calling")
		
	case "machine_learning", "ai-research":
		// ML workloads often benefit from GPUs
		req.MinCPUs = 4
		req.RecommendedCPUs = 16
		req.MinMemoryGB = 16
		req.RecommendedMemoryGB = 64
		req.GPURequired = true
		req.MinGPUs = 1
		req.MinStorageGB = 100
		req.RecommendedStorageGB = 500
		req.Reasoning = append(req.Reasoning,
			"Machine learning workloads benefit significantly from GPU acceleration")
		
	case "climate", "climate-modeling":
		// Climate modeling requires high CPU and network performance
		req.MinCPUs = 16
		req.RecommendedCPUs = 48
		req.MinMemoryGB = 32
		req.RecommendedMemoryGB = 96
		req.MinStorageGB = 1000
		req.RecommendedStorageGB = 2000
		req.Reasoning = append(req.Reasoning,
			"Climate modeling requires high CPU count for parallel simulations")
		
	case "chemistry", "materials_science":
		// Computational chemistry is CPU and memory intensive
		req.MinCPUs = 8
		req.RecommendedCPUs = 32
		req.MinMemoryGB = 16
		req.RecommendedMemoryGB = 128
		req.MinStorageGB = 200
		req.RecommendedStorageGB = 1000
		req.Reasoning = append(req.Reasoning,
			"Computational chemistry requires high memory for molecular calculations")
		
	case "astronomy":
		// Astronomy image processing and analysis
		req.MinCPUs = 8
		req.RecommendedCPUs = 16
		req.MinMemoryGB = 32
		req.RecommendedMemoryGB = 64
		req.MinStorageGB = 1000
		req.RecommendedStorageGB = 5000
		req.Reasoning = append(req.Reasoning,
			"Astronomy data processing requires substantial storage for image files")
		
	default:
		// General research workload
		req.MinCPUs = 4
		req.RecommendedCPUs = 8
		req.MinMemoryGB = 8
		req.RecommendedMemoryGB = 32
		req.MinStorageGB = 100
		req.RecommendedStorageGB = 500
		req.Reasoning = append(req.Reasoning,
			"General research workload with moderate resource requirements")
	}
}

// analyzeDataBasedRequirements adjusts requirements based on data characteristics
func (ra *ResourceAnalyzer) analyzeDataBasedRequirements(dataPattern *data.DataPattern, req *ResourceRequirement) {
	if dataPattern == nil {
		return
	}
	
	totalSizeGB := float64(dataPattern.TotalSize) / (1024 * 1024 * 1024)
	
	// Adjust storage requirements based on data size
	if totalSizeGB > float64(req.RecommendedStorageGB) {
		req.RecommendedStorageGB = int(totalSizeGB * 1.5) // 50% buffer
		req.Reasoning = append(req.Reasoning,
			fmt.Sprintf("Increased storage to accommodate %.1f GB dataset", totalSizeGB))
	}
	
	// Adjust memory requirements for large datasets
	if totalSizeGB > 100 {
		memoryMultiplier := 1.0 + (totalSizeGB / 1000) // +1GB RAM per 1TB data
		req.RecommendedMemoryGB = req.RecommendedMemoryGB * memoryMultiplier
		req.Reasoning = append(req.Reasoning,
			"Increased memory for large dataset processing")
	}
	
	// Adjust CPU requirements based on file count
	if dataPattern.TotalFiles > 10000 {
		cpuMultiplier := 1.0 + float64(dataPattern.TotalFiles) / 100000 // +1 CPU per 100k files
		req.RecommendedCPUs = int(float64(req.RecommendedCPUs) * cpuMultiplier)
		req.Reasoning = append(req.Reasoning,
			fmt.Sprintf("Increased CPU count for processing %d files", dataPattern.TotalFiles))
	}
	
	// Check for small file penalty
	if dataPattern.FileSizes.SmallFiles.CountUnder1MB > 1000 {
		req.RecommendedCPUs = int(float64(req.RecommendedCPUs) * 1.3)
		req.Reasoning = append(req.Reasoning,
			"Increased CPU for small file handling overhead")
	}
}

// applyDomainProfileOptimizations applies domain profile specific optimizations
func (ra *ResourceAnalyzer) applyDomainProfileOptimizations(profile *data.ResearchDomainProfile, req *ResourceRequirement) {
	// Apply transfer optimization preferences
	if profile.TransferOptimization.OptimalConcurrency > req.RecommendedCPUs {
		req.RecommendedCPUs = profile.TransferOptimization.OptimalConcurrency
		req.Reasoning = append(req.Reasoning,
			fmt.Sprintf("Adjusted CPU count for optimal concurrency (%d)", 
				profile.TransferOptimization.OptimalConcurrency))
	}
	
	// Apply security requirements
	if profile.SecurityRequirements.EncryptionRequired {
		req.RecommendedCPUs = int(float64(req.RecommendedCPUs) * 1.1) // 10% overhead for encryption
		req.Reasoning = append(req.Reasoning,
			"Increased CPU for encryption overhead")
	}
	
	// Apply compliance requirements
	if len(profile.ComplianceSettings.Framework) > 0 {
		req.RecommendedMemoryGB = req.RecommendedMemoryGB * 1.2 // 20% overhead for compliance logging
		req.Reasoning = append(req.Reasoning,
			fmt.Sprintf("Increased memory for %s compliance requirements", 
				profile.ComplianceSettings.Framework))
	}
}

// applyHintsAndConstraints applies user-provided hints and constraints
func (ra *ResourceAnalyzer) applyHintsAndConstraints(hints DomainHints, req *ResourceRequirement) {
	// Apply budget constraints
	if hints.BudgetConstraint > 0 && hints.BudgetConstraint < 500 {
		// Reduce requirements for tight budgets
		req.RecommendedCPUs = int(float64(req.RecommendedCPUs) * 0.7)
		req.RecommendedMemoryGB = req.RecommendedMemoryGB * 0.7
		req.Reasoning = append(req.Reasoning,
			fmt.Sprintf("Reduced requirements due to budget constraint ($%.0f)", hints.BudgetConstraint))
	}
	
	// Apply performance hints
	for _, hint := range hints.PerformanceHints {
		switch strings.ToLower(hint) {
		case "high_performance", "performance_critical":
			req.RecommendedCPUs = int(float64(req.RecommendedCPUs) * 1.5)
			req.RecommendedMemoryGB = req.RecommendedMemoryGB * 1.5
			req.Reasoning = append(req.Reasoning,
				"Increased resources for high performance requirements")
		case "cost_optimized", "budget_friendly":
			req.RecommendedCPUs = int(float64(req.RecommendedCPUs) * 0.8)
			req.RecommendedMemoryGB = req.RecommendedMemoryGB * 0.8
			req.Reasoning = append(req.Reasoning,
				"Reduced resources for cost optimization")
		case "gpu_intensive":
			req.GPURequired = true
			req.MinGPUs = 2
			req.Reasoning = append(req.Reasoning,
				"Multiple GPUs required for GPU-intensive workloads")
		}
	}
}

// generateNetworkRequirements determines network performance needs
func (ra *ResourceAnalyzer) generateNetworkRequirements(domain string, dataPattern *data.DataPattern, req *ResourceRequirement) {
	req.NetworkRequirements = NetworkRequirements{
		MinBandwidthMbps: 1000, // 1 Gbps default
		LatencyRequirement: "standard",
		PlacementGroupNeeded: false,
		EFARequired: false,
		MultiAZRequired: false,
	}
	
	// Adjust based on domain
	switch domain {
	case "climate", "climate-modeling":
		req.NetworkRequirements.MinBandwidthMbps = 10000 // 10 Gbps for HPC
		req.NetworkRequirements.PlacementGroupNeeded = true
		req.NetworkRequirements.EFARequired = true
		req.NetworkRequirements.LatencyRequirement = "low"
		
	case "machine_learning", "ai-research":
		if req.GPURequired && req.MinGPUs > 1 {
			req.NetworkRequirements.MinBandwidthMbps = 25000 // 25 Gbps for multi-GPU
			req.NetworkRequirements.PlacementGroupNeeded = true
			req.NetworkRequirements.EFARequired = true
		}
		
	case "genomics":
		if dataPattern != nil && dataPattern.TotalSize > 1024*1024*1024*100 { // > 100GB
			req.NetworkRequirements.MinBandwidthMbps = 5000 // 5 Gbps for large datasets
		}
	}
}

// generateIORequirements determines storage I/O needs
func (ra *ResourceAnalyzer) generateIORequirements(domain string, dataPattern *data.DataPattern, req *ResourceRequirement) {
	req.IORequirements = IORequirements{
		MinIOPS: 3000,
		ThroughputMBps: 125,
		LatencyRequirement: "standard",
		RandomAccessNeeded: false,
		SequentialPattern: true,
	}
	
	// Adjust based on data patterns
	if dataPattern != nil {
		if dataPattern.FileSizes.SmallFiles.CountUnder1MB > 1000 {
			req.IORequirements.MinIOPS = 16000 // High IOPS for small files
			req.IORequirements.RandomAccessNeeded = true
		}
		
		if dataPattern.TotalSize > 1024*1024*1024*1024 { // > 1TB
			req.IORequirements.ThroughputMBps = 1000 // 1 GB/s for large datasets
		}
	}
	
	// Adjust based on domain
	switch domain {
	case "genomics":
		req.IORequirements.MinIOPS = 8000
		req.IORequirements.ThroughputMBps = 500
		req.IORequirements.RandomAccessNeeded = true // BAM file access patterns
		
	case "astronomy":
		req.IORequirements.ThroughputMBps = 2000 // High throughput for image processing
		req.IORequirements.SequentialPattern = true
		
	case "machine_learning":
		req.IORequirements.MinIOPS = 16000 // Fast data loading for training
		req.IORequirements.ThroughputMBps = 1000
	}
}

// FindOptimalInstance finds the best instance type for given requirements
func (ra *ResourceAnalyzer) FindOptimalInstance(req *ResourceRequirement) []string {
	var candidates []string
	
	for instanceType, spec := range ra.instanceSpecs {
		if ra.meetsRequirements(spec, req) {
			candidates = append(candidates, instanceType)
		}
	}
	
	// Sort by cost-effectiveness (this would use actual pricing data)
	// For now, we'll prioritize by generation and family
	return ra.sortByCostEffectiveness(candidates)
}

// meetsRequirements checks if an instance spec meets the requirements
func (ra *ResourceAnalyzer) meetsRequirements(spec InstanceSpec, req *ResourceRequirement) bool {
	if spec.vCPUs < req.MinCPUs {
		return false
	}
	
	if spec.MemoryGB < req.MinMemoryGB {
		return false
	}
	
	if req.GPURequired && spec.GPUs < req.MinGPUs {
		return false
	}
	
	return true
}

// sortByCostEffectiveness sorts instances by cost-effectiveness
func (ra *ResourceAnalyzer) sortByCostEffectiveness(instances []string) []string {
	// This would typically use actual pricing data
	// For now, we'll use a simple heuristic based on generation and family
	
	type instanceRank struct {
		name string
		score float64
	}
	
	var ranked []instanceRank
	
	for _, instance := range instances {
		spec := ra.instanceSpecs[instance]
		
		// Score based on generation (newer is better)
		score := float64(spec.GenerationNumber) * 10
		
		// Score based on baseline performance
		score += spec.BaselinePerformance * 5
		
		// Prefer compute-optimized for general workloads
		if spec.InstanceFamily == "c6i" || spec.InstanceFamily == "c6a" {
			score += 3
		}
		
		ranked = append(ranked, instanceRank{name: instance, score: score})
	}
	
	// Sort by score (descending)
	for i := 0; i < len(ranked)-1; i++ {
		for j := i + 1; j < len(ranked); j++ {
			if ranked[i].score < ranked[j].score {
				ranked[i], ranked[j] = ranked[j], ranked[i]
			}
		}
	}
	
	var result []string
	for _, r := range ranked {
		result = append(result, r.name)
	}
	
	return result
}

// GetInstanceSpec returns the specification for an instance type
func (ra *ResourceAnalyzer) GetInstanceSpec(instanceType string) (InstanceSpec, bool) {
	spec, exists := ra.instanceSpecs[instanceType]
	return spec, exists
}