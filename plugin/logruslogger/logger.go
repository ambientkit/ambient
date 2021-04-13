package logruslogger

import (
	"github.com/josephspurrier/ambient"
	"github.com/sirupsen/logrus"
)

// Logger represents a logger.
type Logger struct {
	log *logrus.Logger

	appName    string
	appVersion string
}

// NewLogger returns a new logger with a default log level of error.
func NewLogger(appName string, appVersion string) *Logger {
	var base = logrus.New()
	//base.SetFormatter(&logrus.JSONFormatter{})
	base.Level = logrus.InfoLevel

	return &Logger{
		log: base,

		appName:    appName,
		appVersion: appVersion,
	}
}

func (l *Logger) logentry() *logrus.Entry {
	standardFields := logrus.Fields{
		"app":     l.appName,
		"version": l.appVersion,
	}

	return l.log.WithFields(standardFields)
}

// SetLogLevel will set the logger output level.
func (l *Logger) SetLogLevel(level ambient.LogLevel) {
	// Set log level temporarily to info.
	l.log.Level = logrus.InfoLevel
	l.logentry().Infoln("log level set to:", level)

	switch level {
	case ambient.LogLevelDebug:
		l.log.Level = logrus.DebugLevel
	case ambient.LogLevelInfo:
		l.log.Level = logrus.InfoLevel
	case ambient.LogLevelWarn:
		l.log.Level = logrus.WarnLevel
	case ambient.LogLevelError:
		l.log.Level = logrus.ErrorLevel
	case ambient.LogLevelFatal:
		l.log.Level = logrus.FatalLevel
	default:
		l.log.Level = logrus.InfoLevel
	}
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

// Warn is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Warn(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Warnln(v...)
	} else {
		l.logentry().Warnf(format, v...)
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
