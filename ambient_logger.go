package ambient

import (
	"context"
	"os"
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// AppLogger represents the log service for the app.
type AppLogger interface {
	Logger

	// Fatal is equivalent to log.Printf() + "\n" if format is not empty.
	// It's equivalent to Println() if format is empty. It's followed by a call
	// to os.Exit(1).
	// Fatal is reserved for the app level only.
	Fatal(format string, v ...interface{})
	// SetLogLevel sets the logger output level.
	SetLogLevel(level LogLevel)
	// Named returns a new logger with the appended name, linked to the existing
	// logger.
	Named(name string) AppLogger
	// Name returns the name of the logger.
	Name() string
	// SetTracerProvider sets the OpenTelemetry tracer provider.
	SetTracerProvider(tp *sdktrace.TracerProvider)
}

// Logger represents the log service for the plugins.
type Logger interface {
	// Log is equivalent to log.Printf() + "\n" if format is not empty.
	// It's equivalent to Println() if format is empty.
	Log(level LogLevel, format string, v ...interface{})
	// Debug is equivalent to log.Printf() + "\n" if format is not empty.
	// It's equivalent to Println() if format is empty.
	Debug(format string, v ...interface{})
	// Info is equivalent to log.Printf() + "\n" if format is not empty.
	// It's equivalent to Println() if format is empty.
	Info(format string, v ...interface{})
	// Warn is equivalent to log.Printf() + "\n" if format is not empty.
	// It's equivalent to Println() if format is empty.
	Warn(format string, v ...interface{})
	// Error is equivalent to log.Printf() + "\n" if format is not empty.
	// It's equivalent to Println() if format is empty.
	Error(format string, v ...interface{})
	// For returns a context-aware logger to support OpenTracing.
	For(ctx context.Context) Logger
	// Trace returns a context and span to support OpenTracing.
	Trace(ctx context.Context, operation string) (context.Context, trace.Span)
}

// LogLevel is a log level.
type LogLevel int

const (
	// LogLevelDebug is for debugging output. It's very verbose.
	LogLevelDebug LogLevel = iota
	// LogLevelInfo is for informational messages. It shows messages on services
	// starting, stopping, and users logging in.
	LogLevelInfo
	// LogLevelWarn is for behavior that may need to be fixed. It shows
	// permission warnings for plugins.
	LogLevelWarn
	// LogLevelError is for messages when something is wrong with the
	// app and it needs to be corrected.
	LogLevelError
	// LogLevelFatal is for messages when the app cannot continue and
	// will halt.
	LogLevelFatal
)

// EnvLogLevel returns the log level from the AMB_LOGLEVEL environment variable.
func EnvLogLevel() LogLevel {
	ll := os.Getenv("AMB_LOGLEVEL")
	switch true {
	case strings.EqualFold(ll, "FATAL") || strings.EqualFold(ll, "4"):
		return LogLevelFatal
	case strings.EqualFold(ll, "ERROR") || strings.EqualFold(ll, "3"):
		return LogLevelError
	case strings.EqualFold(ll, "WARN") || strings.EqualFold(ll, "2"):
		return LogLevelWarn
	case strings.EqualFold(ll, "INFO") || strings.EqualFold(ll, "1"):
		return LogLevelInfo
	case strings.EqualFold(ll, "DEBUG") || strings.EqualFold(ll, "0"):
		return LogLevelDebug
	}
	return LogLevelInfo
}
