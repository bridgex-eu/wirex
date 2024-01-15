package wirex

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middleware in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func with(key string, val any) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, val))

			next.ServeHTTP(w, r)
		})
	}
}

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedWriter := NewResponseWriter(w)
			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)
			slog.Info("request", "method", r.Method, "path", r.URL.EscapedPath(), "status", wrappedWriter.StatusCode, "duration", duration)
		})
	}
}

// ResponseWriter is a custom http.ResponseWriter that captures the status code
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	// Default the status code to 200, as that is what net/http defaults to
	return &ResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code and calls the original WriteHeader
func (rw *ResponseWriter) WriteHeader(code int) {
	slog.Info("Write header", "code", code)
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
