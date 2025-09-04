package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var writeDocumentCmd = &cobra.Command{
	Use:   "document VDB_NAME COLLECTION_NAME DOC_NAME",
	Short: "Write a document to a collection of a vector database",
	Long: `Write a document to a collection of a vector database.
	
Usage:
  maestro write document VDB_NAME COLLECTION_NAME DOC_NAME --file-name=FILE_NAME [options]
  maestro write document VDB_NAME COLLECTION_NAME DOC_NAME --doc-file-name=FILE_NAME [options]

Examples:
	maestro write document my-database my-collection my-doc --file-name=document.txt
	# NOTE: --embed is deprecated and ignored; embedding is per collection`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		docName := args[2]
		return writeDocument(vdbName, collectionName, docName)
	},
}

var writeVdbDocCmd = &cobra.Command{
	Use:   "vdb-doc VDB_NAME COLLECTION_NAME DOC_NAME",
	Short: "Write a document to a collection of a vector database (alias for document)",
	Long: `Write a document to a collection of a vector database.
	
Usage:
  maestro write vdb-doc VDB_NAME COLLECTION_NAME DOC_NAME --file-name=FILE_NAME [options]

Examples:
	maestro write vdb-doc my-database my-collection my-doc --file-name=document.txt
	# NOTE: --embed is deprecated and ignored; embedding is per collection`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		docName := args[2]
		return writeDocument(vdbName, collectionName, docName)
	},
}

var writeDocCmd = &cobra.Command{
	Use:   "doc VDB_NAME COLLECTION_NAME DOC_NAME",
	Short: "Write a document to a collection of a vector database (alias for document)",
	Long: `Write a document to a collection of a vector database.
	
Usage:
  maestro write doc VDB_NAME COLLECTION_NAME DOC_NAME --file-name=FILE_NAME [options]

Examples:
	maestro write doc my-database my-collection my-doc --file-name=document.txt
	# NOTE: --embed is deprecated and ignored; embedding is per collection`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		docName := args[2]
		return writeDocument(vdbName, collectionName, docName)
	},
}

// Flags for document writing
var (
	writeDocumentFileName  string
	writeDocumentEmbedding string
)

// getWriteDocumentFileName returns the file name from either flag
func getWriteDocumentFileName() string {
	if writeDocumentFileName != "" {
		return writeDocumentFileName
	}
	return ""
}

func init() {
	// Add flags for document writing to all document writing commands
	commands := []*cobra.Command{writeDocumentCmd, writeVdbDocCmd, writeDocCmd}
	for _, cmd := range commands {
		cmd.Flags().StringVar(&writeDocumentFileName, "file-name", "", "File name containing the document content")
		cmd.Flags().StringVar(&writeDocumentFileName, "doc-file-name", "", "File name containing the document content (alias for file-name)")
		// DEPRECATED: embedding per document is ignored; kept temporarily for compatibility
		cmd.Flags().StringVar(&writeDocumentEmbedding, "embed", "default", "(DEPRECATED) Embedding model for the document (ignored; embedding is per collection)")
	}
}

func writeDocument(vdbName, collectionName, docName string) error {
	if verbose && !silent {
		fmt.Printf("Writing document '%s' to collection '%s' of vector database '%s'...\n", docName, collectionName, vdbName)
	}

	// Validate that we have a file name
	fileName := getWriteDocumentFileName()
	if fileName == "" {
		return fmt.Errorf("file name is required (use --file-name or --doc-file-name)")
	}

	// Check if file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", fileName)
	}

	if dryRun {
		if !silent {
			fmt.Printf("[DRY RUN] Would write document '%s' to collection '%s' of vector database '%s' from file '%s'\n", docName, collectionName, vdbName, fileName)
		}
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

	// Check if the collection exists
	collectionsResult, err := client.ListCollections(vdbName)
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	// Simple check if collection exists in the result string
	if !strings.Contains(strings.ToLower(collectionsResult), strings.ToLower(collectionName)) {
		return fmt.Errorf("collection '%s' does not exist in vector database '%s'. Please create it first", collectionName, vdbName)
	}

	// Deprecated: warn if user passed a non-default embedding for document writes
	if writeDocumentEmbedding != "default" && !silent {
		fmt.Printf("Warning: --embed is deprecated and ignored on writes; embedding is configured per collection.\n")
	}

	// Check if document already exists (simple check by listing documents)
	documentsResult, err := client.ListDocumentsInCollection(vdbName, collectionName)
	if err != nil {
		// If we can't list documents, we'll proceed anyway
		if verbose {
			fmt.Printf("Warning: Could not check for existing documents: %v\n", err)
		}
	} else {
		// Simple check if document name exists in the result
		if strings.Contains(strings.ToLower(documentsResult), strings.ToLower(docName)) {
			return fmt.Errorf("document '%s' already exists in collection '%s' of vector database '%s'", docName, collectionName, vdbName)
		}
	}

	// Call the MCP server to write the document with panic recovery
	var writeErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				writeErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		writeErr = client.WriteDocument(vdbName, collectionName, docName, fileName, writeDocumentEmbedding)
	}()

	if writeErr != nil {
		return fmt.Errorf("failed to write document '%s' to collection '%s' of vector database '%s': %w", docName, collectionName, vdbName, writeErr)
	}

	if !silent {
		fmt.Printf("✅ Document '%s' written successfully to collection '%s' of vector database '%s'\n", docName, collectionName, vdbName)
	}

	if verbose {
		fmt.Printf("Document writing completed successfully for collection '%s' in vector database '%s'\n", collectionName, vdbName)
	}

	return nil
}
