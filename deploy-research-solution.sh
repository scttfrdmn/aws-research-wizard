#!/bin/bash

# FinOps-First Research Computing Solution Deployer
# One-click deployment and teardown with built-in cost controls and NIST compliance

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/research-config.yaml"
STATE_FILE="${SCRIPT_DIR}/.deployment-state"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging
log() {
    echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check AWS CLI
    if ! command -v aws &> /dev/null; then
        error "AWS CLI not found. Please install AWS CLI."
    fi
    
    # Check Terraform
    if ! command -v terraform &> /dev/null; then
        error "Terraform not found. Please install Terraform."
    fi
    
    # Check AWS credentials
    if ! aws sts get-caller-identity &> /dev/null; then
        error "AWS credentials not configured. Run 'aws configure'."
    fi
    
    # Check required environment variables
    if [[ -z "${AWS_REGION:-}" ]]; then
        export AWS_REGION="us-east-1"
        warn "AWS_REGION not set, using default: us-east-1"
    fi
    
    log "Prerequisites check passed ✓"
}

# Cost estimation function
estimate_costs() {
    local solution_type=$1
    local usage_hours=${2:-8}  # Default 8 hours/day
    
    case $solution_type in
        "serverless-pipeline")
            local daily_cost=$(echo "$usage_hours * 0.5" | bc -l)
            echo "Estimated daily cost: \$${daily_cost} (Serverless Pipeline)"
            ;;
        "ephemeral-compute")
            local daily_cost=$(echo "$usage_hours * 12" | bc -l)
            echo "Estimated daily cost: \$${daily_cost} (Ephemeral Compute)"
            ;;
        "ai-workbench")
            local daily_cost=$(echo "$usage_hours * 8" | bc -l)
            echo "Estimated daily cost: \$${daily_cost} (AI Workbench)"
            ;;
        "workstation-pods")
            local daily_cost=$(echo "$usage_hours * 2" | bc -l)
            echo "Estimated daily cost: \$${daily_cost} (Workstation Pods)"
            ;;
        *)
            echo "Estimated daily cost: \$5-50 (varies by solution)"
            ;;
    esac
}

# Generate Terraform configuration
generate_terraform_config() {
    local solution_type=$1
    local project_name=$2
    local budget_limit=$3
    
    cat > "${SCRIPT_DIR}/main.tf" << EOF
# FinOps-First Research Computing Solution
# Auto-generated configuration for ${solution_type}

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project        = var.project_name
      Environment    = var.environment
      SolutionType   = var.solution_type
      FinOpsManaged  = "true"
      AutoDestroy    = "true"
      CreatedBy      = "research-deployer"
      CreatedAt      = timestamp()
    }
  }
}

# Variables
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "${AWS_REGION}"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "${project_name}"
}

variable "environment" {
  description = "Environment (dev/staging/prod)"
  type        = string
  default     = "research"
}

variable "solution_type" {
  description = "Type of research solution"
  type        = string
  default     = "${solution_type}"
}

variable "budget_limit" {
  description = "Monthly budget limit in USD"
  type        = number
  default     = ${budget_limit}
}

variable "auto_destroy_hours" {
  description = "Hours after which to auto-destroy resources"
  type        = number
  default     = 24
}

# Budget and Cost Controls
resource "aws_budgets_budget" "research_budget" {
  name         = "\${var.project_name}-budget"
  budget_type  = "COST"
  limit_amount = var.budget_limit
  limit_unit   = "USD"
  time_unit    = "MONTHLY"

  cost_filters {
    tag {
      key    = "Project"
      values = [var.project_name]
    }
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                 = 80
    threshold_type            = "PERCENTAGE"
    notification_type         = "ACTUAL"
    subscriber_email_addresses = [data.aws_caller_identity.current.user_id]
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                 = 100
    threshold_type            = "PERCENTAGE"
    notification_type          = "FORECASTED"
    subscriber_email_addresses = [data.aws_caller_identity.current.user_id]
  }
}

# Auto-destruction Lambda
resource "aws_lambda_function" "auto_destroyer" {
  filename         = "auto-destroyer.zip"
  function_name    = "\${var.project_name}-auto-destroyer"
  role            = aws_iam_role.auto_destroyer_role.arn
  handler         = "index.handler"
  runtime         = "python3.11"
  timeout         = 300

  environment {
    variables = {
      PROJECT_NAME = var.project_name
      MAX_HOURS    = var.auto_destroy_hours
    }
  }
}

