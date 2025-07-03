# Infrastructure as Code Evaluation: Terraform vs Pulumi vs CloudFormation

## Executive Summary

**Recommendation: Terraform** for AWS Research Wizard infrastructure deployment, replacing CloudFormation for multi-cloud support and better developer experience.

**LocalStack Pro Integration:** Enables safe, cost-free development testing with full AWS service compatibility.

---

## 🔍 Current State Analysis

### **CloudFormation Issues (Why We're Replacing)**
- ❌ **AWS-only lock-in** - No multi-cloud support
- ❌ **Complex syntax** - JSON/YAML verbosity and limitations
- ❌ **Limited programming constructs** - No loops, functions, or complex logic
- ❌ **Poor debugging** - Unclear error messages and rollback issues
- ❌ **Deployment hanging** - Current blocking issue in Research Wizard
- ❌ **Limited modularity** - Difficult to create reusable components

### **Multi-Cloud Requirements**
- 🎯 **AWS** (primary) - Current focus with 27 research domains
- 🎯 **Azure** (secondary) - For institutional partnerships
- 🎯 **GCP** (tertiary) - For specific research requirements
- 🎯 **Kubernetes** - For container-based deployments

---

## ⚖️ Technology Comparison

### **Terraform**

#### **✅ Advantages**
- **Multi-cloud native** - 3000+ providers (AWS, Azure, GCP, K8s)
- **Mature ecosystem** - 8+ years, huge community, extensive modules
- **Declarative HCL** - Clean, readable infrastructure code
- **State management** - Robust state tracking and locking
- **Plan/Apply workflow** - Safe preview before changes
- **LocalStack integration** - Excellent development support
- **Go ecosystem fit** - Can be embedded or called from Go applications

#### **❌ Disadvantages**
- **Learning curve** - HCL syntax and Terraform concepts
- **State file complexity** - Requires careful state management
- **Large binary** - ~200MB download

#### **Research Wizard Fit: 9/10**

### **Pulumi**

#### **✅ Advantages**
- **Real programming languages** - TypeScript, Python, Go, .NET
- **Multi-cloud native** - Same providers as Terraform
- **Better abstraction** - Component-based architecture
- **Superior testing** - Unit tests with real programming languages
- **Great Go integration** - Native Go SDK

#### **❌ Disadvantages**
- **Newer ecosystem** - Smaller community, fewer examples
- **SaaS dependency** - Pulumi Cloud for state (can be self-hosted)
- **Complexity** - More moving parts than Terraform
- **LocalStack support** - Less mature than Terraform

#### **Research Wizard Fit: 7/10**

### **CloudFormation (Current)**

#### **✅ Advantages**
- **Native AWS integration** - Perfect AWS service coverage
- **No additional tools** - Built into AWS

#### **❌ Disadvantages**
- **AWS lock-in** - Zero multi-cloud support
- **Deployment issues** - Current hanging problem
- **Poor syntax** - Verbose and limited
- **No LocalStack Pro** - Can't use for development

#### **Research Wizard Fit: 3/10**

---

## 🏗️ Terraform Implementation Plan

### **Architecture Design**

#### **1. Module Structure**
```
terraform/
├── modules/
│   ├── research-environment/     # Core research environment
│   ├── domain-specific/          # Domain-specific resources
│   ├── networking/               # VPC, security groups
│   ├── compute/                  # EC2, auto-scaling
│   ├── storage/                  # S3, EFS, EBS
│   └── monitoring/               # CloudWatch, alerts
├── environments/
│   ├── localstack/               # LocalStack development
│   ├── dev/                      # Development AWS
│   ├── staging/                  # Staging AWS
│   └── prod/                     # Production AWS
└── domains/
    ├── genomics.tf               # Domain-specific configs
    ├── climate_modeling.tf
    └── ...                       # All 27 domains
```

#### **2. Core Research Environment Module**
```hcl
# modules/research-environment/main.tf
resource "aws_instance" "research_node" {
  ami           = var.ami_id
  instance_type = var.instance_type

  vpc_security_group_ids = [aws_security_group.research.id]
  subnet_id              = var.subnet_id

  user_data = templatefile("${path.module}/user_data.sh", {
    domain_name    = var.domain_name
    spack_packages = var.spack_packages
    s3_bucket      = var.s3_bucket
  })

  tags = {
    Name   = "research-${var.domain_name}-${var.instance_id}"
    Domain = var.domain_name
    Owner  = var.owner
    Budget = var.monthly_budget
  }
}
```

