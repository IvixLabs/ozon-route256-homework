package interceptor

import (
	"context"
	"google.golang.org/grpc/stats"
	"route256/logger/pkg/logger"
)

type Logger struct {
	logger logger.Logger
}

func NewLogger(logger logger.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (l *Logger) HandleRPC(ctx context.Context, stat stats.RPCStats) {
}
func (l *Logger) TagConn(ctx context.Context, tag *stats.ConnTagInfo) context.Context {

	return logger.ToContext(ctx, l.logger)
}
func (l *Logger) HandleConn(ctx context.Context, conn stats.ConnStats) {
}
