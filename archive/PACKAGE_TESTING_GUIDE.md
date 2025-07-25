# Package Testing Guide - AWS Research Wizard v0.3.1

## 📦 **Package Distribution Testing**

This guide provides instructions for testing the package distributions before submitting to official repositories.

---

## 🍺 **Homebrew Testing**

### **Setup Test Tap**
```bash
# Create local tap for testing
mkdir -p /usr/local/Homebrew/Library/Taps/scttfrdmn/homebrew-tap
cp packaging/homebrew/aws-research-wizard.rb /usr/local/Homebrew/Library/Taps/scttfrdmn/homebrew-tap/

# Test local formula
brew install --formula /usr/local/Homebrew/Library/Taps/scttfrdmn/homebrew-tap/aws-research-wizard.rb
```

### **Test Installation**
```bash
# Test basic functionality
aws-research-wizard --version
aws-research-wizard --help
aws-research-wizard config --help

# Test dependencies
terraform --version
aws --version
```

### **Cleanup**
```bash
brew uninstall aws-research-wizard
rm -rf /usr/local/Homebrew/Library/Taps/scttfrdmn/homebrew-tap
```

---

## 🍫 **Chocolatey Testing**

### **Setup Test Environment**
```powershell
# Create test package
choco pack packaging/chocolatey/aws-research-wizard.nuspec

# Test local package
choco install aws-research-wizard.0.3.1.nupkg --source .
```

### **Test Installation**
```powershell
# Test basic functionality
aws-research-wizard --version
aws-research-wizard --help
aws-research-wizard config --help

# Test dependencies
aws --version
```

### **Cleanup**
```powershell
choco uninstall aws-research-wizard
Remove-Item aws-research-wizard.0.3.1.nupkg
```

---

## 🐍 **Conda Testing**

### **Setup Test Environment**
```bash
# Create test environment
conda create -n test-aws-research-wizard python=3.9
conda activate test-aws-research-wizard

# Build conda package
conda build packaging/conda/meta.yaml

# Test local package
conda install --use-local aws-research-wizard
```

### **Test Installation**
```bash
# Test basic functionality
aws-research-wizard --version
aws-research-wizard --help
aws-research-wizard config --help

# Test dependencies
terraform --version
aws --version
```

### **Cleanup**
```bash
conda deactivate
conda env remove -n test-aws-research-wizard
```

---

## ✅ **Testing Checklist**

### **All Platforms**
- [ ] **Version Check**: `aws-research-wizard --version` shows v0.3.1
- [ ] **Help Command**: `aws-research-wizard --help` displays usage
- [ ] **Config Command**: `aws-research-wizard config --help` works
- [ ] **Deploy Command**: `aws-research-wizard deploy --help` works
- [ ] **Monitor Command**: `aws-research-wizard monitor --help` works

### **Dependencies**
- [ ] **AWS CLI**: `aws --version` works
- [ ] **Terraform**: `terraform --version` works
- [ ] **Auto-installation**: Dependencies installed automatically

### **Integration Tests**
- [ ] **Config List**: `aws-research-wizard config list` shows domains
- [ ] **Deploy Validate**: `aws-research-wizard deploy validate --domain genomics` works
- [ ] **Error Handling**: Invalid commands show helpful error messages

---

## 🚀 **Submission Process**

### **Homebrew Submission**
1. **Fork homebrew-core**: https://github.com/Homebrew/homebrew-core
2. **Create Formula**: Add `aws-research-wizard.rb` to `Formula/`
3. **Test Formula**: Run `brew install aws-research-wizard` locally
4. **Submit PR**: Create pull request with formula

### **Chocolatey Submission**
1. **Package Ready**: Ensure `.nuspec` and install scripts are correct
2. **Test Package**: Verify local installation works
3. **Submit**: Upload to https://chocolatey.org/packages/upload
4. **Moderation**: Wait for community moderation approval

### **Conda Submission**
1. **Fork conda-forge**: https://github.com/conda-forge/staged-recipes
2. **Create Recipe**: Add `meta.yaml` to `recipes/aws-research-wizard/`
3. **Test Recipe**: Run `conda build` locally
4. **Submit PR**: Create pull request with recipe

---

## 🧪 **Test Results Template**

### **Test Environment**
- **OS**: [macOS/Windows/Linux]
- **Version**: [OS version]
- **Package Manager**: [Homebrew/Chocolatey/Conda]
- **Date**: [Test date]

### **Installation Results**
- **Success**: [Yes/No]
- **Time**: [Installation time]
- **Errors**: [Any errors encountered]

### **Functionality Tests**
- **Version Check**: ✅/❌
- **Help Command**: ✅/❌
- **Config Command**: ✅/❌
- **Deploy Command**: ✅/❌
- **Monitor Command**: ✅/❌

### **Dependencies**
- **AWS CLI**: ✅/❌
- **Terraform**: ✅/❌
- **Auto-install**: ✅/❌

### **Notes**
[Any additional observations or issues]

---

## 📋 **Known Issues**

### **SHA256 Checksums**
- **Status**: Placeholder values in package files
- **Action**: Update with actual checksums from release
- **Priority**: High

### **Package URLs**
- **Status**: Point to v0.3.1 release
- **Action**: Verify all URLs are correct
- **Priority**: High

### **Dependencies**
- **Status**: Some platforms may need manual dependency installation
- **Action**: Test automatic dependency resolution
- **Priority**: Medium

---

## 🔧 **Troubleshooting**

### **Homebrew Issues**
- **Formula not found**: Check tap installation
- **Dependencies missing**: Run `brew install aws-cli terraform`
- **Permission denied**: Check file permissions

### **Chocolatey Issues**
- **Package not found**: Check `.nuspec` file path
- **Admin required**: Run PowerShell as administrator
- **Dependencies missing**: Install manually if needed

### **Conda Issues**
- **Build failed**: Check `meta.yaml` syntax
- **Dependencies missing**: Add to requirements section
- **Environment issues**: Use clean conda environment

---

**Testing Status**: 🧪 **Ready for Testing**

*All package distributions are ready for community testing. Please follow this guide to validate installations before submitting to official repositories.*
