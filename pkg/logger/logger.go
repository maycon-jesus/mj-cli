package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"strings"
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
	file   io.WriteCloser
	attrs  Metadata
	groups []string
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
	data["time"] = time.Now().Format("2006-01-02T15:04:05Z07:00")

	maps.Copy(data, l.attrs)
	namespace := l.getNamespace()
	for _, meta := range metadata {
		for key, value := range meta {
			data[namespace+key] = value
		}
	}

	jsonData, _ := json.Marshal(data)
	jsonData = append(jsonData, '\n')
	l.file.Write(jsonData)
	fmt.Print(string(jsonData))
}

func New(writer io.WriteCloser) *Logger {
	return &Logger{file: writer, attrs: make(Metadata), groups: []string{}}
}

func NewWithTemporaryFile(appName string) (*Logger, error) {
	fileName := fmt.Sprintf("%s-log-*", appName)
	logFile, err := os.CreateTemp("", fileName)
	fileWriter := NewFileWriter(logFile)
	if err != nil {
		return nil, err
	}
	return New(fileWriter), nil
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

func (l *Logger) getNamespace() string {
	strBuilder := strings.Builder{}
	for _, group := range l.groups {
		strBuilder.WriteString(group)
		strBuilder.WriteString(".")
	}
	return strBuilder.String()
}

func (l *Logger) WithAttrs(attrs Metadata) *Logger {
	newLogger := &Logger{
		file:   l.file,
		attrs:  maps.Clone(l.attrs),
		groups: append([]string{}, l.groups...),
	}

	namespace := newLogger.getNamespace()
	for key, value := range attrs {
		newLogger.attrs[namespace+key] = value
	}
	return newLogger
}

func (l *Logger) WithGroup(group string) *Logger {
	newLogger := &Logger{
		file:   l.file,
		attrs:  maps.Clone(l.attrs),
		groups: append(append([]string{}, l.groups...), group),
	}
	return newLogger
}

func (l *Logger) Close() {
	l.file.Close()
}
