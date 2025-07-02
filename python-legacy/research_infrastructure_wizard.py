#!/usr/bin/env python3
"""
Research Infrastructure Wizard
Domain-aware intelligent infrastructure recommendations for research workloads
"""

import json
import math
from typing import Dict, List, Tuple, Optional, Any
from dataclasses import dataclass, asdict
from enum import Enum
import argparse
import logging

class Priority(Enum):
    COST = "cost"
    PERFORMANCE = "performance"
    DEADLINE = "deadline"
    BALANCED = "balanced"

class WorkloadSize(Enum):
    SMALL = "small"
    MEDIUM = "medium"
    LARGE = "large"
    MASSIVE = "massive"

@dataclass
class WorkloadCharacteristics:
    domain: str
    primary_tools: List[str]
    problem_size: WorkloadSize
    priority: Priority
    deadline_hours: Optional[int]
    budget_limit: Optional[float]
    data_size_gb: int
    parallel_scaling: str  # "none", "linear", "sublinear", "embarrassing"
    gpu_requirement: str   # "none", "optional", "required", "multi_gpu"
    memory_intensity: str  # "low", "medium", "high", "extreme"
    io_pattern: str       # "sequential", "random", "streaming", "burst"
    collaboration_users: int

@dataclass
class InfrastructureRecommendation:
    instance_type: str
    instance_count: int
    storage_config: Dict[str, Any]
    network_config: Dict[str, Any]
    estimated_cost: Dict[str, float]
    estimated_runtime: Dict[str, float]
    optimization_rationale: List[str]
    alternative_configs: List[Dict[str, Any]]
    deployment_template: str
    monitoring_setup: Dict[str, Any]

