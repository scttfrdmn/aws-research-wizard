package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spack-go/spack-manager/pkg/manager"
	"github.com/spack-go/spack-manager/pkg/tui"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version", "--version", "-v":
		fmt.Printf("spack-manager version %s\n", version)
	case "help", "--help", "-h":
		printHelp()
	case "tui":
		runTUI()
	case "env":
		handleEnvironmentCommands()
	case "install":
		handleInstallCommand()
	case "list":
		handleListCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`spack-manager - A Go-based Spack package manager

Usage:
  spack-manager <command> [arguments]

Commands:
  tui                    Launch interactive Terminal User Interface
  env <subcommand>       Environment management commands
  install <env>          Install all packages in an environment
  list                   List all environments
  version               Show version information
  help                  Show this help message

Examples:
  spack-manager tui                    # Launch interactive TUI
  spack-manager env create my-env      # Create new environment
  spack-manager env list               # List environments
  spack-manager install my-env         # Install environment packages
  spack-manager list                   # List all environments

For more information, use 'spack-manager help' or run the TUI with 'spack-manager tui'.
`)
}

func printHelp() {
	fmt.Printf(`spack-manager - A Go-based Spack package manager

DESCRIPTION:
  spack-manager provides a modern interface for managing Spack package installations
  with real-time progress tracking and an interactive Terminal User Interface.

COMMANDS:
  tui
    Launch the interactive Terminal User Interface for full Spack management.
    This is the recommended way to use spack-manager.

  env <subcommand>
    Environment management commands:

    create <name>         Create a new Spack environment
    list                  List all available environments
    delete <name>         Delete an environment
    info <name>           Show detailed environment information
    add <name> <package>  Add a package to an environment

  install <environment>
    Install all packages in the specified environment with progress tracking.

  list
    List all available Spack environments.

CONFIGURATION:
  spack-manager looks for Spack installation in the following locations:
  - $SPACK_ROOT/bin/spack
  - /opt/spack/bin/spack
  - spack command in $PATH

EXAMPLES:
  # Launch interactive TUI (recommended)
  spack-manager tui

  # Create and set up a new environment
  spack-manager env create genomics
  spack-manager env add genomics gcc@11.3.0
  spack-manager env add genomics python@3.11
  spack-manager install genomics

  # List all environments
  spack-manager list

  # Get detailed help in TUI
  spack-manager tui  # then press '?' for help

AUTHOR:
  Spack Manager Go contributors

VERSION:
  %s

For more information and source code:
  https://github.com/spack-go/spack-manager
`, version)
}

