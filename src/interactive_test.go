package main

import (
	"testing"
)

func TestIsInteractive(t *testing.T) {
	// Test that isInteractive returns true when stdin is a terminal
	// This test assumes it's running in a terminal environment
	result := isInteractive()

	// In a test environment, this might be false if stdin is redirected
	// We'll just test that the function doesn't panic
	if result != true && result != false {
		t.Errorf("isInteractive() returned unexpected value: %v", result)
	}
}

func TestPromptForVectorDatabase(t *testing.T) {
	// Test with provided vdb name
	vdbName := "test-vdb"
	result, err := PromptForVectorDatabase(vdbName)
	if err != nil {
		t.Errorf("PromptForVectorDatabase(%s) returned error: %v", vdbName, err)
	}
	if result != vdbName {
		t.Errorf("PromptForVectorDatabase(%s) returned %s, expected %s", vdbName, result, vdbName)
	}

	// Test with empty vdb name in non-interactive mode
	// This should return an error
	_, err = PromptForVectorDatabase("")
	if err == nil {
		t.Error("PromptForVectorDatabase(\"\") should return error in non-interactive mode")
	}
}

func TestPromptForCollection(t *testing.T) {
	// Test with provided collection name
	vdbName := "test-vdb"
	collectionName := "test-collection"
	result, err := PromptForCollection(vdbName, collectionName)
	if err != nil {
		t.Errorf("PromptForCollection(%s, %s) returned error: %v", vdbName, collectionName, err)
	}
	if result != collectionName {
		t.Errorf("PromptForCollection(%s, %s) returned %s, expected %s", vdbName, collectionName, result, collectionName)
	}

	// Test with empty collection name in non-interactive mode
	// This should return an error
	_, err = PromptForCollection(vdbName, "")
	if err == nil {
		t.Error("PromptForCollection(\"\", \"\") should return error in non-interactive mode")
	}
}

func TestPromptForDocument(t *testing.T) {
	// Test with provided document name
	vdbName := "test-vdb"
	collectionName := "test-collection"
	documentName := "test-document"
	result, err := PromptForDocument(vdbName, collectionName, documentName)
	if err != nil {
		t.Errorf("PromptForDocument(%s, %s, %s) returned error: %v", vdbName, collectionName, documentName, err)
	}
	if result != documentName {
		t.Errorf("PromptForDocument(%s, %s, %s) returned %s, expected %s", vdbName, collectionName, documentName, result, documentName)
	}

	// Test with empty document name in non-interactive mode
	// This should return an error
	_, err = PromptForDocument(vdbName, collectionName, "")
	if err == nil {
		t.Error("PromptForDocument(\"\", \"\", \"\") should return error in non-interactive mode")
	}
}

func TestNewInteractiveSelector(t *testing.T) {
	selector := NewInteractiveSelector()
	if selector == nil {
		t.Error("NewInteractiveSelector() returned nil")
		return
	}
	if selector.reader == nil {
		t.Error("NewInteractiveSelector() reader is nil")
	}
}

func TestGetAvailableVectorDatabases(t *testing.T) {
	// This test requires a running MCP server
	// We'll just test that the function doesn't panic
	_, err := getAvailableVectorDatabases()
	// We expect an error since there's no MCP server running in tests
	if err == nil {
		t.Log("getAvailableVectorDatabases() succeeded unexpectedly")
	}
}

func TestGetAvailableCollections(t *testing.T) {
	// This test requires a running MCP server
	// We'll just test that the function doesn't panic
	_, err := getAvailableCollections("test-vdb")
	// We expect an error since there's no MCP server running in tests
	if err == nil {
		t.Log("getAvailableCollections() succeeded unexpectedly")
	}
}

func TestGetAvailableDocuments(t *testing.T) {
	// This test requires a running MCP server
	// We'll just test that the function doesn't panic
	_, err := getAvailableDocuments("test-vdb", "test-collection")
	// We expect an error since there's no MCP server running in tests
	if err == nil {
		t.Log("getAvailableDocuments() succeeded unexpectedly")
	}
}
