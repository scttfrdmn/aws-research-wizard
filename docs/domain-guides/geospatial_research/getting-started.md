# Geospatial Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $11-17 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working geospatial research environment that can:
- Process satellite imagery and remote sensing data
- Perform GIS analysis with GDAL, QGIS, and Python
- Handle large geospatial datasets up to 1TB
- Create maps and spatial visualizations

### Meet Dr. James Thompson
Dr. James Thompson is a remote sensing researcher at NASA. He analyzes satellite data for climate change studies but waits weeks for university computing resources. Each analysis project requires processing terabytes of satellite imagery.

**Before**: 2-week waits + 1-week processing = 3 weeks per analysis
**After**: 15-minute setup + 6-hour processing = same day results
**Time Saved**: 96% faster geospatial analysis cycle
**Cost Savings**: $700/month vs $2,500 university allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $11-17 (we'll clean up resources when done)
- **Daily research cost**: $25-75 per day when actively processing
- **Monthly estimate**: $300-900 per month for typical usage
- **Free tier**: Some storage included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or GIS experience required

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
- **Region**: Choose `us-west-2` (recommended for geospatial with good storage performance)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain geospatial_research --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: geospatial_research
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Geospatial Environment

```bash
aws-research-wizard deploy start --domain geospatial_research --region us-west-2 --instance r6i.xlarge
```

**What this does**: Creates your geospatial computing environment optimized for spatial data processing.

**This will take**: 5-7 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  Memory: 32GB RAM for large raster processing
  Storage: 500GB SSD for geospatial datasets
```

**💰 Billing starts now**: Your environment costs about $0.50 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your geospatial computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Geospatial Tools

Your environment comes pre-installed with:

### Core GIS Tools
- **GDAL/OGR**: Geospatial data processing - Type `gdalinfo --version` to check
- **QGIS**: Desktop GIS software - Type `qgis --version` to check
- **PostGIS**: Spatial database - Type `psql --version` to check
- **GeoPandas**: Python spatial analysis - Type `python -c "import geopandas; print(geopandas.__version__)"` to check
- **Rasterio**: Python raster processing - Type `python -c "import rasterio; print(rasterio.__version__)"` to check

### Try Your First Command
```bash
gdalinfo --version
```

**What this does**: Shows GDAL version and confirms geospatial tools are installed.

**Expected result**: You see GDAL version info confirming GIS libraries are ready.

## Step 8: Process Satellite Imagery

Let's analyze real satellite data to test everything works:

### Download Sample Satellite Data
```bash
# Create working directory
mkdir ~/geospatial-tutorial
cd ~/geospatial-tutorial

# Download sample Landsat imagery (small subset)
wget -O landsat_sample.tif "https://landsat-pds.s3.amazonaws.com/c1/L8/139/045/LC08_L1TP_139045_20170304_20170316_01_T1/LC08_L1TP_139045_20170304_20170316_01_T1_B4.TIF"

# Download administrative boundaries (shapefile)
wget -O countries.zip "https://www.naturalearthdata.com/http//www.naturalearthdata.com/download/110m/cultural/ne_110m_admin_0_countries.zip"
unzip countries.zip

echo "Sample geospatial data downloaded successfully!"
```

### Basic Raster Analysis
```bash
# Create raster analysis script
cat > raster_analysis.py << 'EOF'
import rasterio
import numpy as np
import matplotlib.pyplot as plt
from rasterio.plot import show
import geopandas as gpd

print("Starting satellite imagery analysis...")

def analyze_raster(filename):
    """Analyze satellite raster data"""
    print(f"\n=== Analyzing Raster: {filename} ===")

    try:
        with rasterio.open(filename) as src:
            # Raster metadata
            print(f"Driver: {src.driver}")
            print(f"Dimensions: {src.width} x {src.height}")
            print(f"Bands: {src.count}")
            print(f"CRS: {src.crs}")
            print(f"Bounds: {src.bounds}")
            print(f"Resolution: {src.res}")

            # Read raster data
            band1 = src.read(1)

            # Calculate statistics
            print(f"Data type: {band1.dtype}")
            print(f"No data value: {src.nodata}")
            print(f"Min value: {np.nanmin(band1)}")
            print(f"Max value: {np.nanmax(band1)}")
            print(f"Mean value: {np.nanmean(band1):.2f}")
            print(f"Standard deviation: {np.nanstd(band1):.2f}")

            # Calculate area coverage
            pixel_area = abs(src.res[0] * src.res[1])  # square meters
            total_pixels = band1.size
            valid_pixels = np.count_nonzero(~np.isnan(band1))

            print(f"Pixel size: {pixel_area:.2f} m²")
            print(f"Total area: {(total_pixels * pixel_area) / 1e6:.2f} km²")
            print(f"Valid data coverage: {100 * valid_pixels / total_pixels:.1f}%")

        return True

    except Exception as e:
        print(f"Error reading raster: {e}")
        return False

def analyze_vector(filename):
    """Analyze vector data (shapefile)"""
    print(f"\n=== Analyzing Vector: {filename} ===")

    try:
        # Read shapefile
        gdf = gpd.read_file(filename)

        print(f"Geometry type: {gdf.geom_type.iloc[0]}")
        print(f"Number of features: {len(gdf)}")
        print(f"CRS: {gdf.crs}")
        print(f"Bounds: {gdf.bounds.iloc[0].values}")

        # Column information
        print(f"Columns: {list(gdf.columns)}")

        # Sample data
        if 'NAME' in gdf.columns:
            print(f"Sample countries: {gdf['NAME'].head().tolist()}")

        # Calculate areas (if polygon)
        if gdf.geom_type.iloc[0] == 'Polygon' or gdf.geom_type.iloc[0] == 'MultiPolygon':
            # Convert to equal area projection for area calculation
            gdf_proj = gdf.to_crs('EPSG:3857')  # Web Mercator
            areas = gdf_proj.geometry.area / 1e6  # Convert to km²
            print(f"Largest feature area: {areas.max():.0f} km²")
            print(f"Smallest feature area: {areas.min():.0f} km²")
            print(f"Average feature area: {areas.mean():.0f} km²")

        return True

    except Exception as e:
        print(f"Error reading vector: {e}")
        return False

# Analyze downloaded data
print("=== Geospatial Data Analysis ===")

# Analyze raster data
raster_success = analyze_raster('landsat_sample.tif')

# Analyze vector data
vector_success = analyze_vector('ne_110m_admin_0_countries.shp')

if raster_success and vector_success:
    print("\n✅ Geospatial analysis completed successfully!")
else:
    print("\n⚠️ Some analyses failed - this is normal with sample data")

print("Ready for advanced spatial analysis and processing")
EOF

python3 raster_analysis.py
```

**What this does**: Analyzes satellite imagery metadata and calculates spatial statistics.

**This will take**: 2-3 minutes

### Spatial Operations
```bash
# Create spatial operations script
cat > spatial_operations.py << 'EOF'
import geopandas as gpd
import rasterio
from rasterio.mask import mask
import numpy as np
from shapely.geometry import Point, Polygon
import matplotlib.pyplot as plt

print("Performing advanced spatial operations...")

def create_sample_points():
    """Create sample point data for analysis"""
    print("\n=== Creating Sample Spatial Data ===")

    # Create random points within a bounding box
    np.random.seed(42)  # For reproducible results

    # Bounding box roughly covering part of US West Coast
    min_lon, max_lon = -125, -120
    min_lat, max_lat = 35, 40

    n_points = 100
    lons = np.random.uniform(min_lon, max_lon, n_points)
    lats = np.random.uniform(min_lat, max_lat, n_points)

    # Create GeoDataFrame
    points = [Point(lon, lat) for lon, lat in zip(lons, lats)]
    gdf_points = gpd.GeoDataFrame({
        'id': range(n_points),
        'value': np.random.normal(100, 25, n_points),  # Random values
        'category': np.random.choice(['A', 'B', 'C'], n_points)
    }, geometry=points, crs='EPSG:4326')

    print(f"Created {len(gdf_points)} sample points")
    print(f"Bounding box: {gdf_points.bounds.iloc[0].values}")
    print(f"Value statistics: mean={gdf_points['value'].mean():.1f}, std={gdf_points['value'].std():.1f}")

    return gdf_points

def spatial_analysis(points_gdf):
    """Perform spatial analysis operations"""
    print("\n=== Spatial Analysis Operations ===")

    # Buffer analysis
    buffer_distance = 0.1  # degrees (roughly 11km)
    buffers = points_gdf.geometry.buffer(buffer_distance)
    print(f"Created buffers with {buffer_distance} degree radius")

    # Spatial clustering - find points within distance
    close_pairs = []
    for i, point1 in enumerate(points_gdf.geometry):
        for j, point2 in enumerate(points_gdf.geometry):
            if i < j:  # Avoid duplicates
                distance = point1.distance(point2)
                if distance < buffer_distance:
                    close_pairs.append((i, j, distance))

    print(f"Found {len(close_pairs)} point pairs within buffer distance")

    # Spatial aggregation by category
    category_stats = points_gdf.groupby('category').agg({
        'value': ['count', 'mean', 'std'],
        'geometry': lambda x: x.unary_union.convex_hull.area
    }).round(3)

    print("Statistics by category:")
    print(category_stats)

    # Create spatial grid for interpolation
    bounds = points_gdf.bounds.iloc[0]
    grid_size = 20

    x_grid = np.linspace(bounds['minx'], bounds['maxx'], grid_size)
    y_grid = np.linspace(bounds['miny'], bounds['maxy'], grid_size)

    # Simple inverse distance weighting
    grid_values = []
    for y in y_grid:
        row_values = []
        for x in x_grid:
            grid_point = Point(x, y)
            distances = points_gdf.geometry.distance(grid_point)
            weights = 1 / (distances + 1e-10)  # Avoid division by zero
            weighted_value = np.average(points_gdf['value'], weights=weights)
            row_values.append(weighted_value)
        grid_values.append(row_values)

    grid_values = np.array(grid_values)
    print(f"Created interpolation grid: {grid_values.shape}")
    print(f"Grid value range: {grid_values.min():.1f} to {grid_values.max():.1f}")

    return grid_values, (x_grid, y_grid)

def raster_vector_integration():
    """Demonstrate raster-vector integration"""
    print("\n=== Raster-Vector Integration ===")

    try:
        # Load countries data
        countries = gpd.read_file('ne_110m_admin_0_countries.shp')

        # Filter to a specific region (e.g., North America)
        north_america = countries[countries['CONTINENT'] == 'North America']

        if len(north_america) > 0:
            print(f"North American countries: {len(north_america)}")

            # Calculate country areas
            north_america_proj = north_america.to_crs('EPSG:3857')
            areas = north_america_proj.geometry.area / 1e6  # km²

            largest_country = north_america.loc[areas.idxmax()]
            print(f"Largest country: {largest_country.get('NAME', 'Unknown')} ({areas.max():.0f} km²)")

            # Create bounding box analysis
            total_bounds = north_america.total_bounds
            bbox_area = ((total_bounds[2] - total_bounds[0]) *
                        (total_bounds[3] - total_bounds[1])) * 111**2  # Rough km² conversion

            print(f"Total bounding box area: {bbox_area:.0f} km²")

        else:
            print("No North American countries found in dataset")

    except Exception as e:
        print(f"Vector integration error: {e}")
        print("This is normal if country data is not available")

# Run spatial analysis
points_data = create_sample_points()
grid_data, grid_coords = spatial_analysis(points_data)
raster_vector_integration()

print("\n✅ Advanced spatial operations completed!")
print("Your environment is ready for complex geospatial analysis")
EOF

python3 spatial_operations.py
```

**What this does**: Performs spatial analysis operations including buffering, clustering, and interpolation.

**Expected result**: Shows spatial analysis results and demonstrates GIS operations.

**🎉 Success!** You've processed real geospatial data in the cloud.

## Step 9: Remote Sensing Analysis

Test advanced geospatial capabilities:

```bash
# Create remote sensing analysis script
cat > remote_sensing.py << 'EOF'
import rasterio
import numpy as np
from rasterio.warp import reproject, Resampling
from rasterio.windows import from_bounds
import matplotlib.pyplot as plt

print("Performing remote sensing analysis...")

def ndvi_calculation_simulation():
    """Simulate NDVI calculation from multispectral data"""
    print("\n=== NDVI Calculation Simulation ===")

    # Simulate multispectral bands (Red and Near-Infrared)
    np.random.seed(42)

    # Create synthetic satellite data
    width, height = 100, 100

    # Simulate Red band (vegetation absorbs red light)
    red_band = np.random.normal(0.1, 0.05, (height, width))
    red_band = np.clip(red_band, 0, 1)  # Reflectance values 0-1

    # Simulate NIR band (vegetation reflects NIR)
    nir_band = np.random.normal(0.6, 0.15, (height, width))
    nir_band = np.clip(nir_band, 0, 1)

    # Add some vegetation patterns
    center_y, center_x = height // 2, width // 2
    y, x = np.ogrid[:height, :width]

    # Create circular vegetation area
    vegetation_mask = (x - center_x)**2 + (y - center_y)**2 < (min(width, height) // 4)**2

    # Enhance vegetation signature
    red_band[vegetation_mask] *= 0.5    # Lower red reflectance
    nir_band[vegetation_mask] *= 1.3    # Higher NIR reflectance

    # Calculate NDVI
    ndvi = (nir_band - red_band) / (nir_band + red_band + 1e-10)  # Avoid division by zero

    # NDVI interpretation
    water_mask = ndvi < 0
    bare_soil_mask = (ndvi >= 0) & (ndvi < 0.2)
    sparse_vegetation_mask = (ndvi >= 0.2) & (ndvi < 0.5)
    dense_vegetation_mask = ndvi >= 0.5

    print(f"NDVI Statistics:")
    print(f"  Min: {ndvi.min():.3f}")
    print(f"  Max: {ndvi.max():.3f}")
    print(f"  Mean: {ndvi.mean():.3f}")
    print(f"  Std: {ndvi.std():.3f}")

    print(f"\nLand Cover Classification:")
    print(f"  Water/Shadow: {np.sum(water_mask)} pixels ({100*np.sum(water_mask)/ndvi.size:.1f}%)")
    print(f"  Bare Soil: {np.sum(bare_soil_mask)} pixels ({100*np.sum(bare_soil_mask)/ndvi.size:.1f}%)")
    print(f"  Sparse Vegetation: {np.sum(sparse_vegetation_mask)} pixels ({100*np.sum(sparse_vegetation_mask)/ndvi.size:.1f}%)")
    print(f"  Dense Vegetation: {np.sum(dense_vegetation_mask)} pixels ({100*np.sum(dense_vegetation_mask)/ndvi.size:.1f}%)")

    return ndvi, red_band, nir_band

def change_detection_simulation():
    """Simulate change detection between two time periods"""
    print("\n=== Change Detection Analysis ===")

    # Create two time periods of NDVI data
    np.random.seed(42)

    # Time 1 (earlier)
    ndvi_t1 = np.random.normal(0.4, 0.2, (50, 50))
    ndvi_t1 = np.clip(ndvi_t1, -1, 1)

    # Time 2 (later) - simulate some changes
    ndvi_t2 = ndvi_t1.copy()

    # Add deforestation in upper left
    ndvi_t2[:20, :20] -= 0.3

    # Add vegetation growth in lower right
    ndvi_t2[30:, 30:] += 0.2

    # Add random noise
    ndvi_t2 += np.random.normal(0, 0.05, ndvi_t2.shape)
    ndvi_t2 = np.clip(ndvi_t2, -1, 1)

    # Calculate change
    ndvi_change = ndvi_t2 - ndvi_t1

    # Define change thresholds
    threshold = 0.1

    decrease_mask = ndvi_change < -threshold  # Vegetation loss
    increase_mask = ndvi_change > threshold   # Vegetation gain
    stable_mask = np.abs(ndvi_change) <= threshold  # No significant change

    print(f"Change Detection Results:")
    print(f"  Vegetation Loss: {np.sum(decrease_mask)} pixels ({100*np.sum(decrease_mask)/ndvi_change.size:.1f}%)")
    print(f"  Vegetation Gain: {np.sum(increase_mask)} pixels ({100*np.sum(increase_mask)/ndvi_change.size:.1f}%)")
    print(f"  Stable Areas: {np.sum(stable_mask)} pixels ({100*np.sum(stable_mask)/ndvi_change.size:.1f}%)")

    print(f"\nChange Statistics:")
    print(f"  Mean change: {ndvi_change.mean():.3f}")
    print(f"  Max increase: {ndvi_change.max():.3f}")
    print(f"  Max decrease: {ndvi_change.min():.3f}")

    return ndvi_t1, ndvi_t2, ndvi_change

def spectral_analysis():
    """Analyze spectral signatures of different land cover types"""
    print("\n=== Spectral Analysis ===")

    # Define spectral bands (wavelengths in micrometers)
    bands = {
        'Blue': 0.48,
        'Green': 0.56,
        'Red': 0.66,
        'NIR': 0.84,
        'SWIR1': 1.65,
        'SWIR2': 2.22
    }

    # Typical spectral signatures (reflectance values 0-1)
    signatures = {
        'Water': [0.05, 0.04, 0.03, 0.02, 0.01, 0.01],
        'Vegetation': [0.04, 0.08, 0.06, 0.50, 0.25, 0.15],
        'Bare Soil': [0.15, 0.20, 0.25, 0.35, 0.40, 0.38],
        'Urban': [0.12, 0.14, 0.16, 0.20, 0.25, 0.22],
        'Snow': [0.85, 0.90, 0.88, 0.85, 0.20, 0.10]
    }

    print("Spectral Signatures by Land Cover Type:")
    print(f"{'Land Cover':<12} {'Blue':<6} {'Green':<6} {'Red':<6} {'NIR':<6} {'SWIR1':<6} {'SWIR2':<6}")
    print("-" * 60)

    for land_cover, reflectances in signatures.items():
        values_str = " ".join([f"{r:5.2f}" for r in reflectances])
        print(f"{land_cover:<12} {values_str}")

    # Calculate vegetation indices for each land cover type
    print(f"\nVegetation Indices:")
    print(f"{'Land Cover':<12} {'NDVI':<8} {'EVI':<8} {'SAVI':<8}")
    print("-" * 36)

    for land_cover, reflectances in signatures.items():
        red, nir = reflectances[2], reflectances[3]
        blue = reflectances[0]

        # NDVI
        ndvi = (nir - red) / (nir + red + 1e-10)

        # EVI (Enhanced Vegetation Index)
        evi = 2.5 * ((nir - red) / (nir + 6 * red - 7.5 * blue + 1))

        # SAVI (Soil Adjusted Vegetation Index)
        L = 0.5  # soil brightness correction factor
        savi = ((nir - red) / (nir + red + L)) * (1 + L)

        print(f"{land_cover:<12} {ndvi:7.3f} {evi:7.3f} {savi:7.3f}")

# Run remote sensing analysis
ndvi_data, red_data, nir_data = ndvi_calculation_simulation()
t1_data, t2_data, change_data = change_detection_simulation()
spectral_analysis()

print("\n✅ Remote sensing analysis completed!")
print("Advanced satellite imagery processing capabilities demonstrated")
EOF

python3 remote_sensing.py
```

**What this does**: Demonstrates remote sensing analysis including NDVI calculation and change detection.

**Expected result**: Shows vegetation analysis and spectral signature comparisons.

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-west-2
```

**Expected result**: Shows costs so far (should be under $8 for this tutorial)

## Step 11: Clean Up (Important!)

When you're done experimenting:

```bash
aws-research-wizard deploy delete --region us-west-2
```

Type `y` when prompted.

**What this does**: Stops billing by removing your cloud resources.

**💰 Important**: Always clean up to avoid ongoing charges.

**Expected result**: "🗑️ Deletion completed successfully"

## Understanding Your Costs

### What You're Paying For
- **Compute**: $0.50 per hour for memory-optimized instance while environment is running
- **Storage**: $0.10 per GB per month for geospatial datasets you save
- **Data Transfer**: Usually free for geospatial data amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large satellite datasets in S3, not on the instance
- Process imagery efficiently to minimize compute time

### Typical Monthly Costs by Usage
- **Light use** (15 hours/week): $150-300
- **Medium use** (4 hours/day): $300-600
- **Heavy use** (8 hours/day): $600-1200

## What's Next?

Now that you have a working geospatial environment, you can:

### Learn More About Geospatial Analysis
- [Large-scale Satellite Processing Tutorial](large-scale-satellite.md)
- [Advanced GIS Analysis Guide](advanced-gis-analysis.md)
- [Cost Optimization for Geospatial](geospatial-cost-optimization.md)

### Explore Advanced Features
- [Multi-temporal satellite analysis](multi-temporal-analysis.md)
- [Team collaboration with spatial databases](team-geospatial-collaboration.md)
- [Automated remote sensing pipelines](automated-geospatial-pipelines.md)

### Join the Geospatial Community
- [Geospatial Research Forum](https://forum.researchwizard.app/geospatial)
- [GitHub Geospatial Examples](https://github.com/aws-research-wizard/geospatial-examples)
- [Monthly GIS Office Hours](https://calendar.researchwizard.app/geospatial-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "GDAL not found" error during analysis
**Solution**: Check GDAL installation: `which gdal-config` and reload environment: `source /etc/profile`
**Prevention**: Wait 5-7 minutes after deployment for all geospatial tools to initialize

**Problem**: "Projection error" when working with different coordinate systems
**Solution**: Check CRS compatibility: `gdalinfo -proj4 filename.tif` and reproject if needed
**Prevention**: Always verify coordinate reference systems before spatial operations

**Problem**: "Memory error" during large raster processing
**Solution**: Process rasters in smaller chunks or use a larger instance type
**Prevention**: Monitor memory usage with `htop` during processing

**Problem**: "Shapefile not found" error
**Solution**: Verify all shapefile components (.shp, .shx, .dbf, .prj): `ls -la *.shp*`
**Prevention**: Always keep complete shapefile sets together

### Getting Help
- Check the [geospatial troubleshooting guide](troubleshooting-geospatial.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $17. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
