# Genomics Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $8-12 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working genomics research environment that can:
- Process DNA sequence files (FASTQ, SAM/BAM, VCF)
- Run popular tools like BWA, GATK, and SAMtools
- Handle datasets up to 500GB in size
- Cost 60% less than traditional computing clusters

### Meet Dr. Sarah Kim
Dr. Sarah Kim is a genomics researcher at Johns Hopkins. She studies rare genetic diseases but waits 3-5 days for university cluster access. Each analysis takes a week to complete, slowing down her research.

**Before**: 3-day waits + 1-week analysis = 10 days per study
**After**: 15-minute setup + 4-hour analysis = same day results
**Time Saved**: 95% faster research cycle
**Cost Savings**: $400/month vs $1,200 university allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $8-12 (we'll clean up resources when done)
- **Daily research cost**: $15-45 per day when actively using
- **Monthly estimate**: $150-450 per month for typical usage
- **Free tier**: Some storage included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or programming experience required

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
- **Region**: Choose `us-west-2` (recommended for genomics)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain genomics --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: genomics
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Genomics Environment

```bash
aws-research-wizard deploy start --domain genomics --region us-west-2 --instance r6i.large
```

**What this does**: Creates your genomics computing environment in the cloud.

**This will take**: 3-5 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ec2-user@12.34.56.78
  S3 Bucket: genomics-data-1234567
```

**💰 Billing starts now**: Your environment costs about $0.24 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ec2-user@12.34.56.78
```

**What this does**: Connects you to your genomics computer in the cloud.

**Expected result**: You see a command prompt like `[ec2-user@ip-10-0-1-123 ~]$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Genomics Tools

Your environment comes pre-installed with:

### Core Genomics Tools
- **BWA**: DNA sequence alignment - Type `bwa` to start
- **GATK**: Variant discovery - Type `gatk --help` to start
- **SAMtools**: Sequence data processing - Type `samtools` to start
- **FastQC**: Quality control - Type `fastqc --help` to start
- **bcftools**: Variant calling utilities - Type `bcftools` to start

### Try Your First Command
```bash
bwa
```

**What this does**: Shows BWA help and confirms it's installed correctly.

**Expected result**: You see BWA version info and usage instructions.

## Step 8: Run a Simple Analysis

Let's align some DNA sequences to test everything works:

### Download Sample Data
```bash
# Create working directory
mkdir ~/genomics-tutorial
cd ~/genomics-tutorial

# Download sample FASTQ files (small test data)
wget https://s3.amazonaws.com/aws-research-data/genomics/sample_R1.fastq.gz
wget https://s3.amazonaws.com/aws-research-data/genomics/sample_R2.fastq.gz
wget https://s3.amazonaws.com/aws-research-data/genomics/reference.fasta
```

### Index the Reference Genome
```bash
bwa index reference.fasta
```

**What this does**: Prepares the reference genome for fast searching.

**This will take**: 30 seconds

### Run Sequence Alignment
```bash
bwa mem reference.fasta sample_R1.fastq.gz sample_R2.fastq.gz > aligned.sam
```

**What this does**: Aligns your DNA sequences to the reference genome.

**This will take**: 1-2 minutes

### View Results
```bash
samtools view -H aligned.sam | head -5
```

**What you should see**: SAM file headers showing alignment statistics.

**🎉 Success!** You've run your first genomics analysis in the cloud.

## Step 9: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-west-2
```

**Expected result**: Shows costs so far (should be under $2 for this tutorial)

## Step 10: Clean Up (Important!)

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
- **Compute**: $0.24 per hour while environment is running
- **Storage**: $0.023 per GB per month for data you save
- **Data Transfer**: Usually free for genomics data amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 70% savings (advanced)
- Store large datasets in S3, not on the instance
- Monitor costs weekly with the built-in cost tracker

### Typical Monthly Costs by Usage
- **Light use** (8 hours/week): $75-125
- **Medium use** (4 hours/day): $300-450
- **Heavy use** (8 hours/day): $600-900

## What's Next?

Now that you have a working genomics environment, you can:

### Learn More About Genomics Tools
- [GATK Best Practices Pipeline Guide]
- [Large Dataset Processing Tutorial]
- [Cost Optimization for Genomics]

### Explore Advanced Features
- [Multi-sample variant calling]
- [Team collaboration setup]
- [Automated pipeline deployment]

### Join the Genomics Community
- [Genomics Research Forum]
- [GitHub Examples Repository]
- [Monthly Genomics Office Hours]

## Troubleshooting

### Common Issues

**Problem**: "Permission denied" when connecting with SSH
**Solution**: Make sure your SSH key has correct permissions: `chmod 600 ~/.ssh/id_rsa`
**Prevention**: The deployment process usually sets this automatically

**Problem**: "Instance not found" error
**Solution**: Check that your region matches: `aws-research-wizard deploy status --region us-west-2`
**Prevention**: Always specify the same region in all commands

**Problem**: BWA or GATK commands not found
**Solution**: Wait 2-3 more minutes after deployment for software installation to complete
**Prevention**: The "Deployment completed" message means infrastructure is ready, not software

### Getting Help
- Check the [genomics troubleshooting guide]
- Ask in [community forum]
- File an issue on [GitHub]

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $12. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