#### **3. Domain-Specific Configuration**
```hcl
# domains/genomics.tf
module "genomics_environment" {
  source = "../modules/research-environment"

  domain_name    = "genomics"
  instance_type  = "r6i.4xlarge"
  ami_id         = data.aws_ami.spack_optimized.id

  spack_packages = [
    "bwa@0.7.17 %gcc@11.4.0 +pic",
    "gatk@4.4.0.0",
    "star@2.7.10b %gcc@11.4.0 +shared+zlib"
  ]

  storage_config = {
    ebs_size = 500
    efs_enabled = true
  }

  monitoring = {
    cost_alerts = true
    max_monthly_cost = 900
  }
}
```

### **LocalStack Pro Integration**

#### **1. Development Environment**
```hcl
# environments/localstack/main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  region                      = "us-east-1"
  skip_credentials_validation = true
  skip_metadata_api_check    = true
  skip_requesting_account_id  = true

  endpoints {
    ec2            = "http://localhost:4566"
    s3             = "http://localhost:4566"
    cloudformation = "http://localhost:4566"
    iam            = "http://localhost:4566"
  }
}
```

#### **2. LocalStack Startup Script**
```bash
#!/bin/bash
# scripts/start-localstack.sh
export LOCALSTACK_AUTH_TOKEN=ls-dIgo5507-demE-zaWa-8414-VEnUSEre3d8d
export LOCALSTACK_VOLUME_DIR=/tmp/localstack
export SERVICES=ec2,s3,iam,cloudformation,cloudwatch

# Start LocalStack Pro with research-specific configuration
docker run --rm -it \
  -p 4566:4566 \
  -p 4510-4559:4510-4559 \
  -e LOCALSTACK_AUTH_TOKEN=$LOCALSTACK_AUTH_TOKEN \
  -e SERVICES=$SERVICES \
  -e DEBUG=1 \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $LOCALSTACK_VOLUME_DIR:/var/lib/localstack \
  localstack/localstack-pro:latest
```

---

## 🚀 Implementation Roadmap

### **Phase 1: Terraform Foundation (Week 1)**
1. **Replace CloudFormation deployment** with basic Terraform
2. **Set up LocalStack Pro** development environment
3. **Create core research environment module**
4. **Test with 3 pilot domains** (digital_humanities, genomics, mathematical_modeling)

### **Phase 2: Domain Coverage (Week 2)**
1. **Migrate all 27 domains** to Terraform configurations
2. **Implement version override** functionality
3. **Add comprehensive monitoring** and cost controls
4. **Test systematic deployment** across all domains

### **Phase 3: Multi-Cloud & Production (Week 3)**
1. **Design multi-cloud abstraction** layer
2. **Add Azure provider** support for key domains
3. **Implement CI/CD pipeline** with automated testing
4. **Production deployment** validation

---

## 💻 Go Integration Strategy

### **Terraform Provider Approach**
```go
// internal/terraform/provider.go
package terraform

import (
    "context"
    "os/exec"
    "encoding/json"
)

type TerraformProvider struct {
    WorkingDir string
    Environment string // localstack, dev, staging, prod
}

func (tp *TerraformProvider) Deploy(ctx context.Context, domain string, config DeployConfig) error {
    // 1. Generate domain-specific terraform configuration
    configPath := tp.generateDomainConfig(domain, config)

    // 2. Run terraform plan
    planOutput, err := tp.runTerraform("plan", "-var-file", configPath)
    if err != nil {
        return fmt.Errorf("terraform plan failed: %w", err)
    }

    // 3. Show plan to user (if not auto-approve)
    if !config.AutoApprove {
        fmt.Println("Terraform Plan:")
        fmt.Println(planOutput)
        if !confirmDeployment() {
            return fmt.Errorf("deployment cancelled by user")
        }
    }

    // 4. Run terraform apply
    return tp.runTerraform("apply", "-auto-approve", "-var-file", configPath)
}
```

