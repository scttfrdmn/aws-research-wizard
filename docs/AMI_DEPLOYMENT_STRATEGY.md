# AMI Deployment Strategy - Research Environment Snapshots

## Overview

Enable users to create, share, and deploy pre-configured research environments as Amazon Machine Images (AMIs), eliminating the 20-30 minute setup time and providing guaranteed consistency across deployments.

## Problem Statement

Current deployment process requires:
- 20-30 minutes for package installation (54 packages in genomics)
- Complex dependency resolution with potential failures
- Network connectivity for all package downloads
- Repeated compilation of identical software stacks

## Solution: AMI-Based Deployment

### Core Features

#### 1. AMI Creation & Management
```bash
# Create AMI from successful deployment
aws-research-wizard deploy save-ami --stack genomics-env --name "genomics-v1.0" --description "Genomics pack with 54 bioinformatics tools"

# Deploy from existing AMI
aws-research-wizard deploy --from-ami ami-12345678 --instance r5.large

# List available AMIs
aws-research-wizard ami list --domain genomics
```

#### 2. AMI Sharing & Distribution
```bash
# Share AMI with specific AWS account
aws-research-wizard ami share --ami-id ami-12345678 --account-id 123456789012

# Make AMI public (for community sharing)
aws-research-wizard ami publish --ami-id ami-12345678 --public

# Copy AMI to multiple regions
aws-research-wizard ami distribute --ami-id ami-12345678 --regions us-west-2,us-east-1,eu-west-1
```

#### 3. Community AMI Registry
```bash
# Submit AMI to community registry
aws-research-wizard ami submit --ami-id ami-12345678 --registry community

# Search community AMIs
aws-research-wizard ami search --domain genomics --tags "covid-19,variant-calling"

# Install from community registry
aws-research-wizard deploy --from-registry community/genomics-covid-v2.1
```

## Technical Implementation

### AMI Creation Pipeline

#### 1. Post-Deployment AMI Creation
```yaml
# Extension to terraform/environments/aws/main.tf
resource "aws_ami_from_instance" "research_ami" {
  count               = var.create_ami ? 1 : 0
  name                = "${var.domain_name}-${var.ami_version}-${timestamp()}"
  source_instance_id  = aws_instance.research_node.id

  # Wait for setup completion
  depends_on = [null_resource.setup_complete]

  tags = {
    DomainPack = var.domain_name
    Version    = var.ami_version
    Created    = timestamp()
    PackageCount = length(var.spack_packages)
  }
}
```

#### 2. AMI Optimization Script
```bash
#!/bin/bash
# ami-cleanup.sh - Prepare instance for AMI creation

# Remove sensitive data
rm -rf /home/ec2-user/.ssh/*
rm -rf /root/.ssh/*
rm -f /var/log/research-setup.log

# Clean package caches
rm -rf /opt/spack/var/spack/cache/*
conda clean -a -y
yum clean all

# Remove bash history
rm -f /home/ec2-user/.bash_history
rm -f /root/.bash_history

# Clear system logs
truncate -s 0 /var/log/messages
truncate -s 0 /var/log/secure

# Create AMI metadata
cat > /home/ec2-user/ami-info.json << EOF
{
  "domain": "${domain_name}",
  "version": "${ami_version}",
  "packages": ${package_count},
  "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "instance_type": "${instance_type}",
  "region": "${aws_region}"
}
EOF
```

### AMI Registry & Sharing

#### 1. Community AMI Registry
```go
// internal/ami/registry.go
type AMIRegistry struct {
    PublicAMIs    map[string]AMIMetadata
    CommunityAMIs map[string]AMIMetadata
    UserAMIs      map[string]AMIMetadata
}

type AMIMetadata struct {
    AMIID          string            `json:"ami_id"`
    Name           string            `json:"name"`
    Domain         string            `json:"domain"`
    Version        string            `json:"version"`
    Description    string            `json:"description"`
    Tags           []string          `json:"tags"`
    PackageCount   int               `json:"package_count"`
    InstanceTypes  []string          `json:"supported_instances"`
    Regions        []string          `json:"available_regions"`
    CreatedBy      string            `json:"created_by"`
    CreatedAt      time.Time         `json:"created_at"`
    Downloads      int               `json:"downloads"`
    Rating         float64           `json:"rating"`
    TestResults    TestResults       `json:"test_results"`
}
```

#### 2. AMI Sharing Infrastructure
```yaml
# AMI sharing configuration
sharing_config:
  auto_share_regions:
    - us-west-2
    - us-east-1
    - eu-west-1
    - ap-northeast-1

  community_registry:
    enabled: true
    approval_required: true
    moderation_rules:
      - security_scan: true
      - package_verification: true
      - performance_benchmark: true

  public_sharing:
    enabled: true
    requires_approval: true
    cost_optimization: true
```

### Integration with Existing System

#### 1. Extended Terraform Configuration
```hcl
# terraform/environments/aws/variables.tf
variable "create_ami" {
  description = "Create AMI after successful deployment"
  type        = bool
  default     = false
}

variable "ami_version" {
  description = "Version tag for created AMI"
  type        = string
  default     = "1.0"
}

variable "source_ami_id" {
  description = "Use existing AMI instead of building from scratch"
  type        = string
  default     = ""
}
```

