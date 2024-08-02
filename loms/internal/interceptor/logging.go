package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"route256/logger/pkg/logger"
	"time"
)

func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func(start time.Time) {
		logger.Infow(ctx, "grpc request time",
			"method", info.FullMethod,
			"time", time.Since(start),
		)
	}(time.Now())

	return handler(ctx, req)
}
