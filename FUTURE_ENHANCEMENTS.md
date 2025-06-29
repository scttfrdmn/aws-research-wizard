# AWS Research Wizard - Future Enhancements

**Date**: 2024-06-28  
**Status**: Strategic Planning  
**Priority**: High Impact Improvements

## 1. Resumable Wizard Sessions

### Problem Statement
Currently, when a workload is launched through the wizard, the session is ephemeral. Users lose context if they quit the application and need to restart monitoring from scratch. This creates operational friction for long-running research workloads.

### Proposed Solution: Persistent Session Management

#### Core Concepts
```
Session State = {
    session_id: uuid,
    domain: string,
    deployment_config: DeploymentConfig,
    aws_resources: ResourceTracker,
    status: "launching" | "running" | "completed" | "failed",
    created_at: timestamp,
    last_accessed: timestamp
}
```

#### Implementation Architecture

**Session Storage Options:**
1. **Local State File** (MVP)
   ```bash
   ~/.aws-research-wizard/
   ├── sessions/
   │   ├── session-{uuid}.yaml     # Session metadata
   │   ├── session-{uuid}.state    # Runtime state
   │   └── active_sessions.json    # Session registry
   ```

2. **Cloud State Backend** (Advanced)
   ```bash
   # S3-backed state for multi-device access
   s3://aws-research-wizard-state/{account-id}/
   ├── sessions/{session-id}/
   │   ├── metadata.yaml
   │   ├── deployment.json
   │   └── monitoring.state
   ```

#### User Experience Flow

**Initial Launch:**
```bash
# Start new session
./aws-research-wizard-deploy --domain genomics --instance r6i.4xlarge
> 🚀 Starting deployment session: ses_genomics_20240628_1647
> 📋 Session saved. Resume with: --session ses_genomics_20240628_1647
> ⏳ Deploying CloudFormation stack: research-wizard-genomics...
> ✅ Deployment complete! Resources ready.
> 📊 Monitoring dashboard available. Press 'q' to quit and resume later.
```

**Resume Session:**
```bash
# List active sessions
./aws-research-wizard-deploy sessions
> Active Sessions (2):
> 📊 ses_genomics_20240628_1647  - Running   - Genomics Lab (r6i.4xlarge)
> ⏳ ses_climate_20240627_0930   - Launching - Climate Model (c6i.8xlarge)

# Resume specific session
./aws-research-wizard-monitor --session ses_genomics_20240628_1647
> 📊 Resuming session: Genomics Lab deployment
> 🖥️  Instance: i-0abc123def456789 (running)
> 💰 Current cost: $24.56 (12.5 hours runtime)
> 📈 CPU Usage: 45% | Memory: 67% | Status: Healthy
```

**Session Management:**
```bash
# Session lifecycle commands
./aws-research-wizard-deploy sessions list           # List all sessions
./aws-research-wizard-deploy sessions resume {id}    # Resume session
./aws-research-wizard-deploy sessions archive {id}   # Archive completed
./aws-research-wizard-deploy sessions cleanup        # Remove old sessions
```

#### Technical Implementation

**Session Manager Component:**
```go
type SessionManager struct {
    storageBackend SessionStorage
    activeRegistry map[string]*Session
}

type Session struct {
    ID              string
    Domain          string
    StackName       string
    InstanceIDs     []string
    LaunchTime      time.Time
    LastAccessed    time.Time
    Status          SessionStatus
    Config          *DeploymentConfig
    Resources       *ResourceState
    MonitoringData  *MonitoringHistory
}

func (sm *SessionManager) SaveSession(session *Session) error
func (sm *SessionManager) LoadSession(sessionID string) (*Session, error)
func (sm *SessionManager) ListActiveSessions() ([]*Session, error)
func (sm *SessionManager) CleanupExpiredSessions() error
```

**Integration Points:**
- **Deploy Command**: Auto-create sessions on deployment
- **Monitor Command**: Detect and resume sessions automatically
- **Config Command**: Session-aware domain recommendations
- **Cost Tracking**: Continuous cost accumulation per session

#### Benefits
- **Operational Continuity**: No loss of context when resuming work
- **Multi-device Access**: Resume sessions from different machines
- **Historical Tracking**: Complete deployment and cost history
- **Team Collaboration**: Shared session visibility for teams
- **Long-running Workloads**: Perfect for multi-day research projects

---

## 2. Web-hosted Domain Packs

### Problem Statement
Domain pack configurations are currently bundled with the application, creating versioning issues, distribution challenges, and preventing community contributions. Users get locked into specific domain pack versions.

### Proposed Solution: GitHub-hosted Domain Registry

#### Architecture Overview
```
Domain Pack Registry:
https://github.com/aws-research-wizard/domain-packs/
├── domains/
│   ├── genomics/
│   │   ├── v1.0.0/                 # Semantic versioning
│   │   │   ├── domain.yaml
│   │   │   ├── cloudformation.json
│   │   │   └── README.md
│   │   ├── v1.1.0/
│   │   └── latest -> v1.1.0        # Symlink to latest
│   ├── climate_modeling/
│   ├── machine_learning/
│   └── community/                  # Community contributions
│       ├── neuroscience_advanced/
│       └── custom_hpc/
├── registry.yaml                  # Master registry
├── schemas/                        # Validation schemas
└── tools/                         # Pack validation tools
```

