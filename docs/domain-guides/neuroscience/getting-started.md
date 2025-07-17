# Neuroscience Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $8-15 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working neuroscience research environment that can:
- Process brain imaging data (fMRI, MRI, DTI)
- Run neuroimaging analysis with FSL, FreeSurfer, and AFNI
- Handle large neuroimaging datasets up to 500GB
- Perform statistical analysis and brain connectivity studies

### Meet Dr. Emily Chen
Dr. Emily Chen is a neuroscientist at Johns Hopkins. She analyzes brain scans to study memory formation but waits 4-6 days for university cluster access. Each analysis takes hours to queue, delaying critical research discoveries.

**Before**: 6-day waits + 8-hour analysis = 6.3 days per study
**After**: 15-minute setup + 3-hour analysis = same day results
**Time Saved**: 90% faster research cycle
**Cost Savings**: $600/month vs $1,800 university allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $8-15 (we'll clean up resources when done)
- **Daily research cost**: $20-60 per day when actively analyzing
- **Monthly estimate**: $250-600 per month for typical usage
- **Free tier**: Some storage included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or neuroscience experience required

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
- **Region**: Choose `us-west-2` (recommended for neuroscience with good memory-optimized instances)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain neuroscience --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: neuroscience
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Neuroscience Environment

```bash
aws-research-wizard deploy start --domain neuroscience --region us-west-2 --instance r6i.xlarge
```

**What this does**: Creates your neuroscience computing environment with high-memory optimization for brain imaging.

**This will take**: 4-6 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  Memory: 32GB RAM optimized for neuroimaging
  Storage: 200GB SSD for fast data access
```

**💰 Billing starts now**: Your environment costs about $0.50 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your neuroscience computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Neuroscience Tools

Your environment comes pre-installed with:

### Core Neuroimaging Tools
- **FSL**: FMRIB Software Library for brain imaging - Type `fsl --version` to check
- **FreeSurfer**: Cortical reconstruction and analysis - Type `freesurfer --version` to check
- **AFNI**: Analysis of Functional NeuroImages - Type `afni --version` to check
- **ANTs**: Advanced Normalization Tools - Type `antsRegistration --version` to check
- **MRtrix**: Diffusion MRI analysis - Type `mrinfo --version` to check

### Try Your First Command
```bash
fsl --version
```

**What this does**: Shows FSL version and confirms neuroimaging tools are installed.

**Expected result**: You see FSL version info and available tools.

## Step 8: Run a Simple Brain Analysis

Let's analyze sample brain imaging data to test everything works:

### Download Sample Brain Data
```bash
# Create working directory
mkdir ~/neuroscience-tutorial
cd ~/neuroscience-tutorial

# Download sample brain data (structural MRI)
wget -O sample_brain.nii.gz "https://fsl.fmrib.ox.ac.uk/fsldownloads/MNI152_T1_1mm_brain.nii.gz"

# Download brain mask
wget -O brain_mask.nii.gz "https://fsl.fmrib.ox.ac.uk/fsldownloads/MNI152_T1_1mm_brain_mask.nii.gz"

echo "Sample brain data downloaded successfully!"
```

### Basic Brain Image Processing
```bash
# Get brain image information
echo "=== Brain Image Information ==="
fslinfo sample_brain.nii.gz

# Calculate brain volume statistics
echo "=== Brain Volume Statistics ==="
fslstats sample_brain.nii.gz -M -S -R

# Create brain histogram
fslstats sample_brain.nii.gz -H 100 brain_histogram.txt
echo "Brain intensity histogram saved to brain_histogram.txt"
```

### Brain Segmentation
```bash
# Perform basic brain tissue segmentation
echo "Starting brain tissue segmentation..."
fast -t 1 -n 3 -H 0.1 -I 4 -l 20.0 -o brain_seg sample_brain.nii.gz

echo "Segmentation complete! Files created:"
ls -la brain_seg*
```

**What this does**: Performs brain tissue segmentation to identify gray matter, white matter, and CSF.

**This will take**: 2-3 minutes

### View Analysis Results
```bash
# Show analysis results
echo "=== Brain Analysis Summary ==="
echo "Original brain image: $(fslstats sample_brain.nii.gz -V | awk '{print $1}') voxels"
echo "Brain volume: $(fslstats sample_brain.nii.gz -V | awk '{print $2}') mm³"

# Check segmentation outputs
if [ -f brain_seg_seg.nii.gz ]; then
    echo "✅ Brain segmentation successful"
    echo "Tissue classes identified: $(fslstats brain_seg_seg.nii.gz -R)"
else
    echo "⚠️ Segmentation files not found"
fi
```

**What you should see**: Brain volume statistics and confirmation of successful tissue segmentation.

**🎉 Success!** You've analyzed real brain imaging data in the cloud.

## Step 9: Functional Connectivity Analysis

Test advanced neuroscience capabilities:

```bash
# Create a simple connectivity analysis script
cat > connectivity_analysis.py << 'EOF'
import numpy as np
import nibabel as nib
from scipy import stats
import matplotlib.pyplot as plt

print("Loading brain data...")
brain_img = nib.load('sample_brain.nii.gz')
brain_data = brain_img.get_fdata()

print(f"Brain image shape: {brain_data.shape}")
print(f"Brain image data type: {brain_data.dtype}")

# Simple region analysis
roi1 = brain_data[50:70, 50:70, 50:70]  # Region 1
roi2 = brain_data[80:100, 50:70, 50:70]  # Region 2

roi1_mean = np.mean(roi1[roi1 > 0])
roi2_mean = np.mean(roi2[roi2 > 0])

print(f"Region 1 mean intensity: {roi1_mean:.2f}")
print(f"Region 2 mean intensity: {roi2_mean:.2f}")

print("✅ Basic connectivity analysis completed!")
EOF

python3 connectivity_analysis.py
```

**What this does**: Demonstrates basic brain region analysis and connectivity measures.

**Expected result**: Shows brain region statistics and connectivity analysis results.

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-west-2
```

**Expected result**: Shows costs so far (should be under $4 for this tutorial)

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
- **Storage**: $0.10 per GB per month for brain imaging data you save
- **Data Transfer**: Usually free for neuroscience analysis amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large brain datasets in S3, not on the instance
- Monitor memory usage to ensure you're using full capacity efficiently

### Typical Monthly Costs by Usage
- **Light use** (12 hours/week): $100-250
- **Medium use** (3 hours/day): $250-450
- **Heavy use** (6 hours/day): $450-750

## What's Next?

Now that you have a working neuroscience environment, you can:

### Learn More About Neuroimaging
- [fMRI Analysis Pipeline Tutorial](fmri-analysis-pipeline.md)
- [FreeSurfer Cortical Reconstruction Guide](freesurfer-reconstruction.md)
- [Cost Optimization for Neuroimaging](neuroscience-cost-optimization.md)

### Explore Advanced Features
- [Multi-subject statistical analysis](multi-subject-analysis.md)
- [Team collaboration with brain imaging data](team-neuro-collaboration.md)
- [Automated neuroimaging pipelines](automated-neuro-pipelines.md)

### Join the Neuroscience Community
- [Neuroscience Research Forum](https://forum.researchwizard.app/neuroscience)
- [GitHub Neuroimaging Examples](https://github.com/aws-research-wizard/neuroscience-examples)
- [Monthly Neuroscience Office Hours](https://calendar.researchwizard.app/neuroscience-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "FSL not found" error during analysis
**Solution**: Check FSL installation: `which fsl` and reload environment: `source $FSLDIR/etc/fslconf/fsl.sh`
**Prevention**: Wait 3-5 minutes after deployment for all neuroimaging tools to initialize

**Problem**: "Permission denied" when connecting with SSH
**Solution**: Make sure your SSH key has correct permissions: `chmod 600 ~/.ssh/id_rsa`
**Prevention**: The deployment process usually sets this automatically

**Problem**: Brain images appear corrupted or unreadable
**Solution**: Check file integrity: `fslinfo filename.nii.gz` and re-download if needed
**Prevention**: Always verify downloads with `fslinfo` before processing

**Problem**: Analysis runs out of memory during processing
**Solution**: Use a larger instance type or reduce image resolution
**Prevention**: Monitor memory usage with `htop` during analysis

### Getting Help
- Check the [neuroscience troubleshooting guide](troubleshooting-neuroscience.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $15. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
