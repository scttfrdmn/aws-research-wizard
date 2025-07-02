#!/usr/bin/env python3
"""
Atmospheric Chemistry Research Pack
Comprehensive atmospheric chemistry modeling, air quality analysis, and chemical transport modeling for AWS Research Wizard
"""

import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from enum import Enum

class AtmosphericDomain(Enum):
    CHEMICAL_TRANSPORT = "chemical_transport"
    AIR_QUALITY = "air_quality"
    GREENHOUSE_GASES = "greenhouse_gases"
    ATMOSPHERIC_COMPOSITION = "atmospheric_composition"
    OZONE_CHEMISTRY = "ozone_chemistry"
    AEROSOL_CHEMISTRY = "aerosol_chemistry"
    EMISSION_INVENTORY = "emission_inventory"
    CLIMATE_CHEMISTRY = "climate_chemistry"

@dataclass
class AtmosphericWorkload:
    """Atmospheric chemistry workload characteristics"""
    domain: AtmosphericDomain
    model_type: str          # GEOS-Chem, CMAQ, WRF-Chem, etc.
    spatial_resolution: str  # Global, Regional, Urban, Local
    temporal_scale: str      # Real-time, Daily, Monthly, Annual, Climate
    chemical_complexity: str # Basic, Standard, Full, Custom
    emission_sources: List[str]  # Anthropogenic, biogenic, natural, etc.
    data_volume_tb: float    # Expected data volume
    computational_intensity: str  # Light, Moderate, Intensive, Extreme

