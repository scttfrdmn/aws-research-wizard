# AWS Research Wizard - LLM Vision & Examples

## Core Philosophy: True Wizard Intelligence

The AWS Research Wizard should be a **magical bridge** between research intent and infrastructure reality. An LLM provides the intelligence to understand research needs and translate them into optimal AWS deployments.

## The Two-World Problem

### Researchers Know:
- "I need to analyze protein structures"
- "I want to train a neural network on medical images"
- "I need to process climate model outputs"

### Researchers DON'T Know:
- GPU vs CPU instances
- Spack vs conda vs pip
- Network optimization, storage configurations

### Infrastructure People Know:
- Instance types, networking, storage
- Package managers, dependencies
- Cost optimization, security

### Infrastructure People DON'T Know:
- What tools geneticists actually use
- Which datasets are standard for climate research
- How computational chemists structure their workflows

**The LLM Wizard = Universal Translator**

---

## Conversational Flow Examples

### Option 1: "Don't Know, Don't Care" (Full Magic)
```
User: "I want to do protein structure analysis"
Wizard: "Perfect! Setting up your environment now..."
[5 minutes later]
Wizard: "Ready! Your protein analysis environment is at: ssh user@ip"
```

### Option 2: "Show Me What You're Planning"
```
User: "I want to do protein structure analysis - show me your plan first"

Wizard: "Here's what I'm planning:
📋 ANALYSIS PLAN
- Tools: PyMOL, ChimeraX, AlphaFold2
- Data: PDB database subset
- Instance: g4dn.xlarge (GPU for molecular visualization)
- Storage: 100GB for structure files
- Cost: ~$0.50/hour

Does this look right, or would you like me to adjust anything?"

User: "Looks good, deploy it!"
```

### Option 3: "I Have Specific Needs"
```
User: "Protein analysis, but I need GROMACS and only have $10 budget"

Wizard: "Got it! Adjusting for GROMACS molecular dynamics and budget:
- Tools: GROMACS, PyMOL (basic visualization)
- Instance: c5.large (CPU-only for budget)
- Spot instance: 70% cost savings
- Runtime limit: ~20 hours within budget
- Warning: No GPU = slower visualization

Proceed with budget-optimized setup?"
```

### Option 4: "I Have My Own Data"
```
User: "I want to do protein structure analysis and my data is here: s3://my-lab-bucket/protein-structures/"

Wizard: "Perfect! I can see your S3 bucket contains:
📁 DATA ANALYSIS
- 47 PDB files (~2.3GB)
- 12 trajectory files (.dcd format)
- Custom topology files

📋 OPTIMIZED PLAN
- Tools: PyMOL, VMD (for trajectories), MDAnalysis
- Instance: r5.xlarge (sufficient for your dataset size)
- Storage: EBS optimized for S3 transfer
- Data transfer: ~$0.02 (same region)
- Auto-mount your S3 bucket at /data/

Ready to deploy your custom protein analysis environment?"

User: "Yes, but can you also include AMBER for MD simulations?"
Wizard: "Adding AMBER! This will need a larger instance for MD. Updating to c5.2xlarge. New cost: $0.34/hour. Proceed?"
```

### Option 5: "Multi-Location Data"
```
User: "Genomics analysis - my FASTQ files are on Google Drive and reference genomes in s3://broad-references/"

Wizard: "I'll handle data from multiple sources:
🔄 DATA INTEGRATION PLAN
- Your FASTQ files: Download from Google Drive to instance
- Reference genomes: Direct S3 mount (no transfer cost)
- Working storage: 500GB for analysis outputs
- Backup: Results auto-sync to your S3 bucket

⚠️  Note: Google Drive download may take 15-20 minutes
💡 Tip: For future analyses, I can help migrate your data to S3 for faster access

Proceed with setup?"
```

### Option 6: "Collaborative Data Analysis"
```
User: "Climate modeling with shared data - our team data is at s3://climate-lab/models/ and I need to share results with 3 collaborators"

Wizard: "Setting up collaborative climate analysis:
🌍 COLLABORATION PLAN
- Data source: s3://climate-lab/models/ (read access)
- Compute: c5.4xlarge for parallel climate processing
- Tools: CDO, NCO, Python climate stack
- Shared results: New S3 bucket with team access
- Jupyter Lab: Team-accessible with shared notebooks

👥 TEAM ACCESS
- Send collaborators: [jupyter-url] with team password
- Shared workspace: /shared/ directory
- Auto-backup of notebooks and results

Deploy collaborative climate workspace?"
```

---

## User Control Phrases & Transparency Levels

### Level 0: Pure Magic
```
User: "genomics analysis"
Wizard: "Done! Environment ready."
```

### Level 1: Summary
```
Wizard: "Creating genomics environment with BWA+GATK on memory-optimized instance. Ready in 8 minutes."
```

### Level 2: Plan Review
```
Wizard: "Here's my plan: [detailed breakdown]
- Approve and deploy?
- Modify something?
- Explain why I chose this?"
```

### Level 3: Teaching Mode
```
Wizard: "I'm choosing BWA because your data is whole-genome sequencing. BWA-MEM is fastest for reads >70bp. Alternative aligners would be STAR (for RNA-seq) or minimap2 (for long reads). I'm using r5.xlarge because variant calling needs 32GB+ RAM for human genomes..."
```

