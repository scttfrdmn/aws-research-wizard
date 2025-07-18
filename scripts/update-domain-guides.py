#!/usr/bin/env python3
"""
Update all domain guides with:
1. "Using Your Own Data" section (after Step 8)
2. "Extend and Contribute" section (in What's Next?)
"""

import os
import re
import glob

# Domain-specific data examples for the "Using Your Own Data" section
DOMAIN_DATA_EXAMPLES = {
    'agricultural_sciences': {
        'formats': [
            '**Crop yield data** (.csv, .xlsx): Farm management records and harvest data',
            '**Soil samples** (.json, .csv): Chemical composition and nutrient analysis',
            '**Weather station data** (.nc, .csv): Temperature, precipitation, and humidity records',
            '**Satellite imagery** (.tif, .hdf): MODIS, Landsat, and Sentinel agricultural monitoring',
            '**IoT sensor data** (.json, .csv): Real-time field monitoring from connected devices'
        ],
        'example_cmd': 'process_crop_yield.py sample_data.csv',
        'example_replace': 'process_crop_yield.py YOUR_FARM_DATA.csv',
        'tools': 'precision agriculture, crop modeling, farm management systems'
    },
    'astronomy_astrophysics': {
        'formats': [
            '**FITS files** (.fits, .fit): Astronomical images and spectra',
            '**HDF5 data** (.h5, .hdf5): Large telescope survey datasets',
            '**ASCII tables** (.dat, .txt): Photometry and astrometry catalogs',
            '**VOTable format** (.xml, .vot): Virtual Observatory data exchange',
            '**Time series data** (.csv, .json): Variable star and exoplanet observations'
        ],
        'example_cmd': 'ds9 galaxy_image.fits',
        'example_replace': 'ds9 YOUR_OBSERVATION.fits',
        'tools': 'telescopes, spectrographs, photometers'
    },
    'atmospheric_chemistry': {
        'formats': [
            '**NetCDF files** (.nc, .nc4): Atmospheric model output and satellite data',
            '**Chemical data** (.csv, .dat): Species concentrations and reaction rates',
            '**Instrument data** (.hdf, .he5): Satellite and ground-based measurements',
            '**Model output** (.grb, .grib2): Weather and chemistry model predictions',
            '**Time series** (.csv, .json): Long-term atmospheric monitoring data'
        ],
        'example_cmd': 'process_ozone.py ozone_data.nc',
        'example_replace': 'process_ozone.py YOUR_ATMOSPHERIC_DATA.nc',
        'tools': 'atmospheric models, chemical transport models'
    },
    'benchmarking_performance': {
        'formats': [
            '**Performance logs** (.log, .txt): System and application performance data',
            '**Metrics data** (.json, .csv): CPU, memory, network, and storage metrics',
            '**Profiling output** (.prof, .perf): Code profiling and optimization data',
            '**Benchmark results** (.xml, .json): Standard benchmark suite outputs',
            '**Trace files** (.trace, .etl): Execution traces and performance events'
        ],
        'example_cmd': 'analyze_performance.py benchmark_results.json',
        'example_replace': 'analyze_performance.py YOUR_BENCHMARK_DATA.json',
        'tools': 'profilers, benchmarking suites, monitoring tools'
    },
    'chemistry_materials': {
        'formats': [
            '**Structure files** (.pdb, .xyz, .cif): Molecular and crystal structures',
            '**Computational output** (.out, .log): Quantum chemistry calculation results',
            '**Spectroscopy data** (.jdx, .csv): NMR, IR, and mass spectrometry results',
            '**Thermodynamic data** (.dat, .json): Energy, enthalpy, and reaction data',
            '**Materials data** (.vasp, .cp2k): Electronic structure calculation inputs/outputs'
        ],
        'example_cmd': 'gaussian molecule.com',
        'example_replace': 'gaussian YOUR_MOLECULE.com',
        'tools': 'quantum chemistry packages, molecular dynamics'
    },
    'climate_modeling': {
        'formats': [
            '**NetCDF files** (.nc, .nc4): `ncdump -h your_file.nc` to examine structure',
            '**GRIB files** (.grb, .grb2): `wgrib2 your_file.grb2 -V` to view metadata',
            '**CSV/ASCII data**: Direct import with pandas or numpy',
            '**Binary formats**: Use appropriate readers (e.g., fortran unformatted)'
        ],
        'example_cmd': 'cdo sellonlatbox,-140,-60,20,70 2m_temperature.nc north_america_temp.nc',
        'example_replace': 'cdo sellonlatbox,-140,-60,20,70 YOUR_DATA_FILE.nc your_analysis.nc',
        'tools': 'climate models, weather prediction'
    },
    'computational_physics': {
        'formats': [
            '**Simulation output** (.dat, .h5, .vtk): Numerical simulation results',
            '**Parameter files** (.json, .xml, .cfg): Model configuration and input parameters',
            '**Time series data** (.csv, .txt): Physical measurements and calculations',
            '**Grid data** (.mesh, .msh): Finite element and finite difference grids',
            '**Binary output** (.bin, .raw): High-performance simulation data files'
        ],
        'example_cmd': 'python3 analyze_simulation.py results.h5',
        'example_replace': 'python3 analyze_simulation.py YOUR_SIMULATION.h5',
        'tools': 'simulation codes, numerical solvers'
    },
    'cybersecurity_research': {
        'formats': [
            '**Network captures** (.pcap, .pcapng): Network traffic and packet analysis',
            '**Log files** (.log, .json): Security events, firewall, and system logs',
            '**Malware samples** (.exe, .dll): Binary analysis and reverse engineering',
            '**Vulnerability data** (.json, .xml): CVE databases and security assessments',
            '**Threat intelligence** (.csv, .json): IOCs, attack patterns, and signatures'
        ],
        'example_cmd': 'wireshark network_capture.pcap',
        'example_replace': 'wireshark YOUR_NETWORK_DATA.pcap',
        'tools': 'security analysis tools, penetration testing'
    },
    'digital_humanities': {
        'formats': [
            '**Text corpus** (.txt, .xml, .json): Historical documents and literary texts',
            '**Metadata** (.csv, .json): Bibliographic and archival information',
            '**Images** (.jpg, .tif, .pdf): Digitized manuscripts and historical documents',
            '**Database exports** (.sql, .csv): Digital collections and repositories',
            '**Linguistic data** (.conllu, .xml): Annotated texts and linguistic corpora'
        ],
        'example_cmd': 'python3 text_analysis.py corpus_sample.txt',
        'example_replace': 'python3 text_analysis.py YOUR_TEXT_CORPUS.txt',
        'tools': 'text analysis, digital archives'
    },
    'drug_discovery': {
        'formats': [
            '**Molecular structures** (.sdf, .mol2, .pdb): Chemical compounds and proteins',
            '**Assay data** (.csv, .xlsx): Biological activity and screening results',
            '**Pharmacological data** (.json, .xml): ADMET properties and drug interactions',
            '**Protein sequences** (.fasta, .pdb): Target proteins and binding sites',
            '**Chemical databases** (.sdf, .smiles): Compound libraries and virtual screens'
        ],
        'example_cmd': 'rdkit_analysis.py compounds.sdf',
        'example_replace': 'rdkit_analysis.py YOUR_COMPOUNDS.sdf',
        'tools': 'cheminformatics, molecular modeling'
    },
    'finops_economics': {
        'formats': [
            '**Financial data** (.csv, .xlsx): Market data, pricing, and economic indicators',
            '**Cost reports** (.json, .csv): Cloud billing and resource usage data',
            '**Time series** (.csv, .json): Economic forecasts and financial modeling',
            '**Transaction data** (.csv, .parquet): Trading records and financial flows',
            '**Optimization data** (.json, .lp): Linear programming and operations research'
        ],
        'example_cmd': 'python3 cost_analysis.py billing_data.csv',
        'example_replace': 'python3 cost_analysis.py YOUR_FINANCIAL_DATA.csv',
        'tools': 'financial modeling, cost optimization'
    },
    'food_science_nutrition': {
        'formats': [
            '**Nutritional data** (.csv, .xlsx): Food composition and dietary analysis',
            '**Sensory data** (.csv, .json): Consumer testing and food quality metrics',
            '**Microbiological data** (.csv, .txt): Food safety and microbial analysis',
            '**Processing data** (.json, .csv): Food manufacturing and quality control',
            '**Spectroscopy data** (.jdx, .csv): Food authentication and composition analysis'
        ],
        'example_cmd': 'python3 nutrition_analysis.py food_data.csv',
        'example_replace': 'python3 nutrition_analysis.py YOUR_FOOD_DATA.csv',
        'tools': 'food analysis, nutritional modeling'
    },
    'forestry_natural_resources': {
        'formats': [
            '**GIS data** (.shp, .kml, .geojson): Forest boundaries and land use maps',
            '**LiDAR data** (.las, .laz): Forest structure and canopy measurements',
            '**Inventory data** (.csv, .xlsx): Tree measurements and forest surveys',
            '**Satellite imagery** (.tif, .hdf): Remote sensing of forest cover',
            '**Environmental data** (.nc, .csv): Climate and ecological measurements'
        ],
        'example_cmd': 'qgis forest_inventory.shp',
        'example_replace': 'qgis YOUR_FOREST_DATA.shp',
        'tools': 'GIS software, forest inventory systems'
    },
    'genomics': {
        'formats': [
            '**FASTQ files** (.fastq, .fq, .fastq.gz): Raw sequencing reads',
            '**BAM/SAM files** (.bam, .sam): Aligned sequence data',
            '**VCF files** (.vcf, .vcf.gz): Variant call format for genetic variants',
            '**BED files** (.bed): Genomic intervals and annotations',
            '**FASTA files** (.fasta, .fa): Reference genomes and sequences'
        ],
        'example_cmd': 'bwa mem chr20.fasta SRR062634_1.filt.fastq.gz SRR062634_2.filt.fastq.gz > aligned.sam',
        'example_replace': 'bwa mem your_reference.fasta your_sample_R1.fastq.gz your_sample_R2.fastq.gz > your_aligned.sam',
        'tools': 'sequencing platforms, variant callers'
    },
    'geoscience': {
        'formats': [
            '**Seismic data** (.segy, .su): Earthquake and exploration seismology',
            '**GIS data** (.shp, .kml): Geological maps and spatial analysis',
            '**Well logs** (.las, .csv): Subsurface drilling and logging data',
            '**Remote sensing** (.tif, .hdf): Satellite geological observations',
            '**Geochemical data** (.csv, .xlsx): Mineral analysis and geochemistry'
        ],
        'example_cmd': 'python3 earthquake_analysis.py seismic_data.csv',
        'example_replace': 'python3 earthquake_analysis.py YOUR_SEISMIC_DATA.csv',
        'tools': 'geological modeling, seismic processing'
    },
    'geospatial_research': {
        'formats': [
            '**Vector data** (.shp, .kml, .geojson): Points, lines, and polygons with coordinates',
            '**Raster data** (.tif, .img, .nc): Satellite imagery and gridded datasets',
            '**GPS data** (.gpx, .csv): Location tracking and field measurements',
            '**Coordinate data** (.csv, .txt): Latitude/longitude and projected coordinates',
            '**Spatial databases** (.gdb, .sqlite): Complex spatial data collections'
        ],
        'example_cmd': 'qgis study_area.shp',
        'example_replace': 'qgis YOUR_SPATIAL_DATA.shp',
        'tools': 'GIS software, GPS devices'
    },
    'machine_learning': {
        'formats': [
            '**Tabular data** (.csv, .xlsx, .parquet): Structured datasets with features and labels',
            '**Images** (.jpg, .png, .tif): Computer vision and image classification datasets',
            '**Text data** (.txt, .json, .csv): Natural language processing and text mining',
            '**Time series** (.csv, .json): Sequential data for forecasting and analysis',
            '**Model files** (.pkl, .h5, .onnx): Pre-trained models and weights'
        ],
        'example_cmd': 'python3 train_model.py training_data.csv',
        'example_replace': 'python3 train_model.py YOUR_DATASET.csv',
        'tools': 'ML frameworks, data preprocessing'
    },
    'marine_biology_oceanography': {
        'formats': [
            '**CTD data** (.csv, .nc): Temperature, salinity, and depth profiles',
            '**Acoustic data** (.wav, .raw): Marine animal sounds and echolocation',
            '**Satellite data** (.nc, .hdf): Sea surface temperature and ocean color',
            '**Species data** (.csv, .json): Biodiversity surveys and specimen records',
            '**Current data** (.nc, .csv): Ocean circulation and flow measurements'
        ],
        'example_cmd': 'python3 ocean_analysis.py ctd_data.csv',
        'example_replace': 'python3 ocean_analysis.py YOUR_OCEAN_DATA.csv',
        'tools': 'oceanographic instruments, marine surveys'
    },
    'materials_science': {
        'formats': [
            '**Crystal structures** (.cif, .pdb): Atomic arrangements and lattice parameters',
            '**Spectroscopy data** (.csv, .jdx): X-ray, NMR, and other characterization',
            '**Microscopy images** (.tif, .dm3): SEM, TEM, and optical microscopy',
            '**Mechanical data** (.csv, .txt): Stress-strain curves and material properties',
            '**Computational data** (.vasp, .lammps): Simulation inputs and outputs'
        ],
        'example_cmd': 'python3 materials_analysis.py sample_data.cif',
        'example_replace': 'python3 materials_analysis.py YOUR_MATERIAL.cif',
        'tools': 'characterization instruments, simulation codes'
    },
    'mathematical_modeling': {
        'formats': [
            '**Numerical data** (.csv, .dat, .txt): Mathematical solutions and parameters',
            '**Model definitions** (.json, .xml): Equation systems and model specifications',
            '**Simulation output** (.h5, .mat): Results from numerical computations',
            '**Optimization data** (.lp, .json): Linear and nonlinear programming problems',
            '**Statistical data** (.csv, .rdata): Data for statistical modeling and analysis'
        ],
        'example_cmd': 'python3 solve_model.py parameters.json',
        'example_replace': 'python3 solve_model.py YOUR_MODEL_PARAMS.json',
        'tools': 'mathematical software, numerical solvers'
    },
    'neuroscience': {
        'formats': [
            '**Neuroimaging data** (.nii, .dcm): MRI, fMRI, and brain imaging',
            '**Electrophysiology** (.edf, .mat): EEG, MEG, and neural recordings',
            '**Behavioral data** (.csv, .json): Cognitive tests and experimental results',
            '**Spike data** (.nev, .plx): Single-unit and multi-unit neural activity',
            '**Anatomical data** (.swc, .obj): Neural morphology and connectivity'
        ],
        'example_cmd': 'python3 brain_analysis.py fmri_data.nii',
        'example_replace': 'python3 brain_analysis.py YOUR_BRAIN_DATA.nii',
        'tools': 'neuroimaging software, data acquisition systems'
    },
    'quantum_computing': {
        'formats': [
            '**Quantum circuits** (.qasm, .json): Circuit definitions and gate sequences',
            '**Experimental data** (.csv, .json): Quantum device measurements and calibration',
            '**Simulation output** (.h5, .json): Quantum state vectors and probability distributions',
            '**Algorithm parameters** (.json, .yaml): Quantum algorithm configurations',
            '**Device specs** (.json, .xml): Quantum hardware characteristics and noise models'
        ],
        'example_cmd': 'python3 quantum_circuit.py experiment_data.json',
        'example_replace': 'python3 quantum_circuit.py YOUR_QUANTUM_DATA.json',
        'tools': 'quantum simulators, quantum devices'
    },
    'renewable_energy_systems': {
        'formats': [
            '**Energy data** (.csv, .json): Solar, wind, and energy production measurements',
            '**Weather data** (.nc, .csv): Meteorological data for renewable energy forecasting',
            '**Grid data** (.csv, .xml): Electrical grid monitoring and smart meter data',
            '**System parameters** (.json, .cfg): Renewable energy system configurations',
            '**Economic data** (.csv, .xlsx): Energy pricing and financial modeling'
        ],
        'example_cmd': 'python3 energy_analysis.py solar_data.csv',
        'example_replace': 'python3 energy_analysis.py YOUR_ENERGY_DATA.csv',
        'tools': 'energy monitoring systems, grid simulators'
    },
    'social_sciences': {
        'formats': [
            '**Survey data** (.csv, .xlsx, .sav): Questionnaire responses and social research',
            '**Demographic data** (.csv, .json): Population statistics and census information',
            '**Network data** (.gml, .json): Social networks and relationship mapping',
            '**Text data** (.txt, .json): Interview transcripts and qualitative research',
            '**Statistical data** (.csv, .rdata): Experimental and observational study results'
        ],
        'example_cmd': 'python3 social_analysis.py survey_data.csv',
        'example_replace': 'python3 social_analysis.py YOUR_SURVEY_DATA.csv',
        'tools': 'statistical software, survey platforms'
    },
    'sports_science_biomechanics': {
        'formats': [
            '**Motion capture** (.c3d, .csv): 3D movement and biomechanical analysis',
            '**Force data** (.csv, .txt): Ground reaction forces and kinetic measurements',
            '**EMG data** (.csv, .edf): Muscle activation and electromyography',
            '**Video analysis** (.mp4, .avi): Performance analysis and technique assessment',
            '**Physiological data** (.csv, .json): Heart rate, VO2, and metabolic measurements'
        ],
        'example_cmd': 'python3 biomech_analysis.py motion_data.c3d',
        'example_replace': 'python3 biomech_analysis.py YOUR_MOTION_DATA.c3d',
        'tools': 'motion capture systems, force plates'
    },
    'structural_biology': {
        'formats': [
            '**Protein structures** (.pdb, .cif): X-ray crystallography and NMR structures',
            '**Cryo-EM data** (.mrc, .map): Electron microscopy density maps',
            '**Sequence data** (.fasta, .aln): Protein and nucleic acid sequences',
            '**Experimental data** (.csv, .json): Biophysical and biochemical measurements',
            '**Molecular dynamics** (.dcd, .xtc): Simulation trajectories and conformations'
        ],
        'example_cmd': 'pymol protein_structure.pdb',
        'example_replace': 'pymol YOUR_PROTEIN.pdb',
        'tools': 'crystallography software, molecular viewers'
    },
    'visualization_studio': {
        'formats': [
            '**3D models** (.obj, .stl, .ply): Three-dimensional objects and meshes',
            '**Visualization data** (.json, .csv): Data for charts, graphs, and interactive plots',
            '**Image data** (.jpg, .png, .tif): Static images and visual content',
            '**Animation data** (.fbx, .dae): 3D animations and motion graphics',
            '**VR/AR content** (.unity, .blend): Virtual and augmented reality assets'
        ],
        'example_cmd': 'blender 3d_model.obj',
        'example_replace': 'blender YOUR_3D_MODEL.obj',
        'tools': '3D software, visualization platforms'
    }
}

