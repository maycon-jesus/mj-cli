package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	LogLevelTrace LogLevel = iota + 1
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelNone
)

type FileInfoWriter interface {
	io.WriteCloser
	Name() string
}

type Metadata = map[string]interface{}

type Logger struct {
	file         FileInfoWriter
	attrs        Metadata
	groups       []string
	fileLevel    LogLevel
	consoleLevel LogLevel
	mu           *sync.Mutex
	stdout       io.Writer
	closeOnce    *sync.Once
}

func New(file FileInfoWriter, defaultStdout io.Writer) *Logger {
	return &Logger{
		file:         file,
		attrs:        make(Metadata),
		groups:       []string{},
		fileLevel:    LogLevelTrace,
		consoleLevel: LogLevelNone,
		mu:           &sync.Mutex{},
		stdout:       defaultStdout,
		closeOnce:    &sync.Once{},
	}
}

func NewWithTemporaryFile(appName string, defaultStdout io.Writer) (*Logger, error) {
	fileName := fmt.Sprintf("%s-log-*", appName)
	logFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}
	return New(logFile, defaultStdout), nil
}

func (l *Logger) log(level LogLevel, message string, metadata ...Metadata) {
	l.mu.Lock()
	defer l.mu.Unlock()

	data := make(Metadata)

	maps.Copy(data, l.attrs)
	namespace := l.getNamespace()
	for _, meta := range metadata {
		for key, value := range meta {
			data[namespace+key] = value
		}
	}

	// Dados internos do logger
	levelString, err := getLevelString(level)
	if err != nil {
		fmt.Fprintf(l.stdout, "Invalid log level: %v\n", err)
		return
	}
	data["level"] = levelString
	data["message"] = message
	data["time"] = time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")

	_, file, line, ok := runtime.Caller(2)
	if ok {
		data["caller"] = fmt.Sprintf("%s:%d", file, line)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(l.stdout, "Failed to marshal log metadata: %v\n", err)
		return
	}
	jsonData = append(jsonData, '\n')

	if level >= l.fileLevel {
		_, err := l.file.Write(jsonData)
		if err != nil {
			fmt.Fprintf(l.stdout, "Failed to write log to file: %v\n", err)
		}
	}
	if level >= l.consoleLevel {
		l.stdout.Write(jsonData)
	}
}

func (l *Logger) Trace(message string, metadata ...Metadata) {
	l.log(LogLevelTrace, message, metadata...)
}

func (l *Logger) Debug(message string, metadata ...Metadata) {
	l.log(LogLevelDebug, message, metadata...)
}

func (l *Logger) Info(message string, metadata ...Metadata) {
	l.log(LogLevelInfo, message, metadata...)
}

func (l *Logger) Warn(message string, metadata ...Metadata) {
	l.log(LogLevelWarn, message, metadata...)
}

func (l *Logger) Error(message string, metadata ...Metadata) {
	l.log(LogLevelError, message, metadata...)
}

func (l *Logger) RecoverPanic() {
	if r := recover(); r != nil {
		l.Error(fmt.Sprintf("Panic recovered: %v", r))
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
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		file:         l.file,
		attrs:        maps.Clone(l.attrs),
		groups:       append([]string{}, l.groups...),
		fileLevel:    l.fileLevel,
		consoleLevel: l.consoleLevel,
		mu:           l.mu,
		stdout:       l.stdout,
		closeOnce:    l.closeOnce,
	}

	namespace := newLogger.getNamespace()
	for key, value := range attrs {
		newLogger.attrs[namespace+key] = value
	}
	return newLogger
}

func (l *Logger) WithGroup(group string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		file:         l.file,
		attrs:        maps.Clone(l.attrs),
		groups:       append(append([]string{}, l.groups...), group),
		fileLevel:    l.fileLevel,
		consoleLevel: l.consoleLevel,
		mu:           l.mu,
		stdout:       l.stdout,
		closeOnce:    l.closeOnce,
	}
	return newLogger
}

func (l *Logger) SetFileLevel(level string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	levelInt, err := getLevelInt(level)
	if err != nil {
		return err
	}
	l.fileLevel = levelInt
	return nil
}
func (l *Logger) SetConsoleLevel(level string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	levelInt, err := getLevelInt(level)
	if err != nil {
		return err
	}
	l.consoleLevel = levelInt
	return nil
}

func (l *Logger) Name() string {
	return l.file.Name()
}

func (l *Logger) Close() {
	l.closeOnce.Do(func() {
		l.file.Close()
	})
}
