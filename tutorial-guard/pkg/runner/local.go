package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// NewLocalEnvironment creates a new local environment for running tests
func NewLocalEnvironment(config Config) Environment {
	workDir := config.WorkDir
	if workDir == "" {
		workDir = "/tmp/tutorial-guard"
	}

	shell := "/bin/bash"
	if envShell := os.Getenv("SHELL"); envShell != "" {
		shell = envShell
	}

	return &LocalEnvironment{
		WorkDir:   workDir,
		EnvVars:   config.EnvVars,
		Shell:     shell,
		Resources: []Resource{},
	}
}

// Setup prepares the local environment for test execution
func (e *LocalEnvironment) Setup(ctx context.Context) error {
	// Create working directory
	if err := ensureWorkDir(e.WorkDir); err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}

	// Record the working directory as a resource
	e.Resources = append(e.Resources, Resource{
		Type:       "directory",
		Identifier: e.WorkDir,
		CreatedAt:  time.Now(),
		Metadata: map[string]string{
			"purpose": "working_directory",
		},
	})

	return nil
}

// Execute runs a code example in the local environment
func (e *LocalEnvironment) Execute(ctx context.Context, example extractor.Example) (*TestResult, error) {
	startTime := time.Now()

	result := &TestResult{
		ExampleID:   example.ID,
		StartTime:   startTime,
		Environment: "local",
		Metadata:    make(map[string]string),
	}

	// Only execute supported languages
	if !e.isExecutableLanguage(example.Language) {
		result.Success = false
		result.Errors = []string{fmt.Sprintf("Language %s is not executable in local environment", example.Language)}
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Create a temporary script file
	scriptPath, err := e.createScriptFile(example)
	if err != nil {
		result.Success = false
		result.Errors = []string{fmt.Sprintf("Failed to create script file: %v", err)}
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Record script file as a resource
	e.Resources = append(e.Resources, Resource{
		Type:       "file",
		Identifier: scriptPath,
		CreatedAt:  time.Now(),
		Metadata: map[string]string{
			"purpose":    "script",
			"example_id": example.ID,
			"language":   example.Language,
		},
	})

	// Execute the script
	cmd := e.createCommand(example.Language, scriptPath)
	cmd.Dir = e.WorkDir

	// Set environment variables
	cmd.Env = append(os.Environ(), e.EnvVars...)

	// Capture output
	output, err := cmd.CombinedOutput()

	result.Output = string(output)
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if err != nil {
		result.Success = false
		result.ErrorOutput = err.Error()

		// Try to get exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}

		result.Errors = []string{fmt.Sprintf("Command execution failed: %v", err)}
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	// Detect any files or directories created during execution
	e.detectCreatedResources(result, example)

	// Add execution metadata
	result.Metadata["script_path"] = scriptPath
	result.Metadata["working_directory"] = e.WorkDir
	result.Metadata["shell"] = e.Shell

	return result, nil
}

// Cleanup removes resources created during test execution
func (e *LocalEnvironment) Cleanup(ctx context.Context) error {
	var errors []string

	for _, resource := range e.Resources {
		switch resource.Type {
		case "file":
			if err := os.Remove(resource.Identifier); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("Failed to remove file %s: %v", resource.Identifier, err))
			}
		case "directory":
			// Only remove if it's our working directory
			if resource.Identifier == e.WorkDir {
				if err := os.RemoveAll(resource.Identifier); err != nil && !os.IsNotExist(err) {
					errors = append(errors, fmt.Sprintf("Failed to remove directory %s: %v", resource.Identifier, err))
				}
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetResources returns the list of resources created
func (e *LocalEnvironment) GetResources() []Resource {
	return e.Resources
}

// createScriptFile creates a temporary script file for the code example
func (e *LocalEnvironment) createScriptFile(example extractor.Example) (string, error) {
	var extension, content string

	switch example.Language {
	case "bash", "sh", "shell":
		extension = ".sh"
		content = "#!/bin/bash\nset -e\n\n" + example.Code
	case "go":
		extension = ".go"
		// For Go, we need to handle different types of code
		if strings.Contains(example.Code, "package main") {
			content = example.Code
		} else {
			// Wrap in a simple main package
			content = "package main\n\nimport \"fmt\"\n\nfunc main() {\n" + example.Code + "\n}"
		}
	case "python":
		extension = ".py"
		content = "#!/usr/bin/env python3\n\n" + example.Code
	default:
		return "", fmt.Errorf("unsupported language for script creation: %s", example.Language)
	}

	// Create temporary file
	filename := fmt.Sprintf("tutorial_guard_%s_%d%s", example.ID, time.Now().Unix(), extension)
	scriptPath := filepath.Join(e.WorkDir, filename)

	if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
		return "", fmt.Errorf("failed to write script file: %w", err)
	}

	return scriptPath, nil
}

// createCommand creates the appropriate command for executing the script
func (e *LocalEnvironment) createCommand(language, scriptPath string) *exec.Cmd {
	switch language {
	case "bash", "sh", "shell":
		return exec.Command(e.Shell, scriptPath)
	case "go":
		return exec.Command("go", "run", scriptPath)
	case "python":
		return exec.Command("python3", scriptPath)
	default:
		// Fallback to shell execution
		return exec.Command(e.Shell, scriptPath)
	}
}

// isExecutableLanguage checks if a language can be executed locally
func (e *LocalEnvironment) isExecutableLanguage(language string) bool {
	executable := map[string]bool{
		"bash":   true,
		"sh":     true,
		"shell":  true,
		"go":     true,
		"python": true,
		"yaml":   false,
		"json":   false,
		"text":   false,
	}

	return executable[language]
}

// detectCreatedResources attempts to detect files/directories created during execution
func (e *LocalEnvironment) detectCreatedResources(result *TestResult, example extractor.Example) {
	// This is a simple implementation - could be enhanced with filesystem monitoring

	// Check for common patterns in the code that suggest file creation
	code := strings.ToLower(example.Code)

	if strings.Contains(code, "touch ") || strings.Contains(code, "> ") || strings.Contains(code, ">> ") {
		result.Metadata["may_create_files"] = "true"
	}

	if strings.Contains(code, "mkdir") {
		result.Metadata["may_create_directories"] = "true"
	}

	if strings.Contains(code, "wget ") || strings.Contains(code, "curl ") {
		result.Metadata["may_download_files"] = "true"
	}

	if strings.Contains(code, "git clone") {
		result.Metadata["may_clone_repository"] = "true"
	}

	// Could be enhanced to:
	// 1. Take filesystem snapshot before execution
	// 2. Compare with snapshot after execution
	// 3. Record actual created files as resources
}
