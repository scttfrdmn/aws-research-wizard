package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DomainPack represents a research domain configuration
type DomainPack struct {
	Name                       string                            `yaml:"name"`
	Description                string                            `yaml:"description"`
	PrimaryDomains             []string                          `yaml:"primary_domains"`
	TargetUsers                string                            `yaml:"target_users"`
	SpackPackages              map[string]interface{}            `yaml:"spack_packages"`
	SystemPackages             map[string]interface{}            `yaml:"system_packages"`
	PythonPackages             map[string]interface{}            `yaml:"python_packages"`
	RPackages                  map[string]interface{}            `yaml:"r_packages"`
	JuliaPackages              map[string]interface{}            `yaml:"julia_packages"`
	AWSInstanceRecommendations map[string]InstanceRecommendation `yaml:"aws_instance_recommendations"`
	EstimatedCost              EstimatedCost                     `yaml:"estimated_cost"`
	WorkflowOrchestration      WorkflowOrchestration             `yaml:"workflow_orchestration"`
	AWSIntegration             AWSIntegration                    `yaml:"aws_integration"`
}

// InstanceRecommendation represents AWS instance recommendations
type InstanceRecommendation struct {
	UseCase      string  `yaml:"use_case"`
	InstanceType string  `yaml:"instance_type"`
	VCPUs        int     `yaml:"vcpus"`
	MemoryGB     int     `yaml:"memory_gb"`
	StorageGB    int     `yaml:"storage_gb"`
	CostPerHour  float64 `yaml:"cost_per_hour"`
}

// EstimatedCost represents cost breakdown
type EstimatedCost struct {
	Compute float64 `yaml:"compute"`
	Storage float64 `yaml:"storage"`
	Total   float64 `yaml:"total"`
}

// WorkflowOrchestration represents workflow tools
type WorkflowOrchestration struct {
	Tools []WorkflowTool `yaml:"tools"`
}

// WorkflowTool represents a workflow management tool
type WorkflowTool struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	S3Support   bool   `yaml:"s3_support"`
}

// AWSIntegration represents AWS-specific configurations
type AWSIntegration struct {
	DataSources     []string               `yaml:"data_sources"`
	StoragePatterns []string               `yaml:"storage_patterns"`
	OptimizedFor    []string               `yaml:"optimized_for"`
	CostStrategy    map[string]interface{} `yaml:"cost_strategy"`
}

// ConfigLoader handles loading domain configurations
type ConfigLoader struct {
	configRoot string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configRoot string) *ConfigLoader {
	return &ConfigLoader{
		configRoot: configRoot,
	}
}

// LoadAllDomains loads all domain pack configurations
func (cl *ConfigLoader) LoadAllDomains() (map[string]*DomainPack, error) {
	domains := make(map[string]*DomainPack)
	domainsPath := filepath.Join(cl.configRoot, "configs", "domains")

	err := filepath.WalkDir(domainsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		domain, err := cl.LoadDomain(path)
		if err != nil {
			return fmt.Errorf("failed to load domain %s: %w", path, err)
		}

		domainName := strings.TrimSuffix(d.Name(), ".yaml")
		domains[domainName] = domain

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk domains directory: %w", err)
	}

	return domains, nil
}

// LoadDomain loads a single domain pack configuration
func (cl *ConfigLoader) LoadDomain(path string) (*DomainPack, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var domain DomainPack
	err = yaml.Unmarshal(data, &domain)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &domain, nil
}

// GetDomainNames returns a list of all available domain names
func (cl *ConfigLoader) GetDomainNames() ([]string, error) {
	domains, err := cl.LoadAllDomains()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(domains))
	for name := range domains {
		names = append(names, name)
	}

	return names, nil
}
