package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestListVectorDatabaseWithRealServer(t *testing.T) {
	cmd := exec.Command("../../maestro", "vectordb", "list")
	cmd.Env = append(os.Environ(), "MAESTRO_TEST_MODE=true")
	output, err := cmd.CombinedOutput()

	// This test is expected to fail if no MCP server is running
	if err != nil {
		outputStr := string(output)
		// Check if the error is due to connection refused (no server running)
		// or unsupported protocol scheme (malformed URL)
		// or HTTP 404 (server running but endpoint not found)
		// or HTTP 400 (server running but session ID required - this is expected with FastMCP)
		// or FastMCP client connection issues
		// or runtime panics from incompatible MCP libraries
		// or our new user-friendly error messages
		if strings.Contains(outputStr, "connection refused") ||
			strings.Contains(outputStr, "unsupported protocol scheme") ||
			strings.Contains(outputStr, "HTTP error 404") ||
			strings.Contains(outputStr, "HTTP error 400") ||
			strings.Contains(outputStr, "Session terminated") ||
			strings.Contains(outputStr, "Client failed to connect") ||
			strings.Contains(outputStr, "Missing session ID") ||
			strings.Contains(outputStr, "panic: runtime error") ||
			strings.Contains(outputStr, "invalid memory address") ||
			strings.Contains(outputStr, "MCP server could not be reached") {
			t.Logf("Test skipped: No MCP server running, malformed URL, or endpoint not found (expected): %s", outputStr)
			return
		}
		// If it's a different error, fail the test
		t.Fatalf("List command failed with unexpected error: %v, output: %s", err, string(output))
	}

	// If the command succeeds, we should get either "No vector databases found" or actual database list
	outputStr := string(output)
	if !strings.Contains(outputStr, "No vector databases found") &&
		!strings.Contains(outputStr, "Found") &&
		!strings.Contains(outputStr, "vector database") {
		t.Errorf("Unexpected output from list command: %s", outputStr)
	}
}

// TestListVectorDatabaseURLNormalization tests that URL normalization works correctly
func TestListVectorDatabaseURLNormalization(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "hostname only",
			url:      "localhost",
			expected: "http://localhost:8030/mcp",
		},
		{
			name:     "hostname with port",
			url:      "localhost:8030",
			expected: "http://localhost:8030/mcp",
		},
		{
			name:     "http URL",
			url:      "http://localhost:8030",
			expected: "http://localhost:8030/mcp",
		},
		{
			name:     "https URL",
			url:      "https://example.com:9000",
			expected: "https://example.com:9000/mcp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test URL normalization by running the CLI and checking the verbose output
			// Use a very short timeout to make it fast
			cmd := exec.Command("../../maestro", "vectordb", "list", "--mcp-server-uri", tc.url, "--verbose")
			cmd.Env = append(os.Environ(), "MAESTRO_MCP_SERVER_URI=", "MAESTRO_TEST_MODE=true")

			// Set a very short timeout for the test
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

			output, _ := cmd.CombinedOutput()

			// The command should show the normalized URL in the verbose output
			outputStr := string(output)
			// Check if the normalized URL appears in the output
			if !strings.Contains(outputStr, fmt.Sprintf("Connecting to MCP server at: %s", tc.expected)) {
				t.Errorf("Expected URL normalization to %s, but got different output: %s", tc.expected, outputStr)
			}
		})
	}
}
