# AWS Research Wizard - Go Implementation

A high-performance, single-binary implementation of the AWS Research Wizard for simplified distribution and deployment in HPC environments.

## Overview

The Go implementation provides a lightweight, fast alternative to the Python version with:
- **Single Binary Distribution**: No Python installation or dependencies required
- **Cross-Platform Support**: Pre-built binaries for Linux, macOS, and Windows
- **SSH-Optimized**: Perfect for remote HPC and research computing environments
- **Fast Startup**: Sub-second response times for all operations
- **Shared Configuration**: Uses the same YAML configs as the Python version

## Development Guidelines

**⚠️ CRITICAL DEVELOPMENT RULE**: When working on this codebase, **never work around or fake solutions to problems**. Always fix them properly. See [CLAUDE.md](go/CLAUDE.md) for detailed development guidelines and architecture documentation.

## Quick Start

### Download and Install

```bash
# Download for your platform (example for Linux AMD64)
curl -L https://github.com/your-org/aws-research-wizard/releases/latest/download/aws-research-wizard-linux-amd64 -o aws-research-wizard
chmod +x aws-research-wizard
sudo mv aws-research-wizard /usr/local/bin/

# Or install to local bin
mkdir -p ~/bin
mv aws-research-wizard ~/bin/
export PATH="$HOME/bin:$PATH"
```

### Basic Usage

```bash
# List all available research domains
aws-research-wizard list

# Get detailed info about a domain
aws-research-wizard info genomics

# Interactive domain selection and cost calculation
aws-research-wizard

# Cost analysis for specific domain
aws-research-wizard cost climate_modeling
```

## Features

### Interactive Domain Selection
- Browse 18+ research domains with cost estimates
- Rich terminal interface with keyboard navigation
- Real-time cost calculations and optimization suggestions

### Cost Calculator
- Instance type recommendations based on domain requirements
- Spot instance savings calculations (70% savings)
- Reserved instance optimization suggestions
- Monthly and annual cost projections

### Research Domain Support
All Python version domains are supported:
- Genomics & Bioinformatics
- Climate Modeling
- Machine Learning
- Atmospheric Chemistry
- Computational Physics
- Materials Science
- Neuroscience
- And 11 more domains...

## Building from Source

### Prerequisites
- Go 1.21 or later
- Git

### Build Instructions

```bash
# Clone the repository
git clone https://github.com/your-org/aws-research-wizard.git
cd aws-research-wizard/go

# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Development

```bash
# Install development dependencies
make dev-setup

# Run tests
make test

# Run with live reload during development
make dev

# Format code
make fmt

# Run linter
make lint

# Security scan
make security
```

## Configuration

The Go implementation uses the same configuration files as the Python version:

### Shared Configuration Structure
```
configs/
├── domains/           # Domain pack configurations
│   ├── genomics.yaml
│   ├── climate_modeling.yaml
│   └── ...
├── schemas/           # Validation schemas
└── templates/         # AWS templates
```

### Domain Configuration Format
Each domain is defined in YAML with:
- **Spack Packages**: Organized by category (aligners, variant_calling, etc.)
- **AWS Instances**: Performance-matched recommendations
- **Cost Estimates**: Compute, storage, and total monthly costs
- **Workflow Tools**: Nextflow, Snakemake, Cromwell integration

## Performance Comparison

| Metric | Go Implementation | Python Implementation |
|--------|------------------|--------------------- |
| Startup Time | < 0.1s | ~3s |
| Memory Usage | ~15MB | ~150MB |
| Binary Size | ~20MB | ~500MB (with virtualenv) |
| Installation | Single download | pip install + dependencies |
| Cross-compilation | Native | Requires Python on target |

## Command Reference

### Main Commands

```bash
# Interactive mode (default)
aws-research-wizard

# List domains
aws-research-wizard list

# Domain information
aws-research-wizard info [domain]

# Cost analysis
aws-research-wizard cost [domain]
```

### Options

```bash
# Specify configuration directory
aws-research-wizard --config /path/to/configs list

# Set AWS region
aws-research-wizard --region us-west-2 cost genomics

