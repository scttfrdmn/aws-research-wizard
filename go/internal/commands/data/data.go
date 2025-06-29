package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	awsClient "github.com/aws-research-wizard/go/internal/aws"
	"github.com/aws-research-wizard/go/internal/data"
)

// DataCmd represents the data command
var DataCmd = &cobra.Command{
	Use:   "data",
	Short: "Manage research data transfers and AWS Open Data integration",
	Long: `The data command provides comprehensive data management capabilities including:
- S3 transfer optimization with multi-part uploads
- AWS Open Data registry integration and discovery  
- Data pipeline orchestration and monitoring
- Transfer progress tracking and analytics`,
	Example: `  # Search for genomics datasets
  aws-research-wizard data search --domain genomics

  # Download a dataset
  aws-research-wizard data download s3://bucket/dataset ./local/path

  # Create a data pipeline
  aws-research-wizard data pipeline create "My Pipeline" --description "Data processing workflow"

  # Monitor active transfers
  aws-research-wizard data monitor`,
}

var (
	// Global data management components
	s3Manager       *data.S3Manager
	openDataRegistry *data.OpenDataRegistry
	pipelineManager *data.PipelineManager
	transferMonitor *data.TransferMonitor
)

func init() {
	// Add subcommands
	DataCmd.AddCommand(searchCmd)
	DataCmd.AddCommand(downloadCmd)
	DataCmd.AddCommand(uploadCmd)
	DataCmd.AddCommand(pipelineCmd)
	DataCmd.AddCommand(monitorCmd)
	DataCmd.AddCommand(listCmd)
	DataCmd.AddCommand(infoCmd)

	// Global flags
	DataCmd.PersistentFlags().String("region", "us-east-1", "AWS region")
	DataCmd.PersistentFlags().String("config-path", "", "Configuration path for pipelines and settings")
	DataCmd.PersistentFlags().Int("concurrency", 10, "Number of concurrent transfers")
	DataCmd.PersistentFlags().String("part-size", "16MB", "Part size for multipart uploads (e.g., 16MB, 64MB)")
}

// Search command
var searchCmd = &cobra.Command{
	Use:   "search [keywords]",
	Short: "Search AWS Open Data registry",
	Long:  "Search for datasets in the AWS Open Data registry by domain, tags, or keywords",
	Example: `  # Search by domain
  aws-research-wizard data search --domain genomics

  # Search by keywords
  aws-research-wizard data search "climate temperature"

  # Search with tags
  aws-research-wizard data search --tags "genomics,dna"`,
	RunE: runSearch,
}

var downloadCmd = &cobra.Command{
	Use:   "download [s3-uri] [local-path]",
	Short: "Download data with optimized S3 transfers",
	Long:  "Download files or datasets from S3 with multi-part download optimization and progress tracking",
	Example: `  # Download a single file
  aws-research-wizard data download s3://bucket/file.txt ./local/file.txt

  # Download with custom settings
  aws-research-wizard data download s3://bucket/dataset ./data --concurrency 20 --part-size 64MB`,
	Args: cobra.ExactArgs(2),
	RunE: runDownload,
}

var uploadCmd = &cobra.Command{
	Use:   "upload [local-path] [s3-uri]",
	Short: "Upload data with optimized S3 transfers",
	Long:  "Upload files to S3 with multi-part upload optimization and progress tracking",
	Example: `  # Upload a single file
  aws-research-wizard data upload ./local/file.txt s3://bucket/file.txt

  # Upload with custom settings
  aws-research-wizard data upload ./dataset s3://bucket/dataset --concurrency 20 --part-size 64MB`,
	Args: cobra.ExactArgs(2),
	RunE: runUpload,
}

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Manage data processing pipelines",
	Long:  "Create, execute, and manage data processing pipelines for complex workflows",
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor active transfers and pipelines",
	Long:  "Launch an interactive dashboard to monitor active transfers and pipeline execution",
	RunE:  runMonitor,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available datasets and resources",
	Long:  "List datasets from AWS Open Data registry or show recommended datasets for research domains",
	Example: `  # List recommended datasets for genomics
  aws-research-wizard data list --domain genomics

  # List all available domains
  aws-research-wizard data list --domains`,
	RunE: runList,
}

