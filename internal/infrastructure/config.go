package infrastructure

import (
	"context"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Port                string
	MongoURI            string
	Database            string
	JWTSecret           string
	JWTExpireMinutes    int
	GeofenceEnabled     bool
	DefaultRadiusMeters float64
}

func LoadConfig() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:            getEnv("DATABASE", "presensi_db"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpireMinutes:    getEnvAsInt("JWT_EXPIRE_MINUTES", 60*24), // 24 hours default
		GeofenceEnabled:     getEnvAsBool("GEOFENCE_ENABLED", false),
		DefaultRadiusMeters: getEnvAsFloat("DEFAULT_RADIUS_METERS", 100), // 100 meters default
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func ConnectMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
