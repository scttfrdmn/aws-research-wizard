#!/bin/bash

# AWS Research Wizard - Test Coverage Enforcement Script
# Ensures 85%+ overall coverage and 80%+ per-file coverage

set -e

MIN_COVERAGE=85
MIN_FILE_COVERAGE=80

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --min-coverage=*)
            MIN_COVERAGE="${1#*=}"
            shift
            ;;
        --min-file-coverage=*)
            MIN_FILE_COVERAGE="${1#*=}"
            shift
            ;;
        *)
            echo "Unknown option $1"
            exit 1
            ;;
    esac
done

echo "🧪 Running Go test coverage analysis..."
echo "📊 Minimum overall coverage: ${MIN_COVERAGE}%"
echo "📁 Minimum per-file coverage: ${MIN_FILE_COVERAGE}%"

# Change to Go directory
cd go || exit 1

# Create coverage directory
mkdir -p coverage

# Run tests with coverage
echo "Running tests with coverage..."
go test -v -race -coverprofile=coverage/coverage.out -covermode=atomic ./... || {
    echo "❌ Tests failed!"
    exit 1
}

# Generate coverage report
go tool cover -func=coverage/coverage.out > coverage/coverage.txt

# Calculate overall coverage
OVERALL_COVERAGE=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo "📊 Overall coverage: ${OVERALL_COVERAGE}%"

# Check overall coverage threshold
if (( $(echo "$OVERALL_COVERAGE < $MIN_COVERAGE" | bc -l) )); then
    echo "❌ Overall coverage ${OVERALL_COVERAGE}% is below minimum ${MIN_COVERAGE}%"
    echo "📋 Coverage by package:"
    go tool cover -func=coverage/coverage.out
    exit 1
fi

# Check per-file coverage
echo "🔍 Checking per-file coverage..."
FAILED_FILES=""
while IFS= read -r line; do
    if [[ $line == *".go:"* ]] && [[ $line != *"total:"* ]]; then
        FILE=$(echo "$line" | awk '{print $1}')
        COVERAGE=$(echo "$line" | awk '{print $3}' | sed 's/%//')

        if (( $(echo "$COVERAGE < $MIN_FILE_COVERAGE" | bc -l) )); then
            echo "⚠️  File $FILE has coverage ${COVERAGE}% (below ${MIN_FILE_COVERAGE}%)"
            FAILED_FILES="$FAILED_FILES\n  - $FILE: ${COVERAGE}%"
        fi
    fi
done < coverage/coverage.txt

if [[ -n "$FAILED_FILES" ]]; then
    echo "❌ Files below minimum coverage threshold:"
    echo -e "$FAILED_FILES"
    echo ""
    echo "💡 To fix low coverage:"
    echo "   1. Add unit tests for uncovered code paths"
    echo "   2. Use 'go test -cover -coverprofile=cp.out ./...' to see coverage details"
    echo "   3. Use 'go tool cover -html=cp.out' to view coverage in browser"
    exit 1
fi

# Generate HTML coverage report
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

echo "✅ All coverage requirements met!"
echo "📊 Overall coverage: ${OVERALL_COVERAGE}%"
echo "📁 All files meet ${MIN_FILE_COVERAGE}% minimum coverage"
echo "📋 HTML report generated: go/coverage/coverage.html"
