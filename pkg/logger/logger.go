package logger

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"time"
)

const (
	LogLevelTrace = 1 + iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type Logger struct {
	file *os.File
}

type Metadata = map[string]interface{}

func getLevelString(level int) string {
	switch level {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) Log(level int, message string, metadata ...Metadata) {
	data := make(Metadata)
	data["level"] = getLevelString(level)
	data["message"] = message
	data["timestamp"] = time.Now().Format("2006-01-02T15:04:05Z07:00")

	for _, meta := range metadata {
		maps.Copy(data, meta)
	}

	jsonData, _ := json.Marshal(data)
	l.file.WriteString(string(jsonData) + "\n")
	fmt.Println(string(jsonData))
}

func New(file *os.File) *Logger {
	return &Logger{file: file}
}

func NewWithTemporaryFile(appName string) (*Logger, error) {
	fileName := fmt.Sprintf("%s-log-*", appName)
	logFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}
	return New(logFile), nil
}

func (l *Logger) Trace(message string, metadata ...Metadata) {
	l.Log(LogLevelTrace, message, metadata...)
}

func (l *Logger) Debug(message string, metadata ...Metadata) {
	l.Log(LogLevelDebug, message, metadata...)
}

func (l *Logger) Info(message string, metadata ...Metadata) {
	l.Log(LogLevelInfo, message, metadata...)
}

func (l *Logger) Warn(message string, metadata ...Metadata) {
	l.Log(LogLevelWarn, message, metadata...)
}

func (l *Logger) Error(message string, metadata ...Metadata) {
	l.Log(LogLevelError, message, metadata...)
}

func (l *Logger) Fatal(message string, metadata ...Metadata) {
	l.Log(LogLevelFatal, message, metadata...)
}

func (l *Logger) RecoverPanic() {
	if r := recover(); r != nil {
		l.Fatal(fmt.Sprintf("Panic recovered: %v", r))
	}
}

func (l *Logger) Close() {
	l.file.Close()
}
