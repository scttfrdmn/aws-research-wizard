# Genomics

Complete genomics analysis with optimized bioinformatics tools

## Overview

**Domain Pack**: `genomics`
**Version**: 1.0.0
**Categories**: life-sciences
**Maintainers**: AWS Research Wizard Team

## Quick Start

```bash
# Deploy this domain pack
aws-research-wizard deploy --domain genomics --size medium

# Get detailed information
aws-research-wizard config info genomics

# List available workflows
aws-research-wizard workflow list --domain genomics
```

## Software Stack

### Core Packages
- **bwa**: bwa@0.7.17 %gcc@11.4.0 +pic
- **bwa-mem2**: bwa-mem2@2.2.1 %gcc@11.4.0 +sse4
- **bowtie2**: bowtie2@2.5.0 %gcc@11.4.0 +tbb
- **star**: star@2.7.10b %gcc@11.4.0 +shared+zlib
- **hisat2**: hisat2@2.2.1 %gcc@11.4.0 +sse4
- **minimap2**: minimap2@2.26 %gcc@11.4.0 +sse4
- **blast-plus**: blast-plus@2.14.0 %gcc@11.4.0 +pic
- **gatk**: gatk@4.4.0.0 %gcc@11.4.0
- **samtools**: samtools@1.17 %gcc@11.4.0 +curses
- **bcftools**: bcftools@1.17 %gcc@11.4.0 +curses

*And 4 more packages...*

### Optimization Settings
- **Compiler**: gcc@11.4.0
- **Target Architecture**: x86_64_v3
- **Optimization Flags**: -O3 -march=native

## AWS Infrastructure

### Instance Types
- **Small**: `c6i.2xlarge`
- **Medium**: `r6i.4xlarge`
- **Large**: `r6i.8xlarge`

### Storage Configuration
- **Type**: gp3
- **Size**: 500 GB
- **IOPS**: 3000
- **Throughput**: 125 MB/s

## Research Workflows

This domain pack includes 2 pre-configured research workflows:

### Variant Calling

GATK best practices variant calling pipeline

```bash
# Run this workflow
aws-research-wizard workflow run variant_calling --domain genomics
```

- **Input Data**: s3://aws-research-data/1000genomes/samples/
- **Expected Output**: VCF files with called variants

### Rna Seq Analysis

RNA-seq differential expression analysis

```bash
# Run this workflow
aws-research-wizard workflow run rna_seq_analysis --domain genomics
```

- **Input Data**: s3://aws-research-data/tcga/rna-seq/
- **Expected Output**: Gene expression matrices and DEG results


## Cost Estimates

| Workload Size | Estimated Daily Cost |
|---------------|---------------------|
| Small | $5-15/day |
| Medium | $25-75/day |
| Large | $50-200/day |

!!! note "Cost Optimization"
    These estimates assume on-demand pricing. Significant savings are possible with:

    - **Spot Instances**: 70-90% savings for fault-tolerant workloads
    - **Reserved Instances**: 30-60% savings for predictable usage
    - **Savings Plans**: 20-72% savings with flexible commitment

## Example Configuration

```yaml
# genomics-research-config.yaml
domain: genomics
size: medium
aws:
  region: us-east-1
  availability_zone: us-east-1a
compute:
  instance_type: r6i.4xlarge
  instance_count: 1
storage:
  type: gp3
  size_gb: 500
```

## Getting Help

- 📖 **Domain-Specific Documentation**: [tutorials/genomics/](../../tutorials/genomics/)
- 💬 **Community Support**: [GitHub Discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)
- 🐛 **Issues**: [GitHub Issues](https://github.com/aws-research-wizard/aws-research-wizard/issues)

## Related Domain Packs

Other domain packs in **Life Sciences**:
