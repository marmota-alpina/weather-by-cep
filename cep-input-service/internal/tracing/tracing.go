package tracing

import (
	"context"
	"log"

	"github.com/marmota-alpina/weather-by-cep/cep-input-service/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func InitTracer(cfg *config.Config) *sdktrace.TracerProvider {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Erro ao inicializar OTLP exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("cep-input-service"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp
}
