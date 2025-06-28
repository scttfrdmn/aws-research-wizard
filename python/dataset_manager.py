#!/usr/bin/env python3
"""
AWS Open Data Dataset Manager for AWS Research Wizard

This module provides comprehensive management of AWS Open Data datasets,
including discovery, cost estimation, subset creation, and demo workflow generation.

Key Features:
- AWS Open Data Registry integration with 18+ datasets totaling 50+ petabytes
- Cost estimation and optimization for data access patterns
- Demo workflow generation with real datasets
- Dataset filtering and subset creation capabilities
- Integration with domain-specific research pack configurations

Classes:
    DatasetInfo: Data class representing individual dataset metadata
    DatasetSubset: Data class for dataset subset information
    DatasetManager: Main dataset management and workflow generation class

The DatasetManager integrates with the AWS Open Data Registry to provide:
- Real dataset discovery across multiple research domains
- Cost-aware data access planning
- Automated demo workflow generation
- Dataset metadata management and validation

Dependencies:
    - boto3: For AWS service integration
    - yaml: For configuration file parsing
    - pathlib: For cross-platform path handling
    - requests: For API interactions
"""

import os
import yaml
import json
import boto3
import logging
from typing import Dict, List, Any, Optional, Tuple
from pathlib import Path
from dataclasses import dataclass, asdict
from datetime import datetime, timedelta
import requests
from urllib.parse import urlparse

@dataclass
class DatasetInfo:
    """
    Comprehensive data class representing AWS Open Data dataset metadata.
    
    This class encapsulates all relevant information about a dataset in the
    AWS Open Data Registry, including size, format, location, and domain-specific
    metadata that enables intelligent dataset selection and cost estimation.
    
    Attributes:
        name (str): Human-readable dataset name
        description (str): Detailed description of dataset contents and purpose
        location (str): S3 bucket/path location for the dataset
        format (str): Primary data format(s) (e.g., VCF, NetCDF, JPEG)
        update_frequency (str): How frequently the dataset is updated
        domains (List[str]): Research domains that can utilize this dataset
        access_pattern (str): Recommended access pattern for cost optimization
        size_tb (Optional[float]): Dataset size in terabytes
        size_gb (Optional[float]): Dataset size in gigabytes (for smaller datasets)
        temporal_coverage (Optional[str]): Time range covered by the dataset
        spatial_resolution (Optional[str]): Geographic/spatial resolution if applicable
        samples (Optional[Any]): Sample count or sample information
        populations (Optional[int]): Number of populations/samples (for genomics datasets)
        images (Optional[int]): Number of images (for computer vision datasets)
        classes (Optional[int]): Number of classification classes (for ML datasets)
        pages (Optional[str]): Number of pages/documents (for text datasets)
        techniques (Optional[int]): Number of analysis techniques supported
        benchmarks (Optional[List[str]]): Standard benchmarks associated with the dataset
    """
    name: str
    description: str
    location: str
    format: str
    update_frequency: str
    domains: List[str]
    access_pattern: str
    size_tb: Optional[float] = None
    size_gb: Optional[float] = None
    temporal_coverage: Optional[str] = None
    spatial_resolution: Optional[str] = None
    samples: Optional[Any] = None
    populations: Optional[int] = None
    images: Optional[int] = None
    classes: Optional[int] = None
    pages: Optional[str] = None
    techniques: Optional[int] = None
    benchmarks: Optional[List[str]] = None
    
