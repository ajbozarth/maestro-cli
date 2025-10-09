// SPDX-License-Identifier: Apache-2.0
// Copyright © 2025 IBM

package common

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// DefaultLogDir is the default directory for log files
	DefaultLogDir string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		// Check if home directory is writable
		if _, err := os.Stat(homeDir); err == nil {
			info, err := os.Stat(homeDir)
			if err == nil && info.Mode().Perm()&(1<<(uint(7))) != 0 {
				DefaultLogDir = filepath.Join(homeDir, ".maestro", "logs")
			} else {
				DefaultLogDir = "./logs"
			}
		} else {
			DefaultLogDir = "./logs"
		}
	} else {
		DefaultLogDir = "./logs"
	}
}

// generateUUID generates a random UUID-like string
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// If we can't generate random bytes, use timestamp as fallback
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// FileLogger handles logging of workflow and agent activities to files
type FileLogger struct {
	LogDir  string
	loggers map[string]*zap.Logger // Map of workflow IDs to loggers
}

// NewFileLogger creates a new FileLogger instance
func NewFileLogger(logDir string) (*FileLogger, error) {
	dir := logDir
	if dir == "" {
		dir = DefaultLogDir
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	return &FileLogger{
		LogDir:  dir,
		loggers: make(map[string]*zap.Logger),
	}, nil
}

// createLogger creates a new zap logger for a specific workflow
func (l *FileLogger) createLogger(workflowID string) (*zap.Logger, error) {
	logPath := filepath.Join(l.LogDir, fmt.Sprintf("maestro_run_%s.jsonl", workflowID))

	// Create encoder config for JSON format
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       zapcore.OmitKey, // Omit log level as it's not in original format
		NameKey:        zapcore.OmitKey,
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     zapcore.OmitKey, // We'll use custom fields instead of message
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create file for logging
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create core with JSON encoder and file writer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	// Create logger
	return zap.New(core), nil
}

// getLogger gets or creates a logger for the specified workflow
func (l *FileLogger) getLogger(workflowID string) (*zap.Logger, error) {
	if logger, ok := l.loggers[workflowID]; ok {
		return logger, nil
	}

	logger, err := l.createLogger(workflowID)
	if err != nil {
		return nil, err
	}

	l.loggers[workflowID] = logger
	return logger, nil
}

// Close closes all loggers and releases resources
func (l *FileLogger) Close() {
	for _, logger := range l.loggers {
		// Sync ensures all buffered logs are written
		_ = logger.Sync()
	}
	l.loggers = make(map[string]*zap.Logger)
}

// GenerateWorkflowID generates a unique workflow ID
func (l *FileLogger) GenerateWorkflowID() string {
	return generateUUID()
}

// writeJSONLine writes a JSON line to the specified log file
// Kept for backward compatibility with tests
func (l *FileLogger) writeJSONLine(logPath string, data interface{}) error {
	// For backward compatibility with tests, use the direct file approach
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline to log file: %w", err)
	}

	return nil
}

// writeJSONLineWithZap writes a JSON line to the specified log file using zap
// This is an internal method used by the new implementation
func (l *FileLogger) writeJSONLineWithZap(logPath string, data interface{}) error {
	// Extract the workflow ID from the log path
	base := filepath.Base(logPath)
	// Expected format: maestro_run_{workflowID}.jsonl
	workflowID := ""
	prefix := "maestro_run_"
	suffix := ".jsonl"

	if len(base) > len(prefix) && strings.HasPrefix(base, prefix) && strings.HasSuffix(base, suffix) {
		workflowID = base[len(prefix) : len(base)-len(suffix)]
	} else {
		// If we can't extract the workflow ID, create a temporary logger
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       zapcore.OmitKey,
			NameKey:        zapcore.OmitKey,
			CallerKey:      zapcore.OmitKey,
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     zapcore.OmitKey,
			StacktraceKey:  zapcore.OmitKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer file.Close()

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			zap.InfoLevel,
		)

		logger := zap.New(core)
		defer logger.Sync()

		// Convert data to zap fields
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		var fields map[string]interface{}
		if err := json.Unmarshal(jsonData, &fields); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		zapFields := make([]zap.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}

		logger.Info("", zapFields...)
		return nil
	}

	// Get or create a logger for this workflow
	logger, err := l.getLogger(workflowID)
	if err != nil {
		return err
	}

	// Convert data to zap fields
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	var fields map[string]interface{}
	if err := json.Unmarshal(jsonData, &fields); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	logger.Info("", zapFields...)
	return nil
}

