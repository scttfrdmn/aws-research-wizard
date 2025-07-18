---
layout: default
title: AWS Research Wizard
description: Easily run research workloads in the cloud - 27 research domains with 20-minute tutorials
---

<div class="hero">
  <h1>Easily Run Research Workloads in the Cloud</h1>
  <p class="hero-subtitle">Easy-to-use tools that set up research computing for you. No cloud expertise needed. Clear costs. Works for 27 different science fields.</p>

  <div class="hero-buttons">
    <a href="getting-started/" class="btn btn-primary">🚀 Get Started</a>
    <a href="https://github.com/scttfrdmn/aws-research-wizard/releases" class="btn btn-secondary">📥 Download</a>
  </div>
</div>

## What is AWS Research Wizard?

AWS Research Wizard helps you easily run research workloads in the cloud. You get:

- **🎯 Ready-made research setups**: Pre-built environments for genomics, climate modeling, AI/ML, and 24 other fields
- **⚡ One-command setup**: Launch complex research computing with a single command
- **💰 Smart cost control**: Automatically picks the cheapest options and scales up/down as needed
- **🔬 Built for researchers**: Designed specifically for scientific computing and big data research

## Quick Start

Get up and running in 5 minutes:

```bash
# Download and install
wget https://github.com/scttfrdmn/aws-research-wizard/releases/latest/download/aws-research-wizard-linux-amd64.tar.gz
tar -xzf aws-research-wizard-linux-amd64.tar.gz
sudo mv aws-research-wizard /usr/local/bin/

# Set up your AWS account
aws configure

# See what research areas are available
aws-research-wizard config list

# Launch a genomics research environment
aws-research-wizard deploy --domain genomics --size standard
```

## Research Areas

### 🧬 Life Sciences
- **[Genomics & DNA Analysis](domain-guides/genomics/getting-started.md)**: Process DNA sequences with GATK, BWA, STAR tools
- **[Structural Biology](domain-guides/structural_biology/getting-started.md)**: Study protein structures and molecular movements
- **[Neuroscience](domain-guides/neuroscience/getting-started.md)**: Brain imaging analysis and connectivity studies
- **[Drug Discovery](domain-guides/drug_discovery/getting-started.md)**: Find new medicines using computer modeling

### 🌍 Physical Sciences
- **[Climate Modeling](domain-guides/climate_modeling/getting-started.md)**: Weather prediction with WRF and CESM models
- **[Materials Science](domain-guides/materials_science/getting-started.md)**: Design new materials with quantum calculations
- **[Chemistry & Materials](domain-guides/chemistry_materials/getting-started.md)**: Molecular dynamics and quantum chemistry
- **[Astronomy](domain-guides/astronomy_astrophysics/getting-started.md)**: Process telescope data and cosmic simulations

### ⚙️ Engineering
- **[Machine Learning & AI](domain-guides/machine_learning/getting-started.md)**: Train neural networks with PyTorch and TensorFlow
- **[Cybersecurity Research](domain-guides/cybersecurity_research/getting-started.md)**: Security analysis and threat detection
- **[Benchmarking & Performance](domain-guides/benchmarking_performance/getting-started.md)**: System performance analysis and optimization

### 🤖 Computer Science
- **[Geospatial Research](domain-guides/geospatial_research/getting-started.md)**: GIS analysis and satellite data processing
- **[Quantum Computing](domain-guides/quantum_computing/getting-started.md)**: Quantum algorithm development and simulation
- **[Mathematical Modeling](domain-guides/mathematical_modeling/getting-started.md)**: Numerical analysis and optimization

**[See All 27 Research Areas →](domain-guides/index.md)**

## What You Get

### 📦 Ready-to-Use Research Environments
Each research area includes optimized software with all the tools you need. No more spending weeks setting up your computing environment.

### 🚀 Smart Computer Selection
The system analyzes your research needs and picks the most cost-effective cloud computers automatically.

### 💰 Smart Cost Control
- Save 70-90% with spot pricing recommendations
- Right-size computers based on your actual usage
- Automatic scaling and shutdown when not in use
- Real-time cost monitoring and budget alerts

### 🔬 Research Workflows
Built-in workflows for common research tasks:
- Genomics DNA analysis pipelines
- Climate model simulation workflows
- Deep learning training pipelines
- High-throughput screening workflows

## Performance & Scale

- **27 Research Areas** covering major scientific disciplines
- **100% AWS Integration** - works with all supported cloud services
- **8.0 GB/s Transfer Speeds** for handling large datasets
- **50+ Petabytes of Research Data** - access to public research datasets

## Getting Help

- 📖 **[Step-by-Step Tutorials](domain-guides/index.md)**: 20-minute guides for each research area
- 🚀 **[Quick Start Guide](getting-started/)**: Get running in 5 minutes
- 💬 **[Community Forum](https://github.com/scttfrdmn/aws-research-wizard/discussions)**: Ask questions and get help
- 🐛 **[Report Problems](https://github.com/scttfrdmn/aws-research-wizard/issues)**: Bug reports and feature requests

## Contributing & Extending

**🚀 AWS Research Wizard is an open research platform!** We welcome suggestions and contributions:

### Suggest New Features
- **[Request domain packs](https://github.com/scttfrdmn/aws-research-wizard/issues/new?template=domain_request.md)**: Missing your research field? Let us know!
- **[Suggest software](https://github.com/scttfrdmn/aws-research-wizard/issues/new?template=software_request.md)**: Need specific tools added to a domain?
- **[Share workflows](https://github.com/scttfrdmn/aws-research-wizard/discussions)**: Discuss research workflows and best practices

### Contribute Code
- **[Development Guide](contributing/development/)**: Set up your development environment
- **[Creating Research Areas](contributing/domain-packs/)**: Add support for new research domains
- **[Testing](contributing/testing/)**: Help improve reliability and coverage

**Your suggestions drive our development roadmap** - this platform grows based on what researchers need!

## License

AWS Research Wizard is released under the [MIT License](https://github.com/scttfrdmn/aws-research-wizard/blob/main/LICENSE).