#### Registry Manifest Format
```yaml
# registry.yaml
registry_version: "2.0"
last_updated: "2024-06-28T23:00:00Z"
domains:
  genomics:
    name: "Genomics & Bioinformatics Laboratory"
    description: "Complete genomics analysis with optimized bioinformatics tools"
    maintainer: "aws-research-wizard"
    latest_version: "1.2.0"
    versions:
      - "1.0.0"
      - "1.1.0" 
      - "1.2.0"
    categories: ["biology", "bioinformatics", "research"]
    download_url: "https://github.com/aws-research-wizard/domain-packs/releases/download/genomics-v1.2.0/genomics.yaml"
    checksum: "sha256:abc123..."
    
  climate_modeling:
    name: "Climate Simulation & Atmospheric Modeling"
    maintainer: "climate-research-org"
    latest_version: "2.1.0"
    # ... similar structure
    
community_packs:
  neuroscience_advanced:
    name: "Advanced Neuroscience Pipeline"
    maintainer: "brain-research-lab"
    status: "community"
    trust_level: "verified"  # verified, community, experimental
```

#### Implementation Strategy

**Phase 1: GitHub Integration**
```go
type DomainRegistry struct {
    RegistryURL    string
    CacheDir       string
    UpdateInterval time.Duration
}

func (dr *DomainRegistry) FetchRegistry() (*Registry, error)
func (dr *DomainRegistry) DownloadDomain(name, version string) (*DomainPack, error)
func (dr *DomainRegistry) UpdateCache() error
func (dr *DomainRegistry) ListAvailable() ([]*DomainInfo, error)
```

**Phase 2: Advanced Features**
```go
type DomainManager struct {
    registry        *DomainRegistry
    localCache      *LocalCache
    versionResolver *VersionResolver
}

func (dm *DomainManager) InstallDomain(name, version string) error
func (dm *DomainManager) UpdateDomain(name string) error
func (dm *DomainManager) ValidateDomain(pack *DomainPack) error
func (dm *DomainManager) SearchDomains(query string) ([]*DomainInfo, error)
```

#### User Experience

**Discovery and Installation:**
```bash
# Browse available domains
./aws-research-wizard-config browse
> 📚 Available Domain Packs (15 official, 8 community):
> 
> Official Packs:
>   🧬 genomics v1.2.0        - Genomics & Bioinformatics Laboratory
>   🌍 climate_modeling v2.1.0 - Climate Simulation & Atmospheric Modeling
>   🤖 machine_learning v3.0.0 - ML Research & Training Platform
> 
> Community Packs:
>   🧠 neuroscience_advanced v1.0.0 - Advanced Neuroscience Pipeline (verified)
>   🔬 custom_hpc v0.9.0           - Custom HPC Workflows (community)

# Install specific version
./aws-research-wizard-config install genomics@1.2.0
> 📦 Downloading genomics v1.2.0...
> ✅ Installed genomics v1.2.0
> 🔧 Available for deployment

# Update to latest
./aws-research-wizard-config update genomics
> 📦 Updating genomics v1.2.0 -> v1.3.0...
> ✅ Updated successfully

# Search for domains
./aws-research-wizard-config search neuroscience
> 🔍 Found 3 matching domains:
>   🧠 neuroscience v2.0.0           - Basic Neuroscience Research
>   🧠 neuroscience_advanced v1.0.0  - Advanced Neuroscience Pipeline
>   🧠 neuroscience_imaging v1.1.0   - Brain Imaging Analysis
```

**Version Management:**
```bash
# Pin versions for reproducibility
./aws-research-wizard-config pin genomics@1.2.0
> 📌 Pinned genomics to v1.2.0 (no auto-updates)

# Validate domain pack
./aws-research-wizard-config validate genomics
> ✅ Domain pack validation passed
> 🔒 Signature verified
> 💰 Cost estimates current
> 🏗️  CloudFormation template valid
```

#### Community Contribution Flow

**Domain Pack Development:**
```bash
# Create new domain pack
git clone https://github.com/aws-research-wizard/domain-packs.git
cd domain-packs

# Use template
cp -r templates/domain-template domains/my_research_domain/
cd domains/my_research_domain/

# Develop and test
./tools/validate-domain.sh my_research_domain
./tools/test-deployment.sh my_research_domain

# Submit PR
git add domains/my_research_domain/
git commit -m "feat: Add My Research Domain pack"
git push origin feature/my-research-domain
# Create PR for review
```

**Quality Assurance:**
- **Automated Testing**: CI/CD validates all domain packs
- **Security Scanning**: CloudFormation templates security audited
- **Cost Validation**: Automatic cost estimate validation
- **Community Review**: Peer review process for community packs

#### Benefits

**For Users:**
- **Always Current**: Latest domain configurations without app updates
- **Version Control**: Pin specific versions for reproducibility
- **Community Access**: Leverage community-contributed domains
- **Faster Updates**: Domain improvements deployed immediately

**For Maintainers:**
- **Decoupled Releases**: Domain updates independent of app releases
- **Community Contributions**: Scalable domain pack ecosystem
- **Quality Control**: Centralized validation and security scanning
- **Analytics**: Usage metrics for domain popularity

**For Ecosystem:**
- **Extensibility**: Easy third-party domain pack development
- **Standardization**: Common format for research environment definitions
- **Collaboration**: Shared domain packs across organizations
- **Innovation**: Rapid experimentation with new research domains

---

## Implementation Priority

### High Priority (Phase 3)
1. **Resumable Sessions**: Critical for operational workflows
   - Local file-based session storage (MVP)
   - Session management CLI commands
   - Monitor command session integration

### Medium Priority (Phase 4)
2. **Web-hosted Domain Packs**: Important for scalability
   - GitHub-based domain registry
   - Domain download and caching
   - Version management system

### Integration Strategy
Both enhancements complement the existing architecture and can be implemented incrementally without breaking current functionality. They address real operational pain points and position the project for community growth and enterprise adoption.

### Success Metrics
- **Session Management**: Reduced deployment friction, improved user retention
- **Domain Registry**: Increased domain pack diversity, community contributions
- **Combined Impact**: Enhanced user experience, scalable ecosystem growth