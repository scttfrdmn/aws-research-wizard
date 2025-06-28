"""
Unit tests for the core Research Infrastructure Wizard.

This module tests the main wizard functionality including workload analysis,
infrastructure recommendations, cost calculations, and optimization logic.
"""

import pytest
from unittest.mock import Mock, patch, MagicMock
from typing import Dict, Any

# Import the module under test
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..'))

from research_infrastructure_wizard import (
    ResearchInfrastructureWizard,
    WorkloadCharacteristics,
    Priority,
    WorkloadSize,
    InfrastructureRecommendation
)


class TestWorkloadCharacteristics:
    """Test the WorkloadCharacteristics dataclass."""
    
    def test_workload_creation_with_valid_data(self, sample_workload_config):
        """Test creating a valid WorkloadCharacteristics instance."""
        workload = WorkloadCharacteristics(
            domain=sample_workload_config["domain"],
            primary_tools=sample_workload_config["primary_tools"],
            problem_size=WorkloadSize(sample_workload_config["problem_size"]),
            priority=Priority(sample_workload_config["priority"]),
            data_size_gb=sample_workload_config["data_size_gb"],
            collaboration_users=sample_workload_config["collaboration_users"],
            deadline_hours=sample_workload_config["deadline_hours"],
            budget_limit=sample_workload_config["budget_limit"],
            gpu_requirement=sample_workload_config["gpu_requirement"],
            memory_intensity=sample_workload_config["memory_intensity"],
            io_pattern=sample_workload_config["io_pattern"],
            parallel_scaling=sample_workload_config["parallel_scaling"]
        )
        
        assert workload.domain == "genomics"
        assert workload.primary_tools == ["GATK", "BWA", "STAR"]
        assert workload.problem_size == WorkloadSize.LARGE
        assert workload.priority == Priority.PERFORMANCE
        assert workload.data_size_gb == 1000
        assert workload.collaboration_users == 5
    
    def test_workload_creation_with_minimal_data(self):
        """Test creating a WorkloadCharacteristics with minimal required data."""
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK"],
            problem_size=WorkloadSize.SMALL,
            priority=Priority.COST,
            data_size_gb=100,
            collaboration_users=1
        )
        
        assert workload.domain == "genomics"
        assert workload.primary_tools == ["GATK"]
        assert workload.problem_size == WorkloadSize.SMALL
        assert workload.priority == Priority.COST
        assert workload.data_size_gb == 100
        assert workload.collaboration_users == 1
        # Test default values
        assert workload.deadline_hours is None
        assert workload.budget_limit is None
        assert workload.gpu_requirement == "none"
        assert workload.memory_intensity == "medium"
        assert workload.io_pattern == "sequential"
        assert workload.parallel_scaling == "linear"
    
    def test_workload_validation_errors(self):
        """Test workload validation with invalid data."""
        # Test invalid problem size
        with pytest.raises(ValueError):
            WorkloadSize("invalid_size")
        
        # Test invalid priority
        with pytest.raises(ValueError):
            Priority("invalid_priority")
    
    def test_workload_size_enum_values(self):
        """Test WorkloadSize enum values."""
        assert WorkloadSize.SMALL.value == "small"
        assert WorkloadSize.MEDIUM.value == "medium"
        assert WorkloadSize.LARGE.value == "large"
        assert WorkloadSize.MASSIVE.value == "massive"
    
    def test_priority_enum_values(self):
        """Test Priority enum values."""
        assert Priority.COST.value == "cost"
        assert Priority.PERFORMANCE.value == "performance"
        assert Priority.BALANCED.value == "balanced"
        assert Priority.DEADLINE.value == "deadline"