#### 2. CLI Command Extensions
```go
// internal/commands/ami/ami.go
func NewAMICommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "ami",
        Short: "Manage research environment AMIs",
        Long:  "Create, share, and deploy AMIs for research environments",
    }

    cmd.AddCommand(
        createCreateCommand(),
        createListCommand(),
        createShareCommand(),
        createPublishCommand(),
        createSearchCommand(),
        createDeployCommand(),
    )

    return cmd
}
```

## AMI Sharing Ecosystem

### 1. Official AWS Research Wizard AMI Gallery
- **Curated AMIs** - Officially maintained and tested
- **Version management** - Clear versioning and changelog
- **Multi-region availability** - Distributed for global access
- **Security updates** - Regular patching and updates

### 2. Community AMI Marketplace
- **User-contributed AMIs** - Community-created specialized environments
- **Rating system** - User feedback and quality metrics
- **Tagging system** - Searchable by research area, tools, datasets
- **Usage analytics** - Download counts, performance metrics

### 3. Educational AMI Library
- **Workshop AMIs** - Pre-configured for tutorials and classes
- **Course-specific** - Tailored for specific academic programs
- **Reproducible research** - Exact environments for paper reproduction
- **Collaborative projects** - Shared environments for team research

## Security & Compliance

### 1. AMI Security Scanning
```bash
# Automated security scan before sharing
aws-research-wizard ami scan --ami-id ami-12345678
# - Vulnerability assessment
# - Malware detection
# - Compliance checking
# - License validation
```

### 2. Access Control
```yaml
# AMI access control policy
access_control:
  private_amis:
    - owner_only: true
    - explicit_sharing: true

  community_amis:
    - moderation_required: true
    - security_scan: true
    - performance_validation: true

  public_amis:
    - approval_workflow: true
    - regular_security_updates: true
    - usage_monitoring: true
```

## Cost Management

### 1. AMI Storage Optimization
- **Incremental snapshots** - Only store changes between versions
- **Automated cleanup** - Remove old AMIs based on policy
- **Regional optimization** - Store AMIs in cost-effective regions
- **Compression** - Optimize AMI size for storage costs

### 2. Cost Sharing Model
```yaml
# Cost sharing configuration
cost_model:
  official_amis:
    - funded_by: aws_research_wizard_project
    - free_to_users: true

  community_amis:
    - creator_pays: snapshot_storage
    - users_pay: ec2_instance_costs

  enterprise_amis:
    - subscription_model: true
    - premium_support: true
```

## Implementation Roadmap

### Phase 1: Core AMI Creation (Q1 2025)
- [ ] Add AMI creation to deployment pipeline
- [ ] Implement AMI cleanup and optimization
- [ ] Basic AMI metadata and tagging
- [ ] CLI commands for AMI management

### Phase 2: AMI Sharing (Q2 2025)
- [ ] Account-to-account AMI sharing
- [ ] Multi-region AMI distribution
- [ ] AMI registry infrastructure
- [ ] Community AMI submission workflow

### Phase 3: AMI Marketplace (Q3 2025)
- [ ] Public AMI gallery
- [ ] Search and discovery features
- [ ] Rating and review system
- [ ] Usage analytics and metrics

### Phase 4: Advanced Features (Q4 2025)
- [ ] Automated AMI updates
- [ ] Security scanning integration
- [ ] Performance benchmarking
- [ ] Cost optimization tools

## Configuration Examples

### 1. Domain Pack AMI Configuration
```yaml
# configs/ami/genomics.yaml
ami_config:
  domain: genomics
  base_ami: ami-0c02fb55956c7d316  # Amazon Linux 2
  create_ami: true

  optimization:
    cleanup_build_artifacts: true
    remove_package_caches: true
    compress_logs: true

  sharing:
    auto_share_regions: ["us-west-2", "us-east-1"]
    community_registry: true
    public_sharing: false

  metadata:
    tags: ["genomics", "bioinformatics", "variant-calling"]
    description: "Complete genomics environment with 54 tools"
    supported_instances: ["r5.large", "r5.xlarge", "r5.2xlarge"]
```

### 2. User AMI Creation Workflow
```bash
# 1. Deploy environment normally
aws-research-wizard deploy --domain genomics --instance r5.large

# 2. Customize environment (install additional tools, data, etc.)
ssh -i ~/.ssh/id_rsa ec2-user@instance-ip
# ... customization work ...

# 3. Create AMI from customized environment
aws-research-wizard ami create --stack genomics-env --name "my-genomics-env" --tags "covid-19,custom"

# 4. Share with collaborators
aws-research-wizard ami share --ami-id ami-12345678 --account-id 123456789012

# 5. Submit to community registry
aws-research-wizard ami submit --ami-id ami-12345678 --registry community
```

## Success Metrics

### 1. Performance Metrics
- **Deployment time reduction**: 20-30 minutes → 2-3 minutes
- **Success rate improvement**: 85% → 99%
- **Cost reduction**: 40% less compute time for setup

### 2. Adoption Metrics
- **AMI creation rate**: Track user-created AMIs
- **Community sharing**: Number of shared AMIs
- **Deployment preference**: AMI vs. package installation ratio

### 3. Quality Metrics
- **AMI ratings**: Community feedback scores
- **Performance benchmarks**: Standardized testing results
- **Security compliance**: Clean security scans

This AMI strategy transforms AWS Research Wizard from a deployment tool into a comprehensive research environment ecosystem, enabling rapid collaboration and reproducible research at scale.
