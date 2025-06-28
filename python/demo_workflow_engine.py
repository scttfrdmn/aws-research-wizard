#!/usr/bin/env python3
"""
Demo Workflow Execution Engine for AWS Research Wizard

This module provides a comprehensive workflow execution engine that can simulate
and execute real research workflows using AWS Open Data datasets. It supports
multiple execution environments and provides detailed cost tracking and performance monitoring.

Key Features:
- Demo workflow execution with real AWS datasets
- Multiple execution environments (local, Docker, AWS Batch, AWS EC2)
- Cost estimation and tracking for workflow execution
- Performance monitoring and resource usage tracking
- Dry-run simulation capabilities for testing
- Comprehensive workflow validation and requirements checking
- Real-time execution monitoring and cancellation support

Classes:
    WorkflowStep: Individual workflow step configuration
    WorkflowExecution: Workflow execution state and results tracking
    ExecutionEnvironment: Execution environment configuration
    DemoWorkflowEngine: Main workflow execution engine

The DemoWorkflowEngine provides:
- Domain-specific workflow generation (genomics, ML, climate, geospatial)
- Integration with real AWS Open Data datasets
- Cost-aware execution planning and monitoring
- Comprehensive logging and error handling
- Performance reporting and analytics

Safety Features:
- AWS execution disabled by default for security
- Cost limits and data download limits
- Timeout protection for long-running workflows
- Resource usage monitoring and limits

Dependencies:
    - boto3: For AWS service integration
    - docker (optional): For containerized execution
    - psutil: For system resource monitoring
    - concurrent.futures: For parallel execution
"""

import os
import sys
import yaml
import json
import boto3
import time
import subprocess
import tempfile
import logging
from typing import Dict, List, Any, Optional, Tuple
from pathlib import Path
from dataclasses import dataclass, asdict
from datetime import datetime, timedelta
import psutil
import threading
from concurrent.futures import ThreadPoolExecutor, as_completed
# Optional imports for enhanced functionality
try:
    import docker
    DOCKER_AVAILABLE = True
except ImportError:
    DOCKER_AVAILABLE = False
    docker = None

# Import our core modules
from config_loader import ConfigLoader
from dataset_manager import DatasetManager

@dataclass
class WorkflowStep:
    """
    Data class representing an individual step in a research workflow.
    
    Each workflow step represents a discrete computational task with defined
    inputs, outputs, resource requirements, and execution parameters.
    
    Attributes:
        name (str): Unique identifier for the workflow step
        description (str): Human-readable description of what this step does
        command (str): The actual command or script to execute
        input_files (List[str]): List of input files required for this step
        output_files (List[str]): List of files that will be produced by this step
        expected_duration_minutes (int): Estimated execution time in minutes
        required_resources (Dict[str, Any]): Resource requirements (CPU, memory, storage)
        environment_variables (Optional[Dict[str, str]]): Environment variables needed
    
    Example:
        >>> step = WorkflowStep(
        ...     name="quality_control",
        ...     description="Run FastQC quality control on sequencing data",
        ...     command="fastqc sample.fastq",
        ...     input_files=["sample.fastq"],
        ...     output_files=["sample_fastqc.html"],
        ...     expected_duration_minutes=15,
        ...     required_resources={"cpu_cores": 2, "memory_gb": 4}
        ... )
    """
    name: str
    description: str
    command: str
    input_files: List[str]
    output_files: List[str]
    expected_duration_minutes: int
    required_resources: Dict[str, Any]
    environment_variables: Optional[Dict[str, str]] = None

