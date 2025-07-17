#!/usr/bin/env python3
"""
Configuration Version Management Tool
Handles versioning, validation, and migration of configuration files
"""

import os
import sys
import yaml
import argparse
# import semver  # Commented out to avoid dependency
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from datetime import datetime

class ConfigVersionManager:
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
        self.configs_dir = root_dir / "configs"
        self.domains_dir = self.configs_dir / "domains"
        self.schemas_dir = self.configs_dir / "schemas"
        self.schema_file = self.schemas_dir / "domain_pack_schema.yaml"
        self.compatibility_file = self.configs_dir / "VERSION_COMPATIBILITY.yaml"

    def get_schema_version(self) -> str:
        """Get current schema version"""
        try:
            with open(self.schema_file, 'r') as f:
                schema = yaml.safe_load(f)
                return schema.get('schema_version', '2.0.0')
        except Exception as e:
            print(f"Error reading schema: {e}")
            return '2.0.0'

    def get_config_version(self, config_path: Path) -> Optional[str]:
        """Get version from a configuration file"""
        try:
            with open(config_path, 'r') as f:
                config = yaml.safe_load(f)
                return config.get('config_version')
        except Exception as e:
            print(f"Error reading config {config_path}: {e}")
            return None

    def add_version_to_config(self, config_path: Path, version: str = None) -> bool:
        """Add version information to a configuration file"""
        try:
            with open(config_path, 'r') as f:
                config = yaml.safe_load(f)

            # Use provided version or current schema version
            if not version:
                version = self.get_schema_version()

            # Add version information
            config['config_version'] = version
            config['config_metadata'] = {
                'created': config.get('config_metadata', {}).get('created', '2024-12-15'),
                'last_modified': datetime.now().strftime('%Y-%m-%d'),
                'schema_compatibility': self.get_schema_version(),
                'change_log': config.get('config_metadata', {}).get('change_log', [
                    {
                        'version': version,
                        'date': datetime.now().strftime('%Y-%m-%d'),
                        'changes': ['Added configuration versioning']
                    }
                ])
            }

            # Write back to file
            with open(config_path, 'w') as f:
                yaml.dump(config, f, default_flow_style=False, sort_keys=False, allow_unicode=True)

            print(f"✓ Added version {version} to {config_path.name}")
            return True

        except Exception as e:
            print(f"✗ Error adding version to {config_path}: {e}")
            return False

    def bump_config_version(self, config_path: Path, bump_type: str) -> bool:
        """Bump version in a configuration file"""
        try:
            current_version = self.get_config_version(config_path)
            if not current_version:
                print(f"Config {config_path.name} has no version, adding initial version")
                return self.add_version_to_config(config_path)

            # Parse and bump version
            if bump_type == 'major':
                new_version = semver.bump_major(current_version)
            elif bump_type == 'minor':
                new_version = semver.bump_minor(current_version)
            elif bump_type == 'patch':
                new_version = semver.bump_patch(current_version)
            else:
                print(f"Invalid bump type: {bump_type}")
                return False

            # Load config and update version
            with open(config_path, 'r') as f:
                config = yaml.safe_load(f)

            config['config_version'] = new_version
            if 'config_metadata' in config:
                config['config_metadata']['last_modified'] = datetime.now().strftime('%Y-%m-%d')

                # Add to change log
                if 'change_log' not in config['config_metadata']:
                    config['config_metadata']['change_log'] = []

                config['config_metadata']['change_log'].append({
                    'version': new_version,
                    'date': datetime.now().strftime('%Y-%m-%d'),
                    'changes': [f"Version bump: {bump_type}"]
                })

            # Write back
            with open(config_path, 'w') as f:
                yaml.dump(config, f, default_flow_style=False, sort_keys=False, allow_unicode=True)

            print(f"✓ Bumped {config_path.name} from {current_version} to {new_version}")
            return True

        except Exception as e:
            print(f"✗ Error bumping version for {config_path}: {e}")
            return False

    def validate_version_compatibility(self) -> Tuple[List[str], List[str]]:
        """Validate version compatibility across all configs"""
        errors = []
        warnings = []
        schema_version = self.get_schema_version()

        # Check all domain configurations
        for config_file in self.domains_dir.glob("*.yaml"):
            config_version = self.get_config_version(config_file)

            if not config_version:
                warnings.append(f"{config_file.name}: Missing config_version field")
                continue

            # Check semantic version format
            try:
                semver.parse(config_version)
            except Exception:
                errors.append(f"{config_file.name}: Invalid version format '{config_version}'")
                continue

            # Check compatibility with schema
            try:
                if semver.compare(config_version, "2.0.0") < 0:
                    errors.append(f"{config_file.name}: Version {config_version} is deprecated (minimum 2.0.0)")
                elif semver.compare(config_version, schema_version) > 0:
                    warnings.append(f"{config_file.name}: Version {config_version} is newer than schema {schema_version}")
            except Exception as e:
                errors.append(f"{config_file.name}: Version comparison error: {e}")

        return errors, warnings

    def generate_version_report(self) -> Dict:
        """Generate comprehensive version status report"""
        schema_version = self.get_schema_version()

        # Collect version information
        domain_configs = list(self.domains_dir.glob("*.yaml"))
        version_distribution = {}
        unversioned = []

        for config_file in domain_configs:
            version = self.get_config_version(config_file)
            if version:
                version_distribution[version] = version_distribution.get(version, 0) + 1
            else:
                unversioned.append(config_file.name)

        errors, warnings = self.validate_version_compatibility()

        return {
            'schema_version': schema_version,
            'total_configs': len(domain_configs),
            'version_distribution': version_distribution,
            'unversioned_configs': unversioned,
            'compatibility_errors': len(errors),
            'compatibility_warnings': len(warnings),
            'errors': errors,
            'warnings': warnings
        }

    def add_versions_to_all_configs(self, version: str = None):
        """Add version information to all domain configurations"""
        domain_configs = list(self.domains_dir.glob("*.yaml"))
        success_count = 0

        for config_file in domain_configs:
            if self.get_config_version(config_file) is None:
                if self.add_version_to_config(config_file, version):
                    success_count += 1

        print(f"Added versions to {success_count}/{len(domain_configs)} configurations")

