package logger

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewFileHandler(t *testing.T) {
	file, err := os.CreateTemp("", "test-file-handler-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	handler := NewFileHandler(file)

	if handler.file != file {
		t.Error("file not set correctly")
	}
	if handler.mu == nil {
		t.Error("mutex should not be nil")
	}
	if handler.attrs == nil {
		t.Error("attrs should not be nil")
	}
	if len(handler.attrs) != 0 {
		t.Error("attrs should be empty initially")
	}
	if handler.prefix != "" {
		t.Error("prefix should be empty initially")
	}
}

func TestNewTemporaryFileHandler(t *testing.T) {
	handler, err := NewTemporaryFileHandler("test-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(handler.Name())
	defer handler.Close()

	if !strings.Contains(handler.Name(), "test-app-log-") {
		t.Errorf("file name %q should contain 'test-app-log-'", handler.Name())
	}
}

func TestFileHandler_Enabled(t *testing.T) {
	tests := []struct {
		name         string
		handlerLevel slog.Level
		recordLevel  slog.Level
		want         bool
	}{
		{
			name:         "info enabled at info level",
			handlerLevel: slog.LevelInfo,
			recordLevel:  slog.LevelInfo,
			want:         true,
		},
		{
			name:         "warn enabled at info level",
			handlerLevel: slog.LevelInfo,
			recordLevel:  slog.LevelWarn,
			want:         true,
		},
		{
			name:         "debug disabled at info level",
			handlerLevel: slog.LevelInfo,
			recordLevel:  slog.LevelDebug,
			want:         false,
		},
		{
			name:         "debug enabled at debug level",
			handlerLevel: slog.LevelDebug,
			recordLevel:  slog.LevelDebug,
			want:         true,
		},
		{
			name:         "error enabled at warn level",
			handlerLevel: slog.LevelWarn,
			recordLevel:  slog.LevelError,
			want:         true,
		},
		{
			name:         "info disabled at warn level",
			handlerLevel: slog.LevelWarn,
			recordLevel:  slog.LevelInfo,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newTestFileHandler(t)
			handler.SetLevel(tt.handlerLevel)

			got := handler.Enabled(context.Background(), tt.recordLevel)
			if got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileHandler_Handle(t *testing.T) {
	handler := newTestFileHandler(t)
	handler.SetLevel(slog.LevelDebug)

	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	record := slog.NewRecord(now, slog.LevelInfo, "test message", 0)
	record.AddAttrs(slog.String("key1", "value1"), slog.Int("key2", 42))

	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() error: %v", err)
	}

	data := readJSONLine(t, handler)

	if data["message"] != "test message" {
		t.Errorf("message = %v, want 'test message'", data["message"])
	}
	if data["level"] != "INFO" {
		t.Errorf("level = %v, want 'INFO'", data["level"])
	}
	if data["key1"] != "value1" {
		t.Errorf("key1 = %v, want 'value1'", data["key1"])
	}
	if data["key2"] != float64(42) {
		t.Errorf("key2 = %v, want 42", data["key2"])
	}
	if data["time"] != "2025-01-15T10:30:00Z" {
		t.Errorf("time = %v, want '2025-01-15T10:30:00Z'", data["time"])
	}
}

func TestFileHandler_HandleWithPresetAttrs(t *testing.T) {
	handler := newTestFileHandler(t)
	handler.SetLevel(slog.LevelDebug)

	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("app", "myapp"),
	}).(*FileHandler)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0)
	err := handlerWithAttrs.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() error: %v", err)
	}

	data := readJSONLine(t, handlerWithAttrs)

	if data["app"] != "myapp" {
		t.Errorf("app = %v, want 'myapp'", data["app"])
	}
	if data["message"] != "hello" {
		t.Errorf("message = %v, want 'hello'", data["message"])
	}
}

func TestFileHandler_WithAttrs(t *testing.T) {
	handler := newTestFileHandler(t)

	newHandler := handler.WithAttrs([]slog.Attr{
		slog.String("k1", "v1"),
		slog.Int("k2", 2),
	})

	fh := newHandler.(*FileHandler)

	if fh.attrs["k1"] != "v1" {
		t.Errorf("k1 = %v, want 'v1'", fh.attrs["k1"])
	}
	if fh.attrs["k2"] != int64(2) {
		t.Errorf("k2 = %v, want 2", fh.attrs["k2"])
	}

	// original handler should not be modified
	if len(handler.attrs) != 0 {
		t.Error("original handler attrs should remain empty")
	}
}

