package stock

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"route256/metrics/pkg/tracer"
)

func beginSpan(ctx context.Context, pkg string, methodName string) (context.Context, tracer.Span) {
	ctx, span := tracer.BeginSpan(ctx, "repository.stock."+pkg+"/"+methodName)
	span.SetAttributes(attribute.String("method", methodName))
	return ctx, span
}
