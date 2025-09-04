package main

import (
	"os/exec"
	"testing"
)

// TestDeleteVectorDatabase tests the delete vector-database command
func TestDeleteVectorDatabase(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "delete", "test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteVectorDatabaseWithVerbose tests the delete command with verbose output
func TestDeleteVectorDatabaseWithVerbose(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "delete", "test-db", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete command failed with verbose: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "Deleting vector database: test-db") {
		t.Errorf("Should show verbose output, got: %s", outputStr)
	}
	if !contains(outputStr, "[DRY RUN] Would delete vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteVectorDatabaseWithInvalidResourceType tests invalid resource type
func TestDeleteVectorDatabaseWithInvalidResourceType(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "invalid-action", "test-db")
	output, _ := cmd.CombinedOutput()

	// Should show help for invalid action (Cobra's default behavior)
	outputStr := string(output)
	if !contains(outputStr, "Available Commands") {
		t.Errorf("Error message should mention available commands, got: %s", outputStr)
	}
}

// TestDeleteVectorDatabaseWithEmptyName tests with empty name
func TestDeleteVectorDatabaseWithEmptyName(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "delete", "")
	output, err := cmd.CombinedOutput()

	// Should fail with empty name
	if err == nil {
		t.Error("Delete command should fail with empty name")
	}

	outputStr := string(output)
	if !contains(outputStr, "vector database name is required") {
		t.Errorf("Error message should mention name is required, got: %s", outputStr)
	}
}

