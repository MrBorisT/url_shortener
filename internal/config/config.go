package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	Port       string
	JWTSecret  string
	JWTTTL     time.Duration
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		Port:       os.Getenv("APP_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
	jwt_ttl_str := os.Getenv("JWT_TTL")
	if jwt_ttl_str == "" {
		jwt_ttl_str = "24h"
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing required JWT_SECRET environment variable")
	}

	cfg.JWTTTL = parseDuration(jwt_ttl_str)

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" ||
		cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBSSLMode == "" {
		return nil, fmt.Errorf("missing required database configuration environment variables")
	}

	if cfg.Port == "" {
		cfg.Port = ":8080"
	} else if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	return cfg, nil
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
