# Resume Development Guide

This guide provides step-by-step instructions for resuming AWS Research Wizard development.

## Quick Start Commands

```bash
# Navigate to project
cd /Users/scttfrdmn/src/aws-research-wizard

# Check current state
git status
git log --oneline -5

# Build and test current version
cd go && make build
./build/aws-research-wizard version
./build/aws-research-wizard --help

# Check todo list (use TodoRead tool in Claude Code)
```

## Current Development State

### ✅ **Completed Foundation**
- **Unified Binary**: Single `aws-research-wizard` executable with subcommands
- **AWS Integration**: Full SDK integration with real AWS testing ($62.23 costs verified)
- **CI/CD Pipeline**: GitHub Actions with multi-platform builds and testing
- **Installation System**: Universal installer and package manager templates
- **Documentation**: Comprehensive installation and deployment guides

### 🎯 **Next Priority: Phase 1 Data Management Engine**

**Objective**: Implement S3 transfer optimization and AWS Open Data integration

**Files to Create**:
```
go/internal/data/
├── s3_manager.go           # S3 transfer optimization
├── open_data.go            # AWS Open Data integration  
├── pipeline.go             # Data pipeline orchestration
└── transfer_monitor.go     # Transfer progress tracking

go/internal/commands/data/
└── data.go                 # Data management subcommand
```

**Implementation Steps**:
1. Create S3 transfer manager with multi-part upload optimization
2. Add AWS Open Data registry integration
3. Implement data pipeline orchestration
4. Add transfer progress monitoring
5. Create data management subcommand

## Key Development Context

### **Project Architecture**
- **Language**: Go 1.21+
- **CLI Framework**: Cobra for command structure
- **TUI Framework**: Bubble Tea and Lip Gloss for interactive interfaces  
- **AWS SDK**: v2 with comprehensive service integration
- **Build Tool**: Make with cross-platform support

### **Current Binary Structure**
```bash
aws-research-wizard
├── config      # Domain configuration and cost analysis
├── deploy      # Infrastructure deployment and management
├── monitor     # Real-time monitoring and dashboards  
└── version     # Version information
```

### **Successful Integration Points**
- **AWS Services**: EC2, CloudFormation, CloudWatch, Cost Explorer, IAM, S3
- **Real Data**: Successfully tested with live AWS account
- **Cost Tracking**: Real-time cost analysis ($62.23 actual costs retrieved)
- **Infrastructure**: CloudFormation stack management working
- **Monitoring**: Live dashboard with metrics and alerts

## Development Workflow

### **1. Start Development Session**
```bash
cd /Users/scttfrdmn/src/aws-research-wizard
git checkout main

# Review current state
cat PROJECT_STATUS.md
cat DEPLOYMENT_STRATEGY.md
```

### **2. Use Todo Management**
```bash
# In Claude Code, use TodoRead to check current tasks
# Use TodoWrite to plan new implementation tasks
```

### **3. Implement Features**
```bash
# Create new directories
mkdir -p go/internal/data
mkdir -p go/internal/commands/data

# Start with S3 manager implementation
# Test incrementally with: make build && ./go/build/aws-research-wizard
```

### **4. Testing and Quality**
```bash
# Build and test
make build
./build/aws-research-wizard --help

# Run tests
make test

# Check code quality  
make lint

# Build for all platforms
make build-all
```

### **5. Commit Progress**
```bash
git add .
git commit -m "feat: implement S3 transfer optimization"
git push origin main  # If using remote repository
```

## Available Make Targets

```bash
make build          # Build unified binary (default)
make build-legacy   # Build separate binaries (compatibility)
make build-all      # Cross-platform builds
make install        # Install to ~/bin/
make test           # Run tests
make lint           # Code quality check
make clean          # Clean build artifacts
make run            # Build and run
make help           # Show all targets
```

## Expected Implementation Flow

### **Phase 1A: S3 Transfer Manager** (2-3 hours)
- Implement `go/internal/data/s3_manager.go`
- Multi-part upload optimization
- Progress tracking and resumption
- Bandwidth throttling capabilities

### **Phase 1B: AWS Open Data Integration** (2-3 hours)  
- Implement `go/internal/data/open_data.go`
- Registry discovery and browsing
- Cost-free data access optimization
- Metadata management

### **Phase 1C: Data Pipeline Orchestration** (2-3 hours)
- Implement `go/internal/data/pipeline.go`
- Transfer job scheduling
- Error handling and retry logic
- Monitoring and alerting

### **Phase 1D: Command Integration** (1-2 hours)
- Create `go/internal/commands/data/data.go`
- Add data subcommand to main.go
- Update help documentation
- Integration testing

## Key Files Reference

### **Main Implementation Areas**
- `go/cmd/main.go` - Unified entry point
- `go/internal/aws/` - AWS SDK integration (complete)
- `go/internal/commands/` - Command modules
- `go/internal/tui/` - Terminal UI components
- `go/Makefile` - Build system

### **Configuration and Docs**
- `PROJECT_STATUS.md` - Current state summary
- `DEPLOYMENT_STRATEGY.md` - Distribution strategy
- `INSTALLATION.md` - User installation guide
- `.github/workflows/` - CI/CD automation

### **Current Working Binary**
```bash
./go/build/aws-research-wizard version
# AWS Research Wizard dev
# Built: 2024-06-29_00:22:33
# Commit: dea4fdfdc1c5b5ccafd91b78a06d88335b5bb3d3
# Go version: go1.21+
```

## Success Indicators

### **When Resuming Successfully**
- ✅ Build completes without errors: `make build`
- ✅ All subcommands show help: `aws-research-wizard {config|deploy|monitor} --help`
- ✅ Version information displays correctly: `aws-research-wizard version`
- ✅ Todo list shows current priorities: Use TodoRead

### **Implementation Readiness**
- ✅ AWS SDK integration complete and tested
- ✅ Build system fully functional
- ✅ Command structure established
- ✅ TUI framework operational
- ✅ Documentation comprehensive

## Troubleshooting

### **Build Issues**
```bash
# Clean and rebuild
make clean && make build

# Check Go environment
go version  # Should be 1.21+
go mod tidy

# Verify dependencies
cd go && go mod download
```

### **Git State Issues**
```bash
# Check status
git status
git log --oneline -5

# If uncommitted changes, either commit or stash
git add . && git commit -m "wip: development pause"
# OR
git stash
```

### **Environment Issues**
```bash
# Verify project location
pwd  # Should be /Users/scttfrdmn/src/aws-research-wizard

# Check file permissions
ls -la go/build/aws-research-wizard  # Should be executable
```

## Context for AI Assistance

When working with AI assistance (Claude Code) on this project:

1. **Current State**: Binary consolidation and deployment infrastructure complete
2. **Next Goal**: Implement S3 transfer optimization and AWS Open Data integration  
3. **Architecture**: Go-based CLI with Cobra commands and Bubble Tea TUI
4. **Working Directory**: `/Users/scttfrdmn/src/aws-research-wizard`
5. **Build Command**: `cd go && make build`
6. **Test Command**: `./go/build/aws-research-wizard --help`

The project is in excellent shape for continued development with a solid foundation and clear next steps.