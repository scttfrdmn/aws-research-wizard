package main

import (
	"fmt"
	"log"

	"github.com/spack-go/spack-manager/pkg/manager"
)

func main() {
	// Create a new Spack manager
	config := manager.Config{
		SpackRoot:   "/opt/spack",
		BinaryCache: "https://cache.spack.io",
		WorkDir:     "/tmp/spack-manager-example",
		LogLevel:    "info",
	}

	sm, err := manager.New(config)
	if err != nil {
		log.Fatalf("Failed to create Spack manager: %v", err)
	}

	// Print Spack version
	version, err := sm.GetSpackVersion()
	if err != nil {
		log.Fatalf("Failed to get Spack version: %v", err)
	}
	fmt.Printf("Using Spack version: %s\n", version)

	// Create a new environment
	env := manager.Environment{
		Name: "example-env",
		Packages: []string{
			"gcc@11.3.0",
			"python@3.11",
			"numpy",
		},
	}

	fmt.Printf("Creating environment '%s'...\n", env.Name)
	if err := sm.CreateEnvironment(env); err != nil {
		log.Fatalf("Failed to create environment: %v", err)
	}

	// List all environments
	environments, err := sm.ListEnvironments()
	if err != nil {
		log.Fatalf("Failed to list environments: %v", err)
	}

	fmt.Println("Available environments:")
	for _, e := range environments {
		fmt.Printf("  - %s (%d packages)\n", e.Name, len(e.Packages))
	}

	// Install the environment with progress tracking
	fmt.Printf("Installing environment '%s'...\n", env.Name)

	progressChan := make(chan manager.ProgressUpdate, 100)

	// Monitor progress in a separate goroutine
	go func() {
		for update := range progressChan {
			if update.IsError {
				fmt.Printf("❌ Error: %s\n", update.Message)
			} else if update.IsComplete {
				fmt.Printf("✅ Complete: %s\n", update.Message)
			} else {
				fmt.Printf("📦 Installing %s: %.1f%% - %s\n",
					update.Package, update.Progress*100, update.Message)
			}
		}
	}()

	if err := sm.InstallEnvironment(env.Name, progressChan); err != nil {
		log.Fatalf("Failed to install environment: %v", err)
	}

	fmt.Printf("🎉 Environment '%s' installed successfully!\n", env.Name)

	// Get detailed environment information
	envInfo, err := sm.GetEnvironmentInfo(env.Name)
	if err != nil {
		log.Fatalf("Failed to get environment info: %v", err)
	}

	fmt.Printf("\nEnvironment '%s' details:\n", envInfo.Name)
	fmt.Printf("Installed packages (%d):\n", len(envInfo.Packages))
	for _, pkg := range envInfo.Packages {
		fmt.Printf("  ✓ %s\n", pkg)
	}
}
