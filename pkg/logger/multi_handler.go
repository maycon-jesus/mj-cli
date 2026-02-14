package logger

import (
	"context"
	"log/slog"
	"maps"
	"sync"
)

// MultiHandler is a [slog.Handler] that fans out each log record to multiple
// named handlers. A record is only forwarded to a handler if that handler
// reports [slog.Handler.Enabled] for the record's level.
type MultiHandler struct {
	handlers map[string]slog.Handler
	mu       *sync.RWMutex
}

// NewMultiHandler returns an empty [MultiHandler]. Use [MultiHandler.AddHandler]
// to register handlers before passing it to [slog.New].
func NewMultiHandler() *MultiHandler {
	return &MultiHandler{
		handlers: make(map[string]slog.Handler),
		mu:       &sync.RWMutex{},
	}
}

// AddHandler registers a handler under the given name. If a handler with the
// same name already exists, it is replaced.
func (m *MultiHandler) AddHandler(name string, handler slog.Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[name] = handler
}

// Enabled reports true if at least one registered handler is enabled for level.
func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, handler := range m.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// WithAttrs returns a new [MultiHandler] where every registered handler has
// the given attributes applied.
func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nHandler := m.cloneUnsafe()
	for name, handler := range nHandler.handlers {
		nHandler.handlers[name] = handler.WithAttrs(attrs)
	}
	return nHandler
}

// WithGroup returns a new [MultiHandler] where every registered handler has
// the given group applied.
func (m *MultiHandler) WithGroup(name string) slog.Handler {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nHandler := m.cloneUnsafe()
	for handlerName, handler := range nHandler.handlers {
		nHandler.handlers[handlerName] = handler.WithGroup(name)
	}
	return nHandler
}

// Handle forwards the record to every registered handler that is enabled for
// the record's level. If any handler returns an error, Handle short-circuits
// and returns that error without forwarding to the remaining handlers.
func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, handler := range m.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ slog.Handler = (*MultiHandler)(nil)

func (m *MultiHandler) cloneUnsafe() *MultiHandler {
	return &MultiHandler{
		handlers: maps.Clone(m.handlers),
		mu:       &sync.RWMutex{},
	}
}
