package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/aws-research-wizard/go/internal/aws"
	"github.com/scttfrdmn/aws-research-wizard/go/internal/config"
)

var (
	configRoot   string
	region       string
	stackName    string
	domainName   string
	instanceType string
	dryRun       bool
	timeout      time.Duration
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "aws-research-wizard-deploy",
		Short: "AWS Research Wizard - Infrastructure Deployment Tool",
		Long: `Deploy optimized AWS research environments with pre-configured
domain packs, monitoring, and cost optimization.

This deployment tool provides:
- Terraform infrastructure management
- EC2 instance provisioning
- Security group configuration
- Monitoring setup
- Cost tracking`,
		Run: runInteractiveDeploy,
	}

	// Add flags
	rootCmd.PersistentFlags().StringVar(&configRoot, "config", "", "Configuration root directory")
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&stackName, "stack", "", "Terraform deployment name")
	rootCmd.PersistentFlags().StringVar(&domainName, "domain", "", "Research domain name")
	rootCmd.PersistentFlags().StringVar(&instanceType, "instance", "", "EC2 instance type")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show deployment plan without executing")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Minute, "Deployment timeout")

	// Add subcommands
	rootCmd.AddCommand(
		createDeployCommand(),
		createStatusCommand(),
		createDeleteCommand(),
		createListCommand(),
		createValidateCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runInteractiveDeploy(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Find config root if not specified
	if configRoot == "" {
		configRoot = findConfigRoot()
	}

	fmt.Printf("🚀 AWS Research Wizard - Infrastructure Deployment\n")
	fmt.Printf("Config Root: %s\n", configRoot)
	fmt.Printf("AWS Region: %s\n\n", region)

	// Initialize AWS client
	awsClient, err := aws.NewClient(ctx, region)
	if err != nil {
		log.Fatalf("Failed to initialize AWS client: %v", err)
	}

	// Validate AWS credentials
	if err := awsClient.ValidateCredentials(ctx); err != nil {
		log.Fatalf("AWS credentials validation failed: %v", err)
	}

	fmt.Printf("✅ AWS credentials validated\n\n")

	// Load domain configuration if specified
	if domainName != "" {
		if err := deployDomain(ctx, awsClient); err != nil {
			log.Fatalf("Deployment failed: %v", err)
		}
	} else {
		fmt.Println("Please specify a domain with --domain flag or use subcommands:")
		fmt.Println("  deploy --domain genomics --instance r6i.4xlarge")
		fmt.Println("  status --stack my-research-stack")
		fmt.Println("  list")
	}
}

func deployDomain(ctx context.Context, awsClient *aws.Client) error {
	// Load domain configuration
	loader := config.NewConfigLoader(configRoot)
	domains, err := loader.LoadAllDomains()
	if err != nil {
		return fmt.Errorf("failed to load domains: %w", err)
	}

	domain, exists := domains[domainName]
	if !exists {
		return fmt.Errorf("domain '%s' not found", domainName)
	}

	fmt.Printf("📋 Deploying Domain: %s\n", domain.Name)
	fmt.Printf("Description: %s\n", domain.Description)

	// Select instance type
	selectedInstance := instanceType
	if selectedInstance == "" {
		// Use the first recommended instance if not specified
		for _, rec := range domain.AWSInstanceRecommendations {
			selectedInstance = rec.InstanceType
			break
		}
	}

	if selectedInstance == "" {
		return fmt.Errorf("no instance type specified or available in domain recommendations")
	}

	fmt.Printf("Instance Type: %s\n", selectedInstance)

	// Generate stack name if not provided
	if stackName == "" {
		stackName = fmt.Sprintf("research-wizard-%s", domainName)
	}

	fmt.Printf("Stack Name: %s\n\n", stackName)

	if dryRun {
		fmt.Printf("🔍 DRY RUN - Deployment plan:\n")
		fmt.Printf("  1. Initialize Terraform workspace: %s\n", stackName)
		fmt.Printf("  2. Launch EC2 instance: %s\n", selectedInstance)
		fmt.Printf("  3. Configure security groups\n")
		fmt.Printf("  4. Set up monitoring and alarms\n")
		fmt.Printf("  5. Configure cost tracking\n")
		fmt.Printf("\nTo execute, run without --dry-run flag\n")
		return nil
	}

	// Create infrastructure manager
	infraManager := aws.NewTerraformManager(awsClient, "")

	// Template generation no longer needed for Terraform
	// CloudFormation template generation removed

	// Create stack parameters
	parameters := map[string]string{
		"InstanceType": selectedInstance,
		"DomainName":   domainName,
		"KeyName":      "", // User should specify key pair
	}

	fmt.Printf("🏗️ Deploying Terraform infrastructure...\n")

	// Initialize Terraform
	if err := infraManager.InitTerraform(ctx); err != nil {
		return fmt.Errorf("failed to initialize Terraform: %w", err)
	}

	// Apply deployment
	deploymentInfo, err := infraManager.ApplyDeployment(ctx, parameters)
	if err != nil {
		return fmt.Errorf("failed to apply deployment: %w", err)
	}

	fmt.Printf("✅ Deployment completed successfully!\n")
	finalStackInfo := deploymentInfo

	fmt.Printf("🎉 Deployment completed successfully!\n\n")
	fmt.Printf("Deployment Details:\n")
	fmt.Printf("  Name: %s\n", finalStackInfo.WorkspaceName)
	fmt.Printf("  Status: %s\n", finalStackInfo.Status)
	fmt.Printf("  Created: %s\n", finalStackInfo.CreatedTime.Format(time.RFC3339))

	if len(finalStackInfo.Outputs) > 0 {
		fmt.Printf("\nDeployment Outputs:\n")
		for key, value := range finalStackInfo.Outputs {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	fmt.Printf("\n📊 Next Steps:\n")
	fmt.Printf("  1. Monitor with: aws-research-wizard-monitor --stack %s\n", stackName)
	fmt.Printf("  2. Check costs: aws-research-wizard-deploy status --stack %s\n", stackName)
	fmt.Printf("  3. SSH to instance using outputs above\n")

	return nil
}

// generateCloudFormationTemplate is deprecated - replaced with Terraform
// This function is kept for reference only and should not be used
func generateCloudFormationTemplate(domain *config.DomainPack, instanceType string) (string, error) {
	return "", fmt.Errorf("CloudFormation templates are deprecated - please use Terraform infrastructure in terraform/environments/aws/")
}

func createDeployCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a research environment",
		Run: func(cmd *cobra.Command, args []string) {
			if domainName == "" {
				log.Fatal("Domain name is required. Use --domain flag.")
			}

			ctx := context.Background()
			if configRoot == "" {
				configRoot = findConfigRoot()
			}

			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			if err := deployDomain(ctx, awsClient); err != nil {
				log.Fatalf("Deployment failed: %v", err)
			}
		},
	}
}

func createStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check deployment status",
		Run: func(cmd *cobra.Command, args []string) {
			if stackName == "" {
				log.Fatal("Stack name is required. Use --stack flag.")
			}

			ctx := context.Background()
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			infraManager := aws.NewTerraformManager(awsClient, "")

			stackInfo, err := infraManager.GetDeploymentInfo(ctx)
			if err != nil {
				log.Fatalf("Failed to get stack info: %v", err)
			}

			fmt.Printf("📊 Stack Status: %s\n\n", stackName)
			fmt.Printf("Status: %s\n", stackInfo.Status)
			fmt.Printf("Created: %s\n", stackInfo.CreatedTime.Format(time.RFC3339))

			if stackInfo.UpdatedTime != nil {
				fmt.Printf("Updated: %s\n", stackInfo.UpdatedTime.Format(time.RFC3339))
			}

			if len(stackInfo.Outputs) > 0 {
				fmt.Printf("\nOutputs:\n")
				for key, value := range stackInfo.Outputs {
					fmt.Printf("  %s: %s\n", key, value)
				}
			}
		},
	}
}

func createDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete a research environment",
		Run: func(cmd *cobra.Command, args []string) {
			if stackName == "" {
				log.Fatal("Stack name is required. Use --stack flag.")
			}

			ctx := context.Background()
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			infraManager := aws.NewTerraformManager(awsClient, "")

			fmt.Printf("⚠️  Deleting stack: %s\n", stackName)
			fmt.Printf("This action cannot be undone. Continue? (y/N): ")

			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" {
				fmt.Println("Deletion cancelled.")
				return
			}

			if err := infraManager.DestroyDeployment(ctx, map[string]string{}); err != nil {
				log.Fatalf("Failed to delete stack: %v", err)
			}

			fmt.Printf("🗑️  Stack deletion initiated. Monitor progress with: aws-research-wizard-deploy status --stack %s\n", stackName)
		},
	}
}

func createListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List deployed research environments",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			infraManager := aws.NewTerraformManager(awsClient, "")

			// List instances with research wizard tags
			filters := map[string][]string{
				"tag:CreatedBy":       {"AWS-Research-Wizard"},
				"instance-state-name": {"running", "pending", "stopping", "stopped"},
			}

			instances, err := infraManager.ListInstances(ctx, filters)
			if err != nil {
				log.Fatalf("Failed to list instances: %v", err)
			}

			fmt.Printf("🖥️  Research Environments (%d total):\n\n", len(instances))

			for _, instance := range instances {
				domain := instance.Tags["Domain"]
				if domain == "" {
					domain = "Unknown"
				}

				fmt.Printf("Instance: %s\n", instance.InstanceID)
				fmt.Printf("  Domain: %s\n", domain)
				fmt.Printf("  Type: %s\n", instance.InstanceType)
				fmt.Printf("  State: %s\n", instance.State)
				fmt.Printf("  Public IP: %s\n", instance.PublicIP)
				fmt.Printf("  Launch Time: %s\n", instance.LaunchTime.Format(time.RFC3339))
				fmt.Printf("\n")
			}
		},
	}
}

func createValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate deployment configuration",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			if configRoot == "" {
				configRoot = findConfigRoot()
			}

			fmt.Printf("🔍 Validating configuration...\n\n")

			// Validate AWS credentials
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			if err := awsClient.ValidateCredentials(ctx); err != nil {
				log.Fatalf("AWS credentials validation failed: %v", err)
			}

			fmt.Printf("✅ AWS credentials valid\n")

			// Validate domain configuration
			if domainName != "" {
				loader := config.NewConfigLoader(configRoot)
				domains, err := loader.LoadAllDomains()
				if err != nil {
					log.Fatalf("Failed to load domains: %v", err)
				}

				if _, exists := domains[domainName]; !exists {
					log.Fatalf("Domain '%s' not found", domainName)
				}

				fmt.Printf("✅ Domain configuration valid: %s\n", domainName)
			}

			// Validate region
			zones, err := awsClient.GetAvailabilityZones(ctx)
			if err != nil {
				log.Fatalf("Failed to validate region: %v", err)
			}

			fmt.Printf("✅ Region valid: %s (%d availability zones)\n", region, len(zones))

			fmt.Printf("\n🎉 All validations passed!\n")
		},
	}
}

func findConfigRoot() string {
	// Look for configs directory in current directory and parent directories
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory")
	}

	for {
		configsPath := filepath.Join(currentDir, "configs")
		if _, err := os.Stat(configsPath); err == nil {
			return currentDir
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break // Reached root directory
		}
		currentDir = parent
	}

	log.Fatal("Could not find configs directory. Please specify with --config flag.")
	return ""
}
