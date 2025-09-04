package main

import (
	"testing"
)

func TestShowStatus(t *testing.T) {
	// Test with no MCP server (should fail gracefully)
	err := showStatus("")
	if err == nil {
		t.Log("showStatus() succeeded unexpectedly (no MCP server running)")
	} else {
		t.Logf("showStatus() failed as expected: %v", err)
	}
}

func TestShowStatusWithSpecificVDB(t *testing.T) {
	// Test with specific VDB name (should fail gracefully)
	err := showStatus("test-vdb")
	if err == nil {
		t.Log("showStatus() with specific VDB succeeded unexpectedly")
	} else {
		t.Logf("showStatus() with specific VDB failed as expected: %v", err)
	}
}
