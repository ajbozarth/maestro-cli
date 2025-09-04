package main

import (
	"testing"
)

func TestSuggestCommand(t *testing.T) {
	tests := []struct {
		input       string
		expectCount int
		expectFirst string
	}{
		{"vectordb", 1, "vectordb"},
		{"vdb", 1, "vectordb"},
		{"coll", 1, "collection"},
		{"doc", 1, "document"},
		{"embed", 1, "embedding"},
		{"del", 1, "delete"},
		{"ls", 1, "list"},
		{"rm", 1, "delete"},
		{"add", 1, "create"},
		{"new", 1, "create"},
		{"get", 1, "list"},
		{"show", 1, "list"},
		{"search", 1, "query"},
		{"ask", 1, "query"},
		{"find", 1, "query"},
		{"vectord", 1, "vectordb"},
		{"collect", 1, "collection"},
		{"docum", 1, "document"},
		{"embedd", 1, "embedding"},
		{"xyz", 0, ""},
	}

	for _, tt := range tests {
		suggestions := SuggestCommand(tt.input)
		if len(suggestions) != tt.expectCount {
			t.Errorf("SuggestCommand(%q) returned %d suggestions, expected %d", tt.input, len(suggestions), tt.expectCount)
			// Debug output
			for i, s := range suggestions {
				t.Logf("  Suggestion %d: %s (score: %d)", i+1, s.Command, s.Score)
			}
		}
		if tt.expectCount > 0 && len(suggestions) > 0 {
			if suggestions[0].Command != tt.expectFirst {
				t.Errorf("SuggestCommand(%q) first suggestion was %q, expected %q", tt.input, suggestions[0].Command, tt.expectFirst)
			}
		}
	}
}

func TestSuggestForError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"flag is required", "Missing required flag. Use --help for usage information."},
		{"not found", "Resource not found. Use 'maestro vectordb list' to see available resources."},
		{"already exists", "Resource already exists. Use a different name or delete the existing resource first."},
		{"permission denied", "Insufficient permissions. Check your access rights or use --verbose for details."},
		{"connection refused", "Cannot connect to server. Ensure the MCP server is running and accessible."},
		{"invalid", "Invalid input. Check your command syntax and try again."},
		{"unknown error", ""},
	}

	for _, tt := range tests {
		result := SuggestForError(tt.input)
		if result != tt.expected {
			t.Errorf("SuggestForError(%q) returned %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestCalculateSimilarity(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"vectordb", "vectordb", 80},
		{"vectord", "vectordb", 80},
		{"vectordb", "vectord", 80},
		{"abc", "def", 0},
		{"", "", 80},
		{"a", "ab", 80},
		{"ab", "a", 80},
	}

	for _, tt := range tests {
		result := calculateSimilarity(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("calculateSimilarity(%q, %q) returned %d, expected %d", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestFormatSuggestions(t *testing.T) {
	suggestions := []CommandSuggestion{
		{Command: "vectordb", Score: 100, Category: "exact_match", Example: "maestro vectordb list"},
		{Command: "collection", Score: 80, Category: "similar", Example: "maestro collection list --vdb=my-vdb"},
	}

	result := FormatSuggestions(suggestions)
	expected := "\nDid you mean one of these?\n  1. vectordb (e.g., maestro vectordb list)\n  2. collection (e.g., maestro collection list --vdb=my-vdb)\n"

	if result != expected {
		t.Errorf("FormatSuggestions() returned %q, expected %q", result, expected)
	}
}

func TestFormatSuggestionsEmpty(t *testing.T) {
	suggestions := []CommandSuggestion{}
	result := FormatSuggestions(suggestions)
	if result != "" {
		t.Errorf("FormatSuggestions() with empty suggestions returned %q, expected empty string", result)
	}
}
