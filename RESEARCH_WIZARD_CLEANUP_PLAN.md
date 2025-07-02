# Research Wizard Cleanup Plan: Post-Split Repository Preparation

**Date:** July 2, 2025
**Status:** Ready for Implementation

## 🎯 Summary

This document outlines the specific cleanup tasks needed for the aws-research-wizard repository after Tutorial Guard and Spack Manager Go are split into their own repositories.

## 🔍 Critical Dependencies Found

### 1. Spack Manager Go Dependency in Main Project
**Location:** `/go/go.mod` line 66
**Issue:** `replace github.com/spack-go/spack-manager => ../spack-manager-go`
**Impact:** The main research-wizard Go project currently depends on the local spack-manager-go directory

### 2. Documentation Cross-References
**Files with References:**
- `README.md` - Line 24: `[spack-manager-go](spack-manager-go/)`
- Multiple phase documentation files reference both projects
- Changelog entries mention both projects

## 📋 Cleanup Tasks by Priority

### Priority 1: Critical Dependencies (Must Fix Before Split)

#### A. Resolve Go Module Dependency
```bash
# In go/go.mod, need to either:
# Option 1: Remove dependency if not actually used
# Option 2: Replace with published module after split
# Option 3: Vendor the required functionality

# Check actual usage:
grep -r "spack-manager" go/internal/ go/cmd/
```

#### B. Verify Compilation After Dependency Removal
```bash
cd go/
go mod tidy
go build ./cmd/main.go
go test ./...
```

### Priority 2: Documentation Updates

#### A. Main README.md Updates
**Lines to modify:**
- Line 24: Remove reference to `spack-manager-go/` directory
- Line 24: Update to reference published module or remove entirely
- Any other spack-manager-go directory references

#### B. Documentation File Updates
**Files requiring updates:**
```
/CHANGELOG.md
/PHASE_2_FINAL_RELEASE.md
/PHASE_2_PROJECT_STATUS.md
/go/PHASE_2_PRODUCTION_RELEASE.md
/go/PHASE_2_FINAL_DOCUMENTATION.md
/go/PHASE_2C_COMPLETION_REPORT.md
/go/PHASE_2_COMPLETE_SUMMARY.md
```

**Actions:**
- Update references from local directory to external repository
- Modify installation instructions
- Update integration examples

### Priority 3: Directory Structure Cleanup

#### A. Remove Split Directories
```bash
# After projects are successfully split:
git rm -r tutorial-guard/
git rm -r spack-manager-go/
```

#### B. Update .gitignore if needed
- Remove any ignore patterns specific to split projects
- Ensure no orphaned references remain

## 🔧 Detailed Cleanup Steps

### Step 1: Analyze Spack Manager Go Usage
```bash
# Commands to run:
cd /Users/scttfrdmn/src/aws-research-wizard/go

# Find all imports of spack-manager
grep -r "spack-manager" . --include="*.go"

# Find all imports of the replaced module
grep -r "github.com/spack-go/spack-manager" . --include="*.go"

# Check if functionality can be replaced or removed
```

### Step 2: Resolve Dependency Strategy
**Three Options:**

#### Option A: Remove Dependency (Preferred if minimal usage)
```bash
# Remove from go.mod
sed -i '/github.com\/spack-go\/spack-manager/d' go.mod
sed -i '/replace github.com\/spack-go\/spack-manager/d' go.mod

# Remove imports and replace with alternative implementation
# Update code to use alternative approaches
```

#### Option B: Use Published Module (If functionality needed)
```bash
# After spack-manager-go is published as independent module:
go mod edit -replace=github.com/spack-go/spack-manager=github.com/scttfrdmn/spack-manager-go@latest
go mod tidy
```

#### Option C: Vendor Required Code (If small usage)
```bash
# Copy only needed functionality into internal package
# Remove external dependency
```

### Step 3: Update Documentation References
```bash
# Update README.md
sed -i 's|\[spack-manager-go\](spack-manager-go/)|[spack-manager-go](https://github.com/scttfrdmn/spack-manager-go)|g' README.md

# Update other documentation files
find . -name "*.md" -exec sed -i 's|spack-manager-go/|https://github.com/scttfrdmn/spack-manager-go/|g' {} \;
```

### Step 4: Remove Project Directories
```bash
# Only after successful split to independent repositories:
git rm -r tutorial-guard/
git rm -r spack-manager-go/
```

## 📊 Validation Checklist

### Pre-Split Validation
- [ ] Identify all cross-dependencies between projects
- [ ] Determine resolution strategy for spack-manager-go dependency
- [ ] Test compilation with dependency resolution
- [ ] Document all files requiring updates

### Post-Cleanup Validation
- [ ] All Go code compiles successfully: `go build ./...`
- [ ] All tests pass: `go test ./...`
- [ ] No broken documentation links
- [ ] No references to removed directories
- [ ] Clean git status with no orphaned files

### Final Integration Test
- [ ] Deploy a research environment to verify functionality
- [ ] Test all CLI commands work correctly
- [ ] Verify spack integration still functions (if kept)
- [ ] Confirm no regression in core research-wizard features

## 🔄 Migration Strategy for Spack Integration

### Current Integration Analysis
Based on go.mod, the research-wizard project imports spack-manager functionality. The split needs to handle this gracefully:

**Option 1: Full Separation**
- Remove spack integration from core research-wizard
- Users install spack-manager-go separately if needed
- Update documentation to reflect separate installation

**Option 2: External Module Dependency**
- Publish spack-manager-go as independent module
- Update research-wizard to import published module
- Maintain integration but with external dependency

**Option 3: Focused Integration**
- Keep minimal spack integration in research-wizard
- Remove dependency on separate spack-manager-go
- Focus research-wizard on core infrastructure automation

## 📁 Final Repository Structure

### AWS Research Wizard (Post-Cleanup)
```
aws-research-wizard/
├── README.md                 # Updated without split project references
├── go/                       # Core research infrastructure platform
│   ├── go.mod               # Clean dependencies, no local replace directives
│   ├── internal/            # Core functionality
│   ├── cmd/                 # CLI applications
│   └── docs/                # Platform documentation
├── python/                  # Python components
├── domain-packs/            # Research domain configurations
├── configs/                 # Configuration files
├── docs/                    # Website documentation
└── tests/                   # Test suite
```

## 🚨 Critical Warnings

### 1. Dependency Resolution Must Complete First
**Do not split projects until spack-manager-go dependency is resolved.** The research-wizard build will fail if the directory is removed while the dependency exists.

### 2. Test Thoroughly After Cleanup
Core research-wizard functionality must remain intact. The spack integration was a major feature that needs careful handling during separation.

### 3. Documentation Accuracy Critical
Many documentation files reference the evolution through different phases. Ensure updates maintain historical accuracy while removing current dependencies.

## 📞 Implementation Notes for Resumption

### When You Restart Claude Code:

1. **Start from parent directory:** `/Users/scttfrdmn/src/`
2. **Work on projects in this order:**
   - Analyze spack-manager-go usage first
   - Resolve dependency before any splits
   - Split tutorial-guard (independent)
   - Split spack-manager-go (with dependency resolution)
   - Clean up research-wizard

3. **Have ready:**
   - GitHub CLI authenticated (`gh auth status`)
   - Clean git state in research-wizard
   - Plan for dependency resolution strategy

This cleanup plan ensures the aws-research-wizard repository remains functional and clean after splitting out the two independent projects, while preserving all the substantial research infrastructure capabilities that have been built.