// TestDeleteVectorDatabaseDryRun tests dry run functionality
func TestDeleteVectorDatabaseDryRun(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "delete", "test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete command failed with dry run: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteVectorDatabaseSilent tests silent mode
func TestDeleteVectorDatabaseSilent(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "delete", "test-db", "--silent", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete command failed with silent mode: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if contains(outputStr, "[DRY RUN] Would delete vector database 'test-db'") {
		t.Error("Silent mode should not show dry run message")
	}
}

// TestDeleteVectorDatabaseWithLongName tests with a long name
func TestDeleteVectorDatabaseWithLongName(t *testing.T) {
	longName := "very-long-vector-database-name-that-should-still-work"
	cmd := exec.Command("../maestro", "vdb", "delete", longName, "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete command failed with long name: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete vector database '"+longName+"'") {
		t.Errorf("Should show dry run message with long name, got: %s", outputStr)
	}
}

// TestDeleteVectorDatabaseWithSpecialCharacters tests with special characters in name
func TestDeleteVectorDatabaseWithSpecialCharacters(t *testing.T) {
	specialName := "test-db-with-special-chars_123"
	cmd := exec.Command("../maestro", "vdb", "delete", specialName, "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete command failed with special characters: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete vector database '"+specialName+"'") {
		t.Errorf("Should show dry run message with special characters, got: %s", outputStr)
	}
}

// TestDeleteVectorDatabaseWithVdbShortcut tests the delete command with vdb command
func TestDeleteVectorDatabaseWithVdbShortcut(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "delete", "test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete command with vdb shortcut failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteCollection tests the delete collection command
func TestDeleteCollection(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "delete", "test-collection", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete collection command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete collection 'test-collection' from vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteCollectionWithVerbose tests the delete collection command with verbose output
func TestDeleteCollectionWithVerbose(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "delete", "test-collection", "--vdb=test-db", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete collection command failed with verbose: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "Deleting collection: test-collection from vector database: test-db") {
		t.Errorf("Should show verbose output, got: %s", outputStr)
	}
	if !contains(outputStr, "[DRY RUN] Would delete collection 'test-collection' from vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteCollectionWithInvalidResourceType tests invalid collection resource type
func TestDeleteCollectionWithInvalidResourceType(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "invalid-action", "test-collection", "--vdb=test-db")
	output, err := cmd.CombinedOutput()

	// Should fail with invalid action
	if err == nil {
		t.Error("Delete collection command should fail with invalid action")
	}

	outputStr := string(output)
	if !contains(outputStr, "Available Commands") {
		t.Errorf("Error message should mention available commands, got: %s", outputStr)
	}
}

// TestDeleteCollectionWithMissingCollectionName tests with missing collection name
func TestDeleteCollectionWithMissingCollectionName(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "delete")
	output, err := cmd.CombinedOutput()

	// Should fail with missing collection name
	if err == nil {
		t.Error("Delete collection command should fail with missing collection name")
	}

	outputStr := string(output)
	if !contains(outputStr, "accepts 1 arg") {
		t.Errorf("Error message should mention missing collection name, got: %s", outputStr)
	}
}

// TestDeleteCollectionWithEmptyNames tests with empty names
func TestDeleteCollectionWithEmptyNames(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "delete", "", "--vdb=test-db")
	output, err := cmd.CombinedOutput()

	// Should fail with empty collection name
	if err == nil {
		t.Error("Delete collection command should fail with empty collection name")
	}

	outputStr := string(output)
	if !contains(outputStr, "collection name is required") {
		t.Errorf("Error message should mention collection name is required, got: %s", outputStr)
	}
}

// TestDeleteCollectionWithVdbColAlias tests the delete command with collection command
func TestDeleteCollectionWithVdbColAlias(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "delete", "test-collection", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete collection command with vdb-col alias failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete collection 'test-collection' from vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteCollectionWithColAlias tests the delete command with collection command
func TestDeleteCollectionWithColAlias(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "delete", "test-collection", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete collection command with col alias failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete collection 'test-collection' from vector database 'test-db'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteDocument tests the delete document command
func TestDeleteDocument(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "delete", "test-document", "--vdb=test-vdb", "--collection=test-collection", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete document command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete document 'test-document' from collection 'test-collection' in vector database 'test-vdb'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteDocumentWithVerbose tests the delete document command with verbose output
func TestDeleteDocumentWithVerbose(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "delete", "test-document", "--vdb=test-vdb", "--collection=test-collection", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete document command failed with verbose: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "Deleting document: test-document from collection: test-collection in vector database: test-vdb") {
		t.Errorf("Should show verbose output, got: %s", outputStr)
	}
	if !contains(outputStr, "[DRY RUN] Would delete document 'test-document' from collection 'test-collection' in vector database 'test-vdb'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteDocumentWithInvalidResourceType tests invalid resource type
func TestDeleteDocumentWithInvalidResourceType(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "invalid-action", "test-document", "--vdb=test-vdb", "--collection=test-collection")
	output, err := cmd.CombinedOutput()

	// Should fail with invalid action
	if err == nil {
		t.Error("Delete command should fail with invalid action")
	}

	outputStr := string(output)
	if !contains(outputStr, "Available Commands") {
		t.Errorf("Error message should mention available commands, got: %s", outputStr)
	}
}

// TestDeleteDocumentWithMissingArguments tests with missing arguments
func TestDeleteDocumentWithMissingArguments(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "delete")
	output, err := cmd.CombinedOutput()

	// Should fail with missing arguments
	if err == nil {
		t.Error("Delete command should fail with missing arguments")
	}

	outputStr := string(output)
	if !contains(outputStr, "accepts 1 arg") {
		t.Errorf("Error message should mention missing arguments, got: %s", outputStr)
	}
}

// TestDeleteDocumentWithEmptyNames tests with empty names
func TestDeleteDocumentWithEmptyNames(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "delete", "", "--vdb=test-vdb", "--collection=test-collection")
	output, err := cmd.CombinedOutput()

	// Should fail with empty names
	if err == nil {
		t.Error("Delete command should fail with empty names")
	}

	outputStr := string(output)
	if !contains(outputStr, "document name is required") {
		t.Errorf("Error message should mention name is required, got: %s", outputStr)
	}
}

// TestDeleteDocumentWithVdbDocAlias tests the document command
func TestDeleteDocumentWithVdbDocAlias(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "delete", "test-document", "--vdb=test-vdb", "--collection=test-collection", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete document command with vdb-doc alias failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete document 'test-document' from collection 'test-collection' in vector database 'test-vdb'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteDocumentWithDocAlias tests the document command
func TestDeleteDocumentWithDocAlias(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "delete", "test-document", "--vdb=test-vdb", "--collection=test-collection", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete document command with doc alias failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN] Would delete document 'test-document' from collection 'test-collection' in vector database 'test-vdb'") {
		t.Errorf("Should show dry run message, got: %s", outputStr)
	}
}

// TestDeleteDocumentSilent tests silent mode
func TestDeleteDocumentSilent(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "delete", "test-document", "--vdb=test-vdb", "--collection=test-collection", "--silent", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Delete document command failed with silent mode: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if contains(outputStr, "[DRY RUN] Would delete document 'test-document' from collection 'test-collection' in vector database 'test-vdb'") {
		t.Error("Silent mode should not show dry run message")
	}
}
