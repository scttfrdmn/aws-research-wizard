package intelligence

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DomainPackLoader loads and manages domain pack configurations
type DomainPackLoader struct {
	domainPacksPath string
	cache           map[string]*DomainPackInfo
}

// DomainPackConfig represents the structure of domain-pack.yaml files
type DomainPackConfig struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Version     string       `yaml:"version"`
	Categories  []string     `yaml:"categories"`
	Maintainers []Maintainer `yaml:"maintainers"`

	SpackConfig SpackConfig      `yaml:"spack_config"`
	AWSConfig   AWSConfig        `yaml:"aws_config"`
	Workflows   []WorkflowConfig `yaml:"workflows"`

	CostEstimates map[string]string `yaml:"cost_estimates"`
	Documentation Documentation     `yaml:"documentation"`
}

// Maintainer represents domain pack maintainer information
type Maintainer struct {
	Name         string `yaml:"name"`
	Email        string `yaml:"email"`
	Organization string `yaml:"organization"`
}

// SpackConfig represents the Spack configuration section
type SpackConfig struct {
	Packages     []string `yaml:"packages"`
	Compiler     string   `yaml:"compiler"`
	Target       string   `yaml:"target"`
	Optimization string   `yaml:"optimization"`
}

// AWSConfig represents the AWS configuration section
type AWSConfig struct {
	InstanceTypes map[string]string `yaml:"instance_types"`
	Storage       StorageConfig     `yaml:"storage"`
	Network       NetworkConfig     `yaml:"network"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	Type       string `yaml:"type"`
	SizeGB     int    `yaml:"size_gb"`
	IOPS       int    `yaml:"iops"`
	Throughput int    `yaml:"throughput"`
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	PlacementGroup     bool `yaml:"placement_group"`
	EnhancedNetworking bool `yaml:"enhanced_networking"`
	EFAEnabled         bool `yaml:"efa_enabled"`
}

// WorkflowConfig represents workflow configuration
type WorkflowConfig struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	Script         string `yaml:"script"`
	InputData      string `yaml:"input_data"`
	ExpectedOutput string `yaml:"expected_output"`
}

// Documentation represents documentation links
type Documentation struct {
	GettingStarted string `yaml:"getting_started"`
	Tutorials      string `yaml:"tutorials"`
	BestPractices  string `yaml:"best_practices"`
}

// NewDomainPackLoader creates a new domain pack loader
func NewDomainPackLoader() DomainPackLoaderInterface {
	// Try to find domain-packs directory
	domainPacksPath := "domain-packs"

	// Check if we're in the go subdirectory
	if _, err := os.Stat("../domain-packs"); err == nil {
		domainPacksPath = "../domain-packs"
	}

	// Check if we're in a different location (e.g., tests from internal/intelligence)
	if _, err := os.Stat("../../domain-packs"); err == nil {
		domainPacksPath = "../../domain-packs"
	}
	
	// Check if we're even deeper (e.g., tests from internal/intelligence)
	if _, err := os.Stat("../../../domain-packs"); err == nil {
		domainPacksPath = "../../../domain-packs"
	}

	return &DomainPackLoader{
		domainPacksPath: domainPacksPath,
		cache:           make(map[string]*DomainPackInfo),
	}
}

// LoadDomainPack loads a domain pack configuration by name
func (dpl *DomainPackLoader) LoadDomainPack(domainName string) (*DomainPackInfo, error) {
	// Check cache first
	if cached, exists := dpl.cache[domainName]; exists {
		return cached, nil
	}

	// Find domain pack file
	configPath, err := dpl.findDomainPackConfig(domainName)
	if err != nil {
		return nil, fmt.Errorf("domain pack not found: %s", domainName)
	}

	// Load configuration
	config, err := dpl.loadConfigFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load domain pack config: %w", err)
	}

	// Convert to DomainPackInfo
	info := dpl.convertToDomainPackInfo(config)

	// Cache result
	dpl.cache[domainName] = info

	return info, nil
}

// LoadAllDomainPacks loads all available domain packs
func (dpl *DomainPackLoader) LoadAllDomainPacks() (map[string]*DomainPackInfo, error) {
	domainPacks := make(map[string]*DomainPackInfo)

	domainsPath := filepath.Join(dpl.domainPacksPath, "domains")

	// Walk through category directories
	categories, err := os.ReadDir(domainsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read domains directory: %w", err)
	}

	for _, category := range categories {
		if !category.IsDir() {
			continue
		}

		categoryPath := filepath.Join(domainsPath, category.Name())
		domains, err := os.ReadDir(categoryPath)
		if err != nil {
			continue
		}

		for _, domain := range domains {
			if !domain.IsDir() {
				continue
			}

			configPath := filepath.Join(categoryPath, domain.Name(), "domain-pack.yaml")
			if _, err := os.Stat(configPath); err != nil {
				continue
			}

			config, err := dpl.loadConfigFile(configPath)
			if err != nil {
				continue
			}

			info := dpl.convertToDomainPackInfo(config)
			domainPacks[config.Name] = info
		}
	}

	return domainPacks, nil
}

// findDomainPackConfig locates the configuration file for a domain pack
func (dpl *DomainPackLoader) findDomainPackConfig(domainName string) (string, error) {
	domainsPath := filepath.Join(dpl.domainPacksPath, "domains")

	// Search through category directories
	categories, err := os.ReadDir(domainsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read domains directory: %w", err)
	}

	for _, category := range categories {
		if !category.IsDir() {
			continue
		}

		// Try exact match first
		configPath := filepath.Join(domainsPath, category.Name(), domainName, "domain-pack.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		// Try searching within category
		categoryPath := filepath.Join(domainsPath, category.Name())
		domains, err := os.ReadDir(categoryPath)
		if err != nil {
			continue
		}

		for _, domain := range domains {
			if !domain.IsDir() {
				continue
			}

			configPath := filepath.Join(categoryPath, domain.Name(), "domain-pack.yaml")
			if _, err := os.Stat(configPath); err != nil {
				continue
			}

			// Load config to check name
			config, err := dpl.loadConfigFile(configPath)
			if err != nil {
				continue
			}

			if config.Name == domainName {
				return configPath, nil
			}
		}
	}

	return "", fmt.Errorf("domain pack not found: %s", domainName)
}

// loadConfigFile loads and parses a domain pack configuration file
func (dpl *DomainPackLoader) loadConfigFile(configPath string) (*DomainPackConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config DomainPackConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return &config, nil
}

// convertToDomainPackInfo converts a DomainPackConfig to DomainPackInfo
func (dpl *DomainPackLoader) convertToDomainPackInfo(config *DomainPackConfig) *DomainPackInfo {
	info := &DomainPackInfo{
		Name:          config.Name,
		Version:       config.Version,
		Description:   config.Description,
		Categories:    config.Categories,
		InstanceTypes: config.AWSConfig.InstanceTypes,
		SpackPackages: config.SpackConfig.Packages,
		EstimatedCost: config.CostEstimates,
	}

	// Convert workflows
	for _, wf := range config.Workflows {
		info.Workflows = append(info.Workflows, WorkflowInfo{
			Name:        wf.Name,
			Description: wf.Description,
			InputData:   wf.InputData,
			OutputData:  wf.ExpectedOutput,
		})
	}

	return info
}

// GetAvailableDomains returns a list of all available domain names
func (dpl *DomainPackLoader) GetAvailableDomains() ([]string, error) {
	domainPacks, err := dpl.LoadAllDomainPacks()
	if err != nil {
		return nil, err
	}

	var domains []string
	for name := range domainPacks {
		domains = append(domains, name)
	}

	return domains, nil
}

// ValidateDomainPack validates a domain pack configuration
func (dpl *DomainPackLoader) ValidateDomainPack(domainName string) error {
	_, err := dpl.LoadDomainPack(domainName)
	return err
}

// ClearCache clears the domain pack cache
func (dpl *DomainPackLoader) ClearCache() {
	dpl.cache = make(map[string]*DomainPackInfo)
}
