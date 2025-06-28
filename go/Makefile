# AWS Research Wizard - Go Implementation Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_CONFIG=aws-research-wizard-config
BINARY_DEPLOY=aws-research-wizard-deploy
BINARY_MONITOR=aws-research-wizard-monitor

# Build directory
BUILD_DIR=build

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse HEAD)

# Linker flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: all build clean test deps help install

# Default target
all: deps build

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build all binaries
build: deps
	@echo "Building AWS Research Wizard Go binaries..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CONFIG) ./cmd/config
	@echo "Built $(BUILD_DIR)/$(BINARY_CONFIG)"
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DEPLOY) ./cmd/deploy
	@echo "Built $(BUILD_DIR)/$(BINARY_DEPLOY)"

# Build for multiple platforms
build-all: deps
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CONFIG)-linux-amd64 ./cmd/config
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DEPLOY)-linux-amd64 ./cmd/deploy
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CONFIG)-linux-arm64 ./cmd/config
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DEPLOY)-linux-arm64 ./cmd/deploy
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CONFIG)-darwin-amd64 ./cmd/config
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DEPLOY)-darwin-amd64 ./cmd/deploy
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CONFIG)-darwin-arm64 ./cmd/config
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DEPLOY)-darwin-arm64 ./cmd/deploy
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CONFIG)-windows-amd64.exe ./cmd/config
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DEPLOY)-windows-amd64.exe ./cmd/deploy
	
	@echo "Built binaries for multiple platforms in $(BUILD_DIR)/"

# Install to local bin
install: build
	@echo "Installing to ~/bin..."
	@mkdir -p ~/bin
	cp $(BUILD_DIR)/$(BINARY_CONFIG) ~/bin/aws-research-wizard-config
	cp $(BUILD_DIR)/$(BINARY_DEPLOY) ~/bin/aws-research-wizard-deploy
	@echo "Installed binaries to ~/bin/:"
	@echo "  aws-research-wizard-config (configuration tool)"
	@echo "  aws-research-wizard-deploy (deployment tool)"
	@echo "Add ~/bin to your PATH if not already included"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the config tool
run: build
	./$(BUILD_DIR)/$(BINARY_CONFIG)

# Run with example domain
run-example: build
	./$(BUILD_DIR)/$(BINARY_CONFIG) info genomics

# Format code
fmt:
	@echo "Formatting Go code..."
	$(GOCMD) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Run security scanner
security:
	@echo "Running security scanner..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development tools installed"

# Check dependencies for vulnerabilities
check-deps:
	@echo "Checking dependencies for vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not found. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

# Quick development build and run
dev: build run

# Show help
help:
	@echo "AWS Research Wizard - Go Implementation"
	@echo ""
	@echo "Available targets:"
	@echo "  all          - Install dependencies and build (default)"
	@echo "  build        - Build the config binary"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  install      - Install to ~/bin/"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the config tool"
	@echo "  run-example  - Build and run with example domain"
	@echo "  fmt          - Format Go code"
	@echo "  lint         - Run linter"
	@echo "  security     - Run security scanner"
	@echo "  dev-setup    - Install development tools"
	@echo "  check-deps   - Check dependencies for vulnerabilities"
	@echo "  dev          - Quick build and run for development"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    # Build the binary"
	@echo "  make run                      # Interactive domain selection"
	@echo "  make run-example              # Show genomics domain info"
	@echo "  ./build/aws-research-wizard-config list    # List all domains"
	@echo "  ./build/aws-research-wizard-config cost genomics  # Cost analysis"