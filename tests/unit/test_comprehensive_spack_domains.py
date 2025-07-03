"""
Unit tests for comprehensive Spack domains module.

Tests the research domain configurations, Spack environment generation,
and domain-specific optimizations.
"""

import pytest
from unittest.mock import Mock, patch, MagicMock
from typing import Dict, Any

import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..'))

from comprehensive_spack_domains import (
    get_research_domains,
    generate_spack_environment,
    optimize_for_aws_graviton,
    create_genomics_environment,
    create_climate_modeling_environment,
    create_machine_learning_environment,
    create_hpc_environment
)


class TestResearchDomains:
    """Test research domain functionality."""

    def test_get_research_domains_returns_dict(self):
        """Test that get_research_domains returns a dictionary."""
        domains = get_research_domains()
        assert isinstance(domains, dict)
        assert len(domains) > 0

    def test_research_domains_have_required_fields(self):
        """Test that research domains have required fields."""
        domains = get_research_domains()
        for domain_name, domain_config in domains.items():
            assert isinstance(domain_name, str)
            assert isinstance(domain_config, dict)
            assert 'description' in domain_config
            assert 'spack_packages' in domain_config
            assert 'aws_instance_recommendations' in domain_config

    def test_genomics_domain_exists(self):
        """Test that genomics domain is present."""
        domains = get_research_domains()
        assert 'genomics' in domains
        genomics = domains['genomics']
        assert 'GATK' in str(genomics.get('spack_packages', ''))

    def test_climate_modeling_domain_exists(self):
        """Test that climate modeling domain is present."""
        domains = get_research_domains()
        assert 'climate_modeling' in domains
        climate = domains['climate_modeling']
        assert 'wrf' in str(climate.get('spack_packages', '')).lower()


class TestSpackEnvironmentGeneration:
    """Test Spack environment generation."""

    def test_generate_spack_environment_basic(self):
        """Test basic Spack environment generation."""
        packages = ['python@3.11', 'numpy@1.25.0']
        env = generate_spack_environment('test_env', packages)

        assert isinstance(env, dict)
        assert 'spack' in env
        assert 'specs' in env['spack']
        assert packages == env['spack']['specs']

    def test_generate_spack_environment_with_compilers(self):
        """Test Spack environment generation with compiler specs."""
        packages = ['python@3.11 %gcc@11']
        env = generate_spack_environment('test_env', packages)

        assert isinstance(env, dict)
        assert packages == env['spack']['specs']

    def test_create_genomics_environment(self):
        """Test genomics-specific environment creation."""
        env = create_genomics_environment()

        assert isinstance(env, dict)
        assert 'spack' in env
        # Check for genomics-specific packages
        specs = str(env['spack']['specs'])
        assert 'gatk' in specs.lower() or 'bwa' in specs.lower()

    def test_create_climate_modeling_environment(self):
        """Test climate modeling environment creation."""
        env = create_climate_modeling_environment()

        assert isinstance(env, dict)
        assert 'spack' in env
        # Check for climate-specific packages
        specs = str(env['spack']['specs'])
        assert 'wrf' in specs.lower() or 'netcdf' in specs.lower()

    def test_create_machine_learning_environment(self):
        """Test machine learning environment creation."""
        env = create_machine_learning_environment()

        assert isinstance(env, dict)
        assert 'spack' in env
        # Check for ML-specific packages
        specs = str(env['spack']['specs'])
        assert 'pytorch' in specs.lower() or 'python' in specs.lower()

    def test_create_hpc_environment(self):
        """Test HPC environment creation."""
        env = create_hpc_environment()

        assert isinstance(env, dict)
        assert 'spack' in env
        # Check for HPC-specific packages
        specs = str(env['spack']['specs'])
        assert 'openmpi' in specs.lower() or 'mpi' in specs.lower()


