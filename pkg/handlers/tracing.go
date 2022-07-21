package handlers

import (
	"context"
	"time"

	"github.com/micro-eshop/catalog/pkg/core/common"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
)

func InitPrivder(ctx context.Context) func() {
	log.Infoln("Initializing OpenTelemetry")
	exp, err := newShoppingListStorageExporter(ctx)
	if err != nil {
		handleErr(err, "failed to create exporter")
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tp)
	b3p := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3p))
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.WithError(err).Fatal("failed to shutdown tracing")
		}
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func newResource() *resource.Resource {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("catalog"),
		semconv.ServiceVersionKey.String("v0.1.0"),
		attribute.String("environment", "demo"),
	)
	return r
}

func newShoppingListStorageExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(common.GetEnvOrDefault("OTEL_TRACER_ENDPOINT", "otel-collector:4317")),
		otlptracegrpc.WithDialOption(grpc.WithBackoffMaxDelay(time.Second*10)),
	)
}
