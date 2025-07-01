# Domain Packs

AWS Research Wizard provides pre-configured domain packs for major research disciplines. Each domain pack includes optimized software stacks, AWS infrastructure configurations, and sample workflows.

## Categories

### [Computer Science](computer-science/)

- **[Ai-Research](computer-science/ai-research.md)**: Deep learning and AI research with GPU optimization

### [Life Sciences](life-sciences/)

- **[Genomics](life-sciences/genomics.md)**: Complete genomics analysis with optimized bioinformatics tools

### [Physical Sciences](physical-sciences/)

- **[Climate-Modeling](physical-sciences/climate-modeling.md)**: Weather prediction and climate simulation tools


## Quick Reference

| Domain Pack | Category | Typical Use Cases |
|-------------|----------|-------------------|
| [Ai-Research](computer-science/ai-research.md) | Computer Science | Distributed Training, Hyperparameter Tuning |
| [Genomics](life-sciences/genomics.md) | Life Sciences | Variant Calling, Rna Seq Analysis |
| [Climate-Modeling](physical-sciences/climate-modeling.md) | Physical Sciences | Weather Forecast, Climate Simulation |

## Getting Started

1. **Browse Domain Packs**: Explore the categories above to find domain packs for your research area
2. **Read Documentation**: Each domain pack has detailed documentation with examples
3. **Deploy Environment**: Use the CLI to deploy your chosen domain pack
4. **Run Workflows**: Execute pre-configured research workflows

```bash
# List all available domain packs
aws-research-wizard config list

# Get detailed information about a domain pack
aws-research-wizard config info genomics

# Deploy a research environment
aws-research-wizard deploy --domain genomics --size medium
```
