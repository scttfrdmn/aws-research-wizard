#!/usr/bin/env python3
"""
Compute-Data Coordination Engine for AWS Research Wizard

This module provides intelligent coordination between data availability and compute resource
provisioning to optimize costs and minimize idle time while ensuring data is ready before
compute resources are allocated.

Key Features:
- Data-driven compute scheduling (don't spin up until data is ready)
- AWS cost optimization with egress waiver program integration
- Intelligent staging and pre-processing coordination
- Real-time monitoring of data transfer and compute readiness
- Cost threshold enforcement and budget management
- Multi-stage workflow coordination with dependencies

Classes:
    DataReadinessMonitor: Tracks data availability and transfer status
    ComputeScheduler: Manages compute resource lifecycle based on data events
    CostOptimizer: Handles AWS cost optimization including egress waiver management
    WorkflowCoordinator: Orchestrates multi-stage workflows with data dependencies
    ComputeDataCoordinator: Main coordination engine

The coordinator addresses real-world challenges:
- Free AWS ingress vs. paid egress (with waiver program support)
- Transfer time optimization vs. cost optimization
- Coordination overhead and timing precision
- Research-specific budget constraints and institutional policies

Dependencies:
    - boto3: For AWS service integration
    - asyncio: For asynchronous coordination
    - dataclasses: For structured configuration
"""

import os
import sys
import yaml
import json
import boto3
import asyncio
import logging
import time
from typing import Dict, List, Any, Optional, Union, Callable
from pathlib import Path
from dataclasses import dataclass, asdict, field
from datetime import datetime, timedelta
from enum import Enum
import concurrent.futures
from collections import defaultdict

# Import our core modules
from data_management_engine import DataManagementEngine, DataSource, TransferPriority
from demo_workflow_engine import DemoWorkflowEngine, ExecutionEnvironment, WorkflowExecution

class DataReadinessState(Enum):
    """States of data readiness for compute workloads."""
    UNKNOWN = "unknown"
    TRANSFERRING = "transferring"
    STAGED = "staged"
    VALIDATED = "validated"
    READY = "ready"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"

class ComputeState(Enum):
    """States of compute resource lifecycle."""
    IDLE = "idle"
    PROVISIONING = "provisioning"
    READY = "ready"
    RUNNING = "running"
    SCALING = "scaling"
    TERMINATING = "terminating"
    TERMINATED = "terminated"
    FAILED = "failed"

class CostThresholdType(Enum):
    """Types of cost thresholds for budget management."""
    DAILY = "daily"
    WEEKLY = "weekly"
    MONTHLY = "monthly"
    PER_WORKFLOW = "per_workflow"
    TOTAL_PROJECT = "total_project"

@dataclass
class DataDependency:
    """Definition of data dependency for a workflow."""
    dependency_id: str
    source_id: str
    data_path: str
    required_size_gb: Optional[float] = None
    quality_requirements: Optional[Dict[str, Any]] = None
    preprocessing_required: bool = False
    staging_location: Optional[str] = None
    readiness_timeout_hours: float = 24.0
    priority: TransferPriority = TransferPriority.STANDARD

@dataclass
class ComputeRequirement:
    """Compute resource requirements for a workflow."""
    requirement_id: str
    instance_type: str
    instance_count: int
    estimated_runtime_hours: float
    storage_gb: int
    network_requirements: Optional[Dict[str, Any]] = None
    auto_scaling: bool = False
    spot_instances: bool = False
    placement_group: Optional[str] = None
    estimated_cost_per_hour: float = 0.0

@dataclass
class CostThreshold:
    """Cost threshold configuration for budget management."""
    threshold_type: CostThresholdType
    limit_usd: float
    warning_percent: float = 80.0
    hard_limit: bool = False
    notification_emails: List[str] = field(default_factory=list)
    egress_waiver_applies: bool = True

@dataclass
class EgressWaiverStatus:
    """AWS Global Data Egress Waiver program status."""
    enabled: bool
    monthly_limit_tb: float
    current_usage_tb: float
    remaining_tb: float
    qualifying_institution: bool
    waiver_expiry_date: Optional[datetime] = None
    usage_reporting_required: bool = True
    last_reported: Optional[datetime] = None

