package grpcp

import (
	"fmt"

	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
)

// GRPCLoggerPlugin .
type GRPCLoggerPlugin struct {
	client protodef.LoggerClient
}

// Debug .
func (l *GRPCLoggerPlugin) Debug(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}

	l.client.Debug(context.Background(), &protodef.LogFormat{Format: out})
}

// Info .
func (l *GRPCLoggerPlugin) Info(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}

	l.client.Info(context.Background(), &protodef.LogFormat{Format: out})
}

// Warn .
func (l *GRPCLoggerPlugin) Warn(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}

	l.client.Warn(context.Background(), &protodef.LogFormat{Format: out})
}

// Error .
func (l *GRPCLoggerPlugin) Error(format string, v ...interface{}) {
	out := ""
	if len(format) == 0 {
		out = fmt.Sprintln(v...)
	} else {
		out = fmt.Sprintf(format, v...)
	}
	l.client.Error(context.Background(), &protodef.LogFormat{Format: out})
}
