#!/usr/bin/env python3
"""
Geospatial Research Pack
Comprehensive geospatial, remote sensing, and earth science environments for AWS Research Wizard
"""

import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from enum import Enum

class GeospatialDomain(Enum):
    REMOTE_SENSING = "remote_sensing"
    GIS_ANALYSIS = "gis_analysis"
    GEOPHYSICS = "geophysics"
    ENVIRONMENTAL_MODELING = "environmental_modeling"
    PRECISION_AGRICULTURE = "precision_agriculture"
    DISASTER_RESPONSE = "disaster_response"
    URBAN_PLANNING = "urban_planning"
    MARITIME_COASTAL = "maritime_coastal"

@dataclass
class GeospatialWorkload:
    """Geospatial workload characteristics"""
    domain: GeospatialDomain
    data_sources: List[str]  # Landsat, Sentinel, MODIS, etc.
    processing_type: str     # Classification, change detection, modeling
    spatial_resolution: str  # High (1-10m), Medium (10-100m), Low (>100m)
    temporal_frequency: str  # Real-time, Daily, Weekly, Monthly, Annual
    coverage_area: str       # Local, Regional, National, Global
    data_volume_tb: float    # Expected data volume in TB
    analysis_complexity: str # Simple, Moderate, Complex, Advanced

