// Package logger provides a way to output logging messages.
// Reference: https://dzone.com/articles/logging-levels-what-they-are-and-how-they-help-you
package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger represents a logger.
type Logger struct {
	*logrus.Logger

	appName    string
	appVersion string
}

// NewLogger returns a new logger with a default log level of error.
func NewLogger(appName string, appVersion string) *Logger {
	var base = logrus.New()
	//base.SetFormatter(&logrus.JSONFormatter{})
	base.Level = logrus.InfoLevel

	return &Logger{
		Logger: base,

		appName:    appName,
		appVersion: appVersion,
	}
}

func (l *Logger) logentry() *logrus.Entry {
	standardFields := logrus.Fields{
		"app":     l.appName,
		"version": l.appVersion,
	}

	return l.WithFields(standardFields)
}

// SetLevel will set the logger output level.
func (l *Logger) SetLevel(level logrus.Level) {
	// Set log level temporarily to info.
	l.Level = logrus.InfoLevel
	l.logentry().Println("log level set to:", level)
	l.Level = level
}

// Debug is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Debug(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Debugln(v...)
	} else {
		l.logentry().Debugf(format, v...)
	}
}

// Info is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Info(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Infoln(v...)
	} else {
		l.logentry().Infof(format, v...)
	}
}

// Error is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Error(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Errorln(v...)
	} else {
		l.logentry().Errorf(format, v...)
	}
}

// Fatal is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty. It's followed by a call
// to os.Exit(1).
func (l *Logger) Fatal(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Fatalln(v...)
	} else {
		l.logentry().Fatalf(format, v...)
	}
}
