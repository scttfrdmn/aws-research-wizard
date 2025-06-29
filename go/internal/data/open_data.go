package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenDataRegistry represents the AWS Open Data Registry
type OpenDataRegistry struct {
	client     *http.Client
	s3Manager  *S3Manager
	registryURL string
}

// DatasetInfo represents information about an AWS Open Data dataset
type DatasetInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`
	Domain      string            `json:"domain"`
	License     string            `json:"license"`
	Resources   []DatasetResource `json:"resources"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Size        string            `json:"size,omitempty"`
	Region      string            `json:"region,omitempty"`
}

// DatasetResource represents a resource within a dataset
type DatasetResource struct {
	Type        string `json:"type"`
	ARN         string `json:"arn,omitempty"`
	Region      string `json:"region,omitempty"`
	Bucket      string `json:"bucket,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Description string `json:"description,omitempty"`
}

// SearchParams for filtering datasets
type SearchParams struct {
	Domain   string
	Tags     []string
	Keywords string
	Limit    int
}

// NewOpenDataRegistry creates a new AWS Open Data Registry client
func NewOpenDataRegistry(s3Manager *S3Manager) *OpenDataRegistry {
	return &OpenDataRegistry{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		s3Manager:   s3Manager,
		registryURL: "https://registry.opendata.aws",
	}
}

// SearchDatasets searches for datasets in the AWS Open Data Registry
func (odr *OpenDataRegistry) SearchDatasets(ctx context.Context, params SearchParams) ([]DatasetInfo, error) {
	// Build URL with search parameters
	url := odr.registryURL + "/api/datasets"
	queryParams := make([]string, 0)
	
	if params.Domain != "" {
		queryParams = append(queryParams, fmt.Sprintf("domain=%s", params.Domain))
	}
	
	if params.Keywords != "" {
		queryParams = append(queryParams, fmt.Sprintf("q=%s", params.Keywords))
	}
	
	if params.Limit > 0 {
		queryParams = append(queryParams, fmt.Sprintf("limit=%d", params.Limit))
	}
	
	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := odr.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch datasets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Datasets []DatasetInfo `json:"datasets"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Filter by tags if specified
	if len(params.Tags) > 0 {
		filtered := make([]DatasetInfo, 0)
		for _, dataset := range response.Datasets {
			if odr.hasAnyTag(dataset.Tags, params.Tags) {
				filtered = append(filtered, dataset)
			}
		}
		return filtered, nil
	}

	return response.Datasets, nil
}

