#!/usr/bin/env python3
"""
Data Management Engine for AWS Research Wizard

This module provides comprehensive data access coordination, including:
1. Static data integration from various sources (institutional storage, other clouds)
2. Real-time instrument data ingestion and processing
3. Data transfer optimization and cost management
4. Metadata management and data discovery
5. Integration with compute workflows

Key Features:
- Multi-source data integration (S3, institutional storage, other cloud providers)
- Real-time streaming data ingestion from scientific instruments
- Intelligent data transfer optimization (cost vs. speed)
- Automated data preprocessing and quality control
- Event-driven workflow triggering based on data availability
- Comprehensive data lineage and provenance tracking

Classes:
    DataSource: Abstract base class for different data source types
    StaticDataSource: For existing datasets on various storage systems
    InstrumentDataSource: For real-time instrument data streams
    DataTransferManager: Optimizes data movement and cost management
    DataWorkflowTrigger: Triggers compute workflows based on data events
    DataManagementEngine: Main orchestration class

Dependencies:
    - boto3: For AWS service integration
    - asyncio: For asynchronous data streaming
    - watchdog: For file system monitoring
    - paramiko: For SFTP/SSH data access
    - requests: For HTTP-based data APIs
"""

import os
import sys
import yaml
import json
import boto3
import asyncio
import logging
import hashlib
import threading
import time
from typing import Dict, List, Any, Optional, Union, Callable, AsyncIterator
from pathlib import Path
from dataclasses import dataclass, asdict
from datetime import datetime, timedelta
from enum import Enum
from abc import ABC, abstractmethod
import concurrent.futures
from urllib.parse import urlparse

# Optional imports for enhanced functionality
try:
    import paramiko
    SFTP_AVAILABLE = True
except ImportError:
    SFTP_AVAILABLE = False
    paramiko = None

try:
    from watchdog.observers import Observer
    from watchdog.events import FileSystemEventHandler
    WATCHDOG_AVAILABLE = True
except ImportError:
    WATCHDOG_AVAILABLE = False
    Observer = None
    FileSystemEventHandler = None

try:
    import aioboto3
    ASYNC_AWS_AVAILABLE = True
except ImportError:
    ASYNC_AWS_AVAILABLE = False
    aioboto3 = None

# Import our core modules
from config_loader import ConfigLoader
from demo_workflow_engine import DemoWorkflowEngine, ExecutionEnvironment

class DataSourceType(Enum):
    """Types of data sources supported by the system."""
    AWS_S3 = "aws_s3"
    INSTITUTIONAL_STORAGE = "institutional_storage"
    SFTP_SERVER = "sftp_server"
    HTTP_API = "http_api"
    LOCAL_FILESYSTEM = "local_filesystem"
    REAL_TIME_INSTRUMENT = "real_time_instrument"
    DATABASE = "database"
    OTHER_CLOUD = "other_cloud"

class TransferPriority(Enum):
    """Data transfer priority levels."""
    IMMEDIATE = "immediate"  # Real-time, cost is secondary
    URGENT = "urgent"       # Within hours, balanced cost/speed
    STANDARD = "standard"   # Within days, cost-optimized
    ARCHIVE = "archive"     # When needed, lowest cost

@dataclass
class DataTransferCost:
    """Data transfer cost estimation and optimization."""
    estimated_cost_usd: float
    transfer_time_estimate: str
    bandwidth_required_mbps: float
    storage_cost_monthly: float
    recommended_strategy: str
    cost_breakdown: Dict[str, float]

@dataclass
class DataMetadata:
    """Comprehensive metadata for tracked datasets."""
    dataset_id: str
    source_location: str
    destination_location: Optional[str]
    size_bytes: int
    checksum_md5: Optional[str]
    checksum_sha256: Optional[str]
    file_format: str
    creation_time: datetime
    last_modified: datetime
    data_type: str
    research_domain: str
    access_permissions: Dict[str, Any]
    quality_metrics: Optional[Dict[str, Any]]
    provenance: Dict[str, Any]
    tags: List[str]

@dataclass
class InstrumentConfig:
    """Configuration for real-time instrument data ingestion."""
    instrument_id: str
    instrument_type: str
    data_rate_mb_per_hour: float
    data_format: str
    quality_control_rules: List[Dict[str, Any]]
    preprocessing_steps: List[str]
    buffer_size_mb: int
    trigger_conditions: Dict[str, Any]
    metadata_extraction: Dict[str, str]

class DataSource(ABC):
    """Abstract base class for all data source types."""

    def __init__(self, source_id: str, source_type: DataSourceType, config: Dict[str, Any]):
        self.source_id = source_id
        self.source_type = source_type
        self.config = config
        self.logger = logging.getLogger(f"{__name__}.{source_id}")
        self.metadata_cache: Dict[str, DataMetadata] = {}

    @abstractmethod
    async def list_available_data(self, filters: Optional[Dict[str, Any]] = None) -> List[DataMetadata]:
        """List all available data from this source."""
        pass

    @abstractmethod
    async def get_data_metadata(self, data_path: str) -> DataMetadata:
        """Get detailed metadata for specific data."""
        pass

    @abstractmethod
    async def estimate_transfer_cost(self, data_path: str, destination: str) -> DataTransferCost:
        """Estimate cost and time for data transfer."""
        pass

    @abstractmethod
    async def initiate_transfer(self, data_path: str, destination: str,
                              priority: TransferPriority) -> str:
        """Initiate data transfer and return transfer job ID."""
        pass

    def _calculate_checksum(self, file_path: str, algorithm: str = "md5") -> str:
        """Calculate file checksum for integrity verification."""
        hash_func = hashlib.md5() if algorithm == "md5" else hashlib.sha256()

        with open(file_path, 'rb') as f:
            for chunk in iter(lambda: f.read(4096), b""):
                hash_func.update(chunk)

        return hash_func.hexdigest()

