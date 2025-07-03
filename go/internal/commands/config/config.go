package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/aws-research-wizard/go/internal/aws"
	"github.com/scttfrdmn/aws-research-wizard/go/internal/config"
	"github.com/scttfrdmn/aws-research-wizard/go/internal/tui"
)

// NewConfigCommand creates the config subcommand
func NewConfigCommand() *cobra.Command {
	var configRoot string
	var simple bool

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Domain configuration and cost analysis",
		Long: `Configure research domains with interactive selection, cost estimation,
and instance type recommendations.

Available operations:
- List all available research domains
- Show detailed domain information
- Interactive domain selection with cost analysis
- Instance type recommendations and optimization`,
		Run: func(cmd *cobra.Command, args []string) {
			runInteractiveConfig(cmd, configRoot, simple)
		},
	}

	// Add flags
	configCmd.PersistentFlags().StringVar(&configRoot, "config", "", "Configuration root directory")
	configCmd.PersistentFlags().BoolVar(&simple, "simple", false, "Use simple TUI without advanced features")

	// Add subcommands
	configCmd.AddCommand(
		createListCommand(&configRoot),
		createInfoCommand(&configRoot),
		createCostCommand(&configRoot),
		createSearchCommand(&configRoot),
	)

	return configCmd
}

func runInteractiveConfig(cmd *cobra.Command, configRoot string, simple bool) {
	ctx := context.Background()

	// Find config root if not specified
	if configRoot == "" {
		configRoot = findConfigRoot()
	}

	region, _ := cmd.Flags().GetString("region")

	fmt.Printf("🔬 AWS Research Wizard - Domain Configuration\n")
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

	// Load domains
	loader := config.NewConfigLoader(configRoot)
	domains, err := loader.LoadAllDomains()
	if err != nil {
		log.Fatalf("Failed to load domains: %v", err)
	}

	fmt.Printf("Loaded %d research domains\n\n", len(domains))

	// Run domain selector
	selectedDomain, err := tui.RunDomainSelector(domains)
	if err != nil {
		log.Fatalf("Failed to run domain selector: %v", err)
	}

	if selectedDomain == nil {
		fmt.Println("No domain selected. Exiting.")
		return
	}

	fmt.Printf("\n✅ Selected Domain: %s\n", selectedDomain.Name)
	fmt.Printf("Description: %s\n\n", selectedDomain.Description)

	// Run cost calculator
	fmt.Println("📊 Calculating costs for recommended instances...")
	selectedInstance, estimate, err := tui.RunCostCalculator(selectedDomain, region)
	if err != nil {
		log.Fatalf("Failed to run cost calculator: %v", err)
	}

	if selectedInstance == "" {
		fmt.Println("No instance selected. Configuration complete.")
		return
	}

	// Display final configuration
	fmt.Printf("\n🎯 Final Configuration:\n")
	fmt.Printf("  Domain: %s\n", selectedDomain.Name)
	fmt.Printf("  Instance: %s\n", selectedInstance)
	fmt.Printf("  Cost: $%.3f/hour ($%.0f/month)\n", estimate.HourlyCost, estimate.MonthlyCost)
	fmt.Printf("  Specs: %d vCPUs, %s RAM\n", estimate.VCPUs, estimate.Memory)
	fmt.Printf("  Spot Savings: $%.0f/month (70%% discount)\n", estimate.SpotSavings*24*30.44)

	fmt.Printf("\n💡 Next Steps:\n")
	fmt.Printf("  1. Deploy with: aws-research-wizard deploy --domain %s --instance %s\n", selectedDomain.Name, selectedInstance)
	fmt.Printf("  2. Monitor with: aws-research-wizard monitor\n")
	fmt.Printf("  3. Get help with: aws-research-wizard --help\n")
}

func createListCommand(configRoot *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available research domains",
		Run: func(cmd *cobra.Command, args []string) {
			if *configRoot == "" {
				*configRoot = findConfigRoot()
			}

			loader := config.NewConfigLoader(*configRoot)
			domains, err := loader.LoadAllDomains()
			if err != nil {
				log.Fatalf("Failed to load domains: %v", err)
			}

			fmt.Printf("Available Research Domains (%d total):\n\n", len(domains))

			for name, domain := range domains {
				fmt.Printf("📚 %s\n", name)
				fmt.Printf("   %s\n", domain.Description)
				fmt.Printf("   Target Users: %v\n", domain.TargetUsers)
				fmt.Printf("   Monthly Cost: $%.0f\n\n", domain.EstimatedCost.Total)
			}
		},
	}
}

