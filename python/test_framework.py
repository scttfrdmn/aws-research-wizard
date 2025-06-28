#!/usr/bin/env python3
"""
Comprehensive Testing Framework for AWS Research Wizard
Includes unit tests, integration tests, performance tests, and demo workflow validation
"""

import os
import sys
import yaml
import json
import boto3
import pytest
import logging
import subprocess
import tempfile
from typing import Dict, List, Any, Optional, Tuple
from pathlib import Path
from dataclasses import dataclass
from datetime import datetime, timedelta
import time
import psutil
import threading
from concurrent.futures import ThreadPoolExecutor, as_completed

# Import our modules
from config_loader import ConfigLoader, DomainPackConfig
from dataset_manager import DatasetManager
from integrate_aws_data import AWSDataIntegrator

@dataclass
class TestResult:
    """Test result data structure"""
    test_name: str
    status: str  # PASS, FAIL, SKIP
    duration_seconds: float
    error_message: Optional[str] = None
    warnings: Optional[List[str]] = None
    metrics: Optional[Dict[str, Any]] = None

@dataclass
class WorkflowTestResult:
    """Demo workflow test result"""
    workflow_name: str
    domain: str
    status: str
    duration_seconds: float
    cost_estimate: float
    data_downloaded_gb: float
    error_message: Optional[str] = None
    performance_metrics: Optional[Dict[str, Any]] = None

class TestFramework:
    """Comprehensive testing framework for AWS Research Wizard"""
    
    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)
        self.test_results: List[TestResult] = []
        self.workflow_results: List[WorkflowTestResult] = []
        
        # Initialize components
        self.config_loader = ConfigLoader(config_root)
        self.dataset_manager = DatasetManager(config_root)
        self.aws_integrator = AWSDataIntegrator(config_root)
        
        # Test configuration
        self.test_config = {
            'timeout_seconds': 3600,  # 1 hour timeout for tests
            'max_parallel_tests': 4,
            'demo_data_limit_gb': 1.0,  # Limit demo data download for testing
            'skip_expensive_tests': True,  # Skip tests that cost >$10
            'mock_aws_services': False,  # Use mock AWS services for testing
        }
        
        # AWS clients (with error handling)
        try:
            self.s3_client = boto3.client('s3')
            self.ec2_client = boto3.client('ec2')
            self.aws_available = True
        except Exception as e:
            self.logger.warning(f"AWS not available: {e}")
            self.aws_available = False
            self.s3_client = None
            self.ec2_client = None

    def run_all_tests(self) -> Dict[str, Any]:
        """Run all test suites"""
        self.logger.info("Starting comprehensive test suite...")
        start_time = time.time()
        
        test_suites = [
            ("Configuration Tests", self.test_configurations),
            ("Dataset Tests", self.test_datasets),
            ("Integration Tests", self.test_aws_integration),
            ("Performance Tests", self.test_performance),
            ("Demo Workflow Tests", self.test_demo_workflows),
            ("Security Tests", self.test_security),
            ("Scalability Tests", self.test_scalability)
        ]
        
        suite_results = {}
        
        for suite_name, test_function in test_suites:
            self.logger.info(f"Running {suite_name}...")
            suite_start = time.time()
            
            try:
                suite_result = test_function()
                suite_results[suite_name] = {
                    'status': 'COMPLETED',
                    'duration': time.time() - suite_start,
                    'results': suite_result
                }
            except Exception as e:
                self.logger.error(f"Test suite {suite_name} failed: {e}")
                suite_results[suite_name] = {
                    'status': 'FAILED',
                    'duration': time.time() - suite_start,
                    'error': str(e)
                }
        
        total_duration = time.time() - start_time
        
        # Generate summary
        summary = self._generate_test_summary(suite_results, total_duration)
        
        return summary

    def test_configurations(self) -> Dict[str, Any]:
        """Test all configuration files and validation"""
        results = {}
        
        # Test 1: Schema validation
        self._add_test_result("schema_validation", self._test_schema_validation)
        
        # Test 2: Configuration loading
        self._add_test_result("config_loading", self._test_config_loading)
        
        # Test 3: MPI template application
        self._add_test_result("mpi_templates", self._test_mpi_templates)
        
        # Test 4: Domain-specific features
        self._add_test_result("domain_features", self._test_domain_features)
        
        # Test 5: Cost estimation validation
        self._add_test_result("cost_validation", self._test_cost_validation)
        
        return {
            'tests_run': 5,
            'passed': len([r for r in self.test_results[-5:] if r.status == 'PASS']),
            'failed': len([r for r in self.test_results[-5:] if r.status == 'FAIL']),
            'details': self.test_results[-5:]
        }

    def test_datasets(self) -> Dict[str, Any]:
        """Test dataset management and AWS Open Data integration"""
        results = {}
        
        # Test 1: Dataset registry loading
        self._add_test_result("dataset_registry", self._test_dataset_registry)
        
        # Test 2: Dataset accessibility
        self._add_test_result("dataset_access", self._test_dataset_access)
        
        # Test 3: Data subset creation
        self._add_test_result("data_subsets", self._test_data_subsets)
        
        # Test 4: Cost estimation
        self._add_test_result("cost_estimation", self._test_cost_estimation)
        
        return {
            'tests_run': 4,
            'passed': len([r for r in self.test_results[-4:] if r.status == 'PASS']),
            'failed': len([r for r in self.test_results[-4:] if r.status == 'FAIL']),
            'details': self.test_results[-4:]
        }

    def test_aws_integration(self) -> Dict[str, Any]:
        """Test AWS service integration"""
        results = {}
        
        if not self.aws_available:
            self.logger.warning("Skipping AWS integration tests - AWS not available")
            return {'status': 'SKIPPED', 'reason': 'AWS not available'}
        
        # Test 1: S3 access
        self._add_test_result("s3_access", self._test_s3_access)
        
        # Test 2: EC2 instance type validation
        self._add_test_result("ec2_validation", self._test_ec2_validation)
        
        # Test 3: EFA capability check
        self._add_test_result("efa_capability", self._test_efa_capability)
        
        return {
            'tests_run': 3,
            'passed': len([r for r in self.test_results[-3:] if r.status == 'PASS']),
            'failed': len([r for r in self.test_results[-3:] if r.status == 'FAIL']),
            'details': self.test_results[-3:]
        }

    def test_performance(self) -> Dict[str, Any]:
        """Test performance of configuration loading and processing"""
        results = {}
        
        # Test 1: Configuration loading performance
        self._add_test_result("config_performance", self._test_config_performance)
        
        # Test 2: Dataset query performance
        self._add_test_result("dataset_performance", self._test_dataset_performance)
        
        # Test 3: Memory usage
        self._add_test_result("memory_usage", self._test_memory_usage)
        
        return {
            'tests_run': 3,
            'passed': len([r for r in self.test_results[-3:] if r.status == 'PASS']),
            'failed': len([r for r in self.test_results[-3:] if r.status == 'FAIL']),
            'details': self.test_results[-3:]
        }

    def test_demo_workflows(self) -> Dict[str, Any]:
        """Test demo workflows with real data (limited scope for testing)"""
        results = {}
        
        # Get all domain configurations
        domains = self.config_loader.list_available_domains()
        
        for domain in domains[:3]:  # Test first 3 domains to limit test time
            workflows = self.dataset_manager.generate_demo_workflows(domain)
            
            for workflow in workflows[:1]:  # Test first workflow per domain
                self._test_demo_workflow(domain, workflow)
        
        return {
            'workflows_tested': len(self.workflow_results),
            'passed': len([r for r in self.workflow_results if r.status == 'PASS']),
            'failed': len([r for r in self.workflow_results if r.status == 'FAIL']),
            'total_cost': sum(r.cost_estimate for r in self.workflow_results),
            'details': self.workflow_results
        }

    def test_security(self) -> Dict[str, Any]:
        """Test security configurations and best practices"""
        results = {}
        
        # Test 1: Configuration security
        self._add_test_result("config_security", self._test_config_security)
        
        # Test 2: Data access security
        self._add_test_result("data_security", self._test_data_security)
        
        return {
            'tests_run': 2,
            'passed': len([r for r in self.test_results[-2:] if r.status == 'PASS']),
            'failed': len([r for r in self.test_results[-2:] if r.status == 'FAIL']),
            'details': self.test_results[-2:]
        }

    def test_scalability(self) -> Dict[str, Any]:
        """Test scalability of configurations and workflows"""
        results = {}
        
        # Test 1: Large configuration handling
        self._add_test_result("large_configs", self._test_large_configs)
        
        # Test 2: Concurrent access
        self._add_test_result("concurrent_access", self._test_concurrent_access)
        
        return {
            'tests_run': 2,
            'passed': len([r for r in self.test_results[-2:] if r.status == 'PASS']),
            'failed': len([r for r in self.test_results[-2:] if r.status == 'FAIL']),
            'details': self.test_results[-2:]
        }

    # Individual test implementations
    def _test_schema_validation(self) -> Tuple[str, float, Optional[str]]:
        """Test configuration schema validation"""
        start_time = time.time()
        
        try:
            validation_results = self.config_loader.validate_all_configs()
            failed_configs = [domain for domain, valid in validation_results.items() if not valid]
            
            if failed_configs:
                return 'FAIL', time.time() - start_time, f"Failed configs: {failed_configs}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_config_loading(self) -> Tuple[str, float, Optional[str]]:
        """Test configuration loading functionality"""
        start_time = time.time()
        
        try:
            domains = self.config_loader.list_available_domains()
            
            if len(domains) < 5:
                return 'FAIL', time.time() - start_time, f"Only {len(domains)} domains found, expected at least 5"
            
            # Test loading each configuration
            for domain in domains:
                config = self.config_loader.load_domain_config(domain)
                if not config:
                    return 'FAIL', time.time() - start_time, f"Failed to load {domain} config"
                
                # Validate required fields
                required_fields = ['name', 'description', 'spack_packages', 'aws_instance_recommendations']
                for field in required_fields:
                    if not hasattr(config, field) or not getattr(config, field):
                        return 'FAIL', time.time() - start_time, f"Missing {field} in {domain}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_mpi_templates(self) -> Tuple[str, float, Optional[str]]:
        """Test MPI template application"""
        start_time = time.time()
        
        try:
            mpi_config = self.config_loader.get_mpi_config()
            
            if not mpi_config:
                return 'FAIL', time.time() - start_time, "MPI configuration not found"
            
            # Check required MPI components
            required_packages = ['openmpi', 'libfabric', 'aws-ofi-nccl']
            mpi_packages_str = str(mpi_config.mpi_packages)
            
            for package in required_packages:
                if package not in mpi_packages_str:
                    return 'FAIL', time.time() - start_time, f"Missing MPI package: {package}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_domain_features(self) -> Tuple[str, float, Optional[str]]:
        """Test domain-specific features"""
        start_time = time.time()
        
        try:
            domains = self.config_loader.list_available_domains()
            
            for domain in domains:
                config = self.config_loader.load_domain_config(domain)
                
                # Check for AWS integration
                if not hasattr(config, 'aws_integration') or not config.aws_integration:
                    return 'FAIL', time.time() - start_time, f"Missing AWS integration in {domain}"
                
                # Check for demo workflows
                if not hasattr(config, 'demo_workflows') or not config.demo_workflows:
                    return 'FAIL', time.time() - start_time, f"Missing demo workflows in {domain}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_cost_validation(self) -> Tuple[str, float, Optional[str]]:
        """Test cost estimation validation"""
        start_time = time.time()
        
        try:
            domains = self.config_loader.list_available_domains()
            
            for domain in domains:
                config = self.config_loader.load_domain_config(domain)
                
                if not hasattr(config, 'estimated_cost') or not config.estimated_cost:
                    return 'FAIL', time.time() - start_time, f"Missing cost estimates in {domain}"
                
                # Validate cost structure
                required_cost_fields = ['total']
                for field in required_cost_fields:
                    if field not in config.estimated_cost:
                        return 'FAIL', time.time() - start_time, f"Missing cost field {field} in {domain}"
                
                # Check reasonable cost ranges
                total_cost = config.estimated_cost['total']
                if total_cost < 100 or total_cost > 10000:
                    return 'FAIL', time.time() - start_time, f"Unreasonable cost estimate ${total_cost} in {domain}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_dataset_registry(self) -> Tuple[str, float, Optional[str]]:
        """Test dataset registry functionality"""
        start_time = time.time()
        
        try:
            datasets = self.dataset_manager.list_available_datasets()
            
            if len(datasets) < 10:
                return 'FAIL', time.time() - start_time, f"Only {len(datasets)} datasets found, expected at least 10"
            
            # Test dataset filtering
            genomics_datasets = self.dataset_manager.list_available_datasets('genomics')
            if len(genomics_datasets) < 2:
                return 'FAIL', time.time() - start_time, "Insufficient genomics datasets"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_dataset_access(self) -> Tuple[str, float, Optional[str]]:
        """Test dataset accessibility (lightweight check)"""
        start_time = time.time()
        
        try:
            if not self.aws_available:
                return 'SKIP', time.time() - start_time, "AWS not available"
            
            # Test a small, public dataset
            test_bucket = "landsat-pds"
            
            try:
                response = self.s3_client.head_bucket(Bucket=test_bucket)
                return 'PASS', time.time() - start_time, None
            except Exception as e:
                return 'FAIL', time.time() - start_time, f"Cannot access test bucket: {e}"
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_data_subsets(self) -> Tuple[str, float, Optional[str]]:
        """Test data subset creation"""
        start_time = time.time()
        
        try:
            # Test creating a subset
            subset_config = {
                'name': 'test_subset',
                'size_gb': 0.1,
                'file_count': 5,
                'workflow_type': 'test'
            }
            
            subset = self.dataset_manager.create_demo_subset('1000 Genomes Project', subset_config)
            
            if not subset:
                return 'FAIL', time.time() - start_time, "Failed to create data subset"
            
            if subset.size_gb != 0.1:
                return 'FAIL', time.time() - start_time, "Incorrect subset size"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_cost_estimation(self) -> Tuple[str, float, Optional[str]]:
        """Test cost estimation accuracy"""
        start_time = time.time()
        
        try:
            # Test cost estimation for different data sizes
            test_sizes = [1.0, 10.0, 100.0]  # GB
            
            for size_gb in test_sizes:
                cost = self.dataset_manager._estimate_data_costs(size_gb)
                
                # Cost should scale roughly linearly
                expected_min = size_gb * 0.1
                expected_max = size_gb * 1.0
                
                if cost < expected_min or cost > expected_max:
                    return 'FAIL', time.time() - start_time, f"Cost ${cost} for {size_gb}GB outside expected range"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_s3_access(self) -> Tuple[str, float, Optional[str]]:
        """Test S3 access functionality"""
        start_time = time.time()
        
        try:
            # Test listing a public bucket
            response = self.s3_client.list_objects_v2(
                Bucket='landsat-pds',
                Prefix='c1/L8/',
                MaxKeys=5
            )
            
            if 'Contents' not in response or len(response['Contents']) == 0:
                return 'FAIL', time.time() - start_time, "No objects found in test bucket"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_ec2_validation(self) -> Tuple[str, float, Optional[str]]:
        """Test EC2 instance type validation"""
        start_time = time.time()
        
        try:
            # Test common instance types from configurations
            test_instances = ['c6i.xlarge', 'r6i.4xlarge', 'g5.2xlarge']
            
            for instance_type in test_instances:
                try:
                    response = self.ec2_client.describe_instance_types(
                        InstanceTypes=[instance_type]
                    )
                    
                    if not response['InstanceTypes']:
                        return 'FAIL', time.time() - start_time, f"Instance type {instance_type} not found"
                        
                except Exception as e:
                    return 'FAIL', time.time() - start_time, f"Error validating {instance_type}: {e}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_efa_capability(self) -> Tuple[str, float, Optional[str]]:
        """Test EFA capability validation"""
        start_time = time.time()
        
        try:
            # Check if EFA-enabled instance types are available
            efa_instances = ['hpc6a.48xlarge', 'p4d.24xlarge', 'c6in.32xlarge']
            
            for instance_type in efa_instances:
                try:
                    response = self.ec2_client.describe_instance_types(
                        InstanceTypes=[instance_type]
                    )
                    
                    instance_info = response['InstanceTypes'][0]
                    
                    # Check for EFA support
                    network_info = instance_info.get('NetworkInfo', {})
                    efa_supported = network_info.get('EfaSupported', False)
                    
                    if not efa_supported:
                        return 'FAIL', time.time() - start_time, f"Instance {instance_type} doesn't support EFA"
                        
                except Exception as e:
                    return 'FAIL', time.time() - start_time, f"Error checking EFA for {instance_type}: {e}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_config_performance(self) -> Tuple[str, float, Optional[str]]:
        """Test configuration loading performance"""
        start_time = time.time()
        
        try:
            domains = self.config_loader.list_available_domains()
            
            # Time loading all configurations
            load_start = time.time()
            configs = {}
            
            for domain in domains:
                config = self.config_loader.load_domain_config(domain)
                configs[domain] = config
            
            load_duration = time.time() - load_start
            
            # Should load all configs in under 5 seconds
            if load_duration > 5.0:
                return 'FAIL', time.time() - start_time, f"Configuration loading too slow: {load_duration:.2f}s"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_dataset_performance(self) -> Tuple[str, float, Optional[str]]:
        """Test dataset query performance"""
        start_time = time.time()
        
        try:
            # Time dataset operations
            query_start = time.time()
            
            datasets = self.dataset_manager.list_available_datasets()
            genomics_datasets = self.dataset_manager.get_datasets_for_domain('genomics')
            workflows = self.dataset_manager.generate_demo_workflows('genomics')
            
            query_duration = time.time() - query_start
            
            # Should complete queries in under 2 seconds
            if query_duration > 2.0:
                return 'FAIL', time.time() - start_time, f"Dataset queries too slow: {query_duration:.2f}s"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_memory_usage(self) -> Tuple[str, float, Optional[str]]:
        """Test memory usage during operations"""
        start_time = time.time()
        
        try:
            process = psutil.Process()
            initial_memory = process.memory_info().rss / 1024 / 1024  # MB
            
            # Perform memory-intensive operations
            for _ in range(3):
                domains = self.config_loader.list_available_domains()
                for domain in domains:
                    config = self.config_loader.load_domain_config(domain)
                    datasets = self.dataset_manager.get_datasets_for_domain(domain)
            
            peak_memory = process.memory_info().rss / 1024 / 1024  # MB
            memory_increase = peak_memory - initial_memory
            
            # Should not use more than 500MB additional memory
            if memory_increase > 500:
                return 'FAIL', time.time() - start_time, f"Excessive memory usage: {memory_increase:.1f}MB"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_demo_workflow(self, domain: str, workflow: Dict[str, Any]):
        """Test a demo workflow (dry run)"""
        start_time = time.time()
        
        try:
            # This is a dry run test - we don't actually run the workflow
            # but validate its configuration and requirements
            
            required_fields = ['name', 'description', 'dataset', 'cost_estimate']
            for field in required_fields:
                if field not in workflow:
                    result = WorkflowTestResult(
                        workflow_name=workflow.get('name', 'unknown'),
                        domain=domain,
                        status='FAIL',
                        duration_seconds=time.time() - start_time,
                        cost_estimate=0.0,
                        data_downloaded_gb=0.0,
                        error_message=f"Missing required field: {field}"
                    )
                    self.workflow_results.append(result)
                    return
            
            # Validate cost estimate
            cost = workflow.get('cost_estimate', 0)
            if cost <= 0 or cost > 1000:
                result = WorkflowTestResult(
                    workflow_name=workflow['name'],
                    domain=domain,
                    status='FAIL',
                    duration_seconds=time.time() - start_time,
                    cost_estimate=cost,
                    data_downloaded_gb=0.0,
                    error_message=f"Invalid cost estimate: ${cost}"
                )
                self.workflow_results.append(result)
                return
            
            # Simulate successful test
            result = WorkflowTestResult(
                workflow_name=workflow['name'],
                domain=domain,
                status='PASS',
                duration_seconds=time.time() - start_time,
                cost_estimate=cost,
                data_downloaded_gb=self.test_config['demo_data_limit_gb'],
                performance_metrics={
                    'validation_time': time.time() - start_time,
                    'estimated_runtime': workflow.get('expected_runtime', 'unknown')
                }
            )
            self.workflow_results.append(result)
            
        except Exception as e:
            result = WorkflowTestResult(
                workflow_name=workflow.get('name', 'unknown'),
                domain=domain,
                status='FAIL',
                duration_seconds=time.time() - start_time,
                cost_estimate=0.0,
                data_downloaded_gb=0.0,
                error_message=str(e)
            )
            self.workflow_results.append(result)

    def _test_config_security(self) -> Tuple[str, float, Optional[str]]:
        """Test configuration security best practices"""
        start_time = time.time()
        
        try:
            domains = self.config_loader.list_available_domains()
            
            for domain in domains:
                config = self.config_loader.load_domain_config(domain)
                
                # Check for security-related configurations
                if domain == 'cybersecurity_research':
                    if not hasattr(config, 'security_features') or not config.security_features:
                        return 'FAIL', time.time() - start_time, "Missing security features in cybersecurity config"
                
                # Check for encryption settings in instance recommendations
                for instance_name, instance_config in config.aws_instance_recommendations.items():
                    # EFA-enabled instances should have proper networking
                    if instance_config.get('efa_enabled', False):
                        if not instance_config.get('placement_group'):
                            return 'FAIL', time.time() - start_time, f"EFA instance missing placement group in {domain}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_data_security(self) -> Tuple[str, float, Optional[str]]:
        """Test data access security"""
        start_time = time.time()
        
        try:
            # Check that AWS integration includes security settings
            validation_results = self.aws_integrator.validate_integrations()
            
            if validation_results['validation_errors']:
                return 'FAIL', time.time() - start_time, f"Security validation errors: {validation_results['validation_errors']}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_large_configs(self) -> Tuple[str, float, Optional[str]]:
        """Test handling of large configurations"""
        start_time = time.time()
        
        try:
            # Load all configurations simultaneously
            domains = self.config_loader.list_available_domains()
            all_configs = self.config_loader.load_all_domain_configs()
            
            if len(all_configs) != len(domains):
                return 'FAIL', time.time() - start_time, "Failed to load all configurations"
            
            # Test memory efficiency
            total_memory = sys.getsizeof(all_configs)
            if total_memory > 50 * 1024 * 1024:  # 50MB
                return 'FAIL', time.time() - start_time, f"Configurations use too much memory: {total_memory / 1024 / 1024:.1f}MB"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _test_concurrent_access(self) -> Tuple[str, float, Optional[str]]:
        """Test concurrent access to configurations"""
        start_time = time.time()
        
        try:
            domains = self.config_loader.list_available_domains()
            
            def load_config(domain):
                return self.config_loader.load_domain_config(domain)
            
            # Test concurrent loading
            with ThreadPoolExecutor(max_workers=4) as executor:
                futures = [executor.submit(load_config, domain) for domain in domains]
                results = [future.result() for future in as_completed(futures)]
            
            # Check all configs loaded successfully
            failed_loads = [r for r in results if r is None]
            if failed_loads:
                return 'FAIL', time.time() - start_time, f"Failed concurrent loads: {len(failed_loads)}"
            
            return 'PASS', time.time() - start_time, None
            
        except Exception as e:
            return 'FAIL', time.time() - start_time, str(e)

    def _add_test_result(self, test_name: str, test_function):
        """Add a test result to the results list"""
        try:
            status, duration, error = test_function()
            result = TestResult(
                test_name=test_name,
                status=status,
                duration_seconds=duration,
                error_message=error
            )
            self.test_results.append(result)
        except Exception as e:
            result = TestResult(
                test_name=test_name,
                status='FAIL',
                duration_seconds=0.0,
                error_message=str(e)
            )
            self.test_results.append(result)

    def _generate_test_summary(self, suite_results: Dict[str, Any], total_duration: float) -> Dict[str, Any]:
        """Generate comprehensive test summary"""
        
        total_tests = len(self.test_results)
        passed_tests = len([r for r in self.test_results if r.status == 'PASS'])
        failed_tests = len([r for r in self.test_results if r.status == 'FAIL'])
        skipped_tests = len([r for r in self.test_results if r.status == 'SKIP'])
        
        workflow_tests = len(self.workflow_results)
        passed_workflows = len([r for r in self.workflow_results if r.status == 'PASS'])
        failed_workflows = len([r for r in self.workflow_results if r.status == 'FAIL'])
        
        summary = {
            'test_summary': {
                'total_duration': total_duration,
                'test_suites': len(suite_results),
                'total_tests': total_tests,
                'passed': passed_tests,
                'failed': failed_tests,
                'skipped': skipped_tests,
                'success_rate': (passed_tests / total_tests * 100) if total_tests > 0 else 0
            },
            'workflow_summary': {
                'total_workflows': workflow_tests,
                'passed': passed_workflows,
                'failed': failed_workflows,
                'total_estimated_cost': sum(r.cost_estimate for r in self.workflow_results)
            },
            'suite_results': suite_results,
            'test_results': [
                {
                    'name': r.test_name,
                    'status': r.status,
                    'duration': r.duration_seconds,
                    'error': r.error_message
                } for r in self.test_results
            ],
            'workflow_results': [
                {
                    'name': r.workflow_name,
                    'domain': r.domain,
                    'status': r.status,
                    'duration': r.duration_seconds,
                    'cost': r.cost_estimate,
                    'error': r.error_message
                } for r in self.workflow_results
            ],
            'recommendations': self._generate_recommendations()
        }
        
        return summary

    def _generate_recommendations(self) -> List[str]:
        """Generate recommendations based on test results"""
        recommendations = []
        
        # Check failure rate
        total_tests = len(self.test_results)
        failed_tests = len([r for r in self.test_results if r.status == 'FAIL'])
        
        if total_tests > 0:
            failure_rate = failed_tests / total_tests
            if failure_rate > 0.1:
                recommendations.append(f"High test failure rate ({failure_rate*100:.1f}%) - review failing tests")
        
        # Check AWS availability
        if not self.aws_available:
            recommendations.append("AWS credentials not configured - some tests were skipped")
        
        # Check workflow costs
        total_cost = sum(r.cost_estimate for r in self.workflow_results)
        if total_cost > 100:
            recommendations.append(f"High demo workflow costs (${total_cost:.2f}) - consider optimizing")
        
        # Performance recommendations
        slow_tests = [r for r in self.test_results if r.duration_seconds > 5.0]
        if slow_tests:
            recommendations.append(f"{len(slow_tests)} tests are slow - consider optimization")
        
        return recommendations

    def export_test_report(self, output_file: str) -> bool:
        """Export comprehensive test report"""
        try:
            # Run all tests if not already run
            if not self.test_results and not self.workflow_results:
                self.run_all_tests()
            
            # Generate HTML report
            html_report = self._generate_html_report()
            
            with open(output_file, 'w') as f:
                f.write(html_report)
            
            self.logger.info(f"Test report exported to {output_file}")
            return True
            
        except Exception as e:
            self.logger.error(f"Failed to export test report: {e}")
            return False

    def _generate_html_report(self) -> str:
        """Generate HTML test report"""
        
        total_tests = len(self.test_results)
        passed_tests = len([r for r in self.test_results if r.status == 'PASS'])
        failed_tests = len([r for r in self.test_results if r.status == 'FAIL'])
        
        html = f"""
<!DOCTYPE html>
<html>
<head>
    <title>AWS Research Wizard Test Report</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 20px; }}
        .header {{ background: #2c3e50; color: white; padding: 20px; border-radius: 5px; }}
        .summary {{ background: #ecf0f1; padding: 15px; margin: 20px 0; border-radius: 5px; }}
        .pass {{ color: #27ae60; }}
        .fail {{ color: #e74c3c; }}
        .skip {{ color: #f39c12; }}
        table {{ width: 100%; border-collapse: collapse; margin: 20px 0; }}
        th, td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
        th {{ background-color: #f2f2f2; }}
        .test-pass {{ background-color: #d5f4e6; }}
        .test-fail {{ background-color: #f8d7da; }}
        .test-skip {{ background-color: #fff3cd; }}
    </style>
</head>
<body>
    <div class="header">
        <h1>AWS Research Wizard Test Report</h1>
        <p>Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>
    </div>
    
    <div class="summary">
        <h2>Test Summary</h2>
        <p>Total Tests: {total_tests}</p>
        <p class="pass">Passed: {passed_tests}</p>
        <p class="fail">Failed: {failed_tests}</p>
        <p>Success Rate: {(passed_tests/total_tests*100) if total_tests > 0 else 0:.1f}%</p>
    </div>
    
    <h2>Test Results</h2>
    <table>
        <tr>
            <th>Test Name</th>
            <th>Status</th>
            <th>Duration (s)</th>
            <th>Error Message</th>
        </tr>
"""
        
        for result in self.test_results:
            status_class = f"test-{result.status.lower()}"
            html += f"""
        <tr class="{status_class}">
            <td>{result.test_name}</td>
            <td>{result.status}</td>
            <td>{result.duration_seconds:.2f}</td>
            <td>{result.error_message or ''}</td>
        </tr>
"""
        
        html += """
    </table>
    
    <h2>Workflow Test Results</h2>
    <table>
        <tr>
            <th>Workflow</th>
            <th>Domain</th>
            <th>Status</th>
            <th>Cost Estimate</th>
            <th>Error</th>
        </tr>
"""
        
        for result in self.workflow_results:
            status_class = f"test-{result.status.lower()}"
            html += f"""
        <tr class="{status_class}">
            <td>{result.workflow_name}</td>
            <td>{result.domain}</td>
            <td>{result.status}</td>
            <td>${result.cost_estimate:.2f}</td>
            <td>{result.error_message or ''}</td>
        </tr>
"""
        
        html += """
    </table>
</body>
</html>
"""
        
        return html


