package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

var requestCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "app_cart",
		Name:      "http_request_total_counter",
		Help:      "Total amount of http requests",
	},
	[]string{"handler", "statusCode", "method"},
)

var requestHistogramVec = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "app_cart",
		Name:      "http_request_duration_histogram",
		Buckets:   prometheus.DefBuckets,
		Help:      "Duration of http request",
	},
	[]string{"handler", "statusCode", "method"})

func RequestCounter(next http.Handler, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logWriter := newLoggingResponseWriter(w)

		labelsValues := []string{path, strconv.Itoa(logWriter.statusCode), req.Method}

		defer func(createdAt time.Time) {
			requestHistogramVec.WithLabelValues(labelsValues...).Observe(time.Since(createdAt).Seconds())
		}(time.Now())

		next.ServeHTTP(logWriter, req)

		requestCounterVec.WithLabelValues(labelsValues...).Inc()
	})
}
