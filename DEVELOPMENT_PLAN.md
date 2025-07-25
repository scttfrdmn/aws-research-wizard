# AWS Research Wizard: Official Development Plan
*Version 1.0 - July 18, 2025*

## Vision Statement
Transform from static domain templates to an intelligent LLM-powered wizard that understands research intent, data complexity, and infrastructure optimization - bridging the gap between research needs and AWS capabilities.

## Project Philosophy
The AWS Research Wizard should be **magical** - researchers describe their work in natural language, and get working research environments without needing to understand infrastructure. The wizard serves as a universal translator between research intent and AWS services.

---

## 12-Week Development Roadmap

### Phase 1: Foundation Stabilization (Weeks 1-2)
**Objective**: Fix critical issues and establish reliable baseline for 27 domains

#### Week 1: Infrastructure Fixes
- [ ] Fix 11 configuration errors identified in testing
- [ ] Implement domain-specific package strategies (conda-bioconda for genomics, pip for ML)
- [ ] Add domain-specific disk sizing (20GB → 50GB+ for data domains)
- [ ] Update cost estimates with real testing data
- [ ] Optimize user_data.sh with package strategy selection logic

#### Week 2: Domain Validation Framework
- [ ] Create domain validation testing framework
- [ ] Test 5 representative domains end-to-end (genomics, ML, climate, chemistry, visualization)
- [ ] Document package installation success rates and timing
- [ ] Create domain validation dashboard
- [ ] Establish performance metrics baseline

**Success Criteria**:
- 5 domains deploy reliably in <10 minutes
- 90%+ package installation success rate
- Cost estimates within 20% of actual
- Zero manual infrastructure configuration needed

---

### Phase 2: Research Intent Intelligence (Weeks 3-4)
**Objective**: Add natural language understanding to map research descriptions to domains

#### Week 3: Intent Recognition System
- [ ] Build research keyword database for all 27 domains
- [ ] Implement confidence scoring algorithm for domain matching
- [ ] Create domain suggestion explanations
- [ ] Add CLI research intent parsing (`--research "DNA analysis"`)

#### Week 4: Interactive Domain Selection
- [ ] Implement interactive domain selection workflows
- [ ] Add hybrid domain support for cross-disciplinary research
- [ ] Create domain explanation templates
- [ ] Test with diverse research descriptions

**Enhanced CLI Examples**:
```bash
# Current
aws-research-wizard deploy --domain genomics

# Enhanced
aws-research-wizard deploy --research "I want to analyze DNA variants"
# → Suggests genomics domain with explanation and cost estimate
```

**Success Criteria**:
- 80%+ research intent → domain mapping accuracy
- Support for hybrid domain requests
- Interactive domain selection working
- User can deploy without knowing AWS terminology

---

### Phase 3: Data Intelligence Foundation (Weeks 5-6)
**Objective**: Add data source awareness and intelligent storage/transfer recommendations

#### Week 5: Data Source Detection
- [ ] Implement S3 bucket analysis (size, file types, region)
- [ ] Add Google Drive size estimation
- [ ] Create file format detection system
- [ ] Build data size categorization

#### Week 6: Storage Strategy Intelligence
- [ ] Build transfer tool selection logic (Globus, rclone, DataSync, CargoShip)
- [ ] Implement storage strategy recommendations (S3 tiers, EFS, FSx)
- [ ] Add cost estimation for storage options
- [ ] Create performance prediction models

**Data Intelligence Examples**:
```bash
aws-research-wizard deploy \
  --research "protein structure analysis" \
  --data "s3://my-lab/structures/" \
  --team-size 5

# Analyzes data, recommends EFS for collaboration, estimates costs
```

**Success Criteria**:
- Automatic data source analysis (S3, Google Drive, local)
- Intelligent storage strategy recommendations
- Transfer tool selection based on data characteristics
- Cost optimization for data workflows

---

### Phase 4: Conversational Interface (Weeks 7-8)
**Objective**: Implement conversational deployment flows with multiple transparency levels

#### Week 7: Conversation Framework
- [ ] Implement conversation state management
- [ ] Create interaction modes (Magic, Summary, Plan Review, Interactive, Teaching)
- [ ] Add plan explanation generation
- [ ] Build confirmation workflows

#### Week 8: Advanced Conversations
- [ ] Implement constraint parsing and handling (budget, time, compliance)
- [ ] Add intelligent error recovery
- [ ] Create alternative suggestion engine
- [ ] Build tradeoff explanation system

**Conversation Modes**:
- **Magic Mode**: "Just do it" - no questions asked
- **Plan Review**: Show plan, get approval
- **Interactive**: Full conversation with questions
- **Teaching**: Explain all decisions and reasoning

**Success Criteria**:
- Multiple conversation modes working
- Constraint handling (budget, time, compliance)
- Intelligent error recovery
- User satisfaction with interaction experience

