# Project Split-Out Plan: Research Wizard Repository Separation

**Date:** July 2, 2025
**Status:** Ready for Implementation

## Executive Summary

The aws-research-wizard repository has grown to contain three distinct, production-ready projects that should be separated into their own repositories for better maintainability, independent development, and clear ownership.

## 🎯 Projects to Split Out

### 1. Tutorial Guard (Priority: High)
**Current Location:** `/tutorial-guard/`
**Target Repository:** `tutorial-guard`
**Status:** Production-ready with comprehensive documentation

**Key Features:**
- Multi-provider AI integration (Claude, GPT-4, Gemini)
- Enterprise-grade quality certification system
- 8,000+ lines of production Go code
- Complete documentation and testing framework

### 2. Spack Manager Go (Priority: High)
**Current Location:** `/spack-manager-go/`
**Target Repository:** `spack-manager-go`
**Status:** Standalone library, production-ready

**Key Features:**
- Standalone Spack package management library
- TUI interface for package management
- Complete Go module with examples and tests

### 3. AWS Research Wizard (Priority: Medium)
**Current Location:** `/go/` (main project)
**Target Repository:** Keep as `aws-research-wizard`
**Status:** Core research infrastructure platform

## 📋 Implementation Plan

### Phase 1: Pre-Split Preparation

#### A. Documentation Audit and Cleanup
```bash
# Current working directory: /Users/scttfrdmn/src/aws-research-wizard
# Actions needed when resuming:

1. Review and update cross-project references
2. Create final project status documentation
3. Clean up any hanging git references
4. Document dependencies between projects
```

#### B. Dependency Analysis
- **Tutorial Guard**: Independent - no dependencies on research-wizard
- **Spack Manager Go**: Independent - extracted as standalone library
- **Research Wizard**: May reference tutorial-guard for documentation testing

### Phase 2: Project Extraction

#### A. Tutorial Guard Split
```bash
# Commands to run from new location:

# 1. Create new repository
gh repo create tutorial-guard --private --description "AI-Powered Documentation Validation Platform"

# 2. Extract with git history preservation
git subtree push --prefix=tutorial-guard origin tutorial-guard-branch
# OR use git filter-branch for complete history separation

# 3. Clone as new repository
git clone https://github.com/scttfrdmn/tutorial-guard.git
cd tutorial-guard

# 4. Update module paths and imports
find . -name "*.go" -exec sed -i 's|github.com/aws-research-wizard/tutorial-guard|github.com/scttfrdmn/tutorial-guard|g' {} \;
go mod edit -module github.com/scttfrdmn/tutorial-guard
go mod tidy

# 5. Update documentation paths
# Update README.md, docs/, and any internal references

# 6. Verify build and tests
go build ./...
go test ./...

# 7. Initial commit and push
git add .
git commit -m "feat: Initialize tutorial-guard as independent repository"
git push origin main
```

#### B. Spack Manager Go Split
```bash
# Commands to run from new location:

# 1. Create new repository
gh repo create spack-manager-go --private --description "Standalone Go library for Spack package management"

# 2. Extract with git history
git subtree push --prefix=spack-manager-go origin spack-go-branch

# 3. Clone as new repository
git clone https://github.com/scttfrdmn/spack-manager-go.git
cd spack-manager-go

# 4. Update module paths
find . -name "*.go" -exec sed -i 's|github.com/aws-research-wizard/spack-manager-go|github.com/scttfrdmn/spack-manager-go|g' {} \;
go mod edit -module github.com/scttfrdmn/spack-manager-go
go mod tidy

# 5. Verify independence
go build ./...
go test ./...

# 6. Initial commit and push
git add .
git commit -m "feat: Initialize spack-manager-go as independent repository"
git push origin main
```

### Phase 3: Research Wizard Cleanup

#### A. Remove Split Projects
```bash
# From research-wizard repository:

# 1. Remove tutorial-guard directory
git rm -r tutorial-guard/

# 2. Remove spack-manager-go directory
git rm -r spack-manager-go/

# 3. Update documentation references
# - Update README.md to remove tutorial-guard references
# - Update any documentation that cross-references split projects
# - Clean up any import statements that reference moved projects

# 4. Update go.mod if any internal dependencies exist
go mod tidy
```

