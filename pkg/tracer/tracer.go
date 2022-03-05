package tracer

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Set global trace provider
func setTracerProvider(name, url string) (trace.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(name),
			attribute.String("env", "dev"),
		)),
	)

	return tp, nil
}

func NewTracerProvider(name, url string) (trace.TracerProvider, error) {
	if os.Getenv("jaeger_url") != "" {
		url = os.Getenv("jaeger_url")
	}
	tp, err := setTracerProvider(name, url)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	return tp, nil
}