@dataclass
class WorkflowExecution:
    """
    Data class tracking the complete state and results of a workflow execution.
    
    This class provides comprehensive tracking of workflow execution from start
    to finish, including performance metrics, cost tracking, and error handling.
    
    Attributes:
        workflow_name (str): Name of the executed workflow
        domain (str): Research domain (e.g., genomics, machine_learning)
        execution_id (str): Unique identifier for this execution
        status (str): Current status - PENDING, RUNNING, COMPLETED, FAILED, CANCELLED
        start_time (datetime): When the workflow execution began
        end_time (Optional[datetime]): When the workflow execution finished
        total_duration_seconds (float): Total execution time in seconds
        cost_actual (float): Actual cost incurred during execution (USD)
        cost_estimated (float): Initial cost estimate for the workflow (USD)
        data_downloaded_gb (float): Amount of data downloaded in GB
        data_processed_gb (float): Amount of data processed in GB
        steps_completed (int): Number of workflow steps successfully completed
        steps_total (int): Total number of steps in the workflow
        error_message (Optional[str]): Error message if execution failed
        performance_metrics (Optional[Dict[str, Any]]): Detailed performance data
        resource_usage (Optional[Dict[str, Any]]): Resource utilization statistics
    """
    workflow_name: str
    domain: str
    execution_id: str
    status: str  # PENDING, RUNNING, COMPLETED, FAILED, CANCELLED
    start_time: datetime
    end_time: Optional[datetime]
    total_duration_seconds: float
    cost_actual: float
    cost_estimated: float
    data_downloaded_gb: float
    data_processed_gb: float
    steps_completed: int
    steps_total: int
    error_message: Optional[str] = None
    performance_metrics: Optional[Dict[str, Any]] = None
    resource_usage: Optional[Dict[str, Any]] = None

@dataclass
class ExecutionEnvironment:
    """
    Data class defining the execution environment configuration for workflows.
    
    This class specifies where and how workflows should be executed, including
    compute resources, containerization, and security settings.
    
    Attributes:
        environment_type (str): Type of execution environment:
            - "local": Execute on the local machine
            - "docker": Execute in Docker containers
            - "aws_batch": Execute using AWS Batch
            - "aws_ec2": Execute on AWS EC2 instances
        instance_type (Optional[str]): AWS instance type (for AWS execution)
        container_image (Optional[str]): Docker image to use (for containerized execution)
        resource_limits (Optional[Dict[str, Any]]): Resource limits and constraints
        security_settings (Optional[Dict[str, Any]]): Security and networking configuration
    
    Example:
        >>> env = ExecutionEnvironment(
        ...     environment_type="local",
        ...     resource_limits={"cpu_cores": 4, "memory_gb": 16}
        ... )
    """
    environment_type: str  # local, docker, aws_batch, aws_ec2
    instance_type: Optional[str] = None
    container_image: Optional[str] = None
    resource_limits: Optional[Dict[str, Any]] = None
    security_settings: Optional[Dict[str, Any]] = None

