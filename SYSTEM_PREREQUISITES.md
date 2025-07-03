# AWS Research Wizard - System Prerequisites & Setup Guide

## Overview

The AWS Research Wizard provides comprehensive research computing environments on AWS. This guide covers local system requirements and setup procedures for macOS, Windows 11, and Linux (Ubuntu) systems.

## 🎯 Quick Start Checklist

- [ ] Python 3.9+ installed
- [ ] AWS CLI v2 configured
- [ ] Git installed
- [ ] 8GB+ RAM available
- [ ] 10GB+ free disk space
- [ ] Internet connection for AWS and package downloads

---

## 📋 System Requirements

### Minimum Requirements
- **RAM**: 8GB (16GB recommended for large workloads)
- **Storage**: 10GB free space (50GB+ for local Spack builds)
- **Network**: Reliable internet connection (>10 Mbps recommended)
- **Python**: 3.9 or higher
- **Operating System**:
  - macOS 11.0+ (Big Sur or later)
  - Windows 11 (64-bit)
  - Ubuntu 20.04 LTS or later

### Recommended Requirements
- **RAM**: 16GB+ for optimal performance
- **Storage**: 100GB+ SSD for local development
- **Network**: High-speed connection for large data transfers
- **CPU**: Multi-core processor (4+ cores recommended)

---

# 🍎 macOS Setup Guide

## Prerequisites Installation

### 1. Install Homebrew (Package Manager)
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Install Essential Tools
```bash
# Install Python 3.11 (recommended version)
brew install python@3.11

# Install Git
brew install git

# Install AWS CLI v2
brew install awscli

# Install additional development tools
brew install cmake gcc gfortran
```

### 3. Install Xcode Command Line Tools
```bash
xcode-select --install
```

### 4. Install Python Virtual Environment Tools
```bash
# Install pipenv for virtual environments
pip3 install pipenv

# Alternative: use built-in venv
python3 -m pip install virtualenv
```

## Setup AWS Research Wizard

### 1. Clone Repository
```bash
git clone https://github.com/aws-research-wizard/aws-research-wizard.git
cd aws-research-wizard
```

### 2. Create Virtual Environment
```bash
# Using pipenv (recommended)
pipenv install

# Or using venv
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

### 3. Configure AWS Credentials
```bash
aws configure
# Enter your AWS Access Key ID, Secret Access Key, Region, and Output format
```

### 4. Verify Installation
```bash
# Activate environment
pipenv shell  # or source venv/bin/activate

# Run environment checker
python aws_environment_checker.py

# Start GUI
streamlit run gui_research_wizard.py
```

## macOS-Specific Notes

- **Apple Silicon (M1/M2)**: Some scientific packages may require Rosetta 2
  ```bash
  softwareupdate --install-rosetta
  ```
- **Gatekeeper**: You may need to allow unsigned packages in System Preferences > Security & Privacy
- **Path Configuration**: Add Homebrew to your PATH in `~/.zshrc`:
  ```bash
  echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc
  source ~/.zshrc
  ```

---

# 🪟 Windows 11 Setup Guide

## Prerequisites Installation

### 1. Install Windows Subsystem for Linux (WSL2) - Recommended
```powershell
# Run as Administrator in PowerShell
wsl --install -d Ubuntu-22.04
```

### 2. Install Python (Native Windows)
Download Python 3.11+ from [python.org](https://www.python.org/downloads/windows/)
- ✅ Check "Add Python to PATH"
- ✅ Check "Install pip"

### 3. Install Git
Download Git from [git-scm.com](https://git-scm.com/download/win)
- Use recommended settings during installation

### 4. Install AWS CLI v2
Download AWS CLI v2 from [AWS Documentation](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

### 5. Install Microsoft C++ Build Tools
Download from [Microsoft Visual Studio](https://visualstudio.microsoft.com/visual-cpp-build-tools/)
- Select "C++ build tools" workload

## Option A: Native Windows Setup

### 1. Open PowerShell as Administrator
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 2. Clone Repository
```powershell
git clone https://github.com/aws-research-wizard/aws-research-wizard.git
cd aws-research-wizard
```

### 3. Create Virtual Environment
```powershell
python -m venv venv
venv\Scripts\activate
pip install -r requirements.txt
```

### 4. Configure AWS
```powershell
aws configure
```

### 5. Run Applications
```powershell
# Check environment
python aws_environment_checker.py

