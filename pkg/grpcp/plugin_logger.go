package grpcp

import (
	"fmt"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
)

// GRPCLoggerPlugin .
type GRPCLoggerPlugin struct {
	client protodef.LoggerClient
}

// Log handler.
func (l *GRPCLoggerPlugin) Log(level ambient.LogLevel, format string, v ...interface{}) {
	switch level {
	case ambient.LogLevelDebug:
		l.Debug(format, v...)
	case ambient.LogLevelInfo:
		l.Info(format, v...)
	case ambient.LogLevelWarn:
		l.Warn(format, v...)
	case ambient.LogLevelError:
		l.Error(format, v...)
	default:
		l.Info(format, v...)
	}
}

// Debug handler.
func (l *GRPCLoggerPlugin) Debug(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}

	l.client.Debug(context.Background(), &protodef.LogFormat{Format: out})
}

// Info handler.
func (l *GRPCLoggerPlugin) Info(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}

	l.client.Info(context.Background(), &protodef.LogFormat{Format: out})
}

// Warn handler.
func (l *GRPCLoggerPlugin) Warn(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}

	l.client.Warn(context.Background(), &protodef.LogFormat{Format: out})
}

// Error handler.
func (l *GRPCLoggerPlugin) Error(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}
	l.client.Error(context.Background(), &protodef.LogFormat{Format: out})
}
