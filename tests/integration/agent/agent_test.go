package agent

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestAgentCreate tests the agent create command
func TestAgentCreate(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: fastapi
  description: "Test agent for unit tests"
  model: gpt-4
  tools:
    - name: test-tool
      description: "A test tool"
`

	tempFile := createTempFile(t, "valid-agent-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "agent", "create", tempFile, "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Agent create command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Creating agents from YAML configuration") {
		t.Errorf("Should show agent creation message, got: %s", outputStr)
	}
}

// TestAgentCreateWithNonExistentFile tests with non-existent file
func TestAgentCreateWithNonExistentFile(t *testing.T) {
	cmd := exec.Command("../../../maestro", "agent", "create", "nonexistent.yaml")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with non-existent file
	if err == nil {
		t.Error("Agent create command should fail with non-existent file")
	}

	if !strings.Contains(outputStr, "no such file or directory") {
		t.Errorf("Error message should mention file not found, got: %s", outputStr)
	}
}

// TestAgentCreateWithInvalidYAML tests with invalid YAML
func TestAgentCreateWithInvalidYAML(t *testing.T) {
	// Create an invalid YAML file
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: "fastapi
  description: "Test agent with invalid YAML"
  model: gpt-4
`

	tempFile := createTempFile(t, "invalid-agent-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "agent", "create", tempFile)
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with invalid YAML
	if err == nil {
		t.Error("Agent create command should fail with invalid YAML")
	}

	if !strings.Contains(outputStr, "no valid YAML documents found") {
		t.Errorf("Error message should mention YAML parsing error, got: %s", outputStr)
	}
}

// TestAgentCreateWithInvalidConfig tests with invalid configuration
func TestAgentCreateWithInvalidConfig(t *testing.T) {
	// Create YAML with invalid configuration
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: invalid-framework
  description: "Test agent with invalid framework"
  model: gpt-4
`

	tempFile := createTempFile(t, "invalid-config-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "agent", "create", tempFile)
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test might fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		// For other errors, check if they mention invalid configuration
		if !strings.Contains(outputStr, "invalid") {
			t.Errorf("Error message should mention invalid configuration, got: %s", outputStr)
		}
	}
}

// TestAgentServe tests the agent serve command
func TestAgentServe(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: fastapi
  description: "Test agent for unit tests"
  model: gpt-4
  tools:
    - name: test-tool
      description: "A test tool"
`

	tempFile := createTempFile(t, "valid-agent-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "agent", "serve", tempFile, "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Agent serve command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Agent server started successfully") {
		t.Errorf("Should show agent serving message, got: %s", outputStr)
	}
}

// TestAgentServeWithCustomPort tests the agent serve command with custom port
func TestAgentServeWithCustomPort(t *testing.T) {
	// Create a valid YAML file for testing
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: fastapi
  description: "Test agent for unit tests"
  model: gpt-4
  tools:
    - name: test-tool
      description: "A test tool"
`

	tempFile := createTempFile(t, "valid-agent-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "agent", "serve", tempFile, "--port=8080", "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Agent serve command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Agent server started successfully") {
		t.Errorf("Should show agent serving message, got: %s", outputStr)
	}
}

// TestAgentServeWithSpecificAgent tests the agent serve command with a specific agent name
func TestAgentServeWithSpecificAgent(t *testing.T) {
	// Create a valid YAML file with multiple agents
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent-1
spec:
  framework: fastapi
  description: "Test agent 1"
  model: gpt-4
---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent-2
spec:
  framework: fastapi
  description: "Test agent 2"
  model: gpt-4
`

	tempFile := createTempFile(t, "valid-agents-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "agent", "serve", tempFile, "--agent-name=test-agent-2", "--dry-run")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// This test is expected to fail if no MCP server is running
	if err != nil {
		// Check if the error is due to MCP server not being available
		if strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running (expected): %s", outputStr)
			return
		}
		t.Fatalf("Agent serve command failed with unexpected error: %v, output: %s", err, string(output))
	}

	if !strings.Contains(outputStr, "Agent server started successfully") {
		t.Errorf("Should show agent serving message, got: %s", outputStr)
	}
}

// TestAgentServeWithNonExistentFile tests with non-existent file
func TestAgentServeWithNonExistentFile(t *testing.T) {
	cmd := exec.Command("../../../maestro", "agent", "serve", "nonexistent.yaml")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with non-existent file
	if err == nil {
		t.Error("Agent serve command should fail with non-existent file")
	}

	if !strings.Contains(outputStr, "no such file or directory") {
		t.Errorf("Error message should mention file not found, got: %s", outputStr)
	}
}

// TestAgentServeWithInvalidYAML tests with invalid YAML
func TestAgentServeWithInvalidYAML(t *testing.T) {
	// Create an invalid YAML file
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: Agent
metadata:
  name: test-agent
spec:
  framework: "fastapi
  description: "Test agent with invalid YAML"
  model: gpt-4
`

	tempFile := createTempFile(t, "invalid-agent-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../../../maestro", "agent", "serve", tempFile)
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Should fail with invalid YAML
	if err == nil {
		t.Error("Agent serve command should fail with invalid YAML")
	}

	if !strings.Contains(outputStr, "no valid YAML documents found") {
		t.Errorf("Error message should mention parsing error, got: %s", outputStr)
	}
}

// TestAgentHelp tests the agent help command
func TestAgentHelp(t *testing.T) {
	cmd := exec.Command("../../../maestro", "agent", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run agent help command: %v", err)
	}

	helpOutput := string(output)

	// Check for expected help content
	expectedContent := []string{
		"agent",
		"create",
		"serve",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(helpOutput, expected) {
			t.Errorf("Help output should contain '%s'", expected)
		}
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