class AtmosphericChemistryPack:
    """
    Comprehensive atmospheric chemistry research environments optimized for AWS
    Supports chemical transport modeling, air quality analysis, and atmospheric composition studies
    """

    def __init__(self):
        self.atmospheric_configurations = {
            "geos_chem_global": self._get_geos_chem_config(),
            "air_quality_modeling": self._get_air_quality_config(),
            "chemical_transport": self._get_chemical_transport_config(),
            "atmospheric_composition": self._get_composition_config(),
            "greenhouse_gas_analysis": self._get_ghg_config(),
            "ozone_chemistry": self._get_ozone_config(),
            "aerosol_modeling": self._get_aerosol_config(),
            "emission_processing": self._get_emission_config(),
            "atmospheric_ml_platform": self._get_ml_atmospheric_config()
        }

    def _get_geos_chem_config(self) -> Dict[str, Any]:
        """GEOS-Chem global chemical transport modeling"""
        return {
            "name": "GEOS-Chem Global Chemical Transport Modeling",
            "description": "Global 3D chemical transport modeling with GEOS-Chem",
            "spack_packages": [
                # GEOS-Chem and dependencies
                "geos-chem@14.2.0 %gcc@11.4.0 +openmp +netcdf +kpp +rrtmg",
                "gcpy@1.4.0 %gcc@11.4.0 +python",  # GEOS-Chem Python toolkit
                "hemco@3.6.0 %gcc@11.4.0 +netcdf +openmp",  # Emissions component

                # Core atmospheric chemistry libraries
                "netcdf-c@4.9.2 %gcc@11.4.0 +mpi +parallel-netcdf +hdf5",
                "netcdf-fortran@4.6.1 %gcc@11.4.0",
                "hdf5@1.14.2 %gcc@11.4.0 +mpi +threadsafe +fortran",
                "parallel-netcdf@1.12.3 %gcc@11.4.0",

                # Chemical kinetics and photolysis
                "kpp@2.5.0 %gcc@11.4.0 +openmp",  # Kinetic PreProcessor
                "fast-jx@7.0.1 %gcc@11.4.0",     # Photolysis rates
                "rrtmg@1.0 %gcc@11.4.0 +openmp", # Radiative transfer

                # Meteorological data processing
                "nco@5.1.6 %gcc@11.4.0 +netcdf +openmpi",
                "cdo@2.2.0 %gcc@11.4.0 +netcdf +hdf5 +proj +openmp",
                "ncview@2.1.8 %gcc@11.4.0 +netcdf",

                # Python atmospheric chemistry stack
                "python@3.11.5 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",

                # GEOS-Chem specific Python tools
                "py-xbpch@0.3.5 %gcc@11.4.0",    # Binary punch file reader
                "py-gcpy@1.4.0 %gcc@11.4.0",     # GEOS-Chem Python interface
                "py-pymultitau@0.3.3 %gcc@11.4.0", # Chemical lifetime analysis

                # Atmospheric data analysis
                "py-atmospheric-chemistry@2.1.0 %gcc@11.4.0",
                "py-atmos@0.4.0 %gcc@11.4.0",
                "py-chempy@0.8.3 %gcc@11.4.0",

                # AWS-optimized parallel computing
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers +pmix +pmi +fabrics",
                "libfabric@1.18.1 %gcc@11.4.0 +verbs +mlx +efa",  # EFA support
                "aws-ofi-nccl@1.7.0 %gcc@11.4.0",  # AWS OFI plugin for NCCL
                "ucx@1.14.1 %gcc@11.4.0 +verbs +mlx +ib_hw_tm",  # Unified Communication X
                "py-mpi4py@3.1.4 %gcc@11.4.0",

                # Visualization
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-basemap@1.3.8 %gcc@11.4.0",
                "paraview@5.11.2 %gcc@11.4.0 +python +mpi +osmesa"
            ],
            "aws_atmospheric_data": {
                "geos_meteorology": {
                    "description": "GEOS-FP and GEOS-IT meteorological data",
                    "bucket": "s3://noaa-geos-pds/",
                    "resolution": "0.25° x 0.3125°",
                    "frequency": "3-hourly",
                    "coverage": "Global"
                },
                "nasa_atmospheric_composition": {
                    "description": "NASA satellite atmospheric composition data",
                    "bucket": "s3://nasa-atmospheric-composition/",
                    "instruments": ["OMI", "AIRS", "MOPITT", "TROPOMI", "OCO-2"],
                    "species": ["NO2", "O3", "CO", "CH4", "SO2", "HCHO", "CO2"],
                    "coverage": "Global"
                },
                "epa_emissions": {
                    "description": "EPA National Emissions Inventory",
                    "bucket": "s3://epa-nei-data/",
                    "resolution": "US counties and point sources",
                    "frequency": "Annual",
                    "sectors": ["Mobile", "Point", "Area", "Biogenic"]
                },
                "edgar_emissions": {
                    "description": "EDGAR global emission database",
                    "bucket": "s3://edgar-emissions/",
                    "resolution": "0.1° x 0.1°",
                    "frequency": "Annual",
                    "coverage": "Global anthropogenic emissions"
                }
            },
            "model_configurations": {
                "global_standard": {
                    "resolution": "4° x 5°",
                    "levels": 72,
                    "chemistry": "Standard NOx-Ox-HC-aerosol simulation",
                    "runtime": "1 year simulation: 24-48 hours",
                    "memory": "64-128 GB",
                    "cores": "16-32"
                },
                "global_high_res": {
                    "resolution": "2° x 2.5°",
                    "levels": 72,
                    "chemistry": "Full chemistry with detailed VOC",
                    "runtime": "1 year simulation: 3-7 days",
                    "memory": "256-512 GB",
                    "cores": "64-128"
                },
                "nested_regional": {
                    "resolution": "0.25° x 0.3125°",
                    "levels": 72,
                    "chemistry": "Regional nested simulation",
                    "runtime": "1 month simulation: 6-12 hours",
                    "memory": "128-256 GB",
                    "cores": "32-64"
                }
            },
            "aws_instance_recommendations": {
                "global_standard": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage": "2TB gp3 SSD + 10TB EFS",
                    "cost_per_hour": 1.632,
                    "use_case": "Standard global GEOS-Chem simulations"
                },
                "global_high_performance": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage": "4TB gp3 SSD + 20TB EFS",
                    "cost_per_hour": 2.016,
                    "use_case": "High-resolution global simulations"
                },
                "hpc_cluster": {
                    "instance_type": "c6i.16xlarge",
                    "vcpus": 64,
                    "memory_gb": 128,
                    "count": "4-16 nodes",
                    "storage": "FSx Lustre 2.4TB",
                    "cost_per_hour": "$6.528-26.112 (4-16 nodes)",
                    "use_case": "Ensemble runs, sensitivity studies"
                }
            },
            "research_applications": [
                "Global atmospheric chemistry modeling",
                "Source attribution studies",
                "Policy impact assessment",
                "Climate-chemistry interactions",
                "Air quality forecasting",
                "Atmospheric oxidation capacity",
                "Tropospheric ozone budget",
                "Secondary organic aerosol formation",
                "Mercury cycling",
                "Inverse modeling and data assimilation"
            ],
            "cost_profile": {
                "single_simulation": "$50-400 per year simulation",
                "research_project": "$800-3000/month (multiple scenarios)",
                "operational_forecasting": "$3000-8000/month",
                "large_campaigns": "$8000-25000/month (ensemble studies)"
            }
        }

    def _get_air_quality_config(self) -> Dict[str, Any]:
        """Air quality modeling and analysis platform"""
        return {
            "name": "Air Quality Modeling & Analysis Platform",
            "description": "Regional air quality modeling with CMAQ, WRF-Chem, and monitoring data analysis",
            "spack_packages": [
                # Air quality models
                "cmaq@5.4 %gcc@11.4.0 +mpi +netcdf +ioapi",
                "wrf-chem@4.5.0 %gcc@11.4.0 +netcdf +mpi +chem +kpp",
                "cam-chem@6.3.0 %gcc@11.4.0 +netcdf +mpi +chemistry",

                # Pre/post processing tools
                "mcip@5.3.3 %gcc@11.4.0 +netcdf +ioapi",  # Meteorology-Chemistry Interface
                "smoke@4.8.1 %gcc@11.4.0 +netcdf +ioapi", # Emissions processing
                "ioapi@3.2 %gcc@11.4.0 +netcdf +mpi",     # I/O Applications Programming Interface

                # Weather preprocessing
                "wrf@4.5.0 %gcc@11.4.0 +netcdf +hdf5 +mpi +openmp",
                "wps@4.5 %gcc@11.4.0 +netcdf +grib2",

                # Chemical mechanism processing
                "kpp@2.5.0 %gcc@11.4.0 +openmp",
                "chemmech@1.0 %gcc@11.4.0",

                # Data processing
                "netcdf-c@4.9.2 %gcc@11.4.0 +mpi +parallel-netcdf",
                "nco@5.1.6 %gcc@11.4.0 +netcdf +openmpi",
                "cdo@2.2.0 %gcc@11.4.0 +netcdf +hdf5",

                # Python air quality analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",

                # Air quality specific packages
                "py-aerocom@0.11.0 %gcc@11.4.0",  # Aerosol comparison
                "py-pyaerocom@0.16.0 %gcc@11.4.0", # AeroCom data analysis
                "py-monetio@0.2.3 %gcc@11.4.0",   # Model and observation comparison
                "py-airnow@1.0.0 %gcc@11.4.0",    # EPA AirNow data access

                # Statistical analysis
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0"
            ],
            "monitoring_data_sources": {
                "epa_aqs": {
                    "description": "EPA Air Quality System monitoring data",
                    "bucket": "s3://epa-aqs-data/",
                    "parameters": ["PM2.5", "PM10", "O3", "NO2", "SO2", "CO"],
                    "frequency": "Hourly",
                    "coverage": "US monitoring network"
                },
                "airnow": {
                    "description": "Real-time air quality data",
                    "bucket": "s3://epa-airnow/",
                    "parameters": ["AQI", "PM2.5", "O3"],
                    "frequency": "Hourly",
                    "coverage": "US real-time network"
                },
                "european_eea": {
                    "description": "European Environment Agency air quality data",
                    "bucket": "s3://eea-air-quality/",
                    "parameters": ["PM2.5", "PM10", "NO2", "O3", "SO2"],
                    "coverage": "European monitoring networks"
                },
                "satellite_products": {
                    "description": "Satellite-derived air quality products",
                    "bucket": "s3://nasa-air-quality/",
                    "products": ["MODIS AOD", "OMI NO2", "TROPOMI NO2/SO2"],
                    "coverage": "Global"
                }
            },
            "modeling_domains": {
                "urban_scale": {
                    "resolution": "1-4 km",
                    "extent": "Metropolitan areas",
                    "applications": ["City air quality planning", "Hotspot analysis"],
                    "runtime": "1 month: 4-8 hours"
                },
                "regional_scale": {
                    "resolution": "4-12 km",
                    "extent": "State/province level",
                    "applications": ["Regional air quality planning", "Transport studies"],
                    "runtime": "1 year: 1-3 days"
                },
                "continental_scale": {
                    "resolution": "12-36 km",
                    "extent": "Continental domains",
                    "applications": ["Policy assessment", "Long-range transport"],
                    "runtime": "Multi-year: 3-10 days"
                }
            },
            "aws_instance_recommendations": {
                "urban_modeling": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage": "1TB gp3 SSD",
                    "cost_per_hour": 0.816,
                    "use_case": "Urban air quality modeling (1-4km resolution)"
                },
                "regional_modeling": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage": "2TB gp3 SSD + 5TB EFS",
                    "cost_per_hour": 1.632,
                    "use_case": "Regional air quality modeling (4-12km)"
                },
                "operational_forecasting": {
                    "instance_type": "c6i.16xlarge",
                    "vcpus": 64,
                    "memory_gb": 128,
                    "storage": "4TB gp3 SSD + 10TB EFS",
                    "cost_per_hour": 3.264,
                    "use_case": "Daily operational air quality forecasting"
                }
            },
            "cost_profile": {
                "research_study": "$400-1200/month",
                "regulatory_modeling": "$1200-4000/month",
                "operational_forecasting": "$3000-10000/month"
            }
        }

    def _get_chemical_transport_config(self) -> Dict[str, Any]:
        """Chemical transport modeling and atmospheric dispersion"""
        return {
            "name": "Chemical Transport & Atmospheric Dispersion Modeling",
            "description": "Multi-scale chemical transport modeling and atmospheric dispersion analysis",
            "spack_packages": [
                # Chemical transport models
                "flexpart@10.4 %gcc@11.4.0 +netcdf +mpi",  # Lagrangian transport
                "hysplit@5.2.3 %gcc@11.4.0 +netcdf",       # Hybrid transport
                "stilt@2.1.0 %gcc@11.4.0 +netcdf +r",      # Stochastic transport

                # Eulerian models
                "geos-chem@14.2.0 %gcc@11.4.0 +openmp +netcdf",
                "mozart@4.0 %gcc@11.4.0 +netcdf +mpi",
                "cam-chem@6.3.0 %gcc@11.4.0 +netcdf +mpi",

                # Emission and boundary condition processing
                "hemco@3.6.0 %gcc@11.4.0 +netcdf +openmp",
                "prep-chem-src@1.7 %gcc@11.4.0 +netcdf",
                "edgar-tools@4.3.2 %gcc@11.4.0 +python",

                # Meteorological preprocessing
                "ecmwf-tools@2.1.0 %gcc@11.4.0 +netcdf +grib",
                "ncep-tools@1.3.0 %gcc@11.4.0 +netcdf +grib2",

                # Python transport analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-atmospheric-transport@1.2.0 %gcc@11.4.0",
                "py-lagrangian@0.3.0 %gcc@11.4.0",
                "py-trajectory@2.0.0 %gcc@11.4.0",

                # Data analysis
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0"
            ],
            "transport_applications": [
                "Source-receptor modeling",
                "Pollution transport pathways",
                "Emergency response modeling",
                "Inverse modeling for emissions",
                "Long-range transport studies",
                "Boundary layer transport",
                "Volcanic ash dispersion",
                "Radioactive material transport",
                "Bioaerosol dispersion",
                "Pollen transport modeling"
            ],
            "cost_profile": {
                "research_project": "$300-1000/month",
                "emergency_response": "$1000-3000/month",
                "operational_system": "$3000-8000/month"
            }
        }

    def _get_composition_config(self) -> Dict[str, Any]:
        """Atmospheric composition analysis and satellite data"""
        return {
            "name": "Atmospheric Composition Analysis Platform",
            "description": "Satellite atmospheric composition data analysis and validation",
            "spack_packages": [
                # Satellite data processing
                "harp@1.19 %gcc@11.4.0 +python +netcdf",     # HARP atmospheric data processor
                "coda@2.24 %gcc@11.4.0 +python +hdf5",      # CODA data reader
                "s5p-tools@1.7.0 %gcc@11.4.0 +python",      # Sentinel-5P tools

                # Retrieval algorithms
                "omi-tools@3.0.0 %gcc@11.4.0 +python +hdf5", # OMI data processing
                "modis-tools@6.1.0 %gcc@11.4.0 +python",     # MODIS atmospheric products
                "airs-tools@7.0.0 %gcc@11.4.0 +python",      # AIRS retrieval tools

                # Python atmospheric composition
                "python@3.11.5 %gcc@11.4.0",
                "py-satpy@0.43.0 %gcc@11.4.0",
                "py-pyresample@1.27.1 %gcc@11.4.0",
                "py-pyhdf@0.10.5 %gcc@11.4.0",
                "py-h5py@3.9.0 %gcc@11.4.0",

                # Atmospheric composition analysis
                "py-atmospheric-composition@1.4.0 %gcc@11.4.0",
                "py-tropomi@2.1.0 %gcc@11.4.0",
                "py-omi@1.3.0 %gcc@11.4.0",
                "py-satellite-validation@0.8.0 %gcc@11.4.0",

                # Data processing
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-dask@2023.8.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0"
            ],
            "satellite_instruments": {
                "tropomi": {
                    "platform": "Sentinel-5P",
                    "species": ["NO2", "SO2", "CO", "CH4", "HCHO", "O3", "aerosols"],
                    "resolution": "7x7 km",
                    "frequency": "Daily global coverage"
                },
                "omi": {
                    "platform": "Aura",
                    "species": ["NO2", "SO2", "HCHO", "BrO", "O3"],
                    "resolution": "13x24 km",
                    "frequency": "Daily global coverage"
                },
                "airs": {
                    "platform": "Aqua",
                    "species": ["CO2", "CH4", "CO", "O3", "temperature"],
                    "resolution": "13.5 km",
                    "frequency": "Twice daily"
                },
                "mopitt": {
                    "platform": "Terra",
                    "species": ["CO"],
                    "resolution": "22x22 km",
                    "frequency": "Global coverage every 3 days"
                }
            },
            "cost_profile": {
                "satellite_analysis": "$200-800/month",
                "validation_studies": "$800-2000/month",
                "operational_monitoring": "$2000-5000/month"
            }
        }

    def _get_ghg_config(self) -> Dict[str, Any]:
        """Greenhouse gas analysis and carbon cycle modeling"""
        return {
            "name": "Greenhouse Gas Analysis & Carbon Cycle Platform",
            "description": "Greenhouse gas monitoring, carbon cycle modeling, and emission verification",
            "spack_packages": [
                # Carbon cycle models
                "casa-gfed@3.0 %gcc@11.4.0 +netcdf +mpi",   # Carbon cycle model
                "lpj-guess@4.1 %gcc@11.4.0 +netcdf +mpi",   # Vegetation model
                "vegas@2.6 %gcc@11.4.0 +netcdf +mpi",       # Dynamic vegetation

                # Atmospheric transport for GHG
                "geos-chem@14.2.0 %gcc@11.4.0 +openmp +co2", # CO2 simulation
                "flexpart@10.4 %gcc@11.4.0 +netcdf +co2",    # CO2 transport
                "stilt@2.1.0 %gcc@11.4.0 +netcdf +co2",      # CO2 footprints

                # Inverse modeling
                "geos-chem-adjoint@35 %gcc@11.4.0 +netcdf +openmp",
                "tm5-4dvar@3.0 %gcc@11.4.0 +netcdf +mpi",
                "oco-2-tools@3.0 %gcc@11.4.0 +python",

                # Python GHG analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-greenhouse-gas@1.2.0 %gcc@11.4.0",
                "py-carbon-cycle@0.8.0 %gcc@11.4.0",
                "py-co2-analysis@2.1.0 %gcc@11.4.0",
                "py-ch4-analysis@1.4.0 %gcc@11.4.0",

                # Data processing
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",

                # Statistical analysis
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0"
            ],
            "ghg_data_sources": {
                "oco2_oco3": {
                    "description": "NASA OCO-2/OCO-3 CO2 satellite data",
                    "bucket": "s3://nasa-oco-data/",
                    "species": ["CO2"],
                    "resolution": "1.3x2.25 km footprints",
                    "coverage": "Global CO2 column measurements"
                },
                "gosat_data": {
                    "description": "JAXA GOSAT greenhouse gas data",
                    "bucket": "s3://gosat-ghg-data/",
                    "species": ["CO2", "CH4"],
                    "resolution": "10.5 km diameter",
                    "coverage": "Global greenhouse gas columns"
                },
                "noaa_flask": {
                    "description": "NOAA Global Monitoring Laboratory flask data",
                    "bucket": "s3://noaa-gml-data/",
                    "species": ["CO2", "CH4", "N2O", "SF6"],
                    "frequency": "Weekly to monthly",
                    "coverage": "Global surface network"
                },
                "edgar_ghg": {
                    "description": "EDGAR greenhouse gas emissions",
                    "bucket": "s3://edgar-ghg-emissions/",
                    "species": ["CO2", "CH4", "N2O"],
                    "resolution": "0.1° x 0.1°",
                    "coverage": "Global anthropogenic emissions"
                }
            },
            "cost_profile": {
                "carbon_research": "$400-1200/month",
                "emission_verification": "$1200-3000/month",
                "operational_monitoring": "$3000-8000/month"
            }
        }

    def _get_ozone_config(self) -> Dict[str, Any]:
        """Ozone chemistry and stratospheric modeling"""
        return {
            "name": "Ozone Chemistry & Stratospheric Modeling Platform",
            "description": "Stratospheric ozone chemistry, UV radiation, and ozone depletion studies",
            "spack_packages": [
                # Stratospheric chemistry models
                "slimcat@7.4 %gcc@11.4.0 +netcdf +mpi",      # Stratospheric chemistry
                "tomcat@1.7 %gcc@11.4.0 +netcdf +mpi",       # Chemical transport
                "waccm@6.3 %gcc@11.4.0 +netcdf +mpi +chem",  # Whole atmosphere model

                # UV radiation models
                "tuv@5.4 %gcc@11.4.0 +netcdf",               # Tropospheric UV
                "libradtran@2.0.5 %gcc@11.4.0 +netcdf",     # Radiative transfer

                # Ozone chemistry
                "fast-jx@7.0.1 %gcc@11.4.0 +strat",         # Photolysis rates
                "kpp@2.5.0 %gcc@11.4.0 +strat +openmp",     # Chemical kinetics

                # Python ozone analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-ozone-analysis@1.8.0 %gcc@11.4.0",
                "py-stratospheric@0.6.0 %gcc@11.4.0",
                "py-uv-radiation@1.2.0 %gcc@11.4.0",

                # Data processing
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0"
            ],
            "cost_profile": {
                "ozone_research": "$300-1000/month",
                "climate_ozone": "$1000-3000/month",
                "operational_monitoring": "$2000-5000/month"
            }
        }

    def _get_aerosol_config(self) -> Dict[str, Any]:
        """Aerosol chemistry and microphysics modeling"""
        return {
            "name": "Aerosol Chemistry & Microphysics Platform",
            "description": "Atmospheric aerosol modeling, air quality, and climate interactions",
            "spack_packages": [
                # Aerosol models
                "gocart@2.1 %gcc@11.4.0 +netcdf +mpi",       # Aerosol transport
                "matrix@1.0 %gcc@11.4.0 +netcdf +aerosol",   # Aerosol microphysics
                "mosaic@8.2 %gcc@11.4.0 +netcdf +aerosol",   # Sectional aerosol

                # Aerosol chemistry in CTMs
                "geos-chem@14.2.0 %gcc@11.4.0 +aerosol +soa",
                "cmaq@5.4 %gcc@11.4.0 +aerosol +isorropia",

                # Thermodynamic models
                "isorropia@2.3 %gcc@11.4.0",                 # Aerosol thermodynamics
                "mesa@1.0 %gcc@11.4.0",                      # Multicomponent equilibrium

                # Python aerosol analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-aerosol-analysis@2.1.0 %gcc@11.4.0",
                "py-particulate-matter@1.3.0 %gcc@11.4.0",
                "py-aerocom@0.11.0 %gcc@11.4.0",

                # Data processing
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0"
            ],
            "cost_profile": {
                "aerosol_research": "$400-1200/month",
                "air_quality_modeling": "$1200-3000/month",
                "climate_aerosol": "$2000-6000/month"
            }
        }

    def _get_emission_config(self) -> Dict[str, Any]:
        """Emission inventory processing and analysis"""
        return {
            "name": "Emission Inventory Processing & Analysis Platform",
            "description": "Emission inventory development, processing, and uncertainty analysis",
            "spack_packages": [
                # Emission processing tools
                "smoke@4.8.1 %gcc@11.4.0 +netcdf +ioapi",    # US EPA SMOKE
                "prep-chem-src@1.7 %gcc@11.4.0 +netcdf",     # Emission preprocessing
                "hemco@3.6.0 %gcc@11.4.0 +netcdf +openmp",   # Harvard-NASA emissions

                # Inventory tools
                "edgar-tools@4.3.2 %gcc@11.4.0 +python",
                "ceds-tools@2021 %gcc@11.4.0 +python",
                "htap-tools@2.3 %gcc@11.4.0 +python",

                # Python emission analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-emission-inventory@1.5.0 %gcc@11.4.0",
                "py-uncertainty-analysis@0.9.0 %gcc@11.4.0",
                "py-spatial-allocation@1.1.0 %gcc@11.4.0",

                # Data processing
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python",
                "postgis@3.3.3 %gcc@11.4.0"
            ],
            "cost_profile": {
                "inventory_development": "$200-800/month",
                "policy_analysis": "$800-2000/month",
                "operational_processing": "$2000-5000/month"
            }
        }

    def _get_ml_atmospheric_config(self) -> Dict[str, Any]:
        """Machine learning for atmospheric chemistry"""
        return {
            "name": "Atmospheric Chemistry Machine Learning Platform",
            "description": "AI/ML for atmospheric chemistry, air quality prediction, and pattern recognition",
            "spack_packages": [
                # ML frameworks
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",

                # Atmospheric ML
                "py-atmospheric-ml@1.0.0 %gcc@11.4.0",
                "py-air-quality-ml@0.8.0 %gcc@11.4.0",
                "py-chemistry-ml@0.5.0 %gcc@11.4.0",

                # GPU computing
                "cuda@11.8.0 %gcc@11.4.0",
                "py-cupy@12.2.0 %gcc@11.4.0",

                # Time series analysis
                "py-prophet@1.1.4 %gcc@11.4.0",
                "py-statsforecast@1.6.0 %gcc@11.4.0",

                # Data processing
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-dask@2023.8.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0"
            ],
            "ml_applications": [
                "Air quality forecasting",
                "Emission source identification",
                "Chemical mechanism reduction",
                "Satellite data gap filling",
                "Pattern recognition in atmospheric data",
                "Real-time air quality prediction",
                "Chemical species prediction",
                "Atmospheric model emulation"
            ],
            "cost_profile": {
                "ml_development": "$500-2000/month",
                "operational_ml": "$2000-6000/month",
                "enterprise_platform": "$6000-15000/month"
            }
        }

    def generate_atmospheric_recommendation(self, workload: AtmosphericWorkload) -> Dict[str, Any]:
        """Generate deployment recommendation based on atmospheric workload"""

        # Select configuration based on domain
        domain_mapping = {
            AtmosphericDomain.CHEMICAL_TRANSPORT: "chemical_transport",
            AtmosphericDomain.AIR_QUALITY: "air_quality_modeling",
            AtmosphericDomain.GREENHOUSE_GASES: "greenhouse_gas_analysis",
            AtmosphericDomain.ATMOSPHERIC_COMPOSITION: "atmospheric_composition",
            AtmosphericDomain.OZONE_CHEMISTRY: "ozone_chemistry",
            AtmosphericDomain.AEROSOL_CHEMISTRY: "aerosol_modeling",
            AtmosphericDomain.EMISSION_INVENTORY: "emission_processing",
            AtmosphericDomain.CLIMATE_CHEMISTRY: "geos_chem_global"
        }

        config_name = domain_mapping.get(workload.domain, "geos_chem_global")
        config = self.atmospheric_configurations[config_name]

        return {
            "configuration": config,
            "workload": workload,
            "estimated_cost": self._estimate_atmospheric_cost(workload),
            "deployment_timeline": "2-6 hours automated setup",
            "optimization_recommendations": self._get_atmospheric_optimization(workload)
        }

    def _estimate_atmospheric_cost(self, workload: AtmosphericWorkload) -> Dict[str, float]:
        """Estimate costs for atmospheric chemistry workloads"""

        # Base cost factors
        if workload.computational_intensity == "Light":
            base_cost = 200
        elif workload.computational_intensity == "Moderate":
            base_cost = 800
        elif workload.computational_intensity == "Intensive":
            base_cost = 2000
        else:  # Extreme
            base_cost = 5000

        # Scale by temporal requirements
        if workload.temporal_scale == "Real-time":
            temporal_multiplier = 2.0
        elif workload.temporal_scale == "Daily":
            temporal_multiplier = 1.5
        else:
            temporal_multiplier = 1.0

        # Storage costs
        storage_cost = workload.data_volume_tb * 80  # $80/TB/month

        total_cost = base_cost * temporal_multiplier + storage_cost

        return {
            "compute": base_cost * temporal_multiplier,
            "storage": storage_cost,
            "total": total_cost
        }

    def _get_atmospheric_optimization(self, workload: AtmosphericWorkload) -> List[str]:
        """Get optimization recommendations for atmospheric workloads"""
        recommendations = []

        if workload.computational_intensity == "Intensive":
            recommendations.append("Consider spot instances for 60-80% cost savings")
            recommendations.append("Use auto-scaling for variable computational demands")

        if workload.data_volume_tb > 5.0:
            recommendations.append("Implement intelligent data tiering (S3 IA/Glacier)")
            recommendations.append("Use data compression for atmospheric model output")

        if workload.temporal_scale == "Real-time":
            recommendations.append("Deploy across multiple regions for redundancy")
            recommendations.append("Use AWS Lambda for preprocessing tasks")

        return recommendations