@dataclass
class DatasetSubset:
    """
    Data class representing a subset of a larger dataset for demo workflows.
    
    This class is used to create manageable, cost-effective subsets of large
    datasets for demonstration and testing purposes. It helps researchers
    explore dataset capabilities without downloading entire multi-terabyte datasets.
    
    Attributes:
        dataset_name (str): Name of the parent dataset
        subset_name (str): Descriptive name for this subset
        description (str): Description of what this subset contains
        size_gb (float): Size of the subset in gigabytes
        file_count (int): Number of files in the subset
        s3_paths (List[str]): S3 paths to the files in this subset
        estimated_cost (float): Estimated cost in USD to download and process this subset
    
    Example:
        >>> subset = DatasetSubset(
        ...     dataset_name="1000 Genomes Project",
        ...     subset_name="Chromosome 22 Sample",
        ...     description="Sample data from chromosome 22 for 100 individuals",
        ...     size_gb=5.2,
        ...     file_count=100,
        ...     s3_paths=["s3://1000genomes/chr22/sample1.vcf", ...],
        ...     estimated_cost=0.15
        ... )
    """
    dataset_name: str
    subset_name: str
    description: str
    size_gb: float
    file_count: int
    s3_paths: List[str]
    estimated_cost: float
    workflow_type: str

