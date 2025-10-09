// SPDX-License-Identifier: Apache-2.0
// Copyright © 2025 IBM

package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

// readJSONLine reads the first line from a JSONL file and unmarshals it into the provided interface
func readJSONLine(filePath string, v interface{}) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Find the first newline character
	data := content
	for i, b := range content {
		if b == '\n' {
			data = content[:i]
			break
		}
	}

	// Unmarshal the JSON data
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func TestNewFileLogger(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with custom log directory
	t.Run("CustomLogDir", func(t *testing.T) {
		logger, err := NewFileLogger(tempDir)
		if err != nil {
			t.Fatalf("NewFileLogger failed: %v", err)
		}
		if logger.LogDir != tempDir {
			t.Errorf("Expected LogDir to be %s, got %s", tempDir, logger.LogDir)
		}
		// Check that loggers map is initialized
		if logger.loggers == nil {
			t.Error("Expected loggers map to be initialized")
		}
	})

	// Test with default log directory
	t.Run("DefaultLogDir", func(t *testing.T) {
		// Save original DefaultLogDir
		origDefaultLogDir := DefaultLogDir
		defer func() { DefaultLogDir = origDefaultLogDir }()

		// Set DefaultLogDir to our temp directory for testing
		DefaultLogDir = tempDir

		logger, err := NewFileLogger("")
		if err != nil {
			t.Fatalf("NewFileLogger failed: %v", err)
		}
		if logger.LogDir != tempDir {
			t.Errorf("Expected LogDir to be %s, got %s", tempDir, logger.LogDir)
		}
	})

	// Test error handling when directory creation fails
	t.Run("DirectoryCreationFails", func(t *testing.T) {
		// Create a file with the same name as our intended directory
		// This will cause MkdirAll to fail
		filePath := filepath.Join(tempDir, "cannot_create_dir")
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()

		// Try to create a logger with a directory path that conflicts with the file
		_, err = NewFileLogger(filePath)
		if err == nil {
			t.Error("Expected an error when directory creation fails, got nil")
		}
	})
}

func TestGenerateWorkflowID(t *testing.T) {
	// Create a logger for testing
	tempDir, err := os.MkdirTemp("", "file_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logger, err := NewFileLogger(tempDir)
	if err != nil {
		t.Fatalf("NewFileLogger failed: %v", err)
	}

	// Generate multiple IDs and check they're unique
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id := logger.GenerateWorkflowID()

		// Check format (should be a hex string of length 32)
		matched, err := regexp.MatchString("^[0-9a-f]{32}$", id)
		if err != nil {
			t.Fatalf("Regex match failed: %v", err)
		}
		if !matched {
			t.Errorf("Generated ID %s doesn't match expected format", id)
		}

		// Check uniqueness
		if ids[id] {
			t.Errorf("Generated duplicate ID: %s", id)
		}

		ids[id] = true
	}
}

func TestWriteJSONLine(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logger, err := NewFileLogger(tempDir)
	if err != nil {
		t.Fatalf("NewFileLogger failed: %v", err)
	}

	// Test successful write
	t.Run("SuccessfulWrite", func(t *testing.T) {
		logPath := filepath.Join(tempDir, "test_log.jsonl")
		testData := map[string]string{"key": "value"}

		err := logger.writeJSONLine(logPath, testData)
		if err != nil {
			t.Fatalf("writeJSONLine failed: %v", err)
		}

		// Read the file and verify content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		expectedJSON, _ := json.Marshal(testData)
		expectedContent := string(expectedJSON) + "\n"
		if string(content) != expectedContent {
			t.Errorf("Expected content %q, got %q", expectedContent, string(content))
		}
	})

	// Test error handling for JSON marshaling
	t.Run("JSONMarshalError", func(t *testing.T) {
		logPath := filepath.Join(tempDir, "test_log_marshal_error.jsonl")

		// Create a struct with a channel which cannot be marshaled to JSON
		type UnmarshalableStruct struct {
			Ch chan int
		}
		unmarshalable := UnmarshalableStruct{Ch: make(chan int)}

		err := logger.writeJSONLine(logPath, unmarshalable)
		if err == nil {
			t.Error("Expected an error when marshaling invalid JSON, got nil")
		}
	})

	// Test error handling for file operations
	t.Run("FileOperationError", func(t *testing.T) {
		// Create a directory with the same name as our intended file
		// This will cause OpenFile to fail
		dirPath := filepath.Join(tempDir, "dir_not_file")
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		err = logger.writeJSONLine(dirPath, map[string]string{"key": "value"})
		if err == nil {
			t.Error("Expected an error when file operation fails, got nil")
		}
	})
}

