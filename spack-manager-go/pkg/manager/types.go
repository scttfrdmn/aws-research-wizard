package manager

import (
	"time"
)

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

// Environment represents a Spack environment configuration
type Environment struct {
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

// Config represents the configuration for the Spack manager
type Config struct {
	SpackRoot   string `yaml:"spack_root" json:"spack_root"`
	BinaryCache string `yaml:"binary_cache" json:"binary_cache"`
	WorkDir     string `yaml:"work_dir" json:"work_dir"`
	LogLevel    string `yaml:"log_level" json:"log_level"`
}

// Manager defines the interface for Spack management operations
type Manager interface {
	// Environment management
	CreateEnvironment(env Environment) error
	ListEnvironments() ([]Environment, error)
	DeleteEnvironment(name string) error
	ActivateEnvironment(name string) error

	// Package management
	InstallPackage(envName, packageSpec string, progressChan chan<- ProgressUpdate) error
	InstallEnvironment(envName string, progressChan chan<- ProgressUpdate) error
	UninstallPackage(envName, packageSpec string) error
	ListPackages(envName string) ([]string, error)

	// Information and validation
	GetEnvironmentInfo(name string) (*Environment, error)
	ValidateEnvironment(env Environment) error
	GetSpackVersion() (string, error)

	// Progress tracking
	GetProgressChannel(envName string) (<-chan ProgressUpdate, bool)
}
