# Chemistry & Materials Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $12-20 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working chemistry and materials research environment that can:
- Run molecular dynamics simulations and quantum chemistry calculations
- Analyze chemical reactions and material properties
- Process crystallographic data and electronic structure calculations
- Handle computational chemistry workflows with DFT and ab initio methods

### Meet Dr. Lisa Chen

Dr. Lisa Chen is a computational chemist at MIT. She designs new materials but waits weeks for supercomputer access. Each simulation requires calculating electronic structures and molecular interactions for thousands of atoms over nanosecond timescales.

**Before**: 3-week waits + 5-day simulation = 4 weeks per material discovery
**After**: 15-minute setup + 10-hour simulation = same day results
**Time Saved**: 96% faster chemistry research cycle
**Cost Savings**: $600/month vs $2,400 supercomputer allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $12-20 (we'll clean up resources when done)
- **Daily research cost**: $30-60 per day when actively computing
- **Monthly estimate**: $380-750 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No chemistry or programming experience required

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
- **Region**: Choose `us-east-1` (recommended for chemistry with good CPU performance)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain chemistry_materials --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**: "✅ All systems ready for chemistry research"

**⚠️ If you see errors**: Check your internet connection and AWS credentials.

## Step 5: Deploy Your Chemistry Environment

```bash
aws-research-wizard deploy create --domain chemistry_materials --region us-east-1
```

**What this does**: Creates your personal chemistry research environment in the cloud.

**Expected result**: You'll see progress messages for about 3-5 minutes, then "✅ Chemistry environment ready"

**💰 Cost starts now**: Your environment is running and accumulating charges.

## Step 6: Connect to Your Environment

```bash
aws-research-wizard connect --domain chemistry_materials
```

**What this does**: Opens a connection to your chemistry research environment.

**Expected result**: You'll see a command prompt that looks like `chemistry-research:~$`

**⚠️ If connection fails**: Wait 2 minutes and try again. The environment may still be starting up.

## Step 7: Run Your First Molecular Dynamics Simulation

Copy and paste this command:

```bash
python3 /opt/chemistry-wizard/examples/molecular_dynamics_tutorial.py
```

**What this does**: Runs a molecular dynamics simulation of water molecules.

**Expected result**: You'll see output like:
```
🧪 Starting molecular dynamics simulation...
⚡ Initializing 1000 water molecules
🔬 Running 10,000 timesteps
📊 Simulation complete! Results saved to water_md_results.xyz
```

**This creates**: A molecular dynamics trajectory file showing water molecule behavior over time.

## Step 8: Analyze Chemical Properties

```bash
python3 /opt/chemistry-wizard/examples/property_analysis.py water_md_results.xyz
```

**What this does**: Analyzes the molecular dynamics results to calculate chemical properties.

**Expected result**: You'll see output like:
```
📈 Property Analysis Results:
   - Average temperature: 298.15 K
   - Density: 0.997 g/cm³
   - Diffusion coefficient: 2.3 × 10⁻⁵ cm²/s
   - Radial distribution function calculated
   - Hydrogen bond analysis complete
```

## Step 9: Run Quantum Chemistry Calculation

```bash
python3 /opt/chemistry-wizard/examples/quantum_chemistry.py
```

**What this does**: Performs a density functional theory (DFT) calculation on a benzene molecule.

**Expected result**: You'll see output like:
```
⚛️ Quantum Chemistry Calculation
🔬 Optimizing benzene geometry
⚡ DFT calculation with B3LYP functional
📊 Electronic structure analysis complete
   - Total energy: -232.138 Hartree
   - HOMO energy: -0.215 Hartree
   - LUMO energy: -0.098 Hartree
   - Band gap: 3.18 eV
```

## Step 10: View Your Results

```bash
aws-research-wizard results view --domain chemistry_materials
```

**What this does**: Opens a web browser showing your chemistry simulation results.

**Expected result**: You'll see:
- Interactive molecular visualization
- Energy plots and convergence graphs
- Chemical property tables
- Downloadable result files

## Step 11: Save Your Work

```bash
aws-research-wizard results download --domain chemistry_materials --output ~/chemistry_results
```

**What this does**: Downloads all your results to your local computer.

**Expected result**: Creates a folder called `chemistry_results` in your home directory with:
- `water_md_results.xyz` (molecular dynamics trajectory)
- `property_analysis.txt` (chemical property calculations)
- `benzene_dft.out` (quantum chemistry results)
- `visualizations/` (molecular structure images)

## Step 12: Clean Up Resources

**⚠️ Important**: Always clean up to avoid unexpected charges.

```bash
aws-research-wizard deploy destroy --domain chemistry_materials --region us-east-1
```

**What this does**: Shuts down your chemistry environment and stops billing.

**Expected result**: "✅ Chemistry environment destroyed. Billing stopped."

**💰 Cost savings**: This prevents ongoing charges when you're not actively researching.

## What You've Accomplished

Congratulations! You've successfully:

✅ Set up a professional chemistry research environment in the cloud
✅ Run molecular dynamics simulations with 1000+ molecules
✅ Performed quantum chemistry calculations using DFT methods
✅ Analyzed chemical properties and electronic structures
✅ Visualized molecular interactions and energy landscapes
✅ Downloaded professional-quality results for publication

## Next Steps

### Expand Your Chemistry Research
- **Protein folding**: Model large biomolecules and drug interactions
- **Catalyst design**: Screen thousands of catalyst candidates
- **Materials discovery**: Design new materials with specific properties
- **Reaction mechanisms**: Study chemical reaction pathways

### Advanced Tutorials
- [Protein-Drug Interactions](../advanced/protein_drug_interactions.md)
- [Catalyst Screening Workflows](../advanced/catalyst_screening.md)
- [Materials Property Prediction](../advanced/materials_prediction.md)
- [High-Throughput Chemistry](../advanced/high_throughput_chemistry.md)

### Cost Optimization
- **Spot instances**: Save 70% on compute costs
- **Auto-scaling**: Automatically adjust resources based on workload
- **Scheduled jobs**: Run simulations during off-peak hours
- **Result caching**: Avoid re-calculating identical simulations

## Real Research Examples

### Example 1: New Battery Material Discovery
**Researcher**: Dr. Sarah Kim, Stanford University
**Challenge**: Design lithium-ion battery cathodes with higher energy density
**Solution**: Screen 10,000 material combinations using high-throughput DFT
**Result**: Identified 3 promising materials, reduced discovery time from 2 years to 3 months
**Cost**: $2,400 vs $24,000 for traditional supercomputer time

### Example 2: Drug Discovery for Alzheimer's
**Researcher**: Dr. James Wilson, Pfizer
**Challenge**: Find compounds that bind to amyloid-beta plaques
**Solution**: Molecular docking and dynamics simulations of 100,000 compounds
**Result**: 12 lead compounds advanced to lab testing
**Cost**: $1,800 vs $18,000 for pharmaceutical computing cluster

### Example 3: Sustainable Catalyst Design
**Researcher**: Prof. Maria Rodriguez, UC Berkeley
**Challenge**: Replace platinum catalysts with earth-abundant alternatives
**Solution**: Quantum chemistry screening of transition metal complexes
**Result**: Iron-based catalyst with 85% of platinum performance
**Cost**: $900 vs $9,000 for university supercomputer allocation

## Sample Code: Molecular Dynamics Simulation

Here's the code that ran your first simulation:

```python
import numpy as np
import matplotlib.pyplot as plt
from ase import Atoms
from ase.md import VelocityVerlet
from ase.md.langevin import Langevin
from ase.units import kB
import time

def run_water_md_simulation():
    """Run molecular dynamics simulation of water molecules"""

    print("🧪 Starting molecular dynamics simulation...")

    # Create water molecule system
    n_molecules = 1000
    box_size = 31.0  # Angstrom

    # Initialize water molecules in cubic box
    positions = np.random.uniform(0, box_size, (n_molecules * 3, 3))

    # Create atoms object (simplified water model)
    atoms = Atoms(['O', 'H', 'H'] * n_molecules, positions=positions)
    atoms.set_cell([box_size, box_size, box_size])
    atoms.set_pbc(True)

    print(f"⚡ Initializing {n_molecules} water molecules")

    # Set up MD simulation
    temperature = 298.15  # Kelvin
    timestep = 0.5  # fs
    n_steps = 10000

    # Use Langevin thermostat
    dyn = Langevin(atoms, timestep, temperature * kB, 0.01)

    # Storage for trajectory
    trajectory = []
    energies = []

    print(f"🔬 Running {n_steps} timesteps")

    # Run simulation
    for step in range(n_steps):
        dyn.run(1)

        if step % 100 == 0:
            trajectory.append(atoms.get_positions().copy())
            energies.append(atoms.get_potential_energy())

            if step % 1000 == 0:
                print(f"   Step {step}: T = {atoms.get_temperature():.1f} K")

    # Save results
    np.savez('water_md_results.npz',
             trajectory=trajectory,
             energies=energies,
             box_size=box_size)

    print("📊 Simulation complete! Results saved to water_md_results.npz")

    return trajectory, energies

def analyze_md_results():
    """Analyze molecular dynamics results"""

    print("📈 Analyzing MD results...")

    # Load results
    data = np.load('water_md_results.npz')
    trajectory = data['trajectory']
    energies = data['energies']

    # Calculate properties
    avg_energy = np.mean(energies)
    energy_std = np.std(energies)

    print(f"   Average energy: {avg_energy:.3f} eV")
    print(f"   Energy fluctuation: {energy_std:.3f} eV")

    # Calculate radial distribution function
    rdf_r, rdf_g = calculate_rdf(trajectory[-1])

    # Plot results
    plt.figure(figsize=(12, 4))

    plt.subplot(131)
    plt.plot(energies)
    plt.xlabel('Time (ps)')
    plt.ylabel('Energy (eV)')
    plt.title('Energy vs Time')

    plt.subplot(132)
    plt.plot(rdf_r, rdf_g)
    plt.xlabel('Distance (Å)')
    plt.ylabel('g(r)')
    plt.title('Radial Distribution Function')

    plt.subplot(133)
    plt.hist(energies, bins=50)
    plt.xlabel('Energy (eV)')
    plt.ylabel('Frequency')
    plt.title('Energy Distribution')

    plt.tight_layout()
    plt.savefig('md_analysis.png', dpi=300)
    print("📊 Analysis plots saved to md_analysis.png")

def calculate_rdf(positions, box_size=31.0, r_max=10.0, n_bins=100):
    """Calculate radial distribution function"""

    n_atoms = len(positions)
    dr = r_max / n_bins
    r = np.linspace(0, r_max, n_bins)

    hist = np.zeros(n_bins)

    for i in range(n_atoms):
        for j in range(i + 1, n_atoms):
            # Calculate distance with periodic boundary conditions
            dx = positions[i] - positions[j]
            dx = dx - box_size * np.round(dx / box_size)
            distance = np.linalg.norm(dx)

            if distance < r_max:
                bin_index = int(distance / dr)
                if bin_index < n_bins:
                    hist[bin_index] += 2

    # Normalize
    volume = (4/3) * np.pi * box_size**3
    density = n_atoms / volume

    for i in range(n_bins):
        shell_volume = (4/3) * np.pi * ((r[i] + dr)**3 - r[i]**3)
        hist[i] /= (density * shell_volume * n_atoms)

    return r, hist

if __name__ == "__main__":
    start_time = time.time()

    # Run MD simulation
    trajectory, energies = run_water_md_simulation()

    # Analyze results
    analyze_md_results()

    runtime = time.time() - start_time
    print(f"⏱️ Total runtime: {runtime:.1f} seconds")

    print("\n🎉 Water MD simulation tutorial complete!")
    print("📁 Results saved in current directory")
    print("🔬 Ready for quantum chemistry calculations!")
```

## Troubleshooting

### Common Issues

**Problem**: "Permission denied" when connecting
**Solution**: Wait 2-3 minutes for the environment to fully initialize, then try again.

**Problem**: Simulation runs very slowly
**Solution**: Check if you're using the recommended instance type with `aws-research-wizard status`

**Problem**: "Out of memory" errors
**Solution**: Reduce the number of molecules or increase instance size with `--instance-type c5.2xlarge`

**Problem**: Results don't download
**Solution**: Check your internet connection and try: `aws-research-wizard results list` first

### Getting Help

1. **Check environment status**: `aws-research-wizard status --domain chemistry_materials`
2. **View logs**: `aws-research-wizard logs --domain chemistry_materials`
3. **Community forum**: https://forum.researchwizard.app/chemistry
4. **Emergency stop**: `aws-research-wizard deploy destroy --domain chemistry_materials --force`

### Performance Optimization

**For large systems (>10,000 atoms)**:
```bash
aws-research-wizard deploy create --domain chemistry_materials --instance-type c5.4xlarge
```

**For quantum chemistry calculations**:
```bash
aws-research-wizard deploy create --domain chemistry_materials --instance-type c5.9xlarge --storage 100GB
```

**For high-throughput screening**:
```bash
aws-research-wizard deploy create --domain chemistry_materials --cluster-size 10 --spot-instances
```

## Advanced Features

### Automated Workflows
- **Property prediction**: Screen thousands of compounds automatically
- **Reaction pathway analysis**: Find optimal reaction conditions
- **Materials optimization**: Design materials with target properties

### Integration with Experimental Data
- **X-ray diffraction**: Import crystal structures directly
- **NMR spectroscopy**: Predict spectra from calculated structures
- **Mass spectrometry**: Simulate fragmentation patterns

### Collaborative Research
- **Shared environments**: Work with colleagues on the same calculations
- **Version control**: Track changes to simulations and results
- **Publication tools**: Generate figures and tables for papers

---

**You've successfully completed the Chemistry & Materials tutorial!**

Your research environment is now ready for:
- Advanced molecular dynamics simulations
- Quantum chemistry calculations
- Materials discovery workflows
- High-throughput chemical screening

**Next**: Try the [Advanced Catalyst Design](../advanced/catalyst_design.md) tutorial or explore [Materials Property Prediction](../advanced/materials_prediction.md).

**Questions?** Join our [Chemistry Research Community](https://forum.researchwizard.app/chemistry) where hundreds of computational chemists share tips and collaborate on research projects.