# Domain-specific software suggestions for the "Extend and Contribute" section
DOMAIN_SOFTWARE_SUGGESTIONS = {
    'agricultural_sciences': ['DSSAT', 'APSIM', 'CropSyst', 'AgroClimate', 'FarmBeats'],
    'astronomy_astrophysics': ['IRAF', 'SAOImage DS9', 'CASA', 'AIPS', 'Montage'],
    'atmospheric_chemistry': ['GEOS-Chem', 'WRF-Chem', 'CAMx', 'CMAQ', 'TM5'],
    'benchmarking_performance': ['SPEC CPU', 'Linpack', 'STREAM', 'IOzone', 'NetPerf'],
    'chemistry_materials': ['VASP', 'Quantum ESPRESSO', 'CP2K', 'LAMMPS', 'Materials Studio'],
    'climate_modeling': ['RegCM', 'MM5', 'ICON', 'FMS'],
    'computational_physics': ['OpenFOAM', 'FEniCS', 'COMSOL', 'ANSYS', 'ParaView'],
    'cybersecurity_research': ['Metasploit', 'Nmap', 'Burp Suite', 'YARA', 'Volatility'],
    'digital_humanities': ['TEI', 'Omeka', 'Gephi', 'Voyant Tools', 'ELAN'],
    'drug_discovery': ['Schrödinger Suite', 'MOE', 'OpenEye', 'ChemAxon', 'Pipeline Pilot'],
    'finops_economics': ['FinOps Toolkit', 'CloudHealth', 'Terraform', 'Kubernetes', 'Prometheus'],
    'food_science_nutrition': ['NutriData', 'Sensory Analysis Software', 'FoodCAD', 'ChemSketch'],
    'forestry_natural_resources': ['FUSION', 'LAStools', 'Forest Vegetation Simulator', 'i-Tree'],
    'genomics': ['STAR', 'Cufflinks', 'Trinity', 'MetaPhlAn', 'SPAdes'],
    'geoscience': ['Seismic Unix', 'GMT', 'PETREL', 'ArcGIS', 'GRASS GIS'],
    'geospatial_research': ['PostGIS', 'GDAL/OGR', 'Leaflet', 'OpenLayers', 'GeoServer'],
    'machine_learning': ['XGBoost', 'LightGBM', 'Optuna', 'MLflow', 'Weights & Biases'],
    'marine_biology_oceanography': ['Ocean Data View', 'MATLAB Oceanography', 'Ferret', 'ERDDAP'],
    'materials_science': ['VESTA', 'CrystalMaker', 'Materials Project', 'ASE', 'Pymatgen'],
    'mathematical_modeling': ['COMSOL', 'MATLAB', 'Mathematica', 'SageMath', 'FEniCS'],
    'neuroscience': ['FSL', 'FreeSurfer', 'SPM', 'AFNI', 'Brainstorm'],
    'quantum_computing': ['Quantum Inspire', 'Forest', 'ProjectQ', 'Strawberry Fields'],
    'renewable_energy_systems': ['HOMER', 'PVsyst', 'WindPRO', 'SAM', 'TRNSYS'],
    'social_sciences': ['SPSS', 'Stata', 'NVivo', 'Atlas.ti', 'Gephi'],
    'sports_science_biomechanics': ['Visual3D', 'OpenSim', 'Kinovea', 'SIMI Motion', 'Contemplas'],
    'structural_biology': ['Coot', 'Phenix', 'CCP4', 'Relion', 'ChimeraX'],
    'visualization_studio': ['Blender', 'Unity', 'Unreal Engine', 'Three.js', 'D3.js']
}

