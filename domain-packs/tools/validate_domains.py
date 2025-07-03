#!/usr/bin/env python3
"""
Domain Pack Validation Tool
Validates all domain pack configurations against schemas
"""

import os
import sys
import json
import yaml
import argparse
import logging
from pathlib import Path
from typing import Dict, List, Any, Optional
import jsonschema
from jsonschema import validate, ValidationError

class DomainPackValidator:
    def __init__(self, domain_packs_dir: str = "domain-packs"):
        self.domain_packs_dir = Path(domain_packs_dir)
        self.schemas_dir = self.domain_packs_dir / "schemas"
        self.domains_dir = self.domain_packs_dir / "domains"

        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)

    def load_schema(self, schema_name: str) -> Dict[str, Any]:
        """Load validation schema"""
        schema_path = self.schemas_dir / f"{schema_name}.schema.json"
        if not schema_path.exists():
            raise FileNotFoundError(f"Schema not found: {schema_path}")

        with open(schema_path, 'r') as f:
            return json.load(f)

    def find_domain_packs(self) -> List[Path]:
        """Find all domain pack configuration files"""
        domain_packs = []

        for category_dir in self.domains_dir.iterdir():
            if category_dir.is_dir():
                for domain_dir in category_dir.iterdir():
                    if domain_dir.is_dir():
                        config_file = domain_dir / "domain-pack.yaml"
                        if config_file.exists():
                            domain_packs.append(config_file)

        return domain_packs

    def validate_yaml_syntax(self, yaml_file: Path) -> bool:
        """Validate YAML syntax"""
        try:
            with open(yaml_file, 'r') as f:
                yaml.safe_load(f)
            return True
        except yaml.YAMLError as e:
            self.logger.error(f"YAML syntax error in {yaml_file}: {e}")
            return False

    def validate_domain_pack_config(self, config_file: Path) -> bool:
        """Validate domain pack configuration against schema"""
        try:
            # Load configuration
            with open(config_file, 'r') as f:
                config = yaml.safe_load(f)

            # Load schema
            schema = self.load_schema("domain-pack")

            # Validate against schema
            validate(instance=config, schema=schema)

            # Additional business logic validation
            if not self._validate_business_rules(config, config_file):
                return False

            self.logger.info(f"✅ {config_file.name} passed validation")
            return True

        except ValidationError as e:
            self.logger.error(f"❌ Schema validation failed for {config_file}: {e.message}")
            return False
        except Exception as e:
            self.logger.error(f"❌ Validation error for {config_file}: {e}")
            return False

    def _validate_business_rules(self, config: Dict[str, Any], config_file: Path) -> bool:
        """Validate business logic rules"""
        errors = []

        # Check if spack.yaml exists
        spack_file = config_file.parent / "spack.yaml"
        if not spack_file.exists():
            errors.append(f"Missing spack.yaml file in {config_file.parent}")

        # Validate instance types are real AWS types
        aws_config = config.get("aws_config", {})
        instance_types = aws_config.get("instance_types", {})

        valid_instance_families = [
            "t3", "t4g", "m5", "m6i", "m6a", "c5", "c6i", "c6a",
            "r5", "r6i", "r6a", "x1e", "z1d", "p3", "p4d", "g4dn", "g5"
        ]

        for size, instance_type in instance_types.items():
            family = instance_type.split('.')[0]
            if family not in valid_instance_families:
                errors.append(f"Unknown instance family '{family}' in {instance_type}")

        # Validate storage configuration
        storage = aws_config.get("storage", {})
        storage_type = storage.get("type")
        if storage_type and storage_type not in ["gp3", "gp2", "io1", "io2", "st1", "sc1"]:
            errors.append(f"Invalid storage type: {storage_type}")

        # Check required directories exist
        domain_dir = config_file.parent
        required_dirs = ["workflows", "docs", "examples"]
        for req_dir in required_dirs:
            dir_path = domain_dir / req_dir
            if not dir_path.exists():
                self.logger.warning(f"⚠️  Missing recommended directory: {dir_path}")

        if errors:
            for error in errors:
                self.logger.error(f"❌ Business rule violation: {error}")
            return False

        return True

    def validate_spack_environment(self, spack_file: Path) -> bool:
        """Validate Spack environment file"""
        try:
            with open(spack_file, 'r') as f:
                spack_config = yaml.safe_load(f)

            # Check required sections
            if 'spack' not in spack_config:
                self.logger.error(f"❌ Missing 'spack' section in {spack_file}")
                return False

            spack_section = spack_config['spack']

            if 'specs' not in spack_section:
                self.logger.error(f"❌ Missing 'specs' section in {spack_file}")
                return False

            # Validate specs format
            specs = spack_section['specs']
            if not isinstance(specs, list):
                self.logger.error(f"❌ 'specs' must be a list in {spack_file}")
                return False

            # Check for basic required packages
            spec_strings = [str(spec) for spec in specs]
            has_compiler = any('gcc@' in spec or 'clang@' in spec for spec in spec_strings)

            if not has_compiler:
                self.logger.warning(f"⚠️  No explicit compiler found in {spack_file}")

            self.logger.info(f"✅ {spack_file.name} passed validation")
            return True

        except Exception as e:
            self.logger.error(f"❌ Spack validation error for {spack_file}: {e}")
            return False

    def validate_all(self) -> bool:
        """Validate all domain packs"""
        self.logger.info("🔍 Starting domain pack validation...")

        domain_packs = self.find_domain_packs()
        if not domain_packs:
            self.logger.warning("⚠️  No domain packs found!")
            return True

        total_packs = len(domain_packs)
        passed = 0
        failed = 0

        for config_file in domain_packs:
            self.logger.info(f"Validating {config_file}...")

            # Validate YAML syntax
            if not self.validate_yaml_syntax(config_file):
                failed += 1
                continue

            # Validate domain pack config
            if not self.validate_domain_pack_config(config_file):
                failed += 1
                continue

            # Validate associated spack.yaml
            spack_file = config_file.parent / "spack.yaml"
            if spack_file.exists():
                if not self.validate_spack_environment(spack_file):
                    failed += 1
                    continue

            passed += 1

        # Summary
        self.logger.info(f"")
        self.logger.info(f"📊 Validation Summary:")
        self.logger.info(f"   Total domain packs: {total_packs}")
        self.logger.info(f"   ✅ Passed: {passed}")
        self.logger.info(f"   ❌ Failed: {failed}")

        if failed > 0:
            self.logger.error(f"❌ {failed} domain pack(s) failed validation!")
            return False
        else:
            self.logger.info(f"✅ All domain packs passed validation!")
            return True

def main():
    parser = argparse.ArgumentParser(description="Validate AWS Research Wizard domain packs")
    parser.add_argument("--all", action="store_true", help="Validate all domain packs")
    parser.add_argument("--domain", type=str, help="Validate specific domain pack")
    parser.add_argument("--dir", type=str, default="domain-packs",
                       help="Domain packs directory (default: domain-packs)")

    args = parser.parse_args()

    validator = DomainPackValidator(args.dir)

    if args.all:
        success = validator.validate_all()
        sys.exit(0 if success else 1)
    elif args.domain:
        # Find specific domain pack
        domain_packs = validator.find_domain_packs()
        found = False
        for config_file in domain_packs:
            if args.domain in str(config_file):
                found = True
                success = validator.validate_domain_pack_config(config_file)
                sys.exit(0 if success else 1)

        if not found:
            print(f"Domain pack '{args.domain}' not found!")
            sys.exit(1)
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == "__main__":
    main()
