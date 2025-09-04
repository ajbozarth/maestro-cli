package main

import (
	"fmt"
	"strings"
	"unicode"
)

// CommandSuggestion represents a command suggestion with similarity score
type CommandSuggestion struct {
	Command  string
	Score    int
	Category string
	Example  string
}

// CommonCommandMistakes maps common typos to correct commands
var CommonCommandMistakes = map[string]string{
	"vectordb": "vectordb",
	"vdb":      "vectordb",
	"coll":     "collection",
	"doc":      "document",
	"embed":    "embedding",
	"del":      "delete",
	"ls":       "list",
	"rm":       "delete",
	"add":      "create",
	"new":      "create",
	"get":      "list",
	"show":     "list",
	"search":   "query",
	"ask":      "query",
	"find":     "query",
}

// CommandExamples provides examples for common commands
var CommandExamples = map[string][]string{
	"vectordb": {
		"maestro vectordb list",
		"maestro vectordb create config.yaml",
		"maestro vectordb delete my-vdb",
	},
	"collection": {
		"maestro collection list --vdb=my-vdb",
		"maestro collection create --vdb=my-vdb --name=my-collection",
		"maestro collection delete my-collection --vdb=my-vdb",
	},
	"document": {
		"maestro document list --vdb=my-vdb --collection=my-collection",
		"maestro document create --vdb=my-vdb --collection=my-collection --name=my-doc --file=./data.txt",
		"maestro document delete my-doc --vdb=my-vdb --collection=my-collection",
	},
	"embedding": {
		"maestro embedding list --vdb=my-vdb",
	},
	"query": {
		"maestro query \"What is the main topic?\" --vdb=my-vdb",
		"maestro query \"Find API documentation\" --vdb=my-vdb --collection=docs --doc-limit 10",
	},
	"validate": {
		"maestro validate config.yaml",
		"maestro validate custom-schema.json config.yaml",
	},
}

// CommonErrorPatterns maps error patterns to helpful suggestions
var CommonErrorPatterns = map[string]string{
	"flag is required":   "Missing required flag. Use --help for usage information.",
	"not found":          "Resource not found. Use 'maestro vectordb list' to see available resources.",
	"already exists":     "Resource already exists. Use a different name or delete the existing resource first.",
	"permission denied":  "Insufficient permissions. Check your access rights or use --verbose for details.",
	"connection refused": "Cannot connect to server. Ensure the MCP server is running and accessible.",
	"invalid":            "Invalid input. Check your command syntax and try again.",
}

// SuggestCommand provides command suggestions based on user input
func SuggestCommand(input string) []CommandSuggestion {
	var suggestions []CommandSuggestion
	seen := make(map[string]bool)

	// Normalize input
	input = strings.ToLower(strings.TrimSpace(input))

	// Check for exact matches in common mistakes
	if correct, exists := CommonCommandMistakes[input]; exists {
		suggestions = append(suggestions, CommandSuggestion{
			Command:  correct,
			Score:    100,
			Category: "exact_match",
			Example:  getRandomExample(correct),
		})
		seen[correct] = true
	}

	// Check for partial matches (avoid duplicates)
	for cmd := range CommandExamples {
		if seen[cmd] {
			continue // Skip if already added
		}
		score := calculateSimilarity(input, cmd)
		if score > 70 { // Only suggest if similarity is above 70%
			suggestions = append(suggestions, CommandSuggestion{
				Command:  cmd,
				Score:    score,
				Category: "similar",
				Example:  getRandomExample(cmd),
			})
			seen[cmd] = true
		}
	}

	// Sort by score (highest first)
	sortSuggestions(suggestions)

	// Return top 3 suggestions
	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return suggestions
}

// SuggestForError provides suggestions based on error messages
func SuggestForError(errMsg string) string {
	errMsg = strings.ToLower(errMsg)

	for pattern, suggestion := range CommonErrorPatterns {
		if strings.Contains(errMsg, pattern) {
			return suggestion
		}
	}

	return ""
}

// ShowCommandExamples displays examples for a specific command
func ShowCommandExamples(command string) {
	examples, exists := CommandExamples[command]
	if !exists {
		return
	}

	fmt.Println("\nExamples:")
	for i, example := range examples {
		fmt.Printf("  %d. %s\n", i+1, example)
	}
}

// ShowContextualHelp displays contextual help based on the current operation
func ShowContextualHelp(command, subcommand string) {
	switch command {
	case "vectordb":
		switch subcommand {
		case "list":
			fmt.Println("\n💡 Tip: Use 'maestro vectordb create <config.yaml>' to create a new vector database")
		case "create":
			fmt.Println("\n💡 Tip: Use 'maestro validate <config.yaml>' to validate your configuration before creating")
		case "delete":
			fmt.Println("\n💡 Tip: Use 'maestro vectordb list' to see all available vector databases")
		}
	case "collection":
		switch subcommand {
		case "list":
			fmt.Println("\n💡 Tip: Use 'maestro collection create --vdb=<vdb> --name=<name>' to create a new collection")
		case "create":
			fmt.Println("\n💡 Tip: Use 'maestro document create' to add documents to your new collection")
		case "delete":
			fmt.Println("\n💡 Tip: This will permanently delete the collection and all its documents")
		}
	case "document":
		switch subcommand {
		case "list":
			fmt.Println("\n💡 Tip: Use 'maestro document create' to add new documents to this collection")
		case "create":
			fmt.Println("\n💡 Tip: Use 'maestro query \"<question>\" --vdb=<vdb>' to search your documents")
		case "delete":
			fmt.Println("\n💡 Tip: This will permanently delete the document")
		}
	case "query":
		fmt.Println("\n💡 Tip: Use --doc-limit to control how many documents are considered in the search")
		fmt.Println("💡 Tip: Use --collection to search in a specific collection")
	}
}

// calculateSimilarity calculates similarity between two strings using Levenshtein distance
func calculateSimilarity(s1, s2 string) int {
	// Simple similarity calculation based on common prefixes and character overlap
	if strings.HasPrefix(s2, s1) || strings.HasPrefix(s1, s2) {
		return 80
	}

	// Count common characters
	common := 0
	for _, c1 := range s1 {
		for _, c2 := range s2 {
			if unicode.ToLower(c1) == unicode.ToLower(c2) {
				common++
				break
			}
		}
	}

	// Calculate similarity percentage
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	if maxLen == 0 {
		return 0
	}

	return (common * 100) / maxLen
}

// sortSuggestions sorts suggestions by score in descending order
func sortSuggestions(suggestions []CommandSuggestion) {
	for i := 0; i < len(suggestions)-1; i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[i].Score < suggestions[j].Score {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}
}

// getRandomExample returns a random example for a command
func getRandomExample(command string) string {
	examples, exists := CommandExamples[command]
	if !exists || len(examples) == 0 {
		return ""
	}

	// For simplicity, return the first example
	// In a more sophisticated implementation, you could use a random number generator
	return examples[0]
}

// FormatSuggestions formats suggestions for display
func FormatSuggestions(suggestions []CommandSuggestion) string {
	if len(suggestions) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString("\nDid you mean one of these?\n")

	for i, suggestion := range suggestions {
		result.WriteString(fmt.Sprintf("  %d. %s", i+1, suggestion.Command))
		if suggestion.Example != "" {
			result.WriteString(fmt.Sprintf(" (e.g., %s)", suggestion.Example))
		}
		result.WriteString("\n")
	}

	return result.String()
}
