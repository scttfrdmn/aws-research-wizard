# AWS Research Wizard - Terminal User Interface (TUI) System

A comprehensive terminal-based interface system for configuring, deploying, and monitoring AWS research environments. Perfect for SSH-based access and remote research computing.

## Overview

The TUI system consists of three integrated components:

1. **Research Wizard TUI** - Interactive domain configuration and deployment
2. **AWS Monitor TUI** - Real-time resource monitoring and cost tracking
3. **Domain Dashboard TUI** - Configurable domain-specific monitoring

## Components

### 1. Research Wizard TUI (`tui_research_wizard.py`)

Interactive configuration and deployment interface for research environments.

**Features:**
- Browse and select from 15 research domains
- Interactive instance type selection with cost comparison
- Real-time cost calculator with optimization suggestions
- Spot instance savings calculator (70-90% savings)
- Configuration validation and deployment planning
- SSH-friendly design for remote access

**Usage:**
```bash
# Full interactive interface
python tui_research_wizard.py

# Simple TUI mode (no Textual dependency)
python tui_research_wizard.py --simple

# Custom configuration directory
python tui_research_wizard.py --config /path/to/configs
```

**Key Features:**
- **Domain Selection**: Visual table with cost estimates and user counts
- **Instance Comparison**: Side-by-side cost analysis (hourly/monthly/annual)
- **Cost Optimization**: Intelligent suggestions for Reserved/Spot instances
- **Deployment Planning**: Dry-run validation with resource requirements

### 2. AWS Monitor TUI (`tui_aws_monitor.py`)

Real-time AWS resource monitoring with cost tracking and alerting.

**Features:**
- Live EC2 instance monitoring (CPU, memory, network, disk I/O)
- Real-time cost tracking with daily/monthly projections
- Configurable alerts for resource thresholds
- Cost breakdown by AWS service
- Performance optimization recommendations
- Alert management system

**Usage:**
```bash
# Launch monitoring dashboard
python tui_aws_monitor.py

# Specify region and refresh interval
python tui_aws_monitor.py --region us-west-2 --refresh 15

# Monitor specific region
python tui_aws_monitor.py --region eu-west-1
```

**Dashboard Layout:**
```
┌─ AWS Resource Monitor ─ 2024-01-15 14:30:22 ─ 🟢 Connected ─ 3 instances ─┐
├─ EC2 Instances ──────────────┬─ Cost Tracking ────────────────────────┐
│ Instance ID      CPU   Mem   │ Service      Monthly    Daily         │
│ i-abc123  45%   67%    running│ EC2          $245.67    $8.19         │
│ i-def456  78%   82%    running│ S3           $67.34     $2.24         │
├─ CPU Usage Trends ───────────┼─ Alert Status ─────────────────────────┤
│ i-abc123: ▁▃▅▇█▆▄▂ [45.2%]   │ 🟢 High CPU Usage: 45.2/80.0         │
│ i-def456: ▃▅▇█▇▅▃▁ [78.1%]   │ 🔴 Daily Cost: $52.4/50.0            │
└───────────────────────────────┴─────────────────────────────────────────┘
```

### 3. Domain Dashboard TUI (`tui_domain_dashboard.py`)

Configurable dashboards tailored to specific research domains.

**Features:**
- Domain-specific metric layouts and visualizations
- Real-time workflow monitoring
- Research-specific data analytics
- Customizable alerts and notifications
- Export capabilities for research reporting
- YAML-based configuration system

**Usage:**
```bash
# Interactive domain selection
python tui_domain_dashboard.py

# Launch specific domain dashboard
python tui_domain_dashboard.py --domain genomics

# Custom configuration directory
python tui_domain_dashboard.py --config /path/to/configs
```

**Supported Domains:**
- **Genomics**: Sample processing, variant calling, quality metrics
- **Climate Modeling**: Model runs, forecast accuracy, temperature analysis
- **Neuroscience**: Brain imaging pipelines, connectivity analysis
- **Materials Science**: DFT calculations, MD simulations, energy convergence
- **Astronomy**: Survey data processing, simulation status
- **And more**: All 15 research domains supported

## Domain-Specific Features