func createInfoCommand(configRoot *string) *cobra.Command {
	return &cobra.Command{
		Use:   "info [domain]",
		Short: "Show detailed information about a domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if *configRoot == "" {
				*configRoot = findConfigRoot()
			}

			domainName := args[0]
			loader := config.NewConfigLoader(*configRoot)
			domains, err := loader.LoadAllDomains()
			if err != nil {
				log.Fatalf("Failed to load domains: %v", err)
			}

			domain, exists := domains[domainName]
			if !exists {
				log.Fatalf("Domain '%s' not found", domainName)
			}

			fmt.Printf("🔬 Domain: %s\n\n", domain.Name)
			fmt.Printf("Description: %s\n\n", domain.Description)

			fmt.Printf("Target Users: %s\n", domain.TargetUsers)

			fmt.Printf("\nSpack Package Categories (%d):\n", len(domain.SpackPackages))
			for category, packages := range domain.SpackPackages {
				fmt.Printf("  • %s: %v\n", category, packages)
			}

			fmt.Printf("\nAWS Instance Recommendations:\n")
			for _, rec := range domain.AWSInstanceRecommendations {
				fmt.Printf("  • %s: %s (%d vCPUs, %d GB) - $%.3f/hour\n",
					rec.UseCase, rec.InstanceType, rec.VCPUs, rec.MemoryGB, rec.CostPerHour)
			}

			fmt.Printf("\nEstimated Costs:\n")
			fmt.Printf("  • Compute: $%.0f/month\n", domain.EstimatedCost.Compute)
			fmt.Printf("  • Storage: $%.0f/month\n", domain.EstimatedCost.Storage)
			fmt.Printf("  • Total: $%.0f/month\n", domain.EstimatedCost.Total)
		},
	}
}

func createCostCommand(configRoot *string) *cobra.Command {
	return &cobra.Command{
		Use:   "cost [domain]",
		Short: "Calculate costs for a specific domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if *configRoot == "" {
				*configRoot = findConfigRoot()
			}

			domainName := args[0]
			loader := config.NewConfigLoader(*configRoot)
			domains, err := loader.LoadAllDomains()
			if err != nil {
				log.Fatalf("Failed to load domains: %v", err)
			}

			domain, exists := domains[domainName]
			if !exists {
				log.Fatalf("Domain '%s' not found", domainName)
			}

			region, _ := cmd.Flags().GetString("region")
			fmt.Printf("💰 Cost Analysis: %s\n\n", domain.Name)

			_, _, err = tui.RunCostCalculator(domain, region)
			if err != nil {
				log.Fatalf("Failed to run cost calculator: %v", err)
			}
		},
	}
}

func createSearchCommand(configRoot *string) *cobra.Command {
	return &cobra.Command{
		Use:   "search [query]",
		Short: "Search domains by name or description",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if *configRoot == "" {
				*configRoot = findConfigRoot()
			}

			query := args[0]
			loader := config.NewConfigLoader(*configRoot)
			domains, err := loader.LoadAllDomains()
			if err != nil {
				log.Fatalf("Failed to load domains: %v", err)
			}

			fmt.Printf("🔍 Search results for '%s':\n\n", query)

			found := 0
			for name, domain := range domains {
				// Simple case-insensitive search in name and description
				nameMatch := containsIgnoreCase(name, query)
				descMatch := containsIgnoreCase(domain.Description, query)

				if nameMatch || descMatch {
					found++
					fmt.Printf("📚 %s\n", name)
					fmt.Printf("   %s\n", domain.Description)
					fmt.Printf("   Monthly Cost: $%.0f\n\n", domain.EstimatedCost.Total)
				}
			}

			if found == 0 {
				fmt.Printf("No domains found matching '%s'\n", query)
			} else {
				fmt.Printf("Found %d matching domains\n", found)
			}
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

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		len(substr) > 0 &&
		strings.ToLower(s) != strings.ToLower(s) ||
		strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
