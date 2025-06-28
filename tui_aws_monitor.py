#!/usr/bin/env python3
"""
AWS Resource Monitor Terminal User Interface

This module provides a comprehensive real-time monitoring interface for AWS resources
used in research computing environments. Features include:

1. Real-time EC2 instance monitoring with CPU, memory, and network metrics
2. S3 storage usage and cost tracking
3. Cost analysis and budget monitoring
4. Research workflow status tracking
5. Alert system for resource thresholds
6. Performance optimization recommendations

Key Features:
- Live resource monitoring with auto-refresh
- Interactive charts and graphs in terminal
- Cost breakdown by service and project
- Alert system for resource limits
- SSH-friendly design for remote monitoring
- Export capabilities for reporting

TUI Components:
- Real-time metrics dashboard
- Cost tracking with projections
- Instance management interface
- Storage analytics
- Network monitoring
- Alert management

Dependencies:
    - rich: Advanced terminal formatting
    - textual: Modern TUI framework
    - boto3: AWS SDK for resource monitoring
    - matplotlib: ASCII chart generation
    - psutil: System metrics (when running on EC2)
"""

import os
import sys
import asyncio
import json
import time
import threading
from datetime import datetime, timedelta
from typing import Dict, List, Any, Optional, Tuple
from dataclasses import dataclass
from collections import defaultdict, deque
import logging

# Rich terminal components
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.progress import Progress, SpinnerColumn, TextColumn
from rich.layout import Layout
from rich.text import Text
from rich.align import Align
from rich.columns import Columns
from rich.live import Live
from rich.bar import Bar
from rich import box
from rich.tree import Tree

# AWS SDK
import boto3
from botocore.exceptions import ClientError, NoCredentialsError

# Optional dependencies for enhanced features
try:
    import psutil
    PSUTIL_AVAILABLE = True
except ImportError:
    PSUTIL_AVAILABLE = False

try:
    from textual.app import App, ComposeResult
    from textual.containers import Container, Horizontal, Vertical, ScrollableContainer
    from textual.widgets import (
        Button, Static, DataTable, Log, ProgressBar, 
        Header, Footer, Label, Switch, ListView, ListItem,
        Sparkline, Chart, PlotextPlot
    )
    from textual.reactive import reactive
    from textual.binding import Binding
    from textual.timer import Timer
    TEXTUAL_AVAILABLE = True
except ImportError:
    TEXTUAL_AVAILABLE = False

@dataclass
class ResourceMetrics:
    """Container for AWS resource metrics."""
    timestamp: datetime
    instance_id: str
    instance_type: str
    cpu_utilization: float
    memory_utilization: float
    network_in: float
    network_out: float
    disk_read: float
    disk_write: float
    status: str

@dataclass
class CostMetrics:
    """Container for cost tracking metrics."""
    service: str
    daily_cost: float
    monthly_cost: float
    projected_monthly: float
    currency: str = "USD"

@dataclass
class AlertRule:
    """Alert rule configuration."""
    name: str
    metric: str
    threshold: float
    comparison: str  # "greater", "less", "equal"
    enabled: bool
    last_triggered: Optional[datetime] = None

class AWSResourceMonitor:
    """
    Real-time AWS resource monitoring and cost tracking.
    
    Provides comprehensive monitoring of AWS resources including:
    - EC2 instances with detailed metrics
    - S3 storage usage and costs
    - Overall cost tracking and projections
    - Performance alerts and recommendations
    """
    
    def __init__(self, region: str = "us-east-1", refresh_interval: int = 30):
        self.region = region
        self.refresh_interval = refresh_interval
        self.console = Console()
        
        # Initialize AWS clients
        try:
            self.ec2_client = boto3.client('ec2', region_name=region)
            self.cloudwatch_client = boto3.client('cloudwatch', region_name=region)
            self.s3_client = boto3.client('s3')
            self.ce_client = boto3.client('ce', region_name='us-east-1')  # Cost Explorer only in us-east-1
            self.aws_available = True
        except (NoCredentialsError, ClientError) as e:
            self.console.print(f"[red]AWS credentials not available: {e}[/red]")
            self.aws_available = False
        
        # Monitoring state
        self.instances_cache = {}
        self.metrics_history: Dict[str, deque] = defaultdict(lambda: deque(maxlen=100))
        self.cost_cache = {}
        self.alerts: List[AlertRule] = []
        self.monitoring_active = False
        
        # Setup default alerts
        self._setup_default_alerts()
    
    def _setup_default_alerts(self):
        """Setup default monitoring alerts."""
        default_alerts = [
            AlertRule("High CPU Usage", "cpu_utilization", 80.0, "greater", True),
            AlertRule("High Memory Usage", "memory_utilization", 85.0, "greater", True),
            AlertRule("Daily Cost Limit", "daily_cost", 50.0, "greater", True),
            AlertRule("Monthly Cost Projection", "projected_monthly", 1000.0, "greater", True),
        ]
        self.alerts.extend(default_alerts)
    
    async def start_monitoring(self):
        """Start the monitoring loop."""
        if not self.aws_available:
            self.console.print("[red]Cannot start monitoring: AWS not available[/red]")
            return
        
        self.monitoring_active = True
        
        # Start monitoring tasks
        tasks = [
            asyncio.create_task(self._monitor_instances()),
            asyncio.create_task(self._monitor_costs()),
            asyncio.create_task(self._check_alerts())
        ]
        
        try:
            await asyncio.gather(*tasks)
        except KeyboardInterrupt:
            self.monitoring_active = False
            self.console.print("[yellow]Monitoring stopped by user[/yellow]")
    
    async def _monitor_instances(self):
        """Monitor EC2 instances continuously."""
        while self.monitoring_active:
            try:
                await self._fetch_instance_metrics()
                await asyncio.sleep(self.refresh_interval)
            except Exception as e:
                self.console.print(f"[red]Error monitoring instances: {e}[/red]")
                await asyncio.sleep(60)  # Wait longer on error
    
    async def _monitor_costs(self):
        """Monitor costs continuously."""
        while self.monitoring_active:
            try:
                await self._fetch_cost_metrics()
                await asyncio.sleep(300)  # Update costs every 5 minutes
            except Exception as e:
                self.console.print(f"[red]Error monitoring costs: {e}[/red]")
                await asyncio.sleep(300)
    
    async def _check_alerts(self):
        """Check alert conditions continuously."""
        while self.monitoring_active:
            try:
                await self._evaluate_alerts()
                await asyncio.sleep(60)  # Check alerts every minute
            except Exception as e:
                self.console.print(f"[red]Error checking alerts: {e}[/red]")
                await asyncio.sleep(60)
    
    async def _fetch_instance_metrics(self):
        """Fetch current instance metrics from AWS."""
        if not self.aws_available:
            return
        
        try:
            # Get all running instances
            response = self.ec2_client.describe_instances(
                Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]
            )
            
            current_time = datetime.utcnow()
            
            for reservation in response['Reservations']:
                for instance in reservation['Instances']:
                    instance_id = instance['InstanceId']
                    instance_type = instance['InstanceType']
                    
                    # Fetch CloudWatch metrics
                    metrics = await self._get_cloudwatch_metrics(instance_id, current_time)
                    
                    resource_metric = ResourceMetrics(
                        timestamp=current_time,
                        instance_id=instance_id,
                        instance_type=instance_type,
                        cpu_utilization=metrics.get('CPUUtilization', 0.0),
                        memory_utilization=metrics.get('MemoryUtilization', 0.0),
                        network_in=metrics.get('NetworkIn', 0.0),
                        network_out=metrics.get('NetworkOut', 0.0),
                        disk_read=metrics.get('DiskReadBytes', 0.0),
                        disk_write=metrics.get('DiskWriteBytes', 0.0),
                        status=instance['State']['Name']
                    )
                    
                    # Store in cache and history
                    self.instances_cache[instance_id] = resource_metric
                    self.metrics_history[instance_id].append(resource_metric)
        
        except ClientError as e:
            self.console.print(f"[red]Error fetching instance metrics: {e}[/red]")
    
    async def _get_cloudwatch_metrics(self, instance_id: str, end_time: datetime) -> Dict[str, float]:
        """Get CloudWatch metrics for an instance."""
        start_time = end_time - timedelta(minutes=5)
        metrics = {}
        
        # Define metrics to fetch
        metric_queries = [
            ('CPUUtilization', 'AWS/EC2'),
            ('NetworkIn', 'AWS/EC2'),
            ('NetworkOut', 'AWS/EC2'),
            ('DiskReadBytes', 'AWS/EBS'),
            ('DiskWriteBytes', 'AWS/EBS')
        ]
        
        for metric_name, namespace in metric_queries:
            try:
                response = self.cloudwatch_client.get_metric_statistics(
                    Namespace=namespace,
                    MetricName=metric_name,
                    Dimensions=[
                        {
                            'Name': 'InstanceId',
                            'Value': instance_id
                        }
                    ],
                    StartTime=start_time,
                    EndTime=end_time,
                    Period=300,  # 5 minutes
                    Statistics=['Average']
                )
                
                if response['Datapoints']:
                    # Get the most recent datapoint
                    latest = max(response['Datapoints'], key=lambda x: x['Timestamp'])
                    metrics[metric_name] = latest['Average']
                else:
                    metrics[metric_name] = 0.0
                    
            except ClientError:
                metrics[metric_name] = 0.0
        
        return metrics
    
    async def _fetch_cost_metrics(self):
        """Fetch cost metrics from AWS Cost Explorer."""
        if not self.aws_available:
            return
        
        try:
            end_date = datetime.now().strftime('%Y-%m-%d')
            start_date = (datetime.now() - timedelta(days=30)).strftime('%Y-%m-%d')
            
            # Get cost and usage data
            response = self.ce_client.get_cost_and_usage(
                TimePeriod={
                    'Start': start_date,
                    'End': end_date
                },
                Granularity='DAILY',
                Metrics=['BlendedCost'],
                GroupBy=[
                    {
                        'Type': 'DIMENSION',
                        'Key': 'SERVICE'
                    }
                ]
            )
            
            # Process cost data
            service_costs = defaultdict(float)
            daily_totals = []
            
            for result in response['ResultsByTime']:
                daily_total = 0.0
                for group in result['Groups']:
                    service_name = group['Keys'][0]
                    cost = float(group['Metrics']['BlendedCost']['Amount'])
                    service_costs[service_name] += cost
                    daily_total += cost
                daily_totals.append(daily_total)
            
            # Calculate projections
            if daily_totals:
                avg_daily_cost = sum(daily_totals[-7:]) / min(7, len(daily_totals))  # Last 7 days average
                projected_monthly = avg_daily_cost * 30
            else:
                avg_daily_cost = 0.0
                projected_monthly = 0.0
            
            # Store cost metrics
            for service, total_cost in service_costs.items():
                cost_metric = CostMetrics(
                    service=service,
                    daily_cost=avg_daily_cost if service == 'Total' else total_cost / 30,
                    monthly_cost=total_cost,
                    projected_monthly=projected_monthly if service == 'Total' else total_cost
                )
                self.cost_cache[service] = cost_metric
                
        except ClientError as e:
            # Cost Explorer might not be available in all regions
            self.console.print(f"[yellow]Cost data not available: {e}[/yellow]")
    
    async def _evaluate_alerts(self):
        """Evaluate alert conditions and trigger notifications."""
        triggered_alerts = []
        
        for alert in self.alerts:
            if not alert.enabled:
                continue
            
            current_value = self._get_metric_value(alert.metric)
            if current_value is None:
                continue
            
            # Check threshold
            triggered = False
            if alert.comparison == "greater" and current_value > alert.threshold:
                triggered = True
            elif alert.comparison == "less" and current_value < alert.threshold:
                triggered = True
            elif alert.comparison == "equal" and abs(current_value - alert.threshold) < 0.01:
                triggered = True
            
            if triggered:
                # Avoid duplicate alerts (only trigger once per hour)
                now = datetime.now()
                if (alert.last_triggered is None or 
                    (now - alert.last_triggered).total_seconds() > 3600):
                    
                    alert.last_triggered = now
                    triggered_alerts.append((alert, current_value))
        
        # Process triggered alerts
        for alert, value in triggered_alerts:
            await self._handle_alert(alert, value)
    
    def _get_metric_value(self, metric_name: str) -> Optional[float]:
        """Get current value for a metric."""
        if metric_name in ["cpu_utilization", "memory_utilization"]:
            # Get average across all instances
            if not self.instances_cache:
                return None
            
            values = []
            for instance_metric in self.instances_cache.values():
                if metric_name == "cpu_utilization":
                    values.append(instance_metric.cpu_utilization)
                elif metric_name == "memory_utilization":
                    values.append(instance_metric.memory_utilization)
            
            return sum(values) / len(values) if values else None
        
        elif metric_name in ["daily_cost", "projected_monthly"]:
            # Get total cost metrics
            total_cost = self.cost_cache.get('Total')
            if total_cost:
                if metric_name == "daily_cost":
                    return total_cost.daily_cost
                elif metric_name == "projected_monthly":
                    return total_cost.projected_monthly
        
        return None
    
    async def _handle_alert(self, alert: AlertRule, value: float):
        """Handle triggered alert."""
        message = f"🚨 ALERT: {alert.name} - Current value: {value:.2f} (Threshold: {alert.threshold})"
        self.console.print(f"[bold red]{message}[/bold red]")
        
        # In a real implementation, you could:
        # - Send email notifications
        # - Write to log files
        # - Send to monitoring systems
        # - Trigger automated responses
    
    def create_dashboard_layout(self) -> Layout:
        """Create the main dashboard layout."""
        layout = Layout()
        
        layout.split_column(
            Layout(name="header", size=3),
            Layout(name="main", ratio=1),
            Layout(name="footer", size=3)
        )
        
        layout["main"].split_row(
            Layout(name="left", ratio=2),
            Layout(name="right", ratio=1)
        )
        
        layout["left"].split_column(
            Layout(name="instances", ratio=2),
            Layout(name="metrics", ratio=1)
        )
        
        layout["right"].split_column(
            Layout(name="costs", ratio=1),
            Layout(name="alerts", ratio=1)
        )
        
        return layout
    
    def update_header(self, layout: Layout):
        """Update header with current status."""
        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        aws_status = "🟢 Connected" if self.aws_available else "🔴 Disconnected"
        instance_count = len(self.instances_cache)
        
        header_text = f"[bold]AWS Resource Monitor[/bold] | {current_time} | {aws_status} | {instance_count} instances"
        
        layout["header"].update(Panel(
            Align.center(header_text),
            title="Status",
            border_style="blue"
        ))
    
    def update_instances_panel(self, layout: Layout):
        """Update instances monitoring panel."""
        if not self.instances_cache:
            layout["instances"].update(Panel(
                "[yellow]No instances found or AWS not available[/yellow]",
                title="🖥️ EC2 Instances",
                border_style="yellow"
            ))
            return
        
        instances_table = Table(box=box.ROUNDED)
        instances_table.add_column("Instance ID", style="cyan", width=19)
        instances_table.add_column("Type", style="yellow", width=12)
        instances_table.add_column("CPU %", style="red", justify="right", width=8)
        instances_table.add_column("Memory %", style="blue", justify="right", width=10)
        instances_table.add_column("Network In", style="green", justify="right", width=12)
        instances_table.add_column("Network Out", style="green", justify="right", width=12)
        instances_table.add_column("Status", style="white", width=10)
        
        for instance_id, metrics in self.instances_cache.items():
            # Format network values
            net_in = f"{metrics.network_in / 1024 / 1024:.1f} MB" if metrics.network_in > 0 else "0 MB"
            net_out = f"{metrics.network_out / 1024 / 1024:.1f} MB" if metrics.network_out > 0 else "0 MB"
            
            # Color code CPU usage
            cpu_color = "red" if metrics.cpu_utilization > 80 else "yellow" if metrics.cpu_utilization > 60 else "green"
            cpu_text = f"[{cpu_color}]{metrics.cpu_utilization:.1f}%[/{cpu_color}]"
            
            instances_table.add_row(
                instance_id,
                metrics.instance_type,
                cpu_text,
                f"{metrics.memory_utilization:.1f}%",
                net_in,
                net_out,
                metrics.status
            )
        
        layout["instances"].update(Panel(
            instances_table,
            title="🖥️ EC2 Instances",
            border_style="green"
        ))
    
    def update_metrics_panel(self, layout: Layout):
        """Update metrics history panel."""
        if not self.metrics_history:
            layout["metrics"].update(Panel(
                "[yellow]No metrics history available[/yellow]",
                title="📊 Metrics History",
                border_style="yellow"
            ))
            return
        
        # Create simple ASCII charts for CPU usage
        metrics_text = []
        
        for instance_id, history in list(self.metrics_history.items())[:3]:  # Show first 3 instances
            if len(history) > 1:
                cpu_values = [m.cpu_utilization for m in list(history)[-20:]]  # Last 20 points
                
                # Simple ASCII sparkline
                if cpu_values:
                    max_val = max(cpu_values) if max(cpu_values) > 0 else 1
                    normalized = [int(v / max_val * 8) for v in cpu_values]
                    bars = ''.join(['▁▂▃▄▅▆▇█'[min(val, 7)] for val in normalized])
                    
                    metrics_text.append(f"[cyan]{instance_id}[/cyan]: {bars} [{cpu_values[-1]:.1f}%]")
        
        if not metrics_text:
            metrics_text = ["[yellow]Collecting metrics...[/yellow]"]
        
        layout["metrics"].update(Panel(
            "\n".join(metrics_text),
            title="📊 CPU Usage Trends",
            border_style="cyan"
        ))
    
    def update_costs_panel(self, layout: Layout):
        """Update costs monitoring panel."""
        if not self.cost_cache:
            layout["costs"].update(Panel(
                "[yellow]Cost data not available[/yellow]",
                title="💰 Cost Tracking",
                border_style="yellow"
            ))
            return
        
        cost_table = Table(box=box.SIMPLE)
        cost_table.add_column("Service", style="cyan")
        cost_table.add_column("Monthly", style="green", justify="right")
        cost_table.add_column("Daily Avg", style="yellow", justify="right")
        
        # Sort by monthly cost, descending
        sorted_costs = sorted(self.cost_cache.items(), 
                            key=lambda x: x[1].monthly_cost, reverse=True)
        
        total_monthly = 0
        for service, cost_metric in sorted_costs[:5]:  # Top 5 services
            cost_table.add_row(
                service[:15],  # Truncate long service names
                f"${cost_metric.monthly_cost:.2f}",
                f"${cost_metric.daily_cost:.2f}"
            )
            total_monthly += cost_metric.monthly_cost
        
        # Add total and projection
        if 'Total' in self.cost_cache:
            total_cost = self.cost_cache['Total']
            cost_table.add_separator()
            cost_table.add_row(
                "[bold]Total[/bold]",
                f"[bold]${total_cost.monthly_cost:.2f}[/bold]",
                f"[bold]${total_cost.daily_cost:.2f}[/bold]"
            )
            
            projection_text = f"\n[yellow]Projected Month: ${total_cost.projected_monthly:.2f}[/yellow]"
        else:
            projection_text = ""
        
        layout["costs"].update(Panel(
            str(cost_table) + projection_text,
            title="💰 Cost Tracking",
            border_style="green"
        ))
    
    def update_alerts_panel(self, layout: Layout):
        """Update alerts panel."""
        alert_lines = []
        
        active_alerts = [a for a in self.alerts if a.enabled]
        if not active_alerts:
            alert_lines.append("[yellow]No alerts configured[/yellow]")
        else:
            for alert in active_alerts[:5]:  # Show first 5 alerts
                current_value = self._get_metric_value(alert.metric)
                
                if current_value is not None:
                    status = "🔴" if current_value > alert.threshold else "🟢"
                    alert_lines.append(f"{status} {alert.name}: {current_value:.1f}/{alert.threshold}")
                else:
                    alert_lines.append(f"⚪ {alert.name}: No data")
        
        layout["alerts"].update(Panel(
            "\n".join(alert_lines),
            title="🚨 Alert Status",
            border_style="red"
        ))
    
    def update_footer(self, layout: Layout):
        """Update footer with controls."""
        footer_text = "[bold]Controls:[/bold] [cyan]q[/cyan] Quit | [cyan]r[/cyan] Refresh | [cyan]a[/cyan] Alerts | [cyan]c[/cyan] Costs | [cyan]h[/cyan] Help"
        
        layout["footer"].update(Panel(
            Align.center(footer_text),
            border_style="blue"
        ))
    
    def run_dashboard(self):
        """Run the live monitoring dashboard."""
        layout = self.create_dashboard_layout()
        
        def update_dashboard():
            """Update all dashboard panels."""
            self.update_header(layout)
            self.update_instances_panel(layout)
            self.update_metrics_panel(layout)
            self.update_costs_panel(layout)
            self.update_alerts_panel(layout)
            self.update_footer(layout)
            return layout
        
        if self.aws_available:
            # Start background monitoring
            def run_monitoring():
                asyncio.set_event_loop(asyncio.new_event_loop())
                loop = asyncio.get_event_loop()
                loop.run_until_complete(self.start_monitoring())
            
            monitoring_thread = threading.Thread(target=run_monitoring, daemon=True)
            monitoring_thread.start()
        
        # Run live dashboard
        try:
            with Live(update_dashboard(), refresh_per_second=1, screen=True) as live:
                while True:
                    # Check for keyboard input (simplified)
                    time.sleep(1)
                    live.update(update_dashboard())
        except KeyboardInterrupt:
            self.monitoring_active = False
            self.console.print("\n[green]Monitoring stopped.[/green]")

class AWSMonitorTUI:
    """Terminal User Interface for AWS monitoring."""
    
    def __init__(self, region: str = "us-east-1"):
        self.console = Console()
        self.monitor = AWSResourceMonitor(region)
    
    def run(self):
        """Start the monitoring TUI."""
        self.console.clear()
        self.console.print(Panel.fit(
            "[bold blue]AWS Resource Monitor[/bold blue]\n"
            "Real-time monitoring dashboard",
            title="Welcome",
            border_style="blue"
        ))
        
        if not self.monitor.aws_available:
            self.console.print(Panel(
                "[red]AWS credentials not configured or insufficient permissions.[/red]\n\n"
                "Please ensure you have:\n"
                "• AWS credentials configured (aws configure)\n"
                "• EC2, CloudWatch, and Cost Explorer permissions\n"
                "• Valid AWS region access",
                title="⚠️ AWS Access Required",
                border_style="red"
            ))
            return
        
        while True:
            self._show_main_menu()
            choice = self.console.input("\n[bold]Enter your choice: [/bold]")
            
            if choice == "1":
                self.monitor.run_dashboard()
            elif choice == "2":
                self._instance_management_menu()
            elif choice == "3":
                self._cost_analysis_menu()
            elif choice == "4":
                self._alerts_management_menu()
            elif choice == "5":
                self._settings_menu()
            elif choice == "q":
                self.console.print("[green]Goodbye![/green]")
                break
            else:
                self.console.print("[red]Invalid choice. Please try again.[/red]")
    
    def _show_main_menu(self):
        """Display the main menu."""
        self.console.clear()
        
        menu_table = Table(title="AWS Resource Monitor - Main Menu", 
                          title_style="bold blue")
        menu_table.add_column("Option", style="cyan", width=8)
        menu_table.add_column("Description", style="white")
        menu_table.add_column("Status", style="green")
        
        aws_status = "Connected" if self.monitor.aws_available else "Not Available"
        instance_count = len(self.monitor.instances_cache)
        alert_count = len([a for a in self.monitor.alerts if a.enabled])
        
        menu_table.add_row("1", "Live Dashboard", f"Region: {self.monitor.region}")
        menu_table.add_row("2", "Instance Management", f"{instance_count} instances")
        menu_table.add_row("3", "Cost Analysis", "Real-time tracking")
        menu_table.add_row("4", "Alert Management", f"{alert_count} active alerts")
        menu_table.add_row("5", "Settings", "Configuration")
        menu_table.add_row("q", "Quit", "")
        
        self.console.print(menu_table)
        
        # Show AWS status
        status_color = "green" if self.monitor.aws_available else "red"
        self.console.print(f"\n[{status_color}]AWS Status: {aws_status}[/{status_color}]")
    
    def _instance_management_menu(self):
        """Instance management interface."""
        self.console.clear()
        self.console.print("[bold]Instance Management[/bold]\n")
        
        if not self.monitor.aws_available:
            self.console.print("[red]AWS not available[/red]")
            self.console.input("Press Enter to continue...")
            return
        
        # Fetch fresh instance data
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.run_until_complete(self.monitor._fetch_instance_metrics())
        loop.close()
        
        if not self.monitor.instances_cache:
            self.console.print("[yellow]No running instances found[/yellow]")
            self.console.input("Press Enter to continue...")
            return
        
        # Display detailed instance information
        for instance_id, metrics in self.monitor.instances_cache.items():
            instance_panel = Panel(
                f"[cyan]Type:[/cyan] {metrics.instance_type}\n"
                f"[cyan]Status:[/cyan] {metrics.status}\n"
                f"[cyan]CPU:[/cyan] {metrics.cpu_utilization:.1f}%\n"
                f"[cyan]Memory:[/cyan] {metrics.memory_utilization:.1f}%\n"
                f"[cyan]Network In:[/cyan] {metrics.network_in / 1024 / 1024:.1f} MB\n"
                f"[cyan]Network Out:[/cyan] {metrics.network_out / 1024 / 1024:.1f} MB",
                title=f"🖥️ {instance_id}",
                border_style="green"
            )
            self.console.print(instance_panel)
        
        self.console.input("\nPress Enter to continue...")
    
    def _cost_analysis_menu(self):
        """Cost analysis interface."""
        self.console.clear()
        self.console.print("[bold]Cost Analysis[/bold]\n")
        
        if not self.monitor.aws_available:
            self.console.print("[red]AWS not available[/red]")
            self.console.input("Press Enter to continue...")
            return
        
        # Fetch fresh cost data
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.run_until_complete(self.monitor._fetch_cost_metrics())
        loop.close()
        
        if not self.monitor.cost_cache:
            self.console.print("[yellow]Cost data not available (may take time to populate)[/yellow]")
            self.console.input("Press Enter to continue...")
            return
        
        # Display cost breakdown
        cost_table = Table(title="Cost Breakdown", box=box.ROUNDED)
        cost_table.add_column("Service", style="cyan")
        cost_table.add_column("Daily Average", style="yellow", justify="right")
        cost_table.add_column("Monthly Total", style="green", justify="right")
        cost_table.add_column("Projected", style="red", justify="right")
        
        total_monthly = 0
        for service, cost_metric in self.monitor.cost_cache.items():
            cost_table.add_row(
                service,
                f"${cost_metric.daily_cost:.2f}",
                f"${cost_metric.monthly_cost:.2f}",
                f"${cost_metric.projected_monthly:.2f}"
            )
            total_monthly += cost_metric.monthly_cost
        
        self.console.print(cost_table)
        
        # Cost optimization suggestions
        suggestions = [
            "💡 Consider Reserved Instances for long-running workloads",
            "💡 Use Spot Instances for fault-tolerant batch processing",
            "💡 Enable detailed monitoring for better optimization insights",
            "💡 Review unused EBS volumes and snapshots",
            "💡 Implement auto-scaling to match demand"
        ]
        
        self.console.print(Panel(
            "\n".join(suggestions),
            title="💰 Cost Optimization Suggestions",
            border_style="yellow"
        ))
        
        self.console.input("\nPress Enter to continue...")
    
    def _alerts_management_menu(self):
        """Alert management interface."""
        self.console.clear()
        self.console.print("[bold]Alert Management[/bold]\n")
        
        # Display current alerts
        alerts_table = Table(title="Configured Alerts", box=box.ROUNDED)
        alerts_table.add_column("#", width=3)
        alerts_table.add_column("Name", style="cyan")
        alerts_table.add_column("Metric", style="yellow")
        alerts_table.add_column("Threshold", style="green", justify="right")
        alerts_table.add_column("Status", style="white")
        alerts_table.add_column("Last Triggered", style="red")
        
        for i, alert in enumerate(self.monitor.alerts, 1):
            status = "🟢 Enabled" if alert.enabled else "🔴 Disabled"
            last_triggered = alert.last_triggered.strftime("%H:%M:%S") if alert.last_triggered else "Never"
            
            alerts_table.add_row(
                str(i),
                alert.name,
                alert.metric,
                str(alert.threshold),
                status,
                last_triggered
            )
        
        self.console.print(alerts_table)
        
        self.console.print("\n[bold]Alert Options:[/bold]")
        self.console.print("1. Add new alert")
        self.console.print("2. Toggle alert on/off")
        self.console.print("3. Modify alert threshold")
        self.console.print("4. Delete alert")
        self.console.print("b. Back to main menu")
        
        choice = self.console.input("\n[bold]Enter your choice: [/bold]")
        
        if choice == "1":
            self._add_alert()
        elif choice == "2":
            self._toggle_alert()
        elif choice == "3":
            self._modify_alert_threshold()
        # elif choice == "b":
        #     return
        else:
            self.console.print("[yellow]Feature coming soon![/yellow]")
            self.console.input("Press Enter to continue...")
    
    def _add_alert(self):
        """Add a new alert rule."""
        self.console.print("\n[bold]Add New Alert[/bold]")
        
        # Get alert details from user
        name = self.console.input("Alert name: ")
        
        self.console.print("\nAvailable metrics:")
        self.console.print("1. cpu_utilization")
        self.console.print("2. memory_utilization") 
        self.console.print("3. daily_cost")
        self.console.print("4. projected_monthly")
        
        metric_choice = self.console.input("Select metric (1-4): ")
        metric_map = {
            "1": "cpu_utilization",
            "2": "memory_utilization",
            "3": "daily_cost",
            "4": "projected_monthly"
        }
        
        metric = metric_map.get(metric_choice)
        if not metric:
            self.console.print("[red]Invalid metric selection[/red]")
            return
        
        try:
            threshold = float(self.console.input("Threshold value: "))
        except ValueError:
            self.console.print("[red]Invalid threshold value[/red]")
            return
        
        # Create and add alert
        new_alert = AlertRule(
            name=name,
            metric=metric,
            threshold=threshold,
            comparison="greater",
            enabled=True
        )
        
        self.monitor.alerts.append(new_alert)
        self.console.print(f"[green]✓ Alert '{name}' added successfully[/green]")
        self.console.input("Press Enter to continue...")
    
    def _toggle_alert(self):
        """Toggle alert on/off."""
        if not self.monitor.alerts:
            self.console.print("[yellow]No alerts configured[/yellow]")
            self.console.input("Press Enter to continue...")
            return
        
        try:
            index = int(self.console.input(f"Enter alert number (1-{len(self.monitor.alerts)}): ")) - 1
            if 0 <= index < len(self.monitor.alerts):
                alert = self.monitor.alerts[index]
                alert.enabled = not alert.enabled
                status = "enabled" if alert.enabled else "disabled"
                self.console.print(f"[green]✓ Alert '{alert.name}' {status}[/green]")
            else:
                self.console.print("[red]Invalid alert number[/red]")
        except ValueError:
            self.console.print("[red]Please enter a valid number[/red]")
        
        self.console.input("Press Enter to continue...")
    
    def _modify_alert_threshold(self):
        """Modify alert threshold."""
        if not self.monitor.alerts:
            self.console.print("[yellow]No alerts configured[/yellow]")
            self.console.input("Press Enter to continue...")
            return
        
        try:
            index = int(self.console.input(f"Enter alert number (1-{len(self.monitor.alerts)}): ")) - 1
            if 0 <= index < len(self.monitor.alerts):
                alert = self.monitor.alerts[index]
                self.console.print(f"Current threshold for '{alert.name}': {alert.threshold}")
                
                new_threshold = float(self.console.input("New threshold value: "))
                alert.threshold = new_threshold
                self.console.print(f"[green]✓ Threshold updated to {new_threshold}[/green]")
            else:
                self.console.print("[red]Invalid alert number[/red]")
        except ValueError:
            self.console.print("[red]Please enter valid values[/red]")
        
        self.console.input("Press Enter to continue...")
    
    def _settings_menu(self):
        """Settings and configuration menu."""
        self.console.clear()
        self.console.print("[bold]Settings & Configuration[/bold]\n")
        
        settings_table = Table(title="Current Settings", box=box.ROUNDED)
        settings_table.add_column("Setting", style="cyan")
        settings_table.add_column("Value", style="white")
        settings_table.add_column("Description", style="yellow")
        
        settings_table.add_row("AWS Region", self.monitor.region, "Current monitoring region")
        settings_table.add_row("Refresh Interval", f"{self.monitor.refresh_interval}s", "Data refresh frequency")
        settings_table.add_row("Alert Count", str(len(self.monitor.alerts)), "Number of configured alerts")
        settings_table.add_row("Monitoring Status", 
                             "Active" if self.monitor.monitoring_active else "Inactive",
                             "Background monitoring state")
        
        self.console.print(settings_table)
        
        self.console.print("\n[bold]Configuration Options:[/bold]")
        self.console.print("1. Change AWS region")
        self.console.print("2. Adjust refresh interval")
        self.console.print("3. Export configuration")
        self.console.print("4. Import configuration")
        self.console.print("b. Back to main menu")
        
        choice = self.console.input("\n[bold]Enter your choice: [/bold]")
        
        if choice == "1":
            new_region = self.console.input(f"Enter new region (current: {self.monitor.region}): ")
            if new_region:
                self.monitor.region = new_region
                self.console.print(f"[green]✓ Region changed to {new_region}[/green]")
                self.console.print("[yellow]Note: Restart required for changes to take effect[/yellow]")
        elif choice == "2":
            try:
                new_interval = int(self.console.input(f"Enter refresh interval in seconds (current: {self.monitor.refresh_interval}): "))
                if new_interval >= 10:
                    self.monitor.refresh_interval = new_interval
                    self.console.print(f"[green]✓ Refresh interval changed to {new_interval}s[/green]")
                else:
                    self.console.print("[red]Minimum interval is 10 seconds[/red]")
            except ValueError:
                self.console.print("[red]Please enter a valid number[/red]")
        else:
            self.console.print("[yellow]Feature coming soon![/yellow]")
        
        if choice in ["1", "2"]:
            self.console.input("Press Enter to continue...")

def main():
    """Main entry point for AWS Monitor TUI."""
    import argparse
    
    parser = argparse.ArgumentParser(description="AWS Resource Monitor Terminal Interface")
    parser.add_argument("--region", default="us-east-1", 
                       help="AWS region to monitor")
    parser.add_argument("--refresh", type=int, default=30,
                       help="Refresh interval in seconds")
    
    args = parser.parse_args()
    
    monitor_tui = AWSMonitorTUI(args.region)
    monitor_tui.monitor.refresh_interval = args.refresh
    
    try:
        monitor_tui.run()
    except KeyboardInterrupt:
        print("\nMonitoring stopped by user.")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()