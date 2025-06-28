package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client provides comprehensive AWS service access
type Client struct {
	cfg             aws.Config
	EC2             *ec2.Client
	CloudFormation  *cloudformation.Client
	CloudWatch      *cloudwatch.Client
	CostExplorer    *costexplorer.Client
	IAM             *iam.Client
	S3              *s3.Client
	Region          string
}

// NewClient creates a new AWS client with all required services
func NewClient(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Client{
		cfg:            cfg,
		EC2:            ec2.NewFromConfig(cfg),
		CloudFormation: cloudformation.NewFromConfig(cfg),
		CloudWatch:     cloudwatch.NewFromConfig(cfg),
		CostExplorer:   costexplorer.NewFromConfig(cfg),
		IAM:            iam.NewFromConfig(cfg),
		S3:             s3.NewFromConfig(cfg),
		Region:         region,
	}, nil
}

// ValidateCredentials checks if AWS credentials are properly configured
func (c *Client) ValidateCredentials(ctx context.Context) error {
	_, err := c.EC2.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return fmt.Errorf("failed to validate AWS credentials: %w", err)
	}
	return nil
}

// GetAccountID retrieves the current AWS account ID
func (c *Client) GetAccountID(ctx context.Context) (string, error) {
	result, err := c.IAM.GetUser(ctx, &iam.GetUserInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get account ID: %w", err)
	}
	
	// Extract account ID from ARN
	arn := *result.User.Arn
	// ARN format: arn:aws:iam::ACCOUNT-ID:user/username
	parts := make([]string, 0, 6)
	start := 0
	for i, char := range arn {
		if char == ':' {
			parts = append(parts, arn[start:i])
			start = i + 1
		}
	}
	if start < len(arn) {
		parts = append(parts, arn[start:])
	}
	
	if len(parts) >= 5 {
		return parts[4], nil
	}
	
	return "", fmt.Errorf("failed to parse account ID from ARN: %s", arn)
}

// GetAvailabilityZones retrieves available AZs in the current region
func (c *Client) GetAvailabilityZones(ctx context.Context) ([]string, error) {
	result, err := c.EC2.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe availability zones: %w", err)
	}
	
	zones := make([]string, 0, len(result.AvailabilityZones))
	for _, zone := range result.AvailabilityZones {
		zones = append(zones, *zone.ZoneName)
	}
	
	return zones, nil
}