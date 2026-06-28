package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Host           string
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	MaxConns int
	MaxIdle  int
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	Expiration time.Duration
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:           getEnv("HOST", "0.0.0.0"),
			Port:           getEnv("PORT", "8080"),
			ReadTimeout:    getEnvDuration("READ_TIMEOUT", 30*time.Second),
			WriteTimeout:   getEnvDuration("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:    getEnvDuration("IDLE_TIMEOUT", 60*time.Second),
			MaxHeaderBytes: getEnvInt("MAX_HEADER_BYTES", 1<<20),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DBHOST", ""),
			Port:     getEnv("DBPORT", "5432"),
			User:     getEnv("DBUSER", ""),
			Password: getEnv("DBPASS", ""),
			Name:     getEnv("DBNAME", ""),
			SSLMode:  getEnv("DBSSLMODE", "disable"),
			MaxConns: getEnvInt("DBMAXCONNS", 25),
			MaxIdle:  getEnvInt("DBMAXIDLE", 10),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Issuer:     getEnv("JWT_ISSUER", "tasks-api"),
			Expiration: getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err != nil {
			return duration
		}
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if number, err := strconv.Atoi(value); err != nil {
			return number
		}
	}

	return defaultValue
}
