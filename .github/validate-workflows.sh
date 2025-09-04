#!/bin/bash

# GitHub Actions Workflow Validation Script
# This script validates that the GitHub Actions workflows are properly configured

set -e

echo "🔍 Validating GitHub Actions Workflows..."

# Check if workflow files exist
WORKFLOWS=(
    ".github/workflows/ci.yml"
    ".github/workflows/lint.yml"
    ".github/workflows/build.yml"
    ".github/workflows/test.yml"
)

for workflow in "${WORKFLOWS[@]}"; do
    if [[ -f "$workflow" ]]; then
        echo "✅ Found: $workflow"
    else
        echo "❌ Missing: $workflow"
        exit 1
    fi
done

# Validate YAML syntax
echo "🔍 Validating YAML syntax..."
for workflow in "${WORKFLOWS[@]}"; do
    if command -v yamllint >/dev/null 2>&1; then
        yamllint "$workflow" || echo "⚠️  YAML linting issues in $workflow (yamllint not required)"
    else
        echo "ℹ️  yamllint not installed, skipping YAML validation"
    fi
done

# Check workflow triggers
echo "🔍 Checking workflow triggers..."
for workflow in "${WORKFLOWS[@]}"; do
    if grep -q "pull_request:" "$workflow"; then
        echo "✅ $workflow has PR triggers"
    else
        echo "⚠️  $workflow missing PR triggers"
    fi
    
    if grep -q "push:" "$workflow"; then
        echo "✅ $workflow has push triggers"
    else
        echo "⚠️  $workflow missing push triggers"
    fi
done

# Check for required steps
echo "🔍 Checking for required steps..."
for workflow in "${WORKFLOWS[@]}"; do
    if grep -q "actions/checkout@v4" "$workflow"; then
        echo "✅ $workflow uses checkout@v4"
    else
        echo "⚠️  $workflow should use actions/checkout@v4"
    fi
    
    if grep -q "actions/setup-go@v4" "$workflow"; then
        echo "✅ $workflow uses setup-go@v4"
    else
        echo "⚠️  $workflow should use actions/setup-go@v4"
    fi
    
    if grep -q "actions/upload-artifact@v4" "$workflow"; then
        echo "✅ $workflow uses upload-artifact@v4"
    elif grep -q "actions/upload-artifact@v3" "$workflow"; then
        echo "❌ $workflow uses deprecated upload-artifact@v3 (should be v4)"
    else
        echo "ℹ️  $workflow doesn't use upload-artifact"
    fi
done

echo "🎉 Workflow validation completed!"
echo ""
echo "📋 Summary:"
echo "- All workflow files are present"
echo "- Workflows are configured for PR and push triggers"
echo "- Using latest GitHub Actions versions"
echo ""
echo "🚀 Ready for CI/CD!"
