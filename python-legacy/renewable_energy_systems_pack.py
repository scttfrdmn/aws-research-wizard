#!/usr/bin/env python3
"""
Renewable Energy Systems Research Pack
Comprehensive renewable energy research, grid integration, and energy systems modeling for AWS Research Wizard
"""

import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from enum import Enum

class RenewableEnergyDomain(Enum):
    SOLAR_ENERGY = "solar_energy"
    WIND_ENERGY = "wind_energy"
    ENERGY_STORAGE = "energy_storage"
    GRID_INTEGRATION = "grid_integration"
    HYDROELECTRIC = "hydroelectric"
    GEOTHERMAL = "geothermal"
    BIOMASS_BIOENERGY = "biomass_bioenergy"
    ENERGY_SYSTEMS = "energy_systems"
    SMART_GRID = "smart_grid"

@dataclass
class RenewableEnergyWorkload:
    """Renewable energy research workload characteristics"""
    domain: RenewableEnergyDomain
    analysis_type: str       # Resource Assessment, System Design, Economic Analysis, Grid Impact
    temporal_scale: str      # Real-time, Hourly, Daily, Monthly, Annual, Lifetime
    spatial_scale: str       # Site, Regional, National, Global
    modeling_complexity: str # Basic, Intermediate, Advanced, Comprehensive
    data_sources: List[str]  # Weather, Market, Grid, Satellite, Sensor
    optimization_focus: str  # Cost, Performance, Reliability, Sustainability
    data_volume_tb: float    # Expected data volume

