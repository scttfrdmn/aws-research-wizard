# Real AWS Deployment Testing Plan

## Executive Summary

**CRITICAL CORRECTION:** Previous "deployment testing" was CLI validation only, not actual AWS infrastructure deployment. This plan establishes methodical testing of all 27 research packs against real AWS infrastructure using live AWS credentials and resources.

**Objective:** Systematically deploy, validate, and fix all 27 research domain packs on real AWS infrastructure.

---

## 🎯 Testing Methodology

### **Phase 1: Infrastructure Validation (Priority 1)**
1. **AWS Profile Setup & Validation**
   - Configure AWS CLI with real 'aws' profile credentials
   - Validate IAM permissions for CloudFormation, EC2, S3, etc.
   - Test basic AWS connectivity and resource creation

2. **Baseline Infrastructure Testing**
   - Deploy minimal test stack to validate CloudFormation integration
   - Verify EC2 instance provisioning with different instance types
   - Test security group creation and VPC configuration
   - Validate monitoring and cost tracking setup

### **Phase 2: Domain Pack Deployment (Priority 1)**
1. **Pilot Domain Testing (3-5 domains)**
   - Select representative domains from different categories
   - Deploy end-to-end with full Spack package installation
   - Document errors, failures, and performance issues
   - Establish error tracking and resolution workflow

2. **Systematic Full Deployment (All 27 domains)**
   - Deploy each domain pack sequentially
   - Track deployment success/failure rates
   - Measure deployment times and resource usage
   - Document Spack package installation issues

### **Phase 3: Validation & Optimization (Priority 2)**
1. **Functional Testing**
   - Verify installed software packages work correctly
   - Test domain-specific workflows and demos
   - Validate data movement and AWS Open Data integration
   - Confirm cost calculations match actual AWS billing

2. **Performance & Scaling Testing**
   - Test MPI scaling for HPC domains
   - Validate GPU acceleration for ML/quantum domains
   - Test auto-scaling and cost optimization features
   - Benchmark performance against expectations

---

## 🔧 Version Override Capabilities Analysis

### **Current State: Fixed Versions**
Current domain configs specify exact versions:
```yaml
spack_packages:
  core_aligners:
  - bwa@0.7.17 %gcc@11.4.0 +pic
  - gatk@4.4.0.0
  - star@2.7.10b %gcc@11.4.0 +shared+zlib
```

### **User Version Override Options**

#### **Option 1: Command-Line Override (Recommended)**
```bash
# Override specific package versions during deployment
./aws-research-wizard deploy start --domain genomics \
  --override "bwa@0.8.0,gatk@4.5.0.0,star@2.8.0"

# Override compiler version globally
./aws-research-wizard deploy start --domain genomics \
  --compiler gcc@12.3.0
```

#### **Option 2: Custom Config Files (User-Friendly)**
```bash
# Copy base domain config for customization
./aws-research-wizard config export genomics > my-genomics.yaml

# Edit versions in my-genomics.yaml, then deploy
./aws-research-wizard deploy start --config my-genomics.yaml
```

#### **Option 3: Environment-Specific Configs**
```yaml
# genomics-dev.yaml (development versions)
spack_packages:
  core_aligners:
  - bwa@develop %gcc@11.4.0 +pic

# genomics-stable.yaml (production versions)
spack_packages:
  core_aligners:
  - bwa@0.7.17 %gcc@11.4.0 +pic
```

### **Implementation Priority**
1. **Option 2 (Config Export/Edit)** - Easiest to implement, most user-friendly
2. **Option 1 (CLI Override)** - More advanced, requires parsing infrastructure
3. **Option 3 (Environment Configs)** - Long-term solution for teams

---

## 🚀 Real AWS Testing Execution Plan

### **Pre-Testing Setup**

#### **1. AWS Environment Preparation**
```bash
# Verify AWS CLI configuration
aws sts get-caller-identity --profile aws

# Check required permissions
aws iam get-user --profile aws
aws iam list-attached-user-policies --user-name $(aws sts get-caller-identity --query 'Arn' --output text --profile aws | cut -d'/' -f2) --profile aws

# Test basic resource creation
aws ec2 describe-regions --profile aws
```

