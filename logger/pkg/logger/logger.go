package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Level int8

const (
	InfoLevel  Level = iota
	WarnLevel  Level = iota
	ErrorLevel Level = iota
)

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Level() Level
}

type contextKey struct{}

var loggerContextKey = contextKey{}

var globalNoopLogger = &noop{}

func New(level Level, serviceName string) Logger {
	return newZap(level, serviceName)
}

func ToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func ContextWithLogger(ctx context.Context, level Level, serviceName string) context.Context {
	logger := New(level, serviceName)
	return ToContext(ctx, logger)
}

func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	ctxLogger := FromContext(ctx)

	mkv := getMetricsKeysValues(ctx)
	ctxLogger.Infow(msg, append(mkv, keysAndValues...)...)
}

func Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	ctxLogger := FromContext(ctx)

	mkv := getMetricsKeysValues(ctx)
	ctxLogger.Warnw(msg, append(mkv, keysAndValues...)...)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	ctxLogger := FromContext(ctx)

	mkv := getMetricsKeysValues(ctx)
	ctxLogger.Errorw(msg, append(mkv, keysAndValues...)...)
}

func getMetricsKeysValues(ctx context.Context) []interface{} {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()

	return []interface{}{"traceID", spanCtx.TraceID(), "spanID", spanCtx.SpanID()}
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerContextKey).(Logger)
	if !ok {
		return globalNoopLogger
	}

	return logger
}
