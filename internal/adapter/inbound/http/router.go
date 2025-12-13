package http

import (
	"log/slog"
	"net/http"

	"github.com/okinn/service-presensi/internal/adapter/inbound/http/middleware"
)

type RouterConfig struct {
	PresensiHandler  *PresensiHandler
	AuthHandler      *AuthHandler
	AuditHandler     *AuditHandler
	LocationHandler  *LocationHandler
	AuthMiddleware   *middleware.AuthMiddleware
	AuditMiddleware  *middleware.AuditMiddleware
	Logger           *slog.Logger
	LoginRateLimiter *middleware.LoginRateLimiter
}

func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		Success(w, http.StatusOK, "OK", nil)
	})

	// Auth routes (public)
	mux.HandleFunc("POST /api/auth/register", cfg.AuthHandler.Register)

	// Login with rate limiting
	mux.Handle("POST /api/auth/login", cfg.LoginRateLimiter.Limit(
		http.HandlerFunc(cfg.AuthHandler.Login),
	))

	// Profile route (protected)
	mux.Handle("GET /api/auth/profile", cfg.AuthMiddleware.Authenticate(
		http.HandlerFunc(cfg.AuthHandler.GetProfile),
	))

	// Presensi routes (protected)
	mux.Handle("POST /api/presensi", cfg.AuthMiddleware.Authenticate(
		http.HandlerFunc(cfg.PresensiHandler.Create),
	))
	mux.Handle("GET /api/presensi", cfg.AuthMiddleware.Authenticate(
		http.HandlerFunc(cfg.PresensiHandler.GetAll),
	))
	mux.Handle("GET /api/presensi/{id}", cfg.AuthMiddleware.Authenticate(
		http.HandlerFunc(cfg.PresensiHandler.GetByID),
	))
	mux.Handle("PUT /api/presensi/{id}", cfg.AuthMiddleware.Authenticate(
		http.HandlerFunc(cfg.PresensiHandler.Update),
	))
	mux.Handle("DELETE /api/presensi/{id}", cfg.AuthMiddleware.Authenticate(
		cfg.AuthMiddleware.RequireRole("admin")(
			http.HandlerFunc(cfg.PresensiHandler.Delete),
		),
	))
	mux.Handle("POST /api/presensi/{id}/checkin", cfg.AuthMiddleware.Authenticate(
		http.HandlerFunc(cfg.PresensiHandler.CheckIn),
	))
	mux.Handle("POST /api/presensi/{id}/checkout", cfg.AuthMiddleware.Authenticate(
		http.HandlerFunc(cfg.PresensiHandler.CheckOut),
	))

	// Audit routes (admin only)
	if cfg.AuditHandler != nil {
		mux.Handle("GET /api/audit", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.AuditHandler.GetAll),
			),
		))
		mux.Handle("GET /api/audit/{id}", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.AuditHandler.GetByID),
			),
		))
		mux.Handle("GET /api/audit/entity", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.AuditHandler.GetByEntity),
			),
		))
		mux.Handle("GET /api/audit/user/{user_id}", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.AuditHandler.GetByUser),
			),
		))
	}

	// Location routes (admin only) - Geofencing management
	if cfg.LocationHandler != nil {
		mux.Handle("POST /api/locations", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.LocationHandler.Create),
			),
		))
		mux.Handle("GET /api/locations", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.LocationHandler.GetAll),
			),
		))
		mux.Handle("GET /api/locations/{id}", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.LocationHandler.GetByID),
			),
		))
		mux.Handle("PUT /api/locations/{id}", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.LocationHandler.Update),
			),
		))
		mux.Handle("DELETE /api/locations/{id}", cfg.AuthMiddleware.Authenticate(
			cfg.AuthMiddleware.RequireRole("admin")(
				http.HandlerFunc(cfg.LocationHandler.Delete),
			),
		))
	}

	// Apply global middlewares
	var handler http.Handler = mux

	// CORS middleware
	corsConfig := middleware.DefaultCORSConfig()
	handler = middleware.CORS(corsConfig)(handler)

	// Audit middleware (logs all write operations)
	if cfg.AuditMiddleware != nil {
		handler = cfg.AuditMiddleware.Audit(handler)
	}

	// Logging middleware
	if cfg.Logger != nil {
		handler = middleware.Logging(cfg.Logger)(handler)
		handler = middleware.Recovery(cfg.Logger)(handler)
	}

	// Global rate limiter
	rateLimiter := middleware.NewRateLimiter(middleware.DefaultRateLimiterConfig())
	handler = rateLimiter.Limit(handler)

	return handler
}