# Additional domain pack suggestions
DOMAIN_PACK_SUGGESTIONS = {
    'agricultural_sciences': ['precision agriculture', 'soil science', 'agricultural economics', 'crop breeding'],
    'astronomy_astrophysics': ['exoplanet research', 'cosmology', 'stellar physics', 'galactic astronomy'],
    'atmospheric_chemistry': ['air quality modeling', 'atmospheric physics', 'stratospheric chemistry'],
    'benchmarking_performance': ['cloud performance', 'network benchmarking', 'GPU computing', 'storage optimization'],
    'chemistry_materials': ['computational chemistry', 'polymer science', 'catalysis research', 'nanomaterials'],
    'climate_modeling': ['atmospheric chemistry', 'hydrology', 'oceanography'],
    'computational_physics': ['plasma physics', 'condensed matter', 'particle physics', 'nuclear physics'],
    'cybersecurity_research': ['malware analysis', 'network security', 'digital forensics', 'threat intelligence'],
    'digital_humanities': ['computational linguistics', 'digital archives', 'cultural analytics', 'text mining'],
    'drug_discovery': ['pharmacokinetics', 'toxicology', 'medicinal chemistry', 'clinical data analysis'],
    'finops_economics': ['financial modeling', 'risk analysis', 'algorithmic trading', 'econometrics'],
    'food_science_nutrition': ['food safety', 'sensory analysis', 'nutritional epidemiology', 'food engineering'],
    'forestry_natural_resources': ['wildlife ecology', 'conservation biology', 'natural resource economics'],
    'genomics': ['single-cell genomics', 'epigenomics', 'proteomics', 'metabolomics'],
    'geoscience': ['hydrogeology', 'petroleum geology', 'environmental geology', 'geophysics'],
    'geospatial_research': ['remote sensing', 'cartography', 'spatial analysis', 'location intelligence'],
    'machine_learning': ['deep learning', 'reinforcement learning', 'computer vision', 'natural language processing'],
    'marine_biology_oceanography': ['fisheries science', 'marine ecology', 'coastal engineering', 'ocean modeling'],
    'materials_science': ['biomaterials', 'electronic materials', 'energy materials', 'manufacturing'],
    'mathematical_modeling': ['operations research', 'optimization', 'statistical modeling', 'numerical analysis'],
    'neuroscience': ['computational neuroscience', 'neuroimaging', 'brain-computer interfaces', 'cognitive science'],
    'quantum_computing': ['quantum algorithms', 'quantum cryptography', 'quantum sensing', 'quantum simulation'],
    'renewable_energy_systems': ['energy storage', 'smart grids', 'energy policy', 'carbon capture'],
    'social_sciences': ['computational social science', 'survey research', 'network analysis', 'behavioral economics'],
    'sports_science_biomechanics': ['exercise physiology', 'sports psychology', 'performance analysis', 'injury prevention'],
    'structural_biology': ['protein folding', 'drug design', 'membrane proteins', 'enzyme mechanisms'],
    'visualization_studio': ['scientific visualization', 'data visualization', 'virtual reality', 'augmented reality']
}

