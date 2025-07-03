#!/bin/bash
# AWS Research Wizard - Universal Installation Script
# This script installs the latest version of aws-research-wizard

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="aws-research-wizard"
GITHUB_REPO="aws-research-wizard/aws-research-wizard"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
TEMP_DIR=$(mktemp -d)
VERSION="${VERSION:-latest}"

# Cleanup function
cleanup() {
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Utility functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect system architecture
detect_arch() {
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l)
            echo "arm"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Detect operating system
detect_os() {
    local os
    os=$(uname -s)
    case $os in
        Linux)
            echo "linux"
            ;;
        Darwin)
            echo "darwin"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}

# Get latest version from GitHub API
get_latest_version() {
    if ! command_exists curl && ! command_exists wget; then
        log_error "Either curl or wget is required to download aws-research-wizard"
        exit 1
    fi

    local latest_url="https://api.github.com/repos/$GITHUB_REPO/releases/latest"
    local version

    if command_exists curl; then
        version=$(curl -s "$latest_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        version=$(wget -qO- "$latest_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi

    if [ -z "$version" ]; then
        log_error "Failed to get latest version from GitHub API"
        exit 1
    fi

    echo "$version"
}

# Download file
download_file() {
    local url="$1"
    local output="$2"

    log_info "Downloading from $url"

    if command_exists curl; then
        curl -fsSL "$url" -o "$output"
    elif command_exists wget; then
        wget -q "$url" -O "$output"
    else
        log_error "Either curl or wget is required to download aws-research-wizard"
        exit 1
    fi
}

# Verify checksum
verify_checksum() {
    local file="$1"
    local expected_checksum="$2"

    if ! command_exists sha256sum; then
        log_warning "sha256sum not available, skipping checksum verification"
        return 0
    fi

    local actual_checksum
    actual_checksum=$(sha256sum "$file" | cut -d' ' -f1)

    if [ "$actual_checksum" = "$expected_checksum" ]; then
        log_success "Checksum verification passed"
        return 0
    else
        log_error "Checksum verification failed"
        log_error "Expected: $expected_checksum"
        log_error "Actual:   $actual_checksum"
        return 1
    fi
}

# Check if installation directory is writable
check_install_dir() {
    if [ -w "$INSTALL_DIR" ]; then
        return 0
    elif [ "$EUID" -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Install binary
install_binary() {
    local binary_path="$1"
    local install_path="$INSTALL_DIR/$BINARY_NAME"

    log_step "Installing $BINARY_NAME to $install_path"

    if check_install_dir; then
        cp "$binary_path" "$install_path"
        chmod +x "$install_path"
    else
        log_info "Requesting administrator privileges to install to $INSTALL_DIR"
        sudo cp "$binary_path" "$install_path"
        sudo chmod +x "$install_path"
    fi

    log_success "Successfully installed $BINARY_NAME to $install_path"
}

# Add to PATH
add_to_path() {
    local shell_rc
    local path_line="export PATH=\"$INSTALL_DIR:\$PATH\""

    # Detect shell and set appropriate RC file
    if [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ "$SHELL" = "/bin/bash" ]; then
        shell_rc="$HOME/.bashrc"
    elif [ "$SHELL" = "/bin/zsh" ]; then
        shell_rc="$HOME/.zshrc"
    else
        shell_rc="$HOME/.profile"
    fi

    # Check if INSTALL_DIR is already in PATH
    if echo "$PATH" | grep -q "$INSTALL_DIR"; then
        log_info "$INSTALL_DIR is already in PATH"
        return 0
    fi

    # Check if we already added it to the RC file
    if [ -f "$shell_rc" ] && grep -q "$INSTALL_DIR" "$shell_rc"; then
        log_info "$INSTALL_DIR already added to $shell_rc"
        return 0
    fi

    # Add to PATH in shell RC file
    if [ -w "$shell_rc" ] || [ ! -f "$shell_rc" ]; then
        echo "" >> "$shell_rc"
        echo "# Added by aws-research-wizard installer" >> "$shell_rc"
        echo "$path_line" >> "$shell_rc"
        log_success "Added $INSTALL_DIR to PATH in $shell_rc"
        log_info "Restart your shell or run: source $shell_rc"
    else
        log_warning "Could not add $INSTALL_DIR to PATH automatically"
        log_info "Please add this line to your shell configuration file:"
        log_info "$path_line"
    fi
}

# Verify installation
verify_installation() {
    local install_path="$INSTALL_DIR/$BINARY_NAME"

    if [ -x "$install_path" ]; then
        log_success "Installation verified: $install_path"

        # Test if binary works
        if "$install_path" version >/dev/null 2>&1; then
            local version_output
            version_output=$("$install_path" version 2>/dev/null | head -n1)
            log_success "Binary is functional: $version_output"
        else
            log_warning "Binary installed but may not be functional"
        fi

        return 0
    else
        log_error "Installation verification failed: $install_path not found or not executable"
        return 1
    fi
}

# Print usage instructions
print_usage() {
    echo ""
    log_success "🎉 aws-research-wizard has been successfully installed!"
    echo ""
    echo -e "${CYAN}Getting Started:${NC}"
    echo "  $BINARY_NAME --help                    # Show all available commands"
    echo "  $BINARY_NAME config list              # List available research domains"
    echo "  $BINARY_NAME deploy --help            # Deploy infrastructure help"
    echo "  $BINARY_NAME monitor --help           # Monitoring dashboard help"
    echo ""
    echo -e "${CYAN}Quick Start Example:${NC}"
    echo "  $BINARY_NAME config                   # Interactive domain configuration"
    echo "  $BINARY_NAME deploy --domain genomics # Deploy a genomics environment"
    echo "  $BINARY_NAME monitor                  # Launch monitoring dashboard"
    echo ""
    echo -e "${CYAN}Documentation:${NC}"
    echo "  https://github.com/$GITHUB_REPO"
    echo ""

    # Check if binary is in PATH
    if command_exists "$BINARY_NAME"; then
        log_success "$BINARY_NAME is ready to use!"
    else
        log_warning "$BINARY_NAME is not in your PATH"
        log_info "Either restart your shell or run: export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

# Main installation function
main() {
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                AWS Research Wizard Installer                 ║"
    echo "║          Complete research environment management             ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"

    # Detect system
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)

    log_info "Detected system: $os-$arch"
    log_info "Install directory: $INSTALL_DIR"

    # Get version
    local version
    if [ "$VERSION" = "latest" ]; then
        log_step "Getting latest version from GitHub"
        version=$(get_latest_version)
    else
        version="$VERSION"
    fi

    log_info "Installing version: $version"

    # Construct download URL
    local binary_filename="$BINARY_NAME-$os-$arch"
    local download_url="https://github.com/$GITHUB_REPO/releases/download/$version/$binary_filename"
    local checksum_url="https://github.com/$GITHUB_REPO/releases/download/$version/checksums.txt"

    # Download binary
    log_step "Downloading $BINARY_NAME binary"
    local binary_path="$TEMP_DIR/$binary_filename"
    download_file "$download_url" "$binary_path"

    # Download and verify checksum
    log_step "Verifying checksum"
    local checksums_path="$TEMP_DIR/checksums.txt"
    if download_file "$checksum_url" "$checksums_path" 2>/dev/null; then
        local expected_checksum
        expected_checksum=$(grep "$binary_filename" "$checksums_path" | cut -d' ' -f1)
        if [ -n "$expected_checksum" ]; then
            verify_checksum "$binary_path" "$expected_checksum"
        else
            log_warning "Checksum not found for $binary_filename, skipping verification"
        fi
    else
        log_warning "Could not download checksums, skipping verification"
    fi

    # Make binary executable
    chmod +x "$binary_path"

    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        log_step "Creating install directory: $INSTALL_DIR"
        if check_install_dir; then
            mkdir -p "$INSTALL_DIR"
        else
            sudo mkdir -p "$INSTALL_DIR"
        fi
    fi

    # Install binary
    install_binary "$binary_path"

    # Add to PATH (only for user installs)
    if [ "$INSTALL_DIR" != "/usr/local/bin" ] && [ "$INSTALL_DIR" != "/usr/bin" ]; then
        add_to_path
    fi

    # Verify installation
    verify_installation

    # Print usage instructions
    print_usage
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --help)
            echo "AWS Research Wizard Installer"
            echo ""
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --version VERSION     Install specific version (default: latest)"
            echo "  --install-dir DIR     Installation directory (default: /usr/local/bin)"
            echo "  --help               Show this help message"
            echo ""
            echo "Environment Variables:"
            echo "  VERSION              Version to install (default: latest)"
            echo "  INSTALL_DIR          Installation directory (default: /usr/local/bin)"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Install latest version"
            echo "  $0 --version v1.0.0                  # Install specific version"
            echo "  $0 --install-dir \$HOME/.local/bin    # Install to user directory"
            echo "  INSTALL_DIR=\$HOME/bin $0             # Install using environment variable"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            log_info "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Run main installation
main
