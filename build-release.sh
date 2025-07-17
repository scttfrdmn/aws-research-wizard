#!/bin/bash

# AWS Research Wizard v0.3.0 Release Build Script
# Usage: ./build-release.sh

set -e

echo "🚀 AWS Research Wizard v0.3.0 Release Build"
echo "=============================================="

# Navigate to Go directory
cd go/

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -f aws-research-wizard-*

# Build for multiple platforms
echo "🔨 Building binaries for multiple platforms..."

# Linux (amd64)
echo "  Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o ../releases/aws-research-wizard-linux-amd64 ./cmd/main.go

# macOS (amd64)
echo "  Building for macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -o ../releases/aws-research-wizard-darwin-amd64 ./cmd/main.go

# macOS (arm64)
echo "  Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -o ../releases/aws-research-wizard-darwin-arm64 ./cmd/main.go

# Windows (amd64)
echo "  Building for Windows amd64..."
GOOS=windows GOARCH=amd64 go build -o ../releases/aws-research-wizard-windows-amd64.exe ./cmd/main.go

# Create release directory if it doesn't exist
mkdir -p ../releases/

# Test the appropriate binary for current platform
echo "🧪 Testing binary for current platform..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    TEST_BINARY="../releases/aws-research-wizard-darwin-amd64"
elif [[ "$OSTYPE" == "linux"* ]]; then
    TEST_BINARY="../releases/aws-research-wizard-linux-amd64"
else
    echo "  ⚠️  Cannot test on this platform, skipping binary test"
    TEST_BINARY=""
fi

if [[ -n "$TEST_BINARY" ]]; then
    if $TEST_BINARY --help > /dev/null 2>&1; then
        echo "  ✅ Binary works correctly"
    else
        echo "  ❌ Binary test failed"
        exit 1
    fi
fi

# Create checksums
echo "🔐 Creating checksums..."
cd ../releases/
sha256sum aws-research-wizard-* > checksums.txt

# Create tar archives
echo "📦 Creating release archives..."
tar -czf aws-research-wizard-v0.3.0-linux-amd64.tar.gz aws-research-wizard-linux-amd64
tar -czf aws-research-wizard-v0.3.0-darwin-amd64.tar.gz aws-research-wizard-darwin-amd64
tar -czf aws-research-wizard-v0.3.0-darwin-arm64.tar.gz aws-research-wizard-darwin-arm64
zip aws-research-wizard-v0.3.0-windows-amd64.zip aws-research-wizard-windows-amd64.exe

echo "✅ Release build complete!"
echo ""
echo "📁 Release files created in ./releases/:"
ls -la

echo ""
echo "🎉 AWS Research Wizard v0.3.0 is ready for release!"
echo "📋 Next steps:"
echo "  1. Test binaries on target platforms"
echo "  2. Create GitHub release with these assets"
echo "  3. Deploy documentation site"
echo "  4. Announce release to community"
