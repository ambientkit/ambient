package zaplogger

import (
	"github.com/ambientkit/ambient"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents a logger.
type Logger struct {
	log      *zap.SugaredLogger
	loglevel zap.AtomicLevel

	appName    string
	appVersion string
}

// NewLogger returns a new logger with a default log level of error.
func NewLogger(appName string, appVersion string) *Logger {
	loglevel := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "" // Disable timestamps.
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	base := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		//zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		//zapcore.Lock(os.Stdout),
		loglevel,
	))

	defer base.Sync()
	sugar := base.Sugar()

	return &Logger{
		log:      sugar,
		loglevel: loglevel,

		appName:    appName,
		appVersion: appVersion,
	}
}

func (l *Logger) logentry() *zap.SugaredLogger {
	return l.log.Named(l.appName + " v" + l.appVersion)
}

// SetLogLevel will set the logger output level.
func (l *Logger) SetLogLevel(level ambient.LogLevel) {
	// Set log level temporarily to info.
	l.loglevel.SetLevel(zap.InfoLevel)

	var loglevel zapcore.Level

	switch level {
	case ambient.LogLevelDebug:
		loglevel = zapcore.DebugLevel
		l.logentry().Infof("zaplogger: log level set to: %v", "debug")
	case ambient.LogLevelInfo:
		loglevel = zapcore.InfoLevel
		l.logentry().Infof("zaplogger: log level set to: %v", "info")
	case ambient.LogLevelWarn:
		loglevel = zapcore.WarnLevel
		l.logentry().Infof("zaplogger: log level set to: %v", "warn")
	case ambient.LogLevelError:
		loglevel = zapcore.ErrorLevel
		l.logentry().Infof("zaplogger: log level set to: %v", "error")
	case ambient.LogLevelFatal:
		loglevel = zapcore.FatalLevel
		l.logentry().Infof("zaplogger: log level set to: %v", "fatal")
	default:
		loglevel = zapcore.InfoLevel
		l.logentry().Infof("zaplogger: log level set to: %v", "info")
	}

	l.loglevel.SetLevel(loglevel)
}

// Debug is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Debug(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Debug(v...)
	} else {
		l.logentry().Debugf(format, v...)
	}
}

// Info is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Info(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Info(v...)
	} else {
		l.logentry().Infof(format, v...)
	}
}

// Warn is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Warn(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Warn(v...)
	} else {
		l.logentry().Warnf(format, v...)
	}
}

// Error is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty.
func (l *Logger) Error(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Error(v...)
	} else {
		l.logentry().Errorf(format, v...)
	}
}

// Fatal is equivalent to log.Printf() + "\n" if format is not empty.
// It's equivalent to Println() if format is empty. It's followed by a call
// to os.Exit(1).
func (l *Logger) Fatal(format string, v ...interface{}) {
	if len(format) == 0 {
		l.logentry().Fatal(v...)
	} else {
		l.logentry().Fatalf(format, v...)
	}
}
