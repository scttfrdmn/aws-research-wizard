# Simplified Research Environment for LocalStack Testing
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region                      = "us-east-1"
  access_key                  = "test"
  secret_key                  = "test"
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    ec2 = "http://localhost:4566"
  }
}

# Data sources
data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# Security group for research environment
resource "aws_security_group" "research" {
  name_prefix = "research-${var.domain_name}-"
  description = "Security group for ${var.domain_name} research environment"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH access"
  }

  ingress {
    from_port   = 8080
    to_port     = 8090
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Web interfaces"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }

  tags = {
    Name        = "research-${var.domain_name}-sg"
    Domain      = var.domain_name
    Environment = "localstack"
    Testing     = "true"
  }
}

# User data script for instance initialization
locals {
  user_data = base64encode(templatefile("../../modules/research-environment/user_data.sh", {
    domain_name    = var.domain_name
    spack_packages = join(" ", var.spack_packages)
    environment    = "localstack"
  }))
}

# EC2 instance for research environment
resource "aws_instance" "research_node" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = var.instance_type

  vpc_security_group_ids = [aws_security_group.research.id]

  user_data = local.user_data

  tags = {
    Name        = "research-${var.domain_name}-node"
    Domain      = var.domain_name
    Budget      = var.monthly_budget
    Environment = "localstack"
    Testing     = "true"
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Outputs
output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.research_node.id
}

output "public_ip" {
  description = "Public IP address of the instance"
  value       = aws_instance.research_node.public_ip
}

output "private_ip" {
  description = "Private IP address of the instance"
  value       = aws_instance.research_node.private_ip
}

output "security_group_id" {
  description = "ID of the security group"
  value       = aws_security_group.research.id
}

output "estimated_cost" {
  description = "Estimated monthly cost in USD"
  value       = var.monthly_budget
}
