# AWS Research Wizard - Installation Guide

Multiple installation methods are available to get AWS Research Wizard on your system. Choose the method that works best for your environment.

## Quick Start

### 🐧 Linux / macOS (Universal)
```bash
# One-line installation
curl -fsSL https://raw.githubusercontent.com/aws-research-wizard/aws-research-wizard/main/install.sh | sh

# Or download and inspect first
wget https://raw.githubusercontent.com/aws-research-wizard/aws-research-wizard/main/install.sh
chmod +x install.sh
./install.sh
```

### 🍎 macOS (Homebrew)
```bash
# Coming soon - Homebrew formula
brew install aws-research-wizard
```

### 🪟 Windows (Chocolatey)
```powershell
# Coming soon - Chocolatey package
choco install aws-research-wizard
```

## Platform-Specific Installation

### macOS

#### Homebrew (Recommended)
```bash
# Add tap (once formula is published)
brew tap aws-research-wizard/tap
brew install aws-research-wizard

# Update
brew upgrade aws-research-wizard
```

#### Manual Installation
```bash
# Download latest release
curl -fsSL https://github.com/aws-research-wizard/aws-research-wizard/releases/latest/download/aws-research-wizard-darwin-amd64 -o aws-research-wizard

# For Apple Silicon Macs
curl -fsSL https://github.com/aws-research-wizard/aws-research-wizard/releases/latest/download/aws-research-wizard-darwin-arm64 -o aws-research-wizard

# Make executable and install
chmod +x aws-research-wizard
sudo mv aws-research-wizard /usr/local/bin/

# Verify installation
aws-research-wizard version
```

### Linux

#### Universal Installer (Recommended)
```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/aws-research-wizard/aws-research-wizard/main/install.sh | sh

# Install specific version
curl -fsSL https://raw.githubusercontent.com/aws-research-wizard/aws-research-wizard/main/install.sh | VERSION=v1.0.0 sh

# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/aws-research-wizard/aws-research-wizard/main/install.sh | INSTALL_DIR=$HOME/.local/bin sh
```

#### Distribution Packages

**Ubuntu/Debian (Coming Soon):**
```bash
# Add repository
curl -fsSL https://packages.aws-research-wizard.dev/gpg | sudo apt-key add -
echo "deb https://packages.aws-research-wizard.dev/apt stable main" | sudo tee /etc/apt/sources.list.d/aws-research-wizard.list

# Install
sudo apt update
sudo apt install aws-research-wizard
```

**RedHat/CentOS/Fedora (Coming Soon):**
```bash
# Add repository
sudo yum-config-manager --add-repo https://packages.aws-research-wizard.dev/rpm/aws-research-wizard.repo

# Install
sudo yum install aws-research-wizard
```

**Arch Linux (Coming Soon):**
```bash
# Install from AUR
yay -S aws-research-wizard

# Or using makepkg
git clone https://aur.archlinux.org/aws-research-wizard.git
cd aws-research-wizard
makepkg -si
```

#### Manual Installation
```bash
# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

# Download latest release
curl -fsSL "https://github.com/aws-research-wizard/aws-research-wizard/releases/latest/download/aws-research-wizard-linux-$ARCH" -o aws-research-wizard

# Make executable and install
chmod +x aws-research-wizard
sudo mv aws-research-wizard /usr/local/bin/

# Verify installation
aws-research-wizard version
```

### Windows

#### Chocolatey (Coming Soon)
```powershell
# Install
choco install aws-research-wizard

# Update
choco upgrade aws-research-wizard
```

#### Scoop (Coming Soon)
```powershell
# Add bucket
scoop bucket add aws-research-wizard https://github.com/aws-research-wizard/scoop-bucket

# Install
scoop install aws-research-wizard
```

#### WinGet (Coming Soon)
```powershell
# Install
winget install aws-research-wizard
```

