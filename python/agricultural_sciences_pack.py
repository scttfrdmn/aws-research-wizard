#!/usr/bin/env python3
"""
Agricultural Sciences Research Pack
Comprehensive agricultural research, precision agriculture, and food systems modeling for AWS Research Wizard
"""

import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from enum import Enum

class AgriculturalDomain(Enum):
    CROP_MODELING = "crop_modeling"
    PRECISION_AGRICULTURE = "precision_agriculture"
    SOIL_SCIENCE = "soil_science"
    PLANT_BREEDING = "plant_breeding"
    AGRICULTURAL_ECONOMICS = "agricultural_economics"
    LIVESTOCK_SYSTEMS = "livestock_systems"
    IRRIGATION_MANAGEMENT = "irrigation_management"
    PEST_DISEASE_MANAGEMENT = "pest_disease_management"
    AGRICULTURAL_GENOMICS = "agricultural_genomics"

@dataclass
class AgriculturalWorkload:
    """Agricultural research workload characteristics"""
    domain: AgriculturalDomain
    farm_scale: str          # Field, Farm, Regional, National, Global
    crop_types: List[str]    # Corn, Wheat, Rice, Soybean, Vegetables, etc.
    analysis_type: str       # Yield Prediction, Resource Optimization, Risk Assessment
    temporal_scale: str      # Real-time, Seasonal, Annual, Multi-year, Climate
    data_sources: List[str]  # Satellite, Drone, Sensor, Weather, Market, Genomic
    modeling_approach: str   # Mechanistic, Statistical, Machine Learning, Hybrid
    data_volume_tb: float    # Expected data volume
    computational_intensity: str  # Light, Moderate, Intensive, Extreme

