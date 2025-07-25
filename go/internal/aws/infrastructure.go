package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// TerraformManager handles infrastructure operations via Terraform
type TerraformManager struct {
	client        *Client
	terraformPath string
	workingDir    string
}

// NewTerraformManager creates a new Terraform-based infrastructure manager
func NewTerraformManager(client *Client, terraformPath string) *TerraformManager {
	if terraformPath == "" {
		terraformPath = "terraform" // Default to system PATH
	}

	return &TerraformManager{
		client:        client,
		terraformPath: terraformPath,
		workingDir:    "terraform/environments/aws",
	}
}

// DeploymentStatus represents the status of a Terraform deployment
type DeploymentStatus string

const (
	DeploymentStatusPlanning   DeploymentStatus = "PLANNING"
	DeploymentStatusApplying   DeploymentStatus = "APPLYING"
	DeploymentStatusComplete   DeploymentStatus = "COMPLETE"
	DeploymentStatusFailed     DeploymentStatus = "FAILED"
	DeploymentStatusDestroying DeploymentStatus = "DESTROYING"
	DeploymentStatusDestroyed  DeploymentStatus = "DESTROYED"
)

// DeploymentInfo contains information about a Terraform deployment
type DeploymentInfo struct {
	WorkspaceName string
	Status        DeploymentStatus
	CreatedTime   time.Time
	UpdatedTime   *time.Time
	Outputs       map[string]string
	Resources     []TerraformResource
}

// TerraformResource represents a resource managed by Terraform
type TerraformResource struct {
	Address      string
	Type         string
	Name         string
	ProviderName string
	Values       map[string]interface{}
}

// SetWorkingDirectory sets the Terraform working directory
func (tm *TerraformManager) SetWorkingDirectory(dir string) {
	tm.workingDir = dir
}

// InitTerraform initializes the Terraform working directory
func (tm *TerraformManager) InitTerraform(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, tm.terraformPath, "init")
	cmd.Dir = tm.workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform init failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// PlanDeployment creates a Terraform plan
func (tm *TerraformManager) PlanDeployment(ctx context.Context, variables map[string]string) error {
	args := []string{"plan"}

	// Add variables
	for key, value := range variables {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
	}

	cmd := exec.CommandContext(ctx, tm.terraformPath, args...)
	cmd.Dir = tm.workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform plan failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// ApplyDeployment applies the Terraform configuration
func (tm *TerraformManager) ApplyDeployment(ctx context.Context, variables map[string]string) (*DeploymentInfo, error) {
	args := []string{"apply", "-auto-approve"}

	// Add variables
	for key, value := range variables {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
	}

	cmd := exec.CommandContext(ctx, tm.terraformPath, args...)
	cmd.Dir = tm.workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("terraform apply failed: %w\nOutput: %s", err, output)
	}

	// Get outputs
	outputs, err := tm.GetOutputs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get outputs after apply: %w", err)
	}

	return &DeploymentInfo{
		WorkspaceName: "default",
		Status:        DeploymentStatusComplete,
		CreatedTime:   time.Now(),
		Outputs:       outputs,
	}, nil
}

// TerraformOutput represents the structure of Terraform output
type TerraformOutput struct {
	Value     interface{} `json:"value"`
	Type      interface{} `json:"type"`
	Sensitive bool        `json:"sensitive"`
}

// GetOutputs retrieves Terraform outputs
func (tm *TerraformManager) GetOutputs(ctx context.Context) (map[string]string, error) {
	cmd := exec.CommandContext(ctx, tm.terraformPath, "output", "-json")
	cmd.Dir = tm.workingDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("terraform output failed: %w", err)
	}

	// Parse JSON output
	var terraformOutputs map[string]TerraformOutput
	if err := json.Unmarshal(output, &terraformOutputs); err != nil {
		return nil, fmt.Errorf("failed to parse terraform outputs: %w", err)
	}

	// Convert to string map
	outputs := make(map[string]string)
	for key, value := range terraformOutputs {
		// Convert value to string representation
		if str, ok := value.Value.(string); ok {
			outputs[key] = str
		} else {
			// For non-string values, marshal back to JSON
			jsonBytes, err := json.Marshal(value.Value)
			if err != nil {
				outputs[key] = fmt.Sprintf("%v", value.Value)
			} else {
				outputs[key] = string(jsonBytes)
			}
		}
	}

	return outputs, nil
}

