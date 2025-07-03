# Real AWS Research Environment for Testing
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = "aws"
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

data "aws_availability_zones" "available" {
  state = "available"
}

# VPC and networking
resource "aws_vpc" "research" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name        = "research-${var.domain_name}-vpc"
    Domain      = var.domain_name
    Environment = "aws"
    Testing     = "true"
  }
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.research.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true

  tags = {
    Name        = "research-${var.domain_name}-public"
    Domain      = var.domain_name
    Environment = "aws"
    Testing     = "true"
  }
}

resource "aws_internet_gateway" "research" {
  vpc_id = aws_vpc.research.id

  tags = {
    Name        = "research-${var.domain_name}-igw"
    Domain      = var.domain_name
    Environment = "aws"
    Testing     = "true"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.research.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.research.id
  }

  tags = {
    Name        = "research-${var.domain_name}-public-rt"
    Domain      = var.domain_name
    Environment = "aws"
    Testing     = "true"
  }
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

# Security group for research environment
resource "aws_security_group" "research" {
  name_prefix = "research-${var.domain_name}-"
  description = "Security group for ${var.domain_name} research environment"
  vpc_id      = aws_vpc.research.id

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
    Environment = "aws"
    Testing     = "true"
  }
}

# Key pair for SSH access
resource "aws_key_pair" "research" {
  key_name   = "research-${var.domain_name}-key"
  public_key = file("~/.ssh/id_rsa.pub")

  tags = {
    Name        = "research-${var.domain_name}-key"
    Domain      = var.domain_name
    Environment = "aws"
    Testing     = "true"
  }
}

# User data script for instance initialization
locals {
  user_data = base64encode(templatefile("../../modules/research-environment/user_data.sh", {
    domain_name    = var.domain_name
    spack_packages = join(" ", var.spack_packages)
    environment    = "aws"
  }))
}

# EC2 instance for research environment
resource "aws_instance" "research_node" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = var.instance_type
  key_name      = aws_key_pair.research.key_name
  subnet_id     = aws_subnet.public.id

  vpc_security_group_ids = [aws_security_group.research.id]

  user_data = local.user_data

  tags = {
    Name        = "research-${var.domain_name}-node"
    Domain      = var.domain_name
    Budget      = var.monthly_budget
    Environment = "aws"
    Testing     = "true"
  }

  lifecycle {
    create_before_destroy = true
  }
}

# S3 bucket for research data
resource "aws_s3_bucket" "research_data" {
  bucket = "research-${replace(var.domain_name, "_", "-")}-${random_id.bucket_suffix.hex}"

  tags = {
    Name        = "research-${var.domain_name}-data"
    Purpose     = "Research data storage"
    Domain      = var.domain_name
    Environment = "aws"
    Testing     = "true"
  }
}

resource "aws_s3_bucket_versioning" "research_data" {
  bucket = aws_s3_bucket.research_data.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "research_data" {
  bucket = aws_s3_bucket.research_data.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# CloudWatch log group for monitoring
resource "aws_cloudwatch_log_group" "research_logs" {
  name              = "/aws/research/${var.domain_name}"
  retention_in_days = 7

  tags = {
    Name        = "research-${var.domain_name}-logs"
    Domain      = var.domain_name
    Environment = "aws"
    Testing     = "true"
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

output "s3_bucket" {
  description = "S3 bucket for research data"
  value       = aws_s3_bucket.research_data.bucket
}

output "log_group" {
  description = "CloudWatch log group"
  value       = aws_cloudwatch_log_group.research_logs.name
}

output "ssh_command" {
  description = "SSH command to connect to the instance"
  value       = "ssh -i ~/.ssh/id_rsa ec2-user@${aws_instance.research_node.public_ip}"
}

output "estimated_cost" {
  description = "Estimated monthly cost in USD"
  value       = var.monthly_budget
}

output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.research.id
}

output "subnet_id" {
  description = "Public subnet ID"
  value       = aws_subnet.public.id
}