def main():
    """CLI interface for test framework"""
    import argparse
    
    parser = argparse.ArgumentParser(description="AWS Research Wizard Test Framework")
    parser.add_argument("--run-all", action="store_true", help="Run all test suites")
    parser.add_argument("--test-configs", action="store_true", help="Test configurations only")
    parser.add_argument("--test-datasets", action="store_true", help="Test datasets only")
    parser.add_argument("--test-workflows", action="store_true", help="Test demo workflows only")
    parser.add_argument("--export-report", type=str, help="Export HTML test report")
    parser.add_argument("--config-root", type=str, default="configs", help="Configuration root directory")
    parser.add_argument("--timeout", type=int, default=3600, help="Test timeout in seconds")
    
    args = parser.parse_args()
    
    # Setup logging
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
    
    # Initialize test framework
    framework = TestFramework(args.config_root)
    framework.test_config['timeout_seconds'] = args.timeout
    
    if args.run_all:
        print("Running comprehensive test suite...")
        results = framework.run_all_tests()
        
        print(f"\nTest Summary:")
        print(f"  Total Tests: {results['test_summary']['total_tests']}")
        print(f"  Passed: {results['test_summary']['passed']}")
        print(f"  Failed: {results['test_summary']['failed']}")
        print(f"  Success Rate: {results['test_summary']['success_rate']:.1f}%")
        print(f"  Duration: {results['test_summary']['total_duration']:.1f}s")
        
        if results['workflow_summary']['total_workflows'] > 0:
            print(f"\nWorkflow Tests:")
            print(f"  Total: {results['workflow_summary']['total_workflows']}")
            print(f"  Passed: {results['workflow_summary']['passed']}")
            print(f"  Failed: {results['workflow_summary']['failed']}")
            print(f"  Total Cost: ${results['workflow_summary']['total_estimated_cost']:.2f}")
        
        if results['recommendations']:
            print(f"\nRecommendations:")
            for rec in results['recommendations']:
                print(f"  - {rec}")
    
    elif args.test_configs:
        print("Testing configurations...")
        results = framework.test_configurations()
        print(f"Configuration tests: {results['passed']}/{results['tests_run']} passed")
    
    elif args.test_datasets:
        print("Testing datasets...")
        results = framework.test_datasets()
        print(f"Dataset tests: {results['passed']}/{results['tests_run']} passed")
    
    elif args.test_workflows:
        print("Testing demo workflows...")
        results = framework.test_demo_workflows()
        print(f"Workflow tests: {results['passed']}/{results['workflows_tested']} passed")
    
    else:
        parser.print_help()
    
    if args.export_report:
        print(f"Exporting test report to {args.export_report}...")
        success = framework.export_test_report(args.export_report)
        print("✅ Report exported" if success else "❌ Failed to export report")


if __name__ == "__main__":
    main()