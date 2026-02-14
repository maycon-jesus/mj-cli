package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewMultiHandler(t *testing.T) {
	mh := NewMultiHandler()

	if mh.handlers == nil {
		t.Error("handlers map should not be nil")
	}
	if len(mh.handlers) != 0 {
		t.Error("handlers map should be empty initially")
	}
}

func TestMultiHandler_AddHandler(t *testing.T) {
	mh := NewMultiHandler()
	fh := newTestFileHandler(t)

	mh.AddHandler("file", fh)

	if len(mh.handlers) != 1 {
		t.Errorf("expected 1 handler, got %d", len(mh.handlers))
	}
	if mh.handlers["file"] != fh {
		t.Error("handler not stored correctly")
	}
}

func TestMultiHandler_AddHandlerOverwrite(t *testing.T) {
	mh := NewMultiHandler()
	fh1 := newTestFileHandler(t)
	fh2 := newTestFileHandler(t)

	mh.AddHandler("file", fh1)
	mh.AddHandler("file", fh2)

	if len(mh.handlers) != 1 {
		t.Errorf("expected 1 handler, got %d", len(mh.handlers))
	}
	if mh.handlers["file"] != fh2 {
		t.Error("handler should be overwritten")
	}
}

func TestMultiHandler_Enabled(t *testing.T) {
	tests := []struct {
		name   string
		levels []slog.Level
		check  slog.Level
		want   bool
	}{
		{
			name:   "enabled when one handler accepts",
			levels: []slog.Level{slog.LevelInfo, slog.LevelWarn},
			check:  slog.LevelInfo,
			want:   true,
		},
		{
			name:   "disabled when no handler accepts",
			levels: []slog.Level{slog.LevelWarn, slog.LevelError},
			check:  slog.LevelInfo,
			want:   false,
		},
		{
			name:   "enabled when all handlers accept",
			levels: []slog.Level{slog.LevelDebug, slog.LevelDebug},
			check:  slog.LevelInfo,
			want:   true,
		},
		{
			name:   "disabled with no handlers",
			levels: nil,
			check:  slog.LevelInfo,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := NewMultiHandler()
			for i, level := range tt.levels {
				fh := newTestFileHandler(t)
				fh.SetLevel(level)
				mh.AddHandler(string(rune('a'+i)), fh)
			}

			got := mh.Enabled(context.Background(), tt.check)
			if got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiHandler_Handle(t *testing.T) {
	mh := NewMultiHandler()

	fh1 := newTestFileHandler(t)
	fh1.SetLevel(slog.LevelDebug)
	mh.AddHandler("h1", fh1)

	fh2 := newTestFileHandler(t)
	fh2.SetLevel(slog.LevelDebug)
	mh.AddHandler("h2", fh2)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "multi test", 0)
	record.AddAttrs(slog.String("from", "test"))

	err := mh.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() error: %v", err)
	}

	// Both handlers should have received the record
	data1 := readJSONLine(t, fh1)
	data2 := readJSONLine(t, fh2)

	if data1["message"] != "multi test" {
		t.Errorf("handler 1: message = %v, want 'multi test'", data1["message"])
	}
	if data2["message"] != "multi test" {
		t.Errorf("handler 2: message = %v, want 'multi test'", data2["message"])
	}
}

func TestMultiHandler_HandleRespectsLevel(t *testing.T) {
	mh := NewMultiHandler()

	fhDebug := newTestFileHandler(t)
	fhDebug.SetLevel(slog.LevelDebug)
	mh.AddHandler("debug", fhDebug)

	fhWarn := newTestFileHandler(t)
	fhWarn.SetLevel(slog.LevelWarn)
	mh.AddHandler("warn", fhWarn)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "info message", 0)
	err := mh.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() error: %v", err)
	}

	// Debug handler should have received it
	content1, _ := os.ReadFile(fhDebug.Name())
	if !strings.Contains(string(content1), "info message") {
		t.Error("debug handler should have received the info message")
	}

	// Warn handler should NOT have received it
	content2, _ := os.ReadFile(fhWarn.Name())
	if strings.Contains(string(content2), "info message") {
		t.Error("warn handler should not have received the info message")
	}
}

func TestMultiHandler_HandleNoHandlers(t *testing.T) {
	mh := NewMultiHandler()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "orphan", 0)

	err := mh.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() with no handlers should not error, got: %v", err)
	}
}

func TestMultiHandler_WithAttrs(t *testing.T) {
	mh := NewMultiHandler()

	fh := newTestFileHandler(t)
	fh.SetLevel(slog.LevelDebug)
	mh.AddHandler("file", fh)

	newMh := mh.WithAttrs([]slog.Attr{
		slog.String("env", "test"),
	}).(*MultiHandler)

	// Write through the new handler
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "with attrs", 0)
	err := newMh.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() error: %v", err)
	}

	data := readJSONLine(t, fh)
	if data["env"] != "test" {
		t.Errorf("env = %v, want 'test'", data["env"])
	}
}

func TestMultiHandler_WithAttrsDoesNotModifyOriginal(t *testing.T) {
	mh := NewMultiHandler()
	fh := newTestFileHandler(t)
	mh.AddHandler("file", fh)

	newMh := mh.WithAttrs([]slog.Attr{
		slog.String("extra", "value"),
	}).(*MultiHandler)

	// The original and new should have different handler instances
	if mh.handlers["file"] == newMh.handlers["file"] {
		t.Error("WithAttrs should create new handler instances")
	}
}

func TestMultiHandler_WithGroup(t *testing.T) {
	mh := NewMultiHandler()

	fh := newTestFileHandler(t)
	fh.SetLevel(slog.LevelDebug)
	mh.AddHandler("file", fh)

	grouped := mh.WithGroup("request").(*MultiHandler)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "grouped", 0)
	record.AddAttrs(slog.String("method", "POST"))

	err := grouped.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle() error: %v", err)
	}

	data := readJSONLine(t, fh)
	if data["request.method"] != "POST" {
		t.Errorf("request.method = %v, want 'POST'", data["request.method"])
	}
}

func TestMultiHandler_WithGroupDoesNotModifyOriginal(t *testing.T) {
	mh := NewMultiHandler()
	fh := newTestFileHandler(t)
	mh.AddHandler("file", fh)

	grouped := mh.WithGroup("grp").(*MultiHandler)

	if mh.handlers["file"] == grouped.handlers["file"] {
		t.Error("WithGroup should create new handler instances")
	}
}

func TestMultiHandler_IntegrationWithSlog(t *testing.T) {
	mh := NewMultiHandler()

	fh := newTestFileHandler(t)
	fh.SetLevel(slog.LevelDebug)
	mh.AddHandler("file", fh)

	log := slog.New(mh)
	log.Info("slog integration", "key", "value")

	data := readJSONLine(t, fh)
	if data["message"] != "slog integration" {
		t.Errorf("message = %v, want 'slog integration'", data["message"])
	}
	if data["key"] != "value" {
		t.Errorf("key = %v, want 'value'", data["key"])
	}
}
