package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// Logger is the structured logging contract used across the codebase.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}

type slogLogger struct {
	l *slog.Logger
}

// New builds a JSON structured logger at the given level.
func New(level string) Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return &slogLogger{l: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))}
}

func (s *slogLogger) Debug(msg string, kv ...interface{}) { s.l.Debug(msg, kv...) }
func (s *slogLogger) Info(msg string, kv ...interface{})  { s.l.Info(msg, kv...) }
func (s *slogLogger) Warn(msg string, kv ...interface{})  { s.l.Warn(msg, kv...) }
func (s *slogLogger) Error(msg string, kv ...interface{}) { s.l.Error(msg, kv...) }

func (s *slogLogger) With(kv ...interface{}) Logger {
	return &slogLogger{l: s.l.With(kv...)}
}

// FromContext returns the logger stored in context, or a default one.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return l
	}
	return New("info")
}

type loggerKey struct{}

// WithContext attaches a logger to the context.
func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}