class ResearchInfrastructureWizard:
    def __init__(self):
        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)

        # Domain-specific knowledge base
        self.domain_profiles = self._initialize_domain_profiles()
        self.instance_catalog = self._initialize_instance_catalog()
        self.storage_profiles = self._initialize_storage_profiles()

    def _initialize_domain_profiles(self) -> Dict[str, Dict[str, Any]]:
        """Initialize domain-specific performance and scaling characteristics"""
        return {
            "genomics": {
                "typical_tools": {
                    "gatk": {"cpu_intensive": True, "memory_scaling": "linear", "io_heavy": True},
                    "bwa": {"cpu_intensive": True, "memory_scaling": "sublinear", "io_heavy": True},
                    "star": {"memory_intensive": True, "memory_scaling": "step", "io_heavy": True},
                    "blast": {"cpu_intensive": True, "memory_scaling": "none", "io_heavy": False},
                    "canu": {"memory_intensive": True, "memory_scaling": "linear", "io_heavy": True}
                },
                "scaling_characteristics": {
                    "parallel_efficiency": 0.8,  # How well it scales to multiple cores
                    "memory_overhead": 1.2,      # Memory overhead factor
                    "io_bandwidth_needs": "high",
                    "network_sensitivity": "low"
                },
                "problem_size_indicators": {
                    "small": {"samples": "<10", "genome_size": "<500M", "data_gb": "<50"},
                    "medium": {"samples": "10-100", "genome_size": "500M-3G", "data_gb": "50-500"},
                    "large": {"samples": "100-1000", "genome_size": "3G-10G", "data_gb": "500-5000"},
                    "massive": {"samples": ">1000", "genome_size": ">10G", "data_gb": ">5000"}
                }
            },

            "climate_modeling": {
                "typical_tools": {
                    "wrf": {"cpu_intensive": True, "mpi_scaling": "good", "memory_scaling": "linear"},
                    "cesm": {"cpu_intensive": True, "mpi_scaling": "excellent", "memory_scaling": "linear"},
                    "nco": {"io_intensive": True, "cpu_intensive": False, "memory_scaling": "none"},
                    "cdo": {"cpu_intensive": True, "memory_scaling": "sublinear", "io_heavy": True}
                },
                "scaling_characteristics": {
                    "parallel_efficiency": 0.9,
                    "memory_overhead": 1.1,
                    "io_bandwidth_needs": "very_high",
                    "network_sensitivity": "high"
                },
                "problem_size_indicators": {
                    "small": {"resolution": ">50km", "domain": "regional", "years": "<5"},
                    "medium": {"resolution": "10-50km", "domain": "continental", "years": "5-20"},
                    "large": {"resolution": "1-10km", "domain": "global", "years": "20-100"},
                    "massive": {"resolution": "<1km", "domain": "global", "years": ">100"}
                }
            },

            "ai_ml": {
                "typical_tools": {
                    "pytorch": {"gpu_required": True, "memory_scaling": "batch_size", "distributed": True},
                    "tensorflow": {"gpu_preferred": True, "memory_scaling": "batch_size", "distributed": True},
                    "xgboost": {"cpu_intensive": True, "memory_scaling": "linear", "gpu_optional": True},
                    "scikit_learn": {"cpu_intensive": True, "memory_scaling": "linear", "gpu_none": True}
                },
                "scaling_characteristics": {
                    "gpu_scaling_efficiency": 0.85,
                    "cpu_scaling_efficiency": 0.7,
                    "memory_overhead": 2.0,  # Higher due to model caching
                    "io_bandwidth_needs": "medium",
                    "network_sensitivity": "very_high"  # For distributed training
                },
                "problem_size_indicators": {
                    "small": {"parameters": "<1M", "dataset": "<1GB", "training_time": "<1h"},
                    "medium": {"parameters": "1M-100M", "dataset": "1-100GB", "training_time": "1-24h"},
                    "large": {"parameters": "100M-10B", "dataset": "100GB-1TB", "training_time": "1-7d"},
                    "massive": {"parameters": ">10B", "dataset": ">1TB", "training_time": ">7d"}
                }
            },

            "materials_science": {
                "typical_tools": {
                    "vasp": {"cpu_intensive": True, "memory_intensive": True, "mpi_scaling": "good"},
                    "quantum_espresso": {"cpu_intensive": True, "mpi_scaling": "excellent"},
                    "lammps": {"cpu_intensive": True, "gpu_optional": True, "mpi_scaling": "excellent"},
                    "gromacs": {"cpu_intensive": True, "gpu_preferred": True, "mpi_scaling": "good"}
                },
                "scaling_characteristics": {
                    "parallel_efficiency": 0.85,
                    "memory_overhead": 1.3,
                    "io_bandwidth_needs": "medium",
                    "network_sensitivity": "medium"
                },
                "problem_size_indicators": {
                    "small": {"atoms": "<500", "timesteps": "<100k", "methods": "DFT"},
                    "medium": {"atoms": "500-5000", "timesteps": "100k-1M", "methods": "DFT+MD"},
                    "large": {"atoms": "5000-50000", "timesteps": "1M-10M", "methods": "MD"},
                    "massive": {"atoms": ">50000", "timesteps": ">10M", "methods": "coarse_grain"}
                }
            },

            "neuroscience": {
                "typical_tools": {
                    "fsl": {"cpu_intensive": True, "memory_scaling": "linear", "gpu_optional": True},
                    "freesurfer": {"cpu_intensive": True, "memory_intensive": True, "gpu_none": True},
                    "spm": {"memory_intensive": True, "cpu_moderate": True, "gpu_none": True},
                    "mne": {"cpu_intensive": True, "memory_scaling": "linear", "gpu_optional": True}
                },
                "scaling_characteristics": {
                    "parallel_efficiency": 0.6,  # Many neuroimaging tools don't parallelize well
                    "memory_overhead": 1.5,
                    "io_bandwidth_needs": "high",
                    "network_sensitivity": "low"
                },
                "problem_size_indicators": {
                    "small": {"subjects": "<20", "resolution": "2mm", "sessions": "1"},
                    "medium": {"subjects": "20-100", "resolution": "1mm", "sessions": "1-3"},
                    "large": {"subjects": "100-500", "resolution": "0.5mm", "sessions": "3-10"},
                    "massive": {"subjects": ">500", "resolution": "<0.5mm", "sessions": ">10"}
                }
            }
        }

    def _initialize_instance_catalog(self) -> Dict[str, Dict[str, Any]]:
        """AWS instance catalog with performance characteristics"""
        return {
            # General Purpose - Graviton3
            "c7g.large": {
                "vcpus": 2, "memory": 4, "network": "up_to_12.5", "cost_hour": 0.0363,
                "architecture": "arm64", "optimized_for": ["single_threaded", "cost_sensitive"],
                "graviton_boost": 1.2
            },
            "c7g.xlarge": {
                "vcpus": 4, "memory": 8, "network": "up_to_12.5", "cost_hour": 0.0725,
                "architecture": "arm64", "optimized_for": ["light_parallel", "cost_sensitive"],
                "graviton_boost": 1.2
            },
            "c7g.2xlarge": {
                "vcpus": 8, "memory": 16, "network": "up_to_12.5", "cost_hour": 0.1451,
                "architecture": "arm64", "optimized_for": ["moderate_parallel", "balanced"],
                "graviton_boost": 1.2
            },
            "c7g.4xlarge": {
                "vcpus": 16, "memory": 32, "network": "up_to_12.5", "cost_hour": 0.2902,
                "architecture": "arm64", "optimized_for": ["parallel_workloads", "balanced"],
                "graviton_boost": 1.2
            },
            "c7g.8xlarge": {
                "vcpus": 32, "memory": 64, "network": "12.5", "cost_hour": 0.5803,
                "architecture": "arm64", "optimized_for": ["highly_parallel", "performance"],
                "graviton_boost": 1.2
            },
            "c7g.16xlarge": {
                "vcpus": 64, "memory": 128, "network": "25", "cost_hour": 1.1606,
                "architecture": "arm64", "optimized_for": ["massively_parallel", "performance"],
                "graviton_boost": 1.2
            },

            # Compute Optimized - Intel
            "c6i.large": {
                "vcpus": 2, "memory": 4, "network": "up_to_12.5", "cost_hour": 0.085,
                "architecture": "x86_64", "optimized_for": ["legacy_software", "x86_required"],
                "graviton_boost": 1.0
            },
            "c6i.xlarge": {
                "vcpus": 4, "memory": 8, "network": "up_to_12.5", "cost_hour": 0.17,
                "architecture": "x86_64", "optimized_for": ["legacy_software", "x86_required"],
                "graviton_boost": 1.0
            },
            "c6i.2xlarge": {
                "vcpus": 8, "memory": 16, "network": "up_to_12.5", "cost_hour": 0.34,
                "architecture": "x86_64", "optimized_for": ["legacy_software", "x86_required"],
                "graviton_boost": 1.0
            },
            "c6i.4xlarge": {
                "vcpus": 16, "memory": 32, "network": "up_to_12.5", "cost_hour": 0.68,
                "architecture": "x86_64", "optimized_for": ["legacy_software", "x86_required"],
                "graviton_boost": 1.0
            },
            "c6i.8xlarge": {
                "vcpus": 32, "memory": 64, "network": "12.5", "cost_hour": 1.36,
                "architecture": "x86_64", "optimized_for": ["legacy_software", "performance"],
                "graviton_boost": 1.0
            },

            # Memory Optimized - Graviton3
            "r7g.large": {
                "vcpus": 2, "memory": 16, "network": "up_to_12.5", "cost_hour": 0.0504,
                "architecture": "arm64", "optimized_for": ["memory_intensive", "cost_sensitive"],
                "graviton_boost": 1.2
            },
            "r7g.xlarge": {
                "vcpus": 4, "memory": 32, "network": "up_to_12.5", "cost_hour": 0.1008,
                "architecture": "arm64", "optimized_for": ["memory_intensive", "balanced"],
                "graviton_boost": 1.2
            },
            "r7g.2xlarge": {
                "vcpus": 8, "memory": 64, "network": "up_to_12.5", "cost_hour": 0.2016,
                "architecture": "arm64", "optimized_for": ["memory_intensive", "balanced"],
                "graviton_boost": 1.2
            },
            "r7g.4xlarge": {
                "vcpus": 16, "memory": 128, "network": "up_to_12.5", "cost_hour": 0.4032,
                "architecture": "arm64", "optimized_for": ["memory_intensive", "performance"],
                "graviton_boost": 1.2
            },
            "r7g.8xlarge": {
                "vcpus": 32, "memory": 256, "network": "12.5", "cost_hour": 0.8064,
                "architecture": "arm64", "optimized_for": ["memory_intensive", "performance"],
                "graviton_boost": 1.2
            },
            "r7g.16xlarge": {
                "vcpus": 64, "memory": 512, "network": "25", "cost_hour": 1.6128,
                "architecture": "arm64", "optimized_for": ["memory_intensive", "performance"],
                "graviton_boost": 1.2
            },

            # High Memory
            "x2gd.medium": {
                "vcpus": 1, "memory": 16, "network": "up_to_10", "cost_hour": 0.084,
                "architecture": "arm64", "optimized_for": ["extreme_memory", "small_scale"],
                "graviton_boost": 1.2
            },
            "x2gd.large": {
                "vcpus": 2, "memory": 32, "network": "up_to_10", "cost_hour": 0.168,
                "architecture": "arm64", "optimized_for": ["extreme_memory", "moderate_scale"],
                "graviton_boost": 1.2
            },
            "x2gd.xlarge": {
                "vcpus": 4, "memory": 64, "network": "up_to_10", "cost_hour": 0.336,
                "architecture": "arm64", "optimized_for": ["extreme_memory", "balanced"],
                "graviton_boost": 1.2
            },

            # GPU Instances
            "g5.xlarge": {
                "vcpus": 4, "memory": 16, "gpu": "A10G", "gpu_memory": 24, "network": "up_to_10",
                "cost_hour": 1.006, "architecture": "x86_64",
                "optimized_for": ["ml_inference", "light_training"],
                "graviton_boost": 1.0
            },
            "g5.2xlarge": {
                "vcpus": 8, "memory": 32, "gpu": "A10G", "gpu_memory": 24, "network": "up_to_10",
                "cost_hour": 1.212, "architecture": "x86_64",
                "optimized_for": ["ml_training", "computer_vision"],
                "graviton_boost": 1.0
            },
            "g5.4xlarge": {
                "vcpus": 16, "memory": 64, "gpu": "A10G", "gpu_memory": 24, "network": "up_to_25",
                "cost_hour": 1.624, "architecture": "x86_64",
                "optimized_for": ["ml_training", "multi_modal"],
                "graviton_boost": 1.0
            },
            "g5.12xlarge": {
                "vcpus": 48, "memory": 192, "gpu": "A10G_x4", "gpu_memory": 96, "network": "50",
                "cost_hour": 5.672, "architecture": "x86_64",
                "optimized_for": ["distributed_training", "large_models"],
                "graviton_boost": 1.0
            },
            "p4d.24xlarge": {
                "vcpus": 96, "memory": 1152, "gpu": "A100_x8", "gpu_memory": 320, "network": "400",
                "cost_hour": 32.77, "architecture": "x86_64",
                "optimized_for": ["large_scale_training", "research"],
                "graviton_boost": 1.0
            },
            "p5.48xlarge": {
                "vcpus": 192, "memory": 2048, "gpu": "H100_x8", "gpu_memory": 640, "network": "3200",
                "cost_hour": 98.32, "architecture": "x86_64",
                "optimized_for": ["massive_training", "frontier_research"],
                "graviton_boost": 1.0
            },

            # HPC Optimized
            "hpc7g.4xlarge": {
                "vcpus": 16, "memory": 128, "network": "200", "cost_hour": 1.68,
                "architecture": "arm64", "optimized_for": ["hpc_workloads", "mpi_scaling"],
                "graviton_boost": 1.3, "hpc_optimized": True
            },
            "hpc7g.8xlarge": {
                "vcpus": 32, "memory": 256, "network": "200", "cost_hour": 3.36,
                "architecture": "arm64", "optimized_for": ["hpc_workloads", "mpi_scaling"],
                "graviton_boost": 1.3, "hpc_optimized": True
            },
            "hpc7g.16xlarge": {
                "vcpus": 64, "memory": 512, "network": "200", "cost_hour": 6.72,
                "architecture": "arm64", "optimized_for": ["hpc_workloads", "mpi_scaling"],
                "graviton_boost": 1.3, "hpc_optimized": True
            }
        }

    def _initialize_storage_profiles(self) -> Dict[str, Dict[str, Any]]:
        """Storage configuration profiles for different workload patterns"""
        return {
            "genomics": {
                "primary": {"type": "gp3", "iops": 16000, "throughput": 1000},
                "scratch": {"type": "gp3", "iops": 3000, "throughput": 125},
                "archive": {"type": "s3_ia", "lifecycle": "30_days"}
            },
            "climate_modeling": {
                "primary": {"type": "gp3", "iops": 16000, "throughput": 1000},
                "scratch": {"type": "instance_store", "raid": 0},
                "archive": {"type": "s3_glacier", "lifecycle": "90_days"}
            },
            "ai_ml": {
                "primary": {"type": "gp3", "iops": 16000, "throughput": 1000},
                "datasets": {"type": "fsx_lustre", "throughput": "50_mb_per_tb"},
                "models": {"type": "efs", "throughput": "provisioned"}
            },
            "materials_science": {
                "primary": {"type": "gp3", "iops": 12000, "throughput": 500},
                "scratch": {"type": "instance_store", "raid": 0},
                "archive": {"type": "s3_standard", "lifecycle": "immediate"}
            }
        }

    def create_interactive_wizard(self) -> WorkloadCharacteristics:
        """Interactive wizard for gathering workload requirements"""
        print("🧙‍♂️ Research Infrastructure Wizard")
        print("=====================================")
        print()

        # Domain selection
        domains = list(self.domain_profiles.keys())
        print("1. What is your research domain?")
        for i, domain in enumerate(domains, 1):
            print(f"   {i}. {domain.replace('_', ' ').title()}")

        while True:
            try:
                choice = int(input("\nSelect domain (1-{}): ".format(len(domains))))
                if 1 <= choice <= len(domains):
                    domain = domains[choice - 1]
                    break
            except ValueError:
                pass
            print("Please enter a valid number.")

        # Tool selection
        tools = list(self.domain_profiles[domain]["typical_tools"].keys())
        print(f"\n2. Which tools will you primarily use? (select all that apply)")
        for i, tool in enumerate(tools, 1):
            print(f"   {i}. {tool.upper()}")
        print(f"   {len(tools) + 1}. Other/Multiple")

        primary_tools = []
        tool_input = input(f"\nEnter numbers separated by commas (e.g., 1,2,3): ")
        for choice in tool_input.split(','):
            try:
                idx = int(choice.strip()) - 1
                if 0 <= idx < len(tools):
                    primary_tools.append(tools[idx])
                elif idx == len(tools):
                    primary_tools.append("multiple")
            except ValueError:
                pass

        if not primary_tools:
            primary_tools = [tools[0]]  # Default to first tool

        # Problem size
        print(f"\n3. What is your problem size/scale?")
        size_info = self.domain_profiles[domain]["problem_size_indicators"]
        sizes = list(size_info.keys())
        for i, size in enumerate(sizes, 1):
            indicators = size_info[size]
            indicator_str = ", ".join([f"{k}: {v}" for k, v in indicators.items()])
            print(f"   {i}. {size.title()} ({indicator_str})")

        while True:
            try:
                choice = int(input(f"\nSelect size (1-{len(sizes)}): "))
                if 1 <= choice <= len(sizes):
                    problem_size = WorkloadSize(sizes[choice - 1])
                    break
            except ValueError:
                pass
            print("Please enter a valid number.")

        # Priority
        print(f"\n4. What is your primary optimization priority?")
        priorities = ["Cost (minimize expenses)", "Performance (fastest results)",
                     "Deadline (meet specific timeline)", "Balanced (good compromise)"]
        for i, priority in enumerate(priorities, 1):
            print(f"   {i}. {priority}")

        while True:
            try:
                choice = int(input(f"\nSelect priority (1-{len(priorities)}): "))
                if 1 <= choice <= len(priorities):
                    priority_map = [Priority.COST, Priority.PERFORMANCE, Priority.DEADLINE, Priority.BALANCED]
                    priority = priority_map[choice - 1]
                    break
            except ValueError:
                pass
            print("Please enter a valid number.")

        # Additional parameters
        deadline_hours = None
        if priority == Priority.DEADLINE:
            while True:
                try:
                    deadline_hours = int(input("\nDeadline in hours from now: "))
                    if deadline_hours > 0:
                        break
                except ValueError:
                    pass
                print("Please enter a positive number of hours.")

        budget_limit = None
        if priority in [Priority.COST, Priority.BALANCED]:
            try:
                budget_str = input("\nMaximum budget in USD (press Enter to skip): ")
                if budget_str.strip():
                    budget_limit = float(budget_str.replace('$', '').replace(',', ''))
            except ValueError:
                pass

        # Data size
        while True:
            try:
                data_size_gb = int(input("\nApproximate data size in GB: "))
                if data_size_gb >= 0:
                    break
            except ValueError:
                pass
            print("Please enter a valid number.")

        # Collaboration
        while True:
            try:
                collaboration_users = int(input("\nNumber of concurrent users (including you): "))
                if collaboration_users >= 1:
                    break
            except ValueError:
                pass
            print("Please enter a number >= 1.")

        # Auto-determine technical characteristics based on domain and tools
        scaling_char = self.domain_profiles[domain]["scaling_characteristics"]

        # Determine parallel scaling
        if any(tool in ["mpi", "distributed"] for tool in primary_tools):
            parallel_scaling = "linear"
        elif scaling_char["parallel_efficiency"] > 0.8:
            parallel_scaling = "sublinear"
        elif any(tool in self.domain_profiles[domain]["typical_tools"]
                for tool in primary_tools
                if self.domain_profiles[domain]["typical_tools"].get(tool, {}).get("cpu_intensive", False)):
            parallel_scaling = "embarrassing"
        else:
            parallel_scaling = "none"

        # Determine GPU requirement
        gpu_requirement = "none"
        for tool in primary_tools:
            tool_info = self.domain_profiles[domain]["typical_tools"].get(tool, {})
            if tool_info.get("gpu_required", False):
                gpu_requirement = "required"
                break
            elif tool_info.get("gpu_preferred", False):
                gpu_requirement = "optional"

        # Determine memory intensity
        memory_intensity = "medium"
        for tool in primary_tools:
            tool_info = self.domain_profiles[domain]["typical_tools"].get(tool, {})
            if tool_info.get("memory_intensive", False):
                memory_intensity = "high"
                break
            elif tool_info.get("memory_scaling", "") == "linear":
                memory_intensity = "medium"

        # Determine I/O pattern
        io_pattern = "sequential"
        if domain in ["genomics", "climate_modeling"]:
            io_pattern = "streaming"
        elif domain == "ai_ml":
            io_pattern = "burst"

        return WorkloadCharacteristics(
            domain=domain,
            primary_tools=primary_tools,
            problem_size=problem_size,
            priority=priority,
            deadline_hours=deadline_hours,
            budget_limit=budget_limit,
            data_size_gb=data_size_gb,
            parallel_scaling=parallel_scaling,
            gpu_requirement=gpu_requirement,
            memory_intensity=memory_intensity,
            io_pattern=io_pattern,
            collaboration_users=collaboration_users
        )

    def recommend_infrastructure(self, workload: WorkloadCharacteristics) -> InfrastructureRecommendation:
        """Generate optimal infrastructure recommendation based on workload characteristics"""

        # Filter instances based on requirements
        candidate_instances = self._filter_instances(workload)

        # Score instances based on priority
        scored_instances = self._score_instances(candidate_instances, workload)

        # Select optimal configuration
        optimal_config = self._select_optimal_config(scored_instances, workload)

        # Generate storage configuration
        storage_config = self._generate_storage_config(workload)

        # Generate network configuration
        network_config = self._generate_network_config(workload)

        # Calculate costs and runtime estimates
        cost_estimates = self._calculate_costs(optimal_config, workload)
        runtime_estimates = self._estimate_runtime(optimal_config, workload)

        # Generate alternatives
        alternatives = self._generate_alternatives(scored_instances, workload)

        # Create deployment template
        deployment_template = self._generate_deployment_template(optimal_config, workload)

        # Generate monitoring setup
        monitoring_setup = self._generate_monitoring_config(workload)

        # Generate optimization rationale
        rationale = self._generate_rationale(optimal_config, workload)

        return InfrastructureRecommendation(
            instance_type=optimal_config["instance_type"],
            instance_count=optimal_config["instance_count"],
            storage_config=storage_config,
            network_config=network_config,
            estimated_cost=cost_estimates,
            estimated_runtime=runtime_estimates,
            optimization_rationale=rationale,
            alternative_configs=alternatives,
            deployment_template=deployment_template,
            monitoring_setup=monitoring_setup
        )

    def _filter_instances(self, workload: WorkloadCharacteristics) -> List[Dict[str, Any]]:
        """Filter instances based on workload requirements"""
        candidates = []

        for instance_type, specs in self.instance_catalog.items():
            # GPU requirements
            if workload.gpu_requirement == "required" and "gpu" not in specs:
                continue
            if workload.gpu_requirement == "none" and "gpu" in specs:
                continue

            # Memory requirements
            min_memory = self._estimate_memory_requirement(workload)
            if specs["memory"] < min_memory:
                continue

            # CPU requirements
            min_cpus = self._estimate_cpu_requirement(workload)
            if specs["vcpus"] < min_cpus:
                continue

            candidates.append({
                "instance_type": instance_type,
                "specs": specs
            })

        return candidates

    def _estimate_memory_requirement(self, workload: WorkloadCharacteristics) -> int:
        """Estimate minimum memory requirement in GB"""
        base_memory = {
            WorkloadSize.SMALL: 4,
            WorkloadSize.MEDIUM: 16,
            WorkloadSize.LARGE: 64,
            WorkloadSize.MASSIVE: 256
        }[workload.problem_size]

        # Domain-specific multipliers
        domain_multiplier = {
            "genomics": 2.0,
            "climate_modeling": 1.5,
            "ai_ml": 3.0,
            "materials_science": 1.8,
            "neuroscience": 2.5
        }.get(workload.domain, 1.0)

        # Memory intensity multiplier
        intensity_multiplier = {
            "low": 0.5,
            "medium": 1.0,
            "high": 2.0,
            "extreme": 4.0
        }[workload.memory_intensity]

        return int(base_memory * domain_multiplier * intensity_multiplier)

    def _estimate_cpu_requirement(self, workload: WorkloadCharacteristics) -> int:
        """Estimate minimum CPU requirement"""
        base_cpus = {
            WorkloadSize.SMALL: 2,
            WorkloadSize.MEDIUM: 8,
            WorkloadSize.LARGE: 32,
            WorkloadSize.MASSIVE: 64
        }[workload.problem_size]

        # Adjust for parallel scaling
        if workload.parallel_scaling == "none":
            return min(base_cpus, 4)
        elif workload.parallel_scaling == "sublinear":
            return min(base_cpus, 16)

        return base_cpus

    def _score_instances(self, candidates: List[Dict], workload: WorkloadCharacteristics) -> List[Dict]:
        """Score instances based on optimization priority"""
        scored = []

        for candidate in candidates:
            specs = candidate["specs"]
            score = 0

            # Performance score
            perf_score = self._calculate_performance_score(specs, workload)

            # Cost efficiency score
            cost_score = self._calculate_cost_efficiency_score(specs, workload)

            # Graviton boost
            graviton_boost = specs.get("graviton_boost", 1.0)

            # Priority weighting
            if workload.priority == Priority.COST:
                score = cost_score * 0.7 + perf_score * 0.3
            elif workload.priority == Priority.PERFORMANCE:
                score = perf_score * 0.8 + cost_score * 0.2
            elif workload.priority == Priority.DEADLINE:
                score = perf_score * 0.9 + cost_score * 0.1
            else:  # BALANCED
                score = perf_score * 0.5 + cost_score * 0.5

            # Apply Graviton boost
            if specs["architecture"] == "arm64":
                score *= graviton_boost

            scored.append({
                **candidate,
                "score": score,
                "perf_score": perf_score,
                "cost_score": cost_score
            })

        return sorted(scored, key=lambda x: x["score"], reverse=True)

    def _calculate_performance_score(self, specs: Dict, workload: WorkloadCharacteristics) -> float:
        """Calculate performance score for instance"""
        # Base performance from CPU count and memory
        cpu_score = specs["vcpus"] * 100
        memory_score = specs["memory"] * 10

        # GPU boost if needed
        gpu_score = 0
        if workload.gpu_requirement in ["required", "optional"] and "gpu" in specs:
            gpu_score = specs.get("gpu_memory", 0) * 50

        # Network performance for distributed workloads
        network_score = 0
        if workload.collaboration_users > 1 or workload.parallel_scaling in ["linear", "sublinear"]:
            network_multiplier = float(specs["network"].replace("up_to_", "").replace("up to ", ""))
            network_score = network_multiplier * 20

        return cpu_score + memory_score + gpu_score + network_score

    def _calculate_cost_efficiency_score(self, specs: Dict, workload: WorkloadCharacteristics) -> float:
        """Calculate cost efficiency score (higher is better)"""
        hourly_cost = specs["cost_hour"]
        performance = specs["vcpus"] * specs["memory"]

        if "gpu" in specs:
            performance *= specs.get("gpu_memory", 1) * 2

        # Cost efficiency (performance per dollar)
        efficiency = performance / hourly_cost

        return efficiency

    def _select_optimal_config(self, scored_instances: List[Dict], workload: WorkloadCharacteristics) -> Dict:
        """Select optimal instance configuration"""
        if not scored_instances:
            # Fallback to a reasonable default
            return {
                "instance_type": "c7g.2xlarge",
                "instance_count": 1
            }

        best_instance = scored_instances[0]

        # Determine instance count for distributed workloads
        instance_count = 1
        if workload.parallel_scaling in ["linear", "sublinear"] and workload.problem_size in [WorkloadSize.LARGE, WorkloadSize.MASSIVE]:
            if workload.problem_size == WorkloadSize.LARGE:
                instance_count = min(8, workload.collaboration_users * 2)
            elif workload.problem_size == WorkloadSize.MASSIVE:
                instance_count = min(32, workload.collaboration_users * 4)

        return {
            "instance_type": best_instance["instance_type"],
            "instance_count": instance_count,
            "specs": best_instance["specs"]
        }

    def _generate_storage_config(self, workload: WorkloadCharacteristics) -> Dict[str, Any]:
        """Generate storage configuration based on workload"""
        domain_storage = self.storage_profiles.get(workload.domain, self.storage_profiles["genomics"])

        # Calculate storage sizes
        primary_size = max(100, workload.data_size_gb * 2)  # 2x for processing overhead
        scratch_size = max(50, workload.data_size_gb)

        return {
            "primary": {
                **domain_storage["primary"],
                "size_gb": primary_size
            },
            "scratch": {
                **domain_storage["scratch"],
                "size_gb": scratch_size
            },
            "archive": domain_storage["archive"]
        }

    def _generate_network_config(self, workload: WorkloadCharacteristics) -> Dict[str, Any]:
        """Generate network configuration"""
        return {
            "vpc": {
                "enable_dns_hostnames": True,
                "enable_dns_support": True
            },
            "placement_group": workload.parallel_scaling in ["linear", "sublinear"],
            "enhanced_networking": True,
            "sr_iov": True
        }

    def _calculate_costs(self, config: Dict, workload: WorkloadCharacteristics) -> Dict[str, float]:
        """Calculate cost estimates"""
        specs = config["specs"]
        instance_count = config["instance_count"]

        hourly_cost = specs["cost_hour"] * instance_count

        # Storage costs (monthly)
        storage_cost_monthly = (
            workload.data_size_gb * 0.1 +  # Primary storage (gp3)
            workload.data_size_gb * 0.05   # Scratch storage
        )

        # Estimate runtime based on problem size and tools
        estimated_hours = {
            WorkloadSize.SMALL: 2,
            WorkloadSize.MEDIUM: 8,
            WorkloadSize.LARGE: 24,
            WorkloadSize.MASSIVE: 72
        }[workload.problem_size]

        # Adjust for Graviton boost
        if specs["architecture"] == "arm64":
            estimated_hours *= 0.8  # 20% performance improvement

        total_cost = hourly_cost * estimated_hours + storage_cost_monthly

        return {
            "hourly_compute": hourly_cost,
            "estimated_runtime_hours": estimated_hours,
            "total_compute": hourly_cost * estimated_hours,
            "storage_monthly": storage_cost_monthly,
            "total_estimated": total_cost
        }

    def _estimate_runtime(self, config: Dict, workload: WorkloadCharacteristics) -> Dict[str, float]:
        """Estimate runtime characteristics"""
        base_runtime = {
            WorkloadSize.SMALL: 2,
            WorkloadSize.MEDIUM: 8,
            WorkloadSize.LARGE: 24,
            WorkloadSize.MASSIVE: 72
        }[workload.problem_size]

        # Adjust for parallel scaling
        if workload.parallel_scaling == "linear" and config["instance_count"] > 1:
            scaling_efficiency = 0.9
            speedup = config["instance_count"] * scaling_efficiency
            runtime = base_runtime / speedup
        elif workload.parallel_scaling == "sublinear" and config["instance_count"] > 1:
            scaling_efficiency = 0.7
            speedup = math.sqrt(config["instance_count"]) * scaling_efficiency
            runtime = base_runtime / speedup
        else:
            runtime = base_runtime

        # Graviton performance boost
        if config["specs"]["architecture"] == "arm64":
            runtime *= 0.8

        return {
            "estimated_hours": runtime,
            "confidence": "medium"
        }

    def _generate_alternatives(self, scored_instances: List[Dict], workload: WorkloadCharacteristics) -> List[Dict]:
        """Generate alternative configurations"""
        alternatives = []

        # Cost-optimized alternative
        if workload.priority != Priority.COST:
            cost_optimized = min(scored_instances, key=lambda x: x["specs"]["cost_hour"])
            alternatives.append({
                "name": "Cost Optimized",
                "instance_type": cost_optimized["instance_type"],
                "rationale": "Lowest hourly cost while meeting requirements",
                "trade_offs": "May take longer to complete"
            })

        # Performance-optimized alternative
        if workload.priority != Priority.PERFORMANCE:
            perf_optimized = max(scored_instances, key=lambda x: x["perf_score"])
            alternatives.append({
                "name": "Performance Optimized",
                "instance_type": perf_optimized["instance_type"],
                "rationale": "Fastest completion time",
                "trade_offs": "Higher cost per hour"
            })

        # GPU alternative if not already selected
        gpu_instances = [x for x in scored_instances if "gpu" in x["specs"]]
        if gpu_instances and workload.gpu_requirement != "required":
            gpu_option = gpu_instances[0]
            alternatives.append({
                "name": "GPU Accelerated",
                "instance_type": gpu_option["instance_type"],
                "rationale": "GPU acceleration for compatible workloads",
                "trade_offs": "Higher cost, may not benefit all algorithms"
            })

        return alternatives[:3]  # Limit to 3 alternatives

    def _generate_deployment_template(self, config: Dict, workload: WorkloadCharacteristics) -> str:
        """Generate deployment template/script"""
        return f"""#!/bin/bash
# Research Infrastructure Deployment
# Domain: {workload.domain}
# Tools: {', '.join(workload.primary_tools)}
# Optimization: {workload.priority.value}

# Deploy {config['instance_count']}x {config['instance_type']} instances
./deploy-spack-research.sh \\
  --domain {workload.domain} \\
  --instance-type {config['instance_type']} \\
  --instance-count {config['instance_count']} \\
  --storage-primary {workload.data_size_gb * 2}GB \\
  --storage-scratch {workload.data_size_gb}GB \\
  --tools {','.join(workload.primary_tools)} \\
  --optimization {workload.priority.value}

# Monitor deployment
aws cloudwatch get-metric-statistics \\
  --namespace AWS/EC2 \\
  --metric-name CPUUtilization \\
  --dimensions Name=InstanceId,Value=i-1234567890abcdef0
"""

    def _generate_monitoring_config(self, workload: WorkloadCharacteristics) -> Dict[str, Any]:
        """Generate monitoring configuration"""
        return {
            "cloudwatch_metrics": [
                "CPUUtilization",
                "MemoryUtilization",
                "NetworkIn",
                "NetworkOut",
                "DiskReadOps",
                "DiskWriteOps"
            ],
            "custom_metrics": {
                "domain_specific": f"{workload.domain}_progress",
                "cost_tracking": "hourly_spend",
                "efficiency": "compute_utilization"
            },
            "alerts": [
                {"metric": "CPUUtilization", "threshold": 90, "action": "scale_up"},
                {"metric": "MemoryUtilization", "threshold": 85, "action": "alert"},
                {"metric": "hourly_spend", "threshold": workload.budget_limit or 100, "action": "budget_alert"}
            ]
        }

    def _generate_rationale(self, config: Dict, workload: WorkloadCharacteristics) -> List[str]:
        """Generate human-readable rationale for recommendations"""
        rationale = []

        # Instance selection rationale
        specs = config["specs"]
        rationale.append(f"Selected {config['instance_type']} for optimal {workload.priority.value} priority")

        # Architecture rationale
        if specs["architecture"] == "arm64":
            rationale.append("Graviton3 processor provides 20-40% better price/performance")

        # Memory rationale
        memory_per_core = specs["memory"] / specs["vcpus"]
        if memory_per_core >= 8:
            rationale.append(f"High memory-to-CPU ratio ({memory_per_core:.1f}GB/core) suits {workload.domain} workloads")

        # GPU rationale
        if "gpu" in specs:
            rationale.append(f"GPU acceleration ({specs['gpu']}) for {workload.domain} compute acceleration")

        # Scaling rationale
        if config["instance_count"] > 1:
            rationale.append(f"Multi-instance deployment ({config['instance_count']} nodes) for parallel scaling")

        # Cost rationale
        if workload.priority == Priority.COST:
            rationale.append("Cost-optimized configuration minimizes expenses while meeting requirements")

        return rationale

    def generate_recommendation_report(self, workload: WorkloadCharacteristics,
                                     recommendation: InfrastructureRecommendation) -> str:
        """Generate comprehensive recommendation report"""

        report = []
        report.append("# 🧙‍♂️ Research Infrastructure Recommendation Report")
        report.append("=" * 50)
        report.append("")

        # Workload summary
        report.append("## 📋 Workload Summary")
        report.append(f"**Research Domain**: {workload.domain.replace('_', ' ').title()}")
        report.append(f"**Primary Tools**: {', '.join(workload.primary_tools).upper()}")
        report.append(f"**Problem Size**: {workload.problem_size.value.title()}")
        report.append(f"**Optimization Priority**: {workload.priority.value.title()}")
        report.append(f"**Data Size**: {workload.data_size_gb:,} GB")
        report.append(f"**Concurrent Users**: {workload.collaboration_users}")
        report.append("")

        # Primary recommendation
        report.append("## 🎯 Primary Recommendation")
        report.append(f"**Instance Type**: {recommendation.instance_type}")
        report.append(f"**Instance Count**: {recommendation.instance_count}")
        report.append("")

        # Cost analysis
        report.append("## 💰 Cost Analysis")
        costs = recommendation.estimated_cost
        report.append(f"**Hourly Compute Cost**: ${costs['hourly_compute']:.2f}")
        report.append(f"**Estimated Runtime**: {costs['estimated_runtime_hours']:.1f} hours")
        report.append(f"**Total Compute Cost**: ${costs['total_compute']:.2f}")
        report.append(f"**Monthly Storage Cost**: ${costs['storage_monthly']:.2f}")
        report.append(f"**Total Estimated Cost**: ${costs['total_estimated']:.2f}")
        report.append("")

        # Performance expectations
        report.append("## ⚡ Performance Expectations")
        runtime = recommendation.estimated_runtime
        report.append(f"**Estimated Completion Time**: {runtime['estimated_hours']:.1f} hours")
        report.append(f"**Confidence Level**: {runtime['confidence'].title()}")
        report.append("")

        # Optimization rationale
        report.append("## 🧠 Why This Configuration?")
        for reason in recommendation.optimization_rationale:
            report.append(f"- {reason}")
        report.append("")

        # Storage configuration
        report.append("## 💾 Storage Configuration")
        storage = recommendation.storage_config
        report.append(f"**Primary Storage**: {storage['primary']['size_gb']} GB {storage['primary']['type'].upper()}")
        report.append(f"**Scratch Storage**: {storage['scratch']['size_gb']} GB {storage['scratch']['type'].upper()}")
        report.append(f"**Archive Strategy**: {storage['archive']['type'].upper()}")
        report.append("")

        # Alternative configurations
        if recommendation.alternative_configs:
            report.append("## 🔄 Alternative Configurations")
            for alt in recommendation.alternative_configs:
                report.append(f"### {alt['name']}")
                report.append(f"**Instance**: {alt['instance_type']}")
                report.append(f"**Rationale**: {alt['rationale']}")
                report.append(f"**Trade-offs**: {alt['trade_offs']}")
                report.append("")

        # Deployment instructions
        report.append("## 🚀 Quick Deploy")
        report.append("```bash")
        report.append(recommendation.deployment_template.strip())
        report.append("```")
        report.append("")

        # Monitoring setup
        report.append("## 📊 Monitoring & Optimization")
        monitoring = recommendation.monitoring_setup
        report.append("**Key Metrics to Monitor**:")
        for metric in monitoring["cloudwatch_metrics"]:
            report.append(f"- {metric}")
        report.append("")

        if "alerts" in monitoring:
            report.append("**Automated Alerts**:")
            for alert in monitoring["alerts"]:
                report.append(f"- {alert['metric']} > {alert['threshold']}% → {alert['action']}")
        report.append("")

        # Best practices
        report.append("## 💡 Optimization Tips")
        report.append("- **Spot Instances**: Consider spot instances for 60-90% cost savings on fault-tolerant workloads")
        report.append("- **Auto Scaling**: Enable auto-scaling for variable workloads")
        report.append("- **Data Locality**: Keep data and compute in the same AWS region")
        report.append("- **Resource Scheduling**: Schedule intensive workloads during off-peak hours")
        report.append("")

        return "\n".join(report)

