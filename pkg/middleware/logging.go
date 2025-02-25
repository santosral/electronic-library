package middleware

import (
	"electronic-library/pkg"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Logging() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := pkg.GetLoggerFromContext(r.Context())
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info(
				"request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration", &duration,
				"client_ip", r.RemoteAddr,
			)
		}
	}
}
