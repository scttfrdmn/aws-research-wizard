#!/usr/bin/env python3
"""
Food Science & Nutrition Research Pack
Comprehensive food science, nutrition analysis, and food systems modeling for AWS Research Wizard
"""

import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from enum import Enum

class FoodScienceDomain(Enum):
    FOOD_CHEMISTRY = "food_chemistry"
    NUTRITION_ANALYSIS = "nutrition_analysis"
    FOOD_MICROBIOLOGY = "food_microbiology"
    FOOD_PROCESSING = "food_processing"
    FOOD_SAFETY = "food_safety"
    SENSORY_ANALYSIS = "sensory_analysis"
    FOOD_PACKAGING = "food_packaging"
    FOOD_SYSTEMS = "food_systems"
    NUTRITIONAL_EPIDEMIOLOGY = "nutritional_epidemiology"

@dataclass
class FoodScienceWorkload:
    """Food science and nutrition workload characteristics"""
    domain: FoodScienceDomain
    research_scale: str      # Laboratory, Pilot, Industrial, Population, Global
    analysis_type: str       # Compositional, Nutritional, Safety, Quality, Sensory
    study_design: str        # Experimental, Observational, Clinical Trial, Meta-analysis
    data_sources: List[str]  # Laboratory, Survey, Clinical, Market, Supply Chain
    sample_size: int         # Number of samples/participants
    temporal_scale: str      # Real-time, Daily, Weekly, Longitudinal, Cross-sectional
    data_volume_tb: float    # Expected data volume
    computational_intensity: str  # Light, Moderate, Intensive, Extreme

