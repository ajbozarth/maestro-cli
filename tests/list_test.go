package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestListVectorDatabase(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "vdb", "list", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list vector databases") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListVectorDatabaseWithVerbose(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "vdb", "list", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List command with verbose failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Listing vector databases") {
		t.Errorf("Expected verbose message 'Listing vector databases', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "[DRY RUN] Would list vector databases") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListVectorDatabaseWithInvalidResourceType(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "invalid-action")
	output, _ := cmd.CombinedOutput()

	// Should show help for invalid action
	outputStr := string(output)
	if !strings.Contains(outputStr, "Available Commands") {
		t.Errorf("Expected help output showing available commands, got: %s", outputStr)
	}
}

func TestListVectorDatabaseWithDryRun(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "list", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List command with dry-run failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list vector databases") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListVectorDatabaseWithSilent(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "vdb", "list", "--silent", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List command with silent failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	// Should not contain the default "No vector databases found" message when silent
	if strings.Contains(outputStr, "No vector databases found") {
		t.Errorf("Should not show default message when silent, got: %s", outputStr)
	}
}

func TestListVectorDatabaseHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "list", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List help command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "List all available vector databases") {
		t.Errorf("Expected help message about listing vector databases, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "maestro vectordb list") {
		t.Errorf("Expected usage example, got: %s", outputStr)
	}
}

