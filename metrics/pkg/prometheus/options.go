package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func WithPrometheus() func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		mux.Handle("GET /metrics", promhttp.Handler())
	}
}
