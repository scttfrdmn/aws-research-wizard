# Geoscience Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $13-19 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working geoscience research environment that can:
- Model earthquake simulations and seismic wave propagation
- Analyze geological data and mineral exploration datasets
- Process geophysical surveys and subsurface imaging
- Handle geological modeling and hydrogeological studies

### Meet Dr. Emma Wilson

Dr. Emma Wilson is a seismologist at USGS. She models earthquake hazards but waits weeks for supercomputer access. Each simulation requires processing thousands of seismic stations and geological formations across vast regions.

**Before**: 3-week waits + 5-day simulation = 4 weeks per seismic hazard study
**After**: 15-minute setup + 8-hour simulation = same day results
**Time Saved**: 96% faster geoscience research cycle
**Cost Savings**: $400/month vs $1,600 government allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $13-19 (we'll clean up resources when done)
- **Daily research cost**: $26-50 per day when actively modeling
- **Monthly estimate**: $320-650 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No geoscience or programming experience required

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
- **Region**: Choose `us-west-2` (recommended for geoscience with good geological data access)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain geoscience --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**: "✅ All validations passed"

**⚠️ If validation fails**: Check your internet connection and AWS credentials.

## Step 5: Deploy Your Research Environment

```bash
aws-research-wizard deploy create --domain geoscience --region us-west-2 --instance-type c5.2xlarge
```

**What this does**: Creates a cloud computer with geoscience modeling tools.

**Expected result**: You'll see progress updates for about 7 minutes, then "✅ Environment ready"

**💰 Billing starts now**: About $0.34 per hour ($8.16 per day if left running)

**⚠️ If deploy fails**: Run the command again. AWS sometimes has temporary issues.

## Step 6: Connect to Your Environment

```bash
aws-research-wizard connect --domain geoscience
```

**What this does**: Opens a connection to your cloud research environment.

**Expected result**: You'll see a terminal prompt like `[geoscientist@ip-10-0-1-123 ~]$`

**🎉 Success**: You're now inside your geoscience research environment!

## Step 7: Verify Your Tools

Let's make sure all the geoscience tools are working:

```bash
# Check Python geoscience tools
python3 -c "import numpy, scipy, matplotlib, pandas; print('✅ Data analysis tools ready')"

# Check geophysics libraries
python3 -c "import obspy, gdal; print('✅ Geophysics tools ready')"

# Check numerical modeling tools
python3 -c "import numba, joblib; print('✅ High-performance computing tools ready')"
```

**Expected result**: You should see "✅" messages confirming tools are installed.

**⚠️ If tools are missing**: Run `sudo yum update && pip3 install obspy gdal numba` then try again.

## Your First Geoscience Analysis

Let's run a real geoscience simulation to test everything:

### 1. Download Sample Geological Data

```bash
# Create workspace
mkdir -p ~/geoscience/simulations
cd ~/geoscience/simulations

# Download seismic data
wget -O seismic_data.csv "https://research-data.aws-wizard.com/geoscience/seismic_data_sample.csv"

# Download geological survey data
wget -O geological_survey.csv "https://research-data.aws-wizard.com/geoscience/geological_survey_sample.csv"
```

### 2. Earthquake Simulation and Seismic Analysis

Create this Python script for earthquake analysis:

```bash
cat > earthquake_simulator.py << 'EOF'
#!/usr/bin/env python3
"""
Earthquake Simulation and Seismic Analysis Suite
Models earthquake mechanics, wave propagation, and seismic hazard assessment
"""

import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from scipy import signal, integrate, constants
from numba import jit, prange
import time
import warnings
warnings.filterwarnings('ignore')

print("🌍 Initializing earthquake simulation and seismic analysis...")

# Earthquake mechanics and fault modeling
print("\n⚡ Earthquake Mechanics and Fault Modeling")
print("=" * 40)

# Generate synthetic fault system
np.random.seed(42)

# Fault parameters
fault_params = {
    'length': 100000,  # meters
    'width': 15000,    # meters
    'strike': 45,      # degrees
    'dip': 90,         # degrees (vertical fault)
    'rake': 0,         # degrees (strike-slip)
    'depth': 10000,    # meters
    'mu': 3e10,        # shear modulus (Pa)
    'stress_drop': 3e6 # Pa
}

print(f"Fault System Parameters:")
print(f"  Length: {fault_params['length']/1000:.1f} km")
print(f"  Width: {fault_params['width']/1000:.1f} km")
print(f"  Strike: {fault_params['strike']}°")
print(f"  Dip: {fault_params['dip']}°")
print(f"  Depth: {fault_params['depth']/1000:.1f} km")
print(f"  Stress drop: {fault_params['stress_drop']/1e6:.1f} MPa")

# Calculate earthquake magnitude
fault_area = fault_params['length'] * fault_params['width']  # m²
seismic_moment = fault_params['mu'] * fault_area * (fault_params['stress_drop'] / fault_params['mu'])
magnitude = (2/3) * (np.log10(seismic_moment) - 9.1)  # Moment magnitude

print(f"\nEarthquake Characteristics:")
print(f"  Fault area: {fault_area/1e6:.1f} km²")
print(f"  Seismic moment: {seismic_moment:.2e} N·m")
print(f"  Moment magnitude: {magnitude:.1f}")

# Seismic wave propagation simulation
print(f"\n🌊 Seismic Wave Propagation Simulation")
print("=" * 35)

# Wave propagation parameters
wave_params = {
    'vp': 6000,  # P-wave velocity (m/s)
    'vs': 3500,  # S-wave velocity (m/s)
    'rho': 2700, # Density (kg/m³)
    'q': 100,    # Quality factor
    'dt': 0.01,  # Time step (s)
    'duration': 60 # Simulation duration (s)
}

print(f"Wave Propagation Parameters:")
print(f"  P-wave velocity: {wave_params['vp']} m/s")
print(f"  S-wave velocity: {wave_params['vs']} m/s")
print(f"  Density: {wave_params['rho']} kg/m³")
print(f"  Quality factor: {wave_params['q']}")

# Create seismic stations
n_stations = 50
station_distances = np.random.uniform(10000, 200000, n_stations)  # 10-200 km
station_azimuths = np.random.uniform(0, 360, n_stations)

# Convert to Cartesian coordinates
station_x = station_distances * np.cos(np.radians(station_azimuths))
station_y = station_distances * np.sin(np.radians(station_azimuths))

print(f"Seismic Network:")
print(f"  Number of stations: {n_stations}")
print(f"  Distance range: {station_distances.min()/1000:.1f} - {station_distances.max()/1000:.1f} km")

# Seismic wave simulation
@jit(nopython=True)
def ricker_wavelet(t, fc):
    """Generate Ricker wavelet source time function"""
    return (1 - 2 * (np.pi * fc * t)**2) * np.exp(-(np.pi * fc * t)**2)

@jit(nopython=True)
def calculate_travel_time(distance, velocity):
    """Calculate seismic wave travel time"""
    return distance / velocity

@jit(nopython=True)
def geometric_spreading(distance):
    """Calculate geometric spreading attenuation"""
    return 1 / distance

@jit(nopython=True)
def anelastic_attenuation(distance, velocity, frequency, q):
    """Calculate anelastic attenuation"""
    return np.exp(-np.pi * frequency * distance / (q * velocity))

# Generate synthetic seismograms
time_array = np.arange(0, wave_params['duration'], wave_params['dt'])
source_frequency = 2.0  # Hz
source_time_function = ricker_wavelet(time_array - 5, source_frequency)

seismograms = []
for i in range(n_stations):
    distance = station_distances[i]

    # P-wave arrival
    p_travel_time = calculate_travel_time(distance, wave_params['vp'])
    p_amplitude = geometric_spreading(distance) * 0.5  # Simplified amplitude

    # S-wave arrival
    s_travel_time = calculate_travel_time(distance, wave_params['vs'])
    s_amplitude = geometric_spreading(distance) * 1.0  # S-waves typically larger

    # Create seismogram
    seismogram = np.zeros_like(time_array)

    # Add P-wave
    p_start_idx = int(p_travel_time / wave_params['dt'])
    if p_start_idx < len(time_array):
        p_end_idx = min(p_start_idx + len(source_time_function), len(time_array))
        source_end_idx = p_end_idx - p_start_idx
        seismogram[p_start_idx:p_end_idx] += p_amplitude * source_time_function[:source_end_idx]

    # Add S-wave
    s_start_idx = int(s_travel_time / wave_params['dt'])
    if s_start_idx < len(time_array):
        s_end_idx = min(s_start_idx + len(source_time_function), len(time_array))
        source_end_idx = s_end_idx - s_start_idx
        seismogram[s_start_idx:s_end_idx] += s_amplitude * source_time_function[:source_end_idx]

    # Add noise
    noise_level = 0.1 * np.max(np.abs(seismogram))
    seismogram += np.random.normal(0, noise_level, len(seismogram))

    seismograms.append({
        'station_id': i,
        'distance': distance,
        'azimuth': station_azimuths[i],
        'seismogram': seismogram,
        'p_arrival': p_travel_time,
        's_arrival': s_travel_time
    })

print(f"Seismogram Generation:")
print(f"  Time duration: {wave_params['duration']} seconds")
print(f"  Sampling rate: {1/wave_params['dt']:.0f} Hz")
print(f"  Source frequency: {source_frequency} Hz")

# Seismic data analysis
print(f"\n📊 Seismic Data Analysis")
print("=" * 21)

# Calculate magnitude from seismograms
def calculate_local_magnitude(seismogram, distance, dt):
    """Calculate local magnitude from seismogram"""
    # Maximum amplitude
    max_amplitude = np.max(np.abs(seismogram))

    # Richter magnitude formula (simplified)
    ml = np.log10(max_amplitude * 1000) + 3 * np.log10(distance/1000) - 2.92

    return ml

# Calculate magnitudes for all stations
magnitudes = []
for station in seismograms:
    ml = calculate_local_magnitude(station['seismogram'], station['distance'], wave_params['dt'])
    magnitudes.append(ml)

average_magnitude = np.mean(magnitudes)
magnitude_std = np.std(magnitudes)

print(f"Magnitude Analysis:")
print(f"  Calculated magnitude: {average_magnitude:.1f} ± {magnitude_std:.1f}")
print(f"  Theoretical magnitude: {magnitude:.1f}")
print(f"  Magnitude range: {np.min(magnitudes):.1f} - {np.max(magnitudes):.1f}")

# Earthquake location using travel time differences
def locate_earthquake(stations, arrivals, velocities):
    """Locate earthquake using travel time differences"""
    # Simplified location using least squares
    # This is a basic implementation - real location uses more sophisticated methods

    # Use first station as reference
    ref_station = stations[0]

    # Calculate relative travel times
    relative_times = []
    relative_distances = []

    for i in range(1, min(len(stations), 10)):  # Use first 10 stations
        dt = arrivals[i] - arrivals[0]
        dd = stations[i]['distance'] - stations[0]['distance']
        relative_times.append(dt)
        relative_distances.append(dd)

    # Simple location estimate (centroid of stations weighted by amplitude)
    weights = [1/s['distance'] for s in stations[:10]]
    total_weight = sum(weights)

    epicenter_x = sum(s['distance'] * np.cos(np.radians(s['azimuth'])) * w for s, w in zip(stations[:10], weights)) / total_weight
    epicenter_y = sum(s['distance'] * np.sin(np.radians(s['azimuth'])) * w for s, w in zip(stations[:10], weights)) / total_weight

    return epicenter_x, epicenter_y

# Extract P-wave arrival times
p_arrivals = [s['p_arrival'] for s in seismograms]

# Locate earthquake
epicenter_x, epicenter_y = locate_earthquake(seismograms, p_arrivals, wave_params)
epicenter_distance = np.sqrt(epicenter_x**2 + epicenter_y**2)

print(f"\nEarthquake Location:")
print(f"  Epicenter: ({epicenter_x/1000:.1f}, {epicenter_y/1000:.1f}) km")
print(f"  Distance from origin: {epicenter_distance/1000:.1f} km")
print(f"  Location uncertainty: ~{epicenter_distance/1000*0.1:.1f} km")

# Geological modeling
print(f"\n🗻 Geological Modeling")
print("=" * 19)

# Generate synthetic geological cross-section
nx, nz = 100, 50
dx, dz = 1000, 500  # meters
x = np.linspace(0, nx*dx, nx)
z = np.linspace(0, nz*dz, nz)
X, Z = np.meshgrid(x, z)

# Create geological layers
geology = np.zeros((nz, nx))

# Layer 1: Sedimentary rocks (0-5 km depth)
sediment_thickness = 5000 + 1000 * np.sin(2*np.pi*x/50000)  # Variable thickness
for i in range(nx):
    thickness_idx = int(sediment_thickness[i] / dz)
    geology[:thickness_idx, i] = 1

# Layer 2: Metamorphic rocks (5-15 km depth)
for i in range(nx):
    start_idx = int(sediment_thickness[i] / dz)
    end_idx = int(15000 / dz)
    geology[start_idx:end_idx, i] = 2

# Layer 3: Igneous rocks (>15 km depth)
for i in range(nx):
    start_idx = int(15000 / dz)
    geology[start_idx:, i] = 3

# Add fault zone
fault_x = nx // 2
fault_width = 3
geology[:, fault_x-fault_width:fault_x+fault_width] = 4

# Rock properties
rock_properties = {
    0: {'name': 'Air', 'density': 0, 'vp': 0, 'vs': 0},
    1: {'name': 'Sedimentary', 'density': 2200, 'vp': 3000, 'vs': 1500},
    2: {'name': 'Metamorphic', 'density': 2700, 'vp': 6000, 'vs': 3500},
    3: {'name': 'Igneous', 'density': 2800, 'vp': 6500, 'vs': 3800},
    4: {'name': 'Fault Zone', 'density': 2400, 'vp': 4000, 'vs': 2000}
}

print(f"Geological Model:")
print(f"  Model dimensions: {nx} × {nz} grid")
print(f"  Grid spacing: {dx} × {dz} m")
print(f"  Depth range: 0 - {nz*dz/1000:.1f} km")

for rock_id, props in rock_properties.items():
    if rock_id > 0:
        print(f"  {props['name']}: ρ={props['density']} kg/m³, Vp={props['vp']} m/s")

# Hydrogeological analysis
print(f"\n💧 Hydrogeological Analysis")
print("=" * 25)

# Groundwater flow simulation (simplified)
# Darcy's law: q = -K * dh/dx
permeability = np.ones((nz, nx)) * 1e-12  # m²

# Assign permeabilities based on rock type
for i in range(nz):
    for j in range(nx):
        rock_type = geology[i, j]
        if rock_type == 1:  # Sedimentary
            permeability[i, j] = 1e-11
        elif rock_type == 2:  # Metamorphic
            permeability[i, j] = 1e-13
        elif rock_type == 3:  # Igneous
            permeability[i, j] = 1e-14
        elif rock_type == 4:  # Fault zone
            permeability[i, j] = 1e-10

# Hydraulic head calculation (simplified steady-state)
hydraulic_head = np.zeros((nz, nx))
for i in range(nz):
    # Linear hydraulic gradient
    hydraulic_head[i, :] = 100 - (i * dz / 1000)  # Decreasing with depth

# Groundwater flow velocity
mu_water = 1e-3  # Pa·s
rho_water = 1000  # kg/m³
g = 9.81  # m/s²

# Calculate flow velocities
flow_velocity = np.zeros((nz, nx))
for i in range(1, nz-1):
    for j in range(1, nx-1):
        # Hydraulic gradient
        dh_dx = (hydraulic_head[i, j+1] - hydraulic_head[i, j-1]) / (2 * dx)

        # Darcy velocity
        K = permeability[i, j] * rho_water * g / mu_water  # Hydraulic conductivity
        flow_velocity[i, j] = -K * dh_dx

print(f"Hydrogeological Properties:")
print(f"  Permeability range: {np.min(permeability):.1e} - {np.max(permeability):.1e} m²")
print(f"  Hydraulic head range: {np.min(hydraulic_head):.1f} - {np.max(hydraulic_head):.1f} m")
print(f"  Flow velocity range: {np.min(flow_velocity):.1e} - {np.max(flow_velocity):.1e} m/s")

# Mineral exploration analysis
print(f"\n⛏️ Mineral Exploration Analysis")
print("=" * 29)

# Generate synthetic geophysical survey data
survey_points = 200
survey_x = np.linspace(0, nx*dx, survey_points)
survey_data = []

for x_pos in survey_x:
    # Magnetic survey
    magnetic_anomaly = 500 + 200 * np.sin(2*np.pi*x_pos/30000) + 50 * np.random.normal()

    # Gravity survey
    gravity_anomaly = 100 + 50 * np.sin(2*np.pi*x_pos/40000) + 10 * np.random.normal()

    # Resistivity survey
    resistivity = 1000 + 500 * np.sin(2*np.pi*x_pos/20000) + 100 * np.random.normal()

    # Ore probability based on geological and geophysical data
    ore_probability = 0.1
    if abs(magnetic_anomaly - 500) > 100:  # Magnetic anomaly
        ore_probability += 0.3
    if gravity_anomaly > 120:  # Gravity anomaly
        ore_probability += 0.2
    if resistivity < 800:  # Low resistivity
        ore_probability += 0.2

    survey_data.append({
        'x': x_pos,
        'magnetic': magnetic_anomaly,
        'gravity': gravity_anomaly,
        'resistivity': resistivity,
        'ore_probability': min(ore_probability, 1.0)
    })

survey_df = pd.DataFrame(survey_data)

# Identify potential mineral targets
high_potential = survey_df[survey_df['ore_probability'] > 0.6]
medium_potential = survey_df[(survey_df['ore_probability'] > 0.3) & (survey_df['ore_probability'] <= 0.6)]

print(f"Mineral Exploration Results:")
print(f"  Survey points: {len(survey_df)}")
print(f"  High potential targets: {len(high_potential)}")
print(f"  Medium potential targets: {len(medium_potential)}")
print(f"  Average ore probability: {survey_df['ore_probability'].mean():.3f}")

if len(high_potential) > 0:
    print(f"  Best target location: x = {high_potential.iloc[0]['x']/1000:.1f} km")
    print(f"  Target probability: {high_potential.iloc[0]['ore_probability']:.3f}")

# Generate comprehensive geoscience visualization
plt.figure(figsize=(16, 12))

# Seismogram display
plt.subplot(3, 3, 1)
for i in range(5):  # Show first 5 stations
    station = seismograms[i]
    plt.plot(time_array, station['seismogram'] + i*2, 'k-', linewidth=1)
    plt.text(0, i*2, f'St{i+1}', fontsize=8)
plt.title('Seismogram Records')
plt.xlabel('Time (s)')
plt.ylabel('Station + Amplitude')

# Magnitude vs distance
plt.subplot(3, 3, 2)
distances = [s['distance']/1000 for s in seismograms]
plt.scatter(distances, magnitudes, alpha=0.6)
plt.title('Magnitude vs Distance')
plt.xlabel('Distance (km)')
plt.ylabel('Local Magnitude')
plt.grid(True, alpha=0.3)

# Station map
plt.subplot(3, 3, 3)
station_x_km = [s['distance']/1000 * np.cos(np.radians(s['azimuth'])) for s in seismograms]
station_y_km = [s['distance']/1000 * np.sin(np.radians(s['azimuth'])) for s in seismograms]
plt.scatter(station_x_km, station_y_km, c='red', s=20, alpha=0.7)
plt.scatter([epicenter_x/1000], [epicenter_y/1000], c='yellow', s=100, marker='*', label='Epicenter')
plt.title('Seismic Station Network')
plt.xlabel('X Distance (km)')
plt.ylabel('Y Distance (km)')
plt.legend()
plt.axis('equal')

# Geological cross-section
plt.subplot(3, 3, 4)
plt.imshow(geology, extent=[0, nx*dx/1000, nz*dz/1000, 0], cmap='terrain', aspect='auto')
plt.colorbar(label='Rock Type')
plt.title('Geological Cross-Section')
plt.xlabel('Distance (km)')
plt.ylabel('Depth (km)')

# Hydraulic head
plt.subplot(3, 3, 5)
plt.imshow(hydraulic_head, extent=[0, nx*dx/1000, nz*dz/1000, 0], cmap='Blues', aspect='auto')
plt.colorbar(label='Hydraulic Head (m)')
plt.title('Groundwater Flow')
plt.xlabel('Distance (km)')
plt.ylabel('Depth (km)')

# Geophysical surveys
plt.subplot(3, 3, 6)
plt.plot(survey_df['x']/1000, survey_df['magnetic'], 'r-', label='Magnetic')
plt.plot(survey_df['x']/1000, survey_df['gravity']*5, 'g-', label='Gravity×5')
plt.plot(survey_df['x']/1000, survey_df['resistivity']/10, 'b-', label='Resistivity/10')
plt.title('Geophysical Surveys')
plt.xlabel('Distance (km)')
plt.ylabel('Anomaly Value')
plt.legend()
plt.grid(True, alpha=0.3)

# Ore probability
plt.subplot(3, 3, 7)
plt.plot(survey_df['x']/1000, survey_df['ore_probability'], 'purple', linewidth=2)
plt.axhline(y=0.6, color='red', linestyle='--', label='High Potential')
plt.axhline(y=0.3, color='orange', linestyle='--', label='Medium Potential')
plt.title('Mineral Exploration Targets')
plt.xlabel('Distance (km)')
plt.ylabel('Ore Probability')
plt.legend()
plt.grid(True, alpha=0.3)

# Travel time curves
plt.subplot(3, 3, 8)
distances_km = np.array(distances)
p_times = distances_km / (wave_params['vp']/1000)  # Convert to km/s
s_times = distances_km / (wave_params['vs']/1000)
plt.plot(distances_km, p_times, 'bo-', label='P-waves', markersize=3)
plt.plot(distances_km, s_times, 'ro-', label='S-waves', markersize=3)
plt.title('Travel Time Curves')
plt.xlabel('Distance (km)')
plt.ylabel('Travel Time (s)')
plt.legend()
plt.grid(True, alpha=0.3)

# Permeability distribution
plt.subplot(3, 3, 9)
plt.imshow(np.log10(permeability), extent=[0, nx*dx/1000, nz*dz/1000, 0], cmap='YlOrRd', aspect='auto')
plt.colorbar(label='log10(Permeability)')
plt.title('Permeability Distribution')
plt.xlabel('Distance (km)')
plt.ylabel('Depth (km)')

plt.tight_layout()
plt.savefig('geoscience_analysis.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Geoscience analysis dashboard saved as 'geoscience_analysis.png'")

# Seismic hazard assessment
print(f"\n⚠️ Seismic Hazard Assessment")
print("=" * 26)

# Ground motion prediction
def predict_ground_motion(magnitude, distance):
    """Predict peak ground acceleration using empirical relationship"""
    # Simplified Boore-Atkinson Ground Motion Prediction Equation
    log_pga = -3.512 + 0.904 * magnitude - 1.328 * np.log10(distance + 25)
    pga = 10 ** log_pga  # g (acceleration of gravity)
    return pga

# Calculate ground motion for all stations
ground_motions = []
for station in seismograms:
    pga = predict_ground_motion(magnitude, station['distance']/1000)
    ground_motions.append(pga)

max_pga = max(ground_motions)
min_pga = min(ground_motions)

print(f"Ground Motion Prediction:")
print(f"  Maximum PGA: {max_pga:.3f} g")
print(f"  Minimum PGA: {min_pga:.3f} g")
print(f"  Average PGA: {np.mean(ground_motions):.3f} g")

# Seismic intensity classification
def classify_intensity(pga):
    """Classify seismic intensity based on PGA"""
    if pga < 0.017:
        return "I - Not felt"
    elif pga < 0.014:
        return "II - Weak"
    elif pga < 0.039:
        return "III - Light"
    elif pga < 0.092:
        return "IV - Moderate"
    elif pga < 0.18:
        return "V - Strong"
    elif pga < 0.34:
        return "VI - Very Strong"
    else:
        return "VII+ - Severe"

intensity_at_max = classify_intensity(max_pga)
print(f"  Maximum intensity: {intensity_at_max}")

# Risk assessment
vulnerable_stations = sum(1 for pga in ground_motions if pga > 0.1)
high_risk_stations = sum(1 for pga in ground_motions if pga > 0.2)

print(f"\nSeismic Risk Assessment:")
print(f"  Stations with PGA > 0.1g: {vulnerable_stations}/{len(ground_motions)}")
print(f"  Stations with PGA > 0.2g: {high_risk_stations}/{len(ground_motions)}")
print(f"  Risk level: {'High' if high_risk_stations > 5 else 'Moderate' if vulnerable_stations > 10 else 'Low'}")

# Geological time analysis
print(f"\n📅 Geological Time Analysis")
print("=" * 24)

# Simulate geological history
geological_events = [
    {'time': 500, 'event': 'Cambrian marine deposition', 'rock_type': 'Sedimentary'},
    {'time': 400, 'event': 'Devonian mountain building', 'rock_type': 'Metamorphic'},
    {'time': 300, 'event': 'Permian igneous intrusion', 'rock_type': 'Igneous'},
    {'time': 200, 'event': 'Triassic rifting', 'rock_type': 'Sedimentary'},
    {'time': 100, 'event': 'Cretaceous volcanic activity', 'rock_type': 'Igneous'},
    {'time': 50, 'event': 'Tertiary uplift', 'rock_type': 'Metamorphic'},
    {'time': 2, 'event': 'Quaternary glaciation', 'rock_type': 'Sedimentary'}
]

print(f"Geological History (Million Years Ago):")
for event in geological_events:
    print(f"  {event['time']} Ma: {event['event']} ({event['rock_type']})")

# Calculate erosion rates
current_elevation = 2000  # meters
total_uplift = 4000  # meters
erosion_time = 50e6  # years
erosion_rate = (total_uplift - current_elevation) / erosion_time * 1000  # mm/year

print(f"\nErosion Analysis:")
print(f"  Current elevation: {current_elevation} m")
print(f"  Total uplift: {total_uplift} m")
print(f"  Erosion rate: {erosion_rate:.3f} mm/year")

# Resource assessment
print(f"\n💎 Resource Assessment")
print("=" * 18)

# Estimate resource potential
total_survey_area = (nx * dx) * 1000  # m² × 1000 m width
high_potential_area = len(high_potential) * 1000 * 1000  # Assuming 1km spacing
medium_potential_area = len(medium_potential) * 1000 * 1000

resource_density = 0.1  # tonnes per m³
ore_grade = 0.02  # 2% metal content
metal_price = 8000  # $/tonne

# Economic calculation
high_potential_volume = high_potential_area * 100  # 100m depth
medium_potential_volume = medium_potential_area * 100

high_potential_metal = high_potential_volume * resource_density * ore_grade
medium_potential_metal = medium_potential_volume * resource_density * ore_grade

high_potential_value = high_potential_metal * metal_price
medium_potential_value = medium_potential_metal * metal_price

print(f"Resource Potential:")
print(f"  High potential area: {high_potential_area/1e6:.1f} km²")
print(f"  Medium potential area: {medium_potential_area/1e6:.1f} km²")
print(f"  High potential metal: {high_potential_metal:.0f} tonnes")
print(f"  Medium potential metal: {medium_potential_metal:.0f} tonnes")
print(f"  High potential value: ${high_potential_value/1e6:.1f} million")
print(f"  Medium potential value: ${medium_potential_value/1e6:.1f} million")

# Environmental impact assessment
print(f"\n🌱 Environmental Impact Assessment")
print("=" * 33)

# Assess environmental risks
groundwater_risk = "High" if np.mean(flow_velocity) > 1e-6 else "Moderate" if np.mean(flow_velocity) > 1e-8 else "Low"
seismic_risk = "High" if max_pga > 0.2 else "Moderate" if max_pga > 0.1 else "Low"
slope_stability = "Stable" if np.mean(sediment_thickness) < 8000 else "Unstable"

print(f"Environmental Risk Assessment:")
print(f"  Groundwater contamination risk: {groundwater_risk}")
print(f"  Seismic hazard risk: {seismic_risk}")
print(f"  Slope stability: {slope_stability}")

# Recommend mitigation measures
mitigation_measures = []
if groundwater_risk == "High":
    mitigation_measures.append("Groundwater monitoring and containment systems")
if seismic_risk == "High":
    mitigation_measures.append("Seismic-resistant infrastructure design")
if slope_stability == "Unstable":
    mitigation_measures.append("Slope stabilization and monitoring")

print(f"\nRecommended Mitigation Measures:")
for i, measure in enumerate(mitigation_measures, 1):
    print(f"  {i}. {measure}")

if not mitigation_measures:
    print("  No specific mitigation measures required")

print(f"\n✅ Geoscience analysis complete!")
print(f"Analyzed earthquake magnitude {magnitude:.1f} with {n_stations} seismic stations")
print(f"Modeled {nx}×{nz} geological cross-section with {len(survey_df)} survey points")
print(f"Identified {len(high_potential)} high-potential mineral targets")
EOF

chmod +x earthquake_simulator.py
```

### 3. Run the Earthquake Simulation

```bash
python3 earthquake_simulator.py
```

**Expected output**: You should see comprehensive geoscience analysis results.

## What You've Accomplished

🎉 **Congratulations!** You've successfully:

1. ✅ Created a geoscience research environment in the cloud
2. ✅ Simulated earthquake mechanics and seismic wave propagation
3. ✅ Analyzed geological structures and hydrogeological systems
4. ✅ Conducted mineral exploration and resource assessment
5. ✅ Performed seismic hazard analysis and environmental impact assessment
6. ✅ Generated comprehensive geoscience reports and visualizations

### Real Research Applications

Your environment can now handle:
- **Seismology**: Earthquake simulation, wave propagation, hazard assessment
- **Geology**: Structural analysis, stratigraphic modeling, geological mapping
- **Hydrogeology**: Groundwater flow, contaminant transport, aquifer analysis
- **Mineral exploration**: Geophysical surveys, ore deposit modeling
- **Geohazards**: Landslide analysis, subsidence monitoring, risk assessment

### Next Steps for Advanced Research

```bash
# Install specialized geoscience packages
pip3 install obspy pyproj geopandas rasterio

# Set up geoscience modeling software
conda install -c conda-forge gmt seismic-analysis-tools

# Configure geoscience databases
aws-research-wizard tools install --domain geoscience --advanced
```

### Monthly Cost Estimate

For typical geoscience research usage:
- **Light usage** (25 hours/week): ~$320/month
- **Medium usage** (40 hours/week): ~$480/month
- **Heavy usage** (60 hours/week): ~$650/month

## Clean Up Resources

**Important**: Always clean up to avoid unexpected charges!

```bash
# Exit your research environment
exit

# Destroy the research environment
aws-research-wizard deploy destroy --domain geoscience
```

**Expected result**: "✅ Environment destroyed successfully"

**💰 Billing stops**: No more charges after cleanup

## Troubleshooting

### Common Issues

**Problem**: "ObsPy import failed"
**Solution**:
```bash
pip3 install obspy
# OR use conda
conda install -c conda-forge obspy
```

**Problem**: "GDAL not found" errors
**Solution**:
```bash
sudo yum install gdal gdal-devel
pip3 install gdal
```

**Problem**: Large seismic simulations run slowly
**Solution**:
```bash
# Use more CPU cores
aws-research-wizard deploy create --domain geoscience --instance-type c5.4xlarge
```

**Problem**: Memory errors with geological modeling
**Solution**:
```bash
# Use memory-optimized instance
aws-research-wizard deploy create --domain geoscience --instance-type r5.2xlarge
```

### Getting Help

- **Geoscience Community**: [forum.aws-research-wizard.com/geoscience](https://forum.aws-research-wizard.com/geoscience)
- **Technical Support**: [support@aws-research-wizard.com](mailto:support@aws-research-wizard.com)
- **Sample Data**: [research-data.aws-wizard.com/geoscience](https://research-data.aws-wizard.com/geoscience)

### Emergency Stop

If something goes wrong and you want to stop all charges immediately:

```bash
aws-research-wizard emergency-stop --all
```

This will terminate everything and stop billing within 2 minutes.

---

**🌍 Happy geoscience research!**
*You now have a professional-grade geoscience environment that scales with your geological and seismological research needs.*
