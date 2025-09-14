package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Game     GameConfig
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port int
}

// GameConfig holds game-specific configuration
type GameConfig struct {
	MaxGuesses int
	WordLength int
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Database: DatabaseConfig{
			Host:            getEnvString("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			Name:            getEnvString("DB_NAME", "wordle"),
			User:            getEnvString("DB_USER", "wordle_user"),
			Password:        getEnvString("DB_PASSWORD", "wordle_password"),
			SSLMode:         getEnvString("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", "1h"),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", "15m"),
		},
		Server: ServerConfig{
			Host: getEnvString("HOST", "localhost"),
			Port: getEnvInt("PORT", 8080),
		},
		Game: GameConfig{
			MaxGuesses: getEnvInt("MAX_GUESSES", 6),
			WordLength: getEnvInt("WORD_LENGTH", 5),
		},
	}

	return config, nil
}

// ConnectionString returns a PostgreSQL connection string
func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// DatabaseURL returns a database URL in the format postgres://user:password@host:port/dbname
func (d *DatabaseConfig) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// Address returns the server address in host:port format
func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Helper functions for environment variable parsing

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	// Parse default value
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Hour // fallback
}

