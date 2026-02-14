package logger

import (
	"os"
	"strings"
	"testing"
)

func TestNewLoggerComplete(t *testing.T) {
	logger, err := NewLoggerComplete("testapp", "1.0.0")
	if err != nil {
		t.Fatalf("NewLoggerComplete() error: %v", err)
	}
	defer os.Remove(logger.FileHandler.Name())
	defer logger.FileHandler.Close()

	if logger.Log == nil {
		t.Error("Log should not be nil")
	}
	if logger.FileHandler == nil {
		t.Error("FileHandler should not be nil")
	}
}

func TestNewLoggerComplete_LogFileCreated(t *testing.T) {
	logger, err := NewLoggerComplete("testapp", "2.0.0")
	if err != nil {
		t.Fatalf("NewLoggerComplete() error: %v", err)
	}
	defer os.Remove(logger.FileHandler.Name())
	defer logger.FileHandler.Close()

	name := logger.FileHandler.Name()
	if !strings.Contains(name, "testapp-log-") {
		t.Errorf("log file name %q should contain 'testapp-log-'", name)
	}

	if _, err := os.Stat(name); os.IsNotExist(err) {
		t.Errorf("log file %q should exist", name)
	}
}

func TestNewLoggerComplete_WritesWithAppAttrs(t *testing.T) {
	logger, err := NewLoggerComplete("myapp", "3.0.0")
	if err != nil {
		t.Fatalf("NewLoggerComplete() error: %v", err)
	}
	defer os.Remove(logger.FileHandler.Name())
	defer logger.FileHandler.Close()

	logger.Log.Info("test log entry")

	data := readJSONLine(t, logger.FileHandler)

	if data["app"] != "myapp" {
		t.Errorf("app = %v, want 'myapp'", data["app"])
	}
	if data["version"] != "3.0.0" {
		t.Errorf("version = %v, want '3.0.0'", data["version"])
	}
	if data["message"] != "test log entry" {
		t.Errorf("message = %v, want 'test log entry'", data["message"])
	}
}

func TestNewLoggerComplete_DebugEnvOverridesLevel(t *testing.T) {
	t.Setenv("DEBUG", "debug")

	logger, err := NewLoggerComplete("testapp", "1.0.0")
	if err != nil {
		t.Fatalf("NewLoggerComplete() error: %v", err)
	}
	defer os.Remove(logger.FileHandler.Name())
	defer logger.FileHandler.Close()

	logger.Log.Debug("debug message")

	content, err := os.ReadFile(logger.FileHandler.Name())
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "debug message") {
		t.Error("debug message should be logged when DEBUG=debug")
	}
}

func TestNewLoggerComplete_DebugEnvInvalidFallsBackToDebug(t *testing.T) {
	t.Setenv("DEBUG", "invalid-level")

	logger, err := NewLoggerComplete("testapp", "1.0.0")
	if err != nil {
		t.Fatalf("NewLoggerComplete() error: %v", err)
	}
	defer os.Remove(logger.FileHandler.Name())
	defer logger.FileHandler.Close()

	logger.Log.Debug("fallback debug msg")

	content, err := os.ReadFile(logger.FileHandler.Name())
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "fallback debug msg") {
		t.Error("debug message should be logged when DEBUG has an invalid value (fallback to Debug)")
	}
}

func TestNewLoggerComplete_NoDebugEnvDefaultsToDebug(t *testing.T) {
	t.Setenv("DEBUG", "")

	logger, err := NewLoggerComplete("testapp", "1.0.0")
	if err != nil {
		t.Fatalf("NewLoggerComplete() error: %v", err)
	}
	defer os.Remove(logger.FileHandler.Name())
	defer logger.FileHandler.Close()

	logger.Log.Debug("should appear")

	content, err := os.ReadFile(logger.FileHandler.Name())
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "should appear") {
		t.Error("debug messages should be logged at default Debug level")
	}
}
