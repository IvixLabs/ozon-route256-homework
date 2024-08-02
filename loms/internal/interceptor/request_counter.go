package interceptor

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"time"
)

var requestCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "app_loms",
		Name:      "grpc_request_total_counter",
		Help:      "Total amount of grpc requests",
	},
	[]string{"method"},
)

var requestHistogramVec = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "app_loms",
		Name:      "grpc_request_duration_histogram",
		Buckets:   prometheus.DefBuckets,
		Help:      "Duration of grpc request",
	},
	[]string{"method"})

func RequestCounter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	labelsValues := []string{info.FullMethod}

	defer requestCounterVec.WithLabelValues(labelsValues...).Inc()

	defer func(createdAt time.Time) {
		requestHistogramVec.WithLabelValues(labelsValues...).Observe(time.Since(createdAt).Seconds())
	}(time.Now())

	return handler(ctx, req)
}
