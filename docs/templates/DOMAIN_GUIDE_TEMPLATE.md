# [Domain Name] Research Environment Guide

> **Time to Complete**: 20 minutes
> **Cost**: $5-15 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working [domain] research environment that can:
- [Key capability 1]
- [Key capability 2]
- [Key capability 3]

### Meet [Researcher Name]
[Researcher persona] is a [field] researcher at [institution]. They need to [research goal] but [current challenge]. With AWS Research Wizard, [outcome achieved].

**Before**: [Problem description]
**After**: [Solution achieved in concrete terms]
**Time Saved**: [Specific time savings]
**Cost Savings**: [Specific cost savings]

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $5-15 (we'll clean up resources when done)
- **Daily research cost**: $X-Y per day when actively using
- **Monthly estimate**: $XX-YY per month for typical usage
- **Free tier**: Some services included free for first 12 months

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
Download from: [link to Windows installer]

**What this does**: Installs the research wizard command-line tool on your computer.

**Expected result**: You should see "Installation successful" message.

**⚠️ If you see an error**: [Common error solutions]

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
- **AWS Access Key**: [How to find this]
- **Secret Key**: [How to find this]
- **Region**: Choose `us-west-2` (recommended)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: [Troubleshooting steps]

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain [domain-name] --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: [domain-name]
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Research Environment

```bash
aws-research-wizard deploy start --domain [domain-name] --region us-west-2 --instance t3.small
```

**What this does**: Creates your research computing environment in the cloud.

**This will take**: 3-5 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ec2-user@12.34.56.78
```

**💰 Billing starts now**: Your environment costs about $X per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ec2-user@[your-ip-address]
```

**What this does**: Connects you to your research computer in the cloud.

**Expected result**: You see a command prompt like `[ec2-user@ip-10-0-1-123 ~]$`

**⚠️ If connection fails**: [SSH troubleshooting]

## Step 7: Explore Your Tools

Your environment comes pre-installed with:

### [Domain-Specific Tools List]
- **[Tool 1]**: [What it does] - Type `[command]` to start
- **[Tool 2]**: [What it does] - Type `[command]` to start
- **[Tool 3]**: [What it does] - Type `[command]` to start

### Try Your First Command
```bash
[domain-specific example command]
```

**What this does**: [Explanation of command]

**Expected result**: [What you should see]

## Step 8: Run a Simple Analysis

Let's do a quick [domain-specific task] to test everything works:

### Download Sample Data
```bash
[commands to get sample data]
```

### Run Analysis
```bash
[commands to run simple analysis]
```

### View Results
```bash
[commands to see results]
```

**What you should see**: [Description of expected output]

**🎉 Success!** You've run your first [domain] analysis in the cloud.

## Step 9: Monitor Your Costs

Check your current spending:

```bash
aws-research-wizard monitor costs --region us-west-2
```

**Expected result**: Shows costs so far (should be under $5 for this tutorial)

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
- **Compute**: $X.XX per hour while environment is running
- **Storage**: $X.XX per GB per month for data you save
- **Data Transfer**: Usually free for research amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Monitor costs weekly with the built-in cost tracker

### Typical Monthly Costs by Usage
- **Light use** (4 hours/week): $XX-YY
- **Medium use** (2 hours/day): $XX-YY
- **Heavy use** (8 hours/day): $XX-YY

## What's Next?

Now that you have a working [domain] environment, you can:

### Learn More About [Domain]
- [Link to advanced domain guide]
- [Link to real-world examples]
- [Link to optimization tips]

### Explore Other Features
- [Data management guide]
- [Team collaboration setup]
- [Advanced configuration options]

### Join the Community
- [Community forum link]
- [GitHub discussions]
- [Example repository]

## Troubleshooting

### Common Issues

**Problem**: [Common error 1]
**Solution**: [Step-by-step fix]
**Prevention**: [How to avoid this]

**Problem**: [Common error 2]
**Solution**: [Step-by-step fix]
**Prevention**: [How to avoid this]

### Getting Help
- Check the [troubleshooting guide]
- Ask in [community forum]
- File an issue on [GitHub]

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $15. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box]

**What would you add?** [Text box]

**Rate the clarity (1-5)**: [Star rating]

---

*Last updated: [Date] | Reading level: 9th grade | Tutorial tested: [Date]*
