/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// SimpleLocalRunner executes commands locally on the host system
type SimpleLocalRunner struct {
	config     LocalConfig
	workingDir string
}

// NewLocalRunner creates a new local command runner
func NewLocalRunner(config LocalConfig) Runner {
	if config.Shell == "" {
		config.Shell = "/bin/sh"
	}
	if config.WorkingDirectory == "" {
		config.WorkingDirectory = "/tmp"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &SimpleLocalRunner{
		config:     config,
		workingDir: config.WorkingDirectory,
	}
}

// Execute runs a command locally and returns the result
func (r *SimpleLocalRunner) Execute(ctx context.Context, command string) (*RunResult, error) {
	start := time.Now()
	
	result := &RunResult{
		Command:     command,
		StartTime:   start,
		Environment: "local",
		WorkingDir:  r.workingDir,
		Success:     false,
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, r.config.Timeout)
	defer cancel()

	// Prepare command
	cmd := exec.CommandContext(timeoutCtx, r.config.Shell, "-c", command)
	cmd.Dir = r.workingDir
	
	// Set environment variables
	if r.config.Environment != nil {
		env := os.Environ()
		for key, value := range r.config.Environment {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	// Execute command
	stdout, err := cmd.Output()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Stdout = string(stdout)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			result.Stderr = string(exitErr.Stderr)
		} else {
			result.Error = err
		}
	} else {
		result.Success = true
	}

	return result, nil
}

// GetEnvironment returns the environment type
func (r *SimpleLocalRunner) GetEnvironment() string {
	return "local"
}

// SetWorkingDirectory changes the working directory
func (r *SimpleLocalRunner) SetWorkingDirectory(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	r.workingDir = absPath
	return nil
}

// GetWorkingDirectory returns the current working directory
func (r *SimpleLocalRunner) GetWorkingDirectory() string {
	return r.workingDir
}

// Cleanup performs any necessary cleanup
func (r *SimpleLocalRunner) Cleanup() error {
	// For local runner, optionally clean up temporary files
	if strings.HasPrefix(r.workingDir, "/tmp/") {
		return os.RemoveAll(r.workingDir)
	}
	return nil
}

// IsHealthy checks if the runner is healthy
func (r *SimpleLocalRunner) IsHealthy(ctx context.Context) error {
	// Simple health check - try to execute a basic command
	result, err := r.Execute(ctx, "echo 'health-check'")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	if !result.Success {
		return fmt.Errorf("health check command failed with exit code %d", result.ExitCode)
	}
	return nil
}

// NewDockerRunner creates a new Docker-based runner (placeholder)
func NewDockerRunner(config DockerConfig) Runner {
	// TODO: Implement Docker runner
	// For now, return a local runner as fallback
	localConfig := LocalConfig{
		WorkingDirectory: config.WorkingDirectory,
		Environment:      config.Environment,
		Timeout:          config.Timeout,
	}
	return NewLocalRunner(localConfig)
}

// NewAWSRunner creates a new AWS-based runner (placeholder)
func NewAWSRunner(config AWSConfig) Runner {
	// TODO: Implement AWS runner
	// For now, return a local runner as fallback
	localConfig := LocalConfig{
		WorkingDirectory: "/tmp/aws-runner",
		Environment:      config.Environment,
		Timeout:          config.Timeout,
	}
	return NewLocalRunner(localConfig)
}