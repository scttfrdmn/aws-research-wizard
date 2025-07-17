# Configuration Versioning Strategy

## Overview

The AWS Research Wizard implements semantic versioning for all configuration files to ensure backward compatibility, track breaking changes, and enable safe configuration evolution across the research infrastructure.

## 📋 **Configuration File Categories**

### **1. Domain Pack Configurations**
- **Location**: `configs/domains/*.yaml`
- **Count**: 27 research domain configurations
- **Purpose**: Define research environment specifications, Spack packages, AWS resources
- **Versioning**: Individual per-domain + schema versioning

### **2. Schema Definitions**
- **Location**: `configs/schemas/domain_pack_schema.yaml`
- **Purpose**: Validation schema for domain pack configurations
- **Versioning**: Global schema versions with migration support

### **3. Infrastructure Templates**
- **Location**: `configs/templates/*.yaml`
- **Purpose**: Base AWS infrastructure configurations
- **Versioning**: Template versioning with compatibility matrix

### **4. Demo and Sample Data**
- **Location**: `configs/demo_data/*.yaml`, `configs/data_sources.yaml`
- **Purpose**: Example configurations and data source definitions
- **Versioning**: Synchronized with schema versions

### **5. Build and Tool Configurations**
- **Files**: `.pre-commit-config.yaml`, `.yamllint.yaml`, `mkdocs.yml`, `pyproject.toml`
- **Purpose**: Development workflow and tool configurations
- **Versioning**: Tool-specific versioning with compatibility tracking

## 🏗️ **Semantic Versioning Implementation**

### **Version Format: `MAJOR.MINOR.PATCH`**

#### **MAJOR Version Changes**
- Breaking schema changes requiring configuration migration
- Removal of required fields or complete restructuring
- Incompatible AWS service API changes
- Spack environment format changes

#### **MINOR Version Changes**
- New optional fields or configuration sections
- Additional AWS instance types or regions
- New Spack package categories
- Enhanced validation rules (non-breaking)

#### **PATCH Version Changes**
- Bug fixes in existing configurations
- Typo corrections and documentation updates
- Performance improvements without schema changes
- Security updates that don't change structure

## 📁 **Versioning Implementation Structure**

### **1. Schema Versioning**

```yaml
# configs/schemas/domain_pack_schema.yaml
schema_version: "2.1.0"
schema_metadata:
  version: "2.1.0"
  created: "2025-01-15"
  last_modified: "2025-07-03"
  compatibility:
    min_config_version: "2.0.0"
    max_config_version: "2.x.x"
  breaking_changes:
    - version: "2.0.0"
      description: "Added required aws_instance_recommendations field"
      migration: "Use migration tool: scripts/migrate-v1-to-v2.py"

# Schema definition continues...
type: object
required:
  - config_version  # NEW: Required version field
  - name
  - description
  # ... rest of schema
```

### **2. Domain Configuration Versioning**

```yaml
# configs/domains/digital_humanities.yaml
config_version: "2.1.0"
config_metadata:
  created: "2024-12-15"
  last_modified: "2025-07-03"
  schema_compatibility: "2.1.0"
  change_log:
    - version: "2.1.0"
      date: "2025-07-03"
      changes: ["Added quantum_computing packages", "Updated cost estimates"]
    - version: "2.0.0"
      date: "2025-01-15"
      changes: ["Added AWS instance recommendations", "Restructured spack_packages"]

name: Digital Humanities Research Laboratory
description: |
  Computational platform for text analysis, cultural analytics, digital
  archives, and humanities data science research
# ... rest of configuration
```

### **3. Version Compatibility Matrix**

```yaml
# configs/VERSION_COMPATIBILITY.yaml
compatibility_matrix:
  schema_versions:
    "2.1.0":
      release_date: "2025-07-03"
      domain_config_versions: ["2.1.0", "2.0.0"]
      terraform_versions: ["1.5.x", "1.6.x"]
      spack_versions: ["0.20.x", "0.21.x"]
      breaking_changes: []

    "2.0.0":
      release_date: "2025-01-15"
      domain_config_versions: ["2.0.0"]
      terraform_versions: ["1.4.x", "1.5.x"]
      spack_versions: ["0.19.x", "0.20.x"]
      breaking_changes:
        - "Added required aws_instance_recommendations field"
        - "Restructured spack_packages from flat list to categorized"

deprecated_versions:
  "1.x.x":
    deprecation_date: "2025-01-15"
    end_of_life: "2025-07-01"
    migration_path: "Use scripts/migrate-v1-to-v2.py"
```

## 🔧 **Versioning Tools and Scripts**

### **1. Version Management Script**

```python
# scripts/config-version-manager.py
#!/usr/bin/env python3
"""
Configuration Version Management Tool
Handles versioning, validation, and migration of configuration files
"""

import yaml
import semver
from pathlib import Path
from typing import Dict, List, Optional

class ConfigVersionManager:
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
        self.configs_dir = root_dir / "configs"

    def bump_version(self, config_type: str, bump_type: str):
        """Bump version for configuration type"""
        pass

    def validate_compatibility(self) -> List[str]:
        """Validate version compatibility across all configs"""
        pass

    def migrate_configs(self, from_version: str, to_version: str):
        """Migrate configurations between versions"""
        pass
```

