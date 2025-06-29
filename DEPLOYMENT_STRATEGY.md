# AWS Research Wizard - Multi-Platform Deployment Strategy

**Date**: 2024-06-29  
**Status**: Planning & Implementation  
**Target**: Production-ready distribution across all platforms

## Overview

This document outlines the comprehensive deployment strategy for distributing the unified `aws-research-wizard` binary to researchers and institutions across macOS, Linux, and Windows platforms.

## Platform-Specific Deployment Methods

### 🍎 macOS Distribution

#### 1. Homebrew (Primary Method)
**Target**: macOS developers and researchers  
**Priority**: High

```bash
# User installation
brew install aws-research-wizard

# Developer installation  
brew tap aws-research-wizard/tap
brew install aws-research-wizard
```

**Implementation Requirements:**
- Create Homebrew formula in dedicated tap: `aws-research-wizard/homebrew-tap`
- Formula template with binary download from GitHub releases
- Automated formula updates via GitHub Actions
- Support for both Intel and Apple Silicon (universal binary)

**Formula Structure:**
```ruby
class AwsResearchWizard < Formula
  desc "Complete research environment management for AWS"
  homepage "https://github.com/aws-research-wizard/aws-research-wizard"
  version "1.0.0"
  
  if Hardware::CPU.intel?
    url "https://github.com/aws-research-wizard/aws-research-wizard/releases/download/v1.0.0/aws-research-wizard-darwin-amd64"
    sha256 "..."
  else
    url "https://github.com/aws-research-wizard/aws-research-wizard/releases/download/v1.0.0/aws-research-wizard-darwin-arm64"
    sha256 "..."
  end

  def install
    bin.install "aws-research-wizard"
  end

  test do
    assert_match "AWS Research Wizard", shell_output("#{bin}/aws-research-wizard --version")
  end
end
```

#### 2. MacPorts (Secondary)
**Target**: Academic institutions using MacPorts  
**Priority**: Medium

- Submit Portfile to MacPorts registry
- Automatic dependency resolution
- Integration with institutional package management

### 🐧 Linux Distribution

#### 1. Universal Shell Installer (Primary Method)
**Target**: All Linux distributions  
**Priority**: High

```bash
# One-line installation
curl -fsSL https://install.aws-research-wizard.dev | sh

# Manual installation
wget https://install.aws-research-wizard.dev/install.sh
chmod +x install.sh
./install.sh
```

**Installer Features:**
- Auto-detect architecture (amd64, arm64)
- Install to `/usr/local/bin` or `$HOME/.local/bin`
- Add to PATH automatically
- Verify checksums and signatures
- Support for offline installation
- Uninstall capability

**Installer Script Template:**
```bash
#!/bin/bash
set -e

BINARY_NAME="aws-research-wizard"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
GITHUB_REPO="aws-research-wizard/aws-research-wizard"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

# Get latest version
VERSION=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Download and install
URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/$BINARY_NAME-linux-$ARCH"
echo "Downloading $BINARY_NAME $VERSION for linux-$ARCH..."
curl -fsSL "$URL" -o "$BINARY_NAME"
chmod +x "$BINARY_NAME"

# Install
sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
echo "✅ $BINARY_NAME installed to $INSTALL_DIR/$BINARY_NAME"
echo "Run '$BINARY_NAME --help' to get started"
```

#### 2. Distribution-Specific Packages
**Target**: Native package manager users  
**Priority**: Medium

**Debian/Ubuntu (APT):**
```bash
# Add repository
curl -fsSL https://packages.aws-research-wizard.dev/gpg | sudo apt-key add -
echo "deb https://packages.aws-research-wizard.dev/apt stable main" | sudo tee /etc/apt/sources.list.d/aws-research-wizard.list

# Install
sudo apt update
sudo apt install aws-research-wizard
```

**RedHat/CentOS/Fedora (RPM):**
```bash
# Add repository
sudo yum-config-manager --add-repo https://packages.aws-research-wizard.dev/rpm/aws-research-wizard.repo

# Install
sudo yum install aws-research-wizard
```

**Arch Linux (AUR):**
```bash
# Install via AUR helper
yay -S aws-research-wizard
```

#### 3. Universal Packages
**Target**: Containerized and portable environments  
**Priority**: Medium

**Snap Package:**
```bash
sudo snap install aws-research-wizard
```

**AppImage:**
```bash
# Portable, no installation required
wget https://github.com/aws-research-wizard/aws-research-wizard/releases/latest/download/aws-research-wizard.AppImage
chmod +x aws-research-wizard.AppImage
./aws-research-wizard.AppImage
```

### 🪟 Windows Distribution

#### 1. Chocolatey (Primary Method)
**Target**: Windows developers and power users  
**Priority**: High

```powershell
# Installation
choco install aws-research-wizard

# Upgrade
choco upgrade aws-research-wizard
```

**Package Structure:**
```xml
<!-- aws-research-wizard.nuspec -->
<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>aws-research-wizard</id>
    <version>1.0.0</version>
    <title>AWS Research Wizard</title>
    <authors>AWS Research Wizard Team</authors>
    <description>Complete research environment management for AWS</description>
    <projectUrl>https://github.com/aws-research-wizard/aws-research-wizard</projectUrl>
    <tags>aws research cli infrastructure</tags>
  </metadata>
  <files>
    <file src="tools\**" target="tools" />
  </files>
</package>
```

#### 2. Scoop (Developer-Focused)
**Target**: Windows developers  
**Priority**: Medium

```powershell
# Add bucket
scoop bucket add aws-research-wizard https://github.com/aws-research-wizard/scoop-bucket

# Install
scoop install aws-research-wizard
```