class DataReadinessMonitor:
    """
    Monitors data availability, transfer progress, and validation status
    to determine when data is ready for compute workloads.
    """
    
    def __init__(self, data_engine: DataManagementEngine):
        self.data_engine = data_engine
        self.logger = logging.getLogger(f"{__name__}.DataReadinessMonitor")
        
        # Track data dependencies and their states
        self.data_dependencies: Dict[str, DataDependency] = {}
        self.readiness_state: Dict[str, DataReadinessState] = {}
        self.transfer_progress: Dict[str, Dict[str, Any]] = {}
        self.validation_results: Dict[str, Dict[str, Any]] = {}
        
        # Monitoring configuration
        self.check_interval_seconds = 30
        self.monitoring_active = False
        
    async def register_data_dependency(self, dependency: DataDependency) -> bool:
        """Register a data dependency for monitoring."""
        try:
            self.data_dependencies[dependency.dependency_id] = dependency
            self.readiness_state[dependency.dependency_id] = DataReadinessState.UNKNOWN
            
            # Initiate transfer if needed
            await self._initiate_dependency_transfer(dependency)
            
            self.logger.info(f"Registered data dependency: {dependency.dependency_id}")
            return True
            
        except Exception as e:
            self.logger.error(f"Failed to register data dependency {dependency.dependency_id}: {e}")
            return False
    
    async def _initiate_dependency_transfer(self, dependency: DataDependency):
        """Initiate transfer for a data dependency if not already available."""
        # Check if data is already available in staging location
        if dependency.staging_location:
            try:
                # Check if data exists in staging
                staging_available = await self._check_staging_availability(dependency)
                if staging_available:
                    self.readiness_state[dependency.dependency_id] = DataReadinessState.STAGED
                    return
            except Exception as e:
                self.logger.warning(f"Could not check staging availability: {e}")
        
        # Initiate transfer
        try:
            destination = dependency.staging_location or "s3://research-staging/"
            transfer_id = await self.data_engine.initiate_smart_transfer(
                dependency.source_id,
                dependency.data_path,
                destination,
                dependency.priority
            )
            
            self.transfer_progress[dependency.dependency_id] = {
                'transfer_id': transfer_id,
                'start_time': datetime.now(),
                'status': 'initiated'
            }
            
            self.readiness_state[dependency.dependency_id] = DataReadinessState.TRANSFERRING
            
        except Exception as e:
            self.logger.error(f"Failed to initiate transfer for {dependency.dependency_id}: {e}")
            self.readiness_state[dependency.dependency_id] = DataReadinessState.FAILED
    
    async def _check_staging_availability(self, dependency: DataDependency) -> bool:
        """Check if data is available in staging location."""
        # This would implement actual checking logic
        # For now, return False to trigger transfer
        return False
    
    async def start_monitoring(self):
        """Start continuous monitoring of data dependencies."""
        if self.monitoring_active:
            return
        
        self.monitoring_active = True
        asyncio.create_task(self._monitoring_loop())
        self.logger.info("Started data readiness monitoring")
    
    async def _monitoring_loop(self):
        """Main monitoring loop to check data readiness status."""
        while self.monitoring_active:
            try:
                await self._check_all_dependencies()
                await asyncio.sleep(self.check_interval_seconds)
            except Exception as e:
                self.logger.error(f"Error in monitoring loop: {e}")
                await asyncio.sleep(self.check_interval_seconds)
    
    async def _check_all_dependencies(self):
        """Check status of all registered data dependencies."""
        for dep_id, dependency in self.data_dependencies.items():
            current_state = self.readiness_state[dep_id]
            
            if current_state == DataReadinessState.TRANSFERRING:
                await self._check_transfer_progress(dep_id, dependency)
            elif current_state == DataReadinessState.STAGED:
                await self._validate_staged_data(dep_id, dependency)
            elif current_state == DataReadinessState.VALIDATED:
                await self._finalize_data_readiness(dep_id, dependency)
    
    async def _check_transfer_progress(self, dep_id: str, dependency: DataDependency):
        """Check progress of data transfer."""
        if dep_id not in self.transfer_progress:
            return
        
        transfer_info = self.transfer_progress[dep_id]
        transfer_id = transfer_info['transfer_id']
        
        # Get transfer status from data engine
        status = self.data_engine.transfer_manager.get_transfer_status(transfer_id)
        
        if status:
            transfer_info['status'] = status['status']
            transfer_info['progress_percent'] = status.get('progress_percent', 0)
            
            if status['status'] == 'completed':
                self.readiness_state[dep_id] = DataReadinessState.STAGED
                self.logger.info(f"Data transfer completed for {dep_id}")
            elif status['status'] == 'failed':
                self.readiness_state[dep_id] = DataReadinessState.FAILED
                self.logger.error(f"Data transfer failed for {dep_id}: {status.get('error_message')}")
    
    async def _validate_staged_data(self, dep_id: str, dependency: DataDependency):
        """Validate staged data meets requirements."""
        try:
            validation_result = {
                'timestamp': datetime.now(),
                'size_check': True,
                'quality_check': True,
                'integrity_check': True,
                'format_check': True
            }
            
            # Implement actual validation logic here
            # For now, assume validation passes
            
            self.validation_results[dep_id] = validation_result
            self.readiness_state[dep_id] = DataReadinessState.VALIDATED
            
            self.logger.info(f"Data validation completed for {dep_id}")
            
        except Exception as e:
            self.logger.error(f"Data validation failed for {dep_id}: {e}")
            self.readiness_state[dep_id] = DataReadinessState.FAILED
    
    async def _finalize_data_readiness(self, dep_id: str, dependency: DataDependency):
        """Finalize data readiness after validation."""
        # Perform any final preprocessing if required
        if dependency.preprocessing_required:
            await self._run_preprocessing(dep_id, dependency)
        else:
            self.readiness_state[dep_id] = DataReadinessState.READY
    
    async def _run_preprocessing(self, dep_id: str, dependency: DataDependency):
        """Run preprocessing steps on staged data."""
        try:
            # Implement preprocessing logic
            # For now, simulate preprocessing
            await asyncio.sleep(1)  # Simulate preprocessing time
            
            self.readiness_state[dep_id] = DataReadinessState.READY
            self.logger.info(f"Preprocessing completed for {dep_id}")
            
        except Exception as e:
            self.logger.error(f"Preprocessing failed for {dep_id}: {e}")
            self.readiness_state[dep_id] = DataReadinessState.FAILED
    
    def get_dependency_status(self, dep_id: str) -> Optional[Dict[str, Any]]:
        """Get current status of a data dependency."""
        if dep_id not in self.data_dependencies:
            return None
        
        return {
            'dependency_id': dep_id,
            'state': self.readiness_state[dep_id].value,
            'dependency': self.data_dependencies[dep_id],
            'transfer_progress': self.transfer_progress.get(dep_id),
            'validation_results': self.validation_results.get(dep_id)
        }
    
    def are_dependencies_ready(self, dep_ids: List[str]) -> bool:
        """Check if all specified dependencies are ready."""
        return all(
            self.readiness_state.get(dep_id) == DataReadinessState.READY
            for dep_id in dep_ids
        )
    
    def stop_monitoring(self):
        """Stop data readiness monitoring."""
        self.monitoring_active = False
        self.logger.info("Stopped data readiness monitoring")

