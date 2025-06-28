#!/usr/bin/env python3
"""
Configuration Loader and Validator for AWS Research Wizard

This module provides comprehensive configuration management for research pack configurations,
including loading from YAML files, schema validation, and template application.

Key Features:
- YAML-based configuration loading with schema validation
- MPI template application for HPC optimization
- Comprehensive domain pack configuration management
- JSON export capabilities for configuration data
- Batch validation for all domain configurations

Classes:
    DomainPackConfig: Data class representing a complete research pack configuration
    MPIConfig: Data class for MPI-specific configuration templates
    ConfigLoader: Main configuration management class

Dependencies:
    - jsonschema: Required for configuration validation
    - yaml: For YAML file parsing
    - pathlib: For cross-platform path handling
"""

import os
import yaml
import json
from typing import Dict, List, Any, Optional, Union
from pathlib import Path
from dataclasses import dataclass, asdict
import logging

try:
    import jsonschema
    from jsonschema import validate, ValidationError
except ImportError:
    print("Please install jsonschema: pip install jsonschema")
    exit(1)


@dataclass
class DomainPackConfig:
    """
    Comprehensive data class representing a research pack configuration.

    This class encapsulates all configuration aspects for a specific research domain,
    including software packages, AWS instance recommendations, cost estimates,
    and domain-specific features like MPI optimizations and AWS data integration.

    Attributes:
        name (str): Human-readable name of the research pack
        description (str): Detailed description of the research pack's purpose
        primary_domains (List[str]): List of research domains this pack serves
        target_users (str): Description of intended user base
        spack_packages (Dict[str, List[str]]): Categorized software packages via Spack
        aws_instance_recommendations (Dict[str, Dict[str, Any]]): AWS instance configurations
        estimated_cost (Dict[str, float]): Cost estimates for different usage scenarios
        research_capabilities (List[str]): List of research capabilities enabled
        aws_data_sources (Optional[List[str]]): Available AWS Open Data sources
        demo_workflows (Optional[List[Dict[str, Any]]]): Pre-configured demo workflows
        mpi_optimizations (Optional[Dict[str, Any]]): MPI and EFA optimization settings
        scaling_profiles (Optional[Dict[str, Any]]): Scaling configurations for different workloads
        mpi_environment (Optional[Dict[str, Any]]): MPI environment variables and settings
        mpi_runtime_flags (Optional[Dict[str, Any]]): MPI runtime optimization flags
        security_features (Optional[Dict[str, Any]]): Security configurations and features
        chemistry_features (Optional[Dict[str, Any]]): Chemistry-specific configurations
        agricultural_features (Optional[Dict[str, Any]]): Agriculture-specific configurations
        geospatial_features (Optional[Dict[str, Any]]): Geospatial analysis configurations
        aws_integration (Optional[Dict[str, Any]]): AWS data integration metadata
    """

    name: str
    description: str
    primary_domains: List[str]
    target_users: str
    spack_packages: Dict[str, List[str]]
    aws_instance_recommendations: Dict[str, Dict[str, Any]]
    estimated_cost: Dict[str, float]
    research_capabilities: List[str]
    aws_data_sources: Optional[List[str]] = None
    demo_workflows: Optional[List[Dict[str, Any]]] = None
    mpi_optimizations: Optional[Dict[str, Any]] = None
    scaling_profiles: Optional[Dict[str, Any]] = None
    mpi_environment: Optional[Dict[str, Any]] = None
    mpi_runtime_flags: Optional[Dict[str, Any]] = None
    security_features: Optional[Dict[str, Any]] = None
    chemistry_features: Optional[Dict[str, Any]] = None
    agricultural_features: Optional[Dict[str, Any]] = None
    geospatial_features: Optional[Dict[str, Any]] = None
    marine_features: Optional[Dict[str, Any]] = None
    sports_features: Optional[Dict[str, Any]] = None
    biomechanics_features: Optional[Dict[str, Any]] = None
    aws_integration: Optional[Dict[str, Any]] = None


