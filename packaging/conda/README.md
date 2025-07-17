# Conda Package for AWS Research Wizard

This directory contains the conda package configuration for distributing AWS Research Wizard through conda-forge.

## Why Conda?

Conda is the preferred package manager for the scientific computing community, making it essential for AWS Research Wizard's target audience:

- **Scientific Computing Focus**: Most researchers already use conda for managing Python, R, and other scientific tools
- **Environment Management**: Conda's environment isolation is perfect for research workflows
- **Cross-Platform**: Works consistently across Linux, macOS, and Windows
- **Dependency Resolution**: Handles complex scientific package dependencies automatically
- **Integration**: Seamless integration with Jupyter, Spack, and other research tools

## Target Users

- Computational biologists using conda for bioinformatics workflows
- Climate scientists with existing conda environments
- Machine learning researchers using conda for PyTorch/TensorFlow
- Data scientists with conda-based data analysis stacks
- HPC users who prefer conda for package management

## Installation (Planned)

```bash
# Install from conda-forge (coming in v0.3.0)
conda install -c conda-forge aws-research-wizard

# Or create dedicated environment
conda create -n research-env aws-research-wizard
conda activate research-env
aws-research-wizard config tui
```

## Integration Benefits

- **Spack Integration**: AWS Research Wizard's Spack-powered environments complement conda
- **Jupyter Support**: Easy integration with Jupyter notebooks for research workflows
- **Scientific Libraries**: Works alongside numpy, scipy, pandas, and other scientific packages
- **Reproducible Environments**: Conda environment.yml files for reproducible research setups

## Development Status

- **Current**: Template ready for conda-forge submission
- **Q2 2025**: Conda package available on conda-forge
- **Target**: Make AWS research infrastructure as easy as `conda install`

## Submission Process

1. Create conda-forge recipe based on meta.yaml template
2. Submit pull request to conda-forge/staged-recipes
3. Set up automated builds and updates
4. Integrate with CI/CD pipeline for automatic version bumps
