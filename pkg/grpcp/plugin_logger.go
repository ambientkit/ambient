package grpcp

import (
	"fmt"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

// GRPCLoggerPlugin .
type GRPCLoggerPlugin struct {
	client         protodef.LoggerClient
	appName        string
	tracerProvider *sdktrace.TracerProvider
	ctx            context.Context
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

	if l.ctx != nil {
		l.client.Debug(l.ctx, &protodef.LogFormat{Format: out})
		return
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

	if l.ctx != nil {
		l.client.Info(l.ctx, &protodef.LogFormat{Format: out})
		return
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

// clone returns a copy of the logger.
func (l *GRPCLoggerPlugin) clone() *GRPCLoggerPlugin {
	out := &GRPCLoggerPlugin{
		client:         l.client,
		appName:        l.appName,
		tracerProvider: l.tracerProvider,
	}

	return out
}

// For handler.
func (l *GRPCLoggerPlugin) For(ctx context.Context) ambient.Logger {
	if span := trace.SpanFromContext(ctx); span != nil {
		logger := l.clone()
		logger.ctx = ctx
		// TODO: Determine if these need to be saved.
		//span.SpanContext().TraceID()
		//span.SpanContext().SpanID()
		return logger
	}
	return l
}

// Trace returns a context and an OpenTelemetry span.
func (l *GRPCLoggerPlugin) Trace(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return l.tracerProvider.Tracer(l.appName).Start(ctx, spanName)
}