class TestGravitonOptimization:
    """Test AWS Graviton optimization features."""

    def test_optimize_for_aws_graviton_basic(self):
        """Test basic Graviton optimization."""
        base_env = {
            'spack': {
                'specs': ['python@3.11', 'numpy@1.25.0'],
                'packages': {}
            }
        }

        optimized = optimize_for_aws_graviton(base_env)

        assert isinstance(optimized, dict)
        assert 'spack' in optimized
        # Should have target architecture specifications
        if 'packages' in optimized['spack']:
            packages_config = optimized['spack']['packages']
            assert isinstance(packages_config, dict)

    def test_optimize_for_aws_graviton_with_compilers(self):
        """Test Graviton optimization with compiler configuration."""
        base_env = {
            'spack': {
                'specs': ['gcc@11.4.0', 'python@3.11 %gcc@11.4.0'],
                'packages': {},
                'compilers': []
            }
        }

        optimized = optimize_for_aws_graviton(base_env)

        assert isinstance(optimized, dict)
        # Should maintain or enhance compiler configuration
        if 'compilers' in optimized['spack']:
            assert isinstance(optimized['spack']['compilers'], list)


class TestDomainSpecificConfigurations:
    """Test domain-specific configuration generation."""

    @patch('comprehensive_spack_domains.generate_spack_environment')
    def test_genomics_configuration_generation(self, mock_generate):
        """Test genomics configuration generation."""
        mock_generate.return_value = {'spack': {'specs': ['gatk@4.4.0']}}

        result = create_genomics_environment()

        assert isinstance(result, dict)
        mock_generate.assert_called()

    @patch('comprehensive_spack_domains.generate_spack_environment')
    def test_climate_configuration_generation(self, mock_generate):
        """Test climate modeling configuration generation."""
        mock_generate.return_value = {'spack': {'specs': ['wrf@4.5.0']}}

        result = create_climate_modeling_environment()

        assert isinstance(result, dict)
        mock_generate.assert_called()


class TestEnvironmentValidation:
    """Test Spack environment validation."""

    def test_environment_has_required_structure(self):
        """Test that generated environments have required structure."""
        packages = ['python@3.11']
        env = generate_spack_environment('test', packages)

        # Required top-level structure
        assert 'spack' in env
        assert isinstance(env['spack'], dict)

        # Required spack structure
        spack_config = env['spack']
        assert 'specs' in spack_config
        assert isinstance(spack_config['specs'], list)

    def test_empty_packages_list_handling(self):
        """Test handling of empty packages list."""
        env = generate_spack_environment('test', [])

        assert isinstance(env, dict)
        assert 'spack' in env
        assert env['spack']['specs'] == []

    def test_invalid_package_spec_handling(self):
        """Test handling of invalid package specifications."""
        # This should not crash even with unusual package specs
        packages = ['@invalid@spec', 'python@']
        env = generate_spack_environment('test', packages)

        assert isinstance(env, dict)
        assert 'spack' in env


class TestIntegrationWithAWS:
    """Test integration features with AWS services."""

    def test_aws_instance_recommendations_structure(self):
        """Test AWS instance recommendations structure."""
        domains = get_research_domains()
        for domain_name, domain_config in domains.items():
            aws_recs = domain_config.get('aws_instance_recommendations', {})
            if aws_recs:
                assert isinstance(aws_recs, dict)
                for size, config in aws_recs.items():
                    assert isinstance(config, dict)
                    if 'instance_type' in config:
                        assert isinstance(config['instance_type'], str)
                    if 'vcpus' in config:
                        assert isinstance(config['vcpus'], int)

    def test_cost_estimation_structure(self):
        """Test cost estimation data structure."""
        domains = get_research_domains()
        for domain_name, domain_config in domains.items():
            cost_est = domain_config.get('estimated_cost', {})
            if cost_est:
                assert isinstance(cost_est, dict)
                for cost_type, amount in cost_est.items():
                    if isinstance(amount, (int, float)):
                        assert amount >= 0  # Cost should be non-negative


@pytest.mark.integration
class TestSpackIntegration:
    """Integration tests with Spack (if available)."""

    def test_spack_environment_yaml_validity(self):
        """Test that generated environments produce valid YAML."""
        import yaml

        packages = ['python@3.11']
        env = generate_spack_environment('test', packages)

        # Should be serializable to YAML
        yaml_str = yaml.dump(env)
        assert isinstance(yaml_str, str)
        assert len(yaml_str) > 0

        # Should be deserializable from YAML
        loaded = yaml.safe_load(yaml_str)
        assert loaded == env