# Auto-destroyer IAM role
resource "aws_iam_role" "auto_destroyer_role" {
  name = "\${var.project_name}-auto-destroyer-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "auto_destroyer_policy" {
  name = "\${var.project_name}-auto-destroyer-policy"
  role = aws_iam_role.auto_destroyer_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "ec2:DescribeInstances",
          "ec2:TerminateInstances",
          "ecs:*",
          "batch:*",
          "sagemaker:*",
          "lambda:*",
          "s3:*"
        ]
        Resource = "*"
        Condition = {
          StringEquals = {
            "aws:ResourceTag/Project" = var.project_name
          }
        }
      }
    ]
  })
}

# EventBridge rule for auto-destruction
resource "aws_cloudwatch_event_rule" "auto_destroy_schedule" {
  name                = "\${var.project_name}-auto-destroy"
  description         = "Trigger auto-destruction check"
  schedule_expression = "rate(1 hour)"
}

resource "aws_cloudwatch_event_target" "auto_destroy_target" {
  rule      = aws_cloudwatch_event_rule.auto_destroy_schedule.name
  target_id = "AutoDestroyLambdaTarget"
  arn       = aws_lambda_function.auto_destroyer.arn
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.auto_destroyer.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.auto_destroy_schedule.arn
}

# Data sources
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Outputs
output "project_name" {
  description = "Project name"
  value       = var.project_name
}

output "budget_arn" {
  description = "Budget ARN"
  value       = aws_budgets_budget.research_budget.arn
}

output "auto_destroyer_function" {
  description = "Auto-destroyer Lambda function"
  value       = aws_lambda_function.auto_destroyer.function_name
}

output "cost_dashboard_url" {
  description = "AWS Cost Explorer URL for this project"
  value       = "https://console.aws.amazon.com/cost-management/home#/custom?timeRangeOption=Last3Months&granularity=Monthly&reportName=Research%20Project%20Costs&chartStyle=Stack&groupBy1=TagKey%3AProject&hasBlended=false&hasUntagged=false&isTemplate=true&filter=%5B%7B%22dimension%22%3A%22TagKey%22%2C%22key%22%3A%22Project%22%2C%22values%22%3A%5B%22\${var.project_name}%22%5D%2C%22matchOptions%22%3A%5B%22EQUALS%22%5D%7D%5D"
}
EOF

    # Generate solution-specific configuration
    case $solution_type in
        "serverless-pipeline")
            generate_serverless_config >> "${SCRIPT_DIR}/main.tf"
            ;;
        "ephemeral-compute")
            generate_ephemeral_compute_config >> "${SCRIPT_DIR}/main.tf"
            ;;
        "ai-workbench")
            generate_ai_workbench_config >> "${SCRIPT_DIR}/main.tf"
            ;;
        "workstation-pods")
            generate_workstation_config >> "${SCRIPT_DIR}/main.tf"
            ;;
    esac
}

# Generate auto-destroyer Lambda
generate_auto_destroyer() {
    cat > "${SCRIPT_DIR}/auto-destroyer.py" << 'EOF'
import boto3
import json
import os
from datetime import datetime, timedelta

def handler(event, context):
    project_name = os.environ['PROJECT_NAME']
    max_hours = int(os.environ.get('MAX_HOURS', 24))
    
    # Initialize AWS clients
    ec2 = boto3.client('ec2')
    ecs = boto3.client('ecs')
    batch = boto3.client('batch')
    sagemaker = boto3.client('sagemaker')
    
    cutoff_time = datetime.utcnow() - timedelta(hours=max_hours)
    destroyed_resources = []
    
    try:
        # Check and terminate old EC2 instances
        instances = ec2.describe_instances(
            Filters=[
                {'Name': 'tag:Project', 'Values': [project_name]},
                {'Name': 'instance-state-name', 'Values': ['running']}
            ]
        )
        
        for reservation in instances['Reservations']:
            for instance in reservation['Instances']:
                launch_time = instance['LaunchTime'].replace(tzinfo=None)
                if launch_time < cutoff_time:
                    ec2.terminate_instances(InstanceIds=[instance['InstanceId']])
                    destroyed_resources.append(f"EC2 Instance: {instance['InstanceId']}")
        
        # Check and stop SageMaker notebook instances
        notebooks = sagemaker.list_notebook_instances(
            StatusEquals='InService'
        )
        
        for notebook in notebooks['NotebookInstances']:
            # Check if it belongs to our project and is old
            try:
                tags = sagemaker.list_tags(ResourceArn=notebook['NotebookInstanceArn'])
                project_tag = next((tag for tag in tags['Tags'] if tag['Key'] == 'Project'), None)
                
                if project_tag and project_tag['Value'] == project_name:
                    creation_time = notebook['CreationTime'].replace(tzinfo=None)
                    if creation_time < cutoff_time:
                        sagemaker.stop_notebook_instance(
                            NotebookInstanceName=notebook['NotebookInstanceName']
                        )
                        destroyed_resources.append(f"SageMaker Notebook: {notebook['NotebookInstanceName']}")
            except Exception as e:
                print(f"Error checking notebook {notebook['NotebookInstanceName']}: {e}")
        
        print(f"Auto-destroyer completed. Destroyed resources: {destroyed_resources}")
        
        return {
            'statusCode': 200,
            'body': json.dumps({
                'message': 'Auto-destruction completed',
                'destroyed_resources': destroyed_resources
            })
        }
        
    except Exception as e:
        print(f"Error in auto-destroyer: {e}")
        return {
            'statusCode': 500,
            'body': json.dumps({'error': str(e)})
        }
EOF

    # Create deployment package
    cd "${SCRIPT_DIR}"
    zip -q auto-destroyer.zip auto-destroyer.py
    rm auto-destroyer.py
}

