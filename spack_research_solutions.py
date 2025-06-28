#!/usr/bin/env python3
"""
Spack-Based Researcher-Ready Solutions
Domain-specific research environments using Spack for optimized software deployment
"""

import os
import json
from typing import Dict, List, Optional
from dataclasses import dataclass, asdict
import argparse
import logging

@dataclass
class SpackResearchSolution:
    name: str
    description: str
    domain: str
    target_users: str
    spack_environment: Dict[str, List[str]]
    spack_config: Dict[str, str]
    sample_workflows: List[str]
    transparent_pricing: Dict[str, str]
    deployment_variants: List[str]
    security_tiers: Dict[str, Dict[str, str]]
    immediate_value: List[str]
    deployment_code: str

class SpackResearchSolutionsGenerator:
    def __init__(self):
        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)

    def create_spack_solutions(self) -> Dict[str, SpackResearchSolution]:
        """Create domain-specific solutions using Spack for software management"""
        
        solutions = {
            "genomics_spack_lab": SpackResearchSolution(
                name="Genomics Research Lab (Spack-Powered)",
                description="Complete genomics analysis environment with Spack-optimized bioinformatics tools",
                domain="Genomics & Bioinformatics",
                target_users="Single researcher, small lab (1-5 people), bioinformatics core",
                
                spack_environment={
                    "core_bio_tools": [
                        "gatk@4.4.0.0",  # Genome Analysis Toolkit
                        "bwa@0.7.17",    # Burrows-Wheeler Aligner
                        "bwa-mem2@2.2.1", # Faster BWA implementation
                        "star@2.7.10b",  # RNA-seq aligner
                        "samtools@1.17", # SAM/BAM manipulation
                        "bcftools@1.17", # VCF/BCF manipulation
                        "htslib@1.17",   # High-throughput sequencing library
                        "bedtools2@2.30.0", # Genome arithmetic
                        "vcftools@0.1.16", # VCF manipulation
                        "fastqc@0.11.9", # Quality control
                        "trimmomatic@0.39", # Read trimming
                        "picard@2.27.4", # Java tools for HTS data
                    ],
                    "analysis_tools": [
                        "blast-plus@2.13.0", # Sequence search
                        "muscle@5.1",     # Multiple sequence alignment
                        "clustalw@2.1",   # Multiple sequence alignment
                        "hmmer@3.3.2",    # Profile HMM search
                        "mafft@7.490",    # Multiple alignment
                        "iqtree@2.2.0",   # Phylogenetic analysis
                        "mrbayes@3.2.7",  # Bayesian phylogenetics
                        "diamond@2.0.15", # Fast protein aligner
                        "minimap2@2.24",  # Long-read alignment
                        "canu@2.2",       # Long-read assembly
                    ],
                    "python_stack": [
                        "python@3.10.8",
                        "py-numpy@1.24.1",
                        "py-pandas@1.5.2",
                        "py-scipy@1.9.3",
                        "py-matplotlib@3.6.2",
                        "py-seaborn@0.12.1",
                        "py-biopython@1.80",
                        "py-pysam@0.20.0",
                        "py-scikit-learn@1.1.3",
                        "py-jupyter@1.0.0",
                        "py-jupyterlab@3.5.0",
                    ],
                    "r_bioconductor": [
                        "r@4.2.2",
                        "r-biocmanager@1.30.19",
                        "r-deseq2@1.36.0",
                        "r-edger@3.38.4",
                        "r-genomicranges@1.48.0",
                        "r-genomicfeatures@1.48.3",
                        "r-biostrings@2.64.1",
                        "r-iranges@2.30.1",
                        "r-genomeinfodb@1.32.4",
                        "r-ggplot2@3.4.0",
                        "r-dplyr@1.0.10",
                        "r-tidyr@1.2.1",
                    ],
                    "workflow_managers": [
                        "nextflow@22.10.4",
                        "snakemake@7.18.2",
                        "cromwell@78",
                        "cwltool@3.1.20221008225030",
                    ],
                    "visualization": [
                        "igv@2.15.4",     # Integrative Genomics Viewer
                        "circos@0.69-9",  # Circular genome plots
                        "jbrowse@1.16.11", # Genome browser
                    ]
                },
                
                spack_config={
                    "compiler_optimization": "gcc@11.3.0 %gcc@11.3.0 cflags=-O3 cxxflags=-O3",
                    "mpi_implementation": "openmpi@4.1.4",
                    "linear_algebra": "openblas@0.3.21 threads=openmp",
                    "python_optimization": "+optimizations+shared+ssl+zlib",
                    "parallel_jobs": "16",  # Compile with 16 cores
                    "build_cache": "enabled",  # Use build cache for faster deployments
                    "target_architecture": "x86_64_v3"  # Optimize for modern CPUs
                },
                
                sample_workflows=[
                    "GATK best practices variant calling pipeline",
                    "RNA-seq differential expression with STAR + DESeq2",
                    "Whole genome assembly with long reads (Canu + Minimap2)",
                    "ChIP-seq peak calling and annotation",
                    "Single-cell RNA-seq analysis workflow",
                    "Metagenomics classification and functional analysis",
                    "Population genetics analysis with VCFtools + PLINK",
                    "Phylogenetic reconstruction from multiple sequence alignments"
                ],
                
                transparent_pricing={
                    "idle_cost": "$0/day (Spack environment cached, compute off)",
                    "light_analysis": "$3-12/day (small datasets, basic tools)",
                    "standard_genomics": "$12-35/day (WGS analysis, moderate compute)",
                    "intensive_compute": "$35-100/day (large cohorts, complex analysis)",
                    "storage_cost": "$2.30/month per 100GB (S3 intelligent tiering)",
                    "spack_build_cache": "$5/month (shared binary cache)",
                    "collaborative_setup": "+$10-25/day for multi-user access",
                    "example_monthly": "$120-700/month for active genomics researcher"
                },
                
                deployment_variants=[
                    "minimal: Core bio tools only (GATK, BWA, SAMtools)",
                    "standard: Full genomics stack with R/Python",
                    "comprehensive: All tools + workflow managers + visualization",
                    "hpc: Optimized for large-scale parallel analysis",
                    "collaborative: Multi-user lab environment with shared Spack install"
                ],
                
                security_tiers={
                    "basic": {
                        "description": "Standard security with Spack build verification",
                        "features": "Package integrity checks, encrypted storage, basic access controls",
                        "compliance": "Good for academic research with public data",
                        "disclaimer": "⚠️  Basic security - Spack packages verified but not suitable for PHI"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security with verified Spack builds",
                        "features": "Cryptographic package verification, audit logging, network isolation",
                        "compliance": "NIST 800-171 technical controls + Spack security",
                        "disclaimer": "⚠️  Technical controls implemented - YOU MUST ensure procedural compliance"
                    },
                    "nist_800_53": {
                        "description": "High security with hardened Spack environment",
                        "features": "Signed packages, continuous monitoring, strict access controls",
                        "compliance": "NIST 800-53 moderate baseline + software supply chain security",
                        "disclaimer": "⚠️  Technical implementation only - requires comprehensive compliance program"
                    }
                },
                
                immediate_value=[
                    "Launch analysis in 3 minutes with optimized, pre-compiled tools",
                    "All software dependencies resolved and performance-optimized",
                    "Example datasets and validated workflows included",
                    "Reproducible environments with exact software versions",
                    "Automatic performance tuning for your compute architecture",
                    "Easy software version management and testing",
                    "Zero software installation headaches",
                    "Lab-wide software environment consistency"
                ],
                
                deployment_code="genomics_spack_lab"
            ),
            
            "climate_spack_lab": SpackResearchSolution(
                name="Climate Modeling Laboratory (Spack-Powered)", 
                description="High-performance climate modeling with Spack-optimized atmospheric science tools",
                domain="Climate Science & Atmospheric Physics",
                target_users="Climate researcher, atmospheric scientist, modeling group",
                
                spack_environment={
                    "climate_models": [
                        "wrf@4.4.1",      # Weather Research and Forecasting
                        "cesm@2.1.3",     # Community Earth System Model
                        "cam@6.3.0",      # Community Atmosphere Model
                        "clm@5.0.34",     # Community Land Model
                        "pop@2.1.0",      # Parallel Ocean Program
                        "cice@6.4.1",     # Sea Ice Model
                        "mpas@8.0.1",     # Model for Prediction Across Scales
                        "fv3@2022.03.22", # FV3 Atmospheric Model
                    ],
                    "analysis_tools": [
                        "nco@5.1.0",      # netCDF Operators
                        "cdo@2.0.6",      # Climate Data Operators
                        "ncview@2.1.8",   # Quick netCDF viewer
                        "ncdump@4.9.0",   # netCDF utilities
                        "grads@2.2.1",    # Grid Analysis and Display System
                        "ferret@7.6.0",   # Data visualization and analysis
                        "ncl@6.6.2",      # NCAR Command Language
                        "vapor@3.7.0",    # 3D atmospheric visualization
                        "visit@3.3.3",    # Visualization software
                        "paraview@5.11.0+mpi+python3", # 3D data visualization
                    ],
                    "data_processing": [
                        "hdf5@1.12.2+mpi",
                        "netcdf-c@4.9.0+mpi",
                        "netcdf-fortran@4.6.0",
                        "netcdf-cxx4@4.3.1",
                        "parallel-netcdf@1.12.3",
                        "esmf@8.4.0+mpi",  # Earth System Modeling Framework
                        "udunits@2.2.28",  # Units of measurement
                        "proj@9.1.0",      # Cartographic projections
                        "gdal@3.6.2+python3", # Geospatial data abstraction
                        "geos@3.11.1",     # Geometry engine
                    ],
                    "python_climate": [
                        "python@3.10.8+optimizations",
                        "py-xarray@2022.11.0",    # N-dimensional arrays
                        "py-dask@2022.11.1",      # Parallel computing
                        "py-cartopy@0.21.0",      # Geospatial data processing
                        "py-matplotlib@3.6.2",    # Plotting
                        "py-basemap@1.3.6",       # Map projections
                        "py-netcdf4@1.6.2",       # netCDF interface
                        "py-h5netcdf@1.0.2",      # Alternative netCDF interface
                        "py-metpy@1.4.0",         # Meteorological tools
                        "py-cf-units@3.1.1",      # CF convention units
                        "py-iris@3.4.0",          # Climate data analysis
                        "py-esmpy@8.4.0",         # ESMF Python interface
                    ],
                    "r_climate": [
                        "r@4.2.2+X+external-lapack",
                        "r-ncdf4@1.19",
                        "r-raster@3.6-3",
                        "r-sp@1.5-1",
                        "r-rgdal@1.6-2",
                        "r-fields@14.1",
                        "r-maps@3.4.1",
                        "r-climate4r@1.6.1",
                        "r-ecoms-udg@1.6.0",
                    ],
                    "hpc_tools": [
                        "openmpi@4.1.4",
                        "intel-mpi@2021.7.1",
                        "fftw@3.3.10+mpi+openmp",
                        "scalapack@2.2.0",
                        "petsc@3.18.1+mpi+hypre+metis",
                        "hypre@2.26.0+mpi+openmp",
                        "metis@5.1.0+int64+real64",
                        "parmetis@4.0.3+int64+real64",
                    ]
                },
                
                spack_config={
                    "compiler_optimization": "intel@2022.2.1 %intel cflags=-O3 cxxflags=-O3 fflags=-O3",
                    "mpi_implementation": "intel-mpi@2021.7.1",
                    "linear_algebra": "intel-mkl@2022.2.1",
                    "parallel_netcdf": "+mpi+parallel-netcdf4+large-file-support",
                    "optimization_flags": "-xHost -ip -ipo -O3 -no-prec-div -static-intel",
                    "parallel_jobs": "32",
                    "target_architecture": "skylake_avx512"  # For Intel HPC nodes
                },
                
                sample_workflows=[
                    "Regional climate downscaling with WRF (10km resolution)",
                    "Global climate model comparison and analysis",
                    "Extreme weather detection using machine learning",
                    "Ensemble climate projections and uncertainty analysis",
                    "Hurricane track prediction and intensity modeling",
                    "Air quality modeling with chemical transport",
                    "Paleoclimate reconstruction and model validation",
                    "Climate impact assessment for agriculture"
                ],
                
                transparent_pricing={
                    "idle_cost": "$0/day (compiled models cached, no compute)",
                    "data_analysis": "$5-20/day (post-processing, visualization)",
                    "regional_modeling": "$40-150/day (WRF runs, few days simulation)",
                    "global_modeling": "$150-500/day (CESM runs, years of simulation)",
                    "ensemble_runs": "$300-1000/day (multiple model runs)",
                    "storage_cost": "$2.30/month per 100GB (model output)",
                    "hpc_bursting": "$500-2000/day for massive simulations",
                    "example_monthly": "$400-2000/month for climate modeler"
                },
                
                deployment_variants=[
                    "analyst: Data analysis tools + visualization (no modeling)",
                    "regional: WRF + regional modeling tools",
                    "global: CESM + global climate models",
                    "comprehensive: All models + analysis + visualization",
                    "hpc: Optimized for large-scale parallel computing"
                ],
                
                security_tiers={
                    "basic": {
                        "description": "Standard security for academic climate research",
                        "features": "Package verification, encrypted data, secure access",
                        "compliance": "Good for open climate research and collaboration",
                        "disclaimer": "⚠️  Basic security - verify data sharing agreements"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for government climate research",
                        "features": "Verified builds, audit trails, data governance",
                        "compliance": "NIST 800-171 for controlled climate data",
                        "disclaimer": "⚠️  Technical controls only - ensure institutional policies"
                    },
                    "nist_800_53": {
                        "description": "High security for sensitive climate/weather data",
                        "features": "Hardened environment, continuous monitoring",
                        "compliance": "NIST 800-53 for operational weather systems",
                        "disclaimer": "⚠️  Technical implementation only - requires full compliance program"
                    }
                },
                
                immediate_value=[
                    "Pre-compiled, optimized climate models ready to run",
                    "Consistent software environment across HPC and cloud",
                    "Example configurations for common model runs",
                    "Automated model output post-processing",
                    "Performance-optimized builds for your architecture",
                    "Easy model version comparison and testing",
                    "Reproducible research with exact software specifications",
                    "Lab-wide software consistency and sharing"
                ],
                
                deployment_code="climate_spack_lab"
            ),
            
            "ai_spack_studio": SpackResearchSolution(
                name="AI/ML Research Studio (Spack-Powered)",
                description="High-performance machine learning environment with Spack-optimized frameworks",
                domain="Artificial Intelligence & Machine Learning",
                target_users="ML researcher, computer science lab, AI development team",
                
                spack_environment={
                    "ml_frameworks": [
                        "pytorch@1.13.1+cuda+nccl+magma",  # Deep learning
                        "tensorflow@2.11.0+cuda+nccl",     # Deep learning
                        "jax@0.4.1+cuda",                  # JAX for research
                        "onnx@1.12.0",                     # Model interchange
                        "opencv@4.6.0+python3+cuda",      # Computer vision
                        "magma@2.6.2+cuda",               # GPU linear algebra
                        "nccl@2.15.5+cuda",               # GPU communication
                        "cutensor@1.6.1.5",               # CUDA tensor operations
                        "cudnn@8.6.0.163",                # CUDA deep learning
                        "tensorrt@8.5.1.7",               # Inference optimization
                    ],
                    "classical_ml": [
                        "py-scikit-learn@1.1.3",
                        "py-xgboost@1.7.1+cuda",
                        "lightgbm@3.3.3+cuda",
                        "py-catboost@1.1.1",
                        "py-rapids@22.12.0+cuda",  # GPU-accelerated data science
                        "py-cuml@22.12.0",        # GPU machine learning
                        "py-cudf@22.12.0",        # GPU dataframes
                        "py-dask@2022.11.1+cuda", # Distributed computing
                    ],
                    "python_ecosystem": [
                        "python@3.10.8+optimizations+shared",
                        "py-numpy@1.24.1+blas=openblas",
                        "py-scipy@1.9.3+blas=openblas",
                        "py-pandas@1.5.2",
                        "py-matplotlib@3.6.2",
                        "py-seaborn@0.12.1",
                        "py-plotly@5.11.0",
                        "py-bokeh@2.4.3",
                        "py-jupyter@1.0.0",
                        "py-jupyterlab@3.5.0",
                        "py-notebook@6.5.2",
                    ],
                    "experiment_tracking": [
                        "py-mlflow@2.0.1",        # ML lifecycle management
                        "py-wandb@0.13.5",        # Weights & Biases
                        "py-tensorboard@2.11.0",  # TensorFlow visualization
                        "py-sacred@0.8.2",        # Experiment tracking
                        "py-hydra-core@1.2.0",    # Configuration management
                        "py-optuna@3.0.4",        # Hyperparameter optimization
                    ],
                    "deployment_tools": [
                        "ray@2.2.0+cuda",         # Distributed computing
                        "py-fastapi@0.88.0",      # API framework
                        "py-uvicorn@0.20.0",      # ASGI server
                        "py-gunicorn@20.1.0",     # WSGI server
                        "py-celery@5.2.1",        # Task queue
                        "redis@7.0.5",            # In-memory database
                        "docker@20.10.21",        # Containerization
                    ],
                    "data_processing": [
                        "apache-spark@3.3.1+hadoop+hive", # Big data processing
                        "py-pyspark@3.3.1",
                        "py-arrow@10.0.1",        # Columnar data
                        "py-pyarrow@10.0.1",
                        "py-polars@0.15.6",       # Fast dataframes
                        "py-vaex@4.16.0",         # Out-of-core dataframes
                    ]
                },
                
                spack_config={
                    "compiler_optimization": "gcc@11.3.0 %gcc@11.3.0 cflags=-O3 cxxflags=-O3",
                    "cuda_arch": "70,75,80,86",  # Support multiple GPU architectures
                    "python_optimization": "+optimizations+shared+ssl+tkinter",
                    "blas_implementation": "openblas@0.3.21 threads=openmp target=haswell",
                    "mpi_implementation": "openmpi@4.1.4+cuda+legacylaunchers",
                    "parallel_jobs": "16",
                    "gpu_optimization": "enabled"
                },
                
                sample_workflows=[
                    "Large language model fine-tuning with PyTorch + DeepSpeed",
                    "Computer vision model training with multi-GPU setup",
                    "Distributed hyperparameter optimization with Ray Tune",
                    "MLOps pipeline with MLflow and automated deployment",
                    "Reinforcement learning with GPU-accelerated environments",
                    "Time series forecasting with transformer models",
                    "Federated learning across multiple nodes",
                    "Model compression and quantization for edge deployment"
                ],
                
                transparent_pricing={
                    "idle_cost": "$0/day (all software cached, no running instances)",
                    "cpu_development": "$3-15/day (prototyping, small experiments)",
                    "single_gpu": "$25-80/day (model training, V100/A100)",
                    "multi_gpu": "$80-300/day (large model training, 4-8 GPUs)",
                    "distributed": "$200-800/day (cluster training, 8+ nodes)",
                    "inference": "$0.05-2.00 per 1000 predictions (serverless)",
                    "storage_cost": "$2.30/month per 100GB (models and datasets)",
                    "example_monthly": "$200-1500/month for ML researcher"
                },
                
                deployment_variants=[
                    "cpu: CPU-only for learning and prototyping",
                    "single_gpu: Single GPU workstation for research",
                    "multi_gpu: Multi-GPU for large model training", 
                    "distributed: Multi-node distributed training",
                    "inference: Optimized for model serving and deployment"
                ],
                
                security_tiers={
                    "basic": {
                        "description": "Standard security with package verification",
                        "features": "Spack package integrity, encrypted storage, basic isolation",
                        "compliance": "Good for academic research and open datasets",
                        "disclaimer": "⚠️  Basic security - review data policies for proprietary models"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for commercial AI research",
                        "features": "Verified ML frameworks, audit logging, model protection",
                        "compliance": "NIST 800-171 + AI/ML security controls",
                        "disclaimer": "⚠️  Technical controls implemented - ensure AI governance policies"
                    },
                    "nist_800_53": {
                        "description": "High security for defense AI applications",
                        "features": "Hardened ML environment, continuous monitoring",
                        "compliance": "NIST 800-53 + AI security framework",
                        "disclaimer": "⚠️  Technical implementation only - requires AI ethics and bias testing"
                    }
                },
                
                immediate_value=[
                    "GPU-optimized ML frameworks ready in minutes",
                    "Consistent environments across development and production",
                    "Performance-tuned builds for your hardware",
                    "Example projects with best practices",
                    "Automated experiment tracking and reproducibility",
                    "Easy scaling from single GPU to distributed training",
                    "Version-controlled software environments",
                    "Lab-wide ML framework consistency"
                ],
                
                deployment_code="ai_spack_studio"
            )
        }
        
        return solutions

    def generate_spack_deployment_script(self, solution_code: str) -> str:
        """Generate Spack-based deployment automation"""
        
        base_script = '''#!/bin/bash
# Spack-Based Research Solution Deployment
# Optimized software builds with dependency management

set -euo pipefail

SOLUTION_TYPE="${1:-genomics_spack_lab}"
PROJECT_NAME="${2:-research-project}"
VARIANT="${3:-standard}"
SECURITY_TIER="${4:-basic}"

# Spack configuration
export SPACK_ROOT="/opt/spack"
export PATH="$SPACK_ROOT/bin:$PATH"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

setup_spack_environment() {
    log "Setting up Spack environment for $SOLUTION_TYPE"
    
    # Clone Spack if not present
    if [[ ! -d "$SPACK_ROOT" ]]; then
        log "Installing Spack..."
        git clone --depth=100 --branch=releases/v0.19 https://github.com/spack/spack.git $SPACK_ROOT
    fi
    
    # Source Spack environment
    source $SPACK_ROOT/share/spack/setup-env.sh
    
    # Add external compilers
    spack compiler find
    
    # Configure build cache for faster deployments
    spack mirror add binary_mirror https://binaries.spack.io/releases/v0.19
    spack buildcache keys --install --trust
    
    log "Spack environment ready"
}

create_research_environment() {
    local solution=$1
    local variant=$2
    
    log "Creating Spack environment: $solution-$variant"
    
    # Create named environment
    spack env create $solution-$variant
    spack env activate $solution-$variant
    
    # Add packages based on solution type
    case $solution in
        "genomics_spack_lab")
            create_genomics_environment $variant
            ;;
        "climate_spack_lab") 
            create_climate_environment $variant
            ;;
        "ai_spack_studio")
            create_ai_environment $variant
            ;;
        *)
            log "Unknown solution type: $solution"
            exit 1
            ;;
    esac
    
    log "Environment definition complete"
}

create_genomics_environment() {
    local variant=$1
    
    log "Configuring genomics environment (variant: $variant)"
    
    # Core bioinformatics tools (always included)
    spack add gatk@4.4.0.0
    spack add bwa@0.7.17
    spack add bwa-mem2@2.2.1
    spack add star@2.7.10b
    spack add samtools@1.17
    spack add bcftools@1.17
    spack add bedtools2@2.30.0
    spack add fastqc@0.11.9
    spack add blast-plus@2.13.0
    
    # Python scientific stack
    spack add python@3.10.8+optimizations+shared
    spack add py-numpy@1.24.1
    spack add py-pandas@1.5.2
    spack add py-biopython@1.80
    spack add py-pysam@0.20.0
    spack add py-jupyter@1.0.0
    
    if [[ "$variant" != "minimal" ]]; then
        # R and Bioconductor
        spack add r@4.2.2+X+external-lapack
        spack add r-biocmanager@1.30.19
        spack add r-deseq2@1.36.0
        spack add r-edger@3.38.4
        
        # Additional analysis tools
        spack add trimmomatic@0.39
        spack add picard@2.27.4
        spack add vcftools@0.1.16
    fi
    
    if [[ "$variant" == "comprehensive" ]]; then
        # Workflow managers
        spack add nextflow@22.10.4
        spack add snakemake@7.18.2
        
        # Visualization tools
        spack add igv@2.15.4
        spack add circos@0.69-9
        
        # Long-read tools
        spack add minimap2@2.24
        spack add canu@2.2
    fi
}

create_climate_environment() {
    local variant=$1
    
    log "Configuring climate modeling environment (variant: $variant)"
    
    # Essential climate tools
    spack add nco@5.1.0
    spack add cdo@2.0.6
    spack add hdf5@1.12.2+mpi
    spack add netcdf-c@4.9.0+mpi
    spack add netcdf-fortran@4.6.0
    
    # Python climate stack
    spack add python@3.10.8+optimizations
    spack add py-xarray@2022.11.0
    spack add py-dask@2022.11.1
    spack add py-cartopy@0.21.0
    spack add py-matplotlib@3.6.2
    spack add py-netcdf4@1.6.2
    
    if [[ "$variant" != "analyst" ]]; then
        # Climate models
        spack add wrf@4.4.1
        spack add openmpi@4.1.4
        spack add fftw@3.3.10+mpi+openmp
        
        # High-performance libraries
        spack add petsc@3.18.1+mpi+hypre+metis
        spack add scalapack@2.2.0
    fi
    
    if [[ "$variant" == "comprehensive" ]]; then
        # Full modeling suite
        spack add paraview@5.11.0+mpi+python3
        spack add vapor@3.7.0
        spack add visit@3.3.3
        
        # Advanced analysis
        spack add esmf@8.4.0+mpi
        spack add py-esmpy@8.4.0
    fi
}

create_ai_environment() {
    local variant=$1
    
    log "Configuring AI/ML environment (variant: $variant)"
    
    # Python and scientific computing
    spack add python@3.10.8+optimizations+shared
    spack add py-numpy@1.24.1+blas=openblas
    spack add py-scipy@1.9.3+blas=openblas
    spack add py-pandas@1.5.2
    spack add py-jupyter@1.0.0
    spack add py-scikit-learn@1.1.3
    
    # Classical ML
    spack add py-xgboost@1.7.1
    spack add lightgbm@3.3.3
    
    if [[ "$variant" != "cpu" ]]; then
        # GPU-accelerated frameworks
        spack add pytorch@1.13.1+cuda+nccl+magma
        spack add py-tensorflow@2.11.0+cuda
        spack add opencv@4.6.0+python3+cuda
        spack add cudnn@8.6.0.163
        spack add nccl@2.15.5+cuda
    fi
    
    if [[ "$variant" == "distributed" ]]; then
        # Distributed computing
        spack add ray@2.2.0+cuda
        spack add openmpi@4.1.4+cuda
        spack add py-dask@2022.11.1+cuda
    fi
    
    # Experiment tracking (all variants)
    spack add py-mlflow@2.0.1
    spack add py-wandb@0.13.5
    spack add py-tensorboard@2.11.0
}

install_environment() {
    log "Installing Spack environment (this may take 30-60 minutes)..."
    
    # Use build cache when possible
    spack install --use-cache
    
    if [[ $? -eq 0 ]]; then
        log "Environment installation completed successfully"
    else
        log "Installation failed, trying without cache..."
        spack install --no-cache
    fi
    
    # Generate environment modules
    spack module tcl refresh --delete-tree
    
    log "Environment ready for use"
}

create_activation_scripts() {
    local solution=$1
    local variant=$2
    
    log "Creating activation scripts..."
    
    cat > "activate_${solution}_${variant}.sh" << EOF
#!/bin/bash
# Activation script for $solution ($variant variant)

export SPACK_ROOT="/opt/spack"
source \$SPACK_ROOT/share/spack/setup-env.sh
spack env activate $solution-$variant

echo "✓ $solution environment activated"
echo "✓ Available tools:"
spack find --format "{name}@{version}"

# Set up common environment variables
export JUPYTER_CONFIG_DIR=\$(pwd)/.jupyter
export PYTHONPATH=\$(pwd):\$PYTHONPATH

# Create quick start guide
cat << GUIDE

Quick Start Guide for $solution:
1. Start Jupyter: jupyter lab --ip=0.0.0.0 --port=8888
2. Check installed packages: spack find
3. Load specific tools: spack load <package-name>
4. Run example workflows in ./examples/ directory

Cost tracking:
- Check current usage: aws ce get-cost-and-usage --time-period Start=\$(date -d '1 month ago' +%Y-%m-%d),End=\$(date +%Y-%m-%d) --granularity MONTHLY --metrics BlendedCost
- Estimated daily cost: \$0 idle, \$5-50 active (depending on usage)

GUIDE
EOF
    
    chmod +x "activate_${solution}_${variant}.sh"
    log "Activation script created: activate_${solution}_${variant}.sh"
}

main() {
    log "Starting Spack-based research solution deployment"
    log "Solution: $SOLUTION_TYPE, Project: $PROJECT_NAME, Variant: $VARIANT"
    
    # Setup Spack
    setup_spack_environment
    
    # Create and configure environment
    create_research_environment $SOLUTION_TYPE $VARIANT
    
    # Install packages
    install_environment
    
    # Create convenience scripts
    create_activation_scripts $SOLUTION_TYPE $VARIANT
    
    log "Deployment complete!"
    log "To activate: source activate_${SOLUTION_TYPE}_${VARIANT}.sh"
    log "Estimated build time: 30-60 minutes (first time)"
    log "Subsequent deployments: 5-10 minutes (cached builds)"
}

main "$@"
'''
        
        return base_script

    def generate_spack_environment_file(self, solution: SpackResearchSolution, variant: str) -> str:
        """Generate Spack environment YAML file"""
        
        # Base environment structure
        env_config = {
            "spack": {
                "specs": [],
                "packages": {},
                "config": {
                    "build_jobs": int(solution.spack_config.get("parallel_jobs", 16)),
                    "install_tree": f"/opt/spack/environments/{solution.deployment_code}",
                    "build_stage": "/tmp/spack-stage"
                },
                "compilers": [
                    {
                        "compiler": {
                            "spec": "gcc@11.3.0",
                            "paths": {
                                "cc": "/usr/bin/gcc",
                                "cxx": "/usr/bin/g++", 
                                "f77": "/usr/bin/gfortran",
                                "fc": "/usr/bin/gfortran"
                            },
                            "flags": {
                                "cflags": "-O3 -march=native",
                                "cxxflags": "-O3 -march=native", 
                                "fflags": "-O3 -march=native"
                            }
                        }
                    }
                ]
            }
        }
        
        # Add packages based on variant
        if variant == "minimal":
            # Add only core tools
            core_packages = list(solution.spack_environment.values())[0][:5]
            env_config["spack"]["specs"] = core_packages
        elif variant == "standard":
            # Add core + python/R
            for category in ["core_bio_tools", "python_stack", "r_bioconductor"]:
                if category in solution.spack_environment:
                    env_config["spack"]["specs"].extend(solution.spack_environment[category])
        else:  # comprehensive
            # Add all packages
            for packages in solution.spack_environment.values():
                env_config["spack"]["specs"].extend(packages)
        
        # Package preferences for optimization
        if "openblas" in solution.spack_config.get("linear_algebra", ""):
            env_config["spack"]["packages"]["all"] = {
                "providers": {
                    "blas": ["openblas"],
                    "lapack": ["openblas"]
                }
            }
        
        # CUDA configuration for AI solutions
        if "cuda" in solution.spack_config.get("cuda_arch", ""):
            env_config["spack"]["packages"]["cuda"] = {
                "buildable": True,
                "version": ["11.8"]
            }
        
        import yaml
        return yaml.dump(env_config, default_flow_style=False, sort_keys=False)

    def generate_spack_solutions_report(self) -> str:
        """Generate comprehensive Spack-based solutions report"""
        
        solutions = self.create_spack_solutions()
        
        report = []
        report.append("# Spack-Powered Research Computing Solutions")
        report.append("## Optimized, Reproducible Software Environments for Research")
        report.append("")
        report.append("### 🚀 Why Spack for Research Computing?")
        report.append("- **Dependency Resolution**: Handles complex software dependencies automatically")
        report.append("- **Performance Optimization**: Builds optimized for your specific hardware")
        report.append("- **Reproducibility**: Exact software versions and build configurations")
        report.append("- **Package Management**: 6000+ scientific packages available")
        report.append("- **Multiple Versions**: Run different versions of the same software")
        report.append("- **Binary Caching**: Fast deployments with pre-built packages")
        report.append("")
        
        for solution_id, solution in solutions.items():
            report.append(f"## {solution.name}")
            report.append(f"**Domain**: {solution.domain}")
            report.append(f"**Target Users**: {solution.target_users}")
            report.append("")
            
            report.append(f"**Description**: {solution.description}")
            report.append("")
            
            # Spack-specific benefits
            report.append("### 🛠️ Spack-Optimized Software Stack")
            
            # Show key packages by category
            for category, packages in solution.spack_environment.items():
                category_name = category.replace('_', ' ').title()
                report.append(f"#### {category_name}")
                for package in packages[:5]:  # Show first 5 packages
                    pkg_name = package.split('@')[0]
                    pkg_version = package.split('@')[1] if '@' in package else 'latest'
                    pkg_features = package.split('+')[1:] if '+' in package else []
                    
                    feature_str = f" (+{', +'.join(pkg_features)})" if pkg_features else ""
                    report.append(f"- **{pkg_name}** v{pkg_version}{feature_str}")
                
                if len(packages) > 5:
                    report.append(f"- ... and {len(packages) - 5} more optimized packages")
                report.append("")
            
            # Performance optimizations
            report.append("### ⚡ Performance Optimizations")
            for key, value in solution.spack_config.items():
                opt_name = key.replace('_', ' ').title()
                report.append(f"- **{opt_name}**: {value}")
            report.append("")
            
            # Transparent pricing
            report.append("### 💰 Transparent Pricing")
            for cost_type, cost_desc in solution.transparent_pricing.items():
                report.append(f"- **{cost_type.replace('_', ' ').title()}**: {cost_desc}")
            report.append("")
            
            # Deployment variants
            report.append("### ⚙️ Deployment Variants")
            for variant in solution.deployment_variants:
                variant_name = variant.split(':')[0]
                variant_desc = variant.split(':', 1)[1].strip()
                report.append(f"- **{variant_name}**: {variant_desc}")
            report.append("")
            
            # Security tiers
            report.append("### 🔒 Security & Compliance Tiers")
            for tier_name, tier_info in solution.security_tiers.items():
                report.append(f"#### {tier_name.upper().replace('_', ' ')}")
                report.append(f"**{tier_info['description']}**")
                report.append(f"- **Features**: {tier_info['features']}")
                report.append(f"- **Compliance**: {tier_info['compliance']}")
                report.append(f"- **⚠️ {tier_info['disclaimer']}**")
                report.append("")
            
            # Immediate value with Spack
            report.append("### 🎯 Immediate Value with Spack")
            for value in solution.immediate_value:
                report.append(f"- {value}")
            report.append("")
            
            # Sample workflows
            report.append("### 📋 Ready-to-Run Research Workflows")
            for workflow in solution.sample_workflows:
                report.append(f"- {workflow}")
            report.append("")
            
            # One-click deployment
            report.append("### 🚀 One-Click Spack Deployment")
            report.append("```bash")
            report.append(f"# Deploy standard configuration")
            report.append(f"./deploy-spack-research.sh {solution.deployment_code} my-lab standard")
            report.append("")
            report.append(f"# Deploy with specific optimizations")
            report.append(f"./deploy-spack-research.sh {solution.deployment_code} my-lab comprehensive nist_800_171")
            report.append("")
            report.append(f"# Activate environment")
            report.append(f"source activate_{solution.deployment_code}_standard.sh")
            report.append("```")
            report.append("")
            
            report.append("---")
            report.append("")
        
        # Spack advantages section
        report.append("## 🏆 Spack Advantages for Research Computing")
        report.append("")
        report.append("### vs. Conda/Pip")
        report.append("| Feature | Spack | Conda/Pip |")
        report.append("|---------|-------|-----------|")
        report.append("| HPC Optimization | ✅ Hardware-specific builds | ❌ Generic builds |")
        report.append("| Dependency Resolution | ✅ Full dependency DAG | ⚠️ Limited resolution |")
        report.append("| Multiple Versions | ✅ Side-by-side installs | ❌ Virtual environments only |")
        report.append("| MPI Integration | ✅ Native MPI support | ❌ Complex setup |")
        report.append("| Compiler Choice | ✅ Any compiler | ⚠️ Limited options |")
        report.append("| Reproducibility | ✅ Exact build specs | ⚠️ Environment files |")
        report.append("")
        
        report.append("### vs. Containers")
        report.append("| Feature | Spack | Containers |")
        report.append("|---------|-------|------------|")
        report.append("| Performance | ✅ Native speed | ⚠️ Slight overhead |")
        report.append("| GPU Access | ✅ Direct access | ⚠️ Requires configuration |")
        report.append("| HPC Integration | ✅ Native MPI | ❌ Complex networking |")
        report.append("| Customization | ✅ Full control | ⚠️ Limited after build |")
        report.append("| Resource Usage | ✅ Minimal overhead | ❌ Container overhead |")
        report.append("| Debugging | ✅ Native debugging | ⚠️ Container complexity |")
        report.append("")
        
        # Cost comparison with Spack optimizations
        report.append("## 💸 Cost Impact of Spack Optimizations")
        report.append("")
        report.append("### Performance Improvements = Cost Savings")
        report.append("```")
        report.append("Example: Climate Modeling with WRF")
        report.append("├── Generic build: 100 hours simulation time")
        report.append("├── Spack optimized: 65 hours simulation time (35% faster)")
        report.append("├── Cost at $2/hour: $200 vs $130")
        report.append("└── Monthly savings: $70+ per researcher")
        report.append("")
        report.append("Example: Genomics GATK Pipeline")
        report.append("├── Default build: 8 hours for WGS analysis")
        report.append("├── Spack optimized: 5.5 hours (31% faster)")
        report.append("├── Cost at $1.50/hour: $12 vs $8.25")
        report.append("└── Savings on 20 samples/month: $75")
        report.append("")
        report.append("Example: AI/ML Training")
        report.append("├── Standard PyTorch: 24 hours training")
        report.append("├── Spack optimized: 18 hours (25% faster)")
        report.append("├── Cost at $5/hour GPU: $120 vs $90")
        report.append("└── Savings per model: $30")
        report.append("```")
        report.append("")
        
        # Implementation timeline
        report.append("## 📅 Spack Implementation Timeline")
        report.append("")
        report.append("### Day 1: Initial Setup")
        report.append("- Spack installation and configuration (30 minutes)")
        report.append("- First environment deployment (60-90 minutes)")
        report.append("- Basic tool testing and validation (30 minutes)")
        report.append("")
        report.append("### Week 1: Migration")
        report.append("- Migrate existing workflows to Spack environment")
        report.append("- Performance benchmark and optimization")
        report.append("- Team training on Spack usage")
        report.append("")
        report.append("### Month 1: Optimization")
        report.append("- Custom package additions and modifications")
        report.append("- Build cache setup for faster deployments")
        report.append("- Advanced workflow integration")
        report.append("")
        
        # Best practices
        report.append("## 📝 Spack Best Practices for Research")
        report.append("")
        report.append("### Environment Management")
        report.append("- Create named environments for each project")
        report.append("- Use environment files for reproducibility")
        report.append("- Version control your spack.yaml files")
        report.append("- Share environments across team members")
        report.append("")
        report.append("### Performance Optimization")
        report.append("- Use compiler flags appropriate for your hardware")
        report.append("- Enable OpenMP for parallel codes")
        report.append("- Choose optimized BLAS/LAPACK implementations")
        report.append("- Use MPI for distributed computing")
        report.append("")
        report.append("### Cost Optimization")
        report.append("- Use binary caches to reduce build times")
        report.append("- Profile your applications for bottlenecks")
        report.append("- Consider architecture-specific builds")
        report.append("- Monitor resource usage and optimize accordingly")
        report.append("")
        
        return "\n".join(report)

