#!/usr/bin/env python3
"""
S3 Transfer Optimizer for AWS Research Wizard

This module provides intelligent S3 transfer capabilities using high-performance tools
including s5cmd, rclone, and optimized AWS CLI configurations. It handles:

1. Tool Selection: Automatically selects the best transfer tool based on workload
2. Storage Tier Intelligence: Optimizes S3 storage class based on access patterns
3. Performance Optimization: Parallel transfers, multipart uploads, compression
4. Cost Optimization: Intelligent tiering, transfer acceleration decisions

Key Features:
- s5cmd integration for high-performance parallel transfers (32x faster than s3cmd)
- rclone support for multi-cloud scenarios and complex sync operations
- Intelligent storage class selection (Standard, IA, Glacier, Deep Archive)
- Automated transfer acceleration and multipart upload optimization
- Research-specific optimizations for large datasets and archival patterns

Performance Benchmarks:
- s5cmd: Up to 4.3 GB/s download on 40Gbps links
- Parallel operations: Thousands of S3 operations per second
- Cost savings: Up to 85% with intelligent tiering for archival data

Classes:
    TransferTool: Enum of available transfer tools
    StorageClass: S3 storage class options with cost/performance characteristics
    S3TransferOptimizer: Main optimization and transfer coordination class
    TransferStrategy: Data class defining transfer approach and settings
    
Dependencies:
    - s5cmd: High-performance parallel S3 tool
    - rclone: Multi-cloud sync and transfer tool
    - boto3: AWS SDK for Python
    - psutil: System performance monitoring
"""

import os
import sys
import json
import boto3
import asyncio
import logging
import subprocess
import time
from typing import Dict, List, Any, Optional, Union
from pathlib import Path
from dataclasses import dataclass, asdict
from datetime import datetime, timedelta
from enum import Enum
import psutil
import shutil

class TransferTool(Enum):
    """Available S3 transfer tools with performance characteristics."""
    S5CMD = "s5cmd"              # 32x faster than s3cmd, best for bulk operations
    RCLONE = "rclone"            # Best for multi-cloud, complex sync scenarios
    AWS_CLI = "aws_cli"          # Good baseline, universal compatibility
    AWS_CLI_OPTIMIZED = "aws_cli_optimized"  # Tuned AWS CLI with parallel settings

class StorageClass(Enum):
    """S3 storage classes with cost and access characteristics."""
    STANDARD = "STANDARD"                    # $0.023/GB/month, immediate access
    INTELLIGENT_TIERING = "INTELLIGENT_TIERING"  # $0.0125/GB/month, automatic optimization
    STANDARD_IA = "STANDARD_IA"             # $0.0125/GB/month, 30-day minimum
    ONEZONE_IA = "ONEZONE_IA"               # $0.01/GB/month, single AZ
    GLACIER = "GLACIER"                      # $0.004/GB/month, 1-5 minute retrieval
    DEEP_ARCHIVE = "DEEP_ARCHIVE"            # $0.00099/GB/month, 12-hour retrieval

@dataclass
class TransferStrategy:
    """Configuration for optimized S3 transfer strategy."""
    tool: TransferTool
    storage_class: StorageClass
    parallel_workers: int
    multipart_threshold_mb: int
    chunk_size_mb: int
    enable_compression: bool
    use_transfer_acceleration: bool
    estimated_cost_per_gb: float
    estimated_time_hours: float
    optimization_notes: List[str]

