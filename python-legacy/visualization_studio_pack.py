#!/usr/bin/env python3
"""
Visualization Studio Pack
Comprehensive scientific visualization environments for AWS Research Wizard
"""

import json
from typing import Dict, List, Any
from dataclasses import dataclass

@dataclass
class VisualizationConfig:
    """Configuration for scientific visualization environments"""
    domain_focus: str
    visualization_tools: List[str]
    gpu_acceleration: bool
    remote_display: bool
    collaboration_features: bool
    data_size_gb: int
    concurrent_users: int

class VisualizationStudioPack:
    """
    Scientific Visualization Studio configurations optimized for AWS
    Supports interactive visualization, remote rendering, and collaborative analysis
    """

    def __init__(self):
        self.visualization_domains = {
            "general_purpose": self._get_general_purpose_config(),
            "high_performance": self._get_hpc_visualization_config(),
            "interactive_jupyter": self._get_jupyter_visualization_config(),
            "collaborative_studio": self._get_collaborative_config(),
            "gpu_accelerated": self._get_gpu_visualization_config(),
            "medical_imaging": self._get_medical_imaging_config(),
            "geospatial": self._get_geospatial_config(),
            "molecular_visualization": self._get_molecular_config()
        }

    def _get_general_purpose_config(self) -> Dict[str, Any]:
        """General-purpose scientific visualization environment"""
        return {
            "name": "General Purpose Visualization Studio",
            "description": "Balanced visualization environment for diverse research needs",
            "spack_packages": [
                # Core visualization frameworks
                "paraview@5.11.2 %gcc@11.4.0 +python +mpi +osmesa +qt",
                "visit@3.3.3 %gcc@11.4.0 +python +mpi +gui",
                "vtk@9.2.6 %gcc@11.4.0 +python +mpi +opengl2",

                # Python visualization stack
                "python@3.11.5 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-mayavi@4.8.1 %gcc@11.4.0",

                # Data processing
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",

                # Jupyter ecosystem
                "py-jupyter@1.0.0 %gcc@11.4.0",
                "py-jupyterlab@4.0.5 %gcc@11.4.0",
                "py-ipywidgets@8.1.0 %gcc@11.4.0",

                # Graphics libraries
                "mesa@23.1.6 %gcc@11.4.0 +opengl +osmesa",
                "freeglut@3.4.0 %gcc@11.4.0",
                "glfw@3.3.8 %gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "small_team": {
                    "instance_type": "g4dn.xlarge",
                    "vcpus": 4,
                    "memory_gb": 16,
                    "gpu": "1x NVIDIA T4",
                    "storage": "125GB NVMe SSD",
                    "cost_per_hour": 0.526,
                    "use_case": "1-2 users, moderate datasets"
                },
                "research_group": {
                    "instance_type": "g4dn.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 64,
                    "gpu": "1x NVIDIA T4",
                    "storage": "900GB NVMe SSD",
                    "cost_per_hour": 1.204,
                    "use_case": "3-8 users, large datasets"
                },
                "high_performance": {
                    "instance_type": "g5.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 64,
                    "gpu": "1x NVIDIA A10G",
                    "storage": "600GB NVMe SSD",
                    "cost_per_hour": 1.624,
                    "use_case": "Real-time rendering, complex visualizations"
                }
            },
            "deployment_features": [
                "JupyterHub multi-user environment",
                "Remote desktop access (VNC/RDP)",
                "Shared project storage (EFS)",
                "Auto-scaling based on usage",
                "GPU-accelerated rendering",
                "Web-based visualization access"
            ],
            "cost_profile": {
                "idle_cost": "$0 (auto-shutdown)",
                "small_usage": "$50-150/month (occasional visualization)",
                "regular_usage": "$200-600/month (active research)",
                "intensive_usage": "$600-1200/month (daily rendering)"
            }
        }

    def _get_hpc_visualization_config(self) -> Dict[str, Any]:
        """High-performance visualization for large-scale simulations"""
        return {
            "name": "HPC Visualization Cluster",
            "description": "Distributed visualization for massive scientific datasets",
            "spack_packages": [
                # Parallel visualization tools
                "paraview@5.11.2 %gcc@11.4.0 +python +mpi +osmesa +adios2 +catalyst",
                "visit@3.3.3 %gcc@11.4.0 +python +mpi +parallel +hdf5",
                "vtk@9.2.6 %gcc@11.4.0 +python +mpi +opengl2 +ioxml",

                # In-situ visualization
                "catalyst@2.0.0 %gcc@11.4.0 +python +mpi",
                "ascent@0.9.1 %gcc@11.4.0 +python +mpi +openmp",

                # Parallel I/O
                "adios2@2.9.1 %gcc@11.4.0 +python +mpi +hdf5",
                "hdf5@1.14.2 %gcc@11.4.0 +mpi +threadsafe +fortran",
                "parallel-netcdf@1.12.3 %gcc@11.4.0",

                # MPI and parallel libraries
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers",
                "petsc@3.19.4 %gcc@11.4.0 +mpi +hypre",

                # Remote rendering
                "ospray@2.12.0 %gcc@11.4.0 +mpi +openimageio",
                "embree@4.1.0 %gcc@11.4.0",

                # Python parallel computing
                "py-mpi4py@3.1.4 %gcc@11.4.0",
                "py-dask@2023.8.0 %gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "visualization_cluster": {
                    "head_node": {
                        "instance_type": "g5.2xlarge",
                        "vcpus": 8,
                        "memory_gb": 32,
                        "gpu": "1x NVIDIA A10G",
                        "role": "Interactive visualization and job management"
                    },
                    "render_nodes": {
                        "instance_type": "g5.xlarge",
                        "vcpus": 4,
                        "memory_gb": 16,
                        "gpu": "1x NVIDIA A10G",
                        "count": "2-8 nodes",
                        "role": "Distributed rendering"
                    },
                    "storage": "FSx Lustre 1.2TB (200MB/s)"
                },
                "cost_per_hour": "$3.24-12.96 (2-8 render nodes)",
                "scaling": "Auto-scale render nodes based on workload"
            },
            "specialized_features": [
                "In-situ visualization during simulation",
                "Distributed parallel rendering",
                "Cinema database generation",
                "Remote collaborative viewing",
                "Automatic data streaming",
                "GPU cluster rendering"
            ],
            "cost_profile": {
                "idle_cost": "$0 (complete shutdown)",
                "burst_rendering": "$50-200/day intensive sessions",
                "continuous_workflow": "$500-2000/month for active projects",
                "large_campaigns": "$2000-8000/month for major simulations"
            }
        }

    def _get_jupyter_visualization_config(self) -> Dict[str, Any]:
        """Interactive Jupyter-based visualization environment"""
        return {
            "name": "Interactive Jupyter Visualization Studio",
            "description": "Web-based interactive visualization and analysis platform",
            "spack_packages": [
                # Jupyter ecosystem
                "py-jupyterlab@4.0.5 %gcc@11.4.0",
                "py-jupyter@1.0.0 %gcc@11.4.0",
                "py-ipywidgets@8.1.0 %gcc@11.4.0",
                "py-voila@0.5.0 %gcc@11.4.0",
                "py-jupyterlab-widgets@3.0.8 %gcc@11.4.0",

                # Interactive visualization
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                "py-altair@5.0.1 %gcc@11.4.0",
                "py-holoviews@1.17.1 %gcc@11.4.0",
                "py-panel@1.2.1 %gcc@11.4.0",
                "py-datashader@0.15.1 %gcc@11.4.0",

                # 3D visualization
                "py-mayavi@4.8.1 %gcc@11.4.0",
                "py-pyvista@0.42.2 %gcc@11.4.0",
                "py-k3d@2.15.2 %gcc@11.4.0",
                "py-pythreejs@2.4.2 %gcc@11.4.0",

                # Scientific computing
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-dask@2023.8.0 %gcc@11.4.0",

                # Machine learning visualization
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-plotnine@0.12.2 %gcc@11.4.0"
            ],
            "deployment_architecture": {
                "jupyterhub_config": {
                    "multi_user": True,
                    "authentication": "OAuth2/LDAP",
                    "spawner": "KubernetesSpawner",
                    "shared_storage": "EFS"
                },
                "container_orchestration": "EKS (Kubernetes)",
                "auto_scaling": "Horizontal Pod Autoscaler",
                "gpu_scheduling": "NVIDIA GPU Operator"
            },
            "aws_instance_recommendations": {
                "jupyterhub_hub": {
                    "instance_type": "t3.medium",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "role": "JupyterHub orchestration"
                },
                "user_notebooks": {
                    "cpu_instances": "t3.large to c5.2xlarge",
                    "gpu_instances": "g4dn.xlarge to g5.4xlarge",
                    "memory_instances": "r5.large to r5.4xlarge",
                    "scaling": "0 to 20 instances based on demand"
                }
            },
            "cost_profile": {
                "base_cost": "$35/month (always-on hub)",
                "per_user_active": "$0.50-4.00/hour based on instance size",
                "typical_team": "$150-500/month (5-15 users)",
                "large_class": "$300-1000/month (20-50 users)"
            }
        }

    def _get_collaborative_config(self) -> Dict[str, Any]:
        """Multi-user collaborative visualization environment"""
        return {
            "name": "Collaborative Research Visualization Platform",
            "description": "Shared visualization environment for research teams",
            "spack_packages": [
                # Collaborative tools
                "paraview@5.11.2 %gcc@11.4.0 +python +mpi +qt +web",
                "trame@2.5.1 %gcc@11.4.0",
                "py-trame-vuetify@2.3.1 %gcc@11.4.0",
                "py-trame-vtk@2.5.8 %gcc@11.4.0",

                # Web-based visualization
                "vtk@9.2.6 %gcc@11.4.0 +python +web +opengl2",
                "py-panel@1.2.1 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",

                # Collaboration infrastructure
                "apache-server@2.4.57 %gcc@11.4.0",
                "nginx@1.25.2 %gcc@11.4.0",
                "py-tornado@6.3.3 %gcc@11.4.0",

                # Version control and sharing
                "git@2.41.0 %gcc@11.4.0",
                "git-lfs@3.4.0 %gcc@11.4.0",
                "py-nbdime@3.2.1 %gcc@11.4.0"
            ],
            "collaboration_features": [
                "Real-time shared visualization sessions",
                "Annotation and markup tools",
                "Session recording and playback",
                "Multi-user cursor tracking",
                "Integrated video conferencing",
                "Shared notebook environments",
                "Version controlled visualization scripts",
                "Export to presentation formats"
            ],
            "aws_services_integration": [
                "AWS WorkSpaces for persistent desktops",
                "Amazon Connect for integrated voice/video",
                "S3 for shared data storage",
                "CloudFront for global content delivery",
                "WAF for security",
                "Cognito for user authentication"
            ],
            "cost_profile": {
                "small_team": "$200-500/month (3-8 users)",
                "research_group": "$500-1200/month (8-20 users)",
                "department": "$1000-3000/month (20-50 users)",
                "additional_storage": "$0.023/GB-month (S3 Standard)"
            }
        }

    def _get_gpu_visualization_config(self) -> Dict[str, Any]:
        """High-end GPU-accelerated visualization"""
        return {
            "name": "GPU-Accelerated Visualization Workstation",
            "description": "High-performance GPU rendering for complex visualizations",
            "spack_packages": [
                # GPU-accelerated visualization
                "paraview@5.11.2 %gcc@11.4.0 +python +mpi +osmesa +cuda",
                "visit@3.3.3 %gcc@11.4.0 +python +cuda +opengl",
                "vtk@9.2.6 %gcc@11.4.0 +python +cuda +opengl2",

                # GPU computing
                "cuda@11.8.0 %gcc@11.4.0",
                "cudnn@8.9.3.28-11.8 %gcc@11.4.0",

                # GPU-accelerated libraries
                "py-cupy@12.2.0 %gcc@11.4.0",
                "py-numba@0.57.1 %gcc@11.4.0",
                "openvdb@10.0.1 %gcc@11.4.0 +cuda",

                # Ray tracing and rendering
                "ospray@2.12.0 %gcc@11.4.0 +openimageio",
                "embree@4.1.0 %gcc@11.4.0",
                "optix@7.7.0 %gcc@11.4.0",

                # Machine learning visualization
                "py-tensorboard@2.13.0 %gcc@11.4.0",
                "py-wandb@0.15.8 %gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "workstation": {
                    "instance_type": "g5.12xlarge",
                    "vcpus": 48,
                    "memory_gb": 192,
                    "gpu": "4x NVIDIA A10G (24GB each)",
                    "storage": "3.8TB NVMe SSD",
                    "cost_per_hour": 5.672,
                    "use_case": "Professional visualization workstation"
                },
                "render_farm": {
                    "instance_type": "p4d.24xlarge",
                    "vcpus": 96,
                    "memory_gb": 1152,
                    "gpu": "8x NVIDIA A100 (40GB each)",
                    "cost_per_hour": 32.772,
                    "use_case": "Large-scale rendering and AI visualization"
                }
            },
            "optimization_features": [
                "GPU memory optimization",
                "Multi-GPU parallel rendering",
                "CUDA-accelerated algorithms",
                "Real-time ray tracing",
                "AI-enhanced visualization",
                "Distributed GPU clusters"
            ]
        }

    def _get_medical_imaging_config(self) -> Dict[str, Any]:
        """Specialized medical imaging visualization"""
        return {
            "name": "Medical Imaging Visualization Suite",
            "description": "DICOM processing and medical image visualization",
            "spack_packages": [
                # Medical imaging tools
                "itk@5.3.0 %gcc@11.4.0 +python +fftw",
                "vtk@9.2.6 %gcc@11.4.0 +python +opengl2 +imaging",
                "py-simpleitk@2.2.1 %gcc@11.4.0",
                "py-pydicom@2.4.2 %gcc@11.4.0",

                # Image processing
                "py-scikit-image@0.21.0 %gcc@11.4.0",
                "py-opencv@4.8.0 %gcc@11.4.0",
                "py-nibabel@5.1.0 %gcc@11.4.0",

                # 3D visualization
                "paraview@5.11.2 %gcc@11.4.0 +python +medical",
                "py-mayavi@4.8.1 %gcc@11.4.0",
                "py-pyvista@0.42.2 %gcc@11.4.0"
            ],
            "compliance_features": [
                "HIPAA compliance controls",
                "PHI data encryption",
                "Audit logging",
                "Access controls and authentication",
                "Data anonymization tools"
            ]
        }

    def _get_geospatial_config(self) -> Dict[str, Any]:
        """Geospatial data visualization"""
        return {
            "name": "Geospatial Visualization Platform",
            "description": "GIS and earth science data visualization",
            "spack_packages": [
                # GIS tools
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql",
                "gdal@3.7.1 %gcc@11.4.0 +python +postgresql +netcdf",
                "proj@9.2.1 %gcc@11.4.0",
                "geos@3.12.0 %gcc@11.4.0",

                # Python geospatial
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",

                # Climate and earth science
                "ncview@2.1.8 %gcc@11.4.0",
                "nco@5.1.6 %gcc@11.4.0",
                "cdo@2.2.0 %gcc@11.4.0"
            ]
        }

    def _get_molecular_config(self) -> Dict[str, Any]:
        """Molecular and structural biology visualization"""
        return {
            "name": "Molecular Visualization Suite",
            "description": "Protein and molecular structure visualization",
            "spack_packages": [
                # Molecular visualization
                "pymol@2.5.0 %gcc@11.4.0 +python",
                "vmd@1.9.3 %gcc@11.4.0 +python +cuda",
                "chimera@1.16 %gcc@11.4.0",

                # Structure analysis
                "py-mdanalysis@2.5.0 %gcc@11.4.0",
                "py-nglview@3.0.8 %gcc@11.4.0",
                "py-py3dmol@2.0.4 %gcc@11.4.0"
            ]
        }

    def generate_deployment_config(self, config_name: str,
                                 team_size: int = 5,
                                 data_size_gb: int = 1000,
                                 gpu_required: bool = True) -> Dict[str, Any]:
        """Generate AWS deployment configuration for visualization environment"""

        if config_name not in self.visualization_domains:
            raise ValueError(f"Unknown configuration: {config_name}")

        config = self.visualization_domains[config_name]

        # Determine instance recommendations based on team size and requirements
        if team_size <= 3:
            scale = "small_team"
        elif team_size <= 15:
            scale = "research_group"
        else:
            scale = "department"

        deployment = {
            "configuration": config,
            "team_size": team_size,
            "data_size_gb": data_size_gb,
            "estimated_monthly_cost": self._calculate_cost(config, team_size, data_size_gb),
            "deployment_timeline": "2-4 hours automated setup",
            "support_features": [
                "Auto-scaling based on usage",
                "Automated backups to S3",
                "CloudWatch monitoring",
                "SNS alerting",
                "IAM role-based access"
            ]
        }

        return deployment

    def _calculate_cost(self, config: Dict[str, Any], team_size: int, data_size_gb: int) -> Dict[str, float]:
        """Calculate estimated monthly costs"""
        # Simplified cost calculation - would be more sophisticated in production
        base_cost = 100  # Base infrastructure
        user_cost = team_size * 50  # Per user cost
        storage_cost = data_size_gb * 0.023  # S3 storage cost

        return {
            "base_infrastructure": base_cost,
            "user_scaling": user_cost,
            "storage": storage_cost,
            "total_estimated": base_cost + user_cost + storage_cost
        }

    def list_available_configurations(self) -> List[str]:
        """List all available visualization configurations"""
        return list(self.visualization_domains.keys())

    def get_configuration_summary(self, config_name: str) -> str:
        """Get a summary of a specific configuration"""
        if config_name not in self.visualization_domains:
            return f"Configuration '{config_name}' not found"

        config = self.visualization_domains[config_name]
        return f"{config['name']}: {config['description']}"

