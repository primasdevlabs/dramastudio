package tracing

import (
	"context"
	"io"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is the minimal span contract used by application code.
type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, func())
}

type otelTracer struct {
	t trace.Tracer
}

// New returns an OTel-backed tracer for the given instrumentation scope.
func New(scope string) Tracer {
	return &otelTracer{t: otel.Tracer(scope)}
}

func (t *otelTracer) StartSpan(ctx context.Context, name string) (context.Context, func()) {
	ctx, span := t.t.Start(ctx, name)
	return ctx, func() { span.End() }
}

// InitProvider installs a global TracerProvider. exporter "stdout" emits
// spans to w; anything else disables export (no-op provider).
func InitProvider(ctx context.Context, serviceName, exporter string, w io.Writer) (func(context.Context) error, error) {
	if exporter != "stdout" {
		otel.SetTracerProvider(trace.NewNoopTracerProvider())
		return func(context.Context) error { return nil }, nil
	}
	exp, err := stdouttrace.New(stdouttrace.WithWriter(w))
	if err != nil {
		return nil, err
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(serviceName)),
	)
	if err != nil {
		log.Printf("tracing: resource merge failed: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}
