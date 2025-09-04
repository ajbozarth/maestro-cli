package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestWriteDocument(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t, "test-*.txt", "This is a test document content")
	defer os.Remove(tempFile)

	// Test cases
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "write document with valid parameters",
			args:        []string{"document", "create", "--name=test_doc", "--file=" + tempFile, "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: false,
		},
		{
			name:        "write document with non-existent database (dry-run succeeds)",
			args:        []string{"document", "create", "--name=test_doc", "--file=" + tempFile, "--vdb=non_existent_db", "--collection=test_collection", "--dry-run"},
			expectError: false, // Dry-run doesn't validate existence
		},
		{
			name:        "write document with non-existent collection (dry-run succeeds)",
			args:        []string{"document", "create", "--name=test_doc", "--file=" + tempFile, "--vdb=test_local_milvus", "--collection=non_existent_collection", "--dry-run"},
			expectError: false, // Dry-run doesn't validate existence
		},
		{
			name:        "write document with missing file name",
			args:        []string{"document", "create", "--name=test_doc", "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: true,
			errorMsg:    "--file flag is required",
		},
		{
			name:        "write document with non-existent file",
			args:        []string{"document", "create", "--name=test_doc", "--file=nonexistent.txt", "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: true,
			errorMsg:    "file not found",
		},
		{
			name:        "write document using document command",
			args:        []string{"document", "create", "--name=test_doc", "--file=" + tempFile, "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: false,
		},
		{
			name:        "write document using document command",
			args:        []string{"document", "create", "--name=test_doc", "--file=" + tempFile, "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
			expectError: false,
		},
		{
			name:        "write document with file flag",
			args:        []string{"document", "create", "--name=test_doc", "--file=" + tempFile, "--vdb=test_local_milvus", "--collection=test_collection", "--dry-run"},
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

func TestWriteDocumentWithVerbose(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t, "test-*.txt", "This is a test document content")
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "document", "create", "--name=test_doc", "--file="+tempFile, "--vdb=test_local_milvus", "--collection=test_collection", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Write command with verbose failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	// Check for either the progress indicator message or the dry-run message
	if !strings.Contains(outputStr, "Creating document") && !strings.Contains(outputStr, "Dry run completed") && !strings.Contains(outputStr, "[DRY RUN] Would create document") {
		t.Errorf("Expected verbose message 'Creating document', 'Dry run completed', or dry-run message, got: %s", outputStr)
	}
}

func TestWriteDocumentWithSilent(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t, "test-*.txt", "This is a test document content")
	defer os.Remove(tempFile)

	cmd := exec.Command("../maestro", "document", "create", "--name=test_doc", "--file="+tempFile, "--vdb=test_local_milvus", "--collection=test_collection", "--silent", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Write command with silent failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "[DRY RUN] Would write document") {
		t.Error("Silent mode should not show dry-run message")
	}
}

func TestWriteDocumentHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "create", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Document create help command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Create a document") {
		t.Errorf("Expected help message about creating documents, got: %s", outputStr)
	}
}

func TestWriteDocumentSubcommandHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "create", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Document create help command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Create a document") {
		t.Errorf("Expected help message about creating documents, got: %s", outputStr)
	}
}

func TestWriteDocumentWithInvalidArguments(t *testing.T) {
	// Test cases for invalid arguments
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "write document with missing arguments",
			args:        []string{"document", "create", "--vdb=test_local_milvus", "--collection=test_collection"},
			expectError: true,
			errorMsg:    "--name flag is required",
		},
		{
			name:        "write document with missing file",
			args:        []string{"document", "create", "--name=test_doc", "--vdb=test_local_milvus", "--collection=test_collection"},
			expectError: true,
			errorMsg:    "--file flag is required",
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
