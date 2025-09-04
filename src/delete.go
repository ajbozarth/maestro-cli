package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CollectionNotFoundError is a custom error type for collection not found errors
type CollectionNotFoundError struct {
	CollectionName string
	VDBName        string
}

func (e *CollectionNotFoundError) Error() string {
	return fmt.Sprintf("Collection '%s' not found in vector database '%s'", e.CollectionName, e.VDBName)
}

// DocumentNotFoundError is a custom error type for document not found errors
type DocumentNotFoundError struct {
	DocumentName   string
	CollectionName string
	VDBName        string
}

func (e *DocumentNotFoundError) Error() string {
	return fmt.Sprintf("Document '%s' not found in collection '%s' of vector database '%s'", e.DocumentName, e.CollectionName, e.VDBName)
}

// confirmDestructiveOperation prompts the user for confirmation before performing a destructive operation
func confirmDestructiveOperation(operation, resourceName string) error {
	// Skip confirmation if --force flag is used
	if force {
		return nil
	}

	// Skip confirmation in dry-run mode
	if dryRun {
		return nil
	}

	// Skip confirmation in silent mode
	if silent {
		return nil
	}

	fmt.Printf("⚠️  Are you sure you want to %s '%s'? This action cannot be undone. [y/N]: ", operation, resourceName)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		return fmt.Errorf("operation cancelled by user")
	}

	return nil
}

func deleteVectorDatabase(name string) error {
	if verbose {
		fmt.Printf("Deleting vector database: %s\n", name)
	}

	if dryRun {
		if !silent {
			fmt.Printf("[DRY RUN] Would delete vector database '%s'\n", name)
		}
		return nil
	}

	// Validate the name
	if name == "" {
		return fmt.Errorf("vector database name is required")
	}

	// Request confirmation for destructive operation
	if err := confirmDestructiveOperation("delete", fmt.Sprintf("vector database '%s'", name)); err != nil {
		return err
	}

	// Perform the deletion
	if err := performVectorDatabaseDeletion(name); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}

	if !silent {
		fmt.Printf("✅ Vector database '%s' deleted successfully\n", name)
	}

	return nil
}

func deleteCollection(vdbName, collectionName string) error {
	if verbose {
		fmt.Printf("Deleting collection: %s from vector database: %s\n", collectionName, vdbName)
	}

	if dryRun {
		if !silent {
			fmt.Printf("[DRY RUN] Would delete collection '%s' from vector database '%s'\n", collectionName, vdbName)
		}
		return nil
	}

	// Validate the names
	if vdbName == "" {
		return fmt.Errorf("vector database name is required")
	}
	if collectionName == "" {
		return fmt.Errorf("collection name is required")
	}

	// Request confirmation for destructive operation
	if err := confirmDestructiveOperation("delete", fmt.Sprintf("collection '%s' from vector database '%s'", collectionName, vdbName)); err != nil {
		return err
	}

	// Perform the deletion
	if err := performCollectionDeletion(vdbName, collectionName); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}

	if !silent {
		fmt.Printf("✅ Collection '%s' deleted successfully from vector database '%s'\n", collectionName, vdbName)
	}

	return nil
}

func deleteDocument(vdbName, collectionName, docName string) error {
	if verbose {
		fmt.Printf("Deleting document: %s from collection: %s in vector database: %s\n", docName, collectionName, vdbName)
	}

	if dryRun {
		if !silent {
			fmt.Printf("[DRY RUN] Would delete document '%s' from collection '%s' in vector database '%s'\n", docName, collectionName, vdbName)
		}
		return nil
	}

	// Validate the names
	if vdbName == "" {
		return fmt.Errorf("vector database name is required")
	}
	if collectionName == "" {
		return fmt.Errorf("collection name is required")
	}
	if docName == "" {
		return fmt.Errorf("document name is required")
	}

	// Request confirmation for destructive operation
	if err := confirmDestructiveOperation("delete", fmt.Sprintf("document '%s' from collection '%s' in vector database '%s'", docName, collectionName, vdbName)); err != nil {
		return err
	}

	// Perform the deletion
	if err := performDocumentDeletion(vdbName, collectionName, docName); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}

	if !silent {
		fmt.Printf("✅ Document '%s' deleted successfully from collection '%s' in vector database '%s'\n", docName, collectionName, vdbName)
	}

	return nil
}

