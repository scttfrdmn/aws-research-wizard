#!/bin/bash

# AWS Research Wizard v0.3.0 GitHub Release Creation Script
# Usage: ./create-release.sh

set -e

VERSION="0.3.0"
RELEASE_NAME="v${VERSION} - Distribution & Documentation Release"
RELEASE_NOTES="RELEASE_NOTES_v${VERSION}.md"

echo "🚀 Creating GitHub Release for AWS Research Wizard v${VERSION}"
echo "=============================================================="

# Check if release notes exist
if [ ! -f "$RELEASE_NOTES" ]; then
    echo "❌ Release notes file not found: $RELEASE_NOTES"
    exit 1
fi

# Check if releases directory exists
if [ ! -d "releases" ]; then
    echo "❌ Releases directory not found. Please run ./build-release.sh first."
    exit 1
fi

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed. Please install it first:"
    echo "   https://github.com/cli/cli#installation"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Please authenticate with GitHub CLI first:"
    echo "   gh auth login"
    exit 1
fi

# Create the release
echo "📝 Creating GitHub release..."
gh release create "v${VERSION}" \
    --title "$RELEASE_NAME" \
    --notes-file "$RELEASE_NOTES" \
    --draft \
    releases/aws-research-wizard-v${VERSION}-linux-amd64.tar.gz \
    releases/aws-research-wizard-v${VERSION}-darwin-amd64.tar.gz \
    releases/aws-research-wizard-v${VERSION}-darwin-arm64.tar.gz \
    releases/aws-research-wizard-v${VERSION}-windows-amd64.zip \
    releases/checksums.txt

echo "✅ Draft release created successfully!"
echo ""
echo "🎯 Next steps:"
echo "  1. Review the draft release at: https://github.com/$(gh repo view --json owner,name -q '.owner.login + "/" + .name')/releases"
echo "  2. Test the release binaries"
echo "  3. Deploy documentation site"
echo "  4. Publish the release when ready"
echo ""
echo "📋 To publish the release:"
echo "   gh release edit v${VERSION} --draft=false"
echo ""
echo "🌟 To announce the release:"
echo "   - Update social media"
echo "   - Send to mailing lists"
echo "   - Post in relevant forums"
echo "   - Update documentation site"
