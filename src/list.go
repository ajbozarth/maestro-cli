package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func listVectorDatabases() error {
	if verbose {
		fmt.Println("Listing vector databases...")
	}

	if dryRun {
		fmt.Println("[DRY RUN] Would list vector databases")
		return nil
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

	// Call the MCP server to list databases with panic recovery
	var databases []DatabaseInfo
	var listErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				listErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		databases, listErr = client.ListDatabases()
	}()

	if listErr != nil {
		return fmt.Errorf("failed to list vector databases: %w", listErr)
	}

	// Display results
	if len(databases) == 0 {
		if !silent {
			fmt.Println("No vector databases found")
		}
		return nil
	}

	if !silent {
		fmt.Printf("Found %d vector database(s):\n\n", len(databases))
	}

	for i, db := range databases {
		if !silent {
			fmt.Printf("%d. %s (%s)\n", i+1, db.Name, db.Type)
			fmt.Printf("   Collection: %s\n", db.Collection)
			fmt.Printf("   Documents: %d\n", db.DocumentCount)
			fmt.Println()
		}
	}

	if verbose {
		fmt.Println("Vector database listing completed successfully")
	}

	return nil
}

func listEmbeddings(vdbName string) error {
	if verbose {
		fmt.Printf("Listing embeddings for vector database '%s'...\n", vdbName)
	}

	if dryRun {
		fmt.Printf("[DRY RUN] Would list embeddings for vector database '%s'\n", vdbName)
		return nil
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

	// Check if the database exists first
	exists, err := client.DatabaseExists(vdbName)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("vector database '%s' does not exist. Please create it first", vdbName)
	}

	// Call the MCP server to get embeddings with panic recovery
	var embeddingsResult string
	var embeddingsErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				embeddingsErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		embeddingsResult, embeddingsErr = client.GetSupportedEmbeddings(vdbName)
	}()

	if embeddingsErr != nil {
		return fmt.Errorf("failed to get embeddings for vector database '%s': %w", vdbName, embeddingsErr)
	}

	// Display results
	if !silent {
		fmt.Println(embeddingsResult)
	}

	if verbose {
		fmt.Printf("Embeddings listing completed successfully for vector database '%s'\n", vdbName)
	}

	return nil
}

func listCollections(vdbName string) error {
	if verbose {
		fmt.Printf("Listing collections for vector database '%s'...\n", vdbName)
	}

	if dryRun {
		fmt.Printf("[DRY RUN] Would list collections for vector database '%s'\n", vdbName)
		return nil
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

	// Check if the database exists first
	exists, err := client.DatabaseExists(vdbName)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("vector database '%s' does not exist. Please create it first", vdbName)
	}

	// Call the MCP server to get collections with panic recovery
	var collectionsResult string
	var collectionsErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				collectionsErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		collectionsResult, collectionsErr = client.ListCollections(vdbName)
	}()

	if collectionsErr != nil {
		return fmt.Errorf("failed to get collections for vector database '%s': %w", vdbName, collectionsErr)
	}

	// Display results
	if !silent {
		fmt.Println(collectionsResult)
	}

	if verbose {
		fmt.Printf("Collections listing completed successfully for vector database '%s'\n", vdbName)
	}

	return nil
}

// showCollectionInfo retrieves and displays detailed collection info
func showCollectionInfo(vdbName, collectionName string) error {
	if verbose {
		fmt.Printf("Retrieving info for collection '%s' in vector database '%s'...\n", collectionName, vdbName)
	}

	if dryRun {
		fmt.Printf("[DRY RUN] Would retrieve info for collection '%s' in vector database '%s'\n", collectionName, vdbName)
		return nil
	}

	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		return fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	if verbose {
		fmt.Printf("Connecting to MCP server at: %s\n", serverURI)
	}

	client, err := NewMCPClient(serverURI)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	// Validate DB exists
	exists, err := client.DatabaseExists(vdbName)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("vector database '%s' does not exist. Please create it first", vdbName)
	}

	var infoStr string
	var infoErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				infoErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		infoStr, infoErr = client.GetCollectionInfo(vdbName, collectionName)
	}()
	if infoErr != nil {
		return fmt.Errorf("failed to get collection info for '%s': %w", collectionName, infoErr)
	}

	if !silent {
		fmt.Println(infoStr)
	}
	if verbose {
		fmt.Printf("Collection info retrieved successfully for '%s' in '%s'\n", collectionName, vdbName)
	}
	return nil
}

func listDocuments(vdbName, collectionName string) error {
	if verbose {
		fmt.Printf("Listing documents in collection '%s' for vector database '%s'...\n", collectionName, vdbName)
	}

	if dryRun {
		fmt.Printf("[DRY RUN] Would list documents in collection '%s' for vector database '%s'\n", collectionName, vdbName)
		return nil
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

	// Check if the database exists first
	exists, err := client.DatabaseExists(vdbName)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("vector database '%s' does not exist. Please create it first", vdbName)
	}

	// Call the MCP server to get documents with panic recovery
	var documentsResult string
	var documentsErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				documentsErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		documentsResult, documentsErr = client.ListDocumentsInCollection(vdbName, collectionName)
	}()

	if documentsErr != nil {
		return fmt.Errorf("failed to get documents for collection '%s' in vector database '%s': %w", collectionName, vdbName, documentsErr)
	}

	// Display results
	if !silent {
		fmt.Println(documentsResult)
	}

	if verbose {
		fmt.Printf("Documents listing completed successfully for collection '%s' in vector database '%s'\n", collectionName, vdbName)
	}

	return nil
}

// listChunkingStrategies calls the MCP tool to retrieve supported chunking strategies
func listChunkingStrategies() error {
	if verbose {
		fmt.Println("Listing supported chunking strategies...")
	}

	if dryRun {
		fmt.Println("[DRY RUN] Would list chunking strategies")
		return nil
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

	// Use the raw call to the tool (no db context required)
	resp, err := client.callMCPServer("get_supported_chunking_strategies", nil)
	if err != nil {
		return fmt.Errorf("failed to call MCP tool: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("MCP server error: %s", resp.Error.Message)
	}
	if resp.Result == nil {
		return fmt.Errorf("no response from MCP server")
	}

	// Support both plain string responses and parsed JSON objects
	if resultStr, ok := resp.Result.(string); ok {
		fmt.Println(resultStr)
		return nil
	}

	// Attempt to pretty-print JSON-like results without HTML escaping
	// Note: resp.Result may already be a map[string]interface{} or []interface{}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resp.Result); err == nil {
		return nil
	}

	return fmt.Errorf("unexpected response format from MCP server: %T", resp.Result)
}