class FoodScienceNutritionPack:
    """
    Comprehensive food science and nutrition research environments optimized for AWS
    Supports food analysis, nutrition research, and food systems modeling
    """

    def __init__(self):
        self.food_science_configurations = {
            "food_chemistry_lab": self._get_food_chemistry_config(),
            "nutrition_analysis": self._get_nutrition_analysis_config(),
            "food_microbiology": self._get_food_microbiology_config(),
            "food_processing": self._get_food_processing_config(),
            "food_safety_platform": self._get_food_safety_config(),
            "sensory_analysis": self._get_sensory_analysis_config(),
            "food_packaging": self._get_food_packaging_config(),
            "food_systems_modeling": self._get_food_systems_config(),
            "nutritional_epidemiology": self._get_nutritional_epidemiology_config()
        }

    def _get_food_chemistry_config(self) -> Dict[str, Any]:
        """Food chemistry analysis and compositional analysis platform"""
        return {
            "name": "Food Chemistry & Compositional Analysis Platform",
            "description": "Comprehensive food composition analysis, chemical characterization, and quality assessment",
            "spack_packages": [
                # Analytical chemistry software
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-chemometrics@1.4.2 %gcc@11.4.0",        # Chemometric analysis
                "r-foodanalysis@1.2.0 %gcc@11.4.0",        # Food analysis tools
                "r-mixomics@6.20.0 %gcc@11.4.0",           # Multivariate analysis

                # Python analytical tools
                "python@3.11.5 %gcc@11.4.0",
                "py-foodpy@2.1.0 %gcc@11.4.0",             # Food analysis toolkit
                "py-chempy@0.8.3 %gcc@11.4.0",             # Chemical analysis
                "py-pychemistry@1.5.0 %gcc@11.4.0",        # Chemistry calculations
                "py-food-composition@2.3.0 %gcc@11.4.0",   # Food composition analysis

                # Spectroscopy and analytical data
                "py-spectral@0.23.1 %gcc@11.4.0",          # Spectral analysis
                "py-nmr-analysis@1.4.0 %gcc@11.4.0",       # NMR data processing
                "py-ms-analysis@2.1.0 %gcc@11.4.0",        # Mass spectrometry
                "py-ftir@1.2.0 %gcc@11.4.0",               # FTIR analysis

                # Chromatography data processing
                "py-chromatography@1.8.0 %gcc@11.4.0",     # Chromatographic analysis
                "py-hplc@2.0.0 %gcc@11.4.0",               # HPLC data processing
                "py-gc-ms@1.7.0 %gcc@11.4.0",              # GC-MS analysis
                "py-lc-ms@2.2.0 %gcc@11.4.0",              # LC-MS analysis

                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Machine learning for food analysis
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",
                "py-xgboost@1.7.6 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",

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
                "development": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 50,
                    "cost_per_hour": 0.085,
                    "use_case": "Development and small-scale food analysis"
                },
                "analytical_lab": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 300,
                    "cost_per_hour": 0.51,
                    "use_case": "Laboratory-scale food chemistry analysis"
                },
                "high_throughput": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "High-throughput food composition analysis"
                },
                "industrial_scale": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Industrial-scale food quality assessment"
                }
            },
            "estimated_cost": {
                "compute": 400,
                "storage": 100,
                "data_transfer": 50,
                "total": 550
            },
            "research_capabilities": [
                "Food composition and nutritional analysis",
                "Chemical characterization and identification",
                "Spectroscopic data analysis (NMR, FTIR, MS)",
                "Chromatographic data processing (HPLC, GC-MS)",
                "Food authenticity and adulteration detection",
                "Quality control and assurance protocols",
                "Chemometric analysis and pattern recognition",
                "Regulatory compliance reporting"
            ],
            "aws_data_sources": [
                "USDA Food Data Central",
                "Food composition databases",
                "Analytical instrument data streams",
                "Quality control laboratory systems"
            ]
        }

    def _get_nutrition_analysis_config(self) -> Dict[str, Any]:
        """Nutrition analysis and dietary assessment platform"""
        return {
            "name": "Nutrition Analysis & Dietary Assessment Platform",
            "description": "Comprehensive nutrition analysis, dietary assessment, and nutritional modeling",
            "spack_packages": [
                # Nutrition analysis software
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-nutrition@1.5.0 %gcc@11.4.0",           # Nutrition analysis
                "r-dietr@0.9.2 %gcc@11.4.0",               # Dietary analysis
                "r-foodcomposition@1.1.0 %gcc@11.4.0",     # Food composition
                "r-nhanesdata@1.0 %gcc@11.4.0",            # NHANES data analysis

                # Python nutrition tools
                "python@3.11.5 %gcc@11.4.0",
                "py-nutrition@2.4.0 %gcc@11.4.0",          # Nutrition analysis toolkit
                "py-dietary-analysis@1.8.0 %gcc@11.4.0",   # Dietary data analysis
                "py-food-diary@1.3.0 %gcc@11.4.0",         # Food diary processing
                "py-nutrient-analysis@2.1.0 %gcc@11.4.0",  # Nutrient calculation

                # Dietary assessment tools
                "py-diet-quality@1.6.0 %gcc@11.4.0",       # Diet quality indices
                "py-meal-planning@2.0.0 %gcc@11.4.0",      # Meal planning algorithms
                "py-portion-size@1.4.0 %gcc@11.4.0",       # Portion size estimation
                "py-nutrition-facts@1.2.0 %gcc@11.4.0",    # Nutrition facts analysis

                # Food database integration
                "py-usda-food-data@1.5.0 %gcc@11.4.0",     # USDA Food Data Central
                "py-fooddata@2.1.0 %gcc@11.4.0",           # Food database access
                "py-recipe-analysis@1.7.0 %gcc@11.4.0",    # Recipe nutritional analysis

                # Biomarker analysis
                "py-biomarkers@1.9.0 %gcc@11.4.0",         # Nutritional biomarkers
                "py-metabolomics@2.2.0 %gcc@11.4.0",       # Metabolomic analysis
                "py-nutrigenomics@1.1.0 %gcc@11.4.0",      # Nutrigenomics analysis

                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Machine learning for nutrition
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",

                # Survey data analysis
                "py-survey@3.1.0 %gcc@11.4.0",             # Survey data processing
                "py-sampling@2.0.0 %gcc@11.4.0",           # Survey sampling methods

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",

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
                "dietary_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Individual dietary analysis and assessment"
                },
                "population_studies": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Population-level nutrition studies and surveys"
                },
                "large_cohorts": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Large cohort studies and longitudinal nutrition research"
                }
            },
            "estimated_cost": {
                "compute": 350,
                "storage": 80,
                "data_transfer": 40,
                "total": 470
            },
            "research_capabilities": [
                "Dietary intake assessment and analysis",
                "Nutrient adequacy and deficiency evaluation",
                "Diet quality index calculation",
                "Food frequency questionnaire analysis",
                "24-hour dietary recall processing",
                "Nutritional biomarker analysis",
                "Population nutrition surveillance",
                "Dietary pattern identification"
            ]
        }

    def _get_food_microbiology_config(self) -> Dict[str, Any]:
        """Food microbiology and safety analysis platform"""
        return {
            "name": "Food Microbiology & Safety Analysis Platform",
            "description": "Microbial analysis, pathogen detection, and food safety assessment",
            "spack_packages": [
                # Microbiology analysis
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-microbiology@1.3.0 %gcc@11.4.0",        # Microbial data analysis
                "r-food-microbiology@2.1.0 %gcc@11.4.0",   # Food microbiology
                "r-predictive-micro@1.5.0 %gcc@11.4.0",    # Predictive microbiology

                # Python microbiology tools
                "python@3.11.5 %gcc@11.4.0",
                "py-microbiology@2.7.0 %gcc@11.4.0",       # Microbial analysis
                "py-pathogen-detection@1.9.0 %gcc@11.4.0", # Pathogen identification
                "py-food-safety@2.3.0 %gcc@11.4.0",        # Food safety analysis
                "py-predictive-micro@1.6.0 %gcc@11.4.0",   # Predictive microbiology

                # Genomic analysis for microbes
                "py-biopython@1.81 %gcc@11.4.0",
                "py-pysam@0.21.0 %gcc@11.4.0",
                "py-scikit-bio@0.5.8 %gcc@11.4.0",
                "py-qiime2@2023.7.0 %gcc@11.4.0",          # Microbiome analysis

                # Machine learning for pathogen detection
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",
                "py-lightgbm@4.0.0 %gcc@11.4.0",

                # Time series for microbial growth
                "py-prophet@1.1.4 %gcc@11.4.0",
                "py-statsforecast@1.6.0 %gcc@11.4.0",

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
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "lab_testing": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 100,
                    "cost_per_hour": 0.085,
                    "use_case": "Laboratory microbial testing and analysis"
                },
                "pathogen_detection": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 300,
                    "cost_per_hour": 0.51,
                    "use_case": "Pathogen detection and identification"
                },
                "genomic_analysis": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Microbial genomics and microbiome analysis"
                }
            },
            "estimated_cost": {
                "compute": 300,
                "storage": 80,
                "data_transfer": 30,
                "total": 410
            },
            "research_capabilities": [
                "Microbial identification and characterization",
                "Pathogen detection and quantification",
                "Predictive microbiology modeling",
                "Food spoilage and shelf-life prediction",
                "Microbiome analysis and diversity assessment",
                "Antimicrobial resistance profiling",
                "Food safety risk assessment",
                "HACCP system development and validation"
            ]
        }

    def _get_food_processing_config(self) -> Dict[str, Any]:
        """Food processing and engineering platform"""
        return {
            "name": "Food Processing & Engineering Platform",
            "description": "Food processing optimization, engineering analysis, and process modeling",
            "spack_packages": [
                # Process modeling and simulation
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-food-processing@1.8.0 %gcc@11.4.0",     # Food processing analysis
                "r-process-optimization@2.1.0 %gcc@11.4.0", # Process optimization

                # Python process engineering
                "python@3.11.5 %gcc@11.4.0",
                "py-food-engineering@2.5.0 %gcc@11.4.0",   # Food engineering calculations
                "py-process-modeling@1.9.0 %gcc@11.4.0",   # Process modeling
                "py-heat-transfer@1.6.0 %gcc@11.4.0",      # Heat transfer calculations
                "py-mass-transfer@1.4.0 %gcc@11.4.0",      # Mass transfer modeling

                # CFD and thermal processing
                "openfoam@10.0 %gcc@11.4.0 +paraview +scotch",
                "py-fenics@2023.1.0 %gcc@11.4.0",          # Finite element analysis
                "py-thermal-processing@2.2.0 %gcc@11.4.0", # Thermal processing

                # Optimization frameworks
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                "py-pyomo@6.6.1 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",

                # Machine learning for process optimization
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",

                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "paraview@5.11.0 %gcc@11.4.0 +python +qt",

                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",

                # AWS-optimized parallel computing with EFA support
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers +pmix +pmi +fabrics",
                "libfabric@1.18.1 %gcc@11.4.0 +verbs +mlx +efa",  # EFA support
                "aws-ofi-nccl@1.7.0 %gcc@11.4.0",  # AWS OFI plugin
                "ucx@1.14.1 %gcc@11.4.0 +verbs +mlx +ib_hw_tm",  # Unified Communication X
                "py-mpi4py@3.1.4 %gcc@11.4.0",
                "slurm@23.02.5 %gcc@11.4.0 +pmix +numa",  # Cluster management

                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "process_design": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Food process design and optimization"
                },
                "cfd_simulation": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage_gb": 500,
                    "cost_per_hour": 1.36,
                    "use_case": "CFD simulation and thermal processing modeling"
                },
                "large_scale_modeling": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Large-scale process modeling and optimization"
                }
            },
            "estimated_cost": {
                "compute": 500,
                "storage": 120,
                "data_transfer": 60,
                "total": 680
            },
            "research_capabilities": [
                "Food process modeling and simulation",
                "Thermal processing optimization",
                "Heat and mass transfer analysis",
                "CFD modeling of food processing equipment",
                "Process control and automation design",
                "Energy efficiency optimization",
                "Scale-up from laboratory to industrial processes",
                "Quality attribute preservation modeling"
            ]
        }

    def _get_food_safety_config(self) -> Dict[str, Any]:
        """Food safety and risk assessment platform"""
        return {
            "name": "Food Safety & Risk Assessment Platform",
            "description": "Food safety analysis, risk assessment, and regulatory compliance",
            "spack_packages": [
                # Risk assessment tools
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-food-safety@2.0.0 %gcc@11.4.0",         # Food safety analysis
                "r-risk-assessment@1.7.0 %gcc@11.4.0",     # Risk assessment
                "r-monte-carlo@2.3.0 %gcc@11.4.0",         # Monte Carlo simulation

                # Python food safety tools
                "python@3.11.5 %gcc@11.4.0",
                "py-food-safety@2.6.0 %gcc@11.4.0",        # Food safety toolkit
                "py-risk-analysis@1.8.0 %gcc@11.4.0",      # Risk analysis
                "py-haccp@1.5.0 %gcc@11.4.0",              # HACCP system analysis
                "py-contamination@2.1.0 %gcc@11.4.0",      # Contamination modeling

                # Regulatory compliance
                "py-fda-compliance@1.3.0 %gcc@11.4.0",     # FDA compliance tools
                "py-codex@1.1.0 %gcc@11.4.0",              # Codex Alimentarius
                "py-regulatory@2.0.0 %gcc@11.4.0",         # Regulatory analysis

                # Traceability and supply chain
                "py-traceability@2.4.0 %gcc@11.4.0",       # Food traceability
                "py-supply-chain@1.7.0 %gcc@11.4.0",       # Supply chain analysis
                "py-blockchain-food@1.2.0 %gcc@11.4.0",    # Blockchain for food safety

                # Machine learning for safety prediction
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
                "py-dash@2.13.0 %gcc@11.4.0",

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
                "safety_analysis": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 100,
                    "cost_per_hour": 0.085,
                    "use_case": "Food safety data analysis and reporting"
                },
                "risk_assessment": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 300,
                    "cost_per_hour": 0.51,
                    "use_case": "Quantitative microbial risk assessment"
                },
                "supply_chain": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Supply chain risk analysis and traceability"
                }
            ],
            "estimated_cost": {
                "compute": 250,
                "storage": 60,
                "data_transfer": 30,
                "total": 340
            },
            "research_capabilities": [
                "Quantitative microbial risk assessment (QMRA)",
                "HACCP system development and validation",
                "Food safety management system design",
                "Contamination source tracking",
                "Supply chain risk analysis",
                "Regulatory compliance assessment",
                "Food safety training and education platforms",
                "Crisis management and recall procedures"
            ]
        }

    def _get_sensory_analysis_config(self) -> Dict[str, Any]:
        """Sensory analysis and consumer research platform"""
        return {
            "name": "Sensory Analysis & Consumer Research Platform",
            "description": "Sensory evaluation, consumer testing, and product development",
            "spack_packages": [
                # Sensory analysis software
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-sensory@1.9.0 %gcc@11.4.0",             # Sensory data analysis
                "r-consumer@2.1.0 %gcc@11.4.0",            # Consumer research
                "r-hedonic@1.5.0 %gcc@11.4.0",             # Hedonic analysis
                "r-discrimination@1.3.0 %gcc@11.4.0",      # Discrimination tests

                # Python sensory tools
                "python@3.11.5 %gcc@11.4.0",
                "py-sensory@2.4.0 %gcc@11.4.0",            # Sensory analysis toolkit
                "py-consumer-research@1.7.0 %gcc@11.4.0",  # Consumer research tools
                "py-taste-panel@1.5.0 %gcc@11.4.0",        # Taste panel analysis
                "py-descriptive@2.0.0 %gcc@11.4.0",        # Descriptive analysis

                # Survey and experimental design
                "py-survey@3.1.0 %gcc@11.4.0",
                "py-experimental-design@2.2.0 %gcc@11.4.0",
                "py-randomization@1.4.0 %gcc@11.4.0",

                # Machine learning for sensory prediction
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",

                # Natural language processing for consumer feedback
                "py-nltk@3.8.1 %gcc@11.4.0",
                "py-spacy@3.6.1 %gcc@11.4.0",
                "py-transformers@4.33.2 %gcc@11.4.0",
                "py-sentiment-analysis@2.1.0 %gcc@11.4.0",

                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-wordcloud@1.9.2 %gcc@11.4.0",

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
                "panel_analysis": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 50,
                    "cost_per_hour": 0.085,
                    "use_case": "Small sensory panel analysis"
                },
                "consumer_studies": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Consumer research and preference mapping"
                },
                "large_scale_testing": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Large-scale consumer testing and market research"
                }
            },
            "estimated_cost": {
                "compute": 200,
                "storage": 40,
                "data_transfer": 20,
                "total": 260
            },
            "research_capabilities": [
                "Descriptive sensory analysis",
                "Consumer preference testing",
                "Discrimination and threshold testing",
                "Hedonic scaling and preference mapping",
                "Texture profile analysis",
                "Flavor profile development",
                "Product optimization and reformulation",
                "Sensory-instrumental correlations"
            ]
        }

    def _get_food_packaging_config(self) -> Dict[str, Any]:
        """Food packaging research and development platform"""
        return {
            "name": "Food Packaging Research & Development Platform",
            "description": "Packaging material testing, barrier properties, and sustainability analysis",
            "spack_packages": [
                # Materials science tools
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-packaging@1.6.0 %gcc@11.4.0",           # Packaging analysis
                "r-materials@2.1.0 %gcc@11.4.0",           # Materials testing

                # Python packaging tools
                "python@3.11.5 %gcc@11.4.0",
                "py-packaging-science@2.3.0 %gcc@11.4.0",  # Packaging science tools
                "py-barrier-properties@1.8.0 %gcc@11.4.0", # Barrier property analysis
                "py-shelf-life@2.0.0 %gcc@11.4.0",         # Shelf-life modeling
                "py-sustainability@1.5.0 %gcc@11.4.0",     # Sustainability analysis

                # Polymer and materials modeling
                "py-polymer-modeling@1.7.0 %gcc@11.4.0",   # Polymer modeling
                "py-materials-analysis@2.2.0 %gcc@11.4.0", # Materials analysis
                "py-barrier-modeling@1.4.0 %gcc@11.4.0",   # Barrier modeling

                # Life cycle assessment
                "py-lca@2.5.0 %gcc@11.4.0",                # Life cycle assessment
                "py-carbon-footprint@1.6.0 %gcc@11.4.0",   # Carbon footprint analysis
                "py-environmental@2.1.0 %gcc@11.4.0",      # Environmental impact

                # Machine learning for packaging
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",

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
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "materials_testing": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 100,
                    "cost_per_hour": 0.085,
                    "use_case": "Packaging materials testing and analysis"
                },
                "package_design": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Package design and optimization"
                },
                "sustainability": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 300,
                    "cost_per_hour": 0.51,
                    "use_case": "Sustainability and life cycle assessment"
                }
            ],
            "estimated_cost": {
                "compute": 180,
                "storage": 50,
                "data_transfer": 25,
                "total": 255
            },
            "research_capabilities": [
                "Barrier property testing and modeling",
                "Package design and optimization",
                "Shelf-life prediction modeling",
                "Migration testing and analysis",
                "Sustainability and LCA assessment",
                "Smart packaging development",
                "Active and intelligent packaging systems",
                "Recyclability and biodegradability testing"
            ]
        }

    def _get_food_systems_config(self) -> Dict[str, Any]:
        """Food systems modeling and analysis platform"""
        return {
            "name": "Food Systems Modeling & Analysis Platform",
            "description": "Food systems sustainability, supply chain analysis, and policy modeling",
            "spack_packages": [
                # Systems modeling
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-food-systems@2.2.0 %gcc@11.4.0",        # Food systems analysis
                "r-supply-chain@1.8.0 %gcc@11.4.0",        # Supply chain modeling
                "r-sustainability@2.0.0 %gcc@11.4.0",      # Sustainability metrics

                # Python systems tools
                "python@3.11.5 %gcc@11.4.0",
                "py-food-systems@2.7.0 %gcc@11.4.0",       # Food systems modeling
                "py-supply-chain@2.3.0 %gcc@11.4.0",       # Supply chain analysis
                "py-food-security@1.9.0 %gcc@11.4.0",      # Food security analysis
                "py-policy-modeling@2.1.0 %gcc@11.4.0",    # Policy impact modeling

                # Network analysis for food systems
                "py-networkx@3.1 %gcc@11.4.0",
                "py-igraph@0.10.6 %gcc@11.4.0",
                "py-food-networks@1.5.0 %gcc@11.4.0",

                # Geospatial analysis
                "py-geopandas@0.13.2 %gcc@11.4.0",
                "py-rasterio@1.3.8 %gcc@11.4.0",
                "py-fiona@1.9.4 %gcc@11.4.0",
                "py-shapely@2.0.1 %gcc@11.4.0",

                # Economic modeling
                "py-economic-modeling@1.7.0 %gcc@11.4.0",
                "py-market-analysis@2.0.0 %gcc@11.4.0",
                "py-price-modeling@1.4.0 %gcc@11.4.0",

                # Machine learning for systems
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",

                # Optimization
                "py-cvxpy@1.3.2 %gcc@11.4.0",
                "py-pulp@2.7.0 %gcc@11.4.0",
                "py-pyomo@6.6.1 %gcc@11.4.0",

                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-folium@0.14.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",

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
                "systems_analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Food systems analysis and modeling"
                },
                "supply_chain": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Complex supply chain and network analysis"
                },
                "global_modeling": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Global food systems and policy modeling"
                }
            },
            "estimated_cost": {
                "compute": 450,
                "storage": 100,
                "data_transfer": 50,
                "total": 600
            },
            "research_capabilities": [
                "Food systems sustainability assessment",
                "Supply chain resilience analysis",
                "Food security and access modeling",
                "Policy impact assessment",
                "Market dynamics and price modeling",
                "Climate change impact on food systems",
                "Urban food systems analysis",
                "Food waste and loss quantification"
            ]
        }

    def _get_nutritional_epidemiology_config(self) -> Dict[str, Any]:
        """Nutritional epidemiology and population health platform"""
        return {
            "name": "Nutritional Epidemiology & Population Health Platform",
            "description": "Population nutrition studies, dietary patterns, and health outcomes analysis",
            "spack_packages": [
                # Epidemiology software
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-epidemiology@1.8.0 %gcc@11.4.0",        # Epidemiological analysis
                "r-survival@3.5-7 %gcc@11.4.0",            # Survival analysis
                "r-meta-analysis@2.1.0 %gcc@11.4.0",       # Meta-analysis
                "r-cohort@1.5.0 %gcc@11.4.0",              # Cohort studies

                # Python epidemiology tools
                "python@3.11.5 %gcc@11.4.0",
                "py-epidemiology@2.6.0 %gcc@11.4.0",       # Epidemiological analysis
                "py-nutrition-epi@1.9.0 %gcc@11.4.0",      # Nutritional epidemiology
                "py-dietary-patterns@2.2.0 %gcc@11.4.0",   # Dietary pattern analysis
                "py-health-outcomes@1.7.0 %gcc@11.4.0",    # Health outcomes modeling

                # Causal inference
                "py-causal-inference@1.5.0 %gcc@11.4.0",
                "py-propensity-score@2.0.0 %gcc@11.4.0",
                "py-instrumental-variables@1.3.0 %gcc@11.4.0",

                # Survival analysis
                "py-lifelines@0.27.7 %gcc@11.4.0",
                "py-survival-analysis@2.1.0 %gcc@11.4.0",

                # Machine learning for epidemiology
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-pytorch@2.0.1 %gcc@11.4.0",

                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",

                # Survey data analysis
                "py-survey@3.1.0 %gcc@11.4.0",
                "py-sampling@2.0.0 %gcc@11.4.0",
                "py-weighting@1.4.0 %gcc@11.4.0",

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
                "cohort_studies": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 300,
                    "cost_per_hour": 0.51,
                    "use_case": "Cohort and case-control studies"
                },
                "population_analysis": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Population-level nutrition and health analysis"
                },
                "meta_analysis": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Large-scale meta-analyses and systematic reviews"
                }
            },
            "estimated_cost": {
                "compute": 500,
                "storage": 150,
                "data_transfer": 75,
                "total": 725
            },
            "research_capabilities": [
                "Prospective cohort study analysis",
                "Case-control and cross-sectional studies",
                "Dietary pattern identification and analysis",
                "Disease risk factor assessment",
                "Meta-analysis and systematic reviews",
                "Causal inference in nutrition research",
                "Nutritional biomarker validation",
                "Population health surveillance"
            ]
        }

    def generate_food_science_recommendation(self, workload: FoodScienceWorkload) -> Dict[str, Any]:
        """Generate optimized AWS infrastructure recommendation for food science research"""

        # Select appropriate configuration based on domain
        domain_config_map = {
            FoodScienceDomain.FOOD_CHEMISTRY: "food_chemistry_lab",
            FoodScienceDomain.NUTRITION_ANALYSIS: "nutrition_analysis",
            FoodScienceDomain.FOOD_MICROBIOLOGY: "food_microbiology",
            FoodScienceDomain.FOOD_PROCESSING: "food_processing",
            FoodScienceDomain.FOOD_SAFETY: "food_safety_platform",
            FoodScienceDomain.SENSORY_ANALYSIS: "sensory_analysis",
            FoodScienceDomain.FOOD_PACKAGING: "food_packaging",
            FoodScienceDomain.FOOD_SYSTEMS: "food_systems_modeling",
            FoodScienceDomain.NUTRITIONAL_EPIDEMIOLOGY: "nutritional_epidemiology"
        }

        config_name = domain_config_map.get(workload.domain, "food_chemistry_lab")
        base_config = self.food_science_configurations[config_name].copy()

        # Adjust configuration based on workload characteristics
        self._optimize_for_scale(base_config, workload)
        self._optimize_for_sample_size(base_config, workload)
        self._optimize_for_computational_intensity(base_config, workload)

        # Generate cost estimates
        base_config["estimated_cost"] = self._calculate_food_science_costs(workload, base_config)

        # Add optimization recommendations
        base_config["optimization_recommendations"] = self._generate_optimization_recommendations(workload)

        return {
            "configuration": base_config,
            "workload_analysis": {
                "domain": workload.domain.value,
                "research_scale": workload.research_scale,
                "analysis_type": workload.analysis_type,
                "sample_size": workload.sample_size,
                "computational_requirements": workload.computational_intensity,
                "data_volume": f"{workload.data_volume_tb} TB"
            },
            "deployment_recommendations": self._generate_deployment_recommendations(workload),
            "estimated_cost": base_config["estimated_cost"]
        }

    def _optimize_for_scale(self, config: Dict[str, Any], workload: FoodScienceWorkload):
        """Optimize configuration based on research scale"""
        scale_multipliers = {
            "Laboratory": 1.0,
            "Pilot": 1.5,
            "Industrial": 3.0,
            "Population": 6.0,
            "Global": 12.0
        }

        multiplier = scale_multipliers.get(workload.research_scale, 1.0)

        # Adjust instance recommendations based on scale
        if "aws_instance_recommendations" in config:
            for instance_config in config["aws_instance_recommendations"].values():
                if multiplier > 3.0:
                    # Scale up storage for large-scale studies
                    instance_config["storage_gb"] = int(instance_config["storage_gb"] * multiplier)

    def _optimize_for_sample_size(self, config: Dict[str, Any], workload: FoodScienceWorkload):
        """Optimize configuration based on sample size"""
        if workload.sample_size > 10000:
            # Add big data processing tools for large sample sizes
            if "spack_packages" in config:
                config["spack_packages"].extend([
                    "py-dask@2023.7.1 %gcc@11.4.0",
                    "py-ray@2.6.1 %gcc@11.4.0"
                ])

    def _optimize_for_computational_intensity(self, config: Dict[str, Any], workload: FoodScienceWorkload):
        """Optimize configuration based on computational intensity"""
        if workload.computational_intensity in ["Intensive", "Extreme"]:
            # Upgrade to compute-optimized instances for intensive workloads
            if "aws_instance_recommendations" in config:
                for key, instance_config in config["aws_instance_recommendations"].items():
                    if workload.computational_intensity == "Extreme":
                        if "meta" in key.lower() or "population" in key.lower():
                            instance_config["instance_type"] = "r6i.16xlarge"
                            instance_config["cost_per_hour"] = 4.10

    def _calculate_food_science_costs(self, workload: FoodScienceWorkload, config: Dict[str, Any]) -> Dict[str, float]:
        """Calculate estimated costs for food science research infrastructure"""
        base_compute = 300
        base_storage = 80
        base_data_transfer = 40

        # Scale costs based on research scale
        scale_multipliers = {"Laboratory": 1.0, "Pilot": 1.5, "Industrial": 3.0, "Population": 6.0, "Global": 12.0}
        multiplier = scale_multipliers.get(workload.research_scale, 1.0)

        # Adjust for sample size
        sample_multiplier = 1.0 + (workload.sample_size / 10000.0)

        # Adjust for computational intensity
        intensity_multipliers = {"Light": 0.5, "Moderate": 1.0, "Intensive": 2.0, "Extreme": 4.0}
        intensity_mult = intensity_multipliers.get(workload.computational_intensity, 1.0)

        compute_cost = base_compute * multiplier * sample_multiplier * intensity_mult
        storage_cost = base_storage * (1 + workload.data_volume_tb / 5.0)
        data_transfer_cost = base_data_transfer * multiplier

        return {
            "compute": compute_cost,
            "storage": storage_cost,
            "data_transfer": data_transfer_cost,
            "total": compute_cost + storage_cost + data_transfer_cost
        }

    def _generate_optimization_recommendations(self, workload: FoodScienceWorkload) -> List[str]:
        """Generate optimization recommendations for food science workloads"""
        recommendations = []

        if workload.research_scale in ["Population", "Global"]:
            recommendations.append("Consider using Spot Instances for batch processing to reduce costs by 60-90%")
            recommendations.append("Implement auto-scaling for variable analytical workloads")

        if workload.sample_size > 5000:
            recommendations.append("Use S3 Intelligent Tiering for automatic cost optimization of large datasets")
            recommendations.append("Consider AWS Batch for parallel processing of large sample batches")

        if "Real-time" in workload.temporal_scale:
            recommendations.append("Use AWS IoT Core for real-time food safety monitoring")
            recommendations.append("Consider Amazon Kinesis for streaming analytical data")

        if workload.computational_intensity == "Extreme":
            recommendations.append("Use GPU instances for machine learning in food analysis")
            recommendations.append("Consider AWS ParallelCluster for intensive computational workflows")

        if "Longitudinal" in workload.temporal_scale:
            recommendations.append("Implement automated backup strategies for long-term studies")
            recommendations.append("Use Amazon Glacier for cost-effective long-term data archival")

        return recommendations

    def _generate_deployment_recommendations(self, workload: FoodScienceWorkload) -> Dict[str, Any]:
        """Generate deployment recommendations for food science research"""
        return {
            "deployment_strategy": "multi-tier" if workload.research_scale in ["Population", "Global"] else "single-tier",
            "backup_strategy": "automated_daily_snapshots",
            "monitoring": ["CloudWatch for infrastructure", "Custom dashboards for food science metrics"],
            "security": ["VPC with private subnets", "IAM roles for service access", "Data encryption at rest and in transit", "HIPAA compliance for human studies"],
            "disaster_recovery": "cross-region backup" if workload.research_scale == "Global" else "single-region backup"
        }

if __name__ == "__main__":
    # Example usage
    pack = FoodScienceNutritionPack()

    # Example workload
    workload = FoodScienceWorkload(
        domain=FoodScienceDomain.NUTRITION_ANALYSIS,
        research_scale="Population",
        analysis_type="Nutritional",
        study_design="Observational",
        data_sources=["Survey", "Clinical"],
        sample_size=5000,
        temporal_scale="Longitudinal",
        data_volume_tb=2.0,
        computational_intensity="Moderate"
    )

    recommendation = pack.generate_food_science_recommendation(workload)
    print(json.dumps(recommendation, indent=2))