# Start GUI
streamlit run gui_research_wizard.py
```

## Option B: WSL2 Ubuntu Setup (Recommended)

### 1. Open WSL2 Ubuntu Terminal
```bash
# Update package manager
sudo apt update && sudo apt upgrade -y

# Install Python and pip
sudo apt install python3.11 python3.11-venv python3-pip -y

# Install build essentials
sudo apt install build-essential cmake gfortran git curl -y
```

### 2. Install AWS CLI v2
```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

### 3. Setup Research Wizard
```bash
# Clone repository
git clone https://github.com/aws-research-wizard/aws-research-wizard.git
cd aws-research-wizard

# Create virtual environment
python3.11 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Configure AWS
aws configure

# Run checker
python aws_environment_checker.py
```

## Windows-Specific Notes

- **WSL2 vs Native**: WSL2 provides better compatibility with Linux-based scientific tools
- **File System**: Use WSL2 file system (`/home/username/`) for better performance
- **Windows Terminal**: Install Windows Terminal for better command-line experience
- **Docker Desktop**: Consider installing for container-based workflows

---

# 🐧 Ubuntu Linux Setup Guide

## Prerequisites Installation

### 1. Update Package Manager
```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Install Essential Packages
```bash
# Install Python 3.11 and development tools
sudo apt install python3.11 python3.11-venv python3.11-dev python3-pip -y

# Install build essentials
sudo apt install build-essential cmake gfortran git curl wget -y

# Install additional scientific computing dependencies
sudo apt install libblas-dev liblapack-dev libopenmpi-dev -y
```

### 3. Install AWS CLI v2
```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
rm -rf aws awscliv2.zip
```

### 4. Setup Python Environment
```bash
# Install virtual environment tools
python3.11 -m pip install --user virtualenv pipenv

# Add to PATH (add to ~/.bashrc for persistence)
export PATH="$HOME/.local/bin:$PATH"
```

## Setup AWS Research Wizard

### 1. Clone Repository
```bash
git clone https://github.com/aws-research-wizard/aws-research-wizard.git
cd aws-research-wizard
```

### 2. Create Virtual Environment
```bash
# Create virtual environment
python3.11 -m venv venv
source venv/bin/activate

# Upgrade pip and install requirements
pip install --upgrade pip
pip install -r requirements.txt
```

### 3. Configure AWS Credentials
```bash
aws configure
# Enter your credentials and preferred region
```

### 4. Verify Installation
```bash
# Run environment checker
python aws_environment_checker.py

# Start GUI (will open in browser)
streamlit run gui_research_wizard.py
```

## Ubuntu-Specific Notes

- **Snap Packages**: Some tools are available via snap: `sudo snap install code`
- **Additional Repositories**: Consider adding universe and multiverse repositories
- **Firewall**: Configure UFW if needed: `sudo ufw allow 8501` (for Streamlit)
- **GPU Support**: Install NVIDIA drivers if using GPU instances

---

# 🔧 Advanced Configuration

## Spack Installation (Optional)

For advanced users who want to build packages locally:

### macOS
```bash
git clone -c feature.manyFiles=true https://github.com/spack/spack.git
export SPACK_ROOT=$PWD/spack
export PATH=$SPACK_ROOT/bin:$PATH
```

### Linux/WSL2
```bash
git clone -c feature.manyFiles=true https://github.com/spack/spack.git
export SPACK_ROOT=$PWD/spack
export PATH=$SPACK_ROOT/bin:$PATH
echo 'export SPACK_ROOT=$HOME/spack' >> ~/.bashrc
echo 'export PATH=$SPACK_ROOT/bin:$PATH' >> ~/.bashrc
```

## Docker Configuration (Optional)

### macOS
```bash
brew install docker
# Or download Docker Desktop from docker.com
```

### Windows 11
Download Docker Desktop from [docker.com](https://www.docker.com/products/docker-desktop)

### Ubuntu
```bash
# Install Docker Engine
sudo apt install docker.io docker-compose -y
sudo systemctl enable docker
sudo usermod -aG docker $USER
# Log out and back in for group changes
```

---

# 🔐 Security Configuration

## AWS Credentials Best Practices

### 1. Use IAM Roles (Recommended)
- Create IAM roles for research workloads
- Use temporary credentials when possible
- Avoid hard-coding credentials

### 2. Configure AWS CLI Profiles
```bash
# Create multiple profiles for different environments
aws configure --profile research
aws configure --profile development
aws configure --profile production
```

### 3. Enable MFA
```bash
# Configure MFA for enhanced security
aws configure set mfa_serial arn:aws:iam::ACCOUNT:mfa/USERNAME --profile research
```

### 4. Set Up Credential Files
```bash
# Edit ~/.aws/credentials
[research]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY
region = us-east-1

