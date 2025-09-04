package main

import (
	"os"
	"os/exec"
	"testing"
)

// TestValidateWithValidYAML tests validation with a valid YAML file
func TestValidateWithValidYAML(t *testing.T) {
	// Create a valid YAML file that matches the vector-database-schema.json
	validYAML := `---
apiVersion: maestro/v1alpha1
kind: VectorDatabase
metadata:
  name: test-milvus
spec:
  type: milvus
  uri: localhost:19530
  collection_name: test_collection
  embedding: text-embedding-3-small
  mode: local
`

	tempFile := createTempFile(t, "valid-*.yaml", validYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "validate", tempFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Validate command failed with valid YAML: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "✅ Validation successful") {
		t.Errorf("Should show success message for valid YAML, got: %s", outputStr)
	}
}

// TestValidateWithInvalidYAML tests validation with an invalid YAML file
func TestValidateWithInvalidYAML(t *testing.T) {
	// Create an invalid YAML file (missing closing quote)
	invalidYAML := `---
apiVersion: maestro/v1alpha1
kind: VectorDatabase
metadata:
  name: test-milvus
spec:
  type: milvus
  uri: localhost:19530
  collection_name: test_collection
  embedding: text-embedding-3-small
  mode: local
  # Missing closing quote - this will cause a YAML parsing error
  api_version: "v1`

	tempFile := createTempFile(t, "invalid-*.yaml", invalidYAML)
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "validate", tempFile)
	output, err := cmd.CombinedOutput()

	// The command should fail (exit code != 0)
	if err == nil {
		t.Errorf("Expected validation to fail for invalid YAML, but it succeeded. Output: %s", string(output))
	}

	// Check that the error message contains expected content
	outputStr := string(output)
	if !contains(outputStr, "invalid yaml format") && !contains(outputStr, "validation failed") {
		t.Errorf("Expected error message about invalid YAML, got: %s", outputStr)
	}
}

// TestValidateWithSchemaFile tests validation with both YAML and schema files
func TestValidateWithSchemaFile(t *testing.T) {
	// Create a simple schema file for testing
	schemaJSON := `{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"version": {"type": "string"},
			"database": {
				"type": "object",
				"properties": {
					"type": {"type": "string"},
					"host": {"type": "string"},
					"port": {"type": "integer"}
				},
				"required": ["type", "host", "port"]
			}
		},
		"required": ["name", "version", "database"]
	}`

	// Create a valid YAML file that matches the simple schema
	validYAML := `---
name: test-config
version: "1.0"
database:
  type: milvus
  host: localhost
  port: 19530
`

	schemaFile := createTempFile(t, "schema-*.json", schemaJSON)
	defer os.Remove(schemaFile)

	yamlFile := createTempFile(t, "config-*.yaml", validYAML)
	defer os.Remove(yamlFile)

	cmd := exec.Command("../maestro", "validate", schemaFile, yamlFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Validate command failed with schema: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "✅ Validation successful") {
		t.Errorf("Should show success message for valid YAML with schema, got: %s", outputStr)
	}
}

// TestValidateWithNonExistentSchema tests validation with non-existent schema file
func TestValidateWithNonExistentSchema(t *testing.T) {
	// Create a valid YAML file
	validYAML := `---
name: test-config
version: 1.0
`

	yamlFile := createTempFile(t, "config-*.yaml", validYAML)
	defer os.Remove(yamlFile)

	cmd := exec.Command("../maestro", "validate", "nonexistent-schema.json", yamlFile)
	output, err := cmd.CombinedOutput()

	// Should fail with non-existent schema file
	if err == nil {
		t.Error("Validate command should fail with non-existent schema file")
	}

	outputStr := string(output)
	if !contains(outputStr, "not found") {
		t.Errorf("Error message should mention 'not found', got: %s", outputStr)
	}
}

// TestValidateWithMultipleFlags tests validation with multiple flags
func TestValidateWithMultipleFlags(t *testing.T) {
	// Create a valid YAML file
	validYAML := `---
name: test-config
version: 1.0
`

	yamlFile := createTempFile(t, "config-*.yaml", validYAML)
	defer os.Remove(yamlFile)

	cmd := exec.Command("../maestro", "validate", "--verbose", "--dry-run", yamlFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Validate command failed with multiple flags: %v, output: %s", err, string(output))
	}

	outputStr := string(output)

	// Should show both verbose and dry-run messages
	if !contains(outputStr, "Validating YAML file") {
		t.Errorf("Verbose flag should show validation message, got: %s", outputStr)
	}

	if !contains(outputStr, "[DRY RUN]") {
		t.Errorf("Dry run flag should show [DRY RUN] message, got: %s", outputStr)
	}
}

// TestValidateArgumentOrder tests the argument order handling
func TestValidateArgumentOrder(t *testing.T) {
	// Create test files
	schemaContent := `{"type": "object"}`
	yamlContent := `---
name: test
`

	schemaFile := createTempFile(t, "schema-*.json", schemaContent)
	defer os.Remove(schemaFile)

	yamlFile := createTempFile(t, "config-*.yaml", yamlContent)
	defer os.Remove(yamlFile)

	// Test both argument orders
	testCases := []struct {
		name       string
		args       []string
		shouldFail bool
	}{
		{
			name:       "YAML first, then schema",
			args:       []string{yamlFile, schemaFile},
			shouldFail: true, // This should fail because it expects schema first
		},
		{
			name:       "Schema first, then YAML",
			args:       []string{schemaFile, yamlFile},
			shouldFail: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("../maestro", append([]string{"validate"}, tc.args...)...)
			output, err := cmd.CombinedOutput()

			outputStr := string(output)

			if tc.shouldFail {
				if err == nil {
					t.Error("Should fail but didn't")
				}
			} else {
				if err != nil {
					t.Errorf("Should succeed but failed: %v, output: %s", err, outputStr)
				}
			}
		})
	}
}

// Helper function to create temporary files with specific content
func createTempFile(t *testing.T, pattern, content string) string {
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