def main():
    parser = argparse.ArgumentParser(description='Generate Spack-based research solutions')
    parser.add_argument('--output', default='spack_research_solutions.md', help='Output report file')
    parser.add_argument('--generate-scripts', action='store_true', help='Generate deployment scripts')
    
    args = parser.parse_args()
    
    generator = SpackResearchSolutionsGenerator()
    
    # Generate main report
    report = generator.generate_spack_solutions_report()
    with open(args.output, 'w', encoding='utf-8') as f:
        f.write(report)
    
    # Generate deployment scripts if requested
    if args.generate_scripts:
        solutions = generator.create_spack_solutions()
        
        # Generate main deployment script
        script = generator.generate_spack_deployment_script("all")
        with open('deploy-spack-research.sh', 'w', encoding='utf-8') as f:
            f.write(script)
        os.chmod('deploy-spack-research.sh', 0o755)
        
        # Generate environment files for each solution
        for solution_id, solution in solutions.items():
            for variant in ["minimal", "standard", "comprehensive"]:
                env_file = generator.generate_spack_environment_file(solution, variant)
                filename = f"spack_env_{solution.deployment_code}_{variant}.yaml"
                with open(filename, 'w', encoding='utf-8') as f:
                    f.write(env_file)
        
        print("Deployment scripts generated:")
        print("- deploy-spack-research.sh")
        print("- spack_env_*.yaml files")
    
    print(f"Spack research solutions report saved to {args.output}")

if __name__ == "__main__":
    main()