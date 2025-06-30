package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// PricingCalculator handles AWS cost calculations
type PricingCalculator struct {
	ec2Client *ec2.Client
	region    string
}

// CostEstimate represents cost breakdown for an instance
type CostEstimate struct {
	InstanceType    string
	VCPUs           int32
	Memory          string
	HourlyCost      float64
	MonthlyCost     float64
	AnnualCost      float64
	SpotSavings     float64
	ReservedSavings float64
}

// NewPricingCalculator creates a new pricing calculator
func NewPricingCalculator(region string) (*PricingCalculator, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &PricingCalculator{
		ec2Client: ec2.NewFromConfig(cfg),
		region:    region,
	}, nil
}

// GetInstanceTypes retrieves available instance types
func (pc *PricingCalculator) GetInstanceTypes(ctx context.Context) ([]types.InstanceTypeInfo, error) {
	input := &ec2.DescribeInstanceTypesInput{}

	result, err := pc.ec2Client.DescribeInstanceTypes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance types: %w", err)
	}

	return result.InstanceTypes, nil
}

// CalculateCost estimates costs for a given instance type
func (pc *PricingCalculator) CalculateCost(instanceType string) (*CostEstimate, error) {
	// This is a simplified cost calculation using approximate pricing
	// In production, you'd use the AWS Pricing API or maintain a pricing table

	baseCosts := map[string]float64{
		// General Purpose
		"t3.micro":   0.0104,
		"t3.small":   0.0208,
		"t3.medium":  0.0416,
		"t3.large":   0.0832,
		"t3.xlarge":  0.1664,
		"t3.2xlarge": 0.3328,

		// Compute Optimized
		"c6i.large":    0.085,
		"c6i.xlarge":   0.17,
		"c6i.2xlarge":  0.34,
		"c6i.4xlarge":  0.68,
		"c6i.8xlarge":  1.36,
		"c6i.12xlarge": 2.04,
		"c6i.16xlarge": 2.72,
		"c6i.24xlarge": 4.08,

		// Memory Optimized
		"r6i.large":    0.126,
		"r6i.xlarge":   0.252,
		"r6i.2xlarge":  0.504,
		"r6i.4xlarge":  1.008,
		"r6i.8xlarge":  2.016,
		"r6i.12xlarge": 3.024,
		"r6i.16xlarge": 4.032,
		"r6i.24xlarge": 6.048,

		// GPU Instances
		"p4d.24xlarge": 32.7726,
		"p3.2xlarge":   3.06,
		"p3.8xlarge":   12.24,
		"p3.16xlarge":  24.48,

		// High Performance Computing
		"hpc6a.48xlarge":  2.88,
		"hpc6id.32xlarge": 3.456,
	}

	hourlyCost, exists := baseCosts[instanceType]
	if !exists {
		// Fallback estimation based on instance size
		hourlyCost = pc.estimateCostFromInstanceType(instanceType)
	}

	// Calculate monthly and annual costs
	monthlyCost := hourlyCost * 24 * 30.44 // Average days per month
	annualCost := hourlyCost * 24 * 365

	// Estimate savings
	spotSavings := hourlyCost * 0.7     // ~70% savings
	reservedSavings := hourlyCost * 0.4 // ~40% savings

	// Extract vCPUs and memory from instance type (simplified)
	vcpus, memory := pc.parseInstanceSpecs(instanceType)

	return &CostEstimate{
		InstanceType:    instanceType,
		VCPUs:           vcpus,
		Memory:          memory,
		HourlyCost:      hourlyCost,
		MonthlyCost:     monthlyCost,
		AnnualCost:      annualCost,
		SpotSavings:     spotSavings,
		ReservedSavings: reservedSavings,
	}, nil
}