#### **2. Testing Infrastructure Setup**
```bash
# Set up dedicated testing environment
export AWS_PROFILE=aws
export AWS_REGION=us-east-1
export TESTING_PREFIX="research-wizard-test"

# Create testing S3 bucket for artifacts
aws s3 mb s3://${TESTING_PREFIX}-artifacts --profile aws
```

### **Domain Testing Sequence (Ordered by Complexity)**

#### **Phase 1: Simple Domains (Low Risk)**
1. **Digital Humanities** - Minimal compute, basic packages
2. **Mathematical Modeling** - Standard compute, well-established packages
3. **Economics & Finance** - Standard statistical packages

#### **Phase 2: Medium Complexity Domains**
4. **Genomics** - Established bioinformatics tools, good Spack support
5. **Climate Modeling** - HPC components, MPI testing
6. **Structural Biology** - Mixed CPU/GPU workloads

#### **Phase 3: High Complexity Domains**
7. **Machine Learning** - GPU requirements, large packages
8. **Quantum Computing** - Specialized packages, memory requirements
9. **Drug Discovery** - Complex dependencies, commercial software

#### **Phase 4: All Remaining Domains (Complete Coverage)**
10-27. All other domains systematically

### **Per-Domain Testing Protocol**

#### **1. Pre-Deployment Validation**
```bash
# Validate domain configuration
./aws-research-wizard config info [domain] --validate

# Check AWS resource requirements
./aws-research-wizard deploy start --domain [domain] --dry-run --instance [type]
```

#### **2. Actual Deployment**
```bash
# Deploy with monitoring
./aws-research-wizard deploy start \
  --domain [domain] \
  --instance [recommended-type] \
  --stack test-[domain]-$(date +%Y%m%d-%H%M) \
  --timeout 60m \
  --debug

# Monitor deployment progress
./aws-research-wizard deploy status --stack test-[domain]-*
```

#### **3. Post-Deployment Validation**
```bash
# SSH into instance and validate packages
ssh -i ~/.ssh/research-key.pem ec2-user@[instance-ip]

# Run domain-specific validation commands
spack find
spack load [key-packages]
[domain-specific-tests]

# Test demo workflows if available
./aws-research-wizard data demo --domain [domain]
```

#### **4. Cleanup & Documentation**
```bash
# Collect logs and artifacts
./aws-research-wizard deploy logs --stack test-[domain]-* > logs/[domain]-deployment.log

# Clean up resources
./aws-research-wizard deploy delete --stack test-[domain]-* --force

# Document results
echo "[domain]: SUCCESS/FAILURE - [notes]" >> deployment-results.txt
```

---

## 📊 Error Tracking & Resolution

### **Error Categories & Tracking**

#### **1. Infrastructure Errors**
- CloudFormation stack failures
- EC2 instance provisioning issues
- Security group/VPC configuration problems
- IAM permission issues

#### **2. Spack Package Errors**
- Package compilation failures
- Dependency resolution conflicts
- Version compatibility issues
- Architecture-specific build problems

#### **3. Domain-Specific Errors**
- Missing dependencies in domain configs
- Incorrect package variants or flags
- Resource requirement mismatches
- Demo workflow failures

### **Error Resolution Workflow**

#### **1. Error Detection**
```bash
# Automated error detection during deployment
./aws-research-wizard deploy start [options] 2>&1 | tee deployment.log

# Parse logs for common error patterns
grep -E "(ERROR|FAILED|Exception)" deployment.log
```

#### **2. Error Classification**
- **P0 (Blocker)**: Deployment completely fails
- **P1 (Critical)**: Deployment succeeds but key packages fail
- **P2 (Major)**: Minor packages fail, domain partially functional
- **P3 (Minor)**: Documentation/demo issues only

