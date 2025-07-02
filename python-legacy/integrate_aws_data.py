#!/usr/bin/env python3
"""
AWS Open Data Integration Script
Updates domain configurations with real dataset workflows and demo data
"""

import os
import yaml
import json
import logging
from pathlib import Path
from typing import Dict, List, Any
from dataset_manager import DatasetManager
from config_loader import ConfigLoader

class AWSDataIntegrator:
    """Integrates AWS Open Data into research pack configurations"""

    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)

        # Initialize managers
        self.dataset_manager = DatasetManager(config_root)
        self.config_loader = ConfigLoader(config_root)

        # Domain mappings
        self.domain_mappings = {
            'genomics': 'genomics',
            'climate_modeling': 'climate_modeling',
            'machine_learning': 'machine_learning',
            'geospatial_research': 'geospatial_research',
            'agricultural_sciences': 'agricultural_sciences',
            'atmospheric_chemistry': 'atmospheric_chemistry',
            'cybersecurity_research': 'cybersecurity',
            'benchmarking_performance': 'benchmarking'
        }

    def update_domain_configurations(self) -> Dict[str, bool]:
        """Update all domain configurations with AWS Open Data workflows"""
        results = {}

        for config_name, domain_key in self.domain_mappings.items():
            try:
                self.logger.info(f"Updating {config_name} with AWS Open Data...")

                # Load existing configuration
                config_file = self.config_root / "domains" / f"{config_name}.yaml"
                if not config_file.exists():
                    self.logger.warning(f"Configuration file not found: {config_file}")
                    results[config_name] = False
                    continue

                with open(config_file, 'r') as f:
                    config_data = yaml.safe_load(f)

                # Generate workflows for this domain
                workflows = self.dataset_manager.generate_demo_workflows(domain_key)
                datasets = self.dataset_manager.get_datasets_for_domain(domain_key)

                # Update AWS data sources
                if datasets:
                    aws_data_sources = [
                        f"{dataset.name} - {dataset.description}"
                        for dataset in datasets
                    ]
                    config_data['aws_data_sources'] = aws_data_sources[:10]  # Limit to top 10

                # Update demo workflows with real data
                if workflows:
                    demo_workflows = []
                    for workflow in workflows:
                        demo_workflow = {
                            'name': workflow['name'],
                            'description': workflow['description'],
                            'dataset': workflow['dataset'],
                            'expected_runtime': workflow['expected_runtime'],
                            'cost_estimate': workflow['cost_estimate']
                        }
                        demo_workflows.append(demo_workflow)

                    # Keep existing demo workflows and add new ones
                    existing_workflows = config_data.get('demo_workflows', [])
                    config_data['demo_workflows'] = existing_workflows + demo_workflows

                # Add dataset integration metadata
                config_data['aws_integration'] = {
                    'datasets_available': len(datasets),
                    'demo_workflows_available': len(workflows),
                    'total_data_volume_tb': sum(dataset.size_tb or 0 for dataset in datasets),
                    'integration_date': '2023-12-01',
                    'data_access_patterns': {
                        'cost_optimized': 'Use S3 Intelligent Tiering',
                        'performance_optimized': 'Access from same AWS region',
                        'security': 'Data encrypted in transit and at rest'
                    }
                }

                # Write updated configuration
                with open(config_file, 'w') as f:
                    yaml.dump(config_data, f, default_flow_style=False, sort_keys=False, indent=2)

                self.logger.info(f"Updated {config_name}: {len(datasets)} datasets, {len(workflows)} workflows")
                results[config_name] = True

            except Exception as e:
                self.logger.error(f"Failed to update {config_name}: {e}")
                results[config_name] = False

        return results

    def create_data_access_guide(self) -> bool:
        """Create a comprehensive data access guide"""
        try:
            guide_file = self.config_root / "demo_data" / "data_access_guide.md"

            # Collect all datasets
            all_datasets = []
            for domain in self.domain_mappings.values():
                datasets = self.dataset_manager.get_datasets_for_domain(domain)
                all_datasets.extend(datasets)

            # Remove duplicates
            unique_datasets = {}
            for dataset in all_datasets:
                unique_datasets[dataset.name] = dataset

            guide_content = self._generate_access_guide_content(list(unique_datasets.values()))

            with open(guide_file, 'w') as f:
                f.write(guide_content)

            self.logger.info(f"Created data access guide: {guide_file}")
            return True

        except Exception as e:
            self.logger.error(f"Failed to create data access guide: {e}")
            return False

    def _generate_access_guide_content(self, datasets: List) -> str:
        """Generate markdown content for data access guide"""
        content = """# AWS Open Data Access Guide

This guide provides information on accessing and using AWS Open Data datasets for research workflows.

## Overview

The AWS Research Wizard integrates with the AWS Open Data Registry to provide access to petabytes of publicly available datasets. These datasets are hosted on AWS and can be accessed without egress charges when used within AWS.

## Available Datasets

"""

        # Group datasets by category
        categories = {}
        for dataset in datasets:
            category = "General"
            if hasattr(dataset, 'domains') and dataset.domains:
                category = dataset.domains[0].replace('_', ' ').title()

            if category not in categories:
                categories[category] = []
            categories[category].append(dataset)

        for category, category_datasets in sorted(categories.items()):
            content += f"### {category}\n\n"

            for dataset in sorted(category_datasets, key=lambda x: x.name):
                size_display = f"{dataset.size_tb} TB" if dataset.size_tb else "Size varies"
                content += f"#### {dataset.name}\n"
                content += f"- **Description**: {dataset.description}\n"
                content += f"- **Location**: `{dataset.location}`\n"
                content += f"- **Size**: {size_display}\n"
                content += f"- **Format**: {dataset.format}\n"
                content += f"- **Update Frequency**: {dataset.update_frequency}\n"
                if hasattr(dataset, 'temporal_coverage') and dataset.temporal_coverage:
                    content += f"- **Temporal Coverage**: {dataset.temporal_coverage}\n"
                if hasattr(dataset, 'spatial_resolution') and dataset.spatial_resolution:
                    content += f"- **Spatial Resolution**: {dataset.spatial_resolution}\n"
                content += "\n"

        content += """## Data Access Patterns

### Cost Optimization
- Use S3 Intelligent Tiering for datasets larger than 128KB
- Access data from the same AWS region to minimize transfer costs
- Some datasets use requester-pays model - check before large downloads

### Performance Optimization
- Use multipart downloads for files larger than 100MB
- Parallelize access across S3 prefixes for maximum throughput
- Implement local caching for frequently accessed data
- Use spot instances for batch processing to reduce costs

### Security Considerations
- Most datasets are publicly accessible
- Some datasets may require data use agreements
- All data is encrypted in transit and at rest
- Follow your organization's data governance policies

## Example Access Code

### Python with Boto3
```python
import boto3

# Initialize S3 client
s3 = boto3.client('s3')

# List objects in a dataset
response = s3.list_objects_v2(
    Bucket='dataset-bucket',
    Prefix='path/to/data/',
    MaxKeys=100
)

# Download a file
s3.download_file(
    Bucket='dataset-bucket',
    Key='path/to/file.nc',
    Filename='local_file.nc'
)
```

### AWS CLI
```bash
# List dataset contents
aws s3 ls s3://dataset-bucket/path/to/data/ --no-sign-request

# Download a file
aws s3 cp s3://dataset-bucket/path/to/file.nc local_file.nc --no-sign-request

# Sync a directory
aws s3 sync s3://dataset-bucket/path/to/data/ ./local_data/ --no-sign-request
```

## Cost Estimation

Use the AWS Simple Monthly Calculator to estimate costs:
- S3 storage: ~$0.023 per GB/month
- Data transfer: Free within same region, $0.09 per GB out to internet
- Compute: Varies by instance type and usage

## Support and Resources

- [AWS Open Data Registry](https://registry.opendata.aws/)
- [AWS Data Documentation](https://docs.aws.amazon.com/datasets/)
- [AWS Research Credits](https://aws.amazon.com/research-credits/)
- [AWS Cost Calculator](https://calculator.aws/)

## Getting Started

1. Set up AWS credentials or use IAM roles
2. Choose appropriate instance types for your workload
3. Test with small data subsets first
4. Implement proper error handling and retry logic
5. Monitor costs and usage with CloudWatch

For specific workflow examples, see the demo workflows in each research pack configuration.
"""

        return content

    def validate_integrations(self) -> Dict[str, Any]:
        """Validate all AWS data integrations"""
        validation_results = {
            'configurations_updated': 0,
            'total_datasets': 0,
            'total_workflows': 0,
            'domains_with_data': [],
            'validation_errors': []
        }

        try:
            # Check each domain configuration
            for config_name in self.domain_mappings.keys():
                config_file = self.config_root / "domains" / f"{config_name}.yaml"

                if config_file.exists():
                    with open(config_file, 'r') as f:
                        config_data = yaml.safe_load(f)

                    # Check for AWS integration metadata
                    if 'aws_integration' in config_data:
                        validation_results['configurations_updated'] += 1
                        validation_results['total_datasets'] += config_data['aws_integration'].get('datasets_available', 0)
                        validation_results['total_workflows'] += config_data['aws_integration'].get('demo_workflows_available', 0)
                        validation_results['domains_with_data'].append(config_name)

                    # Validate demo workflows have required fields
                    demo_workflows = config_data.get('demo_workflows', [])
                    for workflow in demo_workflows:
                        required_fields = ['name', 'description', 'expected_runtime', 'cost_estimate']
                        missing_fields = [field for field in required_fields if field not in workflow]
                        if missing_fields:
                            validation_results['validation_errors'].append(
                                f"{config_name}: Workflow '{workflow.get('name', 'unknown')}' missing fields: {missing_fields}"
                            )

        except Exception as e:
            validation_results['validation_errors'].append(f"Validation error: {e}")

        return validation_results

    def export_integration_summary(self, output_file: str) -> bool:
        """Export integration summary to JSON"""
        try:
            summary = {
                'integration_summary': self.validate_integrations(),
                'domain_statistics': {},
                'dataset_catalog': []
            }

            # Collect domain statistics
            for config_name, domain_key in self.domain_mappings.items():
                datasets = self.dataset_manager.get_datasets_for_domain(domain_key)
                workflows = self.dataset_manager.generate_demo_workflows(domain_key)

                summary['domain_statistics'][config_name] = {
                    'datasets_available': len(datasets),
                    'demo_workflows': len(workflows),
                    'total_data_volume_tb': sum(dataset.size_tb or 0 for dataset in datasets),
                    'primary_datasets': [dataset.name for dataset in datasets[:5]]
                }

            # Create dataset catalog
            all_datasets = self.dataset_manager.list_available_datasets()
            summary['dataset_catalog'] = all_datasets

            with open(output_file, 'w') as f:
                json.dump(summary, f, indent=2, default=str)

            self.logger.info(f"Exported integration summary to {output_file}")
            return True

        except Exception as e:
            self.logger.error(f"Failed to export integration summary: {e}")
            return False