class RenewableEnergyResearchPack:
    """
    Comprehensive renewable energy research environments optimized for AWS
    Supports energy systems modeling, resource assessment, and grid integration studies
    """

    def __init__(self):
        self.renewable_energy_configurations = {
            "solar_energy_research": self._get_solar_energy_config(),
            "wind_energy_research": self._get_wind_energy_config(),
            "energy_storage_systems": self._get_energy_storage_config(),
            "grid_integration_platform": self._get_grid_integration_config(),
            "hydroelectric_systems": self._get_hydroelectric_config(),
            "geothermal_research": self._get_geothermal_config(),
            "biomass_bioenergy": self._get_biomass_config(),
            "energy_systems_modeling": self._get_energy_systems_config(),
            "smart_grid_analytics": self._get_smart_grid_config()
        }

    def _get_solar_energy_config(self) -> Dict[str, Any]:
        """Solar energy research and photovoltaic systems modeling"""
        return {
            "name": "Solar Energy Research & Photovoltaic Systems Platform",
            "description": "Comprehensive solar energy research, PV modeling, and solar resource assessment",
            "spack_packages": [
                # Solar energy modeling tools
                "sam@2023.12.17 %gcc@11.4.0 +python",  # System Advisor Model
                "pvlib@0.10.2 %gcc@11.4.0 +python",    # Photovoltaic modeling
                "solarsim@3.0.0 %gcc@11.4.0 +python +netcdf",  # Solar simulation
                "heliostat@2.3.0 %gcc@11.4.0 +python", # Concentrated solar power

                # Solar irradiance and weather data
                "nrel-psm@3.0.0 %gcc@11.4.0 +python",  # Physical Solar Model
                "solargis@2.1.0 %gcc@11.4.0 +python",  # Solar resource data
                "meteonorm@8.1 %gcc@11.4.0 +python",   # Meteorological database

                # PV device modeling
                "pc1d@6.2 %gcc@11.4.0",                # Solar cell modeling
                "scaps@3.3.09 %gcc@11.4.0",            # Solar cell capacitance simulator
                "gpvdm@7.0.15 %gcc@11.4.0 +openmp",    # General-purpose PV device model

                # Python solar analysis stack
                "python@3.11.5 %gcc@11.4.0",
                "py-pvlib@0.10.2 %gcc@11.4.0",
                "py-pysam@4.2.0 %gcc@11.4.0",          # SAM Python wrapper
                "py-solar-data-tools@1.1.0 %gcc@11.4.0", # Solar data analysis
                "py-rdtools@2.1.6 %gcc@11.4.0",        # Renewable energy data analysis

                # Satellite data processing
                "py-satpy@0.43.0 %gcc@11.4.0",
                "py-pyresample@1.27.1 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",

                # Optimization and machine learning
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",

                # Data analysis and visualization
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",

                # Geographic analysis
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",

                # Time series analysis
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-prophet@1.1.4 %gcc@11.4.0",

                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +python",
                "timescaledb@2.11.2 %gcc@11.4.0",      # Time-series database
                "influxdb@2.7.1 %gcc@11.4.0"           # Time-series database
            ],
            "solar_data_sources": {
                "nrel_nsrdb": {
                    "description": "National Solar Radiation Database",
                    "bucket": "s3://nrel-pds-nsrdb/",
                    "resolution": "4km spatial, hourly temporal",
                    "coverage": "Americas, Asia, Australia",
                    "parameters": ["GHI", "DNI", "DHI", "Temperature", "Wind"]
                },
                "solargis_data": {
                    "description": "Global solar irradiance database",
                    "bucket": "s3://solargis-aws/",
                    "resolution": "1-10km spatial, 15min-hourly temporal",
                    "coverage": "Global",
                    "parameters": ["Solar irradiance", "PV power output", "Weather"]
                },
                "satellite_irradiance": {
                    "description": "Satellite-derived solar irradiance",
                    "sources": ["GOES", "MSG", "Himawari"],
                    "resolution": "1-4km spatial, 15min temporal",
                    "coverage": "Regional coverage by satellite"
                },
                "ground_measurements": {
                    "description": "Ground-based solar measurement networks",
                    "networks": ["BSRN", "SURFRAD", "SOLRAD"],
                    "quality": "Research grade",
                    "availability": "Site-specific, long-term records"
                }
            },
            "research_capabilities": [
                "Solar resource assessment and mapping",
                "Photovoltaic system performance modeling",
                "Concentrated solar power (CSP) analysis",
                "Solar tracking system optimization",
                "PV degradation and reliability analysis",
                "Solar forecasting and variability studies",
                "Agrivoltaics and dual land use",
                "Building-integrated photovoltaics (BIPV)",
                "Floating solar (floatovoltaics) analysis",
                "Solar energy economics and policy analysis"
            ],
            "aws_instance_recommendations": {
                "solar_analysis": {
                    "instance_type": "m6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 32,
                    "storage": "1TB gp3 SSD",
                    "cost_per_hour": 0.384,
                    "use_case": "Individual solar energy research projects"
                },
                "large_scale_assessment": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage": "2TB gp3 SSD + 5TB EFS",
                    "cost_per_hour": 1.008,
                    "use_case": "Regional solar resource assessment"
                },
                "ml_solar_forecasting": {
                    "instance_type": "p3.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 61,
                    "gpu": "1x NVIDIA V100",
                    "cost_per_hour": 3.06,
                    "use_case": "Machine learning for solar forecasting"
                }
            ],
            "cost_profile": {
                "academic_research": "$300-1000/month",
                "solar_developer": "$1000-3000/month",
                "utility_planning": "$2000-6000/month",
                "research_institution": "$3000-10000/month"
            }
        }

    def _get_wind_energy_config(self) -> Dict[str, Any]:
        """Wind energy research and wind resource assessment"""
        return {
            "name": "Wind Energy Research & Resource Assessment Platform",
            "description": "Comprehensive wind energy research, turbine modeling, and wind resource analysis",
            "spack_packages": [
                # Wind energy modeling
                "openfast@3.5.0 %gcc@11.4.0 +openmp +netcdf",  # Wind turbine simulation
                "windse@1.0.0 %gcc@11.4.0 +openmp +petsc",     # Wind farm optimization
                "sowfa@2.6.0 %gcc@11.4.0 +openmp +openfoam",   # Wind farm LES
                "wrf@4.5.0 %gcc@11.4.0 +netcdf +hdf5 +mpi +wind", # Weather modeling for wind

                # Wind resource analysis
                "windtoolkit@1.0.0 %gcc@11.4.0 +python +netcdf", # NREL Wind Toolkit
                "wasp@12.7 %gcc@11.4.0 +fortran",              # Wind Atlas Analysis
                "meteodyn@6.1 %gcc@11.4.0 +mpi",               # CFD wind modeling

                # Computational fluid dynamics
                "openfoam@11 %gcc@11.4.0 +mpi +metis +scotch", # CFD simulation
                "su2@7.5.1 %gcc@11.4.0 +mpi +openmp",          # CFD suite
                "fenics@2019.1.0 %gcc@11.4.0 +mpi +petsc",     # Finite element

                # Python wind analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-windpowerlib@0.2.1 %gcc@11.4.0",           # Wind power calculations
                "py-wind-toolkit@1.0.0 %gcc@11.4.0",           # NREL Wind Toolkit tools
                "py-floris@3.4.0 %gcc@11.4.0",                 # Wake modeling
                "py-pyopenfast@1.0.0 %gcc@11.4.0",            # OpenFAST Python wrapper

                # Atmospheric modeling
                "py-metpy@1.5.1 %gcc@11.4.0",                  # Meteorological analysis
                "py-wrf-python@1.3.4 %gcc@11.4.0",            # WRF data analysis
                "py-atmospheric-models@1.2.0 %gcc@11.4.0",

                # Time series and forecasting
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-prophet@1.1.4 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",

                # Data processing
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-cartopy@0.21.1 %gcc@11.4.0",
                "paraview@5.11.2 %gcc@11.4.0 +python +mpi",

                # Geographic analysis
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",

                # Parallel computing
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers",
                "petsc@3.19.4 %gcc@11.4.0 +mpi +hypre"
            ],
            "wind_data_sources": {
                "nrel_wind_toolkit": {
                    "description": "NREL Wind Integration National Dataset",
                    "bucket": "s3://nrel-pds-wtk/",
                    "resolution": "2km spatial, 5min temporal",
                    "coverage": "Continental United States",
                    "parameters": ["Wind speed", "Wind direction", "Temperature", "Pressure"]
                },
                "era5_reanalysis": {
                    "description": "ERA5 atmospheric reanalysis",
                    "bucket": "s3://era5-pds/",
                    "resolution": "31km spatial, hourly temporal",
                    "coverage": "Global",
                    "parameters": ["10m wind", "100m wind", "Surface pressure", "Temperature"]
                },
                "merra2_wind": {
                    "description": "NASA MERRA-2 wind data",
                    "bucket": "s3://nasa-merra2/",
                    "resolution": "50km spatial, hourly temporal",
                    "coverage": "Global",
                    "parameters": ["Wind speed", "Wind direction", "Height profiles"]
                },
                "global_wind_atlas": {
                    "description": "DTU Global Wind Atlas",
                    "source": "Technical University of Denmark",
                    "resolution": "250m-1km spatial",
                    "coverage": "Global wind resource maps",
                    "parameters": ["Wind speed", "Wind power density", "Weibull parameters"]
                }
            },
            "research_capabilities": [
                "Wind resource assessment and mapping",
                "Wind turbine aerodynamic analysis",
                "Wind farm layout optimization",
                "Wake effect modeling and mitigation",
                "Wind forecasting and variability studies",
                "Offshore wind resource analysis",
                "Extreme wind event analysis",
                "Wind-grid integration studies",
                "Turbine structural analysis",
                "Wind energy economics and policy"
            ],
            "aws_instance_recommendations": {
                "wind_analysis": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage": "2TB gp3 SSD",
                    "cost_per_hour": 0.816,
                    "use_case": "Wind resource analysis and turbine modeling"
                },
                "cfd_simulation": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage": "4TB gp3 SSD",
                    "cost_per_hour": 1.632,
                    "use_case": "Computational fluid dynamics simulations"
                },
                "hpc_wind_cluster": {
                    "instance_type": "c6i.16xlarge",
                    "vcpus": 64,
                    "memory_gb": 128,
                    "count": "4-16 nodes",
                    "storage": "FSx Lustre 2.4TB",
                    "cost_per_hour": "$6.528-26.112 (4-16 nodes)",
                    "use_case": "Large-scale wind farm simulations"
                }
            },
            "cost_profile": {
                "academic_research": "$400-1200/month",
                "wind_developer": "$1200-4000/month",
                "utility_planning": "$2000-6000/month",
                "research_institution": "$4000-12000/month"
            }
        }

    def _get_energy_storage_config(self) -> Dict[str, Any]:
        """Energy storage systems research and battery modeling"""
        return {
            "name": "Energy Storage Systems Research Platform",
            "description": "Battery modeling, energy storage optimization, and grid-scale storage analysis",
            "spack_packages": [
                # Battery and energy storage modeling
                "pybamm@23.9 %gcc@11.4.0 +python +casadi",     # Battery modeling
                "cantera@3.0.0 %gcc@11.4.0 +python +sundials", # Chemical kinetics
                "comsol@6.1 %gcc@11.4.0 +battery +multiphysics", # Multiphysics simulation
                "ampere@2.1.0 %gcc@11.4.0 +python",            # Battery analytics

                # Electrochemical modeling
                "sundials@6.6.2 %gcc@11.4.0 +mpi +openmp +petsc", # ODE/DAE solver
                "casadi@3.6.3 %gcc@11.4.0 +python +ipopt",     # Optimal control
                "ipopt@3.14.12 %gcc@11.4.0 +mumps +metis",     # Nonlinear optimization

                # Grid storage analysis
                "homer@1.8.1 %gcc@11.4.0 +python",             # Hybrid system optimization
                "sam@2023.12.17 %gcc@11.4.0 +storage",         # Storage system modeling
                "gridlab-d@5.1.0 %gcc@11.4.0 +mysql",         # Grid simulation

                # Python energy storage
                "python@3.11.5 %gcc@11.4.0",
                "py-pybamm@23.9 %gcc@11.4.0",
                "py-battery-analyzer@1.2.0 %gcc@11.4.0",
                "py-energy-storage@0.8.0 %gcc@11.4.0",
                "py-electrochemistry@2.1.0 %gcc@11.4.0",

                # Materials modeling
                "py-pymatgen@2023.8.10 %gcc@11.4.0",           # Materials analysis
                "py-atomate@1.0.3 %gcc@11.4.0",                # Materials workflows
                "py-fireworks@2.0.3 %gcc@11.4.0",              # Workflow management

                # Optimization and machine learning
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-cvxpy@1.3.2 %gcc@11.4.0",                  # Convex optimization

                # Time series and control
                "py-control@0.9.4 %gcc@11.4.0",                # Control systems
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-prophet@1.1.4 %gcc@11.4.0",

                # Data analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",

                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +python",
                "timescaledb@2.11.2 %gcc@11.4.0",
                "redis@7.2.0 %gcc@11.4.0"                       # Real-time data
            ],
            "storage_technologies": {
                "lithium_ion_batteries": {
                    "applications": ["Grid storage", "EVs", "Residential"],
                    "modeling_tools": ["PyBaMM", "COMSOL", "Newman models"],
                    "research_areas": ["Degradation", "Thermal management", "Safety"]
                },
                "flow_batteries": {
                    "applications": ["Grid-scale storage", "Long duration"],
                    "modeling_tools": ["CFD", "Electrochemical models"],
                    "research_areas": ["Electrolyte chemistry", "Stack design", "Economics"]
                },
                "compressed_air": {
                    "applications": ["Utility-scale storage", "Long duration"],
                    "modeling_tools": ["Thermodynamics", "System modeling"],
                    "research_areas": ["Efficiency", "Heat management", "Underground storage"]
                },
                "pumped_hydro": {
                    "applications": ["Grid balancing", "Peak shaving"],
                    "modeling_tools": ["Hydraulic modeling", "Optimization"],
                    "research_areas": ["Site assessment", "Environmental impact", "Economics"]
                }
            },
            "research_capabilities": [
                "Battery cell and pack modeling",
                "Electrochemical impedance analysis",
                "Battery management system optimization",
                "Grid-scale storage integration",
                "Energy storage economics and planning",
                "Battery degradation and lifetime analysis",
                "Thermal management optimization",
                "Safety and failure analysis",
                "Hybrid energy system design",
                "Storage control and dispatch optimization"
            ],
            "aws_instance_recommendations": {
                "battery_modeling": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage": "1TB gp3 SSD",
                    "cost_per_hour": 0.504,
                    "use_case": "Individual battery modeling and analysis"
                },
                "multiphysics_simulation": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage": "2TB gp3 SSD",
                    "cost_per_hour": 1.632,
                    "use_case": "Complex multiphysics battery simulations"
                },
                "grid_storage_optimization": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage": "2TB gp3 SSD + TimescaleDB",
                    "cost_per_hour": 1.008,
                    "use_case": "Grid-scale storage analysis and optimization"
                }
            },
            "cost_profile": {
                "academic_research": "$400-1200/month",
                "battery_developer": "$1200-3000/month",
                "utility_storage": "$2000-6000/month",
                "national_lab": "$4000-12000/month"
            }
        }

    def _get_grid_integration_config(self) -> Dict[str, Any]:
        """Grid integration and power systems analysis"""
        return {
            "name": "Grid Integration & Power Systems Analysis Platform",
            "description": "Power grid modeling, renewable integration, and smart grid analytics",
            "spack_packages": [
                # Power systems simulation
                "gridlab-d@5.1.0 %gcc@11.4.0 +mysql +fncs",    # Grid simulation
                "matpower@7.1 %gcc@11.4.0 +octave",            # Power flow analysis
                "pandapower@2.13.1 %gcc@11.4.0 +python",       # Power system analysis
                "pypower@5.1.16 %gcc@11.4.0 +python",          # Power flow

                # Grid optimization
                "powermodels@0.21.1 %gcc@11.4.0 +julia",       # Power optimization
                "juniper@0.9.1 %gcc@11.4.0 +julia",            # Nonlinear optimization
                "ipopt@3.14.12 %gcc@11.4.0 +mumps +hsl",       # Interior point

                # Renewable integration
                "plexos@9.0 %gcc@11.4.0 +commercial",          # Energy planning
                "homer@1.8.1 %gcc@11.4.0 +grid +python",       # Hybrid systems
                "sam@2023.12.17 %gcc@11.4.0 +grid +wind +solar", # System modeling

                # Real-time grid analytics
                "opendsme@1.0.0 %gcc@11.4.0 +python",          # Distribution management
                "griddyn@1.0.0 %gcc@11.4.0 +sundials +klu",    # Dynamic simulation
                "dynawo@1.6.0 %gcc@11.4.0 +sundials +suitesparse", # Hybrid simulation

                # Python power systems
                "python@3.11.5 %gcc@11.4.0",
                "py-pandapower@2.13.1 %gcc@11.4.0",
                "py-pypower@5.1.16 %gcc@11.4.0",
                "py-powersimdata@0.4.0 %gcc@11.4.0",
                "py-grid2op@1.9.6 %gcc@11.4.0",                # Grid optimization

                # Market analysis
                "py-powersimdata@0.4.0 %gcc@11.4.0",
                "py-revruns@0.6.0 %gcc@11.4.0",                # Renewable energy analysis
                "py-energy-markets@1.2.0 %gcc@11.4.0",

                # Machine learning for grid
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",

                # Time series and forecasting
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-prophet@1.1.4 %gcc@11.4.0",
                "py-arch@6.2.0 %gcc@11.4.0",                   # GARCH models

                # Optimization
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",                   # Linear programming
                "py-pyomo@6.6.2 %gcc@11.4.0",                  # Optimization modeling

                # Data processing
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",

                # Database and streaming
                "postgresql@15.4 %gcc@11.4.0 +python",
                "timescaledb@2.11.2 %gcc@11.4.0",
                "apache-kafka@2.13-3.5.0 %gcc@11.4.0",
                "redis@7.2.0 %gcc@11.4.0"
            ],
            "grid_data_sources": {
                "eia_data": {
                    "description": "US Energy Information Administration data",
                    "bucket": "s3://eia-open-data/",
                    "coverage": "US electricity generation, consumption, prices",
                    "frequency": "Hourly to annual",
                    "api": "EIA API for real-time access"
                },
                "entso_e": {
                    "description": "European Network of Transmission System Operators",
                    "source": "ENTSO-E Transparency Platform",
                    "coverage": "European electricity markets",
                    "data": ["Generation", "Load", "Cross-border flows", "Prices"]
                },
                "pjm_market_data": {
                    "description": "PJM Interconnection market data",
                    "source": "PJM Data Miner",
                    "coverage": "PJM electricity market",
                    "data": ["LMP", "Load", "Generation", "Transmission"]
                },
                "scada_data": {
                    "description": "Supervisory Control and Data Acquisition",
                    "source": "Utility SCADA systems",
                    "coverage": "Real-time grid operations",
                    "frequency": "Seconds to minutes"
                }
            },
            "research_capabilities": [
                "Power flow and stability analysis",
                "Renewable energy integration studies",
                "Grid planning and expansion optimization",
                "Electricity market analysis and modeling",
                "Demand response and load management",
                "Grid reliability and resilience assessment",
                "Transmission and distribution optimization",
                "Smart grid and advanced metering analytics",
                "Energy storage integration and control",
                "Cybersecurity for power systems"
            ],
            "aws_instance_recommendations": {
                "grid_analysis": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage": "2TB gp3 SSD",
                    "cost_per_hour": 0.816,
                    "use_case": "Power systems analysis and optimization"
                },
                "real_time_grid": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage": "2TB gp3 SSD + TimescaleDB + Kafka",
                    "cost_per_hour": 1.008,
                    "use_case": "Real-time grid analytics and control"
                },
                "large_scale_planning": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage": "4TB gp3 SSD",
                    "cost_per_hour": 1.632,
                    "use_case": "Large-scale grid planning and optimization"
                }
            },
            "cost_profile": {
                "academic_research": "$500-1500/month",
                "utility_planning": "$1500-5000/month",
                "iso_rto": "$3000-10000/month",
                "national_lab": "$5000-15000/month"
            }
        }

    def _get_hydroelectric_config(self) -> Dict[str, Any]:
        """Hydroelectric power systems and water resource modeling"""
        return {
            "name": "Hydroelectric Power & Water Resource Systems",
            "description": "Hydroelectric system modeling, dam optimization, and water-energy nexus analysis",
            "spack_packages": [
                # Hydrological modeling
                "hec-ras@6.3.1 %gcc@11.4.0 +fortran",          # River analysis
                "hec-hms@4.11 %gcc@11.4.0 +java",              # Hydrologic modeling
                "swat@2012.664 %gcc@11.4.0 +fortran",          # Watershed modeling
                "vic@5.1.0 %gcc@11.4.0 +netcdf +mpi",          # Variable Infiltration Capacity

                # Hydropower modeling
                "homer@1.8.1 %gcc@11.4.0 +hydro +python",      # Hydropower optimization
                "retscreen@8.0 %gcc@11.4.0 +hydro",            # Renewable energy analysis
                "hydropower-toolkit@1.0 %gcc@11.4.0 +python",  # Hydropower assessment

                # Computational fluid dynamics
                "openfoam@11 %gcc@11.4.0 +mpi +turbulence",    # CFD for turbines
                "flow-3d@12.0 %gcc@11.4.0 +mpi +commercial",   # Hydraulic modeling
                "fluent@2023r1 %gcc@11.4.0 +commercial",       # CFD simulation

                # Water systems optimization
                "heclib@7.0 %gcc@11.4.0 +fortran",             # HEC data management
                "modsim@8.1 %gcc@11.4.0 +fortran",             # Water allocation
                "riverware@8.2 %gcc@11.4.0 +commercial",       # River operations

                # Python hydro tools
                "python@3.11.5 %gcc@11.4.0",
                "py-hydroeval@0.1.0 %gcc@11.4.0",              # Hydrological evaluation
                "py-hydrostats@0.0.15 %gcc@11.4.0",            # Hydrological statistics
                "py-pyflowline@0.2.0 %gcc@11.4.0",            # River network
                "py-pyet@1.3.1 %gcc@11.4.0",                   # Evapotranspiration

                # Climate and weather data
                "py-xarray@2023.7.0 %gcc@11.4.0",
                "py-netcdf4@1.6.4 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",

                # Optimization and analysis
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pyomo@6.6.2 %gcc@11.4.0",

                # Time series analysis
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                "py-prophet@1.1.4 %gcc@11.4.0",

                # Data processing
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",

                # Parallel computing
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers",
                "petsc@3.19.4 %gcc@11.4.0 +mpi"
            ],
            "hydro_data_sources": {
                "usgs_streamflow": {
                    "description": "USGS National Water Information System",
                    "bucket": "s3://usgs-nwis-data/",
                    "coverage": "US streamflow, water quality, groundwater",
                    "frequency": "Real-time to daily",
                    "parameters": ["Streamflow", "Stage", "Temperature", "Quality"]
                },
                "nldas_forcing": {
                    "description": "North American Land Data Assimilation System",
                    "bucket": "s3://nldas-data/",
                    "resolution": "12.5km spatial, hourly temporal",
                    "coverage": "North America",
                    "parameters": ["Precipitation", "Temperature", "Radiation", "Wind"]
                },
                "global_reservoirs": {
                    "description": "Global Reservoir and Dam Database",
                    "source": "GRanD database",
                    "coverage": "Global dam and reservoir information",
                    "parameters": ["Capacity", "Operation", "Purpose", "Coordinates"]
                },
                "hydrosheds": {
                    "description": "Hydrological data and maps",
                    "source": "HydroSHEDS",
                    "coverage": "Global hydrographic information",
                    "data": ["Watersheds", "River networks", "Flow directions"]
                }
            },
            "research_capabilities": [
                "Hydropower potential assessment",
                "Dam and reservoir optimization",
                "Turbine design and efficiency analysis",
                "Environmental flow assessment",
                "Flood control and dam safety",
                "Pumped storage hydropower analysis",
                "Small hydropower and run-of-river systems",
                "Climate change impact on hydropower",
                "Fish passage and environmental impacts",
                "Water-energy nexus optimization"
            ],
            "cost_profile": {
                "academic_research": "$400-1200/month",
                "hydro_developer": "$1200-3000/month",
                "utility_operations": "$2000-6000/month",
                "water_agency": "$3000-8000/month"
            }
        }

    def _get_geothermal_config(self) -> Dict[str, Any]:
        """Geothermal energy systems and subsurface modeling"""
        return {
            "name": "Geothermal Energy Systems & Subsurface Modeling",
            "description": "Geothermal resource assessment, reservoir modeling, and heat pump systems",
            "spack_packages": [
                # Geothermal reservoir modeling
                "tough3@1.0 %gcc@11.4.0 +mpi +petsc",          # Reservoir simulation
                "geophires@3.0 %gcc@11.4.0 +python",           # Geothermal economics
                "feflow@8.0 %gcc@11.4.0 +commercial +mpi",     # Finite element groundwater
                "modflow@6.4.2 %gcc@11.4.0 +mpi",              # Groundwater modeling

                # Heat transfer and thermodynamics
                "comsol@6.1 %gcc@11.4.0 +multiphysics +heat",  # Multiphysics
                "openfoam@11 %gcc@11.4.0 +mpi +heat",          # CFD with heat transfer
                "fenics@2019.1.0 %gcc@11.4.0 +mpi +petsc",     # Finite element

                # Subsurface characterization
                "gmt@6.4.0 %gcc@11.4.0 +netcdf +gdal",         # Geospatial analysis
                "petrel@2023.1 %gcc@11.4.0 +commercial",       # Reservoir modeling
                "gocad@2023.1 %gcc@11.4.0 +commercial",        # 3D geological modeling

                # Python geothermal tools
                "python@3.11.5 %gcc@11.4.0",
                "py-geophires@3.0 %gcc@11.4.0",
                "py-pygfunction@2.2.2 %gcc@11.4.0",            # Ground heat exchanger
                "py-geothermal@1.0.0 %gcc@11.4.0",
                "py-coolprop@6.4.3 %gcc@11.4.0",               # Thermodynamic properties

                # Subsurface modeling
                "py-gempy@3.0.0 %gcc@11.4.0",                  # 3D geological modeling
                "py-pyvista@0.42.2 %gcc@11.4.0",               # 3D visualization
                "py-subsurface@0.2.0 %gcc@11.4.0",

                # Heat pump modeling
                "py-hpxml@1.7.0 %gcc@11.4.0",                  # Heat pump modeling
                "py-energyplus@23.2.0 %gcc@11.4.0",            # Building energy

                # Data analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",

                # Optimization
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pyomo@6.6.2 %gcc@11.4.0",
                "py-deap@1.4.1 %gcc@11.4.0",                   # Evolutionary algorithms

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "paraview@5.11.2 %gcc@11.4.0 +python",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python +postgis"
            ],
            "geothermal_data_sources": {
                "usgs_geothermal": {
                    "description": "USGS Geothermal Resources Database",
                    "source": "USGS National Geothermal Data System",
                    "coverage": "US geothermal resources",
                    "data": ["Temperature logs", "Heat flow", "Geochemistry"]
                },
                "global_heat_flow": {
                    "description": "Global heat flow database",
                    "source": "International Heat Flow Commission",
                    "coverage": "Global heat flow measurements",
                    "parameters": ["Surface heat flux", "Thermal conductivity"]
                },
                "temperature_logs": {
                    "description": "Borehole temperature measurements",
                    "sources": ["SMU Geothermal Lab", "State geological surveys"],
                    "coverage": "Regional to local",
                    "depth": "Surface to 10+ km depth"
                },
                "geological_surveys": {
                    "description": "State and national geological survey data",
                    "sources": ["USGS", "State geological surveys"],
                    "data": ["Geology", "Geophysics", "Geochemistry", "Wells"]
                }
            },
            "research_capabilities": [
                "Geothermal resource assessment and exploration",
                "Enhanced geothermal systems (EGS) modeling",
                "Ground source heat pump design",
                "Geothermal reservoir characterization",
                "Thermal performance optimization",
                "Economic analysis and project evaluation",
                "Environmental impact assessment",
                "Direct use geothermal applications",
                "Co-production with oil and gas",
                "Deep geothermal system analysis"
            ],
            "cost_profile": {
                "academic_research": "$300-1000/month",
                "geothermal_developer": "$1000-3000/month",
                "consulting_firm": "$2000-5000/month",
                "national_lab": "$3000-8000/month"
            }
        }

    def _get_biomass_config(self) -> Dict[str, Any]:
        """Biomass and bioenergy systems modeling"""
        return {
            "name": "Biomass & Bioenergy Systems Platform",
            "description": "Biomass resource assessment, bioenergy conversion, and sustainability analysis",
            "spack_packages": [
                # Biomass assessment and modeling
                "polysys@4.0 %gcc@11.4.0 +fortran",            # Agricultural sector model
                "gtbiom@3.0 %gcc@11.4.0 +python +gdal",        # Biomass assessment
                "bipower@2.1 %gcc@11.4.0 +python",             # Bioenergy potential

                # Process modeling
                "aspen-plus@14.0 %gcc@11.4.0 +commercial",     # Process simulation
                "chemcad@8.0 %gcc@11.4.0 +commercial",         # Chemical process
                "superpro@12.0 %gcc@11.4.0 +commercial",       # Bioprocess design

                # Life cycle assessment
                "brightway2@2.4.4 %gcc@11.4.0 +python",        # LCA framework
                "simapro@9.5 %gcc@11.4.0 +commercial",         # LCA software
                "openlca@1.11.0 %gcc@11.4.0 +java",            # Open source LCA

                # Biomass supply chain
                "opensees@3.6.0 %gcc@11.4.0 +mpi",             # Structural analysis
                "anylogic@8.8 %gcc@11.4.0 +java +commercial",  # System modeling

                # Python biomass tools
                "python@3.11.5 %gcc@11.4.0",
                "py-biomass-assessment@1.2.0 %gcc@11.4.0",
                "py-bioenergy@0.8.0 %gcc@11.4.0",
                "py-biorefinery@2.1.0 %gcc@11.4.0",
                "py-thermosteam@0.35.2 %gcc@11.4.0",           # Chemical process

                # Agricultural modeling
                "py-crop-modeling@1.1.0 %gcc@11.4.0",
                "py-pycrop2ml@1.0.0 %gcc@11.4.0",              # Crop modeling
                "py-agpy@0.3.2 %gcc@11.4.0",                   # Agricultural analysis

                # Remote sensing for biomass
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-satpy@0.43.0 %gcc@11.4.0",
                "py-xarray@2023.7.0 %gcc@11.4.0",

                # Economic and optimization
                "py-pyomo@6.6.2 %gcc@11.4.0",
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",

                # Data analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python +postgis"
            ],
            "biomass_data_sources": {
                "nrel_atlas": {
                    "description": "NREL Biomass Resource Atlas",
                    "bucket": "s3://nrel-pds-biomass/",
                    "coverage": "US biomass resource potential",
                    "types": ["Agricultural residues", "Forest residues", "Energy crops"]
                },
                "cropland_data": {
                    "description": "USDA Cropland Data Layer",
                    "bucket": "s3://usda-nass-aws/",
                    "resolution": "30m spatial, annual temporal",
                    "coverage": "Continental US crop types"
                },
                "forest_inventory": {
                    "description": "USDA Forest Service Forest Inventory",
                    "source": "Forest Inventory and Analysis (FIA)",
                    "coverage": "US forest biomass and volume",
                    "frequency": "Annual plot measurements"
                },
                "global_biomass": {
                    "description": "Global biomass and carbon maps",
                    "sources": ["ESA Biomass", "NASA GEDI", "MODIS"],
                    "coverage": "Global forest and agricultural biomass"
                }
            },
            "research_capabilities": [
                "Biomass resource assessment and mapping",
                "Bioenergy conversion pathway analysis",
                "Supply chain optimization and logistics",
                "Life cycle assessment and sustainability",
                "Techno-economic analysis of biorefineries",
                "Agricultural residue availability",
                "Forest biomass and residue quantification",
                "Energy crop production modeling",
                "Biopower plant optimization",
                "Policy and market analysis"
            ],
            "cost_profile": {
                "academic_research": "$300-1000/month",
                "bioenergy_developer": "$1000-3000/month",
                "consulting_firm": "$2000-5000/month",
                "government_agency": "$3000-8000/month"
            }
        }

    def _get_energy_systems_config(self) -> Dict[str, Any]:
        """Integrated energy systems modeling and analysis"""
        return {
            "name": "Integrated Energy Systems Modeling Platform",
            "description": "Comprehensive energy systems analysis, planning, and optimization",
            "spack_packages": [
                # Energy systems modeling
                "temoa@3.0 %gcc@11.4.0 +python +pyomo",        # Energy optimization
                "times@4.7.0 %gcc@11.4.0 +gams",               # Energy systems
                "message@3.6.0 %gcc@11.4.0 +python +gams",     # Energy planning
                "osemosys@2022.1 %gcc@11.4.0 +python +glpk",   # Energy optimization

                # Economic and policy analysis
                "gtap@10.2 %gcc@11.4.0 +fortran +commercial",  # Global trade
                "gcam@6.0 %gcc@11.4.0 +cpp +xml",              # Global change
                "image@3.2 %gcc@11.4.0 +fortran",              # Integrated assessment

                # Power systems integration
                "plexos@9.0 %gcc@11.4.0 +commercial",
                "gridlab-d@5.1.0 %gcc@11.4.0 +mysql",
                "powerworld@23 %gcc@11.4.0 +commercial",

                # Python energy systems
                "python@3.11.5 %gcc@11.4.0",
                "py-pyomo@6.6.2 %gcc@11.4.0",
                "py-pypsa@0.24.0 %gcc@11.4.0",                 # Power system analysis
                "py-calliope@0.6.10 %gcc@11.4.0",              # Energy systems
                "py-urbs@0.7.5 %gcc@11.4.0",                   # Urban energy systems

                # Economic modeling
                "py-computable@1.0.0 %gcc@11.4.0",             # General equilibrium
                "py-energy-economics@2.1.0 %gcc@11.4.0",
                "py-investment-analysis@1.3.0 %gcc@11.4.0",

                # Climate and policy
                "py-climate-economics@1.8.0 %gcc@11.4.0",
                "py-carbon-pricing@0.9.0 %gcc@11.4.0",
                "py-policy-analysis@1.2.0 %gcc@11.4.0",

                # Optimization
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "gurobi@10.0.3 %gcc@11.4.0 +python +commercial",

                # Machine learning
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",

                # Data analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",

                # Database
                "postgresql@15.4 %gcc@11.4.0 +python"
            ],
            "research_capabilities": [
                "Energy system pathway analysis",
                "Technology assessment and selection",
                "Energy-economy-environment modeling",
                "Carbon pricing and policy analysis",
                "Energy security and resilience assessment",
                "Renewable energy transition planning",
                "Sector coupling and integration",
                "Investment and financial analysis",
                "Climate change mitigation strategies",
                "Energy justice and equity analysis"
            ],
            "cost_profile": {
                "academic_research": "$600-1800/month",
                "policy_analysis": "$1800-5000/month",
                "consulting_firm": "$3000-8000/month",
                "government_agency": "$5000-15000/month"
            }
        }

    def _get_smart_grid_config(self) -> Dict[str, Any]:
        """Smart grid analytics and IoT energy systems"""
        return {
            "name": "Smart Grid Analytics & IoT Energy Systems",
            "description": "Smart grid data analytics, IoT integration, and advanced metering infrastructure",
            "spack_packages": [
                # Smart grid platforms
                "gridspice@1.0 %gcc@11.4.0 +python +mongodb",  # Smart grid simulation
                "opendss@9.6.1.3 %gcc@11.4.0 +python",         # Distribution system
                "derms@2.1.0 %gcc@11.4.0 +python +redis",      # DER management

                # IoT and data streaming
                "apache-kafka@2.13-3.5.0 %gcc@11.4.0",
                "apache-spark@3.4.1 %gcc@11.4.0 +hadoop +scala",
                "apache-storm@2.5.0 %gcc@11.4.0 +java",
                "mqtt@5.0 %gcc@11.4.0 +python",                # IoT messaging

                # Time series databases
                "influxdb@2.7.1 %gcc@11.4.0",
                "timescaledb@2.11.2 %gcc@11.4.0",
                "opentsdb@2.4.1 %gcc@11.4.0 +java",

                # Machine learning and analytics
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",

                # Python smart grid tools
                "python@3.11.5 %gcc@11.4.0",
                "py-pandapower@2.13.1 %gcc@11.4.0",
                "py-grid2op@1.9.6 %gcc@11.4.0",
                "py-smartgrid@1.2.0 %gcc@11.4.0",
                "py-ami-analytics@0.8.0 %gcc@11.4.0",          # Advanced metering

                # Demand response and optimization
                "py-demand-response@1.5.0 %gcc@11.4.0",
                "py-load-forecasting@2.1.0 %gcc@11.4.0",
                "py-energy-disaggregation@1.3.0 %gcc@11.4.0",  # Non-intrusive load monitoring

                # Cybersecurity for smart grid
                "py-grid-security@1.0.0 %gcc@11.4.0",
                "py-scada-security@0.7.0 %gcc@11.4.0",

                # Data processing
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-dask@2023.8.0 %gcc@11.4.0",                # Parallel computing

                # Real-time analytics
                "py-streamz@0.6.4 %gcc@11.4.0",                # Real-time processing
                "py-kafka-python@2.0.2 %gcc@11.4.0",
                "py-paho-mqtt@1.6.1 %gcc@11.4.0",              # MQTT client

                # Visualization and dashboards
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                "grafana@10.1.1 %gcc@11.4.0",

                # Database systems
                "postgresql@15.4 %gcc@11.4.0 +python",
                "mongodb@7.0.0 %gcc@11.4.0",
                "redis@7.2.0 %gcc@11.4.0"
            ],
            "smart_grid_applications": [
                "Advanced metering infrastructure (AMI) analytics",
                "Demand response optimization",
                "Distributed energy resource management",
                "Grid edge analytics and control",
                "Non-intrusive load monitoring",
                "Predictive maintenance for grid assets",
                "Real-time grid state estimation",
                "Cybersecurity monitoring and threat detection",
                "Energy efficiency program analytics",
                "Electric vehicle integration and management"
            ],
            "aws_integration": [
                "IoT Core for device management",
                "Kinesis for real-time data streaming",
                "Lambda for serverless processing",
                "Timestream for time-series data",
                "SageMaker for ML model deployment",
                "QuickSight for business intelligence"
            ],
            "cost_profile": {
                "utility_pilot": "$1000-3000/month",
                "distribution_utility": "$3000-8000/month",
                "smart_city": "$5000-15000/month",
                "grid_operator": "$10000-30000/month"
            }
        }

    def generate_renewable_energy_recommendation(self, workload: RenewableEnergyWorkload) -> Dict[str, Any]:
        """Generate deployment recommendation based on renewable energy workload"""

        # Select configuration based on domain
        domain_mapping = {
            RenewableEnergyDomain.SOLAR_ENERGY: "solar_energy_research",
            RenewableEnergyDomain.WIND_ENERGY: "wind_energy_research",
            RenewableEnergyDomain.ENERGY_STORAGE: "energy_storage_systems",
            RenewableEnergyDomain.GRID_INTEGRATION: "grid_integration_platform",
            RenewableEnergyDomain.HYDROELECTRIC: "hydroelectric_systems",
            RenewableEnergyDomain.GEOTHERMAL: "geothermal_research",
            RenewableEnergyDomain.BIOMASS_BIOENERGY: "biomass_bioenergy",
            RenewableEnergyDomain.ENERGY_SYSTEMS: "energy_systems_modeling",
            RenewableEnergyDomain.SMART_GRID: "smart_grid_analytics"
        }

        config_name = domain_mapping.get(workload.domain, "energy_systems_modeling")
        config = self.renewable_energy_configurations[config_name]

        return {
            "configuration": config,
            "workload": workload,
            "estimated_cost": self._estimate_renewable_energy_cost(workload),
            "data_requirements": self._get_data_requirements(workload),
            "deployment_timeline": "1-4 hours automated setup",
            "optimization_recommendations": self._get_renewable_energy_optimization(workload)
        }

    def _estimate_renewable_energy_cost(self, workload: RenewableEnergyWorkload) -> Dict[str, float]:
        """Estimate costs for renewable energy workloads"""

        # Base cost factors
        if workload.modeling_complexity == "Basic":
            base_cost = 300
        elif workload.modeling_complexity == "Intermediate":
            base_cost = 800
        elif workload.modeling_complexity == "Advanced":
            base_cost = 2000
        else:  # Comprehensive
            base_cost = 5000

        # Scale by spatial scope
        if workload.spatial_scale == "Site":
            spatial_multiplier = 1.0
        elif workload.spatial_scale == "Regional":
            spatial_multiplier = 2.0
        elif workload.spatial_scale == "National":
            spatial_multiplier = 4.0
        else:  # Global
            spatial_multiplier = 8.0

        # Temporal requirements
        if workload.temporal_scale == "Real-time":
            temporal_multiplier = 2.0
        elif workload.temporal_scale in ["Hourly", "Daily"]:
            temporal_multiplier = 1.5
        else:
            temporal_multiplier = 1.0

        # Storage costs
        storage_cost = workload.data_volume_tb * 80  # $80/TB/month

        total_cost = base_cost * spatial_multiplier * temporal_multiplier + storage_cost

        return {
            "compute": base_cost * spatial_multiplier * temporal_multiplier,
            "storage": storage_cost,
            "total": total_cost
        }

    def _get_data_requirements(self, workload: RenewableEnergyWorkload) -> List[str]:
        """Get data requirements based on workload"""
        requirements = []

        if "Weather" in workload.data_sources:
            requirements.append("High-resolution meteorological data access")

        if "Satellite" in workload.data_sources:
            requirements.append("Satellite imagery processing capabilities")

        if "Market" in workload.data_sources:
            requirements.append("Real-time energy market data feeds")

        if workload.temporal_scale == "Real-time":
            requirements.append("Low-latency data streaming infrastructure")

        return requirements

    def _get_renewable_energy_optimization(self, workload: RenewableEnergyWorkload) -> List[str]:
        """Get optimization recommendations for renewable energy workloads"""
        recommendations = []

        if workload.modeling_complexity == "Comprehensive":
            recommendations.append("Consider HPC instances for complex simulations")
            recommendations.append("Use spot instances for batch processing (60-80% savings)")

        if workload.data_volume_tb > 5.0:
            recommendations.append("Implement intelligent data tiering strategies")
            recommendations.append("Use S3 Intelligent Tiering for cost optimization")

        if workload.temporal_scale == "Real-time":
            recommendations.append("Deploy across multiple AZs for high availability")
            recommendations.append("Use CloudFront for global data distribution")

        if workload.spatial_scale == "Global":
            recommendations.append("Consider multi-region deployment")
            recommendations.append("Use AWS Global Accelerator for performance")

        return recommendations

    def list_configurations(self) -> List[str]:
        """List all available renewable energy configurations"""
        return list(self.renewable_energy_configurations.keys())

    def get_configuration_details(self, config_name: str) -> Dict[str, Any]:
        """Get detailed configuration information"""
        if config_name not in self.renewable_energy_configurations:
            raise ValueError(f"Configuration '{config_name}' not found")
        return self.renewable_energy_configurations[config_name]