// GetDataset retrieves detailed information about a specific dataset
func (odr *OpenDataRegistry) GetDataset(ctx context.Context, datasetID string) (*DatasetInfo, error) {
	url := fmt.Sprintf("%s/api/datasets/%s", odr.registryURL, datasetID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := odr.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dataset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("dataset not found: %s", datasetID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var dataset DatasetInfo
	if err := json.Unmarshal(body, &dataset); err != nil {
		return nil, fmt.Errorf("failed to parse dataset: %w", err)
	}

	return &dataset, nil
}

// ListDomains returns available research domains in the registry
func (odr *OpenDataRegistry) ListDomains(ctx context.Context) ([]string, error) {
	url := odr.registryURL + "/api/domains"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := odr.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domains: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Domains []string `json:"domains"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse domains: %w", err)
	}

	return response.Domains, nil
}

// GetDatasetSize estimates the size of a dataset by sampling S3 objects
func (odr *OpenDataRegistry) GetDatasetSize(ctx context.Context, dataset *DatasetInfo) (int64, error) {
	var totalSize int64

	for _, resource := range dataset.Resources {
		if resource.Type == "s3" && resource.Bucket != "" {
			size, err := odr.estimateS3Size(ctx, resource.Bucket, resource.Prefix)
			if err != nil {
				return 0, fmt.Errorf("failed to estimate size for %s/%s: %w", resource.Bucket, resource.Prefix, err)
			}
			totalSize += size
		}
	}

	return totalSize, nil
}

// CreateDatasetAccess configures access to a dataset for optimal cost-free usage
func (odr *OpenDataRegistry) CreateDatasetAccess(ctx context.Context, dataset *DatasetInfo, localPath string) (*DatasetAccessConfig, error) {
	config := &DatasetAccessConfig{
		DatasetID:   dataset.ID,
		Name:        dataset.Name,
		LocalPath:   localPath,
		Resources:   make([]AccessResource, 0, len(dataset.Resources)),
		CreatedAt:   time.Now(),
	}

	for _, resource := range dataset.Resources {
		accessResource := AccessResource{
			Type:        resource.Type,
			ARN:         resource.ARN,
			Region:      resource.Region,
			Bucket:      resource.Bucket,
			Prefix:      resource.Prefix,
			Description: resource.Description,
		}

		// Configure optimal access patterns based on resource type
		if resource.Type == "s3" {
			accessResource.AccessPattern = "streaming" // Default to streaming for cost efficiency
			accessResource.CachePolicy = "metadata"    // Cache metadata only
		}

		config.Resources = append(config.Resources, accessResource)
	}

	return config, nil
}

// DatasetAccessConfig defines how to access a dataset optimally
type DatasetAccessConfig struct {
	DatasetID   string           `json:"dataset_id"`
	Name        string           `json:"name"`
	LocalPath   string           `json:"local_path"`
	Resources   []AccessResource `json:"resources"`
	CreatedAt   time.Time        `json:"created_at"`
}

// AccessResource defines how to access a specific resource
type AccessResource struct {
	Type          string `json:"type"`
	ARN           string `json:"arn,omitempty"`
	Region        string `json:"region,omitempty"`
	Bucket        string `json:"bucket,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	Description   string `json:"description,omitempty"`
	AccessPattern string `json:"access_pattern"` // "streaming", "batch", "random"
	CachePolicy   string `json:"cache_policy"`   // "none", "metadata", "partial", "full"
}

// DownloadDatasetSample downloads a sample of the dataset for exploration
func (odr *OpenDataRegistry) DownloadDatasetSample(ctx context.Context, dataset *DatasetInfo, localPath string, maxFiles int) error {
	for _, resource := range dataset.Resources {
		if resource.Type == "s3" && resource.Bucket != "" {
			// List a few sample files
			objects, err := odr.s3Manager.ListObjects(ctx, resource.Bucket, resource.Prefix, int32(maxFiles))
			if err != nil {
				return fmt.Errorf("failed to list sample objects: %w", err)
			}

			// Download up to maxFiles
			for i, obj := range objects {
				if i >= maxFiles {
					break
				}

				localFile := fmt.Sprintf("%s/%s", localPath, *obj.Key)
				err := odr.s3Manager.DownloadFile(ctx, resource.Bucket, *obj.Key, localFile, nil)
				if err != nil {
					return fmt.Errorf("failed to download sample file %s: %w", *obj.Key, err)
				}
			}
		}
	}

	return nil
}

// estimateS3Size estimates the total size of objects in an S3 bucket/prefix
func (odr *OpenDataRegistry) estimateS3Size(ctx context.Context, bucket, prefix string) (int64, error) {
	// List first 1000 objects to get size estimate
	objects, err := odr.s3Manager.ListObjects(ctx, bucket, prefix, 1000)
	if err != nil {
		return 0, err
	}

	var totalSize int64
	for _, obj := range objects {
		if obj.Size != nil {
			totalSize += *obj.Size
		}
	}

	return totalSize, nil
}

// hasAnyTag checks if dataset has any of the specified tags
func (odr *OpenDataRegistry) hasAnyTag(datasetTags, searchTags []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range datasetTags {
		tagSet[strings.ToLower(tag)] = true
	}

	for _, searchTag := range searchTags {
		if tagSet[strings.ToLower(searchTag)] {
			return true
		}
	}

	return false
}

// GetRecommendedDatasets returns datasets recommended for a specific research domain
func (odr *OpenDataRegistry) GetRecommendedDatasets(ctx context.Context, domain string) ([]DatasetInfo, error) {
	// Define domain-specific tags and keywords
	domainConfig := map[string]SearchParams{
		"genomics": {
			Domain:   "genomics",
			Tags:     []string{"genomics", "bioinformatics", "dna", "rna", "sequencing"},
			Keywords: "genome sequence dna rna",
			Limit:    10,
		},
		"climate": {
			Domain:   "climate",
			Tags:     []string{"climate", "weather", "meteorology", "atmosphere"},
			Keywords: "climate weather temperature precipitation",
			Limit:    10,
		},
		"machine-learning": {
			Domain:   "machine learning",
			Tags:     []string{"machine learning", "ai", "neural networks", "deep learning"},
			Keywords: "machine learning artificial intelligence dataset",
			Limit:    10,
		},
		"astronomy": {
			Domain:   "astronomy",
			Tags:     []string{"astronomy", "space", "cosmic", "stellar"},
			Keywords: "astronomy space cosmic stellar",
			Limit:    10,
		},
	}

	params, exists := domainConfig[strings.ToLower(domain)]
	if !exists {
		// Fallback to generic search
		params = SearchParams{
			Keywords: domain,
			Limit:    10,
		}
	}

	return odr.SearchDatasets(ctx, params)
}