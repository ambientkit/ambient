package hclogadapter

import (
	"io"
	"log"

	"github.com/ambientkit/ambient"
	"github.com/hashicorp/go-hclog"
)

// Logger describes the interface that must be implemeted by all loggers.
type Logger struct {
	name string
	log  ambient.AppLogger
}

// New returns a new Ambient logger for plugins.
func New(name string, logger ambient.AppLogger) *Logger {
	return &Logger{
		name: name,
		log:  logger,
	}
}

// Log sends a message to the logger.
// Args are alternating key, val pairs.
// keys must be strings.
// vals can be any type, but display is implementation specific.
// Emit a message and key/value pairs at a provided log level.
func (l *Logger) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.Debug, hclog.Trace, hclog.NoLevel:
		l.log.Debug("%v: %v | %v", l.name, msg, args)
	case hclog.Info:
		l.log.Info("%v: %v | %v", l.name, msg, args)
	case hclog.Warn:
		l.log.Warn("%v: %v | %v", l.name, msg, args)
	case hclog.Error, hclog.Off:
		l.log.Error("%v: %v | %v", l.name, msg, args)
	default:
		l.log.Info("%v: %v | %v", l.name, msg, args)
	}
}

// Trace emits a message and key/value pairs at the TRACE level.
func (l *Logger) Trace(msg string, args ...interface{}) {
	l.Log(hclog.Trace, msg, args)
}

// Debug emits a message and key/value pairs at the DEBUG level.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Log(hclog.Debug, msg, args)
}

// Info emits a message and key/value pairs at the INFO level.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Log(hclog.Info, msg, args)
}

// Warn emits a message and key/value pairs at the WARN level.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Log(hclog.Warn, msg, args)
}

// Error emits a message and key/value pairs at the ERROR level.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Log(hclog.Error, msg, args)
}

// IsTrace indicates if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func (l *Logger) IsTrace() bool {
	return true
}

// IsDebug indicates if DEBUG logs would be emitted. This and the other Is* guards.
func (l *Logger) IsDebug() bool {
	return true
}

// IsInfo indicates if INFO logs would be emitted. This and the other Is* guards.
func (l *Logger) IsInfo() bool {
	return true
}

// IsWarn indicates if WARN logs would be emitted. This and the other Is* guards.
func (l *Logger) IsWarn() bool {
	return true
}

// IsError indicates if ERROR logs would be emitted. This and the other Is* guards.
func (l *Logger) IsError() bool {
	return true
}

// ImpliedArgs returns With key/value pairs.
func (l *Logger) ImpliedArgs() []interface{} {
	return nil
}

// With creates a sublogger that will always have the given key/value pairs.
func (l *Logger) With(args ...interface{}) hclog.Logger {
	return l
}

// Name returns the Name of the logger.
func (l *Logger) Name() string {
	return l.name
}

// Named creates a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (l *Logger) Named(name string) hclog.Logger {
	return &Logger{
		name: l.name + "." + name,
		log:  l.log,
	}
}

// ResetNamed creates a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (l *Logger) ResetNamed(name string) hclog.Logger {
	return &Logger{
		name: name,
		log:  l.log,
	}
}

// SetLevel updates the level. This should affect all related loggers as well,
// unless they were created with IndependentLevels. If an
// implementation cannot update the level on the fly, it should no-op.
func (l *Logger) SetLevel(level hclog.Level) {

}

// StandardLogger returns a value that conforms to the stdlib log.Logger interface.
func (l *Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return nil
}

// StandardWriter returns a value that conforms to io.Writer, which can be passed into log.SetOutput().
func (l *Logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	//return log.Logger
	return nil
}
