#!/usr/bin/env python3
"""
Researcher-Ready Solutions Generator
Domain-specific research environments with pre-installed tools and transparent pricing
"""

import os
import json
from typing import Dict, List, Optional
from dataclasses import dataclass, asdict
import argparse
import logging

@dataclass
class ResearcherSolution:
    name: str
    description: str
    domain: str
    target_users: str
    pre_installed_tools: List[str]
    sample_workflows: List[str]
    transparent_pricing: Dict[str, str]
    deployment_variants: List[str]
    security_tiers: Dict[str, Dict[str, str]]
    immediate_value: List[str]
    deployment_code: str

class ResearcherSolutionsGenerator:
    def __init__(self):
        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)

    def create_researcher_solutions(self) -> Dict[str, ResearcherSolution]:
        """Create domain-specific solutions optimized for single researcher/lab use"""

        solutions = {
            "genomics_workbench": ResearcherSolution(
                name="Genomics Research Workbench",
                description="Complete genomics analysis environment with pre-installed pipelines and reference data",
                domain="Genomics & Bioinformatics",
                target_users="Single researcher, small lab (1-5 people), bioinformatics core",

                pre_installed_tools=[
                    "GATK 4.x (Genome Analysis Toolkit)",
                    "BWA-MEM2 (sequence alignment)",
                    "SAMtools/BCFtools (file manipulation)",
                    "STAR (RNA-seq alignment)",
                    "DESeq2/edgeR (differential expression)",
                    "BLAST+ (sequence searching)",
                    "FastQC (quality control)",
                    "MultiQC (report aggregation)",
                    "IGV (genome browser)",
                    "R/Bioconductor (200+ packages)",
                    "Python scientific stack (BioPython, pandas, etc)",
                    "Nextflow/Snakemake (workflow management)",
                    "Reference genomes (human, mouse, common model organisms)",
                    "Annotation databases (GENCODE, RefSeq, Ensembl)",
                    "Jupyter Lab with genomics kernels"
                ],

                sample_workflows=[
                    "Whole genome sequencing variant calling (GATK best practices)",
                    "RNA-seq differential expression analysis",
                    "ChIP-seq peak calling and annotation",
                    "Single-cell RNA-seq analysis (Seurat/Scanpy)",
                    "Metagenomics classification and assembly",
                    "Population genetics analysis (PLINK, VCFtools)",
                    "Structural variant detection",
                    "Pharmacogenomics analysis"
                ],

                transparent_pricing={
                    "idle_cost": "$0/day (complete shutdown when not in use)",
                    "light_analysis": "$5-15/day (1-4 hours, small datasets <10GB)",
                    "standard_analysis": "$15-40/day (4-8 hours, medium datasets 10-100GB)",
                    "intensive_analysis": "$40-120/day (8+ hours, large datasets >100GB)",
                    "storage_cost": "$2.30/month per 100GB of data stored",
                    "data_transfer": "$0.09/GB for downloads (uploads free)",
                    "gpu_acceleration": "+$20-80/day for ML workloads",
                    "example_monthly": "$150-800/month for active researcher"
                },

                deployment_variants=[
                    "basic: CPU-only, 16GB RAM, good for most analysis",
                    "compute: High-CPU, 64GB RAM, for large cohort studies",
                    "memory: High-memory, 256GB RAM, for assembly and large matrices",
                    "gpu: GPU-enabled, for deep learning and image analysis",
                    "collaborative: Multi-user setup for lab sharing"
                ],

                security_tiers={
                    "basic": {
                        "description": "Standard AWS security with encryption",
                        "features": "HTTPS, encrypted storage, basic access controls",
                        "compliance": "Good for non-sensitive research data",
                        "disclaimer": "⚠️  Basic security only - NOT suitable for PHI or sensitive data"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for controlled unclassified information",
                        "features": "MFA, audit logging, network isolation, enhanced encryption",
                        "compliance": "Implements NIST 800-171 controls",
                        "disclaimer": "⚠️  Technical controls implemented - YOU MUST ensure procedural compliance, training, and documentation for full NIST 800-171 compliance"
                    },
                    "nist_800_53": {
                        "description": "High security for sensitive research and federal compliance",
                        "features": "Continuous monitoring, advanced threat detection, strict access controls",
                        "compliance": "Implements NIST 800-53 moderate baseline",
                        "disclaimer": "⚠️  Technical controls implemented - YOU MUST ensure complete compliance program including policies, procedures, training, and ongoing assessment for NIST 800-53 compliance"
                    }
                },

                immediate_value=[
                    "Launch analysis in 2 minutes - no software installation",
                    "Reference genomes pre-downloaded and indexed",
                    "Example datasets and tutorials included",
                    "Common workflows pre-configured and tested",
                    "Automatic result visualization and reporting",
                    "Collaboration features for lab sharing",
                    "Automatic backup of analysis results",
                    "Cost tracking by project/grant"
                ],

                deployment_code="genomics_workbench"
            ),

            "climate_modeling_lab": ResearcherSolution(
                name="Climate & Atmospheric Modeling Laboratory",
                description="Complete climate modeling environment with pre-configured models and analysis tools",
                domain="Climate Science & Atmospheric Physics",
                target_users="Climate researcher, atmospheric scientist, environmental modeling group",

                pre_installed_tools=[
                    "WRF (Weather Research and Forecasting Model)",
                    "CESM (Community Earth System Model)",
                    "GFDL models (climate simulation)",
                    "NCO/CDO (netCDF manipulation)",
                    "Python climate stack (xarray, cartopy, metpy, etc)",
                    "R climate packages (raster, ncdf4, climate4R)",
                    "MATLAB climate toolboxes",
                    "ParaView (3D visualization)",
                    "Panoply (netCDF viewer)",
                    "GrADS (atmospheric data analysis)",
                    "Ferret (climate data analysis)",
                    "VAPOR (3D atmospheric visualization)",
                    "Jupyter Lab with climate kernels",
                    "Climate data servers (THREDDS, ERDDAP)",
                    "Reanalysis datasets (ERA5, NCEP, etc)"
                ],

                sample_workflows=[
                    "Regional climate downscaling with WRF",
                    "Global climate model analysis and comparison",
                    "Extreme weather event detection and attribution",
                    "Climate change impact assessment",
                    "Atmospheric transport modeling",
                    "Hurricane/typhoon track prediction",
                    "Air quality modeling and forecasting",
                    "Paleoclimate reconstruction and analysis"
                ],

                transparent_pricing={
                    "idle_cost": "$0/day (models and data cached but compute off)",
                    "data_analysis": "$8-25/day (exploratory analysis, visualization)",
                    "model_runs": "$50-200/day (regional simulations, 1-10 year runs)",
                    "large_simulations": "$200-800/day (global models, century runs)",
                    "storage_cost": "$2.30/month per 100GB (model output and datasets)",
                    "data_access": "Free access to major reanalysis datasets",
                    "hpc_bursting": "$500-2000/day for extreme-scale simulations",
                    "example_monthly": "$300-1500/month for active modeler"
                },

                deployment_variants=[
                    "analyst: Focus on data analysis and visualization",
                    "modeler: Includes model compilation and moderate compute",
                    "hpc: High-performance setup for large simulations",
                    "ensemble: Multi-model comparison and ensemble runs",
                    "realtime: Operational forecasting capabilities"
                ],

                security_tiers={
                    "basic": {
                        "description": "Standard security for academic research",
                        "features": "Encrypted data, secure access, basic monitoring",
                        "compliance": "Suitable for open research and publication",
                        "disclaimer": "⚠️  Basic security - verify data sharing agreements for collaborative projects"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for government-funded research",
                        "features": "Strong access controls, audit trails, data governance",
                        "compliance": "NIST 800-171 technical controls implemented",
                        "disclaimer": "⚠️  Technical implementation only - ensure your institution has proper policies and training for full compliance"
                    },
                    "nist_800_53": {
                        "description": "High security for sensitive climate data",
                        "features": "Continuous monitoring, advanced security, strict controls",
                        "compliance": "NIST 800-53 moderate baseline implemented",
                        "disclaimer": "⚠️  Technical controls only - requires comprehensive compliance program for full NIST 800-53 compliance"
                    }
                },

                immediate_value=[
                    "Pre-configured climate models ready to run",
                    "Major datasets already cached and accessible",
                    "Example simulations and analysis notebooks",
                    "Automated visualization and reporting",
                    "Collaboration tools for multi-institution projects",
                    "Publication-quality figure generation",
                    "Automatic metadata and provenance tracking",
                    "Grant proposal cost estimation tools"
                ],

                deployment_code="climate_modeling_lab"
            ),

            "ai_research_studio": ResearcherSolution(
                name="AI/ML Research Studio",
                description="Complete machine learning research environment with latest frameworks and pre-trained models",
                domain="Artificial Intelligence & Machine Learning",
                target_users="ML researcher, computer science lab, AI startup team",

                pre_installed_tools=[
                    "PyTorch (latest stable + nightly builds)",
                    "TensorFlow/Keras (2.x and legacy 1.x)",
                    "JAX/Flax (Google's ML framework)",
                    "Hugging Face Transformers (1000+ pre-trained models)",
                    "OpenCV (computer vision)",
                    "scikit-learn (classical ML)",
                    "XGBoost/LightGBM (gradient boosting)",
                    "Ray (distributed computing)",
                    "Weights & Biases (experiment tracking)",
                    "MLflow (model lifecycle management)",
                    "Jupyter Lab + VS Code integration",
                    "CUDA toolkit (latest drivers)",
                    "cuDNN (optimized deep learning primitives)",
                    "Apex (mixed precision training)",
                    "Common datasets (ImageNet, COCO, GLUE, etc)"
                ],

                sample_workflows=[
                    "Large language model fine-tuning (BERT, GPT, T5)",
                    "Computer vision model training (ResNet, YOLO, transformers)",
                    "Reinforcement learning (stable-baselines3, Ray RLlib)",
                    "Generative AI (GANs, VAEs, diffusion models)",
                    "Time series forecasting and anomaly detection",
                    "Federated learning experiments",
                    "Model optimization and quantization",
                    "Multi-modal learning (vision + language)"
                ],

                transparent_pricing={
                    "idle_cost": "$0/day (notebooks shut down automatically)",
                    "development": "$5-20/day (CPU development, small experiments)",
                    "gpu_training": "$30-150/day (single GPU training, V100/A100)",
                    "multi_gpu": "$100-500/day (multi-GPU training for large models)",
                    "inference": "$0.10-5.00 per 1000 predictions (serverless)",
                    "storage_cost": "$2.30/month per 100GB (models and datasets)",
                    "model_serving": "$20-200/day for production endpoints",
                    "example_monthly": "$200-1000/month for active researcher"
                },

                deployment_variants=[
                    "starter: CPU-only for learning and small experiments",
                    "researcher: Single GPU for standard research",
                    "advanced: Multi-GPU for large model training",
                    "production: Includes model serving and monitoring",
                    "collaborative: Team workspace with shared resources"
                ],

                security_tiers={
                    "basic": {
                        "description": "Standard security for academic research",
                        "features": "Encrypted storage, HTTPS access, basic isolation",
                        "compliance": "Good for open research and model sharing",
                        "disclaimer": "⚠️  Basic security - review data policies for proprietary datasets"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for commercial ML research",
                        "features": "Strong access controls, audit logging, model protection",
                        "compliance": "NIST 800-171 controls for IP protection",
                        "disclaimer": "⚠️  Technical controls implemented - ensure proper data governance and IP policies for compliance"
                    },
                    "nist_800_53": {
                        "description": "High security for sensitive AI research",
                        "features": "Advanced monitoring, model security, strict access controls",
                        "compliance": "NIST 800-53 for government/defense AI research",
                        "disclaimer": "⚠️  Technical implementation only - requires full compliance program including AI ethics and bias testing"
                    }
                },

                immediate_value=[
                    "Latest ML frameworks pre-installed and configured",
                    "1000+ pre-trained models ready to fine-tune",
                    "Example projects and tutorials for common tasks",
                    "Automatic hyperparameter tuning",
                    "Model performance monitoring and comparison",
                    "One-click model deployment and serving",
                    "Collaboration features for team research",
                    "Automatic experiment tracking and reproducibility"
                ],

                deployment_code="ai_research_studio"
            ),

            "materials_simulation_lab": ResearcherSolution(
                name="Materials Science Simulation Laboratory",
                description="Comprehensive materials modeling environment with quantum chemistry and molecular dynamics tools",
                domain="Materials Science & Computational Chemistry",
                target_users="Materials scientist, computational chemist, physics researcher",

                pre_installed_tools=[
                    "VASP (Vienna Ab initio Simulation Package)",
                    "Quantum ESPRESSO (DFT calculations)",
                    "LAMMPS (molecular dynamics)",
                    "GROMACS (biomolecular simulations)",
                    "AMBER (molecular dynamics)",
                    "Gaussian (quantum chemistry)",
                    "OpenMX (density functional theory)",
                    "ASE (Atomic Simulation Environment)",
                    "pymatgen (materials analysis)",
                    "OVITO (visualization)",
                    "VMD (molecular visualization)",
                    "Crystal Maker (structure visualization)",
                    "Phonopy (phonon calculations)",
                    "SeeK-path (band structure)",
                    "Materials Project API access"
                ],

                sample_workflows=[
                    "DFT electronic structure calculations",
                    "Molecular dynamics of polymers and proteins",
                    "Phase transition and thermodynamics",
                    "Mechanical properties prediction",
                    "Battery materials screening",
                    "Catalyst design and optimization",
                    "Defect formation energy calculations",
                    "High-throughput materials discovery"
                ],

                transparent_pricing={
                    "idle_cost": "$0/day (simulations paused, data preserved)",
                    "small_systems": "$10-30/day (DFT on <100 atoms)",
                    "medium_systems": "$50-150/day (MD simulations, 1000s of atoms)",
                    "large_systems": "$200-600/day (extended simulations, large cells)",
                    "hpc_calculations": "$500-2000/day (massive parallel simulations)",
                    "storage_cost": "$2.30/month per 100GB (structures and trajectories)",
                    "licensing": "Open source tools included, commercial licenses separate",
                    "example_monthly": "$400-2000/month for computational materials researcher"
                },

                deployment_variants=[
                    "quantum: Focus on DFT and electronic structure",
                    "molecular: Emphasis on MD and classical simulations",
                    "screening: High-throughput materials discovery",
                    "multiscale: Integration of quantum and classical methods",
                    "collaborative: Multi-user lab environment"
                ],

                security_tiers={
                    "basic": {
                        "description": "Standard academic research security",
                        "features": "Encrypted data, secure access, computation logging",
                        "compliance": "Suitable for fundamental research",
                        "disclaimer": "⚠️  Basic security - check export control requirements for advanced materials research"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for sponsored research",
                        "features": "Access controls, audit trails, data protection",
                        "compliance": "NIST 800-171 for industry partnerships",
                        "disclaimer": "⚠️  Technical controls implemented - ensure proper IP protection and export control compliance"
                    },
                    "nist_800_53": {
                        "description": "High security for defense materials research",
                        "features": "Strict monitoring, advanced security, export control integration",
                        "compliance": "NIST 800-53 for government/defense projects",
                        "disclaimer": "⚠️  Technical implementation only - requires comprehensive security program including export control and classification management"
                    }
                },

                immediate_value=[
                    "Pre-compiled simulation packages with optimizations",
                    "Common crystal structures and force fields included",
                    "Example calculations and input files",
                    "Automated job management and queuing",
                    "Real-time progress monitoring and visualization",
                    "Publication-ready figure generation",
                    "Integration with materials databases",
                    "Collaborative analysis and sharing tools"
                ],

                deployment_code="materials_simulation_lab"
            ),

            "digital_humanities_workspace": ResearcherSolution(
                name="Digital Humanities Research Workspace",
                description="Integrated environment for text analysis, data visualization, and digital scholarship",
                domain="Digital Humanities & Social Sciences",
                target_users="Digital humanities scholar, social scientist, historian with digital methods",

                pre_installed_tools=[
                    "R + humanities packages (quanteda, tidytext, stylo)",
                    "Python text analysis (NLTK, spaCy, gensim, transformers)",
                    "Voyant Tools (web-based text analysis)",
                    "OpenRefine (data cleaning and transformation)",
                    "Gephi (network analysis and visualization)",
                    "QGIS (geographic information systems)",
                    "Palladio (humanities data visualization)",
                    "Tropy (photo management and annotation)",
                    "Omeka S (digital collections)",
                    "TEI Publisher (text encoding)",
                    "Jupyter Lab with humanities kernels",
                    "Tableau Public (data visualization)",
                    "NodeGoat (data management)",
                    "HathiTrust Research Center tools",
                    "Internet Archive API access"
                ],

                sample_workflows=[
                    "Large-scale text analysis and topic modeling",
                    "Historical newspaper digitization and analysis",
                    "Social network analysis of historical figures",
                    "Geographic visualization of cultural phenomena",
                    "Digital edition creation and publication",
                    "Computational stylometry and authorship attribution",
                    "Sentiment analysis of historical documents",
                    "Interactive timeline and map creation"
                ],

                transparent_pricing={
                    "idle_cost": "$0/day (workspaces sleep when not in use)",
                    "light_research": "$3-8/day (text analysis, small datasets)",
                    "intensive_analysis": "$15-40/day (large corpora, complex visualizations)",
                    "collaborative_project": "$25-75/day (team workspace, shared resources)",
                    "storage_cost": "$2.30/month per 100GB (texts, images, databases)",
                    "publishing": "$5-20/month for web publication hosting",
                    "api_usage": "Varies by external service (HathiTrust, etc)",
                    "example_monthly": "$100-600/month for active digital humanities project"
                },

                deployment_variants=[
                    "textual: Focus on text analysis and corpus linguistics",
                    "spatial: Emphasis on GIS and geographic analysis",
                    "visual: Image analysis and multimedia scholarship",
                    "collaborative: Multi-user project workspace",
                    "publishing: Includes web publishing and dissemination tools"
                ],

                security_tiers={
                    "basic": {
                        "description": "Standard security for open humanities research",
                        "features": "Basic access controls, encrypted storage",
                        "compliance": "Appropriate for most humanities research",
                        "disclaimer": "⚠️  Basic security - review copyright and privacy requirements for your materials"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for sensitive cultural materials",
                        "features": "Strong access controls, audit logging, data governance",
                        "compliance": "NIST 800-171 for protected cultural heritage",
                        "disclaimer": "⚠️  Technical controls implemented - ensure proper cultural sensitivity and privacy protections"
                    },
                    "nist_800_53": {
                        "description": "High security for government archives and sensitive collections",
                        "features": "Advanced monitoring, strict access controls, classification support",
                        "compliance": "NIST 800-53 for federal cultural institutions",
                        "disclaimer": "⚠️  Technical implementation only - requires full compliance program including cultural protocols and classification procedures"
                    }
                },

                immediate_value=[
                    "Pre-configured text analysis pipelines",
                    "Sample datasets and example projects",
                    "Templates for common digital humanities methods",
                    "Integrated visualization and publishing tools",
                    "Collaboration features for team projects",
                    "Automatic citation and metadata management",
                    "Web publishing capabilities",
                    "Integration with major digital libraries"
                ],

                deployment_code="digital_humanities_workspace"
            ),

            "neuroscience_analysis_platform": ResearcherSolution(
                name="Neuroscience Data Analysis Platform",
                description="Comprehensive environment for neuroimaging, electrophysiology, and computational neuroscience",
                domain="Neuroscience & Brain Imaging",
                target_users="Neuroscientist, cognitive researcher, medical imaging specialist",

                pre_installed_tools=[
                    "FSL (FMRIB Software Library)",
                    "FreeSurfer (brain surface reconstruction)",
                    "AFNI (Analysis of Functional NeuroImages)",
                    "SPM (Statistical Parametric Mapping)",
                    "ANTs (Advanced Normalization Tools)",
                    "Nipype (neuroimaging pipelines)",
                    "MNE-Python (EEG/MEG analysis)",
                    "EEGLAB (EEG processing)",
                    "FieldTrip (MEG/EEG/LFP analysis)",
                    "Brian2 (spiking neural networks)",
                    "NEURON (computational neuroscience)",
                    "3D Slicer (medical image analysis)",
                    "ITK-SNAP (segmentation)",
                    "Connectome Workbench (brain connectivity)",
                    "Python neuroscience stack (nilearn, nibabel, etc)"
                ],

                sample_workflows=[
                    "fMRI preprocessing and statistical analysis",
                    "Structural brain imaging and morphometry",
                    "Diffusion tensor imaging and tractography",
                    "EEG/MEG source reconstruction",
                    "Brain network connectivity analysis",
                    "Single-cell electrophysiology analysis",
                    "Computational modeling of neural circuits",
                    "Machine learning for brain decoding"
                ],

                transparent_pricing={
                    "idle_cost": "$0/day (imaging data preserved, compute off)",
                    "basic_analysis": "$8-25/day (single subject analysis)",
                    "group_studies": "$25-75/day (multi-subject statistics)",
                    "advanced_modeling": "$50-200/day (computational neuroscience)",
                    "gpu_acceleration": "$75-300/day (deep learning, real-time analysis)",
                    "storage_cost": "$2.30/month per 100GB (imaging data and results)",
                    "large_datasets": "+$50-150/day for population studies",
                    "example_monthly": "$250-1200/month for neuroimaging researcher"
                },

                deployment_variants=[
                    "imaging: Focus on MRI/fMRI analysis pipelines",
                    "electrophysiology: EEG/MEG and single-cell analysis",
                    "computational: Neural modeling and simulation",
                    "clinical: HIPAA-ready for patient data",
                    "collaborative: Multi-site neuroimaging studies"
                ],

                security_tiers={
                    "basic": {
                        "description": "Standard security for animal and anonymized human research",
                        "features": "Encrypted storage, secure access, basic compliance",
                        "compliance": "Suitable for de-identified research data",
                        "disclaimer": "⚠️  Basic security - NOT suitable for identifiable patient data or PHI"
                    },
                    "nist_800_171": {
                        "description": "Enhanced security for sensitive research data",
                        "features": "Strong access controls, audit logging, data governance",
                        "compliance": "NIST 800-171 for controlled research data",
                        "disclaimer": "⚠️  Technical controls implemented - ensure proper IRB approval and data use agreements for compliance"
                    },
                    "nist_800_53": {
                        "description": "High security for clinical data and federal research",
                        "features": "Advanced monitoring, strict controls, HIPAA alignment",
                        "compliance": "NIST 800-53 + HIPAA technical safeguards",
                        "disclaimer": "⚠️  Technical implementation only - requires comprehensive compliance program including HIPAA policies, training, and business associate agreements"
                    }
                },

                immediate_value=[
                    "Pre-configured neuroimaging pipelines",
                    "Standard brain atlases and templates included",
                    "Example datasets and analysis notebooks",
                    "Automated quality control and reporting",
                    "Publication-ready visualization tools",
                    "Collaboration features for multi-site studies",
                    "Integration with neuroimaging databases",
                    "BIDS (Brain Imaging Data Structure) compliance"
                ],

                deployment_code="neuroscience_analysis_platform"
            )
        }

        return solutions

    def generate_deployment_script(self, solution_code: str) -> str:
        """Generate domain-specific deployment script"""

        deployment_scripts = {
            "genomics_workbench": self._generate_genomics_deployment(),
            "climate_modeling_lab": self._generate_climate_deployment(),
            "ai_research_studio": self._generate_ai_deployment(),
            "materials_simulation_lab": self._generate_materials_deployment(),
            "digital_humanities_workspace": self._generate_humanities_deployment(),
            "neuroscience_analysis_platform": self._generate_neuroscience_deployment()
        }

        return deployment_scripts.get(solution_code, "# Deployment script not found")

    def _generate_genomics_deployment(self) -> str:
        return '''
# Genomics Research Workbench Deployment
# Pre-configured with GATK, BWA, STAR, R/Bioconductor, and reference genomes

variable "genomics_tools" {
  description = "Genomics software stack configuration"
  type = object({
    include_gatk = bool
    include_star = bool
    include_reference_genomes = list(string)
    ram_gb = number
    storage_gb = number
  })
  default = {
    include_gatk = true
    include_star = true
    include_reference_genomes = ["hg38", "mm10"]
    ram_gb = 32
    storage_gb = 500
  }
}

# Custom AMI with genomics tools
data "aws_ami" "genomics_workbench" {
  most_recent = true
  owners      = ["self"]

  filter {
    name   = "name"
    values = ["genomics-workbench-*"]
  }
}

# SageMaker notebook instance for interactive analysis
resource "aws_sagemaker_notebook_instance" "genomics_workbench" {
  name                    = "${var.project_name}-genomics-workbench"
  role_arn               = aws_iam_role.sagemaker_role.arn
  instance_type          = var.genomics_tools.ram_gb <= 16 ? "ml.t3.xlarge" : "ml.m5.2xlarge"
  platform_identifier    = "notebook-al2-v1"
  volume_size            = var.genomics_tools.storage_gb

  lifecycle_config_name = aws_sagemaker_notebook_instance_lifecycle_configuration.genomics_setup.name

  tags = {
    Tools = "GATK,BWA,STAR,R-Bioconductor"
    Domain = "Genomics"
  }
}

# Lifecycle configuration for tool installation
resource "aws_sagemaker_notebook_instance_lifecycle_configuration" "genomics_setup" {
  name = "${var.project_name}-genomics-setup"

  on_create = base64encode(<<-EOF
    #!/bin/bash
    set -e

    # Install conda environment with genomics tools
    sudo -u ec2-user -i <<'USEREOF'
    source /home/ec2-user/anaconda3/bin/activate
    conda create -n genomics python=3.9 -y
    conda activate genomics

    # Install bioinformatics tools
    conda install -c bioconda gatk4 bwa star samtools bcftools fastqc multiqc -y
    conda install -c bioconda blast igv jupyter -y
    conda install -c conda-forge pandas numpy scipy matplotlib seaborn -y

    # Install R and Bioconductor
    conda install -c conda-forge r-base r-essentials -y
    Rscript -e "install.packages('BiocManager'); BiocManager::install(c('DESeq2', 'edgeR', 'GenomicRanges'))"

    # Download reference genomes
    mkdir -p /home/ec2-user/SageMaker/references
    cd /home/ec2-user/SageMaker/references

    # Download human reference (example)
    if [[ "${join(",", var.genomics_tools.include_reference_genomes)}" =~ "hg38" ]]; then
      wget -q ftp://ftp.ncbi.nlm.nih.gov/genomes/all/GCA/000/001/405/GCA_000001405.15_GRCh38/seqs_for_alignment_pipelines.ucsc_ids/GCA_000001405.15_GRCh38_no_alt_analysis_set.fna.gz
      gunzip *.fna.gz
      bwa index *.fna
    fi

    # Install Jupyter kernels
    python -m ipykernel install --user --name genomics --display-name "Genomics Analysis"

    # Create example notebooks
    mkdir -p /home/ec2-user/SageMaker/examples
    cat > /home/ec2-user/SageMaker/examples/variant_calling_example.ipynb << 'NOTEBOOK'
{
 "cells": [
  {
   "cell_type": "markdown",
   "metadata": {},
   "source": [
    "# GATK Variant Calling Pipeline\\n",
    "Example workflow for whole genome sequencing variant calling"
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "source": [
    "# This notebook demonstrates:\\n",
    "# 1. Quality control with FastQC\\n",
    "# 2. Alignment with BWA-MEM\\n",
    "# 3. Variant calling with GATK\\n",
    "# 4. Annotation and filtering\\n"
   ]
  }
 ]
}
NOTEBOOK

USEREOF

    echo "Genomics workbench setup complete"
    EOF
  )
}

# S3 bucket for genomics data
resource "aws_s3_bucket" "genomics_data" {
  bucket = "${var.project_name}-genomics-data-${random_id.bucket_suffix.hex}"
}

# Cost tracking dashboard
resource "aws_cloudwatch_dashboard" "genomics_costs" {
  dashboard_name = "${var.project_name}-genomics-costs"

  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/SageMaker", "NotebookInstanceRunning", "NotebookInstanceName", aws_sagemaker_notebook_instance.genomics_workbench.name]
          ]
          period = 300
          stat   = "Average"
          region = data.aws_region.current.name
          title  = "Genomics Workbench Usage"
        }
      }
    ]
  })
}

output "genomics_workbench_url" {
  description = "Jupyter notebook URL for genomics analysis"
  value       = "https://${aws_sagemaker_notebook_instance.genomics_workbench.name}.notebook.${data.aws_region.current.name}.sagemaker.aws/"
}

output "genomics_cost_estimate" {
  description = "Daily cost estimate for genomics workbench"
  value       = "Idle: $0/day, Active: $15-40/day depending on usage"
}
'''

    def generate_transparent_pricing_calculator(self) -> str:
        """Generate cost calculator for transparent pricing"""
        return '''
<!DOCTYPE html>
<html>
<head>
    <title>Research Solution Cost Calculator</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .solution { border: 1px solid #ddd; margin: 20px 0; padding: 15px; border-radius: 5px; }
        .cost-breakdown { background-color: #f9f9f9; padding: 10px; margin: 10px 0; }
        .warning { background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; margin: 10px 0; }
        input, select { margin: 5px; padding: 5px; }
        .result { font-weight: bold; color: #27ae60; }
    </style>
</head>
<body>
    <h1>Research Computing Solution Cost Calculator</h1>

    <div class="solution">
        <h2>Genomics Research Workbench</h2>
        <label>Daily usage hours: <input type="number" id="genomics_hours" value="4" min="0" max="24"></label><br>
        <label>Dataset size:
            <select id="genomics_size">
                <option value="small">Small (&lt;10GB)</option>
                <option value="medium">Medium (10-100GB)</option>
                <option value="large">Large (&gt;100GB)</option>
            </select>
        </label><br>
        <label>Analysis type:
            <select id="genomics_type">
                <option value="basic">Basic analysis</option>
                <option value="standard">Standard pipeline</option>
                <option value="intensive">Intensive compute</option>
            </select>
        </label><br>
        <div class="cost-breakdown" id="genomics_cost"></div>
    </div>

    <div class="solution">
        <h2>Climate Modeling Laboratory</h2>
        <label>Simulation hours/day: <input type="number" id="climate_hours" value="6" min="0" max="24"></label><br>
        <label>Model complexity:
            <select id="climate_complexity">
                <option value="regional">Regional model</option>
                <option value="global">Global model</option>
                <option value="ensemble">Ensemble runs</option>
            </select>
        </label><br>
        <div class="cost-breakdown" id="climate_cost"></div>
    </div>

    <div class="solution">
        <h2>AI/ML Research Studio</h2>
        <label>Development hours/day: <input type="number" id="ai_hours" value="5" min="0" max="24"></label><br>
        <label>Compute type:
            <select id="ai_compute">
                <option value="cpu">CPU only</option>
                <option value="single_gpu">Single GPU</option>
                <option value="multi_gpu">Multi-GPU</option>
            </select>
        </label><br>
        <div class="cost-breakdown" id="ai_cost"></div>
    </div>

    <div class="warning">
        <strong>Important:</strong> These are estimates only. Actual costs depend on specific usage patterns, data sizes, and compute requirements. All solutions have $0 idle cost when not in use.
    </div>

    <script>
        function calculateGenomicsCost() {
            const hours = document.getElementById('genomics_hours').value;
            const size = document.getElementById('genomics_size').value;
            const type = document.getElementById('genomics_type').value;

            let baseRate = 3; // Base hourly rate
            if (size === 'medium') baseRate *= 1.5;
            if (size === 'large') baseRate *= 3;
            if (type === 'standard') baseRate *= 1.5;
            if (type === 'intensive') baseRate *= 3;

            const dailyCost = hours * baseRate;
            const monthlyCost = dailyCost * 22; // Assuming 22 working days

            document.getElementById('genomics_cost').innerHTML =
                `<strong>Daily cost: $${dailyCost.toFixed(2)}</strong><br>
                 Monthly estimate: $${monthlyCost.toFixed(2)}<br>
                 Storage: $2.30/month per 100GB<br>
                 Idle cost: $0/day`;
        }

        function calculateClimateCost() {
            const hours = document.getElementById('climate_hours').value;
            const complexity = document.getElementById('climate_complexity').value;

            let baseRate = 8; // Base hourly rate
            if (complexity === 'global') baseRate *= 2;
            if (complexity === 'ensemble') baseRate *= 4;

            const dailyCost = hours * baseRate;
            const monthlyCost = dailyCost * 22;

            document.getElementById('climate_cost').innerHTML =
                `<strong>Daily cost: $${dailyCost.toFixed(2)}</strong><br>
                 Monthly estimate: $${monthlyCost.toFixed(2)}<br>
                 Storage: $2.30/month per 100GB<br>
                 Idle cost: $0/day`;
        }

        function calculateAICost() {
            const hours = document.getElementById('ai_hours').value;
            const compute = document.getElementById('ai_compute').value;

            let baseRate = 2; // CPU base rate
            if (compute === 'single_gpu') baseRate = 15;
            if (compute === 'multi_gpu') baseRate = 50;

            const dailyCost = hours * baseRate;
            const monthlyCost = dailyCost * 22;

            document.getElementById('ai_cost').innerHTML =
                `<strong>Daily cost: $${dailyCost.toFixed(2)}</strong><br>
                 Monthly estimate: $${monthlyCost.toFixed(2)}<br>
                 Model storage: $2.30/month per 100GB<br>
                 Idle cost: $0/day`;
        }

        // Initial calculations
        calculateGenomicsCost();
        calculateClimateCost();
        calculateAICost();

        // Update on changes
        document.getElementById('genomics_hours').addEventListener('input', calculateGenomicsCost);
        document.getElementById('genomics_size').addEventListener('change', calculateGenomicsCost);
        document.getElementById('genomics_type').addEventListener('change', calculateGenomicsCost);
        document.getElementById('climate_hours').addEventListener('input', calculateClimateCost);
        document.getElementById('climate_complexity').addEventListener('change', calculateClimateCost);
        document.getElementById('ai_hours').addEventListener('input', calculateAICost);
        document.getElementById('ai_compute').addEventListener('change', calculateAICost);
    </script>
</body>
</html>
'''

    def generate_researcher_solutions_report(self) -> str:
        """Generate comprehensive report of researcher-ready solutions"""

        solutions = self.create_researcher_solutions()

        report = []
        report.append("# Researcher-Ready Computing Solutions")
        report.append("## Domain-Specific Environments with Pre-Installed Tools and Transparent Pricing")
        report.append("")
        report.append("### 🎯 Design Philosophy")
        report.append("- **Immediately Useful**: Pre-configured with domain-specific tools and workflows")
        report.append("- **Transparent Pricing**: Clear cost breakdown by usage pattern")
        report.append("- **Single Researcher/Lab Scale**: Optimized for 1-5 person research groups")
        report.append("- **Security Options**: Basic, NIST 800-171, or NIST 800-53 compliance levels")
        report.append("- **Zero Learning Curve**: Launch and start research in minutes")
        report.append("")

        for solution_id, solution in solutions.items():
            report.append(f"## {solution.name}")
            report.append(f"**Domain**: {solution.domain}")
            report.append(f"**Target Users**: {solution.target_users}")
            report.append("")

            # Description and immediate value
            report.append(f"**Description**: {solution.description}")
            report.append("")
            report.append("### 🚀 Immediate Value Proposition")
            for value in solution.immediate_value:
                report.append(f"- {value}")
            report.append("")

            # Pre-installed tools
            report.append("### 🛠️ Pre-Installed Research Tools")
            for tool in solution.pre_installed_tools[:10]:  # Show first 10
                report.append(f"- {tool}")
            if len(solution.pre_installed_tools) > 10:
                report.append(f"- ... and {len(solution.pre_installed_tools) - 10} more tools")
            report.append("")

            # Sample workflows
            report.append("### 📋 Ready-to-Run Workflows")
            for workflow in solution.sample_workflows:
                report.append(f"- {workflow}")
            report.append("")

            # Transparent pricing
            report.append("### 💰 Transparent Pricing")
            for cost_type, cost_desc in solution.transparent_pricing.items():
                report.append(f"- **{cost_type.replace('_', ' ').title()}**: {cost_desc}")
            report.append("")

            # Deployment variants
            report.append("### ⚙️ Deployment Options")
            for variant in solution.deployment_variants:
                report.append(f"- **{variant.split(':')[0]}**: {variant.split(':', 1)[1].strip()}")
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

            # Deployment command
            report.append("### 🚀 One-Click Deployment")
            report.append("```bash")
            report.append(f"# Basic deployment")
            report.append(f"./deploy-research-solution.sh deploy {solution.deployment_code} my-research-lab")
            report.append("")
            report.append(f"# With specific variant and security tier")
            report.append(f"./deploy-research-solution.sh deploy {solution.deployment_code} my-lab \\")
            report.append(f"  --variant researcher --security nist_800_171 --budget 500")
            report.append("```")
            report.append("")

            report.append("---")
            report.append("")

        # Cost comparison section
        report.append("## 💸 Cost Comparison: Traditional vs Researcher-Ready Solutions")
        report.append("")
        report.append("### Traditional Academic Computing Costs (Annual)")
        report.append("```")
        report.append("Genomics Lab:")
        report.append("├── Dedicated servers: $25,000")
        report.append("├── Software licenses: $15,000")
        report.append("├── Storage: $8,000")
        report.append("├── IT support: $12,000")
        report.append("└── Total: $60,000/year")
        report.append("")
        report.append("Climate Modeling Group:")
        report.append("├── HPC cluster: $40,000")
        report.append("├── Software: $20,000")
        report.append("├── Storage: $15,000")
        report.append("├── Maintenance: $18,000")
        report.append("└── Total: $93,000/year")
        report.append("")
        report.append("AI Research Lab:")
        report.append("├── GPU servers: $35,000")
        report.append("├── Software/cloud: $25,000")
        report.append("├── Storage: $10,000")
        report.append("├── Support: $15,000")
        report.append("└── Total: $85,000/year")
        report.append("```")
        report.append("")

        report.append("### Researcher-Ready Solution Costs (Annual)")
        report.append("```")
        report.append("Genomics Workbench:")
        report.append("├── Active usage (150 days): $4,500")
        report.append("├── Storage (500GB): $138")
        report.append("├── No software licenses: $0")
        report.append("├── No IT support needed: $0")
        report.append("└── Total: $4,638/year (92% savings)")
        report.append("")
        report.append("Climate Modeling Lab:")
        report.append("├── Simulation time (100 days): $12,000")
        report.append("├── Storage (1TB): $276")
        report.append("├── No infrastructure: $0")
        report.append("├── No maintenance: $0")
        report.append("└── Total: $12,276/year (87% savings)")
        report.append("")
        report.append("AI Research Studio:")
        report.append("├── Development (200 days): $8,000")
        report.append("├── Training (50 days GPU): $6,000")
        report.append("├── Storage (200GB): $55")
        report.append("├── No hardware refresh: $0")
        report.append("└── Total: $14,055/year (83% savings)")
        report.append("```")
        report.append("")

        # Implementation roadmap
        report.append("## 🗓️ Implementation Roadmap for Research Groups")
        report.append("")
        report.append("### Week 1: Setup and Orientation")
        report.append("- **Day 1**: Deploy first solution relevant to your research")
        report.append("- **Day 2-3**: Team onboarding and training on tools")
        report.append("- **Day 4-5**: Migrate first analysis workflow")
        report.append("")
        report.append("### Week 2-3: Full Migration")
        report.append("- **Week 2**: Deploy additional solutions as needed")
        report.append("- **Week 3**: Establish collaborative workflows")
        report.append("")
        report.append("### Month 2+: Optimization and Scaling")
        report.append("- **Month 2**: Cost optimization and usage analysis")
        report.append("- **Month 3+**: Advanced workflows and automation")
        report.append("")

        # FAQ section
        report.append("## ❓ Frequently Asked Questions")
        report.append("")
        report.append("### Q: What happens to my data when I'm not using the solution?")
        report.append("**A**: Your data is safely stored in S3 with multiple backups. Compute resources shut down to $0 cost, but your data persists and is instantly available when you restart.")
        report.append("")
        report.append("### Q: Can I collaborate with researchers at other institutions?")
        report.append("**A**: Yes! Each solution includes collaboration features and can be configured for secure multi-institutional access.")
        report.append("")
        report.append("### Q: Do the NIST compliance options make me fully compliant?")
        report.append("**A**: No! The technical security controls are implemented, but you must ensure proper policies, procedures, training, and documentation for full compliance. Consult your institution's compliance office.")
        report.append("")
        report.append("### Q: Can I install additional software?")
        report.append("**A**: Yes! You have full control to install additional tools. Changes can be saved as custom images for future deployments.")
        report.append("")
        report.append("### Q: What if I exceed my budget?")
        report.append("**A**: Built-in budget controls will alert you at 50%, 80%, and 100% of your limit. You can set automatic shutdown at budget limits if desired.")
        report.append("")
        report.append("### Q: How do I get support?")
        report.append("**A**: Each solution includes documentation, tutorials, and community forums. Priority support is available for sponsored research.")
        report.append("")

        return "\n".join(report)

def main():
    parser = argparse.ArgumentParser(description='Generate researcher-ready solutions with domain-specific tools')
    parser.add_argument('--output', default='researcher_ready_solutions.md', help='Output report file')
    parser.add_argument('--cost-calculator', action='store_true', help='Generate cost calculator HTML')

    args = parser.parse_args()

    generator = ResearcherSolutionsGenerator()

    # Generate main report
    report = generator.generate_researcher_solutions_report()
    with open(args.output, 'w', encoding='utf-8') as f:
        f.write(report)

    # Generate cost calculator if requested
    if args.cost_calculator:
        calculator = generator.generate_transparent_pricing_calculator()
        with open('cost_calculator.html', 'w', encoding='utf-8') as f:
            f.write(calculator)
        print("Cost calculator saved to cost_calculator.html")

    print(f"Researcher-ready solutions report saved to {args.output}")

if __name__ == "__main__":
    main()