class StaticDataSource(DataSource):
    """
    Data source for static datasets stored on various systems.

    Supports:
    - AWS S3 buckets (including cross-account access)
    - Institutional storage systems (NFS, Lustre, GPFS)
    - SFTP/SSH accessible storage
    - HTTP-based data APIs
    - Other cloud providers (Azure Blob, Google Cloud Storage)
    """

    def __init__(self, source_id: str, source_config: Dict[str, Any]):
        source_type = DataSourceType(source_config.get("type", "aws_s3"))
        super().__init__(source_id, source_type, source_config)

        # Initialize source-specific clients
        if self.source_type == DataSourceType.AWS_S3:
            self.s3_client = boto3.client('s3')
        elif self.source_type == DataSourceType.SFTP_SERVER and SFTP_AVAILABLE:
            self._init_sftp_client()

    def _init_sftp_client(self):
        """Initialize SFTP client for secure file transfer."""
        self.ssh_client = paramiko.SSHClient()
        self.ssh_client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

        try:
            self.ssh_client.connect(
                hostname=self.config['hostname'],
                username=self.config['username'],
                key_filename=self.config.get('key_file'),
                password=self.config.get('password')
            )
            self.sftp_client = self.ssh_client.open_sftp()
            self.logger.info(f"Connected to SFTP server: {self.config['hostname']}")
        except Exception as e:
            self.logger.error(f"Failed to connect to SFTP server: {e}")
            raise

    async def list_available_data(self, filters: Optional[Dict[str, Any]] = None) -> List[DataMetadata]:
        """List available datasets from static storage."""
        if self.source_type == DataSourceType.AWS_S3:
            return await self._list_s3_data(filters)
        elif self.source_type == DataSourceType.SFTP_SERVER:
            return await self._list_sftp_data(filters)
        elif self.source_type == DataSourceType.HTTP_API:
            return await self._list_http_api_data(filters)
        else:
            raise NotImplementedError(f"Listing not implemented for {self.source_type}")

    async def _list_s3_data(self, filters: Optional[Dict[str, Any]] = None) -> List[DataMetadata]:
        """List data from S3 bucket."""
        bucket = self.config['bucket']
        prefix = self.config.get('prefix', '')

        datasets = []

        try:
            paginator = self.s3_client.get_paginator('list_objects_v2')

            for page in paginator.paginate(Bucket=bucket, Prefix=prefix):
                for obj in page.get('Contents', []):
                    # Apply filters if provided
                    if filters and not self._apply_filters(obj, filters):
                        continue

                    metadata = DataMetadata(
                        dataset_id=f"{self.source_id}:{obj['Key']}",
                        source_location=f"s3://{bucket}/{obj['Key']}",
                        destination_location=None,
                        size_bytes=obj['Size'],
                        checksum_md5=obj.get('ETag', '').strip('"'),
                        checksum_sha256=None,
                        file_format=self._detect_file_format(obj['Key']),
                        creation_time=obj['LastModified'],
                        last_modified=obj['LastModified'],
                        data_type=self._infer_data_type(obj['Key']),
                        research_domain=self.config.get('research_domain', 'unknown'),
                        access_permissions={'public_read': True},
                        quality_metrics=None,
                        provenance={'source': f"s3://{bucket}", 'ingestion_time': datetime.now()},
                        tags=self._extract_tags(obj['Key'])
                    )

                    datasets.append(metadata)
                    self.metadata_cache[metadata.dataset_id] = metadata

        except Exception as e:
            self.logger.error(f"Error listing S3 data: {e}")
            raise

        return datasets

    async def get_data_metadata(self, data_path: str) -> DataMetadata:
        """Get detailed metadata for specific dataset."""
        dataset_id = f"{self.source_id}:{data_path}"

        # Check cache first
        if dataset_id in self.metadata_cache:
            return self.metadata_cache[dataset_id]

        # Fetch fresh metadata
        if self.source_type == DataSourceType.AWS_S3:
            return await self._get_s3_metadata(data_path)
        else:
            raise NotImplementedError(f"Metadata not implemented for {self.source_type}")

    async def estimate_transfer_cost(self, data_path: str, destination: str) -> DataTransferCost:
        """Estimate cost for transferring static data."""
        metadata = await self.get_data_metadata(data_path)
        size_gb = metadata.size_bytes / (1024**3)

        # Cost calculations based on AWS pricing (simplified)
        if destination.startswith('s3://'):
            # S3 to S3 transfer
            transfer_cost = 0.0  # Free within same region
            storage_cost = size_gb * 0.023  # Standard S3 storage per month
            bandwidth_mbps = min(1000, size_gb * 8 / 3600)  # Assume 1 hour for small files
            transfer_time = f"{max(1, size_gb / 10)} hours"  # 10 GB/hour estimate
            strategy = "Direct S3 copy"
        else:
            # S3 to EC2/external
            transfer_cost = size_gb * 0.09 if size_gb > 1 else 0.0  # First 1GB free
            storage_cost = size_gb * 0.0125  # EBS GP3 storage
            bandwidth_mbps = min(500, size_gb * 8 / 7200)  # 2 hour estimate
            transfer_time = f"{max(1, size_gb / 5)} hours"  # 5 GB/hour estimate
            strategy = "Optimized download with parallel streams"

        return DataTransferCost(
            estimated_cost_usd=transfer_cost + storage_cost,
            transfer_time_estimate=transfer_time,
            bandwidth_required_mbps=bandwidth_mbps,
            storage_cost_monthly=storage_cost,
            recommended_strategy=strategy,
            cost_breakdown={
                'transfer': transfer_cost,
                'storage': storage_cost,
                'compute': 0.0
            }
        )

    async def initiate_transfer(self, data_path: str, destination: str,
                              priority: TransferPriority) -> str:
        """Initiate transfer of static data."""
        transfer_id = f"transfer_{int(time.time())}_{hash(data_path) % 10000}"

        if self.source_type == DataSourceType.AWS_S3:
            return await self._initiate_s3_transfer(data_path, destination, priority, transfer_id)
        else:
            raise NotImplementedError(f"Transfer not implemented for {self.source_type}")

    async def _initiate_s3_transfer(self, data_path: str, destination: str,
                                   priority: TransferPriority, transfer_id: str) -> str:
        """Initiate S3-based data transfer."""
        bucket = self.config['bucket']

        if destination.startswith('s3://'):
            # S3 to S3 copy
            dest_parts = destination.replace('s3://', '').split('/', 1)
            dest_bucket = dest_parts[0]
            dest_key = dest_parts[1] if len(dest_parts) > 1 else data_path

            copy_source = {'Bucket': bucket, 'Key': data_path}
            self.s3_client.copy_object(
                CopySource=copy_source,
                Bucket=dest_bucket,
                Key=dest_key
            )
        else:
            # S3 to local/EBS download
            local_path = Path(destination) / Path(data_path).name
            local_path.parent.mkdir(parents=True, exist_ok=True)

            self.s3_client.download_file(bucket, data_path, str(local_path))

        self.logger.info(f"Initiated transfer {transfer_id}: {data_path} -> {destination}")
        return transfer_id

    def _detect_file_format(self, filename: str) -> str:
        """Detect file format from filename extension."""
        ext = Path(filename).suffix.lower()
        format_map = {
            '.nc': 'NetCDF',
            '.h5': 'HDF5',
            '.hdf5': 'HDF5',
            '.tif': 'GeoTIFF',
            '.tiff': 'GeoTIFF',
            '.fastq': 'FASTQ',
            '.fq': 'FASTQ',
            '.bam': 'BAM',
            '.vcf': 'VCF',
            '.csv': 'CSV',
            '.json': 'JSON',
            '.parquet': 'Parquet'
        }
        return format_map.get(ext, 'Unknown')

    def _infer_data_type(self, filename: str) -> str:
        """Infer research data type from filename patterns."""
        filename_lower = filename.lower()

        if any(term in filename_lower for term in ['temperature', 'climate', 'weather']):
            return 'climate_data'
        elif any(term in filename_lower for term in ['genome', 'dna', 'rna', 'seq']):
            return 'genomic_data'
        elif any(term in filename_lower for term in ['satellite', 'landsat', 'sentinel']):
            return 'remote_sensing'
        elif any(term in filename_lower for term in ['ocean', 'marine', 'sea']):
            return 'oceanographic'
        else:
            return 'research_data'

    def _extract_tags(self, filename: str) -> List[str]:
        """Extract relevant tags from filename for categorization."""
        tags = []
        filename_lower = filename.lower()

        # Date patterns
        import re
        if re.search(r'\d{4}', filename):
            tags.append('time_series')

        # Geographic indicators
        if any(geo in filename_lower for geo in ['global', 'world', 'earth']):
            tags.append('global_scale')
        elif any(geo in filename_lower for geo in ['region', 'local', 'site']):
            tags.append('regional_scale')

        # Quality indicators
        if any(qual in filename_lower for qual in ['qc', 'quality', 'validated']):
            tags.append('quality_controlled')

        return tags

    def _apply_filters(self, obj: Dict[str, Any], filters: Dict[str, Any]) -> bool:
        """Apply filters to determine if object should be included."""
        if 'min_size' in filters and obj['Size'] < filters['min_size']:
            return False
        if 'max_size' in filters and obj['Size'] > filters['max_size']:
            return False
        if 'file_extension' in filters:
            ext = Path(obj['Key']).suffix.lower()
            if ext not in filters['file_extension']:
                return False

        return True

