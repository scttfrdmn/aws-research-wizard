# AWS Research Wizard

<div class="hero">
  <h1>Run Research Workloads on AWS</h1>
  <p class="hero-subtitle">Execute complex research workloads effortlessly with pre-configured domain packs that handle all the AWS complexity for you.</p>

  <div class="hero-buttons">
    <a href="getting-started/" class="btn btn-primary">🚀 Get Started</a>
    <a href="https://github.com/aws-research-wizard/aws-research-wizard/releases" class="btn btn-secondary">📥 Download</a>
  </div>
</div>

## What is AWS Research Wizard?

AWS Research Wizard is a comprehensive platform that simplifies research computing on AWS by providing:

- **🎯 Pre-configured Domain Packs**: Ready-to-use research environments for genomics, climate modeling, AI/ML, and more
- **⚡ One-Command Deployment**: Deploy complex research infrastructure with a single command
- **💰 Cost Optimization**: Intelligent instance selection and automatic scaling to minimize costs
- **🔬 Research-Focused**: Optimized for scientific computing workflows and data-intensive research

## Quick Start

Get up and running in 5 minutes:

```bash
# Download and install
wget https://github.com/aws-research-wizard/releases/latest/aws-research-wizard-linux-amd64.tar.gz
tar -xzf aws-research-wizard-linux-amd64.tar.gz
sudo mv aws-research-wizard /usr/local/bin/

# Configure AWS credentials
aws configure

# Browse available domain packs
aws-research-wizard config list

# Deploy a research environment
aws-research-wizard deploy --domain genomics --size standard
```

## Research Domains

### 🧬 Life Sciences
- **[Genomics & Bioinformatics](domain-packs/life-sciences/genomics/)**: GATK, BWA, STAR, RNA-seq analysis
- **[Structural Biology](domain-packs/life-sciences/structural-biology/)**: Molecular modeling, protein structure analysis
- **[Neuroscience](domain-packs/life-sciences/neuroscience/)**: Brain imaging, connectivity analysis
- **[Drug Discovery](domain-packs/life-sciences/drug-discovery/)**: Molecular dynamics, virtual screening

### 🌍 Physical Sciences
- **[Climate Modeling](domain-packs/physical-sciences/climate-modeling/)**: WRF, CESM, weather forecasting
- **[Materials Science](domain-packs/physical-sciences/materials-science/)**: DFT calculations, molecular dynamics
- **[Chemistry](domain-packs/physical-sciences/chemistry/)**: Quantum chemistry, reaction modeling
- **[Astronomy](domain-packs/physical-sciences/astronomy/)**: Survey data processing, cosmological simulations

### ⚙️ Engineering
- **[CFD Engineering](domain-packs/engineering/cfd-engineering/)**: Fluid dynamics, aerodynamics simulation
- **[Mechanical Engineering](domain-packs/engineering/mechanical-engineering/)**: FEA, structural analysis
- **[Aerospace Engineering](domain-packs/engineering/aerospace-engineering/)**: Flight simulation, propulsion modeling

### 🤖 Computer Science
- **[AI/ML Research](domain-packs/computer-science/ai-research/)**: PyTorch, TensorFlow, distributed training
- **[HPC Development](domain-packs/computer-science/hpc-development/)**: Parallel computing, performance optimization
- **[Data Science](domain-packs/computer-science/data-science/)**: Large-scale analytics, visualization
- **[Quantum Computing](domain-packs/computer-science/quantum-computing/)**: Quantum algorithms, simulation

## Key Features

### 📦 Pre-Configured Environments
Each domain pack includes optimized software stacks with research-specific tools, libraries, and configurations. No more spending weeks setting up your computing environment.

### 🚀 Intelligent Infrastructure
Smart instance selection based on your workload characteristics. The system analyzes your research requirements and recommends the most cost-effective AWS infrastructure.

### 💰 Cost Optimization
- Spot instance recommendations for 70-90% savings
- Right-sizing based on actual resource usage
- Automatic scaling and shutdown policies
- Cost monitoring and budget alerts

### 🔬 Research Workflows
Built-in workflows for common research tasks:
- Genomics variant calling pipelines
- Climate model simulation workflows
- Deep learning training pipelines
- High-throughput screening workflows

## Performance & Scale

- **18+ Research Domain Packs** covering major scientific disciplines
- **100% AWS Integration Success** rate across all supported services
- **8.0 GB/s Peak Transfer Speeds** for large dataset handling
- **50+ PB AWS Open Data Access** for public research datasets

## Getting Help

- 📖 **[Documentation](getting-started/)**: Comprehensive guides and tutorials
- 🚀 **[Quick Start Guide](getting-started/)**: Get running in 5 minutes
- 💬 **[GitHub Discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)**: Community support
- 🐛 **[Issue Tracker](https://github.com/aws-research-wizard/aws-research-wizard/issues)**: Bug reports and feature requests

## Contributing

AWS Research Wizard is open source and welcomes contributions:

- **[Development Guide](contributing/development/)**: Set up your development environment
- **[Creating Domain Packs](contributing/domain-packs/)**: Add support for new research domains
- **[Testing](contributing/testing/)**: Help improve reliability and coverage

## License

AWS Research Wizard is released under the [MIT License](https://github.com/aws-research-wizard/aws-research-wizard/blob/main/LICENSE).