def main():
    """CLI interface for renewable energy research pack"""
    import argparse

    parser = argparse.ArgumentParser(description="AWS Research Wizard - Renewable Energy Systems Pack")
    parser.add_argument("--list", action="store_true", help="List available configurations")
    parser.add_argument("--config", type=str, help="Show configuration details")
    parser.add_argument("--domain", type=str, choices=[d.value for d in RenewableEnergyDomain],
                       help="Renewable energy domain")
    parser.add_argument("--analysis-type", type=str, choices=["Resource Assessment", "System Design", "Economic Analysis", "Grid Impact"],
                       default="System Design", help="Analysis type")
    parser.add_argument("--spatial-scale", type=str, choices=["Site", "Regional", "National", "Global"],
                       default="Regional", help="Spatial scale")
    parser.add_argument("--output", type=str, help="Output file for recommendation")

    args = parser.parse_args()

    energy_pack = RenewableEnergyResearchPack()

    if args.list:
        print("Available Renewable Energy Research Configurations:")
        for config_name in energy_pack.list_configurations():
            config = energy_pack.get_configuration_details(config_name)
            print(f"  {config_name}: {config['description']}")

    elif args.config:
        try:
            config = energy_pack.get_configuration_details(args.config)
            print(json.dumps(config, indent=2))
        except ValueError as e:
            print(f"Error: {e}")

    elif args.domain:
        workload = RenewableEnergyWorkload(
            domain=RenewableEnergyDomain(args.domain),
            analysis_type=args.analysis_type,
            temporal_scale="Daily",
            spatial_scale=args.spatial_scale,
            modeling_complexity="Intermediate",
            data_sources=["Weather", "Market"],
            optimization_focus="Cost",
            data_volume_tb=1.0
        )

        recommendation = energy_pack.generate_renewable_energy_recommendation(workload)

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
