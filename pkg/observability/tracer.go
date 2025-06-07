package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TraceProviderConfig struct {
	AppName           string
	AppVersion        string
	AppEnv            string
	ObservabilityMode string
	OtlpEndpoint      string
}

func InitTracerProvider(config TraceProviderConfig) (*sdktrace.TracerProvider, error) {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	if config.ObservabilityMode == "console" {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	} else if config.ObservabilityMode == "otlp" {
		exporter, err = otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.OtlpEndpoint),
		)
	}

	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource(config.AppName, config.AppVersion, config.AppEnv)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider, nil
}