def extract_domain_name(file_path):
    """Extract domain name from file path"""
    return os.path.basename(os.path.dirname(file_path))

def create_custom_data_section(domain_name):
    """Create the 'Using Your Own Data' section for a specific domain"""
    domain_info = DOMAIN_DATA_EXAMPLES.get(domain_name, {
        'formats': [
            '**Data files** (.csv, .json, .txt): Your research datasets',
            '**Analysis files** (.dat, .h5): Processed data and results',
            '**Configuration files** (.json, .yaml): Analysis parameters'
        ],
        'example_cmd': 'python3 analysis.py sample_data.csv',
        'example_replace': 'python3 analysis.py YOUR_DATA.csv',
        'tools': 'research software and instruments'
    })

    title = domain_name.replace('_', ' ').title()

    section = f"""## Step 9: Using Your Own {title} Data

Instead of the tutorial data, you can analyze your own {domain_name.replace('_', ' ')} datasets:

### Upload Your Data
```bash
# Option 1: Upload from your local computer
scp -i ~/.ssh/id_rsa your_data_file.* ec2-user@12.34.56.78:~/{domain_name}-tutorial/

# Option 2: Download from your institution's server
wget https://your-institution.edu/data/research_data.csv

# Option 3: Access your AWS S3 bucket
aws s3 cp s3://your-research-bucket/{domain_name}-data/ . --recursive
```

### Common Data Formats Supported
"""

    for format_item in domain_info['formats']:
        section += f"- {format_item}\n"

    section += f"""
### Replace Tutorial Commands
Simply substitute your filenames in any tutorial command:
```bash
# Instead of tutorial data:
{domain_info['example_cmd']}

# Use your data:
{domain_info['example_replace']}
```

### Data Size Considerations
- **Small datasets** (<10 GB): Process directly on the instance
- **Large datasets** (10-100 GB): Use S3 for storage, process in chunks
- **Very large datasets** (>100 GB): Consider multi-node setup or data preprocessing

"""
    return section