func performDocumentDeletion(vdbName, collectionName, docName string) error {
	if verbose {
		fmt.Printf("Deleting document '%s' from collection '%s' in vector database '%s'\n", docName, collectionName, vdbName)
	}

	// Get MCP server URI
	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		return fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	if verbose {
		fmt.Printf("Connecting to MCP server at: %s\n", serverURI)
	}

	// Create MCP client
	client, err := NewMCPClient(serverURI)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	// Check if database exists before deleting document
	exists, err := client.DatabaseExists(vdbName)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("vector database '%s' does not exist", vdbName)
	}

	// Document existence check is now handled by the MCP server

	// Call the MCP server to delete the document with panic recovery
	var deleteErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				deleteErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		deleteErr = client.DeleteDocumentFromCollection(vdbName, collectionName, docName)
	}()

	if deleteErr != nil {
		// Check if it's a document not found error and provide cleaner output
		errMsg := deleteErr.Error()
		if strings.Contains(errMsg, "not found in collection") {
			// Extract the document name, collection name and vdb name from the error message
			// Error format: "Document 'doc_name' not found in collection 'collection_name' of vector database 'vdb_name'"
			return &DocumentNotFoundError{DocumentName: docName, CollectionName: collectionName, VDBName: vdbName}
		}
		// For other errors, wrap with context
		return fmt.Errorf("failed to delete document: %w", deleteErr)
	}

	if verbose {
		fmt.Println("Document deletion completed successfully")
	}

	return nil
}

func performCollectionDeletion(vdbName, collectionName string) error {
	if verbose {
		fmt.Printf("Deleting collection '%s' from vector database '%s'\n", collectionName, vdbName)
	}

	// Get MCP server URI
	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		return fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	if verbose {
		fmt.Printf("Connecting to MCP server at: %s\n", serverURI)
	}

	// Create MCP client
	client, err := NewMCPClient(serverURI)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	// Check if database exists before deleting collection
	exists, err := client.DatabaseExists(vdbName)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("vector database '%s' does not exist", vdbName)
	}

	// Collection existence check is now handled by the MCP server

	// Call the MCP server to delete the collection with panic recovery
	var deleteErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				deleteErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		deleteErr = client.DeleteCollection(vdbName, collectionName)
	}()

	if deleteErr != nil {
		// Check if it's a collection not found error and provide cleaner output
		errMsg := deleteErr.Error()
		if strings.Contains(errMsg, "not found in vector database") {
			// Extract the collection name and vdb name from the error message
			// Error format: "Collection 'collection_name' not found in vector database 'vdb_name'"
			return &CollectionNotFoundError{CollectionName: collectionName, VDBName: vdbName}
		}
		// For other errors, wrap with context
		return fmt.Errorf("failed to delete collection: %w", deleteErr)
	}

	if verbose {
		fmt.Println("Collection deletion completed successfully")
	}

	return nil
}

func performVectorDatabaseDeletion(name string) error {
	if verbose {
		fmt.Printf("Deleting vector database '%s'\n", name)
	}

	// Get MCP server URI
	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		return fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	if verbose {
		fmt.Printf("Connecting to MCP server at: %s\n", serverURI)
	}

	// Create MCP client
	client, err := NewMCPClient(serverURI)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	// Check if database exists before deleting
	exists, err := client.DatabaseExists(name)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("vector database '%s' does not exist", name)
	}

	// Call the MCP server to delete the database with panic recovery
	var deleteErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				deleteErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		deleteErr = client.DeleteVectorDatabase(name)
	}()

	if deleteErr != nil {
		return fmt.Errorf("failed to delete vector database: %w", deleteErr)
	}

	if verbose {
		fmt.Println("Vector database deletion completed successfully")
	}

	return nil
}
