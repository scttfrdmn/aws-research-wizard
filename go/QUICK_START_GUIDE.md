# AWS Research Wizard - Quick Start Guide

## 🚀 Get Started in 5 Minutes

The AWS Research Wizard helps you set up powerful research computers in the cloud. It finds the best options for your work and saves you money.

## ⚡ Easy Installation

### 1. Download and Install
```bash
# Download the program
wget https://github.com/aws-research-wizard/releases/latest/aws-research-wizard-linux-amd64
chmod +x aws-research-wizard-linux-amd64
sudo mv aws-research-wizard-linux-amd64 /usr/local/bin/aws-research-wizard

# Check that it works
aws-research-wizard --version
```

### 2. Connect to AWS
```bash
# Set up your AWS account connection
aws configure
# Enter your AWS Access Key, Secret Key, Region, and Output format

# Test that everything is connected
aws-research-wizard data diagnose
```

## 🎯 How to Use It

### Look at Your Research Data
```bash
# Get smart recommendations for your data
aws-research-wizard data analyze /path/to/your/research/data --verbose

# This shows you:
# - What types of files you have
# - How big they are
# - Best ways to move them
# - What it will cost
```

### Create a Setup for Your Research
```bash
# Make a configuration file for your research area
aws-research-wizard data generate /path/to/your/data \
  --domain genomics \
  --output research-config.yaml

# Available research areas:
# genomics, climate, machine_learning, astronomy, geospatial, chemistry
```

### Upload Your Data
```bash
# Upload a single file with progress bar
aws-research-wizard data upload /path/to/file.dat s3://your-bucket/file.dat

# Upload lots of files using a workflow
aws-research-wizard data workflow run \
  --config research-config.yaml \
  --workflow upload_data
```

## 📊 Research Area Examples

### Genomics Research
```bash
# Works great with FASTQ, BAM, VCF files
aws-research-wizard data generate /data/genomics \
  --domain genomics \
  --output genomics-config.yaml

# Special features: Bundles small files, encrypts data, saves money
```

### Climate Science
```bash
# Works great with NetCDF, GRIB files
aws-research-wizard data generate /data/climate \
  --domain climate \
  --output climate-config.yaml

# Special features: Handles huge files, organizes time-series data
```

### Machine Learning
```bash
# Works great with training datasets
aws-research-wizard data generate /data/ml \
  --domain machine_learning \
  --output ml-config.yaml

# Special features: GPU optimization, model versioning, handles massive datasets
```

## 🔧 Common Commands

| Command | What It Does | Example |
|---------|-------------|---------|
| `analyze` | Look at your data patterns | `aws-research-wizard data analyze /data` |
| `generate` | Create optimized setup | `aws-research-wizard data generate /data --domain genomics` |
| `validate` | Check your configuration | `aws-research-wizard data validate config.yaml` |
| `upload` | Upload files | `aws-research-wizard data upload file.txt s3://bucket/` |
| `download` | Download files | `aws-research-wizard data download s3://bucket/file.txt` |
| `workflow` | Run workflows | `aws-research-wizard data workflow run --config config.yaml` |
| `monitor` | Watch transfers | `aws-research-wizard data monitor` |
| `diagnose` | Check system health | `aws-research-wizard data diagnose` |

## 💡 Helpful Tips

### 1. Always Check Your Data First
Look at your data first to get personalized recommendations:
```bash
aws-research-wizard data analyze /your/data --verbose
```

### 2. Test Before Running
Test workflows before actually doing them:
```bash
aws-research-wizard data workflow run --config config.yaml --dry-run
```

### 3. Watch Your Progress
See transfers happening in real-time:
```bash
aws-research-wizard data monitor --refresh 1s
```

### 4. Check System Health
Regular check-ups make sure everything is working:
```bash
aws-research-wizard data diagnose
```

## 🆘 Need Help?

### Quick Problem Check
```bash
# Check system health
aws-research-wizard data diagnose --verbose

# Check your configuration file
aws-research-wizard data validate your-config.yaml --verbose
```

### Get Help with Commands
```bash
# General help
aws-research-wizard --help

# Help with a specific command
aws-research-wizard data analyze --help
```

### Common Problems

**AWS Account Connection Error:**
```bash
aws configure  # Set up credentials again
aws sts get-caller-identity  # Test that it works
```

**Permission Denied:**
```bash
aws-research-wizard data diagnose --verbose  # Check what's wrong
```

**Slow Performance:**
```bash
aws-research-wizard data analyze /data --performance-report
```

## 📋 Complete Example

Here's a complete example for genomics research:

```bash
# 1. Look at your genomics data
aws-research-wizard data analyze /data/genomics --verbose

# 2. Create an optimized setup
aws-research-wizard data generate /data/genomics \
  --domain genomics \
  --output genomics-project.yaml

# 3. Check the configuration
aws-research-wizard data validate genomics-project.yaml

# 4. Test it first (dry-run)
aws-research-wizard data workflow run \
  --config genomics-project.yaml \
  --workflow upload_data \
  --dry-run

# 5. Do the actual upload
aws-research-wizard data workflow run \
  --config genomics-project.yaml \
  --workflow upload_data

# 6. Watch the progress
aws-research-wizard data monitor
```

## 🎯 What's Next?

- **Read the Full Guide**: See `PRODUCTION_DEPLOYMENT_GUIDE.md` for detailed setup
- **Try Examples**: Check the `examples/` directory for research-specific configurations
- **Join the Community**: Get support and share tips with other researchers
- **Use Advanced Features**: Explore enterprise features for bigger projects

---

**You're now ready to run your research in the cloud with AWS Research Wizard!** 🚀
