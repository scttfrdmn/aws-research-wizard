"""
Shared pytest configuration and fixtures for AWS Research Wizard tests.

This module provides common test fixtures, mock objects, and configuration
that can be used across all test modules. It includes AWS service mocks,
test data generators, and utility functions for testing.
"""

import json
import os
import tempfile
from typing import Any, Dict, Generator, List
from unittest.mock import MagicMock, Mock, patch

import boto3
import pytest
from moto import mock_aws
from moto.core import DEFAULT_ACCOUNT_ID as ACCOUNT_ID


@pytest.fixture(scope="session")
def aws_credentials() -> None:
    """Mock AWS credentials for testing."""
    os.environ.setdefault("AWS_ACCESS_KEY_ID", "testing")
    os.environ.setdefault("AWS_SECRET_ACCESS_KEY", "testing")
    os.environ.setdefault("AWS_SECURITY_TOKEN", "testing")
    os.environ.setdefault("AWS_SESSION_TOKEN", "testing")
    os.environ.setdefault("AWS_DEFAULT_REGION", "us-east-1")


@pytest.fixture
def temp_dir() -> Generator[str, None, None]:
    """Create a temporary directory for test files."""
    with tempfile.TemporaryDirectory() as tmp_dir:
        yield tmp_dir


@pytest.fixture
def sample_workload_config() -> Dict[str, Any]:
    """Sample workload configuration for testing."""
    return {
        "domain": "genomics",
        "primary_tools": ["GATK", "BWA", "STAR"],
        "problem_size": "large",
        "priority": "performance",
        "data_size_gb": 1000,
        "collaboration_users": 5,
        "deadline_hours": 72,
        "budget_limit": 5000,
        "gpu_requirement": "none",
        "memory_intensity": "high",
        "io_pattern": "random",
        "parallel_scaling": "linear"
    }


@pytest.fixture
def sample_aws_instance_types() -> List[Dict[str, Any]]:
    """Sample AWS instance type data for testing."""
    return [
        {
            "InstanceType": "c6i.large",
            "VCpuInfo": {"DefaultVCpus": 2},
            "MemoryInfo": {"SizeInMiB": 4096},
            "ProcessorInfo": {"SustainedClockSpeedInGhz": 3.5},
            "NetworkInfo": {"NetworkPerformance": "Up to 12.5 Gigabit"},
            "EbsInfo": {"EbsOptimizedSupport": "default"},
        },
        {
            "InstanceType": "r6i.xlarge",
            "VCpuInfo": {"DefaultVCpus": 4},
            "MemoryInfo": {"SizeInMiB": 32768},
            "ProcessorInfo": {"SustainedClockSpeedInGhz": 3.5},
            "NetworkInfo": {"NetworkPerformance": "Up to 12.5 Gigabit"},
            "EbsInfo": {"EbsOptimizedSupport": "default"},
        },
        {
            "InstanceType": "g5.2xlarge",
            "VCpuInfo": {"DefaultVCpus": 8},
            "MemoryInfo": {"SizeInMiB": 32768},
            "ProcessorInfo": {"SustainedClockSpeedInGhz": 3.4},
            "NetworkInfo": {"NetworkPerformance": "Up to 25 Gigabit"},
            "GpuInfo": {
                "Gpus": [
                    {"Name": "A10G", "Manufacturer": "NVIDIA", "Count": 1, "MemoryInfo": {"SizeInMiB": 24576}}
                ]
            },
        },
    ]