@dataclass
class MPIConfig:
    """
    MPI configuration template data class for high-performance computing optimization.

    This class encapsulates MPI (Message Passing Interface) and EFA (Elastic Fabric Adapter)
    configurations that are applied to research packs requiring HPC capabilities.

    Attributes:
        mpi_packages (Dict[str, List[str]]): MPI software packages organized by category
        efa_environment (Dict[str, Dict[str, str]]): EFA-specific environment variables
        mpi_runtime_flags (Dict[str, List[str]]): Runtime flags for MPI optimization
        efa_instance_types (Dict[str, List[Dict[str, Any]]]): EFA-enabled instance configurations
        placement_groups (Dict[str, Dict[str, Any]]): AWS placement group configurations
        scaling_profiles (Dict[str, Dict[str, Any]]): Scaling profiles for different workload sizes
    """

    mpi_packages: Dict[str, List[str]]
    efa_environment: Dict[str, Dict[str, str]]
    mpi_runtime_flags: Dict[str, List[str]]
    efa_instance_types: Dict[str, List[Dict[str, Any]]]
    placement_groups: Dict[str, Dict[str, Any]]
    scaling_profiles: Dict[str, Dict[str, Any]]


class ConfigLoader:
    """
    Primary configuration management class for AWS Research Wizard.

    This class handles loading, validation, and management of research pack configurations
    from YAML files. It provides schema validation, template application, and comprehensive
    configuration management capabilities.

    The ConfigLoader supports:
    - Loading individual domain configurations from YAML files
    - Batch loading and validation of all configurations
    - JSON Schema validation for configuration integrity
    - MPI template application for HPC-optimized configurations
    - Export capabilities for configuration data

    Directory Structure Expected:
        config_root/
        ├── domains/          # Domain-specific configuration files
        │   ├── genomics.yaml
        │   ├── machine_learning.yaml
        │   └── ...
        ├── schemas/          # JSON Schema validation files
        │   ├── domain_pack_schema.yaml
        │   └── ...
        └── templates/        # Configuration templates
            ├── aws_mpi_base.yaml
            └── ...

    Attributes:
        config_root (Path): Root directory for configuration files
        logger (logging.Logger): Logger instance for this class
        schemas (Dict[str, Dict]): Loaded JSON schemas for validation
        templates (Dict[str, Dict]): Loaded configuration templates

    Example:
        >>> loader = ConfigLoader("configs")
        >>> config = loader.load_domain_config("genomics")
        >>> if config:
        ...     print(f"Loaded {config.name} with {len(config.spack_packages)} package categories")
    """

    def __init__(self, config_root: str = "configs"):
        """
        Initialize the ConfigLoader with the specified configuration root directory.

        Args:
            config_root (str): Path to the root configuration directory. Defaults to "configs".
                              This directory should contain domains/, schemas/, and templates/ subdirectories.

        Raises:
            FileNotFoundError: If the config_root directory doesn't exist
            PermissionError: If the config_root directory is not readable
        """
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)

        # Load schemas for configuration validation
        self.schemas = self._load_schemas()

        # Load templates for configuration enhancement
        self.templates = self._load_templates()

    def _load_schemas(self) -> Dict[str, Dict]:
        """
        Load JSON schema files for configuration validation.

        Scans the schemas/ directory for YAML files containing JSON schemas
        and loads them into memory for configuration validation.

        Returns:
            Dict[str, Dict]: Dictionary mapping schema names to schema definitions

        Note:
            Schema files should be named descriptively (e.g., domain_pack_schema.yaml)
            as the filename (without extension) becomes the schema identifier.
        """
        schemas = {}
        schema_dir = self.config_root / "schemas"

        if schema_dir.exists():
            for schema_file in schema_dir.glob("*.yaml"):
                try:
                    with open(schema_file, "r") as f:
                        schema_name = schema_file.stem
                        schemas[schema_name] = yaml.safe_load(f)
                        self.logger.info(f"Loaded schema: {schema_name}")
                except Exception as e:
                    self.logger.error(f"Failed to load schema {schema_file}: {e}")

        return schemas

    def _load_templates(self) -> Dict[str, Dict]:
        """
        Load configuration templates for domain pack enhancement.

        Templates provide reusable configuration components that can be applied
        to multiple domain configurations. The primary template is aws_mpi_base.yaml
        which provides MPI and EFA optimizations for HPC workloads.

        Returns:
            Dict[str, Dict]: Dictionary mapping template names to template definitions

        Templates include:
            - aws_mpi_base: MPI and EFA configurations for high-performance computing
            - Additional templates can be added for other common configuration patterns
        """
        templates = {}
        template_dir = self.config_root / "templates"

        if template_dir.exists():
            for template_file in template_dir.glob("*.yaml"):
                try:
                    with open(template_file, "r") as f:
                        template_name = template_file.stem
                        templates[template_name] = yaml.safe_load(f)
                        self.logger.info(f"Loaded template: {template_name}")
                except Exception as e:
                    self.logger.error(f"Failed to load template {template_file}: {e}")

        return templates

    def validate_config(
        self, config: Dict[str, Any], schema_name: str = "domain_pack_schema"
    ) -> bool:
        """Validate configuration against schema"""
        if schema_name not in self.schemas:
            self.logger.warning(f"Schema {schema_name} not found, skipping validation")
            return True

        try:
            validate(instance=config, schema=self.schemas[schema_name])
            self.logger.info(f"Configuration validation passed for schema: {schema_name}")
            return True
        except ValidationError as e:
            self.logger.error(f"Configuration validation failed: {e.message}")
            self.logger.error(f"Failed at path: {' -> '.join(str(p) for p in e.absolute_path)}")
            return False

    def load_domain_config(self, domain_name: str) -> Optional[DomainPackConfig]:
        """Load and validate domain configuration"""
        config_file = self.config_root / "domains" / f"{domain_name}.yaml"

        if not config_file.exists():
            self.logger.error(f"Domain config file not found: {config_file}")
            return None

        try:
            with open(config_file, "r") as f:
                config_data = yaml.safe_load(f)

            # Validate against schema
            if not self.validate_config(config_data, "domain_pack_schema"):
                return None

            # Apply MPI base template if MPI optimizations are enabled
            if config_data.get("mpi_optimizations", {}).get("efa_enabled", False):
                config_data = self._apply_mpi_template(config_data)

            # Convert to dataclass
            domain_config = DomainPackConfig(**config_data)
            self.logger.info(f"Successfully loaded domain config: {domain_name}")

            return domain_config

        except Exception as e:
            self.logger.error(f"Failed to load domain config {domain_name}: {e}")
            return None

    def _apply_mpi_template(self, config: Dict[str, Any]) -> Dict[str, Any]:
        """Apply MPI template optimizations to domain config"""
        if "aws_mpi_base" not in self.templates:
            self.logger.warning("MPI base template not found, skipping MPI optimizations")
            return config

        mpi_template = self.templates["aws_mpi_base"]

        # Add MPI packages to spack_packages
        if "mpi_packages" not in config["spack_packages"]:
            config["spack_packages"]["mpi_packages"] = []

        # Merge MPI packages from template
        for category, packages in mpi_template["mpi_packages"].items():
            if category not in config["spack_packages"]:
                config["spack_packages"][category] = packages
            else:
                # Merge without duplicates
                existing = set(config["spack_packages"][category])
                for pkg in packages:
                    if pkg not in existing:
                        config["spack_packages"][category].append(pkg)

        # Apply EFA optimizations to instance recommendations
        for instance_name, instance_config in config["aws_instance_recommendations"].items():
            if instance_config.get("efa_enabled", False):
                # Find matching EFA instance type
                instance_type = instance_config["instance_type"]
                for category, instances in mpi_template["efa_instance_types"].items():
                    for efa_instance in instances:
                        if efa_instance["instance_type"] == instance_type:
                            # Apply EFA-specific settings
                            instance_config.update(
                                {
                                    "network_performance": efa_instance["network_performance"],
                                    "placement_group": "cluster",
                                    "enhanced_networking": "sr-iov",
                                }
                            )
                            break

        # Add MPI environment and runtime settings
        config["mpi_environment"] = mpi_template["efa_environment"]
        config["mpi_runtime_flags"] = mpi_template["mpi_runtime_flags"]

        self.logger.info("Applied MPI template optimizations")
        return config

    def load_all_domain_configs(self) -> Dict[str, DomainPackConfig]:
        """Load all available domain configurations"""
        configs = {}
        domain_dir = self.config_root / "domains"

        if not domain_dir.exists():
            self.logger.error(f"Domain config directory not found: {domain_dir}")
            return configs

        for config_file in domain_dir.glob("*.yaml"):
            domain_name = config_file.stem
            config = self.load_domain_config(domain_name)
            if config:
                configs[domain_name] = config

        self.logger.info(f"Loaded {len(configs)} domain configurations")
        return configs

    def get_mpi_config(self) -> Optional[MPIConfig]:
        """Get MPI configuration from template"""
        if "aws_mpi_base" not in self.templates:
            return None

        template = self.templates["aws_mpi_base"]
        return MPIConfig(
            mpi_packages=template["mpi_packages"],
            efa_environment=template["efa_environment"],
            mpi_runtime_flags=template["mpi_runtime_flags"],
            efa_instance_types=template["efa_instance_types"],
            placement_groups=template["placement_groups"],
            scaling_profiles=template["scaling_profiles"],
        )

    def list_available_domains(self) -> List[str]:
        """List all available domain configurations"""
        domain_dir = self.config_root / "domains"
        if not domain_dir.exists():
            return []

        return [f.stem for f in domain_dir.glob("*.yaml")]

    def export_config_to_json(self, domain_name: str, output_file: str) -> bool:
        """Export domain configuration to JSON"""
        config = self.load_domain_config(domain_name)
        if not config:
            return False

        try:
            config_dict = asdict(config)
            with open(output_file, "w") as f:
                json.dump(config_dict, f, indent=2, default=str)

            self.logger.info(f"Exported config to {output_file}")
            return True

        except Exception as e:
            self.logger.error(f"Failed to export config: {e}")
            return False

    def validate_all_configs(self) -> Dict[str, bool]:
        """Validate all domain configurations"""
        results = {}
        for domain_name in self.list_available_domains():
            config = self.load_domain_config(domain_name)
            results[domain_name] = config is not None

        return results


