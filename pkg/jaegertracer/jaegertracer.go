package jaegertracer

import (
	"context"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/envdetect"
	"github.com/ambientkit/ambient/pkg/requestuuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// Provider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func Provider(log ambient.Logger, url string, service string) (*sdktrace.TracerProvider, func(ctx context.Context), error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, nil, err
	}

	uuid, _ := requestuuid.UUID()

	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.Bool("local.dev", envdetect.RunningLocalDev()),
			attribute.String("process.id", uuid),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	// Cleanly shutdown and flush telemetry when the application exits.
	f := func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Error(err.Error())
		}
	}

	return tp, f, nil
}