def create_extend_section(domain_name):
    """Create the 'Extend and Contribute' section for a specific domain"""
    software_suggestions = DOMAIN_SOFTWARE_SUGGESTIONS.get(domain_name, ['Additional software packages'])
    domain_pack_suggestions = DOMAIN_PACK_SUGGESTIONS.get(domain_name, ['related research domains'])

    # Format software suggestions
    software_list = ', '.join(software_suggestions[:5])  # Limit to 5 examples

    # Format domain pack suggestions
    domain_list = ', '.join(domain_pack_suggestions[:4])  # Limit to 4 examples

    section = f"""### Extend and Contribute
**🚀 Help us expand AWS Research Wizard!**

**Missing a tool or domain?** We welcome suggestions for:
- **New {domain_name.replace('_', ' ')} software** (e.g., {software_list})
- **Additional domain packs** (e.g., {domain_list})
- **New data sources** or tutorials for specific research workflows

**How to contribute:**
- [Request new features](https://github.com/aws-research-wizard/aws-research-wizard/issues/new?template=feature_request.md)
- [Suggest domain packs](https://github.com/aws-research-wizard/aws-research-wizard/discussions/categories/domain-suggestions)
- [Share your configurations](https://forum.researchwizard.app/share-configs)
- [Join development discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)

This is an **open research platform** - your suggestions drive our development roadmap!

"""
    return section

