package middleware

import (
	"log/slog"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(statusCode int) {
	sr.statusCode = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info("Request received",
				"method", req.Method,
				"url", req.URL.Path,
			)

			sr := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(sr, req)

			logger.Info("Request completed",
				"method", req.Method,
				"url", req.URL.Path,
				"status", sr.statusCode,
			)
		})
	}
}
