# Research Environment Module
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
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

  tags = merge(var.tags, {
    Name = "research-${var.domain_name}-sg"
  })
}

# User data script for instance initialization
locals {
  user_data = base64encode(templatefile("${path.module}/user_data.sh", {
    domain_name    = var.domain_name
    spack_packages = join(" ", var.spack_packages)
    environment    = var.environment
  }))
}

# EC2 instance for research environment
resource "aws_instance" "research_node" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = var.instance_type

  vpc_security_group_ids = [aws_security_group.research.id]

  user_data = local.user_data

  tags = merge(var.tags, {
    Name   = "research-${var.domain_name}-node"
    Domain = var.domain_name
    Budget = var.monthly_budget
  })

  lifecycle {
    create_before_destroy = true
  }
}

# S3 bucket for research data
resource "aws_s3_bucket" "research_data" {
  bucket = "research-${var.domain_name}-${random_id.bucket_suffix.hex}"

  tags = merge(var.tags, {
    Name    = "research-${var.domain_name}-data"
    Purpose = "Research data storage"
  })
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# CloudWatch log group for monitoring
resource "aws_cloudwatch_log_group" "research_logs" {
  name              = "/aws/research/${var.domain_name}"
  retention_in_days = 7

  tags = merge(var.tags, {
    Name = "research-${var.domain_name}-logs"
  })
}
