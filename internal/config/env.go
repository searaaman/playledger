package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds every setting the app reads from the environment.
// All os.Getenv calls live here so handlers and services never touch the environment directly.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ServerPort string

	JWTSecret string
	JWTTTL    time.Duration

	CORSOrigins       []string
	AllowRegistration bool
}

var Cfg Config

// Load reads configuration from the environment. A .env file in the working
// directory is loaded first if present; real environment variables win over it.
func Load() error {
	_ = godotenv.Load()

	ttlHours, err := strconv.Atoi(getEnv("JWT_TTL_HOURS", "168"))
	if err != nil {
		return errors.New("JWT_TTL_HOURS must be a whole number")
	}

	Cfg = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "playledger"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		ServerPort: getEnv("SERVER_PORT", "8080"),

		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTTTL:    time.Duration(ttlHours) * time.Hour,

		CORSOrigins:       strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
		AllowRegistration: getEnv("ALLOW_REGISTRATION", "true") == "true",
	}

	if Cfg.JWTSecret == "" {
		return errors.New("JWT_SECRET must be set")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
