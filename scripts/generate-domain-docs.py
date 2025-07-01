#!/usr/bin/env python3
"""
Domain Pack Documentation Generator
Automatically generates documentation for all domain packs
"""

import os
import sys
import yaml
import json
from pathlib import Path
from typing import Dict, List, Any, Optional
import logging

class DomainDocsGenerator:
    def __init__(self, domain_packs_dir: str = "domain-packs", docs_dir: str = "docs"):
        self.domain_packs_dir = Path(domain_packs_dir)
        self.docs_dir = Path(docs_dir)
        self.output_dir = self.docs_dir / "domain-packs"

        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)

    def find_domain_packs(self) -> Dict[str, List[Path]]:
        """Find all domain packs organized by category"""
        categories = {}
        domains_dir = self.domain_packs_dir / "domains"

        for category_dir in domains_dir.iterdir():
            if category_dir.is_dir():
                category_name = category_dir.name
                categories[category_name] = []

                for domain_dir in category_dir.iterdir():
                    if domain_dir.is_dir():
                        config_file = domain_dir / "domain-pack.yaml"
                        if config_file.exists():
                            categories[category_name].append(config_file)

        return categories

    def load_domain_config(self, config_file: Path) -> Optional[Dict[str, Any]]:
        """Load domain pack configuration"""
        try:
            with open(config_file, 'r') as f:
                return yaml.safe_load(f)
        except Exception as e:
            self.logger.error(f"Error loading {config_file}: {e}")
            return None

    def generate_category_index(self, category: str, domain_configs: List[Dict[str, Any]]) -> str:
        """Generate category index page"""
        category_title = category.replace('-', ' ').title()

        content = f"""# {category_title}

Research domain packs for {category_title.lower()} research computing.

## Available Domain Packs

"""

        for config in sorted(domain_configs, key=lambda x: x['name']):
            name = config['name']
            description = config['description']
            cost_range = config.get('cost_estimates', {}).get('medium_workload', 'Contact for pricing')

            content += f"""### [{config['name'].title()}]({name}.md)

{description}

- **Typical Cost**: {cost_range}
- **Instance Types**: {', '.join(config.get('aws_config', {}).get('instance_types', {}).values())}
- **Categories**: {', '.join(config.get('categories', []))}

"""

        return content

    def generate_domain_page(self, config: Dict[str, Any], config_file: Path) -> str:
        """Generate individual domain pack documentation page"""
        name = config['name']
        title = name.replace('-', ' ').title()
        description = config['description']

        content = f"""# {title}

{description}

## Overview

**Domain Pack**: `{name}`
**Version**: {config.get('version', '1.0.0')}
**Categories**: {', '.join(config.get('categories', []))}
**Maintainers**: {', '.join([m['name'] for m in config.get('maintainers', [])])}

## Quick Start

```bash
# Deploy this domain pack
aws-research-wizard deploy --domain {name} --size medium

# Get detailed information
aws-research-wizard config info {name}

# List available workflows
aws-research-wizard workflow list --domain {name}
```

## Software Stack

### Core Packages
"""

        # Add software packages
        spack_packages = config.get('spack_config', {}).get('packages', [])
        if spack_packages:
            for package in spack_packages[:10]:  # Show first 10
                package_name = package.split('@')[0].split('%')[0]
                content += f"- **{package_name}**: {package}\n"

            if len(spack_packages) > 10:
                content += f"\n*And {len(spack_packages) - 10} more packages...*\n"

        content += f"""
### Optimization Settings
- **Compiler**: {config.get('spack_config', {}).get('compiler', 'gcc@11.4.0')}
- **Target Architecture**: {config.get('spack_config', {}).get('target', 'x86_64_v3')}
- **Optimization Flags**: {config.get('spack_config', {}).get('optimization', '-O3')}

## AWS Infrastructure

### Instance Types
"""

        # Add AWS configuration
        aws_config = config.get('aws_config', {})
        instance_types = aws_config.get('instance_types', {})

        for size, instance_type in instance_types.items():
            content += f"- **{size.title()}**: `{instance_type}`\n"

        storage = aws_config.get('storage', {})
        if storage:
            content += f"""
### Storage Configuration
- **Type**: {storage.get('type', 'gp3')}
- **Size**: {storage.get('size_gb', 500)} GB
- **IOPS**: {storage.get('iops', 3000)}
- **Throughput**: {storage.get('throughput', 125)} MB/s
"""

        # Add workflows
        workflows = config.get('workflows', [])
        if workflows:
            content += f"""
## Research Workflows

This domain pack includes {len(workflows)} pre-configured research workflows:

"""
            for workflow in workflows:
                content += f"""### {workflow['name'].replace('_', ' ').title()}

{workflow['description']}

```bash
# Run this workflow
aws-research-wizard workflow run {workflow['name']} --domain {name}
```

- **Input Data**: {workflow.get('input_data', 'User-provided')}
- **Expected Output**: {workflow.get('expected_output', 'Processed results')}

"""

        # Add cost estimates
        cost_estimates = config.get('cost_estimates', {})
        if cost_estimates:
            content += f"""
## Cost Estimates

| Workload Size | Estimated Daily Cost |
|---------------|---------------------|
| Small | {cost_estimates.get('small_workload', 'Contact for pricing')} |
| Medium | {cost_estimates.get('medium_workload', 'Contact for pricing')} |
| Large | {cost_estimates.get('large_workload', 'Contact for pricing')} |

!!! note "Cost Optimization"
    These estimates assume on-demand pricing. Significant savings are possible with:

    - **Spot Instances**: 70-90% savings for fault-tolerant workloads
    - **Reserved Instances**: 30-60% savings for predictable usage
    - **Savings Plans**: 20-72% savings with flexible commitment
"""

        # Add example configuration
        content += f"""
## Example Configuration

```yaml
# {name}-research-config.yaml
domain: {name}
size: medium
aws:
  region: us-east-1
  availability_zone: us-east-1a
compute:
  instance_type: {instance_types.get('medium', 'c6i.4xlarge')}
  instance_count: 1
storage:
  type: {storage.get('type', 'gp3')}
  size_gb: {storage.get('size_gb', 500)}
```

## Getting Help

- 📖 **Domain-Specific Documentation**: [tutorials/{name}/](../../tutorials/{name}/)
- 💬 **Community Support**: [GitHub Discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)
- 🐛 **Issues**: [GitHub Issues](https://github.com/aws-research-wizard/aws-research-wizard/issues)

## Related Domain Packs

"""

        # Add related domains from same category
        categories = config.get('categories', [])
        if categories:
            content += f"Other domain packs in **{categories[0].replace('-', ' ').title()}**:\n\n"
            # This would be filled in during generation when we have access to all domains

        return content

    def generate_all_docs(self):
        """Generate documentation for all domain packs"""
        self.logger.info("🔍 Generating domain pack documentation...")

        # Create output directory
        self.output_dir.mkdir(parents=True, exist_ok=True)

        # Find all domain packs
        categories = self.find_domain_packs()

        if not categories:
            self.logger.warning("⚠️  No domain packs found!")
            return

        total_domains = sum(len(domains) for domains in categories.values())
        processed = 0

        # Generate documentation for each category
        for category, config_files in categories.items():
            self.logger.info(f"📝 Processing category: {category}")

            # Create category directory
            category_dir = self.output_dir / category
            category_dir.mkdir(exist_ok=True)

            # Load all configurations for this category
            domain_configs = []
            for config_file in config_files:
                config = self.load_domain_config(config_file)
                if config:
                    domain_configs.append(config)

                    # Generate individual domain page
                    domain_content = self.generate_domain_page(config, config_file)
                    domain_file = category_dir / f"{config['name']}.md"

                    with open(domain_file, 'w') as f:
                        f.write(domain_content)

                    processed += 1
                    self.logger.info(f"  ✅ Generated {config['name']}.md")

            # Generate category index
            if domain_configs:
                category_content = self.generate_category_index(category, domain_configs)
                category_index = category_dir / "index.md"

                with open(category_index, 'w') as f:
                    f.write(category_content)

                self.logger.info(f"  ✅ Generated {category}/index.md")

        # Generate main domain packs index
        self.generate_main_index(categories)

        self.logger.info(f"")
        self.logger.info(f"📊 Documentation Generation Summary:")
        self.logger.info(f"   Categories: {len(categories)}")
        self.logger.info(f"   Domain Packs: {processed}")
        self.logger.info(f"   Output Directory: {self.output_dir}")
        self.logger.info(f"✅ Domain pack documentation generated successfully!")

    def generate_main_index(self, categories: Dict[str, List[Path]]):
        """Generate main domain packs index page"""
        content = """# Domain Packs

AWS Research Wizard provides pre-configured domain packs for major research disciplines. Each domain pack includes optimized software stacks, AWS infrastructure configurations, and sample workflows.

## Categories

"""

        for category, config_files in categories.items():
            category_title = category.replace('-', ' ').title()
            content += f"### [{category_title}]({category}/)\n\n"

            # Add brief description for each domain in category
            for config_file in sorted(config_files)[:3]:  # Show first 3
                config = self.load_domain_config(config_file)
                if config:
                    name = config['name']
                    description = config['description']
                    content += f"- **[{name.title()}]({category}/{name}.md)**: {description}\n"

            if len(config_files) > 3:
                content += f"- *And {len(config_files) - 3} more...*\n"

            content += "\n"

        content += """
## Quick Reference

| Domain Pack | Category | Typical Use Cases |
|-------------|----------|-------------------|
"""

        # Add quick reference table
        for category, config_files in categories.items():
            for config_file in config_files:
                config = self.load_domain_config(config_file)
                if config:
                    name = config['name']
                    category_name = category.replace('-', ' ').title()
                    workflows = config.get('workflows', [])
                    use_cases = ', '.join([w['name'].replace('_', ' ').title() for w in workflows[:2]])
                    if len(workflows) > 2:
                        use_cases += "..."

                    content += f"| [{name.title()}]({category}/{name}.md) | {category_name} | {use_cases} |\n"

        content += """
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
"""

        index_file = self.output_dir / "index.md"
        with open(index_file, 'w') as f:
            f.write(content)

        self.logger.info(f"✅ Generated main index: {index_file}")

def main():
    import argparse

    parser = argparse.ArgumentParser(description="Generate domain pack documentation")
    parser.add_argument("--domain-packs-dir", default="domain-packs",
                       help="Domain packs directory (default: domain-packs)")
    parser.add_argument("--docs-dir", default="docs",
                       help="Documentation output directory (default: docs)")

    args = parser.parse_args()

    generator = DomainDocsGenerator(args.domain_packs_dir, args.docs_dir)
    generator.generate_all_docs()

if __name__ == "__main__":
    main()
