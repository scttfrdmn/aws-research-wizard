# AWS Open Data Access Guide

This guide provides information on accessing and using AWS Open Data datasets for research workflows.

## Overview

The AWS Research Wizard integrates with the AWS Open Data Registry to provide access to petabytes of publicly available datasets. These datasets are hosted on AWS and can be accessed without egress charges when used within AWS.

## Available Datasets

### Agricultural Sciences

#### USDA NASS Cropland Data Layer
- **Description**: Crop-specific land cover classification
- **Location**: `s3://usda-nass-cdl/`
- **Size**: 45 TB
- **Format**: GeoTIFF
- **Update Frequency**: Annual
- **Temporal Coverage**: 2008-present
- **Spatial Resolution**: 30m

### Benchmarking Performance

#### MLPerf Training Datasets
- **Description**: Standard machine learning benchmark datasets
- **Location**: `s3://mlperf-datasets/`
- **Size**: 12 TB
- **Format**: TFRecord, Parquet
- **Update Frequency**: Versioned

### Climate Modeling

#### ERA5 Reanalysis Data
- **Description**: ECMWF ERA5 atmospheric reanalysis data
- **Location**: `s3://era5-pds/`
- **Size**: 2500 TB
- **Format**: NetCDF
- **Update Frequency**: Daily
- **Temporal Coverage**: 1979-present
- **Spatial Resolution**: 0.25° x 0.25°

#### NASA Global Precipitation Measurement
- **Description**: Global precipitation satellite observations
- **Location**: `s3://gpm-pds/`
- **Size**: 280 TB
- **Format**: HDF5, NetCDF
- **Update Frequency**: Real-time
- **Temporal Coverage**: 2014-present
- **Spatial Resolution**: 0.1° x 0.1°

#### NASA MERRA-2 Reanalysis
- **Description**: Modern-Era Retrospective analysis for Research and Applications
- **Location**: `s3://nasa-merra2/`
- **Size**: 850 TB
- **Format**: NetCDF4
- **Update Frequency**: Monthly
- **Temporal Coverage**: 1980-present
- **Spatial Resolution**: 0.5° x 0.625°

#### NOAA Global Forecast System
- **Description**: Global numerical weather prediction data
- **Location**: `s3://noaa-gfs-bdp-pds/`
- **Size**: 450 TB
- **Format**: GRIB2
- **Update Frequency**: 4x daily
- **Temporal Coverage**: 2015-present
- **Spatial Resolution**: 0.25° x 0.25°

### Genomics

#### 1000 Genomes Project
- **Description**: International genome sequencing consortium data
- **Location**: `s3://1000genomes/`
- **Size**: 260 TB
- **Format**: VCF, BAM, CRAM
- **Update Frequency**: Static

#### Genome Aggregation Database
- **Description**: Population genomics variant database
- **Location**: `s3://gnomad-public-us-east-1/`
- **Size**: 45 TB
- **Format**: VCF, Parquet
- **Update Frequency**: Versioned releases

#### NCBI Sequence Read Archive
- **Description**: Public sequencing data repository
- **Location**: `s3://sra-pub-run-odp/`
- **Size**: 15000 TB
- **Format**: SRA, FASTQ
- **Update Frequency**: Daily

### Geospatial Research

#### Landsat Collection 2
- **Description**: Landsat satellite imagery archive
- **Location**: `s3://usgs-landsat/`
- **Size**: 1800 TB
- **Format**: Cloud-Optimized GeoTIFF
- **Update Frequency**: Daily
- **Temporal Coverage**: 1972-present
- **Spatial Resolution**: 15-100m

#### MODIS Nadir BRDF-Adjusted Reflectance
- **Description**: MODIS daily surface reflectance product
- **Location**: `s3://modis-pds/`
- **Size**: 120 TB
- **Format**: HDF4, COG
- **Update Frequency**: Daily
- **Temporal Coverage**: 2000-present
- **Spatial Resolution**: 500m

#### Sentinel-2 Cloud-Optimized GeoTIFFs
- **Description**: ESA Sentinel-2 Level-2A atmospheric corrected imagery
- **Location**: `s3://sentinel-cogs/`
- **Size**: 3200 TB
- **Format**: Cloud-Optimized GeoTIFF
- **Update Frequency**: Daily
- **Temporal Coverage**: 2015-present
- **Spatial Resolution**: 10-60m

### Machine Learning

#### Common Crawl Web Corpus
- **Description**: Petabyte-scale web crawl data for NLP
- **Location**: `s3://commoncrawl/`
- **Size**: 380 TB
- **Format**: WARC, WET, WAT
- **Update Frequency**: Monthly

#### ImageNet Object Localization Challenge
- **Description**: Large-scale image classification dataset
- **Location**: `s3://imagenet/`
- **Size**: 150 TB
- **Format**: JPEG, XML
- **Update Frequency**: Static

#### Open Images Dataset V7
- **Description**: Annotated image dataset for computer vision
- **Location**: `s3://open-images-dataset/`
- **Size**: 560 TB
- **Format**: JPEG, CSV
- **Update Frequency**: Versioned releases

## Data Access Patterns

### Cost Optimization
- Use S3 Intelligent Tiering for datasets larger than 128KB
- Access data from the same AWS region to minimize transfer costs
- Some datasets use requester-pays model - check before large downloads

### Performance Optimization
- Use multipart downloads for files larger than 100MB
- Parallelize access across S3 prefixes for maximum throughput
- Implement local caching for frequently accessed data
- Use spot instances for batch processing to reduce costs

### Security Considerations
- Most datasets are publicly accessible
- Some datasets may require data use agreements
- All data is encrypted in transit and at rest
- Follow your organization's data governance policies

## Example Access Code

### Python with Boto3
```python
import boto3

# Initialize S3 client
s3 = boto3.client('s3')

# List objects in a dataset
response = s3.list_objects_v2(
    Bucket='dataset-bucket',
    Prefix='path/to/data/',
    MaxKeys=100
)

# Download a file
s3.download_file(
    Bucket='dataset-bucket',
    Key='path/to/file.nc',
    Filename='local_file.nc'
)
```

### AWS CLI
```bash
# List dataset contents
aws s3 ls s3://dataset-bucket/path/to/data/ --no-sign-request

# Download a file
aws s3 cp s3://dataset-bucket/path/to/file.nc local_file.nc --no-sign-request

# Sync a directory
aws s3 sync s3://dataset-bucket/path/to/data/ ./local_data/ --no-sign-request
```

## Cost Estimation

Use the AWS Simple Monthly Calculator to estimate costs:
- S3 storage: ~$0.023 per GB/month
- Data transfer: Free within same region, $0.09 per GB out to internet
- Compute: Varies by instance type and usage

## Support and Resources

- [AWS Open Data Registry](https://registry.opendata.aws/)
- [AWS Data Documentation](https://docs.aws.amazon.com/datasets/)
- [AWS Research Credits](https://aws.amazon.com/research-credits/)
- [AWS Cost Calculator](https://calculator.aws/)

## Getting Started

1. Set up AWS credentials or use IAM roles
2. Choose appropriate instance types for your workload
3. Test with small data subsets first
4. Implement proper error handling and retry logic
5. Monitor costs and usage with CloudWatch

For specific workflow examples, see the demo workflows in each research pack configuration.
