package middleware

import (
	"net/http"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	responseStatus int
	responseSize   int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseSize += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseStatus = statusCode
}

func NewLoggingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lw := &loggingResponseWriter{
				ResponseWriter: w,
				responseStatus: http.StatusOK,
			}

			next.ServeHTTP(lw, r)

			duration := time.Since(start)

			logger.Log.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", lw.responseStatus),
				zap.Int("size", lw.responseSize),
				zap.Duration("duration", duration),
			)
		})
	}
}
