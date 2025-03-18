package telemetry

import (
	"context"
	"fizzbuzz-server/internal/config"
	"log"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracerProvider *sdktrace.TracerProvider

// InitTracer initializes the OpenTelemetry tracer
func InitTracer() func(context.Context) error {
	cfg := config.Get()

	// Parse resource attributes
	attrs := parseResourceAttributes(cfg.Telemetry.ResourceAttributes)

	// Create a new resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(attrs...),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.Telemetry.ServiceName),
		),
	)
	if err != nil {
		log.Printf("Failed to create resource: %v", err)
		return func(context.Context) error { return nil }
	}

	// Set up a connection to the OTLP exporter
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, cfg.Telemetry.OTLPEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("Failed to connect to OTLP endpoint: %v", err)
		return func(context.Context) error { return nil }
	}

	// Create the OTLP exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Printf("Failed to create OTLP exporter: %v", err)
		return func(context.Context) error { return nil }
	}

	// Create the tracer provider
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set the global tracer provider and propagator
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return a function to shut down the tracer provider
	return func(ctx context.Context) error {
		return tracerProvider.Shutdown(ctx)
	}
}

// parseResourceAttributes parses a comma-separated list of key=value pairs into resource attributes
func parseResourceAttributes(attrString string) []attribute.KeyValue {
	if attrString == "" {
		return nil
	}

	var attrs []attribute.KeyValue
	pairs := strings.Split(attrString, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		attrs = append(attrs, attribute.String(key, value))
	}
	return attrs
}