def main():
    parser = argparse.ArgumentParser(description="Configuration Version Management Tool")
    parser.add_argument('--root', type=Path, default=Path.cwd(), help="Root directory of project")

    subparsers = parser.add_subparsers(dest='command', help='Available commands')

    # Add versions command
    add_parser = subparsers.add_parser('add-versions', help='Add version info to all configs')
    add_parser.add_argument('--version', help='Version to assign (default: current schema version)')

    # Bump version command
    bump_parser = subparsers.add_parser('bump', help='Bump version for specific config')
    bump_parser.add_argument('config', help='Configuration file name')
    bump_parser.add_argument('type', choices=['major', 'minor', 'patch'], help='Version bump type')

    # Validate command
    validate_parser = subparsers.add_parser('validate', help='Validate version compatibility')

    # Report command
    report_parser = subparsers.add_parser('report', help='Generate version status report')

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        return 1

    manager = ConfigVersionManager(args.root)

    if args.command == 'add-versions':
        manager.add_versions_to_all_configs(args.version)

    elif args.command == 'bump':
        config_path = manager.domains_dir / args.config
        if not config_path.exists():
            print(f"Configuration file not found: {config_path}")
            return 1
        manager.bump_config_version(config_path, args.type)

    elif args.command == 'validate':
        errors, warnings = manager.validate_version_compatibility()

        if warnings:
            print("WARNINGS:")
            for warning in warnings:
                print(f"  ⚠️  {warning}")

        if errors:
            print("ERRORS:")
            for error in errors:
                print(f"  ❌ {error}")
            return 1
        else:
            print("✅ All configurations have valid versions")

    elif args.command == 'report':
        report = manager.generate_version_report()

        print(f"Configuration Version Report")
        print(f"=" * 40)
        print(f"Schema Version: {report['schema_version']}")
        print(f"Total Configurations: {report['total_configs']}")
        print(f"Version Distribution: {report['version_distribution']}")
        print(f"Unversioned: {len(report['unversioned_configs'])}")
        print(f"Compatibility Errors: {report['compatibility_errors']}")
        print(f"Compatibility Warnings: {report['compatibility_warnings']}")

        if report['unversioned_configs']:
            print(f"\nUnversioned Configurations:")
            for config in report['unversioned_configs']:
                print(f"  - {config}")

    return 0

if __name__ == "__main__":
    sys.exit(main())