## Data Intelligence Examples

### Data Location Understanding
- **S3 buckets**: Direct mounting, cost optimization
- **Google Drive/Dropbox**: Download and transfer strategies
- **FTP/HTTP**: Automated download with retry logic
- **Institutional databases**: API access where available
- **Local uploads**: Optimized transfer paths

### Data Format Recognition
- **FASTQ/BAM files**: → Genomics pipeline
- **Images (.jpg, .tiff, .dicom)**: → Computer vision tools
- **Molecular files (.pdb, .mol2)**: → Structural biology tools
- **Climate data (.nc, .grib)**: → Earth science tools
- **CSV/Excel**: → Statistical analysis tools

### Data Size Optimization
- **Small datasets (<1GB)**: Standard instances
- **Medium datasets (1-100GB)**: Optimized storage
- **Large datasets (>100GB)**: High-throughput instances, parallel processing
- **Streaming data**: Real-time processing instances

---

## Communication Style Adaptation

**Quick Deploy:**
- "Just do it"
- "Don't care about details"
- "Make it work"

**Review First:**
- "Show me the plan"
- "What are you going to set up?"
- "Let me see before deploying"

**Collaborative:**
- "I need GROMACS specifically"
- "Keep it under $20"
- "Use CPU only"
- "I prefer conda over Spack"

**Learning:**
- "Explain your choices"
- "Why that instance type?"
- "Teach me about the tools"

---

## The Magic: Intelligent Data Integration

The wizard should understand:
1. **Where your data lives** (S3, Google Drive, local, FTP, Globus endpoints)
2. **What format it's in** (FASTQ, images, CSV, molecular structures)
3. **How much you have** (affects instance sizing and storage strategy)
4. **Access patterns** (streaming, batch, random access, archive)
5. **Who needs access** (collaboration requirements)
6. **Where results should go** (backup, sharing, publication, long-term archive)
7. **Optimal AWS storage** (S3 tiers, EFS, FSx performance characteristics)
8. **Best transfer tools** (Globus, rsync, rclone, CargoShip, AWS DataSync)

This data intelligence, combined with research domain knowledge, creates true **wizard-level** assistance that handles the complete research workflow from data to results.

---

## Advanced Data Intelligence Examples

### Tool Selection by Use Case

**Large Dataset Transfer (>10TB):**
```
User: "I have 50TB of genomics data at NCBI that I need to analyze"

Wizard: "For 50TB transfer, I recommend:
🚀 PRIMARY: AWS DataSync + Snowball Edge (if >100TB)
📡 ALTERNATIVE: Globus endpoint-to-endpoint (fastest for academic)
💾 STORAGE: S3 Intelligent Tiering → Glacier for long-term
⚡ COMPUTE: FSx Lustre for high-throughput genomics analysis
💰 ESTIMATED: $2,300 storage + $400 transfer costs"
```

**Real-time Collaboration:**
```
User: "My team needs shared access to molecular dynamics trajectories"

Wizard: "For collaborative MD analysis:
📂 STORAGE: EFS with General Purpose mode
🔄 SYNC: rclone bidirectional sync for team laptops
🖥️  ACCESS: Multi-mount across compute instances
⚡ PERFORMANCE: 3GB/s throughput for trajectory reading
👥 TEAM: 5 users × shared /trajectories/ mount"
```

**Archive and Compliance:**
```
User: "Clinical genomics data - need 7-year retention, HIPAA compliance"

Wizard: "Clinical data compliance setup:
🔒 STORAGE: S3 with encryption at rest + in transit
📋 LIFECYCLE: Standard → IA (30d) → Glacier (1yr) → Deep Archive
🛡️  COMPLIANCE: CloudTrail logging, access controls
🔐 TRANSFER: CargoShip for secure, tracked uploads
📊 MONITORING: Data integrity checks, audit logs"
```

### Storage Strategy Intelligence

**Performance-Optimized:**
```yaml
# High-throughput genomics
storage_strategy:
  primary: "fsx_lustre"      # 1GB/s+ throughput
  cache: "ebs_gp3"          # Hot working data
  archive: "s3_standard"    # Results backup
  transfer_tool: "aws_datasync"

reasoning: "BWA alignment needs sustained high I/O for reference genome access"
```

**Cost-Optimized:**
```yaml
# Infrequent access climate data
storage_strategy:
  primary: "s3_intelligent_tiering"  # Auto-optimization
  working: "efs_standard"            # Shared access
  archive: "s3_glacier_deep"         # Long-term models
  transfer_tool: "rclone"            # Free, efficient

reasoning: "Climate models accessed monthly, optimize for storage cost"
```

**Collaboration-Optimized:**
```yaml
# Multi-user protein structure analysis
storage_strategy:
  shared: "efs_general_purpose"      # Multi-mount
  personal: "ebs_gp3"               # Individual workspaces
  results: "s3_standard"            # Publication ready
  transfer_tool: "globus"           # Academic standard

reasoning: "Team needs concurrent access to PDB database and shared results"
```

---

*This document captures the vision for an LLM-powered research wizard that truly bridges the gap between research intent and infrastructure complexity.*
