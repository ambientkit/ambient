package ambient

import (
	"os"
	"strings"
)

// AppLogger represents the log service for the app.
type AppLogger interface {
	Logger

	// Fatal is reserved for the app level only.
	Fatal(format string, v ...interface{})
	SetLogLevel(level LogLevel)
}

// Logger represents the log service for the plugins.
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
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