// DestroyDeployment destroys the Terraform-managed infrastructure
func (tm *TerraformManager) DestroyDeployment(ctx context.Context, variables map[string]string) error {
	args := []string{"destroy", "-auto-approve"}

	// Add variables
	for key, value := range variables {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
	}

	cmd := exec.CommandContext(ctx, tm.terraformPath, args...)
	cmd.Dir = tm.workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform destroy failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// GetDeploymentInfo retrieves information about the current deployment
func (tm *TerraformManager) GetDeploymentInfo(ctx context.Context) (*DeploymentInfo, error) {
	// Check if state file exists
	stateFile := filepath.Join(tm.workingDir, "terraform.tfstate")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("no deployment found - state file does not exist")
	}

	// Get outputs
	outputs, err := tm.GetOutputs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get outputs: %w", err)
	}

	// Get state info
	showCmd := exec.CommandContext(ctx, tm.terraformPath, "show", "-json")
	showCmd.Dir = tm.workingDir

	stateOutput, err := showCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("terraform show failed: %w", err)
	}

	// Parse state
	resources, createdTime, err := tm.parseStateJSON(stateOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to parse terraform state: %w", err)
	}

	return &DeploymentInfo{
		WorkspaceName: "default",
		Status:        DeploymentStatusComplete,
		CreatedTime:   createdTime,
		Outputs:       outputs,
		Resources:     resources,
	}, nil
}

// TerraformState represents the structure of Terraform state JSON
type TerraformState struct {
	FormatVersion string `json:"format_version"`
	Values        struct {
		RootModule struct {
			Resources []struct {
				Address      string                 `json:"address"`
				Type         string                 `json:"type"`
				Name         string                 `json:"name"`
				ProviderName string                 `json:"provider_name"`
				Values       map[string]interface{} `json:"values"`
			} `json:"resources"`
		} `json:"root_module"`
	} `json:"values"`
}

// parseStateJSON parses Terraform state JSON and extracts resources and creation time
func (tm *TerraformManager) parseStateJSON(stateOutput []byte) ([]TerraformResource, time.Time, error) {
	var state TerraformState
	if err := json.Unmarshal(stateOutput, &state); err != nil {
		return nil, time.Now(), fmt.Errorf("failed to unmarshal terraform state: %w", err)
	}

	var resources []TerraformResource
	var createdTime time.Time = time.Now() // Default to now if not found

	for _, resource := range state.Values.RootModule.Resources {
		tfResource := TerraformResource{
			Address:      resource.Address,
			Type:         resource.Type,
			Name:         resource.Name,
			ProviderName: resource.ProviderName,
			Values:       resource.Values,
		}
		resources = append(resources, tfResource)

		// Try to extract creation time from resource values
		if createdTimeStr, ok := resource.Values["time_created"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, createdTimeStr); err == nil {
				if createdTime.IsZero() || parsed.Before(createdTime) {
					createdTime = parsed
				}
			}
		}
	}

	return resources, createdTime, nil
}

// ValidateTerraformConfiguration validates the Terraform configuration
func (tm *TerraformManager) ValidateTerraformConfiguration(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, tm.terraformPath, "validate")
	cmd.Dir = tm.workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform validate failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// InstanceInfo contains EC2 instance information (kept for compatibility)
type InstanceInfo struct {
	InstanceID       string
	InstanceType     string
	State            string
	PublicIP         string
	PrivateIP        string
	AvailabilityZone string
	LaunchTime       time.Time
	Tags             map[string]string
}

// ListInstances lists EC2 instances with optional filtering
func (tm *TerraformManager) ListInstances(ctx context.Context, filters map[string][]string) ([]InstanceInfo, error) {
	// Convert filters to EC2 format
	ec2Filters := make([]ec2types.Filter, 0, len(filters))
	for name, values := range filters {
		ec2Filters = append(ec2Filters, ec2types.Filter{
			Name:   &name,
			Values: values,
		})
	}

	input := &ec2.DescribeInstancesInput{
		Filters: ec2Filters,
	}

	result, err := tm.client.EC2.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	var instances []InstanceInfo
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// Extract tags
			tags := make(map[string]string)
			for _, tag := range instance.Tags {
				if tag.Key != nil && tag.Value != nil {
					tags[*tag.Key] = *tag.Value
				}
			}

			instanceInfo := InstanceInfo{
				InstanceID:       *instance.InstanceId,
				InstanceType:     string(instance.InstanceType),
				State:            string(instance.State.Name),
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				LaunchTime:       *instance.LaunchTime,
				Tags:             tags,
			}

			if instance.PublicIpAddress != nil {
				instanceInfo.PublicIP = *instance.PublicIpAddress
			}

			if instance.PrivateIpAddress != nil {
				instanceInfo.PrivateIP = *instance.PrivateIpAddress
			}

			instances = append(instances, instanceInfo)
		}
	}

	return instances, nil
}

// TerminateInstance terminates an EC2 instance
func (tm *TerraformManager) TerminateInstance(ctx context.Context, instanceID string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	_, err := tm.client.EC2.TerminateInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to terminate instance %s: %w", instanceID, err)
	}

	return nil
}

// GetDefaultVPC gets the default VPC for the region
func (tm *TerraformManager) GetDefaultVPC(ctx context.Context) (string, error) {
	input := &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			{
				Name:   func() *string { s := "is-default"; return &s }(),
				Values: []string{"true"},
			},
		},
	}

	result, err := tm.client.EC2.DescribeVpcs(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to describe VPCs: %w", err)
	}

	if len(result.Vpcs) == 0 {
		return "", fmt.Errorf("no default VPC found")
	}

	return *result.Vpcs[0].VpcId, nil
}
