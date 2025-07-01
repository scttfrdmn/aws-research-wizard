package data

import (
	"bufio"
	"context"
	"encoding/json"
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

	"github.com/aws/aws-sdk-go/service/s3"
)

// SpackManager handles Spack package installations and environment management
type SpackManager struct {
	spackRoot     string
	binaryCache   string
	awsRegion     string
	s3Client      *s3.S3
	installMutex  sync.Mutex
	progressChans map[string]chan ProgressUpdate
	mu            sync.RWMutex
}

// ProgressUpdate represents installation progress information
type ProgressUpdate struct {
	Package    string    `json:"package"`
	Stage      string    `json:"stage"`
	Progress   float64   `json:"progress"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	IsError    bool      `json:"is_error"`
	IsComplete bool      `json:"is_complete"`
}

// SpackEnvironment represents a Spack environment configuration
type SpackEnvironment struct {
	Name        string                 `yaml:"name" json:"name"`
	Packages    []string               `yaml:"packages" json:"packages"`
	Compilers   []string               `yaml:"compilers" json:"compilers"`
	Config      map[string]interface{} `yaml:"config" json:"config"`
	View        bool                   `yaml:"view" json:"view"`
	Concretized bool                   `yaml:"concretized" json:"concretized"`
}

// InstallationResult contains the result of a package installation
type InstallationResult struct {
	Package      string        `json:"package"`
	Success      bool          `json:"success"`
	Duration     time.Duration `json:"duration"`
	ErrorMessage string        `json:"error_message,omitempty"`
	LogPath      string        `json:"log_path,omitempty"`
}

// NewSpackManager creates a new SpackManager instance
func NewSpackManager(spackRoot, binaryCache, awsRegion string, s3Client *s3.S3) *SpackManager {
	return &SpackManager{
		spackRoot:     spackRoot,
		binaryCache:   binaryCache,
		awsRegion:     awsRegion,
		s3Client:      s3Client,
		progressChans: make(map[string]chan ProgressUpdate),
	}
}

// InstallEnvironment installs all packages in a Spack environment with progress tracking
func (sm *SpackManager) InstallEnvironment(ctx context.Context, env *SpackEnvironment) (*chan ProgressUpdate, error) {
	sm.installMutex.Lock()
	defer sm.installMutex.Unlock()

	// Create progress channel
	progressChan := make(chan ProgressUpdate, 100)

	sm.mu.Lock()
	sm.progressChans[env.Name] = progressChan
	sm.mu.Unlock()

	go func() {
		defer close(progressChan)
		defer func() {
			sm.mu.Lock()
			delete(sm.progressChans, env.Name)
			sm.mu.Unlock()
		}()

		// Create environment
		if err := sm.createEnvironment(ctx, env, progressChan); err != nil {
			progressChan <- ProgressUpdate{
				Package:   "environment",
				Stage:     "creation",
				Message:   fmt.Sprintf("Failed to create environment: %v", err),
				Timestamp: time.Now(),
				IsError:   true,
			}
			return
		}

		// Install packages
		totalPackages := len(env.Packages)
		for i, pkg := range env.Packages {
			select {
			case <-ctx.Done():
				progressChan <- ProgressUpdate{
					Package:   pkg,
					Stage:     "cancelled",
					Message:   "Installation cancelled",
					Timestamp: time.Now(),
					IsError:   true,
				}
				return
			default:
				progress := float64(i) / float64(totalPackages) * 100
				if err := sm.installPackage(ctx, env.Name, pkg, progressChan, progress); err != nil {
					progressChan <- ProgressUpdate{
						Package:   pkg,
						Stage:     "failed",
						Message:   fmt.Sprintf("Installation failed: %v", err),
						Timestamp: time.Now(),
						IsError:   true,
					}
					continue
				}
			}
		}

		// Final progress update
		progressChan <- ProgressUpdate{
			Package:    "environment",
			Stage:      "complete",
			Progress:   100.0,
			Message:    fmt.Sprintf("Environment %s installation complete", env.Name),
			Timestamp:  time.Now(),
			IsComplete: true,
		}
	}()

	return &progressChan, nil
}

// createEnvironment creates a new Spack environment
func (sm *SpackManager) createEnvironment(ctx context.Context, env *SpackEnvironment, progressChan chan ProgressUpdate) error {
	progressChan <- ProgressUpdate{
		Package:   "environment",
		Stage:     "creating",
		Progress:  5.0,
		Message:   fmt.Sprintf("Creating environment %s", env.Name),
		Timestamp: time.Now(),
	}

	// Create environment directory
	envPath := filepath.Join(sm.spackRoot, "var", "spack", "environments", env.Name)
	if err := os.MkdirAll(envPath, 0755); err != nil {
		return fmt.Errorf("failed to create environment directory: %w", err)
	}

	// Generate spack.yaml
	envConfig := map[string]interface{}{
		"spack": map[string]interface{}{
			"specs": env.Packages,
			"view":  env.View,
			"config": map[string]interface{}{
				"install_tree": map[string]interface{}{
					"root": "$spack/opt/spack",
				},
				"build_stage": []string{"$tempdir/$user/spack-stage"},
			},
		},
	}

	if sm.binaryCache != "" {
		envConfig["spack"].(map[string]interface{})["config"].(map[string]interface{})["binary_index_root"] = sm.binaryCache
		envConfig["spack"].(map[string]interface{})["mirrors"] = map[string]interface{}{
			"aws-binary-cache": sm.binaryCache,
		}
	}

	// Write spack.yaml
	envFile := filepath.Join(envPath, "spack.yaml")
	if err := sm.writeYAMLFile(envFile, envConfig); err != nil {
		return fmt.Errorf("failed to write environment config: %w", err)
	}

	progressChan <- ProgressUpdate{
		Package:   "environment",
		Stage:     "created",
		Progress:  10.0,
		Message:   fmt.Sprintf("Environment %s created successfully", env.Name),
		Timestamp: time.Now(),
	}

	return nil
}

// installPackage installs a single package with progress monitoring
func (sm *SpackManager) installPackage(ctx context.Context, envName, pkg string, progressChan chan ProgressUpdate, baseProgress float64) error {
	startTime := time.Now()

	progressChan <- ProgressUpdate{
		Package:   pkg,
		Stage:     "starting",
		Progress:  baseProgress,
		Message:   fmt.Sprintf("Starting installation of %s", pkg),
		Timestamp: startTime,
	}

	// Check if package is already installed
	if installed, err := sm.isPackageInstalled(envName, pkg); err != nil {
		return fmt.Errorf("failed to check if package is installed: %w", err)
	} else if installed {
		progressChan <- ProgressUpdate{
			Package:   pkg,
			Stage:     "skipped",
			Progress:  baseProgress + 10,
			Message:   fmt.Sprintf("Package %s already installed", pkg),
			Timestamp: time.Now(),
		}
		return nil
	}

	// Run spack install with progress monitoring
	cmd := exec.CommandContext(ctx, "spack", "-e", envName, "install", "-v", pkg)
	cmd.Dir = sm.spackRoot

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start spack install: %w", err)
	}

	// Monitor progress from stdout and stderr
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sm.monitorInstallationOutput(stdout, pkg, progressChan, baseProgress, false)
	}()

	go func() {
		defer wg.Done()
		sm.monitorInstallationOutput(stderr, pkg, progressChan, baseProgress, true)
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		duration := time.Since(startTime)
		progressChan <- ProgressUpdate{
			Package:   pkg,
			Stage:     "failed",
			Progress:  baseProgress,
			Message:   fmt.Sprintf("Installation failed after %v: %v", duration, err),
			Timestamp: time.Now(),
			IsError:   true,
		}
		return fmt.Errorf("spack install failed: %w", err)
	}

	duration := time.Since(startTime)
	progressChan <- ProgressUpdate{
		Package:   pkg,
		Stage:     "completed",
		Progress:  baseProgress + 10,
		Message:   fmt.Sprintf("Package %s installed successfully in %v", pkg, duration),
		Timestamp: time.Now(),
	}

	return nil
}

// monitorInstallationOutput monitors spack install output for progress information
func (sm *SpackManager) monitorInstallationOutput(reader io.Reader, pkg string, progressChan chan ProgressUpdate, baseProgress float64, isError bool) {
	scanner := bufio.NewScanner(reader)

	// Regex patterns for different installation stages
	patterns := map[string]*regexp.Regexp{
		"fetch":     regexp.MustCompile(`==> Fetching`),
		"stage":     regexp.MustCompile(`==> Staging`),
		"patch":     regexp.MustCompile(`==> Patching`),
		"configure": regexp.MustCompile(`==> Configuring`),
		"build":     regexp.MustCompile(`==> Building`),
		"install":   regexp.MustCompile(`==> Installing`),
		"progress":  regexp.MustCompile(`\[(\d+)%\]`),
	}

	var currentStage string
	var stageProgress float64

	for scanner.Scan() {
		line := scanner.Text()

		// Check for stage changes
		for stage, pattern := range patterns {
			if stage == "progress" {
				continue
			}
			if pattern.MatchString(line) {
				currentStage = stage
				stageProgress = 0
				progressChan <- ProgressUpdate{
					Package:   pkg,
					Stage:     currentStage,
					Progress:  baseProgress + stageProgress,
					Message:   fmt.Sprintf("%s: %s", strings.Title(currentStage), line),
					Timestamp: time.Now(),
					IsError:   isError,
				}
				break
			}
		}

		// Check for progress percentage
		if match := patterns["progress"].FindStringSubmatch(line); len(match) > 1 {
			if percent, err := strconv.ParseFloat(match[1], 64); err == nil {
				stageProgress = percent / 10.0 // Each stage is roughly 10% of total progress
				progressChan <- ProgressUpdate{
					Package:   pkg,
					Stage:     currentStage,
					Progress:  baseProgress + stageProgress,
					Message:   line,
					Timestamp: time.Now(),
					IsError:   isError,
				}
			}
		}

		// Send any error lines
		if isError && strings.Contains(line, "error") {
			progressChan <- ProgressUpdate{
				Package:   pkg,
				Stage:     "error",
				Progress:  baseProgress + stageProgress,
				Message:   line,
				Timestamp: time.Now(),
				IsError:   true,
			}
		}
	}
}

// isPackageInstalled checks if a package is already installed in the environment
func (sm *SpackManager) isPackageInstalled(envName, pkg string) (bool, error) {
	cmd := exec.Command("spack", "-e", envName, "find", pkg)
	cmd.Dir = sm.spackRoot

	output, err := cmd.Output()
	if err != nil {
		return false, nil // Package not found
	}

	return len(output) > 0 && !strings.Contains(string(output), "No package matches"), nil
}

// ValidateEnvironment validates a Spack environment configuration
func (sm *SpackManager) ValidateEnvironment(env *SpackEnvironment) error {
	if env.Name == "" {
		return fmt.Errorf("environment name cannot be empty")
	}

	if len(env.Packages) == 0 {
		return fmt.Errorf("environment must contain at least one package")
	}

	// Validate package specifications
	for _, pkg := range env.Packages {
		if err := sm.validatePackageSpec(pkg); err != nil {
			return fmt.Errorf("invalid package spec %s: %w", pkg, err)
		}
	}

	return nil
}

// validatePackageSpec validates a single package specification
func (sm *SpackManager) validatePackageSpec(spec string) error {
	if spec == "" {
		return fmt.Errorf("package specification cannot be empty")
	}

	// Basic validation - could be enhanced with more sophisticated parsing
	if strings.Contains(spec, " ") {
		return fmt.Errorf("package specification contains spaces")
	}

	return nil
}

// GetEnvironmentStatus returns the current status of a Spack environment
func (sm *SpackManager) GetEnvironmentStatus(envName string) (*SpackEnvironment, error) {
	envPath := filepath.Join(sm.spackRoot, "var", "spack", "environments", envName)
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("environment %s does not exist", envName)
	}

	// Read spack.yaml
	envFile := filepath.Join(envPath, "spack.yaml")
	data, err := os.ReadFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment config: %w", err)
	}

	// Parse environment configuration
	var envConfig map[string]interface{}
	if err := json.Unmarshal(data, &envConfig); err != nil {
		return nil, fmt.Errorf("failed to parse environment config: %w", err)
	}

	// Extract packages
	spackConfig := envConfig["spack"].(map[string]interface{})
	specs := spackConfig["specs"].([]interface{})

	var packages []string
	for _, spec := range specs {
		packages = append(packages, spec.(string))
	}

	return &SpackEnvironment{
		Name:     envName,
		Packages: packages,
		View:     spackConfig["view"].(bool),
	}, nil
}

// CleanupEnvironment removes a Spack environment
func (sm *SpackManager) CleanupEnvironment(envName string) error {
	envPath := filepath.Join(sm.spackRoot, "var", "spack", "environments", envName)
	return os.RemoveAll(envPath)
}

// GetProgressChannel returns the progress channel for an ongoing installation
func (sm *SpackManager) GetProgressChannel(envName string) (chan ProgressUpdate, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	ch, exists := sm.progressChans[envName]
	return ch, exists
}

// writeYAMLFile writes data to a YAML file
func (sm *SpackManager) writeYAMLFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert to JSON first, then to YAML-like format (simplified)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Simple JSON to YAML conversion for this use case
	yamlData := strings.ReplaceAll(string(jsonData), "\"", "")
	yamlData = strings.ReplaceAll(yamlData, "{", "")
	yamlData = strings.ReplaceAll(yamlData, "}", "")
	yamlData = strings.ReplaceAll(yamlData, "[", "")
	yamlData = strings.ReplaceAll(yamlData, "]", "")
	yamlData = strings.ReplaceAll(yamlData, ",", "")

	_, err = file.WriteString(yamlData)
	return err
}