class AgriculturalSciencesPack:
    """
    Comprehensive agricultural sciences research environments optimized for AWS
    Supports crop modeling, precision agriculture, soil science, and agricultural genomics
    """
    
    def __init__(self):
        self.agricultural_configurations = {
            "crop_modeling_platform": self._get_crop_modeling_config(),
            "precision_agriculture": self._get_precision_agriculture_config(),
            "soil_science_laboratory": self._get_soil_science_config(),
            "plant_breeding_genomics": self._get_plant_breeding_config(),
            "agricultural_economics": self._get_agricultural_economics_config(),
            "livestock_systems": self._get_livestock_systems_config(),
            "irrigation_management": self._get_irrigation_config(),
            "pest_disease_modeling": self._get_pest_disease_config(),
            "agricultural_ml_platform": self._get_agricultural_ml_config()
        }
    
    def _get_crop_modeling_config(self) -> Dict[str, Any]:
        """Crop growth modeling and yield prediction platform"""
        return {
            "name": "Crop Modeling & Yield Prediction Platform",
            "description": "Comprehensive crop growth modeling, yield prediction, and climate impact assessment",
            "spack_packages": [
                # Crop modeling frameworks
                "dssat@4.8.2 %gcc@11.4.0 +fortran +netcdf",  # Decision Support System for Agrotechnology Transfer
                "apsim@2023.05.7336 %gcc@11.4.0 +mono +sqlite", # Agricultural Production Systems sIMulator
                "stics@10.1.0 %gcc@11.4.0 +fortran",        # Simulateur mulTIdisciplinaire pour les Cultures Standard
                "cropgrow@3.1.0 %gcc@11.4.0 +python",      # Generic crop growth model
                
                # Climate and weather models
                "wofost@7.2.1 %gcc@11.4.0 +python +netcdf", # World Food Studies crop model
                "pcse@5.5.3 %gcc@11.4.0 +python",          # Python Crop Simulation Environment
                "aquacrop@7.0 %gcc@11.4.0 +python",        # AquaCrop water-driven crop model
                
                # Soil-plant-atmosphere models
                "swap@4.2.0 %gcc@11.4.0 +fortran +netcdf", # Soil-Water-Atmosphere-Plant model
                "hydrus@3.05 %gcc@11.4.0 +fortran",        # Soil water flow and solute transport
                "century@5.0 %gcc@11.4.0 +fortran",        # Soil organic matter model
                
                # Python agricultural modeling
                "python@3.11.5 %gcc@11.4.0",
                "py-croppy@2.3.0 %gcc@11.4.0",             # Crop modeling Python package
                "py-pydssat@0.5.0 %gcc@11.4.0",            # DSSAT Python interface
                "py-pcse@5.5.3 %gcc@11.4.0",               # Python Crop Simulation Environment
                
                # Agricultural data processing
                "py-agpy@1.2.0 %gcc@11.4.0",               # Agricultural Python tools
                "py-agrometeorology@0.8.0 %gcc@11.4.0",    # Agricultural meteorology
                "py-crop-calendar@1.1.0 %gcc@11.4.0",      # Crop calendar tools
                
                # Geospatial and remote sensing
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-earthengine-api@0.1.364 %gcc@11.4.0",  # Google Earth Engine
                "py-sentinelsat@1.2.1 %gcc@11.4.0",        # Sentinel satellite data
                
                # Climate data and analysis
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-netcdf4@1.6.4 %gcc@11.4.0",
                "py-cftime@1.6.2 %gcc@11.4.0",
                "py-xclim@0.45.0 %gcc@11.4.0",             # Climate data processing
                
                # Statistical analysis and optimization
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Machine learning for agriculture
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                
                # AWS-optimized parallel computing with EFA support
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers +pmix +pmi +fabrics",
                "libfabric@1.18.1 %gcc@11.4.0 +verbs +mlx +efa",  # EFA support
                "aws-ofi-nccl@1.7.0 %gcc@11.4.0",  # AWS OFI plugin for NCCL
                "ucx@1.14.1 %gcc@11.4.0 +verbs +mlx +ib_hw_tm",  # Unified Communication X
                "py-mpi4py@3.1.4 %gcc@11.4.0",
                "py-dask@2023.7.1 %gcc@11.4.0",
                "slurm@23.02.5 %gcc@11.4.0 +pmix +numa",  # Slurm for cluster management
                
                # Database and data management
                "sqlite@3.42.0 %gcc@11.4.0",
                "postgresql@15.4 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "gfortran@11.4.0"
            ],
            "aws_instance_recommendations": {
                "development": {
                    "instance_type": "c6i.xlarge",
                    "vcpus": 4,
                    "memory_gb": 8,
                    "storage_gb": 100,
                    "cost_per_hour": 0.17,
                    "use_case": "Model development and small-scale testing"
                },
                "research_workstation": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage_gb": 500,
                    "cost_per_hour": 0.68,
                    "use_case": "Single farm or field-scale crop modeling"
                },
                "regional_analysis": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Regional crop modeling and multi-year climate analysis"
                },
                "national_modeling": {
                    "instance_type": "r6i.16xlarge",
                    "vcpus": 64,
                    "memory_gb": 512,
                    "storage_gb": 5000,
                    "cost_per_hour": 4.10,
                    "use_case": "National-scale crop modeling and climate impact assessment"
                }
            },
            "estimated_cost": {
                "compute": 500,
                "storage": 100,
                "data_transfer": 50,
                "total": 650
            },
            "research_capabilities": [
                "Crop growth modeling with DSSAT, APSIM, and STICS",
                "Climate impact assessment on agricultural systems",
                "Yield prediction and forecasting",
                "Soil-plant-atmosphere interaction modeling",
                "Multi-model ensemble simulations",
                "Satellite data integration for crop monitoring",
                "Regional and national scale agricultural analysis",
                "Agricultural risk assessment and adaptation planning"
            ],
            "aws_data_sources": [
                "NASA Harvest (crop yield data)",
                "USDA NASS (agricultural statistics)",
                "NOAA Climate Data Online",
                "Landsat and Sentinel satellite imagery",
                "Global weather station networks"
            ]
        }
    
    def _get_precision_agriculture_config(self) -> Dict[str, Any]:
        """Precision agriculture and smart farming platform"""
        return {
            "name": "Precision Agriculture & Smart Farming Platform",
            "description": "IoT sensor integration, variable rate applications, and precision agriculture analytics",
            "spack_packages": [
                # Precision agriculture frameworks
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql +grass",
                "grass@8.3.0 %gcc@11.4.0 +netcdf +postgresql +sqlite",
                "saga@9.1.1 %gcc@11.4.0 +python",
                
                # Remote sensing and image processing
                "gdal@3.7.2 %gcc@11.4.0 +python +postgresql +netcdf +hdf5",
                "opencv@4.8.0 %gcc@11.4.0 +python +contrib +vtk",
                "orfeo-toolbox@9.0.0 %gcc@11.4.0 +python +qt",
                
                # Drone and UAV data processing
                "opendronemap@3.3.0 %gcc@11.4.0 +python +opencv",
                "meshlab@2023.12 %gcc@11.4.0 +qt",
                "cloudcompare@2.13.1 %gcc@11.4.0 +qt",
                
                # IoT and sensor data processing
                "python@3.11.5 %gcc@11.4.0",
                "py-paho-mqtt@1.6.1 %gcc@11.4.0",          # MQTT client
                "py-influxdb@5.3.1 %gcc@11.4.0",           # Time series database
                "py-pyserial@3.5 %gcc@11.4.0",             # Serial communication
                
                # Precision agriculture Python tools
                "py-precision-agriculture@1.5.0 %gcc@11.4.0",
                "py-farmpy@0.3.0 %gcc@11.4.0",             # Farm management tools
                "py-agtech@2.1.0 %gcc@11.4.0",             # Agricultural technology
                
                # Variable rate application
                "py-shapely@2.0.1 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                
                # Machine learning for precision agriculture
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                
                # Image processing and computer vision
                "py-scikit-image@0.21.0 %gcc@11.4.0",
                "py-pillow@10.0.0 %gcc@11.4.0",
                "py-opencv@4.8.0 %gcc@11.4.0",
                
                # Geospatial analysis
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-rasterstats@0.19.0 %gcc@11.4.0",
                "py-earthengine-api@0.1.364 %gcc@11.4.0",
                
                # Data processing and analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                
                # Visualization and dashboards
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",
                "py-streamlit@1.25.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "redis@7.0.12 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "development": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 50,
                    "cost_per_hour": 0.085,
                    "use_case": "Development and testing of precision agriculture applications"
                },
                "field_operations": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Real-time field data processing and variable rate applications"
                },
                "farm_analysis": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 1000,
                    "cost_per_hour": 1.02,
                    "use_case": "Multi-field analysis and farm-scale optimization"
                },
                "regional_precision": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Regional precision agriculture analysis and benchmarking"
                }
            },
            "estimated_cost": {
                "compute": 300,
                "storage": 150,
                "iot_ingestion": 100,
                "data_transfer": 75,
                "total": 625
            },
            "research_capabilities": [
                "IoT sensor data integration and processing",
                "Drone and satellite imagery analysis",
                "Variable rate application mapping",
                "Yield monitoring and analysis",
                "Soil sampling optimization",
                "Prescription map generation",
                "Field boundary delineation",
                "Crop health monitoring with NDVI and other indices"
            ]
        }
    
    def _get_soil_science_config(self) -> Dict[str, Any]:
        """Soil science research and modeling platform"""
        return {
            "name": "Soil Science Research & Modeling Laboratory",
            "description": "Soil physics, chemistry, biology modeling and digital soil mapping",
            "spack_packages": [
                # Soil modeling frameworks
                "hydrus@3.05 %gcc@11.4.0 +fortran",        # Soil water flow modeling
                "swap@4.2.0 %gcc@11.4.0 +fortran +netcdf", # Soil-Water-Atmosphere-Plant
                "century@5.0 %gcc@11.4.0 +fortran",        # Soil organic matter
                "rothc@26.3 %gcc@11.4.0 +fortran",         # Rothamsted Carbon model
                "dndc@9.5 %gcc@11.4.0 +fortran",           # Soil greenhouse gas emissions
                
                # Digital soil mapping
                "saga@9.1.1 %gcc@11.4.0 +python",
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql",
                "grass@8.3.0 %gcc@11.4.0 +netcdf +postgresql",
                
                # Python soil science tools
                "python@3.11.5 %gcc@11.4.0",
                "py-soilgrids@1.0.0 %gcc@11.4.0",          # Global soil information
                "py-pysolar@0.10 %gcc@11.4.0",             # Solar calculations for soil
                "py-pysoil@2.1.0 %gcc@11.4.0",             # Soil science tools
                
                # Geochemical modeling
                "phreeqc@3.7.3 %gcc@11.4.0",               # Geochemical calculations
                "py-phreeqpython@1.5.2 %gcc@11.4.0",       # PHREEQC Python interface
                "py-pygeochem@0.1.4 %gcc@11.4.0",          # Geochemical data analysis
                
                # Geospatial and remote sensing
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Machine learning for soil mapping
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                "py-randomforest@1.0.0 %gcc@11.4.0",
                
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
                "soil_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Soil sample analysis and laboratory data processing"
                },
                "soil_modeling": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Soil process modeling and digital soil mapping"
                },
                "regional_mapping": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Large-scale digital soil mapping and prediction"
                }
            },
            "estimated_cost": {
                "compute": 400,
                "storage": 80,
                "data_transfer": 30,
                "total": 510
            },
            "research_capabilities": [
                "Soil water flow and transport modeling",
                "Soil carbon dynamics and greenhouse gas emissions",
                "Digital soil mapping and prediction",
                "Soil geochemistry and nutrient cycling",
                "Soil erosion and conservation modeling",
                "Soil-plant-atmosphere interactions"
            ]
        }
    
    def _get_plant_breeding_config(self) -> Dict[str, Any]:
        """Plant breeding and agricultural genomics platform"""
        return {
            "name": "Plant Breeding & Agricultural Genomics Platform",
            "description": "Genomic selection, QTL mapping, and breeding program optimization",
            "spack_packages": [
                # Genomics and breeding software
                "tassel@5.2.88 %gcc@11.4.0 +java",         # Trait Analysis by aSSociation, Evolution and Linkage
                "gapit@3.1.0 %gcc@11.4.0 +r",              # Genome Association and Prediction Integrated Tool
                "asreml@4.2 %gcc@11.4.0 +fortran",         # Mixed model analysis
                "blupf90@1.66 %gcc@11.4.0 +fortran",       # Mixed model equations
                
                # R genetics packages (via Spack)
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-qtl@1.66 %gcc@11.4.0",                  # QTL mapping
                "r-rrblup@4.6.3 %gcc@11.4.0",              # Ridge regression BLUP
                "r-sommer@4.2.1 %gcc@11.4.0",              # Solving mixed model equations
                "r-asreml@4.2.0 %gcc@11.4.0",              # ASReml-R package
                
                # Python genomics tools
                "python@3.11.5 %gcc@11.4.0",
                "py-biopython@1.81 %gcc@11.4.0",
                "py-pysam@0.21.0 %gcc@11.4.0",
                "py-scikit-allel@1.3.7 %gcc@11.4.0",       # Population genetics analysis
                "py-pyvcf@0.6.8 %gcc@11.4.0",              # VCF file processing
                
                # Breeding program optimization
                "py-breeding-programs@2.0.0 %gcc@11.4.0",
                "py-genomic-selection@1.5.0 %gcc@11.4.0",
                "py-qtl-analysis@0.8.0 %gcc@11.4.0",
                
                # Machine learning for breeding
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
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
                "breeding_analysis": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 500,
                    "cost_per_hour": 0.51,
                    "use_case": "Small breeding program analysis and QTL mapping"
                },
                "genomic_selection": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 1000,
                    "cost_per_hour": 1.02,
                    "use_case": "Genomic selection and breeding value estimation"
                },
                "population_genomics": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Large-scale population genomics and GWAS"
                }
            },
            "estimated_cost": {
                "compute": 600,
                "storage": 120,
                "data_transfer": 40,
                "total": 760
            },
            "research_capabilities": [
                "Genomic selection and breeding value prediction",
                "QTL mapping and genome-wide association studies",
                "Population genetics and genetic diversity analysis",
                "Breeding program optimization",
                "Marker-assisted selection",
                "Genomic estimated breeding values (GEBV)"
            ]
        }
    
    def _get_agricultural_economics_config(self) -> Dict[str, Any]:
        """Agricultural economics and policy analysis platform"""
        return {
            "name": "Agricultural Economics & Policy Analysis Platform", 
            "description": "Market analysis, policy simulation, and agricultural economics modeling",
            "spack_packages": [
                # Economic modeling frameworks
                "gams@44.3.0 %gcc@11.4.0",                 # General Algebraic Modeling System
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-plm@2.6-2 %gcc@11.4.0",                 # Panel data econometrics
                "r-vars@1.5-10 %gcc@11.4.0",               # Vector autoregression
                
                # Python economics tools
                "python@3.11.5 %gcc@11.4.0",
                "py-pyeconomics@1.2.0 %gcc@11.4.0",
                "py-agricultural-economics@0.9.0 %gcc@11.4.0",
                "py-policy-simulation@2.1.0 %gcc@11.4.0",
                
                # Statistical and econometric analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                
                # Time series analysis
                "py-arch@6.2.0 %gcc@11.4.0",               # ARCH/GARCH models
                "py-pmdarima@2.0.3 %gcc@11.4.0",           # ARIMA models
                "py-prophet@1.1.4 %gcc@11.4.0",            # Forecasting
                
                # Optimization
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                "py-pyomo@6.6.1 %gcc@11.4.0",
                
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
                "policy_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Agricultural policy analysis and market modeling"
                },
                "economic_modeling": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Large-scale economic modeling and simulation"
                }
            },
            "estimated_cost": {
                "compute": 250,
                "storage": 50,
                "data_transfer": 25,
                "total": 325
            },
            "research_capabilities": [
                "Agricultural market analysis and forecasting",
                "Policy impact assessment and simulation",
                "Farm-level economic optimization",
                "Supply chain analysis",
                "Risk management and insurance modeling",
                "Trade and international agricultural economics"
            ]
        }
    
    def _get_livestock_systems_config(self) -> Dict[str, Any]:
        """Livestock systems and animal science platform"""
        return {
            "name": "Livestock Systems & Animal Science Platform",
            "description": "Animal breeding, nutrition modeling, and livestock systems analysis",
            "spack_packages": [
                # Animal breeding software
                "blupf90@1.66 %gcc@11.4.0 +fortran",       # Mixed model equations for animals
                "asreml@4.2 %gcc@11.4.0 +fortran",         # Mixed model analysis
                "dmu@6.5.5 %gcc@11.4.0 +fortran",          # Multivariate mixed models
                
                # Nutrition modeling
                "cncps@6.55 %gcc@11.4.0 +fortran",         # Cornell Net Carbohydrate and Protein System
                "nasem@2016 %gcc@11.4.0 +fortran",         # Nutrient Requirements models
                "inra-feeds@2018 %gcc@11.4.0 +fortran",    # INRA feeding system
                
                # Python animal science tools
                "python@3.11.5 %gcc@11.4.0",
                "py-animal-science@1.8.0 %gcc@11.4.0",
                "py-livestock-systems@2.3.0 %gcc@11.4.0",
                "py-breeding-analysis@1.1.0 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                
                # Machine learning
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",
                
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
                "breeding_analysis": {
                    "instance_type": "r6i.2xlarge", 
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 300,
                    "cost_per_hour": 0.51,
                    "use_case": "Animal breeding value estimation and selection"
                },
                "nutrition_modeling": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage_gb": 500,
                    "cost_per_hour": 0.68,
                    "use_case": "Feed formulation and nutrition optimization"
                }
            },
            "estimated_cost": {
                "compute": 350,
                "storage": 70,
                "data_transfer": 20,
                "total": 440
            },
            "research_capabilities": [
                "Animal breeding value estimation",
                "Feed formulation and nutrition optimization", 
                "Livestock production system modeling",
                "Genetic parameter estimation",
                "Growth curve analysis",
                "Reproductive performance modeling"
            ]
        }
    
    def _get_irrigation_config(self) -> Dict[str, Any]:
        """Irrigation management and water use efficiency platform"""
        return {
            "name": "Irrigation Management & Water Use Efficiency Platform",
            "description": "Irrigation scheduling, water balance modeling, and precision water management",
            "spack_packages": [
                # Irrigation and water modeling
                "cropwat@8.0 %gcc@11.4.0 +fortran",        # Crop water requirements
                "aquacrop@7.0 %gcc@11.4.0 +python",        # Water productivity model
                "swap@4.2.0 %gcc@11.4.0 +fortran +netcdf", # Soil-Water-Atmosphere-Plant
                "hydrus@3.05 %gcc@11.4.0 +fortran",        # Soil water flow
                
                # Python irrigation tools
                "python@3.11.5 %gcc@11.4.0",
                "py-irrigation-scheduling@1.4.0 %gcc@11.4.0",
                "py-evapotranspiration@0.7.0 %gcc@11.4.0",
                "py-water-balance@2.1.0 %gcc@11.4.0",
                
                # Remote sensing for irrigation
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-earthengine-api@0.1.364 %gcc@11.4.0",
                "py-sentinelsat@1.2.1 %gcc@11.4.0",
                
                # Weather data processing
                "py-netcdf4@1.6.4 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-metpy@1.5.0 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                
                # Optimization
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                
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
                "irrigation_scheduling": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 100,
                    "cost_per_hour": 0.085,
                    "use_case": "Real-time irrigation scheduling and monitoring"
                },
                "water_optimization": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 300,
                    "cost_per_hour": 0.34,
                    "use_case": "Water use optimization and efficiency analysis"
                }
            },
            "estimated_cost": {
                "compute": 200,
                "storage": 60,
                "data_transfer": 30,
                "total": 290
            },
            "research_capabilities": [
                "Irrigation scheduling optimization",
                "Crop water requirement estimation",
                "Soil water balance modeling",
                "Evapotranspiration calculation",
                "Water use efficiency analysis",
                "Deficit irrigation strategies"
            ]
        }
    
    def _get_pest_disease_config(self) -> Dict[str, Any]:
        """Pest and disease management modeling platform"""
        return {
            "name": "Pest & Disease Management Modeling Platform",
            "description": "Integrated pest management, disease forecasting, and epidemiological modeling",
            "spack_packages": [
                # Disease and pest modeling
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-epiphy@0.4.0 %gcc@11.4.0",              # Plant disease epidemiology
                "r-agricolae@1.3-5 %gcc@11.4.0",           # Agricultural statistics
                
                # Python pest management tools
                "python@3.11.5 %gcc@11.4.0",
                "py-pest-modeling@1.6.0 %gcc@11.4.0",
                "py-disease-forecasting@2.0.0 %gcc@11.4.0",
                "py-ipm-tools@1.3.0 %gcc@11.4.0",          # Integrated pest management
                
                # Machine learning for pest detection
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-opencv@4.8.0 %gcc@11.4.0",
                
                # Image processing for disease detection
                "py-scikit-image@0.21.0 %gcc@11.4.0",
                "py-pillow@10.0.0 %gcc@11.4.0",
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Epidemiological modeling
                "py-epimodels@0.9.0 %gcc@11.4.0",
                "py-networkx@3.1 %gcc@11.4.0",
                
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
                "disease_monitoring": {
                    "instance_type": "g4dn.xlarge",
                    "vcpus": 4,
                    "memory_gb": 16,
                    "storage_gb": 125,
                    "cost_per_hour": 0.526,
                    "use_case": "AI-powered disease detection and image analysis"
                },
                "epidemiological_modeling": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage_gb": 500,
                    "cost_per_hour": 0.68,
                    "use_case": "Disease spread modeling and forecasting"
                }
            },
            "estimated_cost": {
                "compute": 400,
                "storage": 80,
                "data_transfer": 30,
                "total": 510
            },
            "research_capabilities": [
                "Disease forecasting and risk assessment",
                "Pest population dynamics modeling",
                "AI-powered disease detection from images",
                "Epidemiological spread modeling",
                "Integrated pest management optimization",
                "Pesticide resistance modeling"
            ]
        }
    
    def _get_agricultural_ml_config(self) -> Dict[str, Any]:
        """Agricultural machine learning and AI platform"""
        return {
            "name": "Agricultural Machine Learning & AI Platform",
            "description": "Advanced AI/ML for agriculture, computer vision, and predictive analytics",
            "spack_packages": [
                # Machine learning frameworks
                "python@3.11.5 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                
                # Computer vision for agriculture
                "py-opencv@4.8.0 %gcc@11.4.0",
                "py-scikit-image@0.21.0 %gcc@11.4.0",
                "py-pillow@10.0.0 %gcc@11.4.0",
                "py-torchvision@0.15.2 %gcc@11.4.0",
                
                # Agricultural AI tools
                "py-plantcv@4.1.0 %gcc@11.4.0",            # Plant phenotyping
                "py-agml@1.0.0 %gcc@11.4.0",               # Agricultural machine learning
                "py-crop-yield-prediction@2.3.0 %gcc@11.4.0",
                
                # Deep learning for remote sensing
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-earthengine-api@0.1.364 %gcc@11.4.0",
                "py-satellite-image-deep-learning@1.2.0 %gcc@11.4.0",
                
                # Time series forecasting
                "py-prophet@1.1.4 %gcc@11.4.0",
                "py-statsforecast@1.6.0 %gcc@11.4.0",
                "py-neuralforecast@1.6.4 %gcc@11.4.0",
                
                # Natural language processing for agriculture
                "py-transformers@4.33.2 %gcc@11.4.0",
                "py-spacy@3.6.1 %gcc@11.4.0",
                "py-nltk@3.8.1 %gcc@11.4.0",
                
                # Data processing
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                
                # Model serving and deployment
                "py-fastapi@0.103.0 %gcc@11.4.0",
                "py-mlflow@2.5.0 %gcc@11.4.0",
                "py-wandb@0.15.8 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-streamlit@1.25.0 %gcc@11.4.0",
                
                # Parallel computing
                "py-dask@2023.7.1 %gcc@11.4.0",
                "py-ray@2.6.1 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "redis@7.0.12 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "ml_development": {
                    "instance_type": "g4dn.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 32,
                    "storage_gb": 225,
                    "cost_per_hour": 0.752,
                    "use_case": "ML model development and training"
                },
                "computer_vision": {
                    "instance_type": "g5.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 64,
                    "storage_gb": 600,
                    "cost_per_hour": 1.624,
                    "use_case": "Computer vision and image analysis"
                },
                "large_scale_training": {
                    "instance_type": "p4d.24xlarge",
                    "vcpus": 96,
                    "memory_gb": 1152,
                    "storage_gb": 8000,
                    "cost_per_hour": 32.77,
                    "use_case": "Large-scale model training and hyperparameter tuning"
                }
            },
            "estimated_cost": {
                "compute": 1200,
                "storage": 200,
                "data_transfer": 100,
                "total": 1500
            },
            "research_capabilities": [
                "Crop yield prediction using satellite imagery",
                "Plant disease detection with computer vision",
                "Agricultural chatbots and knowledge systems",
                "Precision agriculture optimization with AI",
                "Market price forecasting",
                "Climate impact prediction on agriculture"
            ]
        }
    
    def generate_agricultural_recommendation(self, workload: AgriculturalWorkload) -> Dict[str, Any]:
        """Generate optimized AWS infrastructure recommendation for agricultural research"""
        
        # Select appropriate configuration based on domain
        domain_config_map = {
            AgriculturalDomain.CROP_MODELING: "crop_modeling_platform",
            AgriculturalDomain.PRECISION_AGRICULTURE: "precision_agriculture", 
            AgriculturalDomain.SOIL_SCIENCE: "soil_science_laboratory",
            AgriculturalDomain.PLANT_BREEDING: "plant_breeding_genomics",
            AgriculturalDomain.AGRICULTURAL_ECONOMICS: "agricultural_economics",
            AgriculturalDomain.LIVESTOCK_SYSTEMS: "livestock_systems",
            AgriculturalDomain.IRRIGATION_MANAGEMENT: "irrigation_management",
            AgriculturalDomain.PEST_DISEASE_MANAGEMENT: "pest_disease_modeling",
            AgriculturalDomain.AGRICULTURAL_GENOMICS: "agricultural_ml_platform"
        }
        
        config_name = domain_config_map.get(workload.domain, "crop_modeling_platform")
        base_config = self.agricultural_configurations[config_name].copy()
        
        # Adjust configuration based on workload characteristics
        self._optimize_for_scale(base_config, workload)
        self._optimize_for_data_volume(base_config, workload)
        self._optimize_for_computational_intensity(base_config, workload)
        
        # Generate cost estimates
        base_config["estimated_cost"] = self._calculate_agricultural_costs(workload, base_config)
        
        # Add optimization recommendations
        base_config["optimization_recommendations"] = self._generate_optimization_recommendations(workload)
        
        return {
            "configuration": base_config,
            "workload_analysis": {
                "domain": workload.domain.value,
                "farm_scale": workload.farm_scale,
                "analysis_type": workload.analysis_type,
                "computational_requirements": workload.computational_intensity,
                "data_volume": f"{workload.data_volume_tb} TB"
            },
            "deployment_recommendations": self._generate_deployment_recommendations(workload),
            "estimated_cost": base_config["estimated_cost"]
        }
    
    def _optimize_for_scale(self, config: Dict[str, Any], workload: AgriculturalWorkload):
        """Optimize configuration based on farm scale"""
        scale_multipliers = {
            "Field": 1.0,
            "Farm": 2.0, 
            "Regional": 4.0,
            "National": 8.0,
            "Global": 16.0
        }
        
        multiplier = scale_multipliers.get(workload.farm_scale, 1.0)
        
        # Adjust instance recommendations based on scale
        if "aws_instance_recommendations" in config:
            for instance_config in config["aws_instance_recommendations"].values():
                if multiplier > 4.0:
                    # Scale up for large geographical areas
                    if "c6i" in instance_config["instance_type"]:
                        instance_config["instance_type"] = instance_config["instance_type"].replace("c6i", "c6i")
                    instance_config["storage_gb"] = int(instance_config["storage_gb"] * multiplier)
    
    def _optimize_for_data_volume(self, config: Dict[str, Any], workload: AgriculturalWorkload):
        """Optimize configuration based on expected data volume"""
        if workload.data_volume_tb > 10.0:
            # Add data processing optimizations for large datasets
            if "spack_packages" in config:
                config["spack_packages"].extend([
                    "py-dask@2023.7.1 %gcc@11.4.0",
                    "py-ray@2.6.1 %gcc@11.4.0"
                ])
    
    def _optimize_for_computational_intensity(self, config: Dict[str, Any], workload: AgriculturalWorkload):
        """Optimize configuration based on computational intensity"""
        if workload.computational_intensity in ["Intensive", "Extreme"]:
            # Upgrade to compute-optimized or GPU instances
            if "aws_instance_recommendations" in config:
                for key, instance_config in config["aws_instance_recommendations"].items():
                    if workload.computational_intensity == "Extreme":
                        if "Computer vision" in workload.analysis_type or "AI" in workload.analysis_type:
                            instance_config["instance_type"] = "g5.2xlarge"
                            instance_config["cost_per_hour"] = 1.624
    
    def _calculate_agricultural_costs(self, workload: AgriculturalWorkload, config: Dict[str, Any]) -> Dict[str, float]:
        """Calculate estimated costs for agricultural research infrastructure"""
        base_compute = 400
        base_storage = 100
        base_data_transfer = 50
        
        # Scale costs based on farm scale
        scale_multipliers = {"Field": 1.0, "Farm": 1.5, "Regional": 3.0, "National": 6.0, "Global": 12.0}
        multiplier = scale_multipliers.get(workload.farm_scale, 1.0)
        
        # Adjust for computational intensity
        intensity_multipliers = {"Light": 0.5, "Moderate": 1.0, "Intensive": 2.0, "Extreme": 4.0}
        intensity_mult = intensity_multipliers.get(workload.computational_intensity, 1.0)
        
        compute_cost = base_compute * multiplier * intensity_mult
        storage_cost = base_storage * (1 + workload.data_volume_tb / 10.0)
        data_transfer_cost = base_data_transfer * multiplier
        
        return {
            "compute": compute_cost,
            "storage": storage_cost,
            "data_transfer": data_transfer_cost,
            "total": compute_cost + storage_cost + data_transfer_cost
        }
    
    def _generate_optimization_recommendations(self, workload: AgriculturalWorkload) -> List[str]:
        """Generate optimization recommendations for agricultural workloads"""
        recommendations = []
        
        if workload.farm_scale in ["Regional", "National", "Global"]:
            recommendations.append("Consider using Spot Instances for batch processing to reduce costs by 60-90%")
            recommendations.append("Implement auto-scaling for seasonal workload variations")
        
        if workload.data_volume_tb > 5.0:
            recommendations.append("Use S3 Intelligent Tiering for automatic cost optimization of large datasets")
            recommendations.append("Consider AWS Batch for parallel processing of large spatial datasets")
        
        if "Real-time" in workload.temporal_scale:
            recommendations.append("Use AWS IoT Core for real-time sensor data ingestion")
            recommendations.append("Consider Amazon Kinesis for streaming agricultural data processing")
        
        if workload.computational_intensity == "Extreme":
            recommendations.append("Use GPU instances for deep learning and computer vision tasks")
            recommendations.append("Consider AWS ParallelCluster for HPC workloads")
        
        return recommendations
    
    def _generate_deployment_recommendations(self, workload: AgriculturalWorkload) -> Dict[str, Any]:
        """Generate deployment recommendations for agricultural research"""
        return {
            "deployment_strategy": "multi-tier" if workload.farm_scale in ["Regional", "National", "Global"] else "single-tier",
            "backup_strategy": "automated_daily_snapshots",
            "monitoring": ["CloudWatch for infrastructure", "Custom dashboards for agricultural metrics"],
            "security": ["VPC with private subnets", "IAM roles for service access", "Data encryption at rest and in transit"],
            "disaster_recovery": "cross-region backup" if workload.farm_scale in ["National", "Global"] else "single-region backup"
        }

if __name__ == "__main__":
    # Example usage
    pack = AgriculturalSciencesPack()
    
    # Example workload
    workload = AgriculturalWorkload(
        domain=AgriculturalDomain.CROP_MODELING,
        farm_scale="Regional",
        crop_types=["Corn", "Soybean"],
        analysis_type="Yield Prediction",
        temporal_scale="Seasonal",
        data_sources=["Satellite", "Weather"],
        modeling_approach="Machine Learning",
        data_volume_tb=2.5,
        computational_intensity="Moderate"
    )
    
    recommendation = pack.generate_agricultural_recommendation(workload)
    print(json.dumps(recommendation, indent=2))