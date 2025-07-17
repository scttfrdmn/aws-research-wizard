#!/usr/bin/env python3
"""
Validate domain configuration YAML files against schema
"""
import os
import sys
import yaml
import jsonschema
from pathlib import Path

def validate_domain_config(config_path, schema_path):
    """Validate a single domain config against schema"""
    try:
        # Load schema
        with open(schema_path, 'r') as f:
            schema = yaml.safe_load(f)

        # Load config
        with open(config_path, 'r') as f:
            config = yaml.safe_load(f)

        # Check version compatibility first
        config_version = config.get('config_version')
        schema_version = schema.get('schema_version')

        if not config_version:
            return False, "Missing config_version field"

        if schema_version:
            try:
                # Simple version comparison without semver dependency
                def compare_versions(v1, v2):
                    def version_tuple(v):
                        # Handle pre-release versions (alpha, beta, rc)
                        base_version = v.split('-')[0]
                        return tuple(map(int, base_version.split(".")))
                    return version_tuple(v1) < version_tuple(v2)

                # Check if config version is compatible with schema
                if compare_versions(config_version, "0.2.0"):
                    return False, f"Config version {config_version} is deprecated (minimum 0.2.0)"
                elif compare_versions(schema_version, config_version):
                    return False, f"Config version {config_version} is newer than schema {schema_version}"
            except Exception as e:
                return False, f"Version validation error: {e}"

        # Standard JSON schema validation
        jsonschema.validate(config, schema)
        return True, None
    except jsonschema.ValidationError as e:
        return False, f"Schema validation error: {e.message}"
    except yaml.YAMLError as e:
        return False, f"YAML syntax error: {e}"
    except Exception as e:
        return False, f"Error: {e}"

def main():
    """Main validation function"""
    root_dir = Path(__file__).parent.parent
    configs_dir = root_dir / "configs" / "domains"
    schema_path = root_dir / "configs" / "schemas" / "domain_pack_schema.yaml"

    if not schema_path.exists():
        print(f"ERROR: Schema file not found: {schema_path}")
        return 1

    if not configs_dir.exists():
        print(f"ERROR: Configs directory not found: {configs_dir}")
        return 1

    # Find all YAML files in configs directory
    yaml_files = list(configs_dir.glob("*.yaml"))
    if not yaml_files:
        print(f"WARNING: No YAML files found in {configs_dir}")
        return 0

    errors = 0
    for config_file in yaml_files:
        is_valid, error_msg = validate_domain_config(config_file, schema_path)
        if is_valid:
            print(f"✓ {config_file.name}")
        else:
            print(f"✗ {config_file.name}: {error_msg}")
            errors += 1

    print(f"\nValidated {len(yaml_files)} domain configurations")
    if errors > 0:
        print(f"Found {errors} validation errors")
        return 1
    else:
        print("All configurations are valid!")
        return 0

if __name__ == "__main__":
    sys.exit(main())
