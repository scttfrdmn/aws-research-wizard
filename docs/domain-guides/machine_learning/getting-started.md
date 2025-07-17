# Machine Learning Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $10-18 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working machine learning research environment that can:
- Train neural networks using PyTorch and TensorFlow
- Run GPU-accelerated deep learning models
- Handle datasets up to 1TB in size
- Scale from single GPU to multi-GPU clusters

### Meet Dr. Maya Patel
Dr. Maya Patel is an AI researcher at Stanford. She trains large language models but waits weeks for GPU access on university clusters. Each training run takes days to schedule, slowing down her research breakthroughs.

**Before**: 2-week waits + 3-day training = 17 days per experiment
**After**: 10-minute setup + 6-hour training = same day results
**Time Saved**: 95% faster research cycle
**Cost Savings**: $800/month vs $2,400 university allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $10-18 (we'll clean up resources when done)
- **Daily research cost**: $25-80 per day when actively training
- **Monthly estimate**: $200-800 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or machine learning experience required

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
- **Region**: Choose `us-west-2` (recommended for ML with good GPU availability)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain machine_learning --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: machine_learning
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your ML Environment

```bash
aws-research-wizard deploy start --domain machine_learning --region us-west-2 --instance g5.xlarge
```

**What this does**: Creates your machine learning computing environment with GPU acceleration.

**This will take**: 4-6 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  GPU Type: NVIDIA A10G (24GB memory)
  Jupyter Lab: http://12.34.56.78:8888
```

**💰 Billing starts now**: Your environment costs about $1.20 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your machine learning computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your ML Tools

Your environment comes pre-installed with:

### Core ML Frameworks
- **PyTorch**: Deep learning framework - Type `python -c "import torch; print(torch.__version__)"` to check
- **TensorFlow**: Google's ML framework - Type `python -c "import tensorflow as tf; print(tf.__version__)"` to check
- **scikit-learn**: Classical ML library - Type `python -c "import sklearn; print(sklearn.__version__)"` to check
- **Jupyter Lab**: Interactive notebooks - Access at `http://your-ip:8888`
- **CUDA**: GPU acceleration - Type `nvidia-smi` to check GPU status

### Try Your First Command
```bash
nvidia-smi
```

**What this does**: Shows your GPU information and confirms CUDA is working.

**Expected result**: You see GPU details including "NVIDIA A10G" and memory usage.

## Step 8: Run a Simple ML Training

Let's train a neural network to test everything works:

### Download Sample Dataset
```bash
# Create working directory
mkdir ~/ml-tutorial
cd ~/ml-tutorial

# Download CIFAR-10 dataset (small image classification dataset)
python3 -c "
import torch
import torchvision
torchvision.datasets.CIFAR10(root='./data', train=True, download=True)
print('Dataset downloaded successfully!')
"
```

### Create Training Script
```bash
cat > simple_training.py << 'EOF'
import torch
import torch.nn as nn
import torchvision
import torchvision.transforms as transforms

# Check if GPU is available
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
print(f'Using device: {device}')

# Load CIFAR-10 dataset
transform = transforms.Compose([
    transforms.ToTensor(),
    transforms.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))
])

trainset = torchvision.datasets.CIFAR10(root='./data', train=True, transform=transform)
trainloader = torch.utils.data.DataLoader(trainset, batch_size=32, shuffle=True)

# Simple CNN model
class SimpleCNN(nn.Module):
    def __init__(self):
        super(SimpleCNN, self).__init__()
        self.conv1 = nn.Conv2d(3, 32, 3, padding=1)
        self.conv2 = nn.Conv2d(32, 64, 3, padding=1)
        self.pool = nn.MaxPool2d(2, 2)
        self.fc1 = nn.Linear(64 * 8 * 8, 512)
        self.fc2 = nn.Linear(512, 10)
        self.relu = nn.ReLU()

    def forward(self, x):
        x = self.pool(self.relu(self.conv1(x)))
        x = self.pool(self.relu(self.conv2(x)))
        x = x.view(-1, 64 * 8 * 8)
        x = self.relu(self.fc1(x))
        x = self.fc2(x)
        return x

# Initialize model and move to GPU
model = SimpleCNN().to(device)
criterion = nn.CrossEntropyLoss()
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)

print('Starting training...')
for epoch in range(2):  # Train for 2 epochs
    for i, (images, labels) in enumerate(trainloader):
        images, labels = images.to(device), labels.to(device)

        optimizer.zero_grad()
        outputs = model(images)
        loss = criterion(outputs, labels)
        loss.backward()
        optimizer.step()

        if i % 100 == 0:
            print(f'Epoch [{epoch+1}/2], Step [{i+1}/{len(trainloader)}], Loss: {loss.item():.4f}')

print('Training completed! 🎉')
EOF
```

### Run Training
```bash
python3 simple_training.py
```

**What this does**: Trains a neural network on your GPU to classify images.

**This will take**: 2-3 minutes

**What you should see**: Training progress with decreasing loss values and GPU utilization.

**🎉 Success!** You've trained your first neural network in the cloud.

## Step 9: Access Jupyter Lab

Open your web browser and go to: `http://your-ip-address:8888`

Replace `your-ip-address` with the IP from Step 5.

**What this gives you**: Interactive notebooks for data science and ML experiments.

**Expected result**: Jupyter Lab interface opens with file browser and notebook options.

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-west-2
```

**Expected result**: Shows costs so far (should be under $5 for this tutorial)

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
- **Compute**: $1.20 per hour for GPU instance while environment is running
- **Storage**: $0.10 per GB per month for data you save
- **Data Transfer**: Usually free for ML training amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 70% savings (advanced)
- Store large datasets in S3, not on the instance
- Monitor GPU utilization to ensure you're using the full capacity

### Typical Monthly Costs by Usage
- **Light use** (10 hours/week): $200-350
- **Medium use** (4 hours/day): $600-900
- **Heavy use** (8 hours/day): $1200-1800

## What's Next?

Now that you have a working ML environment, you can:

### Learn More About Machine Learning
- [Large Model Training Tutorial](large-model-training.md)
- [Multi-GPU Setup Guide](multi-gpu-setup.md)
- [Cost Optimization for ML](ml-cost-optimization.md)

### Explore Advanced Features
- [Distributed training across multiple instances](distributed-training.md)
- [Team collaboration with shared environments](team-collaboration.md)
- [Automated model deployment pipelines](model-deployment.md)

### Join the ML Community
- [ML Research Forum](https://forum.researchwizard.app/ml)
- [GitHub Examples Repository](https://github.com/aws-research-wizard/ml-examples)
- [Monthly ML Office Hours](https://calendar.researchwizard.app/ml-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "CUDA out of memory" error during training
**Solution**: Reduce batch size in your training script: change `batch_size=32` to `batch_size=16`
**Prevention**: Monitor GPU memory usage with `nvidia-smi` before training

**Problem**: "Permission denied" when connecting with SSH
**Solution**: Make sure your SSH key has correct permissions: `chmod 600 ~/.ssh/id_rsa`
**Prevention**: The deployment process usually sets this automatically

**Problem**: Jupyter Lab not accessible in browser
**Solution**: Check security group allows port 8888: `aws-research-wizard deploy status --region us-west-2`
**Prevention**: This should be configured automatically during deployment

**Problem**: PyTorch or TensorFlow not using GPU
**Solution**: Check CUDA installation: `nvidia-smi` and restart Python kernel
**Prevention**: Wait 2-3 minutes after deployment for all software to initialize

### Getting Help
- Check the [ML troubleshooting guide](troubleshooting-ml.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $18. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