---

### Phase 5: LLM Integration (Weeks 9-12)
**Objective**: Add true intelligence with LLM for novel research combinations and expert recommendations

#### Week 9: LLM Foundation
- [ ] Create training dataset from 27 domains
- [ ] Fine-tune LLM on research computing patterns
- [ ] Implement LLM decision engine
- [ ] Add structured output parsing

#### Week 10: Hybrid Domain Intelligence
- [ ] Implement hybrid domain analysis for cross-disciplinary research
- [ ] Add adaptive recommendation engine
- [ ] Create user feedback learning system
- [ ] Build institutional context awareness

#### Week 11: Advanced Data Intelligence
- [ ] Implement LLM data strategy generation
- [ ] Add data lifecycle planning
- [ ] Create cost optimization recommendations
- [ ] Build collaboration pattern recognition

#### Week 12: Production Integration
- [ ] End-to-end LLM integration testing
- [ ] Performance optimization for LLM calls
- [ ] Fallback to rule-based systems
- [ ] Production monitoring and logging
- [ ] User feedback collection and learning

**LLM Capabilities**:
- Handle novel research combinations ("AI-driven protein folding with MD validation")
- Intelligent data strategy generation for complex workflows
- Adaptive learning from user feedback and deployment success
- Expert-level recommendations with detailed reasoning

**Success Criteria**:
- LLM handles novel research combinations
- Adaptive recommendations based on context
- Expert-level data strategy generation
- Production deployment success rate >95%

---

## Data Intelligence Strategy

### Transfer Tool Selection Matrix
- **Globus**: Academic endpoints, large datasets, reliable transfer
- **rclone**: Multi-cloud, scripting, cost-free
- **CargoShip**: Enterprise tracking, compliance, audit trails
- **AWS DataSync**: AWS-native, scheduled sync, bandwidth throttling

### Storage Strategy Framework
- **FSx Lustre**: High-throughput genomics, AI training (1GB/s+)
- **EFS**: Shared notebooks, team collaboration, multi-mount
- **S3 Intelligent Tiering**: Archive workflows, cost optimization
- **EBS GP3**: Individual workspaces, working data

### Data Size Decision Matrix
- **<1GB**: Direct download, local storage
- **1-100GB**: rclone/DataSync, S3 Standard/EFS
- **>100GB**: Globus/Snowball, FSx Lustre/S3 Intelligent Tiering

---

## Technology Stack

### Current Architecture
- **CLI**: Go-based aws-research-wizard
- **Infrastructure**: Terraform with AWS provider
- **Configuration**: YAML domain configs (27 domains)
- **Deployment**: EC2 instances with user_data scripts

### Enhanced Architecture
- **Intelligence Layer**: LLM integration for research understanding
- **Data Layer**: S3/EFS/FSx analysis and optimization
- **Conversation Layer**: Multi-mode interaction system
- **Monitoring Layer**: Deployment success tracking and learning

---

## Success Metrics

### Overall KPIs
- **Deployment Success Rate**: >95%
- **Time to Working Environment**: <10 minutes average
- **Cost Accuracy**: Within 20% of estimates
- **User Satisfaction**: >4.5/5 rating
- **Domain Coverage**: All 27 domains reliably working

### Phase-Specific Metrics
- **Phase 1**: 5 domains working, 90% package success
- **Phase 2**: 80% intent mapping accuracy
- **Phase 3**: Data source analysis working
- **Phase 4**: Multi-mode conversations functional
- **Phase 5**: Novel research combination handling

---

## Risk Mitigation

### Technical Risks
- **LLM Performance**: Maintain rule-based fallbacks
- **AWS Service Limits**: Implement quota monitoring
- **Data Transfer Failures**: Build retry logic and alternatives
- **Cost Overruns**: Implement cost monitoring and warnings

### User Experience Risks
- **Complexity Creep**: Always maintain "magic mode" option
- **Over-Engineering**: Validate against real researcher workflows
- **Performance**: Ensure <30 second response times

### Project Risks
- **Scope Expansion**: Focus on 27 domains excellence vs expansion
- **Domain Knowledge Gaps**: Engage domain experts for validation
- **LLM Dependency**: Ensure critical functions work without LLM

---

## Current Status
- **Genomics Domain**: Tested, optimized with Miniforge (2min vs 90min setup)
- **ML Domain**: Tested, optimized with pip (dramatic speed improvement)
- **Configuration Issues**: 11 critical issues identified and documented
- **Performance Baseline**: Established deployment timing and cost metrics

---

## Next Steps (Immediate)
1. Start Phase 1 Week 1: Fix the 11 configuration errors
2. Implement domain-specific package strategies
3. Test 3 more domains end-to-end for validation
4. Begin planning research intent recognition system

---

*This plan represents the official roadmap for AWS Research Wizard development. All previous planning documents are superseded by this comprehensive plan.*