[development]
aws_access_key_id = YOUR_DEV_ACCESS_KEY
aws_secret_access_key = YOUR_DEV_SECRET_KEY
region = us-west-2
```

---

# 📊 Environment Validation

## Run System Checks

### 1. Basic System Check
```bash
# Python version
python --version

# AWS CLI version
aws --version

# Git version
git --version

# Check available memory
# macOS: system_profiler SPHardwareDataType
# Linux: free -h
# Windows: systeminfo
```

### 2. AWS Environment Check
```bash
# Run comprehensive AWS environment checker
python aws_environment_checker.py --profile research

# Export results for analysis
python aws_environment_checker.py --export environment_check.json
```

### 3. Research Wizard Validation
```bash
# Test import of all modules
python -c "from research_infrastructure_wizard import ResearchInfrastructureWizard; print('✅ Core module working')"

# Test GUI startup
streamlit run gui_research_wizard.py --server.headless true --server.port 8501
```

---

# 🚨 Troubleshooting

## Common Issues and Solutions

### Python/Pip Issues
```bash
# Update pip
python -m pip install --upgrade pip

# Clear pip cache
pip cache purge

# Reinstall virtual environment
rm -rf venv
python -m venv venv
source venv/bin/activate  # Linux/macOS
# venv\Scripts\activate     # Windows
```

### AWS Credential Issues
```bash
# Check current configuration
aws configure list

# Verify credentials
aws sts get-caller-identity

# Test basic AWS access
aws ec2 describe-regions --region us-east-1
```

### Network/Connectivity Issues
```bash
# Test AWS connectivity
curl -I https://aws.amazon.com

# Check DNS resolution
nslookup ec2.amazonaws.com

# Test specific AWS endpoints
aws ec2 describe-availability-zones --region us-east-1
```

### GUI/Streamlit Issues
```bash
# Check if port is in use
# macOS/Linux: netstat -an | grep 8501
# Windows: netstat -an | findstr 8501

# Run on different port
streamlit run gui_research_wizard.py --server.port 8502

# Clear Streamlit cache
streamlit cache clear
```

---

# 📞 Support and Resources

## Documentation
- [AWS Research Wizard Documentation](https://docs.aws-research-wizard.com)
- [AWS CLI Documentation](https://docs.aws.amazon.com/cli/)
- [Spack Documentation](https://spack.readthedocs.io/)

## Community Support
- [GitHub Issues](https://github.com/aws-research-wizard/issues)
- [Community Forum](https://community.aws-research-wizard.com)
- [Stack Overflow](https://stackoverflow.com/questions/tagged/aws-research-wizard)

## Quick Commands Reference

### Start Research Wizard
```bash
# Activate environment
source venv/bin/activate  # Linux/macOS
# venv\Scripts\activate     # Windows

# Run environment check
python aws_environment_checker.py

# Start GUI
streamlit run gui_research_wizard.py
```

### AWS Profile Management
```bash
# List profiles
aws configure list-profiles

# Switch profile
export AWS_PROFILE=research

# Use specific profile
aws s3 ls --profile research
```

### Virtual Environment Management
```bash
# Create new environment
python -m venv research_env

# Activate environment
source research_env/bin/activate  # Linux/macOS
# research_env\Scripts\activate     # Windows

# Deactivate environment
deactivate

# Remove environment
rm -rf research_env
```

---

## 🎉 Ready to Start!

Once you've completed the setup for your operating system, you can begin using the AWS Research Wizard:

1. **Run the environment checker**: `python aws_environment_checker.py`
2. **Start the GUI**: `streamlit run gui_research_wizard.py`
3. **Choose your research domain** and let the wizard optimize your AWS infrastructure!

For additional help, consult the troubleshooting section or reach out to our community support channels.
