#!/usr/bin/env python3
"""
Configurable Domain-Specific Dashboard for AWS Research Wizard

This module provides customizable terminal dashboards tailored to specific research domains.
Each domain can have specialized monitoring views, workflow status tracking, and
domain-specific metrics visualization.

Key Features:
- Domain-specific metric layouts and visualizations
- Configurable dashboard widgets and panels
- Real-time workflow monitoring
- Research-specific data analytics
- Customizable alerts and notifications
- Export capabilities for research reporting

Domain Dashboard Types:
- Genomics: Sample processing, variant calling progress, quality metrics
- Climate: Model runs, data processing, forecast accuracy
- Neuroscience: Brain imaging pipelines, analysis progress
- Materials: DFT calculations, molecular dynamics status
- And more...

Dashboard Components:
- Live workflow status
- Domain-specific metrics
- Resource utilization optimized for domain
- Progress tracking for long-running analyses
- Research output monitoring

Dependencies:
    - rich: Advanced terminal formatting
    - textual: Modern TUI framework
    - pyyaml: Configuration file parsing
    - asyncio: Asynchronous operations
"""

import os
import sys
import yaml
import json
import asyncio
import time
from pathlib import Path
from typing import Dict, List, Any, Optional, Callable
from dataclasses import dataclass, field
from datetime import datetime, timedelta
from collections import defaultdict, deque
import logging

# Rich terminal components
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.progress import Progress, TaskID, SpinnerColumn, TextColumn, BarColumn, TimeRemainingColumn
from rich.layout import Layout
from rich.text import Text
from rich.align import Align
from rich.columns import Columns
from rich.live import Live
from rich.tree import Tree
from rich import box
from rich.status import Status
from rich.json import JSON

# Import our core modules
from config_loader import ConfigLoader
from demo_workflow_engine import DemoWorkflowEngine, WorkflowExecution

@dataclass
class DashboardWidget:
    """Configuration for a dashboard widget."""
    name: str
    type: str  # "table", "progress", "chart", "log", "status"
    title: str
    position: str  # "top-left", "top-right", "bottom-left", "bottom-right", "full"
    refresh_interval: int = 5
    data_source: str = "static"  # "static", "workflow", "aws", "custom"
    config: Dict[str, Any] = field(default_factory=dict)

@dataclass
class DomainDashboardConfig:
    """Complete configuration for a domain-specific dashboard."""
    domain_name: str
    title: str
    description: str
    widgets: List[DashboardWidget]
    layout_config: Dict[str, Any]
    refresh_interval: int = 10
    auto_scroll: bool = True
    color_scheme: str = "default"

class DomainDashboard:
    """
    Configurable dashboard for specific research domains.

    Provides domain-tailored monitoring and visualization including:
    - Workflow progress tracking
    - Domain-specific metric displays
    - Real-time resource monitoring
    - Research output tracking
    """

    def __init__(self, domain_name: str, config_root: str = "configs"):
        self.domain_name = domain_name
        self.config_root = Path(config_root)
        self.console = Console()

        # Load configurations
        self.config_loader = ConfigLoader(str(config_root))
        self.workflow_engine = DemoWorkflowEngine()

        # Dashboard state
        self.is_running = False
        self.widgets_data = {}
        self.workflow_status = {}
        self.metrics_history = defaultdict(lambda: deque(maxlen=50))

        # Load domain configuration and dashboard config
        self.domain_config = self._load_domain_config()
        self.dashboard_config = self._load_dashboard_config()

        # Initialize data collectors
        self.data_collectors = self._setup_data_collectors()

    def _load_domain_config(self) -> Dict[str, Any]:
        """Load domain-specific configuration."""
        try:
            return self.config_loader.load_domain_pack(self.domain_name)
        except Exception as e:
            self.console.print(f"[red]Error loading domain config: {e}[/red]")
            return {}

    def _load_dashboard_config(self) -> DomainDashboardConfig:
        """Load dashboard configuration for the domain."""
        # Try to load custom dashboard config
        dashboard_config_file = self.config_root / "dashboards" / f"{self.domain_name}.yaml"

        if dashboard_config_file.exists():
            with open(dashboard_config_file, 'r') as f:
                config_data = yaml.safe_load(f)
                return self._parse_dashboard_config(config_data)
        else:
            # Generate default dashboard configuration
            return self._generate_default_dashboard_config()

    def _parse_dashboard_config(self, config_data: Dict[str, Any]) -> DomainDashboardConfig:
        """Parse dashboard configuration from YAML."""
        widgets = []
        for widget_data in config_data.get('widgets', []):
            widget = DashboardWidget(
                name=widget_data['name'],
                type=widget_data['type'],
                title=widget_data['title'],
                position=widget_data['position'],
                refresh_interval=widget_data.get('refresh_interval', 5),
                data_source=widget_data.get('data_source', 'static'),
                config=widget_data.get('config', {})
            )
            widgets.append(widget)

        return DomainDashboardConfig(
            domain_name=self.domain_name,
            title=config_data.get('title', f"{self.domain_name.title()} Dashboard"),
            description=config_data.get('description', ''),
            widgets=widgets,
            layout_config=config_data.get('layout', {}),
            refresh_interval=config_data.get('refresh_interval', 10),
            auto_scroll=config_data.get('auto_scroll', True),
            color_scheme=config_data.get('color_scheme', 'default')
        )

    def _generate_default_dashboard_config(self) -> DomainDashboardConfig:
        """Generate default dashboard configuration based on domain type."""

        # Domain-specific default widgets
        domain_widgets = {
            'genomics': [
                DashboardWidget(
                    name="sample_processing",
                    type="progress",
                    title="Sample Processing Pipeline",
                    position="top-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="quality_metrics",
                    type="table",
                    title="Quality Control Metrics",
                    position="top-right",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="variant_calling",
                    type="status",
                    title="Variant Calling Status",
                    position="bottom-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="resource_usage",
                    type="chart",
                    title="Resource Utilization",
                    position="bottom-right",
                    data_source="aws"
                )
            ],
            'climate_modeling': [
                DashboardWidget(
                    name="model_runs",
                    type="progress",
                    title="Climate Model Execution",
                    position="top-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="forecast_accuracy",
                    type="table",
                    title="Forecast Accuracy Metrics",
                    position="top-right",
                    data_source="custom"
                ),
                DashboardWidget(
                    name="data_processing",
                    type="status",
                    title="Data Processing Status",
                    position="bottom-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="temperature_trends",
                    type="chart",
                    title="Temperature Analysis",
                    position="bottom-right",
                    data_source="custom"
                )
            ],
            'neuroscience': [
                DashboardWidget(
                    name="imaging_pipeline",
                    type="progress",
                    title="Brain Imaging Pipeline",
                    position="top-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="subject_status",
                    type="table",
                    title="Subject Processing Status",
                    position="top-right",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="analysis_results",
                    type="status",
                    title="Analysis Results",
                    position="bottom-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="connectivity_metrics",
                    type="chart",
                    title="Brain Connectivity",
                    position="bottom-right",
                    data_source="custom"
                )
            ],
            'materials_science': [
                DashboardWidget(
                    name="dft_calculations",
                    type="progress",
                    title="DFT Calculations",
                    position="top-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="convergence_status",
                    type="table",
                    title="Convergence Monitoring",
                    position="top-right",
                    data_source="custom"
                ),
                DashboardWidget(
                    name="md_simulations",
                    type="status",
                    title="MD Simulations",
                    position="bottom-left",
                    data_source="workflow"
                ),
                DashboardWidget(
                    name="energy_plots",
                    type="chart",
                    title="Energy Convergence",
                    position="bottom-right",
                    data_source="custom"
                )
            ]
        }

        # Default widgets for unspecified domains
        default_widgets = [
            DashboardWidget(
                name="workflow_status",
                type="progress",
                title="Workflow Progress",
                position="top-left",
                data_source="workflow"
            ),
            DashboardWidget(
                name="system_metrics",
                type="table",
                title="System Metrics",
                position="top-right",
                data_source="aws"
            ),
            DashboardWidget(
                name="job_queue",
                type="status",
                title="Job Queue Status",
                position="bottom-left",
                data_source="workflow"
            ),
            DashboardWidget(
                name="resource_chart",
                type="chart",
                title="Resource Usage",
                position="bottom-right",
                data_source="aws"
            )
        ]

        widgets = domain_widgets.get(self.domain_name, default_widgets)

        return DomainDashboardConfig(
            domain_name=self.domain_name,
            title=f"{self.domain_name.replace('_', ' ').title()} Research Dashboard",
            description=f"Real-time monitoring dashboard for {self.domain_name} research workflows",
            widgets=widgets,
            layout_config={},
            refresh_interval=10,
            auto_scroll=True,
            color_scheme='default'
        )

    def _setup_data_collectors(self) -> Dict[str, Callable]:
        """Setup data collection functions for different widget types."""
        collectors = {
            'workflow': self._collect_workflow_data,
            'aws': self._collect_aws_data,
            'custom': self._collect_custom_data,
            'static': self._collect_static_data
        }
        return collectors

    async def _collect_workflow_data(self, widget: DashboardWidget) -> Dict[str, Any]:
        """Collect workflow-related data."""
        if widget.type == "progress":
            return await self._get_workflow_progress_data(widget)
        elif widget.type == "status":
            return await self._get_workflow_status_data(widget)
        elif widget.type == "table":
            return await self._get_workflow_table_data(widget)
        else:
            return {"status": "No data"}

    async def _get_workflow_progress_data(self, widget: DashboardWidget) -> Dict[str, Any]:
        """Get workflow progress data."""
        # Simulate workflow progress based on domain
        domain_workflows = {
            'genomics': [
                {"name": "Quality Control", "progress": 100, "status": "completed"},
                {"name": "Alignment", "progress": 85, "status": "running"},
                {"name": "Variant Calling", "progress": 45, "status": "running"},
                {"name": "Annotation", "progress": 0, "status": "pending"}
            ],
            'climate_modeling': [
                {"name": "Data Preprocessing", "progress": 100, "status": "completed"},
                {"name": "Model Initialization", "progress": 100, "status": "completed"},
                {"name": "Simulation Run", "progress": 65, "status": "running"},
                {"name": "Post-processing", "progress": 0, "status": "pending"}
            ],
            'neuroscience': [
                {"name": "Structural Processing", "progress": 100, "status": "completed"},
                {"name": "Functional Analysis", "progress": 75, "status": "running"},
                {"name": "Connectivity Analysis", "progress": 30, "status": "running"},
                {"name": "Statistical Analysis", "progress": 0, "status": "pending"}
            ]
        }

        workflows = domain_workflows.get(self.domain_name, [
            {"name": "Data Processing", "progress": 60, "status": "running"},
            {"name": "Analysis", "progress": 25, "status": "running"},
            {"name": "Validation", "progress": 0, "status": "pending"}
        ])

        return {"workflows": workflows}

    async def _get_workflow_status_data(self, widget: DashboardWidget) -> Dict[str, Any]:
        """Get workflow status information."""
        # Simulate job queue and status
        current_time = datetime.now()

        status_data = {
            'active_jobs': 3,
            'queued_jobs': 7,
            'completed_today': 12,
            'failed_jobs': 1,
            'estimated_completion': current_time + timedelta(hours=4, minutes=32),
            'cluster_utilization': 78.5
        }

        return status_data

    async def _get_workflow_table_data(self, widget: DashboardWidget) -> Dict[str, Any]:
        """Get workflow data for table display."""
        if self.domain_name == 'genomics':
            return {
                'headers': ['Sample ID', 'Status', 'Quality Score', 'Variants'],
                'rows': [
                    ['SAMPLE_001', '✓ Complete', '98.5', '1,245,678'],
                    ['SAMPLE_002', '🔄 Processing', '97.2', '1,123,456'],
                    ['SAMPLE_003', '⏳ Queued', '-', '-'],
                    ['SAMPLE_004', '⏳ Queued', '-', '-']
                ]
            }
        elif self.domain_name == 'climate_modeling':
            return {
                'headers': ['Model Run', 'Start Time', 'Progress', 'Temperature'],
                'rows': [
                    ['WRF_2024_01', '09:15', '100%', '15.2°C'],
                    ['WRF_2024_02', '10:30', '75%', '14.8°C'],
                    ['WRF_2024_03', '11:45', '45%', '-'],
                    ['WRF_2024_04', '-', '0%', '-']
                ]
            }
        else:
            return {
                'headers': ['Job ID', 'Status', 'Progress', 'Runtime'],
                'rows': [
                    ['JOB_001', 'Complete', '100%', '2h 15m'],
                    ['JOB_002', 'Running', '65%', '1h 22m'],
                    ['JOB_003', 'Queued', '0%', '-'],
                    ['JOB_004', 'Queued', '0%', '-']
                ]
            }

    async def _collect_aws_data(self, widget: DashboardWidget) -> Dict[str, Any]:
        """Collect AWS resource data."""
        # Simulate AWS metrics
        cpu_usage = 45 + (time.time() % 30)  # Simulated varying CPU
        memory_usage = 67 + (time.time() % 20)

        return {
            'cpu_utilization': cpu_usage,
            'memory_utilization': memory_usage,
            'network_in': 125.7,
            'network_out': 89.3,
            'disk_io': 234.5,
            'active_instances': 3,
            'total_cost_today': 47.82
        }

    async def _collect_custom_data(self, widget: DashboardWidget) -> Dict[str, Any]:
        """Collect domain-specific custom data."""
        if self.domain_name == 'climate_modeling' and widget.name == 'forecast_accuracy':
            return {
                'headers': ['Metric', 'Current', 'Target', 'Status'],
                'rows': [
                    ['Temperature RMSE', '1.2°C', '<2.0°C', '✓'],
                    ['Precipitation Bias', '5.3%', '<10%', '✓'],
                    ['Wind Speed RMSE', '2.8 m/s', '<3.0 m/s', '✓'],
                    ['Pressure RMSE', '1.1 hPa', '<2.0 hPa', '✓']
                ]
            }
        elif self.domain_name == 'neuroscience' and widget.name == 'connectivity_metrics':
            return {
                'default_mode_network': 0.75,
                'executive_network': 0.68,
                'salience_network': 0.72,
                'visual_network': 0.81,
                'motor_network': 0.69
            }
        elif self.domain_name == 'materials_science' and widget.name == 'convergence_status':
            return {
                'headers': ['Calculation', 'Energy (eV)', 'Force (eV/Å)', 'Converged'],
                'rows': [
                    ['Si_bulk', '-5.423', '0.001', '✓'],
                    ['GaAs_surface', '-8.234', '0.015', '🔄'],
                    ['Al2O3_slab', '-12.567', '0.025', '🔄'],
                    ['TiO2_interface', '-15.789', '0.045', '❌']
                ]
            }
        else:
            # Generic custom data
            return {
                'metric_1': 78.5,
                'metric_2': 92.1,
                'metric_3': 65.7,
                'status': 'operational'
            }

    async def _collect_static_data(self, widget: DashboardWidget) -> Dict[str, Any]:
        """Collect static configuration data."""
        return {
            'domain': self.domain_name,
            'instance_type': self.domain_config.get('aws_instance_recommendations', {}).get('development', {}).get('instance_type', 'N/A'),
            'software_packages': len(self.domain_config.get('spack_packages', {})),
            'estimated_cost': self.domain_config.get('estimated_cost', {}).get('total', 0)
        }

    def create_dashboard_layout(self) -> Layout:
        """Create the dashboard layout based on configuration."""
        layout = Layout()

        # Header
        layout.split_column(
            Layout(name="header", size=3),
            Layout(name="main", ratio=1),
            Layout(name="footer", size=2)
        )

        # Main area split based on widget positions
        layout["main"].split_column(
            Layout(name="top", ratio=1),
            Layout(name="bottom", ratio=1)
        )

        layout["top"].split_row(
            Layout(name="top-left", ratio=1),
            Layout(name="top-right", ratio=1)
        )

        layout["bottom"].split_row(
            Layout(name="bottom-left", ratio=1),
            Layout(name="bottom-right", ratio=1)
        )

        return layout

    def update_header(self, layout: Layout):
        """Update dashboard header."""
        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

        header_text = (
            f"[bold]{self.dashboard_config.title}[/bold] | "
            f"{current_time} | "
            f"Domain: {self.domain_name} | "
            f"Auto-refresh: {self.dashboard_config.refresh_interval}s"
        )

        layout["header"].update(Panel(
            Align.center(header_text),
            border_style="blue"
        ))

    def update_footer(self, layout: Layout):
        """Update dashboard footer."""
        footer_text = (
            "[bold]Controls:[/bold] [cyan]q[/cyan] Quit | "
            "[cyan]r[/cyan] Refresh | [cyan]p[/cyan] Pause | "
            "[cyan]c[/cyan] Configure | [cyan]h[/cyan] Help"
        )

        layout["footer"].update(Panel(
            Align.center(footer_text),
            border_style="dim"
        ))

    async def update_widget(self, layout: Layout, widget: DashboardWidget):
        """Update a specific widget with fresh data."""
        try:
            # Collect data for the widget
            data_collector = self.data_collectors.get(widget.data_source)
            if data_collector:
                widget_data = await data_collector(widget)
                self.widgets_data[widget.name] = widget_data
            else:
                widget_data = {"error": f"Unknown data source: {widget.data_source}"}

            # Render widget based on type
            if widget.type == "progress":
                rendered_widget = self._render_progress_widget(widget, widget_data)
            elif widget.type == "table":
                rendered_widget = self._render_table_widget(widget, widget_data)
            elif widget.type == "status":
                rendered_widget = self._render_status_widget(widget, widget_data)
            elif widget.type == "chart":
                rendered_widget = self._render_chart_widget(widget, widget_data)
            else:
                rendered_widget = Panel(f"Unknown widget type: {widget.type}",
                                      title=widget.title, border_style="red")

            # Update layout position
            layout[widget.position].update(rendered_widget)

        except Exception as e:
            error_panel = Panel(f"Error: {str(e)}", title=widget.title, border_style="red")
            layout[widget.position].update(error_panel)

    def _render_progress_widget(self, widget: DashboardWidget, data: Dict[str, Any]) -> Panel:
        """Render a progress widget."""
        if 'workflows' in data:
            progress_content = []

            for workflow in data['workflows']:
                name = workflow['name']
                progress = workflow['progress']
                status = workflow['status']

                # Status indicators
                status_indicator = {
                    'completed': '✅',
                    'running': '🔄',
                    'pending': '⏳',
                    'failed': '❌'
                }.get(status, '❓')

                # Progress bar
                bar_length = 20
                filled_length = int(bar_length * progress / 100)
                bar = '█' * filled_length + '░' * (bar_length - filled_length)

                progress_content.append(
                    f"{status_indicator} {name:<20} [{bar}] {progress:3d}%"
                )

            content = "\n".join(progress_content)
        else:
            content = "[yellow]No workflow data available[/yellow]"

        return Panel(content, title=widget.title, border_style="green")

    def _render_table_widget(self, widget: DashboardWidget, data: Dict[str, Any]) -> Panel:
        """Render a table widget."""
        if 'headers' in data and 'rows' in data:
            table = Table(box=box.SIMPLE)

            # Add columns
            for header in data['headers']:
                table.add_column(header, style="cyan")

            # Add rows
            for row in data['rows']:
                table.add_row(*[str(cell) for cell in row])

            return Panel(table, title=widget.title, border_style="cyan")
        else:
            return Panel("[yellow]No table data available[/yellow]",
                        title=widget.title, border_style="yellow")

    def _render_status_widget(self, widget: DashboardWidget, data: Dict[str, Any]) -> Panel:
        """Render a status widget."""
        if isinstance(data, dict):
            status_lines = []

            for key, value in data.items():
                if key == 'estimated_completion':
                    if isinstance(value, datetime):
                        value = value.strftime("%H:%M")

                # Format key names
                display_key = key.replace('_', ' ').title()

                # Color code based on values
                if isinstance(value, (int, float)):
                    if value > 80:
                        color = "red"
                    elif value > 60:
                        color = "yellow"
                    else:
                        color = "green"
                    status_lines.append(f"[cyan]{display_key}:[/cyan] [{color}]{value}[/{color}]")
                else:
                    status_lines.append(f"[cyan]{display_key}:[/cyan] {value}")

            content = "\n".join(status_lines)
        else:
            content = "[yellow]No status data available[/yellow]"

        return Panel(content, title=widget.title, border_style="yellow")

    def _render_chart_widget(self, widget: DashboardWidget, data: Dict[str, Any]) -> Panel:
        """Render a chart widget with ASCII visualization."""
        if isinstance(data, dict) and data:
            chart_lines = []

            # Simple bar chart for metrics
            if self.domain_name == 'neuroscience' and isinstance(data, dict):
                # Network connectivity visualization
                for network, value in data.items():
                    if isinstance(value, (int, float)):
                        bar_length = int(value * 20)  # Scale to 20 chars
                        bar = '█' * bar_length + '░' * (20 - bar_length)
                        chart_lines.append(f"{network:<20} [{bar}] {value:.2f}")
            else:
                # Generic metrics chart
                for key, value in data.items():
                    if isinstance(value, (int, float)):
                        # Simple horizontal bar
                        if value <= 100:  # Percentage values
                            bar_length = int(value * 20 / 100)
                        else:  # Other values, normalize to max
                            max_val = max([v for v in data.values() if isinstance(v, (int, float))])
                            bar_length = int(value * 20 / max_val) if max_val > 0 else 0

                        bar = '█' * bar_length + '░' * (20 - bar_length)
                        chart_lines.append(f"{key:<15} [{bar}] {value}")

            content = "\n".join(chart_lines) if chart_lines else "[yellow]No chart data[/yellow]"
        else:
            content = "[yellow]No chart data available[/yellow]"

        return Panel(content, title=widget.title, border_style="magenta")

    async def run_dashboard(self):
        """Run the live dashboard."""
        self.is_running = True
        layout = self.create_dashboard_layout()

        def update_dashboard():
            """Update all dashboard components."""
            self.update_header(layout)
            self.update_footer(layout)
            return layout

        # Start data collection in background
        async def collect_data():
            while self.is_running:
                try:
                    # Update all widgets
                    for widget in self.dashboard_config.widgets:
                        await self.update_widget(layout, widget)

                    await asyncio.sleep(self.dashboard_config.refresh_interval)
                except Exception as e:
                    self.console.print(f"[red]Data collection error: {e}[/red]")
                    await asyncio.sleep(5)

        # Start background data collection
        data_task = asyncio.create_task(collect_data())

        try:
            # Initial update
            for widget in self.dashboard_config.widgets:
                await self.update_widget(layout, widget)

            # Run live dashboard
            with Live(update_dashboard(), refresh_per_second=1, screen=True) as live:
                while self.is_running:
                    await asyncio.sleep(0.1)
                    live.update(update_dashboard())

        except KeyboardInterrupt:
            self.is_running = False
            self.console.print("\n[green]Dashboard stopped.[/green]")
        finally:
            data_task.cancel()

class DomainDashboardTUI:
    """Terminal User Interface for domain-specific dashboards."""

    def __init__(self, config_root: str = "configs"):
        self.config_root = config_root
        self.console = Console()
        self.config_loader = ConfigLoader(config_root)

        # Load available domains
        self.available_domains = self._get_available_domains()

    def _get_available_domains(self) -> List[str]:
        """Get list of available research domains."""
        domains_dir = Path(self.config_root) / "domains"
        if not domains_dir.exists():
            return []

        domains = []
        for domain_file in domains_dir.glob("*.yaml"):
            domains.append(domain_file.stem)

        return sorted(domains)

    def run(self):
        """Start the domain dashboard TUI."""
        self.console.clear()
        self.console.print(Panel.fit(
            "[bold blue]Domain-Specific Research Dashboards[/bold blue]\n"
            "Configurable monitoring for research domains",
            title="Welcome",
            border_style="blue"
        ))

        while True:
            self._show_main_menu()
            choice = self.console.input("\n[bold]Enter your choice: [/bold]")

            if choice == "1":
                asyncio.run(self._select_and_run_dashboard())
            elif choice == "2":
                self._configure_dashboard()
            elif choice == "3":
                self._create_custom_dashboard()
            elif choice == "4":
                self._export_dashboard_config()
            elif choice == "q":
                self.console.print("[green]Goodbye![/green]")
                break
            else:
                self.console.print("[red]Invalid choice. Please try again.[/red]")

    def _show_main_menu(self):
        """Display the main menu."""
        self.console.clear()

        menu_table = Table(title="Domain Dashboard Manager", title_style="bold blue")
        menu_table.add_column("Option", style="cyan", width=8)
        menu_table.add_column("Description", style="white")
        menu_table.add_column("Status", style="green")

        menu_table.add_row("1", "Launch Domain Dashboard", f"{len(self.available_domains)} domains available")
        menu_table.add_row("2", "Configure Dashboard", "Customize widgets and layout")
        menu_table.add_row("3", "Create Custom Dashboard", "Build new dashboard configuration")
        menu_table.add_row("4", "Export Configuration", "Save dashboard settings")
        menu_table.add_row("q", "Quit", "")

        self.console.print(menu_table)

    async def _select_and_run_dashboard(self):
        """Select domain and run its dashboard."""
        if not self.available_domains:
            self.console.print("[red]No research domains found[/red]")
            self.console.input("Press Enter to continue...")
            return

        self.console.print("\n[bold]Available Research Domains:[/bold]")

        domain_table = Table(box=box.ROUNDED)
        domain_table.add_column("#", width=3)
        domain_table.add_column("Domain", style="cyan")
        domain_table.add_column("Description", style="white", max_width=50)

        for i, domain in enumerate(self.available_domains, 1):
            try:
                config = self.config_loader.load_domain_pack(domain)
                description = config.get('description', 'No description available')[:50]
            except:
                description = "Configuration not available"

            domain_table.add_row(str(i), domain.replace('_', ' ').title(), description)

        self.console.print(domain_table)

        try:
            choice = int(self.console.input(f"\nSelect domain (1-{len(self.available_domains)}): "))
            if 1 <= choice <= len(self.available_domains):
                selected_domain = self.available_domains[choice - 1]

                self.console.print(f"[green]Starting dashboard for {selected_domain}...[/green]")
                self.console.print("[dim]Press Ctrl+C to stop the dashboard[/dim]")

                # Run the domain dashboard
                dashboard = DomainDashboard(selected_domain, self.config_root)
                await dashboard.run_dashboard()
            else:
                self.console.print("[red]Invalid selection[/red]")
        except ValueError:
            self.console.print("[red]Please enter a valid number[/red]")
        except KeyboardInterrupt:
            self.console.print("\n[yellow]Dashboard stopped[/yellow]")

        self.console.input("Press Enter to continue...")

    def _configure_dashboard(self):
        """Configure dashboard widgets and layout."""
        self.console.print("\n[bold]Dashboard Configuration[/bold]")

        if not self.available_domains:
            self.console.print("[red]No domains available to configure[/red]")
            self.console.input("Press Enter to continue...")
            return

        self.console.print("Available domains:")
        for i, domain in enumerate(self.available_domains, 1):
            self.console.print(f"  {i}. {domain.replace('_', ' ').title()}")

        try:
            choice = int(self.console.input(f"\nSelect domain to configure (1-{len(self.available_domains)}): "))
            if 1 <= choice <= len(self.available_domains):
                selected_domain = self.available_domains[choice - 1]
                self._configure_domain_dashboard(selected_domain)
            else:
                self.console.print("[red]Invalid selection[/red]")
        except ValueError:
            self.console.print("[red]Please enter a valid number[/red]")

        self.console.input("Press Enter to continue...")

    def _configure_domain_dashboard(self, domain: str):
        """Configure a specific domain dashboard."""
        self.console.print(f"\n[bold]Configuring dashboard for {domain}[/bold]")

        # Load existing dashboard config
        dashboard = DomainDashboard(domain, self.config_root)
        config = dashboard.dashboard_config

        # Show current configuration
        config_table = Table(title="Current Dashboard Configuration", box=box.ROUNDED)
        config_table.add_column("Widget", style="cyan")
        config_table.add_column("Type", style="yellow")
        config_table.add_column("Position", style="green")
        config_table.add_column("Data Source", style="white")

        for widget in config.widgets:
            config_table.add_row(
                widget.title,
                widget.type,
                widget.position,
                widget.data_source
            )

        self.console.print(config_table)

        self.console.print("\n[bold]Configuration Options:[/bold]")
        self.console.print("1. Add new widget")
        self.console.print("2. Modify existing widget")
        self.console.print("3. Remove widget")
        self.console.print("4. Change refresh interval")
        self.console.print("5. Export configuration")

        choice = self.console.input("\nEnter your choice: ")

        if choice == "4":
            try:
                new_interval = int(self.console.input(f"Current refresh interval: {config.refresh_interval}s\nNew interval (seconds): "))
                if new_interval >= 1:
                    config.refresh_interval = new_interval
                    self.console.print(f"[green]✓ Refresh interval updated to {new_interval}s[/green]")
                else:
                    self.console.print("[red]Interval must be at least 1 second[/red]")
            except ValueError:
                self.console.print("[red]Please enter a valid number[/red]")
        else:
            self.console.print("[yellow]Configuration feature coming soon![/yellow]")

    def _create_custom_dashboard(self):
        """Create a custom dashboard configuration."""
        self.console.print("\n[bold]Create Custom Dashboard[/bold]")

        # Get basic information
        domain_name = self.console.input("Domain name (e.g., 'custom_research'): ").strip()
        if not domain_name:
            self.console.print("[red]Domain name is required[/red]")
            return

        title = self.console.input(f"Dashboard title (default: '{domain_name.title()} Dashboard'): ").strip()
        if not title:
            title = f"{domain_name.title()} Dashboard"

        description = self.console.input("Dashboard description: ").strip()

        # Create basic configuration
        custom_config = {
            'title': title,
            'description': description,
            'refresh_interval': 10,
            'auto_scroll': True,
            'color_scheme': 'default',
            'widgets': [
                {
                    'name': 'workflow_status',
                    'type': 'progress',
                    'title': 'Workflow Progress',
                    'position': 'top-left',
                    'data_source': 'workflow'
                },
                {
                    'name': 'system_metrics',
                    'type': 'table',
                    'title': 'System Metrics',
                    'position': 'top-right',
                    'data_source': 'aws'
                },
                {
                    'name': 'job_status',
                    'type': 'status',
                    'title': 'Job Status',
                    'position': 'bottom-left',
                    'data_source': 'workflow'
                },
                {
                    'name': 'resource_usage',
                    'type': 'chart',
                    'title': 'Resource Usage',
                    'position': 'bottom-right',
                    'data_source': 'aws'
                }
            ]
        }

        # Save configuration
        dashboards_dir = Path(self.config_root) / "dashboards"
        dashboards_dir.mkdir(exist_ok=True)

        config_file = dashboards_dir / f"{domain_name}.yaml"

        with open(config_file, 'w') as f:
            yaml.dump(custom_config, f, default_flow_style=False, sort_keys=False)

        self.console.print(f"[green]✓ Custom dashboard configuration saved to {config_file}[/green]")
        self.console.print("[dim]You can now select this dashboard from the main menu[/dim]")

    def _export_dashboard_config(self):
        """Export dashboard configuration."""
        self.console.print("\n[bold]Export Dashboard Configuration[/bold]")

        dashboards_dir = Path(self.config_root) / "dashboards"
        if not dashboards_dir.exists() or not list(dashboards_dir.glob("*.yaml")):
            self.console.print("[yellow]No dashboard configurations found[/yellow]")
            self.console.input("Press Enter to continue...")
            return

        self.console.print("Available dashboard configurations:")
        config_files = list(dashboards_dir.glob("*.yaml"))

        for i, config_file in enumerate(config_files, 1):
            self.console.print(f"  {i}. {config_file.stem}")

        try:
            choice = int(self.console.input(f"\nSelect configuration to export (1-{len(config_files)}): "))
            if 1 <= choice <= len(config_files):
                selected_file = config_files[choice - 1]

                export_path = self.console.input(f"Export path (default: ./{selected_file.name}): ").strip()
                if not export_path:
                    export_path = f"./{selected_file.name}"

                # Copy configuration file
                import shutil
                shutil.copy2(selected_file, export_path)

                self.console.print(f"[green]✓ Configuration exported to {export_path}[/green]")
            else:
                self.console.print("[red]Invalid selection[/red]")
        except ValueError:
            self.console.print("[red]Please enter a valid number[/red]")

        self.console.input("Press Enter to continue...")

def main():
    """Main entry point for domain dashboard TUI."""
    import argparse

    parser = argparse.ArgumentParser(description="Domain-Specific Dashboard Terminal Interface")
    parser.add_argument("--config", default="configs",
                       help="Configuration directory path")
    parser.add_argument("--domain",
                       help="Directly launch dashboard for specified domain")

    args = parser.parse_args()

    if args.domain:
        # Direct launch for specific domain
        dashboard = DomainDashboard(args.domain, args.config)
        try:
            asyncio.run(dashboard.run_dashboard())
        except KeyboardInterrupt:
            print("\nDashboard stopped.")
    else:
        # Interactive menu
        dashboard_tui = DomainDashboardTUI(args.config)
        dashboard_tui.run()

if __name__ == "__main__":
    main()
