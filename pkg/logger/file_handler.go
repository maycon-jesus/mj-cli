package logger

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"slices"
	"sync"
)

// FileHandler is a thread-safe [slog.Handler] that writes JSON-encoded log
// records to an [os.File], one JSON object per line.
//
// Each record is serialized as a JSON object containing the handler's
// accumulated attributes, the record-level attributes, and the standard
// "level", "time", and "message" fields. Groups created via [FileHandler.WithGroup]
// are represented as nested JSON objects.
//
// Cloned handlers (via [FileHandler.WithAttrs] and [FileHandler.WithGroup])
// share the same file and mutex, so concurrent writes from different handler
// instances are safe.
type FileHandler struct {
	file         *os.File
	mu           *sync.RWMutex
	attrs        map[string]any
	currentGroup map[string]any
	prefix       []string
	level        slog.Level
}

// NewFileHandler returns a [FileHandler] that writes to file.
// The default level is [slog.LevelDebug].
func NewFileHandler(file *os.File) *FileHandler {
	attrs := make(map[string]any)
	return &FileHandler{
		file:         file,
		mu:           &sync.RWMutex{},
		attrs:        attrs,
		currentGroup: attrs,
		level:        slog.LevelDebug,
		prefix:       make([]string, 0),
	}
}

// NewTemporaryFileHandler creates a [FileHandler] backed by a temporary file
// in the OS temp directory. The file is named "<appName>-log-*".
func NewTemporaryFileHandler(appName string) (*FileHandler, error) {
	fileName := appName + "-log-*"
	logFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}
	return NewFileHandler(logFile), nil
}

// Enabled reports whether the handler is configured to log at the given level.
func (h *FileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return level >= h.level
}

// WithAttrs returns a new [FileHandler] whose output includes the given attributes.
// The returned handler shares the same file and mutex as the receiver.
func (h *FileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.RLock()
	defer h.mu.RUnlock()

	nHandler := h.cloneUnsafe()

	for _, attr := range attrs {
		nHandler.currentGroup[attr.Key] = attr.Value.Any()
	}
	return nHandler
}

// WithGroup returns a new [FileHandler] that nests subsequent attributes under
// a JSON object keyed by name. The returned handler shares the same file and mutex.
func (h *FileHandler) WithGroup(name string) slog.Handler {
	h.mu.RLock()
	defer h.mu.RUnlock()
	nHandler := h.cloneUnsafe()
	nHandler.prefix = append(nHandler.prefix, name)
	nHandler.currentGroup[name] = make(map[string]any)
	nHandler.currentGroup = nHandler.currentGroup[name].(map[string]any)
	return nHandler
}

// Handle serializes the record as a single JSON line and writes it to the file.
// The output merges the handler's stored attributes with the record's attributes,
// plus the standard "level", "time" (UTC, RFC 3339), and "message" fields.
func (h *FileHandler) Handle(ctx context.Context, record slog.Record) error {
	data := deepCloneMap(h.attrs)

	dataCurrentGroup := data
	for _, group := range h.prefix {
		if _, ok := dataCurrentGroup[group]; !ok {
			dataCurrentGroup[group] = make(map[string]any)
		}
		dataCurrentGroup = dataCurrentGroup[group].(map[string]any)
	}

	record.Attrs(func(attr slog.Attr) bool {
		dataCurrentGroup[attr.Key] = attr.Value.Any()
		return true
	})

	data["level"] = record.Level.String()
	data["time"] = record.Time.UTC().Format("2006-01-02T15:04:05Z07:00")
	data["message"] = record.Message

	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err = h.file.Write(append(dataJson, '\n'))
	return err
}

// Close closes the underlying file.
func (h *FileHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.file.Close()
}

// Name returns the file path of the underlying log file.
func (h *FileHandler) Name() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.file.Name()
}

// SetLevel changes the minimum level at which records are logged.
func (h *FileHandler) SetLevel(level slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.level = level
}

func (h *FileHandler) cloneUnsafe() *FileHandler {
	attrs := deepCloneMap(h.attrs)
	prefix := slices.Clone(h.prefix)

	currentGroup := attrs

	for _, group := range h.prefix {
		currentGroup = currentGroup[group].(map[string]any)
	}

	return &FileHandler{
		file:         h.file,
		mu:           h.mu,
		attrs:        attrs,
		currentGroup: currentGroup,
		prefix:       prefix,
		level:        h.level,
	}
}

var _ slog.Handler = (*FileHandler)(nil)

func deepCloneMap(m map[string]any) map[string]any {
	clone := make(map[string]any)
	for k, v := range m {
		if nestedMap, ok := v.(map[string]any); ok {
			clone[k] = deepCloneMap(nestedMap)
		} else {
			clone[k] = v
		}
	}
	return clone
}
