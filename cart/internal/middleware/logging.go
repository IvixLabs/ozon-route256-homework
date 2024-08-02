package middleware

import (
	"net/http"
	"route256/logger/pkg/logger"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			ctx := req.Context()

			logger.Infow(ctx, "http request time",
				"method", req.Method,
				"uri", req.RequestURI,
				"time", time.Since(start),
			)
		}(time.Now())

		next.ServeHTTP(w, req)
	})
}
