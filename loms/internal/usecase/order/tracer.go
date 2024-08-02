package order

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"route256/metrics/pkg/tracer"
)

func beginSpan(ctx context.Context, methodName string) (context.Context, tracer.Span) {
	ctx, span := tracer.BeginSpan(ctx, "usecase.order/"+methodName)
	span.SetAttributes(attribute.String("method", methodName))
	return ctx, span
}