class GeospatialResearchPack:
    """
    Comprehensive geospatial research environments optimized for AWS
    Supports remote sensing, GIS analysis, geophysics, and environmental modeling
    """

    def __init__(self):
        self.geospatial_configurations = {
            "remote_sensing_lab": self._get_remote_sensing_config(),
            "gis_analysis_studio": self._get_gis_analysis_config(),
            "geophysics_workstation": self._get_geophysics_config(),
            "environmental_modeling": self._get_environmental_config(),
            "precision_agriculture": self._get_agriculture_config(),
            "disaster_response": self._get_disaster_config(),
            "urban_planning": self._get_urban_planning_config(),
            "maritime_coastal": self._get_maritime_config(),
            "geospatial_ml_platform": self._get_ml_platform_config()
        }

    def _get_remote_sensing_config(self) -> Dict[str, Any]:
        """Remote sensing and satellite imagery analysis"""
        return {
            "name": "Remote Sensing & Earth Observation Laboratory",
            "description": "Satellite imagery processing, analysis, and machine learning",
            "spack_packages": [
                # Core remote sensing tools
                "gdal@3.7.1 %gcc@11.4.0 +python +netcdf +hdf5 +geos +proj +curl",
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql +grass +saga",
                "grass@8.3.0 %gcc@11.4.0 +python +postgresql +netcdf +geos",
                "saga-gis@9.1.1 %gcc@11.4.0 +python +postgresql",

                # Satellite data processing
                "otb@9.0.0 %gcc@11.4.0 +python +opencv +fftw",  # Orfeo ToolBox
                "pktools@2.6.7 %gcc@11.4.0 +python +gdal",
                "snap@10.0.0 %gcc@11.4.0 +python",  # ESA SNAP

                # Python remote sensing ecosystem
                "python@3.11.5 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-dask@2023.8.0 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",

                # Earth observation specific
                "py-earthpy@0.9.4 %gcc@11.4.0",
                "py-satpy@0.43.0 %gcc@11.4.0",
                "py-pyresample@1.27.1 %gcc@11.4.0",
                "py-pystac@1.8.2 %gcc@11.4.0",
                "py-stackstac@0.5.0 %gcc@11.4.0",

                # Machine learning for RS
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-keras@2.13.1 %gcc@11.4.0",

                # Image processing
                "py-opencv@4.8.0 %gcc@11.4.0",
                "py-scikit-image@0.21.0 %gcc@11.4.0",
                "py-pillow@10.0.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",

                # Data formats and I/O
                "hdf5@1.14.2 %gcc@11.4.0 +mpi +threadsafe +fortran",
                "netcdf-c@4.9.2 %gcc@11.4.0 +mpi +parallel-netcdf",
                "netcdf-fortran@4.6.1 %gcc@11.4.0",

                # Parallel processing
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers",
                "py-mpi4py@3.1.4 %gcc@11.4.0"
            ],
            "aws_data_sources": {
                "landsat": {
                    "description": "Landsat 8/9 imagery (30m resolution)",
                    "bucket": "s3://landsat-pds/",
                    "cost": "Free access",
                    "update_frequency": "16 days",
                    "coverage": "Global"
                },
                "sentinel_2": {
                    "description": "Sentinel-2 optical imagery (10m resolution)",
                    "bucket": "s3://sentinel-s2-l1c/",
                    "cost": "Free access",
                    "update_frequency": "5 days",
                    "coverage": "Global land"
                },
                "sentinel_1": {
                    "description": "Sentinel-1 SAR imagery (10m resolution)",
                    "bucket": "s3://sentinel-s1-l1c/",
                    "cost": "Free access",
                    "update_frequency": "6-12 days",
                    "coverage": "Global"
                },
                "modis": {
                    "description": "MODIS Terra/Aqua (250m-1km resolution)",
                    "bucket": "s3://modis-pds/",
                    "cost": "Free access",
                    "update_frequency": "Daily",
                    "coverage": "Global"
                },
                "viirs": {
                    "description": "VIIRS Day/Night Band",
                    "bucket": "s3://noaa-viirs-pds/",
                    "cost": "Free access",
                    "update_frequency": "Daily",
                    "coverage": "Global"
                }
            },
            "aws_instance_recommendations": {
                "small_projects": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage": "500GB gp3 SSD",
                    "cost_per_hour": 0.504,
                    "use_case": "Regional analysis, small datasets (<100GB)"
                },
                "research_group": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage": "1TB gp3 SSD + 2TB EFS",
                    "cost_per_hour": 1.008,
                    "use_case": "Multi-temporal analysis, moderate datasets (100GB-1TB)"
                },
                "large_scale": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage": "2TB gp3 SSD + 10TB EFS",
                    "cost_per_hour": 2.016,
                    "use_case": "Continental analysis, large datasets (1-10TB)"
                },
                "gpu_ml": {
                    "instance_type": "p3.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 61,
                    "gpu": "1x NVIDIA V100",
                    "cost_per_hour": 3.06,
                    "use_case": "Deep learning on satellite imagery"
                }
            },
            "specialized_workflows": [
                "Land cover classification",
                "Change detection analysis",
                "Vegetation index time series",
                "Disaster damage assessment",
                "Urban growth monitoring",
                "Agricultural monitoring",
                "Forest fire detection",
                "Water quality assessment",
                "Flood mapping",
                "Drought monitoring"
            ],
            "cost_profile": {
                "small_research": "$300-800/month (regional studies)",
                "university_lab": "$800-2000/month (multi-project lab)",
                "operational_monitoring": "$1500-5000/month (continuous monitoring)",
                "commercial_services": "$3000-10000/month (large-scale analysis)"
            }
        }

    def _get_gis_analysis_config(self) -> Dict[str, Any]:
        """Geographic Information Systems analysis and spatial computing"""
        return {
            "name": "GIS Analysis & Spatial Computing Studio",
            "description": "Advanced GIS analysis, spatial modeling, and geographic data science",
            "spack_packages": [
                # Core GIS software
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql +grass +saga",
                "grass@8.3.0 %gcc@11.4.0 +python +postgresql +netcdf +geos +proj",
                "saga-gis@9.1.1 %gcc@11.4.0 +python +postgresql",

                # Spatial libraries
                "gdal@3.7.1 %gcc@11.4.0 +python +postgresql +geos +proj +curl",
                "geos@3.12.0 %gcc@11.4.0",
                "proj@9.2.1 %gcc@11.4.0",
                "spatialindex@1.9.3 %gcc@11.4.0",
                "libspatialite@5.0.1 %gcc@11.4.0",

                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +python",
                "postgis@3.3.3 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",

                # Python spatial ecosystem
                "python@3.11.5 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-shapely@2.0.1 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                "py-pyproj@3.6.0 %gcc@11.4.0",
                "py-rtree@1.0.1 %gcc@11.4.0",

                # Spatial analysis
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",

                # Network analysis
                "py-networkx@3.1 %gcc@11.4.0",
                "py-osmnx@1.6.0 %gcc@11.4.0",
                "py-momepy@0.6.0 %gcc@11.4.0",

                # Visualization
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-contextily@1.3.0 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",

                # Web mapping
                "mapserver@8.0.1 %gcc@11.4.0 +python +postgresql",
                "geoserver@2.23.2 %gcc@11.4.0"
            ],
            "aws_services_integration": [
                "Amazon Location Service for geocoding/routing",
                "S3 for spatial data storage",
                "RDS for PostGIS databases",
                "Lambda for spatial processing functions",
                "API Gateway for spatial web services",
                "CloudFront for map tile distribution"
            ],
            "analysis_capabilities": [
                "Spatial statistics and modeling",
                "Network analysis and routing",
                "Spatial clustering and hotspot analysis",
                "Geostatistical analysis (kriging, interpolation)",
                "Watershed and terrain analysis",
                "Accessibility and service area analysis",
                "Land suitability modeling",
                "Spatial decision support systems"
            ],
            "aws_instance_recommendations": {
                "desktop_gis": {
                    "instance_type": "m6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 32,
                    "storage": "500GB gp3 SSD",
                    "cost_per_hour": 0.384,
                    "use_case": "Individual GIS analysis, small datasets"
                },
                "spatial_analysis": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage": "1TB gp3 SSD + PostGIS RDS",
                    "cost_per_hour": 0.504,
                    "use_case": "Advanced spatial analysis, medium datasets"
                },
                "enterprise_gis": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage": "2TB gp3 SSD + Multi-AZ PostGIS",
                    "cost_per_hour": 1.008,
                    "use_case": "Enterprise GIS, large spatial databases"
                }
            },
            "cost_profile": {
                "individual_research": "$200-500/month",
                "small_organization": "$500-1200/month",
                "enterprise_deployment": "$1200-3000/month"
            }
        }

    def _get_geophysics_config(self) -> Dict[str, Any]:
        """Geophysics and subsurface modeling"""
        return {
            "name": "Geophysics & Subsurface Modeling Workstation",
            "description": "Seismic processing, geological modeling, and geophysical interpretation",
            "spack_packages": [
                # Seismic processing
                "seismic-unix@44R3 %gcc@11.4.0 +x11 +fftw",
                "madagascar@3.0 %gcc@11.4.0 +python +fftw",
                "opendt@2.0 %gcc@11.4.0 +mpi",

                # Geological modeling
                "gempy@3.0.0 %gcc@11.4.0 +python",
                "pygslib@0.0.0.6 %gcc@11.4.0",

                # Geophysical libraries
                "fftw@3.3.10 %gcc@11.4.0 +mpi +openmp",
                "openblas@0.3.24 %gcc@11.4.0 +openmp",
                "scalapack@2.2.0 %gcc@11.4.0",

                # Python geophysics
                "python@3.11.5 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",

                # Specialized geophysics packages
                "py-obspy@1.4.0 %gcc@11.4.0",  # Seismology
                "py-pyrocko@2023.06.16 %gcc@11.4.0",  # Seismological toolkit
                "py-simpeg@0.20.0 %gcc@11.4.0",  # Simulation and Parameter Estimation
                "py-fatiando@0.6 %gcc@11.4.0",  # Geophysical modeling

                # Visualization
                "py-mayavi@4.8.1 %gcc@11.4.0",
                "paraview@5.11.2 %gcc@11.4.0 +python +mpi",
                "visit@3.3.3 %gcc@11.4.0 +python",

                # Parallel computing
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers",
                "py-mpi4py@3.1.4 %gcc@11.4.0"
            ],
            "geophysical_methods": [
                "Seismic reflection/refraction processing",
                "Gravity and magnetic modeling",
                "Electrical resistivity tomography",
                "Ground-penetrating radar processing",
                "Magnetotelluric analysis",
                "Potential field modeling",
                "Seismic tomography",
                "Earthquake location and analysis"
            ],
            "aws_instance_recommendations": {
                "seismic_processing": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage": "4TB gp3 SSD",
                    "cost_per_hour": 1.632,
                    "use_case": "Seismic data processing, large surveys"
                },
                "geological_modeling": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage": "2TB gp3 SSD",
                    "cost_per_hour": 1.008,
                    "use_case": "3D geological modeling, reservoir simulation"
                },
                "hpc_cluster": {
                    "instance_type": "c6i.16xlarge",
                    "vcpus": 64,
                    "memory_gb": 128,
                    "count": "2-8 nodes",
                    "storage": "FSx Lustre 2.4TB",
                    "cost_per_hour": "$3.264-13.056 (2-8 nodes)",
                    "use_case": "Large-scale seismic modeling, tomography"
                }
            },
            "cost_profile": {
                "academic_research": "$400-1200/month",
                "industry_consulting": "$1200-4000/month",
                "large_scale_surveys": "$3000-12000/month"
            }
        }

    def _get_environmental_config(self) -> Dict[str, Any]:
        """Environmental monitoring and climate modeling"""
        return {
            "name": "Environmental Monitoring & Climate Analysis Platform",
            "description": "Environmental data analysis, climate modeling, and ecosystem monitoring",
            "spack_packages": [
                # Climate data processing
                "cdo@2.2.0 %gcc@11.4.0 +netcdf +hdf5 +proj +curl",
                "nco@5.1.6 %gcc@11.4.0 +netcdf +openmpi",
                "ncview@2.1.8 %gcc@11.4.0 +netcdf",

                # Environmental modeling
                "wrf@4.5.0 %gcc@11.4.0 +netcdf +hdf5 +mpi +openmp",
                "cesm@2.1.3 %gcc@11.4.0 +netcdf +pnetcdf +mpi",

                # Hydrology
                "grass@8.3.0 %gcc@11.4.0 +python +postgresql +netcdf +geos",
                "saga-gis@9.1.1 %gcc@11.4.0 +python +postgresql",

                # Python climate/environmental stack
                "python@3.11.5 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-dask@2023.8.0 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",

                # Climate analysis
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-metpy@1.5.1 %gcc@11.4.0",
                "py-windspharm@1.7.0 %gcc@11.4.0",
                "py-climate-indices@1.0.13 %gcc@11.4.0",

                # Remote sensing for environment
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-satpy@0.43.0 %gcc@11.4.0",
                "py-pyresample@1.27.1 %gcc@11.4.0",

                # Statistical analysis
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",

                # Data formats
                "netcdf-c@4.9.2 %gcc@11.4.0 +mpi +parallel-netcdf",
                "hdf5@1.14.2 %gcc@11.4.0 +mpi +threadsafe",
                "gdal@3.7.1 %gcc@11.4.0 +python +netcdf +hdf5"
            ],
            "environmental_applications": [
                "Climate change impact assessment",
                "Air quality monitoring and modeling",
                "Water quality assessment",
                "Ecosystem health monitoring",
                "Biodiversity analysis",
                "Carbon footprint analysis",
                "Renewable energy resource assessment",
                "Environmental risk assessment",
                "Habitat suitability modeling",
                "Conservation planning"
            ],
            "aws_instance_recommendations": {
                "environmental_analysis": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage": "1TB gp3 SSD",
                    "cost_per_hour": 0.504,
                    "use_case": "Environmental data analysis, moderate datasets"
                },
                "climate_modeling": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage": "2TB gp3 SSD",
                    "cost_per_hour": 0.816,
                    "use_case": "Regional climate modeling, WRF simulations"
                },
                "large_scale_monitoring": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage": "4TB gp3 SSD + 10TB EFS",
                    "cost_per_hour": 2.016,
                    "use_case": "Continental monitoring, large datasets"
                }
            },
            "cost_profile": {
                "research_project": "$300-800/month",
                "monitoring_program": "$800-2000/month",
                "operational_system": "$2000-6000/month"
            }
        }

    def _get_agriculture_config(self) -> Dict[str, Any]:
        """Precision agriculture and agricultural monitoring"""
        return {
            "name": "Precision Agriculture & Crop Monitoring Platform",
            "description": "Agricultural remote sensing, crop monitoring, and precision farming",
            "spack_packages": [
                # Core GIS and remote sensing
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql +grass",
                "grass@8.3.0 %gcc@11.4.0 +python +postgresql +netcdf",
                "gdal@3.7.1 %gcc@11.4.0 +python +netcdf +geos +proj",

                # Python ecosystem
                "python@3.11.5 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",

                # Agricultural remote sensing
                "py-satpy@0.43.0 %gcc@11.4.0",
                "py-pyresample@1.27.1 %gcc@11.4.0",
                "py-stackstac@0.5.0 %gcc@11.4.0",

                # Machine learning
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-keras@2.13.1 %gcc@11.4.0",

                # Image processing
                "py-opencv@4.8.0 %gcc@11.4.0",
                "py-scikit-image@0.21.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",

                # Statistical analysis
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python",
                "postgis@3.3.3 %gcc@11.4.0"
            ],
            "agricultural_applications": [
                "Crop yield prediction",
                "Vegetation health monitoring (NDVI, EVI)",
                "Soil moisture estimation",
                "Pest and disease detection",
                "Irrigation optimization",
                "Variable rate application mapping",
                "Field boundary delineation",
                "Crop type classification",
                "Harvest timing optimization",
                "Carbon sequestration monitoring"
            ],
            "data_sources": [
                "Landsat 8/9 (30m resolution)",
                "Sentinel-2 (10m resolution)",
                "MODIS (250m resolution)",
                "Drone/UAV imagery (cm resolution)",
                "Weather station data",
                "Soil sensor networks",
                "Yield monitor data",
                "Farm management systems"
            ],
            "aws_instance_recommendations": {
                "farm_analysis": {
                    "instance_type": "m6i.xlarge",
                    "vcpus": 4,
                    "memory_gb": 16,
                    "storage": "500GB gp3 SSD",
                    "cost_per_hour": 0.192,
                    "use_case": "Individual farm analysis, small datasets"
                },
                "regional_monitoring": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage": "1TB gp3 SSD",
                    "cost_per_hour": 0.504,
                    "use_case": "Regional crop monitoring, medium datasets"
                },
                "agtech_platform": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage": "2TB gp3 SSD + RDS",
                    "cost_per_hour": 0.816,
                    "use_case": "AgTech platform, real-time monitoring"
                }
            },
            "cost_profile": {
                "individual_farm": "$150-400/month",
                "agricultural_consultant": "$400-1000/month",
                "agtech_company": "$1000-3000/month"
            }
        }

    def _get_disaster_config(self) -> Dict[str, Any]:
        """Disaster response and emergency management"""
        return {
            "name": "Disaster Response & Emergency Management System",
            "description": "Real-time disaster monitoring, damage assessment, and emergency response",
            "spack_packages": [
                # Real-time processing
                "kafka@2.13-3.5.0 %gcc@11.4.0",
                "py-kafka@2.0.2 %gcc@11.4.0",

                # Core GIS
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql",
                "gdal@3.7.1 %gcc@11.4.0 +python +netcdf +geos",
                "grass@8.3.0 %gcc@11.4.0 +python +postgresql",

                # Python ecosystem
                "python@3.11.5 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",

                # Remote sensing
                "py-satpy@0.43.0 %gcc@11.4.0",
                "py-pyresample@1.27.1 %gcc@11.4.0",

                # Machine learning for change detection
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-opencv@4.8.0 %gcc@11.4.0",

                # Web mapping and visualization
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python",
                "postgis@3.3.3 %gcc@11.4.0",
                "redis@7.2.0 %gcc@11.4.0"
            ],
            "disaster_applications": [
                "Flood mapping and monitoring",
                "Wildfire detection and tracking",
                "Earthquake damage assessment",
                "Hurricane/typhoon impact analysis",
                "Landslide susceptibility mapping",
                "Tsunami inundation modeling",
                "Infrastructure damage assessment",
                "Evacuation route planning",
                "Emergency resource allocation",
                "Real-time hazard monitoring"
            ],
            "real_time_capabilities": [
                "Satellite imagery processing (< 1 hour)",
                "Change detection algorithms",
                "Automated damage assessment",
                "Emergency alert systems",
                "Mobile data collection",
                "Crowdsourced data integration",
                "Social media monitoring",
                "Multi-source data fusion"
            ],
            "aws_services_integration": [
                "Kinesis for real-time data streams",
                "Lambda for automated processing",
                "SNS for emergency notifications",
                "SQS for message queuing",
                "CloudWatch for monitoring",
                "S3 for data archiving",
                "Ground Station for satellite downlinks"
            ],
            "cost_profile": {
                "research_project": "$500-1200/month",
                "emergency_agency": "$1200-3000/month",
                "operational_system": "$3000-8000/month"
            }
        }

    def _get_urban_planning_config(self) -> Dict[str, Any]:
        """Urban planning and smart city analysis"""
        return {
            "name": "Urban Planning & Smart City Analytics Platform",
            "description": "Urban analysis, city planning, and smart city data integration",
            "spack_packages": [
                # Core GIS and analysis
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql +grass",
                "grass@8.3.0 %gcc@11.4.0 +python +postgresql +netcdf",
                "gdal@3.7.1 %gcc@11.4.0 +python +postgresql +geos",

                # Python ecosystem
                "python@3.11.5 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",

                # Urban analysis
                "py-osmnx@1.6.0 %gcc@11.4.0",  # Street network analysis
                "py-momepy@0.6.0 %gcc@11.4.0",  # Urban morphology
                "py-urbanaccess@0.2.2 %gcc@11.4.0",  # Accessibility analysis

                # Network analysis
                "py-networkx@3.1 %gcc@11.4.0",
                "py-igraph@0.10.6 %gcc@11.4.0",

                # Machine learning
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python",
                "postgis@3.3.3 %gcc@11.4.0"
            ],
            "urban_applications": [
                "Land use planning and zoning",
                "Transportation network analysis",
                "Population density modeling",
                "Urban heat island analysis",
                "Green space accessibility",
                "Walkability and bikeability assessment",
                "Housing affordability analysis",
                "Infrastructure capacity planning",
                "Economic development analysis",
                "Environmental justice mapping"
            ],
            "data_sources": [
                "Census and demographic data",
                "OpenStreetMap data",
                "Satellite imagery",
                "Mobile phone data",
                "Social media data",
                "IoT sensor networks",
                "Municipal databases",
                "Real estate data"
            ],
            "cost_profile": {
                "city_planning_dept": "$800-2000/month",
                "consulting_firm": "$1500-4000/month",
                "smart_city_platform": "$3000-8000/month"
            }
        }

    def _get_maritime_config(self) -> Dict[str, Any]:
        """Maritime and coastal zone analysis"""
        return {
            "name": "Maritime & Coastal Zone Analysis Platform",
            "description": "Ocean monitoring, coastal management, and maritime spatial planning",
            "spack_packages": [
                # Marine data processing
                "nco@5.1.6 %gcc@11.4.0 +netcdf +openmpi",
                "cdo@2.2.0 %gcc@11.4.0 +netcdf +hdf5",
                "ncview@2.1.8 %gcc@11.4.0 +netcdf",

                # Core GIS
                "qgis@3.32.2 %gcc@11.4.0 +python +postgresql",
                "gdal@3.7.1 %gcc@11.4.0 +python +netcdf +geos",
                "grass@8.3.0 %gcc@11.4.0 +python +postgresql",

                # Python ecosystem
                "python@3.11.5 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",

                # Oceanographic analysis
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-gsw@3.6.16 %gcc@11.4.0",  # Gibbs SeaWater

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python",
                "postgis@3.3.3 %gcc@11.4.0"
            ],
            "maritime_applications": [
                "Sea level rise modeling",
                "Coastal erosion analysis",
                "Marine protected area planning",
                "Shipping route optimization",
                "Fisheries stock assessment",
                "Ocean current analysis",
                "Tsunami risk assessment",
                "Offshore wind farm planning",
                "Marine pollution tracking",
                "Habitat suitability modeling"
            ],
            "data_sources": [
                "NOAA/NCEI oceanographic data",
                "COPERNICUS Marine Service",
                "NASA Ocean Color data",
                "GEBCO bathymetry",
                "AIS ship tracking data",
                "Tide gauge measurements",
                "Satellite altimetry",
                "Coastal imagery"
            ],
            "cost_profile": {
                "coastal_research": "$400-1000/month",
                "marine_agency": "$1000-2500/month",
                "maritime_industry": "$2000-5000/month"
            }
        }

    def _get_ml_platform_config(self) -> Dict[str, Any]:
        """Geospatial machine learning platform"""
        return {
            "name": "Geospatial Machine Learning Platform",
            "description": "AI/ML for geospatial analysis, computer vision, and spatial prediction",
            "spack_packages": [
                # Core ML frameworks
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-keras@2.13.1 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",

                # Computer vision
                "py-opencv@4.8.0 %gcc@11.4.0",
                "py-scikit-image@0.21.0 %gcc@11.4.0",
                "py-pillow@10.0.0 %gcc@11.4.0",

                # Deep learning for geospatial
                "py-segmentation-models@1.0.1 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",

                # GPU computing
                "cuda@11.8.0 %gcc@11.4.0",
                "cudnn@8.9.3.28-11.8 %gcc@11.4.0",
                "py-cupy@12.2.0 %gcc@11.4.0",

                # ML utilities
                "py-mlflow@2.5.0 %gcc@11.4.0",
                "py-wandb@0.15.8 %gcc@11.4.0",
                "py-optuna@3.3.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0"
            ],
            "ml_applications": [
                "Land cover classification",
                "Object detection in satellite imagery",
                "Semantic segmentation",
                "Change detection",
                "Crop yield prediction",
                "Forest fire risk modeling",
                "Urban growth prediction",
                "Species distribution modeling",
                "Disaster damage assessment",
                "Real-time monitoring systems"
            ],
            "aws_ml_services": [
                "SageMaker for model training/deployment",
                "Rekognition for image analysis",
                "Comprehend for text analysis",
                "Forecast for time series prediction",
                "Personalize for recommendation systems",
                "Ground Truth for data labeling"
            ],
            "aws_instance_recommendations": {
                "model_development": {
                    "instance_type": "p3.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 61,
                    "gpu": "1x NVIDIA V100",
                    "cost_per_hour": 3.06,
                    "use_case": "Model development and training"
                },
                "large_scale_training": {
                    "instance_type": "p4d.24xlarge",
                    "vcpus": 96,
                    "memory_gb": 1152,
                    "gpu": "8x NVIDIA A100",
                    "cost_per_hour": 32.772,
                    "use_case": "Large-scale model training"
                },
                "inference": {
                    "instance_type": "g4dn.xlarge",
                    "vcpus": 4,
                    "memory_gb": 16,
                    "gpu": "1x NVIDIA T4",
                    "cost_per_hour": 0.526,
                    "use_case": "Model inference and deployment"
                }
            },
            "cost_profile": {
                "research_project": "$500-2000/month",
                "commercial_ml": "$2000-8000/month",
                "enterprise_platform": "$8000-25000/month"
            }
        }

    def generate_deployment_recommendation(self, workload: GeospatialWorkload) -> Dict[str, Any]:
        """Generate deployment recommendation based on workload characteristics"""

        # Select appropriate configuration based on domain
        if workload.domain == GeospatialDomain.REMOTE_SENSING:
            config = self.geospatial_configurations["remote_sensing_lab"]
        elif workload.domain == GeospatialDomain.GIS_ANALYSIS:
            config = self.geospatial_configurations["gis_analysis_studio"]
        elif workload.domain == GeospatialDomain.GEOPHYSICS:
            config = self.geospatial_configurations["geophysics_workstation"]
        elif workload.domain == GeospatialDomain.ENVIRONMENTAL_MODELING:
            config = self.geospatial_configurations["environmental_modeling"]
        elif workload.domain == GeospatialDomain.PRECISION_AGRICULTURE:
            config = self.geospatial_configurations["precision_agriculture"]
        elif workload.domain == GeospatialDomain.DISASTER_RESPONSE:
            config = self.geospatial_configurations["disaster_response"]
        elif workload.domain == GeospatialDomain.URBAN_PLANNING:
            config = self.geospatial_configurations["urban_planning"]
        elif workload.domain == GeospatialDomain.MARITIME_COASTAL:
            config = self.geospatial_configurations["maritime_coastal"]
        else:
            config = self.geospatial_configurations["gis_analysis_studio"]  # Default

        # Adjust instance recommendation based on data volume and complexity
        instance_recommendations = config["aws_instance_recommendations"]

        if workload.data_volume_tb < 0.1:
            recommended_instance = list(instance_recommendations.keys())[0]
        elif workload.data_volume_tb < 1.0:
            recommended_instance = list(instance_recommendations.keys())[min(1, len(instance_recommendations) - 1)]
        else:
            recommended_instance = list(instance_recommendations.keys())[-1]

        return {
            "configuration": config,
            "workload": workload,
            "recommended_instance": recommended_instance,
            "instance_details": instance_recommendations[recommended_instance],
            "estimated_monthly_cost": self._estimate_monthly_cost(workload, instance_recommendations[recommended_instance]),
            "deployment_timeline": "1-3 hours automated setup",
            "optimization_recommendations": self._get_optimization_recommendations(workload)
        }

    def _estimate_monthly_cost(self, workload: GeospatialWorkload, instance_config: Dict[str, Any]) -> Dict[str, float]:
        """Estimate monthly costs based on workload and instance configuration"""

        # Base compute cost (assuming 8 hours/day, 20 days/month for typical research)
        base_hours = 160  # 8 hours * 20 days

        if workload.temporal_frequency == "Real-time":
            usage_multiplier = 3.0  # 24/7 operation
        elif workload.temporal_frequency == "Daily":
            usage_multiplier = 1.5
        elif workload.temporal_frequency == "Weekly":
            usage_multiplier = 1.0
        else:  # Monthly, Annual
            usage_multiplier = 0.5

        compute_hours = base_hours * usage_multiplier
        compute_cost = compute_hours * instance_config.get("cost_per_hour", 1.0)

        # Storage cost (based on data volume)
        storage_cost = workload.data_volume_tb * 1000 * 0.08  # $0.08/GB-month for gp3

        # Data transfer cost (estimated)
        transfer_cost = workload.data_volume_tb * 50  # $0.05/GB estimate

        return {
            "compute": compute_cost,
            "storage": storage_cost,
            "data_transfer": transfer_cost,
            "total": compute_cost + storage_cost + transfer_cost
        }

    def _get_optimization_recommendations(self, workload: GeospatialWorkload) -> List[str]:
        """Get optimization recommendations based on workload"""
        recommendations = []

        if workload.data_volume_tb > 1.0:
            recommendations.append("Consider using S3 Intelligent Tiering for cost optimization")
            recommendations.append("Use CloudFront for global data distribution")

        if workload.temporal_frequency == "Real-time":
            recommendations.append("Implement auto-scaling for variable workloads")
            recommendations.append("Use spot instances for batch processing")

        if workload.coverage_area == "Global":
            recommendations.append("Deploy in multiple AWS regions for better performance")
            recommendations.append("Use AWS Global Accelerator for improved connectivity")

        if workload.analysis_complexity == "Advanced":
            recommendations.append("Consider GPU instances for ML/AI workloads")
            recommendations.append("Implement distributed processing with EMR or Batch")

        return recommendations

    def list_configurations(self) -> List[str]:
        """List all available geospatial configurations"""
        return list(self.geospatial_configurations.keys())

    def get_configuration_details(self, config_name: str) -> Dict[str, Any]:
        """Get detailed configuration information"""
        if config_name not in self.geospatial_configurations:
            raise ValueError(f"Configuration '{config_name}' not found")
        return self.geospatial_configurations[config_name]

