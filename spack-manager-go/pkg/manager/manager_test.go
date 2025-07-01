package manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "spack-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a fake spack binary
	spackBin := filepath.Join(tempDir, "bin", "spack")
	if err := os.MkdirAll(filepath.Dir(spackBin), 0755); err != nil {
		t.Fatalf("Failed to create bin dir: %v", err)
	}

	// Create a simple script that outputs version
	script := `#!/bin/bash
if [ "$1" = "version" ]; then
    echo "0.20.1"
    exit 0
fi
exit 1
`
	if err := os.WriteFile(spackBin, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to write fake spack binary: %v", err)
	}

	config := Config{
		SpackRoot: tempDir,
		WorkDir:   filepath.Join(tempDir, "work"),
		LogLevel:  "info",
	}

	sm, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create SpackManager: %v", err)
	}

	if sm.config.SpackRoot != tempDir {
		t.Errorf("Expected SpackRoot %s, got %s", tempDir, sm.config.SpackRoot)
	}
}

func TestValidateEnvironment(t *testing.T) {
	config := Config{
		SpackRoot: "/fake/path", // Won't be used for validation test
		WorkDir:   "/tmp",
		LogLevel:  "info",
	}

	sm := &SpackManager{config: config}

	tests := []struct {
		name    string
		env     Environment
		wantErr bool
	}{
		{
			name: "valid environment",
			env: Environment{
				Name:     "test-env",
				Packages: []string{"gcc@11.3.0"},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			env: Environment{
				Name:     "",
				Packages: []string{"gcc@11.3.0"},
			},
			wantErr: true,
		},
		{
			name: "no packages",
			env: Environment{
				Name:     "test-env",
				Packages: []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateEnvironment(tt.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