#### 3. WinGet (Official Microsoft)
**Target**: Windows 10+ users  
**Priority**: Medium

```powershell
# Installation
winget install aws-research-wizard
```

#### 4. MSI Installer (Enterprise)
**Target**: Enterprise environments  
**Priority**: Medium

- Professional MSI installer with GUI
- Group Policy deployment support
- Silent installation options
- Registry integration for PATH

### 🌐 Universal Distribution Methods

#### 1. GitHub Releases (Foundation)
**Target**: All platforms, direct download  
**Priority**: Critical

**Release Assets:**
```
aws-research-wizard-v1.0.0-linux-amd64
aws-research-wizard-v1.0.0-linux-arm64
aws-research-wizard-v1.0.0-darwin-amd64
aws-research-wizard-v1.0.0-darwin-arm64
aws-research-wizard-v1.0.0-windows-amd64.exe
checksums.txt
checksums.txt.sig
```

#### 2. Docker Container
**Target**: Containerized environments  
**Priority**: Medium

```dockerfile
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY aws-research-wizard /usr/local/bin/
ENTRYPOINT ["aws-research-wizard"]
```

```bash
# Usage
docker run --rm -v ~/.aws:/root/.aws aws-research-wizard/aws-research-wizard config list
```

#### 3. Go Module Installation
**Target**: Go developers  
**Priority**: Low

```bash
go install github.com/aws-research-wizard/aws-research-wizard/cmd@latest
```

## Implementation Roadmap

### Phase 1: Foundation (Week 1)
- ✅ Unified binary consolidation (COMPLETED)
- 🎯 GitHub Actions CI/CD pipeline
- 🎯 Automated release generation
- 🎯 Cross-platform build automation
- 🎯 Release asset signing and verification

### Phase 2: Core Distribution (Week 2)
- 🎯 Universal shell installer for Linux
- 🎯 Homebrew tap and formula
- 🎯 Chocolatey package
- 🎯 GitHub releases with all platforms

### Phase 3: Extended Distribution (Week 3)
- 🎯 Debian/Ubuntu APT packages
- 🎯 RedHat/CentOS RPM packages
- 🎯 Arch Linux AUR package
- 🎯 Scoop bucket

### Phase 4: Enterprise & Containers (Week 4)
- 🎯 Windows MSI installer
- 🎯 Docker container images
- 🎯 Snap package
- 🎯 WinGet manifest

## Security Considerations

### Binary Signing
- **Code Signing**: Sign binaries for Windows and macOS
- **GPG Signatures**: Sign release assets and checksums
- **Reproducible Builds**: Ensure build reproducibility
- **Supply Chain Security**: SLSA compliance

### Distribution Security
- **HTTPS Only**: All downloads over secure connections
- **Checksum Verification**: SHA256 checksums for all binaries
- **Package Repository Security**: Secure package repositories
- **Update Mechanisms**: Secure automatic updates

## Automated Release Pipeline

### GitHub Actions Workflow
```yaml
name: Release
on:
  push:
    tags: ['v*']

jobs:
  build:
    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            goos: linux
            goarch: amd64
          - os: ubuntu-latest
            goos: linux
            goarch: arm64
          - os: macos-latest
            goos: darwin
            goarch: amd64
          - os: macos-latest
            goos: darwin
            goarch: arm64
          - os: windows-latest
            goos: windows
            goarch: amd64
    
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - name: Build binary
        run: |
          GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} \
          go build -ldflags "-X main.version=${{ github.ref_name }}" \
          -o aws-research-wizard-${{ matrix.goos }}-${{ matrix.goarch }} ./cmd
      
      - name: Upload to release
        uses: actions/upload-release-asset@v1
        # ... upload logic
```

### Update Distribution
1. **Tag Release**: `git tag v1.0.0 && git push --tags`
2. **GitHub Actions**: Automatically builds all platform binaries
3. **Package Updates**: Automated PR to package repositories
4. **Container Images**: Automatically builds and pushes Docker images
5. **Verification**: Automated testing of installation methods

## Documentation & Support

### Installation Documentation
- **Quick Start Guide**: One-line installation for each platform
- **Detailed Instructions**: Platform-specific installation guides
- **Troubleshooting**: Common installation issues and solutions
- **Enterprise Deployment**: Guide for institutional deployments

### Update Management
- **Automatic Updates**: Built-in update checker and mechanism
- **Version Management**: Semantic versioning with clear release notes
- **Rollback Support**: Easy downgrade for compatibility issues

## Success Metrics

### Distribution Goals
- **macOS**: 80% via Homebrew, 15% direct download, 5% other
- **Linux**: 60% shell installer, 25% package managers, 15% direct download  
- **Windows**: 70% Chocolatey, 20% direct download, 10% other

### User Experience Targets
- **Installation Time**: < 30 seconds on any platform
- **Zero Dependencies**: No additional software required
- **Cross-Platform Consistency**: Identical experience across platforms
- **Automatic Updates**: Seamless version management

## Future Enhancements

### Advanced Distribution
- **Package Signing**: Enhanced security with package signing
- **Delta Updates**: Efficient incremental updates
- **Enterprise Repository**: Private package repository for organizations
- **Integration Testing**: Automated testing across all distribution methods

### User Experience
- **GUI Installer**: Optional graphical installer for non-technical users
- **Shell Completion**: Auto-completion for all major shells
- **Desktop Integration**: Application launcher integration
- **Documentation Portal**: Comprehensive online documentation

This deployment strategy ensures the AWS Research Wizard reaches researchers across all platforms through their preferred installation methods, while maintaining security, reliability, and ease of use.