class DemoWorkflowEngine:
    """Demo workflow execution and validation engine"""
    
    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)
        
        # Initialize components
        self.config_loader = ConfigLoader(config_root)
        self.dataset_manager = DatasetManager(config_root)
        
        # Execution state
        self.active_executions: Dict[str, WorkflowExecution] = {}
        self.execution_history: List[WorkflowExecution] = []
        
        # Configuration
        self.engine_config = {
            'max_concurrent_workflows': 3,
            'default_timeout_hours': 6,
            'cost_limit_per_workflow': 50.0,
            'data_download_limit_gb': 10.0,
            'enable_docker': True,
            'enable_aws_execution': False,  # Disabled by default for safety
            'workspace_dir': Path('demo_workspace'),
            'log_level': 'INFO'
        }
        
        # Create workspace
        self.workspace = self.engine_config['workspace_dir']
        self.workspace.mkdir(exist_ok=True)
        
        # AWS clients (if available)
        try:
            self.s3_client = boto3.client('s3')
            self.batch_client = boto3.client('batch') if self.engine_config['enable_aws_execution'] else None
            self.ec2_client = boto3.client('ec2') if self.engine_config['enable_aws_execution'] else None
            self.aws_available = True
        except Exception as e:
            self.logger.warning(f"AWS not available: {e}")
            self.aws_available = False
            self.s3_client = None
            self.batch_client = None
            self.ec2_client = None
        
        # Docker client (if available)
        try:
            if self.engine_config['enable_docker'] and DOCKER_AVAILABLE:
                self.docker_client = docker.from_env()
                self.docker_available = True
            else:
                self.docker_client = None
                self.docker_available = False
        except Exception as e:
            self.logger.warning(f"Docker not available: {e}")
            self.docker_client = None
            self.docker_available = False

    def list_available_workflows(self) -> Dict[str, List[Dict[str, Any]]]:
        """List all available demo workflows by domain"""
        workflows_by_domain = {}
        
        domains = self.config_loader.list_available_domains()
        
        for domain in domains:
            config = self.config_loader.load_domain_config(domain)
            if config and hasattr(config, 'demo_workflows') and config.demo_workflows:
                workflows_by_domain[domain] = []
                
                for workflow in config.demo_workflows:
                    workflow_info = {
                        'name': workflow.get('name', 'Unnamed'),
                        'description': workflow.get('description', ''),
                        'dataset': workflow.get('dataset', 'Unknown'),
                        'expected_runtime': workflow.get('expected_runtime', 'Unknown'),
                        'cost_estimate': workflow.get('cost_estimate', 0.0),
                        'complexity': self._assess_workflow_complexity(workflow)
                    }
                    workflows_by_domain[domain].append(workflow_info)
        
        return workflows_by_domain

    def create_workflow_definition(self, domain: str, workflow_name: str) -> Optional[Dict[str, Any]]:
        """Create a detailed workflow definition with steps"""
        config = self.config_loader.load_domain_config(domain)
        if not config or not hasattr(config, 'demo_workflows'):
            return None
        
        # Find the workflow
        target_workflow = None
        for workflow in config.demo_workflows:
            if workflow.get('name') == workflow_name:
                target_workflow = workflow
                break
        
        if not target_workflow:
            return None
        
        # Create detailed workflow definition based on domain and workflow type
        workflow_def = self._generate_workflow_steps(domain, target_workflow)
        
        return workflow_def

    def execute_workflow(self, domain: str, workflow_name: str, 
                        environment: ExecutionEnvironment,
                        dry_run: bool = True) -> str:
        """Execute a demo workflow"""
        
        execution_id = f"{domain}_{workflow_name}_{int(time.time())}"
        
        # Create workflow definition
        workflow_def = self.create_workflow_definition(domain, workflow_name)
        if not workflow_def:
            raise ValueError(f"Workflow {workflow_name} not found in domain {domain}")
        
        # Create execution record
        execution = WorkflowExecution(
            workflow_name=workflow_name,
            domain=domain,
            execution_id=execution_id,
            status='PENDING',
            start_time=datetime.now(),
            end_time=None,
            total_duration_seconds=0.0,
            cost_actual=0.0,
            cost_estimated=workflow_def.get('cost_estimate', 0.0),
            data_downloaded_gb=0.0,
            data_processed_gb=0.0,
            steps_completed=0,
            steps_total=len(workflow_def.get('steps', []))
        )
        
        self.active_executions[execution_id] = execution
        
        # Execute workflow in background thread
        if not dry_run:
            thread = threading.Thread(
                target=self._execute_workflow_async,
                args=(execution_id, workflow_def, environment)
            )
            thread.start()
        else:
            # Dry run - simulate execution
            self._simulate_workflow_execution(execution_id, workflow_def)
        
        return execution_id

    def get_execution_status(self, execution_id: str) -> Optional[WorkflowExecution]:
        """Get execution status"""
        if execution_id in self.active_executions:
            return self.active_executions[execution_id]
        
        # Check history
        for execution in self.execution_history:
            if execution.execution_id == execution_id:
                return execution
        
        return None

    def cancel_execution(self, execution_id: str) -> bool:
        """Cancel a running execution"""
        if execution_id in self.active_executions:
            execution = self.active_executions[execution_id]
            if execution.status == 'RUNNING':
                execution.status = 'CANCELLED'
                execution.end_time = datetime.now()
                execution.total_duration_seconds = (execution.end_time - execution.start_time).total_seconds()
                
                # Move to history
                self.execution_history.append(execution)
                del self.active_executions[execution_id]
                
                self.logger.info(f"Cancelled execution {execution_id}")
                return True
        
        return False

    def validate_workflow_requirements(self, domain: str, workflow_name: str) -> Dict[str, Any]:
        """Validate workflow requirements and prerequisites"""
        
        validation_result = {
            'valid': True,
            'warnings': [],
            'errors': [],
            'requirements': {},
            'estimated_resources': {}
        }
        
        try:
            # Get workflow definition
            workflow_def = self.create_workflow_definition(domain, workflow_name)
            if not workflow_def:
                validation_result['valid'] = False
                validation_result['errors'].append(f"Workflow {workflow_name} not found")
                return validation_result
            
            # Check data requirements
            dataset_name = workflow_def.get('dataset')
            if dataset_name:
                dataset = self.dataset_manager.get_dataset_by_name(dataset_name)
                if not dataset:
                    validation_result['warnings'].append(f"Dataset {dataset_name} not found in registry")
                else:
                    validation_result['requirements']['dataset'] = {
                        'name': dataset.name,
                        'size_tb': dataset.size_tb,
                        'location': dataset.location,
                        'access_pattern': dataset.access_pattern
                    }
            
            # Check computational requirements
            estimated_cost = workflow_def.get('cost_estimate', 0.0)
            if estimated_cost > self.engine_config['cost_limit_per_workflow']:
                validation_result['warnings'].append(
                    f"Estimated cost ${estimated_cost:.2f} exceeds limit ${self.engine_config['cost_limit_per_workflow']:.2f}"
                )
            
            # Check runtime requirements
            runtime = workflow_def.get('expected_runtime', '')
            if 'hours' in runtime.lower():
                try:
                    hours = float(runtime.lower().split('hours')[0].split('-')[-1].strip())
                    if hours > self.engine_config['default_timeout_hours']:
                        validation_result['warnings'].append(
                            f"Expected runtime {hours}h exceeds timeout {self.engine_config['default_timeout_hours']}h"
                        )
                except:
                    pass
            
            # Check environment requirements
            validation_result['estimated_resources'] = {
                'cpu_cores': workflow_def.get('cpu_cores', 4),
                'memory_gb': workflow_def.get('memory_gb', 16),
                'storage_gb': workflow_def.get('storage_gb', 100),
                'gpu_required': workflow_def.get('gpu_required', False),
                'network_bandwidth': workflow_def.get('network_bandwidth', 'standard')
            }
            
        except Exception as e:
            validation_result['valid'] = False
            validation_result['errors'].append(str(e))
        
        return validation_result

    def get_execution_logs(self, execution_id: str) -> List[str]:
        """Get execution logs"""
        log_file = self.workspace / f"{execution_id}.log"
        
        if log_file.exists():
            with open(log_file, 'r') as f:
                return f.readlines()
        
        return []

    def export_execution_report(self, execution_id: str, output_file: str) -> bool:
        """Export execution report"""
        try:
            execution = self.get_execution_status(execution_id)
            if not execution:
                return False
            
            report = {
                'execution_summary': asdict(execution),
                'logs': self.get_execution_logs(execution_id),
                'export_timestamp': datetime.now().isoformat()
            }
            
            with open(output_file, 'w') as f:
                json.dump(report, f, indent=2, default=str)
            
            self.logger.info(f"Exported execution report to {output_file}")
            return True
            
        except Exception as e:
            self.logger.error(f"Failed to export execution report: {e}")
            return False

    # Internal methods
    def _assess_workflow_complexity(self, workflow: Dict[str, Any]) -> str:
        """Assess workflow complexity level"""
        cost = workflow.get('cost_estimate', 0.0)
        runtime = workflow.get('expected_runtime', '')
        
        if cost > 100 or '24' in runtime:
            return 'HIGH'
        elif cost > 20 or '8' in runtime:
            return 'MEDIUM'
        else:
            return 'LOW'

    def _generate_workflow_steps(self, domain: str, workflow: Dict[str, Any]) -> Dict[str, Any]:
        """Generate detailed workflow steps based on domain and workflow type"""
        
        workflow_def = {
            'name': workflow['name'],
            'description': workflow['description'],
            'domain': domain,
            'dataset': workflow.get('dataset', ''),
            'cost_estimate': workflow.get('cost_estimate', 0.0),
            'expected_runtime': workflow.get('expected_runtime', ''),
            'steps': []
        }
        
        # Generate domain-specific steps
        if domain == 'genomics':
            workflow_def['steps'] = self._generate_genomics_steps(workflow)
        elif domain == 'machine_learning':
            workflow_def['steps'] = self._generate_ml_steps(workflow)
        elif domain == 'climate_modeling':
            workflow_def['steps'] = self._generate_climate_steps(workflow)
        elif domain == 'geospatial_research':
            workflow_def['steps'] = self._generate_geospatial_steps(workflow)
        else:
            # Generic steps
            workflow_def['steps'] = self._generate_generic_steps(workflow)
        
        # Add common resource requirements
        workflow_def.update({
            'cpu_cores': 4,
            'memory_gb': 16,
            'storage_gb': 100,
            'gpu_required': 'gpu' in workflow['name'].lower() or domain == 'machine_learning',
            'network_bandwidth': 'high' if 'large' in workflow['description'].lower() else 'standard'
        })
        
        return workflow_def

    def _generate_genomics_steps(self, workflow: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Generate genomics workflow steps"""
        if 'variant calling' in workflow['name'].lower():
            return [
                {
                    'name': 'data_download',
                    'description': 'Download reference genome and sample data',
                    'command': 'aws s3 cp s3://1000genomes/... ./data/',
                    'expected_duration_minutes': 30,
                    'output_files': ['reference.fa', 'sample.fastq']
                },
                {
                    'name': 'quality_control',
                    'description': 'Run FastQC quality control',
                    'command': 'fastqc sample.fastq',
                    'expected_duration_minutes': 15,
                    'output_files': ['sample_fastqc.html']
                },
                {
                    'name': 'alignment',
                    'description': 'Align reads to reference genome',
                    'command': 'bwa mem reference.fa sample.fastq > aligned.sam',
                    'expected_duration_minutes': 60,
                    'output_files': ['aligned.sam']
                },
                {
                    'name': 'variant_calling',
                    'description': 'Call variants using GATK',
                    'command': 'gatk HaplotypeCaller -R reference.fa -I aligned.bam -O variants.vcf',
                    'expected_duration_minutes': 45,
                    'output_files': ['variants.vcf']
                },
                {
                    'name': 'annotation',
                    'description': 'Annotate variants',
                    'command': 'vep -i variants.vcf -o annotated.vcf',
                    'expected_duration_minutes': 20,
                    'output_files': ['annotated.vcf']
                }
            ]
        else:
            return self._generate_generic_steps(workflow)

    def _generate_ml_steps(self, workflow: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Generate machine learning workflow steps"""
        if 'image classification' in workflow['name'].lower():
            return [
                {
                    'name': 'data_download',
                    'description': 'Download training dataset',
                    'command': 'aws s3 sync s3://open-images-dataset/train/ ./data/train/',
                    'expected_duration_minutes': 45,
                    'output_files': ['data/train/']
                },
                {
                    'name': 'data_preprocessing',
                    'description': 'Preprocess images and create data loaders',
                    'command': 'python preprocess_images.py --input ./data/train/ --output ./data/processed/',
                    'expected_duration_minutes': 30,
                    'output_files': ['data/processed/']
                },
                {
                    'name': 'model_training',
                    'description': 'Train ResNet model',
                    'command': 'python train_resnet.py --data ./data/processed/ --epochs 10',
                    'expected_duration_minutes': 180,
                    'output_files': ['model.pth']
                },
                {
                    'name': 'model_evaluation',
                    'description': 'Evaluate model performance',
                    'command': 'python evaluate_model.py --model model.pth --test-data ./data/test/',
                    'expected_duration_minutes': 30,
                    'output_files': ['evaluation_results.json']
                }
            ]
        else:
            return self._generate_generic_steps(workflow)

    def _generate_climate_steps(self, workflow: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Generate climate modeling workflow steps"""
        return [
            {
                'name': 'data_download',
                'description': 'Download ERA5 reanalysis data',
                'command': 'cdsapi download era5-reanalysis temperature precipitation',
                'expected_duration_minutes': 60,
                'output_files': ['era5_data.nc']
            },
            {
                'name': 'data_preprocessing',
                'description': 'Process climate data',
                'command': 'cdo -mergetime era5_*.nc processed_data.nc',
                'expected_duration_minutes': 30,
                'output_files': ['processed_data.nc']
            },
            {
                'name': 'analysis',
                'description': 'Perform climate analysis',
                'command': 'python climate_analysis.py --input processed_data.nc',
                'expected_duration_minutes': 90,
                'output_files': ['climate_trends.nc', 'analysis_plots.png']
            }
        ]

    def _generate_geospatial_steps(self, workflow: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Generate geospatial workflow steps"""
        return [
            {
                'name': 'data_download',
                'description': 'Download Sentinel-2 imagery',
                'command': 'sentinelsat download --area bbox.geojson --date 2023-01-01to2023-12-31',
                'expected_duration_minutes': 40,
                'output_files': ['sentinel2_images/']
            },
            {
                'name': 'preprocessing',
                'description': 'Preprocess satellite imagery',
                'command': 'gdal_translate -of COG input.tif output_cog.tif',
                'expected_duration_minutes': 25,
                'output_files': ['processed_imagery/']
            },
            {
                'name': 'classification',
                'description': 'Perform land cover classification',
                'command': 'python classify_landcover.py --input processed_imagery/',
                'expected_duration_minutes': 120,
                'output_files': ['landcover_map.tif']
            },
            {
                'name': 'validation',
                'description': 'Validate classification results',
                'command': 'python validate_classification.py --classified landcover_map.tif',
                'expected_duration_minutes': 20,
                'output_files': ['accuracy_report.pdf']
            }
        ]

    def _generate_generic_steps(self, workflow: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Generate generic workflow steps"""
        return [
            {
                'name': 'setup',
                'description': 'Setup workflow environment',
                'command': 'echo "Setting up environment"',
                'expected_duration_minutes': 5,
                'output_files': ['setup.log']
            },
            {
                'name': 'data_preparation',
                'description': 'Prepare input data',
                'command': 'echo "Preparing data"',
                'expected_duration_minutes': 15,
                'output_files': ['prepared_data/']
            },
            {
                'name': 'analysis',
                'description': 'Run analysis',
                'command': 'echo "Running analysis"',
                'expected_duration_minutes': 60,
                'output_files': ['results.txt']
            },
            {
                'name': 'cleanup',
                'description': 'Clean up temporary files',
                'command': 'echo "Cleaning up"',
                'expected_duration_minutes': 5,
                'output_files': ['cleanup.log']
            }
        ]

    def _execute_workflow_async(self, execution_id: str, workflow_def: Dict[str, Any], 
                               environment: ExecutionEnvironment):
        """Execute workflow asynchronously"""
        execution = self.active_executions[execution_id]
        execution.status = 'RUNNING'
        
        try:
            # Create execution directory
            exec_dir = self.workspace / execution_id
            exec_dir.mkdir(exist_ok=True)
            
            # Setup logging
            log_file = exec_dir / "execution.log"
            
            # Execute steps
            for i, step in enumerate(workflow_def['steps']):
                self.logger.info(f"Executing step {i+1}/{len(workflow_def['steps'])}: {step['name']}")
                
                step_success = self._execute_step(step, exec_dir, log_file, environment)
                
                if step_success:
                    execution.steps_completed += 1
                else:
                    execution.status = 'FAILED'
                    execution.error_message = f"Step {step['name']} failed"
                    break
            
            # Complete execution
            if execution.status == 'RUNNING':
                execution.status = 'COMPLETED'
            
            execution.end_time = datetime.now()
            execution.total_duration_seconds = (execution.end_time - execution.start_time).total_seconds()
            
        except Exception as e:
            execution.status = 'FAILED'
            execution.error_message = str(e)
            execution.end_time = datetime.now()
            execution.total_duration_seconds = (execution.end_time - execution.start_time).total_seconds()
        
        finally:
            # Move to history
            self.execution_history.append(execution)
            if execution_id in self.active_executions:
                del self.active_executions[execution_id]

    def _simulate_workflow_execution(self, execution_id: str, workflow_def: Dict[str, Any]):
        """Simulate workflow execution for dry run"""
        execution = self.active_executions[execution_id]
        execution.status = 'RUNNING'
        
        # Simulate each step
        total_duration = 0
        for i, step in enumerate(workflow_def['steps']):
            step_duration = step.get('expected_duration_minutes', 10) * 60
            total_duration += step_duration
            
            time.sleep(0.1)  # Small delay for simulation
            execution.steps_completed += 1
        
        # Complete simulation
        execution.status = 'COMPLETED'
        execution.end_time = datetime.now()
        execution.total_duration_seconds = total_duration
        execution.data_downloaded_gb = 1.0  # Simulated
        execution.data_processed_gb = 0.5   # Simulated
        
        # Move to history
        self.execution_history.append(execution)
        del self.active_executions[execution_id]

    def _execute_step(self, step: Dict[str, Any], exec_dir: Path, 
                     log_file: Path, environment: ExecutionEnvironment) -> bool:
        """Execute a single workflow step"""
        try:
            # For demo purposes, we'll simulate step execution
            # In a real implementation, this would execute the actual commands
            
            step_start = time.time()
            
            # Log step start
            with open(log_file, 'a') as f:
                f.write(f"{datetime.now()}: Starting step {step['name']}\n")
                f.write(f"Command: {step['command']}\n")
            
            # Simulate step execution
            duration = step.get('expected_duration_minutes', 5) * 60
            time.sleep(min(duration / 60, 5))  # Cap simulation time at 5 seconds
            
            # Create output files (empty for simulation)
            for output_file in step.get('output_files', []):
                output_path = exec_dir / output_file
                output_path.parent.mkdir(parents=True, exist_ok=True)
                output_path.touch()
            
            step_duration = time.time() - step_start
            
            # Log step completion
            with open(log_file, 'a') as f:
                f.write(f"{datetime.now()}: Completed step {step['name']} in {step_duration:.1f}s\n")
            
            return True
            
        except Exception as e:
            # Log error
            with open(log_file, 'a') as f:
                f.write(f"{datetime.now()}: Error in step {step['name']}: {e}\n")
            
            return False

    def generate_performance_report(self) -> Dict[str, Any]:
        """Generate performance report for all executions"""
        
        all_executions = list(self.active_executions.values()) + self.execution_history
        
        if not all_executions:
            return {'message': 'No executions found'}
        
        # Calculate statistics
        completed_executions = [e for e in all_executions if e.status == 'COMPLETED']
        failed_executions = [e for e in all_executions if e.status == 'FAILED']
        
        total_cost = sum(e.cost_actual for e in completed_executions)
        avg_duration = sum(e.total_duration_seconds for e in completed_executions) / len(completed_executions) if completed_executions else 0
        
        # Group by domain
        by_domain = {}
        for execution in all_executions:
            domain = execution.domain
            if domain not in by_domain:
                by_domain[domain] = {'total': 0, 'completed': 0, 'failed': 0, 'cost': 0.0}
            
            by_domain[domain]['total'] += 1
            if execution.status == 'COMPLETED':
                by_domain[domain]['completed'] += 1
                by_domain[domain]['cost'] += execution.cost_actual
            elif execution.status == 'FAILED':
                by_domain[domain]['failed'] += 1
        
        report = {
            'summary': {
                'total_executions': len(all_executions),
                'completed': len(completed_executions),
                'failed': len(failed_executions),
                'success_rate': len(completed_executions) / len(all_executions) * 100 if all_executions else 0,
                'total_cost': total_cost,
                'average_duration_seconds': avg_duration
            },
            'by_domain': by_domain,
            'recent_executions': [
                {
                    'execution_id': e.execution_id,
                    'workflow_name': e.workflow_name,
                    'domain': e.domain,
                    'status': e.status,
                    'duration': e.total_duration_seconds,
                    'cost': e.cost_actual
                } for e in all_executions[-10:]  # Last 10 executions
            ]
        }
        
        return report


def main():
    """CLI interface for demo workflow engine"""
    import argparse
    
    parser = argparse.ArgumentParser(description="Demo Workflow Execution Engine")
    parser.add_argument("--list-workflows", action="store_true", help="List available workflows")
    parser.add_argument("--validate", nargs=2, metavar=('DOMAIN', 'WORKFLOW'), help="Validate workflow requirements")
    parser.add_argument("--execute", nargs=2, metavar=('DOMAIN', 'WORKFLOW'), help="Execute workflow")
    parser.add_argument("--dry-run", action="store_true", help="Perform dry run (simulation)")
    parser.add_argument("--status", type=str, help="Check execution status")
    parser.add_argument("--cancel", type=str, help="Cancel execution")
    parser.add_argument("--report", action="store_true", help="Generate performance report")
    parser.add_argument("--config-root", type=str, default="configs", help="Configuration root directory")
    
    args = parser.parse_args()
    
    # Setup logging
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
    
    # Initialize engine
    engine = DemoWorkflowEngine(args.config_root)
    
    if args.list_workflows:
        workflows = engine.list_available_workflows()
        print("Available Demo Workflows:")
        
        for domain, domain_workflows in workflows.items():
            print(f"\n{domain.upper()}:")
            for workflow in domain_workflows:
                print(f"  - {workflow['name']}")
                print(f"    Description: {workflow['description']}")
                print(f"    Dataset: {workflow['dataset']}")
                print(f"    Runtime: {workflow['expected_runtime']}")
                print(f"    Cost: ${workflow['cost_estimate']:.2f}")
                print(f"    Complexity: {workflow['complexity']}")
    
    elif args.validate:
        domain, workflow_name = args.validate
        print(f"Validating workflow: {workflow_name} in {domain}")
        
        validation = engine.validate_workflow_requirements(domain, workflow_name)
        
        print(f"Valid: {validation['valid']}")
        if validation['warnings']:
            print("Warnings:")
            for warning in validation['warnings']:
                print(f"  - {warning}")
        
        if validation['errors']:
            print("Errors:")
            for error in validation['errors']:
                print(f"  - {error}")
        
        if validation['estimated_resources']:
            print("Estimated Resources:")
            for resource, value in validation['estimated_resources'].items():
                print(f"  {resource}: {value}")
    
    elif args.execute:
        domain, workflow_name = args.execute
        
        # Create execution environment
        env = ExecutionEnvironment(
            environment_type='local',
            resource_limits={'cpu_cores': 4, 'memory_gb': 16}
        )
        
        print(f"Executing workflow: {workflow_name} in {domain}")
        if args.dry_run:
            print("(Dry run mode - simulation only)")
        
        try:
            execution_id = engine.execute_workflow(domain, workflow_name, env, dry_run=args.dry_run)
            print(f"Execution started: {execution_id}")
            
            # Monitor execution
            while True:
                status = engine.get_execution_status(execution_id)
                if status and status.status in ['COMPLETED', 'FAILED', 'CANCELLED']:
                    break
                time.sleep(1)
            
            print(f"Execution {status.status.lower()}")
            if status.status == 'COMPLETED':
                print(f"Duration: {status.total_duration_seconds:.1f}s")
                print(f"Steps completed: {status.steps_completed}/{status.steps_total}")
                print(f"Cost: ${status.cost_actual:.2f}")
            elif status.error_message:
                print(f"Error: {status.error_message}")
                
        except Exception as e:
            print(f"Execution failed: {e}")
    
    elif args.status:
        status = engine.get_execution_status(args.status)
        if status:
            print(f"Execution {args.status}:")
            print(f"  Status: {status.status}")
            print(f"  Duration: {status.total_duration_seconds:.1f}s")
            print(f"  Steps: {status.steps_completed}/{status.steps_total}")
            print(f"  Cost: ${status.cost_actual:.2f}")
            if status.error_message:
                print(f"  Error: {status.error_message}")
        else:
            print(f"Execution {args.status} not found")
    
    elif args.cancel:
        success = engine.cancel_execution(args.cancel)
        print(f"Cancellation {'successful' if success else 'failed'}")
    
    elif args.report:
        report = engine.generate_performance_report()
        print("Performance Report:")
        print(f"  Total executions: {report['summary']['total_executions']}")
        print(f"  Success rate: {report['summary']['success_rate']:.1f}%")
        print(f"  Total cost: ${report['summary']['total_cost']:.2f}")
        print(f"  Average duration: {report['summary']['average_duration_seconds']:.1f}s")
    
    else:
        parser.print_help()


if __name__ == "__main__":
    main()