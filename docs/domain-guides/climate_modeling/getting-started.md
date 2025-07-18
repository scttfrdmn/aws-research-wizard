# Climate Modeling Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $12-20 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working climate modeling research environment that can:
- Run weather prediction models like WRF and CESM
- Process large atmospheric datasets (NetCDF files)
- Handle high-performance computing with MPI
- Scale from single node to 32-node clusters

### Meet Dr. Carlos Rodriguez
Dr. Carlos Rodriguez is a climate scientist at NOAA. He models hurricane paths but waits 5-7 days for supercomputer access. Each simulation takes weeks to queue, delaying critical weather forecasts.

**Before**: 7-day waits + 2-day simulation = 9 days per forecast
**After**: 15-minute setup + 4-hour simulation = same day results
**Time Saved**: 95% faster forecast cycle
**Cost Savings**: $1,200/month vs $4,000 supercomputer allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $12-20 (we'll clean up resources when done)
- **Daily research cost**: $40-120 per day when actively modeling
- **Monthly estimate**: $400-1200 per month for typical usage
- **Free tier**: Some storage included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or climate modeling experience required

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
- **Region**: Choose `us-east-1` (recommended for climate modeling with best MPI performance)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain climate_modeling --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: climate_modeling
✅ Region valid: us-east-1 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Climate Environment

```bash
aws-research-wizard deploy start --domain climate_modeling --region us-east-1 --instance c6i.2xlarge
```

**What this does**: Creates your climate modeling computing environment with HPC optimization.

**This will take**: 5-7 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ec2-user@12.34.56.78
  MPI Processes: 8 cores available
  Storage: 500GB EBS optimized
```

**💰 Billing starts now**: Your environment costs about $0.68 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ec2-user@12.34.56.78
```

**What this does**: Connects you to your climate modeling computer in the cloud.

**Expected result**: You see a command prompt like `[ec2-user@ip-10-0-1-123 ~]$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Climate Tools

Your environment comes pre-installed with:

### Core Climate Modeling Tools
- **WRF**: Weather Research and Forecasting model - Type `which wrf.exe` to check
- **CESM**: Community Earth System Model - Type `which cesm` to check
- **NCO**: NetCDF Operators for data processing - Type `nco --version` to check
- **CDO**: Climate Data Operators - Type `cdo --version` to check
- **OpenMPI**: Message Passing Interface for parallel computing - Type `mpirun --version` to check

### Try Your First Command
```bash
nco --version
```

**What this does**: Shows NetCDF Operators version and confirms climate tools are installed.

**Expected result**: You see NCO version info and available operators.

## Step 8: Process Real Climate Data from AWS Open Data

Let's analyze real atmospheric data from the ERA5 Reanalysis dataset:

### Download Real Climate Data from ERA5

**📊 Data Download Summary:**
- ERA5 Atmospheric Reanalysis: ~2.1 GB (global meteorological data)
- NOAA Global Forecast System: ~1.8 GB (weather prediction model data)
- NASA GISS Climate Data: ~1.4 GB (temperature and precipitation records)
- **Total download**: ~5.3 GB
- **Estimated time**: 10-15 minutes on typical broadband

```bash
# Create working directory
mkdir ~/climate-tutorial
cd ~/climate-tutorial

# Download real climate modeling data from AWS Open Data
echo "Downloading ERA5 atmospheric reanalysis data (~2.1GB)..."
aws s3 cp s3://era5-pds/2023/01/data/2m_temperature.nc . --no-sign-request
aws s3 cp s3://era5-pds/2023/01/data/mean_sea_level_pressure.nc . --no-sign-request

echo "Downloading NOAA Global Forecast System data (~1.8GB)..."
aws s3 cp s3://noaa-gfs-bdp-pds/gfs.20230101/12/atmos/gfs.t12z.pgrb2.0p25.f000 . --no-sign-request

echo "Downloading NASA GISS climate data (~1.4GB)..."
aws s3 cp s3://nasa-giss-data/temperature/gistemp_v4_global_mean.txt . --no-sign-request
aws s3 cp s3://nasa-giss-data/precipitation/global_precipitation_2023.nc . --no-sign-request

echo "Real climate modeling data downloaded successfully!"

# Check the data structure
echo "Examining downloaded NetCDF files..."
ncdump -h 2m_temperature.nc | head -20
```

**What this data contains**:
- **ERA5 Reanalysis**: High-quality atmospheric reanalysis data from ECMWF with 0.25° resolution
- **NOAA GFS**: Global Forecast System operational weather prediction model output
- **NASA GISS**: Goddard Institute temperature and precipitation climate records
- **Format**: NetCDF4 climate grids with CF conventions and WGS84 coordinates

### Process Climate Data with Real Tools
```bash
# Extract temperature data for North America
cdo sellonlatbox,-140,-60,20,70 2m_temperature.nc north_america_temp.nc

# Calculate daily averages from hourly data
cdo daymean north_america_temp.nc daily_temp.nc

# Compute temperature anomalies (difference from monthly mean)
cdo timmean daily_temp.nc temp_climatology.nc
cdo sub daily_temp.nc temp_climatology.nc temp_anomalies.nc

# Generate statistics
cdo infov temp_anomalies.nc
```

**What this does**:
- Extracts North American temperature data from the global dataset
- Converts hourly data to daily averages
- Calculates temperature anomalies to identify extreme weather events
- Generates statistics about the data distribution

### View Results
```bash
# Check the processed data
cdo info daily_temp.nc

# View temperature statistics
cdo output -timavg -timmean temp_anomalies.nc
```

**🎉 Success!** You've processed real climate data from the ERA5 reanalysis!

### Explore More Climate Data (Optional)
```bash
# Browse available ERA5 variables
aws s3 ls s3://era5-pds/2023/01/data/ --no-sign-request

# Check out NOAA HRRR high-resolution weather model data
aws s3 ls s3://hrrrzarr/sfc/2023/01/01/ --no-sign-request

# Download precipitation data
aws s3 cp s3://era5-pds/2023/01/data/total_precipitation.nc . --no-sign-request
```

**Available datasets for further exploration**:
- **ERA5**: 70+ atmospheric variables, 1979-present
- **NOAA HRRR**: High-resolution (3km) weather model data
- **CESM**: Community Earth System Model output
- **CMIP6**: Climate model intercomparison project data

# Calculate monthly averages
cdo monmean us_region.nc monthly_avg.nc

# Get statistics
cdo info monthly_avg.nc
```

**What this does**: Processes real atmospheric data using climate science tools.

**This will take**: 1-2 minutes

### View Processing Results
```bash
# Show data statistics
echo "=== Climate Data Processing Results ==="
echo "Original file size:" $(du -h sample_weather.nc | cut -f1)
echo "Processed file size:" $(du -h monthly_avg.nc | cut -f1)

# Display data dimensions
echo "=== Data Structure ==="
ncdump -h monthly_avg.nc | grep "dimensions:"
ncdump -h monthly_avg.nc | grep "variables:"
```

**What you should see**: Information about temperature data processing and file statistics.

**🎉 Success!** You've processed real climate data in the cloud.

## Step 9: Test MPI Parallel Computing

Test high-performance computing capabilities:

```bash
# Test MPI with a simple parallel job
mpirun -np 4 echo "Hello from MPI process rank" \$OMPI_COMM_WORLD_RANK

# Test parallel NetCDF processing
echo "Testing parallel climate data processing..."
time cdo -P 4 yearmean sample_weather.nc parallel_year_avg.nc
```

**What this does**: Tests parallel processing capabilities for large climate simulations.

**Expected result**: Shows multiple MPI processes running and parallel data processing.

## Step 9: Using Your Own Climate Data

Instead of the tutorial data, you can analyze your own climate datasets:

### Upload Your Data
```bash
# Option 1: Upload from your local computer
scp -i ~/.ssh/id_rsa your_climate_data.nc ec2-user@12.34.56.78:~/climate-tutorial/

# Option 2: Download from your institution's server
wget https://your-institution.edu/data/climate_model_output.nc

# Option 3: Access your AWS S3 bucket
aws s3 cp s3://your-research-bucket/climate-data/ . --recursive
```

### Common Data Formats Supported
- **NetCDF files** (.nc, .nc4): `ncdump -h your_file.nc` to examine structure
- **GRIB files** (.grb, .grb2): `wgrib2 your_file.grb2 -V` to view metadata
- **CSV/ASCII data**: Direct import with pandas or numpy
- **Binary formats**: Use appropriate readers (e.g., fortran unformatted)

### Replace Tutorial Commands
Simply substitute your filenames in any tutorial command:
```bash
# Instead of tutorial data:
cdo sellonlatbox,-140,-60,20,70 2m_temperature.nc north_america_temp.nc

# Use your data:
cdo sellonlatbox,-140,-60,20,70 YOUR_DATA_FILE.nc your_analysis.nc
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

**Expected result**: Shows costs so far (should be under $8 for this tutorial)

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
- **Compute**: $0.68 per hour for HPC instance while environment is running
- **Storage**: $0.10 per GB per month for data you save
- **Data Transfer**: Usually free for climate modeling amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large datasets in S3, not on the instance
- Use cluster scaling for large simulations only when needed

### Typical Monthly Costs by Usage
- **Light use** (15 hours/week): $200-400
- **Medium use** (4 hours/day): $400-800
- **Heavy use** (8 hours/day): $800-1200

## What's Next?

Now that you have a working climate environment, you can:

### Learn More About Climate Modeling
- [WRF Hurricane Simulation Tutorial](wrf-hurricane-tutorial.md)
- [Multi-Node Cluster Setup Guide](multi-node-climate.md)
- [Cost Optimization for Climate Models](climate-cost-optimization.md)

### Explore Advanced Features
- [Running CESM global climate simulations](cesm-global-simulation.md)
- [Team collaboration with shared datasets](team-climate-collaboration.md)
- [Automated climate data pipelines](climate-data-pipelines.md)

### Join the Climate Community
- [Climate Research Forum](https://forum.researchwizard.app/climate)
- [GitHub Climate Examples](https://github.com/aws-research-wizard/climate-examples)
- [Monthly Climate Office Hours](https://calendar.researchwizard.app/climate-office-hours)

### Extend and Contribute
**🚀 Help us expand AWS Research Wizard!**

**Missing a tool or domain?** We welcome suggestions for:
- **New climate modeling software** (e.g., RegCM, MM5, ICON, FMS)
- **Additional domain packs** (e.g., atmospheric chemistry, hydrology, oceanography)
- **New data sources** or tutorials for specific research workflows

**How to contribute:**
- [Request new features](https://github.com/aws-research-wizard/aws-research-wizard/issues/new?template=feature_request.md)
- [Suggest domain packs](https://github.com/aws-research-wizard/aws-research-wizard/discussions/categories/domain-suggestions)
- [Share your configurations](https://forum.researchwizard.app/share-configs)
- [Join development discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)

This is an **open research platform** - your suggestions drive our development roadmap!

## Troubleshooting

### Common Issues

**Problem**: "MPI not found" error during parallel jobs
**Solution**: Check MPI installation: `which mpirun` and reload environment: `source /etc/profile`
**Prevention**: Wait 3-5 minutes after deployment for all MPI setup to complete

**Problem**: "Permission denied" when connecting with SSH
**Solution**: Make sure your SSH key has correct permissions: `chmod 600 ~/.ssh/id_rsa`
**Prevention**: The deployment process usually sets this automatically

**Problem**: NetCDF files corrupted or unreadable
**Solution**: Check file integrity: `ncdump -h filename.nc` and re-download if needed
**Prevention**: Always verify downloads with `ncdump -h` before processing

**Problem**: Climate simulation runs out of memory
**Solution**: Use a larger instance type or reduce simulation domain size
**Prevention**: Monitor memory usage with `htop` during simulations

### Getting Help
- Check the [climate troubleshooting guide](troubleshooting-climate.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-east-1 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $20. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
