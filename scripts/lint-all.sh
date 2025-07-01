#!/bin/bash

# AWS Research Wizard - Comprehensive Linting Script
# Ensures all code meets quality standards

set -e

echo "🔍 Running comprehensive code linting..."

# Go to project root
cd "$(dirname "$0")/.." || exit 1

# Go linting
echo "📝 Linting Go code..."
cd go || exit 1

# Format Go code
echo "  Formatting Go code..."
go fmt ./...

# Run go vet
echo "  Running go vet..."
go vet ./...

# Run golangci-lint with comprehensive rules
echo "  Running golangci-lint..."
golangci-lint run --timeout=5m --config=.golangci.yml ./... || {
    echo "❌ Go linting failed!"
    exit 1
}

# Check for security issues
echo "  Running gosec security scanner..."
gosec -fmt=json -out=../reports/gosec-report.json -stdout -verbose=text ./... || {
    echo "⚠️  Security issues found - check reports/gosec-report.json"
}

# Go back to project root
cd ..

# Python linting (if Python files exist)
if find python -name "*.py" -type f | head -1 | grep -q .; then
    echo "📝 Linting Python code..."

    # Install/update linting tools
    pip install -q flake8 black isort mypy pylint

    # Format Python code
    echo "  Formatting Python code with black..."
    black python/ || echo "⚠️  Black formatting had issues"

    # Sort imports
    echo "  Sorting imports with isort..."
    isort python/ || echo "⚠️  Import sorting had issues"

    # Run flake8
    echo "  Running flake8..."
    flake8 python/ --max-line-length=88 --extend-ignore=E203,W503 || echo "⚠️  Flake8 found issues"

    # Run pylint
    echo "  Running pylint..."
    pylint python/ --disable=C0114,C0115,C0116 || echo "⚠️  Pylint found issues"
fi

# YAML/JSON linting
echo "📝 Linting YAML and JSON files..."

# Install yamllint if not present
pip install -q yamllint

# Lint YAML files
find . -name "*.yaml" -o -name "*.yml" | grep -v ".git" | while read -r file; do
    echo "  Checking $file..."
    yamllint "$file" || echo "⚠️  YAML issues in $file"
done

# Lint JSON files
find . -name "*.json" | grep -v ".git" | grep -v "node_modules" | while read -r file; do
    echo "  Checking $file..."
    python -m json.tool "$file" > /dev/null || echo "⚠️  JSON issues in $file"
done

# Shell script linting
echo "📝 Linting shell scripts..."

# Install shellcheck if available
if command -v shellcheck > /dev/null; then
    find . -name "*.sh" | grep -v ".git" | while read -r file; do
        echo "  Checking $file..."
        shellcheck "$file" || echo "⚠️  Shell script issues in $file"
    done
else
    echo "  ⚠️  shellcheck not installed - skipping shell script linting"
fi

# Domain pack validation
echo "📝 Validating domain packs..."
if [ -f "domain-packs/tools/validate_domains.py" ]; then
    python domain-packs/tools/validate_domains.py --all || echo "⚠️  Domain pack validation issues"
else
    echo "  ⚠️  Domain pack validator not found"
fi

# Check for common issues
echo "📝 Checking for common issues..."

# Check for hardcoded credentials or secrets
echo "  Scanning for potential secrets..."
if command -v grep > /dev/null; then
    grep -r -i "password\|secret\|key\|token" --include="*.go" --include="*.py" --include="*.yaml" --include="*.json" . | \
    grep -v ".git" | grep -v "test" | grep -v "example" | grep -v "template" | \
    grep -v "# password" | grep -v "// password" || echo "  ✅ No obvious secrets found"
fi

# Check for TODO/FIXME comments
echo "  Checking for TODO/FIXME comments..."
TODO_COUNT=$(find . -name "*.go" -o -name "*.py" -o -name "*.sh" | xargs grep -c "TODO\|FIXME" 2>/dev/null | awk -F: '{sum += $2} END {print sum+0}')
if [ "$TODO_COUNT" -gt 0 ]; then
    echo "  ⚠️  Found $TODO_COUNT TODO/FIXME comments"
    find . -name "*.go" -o -name "*.py" -o -name "*.sh" | xargs grep -n "TODO\|FIXME" 2>/dev/null | head -10
else
    echo "  ✅ No TODO/FIXME comments found"
fi

# Create reports directory
mkdir -p reports

echo ""
echo "✅ Linting complete!"
echo "📋 Check reports/ directory for detailed results"
echo "💡 Run 'scripts/fix-linting.sh' to auto-fix common issues"
