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
	"github.com/okinn/service-presensi/internal/domain/service"
	"github.com/okinn/service-presensi/internal/infrastructure"
	"github.com/okinn/service-presensi/pkg/jwt"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := infrastructure.LoadConfig()

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

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret, time.Duration(cfg.JWTExpireMinutes)*time.Minute)

	presensiRepo := mongodb.NewPresensiRepository(db)
	userRepo := mongodb.NewUserRepository(db)
	locationRepo := mongodb.NewAllowedLocationRepository(db)

	var locationService *service.LocationService
	if cfg.GeofenceEnabled {
		locationService = service.NewLocationService(locationRepo, cfg.GeofenceEnabled)
		logger.Info("Geofencing enabled", slog.Float64("default_radius_meters", cfg.DefaultRadiusMeters))
	} else {
		logger.Info("Geofencing disabled")
	}

	analyticsRepo := mongodb.NewAnalyticsRepository(db)

	presensiUseCase := usecase.NewPresensiUseCase(presensiRepo, locationService)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtManager)
	userUseCase := usecase.NewUserUseCase(userRepo)
	analyticsUseCase := usecase.NewAnalyticsUseCase(analyticsRepo)

	presensiHandler := httpAdapter.NewPresensiHandler(presensiUseCase)
	authHandler := httpAdapter.NewAuthHandler(authUseCase)
	userHandler := httpAdapter.NewUserHandler(userUseCase)
	locationHandler := httpAdapter.NewLocationHandler(locationRepo)
	analyticsHandler := httpAdapter.NewAnalyticsHandler(analyticsUseCase)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	loginRateLimiter := middleware.NewLoginRateLimiter()

	router := httpAdapter.NewRouter(httpAdapter.RouterConfig{
		PresensiHandler:  presensiHandler,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		LocationHandler:  locationHandler,
		AnalyticsHandler: analyticsHandler,
		AuthMiddleware:   authMiddleware,
		Logger:           logger,
		LoginRateLimiter: loginRateLimiter,
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  12 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  75 * time.Second,
	}

	go func() {
		logger.Info("Server starting", slog.String("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server exited")
}
