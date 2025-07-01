package manager

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SpackManager handles Spack package installations and environment management
type SpackManager struct {
	config        Config
	installMutex  sync.Mutex
	progressChans map[string]chan ProgressUpdate
	mu            sync.RWMutex
}

// New creates a new SpackManager instance
func New(config Config) (*SpackManager, error) {
	// Set default values
	if config.SpackRoot == "" {
		config.SpackRoot = "/opt/spack"
	}
	if config.WorkDir == "" {
		config.WorkDir = "/tmp/spack-manager"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	// Validate Spack installation
	sm := &SpackManager{
		config:        config,
		progressChans: make(map[string]chan ProgressUpdate),
	}

	if err := sm.validateSpackInstallation(); err != nil {
		return nil, fmt.Errorf("invalid Spack installation: %w", err)
	}

	return sm, nil
}

// validateSpackInstallation checks if Spack is properly installed
func (sm *SpackManager) validateSpackInstallation() error {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")
	if _, err := os.Stat(spackBin); os.IsNotExist(err) {
		return fmt.Errorf("spack binary not found at %s", spackBin)
	}

	// Test spack command
	cmd := exec.Command(spackBin, "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run spack command: %w", err)
	}

	return nil
}

// GetSpackVersion returns the installed Spack version
func (sm *SpackManager) GetSpackVersion() (string, error) {
	cmd := exec.Command(filepath.Join(sm.config.SpackRoot, "bin", "spack"), "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get spack version: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateEnvironment creates a new Spack environment
func (sm *SpackManager) CreateEnvironment(env Environment) error {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")

	// Create environment
	cmd := exec.Command(spackBin, "env", "create", env.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create environment %s: %w", env.Name, err)
	}

	// Add packages to environment
	for _, pkg := range env.Packages {
		cmd = exec.Command(spackBin, "-e", env.Name, "add", pkg)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add package %s to environment %s: %w", pkg, env.Name, err)
		}
	}

	return nil
}

// ListEnvironments lists all available Spack environments
func (sm *SpackManager) ListEnvironments() ([]Environment, error) {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")
	cmd := exec.Command(spackBin, "env", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	var environments []Environment
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "=") {
			env := Environment{Name: line}
			environments = append(environments, env)
		}
	}

	return environments, nil
}

// DeleteEnvironment removes a Spack environment
func (sm *SpackManager) DeleteEnvironment(name string) error {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")
	cmd := exec.Command(spackBin, "env", "remove", "-y", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete environment %s: %w", name, err)
	}
	return nil
}

// ActivateEnvironment activates a Spack environment
func (sm *SpackManager) ActivateEnvironment(name string) error {
	// Note: Environment activation is typically done in shell context
	// This method validates that the environment exists
	envs, err := sm.ListEnvironments()
	if err != nil {
		return err
	}

	for _, env := range envs {
		if env.Name == name {
			return nil
		}
	}

	return fmt.Errorf("environment %s not found", name)
}

// GetEnvironmentInfo retrieves detailed information about an environment
func (sm *SpackManager) GetEnvironmentInfo(name string) (*Environment, error) {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")

	// Get environment packages
	cmd := exec.Command(spackBin, "-e", name, "find", "--format", "{name}@{version}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get environment info: %w", err)
	}

	env := &Environment{
		Name:     name,
		Packages: []string{},
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "=") {
			env.Packages = append(env.Packages, line)
		}
	}

	return env, nil
}

// ValidateEnvironment validates an environment configuration
func (sm *SpackManager) ValidateEnvironment(env Environment) error {
	if env.Name == "" {
		return fmt.Errorf("environment name cannot be empty")
	}

	if len(env.Packages) == 0 {
		return fmt.Errorf("environment must have at least one package")
	}

	// TODO: Add more validation logic (package specs, compiler validation, etc.)
	return nil
}

// InstallEnvironment installs all packages in a Spack environment with progress tracking
func (sm *SpackManager) InstallEnvironment(envName string, progressChan chan<- ProgressUpdate) error {
	sm.installMutex.Lock()
	defer sm.installMutex.Unlock()

	// Create internal progress channel for tracking
	internalProgressChan := make(chan ProgressUpdate, 100)

	sm.mu.Lock()
	sm.progressChans[envName] = internalProgressChan
	sm.mu.Unlock()

	defer func() {
		sm.mu.Lock()
		delete(sm.progressChans, envName)
		sm.mu.Unlock()
		close(internalProgressChan)
	}()

	// Forward progress updates to external channel
	go func() {
		for update := range internalProgressChan {
			if progressChan != nil {
				select {
				case progressChan <- update:
				default:
					// Channel is full, skip update
				}
			}
		}
	}()

	// Start installation
	return sm.installEnvironmentWithProgress(context.Background(), envName, internalProgressChan)
}