class DatasetManager:
    """Manages AWS Open Data integration and dataset preparation"""
    
    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.logger = logging.getLogger(__name__)
        
        # Load AWS Open Data registry
        self.datasets = self._load_open_data_registry()
        
        # Initialize AWS clients (with error handling for credentials)
        try:
            self.s3_client = boto3.client('s3')
            self.pricing_client = boto3.client('pricing', region_name='us-east-1')
        except Exception as e:
            self.logger.warning(f"AWS credentials not configured: {e}")
            self.s3_client = None
            self.pricing_client = None
    
    def _load_open_data_registry(self) -> Dict[str, Dict[str, DatasetInfo]]:
        """Load AWS Open Data registry configuration"""
        registry_file = self.config_root / "demo_data" / "aws_open_data_registry.yaml"
        
        if not registry_file.exists():
            self.logger.error(f"Open Data registry file not found: {registry_file}")
            return {}
        
        try:
            with open(registry_file, 'r') as f:
                registry_data = yaml.safe_load(f)
            
            datasets = {}
            for category, category_datasets in registry_data.get('datasets', {}).items():
                datasets[category] = {}
                for dataset_id, dataset_config in category_datasets.items():
                    datasets[category][dataset_id] = DatasetInfo(**dataset_config)
            
            self.logger.info(f"Loaded {sum(len(cat) for cat in datasets.values())} datasets from registry")
            return datasets
            
        except Exception as e:
            self.logger.error(f"Failed to load Open Data registry: {e}")
            return {}
    
    def get_datasets_for_domain(self, domain: str) -> List[DatasetInfo]:
        """Get all datasets available for a specific research domain"""
        matching_datasets = []
        
        for category, category_datasets in self.datasets.items():
            for dataset_id, dataset_info in category_datasets.items():
                if domain in dataset_info.domains:
                    matching_datasets.append(dataset_info)
        
        return matching_datasets
    
    def get_dataset_by_name(self, dataset_name: str) -> Optional[DatasetInfo]:
        """Get dataset information by name"""
        for category, category_datasets in self.datasets.items():
            for dataset_id, dataset_info in category_datasets.items():
                if dataset_info.name == dataset_name or dataset_id == dataset_name:
                    return dataset_info
        return None
    
    def create_demo_subset(self, dataset_name: str, subset_config: Dict[str, Any]) -> Optional[DatasetSubset]:
        """Create a demo-sized subset of a large dataset"""
        dataset = self.get_dataset_by_name(dataset_name)
        if not dataset:
            self.logger.error(f"Dataset not found: {dataset_name}")
            return None
        
        # Extract S3 location components
        s3_location = dataset.location.replace('s3://', '')
        bucket_name = s3_location.split('/')[0]
        prefix = '/'.join(s3_location.split('/')[1:]) if '/' in s3_location else ''
        
        subset = DatasetSubset(
            dataset_name=dataset_name,
            subset_name=subset_config.get('name', f"{dataset_name}_demo"),
            description=subset_config.get('description', f"Demo subset of {dataset.name}"),
            size_gb=subset_config.get('size_gb', 1.0),
            file_count=subset_config.get('file_count', 10),
            s3_paths=[],
            estimated_cost=0.0,
            workflow_type=subset_config.get('workflow_type', 'analysis')
        )
        
        # Generate sample S3 paths based on access pattern
        if self.s3_client:
            try:
                # List objects to get actual paths
                response = self.s3_client.list_objects_v2(
                    Bucket=bucket_name,
                    Prefix=prefix,
                    MaxKeys=subset.file_count
                )
                
                subset.s3_paths = [f"s3://{bucket_name}/{obj['Key']}" 
                                 for obj in response.get('Contents', [])]
                
                # Calculate estimated costs
                subset.estimated_cost = self._estimate_data_costs(subset.size_gb)
                
            except Exception as e:
                self.logger.warning(f"Could not list S3 objects for {dataset_name}: {e}")
                # Generate synthetic paths based on access pattern
                subset.s3_paths = self._generate_sample_paths(dataset, subset.file_count)
                subset.estimated_cost = self._estimate_data_costs(subset.size_gb)
        else:
            # Generate synthetic paths when AWS client not available
            subset.s3_paths = self._generate_sample_paths(dataset, subset.file_count)
            subset.estimated_cost = self._estimate_data_costs(subset.size_gb)
        
        return subset
    
    def _generate_sample_paths(self, dataset: DatasetInfo, file_count: int) -> List[str]:
        """Generate sample S3 paths based on dataset access pattern"""
        base_pattern = dataset.access_pattern
        sample_paths = []
        
        # Replace common pattern variables with sample values
        current_date = datetime.now()
        
        for i in range(file_count):
            path = base_pattern
            
            # Replace date patterns
            if '{year}' in path or '{YYYY}' in path:
                path = path.replace('{year}', str(current_date.year))
                path = path.replace('{YYYY}', str(current_date.year))
            
            if '{month}' in path or '{MM}' in path:
                path = path.replace('{month}', f"{current_date.month:02d}")
                path = path.replace('{MM}', f"{current_date.month:02d}")
            
            if '{day}' in path or '{DD}' in path:
                sample_day = (current_date - timedelta(days=i)).day
                path = path.replace('{day}', f"{sample_day:02d}")
                path = path.replace('{DD}', f"{sample_day:02d}")
            
            # Replace other common patterns
            path = path.replace('{HH}', '00')
            path = path.replace('{doy}', f"{current_date.timetuple().tm_yday:03d}")
            path = path.replace('{SRR_ID}', f"SRR{10000000 + i}")
            path = path.replace('{hash}', f"sample_hash_{i:04d}")
            
            sample_paths.append(path)
        
        return sample_paths[:file_count]
    
    def _estimate_data_costs(self, size_gb: float) -> float:
        """Estimate AWS costs for data transfer and processing"""
        # Rough cost estimates (USD)
        s3_storage_cost = size_gb * 0.023  # Standard storage per GB/month
        data_transfer_cost = size_gb * 0.09  # Data transfer out per GB
        compute_estimate = size_gb * 0.05  # Processing compute estimate
        
        return round(s3_storage_cost + data_transfer_cost + compute_estimate, 2)
    
    def _get_dataset_size_gb(self, dataset: DatasetInfo) -> float:
        """Get dataset size in GB, handling both TB and GB specifications"""
        if dataset.size_tb:
            return dataset.size_tb * 1024
        elif dataset.size_gb:
            return dataset.size_gb
        else:
            return 1.0  # Default 1GB if no size specified
    
    def generate_demo_workflows(self, domain: str) -> List[Dict[str, Any]]:
        """Generate demo workflows for a specific domain using real datasets"""
        workflows = []
        datasets = self.get_datasets_for_domain(domain)
        
        workflow_templates = {
            'genomics': [
                {
                    'name': 'Variant Calling Demo',
                    'description': 'GATK best practices pipeline on 1000 Genomes data',
                    'dataset_filter': ['1000 genomes', 'genomes'],
                    'subset_config': {'size_gb': 5.0, 'file_count': 3},
                    'runtime_hours': '2-4',
                    'workflow_type': 'variant_calling'
                },
                {
                    'name': 'RNA-seq Analysis Demo',
                    'description': 'Differential expression analysis using public data',
                    'dataset_filter': ['ncbi', 'sequence read'],
                    'subset_config': {'size_gb': 2.0, 'file_count': 6},
                    'runtime_hours': '1-2',
                    'workflow_type': 'rnaseq'
                }
            ],
            'climate_modeling': [
                {
                    'name': 'Climate Reanalysis Demo',
                    'description': 'Temperature trend analysis using ERA5 data',
                    'dataset_filter': ['era5_reanalysis'],
                    'subset_config': {'size_gb': 3.5, 'file_count': 12},
                    'runtime_hours': '2-6',
                    'workflow_type': 'climate_analysis'
                },
                {
                    'name': 'Weather Forecast Demo',
                    'description': 'Short-term weather prediction using GFS data',
                    'dataset_filter': ['noaa_gfs'],
                    'subset_config': {'size_gb': 2.0, 'file_count': 8},
                    'runtime_hours': '1-3',
                    'workflow_type': 'forecast'
                }
            ],
            'machine_learning': [
                {
                    'name': 'Image Classification Demo',
                    'description': 'Train ResNet on Open Images subset',
                    'dataset_filter': ['open images', 'images'],
                    'subset_config': {'size_gb': 10.0, 'file_count': 1000},
                    'runtime_hours': '3-6',
                    'workflow_type': 'image_classification'
                },
                {
                    'name': 'NLP Model Training Demo',
                    'description': 'Language model training on Common Crawl subset',
                    'dataset_filter': ['common crawl', 'crawl'],
                    'subset_config': {'size_gb': 5.0, 'file_count': 100},
                    'runtime_hours': '4-8',
                    'workflow_type': 'nlp_training'
                }
            ],
            'geospatial_research': [
                {
                    'name': 'Land Cover Classification Demo',
                    'description': 'Supervised classification of Sentinel-2 imagery',
                    'dataset_filter': ['sentinel2_l2a'],
                    'subset_config': {'size_gb': 8.0, 'file_count': 20},
                    'runtime_hours': '2-4',
                    'workflow_type': 'classification'
                },
                {
                    'name': 'Change Detection Demo',
                    'description': 'Multi-temporal analysis using Landsat data',
                    'dataset_filter': ['landsat_collection2'],
                    'subset_config': {'size_gb': 6.0, 'file_count': 15},
                    'runtime_hours': '3-6',
                    'workflow_type': 'change_detection'
                }
            ],
            'agricultural_sciences': [
                {
                    'name': 'Crop Monitoring Demo',
                    'description': 'NDVI time series analysis for crop health',
                    'dataset_filter': ['modis_mcd43a4', 'usda_nass_cdl'],
                    'subset_config': {'size_gb': 4.0, 'file_count': 30},
                    'runtime_hours': '2-4',
                    'workflow_type': 'crop_monitoring'
                }
            ],
            'atmospheric_chemistry': [
                {
                    'name': 'Air Quality Analysis Demo',
                    'description': 'Chemical transport modeling with real emission data',
                    'dataset_filter': ['nasa_merra2'],
                    'subset_config': {'size_gb': 3.0, 'file_count': 20},
                    'runtime_hours': '4-8',
                    'workflow_type': 'air_quality'
                }
            ]
        }
        
        domain_workflows = workflow_templates.get(domain, [])
        
        for workflow_template in domain_workflows:
            # Find matching datasets
            matching_datasets = []
            for dataset in datasets:
                for filter_name in workflow_template['dataset_filter']:
                    if (filter_name.lower() in dataset.name.lower() or 
                        filter_name.lower() in dataset.location.lower() or
                        any(filter_name.lower() in domain.lower() for domain in dataset.domains)):
                        matching_datasets.append(dataset)
                        break
            
            if matching_datasets:
                primary_dataset = matching_datasets[0]
                subset = self.create_demo_subset(primary_dataset.name, 
                                               workflow_template['subset_config'])
                
                if subset:
                    workflow = {
                        'name': workflow_template['name'],
                        'description': workflow_template['description'],
                        'dataset': primary_dataset.name,
                        'dataset_location': primary_dataset.location,
                        'subset_info': asdict(subset),
                        'expected_runtime': workflow_template['runtime_hours'],
                        'cost_estimate': subset.estimated_cost,
                        'workflow_type': workflow_template['workflow_type']
                    }
                    workflows.append(workflow)
        
        return workflows
    
    def list_available_datasets(self, domain: Optional[str] = None) -> List[Dict[str, Any]]:
        """List all available datasets, optionally filtered by domain"""
        dataset_list = []
        
        for category, category_datasets in self.datasets.items():
            for dataset_id, dataset_info in category_datasets.items():
                if domain is None or domain in dataset_info.domains:
                    dataset_summary = {
                        'id': dataset_id,
                        'category': category,
                        'name': dataset_info.name,
                        'description': dataset_info.description,
                        'size_tb': dataset_info.size_tb or (dataset_info.size_gb / 1024 if dataset_info.size_gb else 0),
                        'format': dataset_info.format,
                        'domains': dataset_info.domains,
                        'location': dataset_info.location
                    }
                    dataset_list.append(dataset_summary)
        
        return sorted(dataset_list, key=lambda x: x['name'])
    
    def export_domain_datasets(self, domain: str, output_file: str) -> bool:
        """Export dataset information for a domain to JSON"""
        try:
            datasets = self.get_datasets_for_domain(domain)
            workflows = self.generate_demo_workflows(domain)
            
            export_data = {
                'domain': domain,
                'datasets': [asdict(dataset) for dataset in datasets],
                'demo_workflows': workflows,
                'export_timestamp': datetime.now().isoformat()
            }
            
            with open(output_file, 'w') as f:
                json.dump(export_data, f, indent=2, default=str)
            
            self.logger.info(f"Exported {len(datasets)} datasets and {len(workflows)} workflows to {output_file}")
            return True
            
        except Exception as e:
            self.logger.error(f"Failed to export datasets: {e}")
            return False