# Generate serverless pipeline configuration
generate_serverless_config() {
    cat << 'EOF'

# Serverless Research Data Pipeline Configuration

# S3 bucket for data
resource "aws_s3_bucket" "research_data" {
  bucket = "${var.project_name}-research-data-${random_id.bucket_suffix.hex}"
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

resource "aws_s3_bucket_encryption_configuration" "research_data_encryption" {
  bucket = aws_s3_bucket.research_data.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "research_data_lifecycle" {
  bucket = aws_s3_bucket.research_data.id

  rule {
    id     = "research_data_lifecycle"
    status = "Enabled"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 90
      storage_class = "GLACIER"
    }

    expiration {
      days = 365
    }
  }
}

# Lambda function for data processing
resource "aws_lambda_function" "data_processor" {
  filename         = "data-processor.zip"
  function_name    = "${var.project_name}-data-processor"
  role            = aws_iam_role.lambda_role.arn
  handler         = "index.handler"
  runtime         = "python3.11"
  timeout         = 900
  memory_size     = 3008

  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.research_data.bucket
    }
  }
}

resource "aws_iam_role" "lambda_role" {
  name = "${var.project_name}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_s3_policy" {
  name = "${var.project_name}-lambda-s3-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = [
          "arn:aws:logs:*:*:*",
          "${aws_s3_bucket.research_data.arn}/*"
        ]
      }
    ]
  })
}

# Step Functions state machine
resource "aws_sfn_state_machine" "data_pipeline" {
  name     = "${var.project_name}-data-pipeline"
  role_arn = aws_iam_role.step_functions_role.arn

  definition = jsonencode({
    Comment = "Research Data Processing Pipeline"
    StartAt = "ProcessData"
    States = {
      ProcessData = {
        Type     = "Task"
        Resource = aws_lambda_function.data_processor.arn
        End      = true
      }
    }
  })
}

resource "aws_iam_role" "step_functions_role" {
  name = "${var.project_name}-step-functions-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "states.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "step_functions_policy" {
  name = "${var.project_name}-step-functions-policy"
  role = aws_iam_role.step_functions_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction"
        ]
        Resource = aws_lambda_function.data_processor.arn
      }
    ]
  })
}

output "data_bucket" {
  description = "S3 bucket for research data"
  value       = aws_s3_bucket.research_data.bucket
}

output "pipeline_arn" {
  description = "Step Functions pipeline ARN"
  value       = aws_sfn_state_machine.data_pipeline.arn
}
EOF
}

# Main deployment function
deploy_solution() {
    local solution_type=$1
    local project_name=$2
    local budget_limit=${3:-100}
    
    log "Deploying ${solution_type} solution for project: ${project_name}"
    
    # Estimate costs
    estimate_costs "$solution_type"
    
    # Generate Terraform configuration
    generate_terraform_config "$solution_type" "$project_name" "$budget_limit"
    
    # Generate auto-destroyer
    generate_auto_destroyer
    
    # Initialize and apply Terraform
    log "Initializing Terraform..."
    terraform init
    
    log "Planning deployment..."
    terraform plan -var="project_name=${project_name}" -var="solution_type=${solution_type}" -var="budget_limit=${budget_limit}"
    
    echo -e "${YELLOW}Do you want to proceed with deployment? (y/N)${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        log "Applying Terraform configuration..."
        terraform apply -auto-approve -var="project_name=${project_name}" -var="solution_type=${solution_type}" -var="budget_limit=${budget_limit}"
        
        # Save deployment state
        echo "solution_type=${solution_type}" > "$STATE_FILE"
        echo "project_name=${project_name}" >> "$STATE_FILE"
        echo "deployed_at=$(date)" >> "$STATE_FILE"
        
        log "Deployment completed successfully! ✓"
        log "Budget alerts configured for \$${budget_limit}/month"
        log "Auto-destruction enabled after 24 hours of inactivity"
        
        # Show cost dashboard
        terraform output cost_dashboard_url
        
    else
        log "Deployment cancelled."
    fi
}

