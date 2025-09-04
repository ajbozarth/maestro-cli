package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InteractiveSelector provides interactive selection functionality
type InteractiveSelector struct {
	reader *bufio.Reader
}

// NewInteractiveSelector creates a new interactive selector
func NewInteractiveSelector() *InteractiveSelector {
	return &InteractiveSelector{
		reader: bufio.NewReader(os.Stdin),
	}
}

// SelectVectorDatabase allows users to interactively select a vector database
func (is *InteractiveSelector) SelectVectorDatabase() (string, error) {
	// Get available vector databases
	databases, err := getAvailableVectorDatabases()
	if err != nil {
		return "", fmt.Errorf("failed to get vector databases: %w", err)
	}

	if len(databases) == 0 {
		return "", fmt.Errorf("no vector databases available")
	}

	// Display options
	fmt.Println("Available vector databases:")
	for i, db := range databases {
		fmt.Printf("%d. %s (%s) - %d documents\n", i+1, db.Name, db.Type, db.DocumentCount)
	}

	// Get user selection
	selection, err := is.getSelection(len(databases))
	if err != nil {
		return "", err
	}

	return databases[selection-1].Name, nil
}

// SelectCollection allows users to interactively select a collection from a vector database
func (is *InteractiveSelector) SelectCollection(vdbName string) (string, error) {
	// Get available collections
	collections, err := getAvailableCollections(vdbName)
	if err != nil {
		return "", fmt.Errorf("failed to get collections: %w", err)
	}

	if len(collections) == 0 {
		return "", fmt.Errorf("no collections available in vector database '%s'", vdbName)
	}

	// Display options
	fmt.Printf("Available collections in '%s':\n", vdbName)
	for i, collection := range collections {
		fmt.Printf("%d. %s\n", i+1, collection)
	}

	// Get user selection
	selection, err := is.getSelection(len(collections))
	if err != nil {
		return "", err
	}

	return collections[selection-1], nil
}

// SelectDocument allows users to interactively select a document from a collection
func (is *InteractiveSelector) SelectDocument(vdbName, collectionName string) (string, error) {
	// Get available documents
	documents, err := getAvailableDocuments(vdbName, collectionName)
	if err != nil {
		return "", fmt.Errorf("failed to get documents: %w", err)
	}

	if len(documents) == 0 {
		return "", fmt.Errorf("no documents available in collection '%s'", collectionName)
	}

	// Display options
	fmt.Printf("Available documents in '%s/%s':\n", vdbName, collectionName)
	for i, document := range documents {
		fmt.Printf("%d. %s\n", i+1, document)
	}

	// Get user selection
	selection, err := is.getSelection(len(documents))
	if err != nil {
		return "", err
	}

	return documents[selection-1], nil
}

// getSelection handles the user input for selection
func (is *InteractiveSelector) getSelection(maxOptions int) (int, error) {
	for {
		fmt.Printf("\nEnter your choice (1-%d): ", maxOptions)
		input, err := is.reader.ReadString('\n')
		if err != nil {
			return 0, fmt.Errorf("failed to read input: %w", err)
		}

		// Clean input
		input = strings.TrimSpace(input)

		// Handle empty input
		if input == "" {
			fmt.Println("Please enter a valid choice.")
			continue
		}

		// Parse selection
		selection, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Invalid input '%s'. Please enter a number between 1 and %d.\n", input, maxOptions)
			continue
		}

		// Validate range
		if selection < 1 || selection > maxOptions {
			fmt.Printf("Invalid choice %d. Please enter a number between 1 and %d.\n", selection, maxOptions)
			continue
		}

		return selection, nil
	}
}

// getAvailableVectorDatabases retrieves the list of available vector databases
func getAvailableVectorDatabases() ([]DatabaseInfo, error) {
	// Get MCP server URI
	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	// Create MCP client
	client, err := NewMCPClient(serverURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	// Call the MCP server to list databases
	var databases []DatabaseInfo
	var listErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				listErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		databases, listErr = client.ListDatabases()
	}()

	if listErr != nil {
		return nil, fmt.Errorf("failed to list vector databases: %w", listErr)
	}

	return databases, nil
}

