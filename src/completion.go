package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CompletionItem represents a completion suggestion
type CompletionItem struct {
	Text        string
	Description string
	Type        string
}

// CompletionProvider provides completion suggestions
type CompletionProvider struct {
}

// NewCompletionProvider creates a new completion provider
func NewCompletionProvider() *CompletionProvider {
	return &CompletionProvider{}
}

// CompleteVectorDatabases provides completion for vector database names
func (cp *CompletionProvider) CompleteVectorDatabases(partial string) ([]CompletionItem, error) {
	// Get available vector databases
	databases, err := getAvailableVectorDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to get vector databases: %w", err)
	}

	var completions []CompletionItem
	for _, db := range databases {
		if strings.HasPrefix(strings.ToLower(db.Name), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        db.Name,
				Description: fmt.Sprintf("%s (%s) - %d documents", db.Name, db.Type, db.DocumentCount),
				Type:        "vectordb",
			})
		}
	}

	return completions, nil
}

// CompleteCollections provides completion for collection names
func (cp *CompletionProvider) CompleteCollections(vdbName, partial string) ([]CompletionItem, error) {
	// Get available collections
	collections, err := getAvailableCollections(vdbName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	var completions []CompletionItem
	for _, collection := range collections {
		if strings.HasPrefix(strings.ToLower(collection), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        collection,
				Description: fmt.Sprintf("Collection in %s", vdbName),
				Type:        "collection",
			})
		}
	}

	return completions, nil
}

// CompleteDocuments provides completion for document names
func (cp *CompletionProvider) CompleteDocuments(vdbName, collectionName, partial string) ([]CompletionItem, error) {
	// Get available documents
	documents, err := getAvailableDocuments(vdbName, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	var completions []CompletionItem
	for _, document := range documents {
		if strings.HasPrefix(strings.ToLower(document), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        document,
				Description: fmt.Sprintf("Document in %s/%s", vdbName, collectionName),
				Type:        "document",
			})
		}
	}

	return completions, nil
}

// CompleteFiles provides completion for file paths
func (cp *CompletionProvider) CompleteFiles(partial string) ([]CompletionItem, error) {
	// Parse the partial path
	dir, file := filepath.Split(partial)
	if dir == "" {
		dir = "."
	}

	// Get the absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var completions []CompletionItem
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless explicitly requested
		if !strings.HasPrefix(file, ".") && strings.HasPrefix(name, ".") {
			continue
		}

		// Check if the entry matches the partial name
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(file)) {
			fullPath := filepath.Join(dir, name)

			// Determine the type and description
			var itemType, description string
			if entry.IsDir() {
				itemType = "directory"
				description = fmt.Sprintf("Directory: %s", fullPath)
			} else {
				itemType = "file"
				// Check if it's a YAML file
				if strings.HasSuffix(strings.ToLower(name), ".yaml") || strings.HasSuffix(strings.ToLower(name), ".yml") {
					description = fmt.Sprintf("YAML file: %s", fullPath)
				} else {
					description = fmt.Sprintf("File: %s", fullPath)
				}
			}

			completions = append(completions, CompletionItem{
				Text:        fullPath,
				Description: description,
				Type:        itemType,
			})
		}
	}

	return completions, nil
}

// CompleteEmbeddings provides completion for embedding model names
func (cp *CompletionProvider) CompleteEmbeddings(partial string) ([]CompletionItem, error) {
	// Common embedding models
	embeddings := []string{
		"text-embedding-3-small",
		"text-embedding-3-large",
		"text-embedding-ada-002",
		"all-MiniLM-L6-v2",
		"all-mpnet-base-v2",
		"multi-qa-MiniLM-L6-cos-v1",
		"paraphrase-multilingual-MiniLM-L12-v2",
	}

	var completions []CompletionItem
	for _, embedding := range embeddings {
		if strings.HasPrefix(strings.ToLower(embedding), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        embedding,
				Description: "Embedding model",
				Type:        "embedding",
			})
		}
	}

	return completions, nil
}

// CompleteChunkingStrategies provides completion for chunking strategy values
func (cp *CompletionProvider) CompleteChunkingStrategies(partial string) ([]CompletionItem, error) {
	strategies := []string{"None", "Fixed", "Sentence", "Semantic"}

	var completions []CompletionItem
	for _, s := range strategies {
		if strings.HasPrefix(strings.ToLower(s), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        s,
				Description: "Chunking strategy",
				Type:        "chunking",
			})
		}
	}
	return completions, nil
}