#### **3. Resolution Process**
1. **Immediate Fix**: Simple config corrections, retry deployment
2. **Investigation Required**: Complex dependency issues, version conflicts
3. **Upstream Issues**: Spack package problems, report to maintainers
4. **Design Changes**: Domain config restructuring needed

---

## 📈 Success Metrics & Validation

### **Deployment Success Criteria**

#### **Minimum Viable Deployment (MVD)**
- ✅ CloudFormation stack creates successfully
- ✅ EC2 instance provisions and boots
- ✅ Spack environment installs without critical errors
- ✅ 80%+ of specified packages install successfully
- ✅ Basic domain functionality verified

#### **Full Success Deployment (FSD)**
- ✅ All MVD criteria met
- ✅ 95%+ of specified packages install successfully
- ✅ All demo workflows execute successfully
- ✅ Performance meets documented expectations
- ✅ Cost calculations match actual AWS billing

### **Testing Milestones**

#### **Week 1: Infrastructure & Pilot Testing**
- [ ] AWS profile and permissions validated
- [ ] 3-5 pilot domains deployed successfully
- [ ] Error tracking system operational
- [ ] Initial issue resolution workflow established

#### **Week 2-3: Full Domain Coverage**
- [ ] All 27 domains attempted
- [ ] Success/failure rates documented
- [ ] Critical issues identified and prioritized
- [ ] Performance benchmarks collected

#### **Week 4: Optimization & Fixes**
- [ ] Critical deployment issues resolved
- [ ] Performance optimizations implemented
- [ ] Documentation updated with findings
- [ ] Production readiness assessment completed

---

## 🔄 Continuous Integration Integration

### **Automated Testing Pipeline**
```yaml
# .github/workflows/aws-deployment-test.yml
name: AWS Domain Pack Testing
on:
  schedule:
    - cron: '0 2 * * 0'  # Weekly Sunday 2 AM
  workflow_dispatch:

jobs:
  test-domains:
    strategy:
      matrix:
        domain: [genomics, climate_modeling, machine_learning, ...]
    steps:
      - name: Deploy Domain
        run: |
          ./aws-research-wizard deploy start \
            --domain ${{ matrix.domain }} \
            --stack ci-test-${{ matrix.domain }}-${{ github.run_id }}

      - name: Validate Deployment
        run: |
          ./test-scripts/validate-domain.sh ${{ matrix.domain }}

      - name: Cleanup
        if: always()
        run: |
          ./aws-research-wizard deploy delete \
            --stack ci-test-${{ matrix.domain }}-${{ github.run_id }} --force
```

---

## 💡 Implementation Recommendations

### **Phase 1 Immediate Actions**
1. **Set up AWS testing profile** with appropriate permissions
2. **Implement config export functionality** for version overrides
3. **Start with 3 pilot domains** (digital humanities, genomics, mathematical modeling)
4. **Establish error tracking spreadsheet/system**

### **Phase 2 Systematic Testing**
1. **Deploy all 27 domains** following the complexity-based sequence
2. **Document all failures and issues** with detailed error logs
3. **Fix critical deployment blockers** that affect multiple domains
4. **Optimize successful deployments** for cost and performance

### **Phase 3 Production Readiness**
1. **Implement automated testing pipeline** for regression detection
2. **Create user documentation** based on real deployment experience
3. **Establish monitoring and alerting** for production deployments
4. **Validate cost models** against actual AWS billing data

---

## 🎯 Success Definition

**Real AWS deployment testing is successful when:**

1. **✅ 90%+ domains deploy successfully** on first attempt
2. **✅ All critical packages install** for successful domains
3. **✅ Demo workflows execute** without errors
4. **✅ Cost calculations accurate** within 10% of actual billing
5. **✅ Performance benchmarks met** for HPC/GPU workloads
6. **✅ Error resolution process** documented and proven
7. **✅ Version override capability** implemented and tested

**This represents true production readiness for the AWS Research Wizard platform.**

---

*Created: January 3, 2025*
*Real AWS deployment testing plan for 27 research domain packs*
*Replacing previous CLI-only validation approach*