### **2. Version Validation Integration**

```python
# scripts/validate-domain-configs.py (Enhanced)
def validate_domain_config(config_path, schema_path):
    """Enhanced validation with version checking"""
    try:
        # Load config and schema
        with open(config_path, 'r') as f:
            config = yaml.safe_load(f)
        with open(schema_path, 'r') as f:
            schema = yaml.safe_load(f)

        # Validate version compatibility
        config_version = config.get('config_version')
        schema_version = schema.get('schema_version')

        if not config_version:
            return False, "Missing config_version field"

        if not is_compatible(config_version, schema_version):
            return False, f"Config version {config_version} incompatible with schema {schema_version}"

        # Standard JSON schema validation
        jsonschema.validate(config, schema)
        return True, None

    except Exception as e:
        return False, str(e)
```

### **3. Migration Scripts**

```python
# scripts/migrate-v1-to-v2.py
#!/usr/bin/env python3
"""
Migration script from v1.x to v2.x configuration format
"""

def migrate_domain_config_v1_to_v2(config_path: Path):
    """Migrate v1 domain config to v2 format"""
    with open(config_path, 'r') as f:
        config = yaml.safe_load(f)

    # Add version metadata
    config['config_version'] = "2.0.0"
    config['config_metadata'] = {
        'created': config.get('created_date', '2024-01-01'),
        'last_modified': datetime.now().isoformat()[:10],
        'schema_compatibility': "2.0.0",
        'migrated_from': config.get('version', '1.0.0')
    }

    # Restructure spack_packages from flat list to categorized
    if 'packages' in config:
        config['spack_packages'] = {
            'default': config['packages']
        }
        del config['packages']

    # Add required aws_instance_recommendations
    if 'aws_instance_recommendations' not in config:
        config['aws_instance_recommendations'] = {
            'development': {
                'instance_type': 't3.medium',
                'vcpus': 2,
                'memory_gb': 4,
                'cost_per_hour': 0.0416
            }
        }

    # Write updated config
    with open(config_path, 'w') as f:
        yaml.safe_dump(config, f, default_flow_style=False, sort_keys=False)
```

## 🔄 **Version Lifecycle Management**

### **Release Process**

1. **Version Planning**
   - Identify breaking vs. non-breaking changes
   - Plan migration path for major versions
   - Update compatibility matrix

2. **Version Implementation**
   - Update schema version in `domain_pack_schema.yaml`
   - Bump configuration versions as needed
   - Create migration scripts for breaking changes

3. **Validation and Testing**
   - Run compatibility validation across all configs
   - Test migration scripts on sample configurations
   - Validate Terraform and Spack integration

4. **Documentation and Release**
   - Update VERSION_COMPATIBILITY.yaml
   - Document breaking changes and migration paths
   - Tag release with semantic version

### **Deprecation Policy**

- **Minor versions**: Supported for 1 year after release
- **Major versions**: 6-month deprecation period before EOL
- **Migration tools**: Provided for all breaking changes
- **Compatibility**: Maintain backward compatibility within major versions

## 📊 **Version Tracking Dashboard**

### **Current Version Status**

```yaml
# configs/VERSION_STATUS.yaml
current_versions:
  schema: "2.1.0"
  infrastructure_templates: "1.3.0"
  demo_data: "2.0.0"
  build_configs: "1.2.0"

domain_configurations:
  total_configs: 27
  version_distribution:
    "2.1.0": 22
    "2.0.0": 5
    "1.x.x": 0  # All migrated

compatibility_status:
  fully_compatible: 27
  needs_migration: 0
  deprecated: 0
```

### **Automated Version Checks**

```yaml
# .github/workflows/config-version-check.yml
name: Configuration Version Validation
on: [push, pull_request]

jobs:
  version-validation:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install dependencies
        run: pip install pyyaml jsonschema semver
      - name: Validate configuration versions
        run: python3 scripts/validate-config-versions.py
      - name: Check compatibility matrix
        run: python3 scripts/check-version-compatibility.py
```

## 🎯 **Benefits of Configuration Versioning**

### **1. Safe Evolution**
- Clear upgrade paths for configuration changes
- Backward compatibility guarantees within major versions
- Automated migration tools for breaking changes

### **2. Research Reproducibility**
- Pin specific configuration versions for research projects
- Track exactly which configuration generated research results
- Enable recreation of historical research environments

### **3. Infrastructure Reliability**
- Version compatibility validation prevents deployment failures
- Clear dependency relationships between configurations and tools
- Rollback capability for problematic configuration changes

### **4. Collaboration and Maintenance**
- Clear change history for all configuration modifications
- Version-aware validation and linting
- Automated compatibility checking in CI/CD pipelines

This comprehensive versioning strategy ensures that the AWS Research Wizard configuration ecosystem remains stable, maintainable, and evolution-friendly while supporting the complex needs of research infrastructure management.
