package mock

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/ambientkit/ambient"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// LoggerPlugin represents an Ambient plugin.
type LoggerPlugin struct {
	log    *Logger
	output io.Writer
}

// NewLoggerPlugin returns an Ambient plugin that provides logging using the standard logger.
func NewLoggerPlugin(optionalWriter io.Writer) *LoggerPlugin {
	return &LoggerPlugin{
		output: optionalWriter,
	}
}

// PluginName returns the plugin name.
func (p *LoggerPlugin) PluginName() string {
	return "mocklogger"
}

// PluginVersion returns the plugin version.
func (p *LoggerPlugin) PluginVersion() string {
	return "1.0.0"
}

// Logger returns a logger.
func (p *LoggerPlugin) Logger(appName string, appVersion string, optionalWriter io.Writer) (ambient.AppLogger, error) {
	// Create the logger.
	p.log = p.NewLogger(appName, appVersion, optionalWriter)

	return p.log, nil
}

// Logger represents a logger.
type Logger struct {
	log *log.Logger

	appName    string
	appVersion string

	tracerProvider *sdktrace.TracerProvider
}

// NewLogger returns a new logger with a default log level of error.
func (p *LoggerPlugin) NewLogger(appName string, appVersion string, optionalWriter io.Writer) *Logger {
	l := log.Default()
	if optionalWriter != nil {
		l.SetOutput(optionalWriter)
	} else if p.output != nil {
		l.SetOutput(p.output)
	}

	return &Logger{
		log: l,

		appName:    appName,
		appVersion: appVersion,
	}
}

// SetLogLevel will set the logger output level.
func (l *Logger) SetLogLevel(level ambient.LogLevel) {}

func (l *Logger) output(format string, v ...interface{}) {
	if len(format) == 0 {
		l.log.Println(v...)
	} else {
		l.log.Printf(format+"\n", v...)
	}
}

// Log is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Log(level ambient.LogLevel, format string, v ...interface{}) {
	switch level {
	case ambient.LogLevelDebug:
		l.Debug(format, v...)
	case ambient.LogLevelInfo:
		l.Info(format, v...)
	case ambient.LogLevelWarn:
		l.Warn(format, v...)
	case ambient.LogLevelError:
		l.Error(format, v...)
	case ambient.LogLevelFatal:
		l.Fatal(format, v...)
	default:
		l.Info(format, v...)
	}
}

// Debug is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.output(format, v...)
}

// Info is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Info(format string, v ...interface{}) {
	l.output(format, v...)
}

// Warn is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.output(format, v...)
}

// Error is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Error(format string, v ...interface{}) {
	l.output(format, v...)
}

// Fatal is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty. It's followed by a call
// to os.Exit(1).
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.output(format, v...)
}

// Name returns the name of the logger.
func (l *Logger) Name() string {
	return l.appName
}

// Named returns a new logger with the appended name, linked to the existing
// logger.
func (l *Logger) Named(name string) ambient.AppLogger {
	return &Logger{
		appName:    fmt.Sprintf("%v.%v", l.appName, name),
		log:        l.log,
		appVersion: l.appVersion,
	}
}

// For returns the context-aware logger.
func (l *Logger) For(ctx context.Context) ambient.Logger {
	// FIXME: Do something with context? May not need it in a mock.
	return l
}

// SetTracerProvider sets the tracer provider.
func (l *Logger) SetTracerProvider(tp *sdktrace.TracerProvider) {
	l.tracerProvider = tp
}

// Trace returns a context and an OpenTelemetry span.
func (l *Logger) Trace(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return l.tracerProvider.Tracer(l.appName).Start(ctx, spanName)
}