def main():
    parser = argparse.ArgumentParser(description='Research Infrastructure Wizard')
    parser.add_argument('--interactive', action='store_true', help='Run interactive wizard')
    parser.add_argument('--domain', help='Research domain')
    parser.add_argument('--tools', help='Comma-separated list of tools')
    parser.add_argument('--size', choices=['small', 'medium', 'large', 'massive'], help='Problem size')
    parser.add_argument('--priority', choices=['cost', 'performance', 'deadline', 'balanced'], help='Optimization priority')
    parser.add_argument('--data-gb', type=int, help='Data size in GB')
    parser.add_argument('--users', type=int, default=1, help='Number of concurrent users')
    parser.add_argument('--output', default='infrastructure_recommendation.md', help='Output report file')

    args = parser.parse_args()

    wizard = ResearchInfrastructureWizard()

    if args.interactive or not all([args.domain, args.tools, args.size, args.priority, args.data_gb]):
        # Run interactive wizard
        workload = wizard.create_interactive_wizard()
    else:
        # Use command line arguments
        workload = WorkloadCharacteristics(
            domain=args.domain,
            primary_tools=args.tools.split(','),
            problem_size=WorkloadSize(args.size),
            priority=Priority(args.priority),
            deadline_hours=None,
            budget_limit=None,
            data_size_gb=args.data_gb,
            parallel_scaling="sublinear",  # Default assumption
            gpu_requirement="optional",   # Default assumption
            memory_intensity="medium",    # Default assumption
            io_pattern="streaming",       # Default assumption
            collaboration_users=args.users
        )

    # Generate recommendation
    print("\n🔄 Analyzing requirements and generating recommendations...")
    recommendation = wizard.recommend_infrastructure(workload)

    # Generate report
    report = wizard.generate_recommendation_report(workload, recommendation)

    # Save report
    with open(args.output, 'w', encoding='utf-8') as f:
        f.write(report)

    # Display summary
    print(f"\n✅ Infrastructure recommendation complete!")
    print(f"📄 Full report saved to: {args.output}")
    print(f"\n🎯 Quick Summary:")
    print(f"   Instance: {recommendation.instance_type} x{recommendation.instance_count}")
    print(f"   Estimated Cost: ${recommendation.estimated_cost['total_estimated']:.2f}")
    print(f"   Estimated Runtime: {recommendation.estimated_runtime['estimated_hours']:.1f} hours")
    print(f"   Hourly Rate: ${recommendation.estimated_cost['hourly_compute']:.2f}/hour")

if __name__ == "__main__":
    main()