def main():
    """CLI interface for dataset management"""
    import argparse
    
    parser = argparse.ArgumentParser(description="AWS Open Data Dataset Manager")
    parser.add_argument("--list-datasets", action="store_true", help="List available datasets")
    parser.add_argument("--domain", type=str, help="Filter by research domain")
    parser.add_argument("--generate-workflows", type=str, help="Generate demo workflows for domain")
    parser.add_argument("--export", type=str, help="Export domain datasets to file")
    parser.add_argument("--output", type=str, help="Output file for export")
    parser.add_argument("--config-root", type=str, default="configs", help="Configuration root directory")
    
    args = parser.parse_args()
    
    # Setup logging
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
    
    # Initialize dataset manager
    manager = DatasetManager(args.config_root)
    
    if args.list_datasets:
        datasets = manager.list_available_datasets(args.domain)
        print(f"Available datasets{' for ' + args.domain if args.domain else ''} ({len(datasets)}):")
        for dataset in datasets:
            print(f"  - {dataset['name']} ({dataset['size_tb']} TB)")
            print(f"    Category: {dataset['category']}")
            print(f"    Domains: {', '.join(dataset['domains'])}")
            print(f"    Location: {dataset['location']}")
            print()
    
    elif args.generate_workflows:
        workflows = manager.generate_demo_workflows(args.generate_workflows)
        print(f"Demo workflows for {args.generate_workflows} ({len(workflows)}):")
        for workflow in workflows:
            print(f"  - {workflow['name']}")
            print(f"    Dataset: {workflow['dataset']}")
            print(f"    Runtime: {workflow['expected_runtime']} hours")
            print(f"    Cost: ${workflow['cost_estimate']}")
            print()
    
    elif args.export:
        if not args.output:
            args.output = f"{args.export}_datasets.json"
        
        success = manager.export_domain_datasets(args.export, args.output)
        if success:
            print(f"Datasets exported to: {args.output}")
        else:
            print(f"Failed to export datasets for: {args.export}")
    
    else:
        parser.print_help()


if __name__ == "__main__":
    main()