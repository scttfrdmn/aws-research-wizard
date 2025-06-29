# AWS Research Wizard - Example Configurations

This directory contains comprehensive example configurations for common research scenarios. Each configuration demonstrates best practices for data management in specific research domains.

## Available Examples

### 🧬 Genomics Research (`genomics-project.yaml`)
**Use Case**: Whole genome sequencing data management  
**Key Features**:
- FASTQ/BAM/VCF file handling with domain-specific optimizations
- Small file bundling for efficiency
- Research-grade encryption and compliance
- Cost-optimized storage tiering
- Genomics-specific quality checks and validation

**Typical Data Volumes**: 2.5TB raw sequencing, 1.8TB alignments, 500GB variants  
**Monthly Budget**: $500

### 🌍 Climate Science (`climate-research-project.yaml`)
**Use Case**: Climate modeling and satellite observations  
**Key Features**:
- NetCDF and GRIB file optimization
- Time-series data organization
- Large file handling for model outputs
- Public data compliance (FAIR principles)
- Temporal partitioning and indexing

**Typical Data Volumes**: 15TB satellite data, 25TB model outputs, 45TB reanalysis  
**Monthly Budget**: $2000

### 🤖 Machine Learning (`machine-learning-project.yaml`)
**Use Case**: Deep learning training and model management  
**Key Features**:
- High-volume training dataset management
- Model artifact versioning and registry
- GPU-optimized data loading patterns
- MLflow experiment tracking integration
- Preprocessing cache management

**Typical Data Volumes**: 50TB training data, 8TB model artifacts, 12TB cache  
**Monthly Budget**: $1500

### 🔭 Astronomy (`astronomy-project.yaml`)
**Use Case**: Large-scale astronomical surveys  
**Key Features**:
- FITS file handling and compression
- Observation metadata extraction
- Public data release workflows
- Long-term preservation strategies
- Observatory-specific quality controls

**Typical Data Volumes**: 75TB raw observations, 60TB calibrated data, 25TB mosaics  
**Monthly Budget**: $3000

### 🗺️ Geospatial Research (`geospatial-project.yaml`)
**Use Case**: GIS data processing and mapping  
**Key Features**:
- Spatial indexing and coordinate system handling
- LiDAR point cloud optimization
- Map tile generation and caching
- Vector data topology validation
- Web service integration (WMS/WFS)

**Typical Data Volumes**: 40TB imagery, 100TB LiDAR, 500GB vector data  
**Monthly Budget**: $1200

### ⚗️ Computational Chemistry (`chemistry-project.yaml`)
**Use Case**: Molecular simulations and quantum calculations  
**Key Features**:
- Molecular structure validation
- MD trajectory compression
- Quantum chemistry calculation management
- Chemical database integration
- Property extraction and analysis

**Typical Data Volumes**: 2TB structures, 50TB simulations, 10TB calculations  
**Monthly Budget**: $800

## How to Use These Examples

### 1. Choose Your Domain
Select the example configuration that best matches your research domain:

```bash
# Copy the relevant example
cp examples/genomics-project.yaml my-project.yaml
```

### 2. Customize for Your Needs
Edit the configuration file to match your specific requirements:

- **Project Information**: Update project name, owner, budget, and tags
- **Data Profiles**: Modify paths, file types, and data characteristics
- **Destinations**: Configure your S3 buckets and storage classes
- **Workflows**: Adjust processing steps and schedules
- **Optimization**: Tune performance and cost settings

### 3. Validate Configuration
Use the validation command to check your configuration:

```bash
aws-research-wizard data validate my-project.yaml --verbose
```

### 4. Test with Dry-Run
Test workflows without executing transfers:

```bash
aws-research-wizard data workflow run --config my-project.yaml --workflow upload_data --dry-run
```

### 5. Execute Workflows
Run your data movement workflows:

```bash
aws-research-wizard data workflow run --config my-project.yaml --workflow upload_data
```

## Configuration Structure

Each example follows the same comprehensive structure:

```yaml
project:           # Project metadata and settings
data_profiles:     # Source data definitions
destinations:      # Target storage configurations  
workflows:         # Data movement and processing workflows
settings:          # Global system settings
optimization:      # Performance and cost optimization
monitoring:        # Alerting and metrics collection
compliance:        # Domain-specific compliance requirements
```

## Domain-Specific Optimizations

Each example includes optimizations tailored to the research domain:

- **File Type Handling**: Domain-specific file format validation and processing
- **Transfer Strategies**: Optimized concurrency and part sizes for typical data patterns
- **Storage Tiering**: Appropriate lifecycle policies for data access patterns
- **Compression**: Domain-appropriate compression algorithms
- **Validation**: Research-specific quality checks and integrity validation
- **Metadata**: Extraction and preservation of domain-relevant metadata

## Best Practices Demonstrated

### Cost Optimization
- Intelligent storage class selection
- Lifecycle policies for automatic tiering
- Small file bundling where appropriate
- Compression strategies for different data types

### Performance Optimization
- Optimal concurrency settings for data patterns
- Part size tuning for different file sizes
- Engine selection based on data characteristics
- Network optimization settings

### Reliability
- Comprehensive error handling and retry logic
- Data integrity verification
- Backup and versioning strategies
- Monitoring and alerting configuration

### Compliance
- Domain-specific metadata standards
- Data retention and lifecycle policies
- Access control and encryption settings
- Audit logging and traceability

## Customization Guidelines

### Paths and Locations
Update all file paths to match your local environment:
```yaml
data_profiles:
  my_data:
    path: "/path/to/your/data"  # Update this
```

### AWS Configuration
Configure your AWS settings:
```yaml
destinations:
  my_storage:
    uri: "s3://your-bucket/prefix/"  # Your S3 bucket
    region: "your-preferred-region"  # Your AWS region
```

### Budget and Alerts
Set appropriate budget limits and alert thresholds:
```yaml
optimization:
  cost_optimization:
    budget_limit: "$your-budget/month"
    cost_alerts: ["70%", "85%", "95%"]
```

### Workflows
Customize workflows for your specific needs:
- Add or remove preprocessing/postprocessing steps
- Adjust schedules for automated workflows
- Modify validation criteria
- Configure custom parameters

## Getting Help

### Generate Custom Configuration
Use the generate command to create configurations automatically:
```bash
aws-research-wizard data generate /path/to/data --domain genomics --output my-config.yaml
```

### Analyze Your Data
Understand your data patterns before configuration:
```bash
aws-research-wizard data analyze /path/to/data --verbose
```

### Diagnose Issues
Troubleshoot problems with the diagnostic tools:
```bash
aws-research-wizard data diagnose --verbose
```

### Recovery
Recover from failed workflows:
```bash
aws-research-wizard data recover --list
aws-research-wizard data recover workflow_id --interactive
```

## Contributing

To contribute additional domain examples or improvements:

1. Follow the existing configuration structure
2. Include comprehensive documentation
3. Test thoroughly with representative data
4. Ensure domain-specific optimizations are included
5. Add appropriate compliance and safety considerations

## Support

For questions or issues with these examples:

- Check the main documentation
- Use the diagnostic tools for troubleshooting
- Review the validation output for configuration issues
- Consult domain-specific best practices in your field

---

These examples provide production-ready starting points for research data management across diverse scientific domains. Customize them to match your specific requirements and infrastructure.