package logger

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
	"os"
	"sync"
)

type FileHandler struct {
	file   *os.File
	mu     *sync.RWMutex
	attrs  map[string]any
	prefix string
	level  slog.Level
}

func NewFileHandler(file *os.File) *FileHandler {
	return &FileHandler{
		file:   file,
		mu:     &sync.RWMutex{},
		attrs:  make(map[string]any),
		prefix: "",
	}
}

func NewTemporaryFileHandler(appName string) (*FileHandler, error) {
	fileName := appName + "-log-*"
	logFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}
	return NewFileHandler(logFile), nil
}

func (h *FileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return level >= h.level
}

func (h *FileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()
	nHandler := h.cloneUnsafe()

	for _, attr := range attrs {
		nHandler.attrs[h.prefix+attr.Key] = attr.Value.Any()
	}
	return nHandler
}

func (h *FileHandler) WithGroup(name string) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()
	nHandler := h.cloneUnsafe()
	nHandler.prefix = h.prefix + name + "."
	return nHandler
}

func (h *FileHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data := make(map[string]any)
	maps.Copy(data, h.attrs)

	record.Attrs(func(attr slog.Attr) bool {
		data[h.prefix+attr.Key] = attr.Value.Any()
		return true
	})

	data["level"] = record.Level.String()
	data["time"] = record.Time.UTC().Format("2006-01-02T15:04:05Z07:00")
	data["message"] = record.Message

	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = h.file.Write(append(dataJson, '\n'))
	return err
}

func (h *FileHandler) Close() error {
	return h.file.Close()
}

func (h *FileHandler) Name() string {
	return h.file.Name()
}

func (h *FileHandler) SetLevel(level slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.level = level
}

func (h *FileHandler) cloneUnsafe() *FileHandler {
	return &FileHandler{
		file:   h.file,
		mu:     h.mu,
		attrs:  maps.Clone(h.attrs),
		prefix: h.prefix,
		level:  h.level,
	}
}

var _ slog.Handler = (*FileHandler)(nil)