var infoCmd = &cobra.Command{
	Use:   "info [dataset-id]",
	Short: "Get detailed information about a dataset",
	Long:  "Retrieve detailed information about a specific dataset including size, resources, and access patterns",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	// Search command flags
	searchCmd.Flags().String("domain", "", "Filter by research domain (e.g., genomics, climate)")
	searchCmd.Flags().StringSlice("tags", []string{}, "Filter by tags")
	searchCmd.Flags().Int("limit", 20, "Maximum number of results")

	// Download command flags
	downloadCmd.Flags().Bool("sample", false, "Download only a sample of the dataset")
	downloadCmd.Flags().Int("max-files", 10, "Maximum files to download (for sample)")

	// Upload command flags
	uploadCmd.Flags().String("storage-class", "STANDARD", "S3 storage class (STANDARD, IA, GLACIER)")

	// List command flags
	listCmd.Flags().String("domain", "", "Show recommended datasets for domain")
	listCmd.Flags().Bool("domains", false, "List all available domains")

	// Monitor command flags
	monitorCmd.Flags().Duration("refresh", 1*time.Second, "Refresh interval for monitoring")

	// Pipeline subcommands
	pipelineCreateCmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new data pipeline",
		Args:  cobra.ExactArgs(1),
		RunE:  runPipelineCreate,
	}
	pipelineCreateCmd.Flags().String("description", "", "Pipeline description")
	pipelineCreateCmd.Flags().Int("max-jobs", 10, "Maximum concurrent jobs")

	pipelineListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all pipelines",
		RunE:  runPipelineList,
	}

	pipelineExecuteCmd := &cobra.Command{
		Use:   "execute [pipeline-id]",
		Short: "Execute a pipeline",
		Args:  cobra.ExactArgs(1),
		RunE:  runPipelineExecute,
	}

	pipelineCmd.AddCommand(pipelineCreateCmd, pipelineListCmd, pipelineExecuteCmd)
}

// Initialize data management components
func initializeDataComponents(cmd *cobra.Command) error {
	region, _ := cmd.Flags().GetString("region")
	configPath, _ := cmd.Flags().GetString("config-path")
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	partSizeStr, _ := cmd.Flags().GetString("part-size")

	// Parse part size
	partSize, err := parseSize(partSizeStr)
	if err != nil {
		return fmt.Errorf("invalid part-size: %w", err)
	}

	// Create AWS client
	ctx := context.Background()
	client, err := awsClient.NewClient(ctx, region)
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	// Initialize S3 manager
	s3Config := &data.S3ManagerConfig{
		PartSize:    partSize,
		Concurrency: concurrency,
		MaxRetries:  3,
	}
	s3Manager = data.NewS3Manager(client.S3, region, s3Config)

	// Initialize Open Data registry
	openDataRegistry = data.NewOpenDataRegistry(s3Manager)

	// Initialize pipeline manager
	if configPath == "" {
		configPath = filepath.Join(os.Getenv("HOME"), ".aws-research-wizard")
	}
	pipelineManager = data.NewPipelineManager(s3Manager, openDataRegistry, configPath)

	// Initialize transfer monitor
	transferMonitor = data.NewTransferMonitor(s3Manager, pipelineManager)

	return nil
}

func runSearch(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	domain, _ := cmd.Flags().GetString("domain")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	limit, _ := cmd.Flags().GetInt("limit")

	keywords := ""
	if len(args) > 0 {
		keywords = strings.Join(args, " ")
	}

	params := data.SearchParams{
		Domain:   domain,
		Tags:     tags,
		Keywords: keywords,
		Limit:    limit,
	}

	ctx := context.Background()
	datasets, err := openDataRegistry.SearchDatasets(ctx, params)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(datasets) == 0 {
		fmt.Println("No datasets found matching the criteria")
		return nil
	}

	fmt.Printf("Found %d datasets:\n\n", len(datasets))
	for _, dataset := range datasets {
		fmt.Printf("ID: %s\n", dataset.ID)
		fmt.Printf("Name: %s\n", dataset.Name)
		fmt.Printf("Domain: %s\n", dataset.Domain)
		fmt.Printf("Description: %s\n", dataset.Description)
		if len(dataset.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(dataset.Tags, ", "))
		}
		fmt.Printf("Resources: %d\n", len(dataset.Resources))
		fmt.Println(strings.Repeat("-", 50))
	}

	return nil
}

