package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	size        int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)

			// Get request ID if exists
			requestID := r.Header.Get("X-Request-ID")

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			// Build log attributes
			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.status),
				slog.Duration("duration", duration),
				slog.Int("size", wrapped.size),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			}

			if requestID != "" {
				attrs = append(attrs, slog.String("request_id", requestID))
			}

			if r.URL.RawQuery != "" {
				attrs = append(attrs, slog.String("query", r.URL.RawQuery))
			}

			// Get user info from context if available
			if userID := GetUserID(r.Context()); userID != "" {
				attrs = append(attrs, slog.String("user_id", userID))
			}

			// Log with appropriate level based on status code
			logAttrs := make([]any, len(attrs))
			for i, attr := range attrs {
				logAttrs[i] = attr
			}

			switch {
			case wrapped.status >= 500:
				logger.Error("HTTP Request", logAttrs...)
			case wrapped.status >= 400:
				logger.Warn("HTTP Request", logAttrs...)
			default:
				logger.Info("HTTP Request", logAttrs...)
			}
		})
	}
}
