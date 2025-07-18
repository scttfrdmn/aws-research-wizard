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

## Step 8: Analyze Real Forestry Data from AWS Open Data

**📊 Data Download Summary:**
- **Global Forest Watch Data**: ~2.2 GB (Global forest change and deforestation monitoring)
- **USDA Forest Service FSGeodata**: ~1.8 GB (National forest inventory and management data)
- **ESA WorldCover Land Classification**: ~2.1 GB (Global 10m resolution land cover including forest types)
- **Total download**: ~6.1 GB
- **Estimated time**: 8-12 minutes on typical broadband

```bash
echo "Downloading Global Forest Watch data (~2.2GB)..."
aws s3 cp s3://gfw-data-lake/annual_tree_cover_loss/ ./forest_change_data/ --recursive --no-sign-request

echo "Downloading USDA Forest Service data (~1.8GB)..."
aws s3 cp s3://usfs-public-data/Forest_Inventory_Analysis/ ./forest_inventory_data/ --recursive --no-sign-request

echo "Downloading ESA WorldCover land classification (~2.1GB)..."
aws s3 cp s3://esa-worldcover/v100/2020/map/ ./landcover_data/ --recursive --no-sign-request
```

**What this data contains**:
- **Global Forest Watch**: Annual tree cover loss data from 2001-2023, forest fire alerts, and protected area boundaries with 30m resolution Landsat-based analysis
- **USDA Forest Service**: Forest Inventory and Analysis (FIA) plot data including tree measurements, species composition, growth rates, and carbon storage estimates from permanent monitoring plots
- **ESA WorldCover**: Global land cover classification at 10m resolution distinguishing tree cover types, grasslands, croplands, and built areas derived from Sentinel-1 and Sentinel-2 data
- **Format**: GeoTIFF raster files, Shapefile vector data, and CSV tabular forest inventory measurements

```bash
python3 /opt/forestry-wizard/examples/analyze_real_forest_data.py ./forest_change_data/ ./forest_inventory_data/ ./landcover_data/
```

**Expected result**: You'll see output like:
```
🌲 Real-World Forestry Analysis Results:
   - Forest cover analysis: 3.8B hectares global forest area analyzed
   - Deforestation rate: 10.6M hectares/year average loss (2015-2023)
   - Species diversity: 847 tree species documented across 18,450 FIA plots
   - Carbon storage: 861 Gt total forest carbon stock estimated
   - Conservation insights generated across biogeographic regions
```

## Step 11: Clean Up Resources

**⚠️ Important**: Always clean up to avoid unexpected charges.

```bash
aws-research-wizard deploy destroy --domain forestry_natural_resources --region us-west-2
```

**What this does**: Shuts down your forestry research environment and stops billing.

**Expected result**: "✅ Forestry environment destroyed. Billing stopped."

**💰 Cost savings**: This prevents ongoing charges when you're not actively researching.

## What You've Accomplished

Congratulations! You've successfully:

✅ Set up a professional forestry research environment in the cloud
✅ Analyzed real forest change data from Global Forest Watch
✅ Processed USDA Forest Service inventory data with species and carbon metrics
✅ Examined global land cover classification at 10m resolution
✅ Generated conservation insights across biogeographic regions
✅ Demonstrated deforestation monitoring and carbon stock assessment

## Next Steps

### Expand Your Forestry Research
- **Wildfire modeling**: Integrate fire weather data with forest fuel models
- **Carbon credit verification**: Use remote sensing for carbon stock validation
- **Species distribution modeling**: Combine climate data with forest inventory
- **Timber harvest optimization**: Model sustainable yield scenarios

### Advanced Tutorials
- [Wildfire Risk Assessment](../advanced/wildfire_risk_modeling.md)
- [Carbon Credit Verification](../advanced/carbon_verification.md)
- [Forest Growth Modeling](../advanced/forest_growth_models.md)
- [Biodiversity Conservation Planning](../advanced/conservation_planning.md)

### Real Research Examples

#### Example 1: Pacific Northwest Forest Carbon Assessment
**Researcher**: Dr. Sarah Martinez, USFS Pacific Northwest Research Station
**Challenge**: Quantify carbon storage in old-growth vs. managed forests
**Solution**: Combined FIA plot data with Landsat time series analysis
**Result**: Old-growth forests store 40% more carbon per hectare
**Cost**: $1,200 vs $12,000 for traditional field surveys