def main():
    """CLI interface for visualization pack"""
    import argparse

    parser = argparse.ArgumentParser(description="AWS Research Wizard - Visualization Studio Pack")
    parser.add_argument("--list", action="store_true", help="List available configurations")
    parser.add_argument("--config", type=str, help="Generate configuration")
    parser.add_argument("--team-size", type=int, default=5, help="Team size")
    parser.add_argument("--data-gb", type=int, default=1000, help="Data size in GB")
    parser.add_argument("--output", type=str, help="Output file for configuration")

    args = parser.parse_args()

    viz_pack = VisualizationStudioPack()

    if args.list:
        print("Available Visualization Configurations:")
        for config_name in viz_pack.list_available_configurations():
            print(f"  {config_name}: {viz_pack.get_configuration_summary(config_name)}")

    elif args.config:
        try:
            deployment = viz_pack.generate_deployment_config(
                args.config,
                args.team_size,
                args.data_gb
            )

            if args.output:
                with open(args.output, 'w') as f:
                    json.dump(deployment, f, indent=2)
                print(f"Configuration saved to {args.output}")
            else:
                print(json.dumps(deployment, indent=2))

        except ValueError as e:
            print(f"Error: {e}")

    else:
        parser.print_help()

if __name__ == "__main__":
    main()
