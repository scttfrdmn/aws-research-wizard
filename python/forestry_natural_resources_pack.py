#!/usr/bin/env python3
"""
Forestry & Natural Resources Research Pack
Comprehensive forest management, ecology, and natural resource modeling for AWS Research Wizard
"""

import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from enum import Enum

class ForestryDomain(Enum):
    FOREST_INVENTORY = "forest_inventory"
    FOREST_GROWTH_MODELING = "forest_growth_modeling"
    FOREST_ECOLOGY = "forest_ecology"
    WILDLIFE_MANAGEMENT = "wildlife_management"
    FIRE_MANAGEMENT = "fire_management"
    CARBON_SEQUESTRATION = "carbon_sequestration"
    WATERSHED_MANAGEMENT = "watershed_management"
    FOREST_ECONOMICS = "forest_economics"
    REMOTE_SENSING_FORESTRY = "remote_sensing_forestry"

@dataclass
class ForestryWorkload:
    """Forestry and natural resources workload characteristics"""
    domain: ForestryDomain
    forest_type: str         # Temperate, Tropical, Boreal, Mixed, Plantation
    management_scale: str    # Stand, Forest, Landscape, Regional, National
    analysis_type: str       # Inventory, Growth Prediction, Ecosystem Services, Conservation
    temporal_scale: str      # Real-time, Annual, Decadal, Long-term (>50 years)
    data_sources: List[str]  # LiDAR, Satellite, Field, Drone, Climate, Socioeconomic
    modeling_approach: str   # Empirical, Process-based, Machine Learning, Hybrid
    data_volume_tb: float    # Expected data volume
    computational_intensity: str  # Light, Moderate, Intensive, Extreme

