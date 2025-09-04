package main

import (
	"testing"
)

func TestNewCompletionProvider(t *testing.T) {
	provider := NewCompletionProvider()
	if provider == nil {
		t.Error("NewCompletionProvider() returned nil")
	}
}

func TestCompleteCommands(t *testing.T) {
	provider := NewCompletionProvider()

	// Test with empty partial
	completions, err := provider.CompleteCommands("")
	if err != nil {
		t.Errorf("CompleteCommands(\"\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteCommands(\"\") returned no completions")
	}

	// Test with partial match
	completions, err = provider.CompleteCommands("vect")
	if err != nil {
		t.Errorf("CompleteCommands(\"vect\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteCommands(\"vect\") returned no completions")
	}

	// Test with no match
	completions, err = provider.CompleteCommands("nonexistent")
	if err != nil {
		t.Errorf("CompleteCommands(\"nonexistent\") returned error: %v", err)
	}
	if len(completions) != 0 {
		t.Errorf("CompleteCommands(\"nonexistent\") returned %d completions, expected 0", len(completions))
	}
}

func TestCompleteSubcommands(t *testing.T) {
	provider := NewCompletionProvider()

	// Test vectordb subcommands
	completions, err := provider.CompleteSubcommands("vectordb", "")
	if err != nil {
		t.Errorf("CompleteSubcommands(\"vectordb\", \"\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteSubcommands(\"vectordb\", \"\") returned no completions")
	}

	// Test with partial match
	completions, err = provider.CompleteSubcommands("vectordb", "li")
	if err != nil {
		t.Errorf("CompleteSubcommands(\"vectordb\", \"li\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteSubcommands(\"vectordb\", \"li\") returned no completions")
	}

	// Test with alias
	completions, err = provider.CompleteSubcommands("vdb", "")
	if err != nil {
		t.Errorf("CompleteSubcommands(\"vdb\", \"\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteSubcommands(\"vdb\", \"\") returned no completions")
	}
}

func TestCompleteFlags(t *testing.T) {
	provider := NewCompletionProvider()

	// Test with empty partial
	completions, err := provider.CompleteFlags("")
	if err != nil {
		t.Errorf("CompleteFlags(\"\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteFlags(\"\") returned no completions")
	}

	// Test with partial match
	completions, err = provider.CompleteFlags("--v")
	if err != nil {
		t.Errorf("CompleteFlags(\"--v\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteFlags(\"--v\") returned no completions")
	}
}

func TestCompleteEmbeddings(t *testing.T) {
	provider := NewCompletionProvider()

	// Test with empty partial
	completions, err := provider.CompleteEmbeddings("")
	if err != nil {
		t.Errorf("CompleteEmbeddings(\"\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteEmbeddings(\"\") returned no completions")
	}

	// Test with partial match
	completions, err = provider.CompleteEmbeddings("text-embedding")
	if err != nil {
		t.Errorf("CompleteEmbeddings(\"text-embedding\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("CompleteEmbeddings(\"text-embedding\") returned no completions")
	}
}

func TestCompleteFiles(t *testing.T) {
	provider := NewCompletionProvider()

	// Test with current directory
	completions, err := provider.CompleteFiles("")
	if err != nil {
		t.Errorf("CompleteFiles(\"\") returned error: %v", err)
	}
	// Should return some files/directories in current directory
	if len(completions) == 0 {
		t.Log("CompleteFiles(\"\") returned no completions (this might be normal in test environment)")
	}
}

func TestFormatCompletions(t *testing.T) {
	completions := []CompletionItem{
		{Text: "test1", Description: "Description 1", Type: "test"},
		{Text: "test2", Description: "Description 2", Type: "test"},
	}

	result := FormatCompletions(completions)
	if result == "" {
		t.Error("FormatCompletions() returned empty string")
	}

	// Test with empty completions
	result = FormatCompletions([]CompletionItem{})
	if result != "No completions available" {
		t.Errorf("FormatCompletions([]) returned '%s', expected 'No completions available'", result)
	}
}

func TestGetCompletionsForContext(t *testing.T) {
	// Test with no args
	completions, err := GetCompletionsForContext([]string{}, "")
	if err != nil {
		t.Errorf("GetCompletionsForContext([], \"\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("GetCompletionsForContext([], \"\") returned no completions")
	}

	// Test with command only
	completions, err = GetCompletionsForContext([]string{"vectordb"}, "")
	if err != nil {
		t.Errorf("GetCompletionsForContext([\"vectordb\"], \"\") returned error: %v", err)
	}
	if len(completions) == 0 {
		t.Error("GetCompletionsForContext([\"vectordb\"], \"\") returned no completions")
	}
}

func TestExtractVDBFromArgs(t *testing.T) {
	// Test with VDB flag
	args := []string{"collection", "list", "--vdb", "test-db"}
	result := extractVDBFromArgs(args)
	if result != "test-db" {
		t.Errorf("extractVDBFromArgs(%v) returned '%s', expected 'test-db'", args, result)
	}

	// Test without VDB flag
	args = []string{"collection", "list"}
	result = extractVDBFromArgs(args)
	if result != "" {
		t.Errorf("extractVDBFromArgs(%v) returned '%s', expected ''", args, result)
	}

	// Test with VDB flag at end (should not match)
	args = []string{"collection", "list", "--vdb"}
	result = extractVDBFromArgs(args)
	if result != "" {
		t.Errorf("extractVDBFromArgs(%v) returned '%s', expected ''", args, result)
	}
}
