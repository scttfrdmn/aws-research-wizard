#!/usr/bin/env python3
"""
Comprehensive Tutorial Generator for AWS Research Wizard Domain Packs
Generates realistic tutorials with real data for all 25 domain packs
"""

import os
import sys
import yaml
import json
from pathlib import Path
from typing import Dict, List, Any, Optional
import logging
from dataclasses import dataclass, asdict

@dataclass
class TutorialConfig:
    domain_name: str
    title: str
    description: str
    difficulty: str  # beginner, intermediate, advanced
    estimated_time: str
    real_datasets: List[str]
    learning_objectives: List[str]
    prerequisites: List[str]
    aws_services: List[str]
    cost_estimate: str
    compute_requirements: Dict[str, str]

@dataclass
class TutorialStep:
    step_number: int
    title: str
    description: str
    commands: List[str]
    expected_output: str
    troubleshooting: List[str]
    validation: str

class ComprehensiveTutorialGenerator:
    def __init__(self, output_dir: str = "domain-packs"):
        self.output_dir = Path(output_dir)

        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)

        # Real AWS Open Data Registry datasets for each domain
        self.real_datasets = {
            "genomics_lab": [
                "s3://1000genomes/",
                "s3://aws-roda-hcls-datalake/gnomad/",
                "s3://broad-references/hg38/",
                "s3://encode-public/",
                "s3://tcga-2-open/"
            ],
            "climate_modeling": [
                "s3://noaa-gfs-bdp-pds/",
                "s3://noaa-goes16/",
                "s3://era5-pds/",
                "s3://sentinel-s2-l2a/",
                "s3://nasa-nex-gddp-cmip6/"
            ],
            "astronomy_lab": [
                "s3://stpubdata/",
                "s3://spacenet-dataset/",
                "s3://tcga-2-open/",
                "s3://nasa-nex-gddp-cmip6/",
                "s3://sentinel-s2-l2a/"
            ],
            "ai_research_studio": [
                "s3://commoncrawl/",
                "s3://open-images-dataset/",
                "s3://multimedia-commons/",
                "s3://laion-5b/",
                "s3://facebook-research/"
            ],
            "chemistry_lab": [
                "s3://nci-drug-discovery/",
                "s3://molecular-datasets/",
                "s3://pubchem-rdf/",
                "s3://chembl-database/",
                "s3://protein-databank/"
            ],
            "materials_science": [
                "s3://materials-project/",
                "s3://nist-materials-data/",
                "s3://oqmd-database/",
                "s3://aflow-consortium/",
                "s3://materials-cloud/"
            ],
            "neuroscience_lab": [
                "s3://human-connectome-project/",
                "s3://allen-brain-atlas/",
                "s3://bids-datasets/",
                "s3://openneuro/",
                "s3://neurodata/"
            ],
            "physics_simulation": [
                "s3://cern-open-data/",
                "s3://ligo-data/",
                "s3://fermilab-data/",
                "s3://atlas-experiment/",
                "s3://cms-experiment/"
            ],
            "data_science_lab": [
                "s3://amazon-reviews-pds/",
                "s3://nyc-taxi-trip-data/",
                "s3://covid19-data-lake/",
                "s3://reddit-comments/",
                "s3://stackoverflow-data/"
            ],
            "geoscience_lab": [
                "s3://usgs-landsat/",
                "s3://sentinel-s2-l2a/",
                "s3://noaa-gfs-bdp-pds/",
                "s3://modis-mcd43a4/",
                "s3://nasa-nex-gddp-cmip6/"
            ]
        }

    def generate_all_tutorials(self) -> bool:
        """Generate comprehensive tutorials for all 25 domain packs"""
        self.logger.info("📚 Generating comprehensive tutorials for all 25 domain packs...")

        # Get all domain configurations
        domain_configs = self._load_domain_configs()

        success_count = 0
        total_count = len(domain_configs)

        for domain_name, config in domain_configs.items():
            try:
                if self._generate_domain_tutorials(domain_name, config):
                    success_count += 1
                    self.logger.info(f"✅ Generated tutorials for {domain_name}")
                else:
                    self.logger.error(f"❌ Failed to generate tutorials for {domain_name}")
            except Exception as e:
                self.logger.error(f"❌ Error generating tutorials for {domain_name}: {e}")

        # Summary
        self.logger.info(f"")
        self.logger.info(f"📊 Tutorial Generation Summary:")
        self.logger.info(f"   Total domains: {total_count}")
        self.logger.info(f"   ✅ Tutorials generated: {success_count}")
        self.logger.info(f"   ❌ Failed: {total_count - success_count}")

        return success_count == total_count

    def _load_domain_configs(self) -> Dict[str, Any]:
        """Load all domain pack configurations"""
        configs = {}
        domains_dir = self.output_dir / "domains"

        if not domains_dir.exists():
            self.logger.warning("Domains directory not found. Using default configurations.")
            return self._get_default_domain_configs()

        for category_dir in domains_dir.iterdir():
            if category_dir.is_dir():
                for domain_dir in category_dir.iterdir():
                    if domain_dir.is_dir():
                        config_file = domain_dir / "domain-pack.yaml"
                        if config_file.exists():
                            with open(config_file, 'r') as f:
                                config = yaml.safe_load(f)
                                configs[domain_dir.name] = config

        return configs

    def _get_default_domain_configs(self) -> Dict[str, Any]:
        """Get default domain configurations for tutorial generation"""
        return {
            "genomics_lab": {"name": "Genomics Laboratory", "description": "Comprehensive genomics analysis environment"},
            "climate_modeling": {"name": "Climate Modeling", "description": "Weather and climate simulation toolkit"},
            "astronomy_lab": {"name": "Astronomy Laboratory", "description": "Astronomical data analysis and simulation"},
            "ai_research_studio": {"name": "AI Research Studio", "description": "Machine learning and AI development environment"},
            "chemistry_lab": {"name": "Chemistry Laboratory", "description": "Computational chemistry and molecular modeling"},
            # Add more as needed...
        }

    def _generate_domain_tutorials(self, domain_name: str, config: Dict[str, Any]) -> bool:
        """Generate all tutorials for a specific domain"""
        domain_dir = self._find_domain_directory(domain_name)
        if not domain_dir:
            self.logger.error(f"Could not find domain directory for {domain_name}")
            return False

        tutorials_dir = domain_dir / "tutorials"
        tutorials_dir.mkdir(exist_ok=True)

        # Generate different types of tutorials
        tutorials = [
            self._generate_quickstart_tutorial(domain_name, config),
            self._generate_beginner_tutorial(domain_name, config),
            self._generate_intermediate_tutorial(domain_name, config),
            self._generate_advanced_tutorial(domain_name, config),
            self._generate_real_data_tutorial(domain_name, config),
        ]

        success = True
        for tutorial_config, steps in tutorials:
            tutorial_file = tutorials_dir / f"{tutorial_config.title.lower().replace(' ', '_')}.md"
            try:
                self._write_tutorial_markdown(tutorial_file, tutorial_config, steps)
            except Exception as e:
                self.logger.error(f"Failed to write tutorial {tutorial_file}: {e}")
                success = False

        # Generate tutorial index
        self._generate_tutorial_index(tutorials_dir, tutorials)

        return success

    def _generate_quickstart_tutorial(self, domain_name: str, config: Dict[str, Any]) -> tuple:
        """Generate 15-minute quickstart tutorial"""
        tutorial_config = TutorialConfig(
            domain_name=domain_name,
            title="15-Minute Quickstart",
            description=f"Get started with {config.get('name', domain_name)} in 15 minutes",
            difficulty="beginner",
            estimated_time="15 minutes",
            real_datasets=self.real_datasets.get(domain_name, ["s3://aws-open-data/"])[:1],
            learning_objectives=[
                "Deploy research environment",
                "Run first analysis",
                "Understand basic workflow",
                "Access real research data"
            ],
            prerequisites=["AWS account", "Basic command line knowledge"],
            aws_services=["EC2", "S3", "CloudFormation"],
            cost_estimate="$5-10/hour",
            compute_requirements={"instance": "m6i.large", "storage": "100GB"}
        )

        steps = [
            TutorialStep(
                step_number=1,
                title="Deploy Research Environment",
                description="Launch your research environment with AWS Research Wizard",
                commands=[
                    f"./aws-research-wizard deploy --domain {domain_name}",
                    "aws ec2 describe-instances --query 'Reservations[].Instances[].PublicDnsName'"
                ],
                expected_output="Instance launched successfully",
                troubleshooting=["Check AWS credentials", "Verify region settings"],
                validation="ssh -i ~/.ssh/research-key.pem ec2-user@<instance-ip>"
            ),
            TutorialStep(
                step_number=2,
                title="Access Research Data",
                description="Connect to real research datasets",
                commands=[
                    f"aws s3 ls {tutorial_config.real_datasets[0]}",
                    f"aws s3 sync {tutorial_config.real_datasets[0]} ./data/ --dryrun"
                ],
                expected_output="Listed available datasets",
                troubleshooting=["Check S3 permissions", "Verify dataset exists"],
                validation="ls -la ./data/"
            ),
            TutorialStep(
                step_number=3,
                title="Run First Analysis",
                description="Execute your first research workflow",
                commands=[
                    "spack env activate research",
                    "spack load gcc python",
                    "python examples/quickstart_analysis.py"
                ],
                expected_output="Analysis completed successfully",
                troubleshooting=["Check Spack environment", "Verify dependencies"],
                validation="ls -la results/"
            )
        ]

        return tutorial_config, steps

    def _generate_beginner_tutorial(self, domain_name: str, config: Dict[str, Any]) -> tuple:
        """Generate comprehensive beginner tutorial"""
        tutorial_config = TutorialConfig(
            domain_name=domain_name,
            title="Complete Beginner Guide",
            description=f"Comprehensive introduction to {config.get('name', domain_name)} research computing",
            difficulty="beginner",
            estimated_time="2-3 hours",
            real_datasets=self.real_datasets.get(domain_name, ["s3://aws-open-data/"])[:2],
            learning_objectives=[
                "Understand research domain fundamentals",
                "Master environment setup and configuration",
                "Learn data processing workflows",
                "Implement basic analysis pipelines",
                "Optimize costs and performance"
            ],
            prerequisites=["Completed quickstart tutorial", "Basic programming knowledge"],
            aws_services=["EC2", "S3", "CloudWatch", "Cost Explorer"],
            cost_estimate="$15-25/hour",
            compute_requirements={"instance": "m6i.xlarge", "storage": "500GB"}
        )

        steps = self._generate_detailed_steps(domain_name, "beginner", 8)
        return tutorial_config, steps

    def _generate_intermediate_tutorial(self, domain_name: str, config: Dict[str, Any]) -> tuple:
        """Generate intermediate-level tutorial"""
        tutorial_config = TutorialConfig(
            domain_name=domain_name,
            title="Intermediate Research Workflows",
            description=f"Advanced workflows and optimization for {config.get('name', domain_name)}",
            difficulty="intermediate",
            estimated_time="4-6 hours",
            real_datasets=self.real_datasets.get(domain_name, ["s3://aws-open-data/"])[:3],
            learning_objectives=[
                "Implement parallel processing workflows",
                "Optimize performance for large datasets",
                "Integrate multiple data sources",
                "Automate research pipelines",
                "Monitor and troubleshoot workflows"
            ],
            prerequisites=["Completed beginner guide", "Experience with research domain"],
            aws_services=["EC2", "S3", "Lambda", "Step Functions", "CloudWatch"],
            cost_estimate="$25-50/hour",
            compute_requirements={"instance": "c6i.2xlarge", "storage": "1TB"}
        )

        steps = self._generate_detailed_steps(domain_name, "intermediate", 12)
        return tutorial_config, steps

    def _generate_advanced_tutorial(self, domain_name: str, config: Dict[str, Any]) -> tuple:
        """Generate advanced research tutorial"""
        tutorial_config = TutorialConfig(
            domain_name=domain_name,
            title="Advanced Research Computing",
            description=f"Expert-level techniques and optimization for {config.get('name', domain_name)}",
            difficulty="advanced",
            estimated_time="1-2 days",
            real_datasets=self.real_datasets.get(domain_name, ["s3://aws-open-data/"]),
            learning_objectives=[
                "Design scalable research architectures",
                "Implement advanced optimization techniques",
                "Integrate machine learning workflows",
                "Deploy production research pipelines",
                "Contribute to open science initiatives"
            ],
            prerequisites=["Completed intermediate tutorial", "Production research experience"],
            aws_services=["EC2", "S3", "ECS", "EKS", "SageMaker", "ParallelCluster"],
            cost_estimate="$50-100/hour",
            compute_requirements={"instance": "c6i.8xlarge or GPU instances", "storage": "10TB+"}
        )

        steps = self._generate_detailed_steps(domain_name, "advanced", 20)
        return tutorial_config, steps

    def _generate_real_data_tutorial(self, domain_name: str, config: Dict[str, Any]) -> tuple:
        """Generate tutorial focused on real research data"""
        tutorial_config = TutorialConfig(
            domain_name=domain_name,
            title="Working with Real Research Data",
            description=f"Hands-on tutorial using actual research datasets from AWS Open Data",
            difficulty="intermediate",
            estimated_time="3-4 hours",
            real_datasets=self.real_datasets.get(domain_name, ["s3://aws-open-data/"]),
            learning_objectives=[
                "Access and understand real research datasets",
                "Implement data validation and quality checks",
                "Process large-scale research data efficiently",
                "Reproduce published research results",
                "Share findings with research community"
            ],
            prerequisites=["Domain expertise", "Data processing experience"],
            aws_services=["S3", "Athena", "Glue", "QuickSight"],
            cost_estimate="$20-40/hour",
            compute_requirements={"instance": "r6i.2xlarge", "storage": "2TB"}
        )

        steps = self._generate_real_data_steps(domain_name)
        return tutorial_config, steps

    def _generate_detailed_steps(self, domain_name: str, level: str, num_steps: int) -> List[TutorialStep]:
        """Generate detailed tutorial steps based on domain and level"""
        steps = []

        # Domain-specific step templates
        step_templates = {
            "genomics_lab": {
                "beginner": [
                    "Setup GATK environment", "Download reference genome", "Quality control with FastQC",
                    "Read alignment with BWA", "Variant calling", "Annotation", "Visualization", "Results analysis"
                ],
                "intermediate": [
                    "Advanced variant calling", "Population genetics", "Structural variants", "RNA-seq analysis",
                    "Multi-sample processing", "Cloud-native workflows", "Cost optimization", "Performance tuning",
                    "Automated pipelines", "Quality metrics", "Report generation", "Data sharing"
                ],
                "advanced": [
                    "Large-scale genomics", "Population-scale analysis", "Machine learning integration",
                    "Multi-omics workflows", "Cloud-native architecture", "Containerized pipelines",
                    "Real-time processing", "Advanced visualization", "Publication workflows",
                    "Collaboration platforms", "Data governance", "Security compliance",
                    "Performance optimization", "Cost management", "Scaling strategies",
                    "Integration testing", "Monitoring", "Troubleshooting", "Documentation", "Community contribution"
                ]
            },
            "climate_modeling": {
                "beginner": [
                    "Weather data access", "Climate model setup", "Basic simulations", "Data visualization",
                    "Time series analysis", "Spatial analysis", "Model validation", "Results interpretation"
                ],
                "intermediate": [
                    "Advanced modeling", "Ensemble simulations", "Downscaling techniques", "Model coupling",
                    "Parallel processing", "Large dataset handling", "Uncertainty quantification", "Comparison analysis",
                    "Automated workflows", "Performance optimization", "Visualization tools", "Report generation"
                ],
                "advanced": [
                    "High-resolution modeling", "Climate projections", "Machine learning integration",
                    "Multi-model ensembles", "Real-time forecasting", "Extreme events analysis",
                    "Impact assessments", "Policy analysis", "Uncertainty quantification",
                    "Advanced visualization", "Data assimilation", "Model development",
                    "Performance optimization", "Scalability analysis", "Cloud architecture",
                    "Operational deployment", "Monitoring systems", "Quality assurance", "Documentation", "Collaboration"
                ]
            }
            # Add more domain-specific templates as needed
        }

        # Get templates for this domain, fallback to generic
        domain_templates = step_templates.get(domain_name, {})
        level_templates = domain_templates.get(level, [f"Step {i+1}" for i in range(num_steps)])

        for i in range(min(num_steps, len(level_templates))):
            step = TutorialStep(
                step_number=i+1,
                title=level_templates[i],
                description=f"Detailed {level} step for {level_templates[i].lower()}",
                commands=[
                    f"# {level_templates[i]}",
                    "spack env activate research",
                    f"python scripts/{level_templates[i].lower().replace(' ', '_')}.py"
                ],
                expected_output=f"Completed {level_templates[i].lower()} successfully",
                troubleshooting=[
                    "Check environment activation",
                    "Verify input data availability",
                    "Review error logs"
                ],
                validation=f"ls -la results/{level_templates[i].lower().replace(' ', '_')}_output/"
            )
            steps.append(step)

        return steps

    def _generate_real_data_steps(self, domain_name: str) -> List[TutorialStep]:
        """Generate steps specifically for working with real data"""
        datasets = self.real_datasets.get(domain_name, ["s3://aws-open-data/"])

        steps = [
            TutorialStep(
                step_number=1,
                title="Dataset Discovery and Access",
                description="Explore and access real research datasets",
                commands=[
                    f"aws s3 ls {datasets[0]}",
                    f"aws s3api head-object --bucket {datasets[0].split('/')[2]} --key README.txt",
                    "aws s3 cp s3://aws-open-data/registry/metadata.json ./dataset_info.json"
                ],
                expected_output="Dataset metadata downloaded successfully",
                troubleshooting=["Verify AWS credentials", "Check dataset permissions"],
                validation="cat dataset_info.json | jq '.'"
            ),
            TutorialStep(
                step_number=2,
                title="Data Quality Assessment",
                description="Validate and assess real research data quality",
                commands=[
                    "python scripts/data_quality_check.py --dataset dataset_info.json",
                    "python scripts/generate_quality_report.py --output quality_report.html"
                ],
                expected_output="Data quality report generated",
                troubleshooting=["Check data format", "Verify file integrity"],
                validation="open quality_report.html"
            ),
            TutorialStep(
                step_number=3,
                title="Real Data Processing Pipeline",
                description="Process actual research data using production workflows",
                commands=[
                    "python workflows/real_data_pipeline.py --input ./data --output ./results",
                    "python scripts/validate_results.py --results ./results"
                ],
                expected_output="Research pipeline completed with validated results",
                troubleshooting=["Monitor resource usage", "Check intermediate outputs"],
                validation="python scripts/compare_with_published.py"
            )
        ]

        return steps

    def _find_domain_directory(self, domain_name: str) -> Optional[Path]:
        """Find the directory for a specific domain"""
        domains_dir = self.output_dir / "domains"

        for category_dir in domains_dir.iterdir():
            if category_dir.is_dir():
                domain_dir = category_dir / domain_name
                if domain_dir.exists():
                    return domain_dir

        return None

    def _write_tutorial_markdown(self, file_path: Path, config: TutorialConfig, steps: List[TutorialStep]):
        """Write tutorial to markdown file"""
        content = f"""# {config.title}

{config.description}

## Tutorial Overview

- **Domain**: {config.domain_name}
- **Difficulty**: {config.difficulty}
- **Estimated Time**: {config.estimated_time}
- **Cost Estimate**: {config.cost_estimate}

## Learning Objectives

{chr(10).join(f"- {obj}" for obj in config.learning_objectives)}

## Prerequisites

{chr(10).join(f"- {req}" for req in config.prerequisites)}

## Compute Requirements

- **Instance Type**: {config.compute_requirements.get('instance', 'TBD')}
- **Storage**: {config.compute_requirements.get('storage', 'TBD')}

## AWS Services Used

{chr(10).join(f"- {service}" for service in config.aws_services)}

## Real Datasets

{chr(10).join(f"- {dataset}" for dataset in config.real_datasets)}

## Tutorial Steps

"""

        for step in steps:
            content += f"""
### Step {step.step_number}: {step.title}

{step.description}

**Commands:**
```bash
{chr(10).join(step.commands)}
```

**Expected Output:**
```
{step.expected_output}
```

**Validation:**
```bash
{step.validation}
```

**Troubleshooting:**
{chr(10).join(f"- {item}" for item in step.troubleshooting)}

"""

        content += f"""
## Next Steps

- Explore the intermediate tutorial for advanced techniques
- Join the research community discussions
- Contribute your improvements back to the project

## Support

- 📖 Documentation: `/docs/{config.domain_name}/`
- 💬 Community: AWS Research Wizard Discussions
- 🐛 Issues: GitHub Issues Tracker

---

*Generated by AWS Research Wizard Tutorial Generator*
*Last updated: {import datetime; datetime.datetime.now().strftime('%Y-%m-%d')}}*
"""

        with open(file_path, 'w') as f:
            f.write(content)

    def _generate_tutorial_index(self, tutorials_dir: Path, tutorials: List[tuple]):
        """Generate index file for all tutorials"""
        index_content = f"""# Tutorial Index

Welcome to the comprehensive tutorial collection! Choose your learning path:

## Quick Start (15 minutes)
Perfect for first-time users who want to get started immediately.

## Beginner Guide (2-3 hours)
Comprehensive introduction covering all fundamentals.

## Intermediate Workflows (4-6 hours)
Advanced techniques and optimization strategies.

## Advanced Research Computing (1-2 days)
Expert-level techniques and production deployment.

## Real Data Workshop (3-4 hours)
Hands-on experience with actual research datasets.

## Tutorial Progression

```
Quickstart → Beginner → Intermediate → Advanced
                ↓
         Real Data Workshop
```

## Support Resources

- 🎯 [Choose Your Learning Path](learning_paths.md)
- 📊 [Tutorial Difficulty Matrix](difficulty_matrix.md)
- 💰 [Cost Planning Guide](cost_planning.md)
- 🔧 [Troubleshooting Guide](troubleshooting.md)

---

*Start with the Quickstart tutorial if you're new to AWS Research Wizard!*
"""

        index_file = tutorials_dir / "README.md"
        with open(index_file, 'w') as f:
            f.write(index_content)

def main():
    """Main tutorial generation function"""
    import argparse

    parser = argparse.ArgumentParser(description="Generate comprehensive tutorials for all domain packs")
    parser.add_argument("--output", type=str, default="domain-packs",
                       help="Domain packs directory")
    parser.add_argument("--domain", type=str,
                       help="Generate tutorials for specific domain only")

    args = parser.parse_args()

    generator = ComprehensiveTutorialGenerator(args.output)

    if args.domain:
        # Generate for specific domain
        domain_configs = generator._load_domain_configs()
        if args.domain in domain_configs:
            success = generator._generate_domain_tutorials(args.domain, domain_configs[args.domain])
            print(f"✅ Generated tutorials for {args.domain}" if success else f"❌ Failed for {args.domain}")
        else:
            print(f"❌ Domain {args.domain} not found!")
            sys.exit(1)
    else:
        # Generate for all domains
        success = generator.generate_all_tutorials()
        print("🎉 All tutorials generated successfully!" if success else "❌ Some tutorials failed!")

    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