### **Configuration Generation**
```go
// internal/config/terraform.go
func (dv *DomainValidator) GenerateTerraformConfig(domain string, overrides map[string]string) (string, error) {
    domainConfig, err := dv.LoadDomainConfig(domain)
    if err != nil {
        return "", err
    }

    // Apply version overrides
    if overrides != nil {
        domainConfig = dv.applyVersionOverrides(domainConfig, overrides)
    }

    // Generate terraform variables
    tfVars := map[string]interface{}{
        "domain_name":    domainConfig.Name,
        "instance_type":  domainConfig.AWSInstanceRecommendations.StandardAnalysis.InstanceType,
        "spack_packages": domainConfig.SpackPackages,
        "monthly_budget": domainConfig.EstimatedCost.Total,
    }

    return dv.renderTerraformTemplate(tfVars)
}
```

---

## 🧪 Testing Strategy

### **1. LocalStack Pro Development**
```bash
# Start LocalStack Pro
./scripts/start-localstack.sh

# Deploy to LocalStack
export TF_VAR_environment=localstack
./aws-research-wizard deploy start --domain genomics --provider terraform

# Validate deployment
aws ec2 describe-instances --endpoint-url=http://localhost:4566
```

### **2. Real AWS Validation**
```bash
# Deploy to development AWS
export TF_VAR_environment=dev
export AWS_PROFILE=aws
./aws-research-wizard deploy start --domain genomics --provider terraform

# Cost monitoring
terraform output estimated_monthly_cost
aws ce get-cost-and-usage --time-period Start=2025-01-01,End=2025-01-31
```

### **3. Multi-Cloud Testing**
```bash
# Test Azure deployment (future)
export TF_VAR_environment=azure-dev
./aws-research-wizard deploy start --domain genomics --provider terraform --cloud azure
```

---

## 📊 Success Metrics

### **Phase 1 Success (LocalStack + Basic Terraform)**
- ✅ LocalStack Pro running with full AWS service mocking
- ✅ Terraform deploys 3 pilot domains successfully
- ✅ No deployment hanging issues
- ✅ Cost estimation accurate within LocalStack
- ✅ Version override functionality working

### **Phase 2 Success (Full Domain Coverage)**
- ✅ All 27 domains deploy via Terraform
- ✅ Real AWS deployment successful for 5+ domains
- ✅ Performance meets documented specifications
- ✅ Cost calculations match actual AWS billing

### **Phase 3 Success (Multi-Cloud Ready)**
- ✅ Infrastructure abstracted for multi-cloud
- ✅ Azure deployment working for key domains
- ✅ CI/CD pipeline operational
- ✅ Production deployment validated

---

## 💡 Implementation Priority

### **Immediate (This Session)**
1. **Set up LocalStack Pro** with auth token ✅
2. **Create basic Terraform module** for research environments
3. **Test simple deployment** (digital_humanities domain)
4. **Validate LocalStack integration** works properly

### **Next Session**
1. **Replace CloudFormation** in Go deployment code
2. **Implement version override** config export functionality
3. **Test with real AWS** using Terraform
4. **Document migration** from CloudFormation

### **This Week**
1. **Complete all 27 domains** Terraform migration
2. **Add comprehensive monitoring** and cost controls
3. **Production readiness** validation
4. **Multi-cloud foundation** architecture

---

## 🎯 Decision: Terraform + LocalStack Pro

**Primary Choice: Terraform**
- ✅ Best multi-cloud support
- ✅ Mature ecosystem and community
- ✅ Excellent LocalStack Pro integration
- ✅ Clean HCL syntax
- ✅ Robust state management
- ✅ Perfect fit for Research Wizard architecture

**Development Environment: LocalStack Pro**
- ✅ Full AWS service compatibility
- ✅ Cost-free development and testing
- ✅ Rapid iteration and debugging
- ✅ CI/CD pipeline integration
- ✅ Safe experimentation

**This approach solves the CloudFormation deployment hanging issue while setting up the foundation for multi-cloud support and significantly better developer experience.**

---

*Created: January 3, 2025*
*Infrastructure as Code evaluation for AWS Research Wizard*
*Recommendation: Terraform + LocalStack Pro for multi-cloud research computing*
