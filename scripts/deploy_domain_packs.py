#!/usr/bin/env python3
"""
Domain Pack Deployment Script
Generates all 25 domain packs for AWS Research Wizard Phase 2
"""

import os
import sys
import yaml
import json
from pathlib import Path
from typing import Dict, List, Any
import logging

# Add the python directory to path to import comprehensive_spack_domains
sys.path.append(str(Path(__file__).parent.parent / "python"))
from comprehensive_spack_domains import ComprehensiveSpackGenerator

class DomainPackDeployer:
    def __init__(self, output_dir: str = "domain-packs"):
        self.output_dir = Path(output_dir)
        self.domains_dir = self.output_dir / "domains"
        
        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)
        
        # Category mapping for organizing domains
        self.category_mapping = {
            # Life Sciences
            "genomics_lab": "life-sciences",
            "structural_biology": "life-sciences", 
            "systems_biology": "life-sciences",
            "neuroscience_lab": "life-sciences",
            "drug_discovery": "life-sciences",
            
            # Physical Sciences
            "climate_modeling": "physical-sciences",
            "materials_science": "physical-sciences",
            "chemistry_lab": "physical-sciences", 
            "physics_simulation": "physical-sciences",
            "astronomy_lab": "physical-sciences",
            "geoscience_lab": "physical-sciences",
            
            # Engineering
            "cfd_engineering": "engineering",
            "mechanical_engineering": "engineering",
            "electrical_engineering": "engineering", 
            "aerospace_engineering": "engineering",
            
            # Computer Science & AI
            "ai_research_studio": "computer-science",
            "hpc_development": "computer-science",
            "data_science_lab": "computer-science",
            "quantum_computing": "computer-science",
            
            # Social Sciences & Humanities
            "digital_humanities": "social-sciences",
            "economics_analysis": "social-sciences",
            "social_science_lab": "social-sciences",
            
            # Interdisciplinary
            "mathematical_modeling": "interdisciplinary",
            "visualization_studio": "interdisciplinary",
            "research_workflow": "interdisciplinary",
        }

    def deploy_all_domains(self) -> bool:
        """Deploy all 25 domain packs"""
        self.logger.info("🚀 Starting Phase 2: Domain Pack System deployment...")
        
        # Create directory structure
        self._create_directory_structure()
        
        # Generate all domain packs (only implemented ones)
        generator = ComprehensiveSpackGenerator()
        
        # Only create implemented domain packs
        implemented_domains = {
            "genomics_lab": generator._create_genomics_pack(),
            "climate_modeling": generator._create_climate_pack(),
            "materials_science": generator._create_materials_pack(),
            "ai_research_studio": generator._create_ai_pack(),
            "neuroscience_lab": generator._create_neuroscience_pack(),
            "physics_simulation": generator._create_physics_pack(),
            "astronomy_lab": generator._create_astronomy_pack(),
        }
        domain_packs = implemented_domains
        
        success_count = 0
        total_count = len(domain_packs)
        
        # Deploy in priority order (as specified in Phase 2 plan)
        priority_order = [
            "genomics_lab",          # 1. Genomics
            "ai_research_studio",    # 2. Machine Learning  
            "climate_modeling",      # 3. Climate Modeling
            "astronomy_lab",         # 4. Astronomy
            "chemistry_lab",         # 5. Chemistry
        ]
        
        # Deploy priority domains first
        for domain_key in priority_order:
            if domain_key in domain_packs:
                if self._deploy_single_domain(domain_key, domain_packs[domain_key]):
                    success_count += 1
                    self.logger.info(f"✅ Priority domain {domain_key} deployed successfully")
                else:
                    self.logger.error(f"❌ Failed to deploy priority domain {domain_key}")
        
        # Deploy remaining domains
        remaining_domains = [k for k in domain_packs.keys() if k not in priority_order]
        for domain_key in remaining_domains:
            if self._deploy_single_domain(domain_key, domain_packs[domain_key]):
                success_count += 1
                self.logger.info(f"✅ Domain {domain_key} deployed successfully")
            else:
                self.logger.error(f"❌ Failed to deploy domain {domain_key}")
        
        # Summary
        self.logger.info(f"")
        self.logger.info(f"📊 Deployment Summary:")
        self.logger.info(f"   Total domains: {total_count}")
        self.logger.info(f"   ✅ Successfully deployed: {success_count}")
        self.logger.info(f"   ❌ Failed: {total_count - success_count}")
        
        if success_count == total_count:
            self.logger.info(f"🎉 Phase 2 Domain Pack System deployment complete!")
            return True
        else:
            self.logger.error(f"❌ Phase 2 deployment incomplete - {total_count - success_count} failures")
            return False

    def _create_directory_structure(self):
        """Create the domain pack directory structure"""
        categories = set(self.category_mapping.values())
        
        for category in categories:
            category_dir = self.domains_dir / category
            category_dir.mkdir(parents=True, exist_ok=True)
            self.logger.info(f"📁 Created category directory: {category}")

    def _deploy_single_domain(self, domain_key: str, domain_pack) -> bool:
        """Deploy a single domain pack"""
        try:
            category = self.category_mapping.get(domain_key, "miscellaneous")
            domain_dir = self.domains_dir / category / domain_key
            domain_dir.mkdir(parents=True, exist_ok=True)
            
            # Create domain-pack.yaml
            domain_config = self._convert_to_domain_config(domain_pack)
            domain_config_path = domain_dir / "domain-pack.yaml"
            with open(domain_config_path, 'w') as f:
                yaml.dump(domain_config, f, default_flow_style=False, indent=2)
            
            # Create spack.yaml environment file
            spack_config = self._create_spack_environment(domain_pack)
            spack_config_path = domain_dir / "spack.yaml"
            with open(spack_config_path, 'w') as f:
                yaml.dump(spack_config, f, default_flow_style=False, indent=2)
            
            # Create additional directories
            for subdir in ["workflows", "docs", "examples"]:
                (domain_dir / subdir).mkdir(exist_ok=True)
            
            # Create basic README
            readme_path = domain_dir / "README.md"
            readme_content = self._generate_readme(domain_pack)
            with open(readme_path, 'w') as f:
                f.write(readme_content)
            
            return True
            
        except Exception as e:
            self.logger.error(f"Error deploying {domain_key}: {e}")
            return False

    def _convert_to_domain_config(self, domain_pack) -> Dict[str, Any]:
        """Convert SpackDomainPack to domain-pack.yaml format"""
        return {
            "name": domain_pack.name,
            "description": domain_pack.description,
            "version": "1.0.0",
            "categories": domain_pack.primary_domains,
            "maintainers": [
                {
                    "name": "AWS Research Wizard",
                    "email": "aws-research-wizard@example.com",
                    "organization": "AWS Research Computing"
                }
            ],
            "spack_config": {
                "packages": self._flatten_spack_packages(domain_pack.spack_packages),
                "compiler": "gcc@11.4.0",
                "target": "neoverse_v1",
                "optimization": "-O3"
            },
            "aws_config": {
                "instance_types": self._get_instance_recommendations(domain_pack),
                "storage": {
                    "type": "gp3",
                    "size_gb": 500,
                    "iops": 3000,
                    "throughput": 125
                },
                "network": {
                    "placement_group": True,
                    "enhanced_networking": True,
                    "efa_enabled": True
                }
            },
            "workflows": [
                {
                    "name": workflow,
                    "description": f"Sample {workflow} workflow",
                    "script": f"workflows/{workflow.lower().replace(' ', '_')}.sh",
                    "input_data": "s3://aws-open-data/",
                    "expected_output": "Processed research data"
                }
                for workflow in domain_pack.sample_workflows[:3]  # Limit to 3 workflows
            ],
            "cost_estimates": domain_pack.cost_profile,
            "documentation": {
                "getting_started": f"docs/{domain_pack.name.lower().replace(' ', '_')}_quickstart.md",
                "tutorials": f"docs/{domain_pack.name.lower().replace(' ', '_')}_tutorials.md",
                "best_practices": f"docs/{domain_pack.name.lower().replace(' ', '_')}_best_practices.md"
            }
        }

    def _create_spack_environment(self, domain_pack) -> Dict[str, Any]:
        """Create Spack environment configuration"""
        flattened_packages = self._flatten_spack_packages(domain_pack.spack_packages)
        
        return {
            "spack": {
                "specs": flattened_packages,
                "view": True,
                "concretizer": {
                    "unify": True,
                    "reuse": True
                },
                "config": {
                    "install_tree": {
                        "root": "$spack/opt/spack"
                    },
                    "build_stage": ["$tempdir/$user/spack-stage"],
                    "build_cache": True
                },
                "mirrors": {
                    "aws-binary-cache": domain_pack.aws_spack_cache.get("primary_cache", "https://cache.spack.io/aws-ahug-east/")
                },
                "packages": {
                    "all": {
                        "compiler": ["gcc@11.4.0"],
                        "target": ["neoverse_v1"]
                    }
                }
            }
        }

    def _flatten_spack_packages(self, spack_packages: Dict[str, List[str]]) -> List[str]:
        """Flatten the hierarchical spack packages into a simple list"""
        flattened = []
        for category, packages in spack_packages.items():
            flattened.extend(packages)
        return flattened

    def _get_instance_recommendations(self, domain_pack) -> Dict[str, str]:
        """Generate AWS instance type recommendations based on domain"""
        # Default recommendations - could be enhanced with domain-specific logic
        return {
            "small": "m6i.large",
            "medium": "m6i.xlarge", 
            "large": "m6i.2xlarge",
            "gpu": "g5.xlarge",
            "hpc": "c6i.4xlarge",
            "memory": "r6i.2xlarge"
        }

    def _generate_readme(self, domain_pack) -> str:
        """Generate README content for domain pack"""
        return f"""# {domain_pack.name}

{domain_pack.description}

## Target Users
{domain_pack.target_users}

## Primary Domains
{', '.join(domain_pack.primary_domains)}

## Deployment Variants
{', '.join(domain_pack.deployment_variants)}

## Sample Workflows
{chr(10).join(f"- {workflow}" for workflow in domain_pack.sample_workflows)}

## Immediate Value
{chr(10).join(f"- {value}" for value in domain_pack.immediate_value)}

## Cost Profile
{chr(10).join(f"- {k}: {v}" for k, v in domain_pack.cost_profile.items())}

## Getting Started

1. Deploy the environment:
   ```bash
   ./aws-research-wizard deploy --domain {domain_pack.name.lower().replace(' ', '_')}
   ```

2. Activate the Spack environment:
   ```bash
   spack env activate {domain_pack.name.lower().replace(' ', '_')}
   ```

3. Install packages:
   ```bash
   spack install
   ```

## AWS Optimization

This domain pack is optimized for AWS infrastructure with:
- Binary cache integration for fast deployment
- Graviton3 processor optimization
- Cost-effective instance recommendations
- Research data integration with AWS Open Data

Generated by AWS Research Wizard Phase 2 Deployment
"""

def main():
    """Main deployment function"""
    import argparse
    
    parser = argparse.ArgumentParser(description="Deploy AWS Research Wizard Domain Packs")
    parser.add_argument("--output", type=str, default="domain-packs", 
                       help="Output directory for domain packs")
    parser.add_argument("--validate", action="store_true",
                       help="Validate after deployment")
    
    args = parser.parse_args()
    
    # Deploy all domain packs
    deployer = DomainPackDeployer(args.output)
    success = deployer.deploy_all_domains()
    
    # Optional validation
    if args.validate and success:
        sys.path.append(str(Path(args.output) / "tools"))
        from validate_domains import DomainPackValidator
        
        validator = DomainPackValidator(args.output)
        validation_success = validator.validate_all()
        
        if not validation_success:
            print("❌ Validation failed after deployment!")
            sys.exit(1)
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()