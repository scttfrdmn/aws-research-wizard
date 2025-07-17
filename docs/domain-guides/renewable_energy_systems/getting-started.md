# Renewable Energy Systems Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $10-16 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working renewable energy systems research environment that can:
- Model solar and wind energy systems performance
- Analyze energy storage and grid integration solutions
- Process weather data and energy forecasting models
- Handle renewable energy optimization and sustainability assessment

### Meet Dr. Carlos Rodriguez

Dr. Carlos Rodriguez is a renewable energy engineer at NREL. He designs solar farms but waits weeks for simulation resources. Each project requires modeling thousands of solar panels and wind turbines under varying weather conditions.

**Before**: 2-week waits + 4-day simulation = 3 weeks per energy system design
**After**: 15-minute setup + 6-hour simulation = same day results
**Time Saved**: 95% faster renewable energy research cycle
**Cost Savings**: $350/month vs $1,400 government allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $10-16 (we'll clean up resources when done)
- **Daily research cost**: $20-40 per day when actively modeling
- **Monthly estimate**: $250-500 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No renewable energy or programming experience required

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
- **Region**: Choose `us-west-2` (recommended for renewable energy with good solar irradiance data)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain renewable_energy_systems --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**: "✅ All validations passed"

**⚠️ If validation fails**: Check your internet connection and AWS credentials.

## Step 5: Deploy Your Research Environment

```bash
aws-research-wizard deploy create --domain renewable_energy_systems --region us-west-2 --instance-type c5.2xlarge
```

**What this does**: Creates a cloud computer with renewable energy modeling tools.

**Expected result**: You'll see progress updates for about 6 minutes, then "✅ Environment ready"

**💰 Billing starts now**: About $0.34 per hour ($8.16 per day if left running)

**⚠️ If deploy fails**: Run the command again. AWS sometimes has temporary issues.

## Step 6: Connect to Your Environment

```bash
aws-research-wizard connect --domain renewable_energy_systems
```

**What this does**: Opens a connection to your cloud research environment.

**Expected result**: You'll see a terminal prompt like `[renewables@ip-10-0-1-123 ~]$`

**🎉 Success**: You're now inside your renewable energy research environment!

## Step 7: Verify Your Tools

Let's make sure all the renewable energy tools are working:

```bash
# Check Python energy modeling tools
python3 -c "import numpy, scipy, matplotlib, pandas; print('✅ Data analysis tools ready')"

# Check renewable energy libraries
python3 -c "import pvlib, windpowerlib; print('✅ Renewable energy tools ready')"

# Check optimization tools
python3 -c "import scipy.optimize, pulp; print('✅ Optimization tools ready')"
```

**Expected result**: You should see "✅" messages confirming tools are installed.

**⚠️ If tools are missing**: Run `sudo yum update && pip3 install pvlib windpowerlib pulp` then try again.

## Your First Renewable Energy Analysis

Let's run a real renewable energy system analysis to test everything:

### 1. Download Sample Energy Data

```bash
# Create workspace
mkdir -p ~/renewable_energy/analysis
cd ~/renewable_energy/analysis

# Download weather data
wget -O weather_data.csv "https://research-data.aws-wizard.com/renewable_energy/weather_data_sample.csv"

# Download load profile
wget -O load_profile.csv "https://research-data.aws-wizard.com/renewable_energy/load_profile_sample.csv"
```

### 2. Solar Energy System Analysis

Create this Python script for solar energy analysis:

```bash
cat > solar_energy_analyzer.py << 'EOF'
#!/usr/bin/env python3
"""
Solar Energy System Analysis Suite
Models solar PV performance, energy storage, and grid integration
"""

import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from scipy import optimize
import time
import warnings
warnings.filterwarnings('ignore')

print("☀️ Initializing solar energy system analysis...")

# Solar energy system modeling
print("\n🔋 Solar Photovoltaic System Analysis")
print("=" * 35)

# Generate synthetic solar irradiance data
np.random.seed(42)
hours = np.arange(0, 24, 0.5)  # 30-minute intervals
days = 365

# Solar irradiance modeling
def solar_irradiance_model(hour, day_of_year, latitude=37.7):
    """Model solar irradiance based on time and location"""
    # Solar declination angle
    declination = 23.45 * np.sin(np.radians(360 * (284 + day_of_year) / 365))

    # Hour angle
    hour_angle = 15 * (hour - 12)  # degrees

    # Solar elevation angle
    elevation = np.arcsin(
        np.sin(np.radians(latitude)) * np.sin(np.radians(declination)) +
        np.cos(np.radians(latitude)) * np.cos(np.radians(declination)) * np.cos(np.radians(hour_angle))
    )

    # Clear sky irradiance
    if elevation > 0:
        # Simplified clear sky model
        clear_sky = 1000 * np.sin(elevation) * (0.7 ** (1 / np.sin(elevation)))
    else:
        clear_sky = 0

    # Add weather effects
    cloud_factor = 0.8 + 0.2 * np.random.random()  # 80-100% of clear sky

    return max(0, clear_sky * cloud_factor)

# Create solar data
solar_data = []
for day in range(1, days + 1):
    for hour in hours:
        irradiance = solar_irradiance_model(hour, day)
        ambient_temp = 25 + 10 * np.sin(2 * np.pi * (hour - 6) / 24) + 5 * np.random.normal()
        wind_speed = 3 + 2 * np.random.exponential()

        solar_data.append({
            'day': day,
            'hour': hour,
            'irradiance': irradiance,
            'ambient_temp': ambient_temp,
            'wind_speed': wind_speed
        })

solar_df = pd.DataFrame(solar_data)

print(f"Solar data generated: {len(solar_df)} data points")
print(f"  Average irradiance: {solar_df['irradiance'].mean():.1f} W/m²")
print(f"  Peak irradiance: {solar_df['irradiance'].max():.1f} W/m²")
print(f"  Average temperature: {solar_df['ambient_temp'].mean():.1f} °C")

# PV system parameters
pv_params = {
    'rated_power': 300,  # Watts per panel
    'efficiency': 0.20,  # 20% efficiency
    'temperature_coeff': -0.004,  # %/°C
    'area': 1.5,  # m² per panel
    'num_panels': 1000,  # Total panels
    'inverter_efficiency': 0.96,
    'system_losses': 0.15  # Cable losses, soiling, etc.
}

print(f"\nPV System Configuration:")
print(f"  System size: {pv_params['rated_power'] * pv_params['num_panels'] / 1000:.0f} kW")
print(f"  Number of panels: {pv_params['num_panels']}")
print(f"  Panel efficiency: {pv_params['efficiency']*100:.1f}%")
print(f"  Total area: {pv_params['area'] * pv_params['num_panels']:.0f} m²")

# PV performance calculation
def calculate_pv_output(irradiance, temp, pv_params):
    """Calculate PV power output"""
    if irradiance <= 0:
        return 0

    # Temperature effect
    temp_effect = 1 + pv_params['temperature_coeff'] * (temp - 25)

    # Power output
    dc_power = (irradiance / 1000) * pv_params['rated_power'] * temp_effect

    # System losses
    dc_power *= (1 - pv_params['system_losses'])

    # Inverter efficiency
    ac_power = dc_power * pv_params['inverter_efficiency']

    return max(0, ac_power)

# Calculate PV output for all data points
solar_df['pv_output'] = solar_df.apply(
    lambda row: calculate_pv_output(row['irradiance'], row['ambient_temp'], pv_params) * pv_params['num_panels'],
    axis=1
)

# Daily and monthly analysis
solar_df['date'] = pd.to_datetime(solar_df['day'], unit='D', origin='2023-01-01')
solar_df['month'] = solar_df['date'].dt.month

daily_generation = solar_df.groupby('day')['pv_output'].sum() / 1000  # kWh per day
monthly_generation = solar_df.groupby('month')['pv_output'].sum() / 1000  # kWh per month

print(f"\nPV Performance Analysis:")
print(f"  Daily generation: {daily_generation.mean():.1f} kWh/day average")
print(f"  Annual generation: {daily_generation.sum():.0f} kWh/year")
print(f"  Capacity factor: {daily_generation.mean() * 24 / (pv_params['rated_power'] * pv_params['num_panels'] / 1000) * 100:.1f}%")

# Peak vs off-peak analysis
solar_df['hour_int'] = solar_df['hour'].astype(int)
peak_hours = solar_df[(solar_df['hour_int'] >= 12) & (solar_df['hour_int'] <= 16)]
off_peak_hours = solar_df[(solar_df['hour_int'] < 6) | (solar_df['hour_int'] > 20)]

peak_generation = peak_hours.groupby('day')['pv_output'].sum().mean() / 1000
off_peak_generation = off_peak_hours.groupby('day')['pv_output'].sum().mean() / 1000

print(f"  Peak hours generation: {peak_generation:.1f} kWh/day")
print(f"  Off-peak generation: {off_peak_generation:.1f} kWh/day")
print(f"  Peak/off-peak ratio: {peak_generation/max(off_peak_generation, 0.1):.1f}")

# Wind energy system analysis
print(f"\n💨 Wind Energy System Analysis")
print("=" * 28)

# Wind turbine parameters
wind_params = {
    'rated_power': 2000,  # kW
    'hub_height': 80,  # meters
    'rotor_diameter': 90,  # meters
    'cut_in_speed': 3,  # m/s
    'rated_speed': 12,  # m/s
    'cut_out_speed': 25,  # m/s
    'num_turbines': 25
}

print(f"Wind System Configuration:")
print(f"  System size: {wind_params['rated_power'] * wind_params['num_turbines'] / 1000:.0f} MW")
print(f"  Number of turbines: {wind_params['num_turbines']}")
print(f"  Hub height: {wind_params['hub_height']} m")
print(f"  Rotor diameter: {wind_params['rotor_diameter']} m")

# Wind power curve
def wind_power_curve(wind_speed, wind_params):
    """Calculate wind turbine power output"""
    if wind_speed < wind_params['cut_in_speed']:
        return 0
    elif wind_speed > wind_params['cut_out_speed']:
        return 0
    elif wind_speed >= wind_params['rated_speed']:
        return wind_params['rated_power']
    else:
        # Cubic power curve approximation
        power_ratio = ((wind_speed - wind_params['cut_in_speed']) /
                      (wind_params['rated_speed'] - wind_params['cut_in_speed'])) ** 3
        return wind_params['rated_power'] * power_ratio

# Generate wind data
wind_data = []
for day in range(1, days + 1):
    for hour in hours:
        # Wind speed varies by time of day and season
        base_wind = 8 + 2 * np.sin(2 * np.pi * day / 365)  # Seasonal variation
        daily_variation = 1 + 0.3 * np.sin(2 * np.pi * hour / 24)
        wind_speed = base_wind * daily_variation + np.random.normal(0, 1)
        wind_speed = max(0, wind_speed)

        # Wind direction (affects turbine efficiency)
        wind_direction = np.random.uniform(0, 360)

        wind_data.append({
            'day': day,
            'hour': hour,
            'wind_speed': wind_speed,
            'wind_direction': wind_direction
        })

wind_df = pd.DataFrame(wind_data)

# Calculate wind power output
wind_df['wind_output'] = wind_df['wind_speed'].apply(
    lambda ws: wind_power_curve(ws, wind_params) * wind_params['num_turbines']
)

print(f"\nWind Performance Analysis:")
wind_daily = wind_df.groupby('day')['wind_output'].sum() / 1000  # MWh per day
wind_monthly = wind_df.groupby('day')['wind_output'].sum().rolling(30).mean() / 1000

print(f"  Daily generation: {wind_daily.mean():.1f} MWh/day average")
print(f"  Annual generation: {wind_daily.sum():.0f} MWh/year")
print(f"  Capacity factor: {wind_daily.mean() * 24 / (wind_params['rated_power'] * wind_params['num_turbines'] / 1000) * 100:.1f}%")
print(f"  Average wind speed: {wind_df['wind_speed'].mean():.1f} m/s")

# Energy storage analysis
print(f"\n🔋 Energy Storage System Analysis")
print("=" * 32)

# Battery storage parameters
battery_params = {
    'capacity': 10000,  # kWh
    'power_rating': 5000,  # kW
    'efficiency': 0.90,  # Round-trip efficiency
    'depth_of_discharge': 0.8,  # Max DOD
    'self_discharge': 0.02,  # %/day
    'cycles_per_day': 1
}

print(f"Battery Storage Configuration:")
print(f"  Capacity: {battery_params['capacity']} kWh")
print(f"  Power rating: {battery_params['power_rating']} kW")
print(f"  Round-trip efficiency: {battery_params['efficiency']*100:.0f}%")
print(f"  Usable capacity: {battery_params['capacity'] * battery_params['depth_of_discharge']} kWh")

# Load profile generation
load_data = []
for day in range(1, days + 1):
    for hour in hours:
        # Typical residential/commercial load profile
        if 6 <= hour <= 9:  # Morning peak
            base_load = 0.8
        elif 17 <= hour <= 21:  # Evening peak
            base_load = 1.0
        elif 22 <= hour <= 6:  # Night
            base_load = 0.4
        else:  # Daytime
            base_load = 0.6

        # Seasonal variation
        seasonal_factor = 1 + 0.2 * np.sin(2 * np.pi * day / 365)

        # Random variation
        load = base_load * seasonal_factor * 8000 + np.random.normal(0, 500)  # kW
        load = max(0, load)

        load_data.append({
            'day': day,
            'hour': hour,
            'load': load
        })

load_df = pd.DataFrame(load_data)

# Combine renewable generation and load
energy_df = pd.merge(solar_df[['day', 'hour', 'pv_output']],
                    wind_df[['day', 'hour', 'wind_output']],
                    on=['day', 'hour'])
energy_df = pd.merge(energy_df, load_df[['day', 'hour', 'load']],
                    on=['day', 'hour'])

# Total renewable generation
energy_df['total_generation'] = energy_df['pv_output'] + energy_df['wind_output']
energy_df['net_load'] = energy_df['load'] - energy_df['total_generation']

print(f"\nEnergy Balance Analysis:")
print(f"  Average load: {energy_df['load'].mean():.0f} kW")
print(f"  Average generation: {energy_df['total_generation'].mean():.0f} kW")
print(f"  Average net load: {energy_df['net_load'].mean():.0f} kW")

# Battery operation simulation
def simulate_battery_operation(energy_data, battery_params):
    """Simulate battery charging/discharging"""
    soc = 0.5  # Start at 50% state of charge
    battery_operation = []

    for _, row in energy_data.iterrows():
        net_load = row['net_load']

        # Charging (excess generation)
        if net_load < 0:
            available_power = min(abs(net_load), battery_params['power_rating'])
            max_charge = (1 - soc) * battery_params['capacity']
            charge = min(available_power * 0.5, max_charge) * battery_params['efficiency']  # 30-min interval
            soc += charge / battery_params['capacity']
            battery_action = charge

        # Discharging (load exceeds generation)
        else:
            available_energy = soc * battery_params['capacity'] * battery_params['depth_of_discharge']
            max_discharge = min(net_load, battery_params['power_rating'])
            discharge = min(max_discharge * 0.5, available_energy) / battery_params['efficiency']  # 30-min interval
            soc -= discharge / battery_params['capacity']
            battery_action = -discharge

        # Self-discharge
        soc *= (1 - battery_params['self_discharge'] / 48)  # Per 30-min interval

        # Keep SOC within bounds
        soc = max(0, min(1, soc))

        battery_operation.append({
            'soc': soc,
            'battery_power': battery_action,
            'grid_power': net_load - battery_action
        })

    return pd.DataFrame(battery_operation)

# Run battery simulation
battery_df = simulate_battery_operation(energy_df, battery_params)
energy_df = pd.concat([energy_df, battery_df], axis=1)

print(f"\nBattery Performance:")
print(f"  Average SOC: {battery_df['soc'].mean()*100:.1f}%")
print(f"  Daily cycles: {battery_params['cycles_per_day']:.1f}")
print(f"  Grid import reduction: {abs(energy_df['grid_power'].mean()) / energy_df['net_load'].mean() * 100:.1f}%")

# Economic analysis
print(f"\n💰 Economic Analysis")
print("=" * 17)

# Cost parameters (USD)
costs = {
    'solar_capex': 1500,  # $/kW
    'wind_capex': 1200,   # $/kW
    'battery_capex': 300, # $/kWh
    'solar_opex': 20,     # $/kW/year
    'wind_opex': 30,      # $/kW/year
    'battery_opex': 5,    # $/kWh/year
    'electricity_price': 0.12  # $/kWh
}

# Calculate system costs
solar_cost = (pv_params['rated_power'] * pv_params['num_panels'] / 1000) * costs['solar_capex']
wind_cost = (wind_params['rated_power'] * wind_params['num_turbines'] / 1000) * costs['wind_capex']
battery_cost = battery_params['capacity'] * costs['battery_capex']

total_capex = solar_cost + wind_cost + battery_cost

# Annual costs
solar_opex = (pv_params['rated_power'] * pv_params['num_panels'] / 1000) * costs['solar_opex']
wind_opex = (wind_params['rated_power'] * wind_params['num_turbines'] / 1000) * costs['wind_opex']
battery_opex = battery_params['capacity'] * costs['battery_opex']

total_opex = solar_opex + wind_opex + battery_opex

# Energy value
annual_generation = (daily_generation.sum() + wind_daily.sum() * 1000) / 1000  # MWh
annual_value = annual_generation * 1000 * costs['electricity_price']  # kWh * $/kWh

print(f"Capital Costs:")
print(f"  Solar system: ${solar_cost/1000:.0f}k")
print(f"  Wind system: ${wind_cost/1000:.0f}k")
print(f"  Battery system: ${battery_cost/1000:.0f}k")
print(f"  Total CAPEX: ${total_capex/1000:.0f}k")

print(f"\nAnnual Costs:")
print(f"  Solar O&M: ${solar_opex/1000:.0f}k")
print(f"  Wind O&M: ${wind_opex/1000:.0f}k")
print(f"  Battery O&M: ${battery_opex/1000:.0f}k")
print(f"  Total OPEX: ${total_opex/1000:.0f}k")

print(f"\nEconomic Performance:")
print(f"  Annual generation: {annual_generation:.0f} MWh")
print(f"  Annual value: ${annual_value/1000:.0f}k")
print(f"  Net annual cash flow: ${(annual_value - total_opex)/1000:.0f}k")

# LCOE calculation
project_life = 25  # years
discount_rate = 0.06

# Present value calculation
pv_capex = total_capex
pv_opex = sum(total_opex / (1 + discount_rate)**year for year in range(1, project_life + 1))
pv_generation = sum(annual_generation / (1 + discount_rate)**year for year in range(1, project_life + 1))

lcoe = (pv_capex + pv_opex) / pv_generation  # $/MWh

print(f"  LCOE: ${lcoe:.2f}/MWh")
print(f"  Payback period: {total_capex / (annual_value - total_opex):.1f} years")

# Grid integration analysis
print(f"\n🔌 Grid Integration Analysis")
print("=" * 26)

# Grid stability metrics
generation_variability = energy_df['total_generation'].std() / energy_df['total_generation'].mean()
load_variability = energy_df['load'].std() / energy_df['load'].mean()
net_load_variability = energy_df['net_load'].std() / energy_df['net_load'].mean()

print(f"Variability Analysis:")
print(f"  Generation variability: {generation_variability:.3f}")
print(f"  Load variability: {load_variability:.3f}")
print(f"  Net load variability: {net_load_variability:.3f}")

# Grid services
# Frequency regulation capability
freq_regulation = min(
    battery_params['power_rating'] * 0.1,  # 10% of battery power
    (pv_params['rated_power'] * pv_params['num_panels'] / 1000) * 0.05  # 5% of solar capacity
)

# Voltage support
voltage_support = energy_df['total_generation'].mean() * 0.1  # 10% of average generation

print(f"\nGrid Services Capability:")
print(f"  Frequency regulation: {freq_regulation:.0f} kW")
print(f"  Voltage support: {voltage_support:.0f} kW")
print(f"  Grid import reduction: {(1 - abs(energy_df['grid_power'].mean()) / energy_df['load'].mean()) * 100:.1f}%")

# Generate comprehensive renewable energy visualization
plt.figure(figsize=(16, 12))

# Daily generation profiles
plt.subplot(3, 3, 1)
sample_day = energy_df[energy_df['day'] == 180]  # Mid-year sample
plt.plot(sample_day['hour'], sample_day['pv_output']/1000, 'orange', label='Solar', linewidth=2)
plt.plot(sample_day['hour'], sample_day['wind_output']/1000, 'blue', label='Wind', linewidth=2)
plt.plot(sample_day['hour'], sample_day['load']/1000, 'red', label='Load', linewidth=2)
plt.title('Daily Generation Profile')
plt.xlabel('Hour of Day')
plt.ylabel('Power (MW)')
plt.legend()
plt.grid(True, alpha=0.3)

# Monthly generation
plt.subplot(3, 3, 2)
monthly_solar = solar_df.groupby('month')['pv_output'].sum() / 1000000  # MWh
monthly_wind = wind_df.groupby(wind_df['day'] // 30 + 1)['wind_output'].sum() / 1000000  # MWh
months = range(1, 13)
plt.bar(months, monthly_solar, alpha=0.7, label='Solar', color='orange')
plt.bar(months, monthly_wind, bottom=monthly_solar, alpha=0.7, label='Wind', color='blue')
plt.title('Monthly Generation')
plt.xlabel('Month')
plt.ylabel('Generation (MWh)')
plt.legend()

# Battery state of charge
plt.subplot(3, 3, 3)
sample_week = energy_df[energy_df['day'].between(180, 186)]
plt.plot(range(len(sample_week)), sample_week['soc']*100, 'green', linewidth=2)
plt.title('Battery State of Charge (Sample Week)')
plt.xlabel('Time (30-min intervals)')
plt.ylabel('SOC (%)')
plt.grid(True, alpha=0.3)

# Generation duration curves
plt.subplot(3, 3, 4)
solar_sorted = np.sort(energy_df['pv_output'])[::-1] / 1000
wind_sorted = np.sort(energy_df['wind_output'])[::-1] / 1000
total_sorted = np.sort(energy_df['total_generation'])[::-1] / 1000
hours_year = np.arange(1, len(solar_sorted) + 1)
plt.plot(hours_year, solar_sorted, 'orange', label='Solar', linewidth=2)
plt.plot(hours_year, wind_sorted, 'blue', label='Wind', linewidth=2)
plt.plot(hours_year, total_sorted, 'black', label='Total', linewidth=2)
plt.title('Generation Duration Curves')
plt.xlabel('Hours per Year')
plt.ylabel('Power (MW)')
plt.legend()
plt.grid(True, alpha=0.3)

# Net load histogram
plt.subplot(3, 3, 5)
plt.hist(energy_df['net_load']/1000, bins=50, alpha=0.7, color='purple')
plt.title('Net Load Distribution')
plt.xlabel('Net Load (MW)')
plt.ylabel('Frequency')
plt.grid(True, alpha=0.3)

# Economic breakdown
plt.subplot(3, 3, 6)
cost_categories = ['Solar CAPEX', 'Wind CAPEX', 'Battery CAPEX', 'Annual OPEX']
cost_values = [solar_cost/1000, wind_cost/1000, battery_cost/1000, total_opex/1000]
colors = ['orange', 'blue', 'green', 'red']
plt.bar(cost_categories, cost_values, color=colors)
plt.title('Cost Breakdown')
plt.ylabel('Cost ($k)')
plt.xticks(rotation=45)

# Grid power flow
plt.subplot(3, 3, 7)
sample_day = energy_df[energy_df['day'] == 180]
plt.plot(sample_day['hour'], sample_day['grid_power']/1000, 'purple', linewidth=2)
plt.axhline(y=0, color='black', linestyle='-', alpha=0.3)
plt.title('Grid Power Flow (Sample Day)')
plt.xlabel('Hour of Day')
plt.ylabel('Grid Power (MW)')
plt.grid(True, alpha=0.3)

# Capacity factors
plt.subplot(3, 3, 8)
cf_solar = daily_generation.mean() * 24 / (pv_params['rated_power'] * pv_params['num_panels'] / 1000) * 100
cf_wind = wind_daily.mean() * 24 / (wind_params['rated_power'] * wind_params['num_turbines'] / 1000) * 100
cf_battery = battery_df['soc'].mean() * 100
systems = ['Solar', 'Wind', 'Battery Utilization']
cf_values = [cf_solar, cf_wind, cf_battery]
colors = ['orange', 'blue', 'green']
plt.bar(systems, cf_values, color=colors)
plt.title('System Performance')
plt.ylabel('Capacity Factor / Utilization (%)')

# Wind power curve
plt.subplot(3, 3, 9)
wind_speeds = np.linspace(0, 30, 100)
power_curve = [wind_power_curve(ws, wind_params) for ws in wind_speeds]
plt.plot(wind_speeds, power_curve, 'blue', linewidth=2)
plt.title('Wind Turbine Power Curve')
plt.xlabel('Wind Speed (m/s)')
plt.ylabel('Power (kW)')
plt.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig('renewable_energy_analysis.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Renewable energy analysis saved as 'renewable_energy_analysis.png'")

# Optimization analysis
print(f"\n🎯 System Optimization")
print("=" * 19)

# Optimize system sizing
def optimize_system_sizing(solar_mult, wind_mult, battery_mult):
    """Optimize renewable energy system sizing"""
    # Scaled system parameters
    solar_capacity = (pv_params['rated_power'] * pv_params['num_panels'] / 1000) * solar_mult
    wind_capacity = (wind_params['rated_power'] * wind_params['num_turbines'] / 1000) * wind_mult
    battery_capacity = battery_params['capacity'] * battery_mult

    # Costs
    total_cost = (solar_capacity * costs['solar_capex'] +
                 wind_capacity * costs['wind_capex'] +
                 battery_capacity * costs['battery_capex'])

    # Generation
    annual_gen = (daily_generation.sum() * solar_mult +
                 wind_daily.sum() * 1000 * wind_mult) / 1000  # MWh

    # LCOE
    lcoe_opt = total_cost / (annual_gen * project_life)

    return lcoe_opt, total_cost, annual_gen

# Grid search optimization
solar_range = np.linspace(0.5, 2.0, 10)
wind_range = np.linspace(0.5, 2.0, 10)
battery_range = np.linspace(0.5, 2.0, 10)

best_lcoe = float('inf')
best_config = None

for s in solar_range:
    for w in wind_range:
        for b in battery_range:
            lcoe_opt, cost, gen = optimize_system_sizing(s, w, b)
            if lcoe_opt < best_lcoe:
                best_lcoe = lcoe_opt
                best_config = (s, w, b, cost, gen)

print(f"Optimization Results:")
print(f"  Best LCOE: ${best_lcoe:.2f}/MWh")
print(f"  Optimal solar scaling: {best_config[0]:.1f}x")
print(f"  Optimal wind scaling: {best_config[1]:.1f}x")
print(f"  Optimal battery scaling: {best_config[2]:.1f}x")
print(f"  Optimal system cost: ${best_config[3]/1000:.0f}k")
print(f"  Optimal annual generation: {best_config[4]:.0f} MWh")

# Sensitivity analysis
print(f"\n📈 Sensitivity Analysis")
print("=" * 20)

# Parameters to analyze
base_params = {
    'solar_cost': costs['solar_capex'],
    'wind_cost': costs['wind_capex'],
    'battery_cost': costs['battery_capex'],
    'electricity_price': costs['electricity_price'],
    'discount_rate': discount_rate
}

# Sensitivity ranges
sensitivity_results = {}
for param, base_value in base_params.items():
    variations = np.linspace(0.5, 1.5, 11)  # ±50% variation
    lcoe_variations = []

    for variation in variations:
        # Modify parameter
        modified_costs = costs.copy()
        modified_rate = discount_rate

        if param == 'solar_cost':
            modified_costs['solar_capex'] = base_value * variation
        elif param == 'wind_cost':
            modified_costs['wind_capex'] = base_value * variation
        elif param == 'battery_cost':
            modified_costs['battery_capex'] = base_value * variation
        elif param == 'electricity_price':
            modified_costs['electricity_price'] = base_value * variation
        elif param == 'discount_rate':
            modified_rate = base_value * variation

        # Recalculate LCOE
        solar_cost_mod = (pv_params['rated_power'] * pv_params['num_panels'] / 1000) * modified_costs['solar_capex']
        wind_cost_mod = (wind_params['rated_power'] * wind_params['num_turbines'] / 1000) * modified_costs['wind_capex']
        battery_cost_mod = battery_params['capacity'] * modified_costs['battery_capex']

        total_capex_mod = solar_cost_mod + wind_cost_mod + battery_cost_mod

        pv_opex_mod = sum(total_opex / (1 + modified_rate)**year for year in range(1, project_life + 1))
        pv_gen_mod = sum(annual_generation / (1 + modified_rate)**year for year in range(1, project_life + 1))

        lcoe_mod = (total_capex_mod + pv_opex_mod) / pv_gen_mod
        lcoe_variations.append(lcoe_mod)

    sensitivity_results[param] = {
        'variations': variations,
        'lcoe': lcoe_variations
    }

print("LCOE Sensitivity (±50% parameter variation):")
for param, results in sensitivity_results.items():
    min_lcoe = min(results['lcoe'])
    max_lcoe = max(results['lcoe'])
    sensitivity = (max_lcoe - min_lcoe) / lcoe * 100
    print(f"  {param}: ±{sensitivity:.1f}% LCOE change")

# Environmental impact
print(f"\n🌍 Environmental Impact")
print("=" * 21)

# CO2 emissions avoided
grid_emission_factor = 0.5  # kg CO2/kWh (typical grid mix)
annual_co2_avoided = annual_generation * 1000 * grid_emission_factor / 1000  # tonnes CO2

# Lifecycle emissions
solar_lifecycle = 0.05  # kg CO2/kWh
wind_lifecycle = 0.02   # kg CO2/kWh
battery_lifecycle = 0.1  # kg CO2/kWh (manufacturing)

lifecycle_emissions = (
    daily_generation.sum() * solar_lifecycle +
    wind_daily.sum() * 1000 * wind_lifecycle +
    battery_params['capacity'] * battery_lifecycle * 365 / 1000  # Annual battery impact
) / 1000  # tonnes CO2

net_co2_avoided = annual_co2_avoided - lifecycle_emissions

print(f"Environmental Impact:")
print(f"  Annual CO2 avoided: {annual_co2_avoided:.0f} tonnes")
print(f"  Lifecycle emissions: {lifecycle_emissions:.0f} tonnes")
print(f"  Net CO2 avoided: {net_co2_avoided:.0f} tonnes")
print(f"  Emission factor: {net_co2_avoided / annual_generation:.3f} kg CO2/kWh")

# Land use
solar_land_use = pv_params['area'] * pv_params['num_panels'] / 10000  # hectares
wind_land_use = wind_params['num_turbines'] * 0.5  # hectares (simplified)
total_land_use = solar_land_use + wind_land_use

print(f"\nLand Use:")
print(f"  Solar land use: {solar_land_use:.1f} hectares")
print(f"  Wind land use: {wind_land_use:.1f} hectares")
print(f"  Total land use: {total_land_use:.1f} hectares")
print(f"  Power density: {(solar_capacity + wind_capacity) / total_land_use:.1f} MW/hectare")

print(f"\n✅ Renewable energy analysis complete!")
print(f"Analyzed solar ({pv_params['num_panels']} panels) and wind ({wind_params['num_turbines']} turbines)")
print(f"Optimized system configuration with {battery_params['capacity']} kWh storage")
print(f"Achieved LCOE of ${lcoe:.2f}/MWh with {net_co2_avoided:.0f} tonnes CO2 avoided annually")
EOF

chmod +x solar_energy_analyzer.py
```

### 3. Run the Solar Energy Analysis

```bash
python3 solar_energy_analyzer.py
```

**Expected output**: You should see comprehensive renewable energy system analysis results.

## What You've Accomplished

🎉 **Congratulations!** You've successfully:

1. ✅ Created a renewable energy systems research environment in the cloud
2. ✅ Modeled solar PV and wind energy system performance
3. ✅ Analyzed energy storage and grid integration solutions
4. ✅ Conducted economic optimization and sensitivity analysis
5. ✅ Assessed environmental impact and sustainability metrics
6. ✅ Generated comprehensive renewable energy reports and visualizations

### Real Research Applications

Your environment can now handle:
- **Solar energy modeling**: PV performance, irradiance analysis, system sizing
- **Wind energy analysis**: Power curves, capacity factors, turbine optimization
- **Energy storage**: Battery modeling, grid integration, peak shaving
- **Grid integration**: Load balancing, frequency regulation, voltage support
- **Economic analysis**: LCOE calculations, project finance, sensitivity analysis

### Next Steps for Advanced Research

```bash
# Install specialized renewable energy packages
pip3 install pvlib windpowerlib energy-storage-toolkit

# Set up renewable energy modeling software
conda install -c conda-forge homer-energy pysam

# Configure renewable energy databases
aws-research-wizard tools install --domain renewable_energy_systems --advanced
```

### Monthly Cost Estimate

For typical renewable energy research usage:
- **Light usage** (20 hours/week): ~$250/month
- **Medium usage** (35 hours/week): ~$380/month
- **Heavy usage** (50 hours/week): ~$500/month

## Clean Up Resources

**Important**: Always clean up to avoid unexpected charges!

```bash
# Exit your research environment
exit

# Destroy the research environment
aws-research-wizard deploy destroy --domain renewable_energy_systems
```

**Expected result**: "✅ Environment destroyed successfully"

**💰 Billing stops**: No more charges after cleanup

## Troubleshooting

### Common Issues

**Problem**: "pvlib import failed"
**Solution**:
```bash
pip3 install pvlib-python windpowerlib
# OR use conda
conda install -c conda-forge pvlib-python
```

**Problem**: "Optimization solver not found"
**Solution**:
```bash
pip3 install pulp scipy
# For advanced optimization
conda install -c conda-forge gurobi
```

**Problem**: Large energy simulations run slowly
**Solution**:
```bash
# Use more CPU cores
aws-research-wizard deploy create --domain renewable_energy_systems --instance-type c5.4xlarge
```

**Problem**: Memory errors with annual simulations
**Solution**:
```bash
# Use memory-optimized instance
aws-research-wizard deploy create --domain renewable_energy_systems --instance-type r5.2xlarge
```

### Getting Help

- **Renewable Energy Community**: [forum.aws-research-wizard.com/renewable-energy](https://forum.aws-research-wizard.com/renewable-energy)
- **Technical Support**: [support@aws-research-wizard.com](mailto:support@aws-research-wizard.com)
- **Sample Data**: [research-data.aws-wizard.com/renewable-energy](https://research-data.aws-wizard.com/renewable-energy)

### Emergency Stop

If something goes wrong and you want to stop all charges immediately:

```bash
aws-research-wizard emergency-stop --all
```

This will terminate everything and stop billing within 2 minutes.

---

**☀️ Happy renewable energy research!**
*You now have a professional-grade renewable energy systems environment that scales with your clean energy research and development needs.*