class ComputeScheduler:
    """
    Manages compute resource lifecycle with intelligent scheduling
    based on data availability and cost optimization.
    """
    
    def __init__(self, workflow_engine: DemoWorkflowEngine):
        self.workflow_engine = workflow_engine
        self.logger = logging.getLogger(f"{__name__}.ComputeScheduler")
        
        # Compute resource tracking
        self.compute_resources: Dict[str, ComputeRequirement] = {}
        self.compute_state: Dict[str, ComputeState] = {}
        self.provisioning_jobs: Dict[str, Dict[str, Any]] = {}
        
        # Scheduling configuration
        self.max_idle_time_minutes = 30
        self.pre_provisioning_enabled = False  # Conservative default
        self.cost_optimization_enabled = True
        
    async def register_compute_requirement(self, requirement: ComputeRequirement) -> bool:
        """Register compute requirements for a workflow."""
        try:
            self.compute_resources[requirement.requirement_id] = requirement
            self.compute_state[requirement.requirement_id] = ComputeState.IDLE
            
            self.logger.info(f"Registered compute requirement: {requirement.requirement_id}")
            return True
            
        except Exception as e:
            self.logger.error(f"Failed to register compute requirement {requirement.requirement_id}: {e}")
            return False
    
    async def schedule_compute_for_data(self, requirement_id: str, 
                                      data_dependencies: List[str],
                                      readiness_monitor: DataReadinessMonitor) -> str:
        """Schedule compute resources based on data dependency readiness."""
        if requirement_id not in self.compute_resources:
            raise ValueError(f"Unknown compute requirement: {requirement_id}")
        
        schedule_id = f"schedule_{requirement_id}_{int(time.time())}"
        
        # Create scheduling job
        schedule_job = {
            'schedule_id': schedule_id,
            'requirement_id': requirement_id,
            'data_dependencies': data_dependencies,
            'status': 'waiting_for_data',
            'created_time': datetime.now(),
            'data_ready_time': None,
            'compute_start_time': None,
            'compute_end_time': None
        }
        
        self.provisioning_jobs[schedule_id] = schedule_job
        
        # Start monitoring for data readiness
        asyncio.create_task(self._monitor_and_provision(schedule_job, readiness_monitor))
        
        self.logger.info(f"Scheduled compute {requirement_id} waiting for data dependencies: {data_dependencies}")
        return schedule_id
    
    async def _monitor_and_provision(self, schedule_job: Dict[str, Any], 
                                   readiness_monitor: DataReadinessMonitor):
        """Monitor data readiness and provision compute when ready."""
        requirement_id = schedule_job['requirement_id']
        data_dependencies = schedule_job['data_dependencies']
        
        # Wait for data to be ready
        while not readiness_monitor.are_dependencies_ready(data_dependencies):
            await asyncio.sleep(30)  # Check every 30 seconds
            
            # Check for timeout
            elapsed = datetime.now() - schedule_job['created_time']
            if elapsed.total_seconds() > 86400:  # 24 hour timeout
                schedule_job['status'] = 'timeout'
                self.logger.error(f"Data readiness timeout for schedule {schedule_job['schedule_id']}")
                return
        
        # Data is ready - provision compute
        schedule_job['data_ready_time'] = datetime.now()
        schedule_job['status'] = 'provisioning_compute'
        
        await self._provision_compute_resources(schedule_job)
    
    async def _provision_compute_resources(self, schedule_job: Dict[str, Any]):
        """Provision compute resources for the workflow."""
        requirement_id = schedule_job['requirement_id']
        requirement = self.compute_resources[requirement_id]
        
        try:
            self.compute_state[requirement_id] = ComputeState.PROVISIONING
            schedule_job['compute_start_time'] = datetime.now()
            
            # Create execution environment
            env = ExecutionEnvironment(
                environment_type='aws_ec2',
                instance_type=requirement.instance_type,
                resource_limits={
                    'instance_count': requirement.instance_count,
                    'storage_gb': requirement.storage_gb,
                    'spot_instances': requirement.spot_instances
                }
            )
            
            # This would integrate with actual AWS provisioning
            # For now, simulate provisioning time
            provisioning_time = min(300, requirement.instance_count * 60)  # 1-5 minutes
            await asyncio.sleep(provisioning_time / 60)  # Simulate in seconds for demo
            
            self.compute_state[requirement_id] = ComputeState.READY
            schedule_job['status'] = 'compute_ready'
            
            self.logger.info(f"Compute resources ready for {requirement_id}")
            
            # Set up idle monitoring
            asyncio.create_task(self._monitor_compute_idle(requirement_id))
            
        except Exception as e:
            self.compute_state[requirement_id] = ComputeState.FAILED
            schedule_job['status'] = 'failed'
            self.logger.error(f"Failed to provision compute for {requirement_id}: {e}")
    
    async def _monitor_compute_idle(self, requirement_id: str):
        """Monitor compute resources for idle time and terminate if needed."""
        idle_start_time = None
        
        while self.compute_state.get(requirement_id) in [ComputeState.READY, ComputeState.RUNNING]:
            # Check if compute is idle
            is_idle = await self._check_compute_idle(requirement_id)
            
            if is_idle:
                if idle_start_time is None:
                    idle_start_time = datetime.now()
                else:
                    idle_duration = datetime.now() - idle_start_time
                    if idle_duration.total_seconds() > (self.max_idle_time_minutes * 60):
                        # Terminate idle compute
                        await self._terminate_compute(requirement_id, "idle_timeout")
                        break
            else:
                idle_start_time = None  # Reset idle timer
            
            await asyncio.sleep(60)  # Check every minute
    
    async def _check_compute_idle(self, requirement_id: str) -> bool:
        """Check if compute resources are idle."""
        # This would implement actual monitoring
        # For now, simulate based on state
        return self.compute_state.get(requirement_id) == ComputeState.READY
    
    async def _terminate_compute(self, requirement_id: str, reason: str):
        """Terminate compute resources."""
        self.compute_state[requirement_id] = ComputeState.TERMINATING
        
        # Implement actual termination logic
        await asyncio.sleep(5)  # Simulate termination time
        
        self.compute_state[requirement_id] = ComputeState.TERMINATED
        self.logger.info(f"Terminated compute {requirement_id} due to: {reason}")
    
    def get_compute_status(self, requirement_id: str) -> Optional[Dict[str, Any]]:
        """Get current status of compute resources."""
        if requirement_id not in self.compute_resources:
            return None
        
        return {
            'requirement_id': requirement_id,
            'state': self.compute_state[requirement_id].value,
            'requirement': self.compute_resources[requirement_id],
            'provisioning_jobs': [
                job for job in self.provisioning_jobs.values()
                if job['requirement_id'] == requirement_id
            ]
        }

