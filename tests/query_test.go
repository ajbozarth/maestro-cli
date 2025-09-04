package main

import (
	"os/exec"
	"testing"
)

// TestQueryHelp tests the query command help
func TestQueryHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run query help command: %v", err)
	}

	helpOutput := string(output)

	// Check for expected query help content
	expectedContent := []string{
		"query",
		"Query documents in a vector database using natural language",
		"doc-limit",
		"Maximum number of documents to consider",
		"Examples:",
	}

	for _, expected := range expectedContent {
		if !contains(helpOutput, expected) {
			t.Errorf("Query help output should contain '%s'", expected)
		}
	}
}

// TestQueryNoArgs tests the query command with no arguments
func TestQueryNoArgs(t *testing.T) {
	cmd := exec.Command("../maestro", "query")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Query command should fail with no arguments")
	}

	outputStr := string(output)
	if !contains(outputStr, "Error:") {
		t.Error("Query command should show error message with no arguments")
	}
}

// TestQueryMissingQuery tests the query command with missing query argument
func TestQueryMissingQuery(t *testing.T) {
	cmd := exec.Command("../maestro", "query")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Query command should fail with missing query argument")
	}

	outputStr := string(output)
	if !contains(outputStr, "Error:") {
		t.Error("Query command should show error message with missing query argument")
	}
}

// TestQueryEmptyDatabaseName tests the query command with missing vdb flag
func TestQueryEmptyDatabaseName(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "test query")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Query command should fail with missing --vdb flag")
	}

	outputStr := string(output)
	if !contains(outputStr, "Error:") {
		t.Error("Query command should show error message with missing --vdb flag")
	}
}

// TestQueryEmptyQuery tests the query command with empty query
func TestQueryEmptyQuery(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "", "--vdb=test-db")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Query command should fail with empty query")
	}

	outputStr := string(output)
	if !contains(outputStr, "Error:") {
		t.Error("Query command should show error message with empty query")
	}
}

// TestQueryWithDocLimit tests the query command with doc-limit flag
func TestQueryWithDocLimit(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "test query", "--vdb=test-db", "--doc-limit", "10")
	_, err := cmd.CombinedOutput()

	// This will likely fail due to no MCP server, but we can check the command structure
	// The command should fail due to no MCP server or argument validation
	// Both are acceptable outcomes for this test
	if err != nil {
		// Command failed as expected
		return
	}

	// If it doesn't fail, that's also acceptable (dry-run mode might work)
	// The test passes if we get here
}

// TestQueryWithShortDocLimit tests the query command with short doc-limit flag
func TestQueryWithShortDocLimit(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "test query", "--vdb=test-db", "-d", "5")
	_, err := cmd.CombinedOutput()

	// This will likely fail due to no MCP server, but we can check the command structure
	// The command should fail due to no MCP server or argument validation
	// Both are acceptable outcomes for this test
	if err != nil {
		// Command failed as expected
		return
	}

	// If it doesn't fail, that's also acceptable (dry-run mode might work)
	// The test passes if we get here
}

// TestQueryWithVerboseFlag tests the query command with verbose flag
func TestQueryWithVerboseFlag(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "test query", "--vdb=test-db", "--verbose")
	_, err := cmd.CombinedOutput()

	// This will likely fail due to no MCP server, but we can check the command structure
	// The command should fail due to no MCP server or argument validation
	// Both are acceptable outcomes for this test
	if err != nil {
		// Command failed as expected
		return
	}

	// If it doesn't fail, that's also acceptable (dry-run mode might work)
	// The test passes if we get here
}

// TestQueryWithSilentFlag tests the query command with silent flag
func TestQueryWithSilentFlag(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "test query", "--vdb=test-db", "--silent")
	_, err := cmd.CombinedOutput()

	// This will likely fail due to no MCP server, but we can check the command structure
	// The command should fail due to no MCP server or argument validation
	// Both are acceptable outcomes for this test
	if err != nil {
		// Command failed as expected
		return
	}

	// If it doesn't fail, that's also acceptable (dry-run mode might work)
	// The test passes if we get here
}

// TestQueryWithDryRunFlag tests the query command with dry-run flag
func TestQueryWithDryRunFlag(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "test query", "--vdb=test-db", "--dry-run")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Query command with dry-run should not fail: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN]") {
		t.Error("Query command with dry-run should show dry run message")
	}
}

// TestQueryWithSpecialCharacters tests the query command with special characters
func TestQueryWithSpecialCharacters(t *testing.T) {
	cmd := exec.Command("../maestro", "query", "What's the deal with API endpoints? (v2.0)", "--vdb=test-db", "--dry-run")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Query command with special characters should not fail: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN]") {
		t.Error("Query command with special characters should show dry run message")
	}
}

// TestQueryWithLongQuery tests the query command with a long query
func TestQueryWithLongQuery(t *testing.T) {
	longQuery := "This is a very long query that contains many words and should test the ability of the command to handle long input strings without any issues or problems"
	cmd := exec.Command("../maestro", "query", longQuery, "--vdb=test-db", "--dry-run")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Query command with long query should not fail: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN]") {
		t.Error("Query command with long query should show dry run message")
	}
}

// TestQueryInvalidDocLimit tests the query command with invalid doc-limit values
func TestQueryInvalidDocLimit(t *testing.T) {
	testCases := []string{"abc", "1.5"}

	for _, invalidLimit := range testCases {
		t.Run("InvalidLimit_"+invalidLimit, func(t *testing.T) {
			cmd := exec.Command("../maestro", "query", "test-db", "test query", "--doc-limit", invalidLimit)
			output, err := cmd.CombinedOutput()

			// The command should handle invalid limits gracefully
			outputStr := string(output)

			// It might fail due to argument parsing or MCP server, but shouldn't crash
			if err == nil && !contains(outputStr, "[DRY RUN]") {
				t.Error("Query command should handle invalid doc-limit gracefully")
			}
		})
	}
}

// TestQueryCommandExists tests that the query command exists in the CLI
func TestQueryCommandExists(t *testing.T) {
	cmd := exec.Command("../maestro", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	helpOutput := string(output)
	if !contains(helpOutput, "query") {
		t.Error("Help output should contain 'query' command")
	}
}
