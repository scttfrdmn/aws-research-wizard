package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// InfrastructureManager handles AWS infrastructure operations
type InfrastructureManager struct {
	client *Client
}

// NewInfrastructureManager creates a new infrastructure manager
func NewInfrastructureManager(client *Client) *InfrastructureManager {
	return &InfrastructureManager{
		client: client,
	}
}

// StackStatus represents CloudFormation stack status
type StackStatus string

const (
	StackStatusCreateInProgress StackStatus = "CREATE_IN_PROGRESS"
	StackStatusCreateComplete   StackStatus = "CREATE_COMPLETE"
	StackStatusCreateFailed     StackStatus = "CREATE_FAILED"
	StackStatusDeleteInProgress StackStatus = "DELETE_IN_PROGRESS"
	StackStatusDeleteComplete   StackStatus = "DELETE_COMPLETE"
	StackStatusDeleteFailed     StackStatus = "DELETE_FAILED"
	StackStatusUpdateInProgress StackStatus = "UPDATE_IN_PROGRESS"
	StackStatusUpdateComplete   StackStatus = "UPDATE_COMPLETE"
	StackStatusUpdateFailed     StackStatus = "UPDATE_FAILED"
)

// StackInfo contains CloudFormation stack information
type StackInfo struct {
	StackName   string
	StackID     string
	Status      StackStatus
	CreatedTime time.Time
	UpdatedTime *time.Time
	Outputs     map[string]string
	Parameters  map[string]string
}

// CreateStack creates a new CloudFormation stack
func (im *InfrastructureManager) CreateStack(ctx context.Context, stackName string, templateBody string, parameters map[string]string) (*StackInfo, error) {
	// Convert parameters to CloudFormation format
	cfParams := make([]types.Parameter, 0, len(parameters))
	for key, value := range parameters {
		cfParams = append(cfParams, types.Parameter{
			ParameterKey:   aws.String(key),
			ParameterValue: aws.String(value),
		})
	}

	input := &cloudformation.CreateStackInput{
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(templateBody),
		Parameters:   cfParams,
		Capabilities: []types.Capability{
			types.CapabilityCapabilityIam,
			types.CapabilityCapabilityNamedIam,
		},
		Tags: []types.Tag{
			{
				Key:   aws.String("CreatedBy"),
				Value: aws.String("AWS-Research-Wizard"),
			},
			{
				Key:   aws.String("Purpose"),
				Value: aws.String("Research-Infrastructure"),
			},
		},
	}

	result, err := im.client.CloudFormation.CreateStack(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create stack: %w", err)
	}

	return &StackInfo{
		StackName: stackName,
		StackID:   *result.StackId,
		Status:    StackStatusCreateInProgress,
	}, nil
}

// GetStackInfo retrieves information about a CloudFormation stack
func (im *InfrastructureManager) GetStackInfo(ctx context.Context, stackName string) (*StackInfo, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	result, err := im.client.CloudFormation.DescribeStacks(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe stack: %w", err)
	}

	if len(result.Stacks) == 0 {
		return nil, fmt.Errorf("stack not found: %s", stackName)
	}

	stack := result.Stacks[0]

	// Extract outputs
	outputs := make(map[string]string)
	for _, output := range stack.Outputs {
		if output.OutputKey != nil && output.OutputValue != nil {
			outputs[*output.OutputKey] = *output.OutputValue
		}
	}

	// Extract parameters
	parameters := make(map[string]string)
	for _, param := range stack.Parameters {
		if param.ParameterKey != nil && param.ParameterValue != nil {
			parameters[*param.ParameterKey] = *param.ParameterValue
		}
	}

	stackInfo := &StackInfo{
		StackName:   *stack.StackName,
		StackID:     *stack.StackId,
		Status:      StackStatus(stack.StackStatus),
		CreatedTime: *stack.CreationTime,
		Outputs:     outputs,
		Parameters:  parameters,
	}

	if stack.LastUpdatedTime != nil {
		stackInfo.UpdatedTime = stack.LastUpdatedTime
	}

	return stackInfo, nil
}

// DeleteStack deletes a CloudFormation stack
func (im *InfrastructureManager) DeleteStack(ctx context.Context, stackName string) error {
	input := &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	}

	_, err := im.client.CloudFormation.DeleteStack(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete stack: %w", err)
	}

	return nil
}

// WaitForStackComplete waits for a stack operation to complete
func (im *InfrastructureManager) WaitForStackComplete(ctx context.Context, stackName string, timeout time.Duration) (*StackInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for stack operation to complete")
		case <-ticker.C:
			stackInfo, err := im.GetStackInfo(ctx, stackName)
			if err != nil {
				return nil, err
			}

			switch stackInfo.Status {
			case StackStatusCreateComplete, StackStatusUpdateComplete, StackStatusDeleteComplete:
				return stackInfo, nil
			case StackStatusCreateFailed, StackStatusUpdateFailed, StackStatusDeleteFailed:
				return stackInfo, fmt.Errorf("stack operation failed with status: %s", stackInfo.Status)
			}
		}
	}
}

// InstanceInfo contains EC2 instance information
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
func (im *InfrastructureManager) ListInstances(ctx context.Context, filters map[string][]string) ([]InstanceInfo, error) {
	// Convert filters to EC2 format
	ec2Filters := make([]ec2types.Filter, 0, len(filters))
	for name, values := range filters {
		ec2Filters = append(ec2Filters, ec2types.Filter{
			Name:   aws.String(name),
			Values: values,
		})
	}

	input := &ec2.DescribeInstancesInput{
		Filters: ec2Filters,
	}

	result, err := im.client.EC2.DescribeInstances(ctx, input)
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
func (im *InfrastructureManager) TerminateInstance(ctx context.Context, instanceID string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	_, err := im.client.EC2.TerminateInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to terminate instance %s: %w", instanceID, err)
	}

	return nil
}

// CreateSecurityGroup creates a new security group
func (im *InfrastructureManager) CreateSecurityGroup(ctx context.Context, groupName, description, vpcID string) (string, error) {
	input := &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(groupName),
		Description: aws.String(description),
		VpcId:       aws.String(vpcID),
	}

	result, err := im.client.EC2.CreateSecurityGroup(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create security group: %w", err)
	}

	return *result.GroupId, nil
}

// GetDefaultVPC gets the default VPC for the region
func (im *InfrastructureManager) GetDefaultVPC(ctx context.Context) (string, error) {
	input := &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("is-default"),
				Values: []string{"true"},
			},
		},
	}

	result, err := im.client.EC2.DescribeVpcs(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to describe VPCs: %w", err)
	}

	if len(result.Vpcs) == 0 {
		return "", fmt.Errorf("no default VPC found")
	}

	return *result.Vpcs[0].VpcId, nil
}