# AWS Research Wizard - Domain Packs Repository

This repository contains pre-configured domain packs for research computing workloads on AWS.

## Structure

```
domain-packs/
├── domains/
│   ├── life-sciences/          # Biology, genomics, neuroscience
│   ├── physical-sciences/      # Physics, chemistry, astronomy
│   ├── engineering/           # CFD, mechanical, electrical
│   ├── computer-science/      # AI/ML, HPC, data science
│   └── social-sciences/       # Economics, digital humanities
├── shared/                    # Common configurations
├── tools/                     # Domain pack management tools
└── schemas/                   # Validation schemas
```

## Domain Packs

Each domain pack contains:
- **spack.yaml**: Spack environment specification
- **aws-config.yaml**: AWS infrastructure configuration
- **workflows/**: Sample research workflows
- **examples/**: Realistic usage examples
- **docs/**: Domain-specific documentation

## Quick Start

```bash
# Install a domain pack
aws-research-wizard install genomics

# List available packs
aws-research-wizard list

# Deploy research environment
aws-research-wizard deploy --domain genomics --size standard
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on adding new domain packs.