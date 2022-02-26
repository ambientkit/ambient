package ambient

import (
	"io"
	"log"
)

// MockLoggerPlugin represents an Ambient plugin.
type MockLoggerPlugin struct {
	log    *MockLogger
	output io.Writer
}

// NewMockLoggerPlugin returns an Ambient plugin that provides logging using the standard logger.
func NewMockLoggerPlugin(optionalWriter io.Writer) *MockLoggerPlugin {
	return &MockLoggerPlugin{
		output: optionalWriter,
	}
}

// PluginName returns the plugin name.
func (p *MockLoggerPlugin) PluginName() string {
	return "mocklogger"
}

// PluginVersion returns the plugin version.
func (p *MockLoggerPlugin) PluginVersion() string {
	return "1.0.0"
}

// Logger returns a logger.
func (p *MockLoggerPlugin) Logger(appName string, appVersion string, optionalWriter io.Writer) (AppLogger, error) {
	// Create the logger.
	p.log = p.NewMockLogger(appName, appVersion, optionalWriter)

	return p.log, nil
}

// MockLogger represents a logger.
type MockLogger struct {
	log *log.Logger

	appName    string
	appVersion string
}

// NewMockLogger returns a new logger with a default log level of error.
func (p *MockLoggerPlugin) NewMockLogger(appName string, appVersion string, optionalWriter io.Writer) *MockLogger {
	l := log.Default()
	if optionalWriter != nil {
		l.SetOutput(optionalWriter)
	} else if p.output != nil {
		l.SetOutput(p.output)
	}

	return &MockLogger{
		log: l,

		appName:    appName,
		appVersion: appVersion,
	}
}

// SetLogLevel will set the logger output level.
func (l *MockLogger) SetLogLevel(level LogLevel) {}

func (l *MockLogger) output(format string, v ...interface{}) {
	if len(format) == 0 {
		l.log.Println(v...)
	} else {
		l.log.Printf(format+"\n", v...)
	}
}

// Debug is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *MockLogger) Debug(format string, v ...interface{}) {
	l.output(format, v...)
}

// Info is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *MockLogger) Info(format string, v ...interface{}) {
	l.output(format, v...)
}

// Warn is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *MockLogger) Warn(format string, v ...interface{}) {
	l.output(format, v...)
}

// Error is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *MockLogger) Error(format string, v ...interface{}) {
	l.output(format, v...)
}

// Fatal is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty. It's followed by a call
// to os.Exit(1).
func (l *MockLogger) Fatal(format string, v ...interface{}) {
	l.output(format, v...)
}
