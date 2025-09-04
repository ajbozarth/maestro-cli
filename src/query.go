package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	docLimit int
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query documents using natural language",
	Long: `Query documents in a vector database using natural language.
	
This command allows you to ask questions about documents stored in a vector database.
The query agent will search through the documents and provide relevant answers.`,
	Example: `  maestro query "What is the main topic of the documents?" --vdb=my-vdb
  maestro query "Find information about API endpoints" --vdb=my-vdb --doc-limit 10`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		vdbName, _ := cmd.Flags().GetString("vdb")
		collectionName, _ := cmd.Flags().GetString("collection")

		// Validate inputs
		if strings.TrimSpace(query) == "" {
			return fmt.Errorf("query cannot be empty")
		}

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		// Use default collection if not specified
		if collectionName == "" {
			// If dry-run, don't interact; leave empty so downstream prints dry-run and exits
			if !dryRun {
				var err error
				collectionName, err = PromptForCollection(vdbName, collectionName)
				if err != nil {
					return fmt.Errorf("failed to select collection: %w", err)
				}
			}
		}

		return queryVectorDatabase(vdbName, query, collectionName)
	},
}

func queryVectorDatabase(dbName, query, collectionName string) error {
	// Initialize progress indicator
	var progress *ProgressIndicator
	if ShouldShowProgress() {
		progress = NewProgressIndicator("Processing query...")
		progress.Start()
	}

	if verbose {
		fmt.Println("Querying vector database...")
	}

	if dryRun {
		if progress != nil {
			progress.Stop("Dry run completed")
		}
		fmt.Println("[DRY RUN] Would query vector database")
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
		progress.Update("Executing query...")
	}

	// Call the query method
	result, err := client.Query(dbName, query, docLimit, collectionName)
	if err != nil {
		if progress != nil {
			progress.StopWithError("Query failed")
		}
		return fmt.Errorf("failed to query vector database: %w", err)
	}

	if progress != nil {
		progress.Stop("Query completed successfully")
	}

	if !silent {
		fmt.Println("OK")
	}
	fmt.Println(result)
	return nil
}

func init() {
	// Add flags to query command
	queryCmd.Flags().String("vdb", "", "Vector database name")
	queryCmd.Flags().String("collection", "", "Collection name to search in")
	queryCmd.Flags().IntVarP(&docLimit, "doc-limit", "d", 5, "Maximum number of documents to consider")
}
