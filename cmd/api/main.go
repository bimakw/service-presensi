package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	httpAdapter "github.com/okinn/service-presensi/internal/adapter/inbound/http"
	"github.com/okinn/service-presensi/internal/adapter/inbound/http/middleware"
	"github.com/okinn/service-presensi/internal/adapter/outbound/mongodb"
	"github.com/okinn/service-presensi/internal/application/usecase"
	"github.com/okinn/service-presensi/internal/infrastructure"
	"github.com/okinn/service-presensi/pkg/jwt"
)

func main() {
	// Load .env file if exists
	_ = godotenv.Load()

	// Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg := infrastructure.LoadConfig()

	// Connect to MongoDB (outbound adapter)
	mongoClient, err := infrastructure.ConnectMongo(cfg.MongoURI)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("Error disconnecting MongoDB", slog.String("error", err.Error()))
		}
	}()

	db := mongoClient.Database(cfg.Database)

	// Initialize JWT Manager
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret, time.Duration(cfg.JWTExpireMinutes)*time.Minute)

	// Initialize layers (Dependency Injection)
	// Outbound adapter: MongoDB repository implements domain port
	presensiRepo := mongodb.NewPresensiRepository(db)
	userRepo := mongodb.NewUserRepository(db)

	// Application layer: Use case depends on domain port (not adapter)
	presensiUseCase := usecase.NewPresensiUseCase(presensiRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtManager)

	// Inbound adapter: HTTP handler depends on use case
	presensiHandler := httpAdapter.NewPresensiHandler(presensiUseCase)
	authHandler := httpAdapter.NewAuthHandler(authUseCase)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	loginRateLimiter := middleware.NewLoginRateLimiter()

	// Setup router (inbound adapter)
	router := httpAdapter.NewRouter(httpAdapter.RouterConfig{
		PresensiHandler:  presensiHandler,
		AuthHandler:      authHandler,
		AuthMiddleware:   authMiddleware,
		Logger:           logger,
		LoginRateLimiter: loginRateLimiter,
	})

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting", slog.String("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server exited")
}