def main():
    """CLI interface for configuration management"""
    import argparse

    parser = argparse.ArgumentParser(description="Research Pack Configuration Manager")
    parser.add_argument("--list", action="store_true", help="List available domain configs")
    parser.add_argument("--validate", action="store_true", help="Validate all configurations")
    parser.add_argument("--load", type=str, help="Load specific domain configuration")
    parser.add_argument("--export", type=str, help="Export domain config to JSON")
    parser.add_argument("--output", type=str, help="Output file for export")
    parser.add_argument(
        "--config-root", type=str, default="configs", help="Configuration root directory"
    )

    args = parser.parse_args()

    # Setup logging
    logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

    # Initialize loader
    loader = ConfigLoader(args.config_root)

    if args.list:
        domains = loader.list_available_domains()
        print(f"Available domain configurations ({len(domains)}):")
        for domain in sorted(domains):
            print(f"  - {domain}")

    elif args.validate:
        results = loader.validate_all_configs()
        print("Configuration validation results:")
        for domain, valid in results.items():
            status = "✅ PASS" if valid else "❌ FAIL"
            print(f"  {domain}: {status}")

    elif args.load:
        config = loader.load_domain_config(args.load)
        if config:
            print(f"Successfully loaded configuration for: {config.name}")
            print(f"Primary domains: {', '.join(config.primary_domains)}")
            print(f"Spack package categories: {len(config.spack_packages)}")
            print(f"Instance recommendations: {len(config.aws_instance_recommendations)}")
        else:
            print(f"Failed to load configuration: {args.load}")

    elif args.export:
        if not args.output:
            args.output = f"{args.export}_config.json"

        success = loader.export_config_to_json(args.export, args.output)
        if success:
            print(f"Configuration exported to: {args.output}")
        else:
            print(f"Failed to export configuration: {args.export}")

    else:
        parser.print_help()


if __name__ == "__main__":
    main()