class CostOptimizer:
    """
    Manages AWS cost optimization including egress waiver program
    integration and budget threshold enforcement.
    """
    
    def __init__(self):
        self.logger = logging.getLogger(f"{__name__}.CostOptimizer")
        
        # Cost tracking
        self.cost_thresholds: Dict[str, CostThreshold] = {}
        self.current_costs: Dict[str, float] = defaultdict(float)
        self.egress_waiver: Optional[EgressWaiverStatus] = None
        
        # AWS clients
        try:
            self.ce_client = boto3.client('ce')  # Cost Explorer
            self.s3_client = boto3.client('s3')
            self.aws_available = True
        except Exception:
            self.aws_available = False
            self.logger.warning("AWS cost monitoring not available")
        
        # Load egress waiver configuration
        self._load_egress_waiver_config()
    
    def _load_egress_waiver_config(self):
        """Load AWS Global Data Egress Waiver configuration."""
        # This would load from configuration file
        self.egress_waiver = EgressWaiverStatus(
            enabled=True,
            monthly_limit_tb=100.0,  # 100 TB monthly limit
            current_usage_tb=0.0,
            remaining_tb=100.0,
            qualifying_institution=True,
            usage_reporting_required=True
        )
    
    def register_cost_threshold(self, threshold_id: str, threshold: CostThreshold):
        """Register a cost threshold for monitoring."""
        self.cost_thresholds[threshold_id] = threshold
        self.logger.info(f"Registered cost threshold: {threshold_id} (${threshold.limit_usd})")
    
    async def estimate_workflow_cost(self, data_dependencies: List[DataDependency],
                                   compute_requirement: ComputeRequirement) -> Dict[str, float]:
        """Estimate total cost for a workflow including data and compute."""
        costs = {
            'data_ingress': 0.0,  # Always free to AWS
            'data_egress': 0.0,
            'data_storage': 0.0,
            'compute': 0.0,
            'total': 0.0
        }
        
        # Calculate data costs
        total_data_gb = 0
        for dep in data_dependencies:
            if dep.required_size_gb:
                total_data_gb += dep.required_size_gb
        
        # Storage costs (staging)
        costs['data_storage'] = total_data_gb * 0.023  # S3 Standard per GB/month
        
        # Egress costs (check waiver eligibility)
        if self.egress_waiver and self.egress_waiver.enabled:
            remaining_waiver_gb = self.egress_waiver.remaining_tb * 1024
            if total_data_gb <= remaining_waiver_gb:
                costs['data_egress'] = 0.0  # Covered by waiver
            else:
                excess_gb = total_data_gb - remaining_waiver_gb
                costs['data_egress'] = excess_gb * 0.09  # Standard egress rate
        else:
            costs['data_egress'] = total_data_gb * 0.09
        
        # Compute costs
        costs['compute'] = (compute_requirement.estimated_cost_per_hour * 
                           compute_requirement.estimated_runtime_hours *
                           compute_requirement.instance_count)
        
        # Apply spot instance discount if applicable
        if compute_requirement.spot_instances:
            costs['compute'] *= 0.3  # ~70% discount for spot instances
        
        costs['total'] = sum(costs.values())
        
        return costs
    
    async def check_cost_thresholds(self, estimated_cost: float, 
                                  threshold_types: List[CostThresholdType]) -> Dict[str, Any]:
        """Check if estimated cost exceeds any thresholds."""
        violations = []
        warnings = []
        
        for threshold_id, threshold in self.cost_thresholds.items():
            if threshold.threshold_type in threshold_types:
                current_period_cost = self.current_costs.get(threshold_id, 0.0)
                projected_cost = current_period_cost + estimated_cost
                
                # Check hard limit
                if threshold.hard_limit and projected_cost > threshold.limit_usd:
                    violations.append({
                        'threshold_id': threshold_id,
                        'limit': threshold.limit_usd,
                        'projected': projected_cost,
                        'type': 'hard_limit'
                    })
                
                # Check warning threshold
                warning_limit = threshold.limit_usd * (threshold.warning_percent / 100)
                if projected_cost > warning_limit:
                    warnings.append({
                        'threshold_id': threshold_id,
                        'warning_limit': warning_limit,
                        'projected': projected_cost,
                        'type': 'warning'
                    })
        
        return {
            'violations': violations,
            'warnings': warnings,
            'can_proceed': len(violations) == 0
        }
    
    async def update_egress_usage(self, egress_gb: float):
        """Update egress waiver usage tracking."""
        if self.egress_waiver and self.egress_waiver.enabled:
            egress_tb = egress_gb / 1024
            self.egress_waiver.current_usage_tb += egress_tb
            self.egress_waiver.remaining_tb = max(0, 
                self.egress_waiver.monthly_limit_tb - self.egress_waiver.current_usage_tb)
            
            # Check if approaching limit
            usage_percent = (self.egress_waiver.current_usage_tb / 
                           self.egress_waiver.monthly_limit_tb * 100)
            
            if usage_percent > 80:
                self.logger.warning(f"Egress waiver usage at {usage_percent:.1f}% of monthly limit")
    
    def get_cost_summary(self) -> Dict[str, Any]:
        """Get comprehensive cost summary and waiver status."""
        return {
            'current_costs': dict(self.current_costs),
            'cost_thresholds': {tid: asdict(threshold) for tid, threshold in self.cost_thresholds.items()},
            'egress_waiver': asdict(self.egress_waiver) if self.egress_waiver else None,
            'optimization_recommendations': self._generate_cost_recommendations()
        }
    
    def _generate_cost_recommendations(self) -> List[str]:
        """Generate cost optimization recommendations."""
        recommendations = []
        
        if self.egress_waiver and self.egress_waiver.enabled:
            usage_percent = (self.egress_waiver.current_usage_tb / 
                           self.egress_waiver.monthly_limit_tb * 100)
            if usage_percent > 50:
                recommendations.append("Consider optimizing data transfer patterns to stay within egress waiver limits")
        
        recommendations.extend([
            "Use spot instances for fault-tolerant workloads (up to 90% savings)",
            "Stage data in the same region as compute to minimize transfer costs",
            "Use S3 Intelligent Tiering for long-term storage optimization",
            "Schedule large workloads during off-peak hours for better resource availability"
        ])
        
        return recommendations

