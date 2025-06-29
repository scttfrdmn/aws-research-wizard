package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/aws-research-wizard/go/internal/commands/config"
	"github.com/aws-research-wizard/go/internal/commands/data"
	"github.com/aws-research-wizard/go/internal/commands/deploy"
	"github.com/aws-research-wizard/go/internal/commands/monitor"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "aws-research-wizard",
		Short: "AWS Research Wizard - Complete research environment management",
		Long: `AWS Research Wizard provides comprehensive tools for configuring, deploying, 
and monitoring optimized AWS research environments across multiple scientific domains.

Key Features:
- Interactive domain configuration with 18+ research packs
- S3 transfer optimization and AWS Open Data integration
- Infrastructure deployment with CloudFormation automation  
- Real-time monitoring with cost tracking and alerts
- Data pipeline orchestration and progress tracking
- Single binary distribution with zero dependencies
- Cross-platform support (Linux, macOS, Windows)

Perfect for:
- Research computing environments
- HPC workload deployment
- Scientific computing on AWS
- Cost-optimized research infrastructure`,
		Version: fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, gitCommit),
	}

	// Global flags
	rootCmd.PersistentFlags().String("region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().String("config-root", "", "Configuration root directory")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")

	// Add subcommands
	rootCmd.AddCommand(
		config.NewConfigCommand(),
		data.DataCmd,
		deploy.NewDeployCommand(), 
		monitor.NewMonitorCommand(),
	)

	// Version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("AWS Research Wizard %s\n", version)
			fmt.Printf("Built: %s\n", buildTime)
			fmt.Printf("Commit: %s\n", gitCommit)
			fmt.Printf("Go version: %s\n", getGoVersion())
		},
	})

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getGoVersion() string {
	// This would be set at build time, but for now return a placeholder
	return "go1.21+"
}

func init() {
	// Set up any global initialization
	if os.Getenv("AWS_RESEARCH_WIZARD_DEBUG") == "true" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}