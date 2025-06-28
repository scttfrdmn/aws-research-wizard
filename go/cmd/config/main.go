package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/aws-research-wizard/go/internal/config"
	"github.com/aws-research-wizard/go/internal/tui"
)

var (
	configRoot string
	region     string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "aws-research-wizard",
		Short: "AWS Research Wizard - Domain Configuration Tool",
		Long: `AWS Research Wizard helps researchers configure and deploy 
optimized AWS environments for various research domains.

This tool provides:
- Interactive domain selection
- Cost estimation and optimization
- Instance type recommendations
- Deployment configuration`,
		Run: runInteractiveConfig,
	}

	// Add flags
	rootCmd.PersistentFlags().StringVar(&configRoot, "config", "", "Configuration root directory (default: find configs/)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")

	// Add subcommands
	rootCmd.AddCommand(
		createListCommand(),
		createInfoCommand(),
		createCostCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runInteractiveConfig(cmd *cobra.Command, args []string) {
	// Find config root if not specified
	if configRoot == "" {
		configRoot = findConfigRoot()
	}

	fmt.Printf("🔬 AWS Research Wizard - Configuration Tool\n")
	fmt.Printf("Config Root: %s\n", configRoot)
	fmt.Printf("AWS Region: %s\n\n", region)

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

func createListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available research domains",
		Run: func(cmd *cobra.Command, args []string) {
			if configRoot == "" {
				configRoot = findConfigRoot()
			}

			loader := config.NewConfigLoader(configRoot)
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

func createInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info [domain]",
		Short: "Show detailed information about a domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if configRoot == "" {
				configRoot = findConfigRoot()
			}

			domainName := args[0]
			loader := config.NewConfigLoader(configRoot)
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

func createCostCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cost [domain]",
		Short: "Calculate costs for a specific domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if configRoot == "" {
				configRoot = findConfigRoot()
			}

			domainName := args[0]
			loader := config.NewConfigLoader(configRoot)
			domains, err := loader.LoadAllDomains()
			if err != nil {
				log.Fatalf("Failed to load domains: %v", err)
			}

			domain, exists := domains[domainName]
			if !exists {
				log.Fatalf("Domain '%s' not found", domainName)
			}

			fmt.Printf("💰 Cost Analysis: %s\n\n", domain.Name)

			_, _, err = tui.RunCostCalculator(domain, region)
			if err != nil {
				log.Fatalf("Failed to run cost calculator: %v", err)
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