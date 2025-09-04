package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var createDocumentCmd = &cobra.Command{
	Use:   "document VDB_NAME COLLECTION_NAME DOC_NAME",
	Short: "Create a document in a collection of a vector database",
	Long: `Create a document in a collection of a vector database.
	
Usage:
  maestro create document VDB_NAME COLLECTION_NAME DOC_NAME --file-name=FILE_NAME [options]
  maestro create document VDB_NAME COLLECTION_NAME DOC_NAME --doc-file-name=FILE_NAME [options]

Examples:
	maestro create document my-database my-collection my-doc --file-name=document.txt
	# NOTE: --embed is deprecated and ignored; embedding is per collection`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		docName := args[2]
		return createDocument(vdbName, collectionName, docName)
	},
}

var createVdbDocCmd = &cobra.Command{
	Use:   "vdb-doc VDB_NAME COLLECTION_NAME DOC_NAME",
	Short: "Create a document in a collection of a vector database (alias for document)",
	Long: `Create a document in a collection of a vector database.
	
Usage:
  maestro create vdb-doc VDB_NAME COLLECTION_NAME DOC_NAME --file-name=FILE_NAME [options]

Examples:
	maestro create vdb-doc my-database my-collection my-doc --file-name=document.txt
	# NOTE: --embed is deprecated and ignored; embedding is per collection`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		docName := args[2]
		return createDocument(vdbName, collectionName, docName)
	},
}

var createDocCmd = &cobra.Command{
	Use:   "doc VDB_NAME COLLECTION_NAME DOC_NAME",
	Short: "Create a document in a collection of a vector database (alias for document)",
	Long: `Create a document in a collection of a vector database.
	
Usage:
  maestro create doc VDB_NAME COLLECTION_NAME DOC_NAME --file-name=FILE_NAME [options]

Examples:
	maestro create doc my-database my-collection my-doc --file-name=document.txt
	# NOTE: --embed is deprecated and ignored; embedding is per collection`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		docName := args[2]
		return createDocument(vdbName, collectionName, docName)
	},
}

// Flags for document creation
var (
	documentFileName  string
	documentEmbedding string
)

// getDocumentFileName returns the file name from either flag
func getDocumentFileName() string {
	if documentFileName != "" {
		return documentFileName
	}
	return ""
}

func init() {
	// Add flags for document creation to all document creation commands
	commands := []*cobra.Command{createDocumentCmd, createVdbDocCmd, createDocCmd}
	for _, cmd := range commands {
		cmd.Flags().StringVar(&documentFileName, "file-name", "", "File name containing the document content")
		cmd.Flags().StringVar(&documentFileName, "doc-file-name", "", "File name containing the document content (alias for file-name)")
		// DEPRECATED: embedding per document is ignored; kept temporarily for compatibility
		cmd.Flags().StringVar(&documentEmbedding, "embed", "default", "(DEPRECATED) Embedding model for the document (ignored; embedding is per collection)")
	}
}

func createDocument(vdbName, collectionName, docName string) error {
	// Initialize progress indicator
	var progress *ProgressIndicator
	if ShouldShowProgress() {
		progress = NewProgressIndicator("Creating document...")
		progress.Start()
	}

	// Validate that we have a file name
	fileName := getDocumentFileName()
	if fileName == "" {
		if progress != nil {
			progress.StopWithError("File name is required")
		}
		return fmt.Errorf("file name is required (use --file-name or --doc-file-name)")
	}

	// Check if file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if progress != nil {
			progress.StopWithError("File not found")
		}
		return fmt.Errorf("file not found: %s", fileName)
	}

	if dryRun {
		if !silent {
			fmt.Printf("[DRY RUN] Would create document '%s' in collection '%s' of vector database '%s' from file '%s'\n", docName, collectionName, vdbName, fileName)
		}
		if progress != nil {
			progress.Stop("Dry run completed")
		}
		return nil
	}

	if progress != nil {
		progress.Update("Connecting to MCP server...")
	}

	// Get MCP server URI
	serverURI, err := getMCPServerURI(mcpServerURI)
	if err != nil {
		if progress != nil {
			progress.StopWithError("Failed to get MCP server URI")
		}
		return fmt.Errorf("failed to get MCP server URI: %w", err)
	}

	if verbose {
		fmt.Printf("Connecting to MCP server at: %s\n", serverURI)
	}

	// Create MCP client
	client, err := NewMCPClient(serverURI)
	if err != nil {
		if progress != nil {
			progress.StopWithError("Failed to create MCP client")
		}
		return fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	if progress != nil {
		progress.Update("Validating database...")
	}

	// Check if the database exists first
	exists, err := client.DatabaseExists(vdbName)
	if err != nil {
		if progress != nil {
			progress.StopWithError("Failed to check database")
		}
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		if progress != nil {
			progress.StopWithError("Database does not exist")
		}
		return fmt.Errorf("vector database '%s' does not exist. Please create it first", vdbName)
	}

	if progress != nil {
		progress.Update("Validating collection...")
	}

	// Check if the collection exists
	collectionsResult, err := client.ListCollections(vdbName)
	if err != nil {
		if progress != nil {
			progress.StopWithError("Failed to list collections")
		}
		return fmt.Errorf("failed to list collections: %w", err)
	}

	// Simple check if collection exists in the result string
	if !strings.Contains(strings.ToLower(collectionsResult), strings.ToLower(collectionName)) {
		if progress != nil {
			progress.StopWithError("Collection does not exist")
		}
		return fmt.Errorf("collection '%s' does not exist in vector database '%s'. Please create it first", collectionName, vdbName)
	}

	// Deprecated: warn if user passed a non-default embedding for document creation
	if documentEmbedding != "default" && !silent {
		fmt.Printf("Warning: --embed is deprecated and ignored on document creation; embedding is configured per collection.\n")
	}

	if progress != nil {
		progress.Update("Checking for existing document...")
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
			if progress != nil {
				progress.StopWithError("Document already exists")
			}
			return fmt.Errorf("document '%s' already exists in collection '%s' of vector database '%s'", docName, collectionName, vdbName)
		}
	}

	if progress != nil {
		progress.Update("Creating document...")
	}

	// Call the MCP server to create the document with panic recovery
	var createErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				createErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		createErr = client.WriteDocument(vdbName, collectionName, docName, fileName, documentEmbedding)
	}()

	if createErr != nil {
		if progress != nil {
			progress.StopWithError("Failed to create document")
		}
		return fmt.Errorf("failed to create document '%s' in collection '%s' of vector database '%s': %w", docName, collectionName, vdbName, createErr)
	}

	if progress != nil {
		progress.Stop("Document created successfully")
	}

	if !silent {
		fmt.Printf("✅ Document '%s' created successfully in collection '%s' of vector database '%s' with embedding '%s'\n", docName, collectionName, vdbName, documentEmbedding)
	}

	if verbose {
		fmt.Printf("Document creation completed successfully for collection '%s' in vector database '%s'\n", collectionName, vdbName)
	}

	return nil
}