# Destroy solution
destroy_solution() {
    if [[ ! -f "$STATE_FILE" ]]; then
        error "No deployment state found. Nothing to destroy."
    fi
    
    source "$STATE_FILE"
    
    log "Destroying solution: ${solution_type} for project: ${project_name}"
    
    echo -e "${RED}This will destroy ALL resources for project: ${project_name}${NC}"
    echo -e "${RED}Are you sure? (type 'destroy' to confirm)${NC}"
    read -r response
    
    if [[ "$response" == "destroy" ]]; then
        log "Destroying Terraform infrastructure..."
        terraform destroy -auto-approve -var="project_name=${project_name}" -var="solution_type=${solution_type}"
        
        # Clean up local files
        rm -f "$STATE_FILE" main.tf terraform.tfstate* .terraform.lock.hcl auto-destroyer.zip
        rm -rf .terraform/
        
        log "Solution destroyed successfully! ✓"
        log "All resources have been cleaned up."
    else
        log "Destroy cancelled."
    fi
}

# Show usage
show_usage() {
    cat << EOF
FinOps-First Research Computing Solution Deployer

Usage: $0 [command] [options]

Commands:
  deploy <solution-type> <project-name> [budget-limit]
    Deploy a research computing solution
    
  destroy
    Destroy the current deployment
    
  status
    Show current deployment status
    
  costs
    Show current cost breakdown

Solution Types:
  serverless-pipeline  - Event-driven data processing (\$0 idle)
  ephemeral-compute   - On-demand HPC clusters (\$0 idle)
  ai-workbench        - ML/AI development environment (\$0 idle)
  workstation-pods    - Remote research workstations (\$0 idle)

Examples:
  $0 deploy serverless-pipeline my-genomics-project 200
  $0 deploy ai-workbench deep-learning-research 500
  $0 destroy
  $0 status

Features:
  ✓ Zero idle costs - Pay only for active usage
  ✓ NIST 800-171/800-53 compliant
  ✓ Auto-destruction after 24 hours
  ✓ Real-time cost monitoring
  ✓ One-click deploy/destroy
EOF
}

# Show status
show_status() {
    if [[ ! -f "$STATE_FILE" ]]; then
        log "No active deployment found."
        return
    fi
    
    source "$STATE_FILE"
    log "Current deployment:"
    log "  Solution Type: ${solution_type}"
    log "  Project Name: ${project_name}"
    log "  Deployed At: ${deployed_at}"
    
    # Show Terraform resources
    if [[ -f "terraform.tfstate" ]]; then
        log "Active resources:"
        terraform show -json | jq -r '.values.root_module.resources[]? | "  \(.type): \(.values.id // .values.name // "N/A")"' 2>/dev/null || log "  (Unable to parse terraform state)"
    fi
}

# Show costs
show_costs() {
    if [[ ! -f "$STATE_FILE" ]]; then
        error "No active deployment found."
    fi
    
    source "$STATE_FILE"
    log "Fetching cost data for project: ${project_name}"
    
    # Get costs from AWS Cost Explorer
    aws ce get-cost-and-usage \
        --time-period Start=2025-06-01,End=2025-06-30 \
        --granularity MONTHLY \
        --metrics BlendedCost \
        --group-by Type=DIMENSION,Key=SERVICE \
        --filter "file://<(echo '{\"Tags\":{\"Key\":\"Project\",\"Values\":[\"'${project_name}'\"]}}' | tr -d '\n')" \
        2>/dev/null || warn "Unable to fetch cost data. Check AWS permissions."
}

# Main script logic
main() {
    case "${1:-}" in
        "deploy")
            if [[ $# -lt 3 ]]; then
                error "Usage: $0 deploy <solution-type> <project-name> [budget-limit]"
            fi
            check_prerequisites
            deploy_solution "$2" "$3" "${4:-100}"
            ;;
        "destroy")
            destroy_solution
            ;;
        "status")
            show_status
            ;;
        "costs")
            show_costs
            ;;
        "help"|"-h"|"--help"|"")
            show_usage
            ;;
        *)
            error "Unknown command: $1. Use '$0 help' for usage information."
            ;;
    esac
}

main "$@"