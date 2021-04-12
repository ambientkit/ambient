// Package pluginlogger provides a way to output logging messages.
// Reference: https://dzone.com/articles/logging-levels-what-they-are-and-how-they-help-you
package pluginlogger

import "github.com/josephspurrier/ambient"

// Logger represents a plugin logger.
type Logger struct {
	log ambient.IAppLogger
}

// NewPluginLogger returns a new logger with a default log level of error.
func NewPluginLogger(logger ambient.IAppLogger) *Logger {
	return &Logger{
		log: logger,
	}
}

// Debug is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log.Debug(format, v...)
}

// Info is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Info(format string, v ...interface{}) {
	l.log.Info(format, v...)
}

// Warn is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log.Warn(format, v...)
}

// Error is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Error(format string, v ...interface{}) {
	l.log.Error(format, v...)
}