func runTUI() {
	sm, err := createSpackManager()
	if err != nil {
		fmt.Printf("Error initializing Spack manager: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Launching Spack Manager TUI...")
	if err := tui.Run(sm); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func handleEnvironmentCommands() {
	if len(os.Args) < 3 {
		fmt.Println("Environment subcommand required")
		fmt.Println("Usage: spack-manager env <create|list|delete|info|add> [arguments]")
		os.Exit(1)
	}

	subcommand := os.Args[2]
	sm, err := createSpackManager()
	if err != nil {
		fmt.Printf("Error initializing Spack manager: %v\n", err)
		os.Exit(1)
	}

	switch subcommand {
	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Environment name required")
			fmt.Println("Usage: spack-manager env create <name>")
			os.Exit(1)
		}
		envName := os.Args[3]
		env := manager.Environment{
			Name:     envName,
			Packages: []string{},
		}
		if err := sm.CreateEnvironment(env); err != nil {
			fmt.Printf("Error creating environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Environment '%s' created successfully\n", envName)

	case "list":
		environments, err := sm.ListEnvironments()
		if err != nil {
			fmt.Printf("Error listing environments: %v\n", err)
			os.Exit(1)
		}

		if len(environments) == 0 {
			fmt.Println("No environments found")
			return
		}

		fmt.Println("Available environments:")
		for _, env := range environments {
			fmt.Printf("  • %s (%d packages)\n", env.Name, len(env.Packages))
		}

	case "delete":
		if len(os.Args) < 4 {
			fmt.Println("Environment name required")
			fmt.Println("Usage: spack-manager env delete <name>")
			os.Exit(1)
		}
		envName := os.Args[3]
		if err := sm.DeleteEnvironment(envName); err != nil {
			fmt.Printf("Error deleting environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Environment '%s' deleted successfully\n", envName)

	case "info":
		if len(os.Args) < 4 {
			fmt.Println("Environment name required")
			fmt.Println("Usage: spack-manager env info <name>")
			os.Exit(1)
		}
		envName := os.Args[3]
		env, err := sm.GetEnvironmentInfo(envName)
		if err != nil {
			fmt.Printf("Error getting environment info: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Environment: %s\n", env.Name)
		fmt.Printf("Packages (%d):\n", len(env.Packages))
		for _, pkg := range env.Packages {
			fmt.Printf("  • %s\n", pkg)
		}

	case "add":
		if len(os.Args) < 5 {
			fmt.Println("Environment name and package required")
			fmt.Println("Usage: spack-manager env add <env-name> <package-spec>")
			os.Exit(1)
		}
		envName := os.Args[3]
		packageSpec := os.Args[4]

		progressChan := make(chan manager.ProgressUpdate, 100)
		go func() {
			for update := range progressChan {
				if update.IsError {
					fmt.Printf("❌ %s\n", update.Message)
				} else {
					fmt.Printf("⚡ %s\n", update.Message)
				}
			}
		}()

		if err := sm.InstallPackage(envName, packageSpec, progressChan); err != nil {
			fmt.Printf("Error adding package: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Package '%s' added to environment '%s'\n", packageSpec, envName)

	default:
		fmt.Printf("Unknown environment subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleInstallCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Environment name required")
		fmt.Println("Usage: spack-manager install <environment>")
		os.Exit(1)
	}

	envName := os.Args[2]
	sm, err := createSpackManager()
	if err != nil {
		fmt.Printf("Error initializing Spack manager: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🚀 Installing environment '%s'...\n", envName)

	progressChan := make(chan manager.ProgressUpdate, 100)

	// Monitor progress in a goroutine
	go func() {
		for update := range progressChan {
			if update.IsError {
				fmt.Printf("❌ %s\n", update.Message)
			} else if update.IsComplete {
				fmt.Printf("✅ %s\n", update.Message)
			} else {
				fmt.Printf("⚡ %s (%.1f%%)\n", update.Message, update.Progress*100)
			}
		}
	}()

	if err := sm.InstallEnvironment(envName, progressChan); err != nil {
		fmt.Printf("❌ Installation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Environment '%s' installed successfully!\n", envName)
}

func handleListCommand() {
	sm, err := createSpackManager()
	if err != nil {
		fmt.Printf("Error initializing Spack manager: %v\n", err)
		os.Exit(1)
	}

	environments, err := sm.ListEnvironments()
	if err != nil {
		fmt.Printf("Error listing environments: %v\n", err)
		os.Exit(1)
	}

	if len(environments) == 0 {
		fmt.Println("No environments found")
		fmt.Println("Create one with: spack-manager env create <name>")
		return
	}

	fmt.Printf("Available environments (%d):\n", len(environments))
	for _, env := range environments {
		fmt.Printf("  • %s", env.Name)
		if len(env.Packages) > 0 {
			fmt.Printf(" (%d packages)", len(env.Packages))
		}
		fmt.Println()
	}
}

func createSpackManager() (manager.Manager, error) {
	// Try to find Spack installation
	spackRoot := os.Getenv("SPACK_ROOT")
	if spackRoot == "" {
		spackRoot = "/opt/spack" // Default location
	}

	config := manager.Config{
		SpackRoot: spackRoot,
		WorkDir:   filepath.Join(os.TempDir(), "spack-manager"),
		LogLevel:  "info",
	}

	// Check for binary cache configuration
	if binaryCache := os.Getenv("SPACK_BINARY_CACHE"); binaryCache != "" {
		config.BinaryCache = binaryCache
	}

	return manager.New(config)
}