@pytest.fixture
def mock_boto3_session():
    """Mock boto3 session for testing."""
    session = Mock()
    session.region_name = "us-east-1"
    session.profile_name = "default"
    
    # Mock STS client
    sts_client = Mock()
    sts_client.get_caller_identity.return_value = {
        "Account": ACCOUNT_ID,
        "UserId": "AIDACKCEVSQ6C2EXAMPLE",
        "Arn": "arn:aws:iam::123456789012:user/test-user"
    }
    
    # Mock EC2 client
    ec2_client = Mock()
    ec2_client.describe_instances.return_value = {"Reservations": []}
    ec2_client.describe_regions.return_value = {
        "Regions": [
            {"RegionName": "us-east-1", "Endpoint": "ec2.us-east-1.amazonaws.com"},
            {"RegionName": "us-west-2", "Endpoint": "ec2.us-west-2.amazonaws.com"}
        ]
    }
    
    # Mock S3 client
    s3_client = Mock()
    s3_client.list_buckets.return_value = {"Buckets": []}
    
    # Client factory method
    def client(service_name, **kwargs):
        if service_name == "sts":
            return sts_client
        elif service_name == "ec2":
            return ec2_client
        elif service_name == "s3":
            return s3_client
        else:
            return Mock()
    
    session.client = client
    return session


@pytest.fixture
def mock_aws_services():
    """Mock AWS services using moto."""
    with mock_aws():
        yield


@pytest.fixture
def sample_spack_environment() -> Dict[str, Any]:
    """Sample Spack environment for testing."""
    return {
        "spack": {
            "specs": ["python@3.11", "gcc@11.4.0", "openmpi@4.1.5"],
            "config": {
                "install_tree": {
                    "root": "/opt/spack",
                    "projections": {"all": "{name}-{version}-{hash:7}"}
                }
            },
            "packages": {
                "all": {
                    "compiler": ["gcc@11.4.0"],
                    "target": ["x86_64"]
                }
            },
            "compilers": [
                {
                    "compiler": {
                        "spec": "gcc@11.4.0",
                        "paths": {
                            "cc": "/usr/bin/gcc",
                            "cxx": "/usr/bin/g++",
                            "f77": "/usr/bin/gfortran",
                            "fc": "/usr/bin/gfortran"
                        },
                        "flags": {},
                        "operating_system": "ubuntu22",
                        "target": "x86_64",
                        "modules": [],
                        "environment": {}
                    }
                }
            ]
        }
    }


@pytest.fixture
def sample_research_pack_config() -> Dict[str, Any]:
    """Sample research pack configuration for testing."""
    return {
        "name": "Test Research Pack",
        "description": "A test research pack for unit testing",
        "spack_packages": [
            "python@3.11.5 %gcc@11.4.0",
            "numpy@1.25.2 %gcc@11.4.0",
            "scipy@1.11.2 %gcc@11.4.0"
        ],
        "aws_instance_recommendations": {
            "small": {
                "instance_type": "c6i.large",
                "vcpus": 2,
                "memory_gb": 4,
                "storage_gb": 100,
                "cost_per_hour": 0.085,
                "use_case": "Development and testing"
            },
            "large": {
                "instance_type": "c6i.4xlarge",
                "vcpus": 16,
                "memory_gb": 32,
                "storage_gb": 500,
                "cost_per_hour": 0.68,
                "use_case": "Production workloads"
            }
        },
        "estimated_cost": {
            "compute": 500,
            "storage": 100,
            "data_transfer": 50,
            "total": 650
        },
        "research_capabilities": [
            "Scientific computing with Python",
            "Numerical analysis and linear algebra",
            "Data processing and visualization"
        ]
    }


@pytest.fixture
def mock_streamlit():
    """Mock Streamlit components for testing GUI."""
    with patch("streamlit.title") as mock_title, \
         patch("streamlit.write") as mock_write, \
         patch("streamlit.selectbox") as mock_selectbox, \
         patch("streamlit.button") as mock_button, \
         patch("streamlit.columns") as mock_columns:
        
        mock_selectbox.return_value = "test_selection"
        mock_button.return_value = False
        mock_columns.return_value = [Mock(), Mock()]
        
        yield {
            "title": mock_title,
            "write": mock_write,
            "selectbox": mock_selectbox,
            "button": mock_button,
            "columns": mock_columns
        }