func TestListVectorDatabaseWithVectorDatabase(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "vdb", "list", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List command with 'vector-databases' failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list vector databases") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListVectorDatabaseWithVdbShortcut(t *testing.T) {
	// Test with the "vdb" command
	cmd := exec.Command("../maestro", "vdb", "list", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List command with vdb failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list vector databases") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListVectorDatabaseWithMultipleFlags(t *testing.T) {
	cmd := exec.Command("../maestro", "vdb", "list", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List command with multiple flags failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Listing vector databases") {
		t.Errorf("Expected verbose message, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "[DRY RUN] Would list vector databases") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

// TestListVectorDatabaseWithRealServer tests the actual MCP server connection
// This test is expected to fail if no MCP server is running, which is acceptable

func TestListEmbeddings(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "embedding", "list", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List embeddings command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list embeddings for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListEmbeddingsWithVerbose(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "embedding", "list", "--vdb=test-db", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List embeddings command with verbose failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Listing embeddings for vector database 'test-db'") {
		t.Errorf("Expected verbose message 'Listing embeddings for vector database 'test-db'', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "[DRY RUN] Would list embeddings for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListEmbeddingsWithEmbedAlias(t *testing.T) {
	// Test that embeds alias no longer exists
	cmd := exec.Command("../maestro", "embeds", "list", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	// Should fail since embeds alias no longer exists
	if err == nil {
		t.Error("Expected command to fail since embeds alias no longer exists")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "unknown command") {
		t.Errorf("Expected error about unknown command, got: %s", outputStr)
	}
}

func TestListEmbeddingsWithVdbEmbedAlias(t *testing.T) {
	// Test that vdb-embeds alias doesn't exist in new CLI structure
	cmd := exec.Command("../maestro", "list", "vdb-embeds", "test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	// Should fail since vdb-embeds alias doesn't exist in new structure
	if err == nil {
		t.Error("Expected command to fail since vdb-embeds alias doesn't exist in new CLI structure")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "unknown command") {
		t.Errorf("Expected error about unknown command, got: %s", outputStr)
	}
}

func TestListEmbeddingsMissingVdbName(t *testing.T) {
	cmd := exec.Command("../maestro", "embedding", "list")
	output, err := cmd.CombinedOutput()

	// Should fail with missing VDB_NAME
	if err == nil {
		t.Error("Expected command to fail with missing --vdb flag")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "failed to select vector database: vector database name is required in non-interactive mode") {
		t.Errorf("Expected error about missing vector database name, got: %s", outputStr)
	}
}

func TestListEmbeddingsHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "embedding", "list", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List embeddings help command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "--vdb string") {
		t.Errorf("Expected help to mention --vdb flag, got: %s", outputStr)
	}
}

func TestListCollections(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "collection", "list", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List collections command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list collections for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListCollectionsWithVerbose(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "collection", "list", "--vdb=test-db", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List collections command with verbose failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Listing collections for vector database 'test-db'") {
		t.Errorf("Expected verbose message 'Listing collections for vector database 'test-db'', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "[DRY RUN] Would list collections for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListCollectionsWithColsAlias(t *testing.T) {
	// Test the collection command
	cmd := exec.Command("../maestro", "collection", "list", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List collection command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list collections for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListCollectionsWithVdbColsAlias(t *testing.T) {
	// Test the collection command
	cmd := exec.Command("../maestro", "collection", "list", "--vdb=test-db", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List collection command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list collections for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListCollectionsMissingVdbName(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "list")
	output, err := cmd.CombinedOutput()

	// Should fail with missing --vdb flag
	if err == nil {
		t.Error("Expected command to fail with missing --vdb flag, but it succeeded")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "failed to select vector database: vector database name is required in non-interactive mode") {
		t.Errorf("Expected error about missing vector database name, got: %s", outputStr)
	}
}

func TestListCollectionsHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "collection", "list", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List collections help command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "--vdb string") {
		t.Errorf("Expected help to mention --vdb flag, got: %s", outputStr)
	}
}

func TestListDocuments(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "document", "list", "--vdb=test-db", "--collection=test-collection", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List documents command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list documents in collection 'test-collection' for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListDocumentsWithVerbose(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "document", "list", "--vdb=test-db", "--collection=test-collection", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List documents command with verbose failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Listing documents in collection 'test-collection' for vector database 'test-db'") {
		t.Errorf("Expected verbose message 'Listing documents in collection 'test-collection' for vector database 'test-db'', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "[DRY RUN] Would list documents in collection 'test-collection' for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListDocumentsWithDocsAlias(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "document", "list", "--vdb=test-db", "--collection=test-collection", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List documents command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list documents in collection 'test-collection' for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListDocumentsWithVdbDocsAlias(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "document", "list", "--vdb=test-db", "--collection=test-collection", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List documents command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list documents in collection 'test-collection' for vector database 'test-db'") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListDocumentsMissingVdbName(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "list")
	output, err := cmd.CombinedOutput()

	// Should fail with missing --vdb flag
	if err == nil {
		t.Error("Expected command to fail with missing --vdb flag, but it succeeded")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "failed to select vector database: vector database name is required in non-interactive mode") {
		t.Errorf("Expected error about missing vector database name, got: %s", outputStr)
	}
}

func TestListDocumentsMissingCollectionName(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "list", "--vdb=test-db")
	output, err := cmd.CombinedOutput()

	// Should fail with missing --collection flag
	if err == nil {
		t.Error("Expected command to fail with missing --collection flag, but it succeeded")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "failed to select collection: collection name is required in non-interactive mode") {
		t.Errorf("Expected error about missing collection name, got: %s", outputStr)
	}
}

func TestListDocumentsHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "document", "list", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("List documents help command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "--vdb string") {
		t.Errorf("Expected help to mention --vdb flag, got: %s", outputStr)
	}
}

func TestListChunkingStrategiesDryRun(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "chunking", "list", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Chunking list command failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "[DRY RUN] Would list chunking strategies") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}

func TestListChunkingStrategiesVerbose(t *testing.T) {
	// Use dry-run mode since we don't have a real MCP server running
	cmd := exec.Command("../maestro", "chunking", "list", "--verbose", "--dry-run")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Chunking list command with verbose failed: %v, output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Listing supported chunking strategies") {
		t.Errorf("Expected verbose message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "[DRY RUN] Would list chunking strategies") {
		t.Errorf("Expected dry-run message, got: %s", outputStr)
	}
}
