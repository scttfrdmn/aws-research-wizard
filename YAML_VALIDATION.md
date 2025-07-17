# YAML Configuration and Validation

## Why YAML Instead of JSON?

### 🎯 **Research Domain Configuration Requirements**

The AWS Research Wizard uses YAML for domain pack configurations instead of JSON for several key reasons:

#### **1. Human Readability for Research Community**
```yaml
# YAML: Easy for researchers to read and edit
spack_packages:
  text_processing:
    - python@3.11.5 %gcc@11.4.0 +optimizations+shared+ssl
    - py-nltk@3.8.1 %gcc@11.4.0
    - py-spacy@3.6.1 %gcc@11.4.0
```

vs

```json
// JSON: Harder to read and edit, especially for complex Spack specs
{
  "spack_packages": {
    "text_processing": [
      "python@3.11.5 %gcc@11.4.0 +optimizations+shared+ssl",
      "py-nltk@3.8.1 %gcc@11.4.0",
      "py-spacy@3.6.1 %gcc@11.4.0"
    ]
  }
}
```

#### **2. Comments and Documentation**
YAML allows inline comments, crucial for documenting complex research configurations:

```yaml
spack_packages:
  # Core Python environment with optimizations for research workloads
  text_processing:
    - python@3.11.5 %gcc@11.4.0 +optimizations+shared+ssl  # SSL support for data downloads
    - py-nltk@3.8.1 %gcc@11.4.0  # Natural Language Toolkit for text analysis
```

JSON does not support comments, making it harder to document why specific package versions or configurations are chosen.

#### **3. Multi-line Strings for Descriptions**
Research domain descriptions often require detailed explanations:

```yaml
description: |
  Computational platform for text analysis, cultural analytics, digital
  archives, and humanities data science research. Supports large-scale
  corpus analysis, topic modeling, and digital scholarship workflows.
```

#### **4. Less Syntax Overhead**
YAML reduces visual clutter, making configurations more approachable for researchers:

```yaml
# YAML: Clean and minimal
aws_instance_recommendations:
  development:
    instance_type: t3.medium
    vcpus: 2
    memory_gb: 4
```

vs

```json
// JSON: More syntax overhead
{
  "aws_instance_recommendations": {
    "development": {
      "instance_type": "t3.medium",
      "vcpus": 2,
      "memory_gb": 4
    }
  }
}
```

#### **5. Research Community Standards**
Many research tools and workflows use YAML:
- **Spack**: Uses YAML for environment specifications
- **Conda**: environment.yml files
- **Docker Compose**: docker-compose.yml
- **Kubernetes**: All manifests are YAML
- **CI/CD**: GitHub Actions, GitLab CI use YAML

## 🔍 **YAML Validation Strategy**

### **Three-Layer Validation Approach**

#### **1. Syntax Validation**
```bash
# Basic YAML syntax checking
python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"
```

#### **2. Style and Convention Linting**
```bash
# Comprehensive style checking with yamllint
yamllint -c .yamllint.yaml configs/domains/*.yaml
```

**Configuration highlights:**
- **Line length**: 120 characters (longer than standard 80 for scientific descriptions)
- **Indentation**: Consistent 2-space indentation
- **Comments**: Required spacing for readability
- **Trailing spaces**: Warnings (not errors) for flexibility

#### **3. Schema Validation**
```bash
# Domain-specific schema validation
python3 scripts/validate-domain-configs.py
```

**Validates:**
- Required fields (name, description, spack_packages, etc.)
- Data types and value constraints
- AWS instance type validity
- Spack package specification format
- Cost estimation completeness

### **Pre-commit Integration**

All YAML files are automatically validated before commit:

```yaml
# .pre-commit-config.yaml
- repo: https://github.com/adrienverge/yamllint
  rev: v1.35.1
  hooks:
    - id: yamllint
      args: ['-c', '.yamllint.yaml']
      types: [yaml]

- repo: local
  hooks:
    - id: domain-pack-validation
      name: Domain Pack Schema Validation
      entry: python3 scripts/validate-domain-configs.py
      files: '^configs/domains/.*\.yaml$'
```

### **Current Validation Results**

**Total Domain Configurations**: 27
**Valid Configurations**: 21 ✅
**Configurations with Issues**: 6 ❌

#### **Common Issues Found:**
1. **Missing required fields**: `primary_domains` field missing in 4 configurations
2. **Invalid enum values**:
   - Network type `'100_Gbps'` not in allowed values `['sr-iov', 'standard']`
   - Placement strategy `'single'` not in allowed values `['cluster', 'partition', 'spread']`

#### **Issues by Configuration:**
- `economics_finance.yaml`: Invalid network type
- `visualization_studio.yaml`: Missing `primary_domains`
- `forestry_natural_resources.yaml`: Missing `primary_domains`
- `renewable_energy_systems.yaml`: Missing `primary_domains`
- `quantum_computing.yaml`: Invalid placement strategy
- `food_science_nutrition.yaml`: Missing `primary_domains`

## 🛠 **Developer Workflow**

### **Adding/Modifying Domain Configurations**

1. **Edit YAML file** in `configs/domains/`
2. **Validate locally**:
   ```bash
   yamllint configs/domains/your-domain.yaml
   python3 scripts/validate-domain-configs.py
   ```
3. **Commit changes** - pre-commit hooks automatically validate
4. **CI/CD pipeline** runs additional validation in testing environment

### **YAML Best Practices for Research Configurations**

#### **✅ Do:**
- Use descriptive comments for complex configurations
- Keep lines under 120 characters when possible
- Use consistent 2-space indentation
- Quote strings that might be interpreted as numbers/booleans
- Use multi-line strings for long descriptions

#### **❌ Don't:**
- Mix tabs and spaces
- Use overly long lines (>150 characters)
- Include sensitive information (AWS keys, passwords)
- Use inconsistent naming conventions

### **Schema Evolution**

The domain pack schema (`configs/schemas/domain_pack_schema.yaml`) is versioned and evolving:

- **v1.0**: Basic domain pack structure
- **v1.1**: Added AWS cost modeling
- **v1.2**: Enhanced Spack package specifications
- **v2.0**: Terraform integration support (planned)

## 🚀 **Integration with Infrastructure**

### **Terraform Variable Generation**
YAML configurations are automatically converted to Terraform variables:

```yaml
# domain.yaml
aws_instance_recommendations:
  standard_analysis:
    instance_type: c6i.2xlarge
    vcpus: 8
    memory_gb: 16
```

Becomes:

```hcl
# terraform variables
variable "instance_type" {
  default = "c6i.2xlarge"
}
```

### **Spack Environment Generation**
YAML package specifications become Spack environment files:

```yaml
# domain.yaml
spack_packages:
  text_processing:
    - python@3.11.5 %gcc@11.4.0 +optimizations
```

Becomes:

```yaml
# spack.yaml
spack:
  specs:
    - python@3.11.5 %gcc@11.4.0 +optimizations
```

This YAML-first approach ensures consistency across the entire research infrastructure deployment pipeline while maintaining readability and maintainability for the research community.
