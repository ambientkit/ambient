// Package logger provides a way to output logging messages.
// Reference: https://dzone.com/articles/logging-levels-what-they-are-and-how-they-help-you
package logger

// Logger represents a plugin logger.
type Logger struct {
	log ILogger
}

// ILogger represents an application logger.
type ILogger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// NewPluginLogger returns a new logger with a default log level of error.
func NewPluginLogger(logger ILogger) *Logger {
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