@pytest.fixture
def environment_check_results() -> List[Dict[str, Any]]:
    """Sample environment check results for testing."""
    return [
        {
            "name": "AWS Credentials",
            "status": "PASS",
            "message": "Successfully authenticated",
            "recommendation": None,
            "details": {"account_id": ACCOUNT_ID}
        },
        {
            "name": "IAM Permissions",
            "status": "WARN",
            "message": "Limited permissions detected",
            "recommendation": "Review IAM policy",
            "details": {"missing_permissions": ["s3:ListAllMyBuckets"]}
        },
        {
            "name": "Service Quotas",
            "status": "FAIL",
            "message": "Insufficient quotas",
            "recommendation": "Request quota increases",
            "details": {"failed_quotas": ["ec2:running-on-demand-instances"]}
        }
    ]


@pytest.fixture
def mock_subprocess():
    """Mock subprocess calls for testing command execution."""
    with patch("subprocess.run") as mock_run, \
         patch("subprocess.Popen") as mock_popen:
        
        # Default successful response
        mock_run.return_value = Mock(
            returncode=0,
            stdout="Mock command output",
            stderr="",
            text=True
        )
        
        yield {"run": mock_run, "popen": mock_popen}


@pytest.fixture
def sample_benchmark_results() -> Dict[str, Any]:
    """Sample benchmark results for testing."""
    return {
        "benchmark_suite": "HPC Challenge",
        "timestamp": "2023-08-15T10:30:00Z",
        "system_info": {
            "instance_type": "c6i.8xlarge",
            "vcpus": 32,
            "memory_gb": 64,
            "network": "25 Gbps"
        },
        "results": {
            "hpl": {
                "performance_gflops": 1250.5,
                "efficiency_percent": 78.2,
                "execution_time_seconds": 3600
            },
            "stream": {
                "copy_bandwidth_gbps": 45.3,
                "scale_bandwidth_gbps": 44.1,
                "add_bandwidth_gbps": 42.8,
                "triad_bandwidth_gbps": 43.5
            },
            "random_access": {
                "gups": 0.125,
                "execution_time_seconds": 300
            }
        },
        "cost_analysis": {
            "compute_cost_usd": 2.72,
            "total_runtime_hours": 2.0,
            "cost_per_gflop": 0.00218
        }
    }


@pytest.fixture(autouse=True)
def reset_environment():
    """Reset environment variables between tests."""
    # Store original values
    original_env = os.environ.copy()
    
    yield
    
    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)


# Custom pytest markers for test categorization
def pytest_configure(config):
    """Configure custom pytest markers."""
    config.addinivalue_line("markers", "unit: mark test as a unit test")
    config.addinivalue_line("markers", "integration: mark test as an integration test")
    config.addinivalue_line("markers", "slow: mark test as slow running")
    config.addinivalue_line("markers", "aws: mark test as requiring AWS credentials")
    config.addinivalue_line("markers", "gui: mark test as testing GUI components")


# Test data generators
def generate_instance_recommendations(count: int = 5) -> List[Dict[str, Any]]:
    """Generate test instance recommendations."""
    instances = []
    for i in range(count):
        instances.append({
            "name": f"test_instance_{i}",
            "instance_type": f"c6i.{2**i}xlarge",
            "vcpus": 2 * (2**i),
            "memory_gb": 4 * (2**i),
            "storage_gb": 100 * (i + 1),
            "cost_per_hour": 0.085 * (2**i),
            "use_case": f"Test use case {i}"
        })
    return instances


def generate_workload_characteristics(**overrides) -> Dict[str, Any]:
    """Generate test workload characteristics with optional overrides."""
    base_config = {
        "domain": "genomics",
        "primary_tools": ["GATK", "BWA"],
        "problem_size": "medium",
        "priority": "balanced",
        "data_size_gb": 500,
        "collaboration_users": 3,
        "deadline_hours": None,
        "budget_limit": None,
        "gpu_requirement": "none",
        "memory_intensity": "medium",
        "io_pattern": "sequential",
        "parallel_scaling": "linear"
    }
    base_config.update(overrides)
    return base_config


def generate_cost_estimates(**overrides) -> Dict[str, float]:
    """Generate test cost estimates with optional overrides."""
    base_costs = {
        "compute": 400.0,
        "storage": 100.0,
        "network": 50.0,
        "data_transfer": 25.0,
        "total": 575.0
    }
    base_costs.update(overrides)
    return base_costs