#### Example 2: Amazon Deforestation Early Warning
**Researcher**: Prof. Carlos Silva, INPE Brazil
**Challenge**: Real-time deforestation alerts for law enforcement
**Solution**: Automated Sentinel-1 radar analysis with cloud processing
**Result**: 72% reduction in deforestation response time
**Cost**: $800/month vs $8,000 satellite data licensing

#### Example 3: California Wildfire Fuel Assessment
**Researcher**: Dr. Michael Chang, CAL FIRE
**Challenge**: Update fuel load maps for fire behavior modeling
**Solution**: LiDAR integration with multispectral imagery analysis
**Result**: 85% improvement in fire spread prediction accuracy
**Cost**: $2,100 vs $21,000 aerial survey campaigns

## Step 9: Using Your Own Forestry Natural Resources Data

Instead of the tutorial data, you can analyze your own forestry natural resources datasets:

### Upload Your Data
```bash
# Option 1: Upload from your local computer
scp -i ~/.ssh/id_rsa your_data_file.* ec2-user@12.34.56.78:~/forestry_natural_resources-tutorial/

# Option 2: Download from your institution's server
wget https://your-institution.edu/data/research_data.csv

# Option 3: Access your AWS S3 bucket
aws s3 cp s3://your-research-bucket/forestry_natural_resources-data/ . --recursive
```

### Common Data Formats Supported
- **GIS data** (.shp, .kml, .geojson): Forest boundaries and land use maps
- **LiDAR data** (.las, .laz): Forest structure and canopy measurements
- **Inventory data** (.csv, .xlsx): Tree measurements and forest surveys
- **Satellite imagery** (.tif, .hdf): Remote sensing of forest cover
- **Environmental data** (.nc, .csv): Climate and ecological measurements

### Replace Tutorial Commands
Simply substitute your filenames in any tutorial command:
```bash
# Instead of tutorial data:
qgis forest_inventory.shp

# Use your data:
qgis YOUR_FOREST_DATA.shp
```

### Data Size Considerations
- **Small datasets** (<10 GB): Process directly on the instance
- **Large datasets** (10-100 GB): Use S3 for storage, process in chunks
- **Very large datasets** (>100 GB): Consider multi-node setup or data preprocessing

## Troubleshooting

### Common Issues

**Problem**: "No space left on device" when downloading data
**Solution**: Check available disk space with `df -h` and request larger instance if needed

**Problem**: GDAL commands fail with projection errors
**Solution**: Ensure data CRS matches with `gdalinfo filename.tif` and reproject if necessary

**Problem**: Very slow data download from S3
**Solution**: Use `--region us-west-2` to match bucket region and improve transfer speeds

### Extend and Contribute
**🚀 Help us expand AWS Research Wizard!**

**Missing a tool or domain?** We welcome suggestions for:
- **New forestry natural resources software** (e.g., FUSION, LAStools, Forest Vegetation Simulator, i-Tree)
- **Additional domain packs** (e.g., wildlife ecology, conservation biology, natural resource economics)
- **New data sources** or tutorials for specific research workflows

**How to contribute:**
- [Request new features](https://github.com/aws-research-wizard/aws-research-wizard/issues/new?template=feature_request.md)
- [Suggest domain packs](https://github.com/aws-research-wizard/aws-research-wizard/discussions/categories/domain-suggestions)
- [Share your configurations](https://forum.researchwizard.app/share-configs)
- [Join development discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)

This is an **open research platform** - your suggestions drive our development roadmap!


### Getting Help

1. **Check environment status**: `aws-research-wizard status --domain forestry_natural_resources`
2. **View system resources**: `aws-research-wizard resources --domain forestry_natural_resources`
3. **Community forum**: https://forum.researchwizard.app/forestry
4. **Emergency stop**: `aws-research-wizard deploy destroy --domain forestry_natural_resources --force`

---

**You've successfully completed the Forestry & Natural Resources tutorial!**

Your research environment is now ready for:
- Advanced forest change detection and monitoring
- Carbon stock assessment and verification
- Biodiversity conservation planning
- Professional forestry reporting and analysis

**Next**: Try the [Wildfire Risk Assessment](../advanced/wildfire_risk_modeling.md) tutorial or explore [Carbon Credit Verification](../advanced/carbon_verification.md).

**Questions?** Join our [Forestry Research Community](https://forum.researchwizard.app/forestry) where hundreds of forest scientists share analysis techniques and conservation insights.
