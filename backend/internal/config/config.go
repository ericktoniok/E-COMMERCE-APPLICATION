package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	ImageStoragePath string
	AdminEmail       string
	AdminPassword    string
	MockMpesaURL     string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://ecom:ecompwd@localhost:5432/ecom?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "change_me"),
		ImageStoragePath: getEnv("IMAGE_STORAGE_PATH", "./storage/images"),
		AdminEmail:       getEnv("ADMIN_EMAIL", "admin@local.test"),
		AdminPassword:    getEnv("ADMIN_PASSWORD", "Admin123!"),
		MockMpesaURL:     getEnv("MOCK_MPESA_URL", "http://mock-mpesa:8090"),
	}

	if err := os.MkdirAll(cfg.ImageStoragePath, 0o755); err != nil {
		log.Printf("warn: creating image storage dir: %v", err)
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