class S3TransferOptimizer:
    """
    Intelligent S3 transfer optimization and coordination system.
    
    This class analyzes transfer requirements and selects optimal tools, settings,
    and storage classes based on data characteristics, access patterns, and cost constraints.
    """
    
    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)
        
        # AWS clients
        self.s3_client = boto3.client('s3')
        self.s3_resource = boto3.resource('s3')
        
        # Tool availability detection
        self.available_tools = self._detect_available_tools()
        
        # Performance profiles for different tools
        self.tool_profiles = {
            TransferTool.S5CMD: {
                'max_throughput_gbps': 40,  # Can saturate 40Gbps links
                'parallel_efficiency': 0.95,
                'best_for': ['bulk_upload', 'bulk_download', 'many_small_files'],
                'overhead_seconds': 1
            },
            TransferTool.RCLONE: {
                'max_throughput_gbps': 10,
                'parallel_efficiency': 0.85,
                'best_for': ['multi_cloud', 'complex_sync', 'encryption'],
                'overhead_seconds': 3
            },
            TransferTool.AWS_CLI: {
                'max_throughput_gbps': 5,
                'parallel_efficiency': 0.70,
                'best_for': ['simple_operations', 'single_files', 'scripting'],
                'overhead_seconds': 2
            }
        }
        
        # Research-specific access patterns
        self.research_access_patterns = {
            'genomics': {
                'initial_access_intensive': True,
                'long_term_archive_likely': True,
                'recommended_tier': StorageClass.INTELLIGENT_TIERING
            },
            'climate_modeling': {
                'periodic_batch_access': True,
                'archive_simulation_outputs': True,
                'recommended_tier': StorageClass.STANDARD_IA
            },
            'machine_learning': {
                'training_phase_intensive': True,
                'model_archive_cold': True,
                'recommended_tier': StorageClass.INTELLIGENT_TIERING
            }
        }
    
    def _detect_available_tools(self) -> Dict[TransferTool, bool]:
        """Detect which transfer tools are available on the system."""
        tools = {}
        
        # Check for s5cmd
        tools[TransferTool.S5CMD] = shutil.which('s5cmd') is not None
        
        # Check for rclone
        tools[TransferTool.RCLONE] = shutil.which('rclone') is not None
        
        # Check for AWS CLI
        tools[TransferTool.AWS_CLI] = shutil.which('aws') is not None
        tools[TransferTool.AWS_CLI_OPTIMIZED] = tools[TransferTool.AWS_CLI]
        
        self.logger.info(f"Available transfer tools: {[tool.value for tool, available in tools.items() if available]}")
        return tools
    
    def analyze_transfer_requirements(self, source_path: str, destination: str, 
                                    data_size_gb: float, file_count: int,
                                    research_domain: str = None,
                                    access_pattern: str = "unknown") -> TransferStrategy:
        """
        Analyze transfer requirements and recommend optimal strategy.
        
        Args:
            source_path: Source path (local, S3, or other cloud)
            destination: Destination S3 path
            data_size_gb: Total data size in gigabytes
            file_count: Number of files to transfer
            research_domain: Research domain for access pattern optimization
            access_pattern: Expected access pattern (frequent, infrequent, archive)
        
        Returns:
            TransferStrategy: Optimized transfer strategy
        """
        
        # Analyze workload characteristics
        avg_file_size_mb = (data_size_gb * 1024) / max(file_count, 1)
        is_bulk_operation = file_count > 1000
        is_large_dataset = data_size_gb > 100
        
        optimization_notes = []
        
        # Select optimal transfer tool
        if is_bulk_operation and self.available_tools.get(TransferTool.S5CMD, False):
            tool = TransferTool.S5CMD
            optimization_notes.append("s5cmd selected for bulk operations (32x faster than s3cmd)")
        elif source_path.startswith(('azure:', 'gcs:', 'gdrive:')) and self.available_tools.get(TransferTool.RCLONE, False):
            tool = TransferTool.RCLONE
            optimization_notes.append("rclone selected for multi-cloud transfer")
        elif is_large_dataset and self.available_tools.get(TransferTool.AWS_CLI, False):
            tool = TransferTool.AWS_CLI_OPTIMIZED
            optimization_notes.append("Optimized AWS CLI selected for large dataset")
        else:
            tool = TransferTool.AWS_CLI
            optimization_notes.append("Standard AWS CLI fallback")
        
        # Determine optimal storage class
        storage_class = self._select_storage_class(research_domain, access_pattern, data_size_gb)
        
        # Calculate parallel workers based on system resources and data characteristics
        cpu_count = psutil.cpu_count()
        available_memory_gb = psutil.virtual_memory().available / (1024**3)
        
        if tool == TransferTool.S5CMD:
            # s5cmd can handle high concurrency efficiently
            parallel_workers = min(cpu_count * 4, 50)
        elif tool == TransferTool.RCLONE:
            # rclone moderate concurrency
            parallel_workers = min(cpu_count * 2, 20)
        else:
            # AWS CLI conservative concurrency
            parallel_workers = min(cpu_count, 10)
        
        # Optimize multipart settings
        if avg_file_size_mb > 100:
            multipart_threshold_mb = 64
            chunk_size_mb = min(64, max(8, int(avg_file_size_mb / 10)))
        else:
            multipart_threshold_mb = 8
            chunk_size_mb = 8
        
        # Compression decision
        enable_compression = avg_file_size_mb > 10 and not any(
            source_path.endswith(ext) for ext in ['.gz', '.bz2', '.zip', '.7z', '.bam', '.cram']
        )
        
        # Transfer acceleration decision (cost vs speed tradeoff)
        use_transfer_acceleration = is_large_dataset and data_size_gb > 1000
        if use_transfer_acceleration:
            optimization_notes.append("Transfer acceleration enabled for large dataset (additional cost)")
        
        # Cost estimation
        storage_cost_per_gb = self._get_storage_cost(storage_class)
        transfer_cost_per_gb = 0.0  # Ingress to AWS is free
        if use_transfer_acceleration:
            transfer_cost_per_gb += 0.04  # Transfer acceleration cost
        
        estimated_cost_per_gb = storage_cost_per_gb + transfer_cost_per_gb
        
        # Time estimation based on tool performance
        tool_profile = self.tool_profiles[tool]
        effective_throughput_gbps = tool_profile['max_throughput_gbps'] * tool_profile['parallel_efficiency']
        estimated_time_hours = (data_size_gb / (effective_throughput_gbps * 125)) + (tool_profile['overhead_seconds'] / 3600)
        
        return TransferStrategy(
            tool=tool,
            storage_class=storage_class,
            parallel_workers=parallel_workers,
            multipart_threshold_mb=multipart_threshold_mb,
            chunk_size_mb=chunk_size_mb,
            enable_compression=enable_compression,
            use_transfer_acceleration=use_transfer_acceleration,
            estimated_cost_per_gb=estimated_cost_per_gb,
            estimated_time_hours=estimated_time_hours,
            optimization_notes=optimization_notes
        )
    
    def _select_storage_class(self, research_domain: str, access_pattern: str, data_size_gb: float) -> StorageClass:
        """Select optimal S3 storage class based on research patterns."""
        
        # Research domain-specific recommendations
        if research_domain in self.research_access_patterns:
            domain_pattern = self.research_access_patterns[research_domain]
            if domain_pattern.get('initial_access_intensive') and access_pattern in ['frequent', 'unknown']:
                return StorageClass.INTELLIGENT_TIERING
            elif domain_pattern.get('archive_simulation_outputs') and data_size_gb > 500:
                return StorageClass.STANDARD_IA
        
        # General pattern-based selection
        if access_pattern == 'frequent':
            return StorageClass.STANDARD
        elif access_pattern == 'infrequent':
            return StorageClass.STANDARD_IA if data_size_gb > 10 else StorageClass.STANDARD
        elif access_pattern == 'archive':
            return StorageClass.GLACIER if data_size_gb > 100 else StorageClass.STANDARD_IA
        elif access_pattern == 'deep_archive':
            return StorageClass.DEEP_ARCHIVE
        else:
            # Default to intelligent tiering for unknown patterns
            return StorageClass.INTELLIGENT_TIERING
    
    def _get_storage_cost(self, storage_class: StorageClass) -> float:
        """Get storage cost per GB per month for storage class."""
        costs = {
            StorageClass.STANDARD: 0.023,
            StorageClass.INTELLIGENT_TIERING: 0.0125,
            StorageClass.STANDARD_IA: 0.0125,
            StorageClass.ONEZONE_IA: 0.01,
            StorageClass.GLACIER: 0.004,
            StorageClass.DEEP_ARCHIVE: 0.00099
        }
        return costs.get(storage_class, 0.023)
    
    async def execute_transfer(self, strategy: TransferStrategy, source_path: str, 
                             destination: str, dry_run: bool = False) -> Dict[str, Any]:
        """
        Execute optimized transfer using the selected strategy.
        
        Args:
            strategy: Transfer strategy from analyze_transfer_requirements
            source_path: Source path or pattern
            destination: Destination S3 path
            dry_run: If True, only show what would be executed
        
        Returns:
            Dict containing transfer results and performance metrics
        """
        
        start_time = time.time()
        
        if strategy.tool == TransferTool.S5CMD:
            result = await self._execute_s5cmd_transfer(strategy, source_path, destination, dry_run)
        elif strategy.tool == TransferTool.RCLONE:
            result = await self._execute_rclone_transfer(strategy, source_path, destination, dry_run)
        else:
            result = await self._execute_aws_cli_transfer(strategy, source_path, destination, dry_run)
        
        execution_time = time.time() - start_time
        
        # Add performance metrics
        result.update({
            'strategy_used': asdict(strategy),
            'execution_time_seconds': execution_time,
            'estimated_vs_actual_time': execution_time / (strategy.estimated_time_hours * 3600) if strategy.estimated_time_hours > 0 else None,
            'optimization_notes': strategy.optimization_notes
        })
        
        return result
    
    async def _execute_s5cmd_transfer(self, strategy: TransferStrategy, source_path: str, 
                                    destination: str, dry_run: bool) -> Dict[str, Any]:
        """Execute transfer using s5cmd for high performance."""
        
        cmd = ['s5cmd']
        
        # Add performance optimizations
        cmd.extend(['--numworkers', str(strategy.parallel_workers)])
        
        if dry_run:
            cmd.append('--dry-run')
        
        # Storage class configuration
        if strategy.storage_class != StorageClass.STANDARD:
            cmd.extend(['--storage-class', strategy.storage_class.value])
        
        # Select operation type
        if os.path.isdir(source_path):
            cmd.extend(['sync', source_path, destination])
        else:
            cmd.extend(['cp', source_path, destination])
        
        self.logger.info(f"Executing s5cmd: {' '.join(cmd)}")
        
        if dry_run:
            return {'command': ' '.join(cmd), 'dry_run': True}
        
        # Execute the command
        process = await asyncio.create_subprocess_exec(
            *cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )
        
        stdout, stderr = await process.communicate()
        
        return {
            'tool': 's5cmd',
            'success': process.returncode == 0,
            'stdout': stdout.decode(),
            'stderr': stderr.decode(),
            'return_code': process.returncode
        }
    
    async def _execute_rclone_transfer(self, strategy: TransferStrategy, source_path: str, 
                                     destination: str, dry_run: bool) -> Dict[str, Any]:
        """Execute transfer using rclone for multi-cloud scenarios."""
        
        cmd = ['rclone', 'sync']
        
        # Add performance optimizations
        cmd.extend(['--transfers', str(strategy.parallel_workers)])
        cmd.extend(['--checkers', str(strategy.parallel_workers)])
        
        if strategy.enable_compression:
            cmd.append('--compress')
        
        if dry_run:
            cmd.append('--dry-run')
        
        # Add source and destination
        cmd.extend([source_path, destination])
        
        # Verbose output for monitoring
        cmd.extend(['-v', '--stats', '30s'])
        
        self.logger.info(f"Executing rclone: {' '.join(cmd)}")
        
        if dry_run:
            return {'command': ' '.join(cmd), 'dry_run': True}
        
        # Execute the command
        process = await asyncio.create_subprocess_exec(
            *cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )
        
        stdout, stderr = await process.communicate()
        
        return {
            'tool': 'rclone',
            'success': process.returncode == 0,
            'stdout': stdout.decode(),
            'stderr': stderr.decode(),
            'return_code': process.returncode
        }
    
    async def _execute_aws_cli_transfer(self, strategy: TransferStrategy, source_path: str, 
                                      destination: str, dry_run: bool) -> Dict[str, Any]:
        """Execute transfer using optimized AWS CLI."""
        
        cmd = ['aws', 's3']
        
        # Configure AWS CLI for performance
        env = os.environ.copy()
        env.update({
            'AWS_CLI_FILE_ENCODING': 'UTF-8',
            'AWS_DEFAULT_OUTPUT': 'json'
        })
        
        # Set multipart threshold and chunk size
        aws_config_dir = Path.home() / '.aws'
        aws_config_dir.mkdir(exist_ok=True)
        
        config_content = f"""
[default]
s3 =
    max_concurrent_requests = {strategy.parallel_workers}
    max_bandwidth = 1GB/s
    multipart_threshold = {strategy.multipart_threshold_mb}MB
    multipart_chunksize = {strategy.chunk_size_mb}MB
    use_accelerate_endpoint = {str(strategy.use_transfer_acceleration).lower()}
"""
        
        with open(aws_config_dir / 'config', 'w') as f:
            f.write(config_content)
        
        # Select operation
        if os.path.isdir(source_path):
            cmd.extend(['sync', source_path, destination])
        else:
            cmd.extend(['cp', source_path, destination])
        
        # Add storage class if not standard
        if strategy.storage_class != StorageClass.STANDARD:
            cmd.extend(['--storage-class', strategy.storage_class.value])
        
        if dry_run:
            cmd.append('--dryrun')
        
        self.logger.info(f"Executing AWS CLI: {' '.join(cmd)}")
        
        if dry_run:
            return {'command': ' '.join(cmd), 'dry_run': True}
        
        # Execute the command
        process = await asyncio.create_subprocess_exec(
            *cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            env=env
        )
        
        stdout, stderr = await process.communicate()
        
        return {
            'tool': 'aws_cli',
            'success': process.returncode == 0,
            'stdout': stdout.decode(),
            'stderr': stderr.decode(),
            'return_code': process.returncode
        }

def demonstrate_s3_optimization():
    """Demonstrate S3 transfer optimization capabilities."""
    
    optimizer = S3TransferOptimizer()
    
    # Example: Large genomics dataset
    genomics_strategy = optimizer.analyze_transfer_requirements(
        source_path='/data/genomics/samples/',
        destination='s3://research-genomics/raw-data/',
        data_size_gb=2500,  # 2.5 TB
        file_count=15000,   # Many FASTQ files
        research_domain='genomics',
        access_pattern='frequent'
    )
    
    print("Genomics Dataset Transfer Strategy:")
    print(f"Tool: {genomics_strategy.tool.value}")
    print(f"Storage Class: {genomics_strategy.storage_class.value}")
    print(f"Parallel Workers: {genomics_strategy.parallel_workers}")
    print(f"Estimated Cost: ${genomics_strategy.estimated_cost_per_gb:.4f}/GB")
    print(f"Estimated Time: {genomics_strategy.estimated_time_hours:.2f} hours")
    print("Optimizations:")
    for note in genomics_strategy.optimization_notes:
        print(f"  - {note}")

if __name__ == "__main__":
    demonstrate_s3_optimization()