package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// NewDockerEnvironment creates a new Docker environment for running tests
func NewDockerEnvironment(config Config) Environment {
	image := "ubuntu:22.04"
	workDir := "/workspace"

	return &DockerEnvironment{
		Image:     image,
		WorkDir:   workDir,
		EnvVars:   config.EnvVars,
		Resources: []Resource{},
	}
}

// Setup prepares the Docker environment for test execution
func (e *DockerEnvironment) Setup(ctx context.Context) error {
	// TODO: Implement Docker container setup
	// This is a placeholder implementation

	e.Resources = append(e.Resources, Resource{
		Type:       "docker-container",
		Identifier: "tutorial-guard-container", // Would be actual container ID
		CreatedAt:  time.Now(),
		Metadata: map[string]string{
			"image":    e.Image,
			"work_dir": e.WorkDir,
			"purpose":  "test_execution",
		},
	})

	return nil
}

// Execute runs a code example in the Docker environment
func (e *DockerEnvironment) Execute(ctx context.Context, example extractor.Example) (*TestResult, error) {
	startTime := time.Now()

	result := &TestResult{
		ExampleID:   example.ID,
		StartTime:   startTime,
		Environment: "docker",
		Metadata:    make(map[string]string),
	}

	// TODO: Implement actual Docker execution
	// This is a placeholder that simulates Docker execution

	// For now, return a simulated successful result
	result.Success = true
	result.ExitCode = 0
	result.Output = fmt.Sprintf("Simulated Docker execution of %s code:\n%s", example.Language, example.Code)
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	result.Metadata["container_id"] = "simulated-container-id"
	result.Metadata["image"] = e.Image
	result.Metadata["work_dir"] = e.WorkDir

	return result, nil
}

// Cleanup removes Docker resources created during test execution
func (e *DockerEnvironment) Cleanup(ctx context.Context) error {
	// TODO: Implement actual Docker cleanup
	// - Stop and remove containers
	// - Remove created volumes
	// - Clean up networks if created

	return nil
}

// GetResources returns the list of Docker resources created
func (e *DockerEnvironment) GetResources() []Resource {
	return e.Resources
}

// Note: This is a minimal Docker implementation. A full implementation would:
// 1. Use the Docker API or docker CLI to create containers
// 2. Mount working directory as volume
// 3. Copy script files into container
// 4. Execute commands inside container
// 5. Capture output and exit codes
// 6. Handle container lifecycle (start, stop, remove)
// 7. Support different base images for different languages
// 8. Handle networking and security considerations
