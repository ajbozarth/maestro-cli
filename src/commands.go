package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VDB Commands
var vdbCmd = &cobra.Command{
	Use:     "vectordb",
	Short:   "Manage vector databases",
	Long:    `Manage vector databases including listing, creating, and deleting them.`,
	Aliases: []string{"vdb"},
	Example: `  maestro vectordb list
  maestro vectordb create config.yaml
  maestro vectordb delete my-vdb`,
}

var vdbListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List vector databases",
	Long:    `List all available vector databases.`,
	Example: `  maestro vectordb list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listVectorDatabases()
	},
}

var vdbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a vector database",
	Long:  `Create a vector database from a YAML configuration file.`,
	Example: `  maestro vectordb create config.yaml
  maestro vectordb create config.yaml --type=milvus --uri=localhost:19530`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		yamlFile := args[0]
		return createVectorDatabase(yamlFile)
	},
}

var vdbDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a vector database",
	Long:  `Delete a vector database by name.`,
	Example: `  maestro vectordb delete my-vdb
  maestro vectordb delete my-vdb --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return deleteVectorDatabase(name)
	},
}

// Collection Commands
var collectionCmd = &cobra.Command{
	Use:     "collection",
	Short:   "Manage collections",
	Long:    `Manage collections within vector databases.`,
	Aliases: []string{"coll"},
	Example: `  maestro collection list --vdb=my-vdb
  maestro collection create --vdb=my-vdb --name=my-collection
  maestro collection delete my-collection --vdb=my-vdb`,
}

var collectionInfoCmd = &cobra.Command{
	Use:     "info",
	Short:   "Show collection info",
	Long:    `Show detailed information about a collection including document count, embedding, and chunking configuration.`,
	Example: `  maestro collection info --vdb=my-vdb --name=my-collection`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vdbName, _ := cmd.Flags().GetString("vdb")
		collectionName, _ := cmd.Flags().GetString("name")

		// Interactive selection if missing
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}
		if collectionName == "" {
			var err error
			collectionName, err = PromptForCollection(vdbName, collectionName)
			if err != nil {
				return fmt.Errorf("failed to select collection: %w", err)
			}
		}

		return showCollectionInfo(vdbName, collectionName)
	},
}

var collectionListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List collections",
	Long:    `List collections in a vector database.`,
	Example: `  maestro collection list --vdb=my-vdb`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vdbName, _ := cmd.Flags().GetString("vdb")

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		return listCollections(vdbName)
	},
}

var collectionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a collection",
	Long:  `Create a collection in a vector database.`,
	Example: `  maestro collection create --vdb=my-vdb --name=my-collection
  maestro collection create --vdb=my-vdb --name=my-collection --embedding=text-embedding-3-small`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vdbName, _ := cmd.Flags().GetString("vdb")
		collectionName, _ := cmd.Flags().GetString("name")

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		if collectionName == "" {
			return fmt.Errorf("--name flag is required")
		}
		return createCollection(vdbName, collectionName)
	},
}

var collectionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a collection",
	Long:  `Delete a collection from a vector database.`,
	Example: `  maestro collection delete my-collection --vdb=my-vdb
  maestro collection delete my-collection --vdb=my-vdb --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		collectionName := args[0]
		vdbName, _ := cmd.Flags().GetString("vdb")

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		return deleteCollection(vdbName, collectionName)
	},
}

// Document Commands
var documentCmd = &cobra.Command{
	Use:     "document",
	Short:   "Manage documents",
	Long:    `Manage documents within collections.`,
	Aliases: []string{"doc"},
	Example: `  maestro document list --vdb=my-vdb --collection=my-collection
  maestro document create --vdb=my-vdb --collection=my-collection --name=my-doc --file=./data.txt
  maestro document delete my-doc --vdb=my-vdb --collection=my-collection`,
}

var documentListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List documents",
	Long:    `List documents in a collection.`,
	Example: `  maestro document list --vdb=my-vdb --collection=my-collection`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vdbName, _ := cmd.Flags().GetString("vdb")
		collectionName, _ := cmd.Flags().GetString("collection")

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		// Use interactive selection if collection name is not provided
		if collectionName == "" {
			var err error
			collectionName, err = PromptForCollection(vdbName, collectionName)
			if err != nil {
				return fmt.Errorf("failed to select collection: %w", err)
			}
		}

		return listDocuments(vdbName, collectionName)
	},
}

var documentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a document",
	Long:  `Create a document in a collection.`,
	Example: `  maestro document create --vdb=my-vdb --collection=my-collection --name=my-doc --file=./data.txt
  maestro document create --vdb=my-vdb --collection=my-collection --name=my-doc --file=./data.txt --embedding=text-embedding-3-small`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vdbName, _ := cmd.Flags().GetString("vdb")
		collectionName, _ := cmd.Flags().GetString("collection")
		documentName, _ := cmd.Flags().GetString("name")
		fileName, _ := cmd.Flags().GetString("file")

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		// Use interactive selection if collection name is not provided
		if collectionName == "" {
			var err error
			collectionName, err = PromptForCollection(vdbName, collectionName)
			if err != nil {
				return fmt.Errorf("failed to select collection: %w", err)
			}
		}

		if documentName == "" {
			return fmt.Errorf("--name flag is required")
		}
		if fileName == "" {
			return fmt.Errorf("--file flag is required")
		}

		// Set the documentFileName variable that createDocument expects
		documentFileName = fileName

		return createDocument(vdbName, collectionName, documentName)
	},
}

var documentDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a document",
	Long:  `Delete a document from a collection.`,
	Example: `  maestro document delete my-doc --vdb=my-vdb --collection=my-collection
  maestro document delete my-doc --vdb=my-vdb --collection=my-collection --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		documentName := args[0]
		vdbName, _ := cmd.Flags().GetString("vdb")
		collectionName, _ := cmd.Flags().GetString("collection")

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		// Use interactive selection if collection name is not provided
		if collectionName == "" {
			var err error
			collectionName, err = PromptForCollection(vdbName, collectionName)
			if err != nil {
				return fmt.Errorf("failed to select collection: %w", err)
			}
		}

		return deleteDocument(vdbName, collectionName, documentName)
	},
}

// Embedding Commands
var embeddingCmd = &cobra.Command{
	Use:     "embedding",
	Short:   "Manage embeddings",
	Long:    `Manage embeddings within vector databases.`,
	Aliases: []string{"embed"},
	Example: `  maestro embedding list --vdb=my-vdb`,
}

var embeddingListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List embeddings",
	Long:    `List embeddings in a vector database.`,
	Example: `  maestro embedding list --vdb=my-vdb`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vdbName, _ := cmd.Flags().GetString("vdb")

		// Use interactive selection if vdb name is not provided
		if vdbName == "" {
			var err error
			vdbName, err = PromptForVectorDatabase(vdbName)
			if err != nil {
				return fmt.Errorf("failed to select vector database: %w", err)
			}
		}

		return listEmbeddings(vdbName)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show system status",
	Long:  `Show a quick overview of the current system status including vector databases, collections, and documents.`,
	Example: `  maestro status
  maestro status --vdb=my-vdb
  maestro status --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vdbName, _ := cmd.Flags().GetString("vdb")
		return showStatus(vdbName)
	},
}

// Chunking Commands
var chunkingCmd = &cobra.Command{
	Use:     "chunking",
	Short:   "Manage chunking strategies",
	Long:    `Discover and manage document chunking strategies.`,
	Aliases: []string{"chunks"},
	Example: `  maestro chunking list`,
}

var chunkingListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List supported chunking strategies",
	Long:    `List supported document chunking strategies and default parameters.`,
	Example: `  maestro chunking list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listChunkingStrategies()
	},
}

// Agent commands - create, serve
var agentCmd = &cobra.Command{
	Use:     "agent",
	Short:   "Manage agents",
	Long:    `Manage agent including creating, and serving.`,
	Aliases: []string{"agent"},
	Example: `  maestro create agents.yaml
  maestro agent create agents.yaml`,
}

