package logger

import (
	"context"
	"log/slog"
	"maps"
)

type MultiHandler struct {
	handlers map[string]slog.Handler
}

func NewMultiHandler() *MultiHandler {
	return &MultiHandler{
		handlers: make(map[string]slog.Handler),
	}
}

func (m *MultiHandler) AddHandler(name string, handler slog.Handler) {
	m.handlers[name] = handler
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range m.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nHandler := m.cloneUnsafe()
	for name, handler := range nHandler.handlers {
		nHandler.handlers[name] = handler.WithAttrs(attrs)
	}
	return nHandler
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	nHandler := m.cloneUnsafe()
	for handlerName, handler := range nHandler.handlers {
		nHandler.handlers[handlerName] = handler.WithGroup(name)
	}
	return nHandler
}

func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range m.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) cloneUnsafe() *MultiHandler {
	return &MultiHandler{
		handlers: maps.Clone(m.handlers),
	}
}
