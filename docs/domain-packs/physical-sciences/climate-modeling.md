# Climate Modeling

Weather prediction and climate simulation tools

## Overview

**Domain Pack**: `climate-modeling`
**Version**: 1.0.0
**Categories**: physical-sciences
**Maintainers**: AWS Research Wizard Team

## Quick Start

```bash
# Deploy this domain pack
aws-research-wizard deploy --domain climate-modeling --size medium

# Get detailed information
aws-research-wizard config info climate-modeling

# List available workflows
aws-research-wizard workflow list --domain climate-modeling
```

## Software Stack

### Core Packages
- **wrf**: wrf@4.5.1 %gcc@11.4.0 +mpi+openmp
- **cesm**: cesm@2.1.3 %gcc@11.4.0 +mpi
- **netcdf-c**: netcdf-c@4.9.2 %gcc@11.4.0 +mpi+parallel-netcdf
- **netcdf-fortran**: netcdf-fortran@4.6.0 %gcc@11.4.0
- **hdf5**: hdf5@1.14.1 %gcc@11.4.0 +mpi+fortran
- **nco**: nco@5.1.5 %gcc@11.4.0 +netcdf
- **cdo**: cdo@2.2.0 %gcc@11.4.0 +netcdf+hdf5
- **grib-api**: grib-api@1.28.0 %gcc@11.4.0
- **eccodes**: eccodes@2.30.2 %gcc@11.4.0
- **openmpi**: openmpi@4.1.5 %gcc@11.4.0 +pmi+slurm

*And 1 more packages...*

### Optimization Settings
- **Compiler**: gcc@11.4.0
- **Target Architecture**: x86_64_v3
- **Optimization Flags**: -O3 -march=native -ffast-math

## AWS Infrastructure

### Instance Types
- **Small**: `c6i.4xlarge`
- **Medium**: `c6i.12xlarge`
- **Large**: `c6a.48xlarge`

### Storage Configuration
- **Type**: gp3
- **Size**: 2000 GB
- **IOPS**: 16000
- **Throughput**: 1000 MB/s

## Research Workflows

This domain pack includes 2 pre-configured research workflows:

### Weather Forecast

Regional weather forecasting with WRF

```bash
# Run this workflow
aws-research-wizard workflow run weather_forecast --domain climate-modeling
```

- **Input Data**: s3://aws-research-data/weather/gfs/
- **Expected Output**: Weather forecast grids in NetCDF format

### Climate Simulation

Long-term climate modeling with CESM

```bash
# Run this workflow
aws-research-wizard workflow run climate_simulation --domain climate-modeling
```

- **Input Data**: s3://aws-research-data/climate/forcing/
- **Expected Output**: Climate model outputs and analysis


## Cost Estimates

| Workload Size | Estimated Daily Cost |
|---------------|---------------------|
| Small | $15-45/day |
| Medium | $50-150/day |
| Large | $70-250/day |

!!! note "Cost Optimization"
    These estimates assume on-demand pricing. Significant savings are possible with:

    - **Spot Instances**: 70-90% savings for fault-tolerant workloads
    - **Reserved Instances**: 30-60% savings for predictable usage
    - **Savings Plans**: 20-72% savings with flexible commitment

## Example Configuration

```yaml
# climate-modeling-research-config.yaml
domain: climate-modeling
size: medium
aws:
  region: us-east-1
  availability_zone: us-east-1a
compute:
  instance_type: c6i.12xlarge
  instance_count: 1
storage:
  type: gp3
  size_gb: 2000
```

## Getting Help

- 📖 **Domain-Specific Documentation**: [tutorials/climate-modeling/](../../tutorials/climate-modeling/)
- 💬 **Community Support**: [GitHub Discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)
- 🐛 **Issues**: [GitHub Issues](https://github.com/aws-research-wizard/aws-research-wizard/issues)

## Related Domain Packs

Other domain packs in **Physical Sciences**:
