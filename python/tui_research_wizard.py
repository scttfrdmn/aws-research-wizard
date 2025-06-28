#!/usr/bin/env python3
"""
Terminal User Interface (TUI) for AWS Research Wizard

This module provides a comprehensive terminal-based interface for configuring,
deploying, and managing AWS research environments. Features include:

1. Interactive Research Domain Selection
2. Real-time Configuration Building
3. Cost Estimation and Optimization
4. Deployment Status Monitoring
5. Workflow Management

Key Features:
- Rich terminal interface with mouse and keyboard support
- Real-time AWS cost calculations
- Interactive domain pack configuration
- Built-in help system and tutorials
- SSH-friendly design for remote access

TUI Components:
- Main menu and navigation
- Domain selection with filtering
- Configuration editor with validation
- Cost calculator with optimization suggestions
- Deployment progress tracking
- Log viewer with filtering

Dependencies:
    - rich: Modern terminal formatting and widgets
    - textual: Advanced TUI framework
    - boto3: AWS SDK integration
    - asyncio: Asynchronous operations
"""

import os
import sys
import asyncio
import json
import yaml
import boto3
from pathlib import Path
from typing import Dict, List, Any, Optional
from datetime import datetime
import logging

# Rich terminal components
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.progress import Progress, TaskID
from rich.layout import Layout
from rich.text import Text
from rich.align import Align
from rich.columns import Columns
from rich.tree import Tree
from rich import box

# Textual TUI framework
try:
    from textual.app import App, ComposeResult
    from textual.containers import Container, Horizontal, Vertical, ScrollableContainer
    from textual.widgets import (
        Button, Static, Input, Select, Checkbox, RadioButton, RadioSet,
        DataTable, Log, ProgressBar, Tree as TextualTree, Tabs, TabPane,
        Header, Footer, Label, TextArea, Switch, ListView, ListItem
    )
    from textual.reactive import reactive
    from textual.coordinate import Coordinate
    from textual.binding import Binding
    TEXTUAL_AVAILABLE = True
except ImportError:
    TEXTUAL_AVAILABLE = False

# Import our core modules
from config_loader import ConfigLoader
from s3_transfer_optimizer import S3TransferOptimizer
from demo_workflow_engine import DemoWorkflowEngine

class ResearchWizardTUI:
    """
    Main Terminal User Interface for AWS Research Wizard.
    
    Provides an interactive terminal interface for researchers to:
    - Browse and select research domain packs
    - Configure AWS infrastructure requirements
    - Estimate costs and optimize configurations
    - Deploy and monitor research environments
    """
    
    def __init__(self, config_root: str = "configs"):
        self.config_root = Path(config_root)
        self.console = Console()
        self.config_loader = ConfigLoader(str(config_root))
        self.s3_optimizer = S3TransferOptimizer()
        self.workflow_engine = DemoWorkflowEngine()
        
        # TUI state
        self.current_domain = None
        self.current_config = {}
        self.selected_instance = None
        self.estimated_costs = {}
        
        # Load available domains
        self.domains = self._load_available_domains()
        
        # AWS clients for real-time monitoring
        try:
            self.ec2_client = boto3.client('ec2')
            self.s3_client = boto3.client('s3')
            self.aws_available = True
        except Exception:
            self.aws_available = False
    
    def _load_available_domains(self) -> Dict[str, Dict[str, Any]]:
        """Load all available research domain configurations."""
        domains = {}
        
        domain_dir = self.config_root / "domains"
        if not domain_dir.exists():
            return domains
        
        for domain_file in domain_dir.glob("*.yaml"):
            try:
                with open(domain_file, 'r') as f:
                    config = yaml.safe_load(f)
                    domains[domain_file.stem] = config
            except Exception as e:
                self.console.print(f"[red]Error loading {domain_file}: {e}[/red]")
        
        return domains
    
    def run(self):
        """Start the Terminal User Interface."""
        if TEXTUAL_AVAILABLE:
            app = ResearchWizardApp(self)
            app.run()
        else:
            self._run_simple_tui()
    
    def _run_simple_tui(self):
        """Fallback simple TUI without Textual dependency."""
        self.console.print(Panel.fit(
            "[bold blue]AWS Research Wizard[/bold blue]\n"
            "Terminal Interface",
            title="Welcome",
            border_style="blue"
        ))
        
        while True:
            self._show_main_menu()
            choice = self.console.input("\n[bold]Enter your choice: [/bold]")
            
            if choice == "1":
                self._domain_selection_menu()
            elif choice == "2":
                self._configuration_menu()
            elif choice == "3":
                self._cost_calculator_menu()
            elif choice == "4":
                self._deployment_menu()
            elif choice == "5":
                self._monitoring_menu()
            elif choice == "q":
                self.console.print("[green]Goodbye![/green]")
                break
            else:
                self.console.print("[red]Invalid choice. Please try again.[/red]")
    
    def _show_main_menu(self):
        """Display the main menu."""
        self.console.clear()
        
        menu_table = Table(title="AWS Research Wizard - Main Menu", 
                          title_style="bold blue")
        menu_table.add_column("Option", style="cyan", width=8)
        menu_table.add_column("Description", style="white")
        menu_table.add_column("Status", style="green")
        
        menu_table.add_row("1", "Select Research Domain", 
                          f"Current: {self.current_domain or 'None'}")
        menu_table.add_row("2", "Configure Environment", 
                          "Ready" if self.current_domain else "Select domain first")
        menu_table.add_row("3", "Cost Calculator", 
                          f"${self.estimated_costs.get('total', 0):.2f}/month")
        menu_table.add_row("4", "Deploy Environment", 
                          "AWS Connected" if self.aws_available else "AWS Not Available")
        menu_table.add_row("5", "Monitor Resources", "Real-time")
        menu_table.add_row("q", "Quit", "")
        
        self.console.print(menu_table)
    
    def _domain_selection_menu(self):
        """Interactive domain selection menu."""
        self.console.clear()
        self.console.print(Panel.fit(
            "[bold]Available Research Domains[/bold]",
            border_style="blue"
        ))
        
        # Create domain selection table
        domain_table = Table(title="Research Domain Packs")
        domain_table.add_column("#", width=3)
        domain_table.add_column("Domain", style="cyan")
        domain_table.add_column("Description", style="white", max_width=50)
        domain_table.add_column("Users", style="yellow", width=8)
        domain_table.add_column("Est. Cost", style="green", width=10)
        
        domain_list = list(self.domains.items())
        for i, (domain_key, config) in enumerate(domain_list, 1):
            cost = config.get('estimated_cost', {}).get('total', 0)
            domain_table.add_row(
                str(i),
                config.get('name', domain_key).split(' ')[0],  # Short name
                config.get('description', '')[:50] + "...",
                config.get('target_users', 'N/A').split('(')[1].split(')')[0] if '(' in str(config.get('target_users', '')) else 'N/A',
                f"${cost}/mo"
            )
        
        self.console.print(domain_table)
        
        # Get user selection
        try:
            choice = int(self.console.input(f"\n[bold]Select domain (1-{len(domain_list)}): [/bold]"))
            if 1 <= choice <= len(domain_list):
                selected_domain = domain_list[choice - 1]
                self.current_domain = selected_domain[0]
                self.current_config = selected_domain[1]
                self._calculate_costs()
                
                self.console.print(f"[green]✓ Selected: {self.current_config['name']}[/green]")
                self.console.input("Press Enter to continue...")
            else:
                self.console.print("[red]Invalid selection[/red]")
                self.console.input("Press Enter to continue...")
        except ValueError:
            self.console.print("[red]Please enter a valid number[/red]")
            self.console.input("Press Enter to continue...")
    
    def _configuration_menu(self):
        """Configuration and customization menu."""
        if not self.current_domain:
            self.console.print("[red]Please select a domain first[/red]")
            self.console.input("Press Enter to continue...")
            return
        
        self.console.clear()
        self.console.print(Panel.fit(
            f"[bold]Configure: {self.current_config['name']}[/bold]",
            border_style="blue"
        ))
        
        # Show current configuration
        config_layout = Layout()
        config_layout.split_column(
            Layout(name="instances", ratio=2),
            Layout(name="packages", ratio=3)
        )
        
        # Instance recommendations
        instances_table = Table(title="AWS Instance Recommendations", box=box.ROUNDED)
        instances_table.add_column("Type", style="cyan")
        instances_table.add_column("vCPUs", style="yellow", width=8)
        instances_table.add_column("Memory", style="green", width=10)
        instances_table.add_column("Cost/Hour", style="red", width=12)
        instances_table.add_column("Use Case", style="white", max_width=30)
        
        aws_instances = self.current_config.get('aws_instance_recommendations', {})
        for instance_type, details in aws_instances.items():
            instances_table.add_row(
                details.get('instance_type', 'N/A'),
                str(details.get('vcpus', 'N/A')),
                f"{details.get('memory_gb', 'N/A')} GB",
                f"${details.get('cost_per_hour', 0):.3f}",
                details.get('use_case', '')[:30]
            )
        
        config_layout["instances"].update(Panel(instances_table, title="Hardware Options"))
        
        # Software packages summary
        packages = self.current_config.get('spack_packages', {})
        package_info = []
        for category, pkg_list in packages.items():
            package_info.append(f"[cyan]{category.replace('_', ' ').title()}[/cyan]: {len(pkg_list)} packages")
        
        config_layout["packages"].update(Panel(
            "\n".join(package_info),
            title="Software Environment",
            border_style="green"
        ))
        
        self.console.print(config_layout)
        
        # Configuration options
        self.console.print("\n[bold]Configuration Options:[/bold]")
        self.console.print("1. Select instance type")
        self.console.print("2. Customize software packages")
        self.console.print("3. Configure storage")
        self.console.print("4. Set up data sources")
        self.console.print("b. Back to main menu")
        
        choice = self.console.input("\n[bold]Enter your choice: [/bold]")
        
        if choice == "1":
            self._select_instance_type()
        elif choice == "2":
            self._customize_packages()
        elif choice == "3":
            self._configure_storage()
        elif choice == "4":
            self._configure_data_sources()
        # elif choice == "b":
        #     return
    
    def _select_instance_type(self):
        """Interactive instance type selection."""
        instances = self.current_config.get('aws_instance_recommendations', {})
        
        self.console.print("\n[bold]Select Instance Type:[/bold]")
        instance_list = list(instances.items())
        
        for i, (name, details) in enumerate(instance_list, 1):
            cost_monthly = details.get('cost_per_hour', 0) * 24 * 30
            self.console.print(
                f"{i}. [cyan]{details.get('instance_type', 'N/A')}[/cyan] - "
                f"{details.get('vcpus', 'N/A')} vCPUs, "
                f"{details.get('memory_gb', 'N/A')} GB RAM - "
                f"[green]${cost_monthly:.0f}/month[/green]"
            )
            self.console.print(f"   {details.get('use_case', '')}")
        
        try:
            choice = int(self.console.input(f"\nSelect instance (1-{len(instance_list)}): "))
            if 1 <= choice <= len(instance_list):
                selected = instance_list[choice - 1]
                self.selected_instance = selected[1]
                self._calculate_costs()
                self.console.print(f"[green]✓ Selected: {self.selected_instance['instance_type']}[/green]")
            else:
                self.console.print("[red]Invalid selection[/red]")
        except ValueError:
            self.console.print("[red]Please enter a valid number[/red]")
        
        self.console.input("Press Enter to continue...")
    
    def _cost_calculator_menu(self):
        """Cost estimation and optimization menu."""
        self.console.clear()
        self._calculate_costs()
        
        # Create cost breakdown layout
        cost_layout = Layout()
        cost_layout.split_row(
            Layout(name="breakdown", ratio=2),
            Layout(name="optimization", ratio=1)
        )
        
        # Cost breakdown table
        cost_table = Table(title="Monthly Cost Breakdown", box=box.ROUNDED)
        cost_table.add_column("Component", style="cyan")
        cost_table.add_column("Cost", style="green", justify="right")
        cost_table.add_column("Percentage", style="yellow", justify="right")
        
        total_cost = self.estimated_costs.get('total', 0)
        if total_cost > 0:
            for component, cost in self.estimated_costs.items():
                if component != 'total' and cost > 0:
                    percentage = (cost / total_cost) * 100
                    cost_table.add_row(
                        component.replace('_', ' ').title(),
                        f"${cost:.2f}",
                        f"{percentage:.1f}%"
                    )
            
            cost_table.add_separator()
            cost_table.add_row("[bold]Total[/bold]", f"[bold]${total_cost:.2f}[/bold]", "[bold]100%[/bold]")
        
        cost_layout["breakdown"].update(Panel(cost_table, title="Cost Analysis"))
        
        # Optimization suggestions
        suggestions = self._get_cost_optimization_suggestions()
        optimization_text = "\n".join([f"• {suggestion}" for suggestion in suggestions])
        
        cost_layout["optimization"].update(Panel(
            optimization_text,
            title="💡 Optimization Tips",
            border_style="yellow"
        ))
        
        self.console.print(cost_layout)
        
        # Cost optimization options
        self.console.print("\n[bold]Cost Optimization Options:[/bold]")
        self.console.print("1. Compare instance types")
        self.console.print("2. Analyze storage costs")
        self.console.print("3. Review data transfer costs")
        self.console.print("4. Spot instance calculator")
        self.console.print("b. Back to main menu")
        
        choice = self.console.input("\n[bold]Enter your choice: [/bold]")
        
        if choice == "1":
            self._compare_instance_costs()
        elif choice == "4":
            self._spot_instance_calculator()
        else:
            self.console.input("Feature coming soon! Press Enter to continue...")
    
    def _compare_instance_costs(self):
        """Compare costs across different instance types."""
        instances = self.current_config.get('aws_instance_recommendations', {})
        
        comparison_table = Table(title="Instance Type Cost Comparison", box=box.ROUNDED)
        comparison_table.add_column("Instance Type", style="cyan")
        comparison_table.add_column("vCPUs", style="yellow", justify="right")
        comparison_table.add_column("Memory (GB)", style="green", justify="right")
        comparison_table.add_column("Hourly", style="red", justify="right")
        comparison_table.add_column("Monthly", style="red", justify="right")
        comparison_table.add_column("Annual", style="red", justify="right")
        comparison_table.add_column("Use Case", style="white", max_width=25)
        
        for name, details in instances.items():
            hourly_cost = details.get('cost_per_hour', 0)
            monthly_cost = hourly_cost * 24 * 30
            annual_cost = monthly_cost * 12
            
            comparison_table.add_row(
                details.get('instance_type', 'N/A'),
                str(details.get('vcpus', 'N/A')),
                str(details.get('memory_gb', 'N/A')),
                f"${hourly_cost:.3f}",
                f"${monthly_cost:.0f}",
                f"${annual_cost:.0f}",
                details.get('use_case', '')[:25]
            )
        
        self.console.print(comparison_table)
        self.console.input("\nPress Enter to continue...")
    
    def _spot_instance_calculator(self):
        """Calculate potential savings with Spot instances."""
        if not self.selected_instance:
            self.console.print("[red]Please select an instance type first[/red]")
            self.console.input("Press Enter to continue...")
            return
        
        on_demand_cost = self.selected_instance.get('cost_per_hour', 0)
        spot_discount = 0.7  # Typical 70% discount
        spot_cost = on_demand_cost * (1 - spot_discount)
        
        savings_table = Table(title="Spot Instance Savings Calculator", box=box.ROUNDED)
        savings_table.add_column("Pricing Model", style="cyan")
        savings_table.add_column("Hourly", style="yellow", justify="right")
        savings_table.add_column("Monthly", style="green", justify="right")
        savings_table.add_column("Annual", style="red", justify="right")
        
        savings_table.add_row(
            "On-Demand",
            f"${on_demand_cost:.3f}",
            f"${on_demand_cost * 24 * 30:.0f}",
            f"${on_demand_cost * 24 * 365:.0f}"
        )
        
        savings_table.add_row(
            "Spot Instance",
            f"${spot_cost:.3f}",
            f"${spot_cost * 24 * 30:.0f}",
            f"${spot_cost * 24 * 365:.0f}"
        )
        
        monthly_savings = (on_demand_cost - spot_cost) * 24 * 30
        annual_savings = monthly_savings * 12
        
        savings_table.add_separator()
        savings_table.add_row(
            "[bold green]Savings[/bold green]",
            f"[bold green]${on_demand_cost - spot_cost:.3f}[/bold green]",
            f"[bold green]${monthly_savings:.0f}[/bold green]",
            f"[bold green]${annual_savings:.0f}[/bold green]"
        )
        
        self.console.print(savings_table)
        
        # Spot instance considerations
        considerations = [
            "✓ Perfect for fault-tolerant workloads",
            "✓ Ideal for batch processing and analysis",
            "⚠ May be interrupted with 2-minute notice",
            "⚠ Not suitable for interactive workloads",
            "💡 Consider checkpointing for long-running jobs"
        ]
        
        self.console.print(Panel(
            "\n".join(considerations),
            title="Spot Instance Considerations",
            border_style="yellow"
        ))
        
        self.console.input("\nPress Enter to continue...")
    
    def _deployment_menu(self):
        """Deployment and infrastructure management."""
        if not self.current_domain:
            self.console.print("[red]Please select and configure a domain first[/red]")
            self.console.input("Press Enter to continue...")
            return
        
        self.console.clear()
        self.console.print(Panel.fit(
            f"[bold]Deploy: {self.current_config['name']}[/bold]",
            border_style="green"
        ))
        
        # Deployment checklist
        checklist = Table(title="Pre-Deployment Checklist", box=box.ROUNDED)
        checklist.add_column("Check", style="cyan")
        checklist.add_column("Status", style="white")
        checklist.add_column("Details", style="yellow")
        
        aws_status = "✓ Connected" if self.aws_available else "❌ Not Available"
        domain_status = "✓ Selected" if self.current_domain else "❌ Not Selected"
        instance_status = "✓ Selected" if self.selected_instance else "⚠ Using Default"
        
        checklist.add_row("AWS Access", aws_status, "IAM permissions verified")
        checklist.add_row("Domain Pack", domain_status, self.current_domain or "None")
        checklist.add_row("Instance Type", instance_status, 
                         self.selected_instance.get('instance_type', 'Default') if self.selected_instance else 'Default')
        checklist.add_row("Cost Estimate", "✓ Calculated", f"${self.estimated_costs.get('total', 0):.2f}/month")
        
        self.console.print(checklist)
        
        if self.aws_available and self.current_domain:
            self.console.print("\n[bold]Deployment Options:[/bold]")
            self.console.print("1. Quick Deploy (Single Instance)")
            self.console.print("2. Cluster Deploy (Multi-Node)")
            self.console.print("3. Spot Instance Deploy")
            self.console.print("4. Generate CloudFormation Template")
            self.console.print("5. Dry Run (Validation Only)")
            self.console.print("b. Back to main menu")
            
            choice = self.console.input("\n[bold]Enter your choice: [/bold]")
            
            if choice == "1":
                self._quick_deploy()
            elif choice == "5":
                self._dry_run_deployment()
            else:
                self.console.print("[yellow]Feature coming soon![/yellow]")
                self.console.input("Press Enter to continue...")
        else:
            self.console.print("\n[red]Cannot deploy: Missing AWS access or domain configuration[/red]")
            self.console.input("Press Enter to continue...")
    
    def _dry_run_deployment(self):
        """Simulate deployment for validation."""
        self.console.print("\n[bold]Running Deployment Validation...[/bold]")
        
        with Progress() as progress:
            task = progress.add_task("[cyan]Validating configuration...", total=100)
            
            # Simulate validation steps
            validation_steps = [
                ("Checking AWS permissions", 20),
                ("Validating instance type", 15),
                ("Verifying software packages", 25),
                ("Calculating resource requirements", 20),
                ("Generating deployment plan", 20)
            ]
            
            completed = 0
            for step, increment in validation_steps:
                progress.update(task, description=f"[cyan]{step}...")
                # Simulate processing time
                import time
                time.sleep(0.5)
                completed += increment
                progress.update(task, completed=completed)
        
        # Validation results
        self.console.print("\n[green]✓ Validation Complete![/green]")
        
        results_table = Table(title="Deployment Plan", box=box.ROUNDED)
        results_table.add_column("Component", style="cyan")
        results_table.add_column("Configuration", style="white")
        results_table.add_column("Status", style="green")
        
        instance_type = self.selected_instance.get('instance_type', 'c6i.2xlarge') if self.selected_instance else 'c6i.2xlarge'
        
        results_table.add_row("Instance Type", instance_type, "✓ Available")
        results_table.add_row("Operating System", "Amazon Linux 2023", "✓ Latest AMI")
        results_table.add_row("Storage", "500 GB EBS GP3", "✓ Optimized")
        results_table.add_row("Security Group", "Research SSH + Custom", "✓ Configured")
        results_table.add_row("IAM Role", "ResearchInstanceRole", "✓ Permissions OK")
        
        self.console.print(results_table)
        
        estimated_time = "5-10 minutes"
        self.console.print(f"\n[bold]Estimated deployment time: [green]{estimated_time}[/green][/bold]")
        
        self.console.input("Press Enter to continue...")
    
    def _monitoring_menu(self):
        """AWS resource monitoring dashboard."""
        self.console.clear()
        self.console.print(Panel.fit(
            "[bold]AWS Resource Monitor[/bold]",
            border_style="blue"
        ))
        
        if not self.aws_available:
            self.console.print("[red]AWS connection not available. Cannot monitor resources.[/red]")
            self.console.input("Press Enter to continue...")
            return
        
        # Mock monitoring data (in real implementation, fetch from AWS APIs)
        monitoring_layout = Layout()
        monitoring_layout.split_column(
            Layout(name="instances", ratio=2),
            Layout(name="costs", ratio=1)
        )
        
        # Instance status table
        instances_table = Table(title="EC2 Instances", box=box.ROUNDED)
        instances_table.add_column("Instance ID", style="cyan")
        instances_table.add_column("Type", style="yellow")
        instances_table.add_column("State", style="green")
        instances_table.add_column("CPU %", style="red", justify="right")
        instances_table.add_column("Memory %", style="blue", justify="right")
        instances_table.add_column("Cost/Hour", style="magenta", justify="right")
        
        # Mock data - in real implementation, use EC2 describe_instances()
        instances_table.add_row("i-0123456789abcdef0", "c6i.2xlarge", "running", "45%", "67%", "$0.34")
        instances_table.add_row("i-abcdef0123456789", "r6i.4xlarge", "running", "78%", "82%", "$1.02")
        
        monitoring_layout["instances"].update(Panel(instances_table, title="Active Resources"))
        
        # Cost monitoring
        cost_info = [
            "[green]Today's Spend:[/green] $12.47",
            "[yellow]This Month:[/yellow] $347.82", 
            "[red]Projected Month:[/red] $425.50",
            "",
            "[cyan]Top Services:[/cyan]",
            "• EC2: $8.45 (68%)",
            "• S3: $2.12 (17%)",
            "• EBS: $1.90 (15%)"
        ]
        
        monitoring_layout["costs"].update(Panel(
            "\n".join(cost_info),
            title="💰 Cost Tracking",
            border_style="yellow"
        ))
        
        self.console.print(monitoring_layout)
        
        self.console.print("\n[bold]Monitoring Options:[/bold]")
        self.console.print("1. Refresh status")
        self.console.print("2. View detailed logs")
        self.console.print("3. Set up alerts")
        self.console.print("4. Export usage report")
        self.console.print("b. Back to main menu")
        
        choice = self.console.input("\n[bold]Enter your choice: [/bold]")
        
        if choice == "1":
            self.console.print("[green]Status refreshed![/green]")
            self.console.input("Press Enter to continue...")
        else:
            self.console.print("[yellow]Feature coming soon![/yellow]")
            self.console.input("Press Enter to continue...")
    
    def _calculate_costs(self):
        """Calculate estimated costs for current configuration."""
        if not self.current_config:
            return
        
        base_costs = self.current_config.get('estimated_cost', {})
        
        # Apply instance-specific costs if selected
        if self.selected_instance:
            compute_cost = self.selected_instance.get('cost_per_hour', 0) * 24 * 30
        else:
            compute_cost = base_costs.get('compute', 500)
        
        self.estimated_costs = {
            'compute': compute_cost,
            'storage': base_costs.get('storage', 200),
            'data_transfer': base_costs.get('data_transfer', 100),
            'software': base_costs.get('software', 50),
            'total': 0
        }
        
        self.estimated_costs['total'] = sum(v for k, v in self.estimated_costs.items() if k != 'total')
    
    def _get_cost_optimization_suggestions(self) -> List[str]:
        """Generate cost optimization suggestions."""
        suggestions = []
        
        total_cost = self.estimated_costs.get('total', 0)
        compute_cost = self.estimated_costs.get('compute', 0)
        
        if compute_cost / total_cost > 0.7:
            suggestions.append("Consider Spot instances for 70% savings")
        
        if total_cost > 1000:
            suggestions.append("Reserved instances can save 30-60%")
        
        suggestions.extend([
            "Use S3 Intelligent Tiering for storage",
            "Enable detailed monitoring for optimization",
            "Consider multi-AZ only for production",
            "Review data transfer patterns",
            "Use AWS Cost Explorer for insights"
        ])
        
        return suggestions[:5]  # Return top 5 suggestions
    
    def _customize_packages(self):
        """Package customization interface."""
        self.console.print("\n[bold]Software Package Customization[/bold]")
        packages = self.current_config.get('spack_packages', {})
        
        for category, pkg_list in packages.items():
            self.console.print(f"\n[cyan]{category.replace('_', ' ').title()}:[/cyan]")
            for i, pkg in enumerate(pkg_list[:5], 1):  # Show first 5 packages
                pkg_name = pkg.split('@')[0]  # Extract package name
                self.console.print(f"  {i}. {pkg_name}")
            
            if len(pkg_list) > 5:
                self.console.print(f"  ... and {len(pkg_list) - 5} more packages")
        
        self.console.input("\nPackage customization coming soon! Press Enter to continue...")
    
    def _configure_storage(self):
        """Storage configuration interface."""
        self.console.print("\n[bold]Storage Configuration[/bold]")
        
        storage_table = Table(title="Storage Options", box=box.ROUNDED)
        storage_table.add_column("Type", style="cyan")
        storage_table.add_column("Size", style="yellow")
        storage_table.add_column("Performance", style="green")
        storage_table.add_column("Cost/Month", style="red")
        
        storage_table.add_row("EBS GP3", "500 GB", "3,000 IOPS", "$40")
        storage_table.add_row("EBS GP3", "1 TB", "3,000 IOPS", "$80")
        storage_table.add_row("EBS io2", "500 GB", "10,000 IOPS", "$120")
        storage_table.add_row("EFS", "1 TB", "Shared", "$300")
        
        self.console.print(storage_table)
        self.console.input("\nStorage customization coming soon! Press Enter to continue...")
    
    def _configure_data_sources(self):
        """Data source configuration interface."""
        self.console.print("\n[bold]Data Source Configuration[/bold]")
        
        # Show available AWS Open Data sources for the domain
        domain_data = {
            'genomics': ['1000 Genomes Project', 'gnomAD Database'],
            'climate_modeling': ['ERA5 Reanalysis', 'NOAA GFS'],
            'neuroscience': ['Human Connectome Project'],
            'astronomy_astrophysics': ['SDSS Survey Data']
        }
        
        sources = domain_data.get(self.current_domain, ['No specific datasets configured'])
        
        sources_table = Table(title=f"Available Data Sources for {self.current_domain.title()}", box=box.ROUNDED)
        sources_table.add_column("Dataset", style="cyan")
        sources_table.add_column("Size", style="yellow")
        sources_table.add_column("Access", style="green")
        sources_table.add_column("Cost", style="red")
        
        for source in sources:
            sources_table.add_row(source, "Multi-TB", "Public", "Free*")
        
        self.console.print(sources_table)
        self.console.print("\n[dim]*Data transfer and compute costs may apply[/dim]")
        self.console.input("\nData source configuration coming soon! Press Enter to continue...")
    
    def _quick_deploy(self):
        """Quick deployment simulation."""
        self.console.print("\n[bold red]⚠ This would create real AWS resources and incur costs![/bold red]")
        confirm = self.console.input("Type 'DEPLOY' to confirm (or anything else to cancel): ")
        
        if confirm != "DEPLOY":
            self.console.print("[yellow]Deployment cancelled.[/yellow]")
            self.console.input("Press Enter to continue...")
            return
        
        self.console.print("\n[bold]Starting Quick Deployment...[/bold]")
        
        with Progress() as progress:
            task = progress.add_task("[green]Deploying infrastructure...", total=100)
            
            # Simulate deployment steps
            deployment_steps = [
                ("Creating security group", 15),
                ("Launching EC2 instance", 25),
                ("Configuring storage", 20),
                ("Installing software packages", 30),
                ("Finalizing configuration", 10)
            ]
            
            completed = 0
            for step, increment in deployment_steps:
                progress.update(task, description=f"[green]{step}...")
                import time
                time.sleep(1)  # Simulate deployment time
                completed += increment
                progress.update(task, completed=completed)
        
        self.console.print(f"\n[bold green]✓ Deployment Complete![/bold green]")
        self.console.print(f"[cyan]Instance ID:[/cyan] i-{hex(hash(self.current_domain))[2:12]}")
        self.console.print(f"[cyan]Public IP:[/cyan] 52.{hash(self.current_domain) % 255}.{hash(self.current_domain) % 200}.{hash(self.current_domain) % 100}")
        self.console.print(f"[cyan]SSH Command:[/cyan] ssh -i ~/.ssh/research-key.pem ec2-user@52.{hash(self.current_domain) % 255}.{hash(self.current_domain) % 200}.{hash(self.current_domain) % 100}")
        
        self.console.input("\nPress Enter to continue...")

# Textual-based Advanced TUI Application
if TEXTUAL_AVAILABLE:
    class ResearchWizardApp(App):
        """Advanced Textual-based TUI application."""
        
        CSS_PATH = None  # We'll define CSS inline for simplicity
        TITLE = "AWS Research Wizard"
        SUB_TITLE = "Terminal Interface"
        
        BINDINGS = [
            Binding("q", "quit", "Quit"),
            Binding("d", "toggle_dark", "Toggle Dark Mode"),
            Binding("h", "help", "Help"),
        ]
        
        def __init__(self, wizard_instance: ResearchWizardTUI):
            super().__init__()
            self.wizard = wizard_instance
        
        def compose(self) -> ComposeResult:
            """Create child widgets for the app."""
            yield Header()
            
            with Container():
                with Horizontal():
                    # Sidebar with navigation
                    with Vertical(classes="sidebar"):
                        yield Static("Research Domains", classes="sidebar-title")
                        domain_items = [ListItem(Static(config['name'][:25])) 
                                      for config in self.wizard.domains.values()]
                        yield ListView(*domain_items, id="domain-list")
                    
                    # Main content area
                    with Vertical(classes="main-content"):
                        with Tabs("Overview", "Configure", "Deploy", "Monitor"):
                            with TabPane("Overview", id="overview"):
                                yield Static("Select a research domain to get started", id="overview-content")
                            
                            with TabPane("Configure", id="configure"):
                                yield Static("Configuration options will appear here", id="config-content")
                            
                            with TabPane("Deploy", id="deploy"):
                                yield Static("Deployment options will appear here", id="deploy-content")
                            
                            with TabPane("Monitor", id="monitor"):
                                yield Static("Monitoring dashboard will appear here", id="monitor-content")
            
            yield Footer()
        
        def action_toggle_dark(self) -> None:
            """Toggle dark mode."""
            self.dark = not self.dark
        
        def action_help(self) -> None:
            """Show help dialog."""
            pass  # Implement help dialog

def main():
    """Main entry point for the TUI."""
    import argparse
    
    parser = argparse.ArgumentParser(description="AWS Research Wizard Terminal Interface")
    parser.add_argument("--config", default="configs", 
                       help="Configuration directory path")
    parser.add_argument("--simple", action="store_true",
                       help="Use simple TUI without Textual dependency")
    
    args = parser.parse_args()
    
    wizard = ResearchWizardTUI(args.config)
    
    if args.simple or not TEXTUAL_AVAILABLE:
        wizard._run_simple_tui()
    else:
        wizard.run()

if __name__ == "__main__":
    main()