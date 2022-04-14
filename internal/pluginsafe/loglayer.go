package pluginsafe

import (
	"context"

	"github.com/ambientkit/ambient"
	"go.opentelemetry.io/otel/trace"
)

// PluginLogger represents a plugin logger.
type PluginLogger struct {
	log ambient.AppLogger
}

// NewPluginLogger returns a new logger with a default log level of error.
func NewPluginLogger(logger ambient.AppLogger) *PluginLogger {
	return &PluginLogger{
		log: logger,
	}
}

// Log is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *PluginLogger) Log(level ambient.LogLevel, format string, v ...interface{}) {
	switch level {
	case ambient.LogLevelDebug:
		l.Debug(format, v...)
	case ambient.LogLevelInfo:
		l.Info(format, v...)
	case ambient.LogLevelWarn:
		l.Warn(format, v...)
	case ambient.LogLevelError:
		l.Error(format, v...)
	default:
		l.Info(format, v...)
	}
}

// Debug is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *PluginLogger) Debug(format string, v ...interface{}) {
	l.log.Debug(format, v...)
}

// Info is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *PluginLogger) Info(format string, v ...interface{}) {
	l.log.Info(format, v...)
}

// Warn is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *PluginLogger) Warn(format string, v ...interface{}) {
	l.log.Warn(format, v...)
}

// Error is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *PluginLogger) Error(format string, v ...interface{}) {
	l.log.Error(format, v...)
}

// For returns a context-aware logger.
func (l *PluginLogger) For(ctx context.Context) ambient.Logger {
	return l.log.For(ctx)
}

// Trace returns a context and an OpenTelemetry span.
func (l *PluginLogger) Trace(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return l.log.Trace(ctx, spanName)
}
