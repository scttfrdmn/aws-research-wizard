# Forestry & Natural Resources Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $9-15 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working forestry and natural resources research environment that can:
- Analyze forest management and biodiversity data
- Process satellite imagery and remote sensing data for forest monitoring
- Model forest growth, carbon sequestration, and ecosystem services
- Handle conservation planning and natural resource assessment

### Meet Dr. James Thompson

Dr. James Thompson is a forest ecologist at USDA Forest Service. He analyzes forest health data but waits weeks for computing resources. Each forest assessment requires processing thousands of satellite images and field measurements.

**Before**: 3-week waits + 1-week analysis = 4 weeks per forest study
**After**: 15-minute setup + 6-hour analysis = same day results
**Time Saved**: 97% faster forest research cycle
**Cost Savings**: $350/month vs $1,400 government allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $9-15 (we'll clean up resources when done)
- **Daily research cost**: $15-35 per day when actively analyzing
- **Monthly estimate**: $180-450 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No forestry or programming experience required

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
- **Region**: Choose `us-west-2` (recommended for forestry with good satellite data access)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain forestry_natural_resources --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**: "✅ All validations passed"

**⚠️ If validation fails**: Check your internet connection and AWS credentials.

## Step 5: Deploy Your Research Environment

```bash
aws-research-wizard deploy create --domain forestry_natural_resources --region us-west-2 --instance-type r5.large
```

**What this does**: Creates a cloud computer with forestry research tools installed.

**Expected result**: You'll see progress updates for about 5 minutes, then "✅ Environment ready"

**💰 Billing starts now**: About $0.13 per hour ($3.12 per day if left running)

**⚠️ If deploy fails**: Run the command again. AWS sometimes has temporary issues.

## Step 6: Connect to Your Environment

```bash
aws-research-wizard connect --domain forestry_natural_resources
```

**What this does**: Opens a connection to your cloud research environment.

**Expected result**: You'll see a terminal prompt like `[forester@ip-10-0-1-123 ~]$`

**🎉 Success**: You're now inside your forestry research environment!

## Step 7: Verify Your Tools

Let's make sure all the forestry tools are working:

```bash
# Check Python geospatial tools
python3 -c "import pandas, numpy, rasterio, geopandas; print('✅ Geospatial tools ready')"

# Check R forestry packages
R --version | head -1

# Check GDAL for satellite data processing
gdal-config --version
```

**Expected result**: You should see "✅" messages confirming tools are installed.

**⚠️ If tools are missing**: Run `sudo yum update && sudo yum install gdal python3-pip R` then try again.

## Your First Forest Analysis

Let's run a real forest analysis to test everything:

### 1. Download Sample Forest Data

```bash
# Create workspace
mkdir -p ~/forest_research/forest_analysis
cd ~/forest_research/forest_analysis

# Download forest inventory data
wget -O forest_inventory.csv "https://research-data.aws-wizard.com/forestry/forest_inventory_sample.csv"

# Download satellite imagery sample
wget -O landsat_forest.tif "https://research-data.aws-wizard.com/forestry/landsat_forest_sample.tif"
```

### 2. Forest Inventory Analysis

Create this Python script for forest analysis:

```bash
cat > forest_analyzer.py << 'EOF'
#!/usr/bin/env python3
"""
Forest & Natural Resources Analysis Suite
Analyzes forest inventory, growth patterns, and ecosystem services
"""

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
from scipy import stats
import warnings
warnings.filterwarnings('ignore')

# Load forest inventory data
print("🌲 Loading forest inventory data...")
data = pd.read_csv('forest_inventory.csv')
print(f"Loaded {len(data)} forest plots")

# Basic forest analysis
print("\n🌳 Forest Stand Analysis")
print("=" * 25)

# Tree species composition
species_composition = data['species'].value_counts()
print("Top 10 Tree Species by Frequency:")
for species, count in species_composition.head(10).items():
    percentage = (count / len(data)) * 100
    print(f"  {species}: {count} plots ({percentage:.1f}%)")

# Forest structure analysis
print(f"\n📊 Forest Structure Metrics")
print("=" * 28)

# Calculate basal area (cross-sectional area of tree stems)
data['basal_area'] = np.pi * (data['dbh_cm'] / 200) ** 2  # DBH in cm to basal area in m²

# Stand-level metrics
structure_stats = {
    'DBH (cm)': data['dbh_cm'].describe(),
    'Height (m)': data['height_m'].describe(),
    'Basal Area (m²)': data['basal_area'].describe(),
    'Crown Width (m)': data['crown_width_m'].describe()
}

for metric, stats_data in structure_stats.items():
    print(f"\n{metric}:")
    print(f"  Mean: {stats_data['mean']:.1f}")
    print(f"  Median: {stats_data['50%']:.1f}")
    print(f"  Range: {stats_data['min']:.1f} - {stats_data['max']:.1f}")

# Age class distribution
data['age_class'] = pd.cut(data['age_years'],
                          bins=[0, 20, 40, 80, 120, 200],
                          labels=['Seedling', 'Sapling', 'Mature', 'Old Growth', 'Ancient'])

age_distribution = data['age_class'].value_counts()
print(f"\nAge Class Distribution:")
for age_class, count in age_distribution.items():
    percentage = (count / len(data)) * 100
    print(f"  {age_class}: {count} trees ({percentage:.1f}%)")

# Carbon sequestration estimation
print(f"\n🌍 Carbon Sequestration Analysis")
print("=" * 30)

# Biomass estimation using allometric equations
def estimate_biomass(dbh, height, species_group):
    """Estimate tree biomass using simplified allometric equations"""
    if species_group in ['Oak', 'Maple', 'Hardwood']:
        # Hardwood equation: Biomass = 0.112 * DBH^2.185 * Height^0.507
        biomass = 0.112 * (dbh ** 2.185) * (height ** 0.507)
    else:
        # Softwood equation: Biomass = 0.095 * DBH^2.129 * Height^0.543
        biomass = 0.095 * (dbh ** 2.129) * (height ** 0.543)
    return biomass

# Group species into hardwood/softwood
hardwood_species = ['Oak', 'Maple', 'Birch', 'Beech', 'Ash']
data['species_group'] = data['species'].apply(
    lambda x: 'Hardwood' if any(hw in x for hw in hardwood_species) else 'Softwood'
)

# Calculate biomass and carbon
data['biomass_kg'] = data.apply(
    lambda row: estimate_biomass(row['dbh_cm'], row['height_m'], row['species_group']), axis=1
)
data['carbon_kg'] = data['biomass_kg'] * 0.47  # Carbon is ~47% of biomass

# Forest-wide carbon analysis
total_carbon = data['carbon_kg'].sum()
carbon_per_hectare = total_carbon / data['plot_area_ha'].sum()

print(f"Total carbon stored: {total_carbon/1000:.1f} tons")
print(f"Carbon density: {carbon_per_hectare:.1f} kg/hectare")

# Carbon by species group
carbon_by_group = data.groupby('species_group')['carbon_kg'].agg(['sum', 'mean']).round(1)
print(f"\nCarbon storage by species group:")
for group, values in carbon_by_group.iterrows():
    print(f"  {group}: Total {values['sum']/1000:.1f} tons, Average {values['mean']:.1f} kg/tree")

# Health and mortality analysis
print(f"\n🏥 Forest Health Assessment")
print("=" * 27)

# Health status distribution
health_status = data['health_status'].value_counts()
print("Tree Health Distribution:")
for status, count in health_status.items():
    percentage = (count / len(data)) * 100
    print(f"  {status}: {count} trees ({percentage:.1f}%)")

# Mortality risk factors
mortality_analysis = data.groupby('health_status').agg({
    'age_years': 'mean',
    'dbh_cm': 'mean',
    'crown_width_m': 'mean'
}).round(1)

print(f"\nHealth status characteristics:")
print(mortality_analysis)

# Generate forest visualization
plt.figure(figsize=(15, 12))

# Species composition pie chart
plt.subplot(3, 3, 1)
top_species = species_composition.head(8)
plt.pie(top_species.values, labels=top_species.index, autopct='%1.1f%%')
plt.title('Species Composition (Top 8)')

# DBH distribution
plt.subplot(3, 3, 2)
plt.hist(data['dbh_cm'], bins=25, alpha=0.7, color='brown')
plt.title('DBH Distribution')
plt.xlabel('DBH (cm)')
plt.ylabel('Frequency')

# Height vs DBH relationship
plt.subplot(3, 3, 3)
plt.scatter(data['dbh_cm'], data['height_m'], alpha=0.6, color='green')
plt.title('Height vs DBH Relationship')
plt.xlabel('DBH (cm)')
plt.ylabel('Height (m)')

# Age class distribution
plt.subplot(3, 3, 4)
age_distribution.plot(kind='bar', color='darkgreen')
plt.title('Age Class Distribution')
plt.ylabel('Number of Trees')
plt.xticks(rotation=45)

# Carbon storage by species group
plt.subplot(3, 3, 5)
carbon_totals = data.groupby('species_group')['carbon_kg'].sum() / 1000
carbon_totals.plot(kind='bar', color=['brown', 'green'])
plt.title('Carbon Storage by Species Group')
plt.ylabel('Carbon (tons)')
plt.xticks(rotation=0)

# Health status distribution
plt.subplot(3, 3, 6)
health_status.plot(kind='bar', color=['red', 'orange', 'yellow', 'green'])
plt.title('Forest Health Status')
plt.ylabel('Number of Trees')
plt.xticks(rotation=45)

# Basal area distribution
plt.subplot(3, 3, 7)
plt.hist(data['basal_area'], bins=20, alpha=0.7, color='purple')
plt.title('Basal Area Distribution')
plt.xlabel('Basal Area (m²)')
plt.ylabel('Frequency')

# Carbon vs Age relationship
plt.subplot(3, 3, 8)
plt.scatter(data['age_years'], data['carbon_kg'], alpha=0.6, color='blue')
plt.title('Carbon Storage vs Tree Age')
plt.xlabel('Age (years)')
plt.ylabel('Carbon (kg)')

# Crown width vs DBH
plt.subplot(3, 3, 9)
plt.scatter(data['dbh_cm'], data['crown_width_m'], alpha=0.6, color='orange')
plt.title('Crown Width vs DBH')
plt.xlabel('DBH (cm)')
plt.ylabel('Crown Width (m)')

plt.tight_layout()
plt.savefig('forest_inventory_analysis.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Forest analysis dashboard saved as 'forest_inventory_analysis.png'")

# Growth prediction model
print(f"\n📈 Forest Growth Modeling")
print("=" * 24)

# Simple growth model based on age and site quality
def predict_dbh_growth(current_dbh, age, site_index):
    """Predict DBH growth using Chapman-Richards model"""
    # Simplified Chapman-Richards growth equation
    max_dbh = site_index * 0.8  # Maximum DBH based on site quality
    growth_rate = 0.05 + (site_index / 100)  # Growth rate modifier

    # Annual growth decreases with age
    annual_growth = growth_rate * (1 - current_dbh / max_dbh) * (1 / (1 + age / 50))
    return max(0, annual_growth)

# Calculate projected growth
data['projected_5yr_growth'] = data.apply(
    lambda row: predict_dbh_growth(row['dbh_cm'], row['age_years'], row['site_index']) * 5,
    axis=1
)

data['projected_dbh_5yr'] = data['dbh_cm'] + data['projected_5yr_growth']

# Growth statistics
avg_growth = data['projected_5yr_growth'].mean()
print(f"Average 5-year DBH growth projection: {avg_growth:.1f} cm")

# Growth by age class
growth_by_age = data.groupby('age_class')['projected_5yr_growth'].mean()
print(f"\nProjected 5-year growth by age class:")
for age_class, growth in growth_by_age.items():
    print(f"  {age_class}: {growth:.1f} cm")

# Stand density analysis
print(f"\n🌲 Stand Density Analysis")
print("=" * 23)

# Calculate trees per hectare
trees_per_ha = len(data) / data['plot_area_ha'].sum()
print(f"Tree density: {trees_per_ha:.0f} trees/hectare")

# Basal area per hectare
basal_area_per_ha = data['basal_area'].sum() / data['plot_area_ha'].sum()
print(f"Basal area: {basal_area_per_ha:.1f} m²/hectare")

# Stand density index (relative measure of crowding)
# Simplified SDI calculation
avg_dbh = data['dbh_cm'].mean()
sdi = trees_per_ha * (avg_dbh / 25) ** 1.6
print(f"Stand Density Index: {sdi:.0f}")

if sdi > 1000:
    density_status = "Overstocked - consider thinning"
elif sdi > 600:
    density_status = "Fully stocked - monitor closely"
elif sdi > 300:
    density_status = "Adequately stocked"
else:
    density_status = "Understocked - natural regeneration or planting needed"

print(f"Density assessment: {density_status}")

print(f"\n✅ Forest inventory analysis complete!")
print(f"Analyzed {len(data)} trees across {data['plot_area_ha'].sum():.1f} hectares")
EOF

chmod +x forest_analyzer.py
```

### 3. Run the Forest Analysis

```bash
python3 forest_analyzer.py
```

**Expected output**: You should see comprehensive forest inventory analysis results.

### 4. Satellite Data Processing Script

```bash
cat > satellite_forest_analyzer.py << 'EOF'
#!/usr/bin/env python3
"""
Satellite Forest Monitoring Tool
Processes remote sensing data for forest change detection
"""

import numpy as np
import matplotlib.pyplot as plt
import warnings
warnings.filterwarnings('ignore')

try:
    import rasterio
    from rasterio.plot import show
    from rasterio.mask import mask
    print("📡 Satellite data processing tools ready")
except ImportError:
    print("⚠️ Installing satellite processing tools...")
    import subprocess
    subprocess.run(['pip3', 'install', 'rasterio', 'geopandas'], check=True)
    import rasterio
    from rasterio.plot import show

# Simulate satellite data processing (since we don't have real Landsat data)
print("\n🛰️ Satellite Forest Monitoring Analysis")
print("=" * 35)

# Generate synthetic NDVI data (Normalized Difference Vegetation Index)
print("Generating synthetic forest NDVI data...")
np.random.seed(42)

# Create a 100x100 pixel forest area
forest_area = np.random.normal(0.7, 0.15, (100, 100))  # Healthy forest NDVI ~0.7
forest_area = np.clip(forest_area, 0, 1)  # NDVI ranges from 0 to 1

# Add some disturbance areas (logging, fire, disease)
# Logging patches (lower NDVI)
forest_area[20:30, 20:40] = np.random.normal(0.3, 0.1, (10, 20))
forest_area[60:75, 70:85] = np.random.normal(0.25, 0.1, (15, 15))

# Disease patches (gradually declining NDVI)
forest_area[40:55, 10:25] = np.random.normal(0.45, 0.1, (15, 15))

# Water bodies (very low NDVI)
forest_area[10:15, 80:95] = np.random.normal(0.1, 0.05, (5, 15))

# Forest health classification
def classify_forest_health(ndvi):
    """Classify forest health based on NDVI values"""
    if ndvi >= 0.6:
        return 'Healthy'
    elif ndvi >= 0.4:
        return 'Stressed'
    elif ndvi >= 0.2:
        return 'Severely Degraded'
    else:
        return 'Non-Forest/Water'

# Apply classification
health_classes = np.vectorize(classify_forest_health)(forest_area)

# Calculate forest health statistics
unique_classes, counts = np.unique(health_classes, return_counts=True)
total_pixels = forest_area.size

print(f"\nForest Health Classification Results:")
print("=" * 35)
for class_name, count in zip(unique_classes, counts):
    percentage = (count / total_pixels) * 100
    print(f"  {class_name}: {count} pixels ({percentage:.1f}%)")

# Change detection simulation
print(f"\n🔍 Forest Change Detection")
print("=" * 25)

# Simulate previous year's data
previous_year = forest_area + np.random.normal(0.05, 0.02, forest_area.shape)
previous_year = np.clip(previous_year, 0, 1)

# Calculate NDVI change
ndvi_change = forest_area - previous_year

# Classify changes
def classify_change(change):
    """Classify forest change based on NDVI difference"""
    if change > 0.1:
        return 'Significant Improvement'
    elif change > 0.05:
        return 'Moderate Improvement'
    elif change > -0.05:
        return 'Stable'
    elif change > -0.15:
        return 'Moderate Decline'
    else:
        return 'Significant Decline'

change_classes = np.vectorize(classify_change)(ndvi_change)

# Change statistics
change_unique, change_counts = np.unique(change_classes, return_counts=True)
print("Forest Change Assessment (Year-over-Year):")
for change_type, count in zip(change_unique, change_counts):
    percentage = (count / total_pixels) * 100
    print(f"  {change_type}: {count} pixels ({percentage:.1f}%)")

# Deforestation hotspot detection
deforestation_threshold = -0.2
deforestation_pixels = np.sum(ndvi_change < deforestation_threshold)
deforestation_area_ha = deforestation_pixels * 0.09  # Assuming 30m pixels = 0.09 ha

print(f"\nDeforestation Alert:")
print(f"  Pixels with severe NDVI decline: {deforestation_pixels}")
print(f"  Estimated deforested area: {deforestation_area_ha:.1f} hectares")

if deforestation_area_ha > 10:
    print(f"  ⚠️ WARNING: Significant deforestation detected!")
else:
    print(f"  ✅ Deforestation within normal limits")

# Forest fragmentation analysis
print(f"\n🧩 Forest Fragmentation Analysis")
print("=" * 30)

# Calculate forest patches (simplified connectivity analysis)
healthy_forest = (forest_area >= 0.6).astype(int)

# Calculate edge-to-interior ratio
from scipy import ndimage

# Find forest edges
forest_edges = healthy_forest - ndimage.binary_erosion(healthy_forest)
edge_pixels = np.sum(forest_edges)
interior_pixels = np.sum(healthy_forest) - edge_pixels

edge_to_interior_ratio = edge_pixels / max(interior_pixels, 1)

print(f"Forest edge pixels: {edge_pixels}")
print(f"Forest interior pixels: {interior_pixels}")
print(f"Edge-to-interior ratio: {edge_to_interior_ratio:.2f}")

if edge_to_interior_ratio > 0.5:
    fragmentation_status = "Highly fragmented"
elif edge_to_interior_ratio > 0.3:
    fragmentation_status = "Moderately fragmented"
else:
    fragmentation_status = "Well-connected"

print(f"Fragmentation assessment: {fragmentation_status}")

# Generate comprehensive satellite analysis visualization
plt.figure(figsize=(16, 12))

# Current NDVI
plt.subplot(3, 3, 1)
ndvi_plot = plt.imshow(forest_area, cmap='RdYlGn', vmin=0, vmax=1)
plt.colorbar(ndvi_plot, label='NDVI')
plt.title('Current Forest NDVI')
plt.axis('off')

# Forest health classification
plt.subplot(3, 3, 2)
health_colors = {'Healthy': 0, 'Stressed': 1, 'Severely Degraded': 2, 'Non-Forest/Water': 3}
health_numeric = np.vectorize(health_colors.get)(health_classes)
health_plot = plt.imshow(health_numeric, cmap='RdYlGn_r')
plt.colorbar(health_plot, ticks=[0,1,2,3], label='Health Class')
plt.title('Forest Health Classification')
plt.axis('off')

# NDVI change detection
plt.subplot(3, 3, 3)
change_plot = plt.imshow(ndvi_change, cmap='RdBu', vmin=-0.3, vmax=0.3)
plt.colorbar(change_plot, label='NDVI Change')
plt.title('Year-over-Year NDVI Change')
plt.axis('off')

# Forest health distribution
plt.subplot(3, 3, 4)
plt.pie(counts, labels=unique_classes, autopct='%1.1f%%',
        colors=['green', 'orange', 'red', 'blue'])
plt.title('Forest Health Distribution')

# Change distribution
plt.subplot(3, 3, 5)
plt.bar(range(len(change_unique)), change_counts, color=['darkred', 'red', 'gray', 'lightgreen', 'green'])
plt.xticks(range(len(change_unique)), change_unique, rotation=45)
plt.title('Change Type Distribution')
plt.ylabel('Pixel Count')

# NDVI histogram
plt.subplot(3, 3, 6)
plt.hist(forest_area.flatten(), bins=30, alpha=0.7, color='green')
plt.axvline(x=0.6, color='red', linestyle='--', label='Health Threshold')
plt.title('NDVI Distribution')
plt.xlabel('NDVI Value')
plt.ylabel('Frequency')
plt.legend()

# Deforestation hotspots
plt.subplot(3, 3, 7)
deforestation_mask = ndvi_change < deforestation_threshold
plt.imshow(deforestation_mask, cmap='Reds')
plt.title('Deforestation Hotspots')
plt.axis('off')

# Forest edges
plt.subplot(3, 3, 8)
plt.imshow(forest_edges, cmap='Blues')
plt.title('Forest Edge Detection')
plt.axis('off')

# Healthy forest connectivity
plt.subplot(3, 3, 9)
plt.imshow(healthy_forest, cmap='Greens')
plt.title('Healthy Forest Connectivity')
plt.axis('off')

plt.tight_layout()
plt.savefig('satellite_forest_analysis.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Satellite analysis dashboard saved as 'satellite_forest_analysis.png'")

# Carbon monitoring from satellite data
print(f"\n🌍 Satellite-Based Carbon Monitoring")
print("=" * 33)

# Estimate biomass from NDVI (simplified relationship)
# Higher NDVI generally correlates with higher biomass
biomass_per_pixel = (forest_area * 150) ** 1.2  # Simplified allometric relationship
carbon_per_pixel = biomass_per_pixel * 0.47  # 47% carbon content

total_carbon_tons = np.sum(carbon_per_pixel) * 0.09 / 1000  # Convert to tons, account for pixel area
avg_carbon_density = np.mean(carbon_per_pixel) * 0.09  # kg/ha

print(f"Estimated total carbon storage: {total_carbon_tons:.1f} tons")
print(f"Average carbon density: {avg_carbon_density:.1f} kg/hectare")

# Carbon change assessment
previous_biomass = (previous_year * 150) ** 1.2
previous_carbon = previous_biomass * 0.47
carbon_change = np.sum((carbon_per_pixel - previous_carbon) * 0.09) / 1000

print(f"Annual carbon change: {carbon_change:+.1f} tons")

if carbon_change > 0:
    print(f"✅ Forest is a net carbon sink (+{carbon_change:.1f} tons/year)")
else:
    print(f"⚠️ Forest is a net carbon source ({carbon_change:.1f} tons/year)")

print(f"\n✅ Satellite forest monitoring complete!")
print(f"Analyzed {total_pixels} pixels covering {total_pixels * 0.09:.1f} hectares")
EOF

chmod +x satellite_forest_analyzer.py
```

### 5. Run Satellite Analysis

```bash
python3 satellite_forest_analyzer.py
```

**Expected output**: Satellite-based forest monitoring with change detection analysis.

### 6. Biodiversity Assessment Script

```bash
cat > biodiversity_analyzer.py << 'EOF'
#!/usr/bin/env python3
"""
Forest Biodiversity Assessment Tool
Analyzes species diversity and ecosystem health indicators
"""

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from scipy import stats
from math import log
import warnings
warnings.filterwarnings('ignore')

# Generate synthetic biodiversity data
print("🦋 Generating forest biodiversity data...")
np.random.seed(42)

# Create sample biodiversity survey data
n_plots = 50
n_species = 85

# Generate species abundance data
species_names = [
    # Trees
    'White Oak', 'Red Maple', 'Eastern Hemlock', 'Yellow Birch', 'American Beech',
    'Sugar Maple', 'Black Cherry', 'White Pine', 'Basswood', 'Red Oak',
    # Understory plants
    'Trout Lily', 'Wild Ginger', 'Bloodroot', 'Mayapple', 'Jack-in-pulpit',
    'Wild Leek', 'Hepatica', 'Spring Beauty', 'Dutchman Breeches', 'Violet',
    # Ferns
    'Christmas Fern', 'Royal Fern', 'Cinnamon Fern', 'Bracken Fern', 'Maidenhair Fern',
    # Birds
    'Wood Thrush', 'Scarlet Tanager', 'Pileated Woodpecker', 'Ovenbird', 'Hermit Thrush',
    'Black-throated Blue Warbler', 'Veery', 'Red-eyed Vireo', 'White-breasted Nuthatch',
    # Mammals
    'White-tailed Deer', 'Black Bear', 'Gray Squirrel', 'Chipmunk', 'Red Squirrel',
    'Flying Squirrel', 'Raccoon', 'Opossum', 'Porcupine', 'Bobcat',
    # Amphibians
    'Red-backed Salamander', 'Spotted Salamander', 'Wood Frog', 'Spring Peeper',
    'American Toad', 'Two-lined Salamander', 'Four-toed Salamander',
    # Invertebrates
    'Forest Wolf Spider', 'Millipede', 'Centipede', 'Ground Beetle', 'Carpenter Ant',
    'Pill Bug', 'Earthworm', 'Snail', 'Slug', 'Harvestman'
] + [f'Species_{i}' for i in range(len([
    'White Oak', 'Red Maple', 'Eastern Hemlock', 'Yellow Birch', 'American Beech',
    'Sugar Maple', 'Black Cherry', 'White Pine', 'Basswood', 'Red Oak',
    'Trout Lily', 'Wild Ginger', 'Bloodroot', 'Mayapple', 'Jack-in-pulpit',
    'Wild Leek', 'Hepatica', 'Spring Beauty', 'Dutchman Breeches', 'Violet',
    'Christmas Fern', 'Royal Fern', 'Cinnamon Fern', 'Bracken Fern', 'Maidenhair Fern',
    'Wood Thrush', 'Scarlet Tanager', 'Pileated Woodpecker', 'Ovenbird', 'Hermit Thrush',
    'Black-throated Blue Warbler', 'Veery', 'Red-eyed Vireo', 'White-breasted Nuthatch',
    'White-tailed Deer', 'Black Bear', 'Gray Squirrel', 'Chipmunk', 'Red Squirrel',
    'Flying Squirrel', 'Raccoon', 'Opossum', 'Porcupine', 'Bobcat',
    'Red-backed Salamander', 'Spotted Salamander', 'Wood Frog', 'Spring Peeper',
    'American Toad', 'Two-lined Salamander', 'Four-toed Salamander',
    'Forest Wolf Spider', 'Millipede', 'Centipede', 'Ground Beetle', 'Carpenter Ant',
    'Pill Bug', 'Earthworm', 'Snail', 'Slug', 'Harvestman'
]), n_species)]

# Create taxonomic groups
taxonomic_groups = {
    'Trees': species_names[:10],
    'Understory Plants': species_names[10:20],
    'Ferns': species_names[20:25],
    'Birds': species_names[25:34],
    'Mammals': species_names[34:44],
    'Amphibians': species_names[44:51],
    'Invertebrates': species_names[51:61],
    'Other': species_names[61:]
}

# Generate abundance data with realistic patterns
biodiversity_data = []

for plot_id in range(1, n_plots + 1):
    # Simulate habitat quality effect on species richness
    habitat_quality = np.random.uniform(0.3, 1.0)

    # Number of species in this plot (influenced by habitat quality)
    n_species_plot = int(np.random.poisson(20 * habitat_quality))

    # Select random species for this plot
    plot_species = np.random.choice(species_names,
                                   size=min(n_species_plot, len(species_names)),
                                   replace=False)

    for species in plot_species:
        # Abundance follows log-normal distribution (common in ecology)
        abundance = max(1, int(np.random.lognormal(2, 1) * habitat_quality))

        # Determine taxonomic group
        taxonomic_group = 'Other'
        for group, group_species in taxonomic_groups.items():
            if species in group_species:
                taxonomic_group = group
                break

        biodiversity_data.append({
            'plot_id': plot_id,
            'species': species,
            'abundance': abundance,
            'taxonomic_group': taxonomic_group,
            'habitat_quality': round(habitat_quality, 2)
        })

df = pd.DataFrame(biodiversity_data)
print(f"Generated biodiversity data for {len(df)} species observations across {n_plots} plots")

# Biodiversity analysis
print(f"\n🌿 Forest Biodiversity Assessment")
print("=" * 30)

# Species richness analysis
plot_richness = df.groupby('plot_id').agg({
    'species': 'nunique',
    'abundance': 'sum',
    'habitat_quality': 'first'
}).rename(columns={'species': 'species_richness', 'abundance': 'total_abundance'})

print(f"Species Richness Statistics:")
print(f"  Mean species per plot: {plot_richness['species_richness'].mean():.1f}")
print(f"  Range: {plot_richness['species_richness'].min()} - {plot_richness['species_richness'].max()} species")
print(f"  Total species across all plots: {df['species'].nunique()}")

# Taxonomic diversity
taxonomic_diversity = df.groupby('taxonomic_group').agg({
    'species': 'nunique',
    'abundance': 'sum'
}).rename(columns={'species': 'species_count', 'abundance': 'total_abundance'})

print(f"\nTaxonomic Group Diversity:")
for group, data in taxonomic_diversity.iterrows():
    print(f"  {group}: {data['species_count']} species, {data['total_abundance']} individuals")

# Calculate diversity indices
def calculate_shannon_diversity(abundances):
    """Calculate Shannon diversity index"""
    total = sum(abundances)
    if total == 0:
        return 0
    proportions = [n/total for n in abundances if n > 0]
    return -sum(p * log(p) for p in proportions)

def calculate_simpson_diversity(abundances):
    """Calculate Simpson diversity index"""
    total = sum(abundances)
    if total == 0:
        return 0
    return 1 - sum((n/total)**2 for n in abundances)

# Calculate diversity indices for each plot
diversity_indices = []
for plot_id in range(1, n_plots + 1):
    plot_data = df[df['plot_id'] == plot_id]
    abundances = plot_data['abundance'].tolist()

    shannon = calculate_shannon_diversity(abundances)
    simpson = calculate_simpson_diversity(abundances)

    diversity_indices.append({
        'plot_id': plot_id,
        'shannon_diversity': shannon,
        'simpson_diversity': simpson,
        'species_richness': len(abundances),
        'habitat_quality': plot_data['habitat_quality'].iloc[0] if len(plot_data) > 0 else 0
    })

diversity_df = pd.DataFrame(diversity_indices)

print(f"\n📊 Diversity Indices Summary")
print("=" * 26)
print(f"Shannon Diversity:")
print(f"  Mean: {diversity_df['shannon_diversity'].mean():.2f}")
print(f"  Range: {diversity_df['shannon_diversity'].min():.2f} - {diversity_df['shannon_diversity'].max():.2f}")

print(f"Simpson Diversity:")
print(f"  Mean: {diversity_df['simpson_diversity'].mean():.2f}")
print(f"  Range: {diversity_df['simpson_diversity'].min():.2f} - {diversity_df['simpson_diversity'].max():.2f}")

# Species abundance patterns
print(f"\n🦎 Species Abundance Patterns")
print("=" * 27)

species_totals = df.groupby('species')['abundance'].sum().sort_values(ascending=False)
print(f"Most abundant species:")
for species, abundance in species_totals.head(10).items():
    print(f"  {species}: {abundance} individuals")

# Rare species analysis
rare_species = species_totals[species_totals <= 5]
print(f"\nRare species (≤5 individuals): {len(rare_species)} species")
print(f"Percentage of rare species: {len(rare_species)/len(species_totals)*100:.1f}%")

# Habitat quality correlation
print(f"\n🏞️ Habitat Quality Relationships")
print("=" * 30)

# Correlation between habitat quality and diversity
habitat_shannon_corr = stats.pearsonr(diversity_df['habitat_quality'],
                                     diversity_df['shannon_diversity'])
habitat_richness_corr = stats.pearsonr(diversity_df['habitat_quality'],
                                      diversity_df['species_richness'])

print(f"Habitat quality vs Shannon diversity: r = {habitat_shannon_corr[0]:.3f}, p = {habitat_shannon_corr[1]:.3f}")
print(f"Habitat quality vs species richness: r = {habitat_richness_corr[0]:.3f}, p = {habitat_richness_corr[1]:.3f}")

# Conservation priority assessment
print(f"\n🛡️ Conservation Priority Assessment")
print("=" * 32)

# Calculate conservation scores based on multiple factors
diversity_df['conservation_score'] = (
    diversity_df['shannon_diversity'] * 0.4 +  # Diversity weight
    diversity_df['species_richness'] / diversity_df['species_richness'].max() * 0.3 +  # Richness weight
    diversity_df['habitat_quality'] * 0.3  # Habitat quality weight
)

# Classify conservation priority
diversity_df['conservation_priority'] = pd.cut(
    diversity_df['conservation_score'],
    bins=[0, 0.4, 0.7, 1.0],
    labels=['Low', 'Medium', 'High']
)

priority_counts = diversity_df['conservation_priority'].value_counts()
print(f"Conservation Priority Distribution:")
for priority, count in priority_counts.items():
    percentage = (count / len(diversity_df)) * 100
    print(f"  {priority} Priority: {count} plots ({percentage:.1f}%)")

# Generate biodiversity visualization
plt.figure(figsize=(16, 12))

# Species richness distribution
plt.subplot(3, 3, 1)
plt.hist(diversity_df['species_richness'], bins=15, alpha=0.7, color='green')
plt.title('Species Richness Distribution')
plt.xlabel('Number of Species')
plt.ylabel('Number of Plots')

# Shannon diversity vs habitat quality
plt.subplot(3, 3, 2)
plt.scatter(diversity_df['habitat_quality'], diversity_df['shannon_diversity'],
           alpha=0.6, color='blue')
plt.title('Shannon Diversity vs Habitat Quality')
plt.xlabel('Habitat Quality')
plt.ylabel('Shannon Diversity Index')

# Taxonomic group composition
plt.subplot(3, 3, 3)
group_counts = taxonomic_diversity['species_count']
plt.pie(group_counts.values, labels=group_counts.index, autopct='%1.1f%%')
plt.title('Taxonomic Group Composition')

# Species abundance rank curve
plt.subplot(3, 3, 4)
sorted_abundances = species_totals.sort_values(ascending=False)
plt.loglog(range(1, len(sorted_abundances) + 1), sorted_abundances.values, 'o-')
plt.title('Species Abundance Rank Curve')
plt.xlabel('Species Rank')
plt.ylabel('Abundance (log scale)')

# Conservation priority map
plt.subplot(3, 3, 5)
priority_colors = {'Low': 0, 'Medium': 1, 'High': 2}
priority_numeric = diversity_df['conservation_priority'].map(priority_colors)
plt.scatter(diversity_df['plot_id'], priority_numeric,
           c=priority_numeric, cmap='RdYlGn', s=100)
plt.title('Conservation Priority by Plot')
plt.xlabel('Plot ID')
plt.ylabel('Priority Level')

# Diversity indices comparison
plt.subplot(3, 3, 6)
plt.scatter(diversity_df['shannon_diversity'], diversity_df['simpson_diversity'],
           alpha=0.6, color='purple')
plt.title('Shannon vs Simpson Diversity')
plt.xlabel('Shannon Diversity')
plt.ylabel('Simpson Diversity')

# Species accumulation by taxonomic group
plt.subplot(3, 3, 7)
group_abundances = df.groupby('taxonomic_group')['abundance'].sum().sort_values(ascending=True)
group_abundances.plot(kind='barh', color='orange')
plt.title('Total Abundance by Group')
plt.xlabel('Total Abundance')

# Habitat quality distribution
plt.subplot(3, 3, 8)
plt.hist(diversity_df['habitat_quality'], bins=10, alpha=0.7, color='brown')
plt.title('Habitat Quality Distribution')
plt.xlabel('Habitat Quality Score')
plt.ylabel('Number of Plots')

# Rare species distribution
plt.subplot(3, 3, 9)
rarity_threshold = [1, 2, 5, 10, 20]
rarity_counts = [sum(species_totals <= t) for t in rarity_threshold]
plt.bar(range(len(rarity_threshold)), rarity_counts, color='red', alpha=0.7)
plt.xticks(range(len(rarity_threshold)), [f'≤{t}' for t in rarity_threshold])
plt.title('Rare Species Distribution')
plt.xlabel('Abundance Threshold')
plt.ylabel('Number of Species')

plt.tight_layout()
plt.savefig('biodiversity_assessment.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Biodiversity assessment dashboard saved as 'biodiversity_assessment.png'")

# Ecosystem health indicators
print(f"\n🌳 Ecosystem Health Indicators")
print("=" * 28)

# Calculate ecosystem health metrics
ecosystem_health = {
    'Species Richness': diversity_df['species_richness'].mean(),
    'Shannon Diversity': diversity_df['shannon_diversity'].mean(),
    'Habitat Quality': diversity_df['habitat_quality'].mean(),
    'Conservation Score': diversity_df['conservation_score'].mean()
}

print(f"Ecosystem Health Metrics:")
for metric, value in ecosystem_health.items():
    if value > 0.7:
        status = "Excellent"
    elif value > 0.5:
        status = "Good"
    elif value > 0.3:
        status = "Fair"
    else:
        status = "Poor"
    print(f"  {metric}: {value:.2f} ({status})")

# Threat assessment
print(f"\nThreat Assessment:")
high_priority_plots = len(diversity_df[diversity_df['conservation_priority'] == 'High'])
threatened_habitats = len(diversity_df[diversity_df['habitat_quality'] < 0.5])

print(f"  High conservation priority plots: {high_priority_plots}/{len(diversity_df)} ({high_priority_plots/len(diversity_df)*100:.1f}%)")
print(f"  Threatened habitats (quality < 0.5): {threatened_habitats}/{len(diversity_df)} ({threatened_habitats/len(diversity_df)*100:.1f}%)")

if threatened_habitats > len(diversity_df) * 0.2:
    print(f"  ⚠️ WARNING: >20% of habitats are threatened!")
else:
    print(f"  ✅ Habitat degradation within acceptable limits")

print(f"\n✅ Biodiversity assessment complete!")
print(f"Analyzed {df['species'].nunique()} species across {n_plots} forest plots")
EOF

chmod +x biodiversity_analyzer.py
```

### 7. Run Biodiversity Assessment

```bash
python3 biodiversity_analyzer.py
```

**Expected output**: Comprehensive biodiversity assessment with conservation priorities.

## What You've Accomplished

🎉 **Congratulations!** You've successfully:

1. ✅ Created a forestry research environment in the cloud
2. ✅ Analyzed forest inventory data and stand structure
3. ✅ Processed satellite imagery for forest monitoring
4. ✅ Conducted biodiversity assessment and conservation planning
5. ✅ Generated comprehensive forest management reports

### Real Research Applications

Your environment can now handle:
- **Forest inventory**: Stand structure, growth modeling, yield tables
- **Remote sensing**: Satellite change detection, NDVI monitoring
- **Biodiversity surveys**: Species richness, diversity indices
- **Carbon accounting**: Biomass estimation, sequestration rates
- **Conservation planning**: Priority area identification, threat assessment

### Next Steps for Advanced Research

```bash
# Install specialized forestry packages
pip3 install rasterio geopandas forestry-toolkit lidar-processing

# Set up forest databases
wget https://fia.fs.usda.gov/fia-datamart/datasets/

# Configure forest modeling tools
aws-research-wizard tools install --domain forestry_natural_resources --advanced
```

### Monthly Cost Estimate

For typical forestry research usage:
- **Light usage** (15 hours/week): ~$180/month
- **Medium usage** (25 hours/week): ~$290/month
- **Heavy usage** (40 hours/week): ~$450/month

## Clean Up Resources

**Important**: Always clean up to avoid unexpected charges!

```bash
# Exit your research environment
exit

# Destroy the research environment
aws-research-wizard deploy destroy --domain forestry_natural_resources
```

**Expected result**: "✅ Environment destroyed successfully"

**💰 Billing stops**: No more charges after cleanup

## Troubleshooting

### Common Issues

**Problem**: "GDAL not found" errors
**Solution**:
```bash
sudo yum install gdal gdal-devel
pip3 install rasterio geopandas
```

**Problem**: "Rasterio import failed"
**Solution**:
```bash
pip3 install --upgrade rasterio
export GDAL_DATA=/usr/share/gdal
```

**Problem**: Large satellite data processing is slow
**Solution**:
```bash
# Use larger instance type for heavy processing
aws-research-wizard deploy create --domain forestry_natural_resources --instance-type r5.xlarge
```

**Problem**: Memory errors with large datasets
**Solution**:
```bash
# Process data in chunks
python3 -c "import gc; gc.collect()"  # Force garbage collection
```

### Getting Help

- **Forestry Community**: [forum.aws-research-wizard.com/forestry](https://forum.aws-research-wizard.com/forestry)
- **Technical Support**: [support@aws-research-wizard.com](mailto:support@aws-research-wizard.com)
- **Sample Data**: [research-data.aws-wizard.com/forestry](https://research-data.aws-wizard.com/forestry)

### Emergency Stop

If something goes wrong and you want to stop all charges immediately:

```bash
aws-research-wizard emergency-stop --all
```

This will terminate everything and stop billing within 2 minutes.

---

**🌲 Happy forest research!**
*You now have a professional-grade forestry research environment that scales with your needs.*