// estimateCostFromInstanceType provides fallback cost estimation
func (pc *PricingCalculator) estimateCostFromInstanceType(instanceType string) float64 {
	// Simple heuristic based on instance type patterns
	parts := strings.Split(instanceType, ".")
	if len(parts) != 2 {
		return 0.10 // Default fallback
	}

	family := parts[0]
	size := parts[1]

	// Base cost multipliers by family
	familyMultipliers := map[string]float64{
		"t3":    0.02,
		"t4g":   0.018,
		"m6i":   0.08,
		"c6i":   0.085,
		"r6i":   0.126,
		"p4d":   15.0,
		"p3":    1.5,
		"hpc6a": 0.06,
	}

	// Size multipliers
	sizeMultipliers := map[string]float64{
		"nano":     0.25,
		"micro":    0.5,
		"small":    1.0,
		"medium":   2.0,
		"large":    4.0,
		"xlarge":   8.0,
		"2xlarge":  16.0,
		"4xlarge":  32.0,
		"8xlarge":  64.0,
		"12xlarge": 96.0,
		"16xlarge": 128.0,
		"24xlarge": 192.0,
		"32xlarge": 256.0,
		"48xlarge": 384.0,
	}

	familyBase := familyMultipliers[family]
	if familyBase == 0 {
		familyBase = 0.08 // Default
	}

	sizeMultiplier := sizeMultipliers[size]
	if sizeMultiplier == 0 {
		sizeMultiplier = 4.0 // Default to large
	}

	return familyBase * sizeMultiplier
}

// parseInstanceSpecs extracts vCPUs and memory info (simplified)
func (pc *PricingCalculator) parseInstanceSpecs(instanceType string) (int32, string) {
	// This is a simplified mapping - in production you'd query the EC2 API
	specs := map[string]struct {
		vcpus  int32
		memory string
	}{
		"t3.micro":       {2, "1 GiB"},
		"t3.small":       {2, "2 GiB"},
		"t3.medium":      {2, "4 GiB"},
		"t3.large":       {2, "8 GiB"},
		"t3.xlarge":      {4, "16 GiB"},
		"t3.2xlarge":     {8, "32 GiB"},
		"c6i.large":      {2, "4 GiB"},
		"c6i.xlarge":     {4, "8 GiB"},
		"c6i.2xlarge":    {8, "16 GiB"},
		"c6i.4xlarge":    {16, "32 GiB"},
		"c6i.8xlarge":    {32, "64 GiB"},
		"c6i.12xlarge":   {48, "96 GiB"},
		"c6i.16xlarge":   {64, "128 GiB"},
		"c6i.24xlarge":   {96, "192 GiB"},
		"r6i.large":      {2, "16 GiB"},
		"r6i.xlarge":     {4, "32 GiB"},
		"r6i.2xlarge":    {8, "64 GiB"},
		"r6i.4xlarge":    {16, "128 GiB"},
		"r6i.8xlarge":    {32, "256 GiB"},
		"r6i.12xlarge":   {48, "384 GiB"},
		"r6i.16xlarge":   {64, "512 GiB"},
		"r6i.24xlarge":   {96, "768 GiB"},
		"p4d.24xlarge":   {96, "1152 GiB"},
		"p3.2xlarge":     {8, "61 GiB"},
		"p3.8xlarge":     {32, "244 GiB"},
		"p3.16xlarge":    {64, "488 GiB"},
		"hpc6a.48xlarge": {96, "384 GiB"},
	}

	if spec, exists := specs[instanceType]; exists {
		return spec.vcpus, spec.memory
	}

	// Fallback parsing from instance type name
	parts := strings.Split(instanceType, ".")
	if len(parts) == 2 {
		size := parts[1]
		// Simple heuristic
		switch size {
		case "micro":
			return 1, "1 GiB"
		case "small":
			return 1, "2 GiB"
		case "medium":
			return 2, "4 GiB"
		case "large":
			return 2, "8 GiB"
		case "xlarge":
			return 4, "16 GiB"
		case "2xlarge":
			return 8, "32 GiB"
		case "4xlarge":
			return 16, "64 GiB"
		case "8xlarge":
			return 32, "128 GiB"
		default:
			// Try to parse number from size
			if strings.Contains(size, "xlarge") {
				multiplierStr := strings.Replace(size, "xlarge", "", 1)
				if multiplier, err := strconv.Atoi(multiplierStr); err == nil {
					// Check for overflow before conversion
					result := multiplier * 4
					if result > 2147483647 {
						return 2147483647, fmt.Sprintf("%d GiB", multiplier*16)
					}
					return int32(result), fmt.Sprintf("%d GiB", multiplier*16)
				}
			}
		}
	}

	return 2, "4 GiB" // Default fallback
}