class InstrumentDataSource(DataSource):
    """
    Data source for real-time instrument data streams.

    Supports:
    - Scientific instrument data feeds (microscopes, telescopes, sequencers)
    - Sensor networks and IoT devices
    - Laboratory equipment with data export capabilities
    - Real-time streaming protocols (TCP, UDP, HTTP streaming)
    """

    def __init__(self, source_id: str, instrument_config: InstrumentConfig):
        config = asdict(instrument_config)
        super().__init__(source_id, DataSourceType.REAL_TIME_INSTRUMENT, config)

        self.instrument_config = instrument_config
        self.data_buffer: List[bytes] = []
        self.buffer_size_current = 0
        self.streaming_active = False
        self.processing_queue = asyncio.Queue()

        # Initialize instrument-specific connections
        self._init_instrument_connection()

    def _init_instrument_connection(self):
        """Initialize connection to scientific instrument."""
        instrument_type = self.instrument_config.instrument_type

        if instrument_type == "tcp_stream":
            self._init_tcp_connection()
        elif instrument_type == "http_stream":
            self._init_http_stream()
        elif instrument_type == "file_watcher":
            self._init_file_watcher()
        else:
            self.logger.warning(f"Unknown instrument type: {instrument_type}")

    def _init_tcp_connection(self):
        """Initialize TCP connection for real-time data streaming."""
        # Implementation would depend on specific instrument protocol
        self.logger.info("Initialized TCP connection for instrument data")

    def _init_http_stream(self):
        """Initialize HTTP streaming connection."""
        # Implementation for HTTP-based data streams
        self.logger.info("Initialized HTTP streaming for instrument data")

    def _init_file_watcher(self):
        """Initialize file system watcher for instrument output files."""
        if not WATCHDOG_AVAILABLE:
            raise ImportError("Watchdog required for file watching functionality")

        class InstrumentFileHandler(FileSystemEventHandler):
            def __init__(self, data_source):
                self.data_source = data_source

            def on_created(self, event):
                if not event.is_directory:
                    self.data_source._handle_new_file(event.src_path)

        self.file_handler = InstrumentFileHandler(self)
        self.observer = Observer()
        watch_path = self.config.get('watch_directory', '/tmp/instrument_data')
        self.observer.schedule(self.file_handler, watch_path, recursive=True)
        self.observer.start()

        self.logger.info(f"Started file watcher for: {watch_path}")

    async def list_available_data(self, filters: Optional[Dict[str, Any]] = None) -> List[DataMetadata]:
        """List recently captured instrument data."""
        # For real-time instruments, this returns recent data buffer contents
        recent_data = []

        # This would be implemented based on the specific buffering strategy
        # For now, return placeholder data
        current_time = datetime.now()

        if self.buffer_size_current > 0:
            metadata = DataMetadata(
                dataset_id=f"{self.source_id}:buffer_{int(time.time())}",
                source_location=f"instrument://{self.instrument_config.instrument_id}/buffer",
                destination_location=None,
                size_bytes=self.buffer_size_current,
                checksum_md5=None,
                checksum_sha256=None,
                file_format=self.instrument_config.data_format,
                creation_time=current_time,
                last_modified=current_time,
                data_type="real_time_instrument",
                research_domain=self.config.get('research_domain', 'unknown'),
                access_permissions={'real_time': True},
                quality_metrics=self._assess_data_quality(),
                provenance={
                    'instrument_id': self.instrument_config.instrument_id,
                    'capture_time': current_time,
                    'data_rate': self.instrument_config.data_rate_mb_per_hour
                },
                tags=['real_time', 'instrument_data']
            )
            recent_data.append(metadata)

        return recent_data

    async def get_data_metadata(self, data_path: str) -> DataMetadata:
        """Get metadata for instrument data stream."""
        # For real-time data, metadata is dynamically generated
        current_time = datetime.now()

        return DataMetadata(
            dataset_id=f"{self.source_id}:{data_path}",
            source_location=f"instrument://{self.instrument_config.instrument_id}/{data_path}",
            destination_location=None,
            size_bytes=self.buffer_size_current,
            checksum_md5=None,
            checksum_sha256=None,
            file_format=self.instrument_config.data_format,
            creation_time=current_time,
            last_modified=current_time,
            data_type="real_time_instrument",
            research_domain=self.config.get('research_domain', 'unknown'),
            access_permissions={'real_time': True},
            quality_metrics=self._assess_data_quality(),
            provenance={
                'instrument_id': self.instrument_config.instrument_id,
                'capture_time': current_time
            },
            tags=['real_time', 'streaming']
        )

    async def estimate_transfer_cost(self, data_path: str, destination: str) -> DataTransferCost:
        """Estimate cost for real-time data transfer."""
        # Real-time data typically has minimal transfer cost but requires streaming infrastructure
        data_rate_mbps = self.instrument_config.data_rate_mb_per_hour / 3600 * 8

        # Estimate based on streaming for 1 hour
        streaming_cost = 0.05  # Placeholder cost for streaming infrastructure
        storage_cost = (self.instrument_config.data_rate_mb_per_hour / 1024) * 0.023  # S3 storage

        return DataTransferCost(
            estimated_cost_usd=streaming_cost + storage_cost,
            transfer_time_estimate="Real-time",
            bandwidth_required_mbps=data_rate_mbps,
            storage_cost_monthly=storage_cost * 24 * 30,  # Monthly estimate
            recommended_strategy="Real-time streaming with buffering",
            cost_breakdown={
                'streaming': streaming_cost,
                'storage': storage_cost,
                'compute': 0.0
            }
        )

    async def initiate_transfer(self, data_path: str, destination: str,
                              priority: TransferPriority) -> str:
        """Initiate real-time data streaming."""
        transfer_id = f"stream_{self.instrument_config.instrument_id}_{int(time.time())}"

        # Start streaming task
        if not self.streaming_active:
            self.streaming_active = True
            asyncio.create_task(self._stream_data_to_destination(destination, transfer_id))

        self.logger.info(f"Initiated real-time streaming {transfer_id} to {destination}")
        return transfer_id

    async def _stream_data_to_destination(self, destination: str, transfer_id: str):
        """Stream real-time data to specified destination."""
        while self.streaming_active:
            try:
                # Wait for data in the processing queue
                data_chunk = await asyncio.wait_for(self.processing_queue.get(), timeout=1.0)

                # Process and send data
                await self._send_data_chunk(data_chunk, destination)

                # Apply quality control
                if not self._passes_quality_control(data_chunk):
                    self.logger.warning(f"Data chunk failed quality control: {transfer_id}")

            except asyncio.TimeoutError:
                # No data received, continue waiting
                continue
            except Exception as e:
                self.logger.error(f"Error in data streaming {transfer_id}: {e}")
                break

    async def _send_data_chunk(self, data_chunk: bytes, destination: str):
        """Send data chunk to destination."""
        if destination.startswith('s3://'):
            # Stream to S3
            await self._stream_to_s3(data_chunk, destination)
        else:
            # Stream to local storage
            await self._stream_to_local(data_chunk, destination)

    async def _stream_to_s3(self, data_chunk: bytes, s3_path: str):
        """Stream data chunk to S3."""
        # Implementation would use S3 streaming upload
        pass

    async def _stream_to_local(self, data_chunk: bytes, local_path: str):
        """Stream data chunk to local filesystem."""
        # Implementation would append to local file
        pass

    def _handle_new_file(self, file_path: str):
        """Handle new file detected by file watcher."""
        self.logger.info(f"New instrument file detected: {file_path}")

        # Add file to processing queue
        try:
            with open(file_path, 'rb') as f:
                data = f.read()
                asyncio.create_task(self.processing_queue.put(data))
        except Exception as e:
            self.logger.error(f"Error reading instrument file {file_path}: {e}")

    def _assess_data_quality(self) -> Dict[str, Any]:
        """Assess quality of current data buffer."""
        if self.buffer_size_current == 0:
            return {'status': 'no_data', 'score': 0.0}

        # Implement quality assessment based on instrument-specific rules
        quality_score = 0.95  # Placeholder

        return {
            'status': 'good' if quality_score > 0.8 else 'warning',
            'score': quality_score,
            'buffer_utilization': self.buffer_size_current / (self.instrument_config.buffer_size_mb * 1024 * 1024),
            'data_rate_actual': self._calculate_actual_data_rate()
        }

    def _passes_quality_control(self, data_chunk: bytes) -> bool:
        """Check if data chunk passes quality control rules."""
        # Implement quality control based on instrument configuration
        for rule in self.instrument_config.quality_control_rules:
            if not self._apply_quality_rule(data_chunk, rule):
                return False
        return True

    def _apply_quality_rule(self, data_chunk: bytes, rule: Dict[str, Any]) -> bool:
        """Apply specific quality control rule."""
        rule_type = rule.get('type')

        if rule_type == 'min_size':
            return len(data_chunk) >= rule['value']
        elif rule_type == 'max_size':
            return len(data_chunk) <= rule['value']
        elif rule_type == 'format_check':
            # Implement format-specific validation
            return True

        return True

    def _calculate_actual_data_rate(self) -> float:
        """Calculate actual data rate from recent buffer activity."""
        # This would track actual data throughput
        return self.instrument_config.data_rate_mb_per_hour  # Placeholder