class WorkflowCoordinator:
    """
    Orchestrates complete workflows with data dependencies and compute coordination.
    """
    
    def __init__(self, data_engine: DataManagementEngine, workflow_engine: DemoWorkflowEngine):
        self.data_engine = data_engine
        self.workflow_engine = workflow_engine
        self.logger = logging.getLogger(f"{__name__}.WorkflowCoordinator")
        
        # Component initialization
        self.readiness_monitor = DataReadinessMonitor(data_engine)
        self.compute_scheduler = ComputeScheduler(workflow_engine)
        self.cost_optimizer = CostOptimizer()
        
        # Workflow tracking
        self.coordinated_workflows: Dict[str, Dict[str, Any]] = {}
    
    async def coordinate_workflow(self, workflow_config: Dict[str, Any]) -> str:
        """Coordinate a complete workflow with data dependencies and compute."""
        workflow_id = f"coordinated_{int(time.time())}"
        
        # Parse workflow configuration
        data_deps = [DataDependency(**dep) for dep in workflow_config['data_dependencies']]
        compute_req = ComputeRequirement(**workflow_config['compute_requirement'])
        cost_thresholds = workflow_config.get('cost_thresholds', [])
        
        # Create coordination record
        coordination = {
            'workflow_id': workflow_id,
            'config': workflow_config,
            'data_dependencies': data_deps,
            'compute_requirement': compute_req,
            'status': 'initializing',
            'created_time': datetime.now(),
            'data_ready_time': None,
            'compute_start_time': None,
            'workflow_start_time': None,
            'completion_time': None,
            'total_cost': 0.0
        }
        
        self.coordinated_workflows[workflow_id] = coordination
        
        try:
            # Step 1: Cost estimation and threshold checking
            estimated_costs = await self.cost_optimizer.estimate_workflow_cost(data_deps, compute_req)
            threshold_check = await self.cost_optimizer.check_cost_thresholds(
                estimated_costs['total'], 
                cost_thresholds
            )
            
            if not threshold_check['can_proceed']:
                coordination['status'] = 'cost_limit_exceeded'
                self.logger.error(f"Workflow {workflow_id} exceeds cost limits: {threshold_check['violations']}")
                return workflow_id
            
            coordination['estimated_costs'] = estimated_costs
            coordination['status'] = 'cost_approved'
            
            # Step 2: Register data dependencies and start monitoring
            await self.readiness_monitor.start_monitoring()
            
            for dep in data_deps:
                await self.readiness_monitor.register_data_dependency(dep)
            
            coordination['status'] = 'waiting_for_data'
            
            # Step 3: Register compute requirements
            await self.compute_scheduler.register_compute_requirement(compute_req)
            
            # Step 4: Schedule compute based on data readiness
            dep_ids = [dep.dependency_id for dep in data_deps]
            schedule_id = await self.compute_scheduler.schedule_compute_for_data(
                compute_req.requirement_id, dep_ids, self.readiness_monitor
            )
            
            coordination['schedule_id'] = schedule_id
            coordination['status'] = 'coordinating'
            
            # Step 5: Monitor coordination and execute workflow when ready
            asyncio.create_task(self._monitor_coordination(workflow_id))
            
            self.logger.info(f"Initiated workflow coordination: {workflow_id}")
            
        except Exception as e:
            coordination['status'] = 'failed'
            coordination['error'] = str(e)
            self.logger.error(f"Failed to coordinate workflow {workflow_id}: {e}")
        
        return workflow_id
    
    async def _monitor_coordination(self, workflow_id: str):
        """Monitor workflow coordination and execute when ready."""
        coordination = self.coordinated_workflows[workflow_id]
        
        while coordination['status'] == 'coordinating':
            # Check if both data and compute are ready
            data_deps = coordination['data_dependencies']
            compute_req = coordination['compute_requirement']
            
            dep_ids = [dep.dependency_id for dep in data_deps]
            data_ready = self.readiness_monitor.are_dependencies_ready(dep_ids)
            
            compute_status = self.compute_scheduler.get_compute_status(compute_req.requirement_id)
            compute_ready = (compute_status and 
                           compute_status['state'] == ComputeState.READY.value)
            
            if data_ready and compute_ready:
                # Both data and compute are ready - execute workflow
                await self._execute_coordinated_workflow(workflow_id)
                break
            
            await asyncio.sleep(30)  # Check every 30 seconds
    
    async def _execute_coordinated_workflow(self, workflow_id: str):
        """Execute the workflow when both data and compute are ready."""
        coordination = self.coordinated_workflows[workflow_id]
        
        try:
            coordination['status'] = 'executing'
            coordination['workflow_start_time'] = datetime.now()
            
            # Create execution environment
            compute_req = coordination['compute_requirement']
            env = ExecutionEnvironment(
                environment_type='aws_ec2',
                instance_type=compute_req.instance_type,
                resource_limits={
                    'instance_count': compute_req.instance_count,
                    'storage_gb': compute_req.storage_gb
                }
            )
            
            # Execute workflow
            workflow_config = coordination['config']
            execution_id = self.workflow_engine.execute_workflow(
                domain=workflow_config['domain'],
                workflow_name=workflow_config['workflow_name'],
                environment=env,
                dry_run=workflow_config.get('dry_run', False)
            )
            
            coordination['execution_id'] = execution_id
            coordination['status'] = 'running'
            
            # Monitor workflow execution
            await self._monitor_workflow_execution(workflow_id, execution_id)
            
        except Exception as e:
            coordination['status'] = 'execution_failed'
            coordination['error'] = str(e)
            self.logger.error(f"Failed to execute coordinated workflow {workflow_id}: {e}")
    
    async def _monitor_workflow_execution(self, workflow_id: str, execution_id: str):
        """Monitor workflow execution and handle completion."""
        coordination = self.coordinated_workflows[workflow_id]
        
        while True:
            execution_status = self.workflow_engine.get_execution_status(execution_id)
            
            if execution_status and execution_status.status in ['COMPLETED', 'FAILED', 'CANCELLED']:
                coordination['completion_time'] = datetime.now()
                coordination['status'] = f"workflow_{execution_status.status.lower()}"
                coordination['total_cost'] = execution_status.cost_actual
                
                # Update cost tracking
                await self.cost_optimizer.update_egress_usage(execution_status.data_downloaded_gb)
                
                # Clean up compute resources
                compute_req = coordination['compute_requirement']
                await self.compute_scheduler._terminate_compute(
                    compute_req.requirement_id, "workflow_completed"
                )
                
                self.logger.info(f"Workflow coordination completed: {workflow_id}")
                break
            
            await asyncio.sleep(60)  # Check every minute
    
    def get_coordination_status(self, workflow_id: str) -> Optional[Dict[str, Any]]:
        """Get status of workflow coordination."""
        if workflow_id not in self.coordinated_workflows:
            return None
        
        coordination = self.coordinated_workflows[workflow_id]
        
        # Add current status of components
        status = dict(coordination)
        
        # Data dependency status
        data_deps = coordination.get('data_dependencies', [])
        status['data_status'] = [
            self.readiness_monitor.get_dependency_status(dep.dependency_id)
            for dep in data_deps
        ]
        
        # Compute status
        compute_req = coordination.get('compute_requirement')
        if compute_req:
            status['compute_status'] = self.compute_scheduler.get_compute_status(
                compute_req.requirement_id
            )
        
        return status