// getAvailableCollections retrieves the list of available collections for a vector database
func getAvailableCollections(vdbName string) ([]string, error) {
	// Get MCP server URI
	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	// Create MCP client
	client, err := NewMCPClient(serverURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	// Call the MCP server to list collections
	var collectionsResult string
	var listErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				listErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		collectionsResult, listErr = client.ListCollections(vdbName)
	}()

	if listErr != nil {
		return nil, fmt.Errorf("failed to list collections: %w", listErr)
	}

	// Parse the collections result into a slice
	// The result is typically a formatted string, we need to extract collection names
	// For now, let's split by newlines and clean up
	lines := strings.Split(strings.TrimSpace(collectionsResult), "\n")
	var collections []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Found") && !strings.HasPrefix(line, "Collections") {
			// Extract collection name from the line
			// Assuming format like "1. collection_name" or just "collection_name"
			if strings.Contains(line, ". ") {
				parts := strings.SplitN(line, ". ", 2)
				if len(parts) > 1 {
					collections = append(collections, strings.TrimSpace(parts[1]))
				}
			} else {
				collections = append(collections, line)
			}
		}
	}

	return collections, nil
}

// getAvailableDocuments retrieves the list of available documents in a collection
func getAvailableDocuments(vdbName, collectionName string) ([]string, error) {
	// Get MCP server URI
	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	// Create MCP client
	client, err := NewMCPClient(serverURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	// Call the MCP server to list documents
	var documentsResult string
	var listErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				listErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		documentsResult, listErr = client.ListDocumentsInCollection(vdbName, collectionName)
	}()

	if listErr != nil {
		return nil, fmt.Errorf("failed to list documents: %w", listErr)
	}

	// Parse the documents result into a slice
	// The result is typically a formatted string, we need to extract document names
	lines := strings.Split(strings.TrimSpace(documentsResult), "\n")
	var documents []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Found") && !strings.HasPrefix(line, "Documents") {
			// Extract document name from the line
			// Assuming format like "1. document_name" or just "document_name"
			if strings.Contains(line, ". ") {
				parts := strings.SplitN(line, ". ", 2)
				if len(parts) > 1 {
					documents = append(documents, strings.TrimSpace(parts[1]))
				}
			} else {
				documents = append(documents, line)
			}
		}
	}

	return documents, nil
}

// PromptForVectorDatabase prompts the user to select a vector database if not provided
func PromptForVectorDatabase(vdbName string) (string, error) {
	if vdbName != "" {
		return vdbName, nil
	}

	// Check if we're in non-interactive mode (e.g., CI/CD)
	if !isInteractive() {
		return "", fmt.Errorf("vector database name is required in non-interactive mode")
	}

	selector := NewInteractiveSelector()
	return selector.SelectVectorDatabase()
}

// PromptForCollection prompts the user to select a collection if not provided
func PromptForCollection(vdbName, collectionName string) (string, error) {
	if collectionName != "" {
		return collectionName, nil
	}

	// Check if we're in non-interactive mode
	if !isInteractive() {
		return "", fmt.Errorf("collection name is required in non-interactive mode")
	}

	selector := NewInteractiveSelector()
	return selector.SelectCollection(vdbName)
}

// PromptForDocument prompts the user to select a document if not provided
func PromptForDocument(vdbName, collectionName, documentName string) (string, error) {
	if documentName != "" {
		return documentName, nil
	}

	// Check if we're in non-interactive mode
	if !isInteractive() {
		return "", fmt.Errorf("document name is required in non-interactive mode")
	}

	selector := NewInteractiveSelector()
	return selector.SelectDocument(vdbName, collectionName)
}

// isInteractive checks if the current environment supports interactive input
func isInteractive() bool {
	// Check if we're in test mode
	if os.Getenv("MAESTRO_TEST_MODE") == "true" {
		return false
	}

	// Check if stdin is a terminal
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	// Check if we're in a terminal (not redirected from file)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
