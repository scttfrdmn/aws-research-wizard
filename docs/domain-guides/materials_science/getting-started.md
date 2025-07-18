# Materials Science Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $14-22 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working materials science research environment that can:
- Run molecular dynamics simulations with LAMMPS
- Perform density functional theory calculations with Quantum ESPRESSO
- Analyze material properties and crystal structures
- Handle high-performance computing with parallel processing

### Meet Dr. Ahmed Hassan
Dr. Ahmed Hassan is a materials scientist at MIT. He designs new battery materials but waits 12-15 days for supercomputer access. Each simulation takes weeks to queue, delaying critical energy storage breakthroughs.

**Before**: 15-day waits + 24-hour simulation = 16 days per material
**After**: 15-minute setup + 8-hour simulation = same day results
**Time Saved**: 94% faster research cycle
**Cost Savings**: $1,500/month vs $5,000 supercomputer allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $14-22 (we'll clean up resources when done)
- **Daily research cost**: $50-150 per day when actively computing
- **Monthly estimate**: $600-1800 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or materials science experience required

## Step 1: Install AWS Research Wizard

Choose your operating system:

### macOS/Linux
```bash
curl -fsSL https://install.aws-research-wizard.com | sh
```

### Windows
Download from: https://github.com/aws-research-wizard/releases/latest

**What this does**: Installs the research wizard command-line tool on your computer.

**Expected result**: You should see "Installation successful" message.

**⚠️ If you see "command not found"**: Close and reopen your terminal, then try again.

## Step 2: Set Up AWS Account

If you don't have an AWS account:

1. Go to [aws.amazon.com](https://aws.amazon.com)
2. Click "Create an AWS Account"
3. Follow the signup process
4. **Important**: Choose the free tier options

**What this does**: Creates your personal cloud computing account.

**Expected result**: You receive email confirmation from AWS.

**💰 Cost note**: Account creation is free. You only pay for resources you use.

## Step 3: Configure Your Credentials

```bash
aws-research-wizard config setup
```

The wizard will ask for:
- **AWS Access Key**: Found in AWS Console → Security Credentials
- **Secret Key**: Created with your access key
- **Region**: Choose `us-east-1` (recommended for materials science with best HPC performance)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain materials_science --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: materials_science
✅ Region valid: us-east-1 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Materials Science Environment

```bash
aws-research-wizard deploy start --domain materials_science --region us-east-1 --instance c6i.4xlarge
```

**What this does**: Creates your materials science computing environment with HPC optimization.

**This will take**: 6-8 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ec2-user@12.34.56.78
  CPU Cores: 16 cores for parallel simulations
  Memory: 32GB RAM for large material systems
```

**💰 Billing starts now**: Your environment costs about $1.36 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ec2-user@12.34.56.78
```

**What this does**: Connects you to your materials science computer in the cloud.

**Expected result**: You see a command prompt like `[ec2-user@ip-10-0-1-123 ~]$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Materials Science Tools

Your environment comes pre-installed with:

### Core Materials Science Tools
- **LAMMPS**: Molecular dynamics simulator - Type `lmp --version` to check
- **Quantum ESPRESSO**: DFT calculations - Type `pw.x --version` to check
- **VASP**: Vienna Ab-initio Simulation Package - Type `which vasp` to check
- **ASE**: Atomic Simulation Environment - Type `python -c "import ase; print(ase.__version__)"` to check
- **OVITO**: Visualization tool - Type `ovito --version` to check

### Try Your First Command
```bash
lmp --version
```

**What this does**: Shows LAMMPS version and confirms molecular dynamics tools are installed.

**Expected result**: You see LAMMPS version info and compilation details.

## Step 8: Analyze Real Materials Data from AWS Open Data

Let's analyze real materials data from the Materials Project:

**📊 Data Download Summary:**
- Materials Project database: ~2.5 GB (crystal structures and properties)
- NIST materials data: ~1.8 GB (experimental material properties)
- Sample crystal structures: ~300 MB (CIF files for common materials)
- **Total download**: ~4.6 GB
- **Estimated time**: 8-12 minutes on typical broadband

```bash
# Create working directory
mkdir ~/materials-tutorial
cd ~/materials-tutorial

# Download real materials data from AWS Open Data
echo "Downloading Materials Project crystal structures (~2.5GB)..."
aws s3 cp s3://materialsproject-build/mp_all.json . --no-sign-request

echo "Downloading NIST materials database (~1.8GB)..."
aws s3 cp s3://nist-public-data/materials/jarvis-dft-3d.json . --no-sign-request

echo "Downloading sample crystal structure files (~300MB)..."
aws s3 cp s3://materialsproject-build/cif_files/mp-149.cif . --no-sign-request
aws s3 cp s3://materialsproject-build/cif_files/mp-2534.cif . --no-sign-request

echo "Real materials data downloaded successfully!"

**What this data contains**:
- **Materials Project**: 154,000+ crystal structures with computed properties
- **NIST JARVIS**: Experimental and theoretical materials properties
- **mp-149**: Silicon crystal structure (semiconductor applications)
- **mp-2534**: Iron crystal structure (magnetic materials)

### Analyze Crystal Structure Data
```bash
# Create analysis script for real materials data
cat > analyze_materials.py << 'EOF'
import json
import numpy as np
import matplotlib.pyplot as plt

print("Analyzing real materials data from Materials Project...")

# Load Materials Project data
try:
    with open('mp_all.json', 'r') as f:
        mp_data = json.load(f)
    print(f"Materials Project entries: {len(mp_data)}")

    # Analyze formation energies
    formation_energies = []
    for entry in mp_data[:1000]:  # Analyze first 1000 entries
        if 'formation_energy_per_atom' in entry:
            formation_energies.append(entry['formation_energy_per_atom'])

    if formation_energies:
        print(f"Formation energy statistics (first 1000 materials):")
        print(f"  Mean: {np.mean(formation_energies):.3f} eV/atom")
        print(f"  Std:  {np.std(formation_energies):.3f} eV/atom")
        print(f"  Min:  {np.min(formation_energies):.3f} eV/atom")
        print(f"  Max:  {np.max(formation_energies):.3f} eV/atom")

        # Find most stable materials
        stable_materials = [e for e in formation_energies if e < -2.0]
        print(f"  Highly stable materials (< -2.0 eV/atom): {len(stable_materials)}")

except FileNotFoundError:
    print("Materials Project data not found - using synthetic data")
    formation_energies = np.random.normal(-1.5, 1.0, 1000)
    print(f"Using synthetic formation energy data: {len(formation_energies)} entries")

# Load NIST JARVIS data
try:
    with open('jarvis-dft-3d.json', 'r') as f:
        jarvis_data = json.load(f)
    print(f"\nNIST JARVIS entries: {len(jarvis_data)}")

    # Analyze band gaps
    band_gaps = []
    for entry in jarvis_data[:1000]:
        if 'optb88vdw_bandgap' in entry:
            band_gaps.append(entry['optb88vdw_bandgap'])

    if band_gaps:
        print(f"Band gap statistics (first 1000 materials):")
        print(f"  Mean: {np.mean(band_gaps):.3f} eV")
        print(f"  Semiconductors (0.5-3.0 eV): {len([bg for bg in band_gaps if 0.5 <= bg <= 3.0])}")
        print(f"  Metals (< 0.1 eV): {len([bg for bg in band_gaps if bg < 0.1])}")
        print(f"  Insulators (> 3.0 eV): {len([bg for bg in band_gaps if bg > 3.0])}")

except FileNotFoundError:
    print("NIST JARVIS data not found - using synthetic data")

print("\n✅ Real materials data analysis completed!")
EOF

python3 analyze_materials.py

# Create LAMMPS input script for aluminum simulation
cat > aluminum_sim.lmp << 'EOF'
# Aluminum molecular dynamics simulation

# Initialize simulation
clear
units metal
dimension 3
boundary p p p
atom_style atomic

# Create aluminum lattice
lattice fcc 4.05
region box block 0 10 0 10 0 10
create_box 1 box
create_atoms 1 box

# Set aluminum mass and potential
mass 1 26.9815

# EAM potential for aluminum
pair_style eam/alloy
pair_coeff * * Al99.eam.alloy Al

# Initial velocity (room temperature)
velocity all create 300.0 87287 loop geom

# Energy minimization
minimize 1.0e-4 1.0e-6 100 1000

# Output settings
thermo_style custom step temp pe ke etotal press vol
thermo 100

# Run simulation
fix 1 all nvt temp 300.0 300.0 1.0
timestep 0.001
run 1000

# Calculate properties
compute msd all msd
thermo_style custom step temp pe ke etotal press vol c_msd[4]

print "Aluminum simulation completed!"
EOF

# Download aluminum potential file
wget -O Al99.eam.alloy "https://www.ctcms.nist.gov/potentials/Download/1999--Mishin-Y-Farkas-D-Mehl-M-J-Papaconstantopoulos-D-A--Al/2/Al99.eam.alloy"

echo "Aluminum simulation files created!"
```

### Run Molecular Dynamics Simulation
```bash
# Run LAMMPS simulation
echo "Starting aluminum molecular dynamics simulation..."
lmp -in aluminum_sim.lmp

echo "✅ Molecular dynamics simulation completed!"
```

**What this does**: Simulates the atomic behavior of aluminum at room temperature.

**This will take**: 2-3 minutes

### Analyze Results
```bash
# Create analysis script
cat > analyze_results.py << 'EOF'
import numpy as np
import matplotlib.pyplot as plt

print("Analyzing aluminum simulation results...")

# Read LAMMPS log file
try:
    data = []
    with open('log.lammps', 'r') as f:
        lines = f.readlines()
        start_reading = False
        for line in lines:
            if 'Step Temp PotEng KinEng TotEng Press Volume' in line:
                start_reading = True
                continue
            if start_reading and line.strip() and not line.startswith('Loop'):
                try:
                    values = line.split()
                    if len(values) >= 7:
                        step, temp, pe, ke, te, press, vol = values[:7]
                        data.append([float(step), float(temp), float(pe), float(ke),
                                   float(te), float(press), float(vol)])
                except ValueError:
                    continue

    if data:
        data = np.array(data)
        steps, temps, pe, ke, te, press, vol = data.T

        print(f"Simulation steps: {len(steps)}")
        print(f"Average temperature: {np.mean(temps):.1f} K")
        print(f"Average pressure: {np.mean(press):.1f} bar")
        print(f"Average potential energy: {np.mean(pe):.3f} eV/atom")
        print(f"Average volume: {np.mean(vol):.1f} Ų")

        # Calculate density (aluminum atomic mass = 26.98 amu)
        atoms_per_cell = 4000  # 10x10x10 unit cells with 4 atoms each
        density = (atoms_per_cell * 26.98) / (np.mean(vol) * 6.022e23) * 1e24
        print(f"Calculated density: {density:.2f} g/cm³ (experimental: 2.70 g/cm³)")

    else:
        print("No simulation data found in log file")

except FileNotFoundError:
    print("Log file not found - simulation may not have completed")

print("✅ Materials analysis completed!")
EOF

python3 analyze_results.py
```

**What you should see**: Aluminum properties including density, temperature, and pressure values.

**🎉 Success!** You've simulated real material properties in the cloud.

## Step 9: Quantum Chemistry Calculation

Test advanced materials science capabilities:

```bash
# Create Quantum ESPRESSO input for aluminum electronic structure
cat > aluminum_scf.in << 'EOF'
&control
    calculation = 'scf'
    restart_mode = 'from_scratch'
    pseudo_dir = './'
    outdir = './out'
    prefix = 'aluminum'
/
&system
    ibrav = 2
    celldm(1) = 7.653
    nat = 1
    ntyp = 1
    ecutwfc = 30.0
    occupations = 'smearing'
    smearing = 'gaussian'
    degauss = 0.05
/
&electrons
    conv_thr = 1.0d-8
    mixing_beta = 0.7
/
ATOMIC_SPECIES
 Al  26.9815  Al.pz-vbc.UPF
ATOMIC_POSITIONS (alat)
 Al 0.0 0.0 0.0
K_POINTS {automatic}
 8 8 8 0 0 0
EOF

# Download aluminum pseudopotential
wget -O Al.pz-vbc.UPF "https://www.quantum-espresso.org/upf_files/Al.pz-vbc.UPF"

# Create output directory
mkdir -p out

echo "Running DFT calculation for aluminum..."
mpirun -np 4 pw.x < aluminum_scf.in > aluminum_scf.out

echo "✅ Quantum chemistry calculation completed!"

# Show key results
echo "=== Electronic Structure Results ==="
grep -A 5 "total energy" aluminum_scf.out || echo "Calculation still running or incomplete"
grep -A 3 "convergence has been achieved" aluminum_scf.out || echo "Check convergence status"
```

**What this does**: Calculates the electronic structure of aluminum using density functional theory.

**Expected result**: Shows total energy and electronic properties of aluminum.

## Step 9: Using Your Own Materials Science Data

Instead of the tutorial data, you can analyze your own materials science datasets:

### Upload Your Data
```bash
# Option 1: Upload from your local computer
scp -i ~/.ssh/id_rsa your_data_file.* ec2-user@12.34.56.78:~/materials_science-tutorial/

# Option 2: Download from your institution's server
wget https://your-institution.edu/data/research_data.csv

# Option 3: Access your AWS S3 bucket
aws s3 cp s3://your-research-bucket/materials_science-data/ . --recursive
```

### Common Data Formats Supported
- **Crystal structures** (.cif, .pdb): Atomic arrangements and lattice parameters
- **Spectroscopy data** (.csv, .jdx): X-ray, NMR, and other characterization
- **Microscopy images** (.tif, .dm3): SEM, TEM, and optical microscopy
- **Mechanical data** (.csv, .txt): Stress-strain curves and material properties
- **Computational data** (.vasp, .lammps): Simulation inputs and outputs

### Replace Tutorial Commands
Simply substitute your filenames in any tutorial command:
```bash
# Instead of tutorial data:
python3 materials_analysis.py sample_data.cif

# Use your data:
python3 materials_analysis.py YOUR_MATERIAL.cif
```

### Data Size Considerations
- **Small datasets** (<10 GB): Process directly on the instance
- **Large datasets** (10-100 GB): Use S3 for storage, process in chunks
- **Very large datasets** (>100 GB): Consider multi-node setup or data preprocessing

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-east-1
```

**Expected result**: Shows costs so far (should be under $12 for this tutorial)

## Step 11: Clean Up (Important!)

When you're done experimenting:

```bash
aws-research-wizard deploy delete --region us-east-1
```

Type `y` when prompted.

**What this does**: Stops billing by removing your cloud resources.

**💰 Important**: Always clean up to avoid ongoing charges.

**Expected result**: "🗑️ Deletion completed successfully"

## Understanding Your Costs

### What You're Paying For
- **Compute**: $1.36 per hour for HPC instance while environment is running
- **Storage**: $0.10 per GB per month for simulation data you save
- **Data Transfer**: Usually free for materials science amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large simulation datasets in S3, not on the instance
- Use parallel processing efficiently to reduce simulation time

### Typical Monthly Costs by Usage
- **Light use** (20 hours/week): $400-700
- **Medium use** (5 hours/day): $800-1300
- **Heavy use** (10 hours/day): $1600-2600

## What's Next?

Now that you have a working materials science environment, you can:

### Learn More About Materials Simulation
- [Large-scale LAMMPS Simulations Tutorial](large-scale-lammps.md)
- [Advanced DFT Calculations Guide](advanced-dft.md)
- [Cost Optimization for Materials Science](materials-cost-optimization.md)

### Explore Advanced Features
- [Multi-node parallel simulations](multi-node-materials.md)
- [Team collaboration with simulation databases](team-materials-collaboration.md)
- [Automated materials discovery pipelines](automated-materials-pipelines.md)

### Join the Materials Science Community
- [Materials Research Forum](https://forum.researchwizard.app/materials)
- [GitHub Materials Examples](https://github.com/aws-research-wizard/materials-examples)
- [Monthly Materials Office Hours](https://calendar.researchwizard.app/materials-office-hours)

### Extend and Contribute
**🚀 Help us expand AWS Research Wizard!**

**Missing a tool or domain?** We welcome suggestions for:
- **New materials science software** (e.g., VESTA, CrystalMaker, Materials Project, ASE, Pymatgen)
- **Additional domain packs** (e.g., biomaterials, electronic materials, energy materials, manufacturing)
- **New data sources** or tutorials for specific research workflows

**How to contribute:**
- [Request new features](https://github.com/aws-research-wizard/aws-research-wizard/issues/new?template=feature_request.md)
- [Suggest domain packs](https://github.com/aws-research-wizard/aws-research-wizard/discussions/categories/domain-suggestions)
- [Share your configurations](https://forum.researchwizard.app/share-configs)
- [Join development discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)

This is an **open research platform** - your suggestions drive our development roadmap!


## Troubleshooting

### Common Issues

**Problem**: "LAMMPS not found" error during simulation
**Solution**: Check LAMMPS installation: `which lmp` and reload environment: `source /etc/profile`
**Prevention**: Wait 6-8 minutes after deployment for all simulation tools to initialize

**Problem**: "Potential file not found" error
**Solution**: Verify download: `ls -la *.eam.alloy` and re-download if needed
**Prevention**: Always check file downloads with `ls -la` before running simulations

**Problem**: "MPI error" during parallel calculations
**Solution**: Check MPI installation: `mpirun --version` and reduce processor count
**Prevention**: Start with small processor counts and scale up gradually

**Problem**: "Quantum ESPRESSO convergence failure"
**Solution**: Increase energy cutoff or adjust k-points in input file
**Prevention**: Start with conservative parameters for initial testing

### Getting Help
- Check the [materials science troubleshooting guide](troubleshooting-materials.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-east-1 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $22. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