class DataTransferManager:
    """
    Manages and optimizes data transfers across different sources and destinations.

    Features:
    - Cost-optimized transfer strategies
    - Parallel transfer coordination
    - Bandwidth management and throttling
    - Transfer job scheduling and prioritization
    - Error handling and retry logic
    """

    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)

        # Transfer job tracking
        self.active_transfers: Dict[str, Dict[str, Any]] = {}
        self.transfer_history: List[Dict[str, Any]] = []

        # Configuration
        self.transfer_config = {
            'max_concurrent_transfers': 5,
            'bandwidth_limit_mbps': 1000,
            'retry_attempts': 3,
            'cost_optimization_enabled': True,
            'transfer_timeout_hours': 24
        }

    async def optimize_transfer_strategy(self, metadata: DataMetadata,
                                       destination: str, priority: TransferPriority) -> Dict[str, Any]:
        """Determine optimal transfer strategy based on cost, speed, and priority."""
        size_gb = metadata.size_bytes / (1024**3)

        strategies = []

        # Direct transfer strategy
        direct_cost = self._estimate_direct_transfer_cost(size_gb, destination)
        strategies.append({
            'name': 'direct_transfer',
            'cost': direct_cost,
            'time_estimate': f"{max(1, size_gb / 10)} hours",
            'reliability': 0.95,
            'complexity': 'low'
        })

        # Multi-part parallel transfer
        if size_gb > 1:
            parallel_cost = direct_cost * 1.1  # Slightly higher cost for coordination
            strategies.append({
                'name': 'parallel_multipart',
                'cost': parallel_cost,
                'time_estimate': f"{max(0.5, size_gb / 20)} hours",
                'reliability': 0.90,
                'complexity': 'medium'
            })

        # Compressed transfer (if applicable)
        if metadata.file_format in ['CSV', 'JSON', 'Text']:
            compressed_cost = direct_cost * 0.7  # Assume 30% compression
            strategies.append({
                'name': 'compressed_transfer',
                'cost': compressed_cost,
                'time_estimate': f"{max(1, size_gb * 0.7 / 10)} hours",
                'reliability': 0.92,
                'complexity': 'medium'
            })

        # Select best strategy based on priority
        if priority == TransferPriority.IMMEDIATE:
            best_strategy = min(strategies, key=lambda x: float(x['time_estimate'].split()[0]))
        elif priority == TransferPriority.ARCHIVE:
            best_strategy = min(strategies, key=lambda x: x['cost'])
        else:
            # Balanced approach
            scored_strategies = []
            for strategy in strategies:
                time_hours = float(strategy['time_estimate'].split()[0])
                score = strategy['cost'] + (time_hours * 10)  # Weight time vs cost
                scored_strategies.append((score, strategy))
            best_strategy = min(scored_strategies, key=lambda x: x[0])[1]

        return best_strategy

    def _estimate_direct_transfer_cost(self, size_gb: float, destination: str) -> float:
        """Estimate cost for direct data transfer."""
        if destination.startswith('s3://'):
            return size_gb * 0.01  # S3 storage cost
        else:
            return size_gb * 0.09  # Data egress cost

    async def schedule_transfer(self, source: DataSource, data_path: str,
                              destination: str, priority: TransferPriority) -> str:
        """Schedule and manage data transfer job."""
        transfer_id = f"xfer_{int(time.time())}_{hash(data_path) % 10000}"

        # Get metadata and optimize strategy
        metadata = await source.get_data_metadata(data_path)
        strategy = await self.optimize_transfer_strategy(metadata, destination, priority)

        # Create transfer job
        transfer_job = {
            'transfer_id': transfer_id,
            'source_id': source.source_id,
            'data_path': data_path,
            'destination': destination,
            'priority': priority,
            'strategy': strategy,
            'metadata': metadata,
            'status': 'scheduled',
            'created_time': datetime.now(),
            'start_time': None,
            'completion_time': None,
            'error_message': None,
            'progress_percent': 0.0
        }

        self.active_transfers[transfer_id] = transfer_job

        # Start transfer asynchronously
        asyncio.create_task(self._execute_transfer(source, transfer_job))

        self.logger.info(f"Scheduled transfer {transfer_id} with strategy: {strategy['name']}")
        return transfer_id

    async def _execute_transfer(self, source: DataSource, transfer_job: Dict[str, Any]):
        """Execute the actual data transfer."""
        transfer_id = transfer_job['transfer_id']

        try:
            transfer_job['status'] = 'running'
            transfer_job['start_time'] = datetime.now()

            # Execute transfer based on strategy
            strategy_name = transfer_job['strategy']['name']

            if strategy_name == 'direct_transfer':
                await self._execute_direct_transfer(source, transfer_job)
            elif strategy_name == 'parallel_multipart':
                await self._execute_parallel_transfer(source, transfer_job)
            elif strategy_name == 'compressed_transfer':
                await self._execute_compressed_transfer(source, transfer_job)

            transfer_job['status'] = 'completed'
            transfer_job['completion_time'] = datetime.now()
            transfer_job['progress_percent'] = 100.0

            self.logger.info(f"Transfer {transfer_id} completed successfully")

        except Exception as e:
            transfer_job['status'] = 'failed'
            transfer_job['error_message'] = str(e)
            self.logger.error(f"Transfer {transfer_id} failed: {e}")

        finally:
            # Move to history
            self.transfer_history.append(transfer_job)
            if transfer_id in self.active_transfers:
                del self.active_transfers[transfer_id]

    async def _execute_direct_transfer(self, source: DataSource, transfer_job: Dict[str, Any]):
        """Execute direct transfer strategy."""
        await source.initiate_transfer(
            transfer_job['data_path'],
            transfer_job['destination'],
            transfer_job['priority']
        )

    async def _execute_parallel_transfer(self, source: DataSource, transfer_job: Dict[str, Any]):
        """Execute parallel multipart transfer strategy."""
        # Implementation would split large files into chunks and transfer in parallel
        await self._execute_direct_transfer(source, transfer_job)  # Simplified for now

    async def _execute_compressed_transfer(self, source: DataSource, transfer_job: Dict[str, Any]):
        """Execute compressed transfer strategy."""
        # Implementation would compress data before transfer
        await self._execute_direct_transfer(source, transfer_job)  # Simplified for now

    def get_transfer_status(self, transfer_id: str) -> Optional[Dict[str, Any]]:
        """Get status of transfer job."""
        if transfer_id in self.active_transfers:
            return self.active_transfers[transfer_id]

        # Check history
        for job in self.transfer_history:
            if job['transfer_id'] == transfer_id:
                return job

        return None

    def get_transfer_statistics(self) -> Dict[str, Any]:
        """Get comprehensive transfer statistics."""
        all_transfers = list(self.active_transfers.values()) + self.transfer_history

        if not all_transfers:
            return {'message': 'No transfers found'}

        completed_transfers = [t for t in all_transfers if t['status'] == 'completed']
        failed_transfers = [t for t in all_transfers if t['status'] == 'failed']

        total_data_gb = sum(t['metadata'].size_bytes / (1024**3) for t in completed_transfers)
        total_cost = sum(t['strategy']['cost'] for t in completed_transfers)

        avg_duration = 0
        if completed_transfers:
            durations = []
            for t in completed_transfers:
                if t['start_time'] and t['completion_time']:
                    duration = (t['completion_time'] - t['start_time']).total_seconds()
                    durations.append(duration)
            avg_duration = sum(durations) / len(durations) if durations else 0

        return {
            'total_transfers': len(all_transfers),
            'completed': len(completed_transfers),
            'failed': len(failed_transfers),
            'success_rate': len(completed_transfers) / len(all_transfers) * 100 if all_transfers else 0,
            'total_data_transferred_gb': total_data_gb,
            'total_cost_usd': total_cost,
            'average_duration_seconds': avg_duration,
            'active_transfers': len(self.active_transfers)
        }

