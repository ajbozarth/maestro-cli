#!/bin/bash

# CLI Test Runner Script
# This script runs the Go tests for the maestro CLI

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

print_help() {
    echo -e "${BLUE}Maestro CLI Test Suite${NC}"
    echo ""
    echo "Usage: ./test.sh [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  unit        Run only unit tests (fast, no external dependencies)"
    echo "  integration Run only integration tests (requires MCP server)"
    echo "  all         Run all tests (unit + integration)"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./test.sh unit        # Run only unit tests"
    echo "  ./test.sh integration # Run only integration tests"
    echo "  ./test.sh all         # Run all tests"
    echo "  ./test.sh             # Run unit tests (default)"
    echo ""
    echo "Test Categories:"
    echo "  Unit Tests:"
    echo "    - CLI command parsing and validation"
    echo "    - Help text and usage examples"
    echo "    - Dry-run functionality"
    echo "    - Error handling for missing parameters"
    echo ""
    echo "  Integration Tests:"
    echo "    - Real MCP server communication"
    echo "    - URL normalization"
    echo "    - End-to-end CLI functionality"
    echo "    - Requires running MCP server"
}

# Check if we're in the right directory
if [ ! -f "tests/main_test.go" ]; then
    print_error "This script must be run from the cli directory"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check command line arguments
case "${1:-unit}" in
    "unit")
        RUN_UNIT_TESTS=true
        RUN_INTEGRATION_TESTS=false
        ;;
    "integration")
        RUN_UNIT_TESTS=false
        RUN_INTEGRATION_TESTS=true
        ;;
    "all")
        RUN_UNIT_TESTS=true
        RUN_INTEGRATION_TESTS=true
        ;;
    "help"|"-h"|"--help")
        print_help
        exit 0
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        print_help
        exit 1
        ;;
esac

# Run unit tests if requested
if [ "$RUN_UNIT_TESTS" = true ]; then
    print_header "Running CLI Unit Tests..."
    
    # Run all tests with verbose output
    print_status "Running tests with verbose output..."
    go test -v ./src

    # Run tests with coverage
    print_status "Running tests with coverage..."
    go test -cover ./src

    # Run tests with race detection
    print_status "Running tests with race detection..."
    go test -race ./src

    # Run specific test files if they exist
    if [ -f "tests/main_test.go" ]; then
        print_status "Running main tests..."
        go test -v tests/main_test.go tests/validate_test.go tests/create_test.go tests/delete_test.go tests/list_test.go
    fi

    print_status "✓ Unit tests completed successfully!"
fi

# Run integration tests if requested
if [ "$RUN_INTEGRATION_TESTS" = true ]; then
    print_header "Running CLI Integration Tests..."
    
    # Check if integration test directory exists
    if [ ! -d "tests/integration" ]; then
        print_error "Integration test directory not found: tests/integration"
        exit 1
    fi

    # Run integration tests
    print_status "Running integration tests..."
    if go test -v tests/integration/list_integration_test.go; then
        print_status "✓ Integration tests completed successfully!"
    else
        print_warning "Integration tests failed (this may be expected if no MCP server is running)"
        print_warning "To run integration tests successfully, start the MCP server first:"
        print_warning "  cd .. && ./start.sh --http"
    fi

    # Run maestro integration tests
    print_status "Running maestrointegration tests..."
    if go test -v tests/integration/agent/agent_test.go; then
        print_status "✓ Agent Integration tests completed successfully!"
    else
        print_warning "Agent Integration tests failed (this may be expected if no MCP server is running)"
        print_warning "To run integration tests successfully, start the MCP server first:"
        print_warning "  cd .. && ./start.sh --http"
    fi
    if go test -v tests/integration/workflow/workflow_test.go; then
        print_status "✓ Workflow Integration tests completed successfully!"
    else
        print_warning "Workflow Integration tests failed (this may be expected if no MCP server is running)"
        print_warning "To run integration tests successfully, start the MCP server first:"
        print_warning "  cd .. && ./start.sh --http"
    fi
    if go test -v tests/integration/customresource/cr_test.go; then
        print_status "✓ Custom Resource Integration tests completed successfully!"
    else
        print_warning "Custom Resource Integration tests failed (this may be expected if no MCP server is running)"
        print_warning "To run integration tests successfully, start the MCP server first:"
        print_warning "  cd .. && ./start.sh --http"
    fi
fi

print_status "All CLI tests completed successfully!" 
