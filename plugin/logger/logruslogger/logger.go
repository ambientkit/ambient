package logruslogger

import (
	"github.com/ambientkit/ambient"
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

	switch level {
	case ambient.LogLevelDebug:
		l.log.Level = logrus.DebugLevel
		l.logentry().Infoln("logruslogger: log level set to:", "debug")
	case ambient.LogLevelInfo:
		l.log.Level = logrus.InfoLevel
		l.logentry().Infoln("logruslogger: log level set to:", "info")
	case ambient.LogLevelWarn:
		l.log.Level = logrus.WarnLevel
		l.logentry().Infoln("logruslogger: log level set to:", "warn")
	case ambient.LogLevelError:
		l.log.Level = logrus.ErrorLevel
		l.logentry().Infoln("logruslogger: log level set to:", "error")
	case ambient.LogLevelFatal:
		l.log.Level = logrus.FatalLevel
		l.logentry().Infoln("loglogrusloggerrus: log level set to:", "fatal")
	default:
		l.log.Level = logrus.InfoLevel
		l.logentry().Infoln("logruslogger: log level set to:", "info")
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