class DataWorkflowTrigger:
    """
    Triggers compute workflows based on data availability and events.

    Features:
    - Event-driven workflow execution
    - Data dependency tracking
    - Conditional workflow triggers
    - Integration with workflow execution engine
    """

    def __init__(self, workflow_engine: DemoWorkflowEngine):
        self.workflow_engine = workflow_engine
        self.logger = logging.getLogger(__name__)

        # Trigger rules and monitoring
        self.trigger_rules: List[Dict[str, Any]] = []
        self.data_monitors: Dict[str, Any] = {}
        self.triggered_workflows: List[Dict[str, Any]] = []

    def add_trigger_rule(self, rule_config: Dict[str, Any]) -> str:
        """Add a new data-driven workflow trigger rule."""
        rule_id = f"rule_{int(time.time())}_{len(self.trigger_rules)}"

        rule = {
            'rule_id': rule_id,
            'name': rule_config['name'],
            'description': rule_config.get('description', ''),
            'data_conditions': rule_config['data_conditions'],
            'workflow_config': rule_config['workflow_config'],
            'trigger_frequency': rule_config.get('trigger_frequency', 'immediate'),
            'enabled': True,
            'created_time': datetime.now(),
            'last_triggered': None,
            'trigger_count': 0
        }

        self.trigger_rules.append(rule)
        self.logger.info(f"Added trigger rule: {rule['name']} ({rule_id})")

        return rule_id

    async def evaluate_triggers(self, data_event: Dict[str, Any]):
        """Evaluate all trigger rules against a data event."""
        for rule in self.trigger_rules:
            if not rule['enabled']:
                continue

            if await self._rule_matches_event(rule, data_event):
                await self._execute_triggered_workflow(rule, data_event)

    async def _rule_matches_event(self, rule: Dict[str, Any], data_event: Dict[str, Any]) -> bool:
        """Check if a trigger rule matches the data event."""
        conditions = rule['data_conditions']

        # Check data type condition
        if 'data_type' in conditions:
            if data_event.get('data_type') != conditions['data_type']:
                return False

        # Check size condition
        if 'min_size_mb' in conditions:
            size_mb = data_event.get('size_bytes', 0) / (1024 * 1024)
            if size_mb < conditions['min_size_mb']:
                return False

        # Check file format condition
        if 'file_format' in conditions:
            if data_event.get('file_format') not in conditions['file_format']:
                return False

        # Check frequency limitation
        if rule['trigger_frequency'] != 'immediate' and rule['last_triggered']:
            time_since_last = datetime.now() - rule['last_triggered']
            if rule['trigger_frequency'] == 'hourly' and time_since_last < timedelta(hours=1):
                return False
            elif rule['trigger_frequency'] == 'daily' and time_since_last < timedelta(days=1):
                return False

        return True

    async def _execute_triggered_workflow(self, rule: Dict[str, Any], data_event: Dict[str, Any]):
        """Execute workflow triggered by data event."""
        workflow_config = rule['workflow_config']

        try:
            # Create execution environment
            env = ExecutionEnvironment(
                environment_type=workflow_config.get('environment_type', 'local'),
                resource_limits=workflow_config.get('resource_limits', {})
            )

            # Execute workflow
            execution_id = self.workflow_engine.execute_workflow(
                domain=workflow_config['domain'],
                workflow_name=workflow_config['workflow_name'],
                environment=env,
                dry_run=workflow_config.get('dry_run', False)
            )

            # Record triggered workflow
            triggered_workflow = {
                'rule_id': rule['rule_id'],
                'execution_id': execution_id,
                'trigger_time': datetime.now(),
                'data_event': data_event,
                'workflow_config': workflow_config
            }

            self.triggered_workflows.append(triggered_workflow)

            # Update rule statistics
            rule['last_triggered'] = datetime.now()
            rule['trigger_count'] += 1

            self.logger.info(f"Triggered workflow {execution_id} from rule {rule['rule_id']}")

        except Exception as e:
            self.logger.error(f"Failed to execute triggered workflow: {e}")