def main():
    """Main execution function"""
    import argparse

    parser = argparse.ArgumentParser(description="AWS Open Data Integration")
    parser.add_argument("--update-configs", action="store_true", help="Update domain configurations")
    parser.add_argument("--create-guide", action="store_true", help="Create data access guide")
    parser.add_argument("--validate", action="store_true", help="Validate integrations")
    parser.add_argument("--export-summary", type=str, help="Export integration summary to file")
    parser.add_argument("--config-root", type=str, default="configs", help="Configuration root directory")

    args = parser.parse_args()

    # Setup logging
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

    # Initialize integrator
    integrator = AWSDataIntegrator(args.config_root)

    if args.update_configs:
        print("Updating domain configurations with AWS Open Data...")
        results = integrator.update_domain_configurations()

        success_count = sum(1 for success in results.values() if success)
        total_count = len(results)

        print(f"Updated {success_count}/{total_count} configurations successfully:")
        for domain, success in results.items():
            status = "✅" if success else "❌"
            print(f"  {status} {domain}")

    if args.create_guide:
        print("Creating data access guide...")
        success = integrator.create_data_access_guide()
        print("✅ Data access guide created" if success else "❌ Failed to create guide")

    if args.validate:
        print("Validating AWS data integrations...")
        results = integrator.validate_integrations()

        print(f"Validation Results:")
        print(f"  Configurations updated: {results['configurations_updated']}")
        print(f"  Total datasets available: {results['total_datasets']}")
        print(f"  Total demo workflows: {results['total_workflows']}")
        print(f"  Domains with data: {len(results['domains_with_data'])}")

        if results['validation_errors']:
            print("  Validation errors:")
            for error in results['validation_errors']:
                print(f"    - {error}")
        else:
            print("  ✅ No validation errors found")

    if args.export_summary:
        print(f"Exporting integration summary to {args.export_summary}...")
        success = integrator.export_integration_summary(args.export_summary)
        print("✅ Summary exported" if success else "❌ Failed to export summary")


if __name__ == "__main__":
    main()