func runDownload(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	s3URI := args[0]
	localPath := args[1]
	sample, _ := cmd.Flags().GetBool("sample")
	maxFiles, _ := cmd.Flags().GetInt("max-files")

	// Parse S3 URI (s3://bucket/key)
	bucket, key, err := parseS3URI(s3URI)
	if err != nil {
		return fmt.Errorf("invalid S3 URI: %w", err)
	}

	ctx := context.Background()

	// Progress callback
	progressCallback := func(progress data.TransferProgress) {
		fmt.Printf("\rProgress: %.1f%% (%s/s) ETA: %s", 
			progress.Percentage, 
			formatBytes(progress.Speed), 
			progress.ETA.Round(time.Second))
	}

	if sample {
		// Download sample files
		objects, err := s3Manager.ListObjects(ctx, bucket, key, int32(maxFiles))
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		fmt.Printf("Downloading %d sample files...\n", len(objects))
		for _, obj := range objects {
			localFile := filepath.Join(localPath, *obj.Key)
			err := s3Manager.DownloadFile(ctx, bucket, *obj.Key, localFile, progressCallback)
			if err != nil {
				return fmt.Errorf("failed to download %s: %w", *obj.Key, err)
			}
			fmt.Printf("\nCompleted: %s\n", *obj.Key)
		}
	} else {
		// Download single file
		err := s3Manager.DownloadFile(ctx, bucket, key, localPath, progressCallback)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		fmt.Printf("\nDownload completed: %s\n", localPath)
	}

	return nil
}

func runUpload(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	localPath := args[0]
	s3URI := args[1]

	// Parse S3 URI
	bucket, key, err := parseS3URI(s3URI)
	if err != nil {
		return fmt.Errorf("invalid S3 URI: %w", err)
	}

	ctx := context.Background()

	// Progress callback
	progressCallback := func(progress data.TransferProgress) {
		fmt.Printf("\rProgress: %.1f%% (%s/s) ETA: %s", 
			progress.Percentage, 
			formatBytes(progress.Speed), 
			progress.ETA.Round(time.Second))
	}

	err = s3Manager.UploadFile(ctx, bucket, key, localPath, progressCallback)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("\nUpload completed: %s\n", s3URI)
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	domain, _ := cmd.Flags().GetString("domain")
	showDomains, _ := cmd.Flags().GetBool("domains")

	ctx := context.Background()

	if showDomains {
		domains, err := openDataRegistry.ListDomains(ctx)
		if err != nil {
			return fmt.Errorf("failed to list domains: %w", err)
		}

		fmt.Println("Available domains:")
		for _, d := range domains {
			fmt.Printf("  - %s\n", d)
		}
		return nil
	}

	if domain != "" {
		datasets, err := openDataRegistry.GetRecommendedDatasets(ctx, domain)
		if err != nil {
			return fmt.Errorf("failed to get recommended datasets: %w", err)
		}

		fmt.Printf("Recommended datasets for %s:\n\n", domain)
		for _, dataset := range datasets {
			fmt.Printf("• %s\n  %s\n\n", dataset.Name, dataset.Description)
		}
		return nil
	}

	return fmt.Errorf("specify --domain or --domains flag")
}

func runInfo(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	datasetID := args[0]
	ctx := context.Background()

	dataset, err := openDataRegistry.GetDataset(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to get dataset info: %w", err)
	}

	fmt.Printf("Dataset: %s\n", dataset.Name)
	fmt.Printf("ID: %s\n", dataset.ID)
	fmt.Printf("Domain: %s\n", dataset.Domain)
	fmt.Printf("Description: %s\n", dataset.Description)
	fmt.Printf("License: %s\n", dataset.License)
	
	if len(dataset.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(dataset.Tags, ", "))
	}

	fmt.Printf("\nResources (%d):\n", len(dataset.Resources))
	for i, resource := range dataset.Resources {
		fmt.Printf("  %d. Type: %s\n", i+1, resource.Type)
		if resource.Bucket != "" {
			fmt.Printf("     Location: s3://%s/%s\n", resource.Bucket, resource.Prefix)
		}
		if resource.Description != "" {
			fmt.Printf("     Description: %s\n", resource.Description)
		}
	}

	// Estimate size
	fmt.Println("\nEstimating dataset size...")
	size, err := openDataRegistry.GetDatasetSize(ctx, dataset)
	if err != nil {
		fmt.Printf("Warning: Could not estimate size: %v\n", err)
	} else {
		fmt.Printf("Estimated size: %s\n", formatBytes(size))
	}

	return nil
}

