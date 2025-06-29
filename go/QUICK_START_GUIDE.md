# AWS Research Wizard - Quick Start Guide

## 🚀 Get Started in 5 Minutes

The AWS Research Wizard helps research institutions optimize their data movement to AWS with intelligent recommendations and automated workflows.

## ⚡ Quick Installation

### 1. Download and Install
```bash
# Download the binary (replace with actual download URL)
wget https://github.com/aws-research-wizard/releases/latest/aws-research-wizard-linux-amd64
chmod +x aws-research-wizard-linux-amd64
sudo mv aws-research-wizard-linux-amd64 /usr/local/bin/aws-research-wizard

# Verify installation
aws-research-wizard --version
```

### 2. Setup AWS Credentials
```bash
# Configure AWS CLI (if not already done)
aws configure
# Enter your AWS Access Key, Secret Key, Region, and Output format

# Test AWS connectivity
aws-research-wizard data diagnose
```

## 🎯 Basic Usage

### Analyze Your Research Data
```bash
# Get intelligent recommendations for your data
aws-research-wizard data analyze /path/to/your/research/data --verbose

# This shows:
# - File type analysis
# - Size distribution 
# - Transfer recommendations
# - Cost estimates
```

### Generate Optimized Configuration
```bash
# Auto-generate configuration for your research domain
aws-research-wizard data generate /path/to/your/data \
  --domain genomics \
  --output research-config.yaml

# Supported domains:
# genomics, climate, machine_learning, astronomy, geospatial, chemistry
```

### Upload Data to S3
```bash
# Quick upload with progress tracking
aws-research-wizard data upload /path/to/file.dat s3://your-bucket/file.dat

# Bulk upload with workflow
aws-research-wizard data workflow run \
  --config research-config.yaml \
  --workflow upload_data
```

## 📊 Research Domain Examples

### Genomics Research
```bash
# Optimized for FASTQ, BAM, VCF files
aws-research-wizard data generate /data/genomics \
  --domain genomics \
  --output genomics-config.yaml

# Features: Small file bundling, encryption, cost optimization
```

### Climate Science  
```bash
# Optimized for NetCDF, GRIB files
aws-research-wizard data generate /data/climate \
  --domain climate \
  --output climate-config.yaml

# Features: Large file handling, time-series organization
```

### Machine Learning
```bash
# Optimized for training datasets
aws-research-wizard data generate /data/ml \
  --domain machine_learning \
  --output ml-config.yaml

# Features: GPU optimization, model versioning, high-volume datasets
```

## 🔧 Common Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `analyze` | Analyze data patterns | `aws-research-wizard data analyze /data` |
| `generate` | Create optimized config | `aws-research-wizard data generate /data --domain genomics` |
| `validate` | Check configuration | `aws-research-wizard data validate config.yaml` |
| `upload` | Upload files | `aws-research-wizard data upload file.txt s3://bucket/` |
| `download` | Download files | `aws-research-wizard data download s3://bucket/file.txt` |
| `workflow` | Run workflows | `aws-research-wizard data workflow run --config config.yaml` |
| `monitor` | Monitor transfers | `aws-research-wizard data monitor` |
| `diagnose` | System health check | `aws-research-wizard data diagnose` |

## 💡 Pro Tips

### 1. Start with Analysis
Always analyze your data first to get personalized recommendations:
```bash
aws-research-wizard data analyze /your/data --verbose
```

### 2. Use Dry-Run
Test workflows before executing:
```bash
aws-research-wizard data workflow run --config config.yaml --dry-run
```

### 3. Monitor Progress
Watch transfers in real-time:
```bash
aws-research-wizard data monitor --refresh 1s
```

### 4. Check System Health
Regular diagnostics ensure everything is working:
```bash
aws-research-wizard data diagnose
```

## 🆘 Need Help?

### Quick Diagnostics
```bash
# Check system health
aws-research-wizard data diagnose --verbose

# Validate your configuration
aws-research-wizard data validate your-config.yaml --verbose
```

### Get Command Help
```bash
# General help
aws-research-wizard --help

# Specific command help
aws-research-wizard data analyze --help
```

### Common Issues

**AWS Credentials Error:**
```bash
aws configure  # Reconfigure credentials
aws sts get-caller-identity  # Test access
```

**Permission Denied:**
```bash
aws-research-wizard data diagnose --verbose  # Check permissions
```

**Slow Performance:**
```bash
aws-research-wizard data analyze /data --performance-report
```

## 📋 Example Workflow

Here's a complete example for genomics research:

```bash
# 1. Analyze your genomics data
aws-research-wizard data analyze /data/genomics --verbose

# 2. Generate optimized configuration
aws-research-wizard data generate /data/genomics \
  --domain genomics \
  --output genomics-project.yaml

# 3. Validate the configuration
aws-research-wizard data validate genomics-project.yaml

# 4. Test with dry-run
aws-research-wizard data workflow run \
  --config genomics-project.yaml \
  --workflow upload_data \
  --dry-run

# 5. Execute the actual upload
aws-research-wizard data workflow run \
  --config genomics-project.yaml \
  --workflow upload_data

# 6. Monitor progress
aws-research-wizard data monitor
```

## 🎯 What's Next?

- **Read the Full Guide**: See `PRODUCTION_DEPLOYMENT_GUIDE.md` for detailed setup
- **Explore Examples**: Check the `examples/` directory for domain-specific configurations
- **Join the Community**: Get support and share best practices
- **Optimize Further**: Use advanced features for enterprise deployments

---

**You're now ready to optimize your research data movement with AWS Research Wizard!** 🚀