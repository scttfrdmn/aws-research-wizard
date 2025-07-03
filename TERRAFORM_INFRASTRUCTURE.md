# Terraform Infrastructure Implementation

## Overview

The AWS Research Wizard now supports Terraform-based infrastructure deployment as a modern alternative to CloudFormation. This implementation provides both LocalStack Pro development testing and real AWS deployment capabilities.

## Architecture

### Infrastructure Components

#### Core Resources
- **VPC**: Dedicated virtual private cloud with DNS support
- **Public Subnet**: Internet-accessible subnet in first AZ
- **Internet Gateway**: Public internet access
- **Route Tables**: Routing configuration for public subnet
- **Security Group**: SSH (22) and web interfaces (8080-8090) access
- **EC2 Instance**: Research computing node with user data initialization
- **Key Pair**: SSH access using generated RSA keypair
- **S3 Bucket**: Research data storage with versioning and encryption
- **CloudWatch Log Group**: Centralized logging with 7-day retention

#### Research Environment Setup
- **Spack Package Manager**: Installed to `/opt/spack`
- **Research Directories**: Created at `/home/ec2-user/research/`
- **Domain-Specific Packages**: Configurable via Terraform variables
- **CloudWatch Agent**: Monitoring and logging (real AWS only)
- **User Data Script**: Automated environment initialization

## Directory Structure

```
terraform/
├── modules/
│   └── research-environment/
│       ├── main.tf           # Modular infrastructure components
│       ├── variables.tf      # Module input variables
│       ├── outputs.tf        # Module outputs
│       └── user_data.sh      # Instance initialization script
└── environments/
    ├── localstack/
    │   ├── main.tf           # LocalStack Pro configuration
    │   └── variables.tf      # LocalStack variables
    └── aws/
        ├── main.tf           # Real AWS configuration
        └── variables.tf      # AWS variables
```

## Deployment Environments

### LocalStack Pro Development

**Purpose**: Local development and testing without AWS costs

**Configuration**:
- Provider endpoints pointed to `localhost:4566`
- Simplified networking (uses LocalStack defaults)
- No S3/CloudWatch for faster testing
- Auth token: `ls-dIgo5507-demE-zaWa-8414-VEnUSEre3d8d`

**Usage**:
```bash
cd terraform/environments/localstack
terraform init
terraform plan
terraform apply
```

**Limitations**:
- S3 operations may hang (LocalStack limitation)
- CloudWatch logs not fully functional
- EC2 networking simplified

### Real AWS Production

**Purpose**: Production deployments with full AWS feature set

**Configuration**:
- Complete VPC with public/private networking
- Full S3 bucket with versioning and encryption
- CloudWatch monitoring and logging
- AWS profile: `aws`
- Region: `us-east-1`

**Usage**:
```bash
cd terraform/environments/aws
terraform init
terraform plan
terraform apply
# Clean up when done:
terraform destroy
```

## Test Results

### LocalStack Pro Deployment ✅

- **Instance Created**: `i-2ad12bb91a51e1bb9`
- **IP Address**: Public: `54.214.173.14`, Private: `10.36.219.232`
- **Security Group**: `sg-7a8a1234712e3c214`
- **Domain**: `digital_humanities` with NLP packages
- **Status**: Deployment and destruction successful

### Real AWS Deployment ✅

- **Instance Created**: `i-0316ba6d35fa79c61` (t3.small)
- **IP Address**: Public: `54.166.226.176`, Private: `10.0.1.11`
- **VPC**: `vpc-048c56864d1d33f6a` with complete networking
- **S3 Bucket**: `research-digital-humanities-ea84b05f`
- **Spack**: Version 1.0.0.dev0 with 8,487 packages
- **SSH Access**: Confirmed working with generated keypair
- **Research Environment**: Fully operational with proper directory structure

## Domain Configuration Testing

### Digital Humanities Domain

**Packages Tested**:
- `python@3.11.5` - Core Python environment
- `py-nltk@3.8.1` - Natural Language Toolkit
- `py-spacy@3.6.1` - Industrial NLP library
- `py-pandas@2.0.3` - Data analysis framework
- `git@2.41.0` - Version control

**Budget**: $750/month as specified in domain config
**Instance Type**: t3.small for development testing
**Status**: Successfully deployed and operational

## Key Improvements over CloudFormation

### Development Workflow
1. **LocalStack Integration**: Test infrastructure changes locally without AWS costs
2. **Faster Iteration**: Simplified LocalStack config for rapid development
3. **Cost Control**: Easy cleanup with `terraform destroy`

### Infrastructure Management
1. **Modular Design**: Reusable modules for consistent deployments
2. **State Management**: Terraform state tracking for change management
3. **Variable Configuration**: Environment-specific settings
4. **Output Values**: Structured access to resource information

### Multi-Cloud Readiness
1. **Provider Abstraction**: Easy to extend to other cloud providers
2. **Environment Separation**: Clear development vs production separation
3. **Configuration Management**: Centralized variable management

## Issues Resolved

### S3 Bucket Naming
- **Problem**: Underscores not allowed in S3 bucket names
- **Solution**: Used `replace(var.domain_name, "_", "-")` for compliant naming

### Versioned Object Cleanup
- **Problem**: S3 buckets with versioning couldn't be deleted via Terraform
- **Solution**: Manual cleanup of all versions before destruction

### SSH Key Management
- **Problem**: No existing SSH keypair for AWS access
- **Solution**: Generated RSA keypair and configured in Terraform

## Security Considerations

### Network Security
- VPC with controlled subnet access
- Security groups with minimal required ports (22, 8080-8090)
- Public subnet only for necessary internet access

### Data Security
- S3 bucket encryption with AES256
- Object versioning enabled for data protection
- CloudWatch logs for audit trail

### Access Control
- SSH key-based authentication
- AWS profile-based credential management
- Resource tagging for environment identification

## Cost Optimization

### Resource Tagging
All resources tagged with:
- `Domain`: Research domain name
- `Environment`: `localstack` or `aws`
- `Testing`: `true` for easy identification
- `Budget`: Monthly budget tracking

### Instance Sizing
- Development: `t3.small` ($0.0208/hour)
- Production: Configurable based on domain requirements
- Automatic termination via Terraform destroy

## Next Steps

1. **Integration with Go CLI**: Update research wizard to use Terraform instead of CloudFormation
2. **Multi-Domain Testing**: Deploy all 27 research domains systematically
3. **Production Modules**: Create production-ready modules with enhanced security
4. **CI/CD Integration**: Automate testing and deployment workflows
5. **Cost Monitoring**: Implement automated cost tracking and alerts

## Troubleshooting

### Common Issues

**LocalStack S3 Hanging**:
- Remove S3 resources from LocalStack config for faster testing
- Use simplified infrastructure for development

**SSH Connection Refused**:
- Verify security group allows port 22 from your IP
- Check instance status and user data completion

**Terraform State Issues**:
- Use `terraform refresh` to sync state with actual resources
- Manual cleanup may be required for S3 buckets with versioning

### Debug Commands

```bash
# Check instance status
aws ec2 describe-instances --instance-ids i-XXXXX --profile aws

# Monitor user data execution
ssh -i ~/.ssh/id_rsa ec2-user@IP "tail -f /var/log/research-setup.log"

# Verify Spack installation
ssh -i ~/.ssh/id_rsa ec2-user@IP "source /opt/spack/share/spack/setup-env.sh && spack --version"

# Check S3 bucket
aws s3 ls s3://bucket-name --profile aws
```

This Terraform infrastructure provides a solid foundation for deploying research environments with both development testing capabilities and production-ready AWS deployments.
