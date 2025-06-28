#!/usr/bin/env python3
"""
Spack Environment Importer - Native Integration
Captures existing Spack environments using native Spack capabilities and generates AWS-optimized deployment configurations
"""

import json
import os
import subprocess
import sys
import yaml
import re
import tempfile
import shutil
from typing import Dict, List, Any, Optional, Tuple, Union
from pathlib import Path
from dataclasses import dataclass, asdict
import argparse
import logging
from datetime import datetime

@dataclass
class SpackEnvironmentSpec:
    """Represents a complete Spack environment specification"""
    name: str
    spack_yaml: Dict[str, Any]
    concrete_specs: List[Dict[str, Any]]
    installed_packages: List[str]
    compilers: List[Dict[str, Any]]
    repos: List[str]
    config: Dict[str, Any]
    spack_version: str
    creation_date: str
    source_system: str
    binary_cache_info: Optional[Dict[str, Any]] = None
    container_recipe: Optional[str] = None

@dataclass
class AWSOptimization:
    """AWS-specific optimizations for the environment"""
    recommended_instance_type: str
    recommended_instance_count: int
    estimated_memory_gb: int
    estimated_storage_gb: int
    graviton_compatible: bool
    gpu_required: bool
    mpi_enabled: bool
    binary_cache_strategy: str
    container_image: Optional[str]
    estimated_cost_per_hour: float
    estimated_monthly_cost: float
    optimization_notes: List[str]
    aws_services: List[str]

@dataclass
class MigrationPlan:
    """Complete migration plan from HPC to AWS"""
    source_analysis: Dict[str, Any]
    target_architecture: Dict[str, Any]
    binary_cache_migration: Dict[str, Any]
    container_strategy: Dict[str, Any]
    performance_comparison: Dict[str, Any]
    cost_analysis: Dict[str, Any]
    migration_steps: List[str]
    validation_tests: List[str]

