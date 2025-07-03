#!/bin/bash
# Research Environment Initialization Script
# Domain: ${domain_name}
# Environment: ${environment}

set -euo pipefail

# Logging setup
exec > >(tee /var/log/research-setup.log)
exec 2>&1

echo "=== Research Environment Setup Started ==="
echo "Domain: ${domain_name}"
echo "Environment: ${environment}"
echo "Timestamp: $(date)"

# Update system
echo "Updating system packages..."
yum update -y

# Install essential tools
echo "Installing essential tools..."
yum install -y \
    git \
    wget \
    curl \
    htop \
    python3 \
    python3-pip \
    gcc \
    gcc-c++ \
    make \
    cmake \
    tar \
    gzip

# Install Spack package manager
echo "Installing Spack package manager..."
cd /opt
git clone -c feature.manyFiles=true https://github.com/spack/spack.git
chown -R ec2-user:ec2-user /opt/spack

# Set up Spack environment for ec2-user
echo "Setting up Spack environment..."
cat >> /home/ec2-user/.bashrc << 'EOF'
# Spack setup
export SPACK_ROOT=/opt/spack
source $SPACK_ROOT/share/spack/setup-env.sh
EOF

# Create research directory structure
echo "Creating research directory structure..."
mkdir -p /home/ec2-user/research/{data,scripts,results,logs}
chown -R ec2-user:ec2-user /home/ec2-user/research

# Install domain-specific packages (if any)
if [ ! -z "${spack_packages}" ]; then
    echo "Installing Spack packages: ${spack_packages}"
    # Note: This will be implemented in the next iteration
    # su - ec2-user -c "source /opt/spack/share/spack/setup-env.sh && spack install ${spack_packages}"
fi

# Set up CloudWatch agent (for real AWS)
if [ "${environment}" != "localstack" ]; then
    echo "Setting up CloudWatch agent..."
    wget https://s3.amazonaws.com/amazoncloudwatch-agent/amazon_linux/amd64/latest/amazon-cloudwatch-agent.rpm
    rpm -U ./amazon-cloudwatch-agent.rpm
fi

# Create completion marker
echo "Setup completed at $(date)" > /home/ec2-user/research/setup-complete.txt
chown ec2-user:ec2-user /home/ec2-user/research/setup-complete.txt

echo "=== Research Environment Setup Completed ==="
