# Domain Pack Testing Guide - Lessons Learned

## Overview

This document captures critical learnings from systematic domain pack testing to ensure consistent, successful deployments across all 27 research domains. Apply these guidelines before testing each new domain.

**Last Updated:** 2025-07-18
**Based on:** Genomics tutorial testing results

---

## Infrastructure Requirements by Domain Type

### Bioinformatics Domains (genomics, marine_biology_oceanography, structural_biology)
- **Minimum Instance:** r5.xlarge (32GB RAM, 4 vCPUs)
- **Disk Space:** 50GB minimum (genomics generated 1.2GB+ data)
- **Package Manager:** Miniforge/Mamba (20x faster than Spack)
- **Key Channels:** bioconda, conda-forge

### Machine Learning Domains (machine_learning, neuroscience, computer_vision)
- **Minimum Instance:** c5.2xlarge (CPU) or g4dn.xlarge (GPU-enabled)
- **Disk Space:** 20GB minimum (much more efficient than expected)
- **Package Manager:** **pip3 PREFERRED** - dramatically faster than Spack/conda
- **Key Packages:** numpy, pandas, scikit-learn, matplotlib (pip ecosystem)

### Climate/Earth Science Domains (climate_modeling, geoscience, atmospheric_chemistry)
- **Minimum Instance:** c5.4xlarge (16 vCPUs for parallel processing)
- **Disk Space:** 40GB minimum (large datasets)
- **Package Manager:** Mix of conda + specialized HPC tools
- **Key Channels:** conda-forge, pyviz

### High-Performance Computing (physics_computational, chemistry_computational)
- **Minimum Instance:** c5n.2xlarge (enhanced networking)
- **Disk Space:** 30GB minimum
- **Package Manager:** Spack (for specialized HPC libraries)
- **Key Focus:** MPI, OpenMP, specialized solvers

---

## Critical Configuration Updates Needed

### 1. Terraform Infrastructure (`terraform/environments/aws/main.tf`)

```hcl
# BEFORE (line ~75)
root_block_device {
  volume_size = 8  # TOO SMALL
}

# AFTER - Domain-specific sizing
root_block_device {
  volume_size = var.domain_disk_size
}
```

**Required Variable Addition:**
```hcl
variable "domain_disk_size" {
  description = "Disk size based on domain requirements"
  type        = number
  default     = 30
}
```

### 2. User Data Script (`terraform/modules/research-environment/user_data.sh`)

**Lines 57-107 Need Domain-Specific Package Management:**

```bash
# Add domain-specific package manager selection
case "${domain_name}" in
  "genomics"|"marine_biology_oceanography"|"structural_biology")
    echo "Using Miniforge for bioinformatics domain"
    PACKAGE_MANAGER="miniforge"
    ;;
  "machine_learning"|"neuroscience")
    echo "Using Miniforge + pip for ML domain"
    PACKAGE_MANAGER="miniforge_pip"
    ;;
  "climate_modeling"|"geoscience")
    echo "Using conda + HPC tools for earth science domain"
    PACKAGE_MANAGER="conda_hpc"
    ;;
  *)
    echo "Using Spack for specialized HPC domain"
    PACKAGE_MANAGER="spack"
    ;;
esac
```

### 3. Domain Configuration Updates

**All 27 `configs/domains/*.yaml` files need:**

```yaml
# Add deployment strategy section
deployment_strategy:
  package_manager: "miniforge"  # or spack, conda_hpc
  disk_size_gb: 50
  min_memory_gb: 32
  preferred_instance_family: "r5"  # or c5, g4dn

# Update realistic cost estimates
estimated_cost:
  setup_time_minutes: 90        # Based on real testing
  workload_time_minutes: 32     # Domain-specific
  compute_cost_per_hour: 0.126  # r5.large actual
  total_cost_estimate: 0.26     # Realistic total
```

---

## Package Manager Strategy Matrix

| Domain Category | Primary | Secondary | Use Case |
|-----------------|---------|-----------|----------|
| **Bioinformatics** | Miniforge/Mamba | pip | Fast bio tool installation |
| **Machine Learning** | **pip3** | conda | Python ecosystem (fastest option) |
| **Earth Sciences** | conda | Spack | Mix of Python + HPC |
| **Pure HPC** | Spack | conda | Specialized libraries |
| **Visualization** | Miniforge | npm | Python + web tools |
| **Social Sciences** | Miniforge | R packages | Python + R ecosystem |