// Workflow commands - run, serve, deploy
var workflowCmd = &cobra.Command{
	Use:     "workflow",
	Short:   "Manage workflow",
	Long:    "Manage workflow including running, serving and deploying.",
	Aliases: []string{"workflow"},
	Example: `  maestro workflow run agents.yaml workflow.yaml
  maestro workflow serve agents.yaml workflow.yaml
  maestro workflow deploy agents.yaml workflow.yaml`,
}

// CustomResource commands - create
var customResourceCmd = &cobra.Command{
	Use:     "customresource",
	Short:   "Manage custom resource",
	Long:    `Manage custom resource including creating.`,
	Aliases: []string{"cr"},
	Example: `  maestro customresource create yaml_file.`,
}

// MetaAgent commands - run
var metaAgentCmd = &cobra.Command{
	Use:     "metaagent",
	Short:   "Manage meta agent",
	Long:    `Manage meta agent including running.`,
	Aliases: []string{"metaagent"},
	Example: `  maestro metaagent run TEXT_FILE.`,
}

func init() {
	// Add flags to collection commands
	collectionListCmd.Flags().String("vdb", "", "Vector database name")
	collectionCreateCmd.Flags().String("vdb", "", "Vector database name")
	collectionCreateCmd.Flags().String("name", "", "Collection name")
	collectionDeleteCmd.Flags().String("vdb", "", "Vector database name")
	collectionInfoCmd.Flags().String("vdb", "", "Vector database name")
	collectionInfoCmd.Flags().String("name", "", "Collection name")

	// Add flags to document commands
	documentListCmd.Flags().String("vdb", "", "Vector database name")
	documentListCmd.Flags().String("collection", "", "Collection name")
	documentCreateCmd.Flags().String("vdb", "", "Vector database name")
	documentCreateCmd.Flags().String("collection", "", "Collection name")
	documentCreateCmd.Flags().String("name", "", "Document name")
	documentCreateCmd.Flags().String("file", "", "File path")
	documentDeleteCmd.Flags().String("vdb", "", "Vector database name")
	documentDeleteCmd.Flags().String("collection", "", "Collection name")

	// Add flags to embedding commands
	embeddingListCmd.Flags().String("vdb", "", "Vector database name")

	// Add flags to status command
	statusCmd.Flags().String("vdb", "", "Vector database name (optional, shows all if not specified)")

	// Add flags to vdb create command for overrides
	vdbCreateCmd.Flags().StringVar(&overrideType, "type", "", "Override the database type (milvus, weaviate)")
	vdbCreateCmd.Flags().StringVar(&overrideURI, "uri", "", "Override the connection URI")
	vdbCreateCmd.Flags().StringVar(&overrideCollectionName, "collection-name", "", "Override the collection name")
	vdbCreateCmd.Flags().StringVar(&overrideEmbedding, "embedding", "", "Override the embedding model")
	vdbCreateCmd.Flags().StringVar(&overrideMode, "mode", "", "Override the deployment mode (local, remote)")

	// Add flags to collection create command
	collectionCreateCmd.Flags().StringVar(&collectionEmbedding, "embedding", "default", "Embedding model to use for the collection")
	// Expose chunking flags here as well (kept in sync with create_collection.go)
	collectionCreateCmd.Flags().StringVar(&collectionChunkStrategy, "chunking-strategy", "None", "Chunking strategy to use for the collection (None, Fixed, Sentence, Semantic)")
	collectionCreateCmd.Flags().IntVar(&collectionChunkSize, "chunk-size", 0, "Chunk size in characters (optional; defaults may vary by strategy)")
	collectionCreateCmd.Flags().IntVar(&collectionChunkOverlap, "chunk-overlap", 0, "Chunk overlap in characters (optional)")
	// Semantic-specific parameters
	collectionCreateCmd.Flags().StringVar(&semanticModel, "semantic-model", "", "Semantic chunking model identifier (e.g., all-MiniLM-L6-v2)")
	collectionCreateCmd.Flags().IntVar(&semanticWindowSize, "semantic-window-size", 0, "Semantic chunking window size (optional)")
	collectionCreateCmd.Flags().Float64Var(&semanticThresholdPercentile, "semantic-threshold-percentile", 0, "Semantic chunking threshold percentile (0-100, optional)")
}
