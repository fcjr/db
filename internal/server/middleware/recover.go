package middleware

import (
	"log/slog"
	"net/http"
)

// WithRecovery is a middleware that recovers from panics to prevent the server from crashing
// in case of unexpected errors. It logs the panic and returns a 500 Internal Server Error response.
func WithRecovery(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Recovered from panic", "error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, req)
		})
	}
}