class DataManagementEngine:
    """
    Main orchestration class for comprehensive data management.

    Integrates all data management components to provide:
    - Unified data source management
    - Intelligent transfer coordination
    - Event-driven workflow triggering
    - Cost optimization and monitoring
    """

    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)

        # Core components
        self.config_loader = ConfigLoader(config_root)
        self.workflow_engine = DemoWorkflowEngine(config_root)
        self.transfer_manager = DataTransferManager(config_root)
        self.workflow_trigger = DataWorkflowTrigger(self.workflow_engine)

        # Data source registry
        self.data_sources: Dict[str, DataSource] = {}
        self.source_configs: Dict[str, Dict[str, Any]] = {}

        # Event processing
        self.event_queue = asyncio.Queue()
        self.event_processor_active = False

        # Load configuration
        self._load_data_source_configs()

    def _load_data_source_configs(self):
        """Load data source configurations from config files."""
        data_config_file = self.config_root / "data_sources.yaml"

        if data_config_file.exists():
            with open(data_config_file, 'r') as f:
                self.source_configs = yaml.safe_load(f)

            self.logger.info(f"Loaded {len(self.source_configs)} data source configurations")
        else:
            self.logger.warning("No data source configuration file found")

    async def register_data_source(self, source_id: str, source_config: Dict[str, Any]) -> bool:
        """Register a new data source."""
        try:
            source_type = DataSourceType(source_config.get("type", "aws_s3"))

            if source_type == DataSourceType.REAL_TIME_INSTRUMENT:
                # Create instrument configuration
                instrument_config = InstrumentConfig(**source_config['instrument_config'])
                data_source = InstrumentDataSource(source_id, instrument_config)
            else:
                data_source = StaticDataSource(source_id, source_config)

            self.data_sources[source_id] = data_source
            self.source_configs[source_id] = source_config

            self.logger.info(f"Registered data source: {source_id} ({source_type.value})")
            return True

        except Exception as e:
            self.logger.error(f"Failed to register data source {source_id}: {e}")
            return False

    async def discover_available_data(self, source_id: Optional[str] = None,
                                    filters: Optional[Dict[str, Any]] = None) -> Dict[str, List[DataMetadata]]:
        """Discover available data across all or specific data sources."""
        discovered_data = {}

        sources_to_search = [source_id] if source_id else list(self.data_sources.keys())

        for sid in sources_to_search:
            if sid in self.data_sources:
                try:
                    data_list = await self.data_sources[sid].list_available_data(filters)
                    discovered_data[sid] = data_list
                    self.logger.info(f"Discovered {len(data_list)} datasets from {sid}")
                except Exception as e:
                    self.logger.error(f"Error discovering data from {sid}: {e}")
                    discovered_data[sid] = []

        return discovered_data

    async def initiate_smart_transfer(self, source_id: str, data_path: str,
                                    destination: str, priority: TransferPriority = TransferPriority.STANDARD) -> str:
        """Initiate intelligent data transfer with optimization."""
        if source_id not in self.data_sources:
            raise ValueError(f"Unknown data source: {source_id}")

        source = self.data_sources[source_id]
        transfer_id = await self.transfer_manager.schedule_transfer(source, data_path, destination, priority)

        # Generate data event for potential workflow triggering
        metadata = await source.get_data_metadata(data_path)
        data_event = {
            'event_type': 'data_transfer_initiated',
            'source_id': source_id,
            'data_path': data_path,
            'destination': destination,
            'transfer_id': transfer_id,
            'data_type': metadata.data_type,
            'file_format': metadata.file_format,
            'size_bytes': metadata.size_bytes,
            'timestamp': datetime.now()
        }

        await self.event_queue.put(data_event)

        return transfer_id

    async def setup_real_time_processing(self, instrument_id: str,
                                       processing_config: Dict[str, Any]) -> str:
        """Set up real-time data processing for an instrument."""
        if instrument_id not in self.data_sources:
            raise ValueError(f"Unknown instrument: {instrument_id}")

        source = self.data_sources[instrument_id]
        if not isinstance(source, InstrumentDataSource):
            raise ValueError(f"Source {instrument_id} is not a real-time instrument")

        # Set up trigger rule for real-time processing
        trigger_rule = {
            'name': f"Real-time processing for {instrument_id}",
            'description': f"Automatic processing of data from {instrument_id}",
            'data_conditions': {
                'data_type': 'real_time_instrument',
                'min_size_mb': processing_config.get('min_trigger_size_mb', 1)
            },
            'workflow_config': {
                'domain': processing_config['domain'],
                'workflow_name': processing_config['workflow_name'],
                'environment_type': processing_config.get('environment_type', 'local'),
                'dry_run': processing_config.get('dry_run', False)
            },
            'trigger_frequency': processing_config.get('trigger_frequency', 'immediate')
        }

        rule_id = self.workflow_trigger.add_trigger_rule(trigger_rule)

        # Start real-time streaming
        destination = processing_config.get('staging_location', 's3://research-staging/')
        transfer_id = await source.initiate_transfer("real_time_stream", destination, TransferPriority.IMMEDIATE)

        self.logger.info(f"Set up real-time processing for {instrument_id}: rule {rule_id}, transfer {transfer_id}")
        return rule_id

    async def start_event_processing(self):
        """Start background event processing for trigger evaluation."""
        if self.event_processor_active:
            return

        self.event_processor_active = True
        asyncio.create_task(self._process_events())
        self.logger.info("Started data event processing")

    async def _process_events(self):
        """Background task to process data events and evaluate triggers."""
        while self.event_processor_active:
            try:
                # Wait for events with timeout
                event = await asyncio.wait_for(self.event_queue.get(), timeout=1.0)

                # Evaluate trigger rules
                await self.workflow_trigger.evaluate_triggers(event)

            except asyncio.TimeoutError:
                # No events, continue
                continue
            except Exception as e:
                self.logger.error(f"Error processing data event: {e}")

    def stop_event_processing(self):
        """Stop background event processing."""
        self.event_processor_active = False
        self.logger.info("Stopped data event processing")

    def get_data_management_status(self) -> Dict[str, Any]:
        """Get comprehensive status of data management system."""
        # Count datasets across all sources
        total_datasets = sum(len(self.data_sources[sid].metadata_cache) for sid in self.data_sources)

        # Transfer statistics
        transfer_stats = self.transfer_manager.get_transfer_statistics()

        # Trigger statistics
        trigger_stats = {
            'total_rules': len(self.workflow_trigger.trigger_rules),
            'enabled_rules': len([r for r in self.workflow_trigger.trigger_rules if r['enabled']]),
            'total_triggered_workflows': len(self.workflow_trigger.triggered_workflows)
        }

        return {
            'data_sources': {
                'total_registered': len(self.data_sources),
                'by_type': self._count_sources_by_type(),
                'total_datasets_cataloged': total_datasets
            },
            'transfers': transfer_stats,
            'triggers': trigger_stats,
            'event_processing': {
                'active': self.event_processor_active,
                'queue_size': self.event_queue.qsize()
            }
        }

    def _count_sources_by_type(self) -> Dict[str, int]:
        """Count data sources by type."""
        type_counts = {}

        for source in self.data_sources.values():
            source_type = source.source_type.value
            type_counts[source_type] = type_counts.get(source_type, 0) + 1

        return type_counts

