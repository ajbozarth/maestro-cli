// Package tool contains integration tests for the tool command
//
// Test Assumptions and Limitations:
//  1. These tests are integration tests that execute the CLI commands as external processes
//  2. Tests use the --dry-run flag when possible to avoid actual resource creation
//  3. Tests are designed to be resilient to different environments (CI, local dev)
//  4. The tests focus on command execution rather than specific output validation
//  5. Some tests may pass even if the underlying functionality has issues, as they
//     primarily test that the command doesn't panic or crash
package tool

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestToolCreate tests the tool create command
func TestToolCreate(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: MCPTool
metadata:
  name: test-tool
spec:
  description: "Test tool for unit tests"
  parameters:
    - name: param1
      description: "A test parameter"
      required: true
      type: string
  returns:
    description: "Test return value"
    type: string
`

	tempFile := createTempFile(t, "valid-tool-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "tool", "create", tempFile, "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Tool create command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Creating MCP tools from YAML configuration") {
		t.Errorf("Should show MCP tools creation message, got: %s", outputStr)
	}
}

// TestToolCreateWithNonExistentFile tests with non-existent file
func TestToolCreateWithNonExistentFile(t *testing.T) {
	cmd := exec.Command("../../../maestro", "tool", "create", "nonexistent.yaml")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with non-existent file
	if err == nil {
		t.Error("Tool create command should fail with non-existent file")
	}

	if !strings.Contains(outputStr, "no such file or directory") {
		t.Errorf("Error message should mention file not found, got: %s", outputStr)
	}
}

// TestToolCreateWithInvalidYAML tests with invalid YAML
func TestToolCreateWithInvalidYAML(t *testing.T) {
	// Create an invalid YAML file
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: MCPTool
metadata:
  name: test-tool
spec:
  description: "Test tool with invalid YAML
  parameters:
    - name: param1
      description: "A test parameter"
      required: true
      type: string
`

	tempFile := createTempFile(t, "invalid-tool-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "tool", "create", tempFile)
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with invalid YAML
	if err == nil {
		t.Error("Tool create command should fail with invalid YAML")
	}

	if !strings.Contains(outputStr, "no valid YAML documents found") {
		t.Errorf("Error message should mention YAML parsing error, got: %s", outputStr)
	}
}

// Helper function to create a temporary file with content
func createTempFile(t *testing.T, pattern string, content string) string {
	tmpfile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}

// Made with Bob
