package main

import (
	"log"

	"github.com/spack-go/spack-manager/pkg/manager"
	"github.com/spack-go/spack-manager/pkg/tui"
)

func main() {
	// Create a Spack manager
	config := manager.Config{
		SpackRoot:   "/opt/spack",
		BinaryCache: "https://cache.spack.io",
		WorkDir:     "/tmp/spack-manager-tui",
		LogLevel:    "info",
	}

	sm, err := manager.New(config)
	if err != nil {
		log.Fatalf("Failed to create Spack manager: %v", err)
	}

	// Launch the TUI
	if err := tui.Run(sm); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