### Genomics Dashboard
```
├─ Sample Processing Pipeline ─┬─ Quality Control Metrics ──────────────┐
│ ✅ Quality Control    [100%] │ Sample ID    Status      Quality  Vars │
│ 🔄 Alignment         [85%]  │ SAMPLE_001   ✓ Complete   98.5   1.2M  │
│ 🔄 Variant Calling   [45%]  │ SAMPLE_002   🔄 Process   97.2   1.1M  │
│ ⏳ Annotation        [0%]   │ SAMPLE_003   ⏳ Queued    -      -     │
```

### Climate Modeling Dashboard
```
├─ Climate Model Execution ────┬─ Forecast Accuracy Metrics ────────────┐
│ ✅ Data Preprocessing [100%] │ Metric           Current   Target      │
│ ✅ Model Init         [100%] │ Temperature RMSE   1.2°C   <2.0°C  ✓  │
│ 🔄 Simulation Run    [65%]  │ Precipitation Bias 5.3%    <10%    ✓  │
│ ⏳ Post-processing   [0%]   │ Wind Speed RMSE    2.8m/s  <3.0m/s ✓  │
```

### Neuroscience Dashboard
```
├─ Brain Imaging Pipeline ─────┬─ Subject Processing Status ────────────┐
│ ✅ Structural Process [100%] │ Subject      fMRI    DTI    Status     │
│ 🔄 Functional Analysis[75%]  │ SUB_001      ✓       ✓      Complete   │
│ 🔄 Connectivity      [30%]  │ SUB_002      🔄      ✓      Processing  │
│ ⏳ Statistical       [0%]   │ SUB_003      ⏳      ⏳     Queued      │
```

## Installation & Dependencies

### Required Dependencies
```bash
pip install rich boto3 pyyaml asyncio
```

### Optional Dependencies (Enhanced Features)
```bash
# For advanced TUI features
pip install textual

# For system monitoring
pip install psutil

# For enhanced chart capabilities
pip install matplotlib plotext
```

### AWS Configuration
```bash
# Configure AWS credentials
aws configure

# Or set environment variables
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1
```

## Configuration

### Domain Dashboard Configuration

Create custom dashboard configurations in `configs/dashboards/`:

```yaml
# configs/dashboards/custom_research.yaml
title: "Custom Research Dashboard"
description: "Tailored monitoring for custom research workflows"
refresh_interval: 10
auto_scroll: true
color_scheme: "default"

widgets:
  - name: "workflow_progress"
    type: "progress"
    title: "Workflow Progress"
    position: "top-left"
    data_source: "workflow"
    refresh_interval: 5

  - name: "system_metrics"
    type: "table"
    title: "System Metrics"
    position: "top-right"
    data_source: "aws"

  - name: "job_status"
    type: "status"
    title: "Job Status"
    position: "bottom-left"
    data_source: "workflow"

  - name: "resource_usage"
    type: "chart"
    title: "Resource Usage"
    position: "bottom-right"
    data_source: "aws"
```

### Widget Types

- **progress**: Progress bars for workflows and pipelines
- **table**: Tabular data display for metrics and status
- **status**: Key-value status information
- **chart**: ASCII charts and visualizations

### Data Sources

- **workflow**: Workflow engine integration
- **aws**: AWS resource metrics
- **custom**: Domain-specific data collectors
- **static**: Configuration-based static data

## Advanced Features

### Cost Optimization Interface

Real-time cost comparison with optimization suggestions:

```
Cost Comparison Table:
Instance Type    Hourly    Monthly    Annual    Savings w/ Spot
c6i.2xlarge     $0.34     $245       $2,940    $2,058 (70%)
r6i.4xlarge     $1.02     $734       $8,813    $6,169 (70%)
p4d.24xlarge    $32.77    $23,594    $283,122  $198,185 (70%)

💡 Optimization Suggestions:
• Consider Spot instances for 70% savings
• Reserved instances can save 30-60%
• Use S3 Intelligent Tiering for storage
• Enable detailed monitoring for optimization
```

### Alert Management System

Configurable alerts for resource monitoring:

```python
# Example alert configuration
alerts = [
    AlertRule("High CPU Usage", "cpu_utilization", 80.0, "greater", True),
    AlertRule("High Memory Usage", "memory_utilization", 85.0, "greater", True),
    AlertRule("Daily Cost Limit", "daily_cost", 50.0, "greater", True),
    AlertRule("Monthly Cost Projection", "projected_monthly", 1000.0, "greater", True),
]
```

### Keyboard Controls

All TUI interfaces support keyboard navigation:

- **q**: Quit application
- **r**: Refresh data
- **h**: Help/Documentation
- **p**: Pause/Resume auto-refresh
- **c**: Configuration menu
- **Arrow Keys**: Navigate menus
- **Enter**: Select items
- **Ctrl+C**: Stop monitoring/return to menu

## SSH and Remote Access

The TUI system is optimized for SSH connections and remote access:

### Screen/Tmux Integration
```bash
# Start persistent monitoring session
screen -S aws-monitor
python tui_aws_monitor.py

# Detach: Ctrl+A, D
# Reattach: screen -r aws-monitor

# With tmux
tmux new-session -s research-dashboard
python tui_domain_dashboard.py --domain genomics
```

### Low Bandwidth Optimization
- ASCII-only interface for minimal bandwidth usage
- Configurable refresh intervals
- Efficient data updates
- Text-based visualizations

## Use Cases

### Research Computing Centers
- Monitor shared HPC resources
- Track multi-user cost allocation
- Real-time workflow monitoring
- Resource optimization recommendations

### Individual Researchers
- Personal AWS resource monitoring
- Domain-specific workflow tracking
- Cost control and optimization
- Remote access to cloud resources

### Multi-Institutional Collaborations
- Shared resource monitoring
- Collaborative workflow tracking
- Cost transparency across institutions
- Standardized monitoring interfaces

## Troubleshooting

### Common Issues

**AWS Credentials Not Found:**
```bash
# Check AWS configuration
aws sts get-caller-identity

# Configure if needed
aws configure
```

**Textual Import Error:**
```bash
# Install optional dependency
pip install textual

# Or use simple mode
python tui_research_wizard.py --simple
```

**Permission Errors:**
```bash
# Verify AWS permissions
aws ec2 describe-instances
aws cloudwatch get-metric-statistics
aws ce get-cost-and-usage
```

### Performance Optimization

**For Large Numbers of Instances:**
- Increase refresh intervals
- Use filtering options
- Consider regional monitoring

**For Remote/Slow Connections:**
- Use simple TUI mode
- Increase refresh intervals
- Minimize concurrent monitoring

## Integration with Other Components

### Workflow Engine Integration
```python
# Integration with demo workflow engine
from demo_workflow_engine import DemoWorkflowEngine

workflow_engine = DemoWorkflowEngine()
workflow_status = workflow_engine.get_execution_status()
```

### Configuration System Integration
```python
# Integration with config loader
from config_loader import ConfigLoader

config_loader = ConfigLoader("configs")
domain_config = config_loader.load_domain_pack("genomics")
```

### S3 Transfer Integration
```python
# Integration with S3 transfer optimizer
from s3_transfer_optimizer import S3TransferOptimizer

s3_optimizer = S3TransferOptimizer()
transfer_strategy = s3_optimizer.analyze_transfer_requirements(...)
```

## Future Enhancements

### Planned Features
- WebSocket-based real-time updates
- Enhanced chart and visualization capabilities
- Multi-region monitoring aggregation
- Custom plugin system for data sources
- Integration with external monitoring systems

### Community Contributions
- Custom domain dashboard templates
- Additional widget types
- Enhanced visualization capabilities
- Integration with popular research tools

## Support and Documentation

### Getting Help
- Run with `--help` flag for command-line options
- Press `h` in any TUI for keyboard shortcuts
- Check logs for detailed error information

### Reporting Issues
- Include TUI component and version information
- Provide AWS region and resource details
- Include error messages and logs

### Contributing
- Follow existing code structure and patterns
- Test with multiple research domains
- Ensure SSH/remote compatibility
- Document new features and configurations

## License

This TUI system is part of the AWS Research Wizard project and follows the same licensing terms.