class ComputeDataCoordinator:
    """
    Main coordination engine that brings together all components
    for intelligent data-compute coordination.
    """
    
    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)
        
        # Initialize core engines
        self.data_engine = DataManagementEngine(config_root)
        self.workflow_engine = DemoWorkflowEngine(config_root)
        
        # Initialize coordinator
        self.workflow_coordinator = WorkflowCoordinator(self.data_engine, self.workflow_engine)
        
        # Load configuration
        self._load_coordination_config()
    
    def _load_coordination_config(self):
        """Load coordination configuration."""
        config_file = self.config_root / "data_sources.yaml"
        
        if config_file.exists():
            with open(config_file, 'r') as f:
                config = yaml.safe_load(f)
                
                # Load global settings
                global_settings = config.get('global_settings', {})
                self._apply_global_settings(global_settings)
    
    def _apply_global_settings(self, settings: Dict[str, Any]):
        """Apply global coordination settings."""
        # Configure egress waiver
        egress_waiver_config = settings.get('egress_waiver', {})
        if egress_waiver_config.get('enabled'):
            waiver_status = EgressWaiverStatus(
                enabled=True,
                monthly_limit_tb=egress_waiver_config.get('waiver_limit_tb_monthly', 100),
                current_usage_tb=0.0,
                remaining_tb=egress_waiver_config.get('waiver_limit_tb_monthly', 100),
                qualifying_institution=egress_waiver_config.get('qualifying_institution', True)
            )
            self.workflow_coordinator.cost_optimizer.egress_waiver = waiver_status
        
        # Configure compute coordination
        compute_config = settings.get('compute_coordination', {})
        if compute_config.get('wait_for_data_availability'):
            self.workflow_coordinator.compute_scheduler.pre_provisioning_enabled = False
        
        self.logger.info("Applied global coordination settings")
    
    async def coordinate_research_workflow(self, workflow_config: Dict[str, Any]) -> str:
        """Main entry point for coordinating research workflows."""
        return await self.workflow_coordinator.coordinate_workflow(workflow_config)
    
    async def get_system_status(self) -> Dict[str, Any]:
        """Get comprehensive system status."""
        data_status = self.data_engine.get_data_management_status()
        cost_status = self.workflow_coordinator.cost_optimizer.get_cost_summary()
        
        return {
            'data_management': data_status,
            'cost_optimization': cost_status,
            'active_coordinations': len(self.workflow_coordinator.coordinated_workflows),
            'coordination_details': {
                wf_id: self.workflow_coordinator.get_coordination_status(wf_id)
                for wf_id in self.workflow_coordinator.coordinated_workflows.keys()
            }
        }

