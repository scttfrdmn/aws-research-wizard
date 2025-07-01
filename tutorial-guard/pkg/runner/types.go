package runner

import (
	"context"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// TestResult represents the result of executing a code example
type TestResult struct {
	ExampleID   string            `json:"example_id"`
	Success     bool              `json:"success"`
	ExitCode    int               `json:"exit_code"`
	Output      string            `json:"output"`
	ErrorOutput string            `json:"error_output"`
	Duration    time.Duration     `json:"duration"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Environment string            `json:"environment"`
	Resources   []Resource        `json:"resources"` // Created resources
	Errors      []string          `json:"errors"`    // Error details
	Warnings    []string          `json:"warnings"`  // Warning messages
	Metadata    map[string]string `json:"metadata"`  // Additional info
}

// TestSuite represents a collection of test results
type TestSuite struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Results     []TestResult      `json:"results"`
	Summary     TestSummary       `json:"summary"`
	Environment string            `json:"environment"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Duration    time.Duration     `json:"duration"`
	Metadata    map[string]string `json:"metadata"`
}

// TestSummary provides aggregate statistics
type TestSummary struct {
	Total    int `json:"total"`
	Passed   int `json:"passed"`
	Failed   int `json:"failed"`
	Skipped  int `json:"skipped"`
	Warnings int `json:"warnings"`
}

// Resource represents a resource created during test execution
type Resource struct {
	Type       string            `json:"type"`       // file, directory, container, aws-resource, etc.
	Identifier string            `json:"identifier"` // path, id, arn, etc.
	Metadata   map[string]string `json:"metadata"`   // additional details
	CreatedAt  time.Time         `json:"created_at"`
}

// Config holds configuration for the test runner
type Config struct {
	Environment string   // docker, local, aws, etc.
	Timeout     string   // timeout per test
	Parallel    int      // number of parallel runners
	Cleanup     bool     // cleanup resources after tests
	Verbose     bool     // verbose output
	WorkDir     string   // working directory for tests
	EnvVars     []string // environment variables
	PreScript   string   // script to run before tests
	PostScript  string   // script to run after tests
	RetryCount  int      // number of retries for failed tests
	FailFast    bool     // stop on first failure
}

// Runner interface defines the contract for executing examples
type Runner interface {
	RunExample(ctx context.Context, example extractor.Example) (*TestResult, error)
	RunExamples(ctx context.Context, examples []extractor.Example) (*TestSuite, error)
	RunExamplesFromFile(filename string) (*TestSuite, error)
	Cleanup() error
}

// Environment interface defines the contract for test environments
type Environment interface {
	Setup(ctx context.Context) error
	Execute(ctx context.Context, example extractor.Example) (*TestResult, error)
	Cleanup(ctx context.Context) error
	GetResources() []Resource
}

// DockerEnvironment runs tests in Docker containers
type DockerEnvironment struct {
	Image       string
	WorkDir     string
	Volumes     []string
	EnvVars     []string
	Network     string
	Resources   []Resource
	ContainerID string
}

// LocalEnvironment runs tests in the local shell
type LocalEnvironment struct {
	WorkDir   string
	EnvVars   []string
	Shell     string
	Resources []Resource
}

// AWSEnvironment runs tests against AWS resources
type AWSEnvironment struct {
	Region     string
	Profile    string
	AssumeRole string
	WorkDir    string
	Resources  []Resource
	S3Bucket   string
	VPC        string
	Subnet     string
}

// TestRunner is the main implementation of the Runner interface
type TestRunner struct {
	config      Config
	environment Environment
}

// ExecutionContext holds context for a single test execution
type ExecutionContext struct {
	WorkDir   string
	EnvVars   map[string]string
	Timeout   time.Duration
	Example   extractor.Example
	Resources []Resource
}

// ValidationRule defines a rule for validating test results
type ValidationRule struct {
	Name        string
	Description string
	Check       func(*TestResult) error
}

// CleanupHandler manages resource cleanup
type CleanupHandler struct {
	Resources []Resource
	Handlers  map[string]func(Resource) error
}