def main():
    """CLI interface for geospatial research pack"""
    import argparse

    parser = argparse.ArgumentParser(description="AWS Research Wizard - Geospatial Research Pack")
    parser.add_argument("--list", action="store_true", help="List available configurations")
    parser.add_argument("--config", type=str, help="Show configuration details")
    parser.add_argument("--domain", type=str, choices=[d.value for d in GeospatialDomain],
                       help="Geospatial domain")
    parser.add_argument("--data-volume", type=float, default=1.0, help="Data volume in TB")
    parser.add_argument("--coverage", type=str, choices=["Local", "Regional", "National", "Global"],
                       default="Regional", help="Coverage area")
    parser.add_argument("--frequency", type=str, choices=["Real-time", "Daily", "Weekly", "Monthly", "Annual"],
                       default="Weekly", help="Temporal frequency")
    parser.add_argument("--output", type=str, help="Output file for recommendation")

    args = parser.parse_args()

    geo_pack = GeospatialResearchPack()

    if args.list:
        print("Available Geospatial Configurations:")
        for config_name in geo_pack.list_configurations():
            config = geo_pack.get_configuration_details(config_name)
            print(f"  {config_name}: {config['description']}")

    elif args.config:
        try:
            config = geo_pack.get_configuration_details(args.config)
            print(json.dumps(config, indent=2))
        except ValueError as e:
            print(f"Error: {e}")

    elif args.domain:
        # Generate recommendation based on workload
        workload = GeospatialWorkload(
            domain=GeospatialDomain(args.domain),
            data_sources=["satellite", "ground_truth"],
            processing_type="analysis",
            spatial_resolution="Medium",
            temporal_frequency=args.frequency,
            coverage_area=args.coverage,
            data_volume_tb=args.data_volume,
            analysis_complexity="Moderate"
        )

        recommendation = geo_pack.generate_deployment_recommendation(workload)

        if args.output:
            with open(args.output, 'w') as f:
                json.dump(recommendation, f, indent=2)
            print(f"Recommendation saved to {args.output}")
        else:
            print(json.dumps(recommendation, indent=2))

    else:
        parser.print_help()

if __name__ == "__main__":
    main()