func TestFileHandler_WithAttrsUsesPrefix(t *testing.T) {
	handler := newTestFileHandler(t)

	grouped := handler.WithGroup("req").(*FileHandler)
	withAttr := grouped.WithAttrs([]slog.Attr{
		slog.String("method", "GET"),
	}).(*FileHandler)

	if withAttr.attrs["req.method"] != "GET" {
		t.Errorf("attrs = %v, want key 'req.method' with value 'GET'", withAttr.attrs)
	}
}

func TestFileHandler_WithGroup(t *testing.T) {
	handler := newTestFileHandler(t)

	grouped := handler.WithGroup("request").(*FileHandler)

	if grouped.prefix != "request." {
		t.Errorf("prefix = %q, want 'request.'", grouped.prefix)
	}

	// original handler should not be modified
	if handler.prefix != "" {
		t.Error("original handler prefix should remain empty")
	}
}

func TestFileHandler_WithGroupNested(t *testing.T) {
	handler := newTestFileHandler(t)

	nested := handler.WithGroup("http").(*FileHandler).WithGroup("request").(*FileHandler)

	if nested.prefix != "http.request." {
		t.Errorf("prefix = %q, want 'http.request.'", nested.prefix)
	}
}

func TestFileHandler_HandleWithGroup(t *testing.T) {
	handler := newTestFileHandler(t)
	handler.SetLevel(slog.LevelDebug)

	grouped := handler.WithGroup("ctx").(*FileHandler)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "grouped msg", 0)
	record.AddAttrs(slog.String("user", "alice"))

	err := grouped.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() error: %v", err)
	}

	data := readJSONLine(t, grouped)

	if data["ctx.user"] != "alice" {
		t.Errorf("ctx.user = %v, want 'alice'", data["ctx.user"])
	}
}

func TestFileHandler_SetLevel(t *testing.T) {
	handler := newTestFileHandler(t)

	handler.SetLevel(slog.LevelWarn)

	if !handler.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("should be enabled for Warn after SetLevel(Warn)")
	}
	if handler.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("should not be enabled for Info after SetLevel(Warn)")
	}
}

func TestFileHandler_Close(t *testing.T) {
	file, err := os.CreateTemp("", "test-close-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	handler := NewFileHandler(file)
	err = handler.Close()
	if err != nil {
		t.Fatalf("Close() error: %v", err)
	}

	// Writing after close should fail
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "after close", 0)
	err = handler.Handle(context.Background(), record)
	if err == nil {
		t.Error("Handle() should return error after Close()")
	}
}

func TestFileHandler_Name(t *testing.T) {
	handler, err := NewTemporaryFileHandler("myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(handler.Name())
	defer handler.Close()

	name := handler.Name()
	if !strings.Contains(name, "myapp-log-") {
		t.Errorf("Name() = %q, should contain 'myapp-log-'", name)
	}
}

func TestFileHandler_ConcurrentAccess(t *testing.T) {
	handler := newTestFileHandler(t)
	handler.SetLevel(slog.LevelDebug)

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			record := slog.NewRecord(time.Now(), slog.LevelInfo, "concurrent", 0)
			record.AddAttrs(slog.Int("i", i))
			_ = handler.Handle(context.Background(), record)
		}(i)
	}
	wg.Wait()

	content, err := os.ReadFile(handler.Name())
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 50 {
		t.Errorf("expected 50 log lines, got %d", len(lines))
	}

	for i, line := range lines {
		var data map[string]any
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			t.Errorf("line %d: invalid JSON: %v", i, err)
		}
	}
}

func TestFileHandler_ImplementsSlogHandler(t *testing.T) {
	var _ slog.Handler = (*FileHandler)(nil)
}

// helpers

func newTestFileHandler(t *testing.T) *FileHandler {
	t.Helper()
	file, err := os.CreateTemp("", "test-handler-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})
	return NewFileHandler(file)
}

func readJSONLine(t *testing.T, handler *FileHandler) map[string]any {
	t.Helper()
	content, err := os.ReadFile(handler.file.Name())
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	lastLine := lines[len(lines)-1]

	var data map[string]any
	if err := json.Unmarshal([]byte(lastLine), &data); err != nil {
		t.Fatalf("failed to parse JSON: %v\nline: %s", err, lastLine)
	}
	return data
}