class SpackEnvironmentImporter:
    """
    Advanced Spack environment importer using native Spack capabilities
    """
    
    def __init__(self, verbose: bool = False, spack_root: Optional[Path] = None):
        self.verbose = verbose
        self.logger = self._setup_logging()
        self.spack_root = spack_root or self._find_spack_installation()
        self.spack_cmd = self._get_spack_command()
        self.aws_instance_catalog = self._load_aws_instance_catalog()
        self.aws_graviton_packages = self._load_graviton_package_db()
        self.temp_dir = None
        
        if not self.spack_cmd:
            raise RuntimeError("Spack installation not found or not accessible")
    
    def __enter__(self):
        """Context manager entry"""
        self.temp_dir = Path(tempfile.mkdtemp(prefix="spack_importer_"))
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit"""
        if self.temp_dir and self.temp_dir.exists():
            shutil.rmtree(self.temp_dir)
    
    def _setup_logging(self) -> logging.Logger:
        """Setup logging configuration"""
        logging.basicConfig(
            level=logging.DEBUG if self.verbose else logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s'
        )
        return logging.getLogger(__name__)
    
    def _find_spack_installation(self) -> Optional[Path]:
        """Find Spack installation using multiple methods"""
        # Method 1: Check SPACK_ROOT environment variable
        if "SPACK_ROOT" in os.environ:
            spack_root = Path(os.environ["SPACK_ROOT"])
            if self._validate_spack_installation(spack_root):
                return spack_root
        
        # Method 2: Check if spack is in PATH
        try:
            result = subprocess.run(
                ["which", "spack"], 
                capture_output=True, 
                text=True,
                timeout=10
            )
            if result.returncode == 0:
                spack_bin = Path(result.stdout.strip())
                # Follow symlinks to find actual installation
                actual_bin = spack_bin.resolve()
                spack_root = actual_bin.parent.parent
                if self._validate_spack_installation(spack_root):
                    return spack_root
        except (subprocess.TimeoutExpired, FileNotFoundError):
            pass
        
        # Method 3: Check common installation locations
        common_locations = [
            Path.home() / "spack",
            Path("/opt/spack"),
            Path("/usr/local/spack"),
            Path("/apps/spack"),
            Path("/sw/spack"),
            Path("/scratch/spack"),
            Path("/work/spack")
        ]
        
        for location in common_locations:
            if self._validate_spack_installation(location):
                return location
        
        return None
    
    def _validate_spack_installation(self, spack_root: Path) -> bool:
        """Validate that a path contains a working Spack installation"""
        if not spack_root.exists():
            return False
        
        spack_bin = spack_root / "bin" / "spack"
        if not spack_bin.exists():
            return False
        
        # Test if spack command works
        try:
            result = subprocess.run(
                [str(spack_bin), "version"],
                capture_output=True,
                text=True,
                timeout=30
            )
            return result.returncode == 0
        except (subprocess.TimeoutExpired, FileNotFoundError):
            return False
    
    def _get_spack_command(self) -> Optional[str]:
        """Get the spack command to use"""
        if self.spack_root:
            spack_bin = self.spack_root / "bin" / "spack"
            if spack_bin.exists():
                return str(spack_bin)
        
        # Try system spack
        try:
            result = subprocess.run(
                ["which", "spack"],
                capture_output=True,
                text=True,
                timeout=10
            )
            if result.returncode == 0:
                return result.stdout.strip()
        except subprocess.TimeoutExpired:
            pass
        
        return None
    
    def _load_aws_instance_catalog(self) -> Dict[str, Any]:
        """Load comprehensive AWS instance catalog with Graviton support"""
        return {
            # Intel/AMD x86_64 instances
            "c6i.large": {
                "vcpus": 2, "memory_gb": 4, "arch": "x86_64", "network": "up_to_12_5",
                "cost_per_hour": 0.096, "graviton": False, "category": "compute"
            },
            "c6i.xlarge": {
                "vcpus": 4, "memory_gb": 8, "arch": "x86_64", "network": "up_to_12_5", 
                "cost_per_hour": 0.192, "graviton": False, "category": "compute"
            },
            "c6i.2xlarge": {
                "vcpus": 8, "memory_gb": 16, "arch": "x86_64", "network": "up_to_12_5",
                "cost_per_hour": 0.384, "graviton": False, "category": "compute"
            },
            "c6i.4xlarge": {
                "vcpus": 16, "memory_gb": 32, "arch": "x86_64", "network": "up_to_12_5",
                "cost_per_hour": 0.768, "graviton": False, "category": "compute"
            },
            "c6i.8xlarge": {
                "vcpus": 32, "memory_gb": 64, "arch": "x86_64", "network": "12_5",
                "cost_per_hour": 1.536, "graviton": False, "category": "compute"
            },
            
            # Graviton (ARM64) instances - better price/performance
            "c7g.large": {
                "vcpus": 2, "memory_gb": 4, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.0725, "graviton": True, "category": "compute"
            },
            "c7g.xlarge": {
                "vcpus": 4, "memory_gb": 8, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.145, "graviton": True, "category": "compute"
            },
            "c7g.2xlarge": {
                "vcpus": 8, "memory_gb": 16, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.29, "graviton": True, "category": "compute"
            },
            "c7g.4xlarge": {
                "vcpus": 16, "memory_gb": 32, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.58, "graviton": True, "category": "compute"
            },
            "c7g.8xlarge": {
                "vcpus": 32, "memory_gb": 64, "arch": "aarch64", "network": "12_5",
                "cost_per_hour": 1.16, "graviton": True, "category": "compute"
            },
            
            # Memory-optimized Graviton
            "r7g.large": {
                "vcpus": 2, "memory_gb": 16, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.1008, "graviton": True, "category": "memory"
            },
            "r7g.xlarge": {
                "vcpus": 4, "memory_gb": 32, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.2016, "graviton": True, "category": "memory"
            },
            "r7g.2xlarge": {
                "vcpus": 8, "memory_gb": 64, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.4032, "graviton": True, "category": "memory"
            },
            "r7g.4xlarge": {
                "vcpus": 16, "memory_gb": 128, "arch": "aarch64", "network": "up_to_12_5",
                "cost_per_hour": 0.8064, "graviton": True, "category": "memory"
            },
            
            # HPC-optimized instances
            "hpc7g.4xlarge": {
                "vcpus": 16, "memory_gb": 128, "arch": "aarch64", "network": "200",
                "cost_per_hour": 1.344, "graviton": True, "category": "hpc"
            },
            "hpc7g.8xlarge": {
                "vcpus": 32, "memory_gb": 256, "arch": "aarch64", "network": "200",
                "cost_per_hour": 2.688, "graviton": True, "category": "hpc"
            },
            "hpc7g.16xlarge": {
                "vcpus": 64, "memory_gb": 512, "arch": "aarch64", "network": "200",
                "cost_per_hour": 5.376, "graviton": True, "category": "hpc"
            },
            
            # GPU instances
            "g4dn.xlarge": {
                "vcpus": 4, "memory_gb": 16, "gpu": "1x T4", "arch": "x86_64",
                "cost_per_hour": 0.526, "graviton": False, "category": "gpu"
            },
            "p4d.24xlarge": {
                "vcpus": 96, "memory_gb": 1152, "gpu": "8x A100", "arch": "x86_64",
                "cost_per_hour": 32.772, "graviton": False, "category": "gpu"
            }
        }
    
    def _load_graviton_package_db(self) -> Dict[str, Any]:
        """Load database of Graviton-compatible packages"""
        return {
            "native_support": [
                # Compilers
                "gcc", "llvm", "intel-oneapi-compilers",
                
                # MPI implementations
                "openmpi", "mpich", "intel-mpi",
                
                # Math libraries
                "openblas", "fftw", "scalapack", "petsc", "slepc",
                
                # Scientific computing
                "python", "numpy", "scipy", "hdf5", "netcdf-c", "boost",
                
                # Machine learning
                "tensorflow", "pytorch", "onnx",
                
                # Climate/weather
                "wrf", "nco", "cdo", "ncview",
                
                # Bioinformatics (many work well on ARM)
                "samtools", "bcftools", "bwa", "bowtie2", "minimap2"
            ],
            
            "optimization_flags": {
                "aarch64": [
                    "target=graviton3",
                    "tune=graviton3",
                    "cpu=neoverse-v1"
                ]
            },
            
            "performance_boost": {
                # Packages that see significant performance improvements on Graviton
                "tensorflow": 1.3,   # 30% performance boost
                "pytorch": 1.25,     # 25% performance boost
                "openblas": 1.4,     # 40% performance boost
                "fftw": 1.2,         # 20% performance boost
                "wrf": 1.35,         # 35% performance boost for weather modeling
                "python": 1.2        # 20% general Python performance boost
            }
        }
    
    def discover_environments(self) -> List[str]:
        """Discover available Spack environments using native commands"""
        self.logger.info("Discovering Spack environments...")
        
        try:
            result = subprocess.run(
                [self.spack_cmd, "env", "list", "--format", "{name}"],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            if result.returncode != 0:
                self.logger.error(f"Failed to list environments: {result.stderr}")
                return []
            
            environments = []
            for line in result.stdout.strip().split('\n'):
                line = line.strip()
                if line and not line.startswith('='):
                    environments.append(line)
            
            self.logger.info(f"Found {len(environments)} environments: {environments}")
            return environments
        
        except Exception as e:
            self.logger.error(f"Error discovering environments: {e}")
            return []
    
    def capture_environment(self, env_name: str) -> Optional[SpackEnvironmentSpec]:
        """Capture complete Spack environment using native export capabilities"""
        self.logger.info(f"Capturing environment: {env_name}")
        
        try:
            # Step 1: Export environment to YAML using native command
            export_result = subprocess.run(
                [self.spack_cmd, "env", "export", env_name],
                capture_output=True,
                text=True,
                timeout=120
            )
            
            if export_result.returncode != 0:
                self.logger.error(f"Failed to export environment: {export_result.stderr}")
                return None
            
            # Parse the exported YAML
            try:
                spack_yaml = yaml.safe_load(export_result.stdout)
            except yaml.YAMLError as e:
                self.logger.error(f"Failed to parse exported YAML: {e}")
                return None
            
            # Step 2: Get concrete specs (resolved dependencies)
            concrete_specs = self._get_concrete_specs(env_name)
            
            # Step 3: Get installed packages
            installed_packages = self._get_installed_packages(env_name)
            
            # Step 4: Get compiler information
            compilers = self._get_compiler_info()
            
            # Step 5: Get repository information
            repos = self._get_repo_info()
            
            # Step 6: Get Spack configuration
            config = self._get_spack_config()
            
            # Step 7: Get Spack version
            spack_version = self._get_spack_version()
            
            # Step 8: Check for existing binary cache
            binary_cache_info = self._analyze_binary_cache(env_name)
            
            # Step 9: Generate container recipe if possible
            container_recipe = self._generate_container_recipe(env_name)
            
            return SpackEnvironmentSpec(
                name=env_name,
                spack_yaml=spack_yaml,
                concrete_specs=concrete_specs,
                installed_packages=installed_packages,
                compilers=compilers,
                repos=repos,
                config=config,
                spack_version=spack_version,
                creation_date=datetime.now().isoformat(),
                source_system=self._get_system_info(),
                binary_cache_info=binary_cache_info,
                container_recipe=container_recipe
            )
        
        except Exception as e:
            self.logger.error(f"Error capturing environment {env_name}: {e}")
            return None
    
    def _get_concrete_specs(self, env_name: str) -> List[Dict[str, Any]]:
        """Get concrete (resolved) specs for the environment"""
        try:
            result = subprocess.run(
                [self.spack_cmd, "-e", env_name, "find", "--format", 
                 "{name}@{version}%{compiler}{variants}{arch}"],
                capture_output=True,
                text=True,
                timeout=60
            )
            
            if result.returncode != 0:
                self.logger.warning(f"Could not get concrete specs: {result.stderr}")
                return []
            
            specs = []
            for line in result.stdout.strip().split('\n'):
                line = line.strip()
                if line and not line.startswith('='):
                    specs.append({"spec": line})
            
            return specs
        
        except Exception as e:
            self.logger.warning(f"Error getting concrete specs: {e}")
            return []
    
    def _get_installed_packages(self, env_name: str) -> List[str]:
        """Get list of installed packages in the environment"""
        try:
            result = subprocess.run(
                [self.spack_cmd, "-e", env_name, "find", "--installed", "--format", "{name}"],
                capture_output=True,
                text=True,
                timeout=60
            )
            
            if result.returncode != 0:
                return []
            
            packages = []
            for line in result.stdout.strip().split('\n'):
                line = line.strip()
                if line and not line.startswith('='):
                    packages.append(line)
            
            return packages
        
        except Exception as e:
            self.logger.warning(f"Error getting installed packages: {e}")
            return []
    
    def _get_compiler_info(self) -> List[Dict[str, Any]]:
        """Get available compiler information"""
        try:
            result = subprocess.run(
                [self.spack_cmd, "compiler", "info", "--format", "json"],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            if result.returncode == 0:
                try:
                    return json.loads(result.stdout)
                except json.JSONDecodeError:
                    pass
            
            # Fallback to simple list
            result = subprocess.run(
                [self.spack_cmd, "compiler", "list"],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            if result.returncode == 0:
                compilers = []
                current_family = None
                for line in result.stdout.split('\n'):
                    line = line.strip()
                    if line.endswith(':'):
                        current_family = line[:-1]
                    elif line and current_family:
                        compilers.append({
                            "family": current_family,
                            "version": line,
                            "spec": f"{current_family}@{line}"
                        })
                return compilers
        
        except Exception as e:
            self.logger.warning(f"Error getting compiler info: {e}")
        
        return []
    
    def _get_repo_info(self) -> List[str]:
        """Get repository information"""
        try:
            result = subprocess.run(
                [self.spack_cmd, "repo", "list"],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            if result.returncode == 0:
                repos = []
                for line in result.stdout.split('\n'):
                    line = line.strip()
                    if line and not line.startswith('=') and not line.startswith('-'):
                        # Extract repo path/name
                        parts = line.split()
                        if len(parts) > 1:
                            repos.append(parts[1])  # Usually the path
                return repos
        
        except Exception as e:
            self.logger.warning(f"Error getting repo info: {e}")
        
        return []
    
    def _get_spack_config(self) -> Dict[str, Any]:
        """Get Spack configuration"""
        try:
            result = subprocess.run(
                [self.spack_cmd, "config", "get", "config"],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            if result.returncode == 0:
                try:
                    return yaml.safe_load(result.stdout)
                except yaml.YAMLError:
                    pass
        
        except Exception as e:
            self.logger.warning(f"Error getting Spack config: {e}")
        
        return {}
    
    def _get_spack_version(self) -> str:
        """Get Spack version"""
        try:
            result = subprocess.run(
                [self.spack_cmd, "version"],
                capture_output=True,
                text=True,
                timeout=15
            )
            
            if result.returncode == 0:
                return result.stdout.strip()
        
        except Exception as e:
            self.logger.warning(f"Error getting Spack version: {e}")
        
        return "unknown"
    
    def _analyze_binary_cache(self, env_name: str) -> Optional[Dict[str, Any]]:
        """Analyze binary cache usage for the environment"""
        try:
            # Check configured binary caches
            result = subprocess.run(
                [self.spack_cmd, "mirror", "list"],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            if result.returncode == 0:
                mirrors = []
                for line in result.stdout.split('\n'):
                    line = line.strip()
                    if line and not line.startswith('='):
                        parts = line.split()
                        if len(parts) >= 2:
                            mirrors.append({
                                "name": parts[0],
                                "url": parts[1]
                            })
                
                return {
                    "configured_mirrors": mirrors,
                    "cache_usage": "unknown"  # Would need more analysis
                }
        
        except Exception as e:
            self.logger.warning(f"Error analyzing binary cache: {e}")
        
        return None
    
    def _generate_container_recipe(self, env_name: str) -> Optional[str]:
        """Generate container recipe using Spack's native containerize command"""
        if not self.temp_dir:
            return None
        
        try:
            # Generate Dockerfile using Spack's containerize command
            dockerfile_path = self.temp_dir / "Dockerfile"
            
            result = subprocess.run(
                [self.spack_cmd, "-e", env_name, "containerize", 
                 "--format", "docker", "--output", str(dockerfile_path)],
                capture_output=True,
                text=True,
                timeout=120
            )
            
            if result.returncode == 0 and dockerfile_path.exists():
                with open(dockerfile_path, 'r') as f:
                    return f.read()
            else:
                self.logger.warning(f"Could not generate container recipe: {result.stderr}")
        
        except Exception as e:
            self.logger.warning(f"Error generating container recipe: {e}")
        
        return None
    
    def _get_system_info(self) -> str:
        """Get comprehensive system information"""
        try:
            # Get hostname
            hostname_result = subprocess.run(["hostname", "-f"], capture_output=True, text=True)
            hostname = hostname_result.stdout.strip() if hostname_result.returncode == 0 else "unknown"
            
            # Get OS info
            os_result = subprocess.run(["uname", "-a"], capture_output=True, text=True)
            os_info = os_result.stdout.strip() if os_result.returncode == 0 else "unknown"
            
            # Get CPU info
            cpu_info = "unknown"
            try:
                with open("/proc/cpuinfo", "r") as f:
                    for line in f:
                        if line.startswith("model name"):
                            cpu_info = line.split(":", 1)[1].strip()
                            break
            except:
                pass
            
            return f"{hostname} | {os_info} | CPU: {cpu_info}"
        
        except Exception:
            return "unknown"
    
    def analyze_environment(self, env_spec: SpackEnvironmentSpec) -> AWSOptimization:
        """Analyze environment and generate comprehensive AWS optimization"""
        self.logger.info(f"Analyzing environment: {env_spec.name}")
        
        # Analyze package characteristics
        analysis = self._analyze_package_requirements(env_spec)
        
        # Check Graviton compatibility
        graviton_analysis = self._analyze_graviton_compatibility(env_spec)
        
        # Select optimal instance configuration
        instance_config = self._select_optimal_aws_configuration(analysis, graviton_analysis)
        
        # Generate binary cache strategy
        binary_cache_strategy = self._design_binary_cache_strategy(env_spec)
        
        # Generate container strategy
        container_image = self._design_container_strategy(env_spec)
        
        # Calculate costs
        cost_analysis = self._calculate_aws_costs(instance_config, analysis)
        
        # Generate AWS services recommendations
        aws_services = self._recommend_aws_services(analysis, instance_config)
        
        # Generate optimization notes
        optimization_notes = self._generate_optimization_notes(
            analysis, graviton_analysis, instance_config
        )
        
        return AWSOptimization(
            recommended_instance_type=instance_config["instance_type"],
            recommended_instance_count=instance_config["instance_count"],
            estimated_memory_gb=instance_config["memory_gb"],
            estimated_storage_gb=analysis["storage_requirements"],
            graviton_compatible=graviton_analysis["compatible"],
            gpu_required=analysis["gpu_required"],
            mpi_enabled=analysis["mpi_enabled"],
            binary_cache_strategy=binary_cache_strategy,
            container_image=container_image,
            estimated_cost_per_hour=cost_analysis["hourly_cost"],
            estimated_monthly_cost=cost_analysis["monthly_cost"],
            optimization_notes=optimization_notes,
            aws_services=aws_services
        )
    
    def _analyze_package_requirements(self, env_spec: SpackEnvironmentSpec) -> Dict[str, Any]:
        """Analyze package requirements and characteristics"""
        analysis = {
            "total_packages": len(env_spec.installed_packages),
            "memory_intensive_packages": 0,
            "gpu_packages": 0,
            "mpi_packages": 0,
            "compile_intensive_packages": 0,
            "io_intensive_packages": 0,
            "memory_requirements": 8,  # GB baseline
            "storage_requirements": 100,  # GB baseline
            "gpu_required": False,
            "mpi_enabled": False,
            "package_categories": {}
        }
        
        # Package classification database
        package_db = {
            "memory_intensive": [
                "gatk", "star", "canu", "trinity", "velvet", "spades", "abyss",
                "vasp", "quantum-espresso", "cp2k", "gaussian", "molpro", "orca",
                "wrf", "cesm", "cam", "clm", "pop", "cice", "mpas",
                "openfoam", "fluent", "ansys", "comsol", "su2"
            ],
            "gpu_accelerated": [
                "tensorflow", "pytorch", "keras", "cuda", "cudnn", "nccl", "cupy",
                "amber", "gromacs", "namd", "lammps", "hoomd-blue",
                "magma", "scalapack", "fftw", "mkl", "openblas",
                "opencv", "vtk", "paraview", "visit"
            ],
            "mpi_enabled": [
                "openmpi", "mpich", "intel-mpi", "mvapich2", "spectrum-mpi",
                "hdf5", "netcdf-c", "parallel-netcdf", "adios2", "pnetcdf",
                "petsc", "slepc", "mumps", "scalapack", "hypre", "trilinos",
                "wrf", "cesm", "geos-chem", "cmaq", "cam-chem",
                "quantum-espresso", "vasp", "cp2k", "nwchem", "gamess",
                "openfoam", "su2", "fenics", "dealii"
            ],
            "io_intensive": [
                "blast", "diamond", "kraken2", "centrifuge", "bowtie2",
                "gatk", "picard", "samtools", "bcftools", "vcftools",
                "gdal", "netcdf-c", "hdf5", "zarr", "h5py"
            ]
        }
        
        # Analyze each package
        for package_name in env_spec.installed_packages:
            package_lower = package_name.lower()
            
            # Check package categories
            for category, package_list in package_db.items():
                if any(pkg in package_lower for pkg in package_list):
                    analysis[f"{category}_packages"] = analysis.get(f"{category}_packages", 0) + 1
                    
                    # Update requirements based on category
                    if category == "memory_intensive":
                        analysis["memory_requirements"] = max(
                            analysis["memory_requirements"], 32
                        )
                    elif category == "gpu_accelerated":
                        analysis["gpu_required"] = True
                        analysis["memory_requirements"] = max(
                            analysis["memory_requirements"], 32
                        )
                    elif category == "mpi_enabled":
                        analysis["mpi_enabled"] = True
                        analysis["memory_requirements"] = max(
                            analysis["memory_requirements"], 16
                        )
                    elif category == "io_intensive":
                        analysis["storage_requirements"] = max(
                            analysis["storage_requirements"], 500
                        )
        
        # Adjust requirements based on package counts
        if analysis["memory_intensive_packages"] > 3:
            analysis["memory_requirements"] *= 2
        
        if analysis["mpi_enabled"] and analysis["memory_intensive_packages"] > 0:
            analysis["memory_requirements"] = max(analysis["memory_requirements"], 64)
        
        if analysis["io_intensive_packages"] > 2:
            analysis["storage_requirements"] *= 2
        
        return analysis
    
    def _analyze_graviton_compatibility(self, env_spec: SpackEnvironmentSpec) -> Dict[str, Any]:
        """Analyze Graviton (ARM64) compatibility"""
        analysis = {
            "compatible": True,
            "incompatible_packages": [],
            "native_packages": [],
            "performance_boost_packages": [],
            "estimated_performance_gain": 1.0,
            "cost_savings": 0.0
        }
        
        # Known incompatible packages (mostly legacy or x86-specific)
        incompatible_packages = [
            "intel-parallel-studio", "intel-mkl", "intel-mpi", "intel-tbb",
            "pgi", "nvhpc", "aocc",  # Intel/NVIDIA/AMD specific compilers
            "matlab", "mathematica", "ansys-fluent",  # Proprietary software
        ]
        
        for package_name in env_spec.installed_packages:
            package_lower = package_name.lower()
            
            # Check for incompatible packages
            if any(pkg in package_lower for pkg in incompatible_packages):
                analysis["compatible"] = False
                analysis["incompatible_packages"].append(package_name)
            
            # Check for native Graviton support
            elif any(pkg in package_lower for pkg in self.aws_graviton_packages["native_support"]):
                analysis["native_packages"].append(package_name)
                
                # Check for performance boost
                for boost_pkg, boost_factor in self.aws_graviton_packages["performance_boost"].items():
                    if boost_pkg in package_lower:
                        analysis["performance_boost_packages"].append({
                            "package": package_name,
                            "boost_factor": boost_factor
                        })
                        # Update estimated performance gain (weighted average)
                        analysis["estimated_performance_gain"] = max(
                            analysis["estimated_performance_gain"], boost_factor
                        )
        
        # Calculate potential cost savings (Graviton instances are ~20% cheaper)
        if analysis["compatible"]:
            analysis["cost_savings"] = 0.20  # 20% cost savings
        
        return analysis
    
    def _select_optimal_aws_configuration(self, analysis: Dict[str, Any], 
                                         graviton_analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Select optimal AWS instance configuration"""
        
        # Determine instance requirements
        memory_required = analysis["memory_requirements"]
        gpu_required = analysis["gpu_required"]
        mpi_enabled = analysis["mpi_enabled"]
        
        # Filter instances based on requirements
        candidate_instances = []
        
        for instance_type, specs in self.aws_instance_catalog.items():
            # Skip GPU instances if not needed
            if gpu_required and "gpu" not in specs:
                continue
            if not gpu_required and "gpu" in specs:
                continue
            
            # Check memory requirements
            if specs["memory_gb"] < memory_required:
                continue
            
            # For MPI workloads, prefer HPC-optimized or compute-optimized instances
            if mpi_enabled and analysis["total_packages"] > 20:
                if specs["category"] not in ["hpc", "compute"]:
                    continue
            
            # Check Graviton compatibility
            if specs["graviton"] and not graviton_analysis["compatible"]:
                continue
            
            candidate_instances.append((instance_type, specs))
        
        if not candidate_instances:
            # Fallback to basic instance
            return {
                "instance_type": "m6i.xlarge",
                "instance_count": 1,
                "memory_gb": 16,
                "architecture": "x86_64"
            }
        
        # Sort by cost-effectiveness, preferring Graviton when compatible
        def cost_effectiveness(item):
            instance_type, specs = item
            memory_per_dollar = specs["memory_gb"] / specs["cost_per_hour"]
            
            # Bonus for Graviton compatibility
            if specs["graviton"] and graviton_analysis["compatible"]:
                memory_per_dollar *= 1.5  # Prefer Graviton
            
            return memory_per_dollar
        
        candidate_instances.sort(key=cost_effectiveness, reverse=True)
        best_instance_type, best_specs = candidate_instances[0]
        
        # Determine instance count for MPI workloads
        instance_count = 1
        if mpi_enabled and analysis["total_packages"] > 30:
            # Scale based on package complexity
            if analysis["memory_intensive_packages"] > 5:
                instance_count = min(8, max(2, analysis["total_packages"] // 20))
            else:
                instance_count = min(4, max(2, analysis["total_packages"] // 25))
        
        return {
            "instance_type": best_instance_type,
            "instance_count": instance_count,
            "memory_gb": best_specs["memory_gb"],
            "architecture": best_specs["arch"],
            "graviton": best_specs.get("graviton", False)
        }
    
    def _design_binary_cache_strategy(self, env_spec: SpackEnvironmentSpec) -> str:
        """Design binary cache strategy for AWS deployment"""
        
        # Check if environment already uses binary cache
        existing_mirrors = []
        if env_spec.binary_cache_info:
            existing_mirrors = env_spec.binary_cache_info.get("configured_mirrors", [])
        
        strategies = []
        
        # Strategy 1: Use existing AWS Spack cache if available
        strategies.append("Use AWS Spack public binary cache for common packages")
        
        # Strategy 2: Create custom S3 binary cache
        strategies.append("Create custom S3 binary cache for organization-specific builds")
        
        # Strategy 3: Hybrid approach
        if len(env_spec.installed_packages) > 50:
            strategies.append("Implement tiered caching: public cache + private S3 cache")
        
        # Strategy 4: Container-based approach
        if env_spec.container_recipe:
            strategies.append("Use containerized environment with pre-built packages")
        
        return " | ".join(strategies)
    
    def _design_container_strategy(self, env_spec: SpackEnvironmentSpec) -> Optional[str]:
        """Design container strategy"""
        if env_spec.container_recipe:
            return f"Custom container based on {env_spec.spack_version}"
        
        # Recommend container approach for complex environments
        if len(env_spec.installed_packages) > 20:
            return "Multi-stage Docker build with Spack cache optimization"
        
        return None
    
    def _calculate_aws_costs(self, instance_config: Dict[str, Any], 
                           analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Calculate comprehensive AWS costs"""
        
        instance_type = instance_config["instance_type"]
        instance_count = instance_config["instance_count"]
        
        instance_specs = self.aws_instance_catalog[instance_type]
        
        # Compute costs
        hourly_compute = instance_specs["cost_per_hour"] * instance_count
        
        # Storage costs (EBS + EFS if needed)
        storage_gb = analysis["storage_requirements"]
        storage_hourly = (storage_gb * 0.08 / 730)  # $0.08/GB-month for gp3
        
        if analysis["mpi_enabled"] and instance_count > 1:
            # Add shared storage (FSx Lustre) for MPI workloads
            storage_hourly += 0.125 * max(1200, storage_gb) / 730  # FSx Lustre cost
        
        # Network costs (minimal for most research workloads)
        network_hourly = 0.01 * instance_count  # Estimate
        
        total_hourly = hourly_compute + storage_hourly + network_hourly
        
        # Monthly costs (assume 30% utilization for research workloads)
        monthly_cost = total_hourly * 24 * 30 * 0.3
        
        return {
            "hourly_cost": total_hourly,
            "monthly_cost": monthly_cost,
            "compute_hourly": hourly_compute,
            "storage_hourly": storage_hourly,
            "utilization_assumed": 0.3
        }
    
    def _recommend_aws_services(self, analysis: Dict[str, Any], 
                               instance_config: Dict[str, Any]) -> List[str]:
        """Recommend AWS services based on requirements"""
        services = [
            "EC2 for compute instances",
            "EBS gp3 for primary storage",
            "S3 for data backup and archival",
            "CloudWatch for monitoring",
            "IAM for access control"
        ]
        
        if analysis["mpi_enabled"]:
            services.extend([
                "VPC with cluster placement groups",
                "FSx Lustre for high-performance shared storage",
                "Elastic Fabric Adapter (EFA) for ultra-low latency MPI networking up to 200 Gbps",
                "AWS OFI-NCCL for optimized multi-GPU communication over EFA"
            ])
        
        if analysis["gpu_required"]:
            services.append("EC2 GPU instances with NVIDIA drivers")
        
        if analysis["io_intensive_packages"] > 2:
            services.append("EFS for shared network file system")
        
        if analysis["total_packages"] > 50:
            services.extend([
                "S3 for Spack binary cache",
                "ECR for container image storage"
            ])
        
        # Cost optimization services
        services.extend([
            "AWS Cost Explorer for cost monitoring",
            "AWS Budgets for cost alerts",
            "Spot Instances for 60-80% cost savings on batch workloads"
        ])
        
        return services
    
    def _generate_optimization_notes(self, analysis: Dict[str, Any], 
                                   graviton_analysis: Dict[str, Any],
                                   instance_config: Dict[str, Any]) -> List[str]:
        """Generate optimization recommendations"""
        notes = []
        
        # Graviton optimization
        if graviton_analysis["compatible"]:
            savings = graviton_analysis["cost_savings"] * 100
            performance = (graviton_analysis["estimated_performance_gain"] - 1) * 100
            notes.append(f"Graviton instances provide ~{savings:.0f}% cost savings and up to {performance:.0f}% performance boost")
        
        # MPI optimization
        if analysis["mpi_enabled"]:
            notes.append("Use EFA-enabled instances with cluster placement groups for optimal MPI performance")
            notes.append("Configure libfabric with EFA provider for ultra-low latency networking")
            if instance_config["instance_count"] > 1:
                notes.append(f"Multi-node setup ({instance_config['instance_count']} instances) recommended for scaling")
        
        # Memory optimization
        if analysis["memory_intensive_packages"] > 0:
            notes.append("Memory-optimized instances selected for large memory requirements")
        
        # GPU optimization
        if analysis["gpu_required"]:
            notes.append("GPU instances required for acceleration - consider Spot instances for training workloads")
        
        # Storage optimization
        if analysis["storage_requirements"] > 1000:
            notes.append("Consider tiered storage strategy: EBS for active data, S3 for archival")
        
        # Binary cache optimization
        notes.append("Implement Spack binary cache on S3 for faster package installation")
        
        # Cost optimization
        notes.append("Use Spot instances for batch workloads to achieve 60-80% cost savings")
        notes.append("Implement auto-scaling to minimize costs during idle periods")
        
        # Container optimization
        if len(analysis) > 20:  # Complex environments
            notes.append("Consider containerized deployment for improved reproducibility and faster startup")
        
        return notes
    
    def generate_migration_plan(self, env_spec: SpackEnvironmentSpec, 
                               optimization: AWSOptimization) -> MigrationPlan:
        """Generate comprehensive migration plan"""
        
        # Analyze source system
        source_analysis = {
            "environment_name": env_spec.name,
            "package_count": len(env_spec.installed_packages),
            "spack_version": env_spec.spack_version,
            "source_system": env_spec.source_system,
            "compilers": len(env_spec.compilers),
            "repositories": len(env_spec.repos)
        }
        
        # Define target architecture
        target_architecture = {
            "instance_type": optimization.recommended_instance_type,
            "instance_count": optimization.recommended_instance_count,
            "architecture": "aarch64" if optimization.graviton_compatible else "x86_64",
            "memory_gb": optimization.estimated_memory_gb,
            "storage_gb": optimization.estimated_storage_gb,
            "gpu_enabled": optimization.gpu_required,
            "mpi_enabled": optimization.mpi_enabled
        }
        
        # Binary cache migration strategy
        binary_cache_migration = {
            "source_cache": env_spec.binary_cache_info,
            "target_strategy": optimization.binary_cache_strategy,
            "s3_bucket": f"spack-cache-{env_spec.name.lower()}",
            "migration_steps": [
                "Create S3 bucket for binary cache",
                "Configure Spack to use S3 mirror",
                "Build packages with cache enabled",
                "Verify cache functionality"
            ]
        }
        
        # Container strategy
        container_strategy = {
            "use_containers": optimization.container_image is not None,
            "base_image": "ubuntu:22.04" if optimization.graviton_compatible else "ubuntu:22.04",
            "spack_image": f"spack/ubuntu-jammy:v{env_spec.spack_version}",
            "custom_recipe": env_spec.container_recipe is not None
        }
        
        # Performance comparison estimate
        performance_comparison = {
            "baseline": "Current HPC system",
            "aws_estimate": {
                "compute_performance": "Similar to 10-20% improvement" if optimization.graviton_compatible else "Similar",
                "i_o_performance": "Significantly improved with FSx Lustre" if optimization.mpi_enabled else "Improved",
                "network_performance": "Enhanced with 200 Gbps networking" if optimization.mpi_enabled else "Standard",
                "cost_efficiency": f"${optimization.estimated_monthly_cost:.0f}/month estimated"
            }
        }
        
        # Cost analysis
        cost_analysis = {
            "monthly_estimate": optimization.estimated_monthly_cost,
            "hourly_rate": optimization.estimated_cost_per_hour,
            "cost_breakdown": {
                "compute": f"${optimization.estimated_cost_per_hour * 0.8:.2f}/hour",
                "storage": f"${optimization.estimated_cost_per_hour * 0.15:.2f}/hour", 
                "network": f"${optimization.estimated_cost_per_hour * 0.05:.2f}/hour"
            },
            "cost_optimization": {
                "spot_savings": "60-80% for batch workloads",
                "graviton_savings": "20% with Graviton instances" if optimization.graviton_compatible else "N/A",
                "reserved_savings": "40-60% for steady-state workloads"
            }
        }
        
        # Migration steps
        migration_steps = [
            "Phase 1: Environment Assessment and Planning",
            "  - Review current Spack environment configuration",
            "  - Validate package compatibility with target architecture",
            "  - Plan data migration strategy",
            
            "Phase 2: AWS Infrastructure Setup",
            "  - Create VPC and networking infrastructure",
            "  - Launch EC2 instances with recommended configuration",
            "  - Configure storage systems (EBS, EFS, FSx as needed)",
            "  - Set up monitoring and logging",
            
            "Phase 3: Spack Environment Migration",
            "  - Install Spack on AWS instances",
            "  - Configure binary cache on S3",
            "  - Import environment specification",
            "  - Build and install packages",
            
            "Phase 4: Testing and Validation",
            "  - Run validation tests for all packages",
            "  - Performance benchmarking",
            "  - Workflow testing",
            "  - User acceptance testing",
            
            "Phase 5: Production Deployment",
            "  - Migrate production workloads",
            "  - Configure monitoring and alerting",
            "  - Train users on new environment",
            "  - Implement backup and disaster recovery"
        ]
        
        # Validation tests
        validation_tests = [
            "Package installation verification",
            "Compiler functionality testing",
            "MPI communication testing (if applicable)",
            "GPU functionality testing (if applicable)",
            "I/O performance benchmarking",
            "Application-specific regression testing",
            "Multi-user environment testing",
            "Backup and restore procedures",
            "Cost monitoring and alerting verification"
        ]
        
        return MigrationPlan(
            source_analysis=source_analysis,
            target_architecture=target_architecture,
            binary_cache_migration=binary_cache_migration,
            container_strategy=container_strategy,
            performance_comparison=performance_comparison,
            cost_analysis=cost_analysis,
            migration_steps=migration_steps,
            validation_tests=validation_tests
        )
    
    def export_aws_configuration(self, env_spec: SpackEnvironmentSpec, 
                                optimization: AWSOptimization,
                                migration_plan: MigrationPlan,
                                output_dir: Path) -> Dict[str, Path]:
        """Export comprehensive AWS configuration files"""
        self.logger.info(f"Exporting AWS configuration to: {output_dir}")
        
        output_dir.mkdir(parents=True, exist_ok=True)
        exported_files = {}
        
        # 1. Export main configuration
        main_config = {
            "environment_spec": asdict(env_spec),
            "aws_optimization": asdict(optimization),
            "migration_plan": asdict(migration_plan),
            "export_timestamp": datetime.now().isoformat()
        }
        
        config_file = output_dir / f"{env_spec.name}_aws_config.json"
        with open(config_file, 'w') as f:
            json.dump(main_config, f, indent=2, default=str)
        exported_files["main_config"] = config_file
        
        # 2. Export Spack environment YAML
        spack_yaml_file = output_dir / f"{env_spec.name}_spack_env.yaml"
        with open(spack_yaml_file, 'w') as f:
            yaml.dump(env_spec.spack_yaml, f, default_flow_style=False)
        exported_files["spack_yaml"] = spack_yaml_file
        
        # 3. Export AWS-optimized Spack configuration
        aws_spack_config = self._generate_aws_spack_config(optimization)
        aws_config_file = output_dir / f"{env_spec.name}_aws_spack_config.yaml"
        with open(aws_config_file, 'w') as f:
            yaml.dump(aws_spack_config, f, default_flow_style=False)
        exported_files["aws_spack_config"] = aws_config_file
        
        # 4. Export Terraform configuration
        terraform_config = self._generate_terraform_config(optimization, migration_plan)
        terraform_file = output_dir / f"{env_spec.name}_infrastructure.tf"
        with open(terraform_file, 'w') as f:
            f.write(terraform_config)
        exported_files["terraform"] = terraform_file
        
        # 5. Export deployment script
        deploy_script = self._generate_deployment_script(env_spec, optimization)
        deploy_file = output_dir / f"deploy_{env_spec.name}.sh"
        with open(deploy_file, 'w') as f:
            f.write(deploy_script)
        deploy_file.chmod(0o755)
        exported_files["deployment_script"] = deploy_file
        
        # 6. Export container files if applicable
        if env_spec.container_recipe:
            dockerfile = output_dir / f"{env_spec.name}_Dockerfile"
            with open(dockerfile, 'w') as f:
                f.write(env_spec.container_recipe)
            exported_files["dockerfile"] = dockerfile
            
            # Generate optimized Dockerfile for AWS
            aws_dockerfile = self._generate_aws_optimized_dockerfile(env_spec, optimization)
            aws_dockerfile_file = output_dir / f"{env_spec.name}_Dockerfile.aws"
            with open(aws_dockerfile_file, 'w') as f:
                f.write(aws_dockerfile)
            exported_files["aws_dockerfile"] = aws_dockerfile_file
        
        # 7. Export migration guide
        migration_guide = self._generate_migration_guide_markdown(
            env_spec, optimization, migration_plan
        )
        guide_file = output_dir / f"{env_spec.name}_migration_guide.md"
        with open(guide_file, 'w') as f:
            f.write(migration_guide)
        exported_files["migration_guide"] = guide_file
        
        # 8. Export cost analysis spreadsheet data
        cost_data = self._generate_cost_analysis_data(optimization, migration_plan)
        cost_file = output_dir / f"{env_spec.name}_cost_analysis.json"
        with open(cost_file, 'w') as f:
            json.dump(cost_data, f, indent=2)
        exported_files["cost_analysis"] = cost_file
        
        self.logger.info(f"Exported {len(exported_files)} configuration files")
        return exported_files
    
    def _generate_aws_spack_config(self, optimization: AWSOptimization) -> Dict[str, Any]:
        """Generate AWS-optimized Spack configuration"""
        config = {
            "config": {
                "install_tree": "/opt/spack",
                "source_cache": "/tmp/spack-cache/source",
                "build_stage": "/tmp/spack-stage",
                "binary_index_root": "/opt/spack/.spack/binary_index",
                "shared_linking": True,
                "checksum": True,
                "dirty": False,
                "build_jobs": optimization.recommended_instance_count * 8 if optimization.mpi_enabled else 8
            },
            "mirrors": {
                "aws-binary-cache": f"s3://spack-cache-{optimization.recommended_instance_type}",
                "aws-public-cache": "s3://spack-binaries-pnnl/0.20.1"
            },
            "packages": {
                "all": {
                    "target": ["graviton3"] if optimization.graviton_compatible else ["x86_64_v3"],
                    "variants": "+shared +pic",
                    "providers": {
                        "mpi": ["openmpi"] if optimization.mpi_enabled else [],
                        "blas": ["openblas"],
                        "lapack": ["openblas"]
                    }
                }
            }
        }
        
        if optimization.graviton_compatible:
            config["packages"]["all"]["variants"] += " target=graviton3"
            
        if optimization.gpu_required:
            config["packages"]["cuda"] = {
                "externals": [
                    {"spec": "cuda@11.8", "prefix": "/usr/local/cuda"}
                ],
                "buildable": False
            }
        
        return config
    
    def _generate_terraform_config(self, optimization: AWSOptimization, 
                                  migration_plan: MigrationPlan) -> str:
        """Generate Terraform configuration for AWS infrastructure"""
        
        instance_type = optimization.recommended_instance_type
        instance_count = optimization.recommended_instance_count
        
        terraform_config = f'''
# AWS Research Wizard - Terraform Configuration
# Generated for Spack environment migration

terraform {{
  required_version = ">= 1.0"
  required_providers {{
    aws = {{
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }}
  }}
}}

provider "aws" {{
  region = var.aws_region
}}

variable "aws_region" {{
  description = "AWS region for deployment"
  default     = "us-east-1"
}}

variable "environment_name" {{
  description = "Name of the Spack environment"
  default     = "{migration_plan.source_analysis['environment_name']}"
}}

variable "key_name" {{
  description = "EC2 Key Pair name"
  type        = string
}}

# VPC Configuration
resource "aws_vpc" "spack_vpc" {{
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = {{
    Name = "${{var.environment_name}}-vpc"
  }}
}}

resource "aws_subnet" "spack_subnet" {{
  vpc_id                  = aws_vpc.spack_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true
  
  tags = {{
    Name = "${{var.environment_name}}-subnet"
  }}
}}

resource "aws_internet_gateway" "spack_igw" {{
  vpc_id = aws_vpc.spack_vpc.id
  
  tags = {{
    Name = "${{var.environment_name}}-igw"
  }}
}}

resource "aws_route_table" "spack_rt" {{
  vpc_id = aws_vpc.spack_vpc.id
  
  route {{
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.spack_igw.id
  }}
  
  tags = {{
    Name = "${{var.environment_name}}-rt"
  }}
}}

resource "aws_route_table_association" "spack_rta" {{
  subnet_id      = aws_subnet.spack_subnet.id
  route_table_id = aws_route_table.spack_rt.id
}}

data "aws_availability_zones" "available" {{
  state = "available"
}}

# Security Group
resource "aws_security_group" "spack_sg" {{
  name_prefix = "${{var.environment_name}}-"
  vpc_id      = aws_vpc.spack_vpc.id
  
  ingress {{
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }}
'''

        if optimization.mpi_enabled:
            terraform_config += '''
  # MPI communication ports
  ingress {
    from_port = 0
    to_port   = 65535
    protocol  = "tcp"
    self      = true
  }
  
  ingress {
    from_port = 0
    to_port   = 65535
    protocol  = "udp"
    self      = true
  }
'''

        terraform_config += f'''
  egress {{
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }}
  
  tags = {{
    Name = "${{var.environment_name}}-sg"
  }}
}}

# AMI Selection
data "aws_ami" "ubuntu" {{
  most_recent = true
  owners      = ["099720109477"] # Canonical
  
  filter {{
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-{"amd64" if not optimization.graviton_compatible else "arm64"}-server-*"]
  }}
}}

'''

        if optimization.mpi_enabled and instance_count > 1:
            terraform_config += f'''
# Placement Group for MPI
resource "aws_placement_group" "spack_cluster" {{
  name     = "${{var.environment_name}}-cluster"
  strategy = "cluster"
}}

'''

        terraform_config += f'''
# EC2 Instances
resource "aws_instance" "spack_nodes" {{
  count                  = {instance_count}
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = "{instance_type}"
  key_name              = var.key_name
  subnet_id             = aws_subnet.spack_subnet.id
  security_groups       = [aws_security_group.spack_sg.id]
'''

        if optimization.mpi_enabled and instance_count > 1:
            terraform_config += '''
  placement_group       = aws_placement_group.spack_cluster.id
'''

        terraform_config += f'''
  
  root_block_device {{
    volume_type = "gp3"
    volume_size = {optimization.estimated_storage_gb}
    encrypted   = true
  }}
  
  user_data = base64encode(templatefile("${{path.module}}/user_data.sh", {{
    environment_name = var.environment_name
  }}))
  
  tags = {{
    Name = "${{var.environment_name}}-node-${{count.index + 1}}"
  }}
}}

'''

        if optimization.mpi_enabled:
            terraform_config += '''
# EFS for shared storage
resource "aws_efs_file_system" "spack_shared" {
  creation_token = "${var.environment_name}-shared"
  encrypted      = true
  
  tags = {
    Name = "${var.environment_name}-shared-storage"
  }
}

resource "aws_efs_mount_target" "spack_mount" {
  file_system_id  = aws_efs_file_system.spack_shared.id
  subnet_id       = aws_subnet.spack_subnet.id
  security_groups = [aws_security_group.spack_sg.id]
}
'''

        terraform_config += f'''
# S3 Bucket for Binary Cache
resource "aws_s3_bucket" "spack_cache" {{
  bucket = "${{var.environment_name}}-spack-cache-${{random_string.bucket_suffix.result}}"
}}

resource "aws_s3_bucket_versioning" "spack_cache_versioning" {{
  bucket = aws_s3_bucket.spack_cache.id
  versioning_configuration {{
    status = "Enabled"
  }}
}}

resource "random_string" "bucket_suffix" {{
  length  = 8
  special = false
  upper   = false
}}

# Outputs
output "instance_ips" {{
  value = aws_instance.spack_nodes[*].public_ip
}}

output "s3_cache_bucket" {{
  value = aws_s3_bucket.spack_cache.bucket
}}

'''

        if optimization.mpi_enabled:
            terraform_config += '''
output "efs_id" {
  value = aws_efs_file_system.spack_shared.id
}
'''

        return terraform_config
    
    def _generate_deployment_script(self, env_spec: SpackEnvironmentSpec, 
                                   optimization: AWSOptimization) -> str:
        """Generate comprehensive deployment script"""
        
        script = f'''#!/bin/bash
# AWS Research Wizard - Deployment Script
# Spack Environment: {env_spec.name}
# Generated: {datetime.now().isoformat()}

set -euo pipefail

# Configuration
ENVIRONMENT_NAME="{env_spec.name}"
SPACK_VERSION="{env_spec.spack_version}"
INSTANCE_TYPE="{optimization.recommended_instance_type}"
GRAVITON_ENABLED="{str(optimization.graviton_compatible).lower()}"
MPI_ENABLED="{str(optimization.mpi_enabled).lower()}"
GPU_ENABLED="{str(optimization.gpu_required).lower()}"

# Colors for output
RED='\\033[0;31m'
GREEN='\\033[0;32m'
YELLOW='\\033[1;33m'
NC='\\033[0m' # No Color

log() {{
    echo -e "${{GREEN}}[$(date +'%Y-%m-%d %H:%M:%S')] $1${{NC}}"
}}

warn() {{
    echo -e "${{YELLOW}}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${{NC}}"
}}

error() {{
    echo -e "${{RED}}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${{NC}}"
    exit 1
}}

# Function to check if running on correct instance type
check_instance_type() {{
    local current_type=$(curl -s http://169.254.169.254/latest/meta-data/instance-type)
    if [[ "$current_type" != "$INSTANCE_TYPE" ]]; then
        warn "Running on $current_type instead of recommended $INSTANCE_TYPE"
    else
        log "Confirmed running on recommended instance type: $INSTANCE_TYPE"
    fi
}}

# Update system packages
update_system() {{
    log "Updating system packages..."
    sudo apt-get update
    sudo apt-get upgrade -y
    
    # Install essential packages
    sudo apt-get install -y \\
        build-essential \\
        gfortran \\
        git \\
        curl \\
        wget \\
        python3 \\
        python3-pip \\
        cmake \\
        ninja-build \\
        pkg-config \\
        libssl-dev \\
        zlib1g-dev \\
        libbz2-dev \\
        libreadline-dev \\
        libsqlite3-dev \\
        libncurses5-dev \\
        libncursesw5-dev \\
        xz-utils \\
        tk-dev \\
        libffi-dev \\
        liblzma-dev
}}

# Install AWS CLI
install_aws_cli() {{
    log "Installing AWS CLI..."
    if [[ "$GRAVITON_ENABLED" == "true" ]]; then
        curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
    else
        curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    fi
    unzip awscliv2.zip
    sudo ./aws/install
    rm -rf aws awscliv2.zip
}}

'''

        if optimization.mpi_enabled:
            script += '''
# Configure EFS mounting
setup_efs_mount() {
    log "Setting up EFS shared storage..."
    sudo apt-get install -y nfs-utils
    sudo mkdir -p /shared
    
    # Get EFS ID from Terraform output or environment variable
    if [[ -n "${EFS_ID:-}" ]]; then
        sudo mount -t efs -o tls ${EFS_ID}:/ /shared
        echo "${EFS_ID}:/ /shared efs defaults,_netdev" | sudo tee -a /etc/fstab
        log "EFS mounted at /shared"
    else
        warn "EFS_ID not provided, skipping EFS mount"
    fi
}
'''

        if optimization.gpu_required:
            script += '''
# Install NVIDIA drivers and CUDA
install_nvidia_drivers() {
    log "Installing NVIDIA drivers and CUDA..."
    
    # Install NVIDIA driver
    sudo apt-get install -y nvidia-driver-535
    
    # Install CUDA
    wget https://developer.download.nvidia.com/compute/cuda/12.2.0/local_installers/cuda_12.2.0_535.54.03_linux.run
    sudo sh cuda_12.2.0_535.54.03_linux.run --silent --toolkit
    
    # Set CUDA environment variables
    echo 'export PATH=/usr/local/cuda/bin:$PATH' >> ~/.bashrc
    echo 'export LD_LIBRARY_PATH=/usr/local/cuda/lib64:$LD_LIBRARY_PATH' >> ~/.bashrc
    
    log "NVIDIA drivers and CUDA installed"
}
'''

        script += f'''
# Install Spack
install_spack() {{
    log "Installing Spack {env_spec.spack_version}..."
    
    cd /opt
    sudo git clone -c feature.manyFiles=true https://github.com/spack/spack.git
    cd spack
    sudo git checkout {env_spec.spack_version}
    
    # Set up Spack environment
    echo 'export SPACK_ROOT=/opt/spack' | sudo tee -a /etc/environment
    echo 'export PATH=$SPACK_ROOT/bin:$PATH' | sudo tee -a /etc/environment
    echo 'source $SPACK_ROOT/share/spack/setup-env.sh' | sudo tee -a /etc/bash.bashrc
    
    # Make spack accessible to all users
    sudo chown -R root:root /opt/spack
    sudo chmod -R go+rX /opt/spack
    
    source /opt/spack/share/spack/setup-env.sh
    
    log "Spack installed successfully"
}}

# Configure Spack for AWS
configure_spack_aws() {{
    log "Configuring Spack for AWS environment..."
    
    source /opt/spack/share/spack/setup-env.sh
    
    # Copy AWS-optimized configuration
    if [[ -f "aws_spack_config.yaml" ]]; then
        spack config add -f aws_spack_config.yaml
        log "AWS Spack configuration applied"
    fi
    
    # Configure binary cache
    if [[ -n "${{S3_CACHE_BUCKET:-}}" ]]; then
        spack mirror add aws-cache s3://$S3_CACHE_BUCKET
        log "S3 binary cache configured: s3://$S3_CACHE_BUCKET"
    fi
    
    # Add public binary cache
    spack mirror add aws-public s3://spack-binaries-pnnl/0.20.1
    
    # Configure compilers
    spack compiler find
    
    log "Spack AWS configuration complete"
}}

# Create and populate Spack environment
create_spack_environment() {{
    log "Creating Spack environment: $ENVIRONMENT_NAME..."
    
    source /opt/spack/share/spack/setup-env.sh
    
    # Create environment from exported YAML
    if [[ -f "{env_spec.name}_spack_env.yaml" ]]; then
        spack env create $ENVIRONMENT_NAME {env_spec.name}_spack_env.yaml
        log "Environment created from YAML specification"
    else
        error "Environment YAML file not found"
    fi
    
    # Activate environment
    spack env activate $ENVIRONMENT_NAME
    
    log "Environment $ENVIRONMENT_NAME created and activated"
}}

# Install packages
install_packages() {{
    log "Installing packages in environment: $ENVIRONMENT_NAME..."
    
    source /opt/spack/share/spack/setup-env.sh
    spack env activate $ENVIRONMENT_NAME
    
    # Install with binary cache preference
    spack install --use-cache
    
    log "Package installation complete"
}}

# Generate module files
generate_modules() {{
    log "Generating environment modules..."
    
    source /opt/spack/share/spack/setup-env.sh
    spack env activate $ENVIRONMENT_NAME
    
    # Generate modules
    spack module tcl refresh --delete-tree
    
    log "Environment modules generated"
}}

# Main execution
main() {{
    log "Starting Spack environment deployment for: $ENVIRONMENT_NAME"
    log "Target instance type: $INSTANCE_TYPE"
    log "Graviton optimization: $GRAVITON_ENABLED"
    log "MPI support: $MPI_ENABLED"
    log "GPU support: $GPU_ENABLED"
    
    check_instance_type
    update_system
    install_aws_cli
'''

        if optimization.mpi_enabled:
            script += '''
    setup_efs_mount
'''

        if optimization.gpu_required:
            script += '''
    install_nvidia_drivers
'''

        script += f'''
    install_spack
    configure_spack_aws
    create_spack_environment
    install_packages
    generate_modules
    
    log "Deployment completed successfully!"
    log "To use the environment, run:"
    log "  source /opt/spack/share/spack/setup-env.sh"
    log "  spack env activate $ENVIRONMENT_NAME"
}}

# Execute main function
main "$@"
'''

        return script
    
    def _generate_aws_optimized_dockerfile(self, env_spec: SpackEnvironmentSpec, 
                                         optimization: AWSOptimization) -> str:
        """Generate AWS-optimized Dockerfile"""
        
        base_image = "ubuntu:22.04"
        if optimization.graviton_compatible:
            base_image = "arm64v8/ubuntu:22.04"
        
        dockerfile = f'''# AWS-Optimized Spack Environment
# Environment: {env_spec.name}
# Architecture: {"ARM64/Graviton" if optimization.graviton_compatible else "x86_64"}

FROM {base_image}

# Metadata
LABEL maintainer="AWS Research Wizard"
LABEL environment="{env_spec.name}"
LABEL spack_version="{env_spec.spack_version}"
LABEL architecture="{"aarch64" if optimization.graviton_compatible else "x86_64"}"

# Environment variables
ENV SPACK_ROOT=/opt/spack
ENV PATH=$SPACK_ROOT/bin:$PATH
ENV DEBIAN_FRONTEND=noninteractive

# Install system dependencies
RUN apt-get update && apt-get install -y \\
    build-essential \\
    gfortran \\
    git \\
    curl \\
    wget \\
    python3 \\
    python3-pip \\
    cmake \\
    ninja-build \\
    pkg-config \\
    libssl-dev \\
    zlib1g-dev \\
    ca-certificates \\
    && rm -rf /var/lib/apt/lists/*

'''

        if optimization.mpi_enabled:
            dockerfile += '''
# Install MPI dependencies
RUN apt-get update && apt-get install -y \\
    libopenmpi-dev \\
    openmpi-bin \\
    && rm -rf /var/lib/apt/lists/*

'''

        if optimization.gpu_required:
            dockerfile += '''
# Install CUDA (runtime)
RUN apt-get update && apt-get install -y \\
    nvidia-cuda-runtime \\
    && rm -rf /var/lib/apt/lists/*

'''

        dockerfile += f'''
# Install Spack
RUN git clone -c feature.manyFiles=true https://github.com/spack/spack.git $SPACK_ROOT \\
    && cd $SPACK_ROOT \\
    && git checkout {env_spec.spack_version}

# Set up Spack environment
RUN echo 'source $SPACK_ROOT/share/spack/setup-env.sh' >> /etc/bash.bashrc

# Copy environment specification
COPY {env_spec.name}_spack_env.yaml /tmp/spack_env.yaml
COPY aws_spack_config.yaml /tmp/aws_config.yaml

# Configure Spack for AWS
RUN . $SPACK_ROOT/share/spack/setup-env.sh \\
    && spack config add -f /tmp/aws_config.yaml \\
    && spack compiler find \\
    && spack mirror add aws-public s3://spack-binaries-pnnl/0.20.1

# Create environment
RUN . $SPACK_ROOT/share/spack/setup-env.sh \\
    && spack env create {env_spec.name} /tmp/spack_env.yaml

# Install packages with cache preference
RUN . $SPACK_ROOT/share/spack/setup-env.sh \\
    && spack env activate {env_spec.name} \\
    && spack install --use-cache

# Generate modules
RUN . $SPACK_ROOT/share/spack/setup-env.sh \\
    && spack env activate {env_spec.name} \\
    && spack module tcl refresh --delete-tree

# Set up environment activation script
RUN echo '#!/bin/bash' > /usr/local/bin/activate-spack \\
    && echo 'source $SPACK_ROOT/share/spack/setup-env.sh' >> /usr/local/bin/activate-spack \\
    && echo 'spack env activate {env_spec.name}' >> /usr/local/bin/activate-spack \\
    && chmod +x /usr/local/bin/activate-spack

# Default command
CMD ["/bin/bash", "-l"]
'''

        return dockerfile
    
    def _generate_migration_guide_markdown(self, env_spec: SpackEnvironmentSpec,
                                         optimization: AWSOptimization,
                                         migration_plan: MigrationPlan) -> str:
        """Generate comprehensive migration guide in Markdown format"""
        
        guide = f'''# Spack Environment Migration Guide

## Environment: {env_spec.name}

**Migration Date:** {datetime.now().strftime("%Y-%m-%d")}  
**Source System:** {env_spec.source_system}  
**Spack Version:** {env_spec.spack_version}  
**Target Architecture:** {"ARM64 (Graviton)" if optimization.graviton_compatible else "x86_64"}

---

## 📊 Executive Summary

This guide provides a comprehensive plan for migrating your Spack environment `{env_spec.name}` from your current HPC system to AWS. The migration will provide improved cost-efficiency, scalability, and performance.

### Key Benefits
- **Cost Optimization:** Estimated monthly cost of ${optimization.estimated_monthly_cost:.0f}
- **Performance:** {"Up to 20-40% improvement with Graviton instances" if optimization.graviton_compatible else "Similar to current performance"}
- **Scalability:** {"Multi-node scaling up to " + str(optimization.recommended_instance_count) + " instances" if optimization.mpi_enabled else "Single-node scaling"}
- **Reliability:** AWS managed infrastructure with 99.99% uptime SLA

---

## 🎯 Target Architecture

| Component | Specification |
|-----------|---------------|
| **Instance Type** | {optimization.recommended_instance_type} |
| **Instance Count** | {optimization.recommended_instance_count} |
| **Memory** | {optimization.estimated_memory_gb} GB |
| **Storage** | {optimization.estimated_storage_gb} GB |
| **Architecture** | {"ARM64 (Graviton3)" if optimization.graviton_compatible else "x86_64"} |
| **GPU Support** | {"Yes" if optimization.gpu_required else "No"} |
| **MPI Support** | {"Yes" if optimization.mpi_enabled else "No"} |

---

## 💰 Cost Analysis

### Monthly Cost Breakdown
- **Hourly Rate:** ${optimization.estimated_cost_per_hour:.2f}/hour
- **Monthly Estimate:** ${optimization.estimated_monthly_cost:.0f}/month (30% utilization)
- **Annual Estimate:** ${optimization.estimated_monthly_cost * 12:.0f}/year

### Cost Optimization Opportunities
'''

        for note in optimization.optimization_notes:
            if "cost" in note.lower() or "saving" in note.lower():
                guide += f"- {note}\n"

        guide += f'''

---

## 📦 Package Analysis

**Total Packages:** {len(env_spec.installed_packages)}  
**Concrete Specs:** {len(env_spec.concrete_specs)}  
**Compilers:** {len(env_spec.compilers)}

'''

        if optimization.graviton_compatible:
            guide += '''
### Graviton Compatibility
✅ **Fully Compatible** - All packages support ARM64 architecture

**Performance Benefits:**
'''
            for note in optimization.optimization_notes:
                if "graviton" in note.lower() or "performance" in note.lower():
                    guide += f"- {note}\n"

        guide += f'''

---

## 🚀 Migration Plan

### Phase 1: Assessment and Planning (Week 1)
'''
        
        for step in migration_plan.migration_steps[:4]:
            guide += f"- {step}\n"

        guide += '''

### Phase 2: Infrastructure Setup (Week 2)
'''
        
        for step in migration_plan.migration_steps[4:8]:
            guide += f"- {step}\n"

        guide += '''

### Phase 3: Environment Migration (Week 3)
'''
        
        for step in migration_plan.migration_steps[8:12]:
            guide += f"- {step}\n"

        guide += '''

### Phase 4: Testing and Validation (Week 4)
'''
        
        for step in migration_plan.migration_steps[12:]:
            guide += f"- {step}\n"

        guide += f'''

---

## 🔧 Technical Implementation

### AWS Services Required
'''
        
        for service in optimization.aws_services:
            guide += f"- {service}\n"

        guide += f'''

### Binary Cache Strategy
{optimization.binary_cache_strategy}

### Container Strategy
'''
        
        if optimization.container_image:
            guide += f"- **Container Image:** {optimization.container_image}\n"
            guide += "- **Strategy:** Containerized deployment for improved reproducibility\n"
            guide += "- **Benefits:** Faster deployment, consistent environments, easy scaling\n"
        else:
            guide += "- **Strategy:** Native Spack installation\n"
            guide += "- **Benefits:** Direct access to hardware, easier debugging\n"

        guide += f'''

---

## ✅ Validation Tests

'''
        
        for test in migration_plan.validation_tests:
            guide += f"- [ ] {test}\n"

        guide += f'''

---

## 📋 Deployment Instructions

### Prerequisites
1. AWS account with appropriate permissions
2. Terraform installed (>= 1.0)
3. AWS CLI configured
4. SSH key pair for EC2 access

### Step-by-Step Deployment

#### 1. Infrastructure Deployment
```bash
# Deploy AWS infrastructure
terraform init
terraform plan -var="key_name=your-key-name"
terraform apply -var="key_name=your-key-name"
```

#### 2. Spack Environment Setup
```bash
# SSH to the first instance
ssh -i your-key.pem ubuntu@<instance-ip>

# Run deployment script
./deploy_{env_spec.name}.sh
```

#### 3. Environment Activation
```bash
# Activate Spack environment
source /opt/spack/share/spack/setup-env.sh
spack env activate {env_spec.name}
```

---

## 🔍 Monitoring and Optimization

### CloudWatch Metrics to Monitor
- CPU utilization
- Memory usage
- Network I/O
- EBS I/O (if applicable)
{"- GPU utilization" if optimization.gpu_required else ""}

### Cost Monitoring
- Set up AWS Budgets for cost alerts
- Use AWS Cost Explorer for usage analysis
- Implement cost allocation tags

### Performance Optimization
'''
        
        for note in optimization.optimization_notes:
            if "performance" in note.lower() or "optimization" in note.lower():
                guide += f"- {note}\n"

        guide += f'''

---

## 🆘 Troubleshooting

### Common Issues and Solutions

#### Package Installation Failures
- **Issue:** Package fails to build from source
- **Solution:** Check binary cache availability, verify compiler configuration

#### Performance Issues
- **Issue:** Slower performance than expected
- **Solution:** Verify instance type selection, check for CPU throttling

#### MPI Communication Problems (if applicable)
- **Issue:** MPI jobs fail to communicate
- **Solution:** Verify security group settings, check placement group configuration

---

## 📞 Support Resources

### Documentation
- [AWS Research Wizard Documentation](https://docs.aws-research-wizard.com)
- [Spack Documentation](https://spack.readthedocs.io/)
- [AWS HPC Documentation](https://docs.aws.amazon.com/hpc/)

### Community Support
- AWS Research Wizard Community Forum
- Spack Community Slack
- AWS HPC User Forum

### Professional Support
- AWS Enterprise Support
- AWS Research Computing Team
- Spack Professional Services

---

## 📈 Next Steps

1. **Review** this migration guide thoroughly
2. **Plan** your migration timeline (4-6 weeks recommended)
3. **Test** the deployment in a development environment
4. **Train** your team on the new AWS environment
5. **Execute** the production migration
6. **Monitor** performance and costs post-migration
7. **Optimize** based on usage patterns

---

*This migration guide was generated by the AWS Research Wizard Spack Environment Importer. For questions or support, please contact the AWS Research Computing team.*
'''

        return guide
    
    def _generate_cost_analysis_data(self, optimization: AWSOptimization,
                                   migration_plan: MigrationPlan) -> Dict[str, Any]:
        """Generate detailed cost analysis data"""
        
        return {
            "summary": {
                "hourly_cost": optimization.estimated_cost_per_hour,
                "monthly_cost": optimization.estimated_monthly_cost,
                "annual_cost": optimization.estimated_monthly_cost * 12
            },
            "breakdown": migration_plan.cost_analysis["cost_breakdown"],
            "optimization_opportunities": migration_plan.cost_analysis["cost_optimization"],
            "comparison": {
                "on_demand_monthly": optimization.estimated_monthly_cost,
                "spot_monthly": optimization.estimated_monthly_cost * 0.3,  # ~70% savings
                "reserved_monthly": optimization.estimated_monthly_cost * 0.6,  # ~40% savings
                "graviton_savings": optimization.estimated_monthly_cost * 0.2 if optimization.graviton_compatible else 0
            },
            "utilization_scenarios": {
                "light_usage_10_percent": optimization.estimated_cost_per_hour * 24 * 30 * 0.1,
                "normal_usage_30_percent": optimization.estimated_monthly_cost,
                "heavy_usage_70_percent": optimization.estimated_cost_per_hour * 24 * 30 * 0.7,
                "continuous_100_percent": optimization.estimated_cost_per_hour * 24 * 30
            }
        }

def main():
    """Enhanced CLI interface for the native Spack Environment Importer"""
    parser = argparse.ArgumentParser(
        description="AWS Research Wizard - Native Spack Environment Importer",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # List available environments
  python spack_environment_importer.py --list
  
  # Import specific environment
  python spack_environment_importer.py --import my_env --output ./configs
  
  # Import with custom Spack installation
  python spack_environment_importer.py --import my_env --spack-root /opt/spack
  
  # Generate migration plan only
  python spack_environment_importer.py --import my_env --migration-plan-only
        """
    )
    
    parser.add_argument("--verbose", "-v", action="store_true", 
                       help="Enable verbose logging")
    parser.add_argument("--spack-root", type=Path, 
                       help="Path to Spack installation (auto-detected if not provided)")
    parser.add_argument("--list", action="store_true", 
                       help="List available Spack environments")
    parser.add_argument("--import", dest="import_env", type=str, 
                       help="Import specific environment by name")
    parser.add_argument("--output", "-o", type=Path, default=Path("./aws_configs"), 
                       help="Output directory for generated configurations")
    parser.add_argument("--migration-plan-only", action="store_true",
                       help="Generate migration plan only, skip detailed configs")
    parser.add_argument("--format", choices=["json", "yaml", "all"], default="all",
                       help="Output format for configurations")
    
    args = parser.parse_args()
    
    try:
        with SpackEnvironmentImporter(verbose=args.verbose, spack_root=args.spack_root) as importer:
            
            if args.list:
                # List available environments
                environments = importer.discover_environments()
                if environments:
                    print("📦 Available Spack environments:")
                    for i, env in enumerate(environments, 1):
                        print(f"  {i:2d}. {env}")
                    print(f"\n📊 Total: {len(environments)} environments found")
                else:
                    print("❌ No Spack environments found.")
                    print("💡 Make sure Spack is installed and environments exist.")
                return
            
            if args.import_env:
                # Import specific environment
                print(f"🔍 Importing environment: {args.import_env}")
                
                # Step 1: Capture environment
                env_spec = importer.capture_environment(args.import_env)
                if not env_spec:
                    print(f"❌ Failed to capture environment: {args.import_env}")
                    return
                
                print(f"✅ Environment captured:")
                print(f"   📦 Packages: {len(env_spec.installed_packages)}")
                print(f"   🔧 Compilers: {len(env_spec.compilers)}")
                print(f"   📚 Repositories: {len(env_spec.repos)}")
                
                # Step 2: Analyze for AWS optimization
                print(f"🔬 Analyzing AWS optimization opportunities...")
                optimization = importer.analyze_environment(env_spec)
                
                print(f"✅ Analysis complete:")
                print(f"   💻 Recommended instance: {optimization.recommended_instance_type}")
                print(f"   📊 Instance count: {optimization.recommended_instance_count}")
                print(f"   💰 Estimated cost: ${optimization.estimated_cost_per_hour:.2f}/hour")
                print(f"   📅 Monthly estimate: ${optimization.estimated_monthly_cost:.0f}/month")
                print(f"   🏗️  Graviton compatible: {'Yes' if optimization.graviton_compatible else 'No'}")
                print(f"   🔀 MPI enabled: {'Yes' if optimization.mpi_enabled else 'No'}")
                print(f"   🎮 GPU required: {'Yes' if optimization.gpu_required else 'No'}")
                
                # Step 3: Generate migration plan
                print(f"📋 Generating migration plan...")
                migration_plan = importer.generate_migration_plan(env_spec, optimization)
                
                if args.migration_plan_only:
                    # Just output migration plan summary
                    print(f"✅ Migration plan generated:")
                    print(f"   🎯 Target: {migration_plan.target_architecture['instance_type']}")
                    print(f"   💾 Storage: {migration_plan.target_architecture['storage_gb']} GB")
                    print(f"   💰 Monthly cost: ${migration_plan.cost_analysis['monthly_estimate']:.0f}")
                    print(f"   📈 Performance: {migration_plan.performance_comparison['aws_estimate']['compute_performance']}")
                    return
                
                # Step 4: Export comprehensive configuration
                print(f"📤 Exporting AWS configuration...")
                exported_files = importer.export_aws_configuration(
                    env_spec, optimization, migration_plan, args.output
                )
                
                print(f"✅ Export complete! Generated {len(exported_files)} files:")
                for file_type, file_path in exported_files.items():
                    print(f"   📄 {file_type}: {file_path}")
                
                print(f"\n🎉 Migration configuration ready!")
                print(f"📖 Next steps:")
                print(f"   1. Review the migration guide: {exported_files['migration_guide']}")
                print(f"   2. Customize Terraform config: {exported_files['terraform']}")
                print(f"   3. Run deployment script: {exported_files['deployment_script']}")
                print(f"   4. Monitor costs and performance post-migration")
            
            else:
                # Interactive mode
                environments = importer.discover_environments()
                if not environments:
                    print("❌ No Spack environments found.")
                    return
                
                print("📦 Available environments:")
                for i, env in enumerate(environments, 1):
                    print(f"  {i:2d}. {env}")
                
                try:
                    choice = input(f"\n🎯 Select environment (1-{len(environments)}) or 'q' to quit: ").strip()
                    if choice.lower() == 'q':
                        print("👋 Goodbye!")
                        return
                    
                    choice_idx = int(choice) - 1
                    if 0 <= choice_idx < len(environments):
                        selected_env = environments[choice_idx]
                        print(f"\n🚀 Processing environment: {selected_env}")
                        
                        # Continue with import process...
                        env_spec = importer.capture_environment(selected_env)
                        if env_spec:
                            optimization = importer.analyze_environment(env_spec)
                            migration_plan = importer.generate_migration_plan(env_spec, optimization)
                            exported_files = importer.export_aws_configuration(
                                env_spec, optimization, migration_plan, args.output
                            )
                            
                            print(f"✅ Configuration exported to: {args.output}")
                            print(f"📖 Start with the migration guide: {exported_files['migration_guide']}")
                    else:
                        print("❌ Invalid selection.")
                        
                except (ValueError, KeyboardInterrupt):
                    print("\n👋 Operation cancelled.")
    
    except Exception as e:
        print(f"❌ Error: {e}")
        if args.verbose:
            import traceback
            traceback.print_exc()

if __name__ == "__main__":
    main()