#### Manual Installation
1. Go to [Releases](https://github.com/aws-research-wizard/aws-research-wizard/releases/latest)
2. Download `aws-research-wizard-windows-amd64.exe`
3. Rename to `aws-research-wizard.exe`
4. Move to a directory in your PATH (e.g., `C:\Windows\System32` or create `C:\tools\aws-research-wizard\`)
5. Add the directory to your PATH environment variable

**PowerShell:**
```powershell
# Download latest release
Invoke-WebRequest -Uri "https://github.com/aws-research-wizard/aws-research-wizard/releases/latest/download/aws-research-wizard-windows-amd64.exe" -OutFile "aws-research-wizard.exe"

# Move to PATH directory (requires admin)
Move-Item "aws-research-wizard.exe" "C:\Windows\System32\"

# Verify installation
aws-research-wizard version
```

## Alternative Installation Methods

### Docker
```bash
# Pull image
docker pull ghcr.io/aws-research-wizard/aws-research-wizard:latest

# Run with AWS credentials
docker run --rm -v ~/.aws:/root/.aws aws-research-wizard/aws-research-wizard config list

# Create alias for easier use
echo 'alias aws-research-wizard="docker run --rm -v ~/.aws:/root/.aws -v $(pwd):/workspace -w /workspace aws-research-wizard/aws-research-wizard"' >> ~/.bashrc
```

### Go Install
```bash
# Install from source (requires Go 1.21+)
go install github.com/aws-research-wizard/aws-research-wizard/cmd@latest

# Note: Binary will be named 'cmd', you may want to rename it
mv $(go env GOPATH)/bin/cmd $(go env GOPATH)/bin/aws-research-wizard
```

### AppImage (Linux)
```bash
# Download AppImage
wget https://github.com/aws-research-wizard/aws-research-wizard/releases/latest/download/aws-research-wizard.AppImage

# Make executable
chmod +x aws-research-wizard.AppImage

# Run
./aws-research-wizard.AppImage --help

# Optional: Install to system
sudo mv aws-research-wizard.AppImage /usr/local/bin/aws-research-wizard
```

### Snap (Linux)
```bash
# Install from Snap Store (coming soon)
sudo snap install aws-research-wizard

# Run
aws-research-wizard --help
```

## Verification

After installation, verify that AWS Research Wizard is working correctly:

```bash
# Check version
aws-research-wizard version

# Show help
aws-research-wizard --help

# List available commands
aws-research-wizard config --help
aws-research-wizard deploy --help
aws-research-wizard monitor --help
```

Expected output:
```
AWS Research Wizard dev
Built: 2024-06-29_12:34:56
Commit: abc123def456
Go version: go1.21+
```

## Troubleshooting

### Binary Not Found
**Problem:** `aws-research-wizard: command not found`

**Solutions:**
1. Check if the binary is in your PATH:
   ```bash
   echo $PATH
   which aws-research-wizard
   ```

2. If installed to a custom directory, add it to PATH:
   ```bash
   export PATH="/path/to/installation:$PATH"
   # Add to ~/.bashrc or ~/.zshrc for persistence
   ```

3. Use the full path:
   ```bash
   /usr/local/bin/aws-research-wizard --help
   ```

### Permission Denied
**Problem:** `Permission denied` when trying to install

**Solutions:**
1. Use sudo for system-wide installation:
   ```bash
   sudo ./install.sh
   ```

2. Install to user directory:
   ```bash
   INSTALL_DIR=$HOME/.local/bin ./install.sh
   # Ensure ~/.local/bin is in your PATH
   ```

3. Check file permissions:
   ```bash
   ls -la /usr/local/bin/aws-research-wizard
   chmod +x /usr/local/bin/aws-research-wizard
   ```

### Architecture Not Supported
**Problem:** `Unsupported architecture` error

**Available architectures:**
- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64

If your architecture isn't supported, please [open an issue](https://github.com/aws-research-wizard/aws-research-wizard/issues).

### Download Fails
**Problem:** Download fails or times out

**Solutions:**
1. Check internet connection
2. Try manual download from [releases page](https://github.com/aws-research-wizard/aws-research-wizard/releases)
3. Use alternative download method (wget instead of curl)
4. Check if corporate firewall is blocking downloads

### Checksum Verification Fails
**Problem:** Checksum verification fails during installation

**Solutions:**
1. Re-download the binary (file may be corrupted)
2. Skip verification (not recommended):
   ```bash
   # Download installer and modify to skip checksum
   wget https://raw.githubusercontent.com/aws-research-wizard/aws-research-wizard/main/install.sh
   # Edit script to comment out checksum verification
   ```

## Updating

### Automatic Updates
AWS Research Wizard will notify you when updates are available:
```bash
aws-research-wizard version --check
```

### Manual Updates

**Homebrew:**
```bash
brew upgrade aws-research-wizard
```

**Chocolatey:**
```powershell
choco upgrade aws-research-wizard
```

**Universal Installer:**
```bash
# Re-run installer to get latest version
curl -fsSL https://raw.githubusercontent.com/aws-research-wizard/aws-research-wizard/main/install.sh | sh
```

**Manual:**
1. Download new version from releases
2. Replace existing binary
3. Verify new version: `aws-research-wizard version`

## Uninstallation

### Remove Binary
```bash
# Find installation location
which aws-research-wizard

# Remove binary (adjust path as needed)
sudo rm /usr/local/bin/aws-research-wizard
```

### Package Managers

**Homebrew:**
```bash
brew uninstall aws-research-wizard
```

**Chocolatey:**
```powershell
choco uninstall aws-research-wizard
```

**APT (Ubuntu/Debian):**
```bash
sudo apt remove aws-research-wizard
```

### Clean Up Configuration
```bash
# Remove any configuration files (if created)
rm -rf ~/.aws-research-wizard
```

## Getting Help

- **Documentation:** [GitHub Repository](https://github.com/aws-research-wizard/aws-research-wizard)
- **Issues:** [Report a Bug](https://github.com/aws-research-wizard/aws-research-wizard/issues)
- **Discussions:** [Community Discussions](https://github.com/aws-research-wizard/aws-research-wizard/discussions)
- **Built-in Help:** `aws-research-wizard --help`

## Next Steps

After installation, check out:
1. [Quick Start Guide](./README.md#quick-start)
2. [Configuration Guide](./docs/configuration.md)
3. [Deployment Examples](./docs/examples.md)
4. [Monitoring Dashboard](./docs/monitoring.md)