// TokenUsage represents token usage information
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

// AgentResponseLog represents a log entry for an agent response
type AgentResponseLog struct {
	LogType    string      `json:"log_type"`
	Timestamp  string      `json:"timestamp"`
	WorkflowID string      `json:"workflow_id"`
	StepIndex  int         `json:"step_index"`
	AgentName  string      `json:"agent_name"`
	Model      string      `json:"model"`
	Input      string      `json:"input"`
	Response   string      `json:"response"`
	ToolUsed   string      `json:"tool_used,omitempty"`
	StartTime  string      `json:"start_time,omitempty"`
	EndTime    string      `json:"end_time,omitempty"`
	DurationMS int64       `json:"duration_ms,omitempty"`
	TokenUsage *TokenUsage `json:"token_usage,omitempty"`
}

// LogAgentResponse logs an agent response
func (l *FileLogger) LogAgentResponse(
	workflowID string,
	stepIndex int,
	agentName string,
	model string,
	inputText string,
	responseText string,
	toolUsed string,
	startTime *time.Time,
	endTime *time.Time,
	durationMS int64,
	tokenUsage *TokenUsage,
) error {
	logPath := filepath.Join(l.LogDir, fmt.Sprintf("maestro_run_%s.jsonl", workflowID))

	var startTimeStr, endTimeStr string
	if startTime != nil {
		startTimeStr = startTime.UTC().Format(time.RFC3339Nano)
	}
	if endTime != nil {
		endTimeStr = endTime.UTC().Format(time.RFC3339Nano)
	}

	data := AgentResponseLog{
		LogType:    "agent_response",
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		WorkflowID: workflowID,
		StepIndex:  stepIndex,
		AgentName:  agentName,
		Model:      model,
		Input:      inputText,
		Response:   responseText,
		ToolUsed:   toolUsed,
		StartTime:  startTimeStr,
		EndTime:    endTimeStr,
		DurationMS: durationMS,
		TokenUsage: tokenUsage,
	}

	return l.writeJSONLine(logPath, data)
}

// WorkflowRunLog represents a log entry for a workflow run
type WorkflowRunLog struct {
	LogType      string   `json:"log_type"`
	Timestamp    string   `json:"timestamp"`
	WorkflowID   string   `json:"workflow_id"`
	WorkflowName string   `json:"workflow_name"`
	Status       string   `json:"status"`
	Prompt       string   `json:"prompt"`
	Output       string   `json:"output"`
	ModelsUsed   []string `json:"models_used"`
	StartTime    string   `json:"start_time,omitempty"`
	EndTime      string   `json:"end_time,omitempty"`
	DurationMS   int64    `json:"duration_ms,omitempty"`
}

// LogWorkflowRun logs a workflow run
func (l *FileLogger) LogWorkflowRun(
	workflowID string,
	workflowName string,
	prompt string,
	output string,
	modelsUsed []string,
	status string,
	startTime *time.Time,
	endTime *time.Time,
	durationMS int64,
) error {
	logPath := filepath.Join(l.LogDir, fmt.Sprintf("maestro_run_%s.jsonl", workflowID))

	var startTimeStr, endTimeStr string
	if startTime != nil {
		startTimeStr = startTime.UTC().Format(time.RFC3339Nano)
	}
	if endTime != nil {
		endTimeStr = endTime.UTC().Format(time.RFC3339Nano)
	}

	data := WorkflowRunLog{
		LogType:      "workflow_summary",
		Timestamp:    time.Now().UTC().Format(time.RFC3339Nano),
		WorkflowID:   workflowID,
		WorkflowName: workflowName,
		Status:       status,
		Prompt:       prompt,
		Output:       output,
		ModelsUsed:   modelsUsed,
		StartTime:    startTimeStr,
		EndTime:      endTimeStr,
		DurationMS:   durationMS,
	}

	return l.writeJSONLine(logPath, data)
}

// Made with Bob
