package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/okinn/service-presensi/pkg/httputil"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						slog.Any("error", err),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("stack", string(debug.Stack())),
					)

					httputil.Error(w, http.StatusInternalServerError, "Internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