class TestResearchInfrastructureWizard:
    """Test the main ResearchInfrastructureWizard class."""
    
    def setup_method(self):
        """Set up test fixtures before each test method."""
        self.wizard = ResearchInfrastructureWizard()
    
    def test_wizard_initialization(self):
        """Test wizard initialization with default values."""
        assert self.wizard.aws_region == "us-east-1"
        assert self.wizard.cost_optimization_enabled is True
        assert self.wizard.performance_target == "balanced"
        assert len(self.wizard.aws_instance_catalog) > 0
    
    def test_wizard_initialization_with_custom_region(self):
        """Test wizard initialization with custom region."""
        wizard = ResearchInfrastructureWizard(aws_region="us-west-2")
        assert wizard.aws_region == "us-west-2"
    
    @patch('research_infrastructure_wizard.ResearchInfrastructureWizard._load_aws_pricing')
    def test_wizard_initialization_with_mocked_pricing(self, mock_pricing):
        """Test wizard initialization with mocked AWS pricing."""
        mock_pricing.return_value = {"c6i.large": 0.085, "r6i.xlarge": 0.255}
        
        wizard = ResearchInfrastructureWizard()
        assert wizard.aws_instance_catalog is not None
        mock_pricing.assert_called_once()
    
    def test_generate_recommendation_with_genomics_workload(self, sample_workload_config):
        """Test generating recommendations for a genomics workload."""
        workload = WorkloadCharacteristics(
            domain=sample_workload_config["domain"],
            primary_tools=sample_workload_config["primary_tools"],
            problem_size=WorkloadSize.LARGE,
            priority=Priority.PERFORMANCE,
            data_size_gb=sample_workload_config["data_size_gb"],
            collaboration_users=sample_workload_config["collaboration_users"],
            deadline_hours=sample_workload_config["deadline_hours"],
            budget_limit=sample_workload_config["budget_limit"],
            gpu_requirement=sample_workload_config["gpu_requirement"],
            memory_intensity=sample_workload_config["memory_intensity"],
            io_pattern=sample_workload_config["io_pattern"],
            parallel_scaling=sample_workload_config["parallel_scaling"]
        )
        
        with patch.object(self.wizard, '_analyze_workload_requirements') as mock_analyze, \
             patch.object(self.wizard, '_select_optimal_instances') as mock_select, \
             patch.object(self.wizard, '_calculate_costs') as mock_costs, \
             patch.object(self.wizard, '_generate_deployment_plan') as mock_deploy:
            
            mock_analyze.return_value = {
                "cpu_requirements": {"min_vcpus": 16, "recommended_vcpus": 32},
                "memory_requirements": {"min_gb": 64, "recommended_gb": 128},
                "storage_requirements": {"min_gb": 1000, "recommended_gb": 2000},
                "network_requirements": {"bandwidth": "high", "latency": "low"}
            }
            
            mock_select.return_value = {
                "primary": {"instance_type": "r6i.4xlarge", "count": 1},
                "worker": {"instance_type": "c6i.2xlarge", "count": 4}
            }
            
            mock_costs.return_value = {
                "compute": 500.0,
                "storage": 150.0,
                "network": 75.0,
                "total": 725.0
            }
            
            mock_deploy.return_value = {
                "infrastructure": {"vpc": True, "subnets": 3},
                "scaling": {"min_nodes": 1, "max_nodes": 10},
                "optimization": ["spot_instances", "auto_scaling"]
            }
            
            recommendation = self.wizard.generate_recommendation(workload)
            
            # Verify all methods were called
            mock_analyze.assert_called_once_with(workload)
            mock_select.assert_called_once()
            mock_costs.assert_called_once()
            mock_deploy.assert_called_once()
            
            # Verify recommendation structure
            assert "workload_analysis" in recommendation
            assert "infrastructure_recommendation" in recommendation
            assert "cost_estimate" in recommendation
            assert "deployment_plan" in recommendation
            assert "optimization_suggestions" in recommendation
    
    def test_analyze_workload_requirements_genomics(self):
        """Test workload requirement analysis for genomics domain."""
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK", "BWA", "STAR"],
            problem_size=WorkloadSize.LARGE,
            priority=Priority.PERFORMANCE,
            data_size_gb=1000,
            collaboration_users=5,
            memory_intensity="high",
            io_pattern="random",
            parallel_scaling="linear"
        )
        
        requirements = self.wizard._analyze_workload_requirements(workload)
        
        assert "cpu_requirements" in requirements
        assert "memory_requirements" in requirements
        assert "storage_requirements" in requirements
        assert "network_requirements" in requirements
        
        # Check that genomics-specific requirements are applied
        assert requirements["memory_requirements"]["min_gb"] >= 32  # High memory for genomics
        assert requirements["storage_requirements"]["iops"] >= 3000  # High IOPS for random I/O
    
    def test_analyze_workload_requirements_machine_learning(self):
        """Test workload requirement analysis for machine learning domain."""
        workload = WorkloadCharacteristics(
            domain="machine_learning",
            primary_tools=["PyTorch", "TensorFlow"],
            problem_size=WorkloadSize.LARGE,
            priority=Priority.PERFORMANCE,
            data_size_gb=500,
            collaboration_users=3,
            gpu_requirement="required",
            memory_intensity="high",
            parallel_scaling="linear"
        )
        
        requirements = self.wizard._analyze_workload_requirements(workload)
        
        # Check GPU-specific requirements
        assert requirements["gpu_requirements"]["required"] is True
        assert requirements["gpu_requirements"]["memory_gb"] >= 16
    
    def test_select_optimal_instances_cost_priority(self):
        """Test instance selection with cost optimization priority."""
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK"],
            problem_size=WorkloadSize.MEDIUM,
            priority=Priority.COST,
            data_size_gb=500,
            collaboration_users=2
        )
        
        requirements = {
            "cpu_requirements": {"min_vcpus": 8, "recommended_vcpus": 16},
            "memory_requirements": {"min_gb": 32, "recommended_gb": 64},
            "storage_requirements": {"min_gb": 500, "recommended_gb": 1000}
        }
        
        with patch.object(self.wizard, '_get_instance_pricing') as mock_pricing:
            mock_pricing.return_value = {
                "c6i.2xlarge": 0.34,
                "r6i.2xlarge": 0.51,
                "m6i.2xlarge": 0.38
            }
            
            instances = self.wizard._select_optimal_instances(workload, requirements)
            
            # Should prioritize cost-effective instances
            assert instances is not None
            assert "primary" in instances
    
    def test_select_optimal_instances_performance_priority(self):
        """Test instance selection with performance optimization priority."""
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK"],
            problem_size=WorkloadSize.LARGE,
            priority=Priority.PERFORMANCE,
            data_size_gb=1000,
            collaboration_users=5
        )
        
        requirements = {
            "cpu_requirements": {"min_vcpus": 16, "recommended_vcpus": 32},
            "memory_requirements": {"min_gb": 64, "recommended_gb": 128},
            "storage_requirements": {"min_gb": 1000, "recommended_gb": 2000}
        }
        
        instances = self.wizard._select_optimal_instances(workload, requirements)
        
        # Should prioritize high-performance instances
        assert instances is not None
        assert "primary" in instances
    
    def test_calculate_costs_basic(self):
        """Test basic cost calculation."""
        instances = {
            "primary": {"instance_type": "r6i.4xlarge", "count": 1},
            "worker": {"instance_type": "c6i.2xlarge", "count": 2}
        }
        
        storage_gb = 1000
        
        with patch.object(self.wizard, '_get_instance_pricing') as mock_pricing, \
             patch.object(self.wizard, '_get_storage_pricing') as mock_storage_pricing:
            
            mock_pricing.return_value = {
                "r6i.4xlarge": 1.02,
                "c6i.2xlarge": 0.34
            }
            mock_storage_pricing.return_value = {"gp3": 0.08}  # per GB per month
            
            costs = self.wizard._calculate_costs(instances, storage_gb)
            
            assert "compute" in costs
            assert "storage" in costs
            assert "network" in costs
            assert "total" in costs
            assert costs["total"] > 0
    
    def test_calculate_costs_with_spot_instances(self):
        """Test cost calculation with spot instance optimization."""
        instances = {
            "primary": {"instance_type": "r6i.4xlarge", "count": 1, "spot": True},
            "worker": {"instance_type": "c6i.2xlarge", "count": 2, "spot": True}
        }
        
        storage_gb = 500
        
        with patch.object(self.wizard, '_get_instance_pricing') as mock_pricing, \
             patch.object(self.wizard, '_get_storage_pricing') as mock_storage_pricing, \
             patch.object(self.wizard, '_get_spot_pricing') as mock_spot_pricing:
            
            mock_pricing.return_value = {
                "r6i.4xlarge": 1.02,
                "c6i.2xlarge": 0.34
            }
            mock_storage_pricing.return_value = {"gp3": 0.08}
            mock_spot_pricing.return_value = {
                "r6i.4xlarge": 0.31,  # ~70% discount
                "c6i.2xlarge": 0.10   # ~70% discount
            }
            
            costs = self.wizard._calculate_costs(instances, storage_gb)
            
            # Spot instances should result in lower compute costs
            assert costs["compute"] > 0
            assert "spot_savings" in costs
    
    def test_generate_deployment_plan_basic(self):
        """Test deployment plan generation."""
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK"],
            problem_size=WorkloadSize.MEDIUM,
            priority=Priority.BALANCED,
            data_size_gb=500,
            collaboration_users=3
        )
        
        instances = {
            "primary": {"instance_type": "r6i.2xlarge", "count": 1},
            "worker": {"instance_type": "c6i.large", "count": 2}
        }
        
        deployment_plan = self.wizard._generate_deployment_plan(workload, instances)
        
        assert "infrastructure" in deployment_plan
        assert "scaling" in deployment_plan
        assert "security" in deployment_plan
        assert "monitoring" in deployment_plan
        assert "backup" in deployment_plan
    
    def test_generate_optimization_suggestions_cost_focused(self):
        """Test optimization suggestions for cost-focused workloads."""
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK"],
            problem_size=WorkloadSize.LARGE,
            priority=Priority.COST,
            data_size_gb=1000,
            collaboration_users=5
        )
        
        cost_estimate = {
            "compute": 800.0,
            "storage": 200.0,
            "network": 100.0,
            "total": 1100.0
        }
        
        suggestions = self.wizard._generate_optimization_suggestions(workload, cost_estimate)
        
        assert len(suggestions) > 0
        # Should include cost optimization suggestions
        cost_suggestions = [s for s in suggestions if "cost" in s.lower() or "spot" in s.lower()]
        assert len(cost_suggestions) > 0
    
    def test_generate_optimization_suggestions_performance_focused(self):
        """Test optimization suggestions for performance-focused workloads."""
        workload = WorkloadCharacteristics(
            domain="machine_learning",
            primary_tools=["PyTorch"],
            problem_size=WorkloadSize.LARGE,
            priority=Priority.PERFORMANCE,
            data_size_gb=500,
            collaboration_users=3,
            gpu_requirement="required"
        )
        
        cost_estimate = {
            "compute": 1500.0,
            "storage": 150.0,
            "network": 75.0,
            "total": 1725.0
        }
        
        suggestions = self.wizard._generate_optimization_suggestions(workload, cost_estimate)
        
        assert len(suggestions) > 0
        # Should include performance optimization suggestions
        perf_suggestions = [s for s in suggestions if "performance" in s.lower() or "gpu" in s.lower()]
        assert len(perf_suggestions) > 0
    
    def test_domain_specific_optimizations_genomics(self):
        """Test genomics-specific optimizations."""
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK", "BWA", "STAR"],
            problem_size=WorkloadSize.LARGE,
            priority=Priority.PERFORMANCE,
            data_size_gb=2000,
            collaboration_users=5,
            memory_intensity="high",
            io_pattern="random"
        )
        
        optimizations = self.wizard._get_domain_specific_optimizations(workload)
        
        assert len(optimizations) > 0
        # Should include genomics-specific recommendations
        genomics_opts = [opt for opt in optimizations if any(
            term in opt.lower() for term in ["memory", "storage", "iops", "parallel"]
        )]
        assert len(genomics_opts) > 0
    
    def test_domain_specific_optimizations_climate_modeling(self):
        """Test climate modeling-specific optimizations."""
        workload = WorkloadCharacteristics(
            domain="climate_modeling",
            primary_tools=["WRF", "CESM"],
            problem_size=WorkloadSize.MASSIVE,
            priority=Priority.PERFORMANCE,
            data_size_gb=5000,
            collaboration_users=10,
            parallel_scaling="linear"
        )
        
        optimizations = self.wizard._get_domain_specific_optimizations(workload)
        
        assert len(optimizations) > 0
        # Should include HPC-specific recommendations
        hpc_opts = [opt for opt in optimizations if any(
            term in opt.lower() for term in ["mpi", "network", "placement", "cluster"]
        )]
        assert len(hpc_opts) > 0
    
    def test_error_handling_invalid_workload(self):
        """Test error handling with invalid workload data."""
        # Test with None workload
        with pytest.raises(ValueError, match="Workload cannot be None"):
            self.wizard.generate_recommendation(None)
        
        # Test with missing required fields
        incomplete_workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=[],  # Empty tools list
            problem_size=WorkloadSize.SMALL,
            priority=Priority.COST,
            data_size_gb=0,  # Invalid size
            collaboration_users=0  # Invalid user count
        )
        
        with pytest.raises(ValueError):
            self.wizard.generate_recommendation(incomplete_workload)
    
    def test_aws_region_specific_pricing(self):
        """Test that different AWS regions affect pricing."""
        wizard_us_east = ResearchInfrastructureWizard(aws_region="us-east-1")
        wizard_eu_west = ResearchInfrastructureWizard(aws_region="eu-west-1")
        
        # Mock different pricing for different regions
        with patch.object(wizard_us_east, '_get_instance_pricing') as mock_us_pricing, \
             patch.object(wizard_eu_west, '_get_instance_pricing') as mock_eu_pricing:
            
            mock_us_pricing.return_value = {"c6i.large": 0.085}
            mock_eu_pricing.return_value = {"c6i.large": 0.095}  # Higher price in EU
            
            us_price = wizard_us_east._get_instance_pricing()["c6i.large"]
            eu_price = wizard_eu_west._get_instance_pricing()["c6i.large"]
            
            assert us_price != eu_price
    
    @pytest.mark.slow
    def test_recommendation_performance_large_workload(self):
        """Test recommendation generation performance with large workload."""
        import time
        
        workload = WorkloadCharacteristics(
            domain="genomics",
            primary_tools=["GATK", "BWA", "STAR", "SAMtools", "BCFtools"],
            problem_size=WorkloadSize.MASSIVE,
            priority=Priority.PERFORMANCE,
            data_size_gb=10000,
            collaboration_users=20,
            deadline_hours=24,
            budget_limit=10000,
            memory_intensity="high",
            io_pattern="mixed",
            parallel_scaling="linear"
        )
        
        start_time = time.time()
        recommendation = self.wizard.generate_recommendation(workload)
        end_time = time.time()
        
        # Should complete within reasonable time (< 5 seconds)
        assert end_time - start_time < 5.0
        assert recommendation is not None
        assert "workload_analysis" in recommendation


class TestResearchDomainEnum:
    """Test the ResearchDomain enumeration."""
    
    def test_research_domain_values(self):
        """Test that all expected research domains are defined."""
        expected_domains = [
            "genomics", "climate_modeling", "materials_science", 
            "machine_learning", "physics_simulation"
        ]
        
        for domain in expected_domains:
            assert hasattr(ResearchDomain, domain.upper())
    
    def test_research_domain_string_conversion(self):
        """Test research domain string conversion."""
        assert ResearchDomain.GENOMICS.value == "genomics"
        assert ResearchDomain.CLIMATE_MODELING.value == "climate_modeling"
        assert ResearchDomain.MACHINE_LEARNING.value == "machine_learning"