// CompleteCommands provides completion for command names
func (cp *CompletionProvider) CompleteCommands(partial string) ([]CompletionItem, error) {
	commands := []string{
		"vectordb", "vdb",
		"collection", "coll",
		"document", "doc",
		"embedding", "embed",
		"chunking", "chunks",
		"status",
		"query",
		"validate",
	}

	var completions []CompletionItem
	for _, cmd := range commands {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        cmd,
				Description: "Command",
				Type:        "command",
			})
		}
	}

	return completions, nil
}

// CompleteSubcommands provides completion for subcommand names
func (cp *CompletionProvider) CompleteSubcommands(command, partial string) ([]CompletionItem, error) {
	var subcommands []string

	switch command {
	case "vectordb", "vdb":
		subcommands = []string{"list", "create", "delete"}
	case "collection", "coll":
		subcommands = []string{"list", "info", "create", "delete"}
	case "document", "doc":
		subcommands = []string{"list", "create", "delete"}
	case "embedding", "embed":
		subcommands = []string{"list"}
	case "chunking", "chunks":
		subcommands = []string{"list"}
	default:
		return nil, nil
	}

	var completions []CompletionItem
	for _, subcmd := range subcommands {
		if strings.HasPrefix(strings.ToLower(subcmd), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        subcmd,
				Description: fmt.Sprintf("Subcommand of %s", command),
				Type:        "subcommand",
			})
		}
	}

	return completions, nil
}

// CompleteFlags provides completion for flag names
func (cp *CompletionProvider) CompleteFlags(partial string) ([]CompletionItem, error) {
	flags := []string{
		"--vdb", "--collection", "--name", "--file", "--embedding",
		"--verbose", "--silent", "--dry-run", "--force",
		"--mcp-server-uri", "--doc-limit",
		"--chunking-strategy", "--chunk-size", "--chunk-overlap",
		"--semantic-model", "--semantic-window-size", "--semantic-threshold-percentile",
		"-h", "--help", "--version",
	}

	var completions []CompletionItem
	for _, flag := range flags {
		if strings.HasPrefix(strings.ToLower(flag), strings.ToLower(partial)) {
			completions = append(completions, CompletionItem{
				Text:        flag,
				Description: "Flag",
				Type:        "flag",
			})
		}
	}

	return completions, nil
}

// FormatCompletions formats completion items for display
func FormatCompletions(completions []CompletionItem) string {
	if len(completions) == 0 {
		return "No completions available"
	}

	var lines []string
	for _, item := range completions {
		lines = append(lines, fmt.Sprintf("%s\t%s", item.Text, item.Description))
	}

	return strings.Join(lines, "\n")
}

// GetCompletionsForContext provides completions based on the current context
func GetCompletionsForContext(args []string, currentWord string) ([]CompletionItem, error) {
	provider := NewCompletionProvider()

	// Determine the context based on arguments
	if len(args) == 0 {
		// No command yet, suggest commands
		return provider.CompleteCommands(currentWord)
	}

	command := args[0]
	if len(args) == 1 {
		// Command given, suggest subcommands
		return provider.CompleteSubcommands(command, currentWord)
	}

	// Check if we're in a flag context
	if strings.HasPrefix(currentWord, "-") {
		return provider.CompleteFlags(currentWord)
	}

	// Check if we're in a flag value context
	if len(args) >= 2 {
		prevArg := args[len(args)-2]
		switch prevArg {
		case "--vdb":
			return provider.CompleteVectorDatabases(currentWord)
		case "--collection":
			// Need to get VDB from previous args
			vdbName := extractVDBFromArgs(args)
			if vdbName != "" {
				return provider.CompleteCollections(vdbName, currentWord)
			}
		case "--name":
			// Could be collection or document name
			vdbName := extractVDBFromArgs(args)
			if vdbName != "" {
				return provider.CompleteCollections(vdbName, currentWord)
			}
		case "--file":
			return provider.CompleteFiles(currentWord)
		case "--embedding":
			return provider.CompleteEmbeddings(currentWord)
		case "--chunking-strategy":
			return provider.CompleteChunkingStrategies(currentWord)
		}
	}

	// Default: suggest files
	return provider.CompleteFiles(currentWord)
}

// extractVDBFromArgs extracts the VDB name from command arguments
func extractVDBFromArgs(args []string) string {
	for i, arg := range args {
		if arg == "--vdb" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