func TestLogAgentResponse(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logger, err := NewFileLogger(tempDir)
	if err != nil {
		t.Fatalf("NewFileLogger failed: %v", err)
	}

	workflowID := "test-workflow-id"
	stepIndex := 1
	agentName := "test-agent"
	model := "test-model"
	inputText := "test input"
	responseText := "test response"
	toolUsed := "test-tool"

	// Test with all parameters provided
	t.Run("AllParametersProvided", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Minute)
		endTime := time.Now()
		durationMS := int64(endTime.Sub(startTime).Milliseconds())
		tokenUsage := &TokenUsage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		}

		err := logger.LogAgentResponse(
			workflowID,
			stepIndex,
			agentName,
			model,
			inputText,
			responseText,
			toolUsed,
			&startTime,
			&endTime,
			durationMS,
			tokenUsage,
		)
		if err != nil {
			t.Fatalf("LogAgentResponse failed: %v", err)
		}

		// Verify log file exists
		logPath := filepath.Join(logger.LogDir, "maestro_run_"+workflowID+".jsonl")
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Fatalf("Log file was not created: %v", err)
		}

		// Read and parse the log file
		var logEntry AgentResponseLog
		err = readJSONLine(logPath, &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse log entry: %v", err)
		}

		// Verify log entry fields
		if logEntry.LogType != "agent_response" {
			t.Errorf("Expected LogType to be 'agent_response', got %s", logEntry.LogType)
		}
		if logEntry.WorkflowID != workflowID {
			t.Errorf("Expected WorkflowID to be %s, got %s", workflowID, logEntry.WorkflowID)
		}
		if logEntry.StepIndex != stepIndex {
			t.Errorf("Expected StepIndex to be %d, got %d", stepIndex, logEntry.StepIndex)
		}
		if logEntry.AgentName != agentName {
			t.Errorf("Expected AgentName to be %s, got %s", agentName, logEntry.AgentName)
		}
		if logEntry.Model != model {
			t.Errorf("Expected Model to be %s, got %s", model, logEntry.Model)
		}
		if logEntry.Input != inputText {
			t.Errorf("Expected Input to be %s, got %s", inputText, logEntry.Input)
		}
		if logEntry.Response != responseText {
			t.Errorf("Expected Response to be %s, got %s", responseText, logEntry.Response)
		}
		if logEntry.ToolUsed != toolUsed {
			t.Errorf("Expected ToolUsed to be %s, got %s", toolUsed, logEntry.ToolUsed)
		}
		if logEntry.DurationMS != durationMS {
			t.Errorf("Expected DurationMS to be %d, got %d", durationMS, logEntry.DurationMS)
		}
		if logEntry.TokenUsage.PromptTokens != tokenUsage.PromptTokens {
			t.Errorf("Expected TokenUsage.PromptTokens to be %d, got %d", tokenUsage.PromptTokens, logEntry.TokenUsage.PromptTokens)
		}
	})

	// Test with nil optional parameters
	t.Run("NilOptionalParameters", func(t *testing.T) {
		// Remove previous log file if it exists
		logPath := filepath.Join(logger.LogDir, "maestro_run_"+workflowID+".jsonl")
		os.Remove(logPath)

		err := logger.LogAgentResponse(
			workflowID,
			stepIndex,
			agentName,
			model,
			inputText,
			responseText,
			toolUsed,
			nil, // startTime
			nil, // endTime
			0,   // durationMS
			nil, // tokenUsage
		)
		if err != nil {
			t.Fatalf("LogAgentResponse failed: %v", err)
		}

		// Verify log file exists
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Fatalf("Log file was not created: %v", err)
		}

		// Read and parse the log file
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		var logEntry AgentResponseLog
		err = json.Unmarshal(content, &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse log entry: %v", err)
		}

		// Verify log entry fields
		if logEntry.StartTime != "" {
			t.Errorf("Expected StartTime to be empty, got %s", logEntry.StartTime)
		}
		if logEntry.EndTime != "" {
			t.Errorf("Expected EndTime to be empty, got %s", logEntry.EndTime)
		}
		if logEntry.TokenUsage != nil {
			t.Errorf("Expected TokenUsage to be nil, got %+v", logEntry.TokenUsage)
		}
	})
}