---

## Pre-Testing Validation Checklist

### Before Deploying Any Domain:

- [ ] **Check disk space requirements** - Review expected data sizes
- [ ] **Validate instance type** - Ensure adequate RAM/CPU for workload
- [ ] **Verify package specifications** - Test package manager compatibility
- [ ] **Review cost estimates** - Update based on realistic performance data
- [ ] **Check network requirements** - Large dataset downloads need bandwidth
- [ ] **Validate tutorial steps** - Ensure no step numbering conflicts

### During Testing:

- [ ] **Monitor disk usage** - `df -h` throughout deployment
- [ ] **Track setup time** - Measure actual infrastructure + package installation
- [ ] **Measure workload time** - Separate setup from actual computation
- [ ] **Document failures** - Capture specific error messages and solutions
- [ ] **Record performance** - CPU usage, memory usage, network I/O

### After Testing:

- [ ] **Update cost estimates** - Replace theoretical with actual costs
- [ ] **Document timing** - Real setup and workload execution times
- [ ] **Note optimizations** - What worked well, what needs improvement
- [ ] **Update recommendations** - Instance types, disk sizes, packages

---

## Known Issues and Solutions

### Issue: Spack Binary Cache Ineffective
**Problem:** Complex package variants don't match binary cache, forcing source compilation
**Solution:** Use Miniforge/Mamba for domains with standard scientific packages
**Affects:** All bioinformatics domains

### Issue: Disk Space Exhaustion
**Problem:** 8GB insufficient for genomics analysis (926MB alignment + 200MB reads)
**Solution:** Increase to 30-50GB based on domain
**Affects:** All data-intensive domains

### Issue: Package Installation Failures
**Problem:** Invalid Spack variant specifications in YAML configs
**Solution:** Validate all package specs before deployment
**Affects:** 11+ packages across multiple domains

### Issue: Instance Type Mismatches
**Problem:** CPU-optimized instances for memory-intensive workloads
**Solution:** Use memory-optimized (r5/r6i) for data analysis domains
**Affects:** Genomics, climate modeling, machine learning

---

## Domain Testing Priority Matrix

### High Priority (Test Next)
1. **machine_learning** - Large user base, GPU requirements
2. **climate_modeling** - Representative of HPC + data domains
3. **neuroscience** - Mixed Python + specialized tools

### Medium Priority
4. **chemistry_computational** - Pure HPC, Spack validation
5. **geoscience** - Earth science representative
6. **materials_science** - Engineering domain testing

### Lower Priority
7-27. **Remaining domains** - Apply learned patterns

---

## Cost Analysis Template

Track these metrics for each domain:

```yaml
domain_metrics:
  test_date: "2025-07-18"
  instance_type: "r5.large"

  setup_phase:
    time_minutes: 90
    cost_dollars: 0.189

  workload_phase:
    time_minutes: 32
    cost_dollars: 0.069

  total:
    time_minutes: 122
    cost_dollars: 0.258

  issues_found: 11
  optimization_notes: "Miniforge 20x faster than Spack"
```

---

## Implementation Action Items

### Immediate (Before Next Domain Test)
1. **Update Terraform disk sizes** - Add domain_disk_size variable
2. **Implement package manager selection** - Domain-specific logic in user_data.sh
3. **Fix instance type recommendations** - All 27 domain configs
4. **Create deployment validation script** - Pre-flight checks

### Short Term (Next 1-2 domains)
5. **Document ML domain testing** - GPU requirements, framework installations
6. **Test climate modeling** - HPC + data workflow validation
7. **Refine cost estimates** - Update all configs with real data

### Long Term (Full 27 domain validation)
8. **Automated testing pipeline** - Systematic validation across all domains
9. **AMI creation strategy** - Fast deployment for validated environments
10. **Performance optimization** - Domain-specific tuning

---

## Usage Notes

- **Reference this document** before testing each new domain pack
- **Update findings** after each domain test completion
- **Apply lessons learned** to avoid repeating known issues
- **Validate assumptions** - Don't assume patterns hold across all domains

This guide evolves with each domain tested, capturing institutional knowledge for reliable, efficient research environment deployment.