def main():
    """CLI interface for compute-data coordinator."""
    import argparse
    
    parser = argparse.ArgumentParser(description="Compute-Data Coordination Engine")
    parser.add_argument("--coordinate-workflow", type=str, help="Coordinate workflow from config file")
    parser.add_argument("--status", action="store_true", help="Show system status")
    parser.add_argument("--workflow-status", type=str, help="Show specific workflow status")
    parser.add_argument("--config-root", type=str, default="configs", help="Configuration root directory")
    
    args = parser.parse_args()
    
    # Setup logging
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
    
    # Initialize coordinator
    coordinator = ComputeDataCoordinator(args.config_root)
    
    async def run_async_commands():
        if args.coordinate_workflow:
            with open(args.coordinate_workflow, 'r') as f:
                workflow_config = yaml.safe_load(f)
            
            workflow_id = await coordinator.coordinate_research_workflow(workflow_config)
            print(f"Workflow coordination initiated: {workflow_id}")
        
        elif args.workflow_status:
            status = coordinator.workflow_coordinator.get_coordination_status(args.workflow_status)
            if status:
                print(f"Workflow {args.workflow_status} status: {status['status']}")
                if 'estimated_costs' in status:
                    print(f"Estimated costs: ${status['estimated_costs']['total']:.2f}")
            else:
                print(f"Workflow {args.workflow_status} not found")
        
        elif args.status:
            status = await coordinator.get_system_status()
            
            print("Compute-Data Coordination System Status:")
            print(f"  Data Sources: {status['data_management']['data_sources']['total_registered']}")
            print(f"  Active Transfers: {status['data_management']['transfers'].get('active_transfers', 0)}")
            print(f"  Active Coordinations: {status['active_coordinations']}")
            
            if status['cost_optimization']['egress_waiver']:
                waiver = status['cost_optimization']['egress_waiver']
                print(f"  Egress Waiver: {waiver['current_usage_tb']:.1f}/{waiver['monthly_limit_tb']} TB used")
        
        else:
            parser.print_help()
    
    # Run async commands
    asyncio.run(run_async_commands())

if __name__ == "__main__":
    main()