func TestLogWorkflowRun(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logger, err := NewFileLogger(tempDir)
	if err != nil {
		t.Fatalf("NewFileLogger failed: %v", err)
	}

	workflowID := "test-workflow-id"
	workflowName := "test-workflow"
	prompt := "test prompt"
	output := "test output"
	modelsUsed := []string{"model1", "model2"}
	status := "completed"

	// Test with all parameters provided
	t.Run("AllParametersProvided", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Minute)
		endTime := time.Now()
		durationMS := int64(endTime.Sub(startTime).Milliseconds())

		err := logger.LogWorkflowRun(
			workflowID,
			workflowName,
			prompt,
			output,
			modelsUsed,
			status,
			&startTime,
			&endTime,
			durationMS,
		)
		if err != nil {
			t.Fatalf("LogWorkflowRun failed: %v", err)
		}

		// Verify log file exists
		logPath := filepath.Join(logger.LogDir, "maestro_run_"+workflowID+".jsonl")
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Fatalf("Log file was not created: %v", err)
		}

		// Read and parse the log file
		var logEntry WorkflowRunLog
		err = readJSONLine(logPath, &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse log entry: %v", err)
		}

		// Verify log entry fields
		if logEntry.LogType != "workflow_summary" {
			t.Errorf("Expected LogType to be 'workflow_summary', got %s", logEntry.LogType)
		}
		if logEntry.WorkflowID != workflowID {
			t.Errorf("Expected WorkflowID to be %s, got %s", workflowID, logEntry.WorkflowID)
		}
		if logEntry.WorkflowName != workflowName {
			t.Errorf("Expected WorkflowName to be %s, got %s", workflowName, logEntry.WorkflowName)
		}
		if logEntry.Status != status {
			t.Errorf("Expected Status to be %s, got %s", status, logEntry.Status)
		}
		if logEntry.Prompt != prompt {
			t.Errorf("Expected Prompt to be %s, got %s", prompt, logEntry.Prompt)
		}
		if logEntry.Output != output {
			t.Errorf("Expected Output to be %s, got %s", output, logEntry.Output)
		}
		if len(logEntry.ModelsUsed) != len(modelsUsed) {
			t.Errorf("Expected ModelsUsed length to be %d, got %d", len(modelsUsed), len(logEntry.ModelsUsed))
		}
		if logEntry.DurationMS != durationMS {
			t.Errorf("Expected DurationMS to be %d, got %d", durationMS, logEntry.DurationMS)
		}
	})

	// Test with nil optional parameters
	t.Run("NilOptionalParameters", func(t *testing.T) {
		// Remove previous log file if it exists
		logPath := filepath.Join(logger.LogDir, "maestro_run_"+workflowID+".jsonl")
		os.Remove(logPath)

		err := logger.LogWorkflowRun(
			workflowID,
			workflowName,
			prompt,
			output,
			modelsUsed,
			status,
			nil, // startTime
			nil, // endTime
			0,   // durationMS
		)
		if err != nil {
			t.Fatalf("LogWorkflowRun failed: %v", err)
		}

		// Verify log file exists
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Fatalf("Log file was not created: %v", err)
		}

		// Read and parse the log file
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		var logEntry WorkflowRunLog
		err = json.Unmarshal(content, &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse log entry: %v", err)
		}

		// Verify log entry fields
		if logEntry.StartTime != "" {
			t.Errorf("Expected StartTime to be empty, got %s", logEntry.StartTime)
		}
		if logEntry.EndTime != "" {
			t.Errorf("Expected EndTime to be empty, got %s", logEntry.EndTime)
		}
	})
}

// Made with Bob

func TestGetLogger(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logger, err := NewFileLogger(tempDir)
	if err != nil {
		t.Fatalf("NewFileLogger failed: %v", err)
	}
	defer logger.Close() // Clean up resources

	// Test getting a new logger
	t.Run("GetNewLogger", func(t *testing.T) {
		workflowID := "test-workflow"
		zapLogger, err := logger.getLogger(workflowID)
		if err != nil {
			t.Fatalf("getLogger failed: %v", err)
		}
		if zapLogger == nil {
			t.Error("Expected non-nil logger")
		}

		// Check that logger was cached
		if cachedLogger, ok := logger.loggers[workflowID]; !ok || cachedLogger != zapLogger {
			t.Error("Logger was not properly cached")
		}
	})

	// Test getting an existing logger
	t.Run("GetExistingLogger", func(t *testing.T) {
		workflowID := "test-workflow-2"
		firstLogger, err := logger.getLogger(workflowID)
		if err != nil {
			t.Fatalf("First getLogger failed: %v", err)
		}

		secondLogger, err := logger.getLogger(workflowID)
		if err != nil {
			t.Fatalf("Second getLogger failed: %v", err)
		}

		if firstLogger != secondLogger {
			t.Error("Expected same logger instance to be returned")
		}
	})
}

func TestClose(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logger, err := NewFileLogger(tempDir)
	if err != nil {
		t.Fatalf("NewFileLogger failed: %v", err)
	}

	// Create some loggers
	workflowIDs := []string{"workflow1", "workflow2", "workflow3"}
	for _, id := range workflowIDs {
		_, err := logger.getLogger(id)
		if err != nil {
			t.Fatalf("Failed to get logger for %s: %v", id, err)
		}
	}

	// Verify loggers exist
	if len(logger.loggers) != len(workflowIDs) {
		t.Errorf("Expected %d loggers, got %d", len(workflowIDs), len(logger.loggers))
	}

	// Close loggers
	logger.Close()

	// Verify loggers map is empty
	if len(logger.loggers) != 0 {
		t.Errorf("Expected empty loggers map after Close, got %d entries", len(logger.loggers))
	}
}
