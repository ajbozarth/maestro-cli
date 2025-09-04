package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCollectionCmd = &cobra.Command{
	Use:   "collection VDB_NAME COLLECTION_NAME",
	Short: "Create a collection in a vector database",
	Long: `Create a collection in a vector database.
	
Usage:
  maestro create collection VDB_NAME COLLECTION_NAME [options]

Examples:
  maestro create collection my-database my-collection
  maestro create collection my-database my-collection --embedding=text-embedding-3-small`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		return createCollection(vdbName, collectionName)
	},
}

var createVdbColCmd = &cobra.Command{
	Use:   "vdb-col VDB_NAME COLLECTION_NAME",
	Short: "Create a collection in a vector database (alias for collection)",
	Long: `Create a collection in a vector database.
	
Usage:
  maestro create vdb-col VDB_NAME COLLECTION_NAME [options]

Examples:
  maestro create vdb-col my-database my-collection
  maestro create vdb-col my-database my-collection --embedding=text-embedding-3-small`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		return createCollection(vdbName, collectionName)
	},
}

var createColCmd = &cobra.Command{
	Use:   "col VDB_NAME COLLECTION_NAME",
	Short: "Create a collection in a vector database (alias for collection)",
	Long: `Create a collection in a vector database.
	
Usage:
  maestro create col VDB_NAME COLLECTION_NAME [options]

Examples:
  maestro create col my-database my-collection
  maestro create col my-database my-collection --embedding=text-embedding-3-small`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		vdbName := args[0]
		collectionName := args[1]
		return createCollection(vdbName, collectionName)
	},
}

// Flags for collection creation - consistent chunking/embedding within a collection
var (
	collectionEmbedding     string
	collectionChunkStrategy string
	collectionChunkSize     int
	collectionChunkOverlap  int
	// Semantic-specific params
	semanticModel               string
	semanticWindowSize          int
	semanticThresholdPercentile float64
)

func init() {
	// Add flags for collection creation to all collection creation commands
	commands := []*cobra.Command{createCollectionCmd, createVdbColCmd, createColCmd}
	for _, cmd := range commands {
		cmd.Flags().StringVar(&collectionEmbedding, "embedding", "default", "Embedding model to use for the collection")
		cmd.Flags().StringVar(&collectionChunkStrategy, "chunking-strategy", "None", "Chunking strategy to use for the collection (None, Fixed, Sentence, Semantic)")
		cmd.Flags().IntVar(&collectionChunkSize, "chunk-size", 0, "Chunk size in characters (optional; defaults may vary by strategy)")
		cmd.Flags().IntVar(&collectionChunkOverlap, "chunk-overlap", 0, "Chunk overlap in characters (optional)")
	}
}

func createCollection(vdbName, collectionName string) error {
	if verbose && !silent {
		fmt.Printf("Creating collection '%s' in vector database '%s'...\n", collectionName, vdbName)
	}

	if dryRun {
		if !silent {
			fmt.Printf("[DRY RUN] Would create collection '%s' in vector database '%s' with embedding '%s'\n", collectionName, vdbName, collectionEmbedding)
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

	// Call the MCP server to create the collection with panic recovery
	var createErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Convert panic to a user-friendly error
				createErr = fmt.Errorf("MCP server could not be reached at %s. Please ensure the server is running and accessible", serverURI)
			}
		}()
		// Build chunking config if provided
		var chunkCfg map[string]interface{} = nil
		if collectionChunkStrategy != "None" {
			params := map[string]interface{}{}
			if collectionChunkSize > 0 {
				params["chunk_size"] = collectionChunkSize
			}
			if collectionChunkOverlap > 0 {
				params["overlap"] = collectionChunkOverlap
			}
			if collectionChunkStrategy == "Semantic" {
				if semanticWindowSize > 0 {
					params["window_size"] = semanticWindowSize
				}
				if semanticThresholdPercentile > 0 {
					params["threshold_percentile"] = semanticThresholdPercentile
				}
				if semanticModel != "" {
					// Use 'model_name' to align with server and semantic chunking API
					params["model_name"] = semanticModel
				}
			}
			chunkCfg = map[string]interface{}{
				"strategy":   collectionChunkStrategy,
				"parameters": params,
			}
		}

		createErr = client.CreateCollectionWithChunking(vdbName, collectionName, collectionEmbedding, chunkCfg)
	}()

	if createErr != nil {
		return fmt.Errorf("failed to create collection '%s' in vector database '%s': %w", collectionName, vdbName, createErr)
	}

	if !silent {
		fmt.Printf("✅ Collection '%s' created successfully in vector database '%s' with embedding '%s'\n", collectionName, vdbName, collectionEmbedding)
	}

	if verbose {
		fmt.Printf("Collection creation completed successfully for vector database '%s'\n", vdbName)
	}

	return nil
}
