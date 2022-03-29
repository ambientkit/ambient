package mock

import (
	"fmt"
	"io"
	"log"

	"github.com/ambientkit/ambient"
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
