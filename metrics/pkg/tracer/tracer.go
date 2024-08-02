package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"route256/common/pkg/env"
)

const (
	EnvAppTracer         = "APP_TRACER"
	EnvAppJaegerEndpoint = "APP_JAEGER_ENDPOINT"
)

type ShutdownTracerProvider interface {
	Shutdown(ctx context.Context) error
}

var globalIsTracer bool

func IsTracer() bool {
	return globalIsTracer
}

func InitTracerProvider(ctx context.Context, serviceName string) {
	if !isEnvTracer() {
		return
	}
	globalIsTracer = true

	endpoint := env.GetEnvVar(EnvAppJaegerEndpoint)

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(endpoint),
	)
	if err != nil {
		log.Panicln(err)
	}

	resourceObj, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		log.Panicln(err)
	}

	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(resourceObj),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func Shutdown(ctx context.Context) error {

	tp, ok := otel.GetTracerProvider().(ShutdownTracerProvider)

	if !ok {
		return nil
	}

	err := tp.Shutdown(ctx)

	if err != nil {
		return err
	}

	return nil
}

type Span interface {
	End(options ...trace.SpanEndOption)
	SetAttributes(kv ...attribute.KeyValue)
}

type noopSpan struct {
}

func (*noopSpan) End(_ ...trace.SpanEndOption) {

}

func (*noopSpan) SetAttributes(_ ...attribute.KeyValue) {

}

func BeginSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, Span) {
	if !IsTracer() {
		return ctx, &noopSpan{}
	}

	tracer := otel.GetTracerProvider().Tracer("")

	ctx, spanObj := tracer.Start(ctx, spanName, opts...)
	return ctx, spanObj
}

func isEnvTracer() bool {
	v, ok := os.LookupEnv(EnvAppTracer)
	if !ok {
		return true
	}

	return v != "disabled"
}