def update_domain_guide(file_path):
    """Update a single domain guide with both sections"""
    domain_name = extract_domain_name(file_path)

    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # Skip files that already have the sections (to avoid duplicates)
    if "Using Your Own" in content and "Extend and Contribute" in content:
        print(f"⏭️  Skipping {domain_name} - already updated")
        return False

    # Try multiple patterns to find where to insert the custom data section
    insertion_patterns = [
        r'## Step (\d+): Monitor Your Costs',
        r'## Understanding Your Costs',
        r'## What\'s Next\?',
        r'## Troubleshooting',
        r'### Getting Help'
    ]

    custom_data_inserted = False
    for pattern in insertion_patterns:
        match = re.search(pattern, content)
        if match:
            custom_data_section = create_custom_data_section(domain_name)
            insertion_point = match.start()
            content = content[:insertion_point] + custom_data_section + content[insertion_point:]
            custom_data_inserted = True
            print(f"✅ Inserted custom data section before '{match.group()}' in {domain_name}")
            break

    if not custom_data_inserted:
        # If no suitable section found, add before the end of file
        custom_data_section = create_custom_data_section(domain_name)
        # Find the last line with content (not just whitespace)
        lines = content.split('\n')
        for i in range(len(lines) - 1, -1, -1):
            if lines[i].strip():
                break
        content = '\n'.join(lines[:i+1]) + '\n\n' + custom_data_section + '\n' + '\n'.join(lines[i+1:])
        print(f"✅ Added custom data section to end of {domain_name}")

    # Try to find where to insert the extend section
    extend_insertion_patterns = [
        r'### Join the .+ Community\n(.*?)(?=\n##|\n### |$)',
        r'### Getting Help',
        r'### Emergency Stop',
        r'## Troubleshooting'
    ]

    extend_inserted = False
    for pattern in extend_insertion_patterns:
        match = re.search(pattern, content, re.DOTALL)
        if match:
            extend_section = create_extend_section(domain_name)
            if '### Join the' in match.group():
                # Insert after community section
                insertion_point = match.end()
                content = content[:insertion_point] + "\n" + extend_section + content[insertion_point:]
            else:
                # Insert before other sections
                insertion_point = match.start()
                content = content[:insertion_point] + extend_section + "\n" + content[insertion_point:]
            extend_inserted = True
            print(f"✅ Inserted extend section near '{pattern}' in {domain_name}")
            break

    if not extend_inserted:
        # Add before troubleshooting or at the end
        extend_section = create_extend_section(domain_name)
        lines = content.split('\n')
        for i in range(len(lines) - 1, -1, -1):
            if lines[i].strip():
                break
        content = '\n'.join(lines[:i+1]) + '\n\n' + extend_section + '\n' + '\n'.join(lines[i+1:])
        print(f"✅ Added extend section to end of {domain_name}")

    # Update step numbers if needed (Monitor Your Costs becomes Step 10, Clean Up becomes Step 11)
    content = re.sub(r'## Step (\d+): Monitor Your Costs', r'## Step 10: Monitor Your Costs', content)
    content = re.sub(r'## Step (\d+): Clean Up', r'## Step 11: Clean Up', content)

    # Write the updated content
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)

    print(f"✅ Updated {domain_name}")
    return True

def main():
    """Update all domain guides"""
    guide_pattern = "/Users/scttfrdmn/src/aws-research-wizard/docs/domain-guides/*/getting-started.md"
    guide_files = glob.glob(guide_pattern)

    print(f"Found {len(guide_files)} domain guides to update")

    updated_count = 0
    for guide_file in sorted(guide_files):
        if update_domain_guide(guide_file):
            updated_count += 1

    print(f"\n🎉 Successfully updated {updated_count} domain guides!")
    print("Added sections:")
    print("  1. 'Using Your Own Data' (Step 9)")
    print("  2. 'Extend and Contribute' (in What's Next?)")

if __name__ == "__main__":
    main()