#### B. Documentation Updates
```bash
# Files to update in research-wizard:

1. README.md - Remove tutorial-guard sections
2. docs/ - Update any cross-references
3. Any integration guides that reference split projects
4. Update CI/CD configs if they test split projects
```

## 🔄 Post-Split Integration

### Repository Relationships
```
┌─────────────────────┐    ┌─────────────────────┐
│   tutorial-guard    │    │  spack-manager-go   │
│  (Independent)      │    │   (Independent)     │
└─────────────────────┘    └─────────────────────┘
                                     │
                                     │ (potential import)
                                     ▼
                            ┌─────────────────────┐
                            │ aws-research-wizard │
                            │  (Core Platform)    │
                            └─────────────────────┘
```

### Cross-Project Integration (Future)
- **Tutorial Guard**: Can be used to validate Research Wizard documentation
- **Spack Manager Go**: Can be imported as library in Research Wizard if needed
- **Research Wizard**: Remains the core infrastructure platform

## 📁 Directory Structure After Split

### Tutorial Guard Repository
```
tutorial-guard/
├── README.md                 # Comprehensive platform documentation
├── PROJECT_COMPLETION.md     # Complete achievement summary
├── LICENSE                   # Proprietary license
├── go.mod                    # Module: github.com/scttfrdmn/tutorial-guard
├── cmd/                      # Command-line applications
├── pkg/                      # Core packages (ai, certification, etc.)
├── docs/                     # Technical documentation
├── examples/                 # Usage examples
└── tests/                    # Test fixtures
```

### Spack Manager Go Repository
```
spack-manager-go/
├── README.md                 # Library documentation
├── LICENSE                   # License file
├── go.mod                    # Module: github.com/scttfrdmn/spack-manager-go
├── cmd/spack-manager/        # CLI application
├── pkg/                      # Core library packages
├── examples/                 # Usage examples
├── docs/                     # API documentation
└── tests/                    # Test suite
```

### AWS Research Wizard Repository (Cleaned)
```
aws-research-wizard/
├── README.md                 # Core platform documentation
├── go/                       # Main Go application
├── python/                   # Python components
├── domain-packs/             # Research domain configurations
├── configs/                  # Configuration files
├── docs/                     # Platform documentation
├── scripts/                  # Utility scripts
└── tests/                    # Test suite
```

## ⚠️ Critical Steps for Resumption

### When You Restart Claude Code:

1. **Change to appropriate directory** (not inside aws-research-wizard)
```bash
cd /Users/scttfrdmn/src/
# This allows creating new repositories alongside research-wizard
```

2. **Have GitHub CLI authenticated**
```bash
gh auth status
# Ensure you can create repositories
```

3. **Prepare for git operations**
```bash
# Ensure clean git state in research-wizard
cd aws-research-wizard
git status
git add . && git commit -m "Pre-split preparation commit"
```

### Information Needed for Split:
- **Repository naming preference**: `tutorial-guard` vs `ai-tutorial-guard`, etc.
- **Visibility**: Private (recommended initially) vs Public
- **License strategy**: Keep proprietary or open source specific projects
- **Module naming**: Preferred Go module paths

## 🎯 Success Criteria

### Tutorial Guard Split Success:
- [ ] Independent repository created with full git history
- [ ] All Go modules build and test successfully
- [ ] Documentation complete and accurate
- [ ] No dependencies on research-wizard
- [ ] CI/CD pipeline functional (if applicable)

### Spack Manager Go Split Success:
- [ ] Independent repository with library focus
- [ ] Clean API documentation
- [ ] Example usage functional
- [ ] No dependencies on research-wizard
- [ ] Importable as external Go module

### Research Wizard Cleanup Success:
- [ ] Split projects cleanly removed
- [ ] No broken references in documentation
- [ ] All remaining functionality intact
- [ ] Clean repository state ready for continued development

## 📞 Next Steps

1. **Exit Claude Code** and navigate to `/Users/scttfrdmn/src/`
2. **Restart Claude Code** from the parent directory
3. **Begin with Tutorial Guard split** (highest value, most complete)
4. **Follow with Spack Manager Go split**
5. **Complete Research Wizard cleanup**
6. **Verify all projects function independently**

This plan ensures clean separation while preserving the substantial work completed in each project, particularly the production-ready Tutorial Guard platform with its comprehensive AI integration and certification system.