// installEnvironmentWithProgress performs the actual installation with progress tracking
func (sm *SpackManager) installEnvironmentWithProgress(ctx context.Context, envName string, progressChan chan<- ProgressUpdate) error {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")

	// Send initial progress
	progressChan <- ProgressUpdate{
		Package:   "environment",
		Stage:     "starting",
		Progress:  0.0,
		Message:   "Starting environment installation",
		Timestamp: time.Now(),
	}

	// Install environment
	cmd := exec.CommandContext(ctx, spackBin, "-e", envName, "install")

	// Set up binary cache if specified
	if sm.config.BinaryCache != "" {
		cmd.Env = append(os.Environ(), "SPACK_BINARY_CACHE="+sm.config.BinaryCache)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start installation: %w", err)
	}

	// Monitor progress
	go sm.monitorInstallationProgress(stdout, stderr, progressChan)

	// Wait for installation to complete
	if err := cmd.Wait(); err != nil {
		progressChan <- ProgressUpdate{
			Package:   "environment",
			Stage:     "failed",
			Progress:  0.0,
			Message:   fmt.Sprintf("Installation failed: %v", err),
			Timestamp: time.Now(),
			IsError:   true,
		}
		return fmt.Errorf("installation failed: %w", err)
	}

	// Send completion progress
	progressChan <- ProgressUpdate{
		Package:    "environment",
		Stage:      "completed",
		Progress:   1.0,
		Message:    "Environment installation completed successfully",
		Timestamp:  time.Now(),
		IsComplete: true,
	}

	return nil
}

// monitorInstallationProgress monitors the installation output and sends progress updates
func (sm *SpackManager) monitorInstallationProgress(stdout, stderr io.Reader, progressChan chan<- ProgressUpdate) {
	// Progress regex patterns
	packageRegex := regexp.MustCompile(`Installing\s+(\S+)`)
	progressRegex := regexp.MustCompile(`(\d+)%`)

	// Monitor stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			if matches := packageRegex.FindStringSubmatch(line); len(matches) > 1 {
				packageName := matches[1]
				progressChan <- ProgressUpdate{
					Package:   packageName,
					Stage:     "installing",
					Progress:  0.0,
					Message:   fmt.Sprintf("Installing %s", packageName),
					Timestamp: time.Now(),
				}
			}

			if matches := progressRegex.FindStringSubmatch(line); len(matches) > 1 {
				if progress, err := strconv.ParseFloat(matches[1], 64); err == nil {
					progressChan <- ProgressUpdate{
						Package:   "current",
						Stage:     "installing",
						Progress:  progress / 100.0,
						Message:   line,
						Timestamp: time.Now(),
					}
				}
			}
		}
	}()

	// Monitor stderr for errors
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			progressChan <- ProgressUpdate{
				Package:   "environment",
				Stage:     "error",
				Progress:  0.0,
				Message:   line,
				Timestamp: time.Now(),
				IsError:   true,
			}
		}
	}()
}

// InstallPackage installs a single package in an environment
func (sm *SpackManager) InstallPackage(envName, packageSpec string, progressChan chan<- ProgressUpdate) error {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")

	// Add package to environment
	cmd := exec.Command(spackBin, "-e", envName, "add", packageSpec)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add package %s: %w", packageSpec, err)
	}

	// Install the package
	cmd = exec.Command(spackBin, "-e", envName, "install", packageSpec)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s: %w", packageSpec, err)
	}

	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Package:    packageSpec,
			Stage:      "completed",
			Progress:   1.0,
			Message:    fmt.Sprintf("Package %s installed successfully", packageSpec),
			Timestamp:  time.Now(),
			IsComplete: true,
		}
	}

	return nil
}

// UninstallPackage removes a package from an environment
func (sm *SpackManager) UninstallPackage(envName, packageSpec string) error {
	spackBin := filepath.Join(sm.config.SpackRoot, "bin", "spack")
	cmd := exec.Command(spackBin, "-e", envName, "remove", packageSpec)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", packageSpec, err)
	}
	return nil
}

// ListPackages lists all packages in an environment
func (sm *SpackManager) ListPackages(envName string) ([]string, error) {
	env, err := sm.GetEnvironmentInfo(envName)
	if err != nil {
		return nil, err
	}
	return env.Packages, nil
}

// GetProgressChannel returns the progress channel for an environment
func (sm *SpackManager) GetProgressChannel(envName string) (<-chan ProgressUpdate, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	ch, exists := sm.progressChans[envName]
	if !exists {
		return nil, false
	}
	return ch, true
}