class ForestryNaturalResourcesPack:
    """
    Comprehensive forestry and natural resources research environments optimized for AWS
    Supports forest management, ecology, remote sensing, and ecosystem modeling
    """
    
    def __init__(self):
        self.forestry_configurations = {
            "forest_inventory_system": self._get_forest_inventory_config(),
            "forest_growth_modeling": self._get_forest_growth_config(),
            "forest_ecology_platform": self._get_forest_ecology_config(),
            "wildlife_management": self._get_wildlife_management_config(),
            "fire_management_system": self._get_fire_management_config(),
            "carbon_sequestration": self._get_carbon_sequestration_config(),
            "watershed_management": self._get_watershed_management_config(),
            "forest_economics": self._get_forest_economics_config(),
            "remote_sensing_forestry": self._get_remote_sensing_config()
        }
    
    def _get_forest_inventory_config(self) -> Dict[str, Any]:
        """Forest inventory and mensuration platform"""
        return {
            "name": "Forest Inventory & Mensuration Platform",
            "description": "Comprehensive forest inventory, biomass estimation, and mensuration analysis",
            "spack_packages": [
                # Forest inventory software
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-forestinventory@1.0.0 %gcc@11.4.0",     # Forest inventory analysis
                "r-raster@3.6-23 %gcc@11.4.0",             # Raster data processing
                "r-sp@2.0-0 %gcc@11.4.0",                  # Spatial data classes
                "r-rgdal@1.6-7 %gcc@11.4.0",               # Geospatial data abstraction
                
                # LiDAR processing
                "pdal@2.5.6 %gcc@11.4.0 +python +hdf5 +laszip", # Point cloud processing
                "liblas@1.8.1 %gcc@11.4.0",                # LAS file format support
                "fusion@4.40 %gcc@11.4.0",                 # LiDAR data analysis
                "cloudcompare@2.13.1 %gcc@11.4.0 +qt",     # Point cloud processing
                
                # Python forestry tools
                "python@3.11.5 %gcc@11.4.0",
                "py-laspy@2.5.1 %gcc@11.4.0",              # LAS file processing
                "py-open3d@0.17.0 %gcc@11.4.0",            # 3D data processing
                "py-forestry@1.8.0 %gcc@11.4.0",           # Forest analysis tools
                "py-treemetrics@2.3.0 %gcc@11.4.0",        # Tree measurement algorithms
                
                # Allometric and biomass estimation
                "py-allometry@1.5.0 %gcc@11.4.0",          # Allometric equations
                "py-biomass-estimation@2.1.0 %gcc@11.4.0", # Biomass calculation
                "py-tree-segmentation@1.2.0 %gcc@11.4.0",  # Individual tree detection
                
                # Geospatial analysis
                "gdal@3.7.2 %gcc@11.4.0 +python +netcdf +hdf5",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                "py-shapely@2.0.1 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Machine learning for forest inventory
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-pyvista@0.42.0 %gcc@11.4.0",           # 3D visualization
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +postgis",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "development": {
                    "instance_type": "c6i.xlarge",
                    "vcpus": 4,
                    "memory_gb": 8,
                    "storage_gb": 100,
                    "cost_per_hour": 0.17,
                    "use_case": "Development and small plot analysis"
                },
                "stand_inventory": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 500,
                    "cost_per_hour": 0.51,
                    "use_case": "Stand-level forest inventory and LiDAR processing"
                },
                "forest_inventory": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 1000,
                    "cost_per_hour": 1.02,
                    "use_case": "Forest-level inventory and biomass estimation"
                },
                "landscape_analysis": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Landscape-scale inventory and carbon assessment"
                }
            },
            "estimated_cost": {
                "compute": 600,
                "storage": 200,
                "data_transfer": 100,
                "total": 900
            },
            "research_capabilities": [
                "LiDAR point cloud processing and analysis",
                "Individual tree detection and segmentation",
                "Forest biomass and carbon stock estimation",
                "Allometric equation development and validation",
                "Forest inventory statistical analysis",
                "Integration of field and remote sensing data",
                "Height, diameter, and volume estimation",
                "Forest structure characterization"
            ],
            "aws_data_sources": [
                "NASA GEDI (Global Ecosystem Dynamics Investigation)",
                "USGS National Map LiDAR",
                "Landsat and Sentinel satellite imagery",
                "Forest Inventory and Analysis (FIA) data",
                "Global Forest Watch datasets"
            ]
        }
    
    def _get_forest_growth_config(self) -> Dict[str, Any]:
        """Forest growth modeling and yield prediction platform"""
        return {
            "name": "Forest Growth Modeling & Yield Prediction Platform",
            "description": "Individual tree and stand growth modeling, yield prediction, and forest dynamics",
            "spack_packages": [
                # Forest growth models
                "fvs@2023.1 %gcc@11.4.0 +fortran",         # Forest Vegetation Simulator
                "organon@9.2 %gcc@11.4.0 +fortran",        # Oregon growth model
                "silva@3.0.8 %gcc@11.4.0 +fortran",        # European forest growth
                "3pg@2.7 %gcc@11.4.0 +fortran",            # Physiological Principles Predicting Growth
                
                # R forest modeling packages
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-fgm@1.2.0 %gcc@11.4.0",                 # Forest growth modeling
                "r-sitree@0.1-19 %gcc@11.4.0",             # Individual tree growth
                "r-treering@1.0.2 %gcc@11.4.0",            # Tree ring analysis
                "r-dplr@1.7.6 %gcc@11.4.0",                # Dendrochronology program library
                
                # Python forest dynamics
                "python@3.11.5 %gcc@11.4.0",
                "py-forest-dynamics@2.5.0 %gcc@11.4.0",    # Forest dynamics modeling
                "py-tree-growth@1.8.0 %gcc@11.4.0",        # Tree growth algorithms
                "py-yield-prediction@2.2.0 %gcc@11.4.0",   # Yield prediction models
                "py-gap-models@1.4.0 %gcc@11.4.0",         # Gap-based forest models
                
                # Climate and environmental data
                "py-netcdf4@1.6.4 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-cftime@1.6.2 %gcc@11.4.0",
                "py-climate-data@1.3.0 %gcc@11.4.0",
                
                # Machine learning for growth prediction
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Optimization for forest management
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                "py-pyomo@6.6.1 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                
                # AWS-optimized parallel computing with EFA support
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers +pmix +pmi +fabrics",
                "libfabric@1.18.1 %gcc@11.4.0 +verbs +mlx +efa",  # EFA support
                "aws-ofi-nccl@1.7.0 %gcc@11.4.0",  # AWS OFI plugin
                "ucx@1.14.1 %gcc@11.4.0 +verbs +mlx +ib_hw_tm",  # Unified Communication X
                "py-mpi4py@3.1.4 %gcc@11.4.0",
                "py-dask@2023.7.1 %gcc@11.4.0",
                "slurm@23.02.5 %gcc@11.4.0 +pmix +numa",  # Cluster management
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "gfortran@11.4.0"
            ],
            "aws_instance_recommendations": {
                "individual_tree": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Individual tree growth modeling"
                },
                "stand_modeling": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Stand-level growth and yield modeling"
                },
                "landscape_simulation": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Landscape-scale forest dynamics simulation"
                }
            },
            "estimated_cost": {
                "compute": 500,
                "storage": 150,
                "data_transfer": 75,
                "total": 725
            },
            "research_capabilities": [
                "Individual tree and stand growth modeling",
                "Forest yield prediction and optimization",
                "Climate change impact on forest growth",
                "Silvicultural treatment effect modeling",
                "Forest dynamics and succession modeling",
                "Integration of process-based and empirical models",
                "Long-term forest productivity assessment",
                "Dendrochronology and tree ring analysis"
            ]
        }
    
    def _get_forest_ecology_config(self) -> Dict[str, Any]:
        """Forest ecology and ecosystem dynamics platform"""
        return {
            "name": "Forest Ecology & Ecosystem Dynamics Platform",
            "description": "Biodiversity analysis, ecosystem services, and ecological modeling",
            "spack_packages": [
                # Ecological modeling frameworks
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-vegan@2.6-4 %gcc@11.4.0",               # Community ecology
                "r-bipartite@2.19 %gcc@11.4.0",            # Network analysis
                "r-ade4@1.7-22 %gcc@11.4.0",               # Multivariate analysis
                "r-picante@1.8.2 %gcc@11.4.0",             # Phylogenetic analysis
                
                # Ecosystem service modeling
                "invest@3.14.0 %gcc@11.4.0 +python",       # InVEST ecosystem services
                "aries@1.8.0 %gcc@11.4.0 +python",         # ARIES ecosystem services
                "solves@4.4.0 %gcc@11.4.0 +python",        # Ecosystem service mapping
                
                # Python ecological tools
                "python@3.11.5 %gcc@11.4.0",
                "py-scikit-bio@0.5.8 %gcc@11.4.0",         # Bioinformatics
                "py-ecology@2.3.0 %gcc@11.4.0",            # Ecological analysis
                "py-biodiversity@1.7.0 %gcc@11.4.0",       # Biodiversity metrics
                "py-ecosystem-services@2.1.0 %gcc@11.4.0", # Ecosystem service valuation
                
                # Species distribution modeling
                "maxent@3.4.4 %gcc@11.4.0 +java",          # Maximum entropy modeling
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-species-distribution@1.5.0 %gcc@11.4.0",
                
                # Spatial analysis
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                
                # Network analysis for ecology
                "py-networkx@3.1 %gcc@11.4.0",
                "py-igraph@0.10.6 %gcc@11.4.0",
                "py-graph-tool@2.45 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +postgis",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "biodiversity_analysis": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 300,
                    "cost_per_hour": 0.51,
                    "use_case": "Biodiversity and community ecology analysis"
                },
                "ecosystem_services": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Ecosystem service modeling and valuation"
                },
                "landscape_ecology": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Landscape-scale ecological modeling"
                }
            },
            "estimated_cost": {
                "compute": 550,
                "storage": 120,
                "data_transfer": 60,
                "total": 730
            },
            "research_capabilities": [
                "Biodiversity assessment and monitoring",
                "Ecosystem service quantification and mapping",
                "Species distribution modeling",
                "Food web and ecological network analysis",
                "Habitat connectivity and fragmentation analysis",
                "Conservation planning and prioritization",
                "Climate change impact on ecosystems",
                "Restoration effectiveness assessment"
            ]
        }
    
    def _get_wildlife_management_config(self) -> Dict[str, Any]:
        """Wildlife management and conservation platform"""
        return {
            "name": "Wildlife Management & Conservation Platform",
            "description": "Wildlife population modeling, habitat analysis, and conservation planning",
            "spack_packages": [
                # Wildlife population modeling
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-rmark@3.0.0 %gcc@11.4.0",               # Mark-recapture analysis
                "r-unmarked@1.4.1 %gcc@11.4.0",            # Hierarchical models
                "r-distance@1.0.8 %gcc@11.4.0",            # Distance sampling
                "r-secr@4.6.9 %gcc@11.4.0",                # Spatially explicit capture-recapture
                
                # Population viability analysis
                "vortex@10.5.5 %gcc@11.4.0",               # Population viability analysis
                "ramas@6.0 %gcc@11.4.0",                   # Risk assessment
                "py-pva@2.1.0 %gcc@11.4.0",                # Population viability in Python
                
                # Python wildlife tools
                "python@3.11.5 %gcc@11.4.0",
                "py-wildlife-analysis@1.9.0 %gcc@11.4.0",   # Wildlife data analysis
                "py-movement-ecology@2.4.0 %gcc@11.4.0",    # Animal movement analysis
                "py-habitat-modeling@1.6.0 %gcc@11.4.0",    # Habitat suitability modeling
                "py-conservation-genetics@1.3.0 %gcc@11.4.0", # Conservation genetics
                
                # Remote sensing for wildlife
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-earthengine-api@0.1.364 %gcc@11.4.0",
                "py-sentinelsat@1.2.1 %gcc@11.4.0",
                
                # Machine learning for wildlife
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                
                # Spatial analysis
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-shapely@2.0.1 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +postgis",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "population_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Wildlife population modeling and analysis"
                },
                "habitat_analysis": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Habitat suitability modeling and landscape analysis"
                },
                "conservation_planning": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Large-scale conservation planning and optimization"
                }
            },
            "estimated_cost": {
                "compute": 450,
                "storage": 100,
                "data_transfer": 50,
                "total": 600
            },
            "research_capabilities": [
                "Wildlife population estimation and monitoring",
                "Habitat suitability and connectivity modeling",
                "Conservation planning and reserve design",
                "Species distribution and range modeling",
                "Population viability analysis",
                "Movement ecology and migration analysis",
                "Human-wildlife conflict analysis",
                "Conservation genetics and population structure"
            ]
        }
    
    def _get_fire_management_config(self) -> Dict[str, Any]:
        """Fire management and wildfire modeling platform"""
        return {
            "name": "Fire Management & Wildfire Modeling Platform",
            "description": "Wildfire behavior modeling, risk assessment, and fire management planning",
            "spack_packages": [
                # Fire behavior models
                "farsite@4.1.055 %gcc@11.4.0 +fortran",    # Fire Area Simulator
                "flammap@6.2.0 %gcc@11.4.0 +fortran",      # Fire mapping and analysis
                "wfds@6.7.7 %gcc@11.4.0 +fortran +mpi",    # Wildland Fire Dynamics Simulator
                "fofem@6.8 %gcc@11.4.0 +fortran",          # First Order Fire Effects Model
                
                # Weather and climate data
                "py-netcdf4@1.6.4 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-metpy@1.5.0 %gcc@11.4.0",
                
                # Python fire modeling
                "python@3.11.5 %gcc@11.4.0",
                "py-fire-behavior@2.7.0 %gcc@11.4.0",      # Fire behavior modeling
                "py-wildfire-risk@1.9.0 %gcc@11.4.0",      # Wildfire risk assessment
                "py-fire-effects@1.5.0 %gcc@11.4.0",       # Fire effects modeling
                "py-fuel-models@2.2.0 %gcc@11.4.0",        # Fuel load modeling
                
                # Remote sensing for fire
                "py-earthengine-api@0.1.364 %gcc@11.4.0",
                "py-sentinelsat@1.2.1 %gcc@11.4.0",
                "py-modis@0.7.3 %gcc@11.4.0",              # MODIS fire products
                "py-viirs@1.2.0 %gcc@11.4.0",              # VIIRS fire detection
                
                # Machine learning for fire prediction
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",
                
                # Spatial analysis
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-shapely@2.0.1 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Optimization for fire management
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                "py-pyomo@6.6.1 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                
                # Parallel computing
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers",
                "py-mpi4py@3.1.4 %gcc@11.4.0",
                "py-dask@2023.7.1 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +postgis",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "gfortran@11.4.0"
            ],
            "aws_instance_recommendations": {
                "fire_modeling": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage_gb": 500,
                    "cost_per_hour": 0.68,
                    "use_case": "Fire behavior modeling and simulation"
                },
                "risk_assessment": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 1000,
                    "cost_per_hour": 1.02,
                    "use_case": "Wildfire risk assessment and mapping"
                },
                "real_time_monitoring": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Real-time fire monitoring and decision support"
                }
            },
            "estimated_cost": {
                "compute": 800,
                "storage": 300,
                "data_transfer": 150,
                "total": 1250
            },
            "research_capabilities": [
                "Wildfire behavior modeling and simulation",
                "Fire risk assessment and mapping",
                "Fuel load and moisture modeling",
                "Fire weather analysis",
                "Prescribed burn planning",
                "Fire effects and ecological impact assessment",
                "Real-time fire monitoring and detection",
                "Fire management decision support systems"
            ]
        }
    
    def _get_carbon_sequestration_config(self) -> Dict[str, Any]:
        """Carbon sequestration and forest carbon modeling platform"""
        return {
            "name": "Carbon Sequestration & Forest Carbon Modeling Platform",
            "description": "Forest carbon dynamics, sequestration assessment, and climate change mitigation",
            "spack_packages": [
                # Carbon cycle models
                "century@5.0 %gcc@11.4.0 +fortran",        # Soil organic matter model
                "casa@2.1 %gcc@11.4.0 +fortran",           # Carbon and nitrogen cycle
                "biome-bgc@4.2 %gcc@11.4.0 +fortran",      # Biogeochemical cycles
                "cbm-cfs3@1.5.8 %gcc@11.4.0 +fortran",     # Carbon Budget Model
                
                # Forest carbon tools
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-forest-carbon@2.1.0 %gcc@11.4.0",       # Forest carbon calculations
                "r-carbonfacts@1.3.0 %gcc@11.4.0",         # Carbon accounting
                
                # Python carbon modeling
                "python@3.11.5 %gcc@11.4.0",
                "py-forest-carbon@2.8.0 %gcc@11.4.0",      # Forest carbon modeling
                "py-carbon-dynamics@1.7.0 %gcc@11.4.0",    # Carbon cycle dynamics
                "py-biomass-carbon@1.4.0 %gcc@11.4.0",     # Biomass to carbon conversion
                "py-soil-carbon@2.3.0 %gcc@11.4.0",        # Soil carbon modeling
                
                # Climate data processing
                "py-netcdf4@1.6.4 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-cftime@1.6.2 %gcc@11.4.0",
                "py-climate-indices@1.0.11 %gcc@11.4.0",
                
                # Remote sensing for carbon
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-earthengine-api@0.1.364 %gcc@11.4.0",
                "py-sentinelsat@1.2.1 %gcc@11.4.0",
                
                # Machine learning for carbon prediction
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Uncertainty analysis
                "py-uncertainties@3.1.7 %gcc@11.4.0",
                "py-monte-carlo@2.0.5 %gcc@11.4.0",
                "py-sensitivity-analysis@1.8.0 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "gfortran@11.4.0"
            ],
            "aws_instance_recommendations": {
                "carbon_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 300,
                    "cost_per_hour": 0.34,
                    "use_case": "Forest carbon stock analysis and accounting"
                },
                "carbon_modeling": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Carbon dynamics modeling and projection"
                },
                "climate_scenarios": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Climate scenario analysis and carbon forecasting"
                }
            },
            "estimated_cost": {
                "compute": 500,
                "storage": 150,
                "data_transfer": 75,
                "total": 725
            },
            "research_capabilities": [
                "Forest carbon stock assessment and monitoring",
                "Carbon sequestration rate calculation",
                "Soil carbon dynamics modeling",
                "Climate change impact on forest carbon",
                "Carbon accounting and reporting (IPCC guidelines)",
                "REDD+ and carbon offset project development",
                "Uncertainty analysis in carbon estimates",
                "Carbon market and policy analysis"
            ]
        }
    
    def _get_watershed_management_config(self) -> Dict[str, Any]:
        """Watershed management and hydrology platform"""
        return {
            "name": "Watershed Management & Forest Hydrology Platform",
            "description": "Watershed modeling, water quality analysis, and forest hydrology",
            "spack_packages": [
                # Hydrological models
                "swat@2012.664 %gcc@11.4.0 +fortran",      # Soil and Water Assessment Tool
                "hspf@12.2 %gcc@11.4.0 +fortran",          # Hydrological Simulation Program
                "mike-she@2023.1 %gcc@11.4.0 +fortran",    # Integrated hydrological modeling
                "rhessys@7.4 %gcc@11.4.0 +fortran",        # Regional Hydro-Ecological Simulation System
                
                # Python hydrology tools
                "python@3.11.5 %gcc@11.4.0",
                "py-hydrology@2.6.0 %gcc@11.4.0",          # Hydrological analysis
                "py-watershed@1.8.0 %gcc@11.4.0",          # Watershed modeling
                "py-water-quality@2.2.0 %gcc@11.4.0",      # Water quality modeling
                "py-streamflow@1.5.0 %gcc@11.4.0",         # Streamflow analysis
                
                # GIS and spatial analysis
                "gdal@3.7.2 %gcc@11.4.0 +python +netcdf +hdf5",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                
                # Climate and weather data
                "py-netcdf4@1.6.4 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-metpy@1.5.0 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                
                # Optimization for watershed management
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +postgis",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "gfortran@11.4.0"
            ],
            "aws_instance_recommendations": {
                "watershed_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Small watershed analysis and modeling"
                },
                "basin_modeling": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "River basin and large watershed modeling"
                },
                "regional_hydrology": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Regional hydrological modeling and forecasting"
                }
            },
            "estimated_cost": {
                "compute": 450,
                "storage": 120,
                "data_transfer": 60,
                "total": 630
            },
            "research_capabilities": [
                "Watershed modeling and water balance analysis",
                "Streamflow prediction and forecasting",
                "Water quality assessment and modeling",
                "Forest management impact on hydrology",
                "Flood risk assessment and management",
                "Groundwater and surface water interactions",
                "Climate change impact on water resources",
                "Best management practice effectiveness"
            ]
        }
    
    def _get_forest_economics_config(self) -> Dict[str, Any]:
        """Forest economics and management optimization platform"""
        return {
            "name": "Forest Economics & Management Optimization Platform",
            "description": "Forest valuation, economic analysis, and management optimization",
            "spack_packages": [
                # Economic optimization
                "gams@44.3.0 %gcc@11.4.0",                 # General Algebraic Modeling System
                "lindo@14.0 %gcc@11.4.0",                  # Linear optimization
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-forest-economics@1.5.0 %gcc@11.4.0",   # Forest economic analysis
                
                # Python economics tools
                "python@3.11.5 %gcc@11.4.0",
                "py-forest-economics@2.4.0 %gcc@11.4.0",   # Forest economic modeling
                "py-forestry-optimization@1.7.0 %gcc@11.4.0", # Forest management optimization
                "py-timber-valuation@1.3.0 %gcc@11.4.0",   # Timber valuation methods
                "py-ecosystem-valuation@2.1.0 %gcc@11.4.0", # Ecosystem service valuation
                
                # Optimization frameworks
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                "py-pyomo@6.6.1 %gcc@11.4.0",
                "py-ortools@9.7.2996 %gcc@11.4.0",
                
                # Statistical and econometric analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                
                # Financial analysis
                "py-finance@1.4.0 %gcc@11.4.0",
                "py-numpy-financial@1.0.0 %gcc@11.4.0",
                "py-quantlib@1.31 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "economic_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Forest economic analysis and valuation"
                },
                "optimization": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 300,
                    "cost_per_hour": 1.02,
                    "use_case": "Forest management optimization and planning"
                }
            },
            "estimated_cost": {
                "compute": 300,
                "storage": 60,
                "data_transfer": 30,
                "total": 390
            },
            "research_capabilities": [
                "Forest valuation and asset assessment",
                "Optimal rotation age and thinning schedules",
                "Forest investment and portfolio analysis",
                "Ecosystem service economic valuation",
                "Timber market analysis and forecasting",
                "Risk assessment and insurance modeling",
                "Policy impact analysis",
                "Sustainable forest management optimization"
            ]
        }
    
    def _get_remote_sensing_config(self) -> Dict[str, Any]:
        """Remote sensing for forestry platform"""
        return {
            "name": "Remote Sensing for Forestry Platform",
            "description": "Satellite and LiDAR analysis for forest monitoring and assessment",
            "spack_packages": [
                # Remote sensing software
                "gdal@3.7.2 %gcc@11.4.0 +python +netcdf +hdf5",
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql",
                "grass@8.3.0 %gcc@11.4.0 +netcdf +postgresql",
                "saga@9.1.1 %gcc@11.4.0 +python",
                "orfeo-toolbox@9.0.0 %gcc@11.4.0 +python",
                
                # LiDAR processing
                "pdal@2.5.6 %gcc@11.4.0 +python +hdf5 +laszip",
                "liblas@1.8.1 %gcc@11.4.0",
                "cloudcompare@2.13.1 %gcc@11.4.0 +qt",
                
                # Python remote sensing
                "python@3.11.5 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-earthengine-api@0.1.364 %gcc@11.4.0",
                "py-sentinelsat@1.2.1 %gcc@11.4.0",
                "py-planetary-computer@0.4.9 %gcc@11.4.0",
                
                # Image processing
                "py-scikit-image@0.21.0 %gcc@11.4.0",
                "py-opencv@4.8.0 %gcc@11.4.0",
                "py-pillow@10.0.0 %gcc@11.4.0",
                
                # Machine learning for remote sensing
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                
                # Geospatial analysis
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-shapely@2.0.1 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                
                # Data processing
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-dask@2023.7.1 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-pyvista@0.42.0 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +postgis",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "image_processing": {
                    "instance_type": "g4dn.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 32,
                    "storage_gb": 500,
                    "cost_per_hour": 0.752,
                    "use_case": "Satellite image processing and analysis"
                },
                "lidar_processing": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 1000,
                    "cost_per_hour": 1.02,
                    "use_case": "LiDAR point cloud processing and forest structure analysis"
                },
                "large_scale_analysis": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Large-scale forest monitoring and change detection"
                }
            },
            "estimated_cost": {
                "compute": 700,
                "storage": 250,
                "data_transfer": 150,
                "total": 1100
            },
            "research_capabilities": [
                "Forest cover mapping and classification",
                "Change detection and deforestation monitoring",
                "Forest structure analysis from LiDAR",
                "Individual tree detection and crown delineation",
                "Forest health and stress assessment",
                "Biomass estimation from remote sensing",
                "Time series analysis of forest dynamics",
                "Integration of multi-sensor data"
            ]
        }
    
    def generate_forestry_recommendation(self, workload: ForestryWorkload) -> Dict[str, Any]:
        """Generate optimized AWS infrastructure recommendation for forestry research"""
        
        # Select appropriate configuration based on domain
        domain_config_map = {
            ForestryDomain.FOREST_INVENTORY: "forest_inventory_system",
            ForestryDomain.FOREST_GROWTH_MODELING: "forest_growth_modeling",
            ForestryDomain.FOREST_ECOLOGY: "forest_ecology_platform",
            ForestryDomain.WILDLIFE_MANAGEMENT: "wildlife_management",
            ForestryDomain.FIRE_MANAGEMENT: "fire_management_system",
            ForestryDomain.CARBON_SEQUESTRATION: "carbon_sequestration",
            ForestryDomain.WATERSHED_MANAGEMENT: "watershed_management",
            ForestryDomain.FOREST_ECONOMICS: "forest_economics",
            ForestryDomain.REMOTE_SENSING_FORESTRY: "remote_sensing_forestry"
        }
        
        config_name = domain_config_map.get(workload.domain, "forest_inventory_system")
        base_config = self.forestry_configurations[config_name].copy()
        
        # Adjust configuration based on workload characteristics
        self._optimize_for_scale(base_config, workload)
        self._optimize_for_data_volume(base_config, workload)
        self._optimize_for_computational_intensity(base_config, workload)
        
        # Generate cost estimates
        base_config["estimated_cost"] = self._calculate_forestry_costs(workload, base_config)
        
        # Add optimization recommendations
        base_config["optimization_recommendations"] = self._generate_optimization_recommendations(workload)
        
        return {
            "configuration": base_config,
            "workload_analysis": {
                "domain": workload.domain.value,
                "management_scale": workload.management_scale,
                "analysis_type": workload.analysis_type,
                "computational_requirements": workload.computational_intensity,
                "data_volume": f"{workload.data_volume_tb} TB"
            },
            "deployment_recommendations": self._generate_deployment_recommendations(workload),
            "estimated_cost": base_config["estimated_cost"]
        }
    
    def _optimize_for_scale(self, config: Dict[str, Any], workload: ForestryWorkload):
        """Optimize configuration based on management scale"""
        scale_multipliers = {
            "Stand": 1.0,
            "Forest": 2.0,
            "Landscape": 4.0,
            "Regional": 8.0,
            "National": 16.0
        }
        
        multiplier = scale_multipliers.get(workload.management_scale, 1.0)
        
        # Adjust instance recommendations based on scale
        if "aws_instance_recommendations" in config:
            for instance_config in config["aws_instance_recommendations"].values():
                if multiplier > 4.0:
                    # Scale up for large management areas
                    instance_config["storage_gb"] = int(instance_config["storage_gb"] * multiplier)
    
    def _optimize_for_data_volume(self, config: Dict[str, Any], workload: ForestryWorkload):
        """Optimize configuration based on expected data volume"""
        if workload.data_volume_tb > 5.0:
            # Add data processing optimizations for large datasets
            if "spack_packages" in config:
                config["spack_packages"].extend([
                    "py-dask@2023.7.1 %gcc@11.4.0",
                    "py-ray@2.6.1 %gcc@11.4.0"
                ])
    
    def _optimize_for_computational_intensity(self, config: Dict[str, Any], workload: ForestryWorkload):
        """Optimize configuration based on computational intensity"""
        if workload.computational_intensity in ["Intensive", "Extreme"]:
            # Upgrade to compute-optimized or GPU instances for intensive workloads
            if "aws_instance_recommendations" in config:
                for key, instance_config in config["aws_instance_recommendations"].items():
                    if workload.computational_intensity == "Extreme":
                        if "LiDAR" in workload.data_sources or "Satellite" in workload.data_sources:
                            instance_config["instance_type"] = "g5.2xlarge"
                            instance_config["cost_per_hour"] = 1.624
    
    def _calculate_forestry_costs(self, workload: ForestryWorkload, config: Dict[str, Any]) -> Dict[str, float]:
        """Calculate estimated costs for forestry research infrastructure"""
        base_compute = 500
        base_storage = 150
        base_data_transfer = 75
        
        # Scale costs based on management scale
        scale_multipliers = {"Stand": 1.0, "Forest": 1.5, "Landscape": 3.0, "Regional": 6.0, "National": 12.0}
        multiplier = scale_multipliers.get(workload.management_scale, 1.0)
        
        # Adjust for computational intensity
        intensity_multipliers = {"Light": 0.5, "Moderate": 1.0, "Intensive": 2.0, "Extreme": 4.0}
        intensity_mult = intensity_multipliers.get(workload.computational_intensity, 1.0)
        
        compute_cost = base_compute * multiplier * intensity_mult
        storage_cost = base_storage * (1 + workload.data_volume_tb / 5.0)
        data_transfer_cost = base_data_transfer * multiplier
        
        return {
            "compute": compute_cost,
            "storage": storage_cost,
            "data_transfer": data_transfer_cost,
            "total": compute_cost + storage_cost + data_transfer_cost
        }
    
    def _generate_optimization_recommendations(self, workload: ForestryWorkload) -> List[str]:
        """Generate optimization recommendations for forestry workloads"""
        recommendations = []
        
        if workload.management_scale in ["Landscape", "Regional", "National"]:
            recommendations.append("Consider using Spot Instances for batch processing to reduce costs by 60-90%")
            recommendations.append("Implement auto-scaling for seasonal analysis workflows")
        
        if workload.data_volume_tb > 10.0:
            recommendations.append("Use S3 Intelligent Tiering for automatic cost optimization of large remote sensing datasets")
            recommendations.append("Consider AWS Batch for parallel processing of LiDAR and satellite imagery")
        
        if "Real-time" in workload.temporal_scale:
            recommendations.append("Use AWS IoT Core for real-time forest sensor data ingestion")
            recommendations.append("Consider Amazon Kinesis for streaming forest monitoring data")
        
        if workload.computational_intensity == "Extreme":
            recommendations.append("Use GPU instances for deep learning and computer vision tasks in remote sensing")
            recommendations.append("Consider AWS ParallelCluster for HPC forestry modeling workloads")
        
        if "Fire" in workload.domain.value:
            recommendations.append("Implement real-time data pipelines for fire weather and satellite detection")
            recommendations.append("Use AWS Lambda for automated fire alert processing")
        
        return recommendations
    
    def _generate_deployment_recommendations(self, workload: ForestryWorkload) -> Dict[str, Any]:
        """Generate deployment recommendations for forestry research"""
        return {
            "deployment_strategy": "multi-tier" if workload.management_scale in ["Regional", "National"] else "single-tier",
            "backup_strategy": "automated_daily_snapshots",
            "monitoring": ["CloudWatch for infrastructure", "Custom dashboards for forest metrics"],
            "security": ["VPC with private subnets", "IAM roles for service access", "Data encryption at rest and in transit"],
            "disaster_recovery": "cross-region backup" if workload.management_scale == "National" else "single-region backup"
        }

if __name__ == "__main__":
    # Example usage
    pack = ForestryNaturalResourcesPack()
    
    # Example workload
    workload = ForestryWorkload(
        domain=ForestryDomain.FOREST_INVENTORY,
        forest_type="Temperate",
        management_scale="Forest",
        analysis_type="Inventory",
        temporal_scale="Annual",
        data_sources=["LiDAR", "Satellite"],
        modeling_approach="Machine Learning",
        data_volume_tb=3.0,
        computational_intensity="Moderate"
    )
    
    recommendation = pack.generate_forestry_recommendation(workload)
    print(json.dumps(recommendation, indent=2))