func runMonitor(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	// Start monitoring in background
	ctx := context.Background()
	go transferMonitor.StartMonitoring(ctx)

	// Create and run TUI
	// model := data.NewMonitorModel(transferMonitor)
	
	// This would typically use Bubble Tea's Program.Run()
	// For now, show a simple text-based monitor
	fmt.Println("Transfer Monitor (Press Ctrl+C to exit)")
	fmt.Println("=====================================")

	for {
		transferSummary := transferMonitor.GetTransferSummary()
		pipelineSummary := transferMonitor.GetPipelineSummary()

		fmt.Printf("\rTransfers - Active: %d | Completed: %d | Failed: %d | Speed: %s/s", 
			transferSummary.ActiveTransfers,
			transferSummary.CompletedTransfers, 
			transferSummary.FailedTransfers,
			formatBytes(transferSummary.AverageSpeed))

		fmt.Printf(" | Pipelines - Active: %d | Jobs: %d/%d", 
			pipelineSummary.ActivePipelines,
			pipelineSummary.CompletedJobs,
			pipelineSummary.TotalJobs)

		time.Sleep(1 * time.Second)
	}
}

func runPipelineCreate(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	name := args[0]
	description, _ := cmd.Flags().GetString("description")
	maxJobs, _ := cmd.Flags().GetInt("max-jobs")

	config := data.PipelineConfig{
		MaxConcurrentJobs: maxJobs,
		RetryAttempts:     3,
		RetryDelay:        5 * time.Second,
		Timeout:           30 * time.Minute,
		WorkingDirectory:  "/tmp",
	}

	pipeline := pipelineManager.CreatePipeline(name, description, config)
	
	fmt.Printf("Created pipeline: %s (ID: %s)\n", pipeline.Name, pipeline.ID)
	return nil
}

func runPipelineList(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	pipelines := pipelineManager.ListPipelines()
	
	if len(pipelines) == 0 {
		fmt.Println("No pipelines found")
		return nil
	}

	fmt.Printf("Found %d pipelines:\n\n", len(pipelines))
	for _, pipeline := range pipelines {
		fmt.Printf("ID: %s\n", pipeline.ID)
		fmt.Printf("Name: %s\n", pipeline.Name)
		fmt.Printf("Status: %s\n", pipeline.Status)
		fmt.Printf("Jobs: %d\n", len(pipeline.Jobs))
		fmt.Printf("Created: %s\n", pipeline.CreatedAt.Format(time.RFC3339))
		fmt.Println(strings.Repeat("-", 40))
	}

	return nil
}

func runPipelineExecute(cmd *cobra.Command, args []string) error {
	if err := initializeDataComponents(cmd); err != nil {
		return err
	}

	pipelineID := args[0]
	ctx := context.Background()

	err := pipelineManager.ExecutePipeline(ctx, pipelineID)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	fmt.Printf("Pipeline execution started: %s\n", pipelineID)
	return nil
}

// Utility functions

func parseS3URI(uri string) (bucket, key string, err error) {
	if !strings.HasPrefix(uri, "s3://") {
		return "", "", fmt.Errorf("URI must start with s3://")
	}

	path := strings.TrimPrefix(uri, "s3://")
	parts := strings.SplitN(path, "/", 2)
	
	if len(parts) < 1 {
		return "", "", fmt.Errorf("invalid S3 URI format")
	}

	bucket = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}

	return bucket, key, nil
}

func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))
	
	multiplier := int64(1)
	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return size * multiplier, nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}