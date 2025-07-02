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
	"time"
)

// Runner defines the interface for executing commands in different environments
type Runner interface {
	// Execute runs a command and returns the result
	Execute(ctx context.Context, command string) (*RunResult, error)
	
	// GetEnvironment returns the current environment type
	GetEnvironment() string
	
	// SetWorkingDirectory changes the working directory for subsequent commands
	SetWorkingDirectory(path string) error
	
	// GetWorkingDirectory returns the current working directory
	GetWorkingDirectory() string
	
	// Cleanup performs any necessary cleanup
	Cleanup() error
	
	// IsHealthy checks if the runner is in a healthy state
	IsHealthy(ctx context.Context) error
}

// RunResult represents the result of a command execution
type RunResult struct {
	Command      string        `json:"command"`
	ExitCode     int           `json:"exit_code"`
	Stdout       string        `json:"stdout"`
	Stderr       string        `json:"stderr"`
	Duration     time.Duration `json:"duration"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Environment  string        `json:"environment"`
	WorkingDir   string        `json:"working_dir"`
	Success      bool          `json:"success"`
	Error        error         `json:"error,omitempty"`
}

// ResourceLimits defines resource constraints for command execution
type ResourceLimits struct {
	MaxMemoryMB    int           `json:"max_memory_mb"`
	MaxDiskMB      int           `json:"max_disk_mb"`
	MaxCPUPercent  float64       `json:"max_cpu_percent"`
	MaxProcesses   int           `json:"max_processes"`
	MaxFileHandles int           `json:"max_file_handles"`
	Timeout        time.Duration `json:"timeout"`
}

// LocalConfig defines configuration for local command execution
type LocalConfig struct {
	WorkingDirectory string            `json:"working_directory"`
	Environment      map[string]string `json:"environment"`
	ResourceLimits   ResourceLimits    `json:"resource_limits"`
	Timeout          time.Duration     `json:"timeout"`
	Shell            string            `json:"shell"` // Default: /bin/sh
}

// DockerConfig defines configuration for Docker-based execution
type DockerConfig struct {
	Image            string            `json:"image"`
	Tag              string            `json:"tag"`
	WorkingDirectory string            `json:"working_directory"`
	Environment      map[string]string `json:"environment"`
	Volumes          []string          `json:"volumes"`
	NetworkMode      string            `json:"network_mode"`
	Timeout          time.Duration     `json:"timeout"`
	CleanupOnExit    bool              `json:"cleanup_on_exit"`
}

// AWSConfig defines configuration for AWS-based execution
type AWSConfig struct {
	Region          string            `json:"region"`
	InstanceType    string            `json:"instance_type"`
	AMI             string            `json:"ami"`
	SecurityGroups  []string          `json:"security_groups"`
	SubnetID        string            `json:"subnet_id"`
	KeyPair         string            `json:"key_pair"`
	Environment     map[string]string `json:"environment"`
	Timeout         time.Duration     `json:"timeout"`
	TerminateOnExit bool              `json:"terminate_on_exit"`
}