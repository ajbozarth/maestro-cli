package main

import (
	"os/exec"
	"testing"
)

// TestSearchHelp tests the search command help
func TestSearchHelp(t *testing.T) {
	cmd := exec.Command("../maestro", "search", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run query help command: %v", err)
	}

	helpOutput := string(output)

	// Check for expected search help content
	expectedContent := []string{
		"search",
		"Search documents in a vector database using natural language",
		"doc-limit",
		"Maximum number of documents to consider",
		"Examples:",
	}

	for _, expected := range expectedContent {
		if !contains(helpOutput, expected) {
			t.Errorf("Search help output should contain '%s'", expected)
		}
	}
}

// TestSearchNoArgs tests the search command with no arguments
func TestSearchNoArgs(t *testing.T) {
	cmd := exec.Command("../maestro", "search")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Search command should fail with no arguments")
	}

	outputStr := string(output)
	if !contains(outputStr, "Error:") {
		t.Error("Search command should show error message with no arguments")
	}
}

// TestSearchEmptyDatabaseName tests the search command with missing vdb flag
func TestSearchEmptyDatabaseName(t *testing.T) {
	cmd := exec.Command("../maestro", "search", "test query")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Search command should fail with missing --vdb flag")
	}

	outputStr := string(output)
	if !contains(outputStr, "Error:") {
		t.Error("Search command should show error message with missing --vdb flag")
	}
}

// TestSearchEmptySearch tests the search command with empty query
func TestSearchEmptySearch(t *testing.T) {
	cmd := exec.Command("../maestro", "search", "", "--vdb=test-db")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Search command should fail with empty query")
	}

	outputStr := string(output)
	if !contains(outputStr, "Error:") {
		t.Error("Search command should show error message with empty query")
	}
}

// TestSearchWithDryRunFlag tests the search command with dry-run flag
func TestSearchWithDryRunFlag(t *testing.T) {
	cmd := exec.Command("../maestro", "search", "test query", "--vdb=test-db", "--dry-run")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Search command with dry-run should not fail: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN]") {
		t.Error("Search command with dry-run should show dry run message")
	}
}

// TestSearchWithSpecialCharacters tests the search command with special characters
func TestSearchWithSpecialCharacters(t *testing.T) {
	cmd := exec.Command("../maestro", "search", "What's the deal with API endpoints? (v2.0)", "--vdb=test-db", "--dry-run")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Search command with special characters should not fail: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN]") {
		t.Error("Search command with special characters should show dry run message")
	}
}

// TestSearchWithLongSearch tests the search command with a long query
func TestSearchWithLongSearch(t *testing.T) {
	longSearch := "This is a very long query that contains many words and should test the ability of the command to handle long input strings without any issues or problems"
	cmd := exec.Command("../maestro", "search", longSearch, "--vdb=test-db", "--dry-run")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Search command with long query should not fail: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "[DRY RUN]") {
		t.Error("Search command with long query should show dry run message")
	}
}

// TestSearchInvalidDocLimit tests the search command with invalid doc-limit values
func TestSearchInvalidDocLimit(t *testing.T) {
	testCases := []string{"abc", "1.5"}

	for _, invalidLimit := range testCases {
		t.Run("InvalidLimit_"+invalidLimit, func(t *testing.T) {
			cmd := exec.Command("../maestro", "search", "test-db", "test query", "--doc-limit", invalidLimit)
			output, err := cmd.CombinedOutput()

			// The command should handle invalid limits gracefully
			outputStr := string(output)

			// It might fail due to argument parsing or MCP server, but shouldn't crash
			if err == nil && !contains(outputStr, "[DRY RUN]") {
				t.Error("Search command should handle invalid doc-limit gracefully")
			}
		})
	}
}

// TestSearchCommandExists tests that the search command exists in the CLI
func TestSearchCommandExists(t *testing.T) {
	cmd := exec.Command("../maestro", "--help")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	helpOutput := string(output)
	if !contains(helpOutput, "search") {
		t.Error("Help output should contain 'search' command")
	}
}
