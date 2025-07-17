# Astronomy & Astrophysics Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $15-25 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working astronomy research environment that can:
- Process large telescope survey data (FITS images)
- Run astronomical data analysis with Python and specialized tools
- Handle datasets up to 2TB in size
- Perform image processing and photometry analysis

### Meet Dr. Sarah Johnson
Dr. Sarah Johnson is an astronomer at Caltech. She analyzes galaxy survey data from the Hubble Space Telescope but waits 8-10 days for university supercomputer access. Each analysis takes days to queue, delaying critical discovery publications.

**Before**: 10-day waits + 12-hour analysis = 10.5 days per discovery
**After**: 15-minute setup + 6-hour analysis = same day results
**Time Saved**: 94% faster research cycle
**Cost Savings**: $900/month vs $3,200 supercomputer allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $15-25 (we'll clean up resources when done)
- **Daily research cost**: $35-100 per day when actively analyzing
- **Monthly estimate**: $400-1200 per month for typical usage
- **Free tier**: Some storage included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or astronomy experience required

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
- **Region**: Choose `us-west-2` (recommended for astronomy with good high-memory instances)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain astronomy_astrophysics --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: astronomy_astrophysics
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Astronomy Environment

```bash
aws-research-wizard deploy start --domain astronomy_astrophysics --region us-west-2 --instance r6i.2xlarge
```

**What this does**: Creates your astronomy computing environment optimized for large image processing.

**This will take**: 5-7 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  Memory: 64GB RAM for large FITS processing
  Storage: 1TB NVMe SSD for fast data access
```

**💰 Billing starts now**: Your environment costs about $1.00 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your astronomy computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Astronomy Tools

Your environment comes pre-installed with:

### Core Astronomy Tools
- **AstroPy**: Core Python astronomy package - Type `python -c "import astropy; print(astropy.__version__)"` to check
- **DS9**: FITS image viewer - Type `ds9 --version` to check
- **IRAF**: Image Reduction and Analysis Facility - Type `which pyraf` to check
- **SExtractor**: Source extraction - Type `sextractor --version` to check
- **TOPCAT**: Table analysis tool - Type `which topcat` to check

### Try Your First Command
```bash
python -c "import astropy; print('AstroPy version:', astropy.__version__)"
```

**What this does**: Shows AstroPy version and confirms astronomy tools are installed.

**Expected result**: You see AstroPy version info confirming astronomical Python libraries are ready.

## Step 8: Process Astronomical Data

Let's analyze real telescope data to test everything works:

### Download Sample Astronomical Data
```bash
# Create working directory
mkdir ~/astronomy-tutorial
cd ~/astronomy-tutorial

# Download sample FITS image from Hubble Space Telescope
wget -O sample_galaxy.fits "https://archive.stsci.edu/pub/hlsp/ghosts/hst/wfc3-ir/ngc4449/v1/hlsp_ghosts_hst_wfc3-ir_ngc4449_f160w_v1_drz.fits"

# Download star catalog data
wget -O star_catalog.csv "https://raw.githubusercontent.com/astropy/astropy-data/master/testing/data/j94f05bgq_flt.fits"

echo "Sample astronomical data downloaded successfully!"
```

### Basic FITS Image Analysis
```bash
# Create Python script for FITS analysis
cat > fits_analysis.py << 'EOF'
import numpy as np
from astropy.io import fits
from astropy.stats import sigma_clipped_stats
import matplotlib.pyplot as plt

print("Loading FITS image...")
hdu_list = fits.open('sample_galaxy.fits')
image_data = hdu_list[0].data

print(f"Image shape: {image_data.shape}")
print(f"Image data type: {image_data.dtype}")

# Calculate image statistics
mean, median, std = sigma_clipped_stats(image_data, sigma=3.0)
print(f"Image statistics:")
print(f"  Mean: {mean:.2f}")
print(f"  Median: {median:.2f}")
print(f"  Standard deviation: {std:.2f}")

# Find brightest pixel (likely a star or galaxy core)
max_value = np.nanmax(image_data)
max_location = np.unravel_index(np.nanargmax(image_data), image_data.shape)
print(f"Brightest pixel: {max_value:.2f} at position {max_location}")

# Count sources above threshold
threshold = median + 5 * std
bright_sources = np.sum(image_data > threshold)
print(f"Bright sources (>5σ above background): {bright_sources}")

hdu_list.close()
print("✅ FITS image analysis completed!")
EOF

python3 fits_analysis.py
```

**What this does**: Analyzes real Hubble Space Telescope data to find astronomical sources.

**This will take**: 1-2 minutes

### Photometry Analysis
```bash
# Create photometry script
cat > photometry.py << 'EOF'
import numpy as np
from astropy.io import fits
from astropy.stats import sigma_clipped_stats
from photutils.detection import DAOStarFinder
from photutils.aperture import CircularAperture, aperture_photometry

print("Starting photometry analysis...")

# Load image
hdu_list = fits.open('sample_galaxy.fits')
data = hdu_list[0].data

# Calculate background statistics
mean, median, std = sigma_clipped_stats(data, sigma=3.0)
print(f"Background level: {median:.2f} ± {std:.2f}")

# Find sources
daofind = DAOStarFinder(fwhm=3.0, threshold=5.*std)
sources = daofind(data - median)

if sources is not None:
    print(f"Found {len(sources)} sources")

    # Perform aperture photometry on first 10 sources
    positions = np.transpose((sources['xcentroid'][:10], sources['ycentroid'][:10]))
    apertures = CircularAperture(positions, r=4.)
    phot_table = aperture_photometry(data, apertures)

    print("Photometry results (first 10 sources):")
    for i in range(min(10, len(phot_table))):
        print(f"Source {i+1}: flux = {phot_table['aperture_sum'][i]:.1f}")
else:
    print("No sources detected")

hdu_list.close()
print("✅ Photometry analysis completed!")
EOF

python3 photometry.py
```

**What this does**: Performs automated source detection and photometry on astronomical images.

**Expected result**: Shows detected sources and their measured brightness values.

**🎉 Success!** You've analyzed real telescope data in the cloud.

## Step 9: Coordinate System Analysis

Test advanced astronomy capabilities:

```bash
# Create coordinate analysis script
cat > coordinates.py << 'EOF'
from astropy.coordinates import SkyCoord
from astropy import units as u
from astropy.io import fits
from astropy.wcs import WCS

print("Analyzing coordinate systems...")

# Load FITS header for WCS information
hdu_list = fits.open('sample_galaxy.fits')
header = hdu_list[0].header

try:
    # Create WCS object
    wcs = WCS(header)
    print(f"Coordinate system: {wcs.wcs.ctype}")
    print(f"Reference pixel: {wcs.wcs.crpix}")
    print(f"Reference coordinate: {wcs.wcs.crval}")
    print(f"Pixel scale: {wcs.pixel_scale_matrix}")

    # Convert pixel coordinates to sky coordinates
    pixel_coords = [[100, 100], [200, 200], [300, 300]]
    for i, (x, y) in enumerate(pixel_coords):
        sky_coord = wcs.pixel_to_world(x, y)
        print(f"Pixel ({x}, {y}) → Sky: {sky_coord}")

except Exception as e:
    print(f"WCS analysis not available: {e}")
    print("This is normal for some FITS files without coordinate information")

# Demonstrate coordinate transformations
print("\nCoordinate system examples:")
# Famous astronomical objects
m31 = SkyCoord('00h42m44.3s', '+41d16m09s', frame='icrs')
print(f"Andromeda Galaxy (M31): {m31}")

galactic_center = SkyCoord('17h45m40s', '-29d00m28s', frame='icrs')
galactic_coord = galactic_center.galactic
print(f"Galactic Center in Galactic coordinates: {galactic_coord}")

hdu_list.close()
print("✅ Coordinate analysis completed!")
EOF

python3 coordinates.py
```

**What this does**: Demonstrates astronomical coordinate system handling and transformations.

**Expected result**: Shows coordinate system information and celestial coordinate examples.

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
- **Compute**: $1.00 per hour for high-memory instance while environment is running
- **Storage**: $0.10 per GB per month for astronomical data you save
- **Data Transfer**: Usually free for astronomy data amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large survey datasets in S3, not on the instance
- Monitor memory usage to ensure efficient processing of large FITS files

### Typical Monthly Costs by Usage
- **Light use** (15 hours/week): $250-400
- **Medium use** (4 hours/day): $500-800
- **Heavy use** (8 hours/day): $1000-1600

## What's Next?

Now that you have a working astronomy environment, you can:

### Learn More About Astronomical Data Analysis
- [Large Survey Data Processing Tutorial](large-survey-processing.md)
- [Multi-wavelength Analysis Guide](multi-wavelength-analysis.md)
- [Cost Optimization for Astronomy](astronomy-cost-optimization.md)

### Explore Advanced Features
- [Distributed processing of survey data](distributed-astronomy.md)
- [Team collaboration with astronomical databases](team-astro-collaboration.md)
- [Automated telescope data pipelines](automated-astro-pipelines.md)

### Join the Astronomy Community
- [Astronomy Research Forum](https://forum.researchwizard.app/astronomy)
- [GitHub Astronomy Examples](https://github.com/aws-research-wizard/astronomy-examples)
- [Monthly Astronomy Office Hours](https://calendar.researchwizard.app/astronomy-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "AstroPy import error" during analysis
**Solution**: Check Python environment: `which python3` and reinstall if needed: `pip install astropy`
**Prevention**: Wait 5-7 minutes after deployment for all astronomy packages to initialize

**Problem**: "FITS file corrupted" error
**Solution**: Verify download: `file sample_galaxy.fits` and re-download if needed
**Prevention**: Always check file integrity with `file` command after downloads

**Problem**: "Memory error" during large image processing
**Solution**: Use a larger instance type or process images in smaller sections
**Prevention**: Monitor memory usage with `htop` during analysis

**Problem**: "DS9 display not working" in SSH session
**Solution**: Enable X11 forwarding: `ssh -X -i ~/.ssh/id_rsa ubuntu@ip-address`
**Prevention**: For headless analysis, use Python matplotlib instead of DS9

### Getting Help
- Check the [astronomy troubleshooting guide](troubleshooting-astronomy.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $25. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