# Get help
aws-research-wizard --help
```

## Use Cases

### HPC Environments
- Single binary deployment without package management
- No Python dependency conflicts
- Fast execution for job submission scripts
- Minimal resource footprint on login nodes

### Remote SSH Access
- Terminal-based interface works over any SSH connection
- No X11 forwarding required
- Works in screen/tmux sessions
- Optimized for low-bandwidth connections

### CI/CD Integration
- Fast execution in automation pipelines
- Cross-platform builds in GitHub Actions
- Docker-friendly single binary
- Consistent behavior across environments

## Distribution

### Release Binaries

Pre-built binaries are available for:
- Linux AMD64/ARM64
- macOS AMD64/ARM64 (Intel/Apple Silicon)
- Windows AMD64

### Docker Usage

```bash
# Pull and run
docker run --rm -it -v $(pwd)/configs:/configs \
  ghcr.io/your-org/aws-research-wizard:latest list

# Or build locally
docker build -t aws-research-wizard .
docker run --rm -it aws-research-wizard list
```

### Package Managers

```bash
# Homebrew (macOS/Linux)
brew install your-org/tap/aws-research-wizard

# Snap (Linux)
sudo snap install aws-research-wizard

# Chocolatey (Windows)
choco install aws-research-wizard
```

## Integration with Python Version

Both implementations can coexist and share configurations:

```bash
# Use Go for quick queries
aws-research-wizard list
aws-research-wizard info genomics

# Use Python for advanced features
python python/domain_tutorial_generator.py
python python/tui_aws_monitor.py
```

### When to Use Each Version

**Go Implementation:**
- Quick domain browsing and cost estimation
- HPC job submission and automation
- Remote SSH environments
- Production deployment scripts
- CI/CD pipelines

**Python Implementation:**
- Jupyter notebook generation
- Advanced TUI monitoring dashboards
- Complex data analysis workflows
- Research tutorial creation
- Development and experimentation

## Architecture

### Code Organization

```
go/
├── cmd/                    # Main applications
│   ├── config/            # Domain configuration tool
│   ├── deploy/            # Deployment automation (future)
│   └── monitor/           # Resource monitoring (future)
├── internal/              # Internal packages
│   ├── config/            # Configuration loading
│   ├── aws/               # AWS SDK integration
│   └── tui/               # Terminal UI components
├── build/                 # Build artifacts
└── Makefile              # Build automation
```

### Key Components

- **Config Loader**: YAML parsing with validation
- **AWS Integration**: Cost calculation and instance recommendations
- **TUI Framework**: Bubble Tea-based interactive interfaces
- **CLI Framework**: Cobra-based command structure

## Contributing

### Development Setup

```bash
# Fork and clone
git clone https://github.com/your-username/aws-research-wizard.git
cd aws-research-wizard/go

# Install development tools
make dev-setup

# Make changes and test
make test
make lint
make build

# Submit pull request
```

### Adding New Features

1. **New Commands**: Add to `cmd/` directory
2. **TUI Components**: Extend `internal/tui/`
3. **AWS Integration**: Enhance `internal/aws/`
4. **Configuration**: Update `internal/config/`

### Code Standards

- Go modules for dependency management
- `gofmt` and `golangci-lint` for code quality
- Comprehensive unit tests
- Security scanning with `gosec`
- Semantic versioning for releases

## Future Roadmap

### Phase 1 (Current)
- ✅ Domain configuration and cost calculation
- ✅ Interactive TUI interfaces
- ✅ Cross-platform builds

### Phase 2 (Planned)
- Deployment automation (`aws-research-wizard deploy`)
- Real-time monitoring (`aws-research-wizard monitor`)
- AWS resource management
- Configuration generation and validation

### Phase 3 (Future)
- Integration with Python workflow engine
- Advanced cost optimization algorithms
- Multi-cloud support (Azure, GCP)
- Plugin system for custom domains

## Support

### Documentation
- [Go Implementation Guide](README_GO.md) (this document)
- [Python Implementation Guide](README.md)
- [TUI User Guide](README_TUI.md)
- [Configuration Schema](configs/schemas/)

### Community
- GitHub Issues for bug reports
- Discussions for feature requests
- Wiki for advanced usage examples
- Contributing guidelines for developers

## License

Same license as the main AWS Research Wizard project.
