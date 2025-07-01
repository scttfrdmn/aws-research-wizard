package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// NewAWSEnvironment creates a new AWS environment for running tests
func NewAWSEnvironment(config Config) Environment {
	return &AWSEnvironment{
		Region:    "us-east-1",
		WorkDir:   "/tmp/tutorial-guard-aws",
		Resources: []Resource{},
	}
}

// Setup prepares the AWS environment for test execution
func (e *AWSEnvironment) Setup(ctx context.Context) error {
	// TODO: Implement AWS environment setup
	// This could include:
	// - Creating temporary S3 bucket for test artifacts
	// - Setting up EC2 instances for execution
	// - Creating IAM roles for testing
	// - Setting up VPC/networking if needed

	e.Resources = append(e.Resources, Resource{
		Type:       "aws-s3-bucket",
		Identifier: "tutorial-guard-test-bucket", // Would be actual bucket name
		CreatedAt:  time.Now(),
		Metadata: map[string]string{
			"region":  e.Region,
			"purpose": "test_artifacts",
		},
	})

	return nil
}

// Execute runs a code example in the AWS environment
func (e *AWSEnvironment) Execute(ctx context.Context, example extractor.Example) (*TestResult, error) {
	startTime := time.Now()

	result := &TestResult{
		ExampleID:   example.ID,
		StartTime:   startTime,
		Environment: "aws",
		Metadata:    make(map[string]string),
	}

	// TODO: Implement actual AWS execution
	// This could involve:
	// - Uploading scripts to S3
	// - Running commands on EC2 instances
	// - Using AWS Systems Manager for command execution
	// - Validating AWS CLI commands
	// - Testing AWS resource creation/deletion

	// For now, simulate AWS command validation
	if example.Language == "bash" && containsAWSCommands(example.Code) {
		result.Success = true
		result.ExitCode = 0
		result.Output = fmt.Sprintf("Simulated AWS execution of:\n%s", example.Code)
		result.Metadata["aws_commands_detected"] = "true"
	} else {
		result.Success = false
		result.Errors = []string{"AWS environment only supports AWS CLI commands"}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	result.Metadata["region"] = e.Region
	result.Metadata["s3_bucket"] = e.S3Bucket

	return result, nil
}

// Cleanup removes AWS resources created during test execution
func (e *AWSEnvironment) Cleanup(ctx context.Context) error {
	// TODO: Implement actual AWS cleanup
	// This is critical for cost control:
	// - Delete S3 buckets and objects
	// - Terminate EC2 instances
	// - Remove IAM roles/policies
	// - Delete VPC resources
	// - Clean up any other billable resources

	return nil
}

// GetResources returns the list of AWS resources created
func (e *AWSEnvironment) GetResources() []Resource {
	return e.Resources
}

// containsAWSCommands checks if the code contains AWS CLI commands
func containsAWSCommands(code string) bool {
	awsCommands := []string{
		"aws s3",
		"aws ec2",
		"aws iam",
		"aws cloudformation",
		"aws lambda",
		"aws configure",
	}

	for _, cmd := range awsCommands {
		if containsIgnoreCase(code, cmd) {
			return true
		}
	}

	return false
}

// Helper function for case-insensitive string checking
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				anySubstring(s, substr)))
}

func anySubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if equalIgnoreCase(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

func equalIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if toLower(a[i]) != toLower(b[i]) {
			return false
		}
	}
	return true
}

func toLower(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// Note: This is a minimal AWS implementation. A full implementation would:
// 1. Use AWS SDK to interact with services
// 2. Create isolated test environments (separate VPC, etc.)
// 3. Implement cost controls and resource limits
// 4. Support multiple AWS regions
// 5. Handle IAM permissions and assume roles
// 6. Validate actual AWS resource creation/deletion
// 7. Integrate with AWS CloudFormation for complex setups
// 8. Implement comprehensive cleanup to avoid charges
