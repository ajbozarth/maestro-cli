package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestRetrieveCollection(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "retrieve collection with valid database",
			args:        []string{"collection", "list", "--vdb=test_local_milvus", "--dry-run"},
			expectError: false,
		},
		{
			name:        "retrieve collection with non-existent database (dry-run succeeds)",
			args:        []string{"collection", "list", "--vdb=non_existent_db", "--dry-run"},
			expectError: false, // Dry-run doesn't validate existence
		},
		{
			name:        "retrieve collection with specific collection name",
			args:        []string{"collection", "list", "--vdb=test_local_milvus", "--dry-run"},
			expectError: false,
		},
		{
			name:        "retrieve collection with invalid resource type",
			args:        []string{"collection", "invalid-action", "--vdb=test_local_milvus", "--dry-run"},
			expectError: true,
			errorMsg:    "Available Commands",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("../maestro", tc.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
				if tc.errorMsg != "" && !strings.Contains(outputStr, tc.errorMsg) {
					t.Errorf("Expected error message containing '%s', got: %s", tc.errorMsg, outputStr)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v, output: %s", err, outputStr)
				}
			}
		})
	}
}

func TestRetrieveDocument(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "retrieve document with missing collection name",
			args:        []string{"document", "list", "--vdb=test_local_milvus", "--dry-run"},
			expectError: true,
			errorMsg:    "collection name is required in non-interactive mode",
		},
		{
			name:        "retrieve document with non-existent database (dry-run succeeds)",
			args:        []string{"document", "list", "--vdb=non_existent_db", "--collection=test_collection", "--dry-run"},
			expectError: false, // Dry-run doesn't validate existence
		},
		{
			name:        "retrieve document with non-existent collection (dry-run succeeds)",
			args:        []string{"document", "list", "--vdb=test_local_milvus", "--collection=non_existent_collection", "--dry-run"},
			expectError: false, // Dry-run doesn't validate existence
		},
		{
			name:        "retrieve document with valid parameters",
			args:        []string{"document", "list", "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: false,
		},
		{
			name:        "retrieve document using document command",
			args:        []string{"document", "list", "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: false,
		},
		{
			name:        "retrieve document using document command",
			args:        []string{"document", "list", "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("../maestro", tc.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
				if tc.errorMsg != "" && !strings.Contains(outputStr, tc.errorMsg) {
					t.Errorf("Expected error message containing '%s', got: %s", tc.errorMsg, outputStr)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v, output: %s", err, outputStr)
				}
			}
		})
	}
}

func TestRetrieveVerboseMode(t *testing.T) {
	// Test verbose mode for retrieve collection
	cmd := exec.Command("../maestro", "collection", "list", "--vdb=test_local_milvus", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Retrieve command with verbose failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Listing collections") {
		t.Errorf("Expected verbose output, got: %s", outputStr)
	}
}