def main():
    """CLI interface for atmospheric chemistry pack"""
    import argparse

    parser = argparse.ArgumentParser(description="AWS Research Wizard - Atmospheric Chemistry Pack")
    parser.add_argument("--list", action="store_true", help="List available configurations")
    parser.add_argument("--config", type=str, help="Show configuration details")
    parser.add_argument("--domain", type=str, choices=[d.value for d in AtmosphericDomain],
                       help="Atmospheric domain")
    parser.add_argument("--model", type=str, default="GEOS-Chem", help="Model type")
    parser.add_argument("--resolution", type=str, choices=["Global", "Regional", "Urban", "Local"],
                       default="Regional", help="Spatial resolution")
    parser.add_argument("--output", type=str, help="Output file for recommendation")

    args = parser.parse_args()

    atmo_pack = AtmosphericChemistryPack()

    if args.list:
        print("Available Atmospheric Chemistry Configurations:")
        for config_name in atmo_pack.atmospheric_configurations.keys():
            config = atmo_pack.atmospheric_configurations[config_name]
            print(f"  {config_name}: {config['description']}")

    elif args.config:
        if args.config in atmo_pack.atmospheric_configurations:
            config = atmo_pack.atmospheric_configurations[args.config]
            print(json.dumps(config, indent=2))
        else:
            print(f"Configuration '{args.config}' not found")

    elif args.domain:
        workload = AtmosphericWorkload(
            domain=AtmosphericDomain(args.domain),
            model_type=args.model,
            spatial_resolution=args.resolution,
            temporal_scale="Daily",
            chemical_complexity="Standard",
            emission_sources=["anthropogenic", "biogenic"],
            data_volume_tb=1.0,
            computational_intensity="Moderate"
        )

        recommendation = atmo_pack.generate_atmospheric_recommendation(workload)

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
