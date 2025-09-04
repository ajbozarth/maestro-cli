package main

import (
	"fmt"
	"os"
	"strings"
)

// showStatus displays a quick overview of the current system status
func showStatus(vdbName string) error {
	// Initialize progress indicator
	var progress *ProgressIndicator
	if ShouldShowProgress() {
		progress = NewProgressIndicator("Gathering system status...")
		progress.Start()
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
		progress.Update("Retrieving vector databases...")
	}

	// Get vector databases
	databases, err := client.ListDatabases()
	if err != nil {
		if progress != nil {
			progress.StopWithError("Failed to list databases")
		}
		return fmt.Errorf("failed to list vector databases: %w", err)
	}

	if progress != nil {
		progress.Stop("Status retrieved successfully")
	}

	// Display status header
	fmt.Println("🔍 Maestro Knowledge System Status")
	fmt.Println("==================================")

	if len(databases) == 0 {
		fmt.Println("❌ No vector databases found")
		return nil
	}

	// Filter by specific VDB if provided
	if vdbName != "" {
		var found bool
		for _, db := range databases {
			if db.Name == vdbName {
				databases = []DatabaseInfo{db}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("vector database '%s' not found", vdbName)
		}
	}

	// Display each database status
	for i, db := range databases {
		if i > 0 {
			fmt.Println()
		}

		fmt.Printf("📊 Vector Database: %s (%s)\n", db.Name, db.Type)
		fmt.Printf("   📁 Collection: %s\n", db.Collection)
		fmt.Printf("   📄 Documents: %d\n", db.DocumentCount)

		// Get collections for this database
		if progress != nil {
			progress.Update(fmt.Sprintf("Retrieving collections for %s...", db.Name))
		}

		collectionsResult, err := client.ListCollections(db.Name)
		if err == nil {
			// Parse collections
			lines := strings.Split(strings.TrimSpace(collectionsResult), "\n")
			var collections []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "Found") && !strings.HasPrefix(line, "Collections") {
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

			if len(collections) > 0 {
				fmt.Printf("   📂 Collections: %s\n", strings.Join(collections, ", "))
			}
		}

		// Get supported embeddings
		if progress != nil {
			progress.Update(fmt.Sprintf("Retrieving embeddings for %s...", db.Name))
		}

		embeddingsResult, err := client.GetSupportedEmbeddings(db.Name)
		if err == nil {
			// Clean up embeddings result
			embeddings := strings.TrimSpace(embeddingsResult)
			if embeddings != "" {
				fmt.Printf("   🧠 Supported Embeddings: %s\n", embeddings)
			}
		}

		// Show health status
		fmt.Printf("   ✅ Status: Online\n")
	}

	// Show summary
	fmt.Println()
	fmt.Println("📈 Summary:")
	fmt.Printf("   • Total Vector Databases: %d\n", len(databases))

	totalDocuments := 0
	for _, db := range databases {
		totalDocuments += db.DocumentCount
	}
	fmt.Printf("   • Total Documents: %d\n", totalDocuments)

	// Show MCP server status
	fmt.Printf("   • MCP Server: %s\n", serverURI)
	fmt.Printf("   • Connection: ✅ Active\n")

	if verbose {
		fmt.Println()
		fmt.Println("🔧 Additional Information:")
		fmt.Printf("   • CLI Version: %s\n", version)
		fmt.Printf("   • Test Mode: %t\n", os.Getenv("MAESTRO_TEST_MODE") == "true")
	}

	return nil
}