def main():
    """CLI interface for data management engine."""
    import argparse

    parser = argparse.ArgumentParser(description="Data Management Engine for AWS Research Wizard")
    parser.add_argument("--register-source", nargs=2, metavar=('ID', 'CONFIG_FILE'),
                       help="Register new data source")
    parser.add_argument("--discover-data", type=str, help="Discover data from source")
    parser.add_argument("--transfer-data", nargs=3, metavar=('SOURCE', 'PATH', 'DEST'),
                       help="Transfer data")
    parser.add_argument("--setup-realtime", nargs=2, metavar=('INSTRUMENT', 'CONFIG_FILE'),
                       help="Setup real-time processing")
    parser.add_argument("--status", action="store_true", help="Show system status")
    parser.add_argument("--config-root", type=str, default="configs", help="Configuration root directory")

    args = parser.parse_args()

    # Setup logging
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

    # Initialize engine
    engine = DataManagementEngine(args.config_root)

    async def run_async_commands():
        if args.register_source:
            source_id, config_file = args.register_source
            with open(config_file, 'r') as f:
                config = yaml.safe_load(f)

            success = await engine.register_data_source(source_id, config)
            print(f"Registration {'successful' if success else 'failed'}")

        elif args.discover_data:
            discovered = await engine.discover_available_data(args.discover_data)

            for source_id, datasets in discovered.items():
                print(f"\n{source_id.upper()}:")
                for dataset in datasets:
                    print(f"  - {dataset.dataset_id}")
                    print(f"    Size: {dataset.size_bytes / (1024**3):.2f} GB")
                    print(f"    Format: {dataset.file_format}")
                    print(f"    Type: {dataset.data_type}")

        elif args.transfer_data:
            source_id, data_path, destination = args.transfer_data

            transfer_id = await engine.initiate_smart_transfer(source_id, data_path, destination)
            print(f"Transfer initiated: {transfer_id}")

        elif args.setup_realtime:
            instrument_id, config_file = args.setup_realtime
            with open(config_file, 'r') as f:
                config = yaml.safe_load(f)

            rule_id = await engine.setup_real_time_processing(instrument_id, config)
            print(f"Real-time processing setup: {rule_id}")

        elif args.status:
            status = engine.get_data_management_status()

            print("Data Management System Status:")
            print(f"  Data Sources: {status['data_sources']['total_registered']}")
            print(f"  Datasets Cataloged: {status['data_sources']['total_datasets_cataloged']}")
            print(f"  Active Transfers: {status['transfers'].get('active_transfers', 0)}")
            print(f"  Trigger Rules: {status['triggers']['total_rules']}")
            print(f"  Event Processing: {'Active' if status['event_processing']['active'] else 'Inactive'}")

        else:
            parser.print_help()

    # Run async commands
    asyncio.run(run_async_commands())

if __name__ == "__main__":
    main()
