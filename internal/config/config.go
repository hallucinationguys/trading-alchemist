package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Server configuration
	Server ServerConfig

	// Database configuration
	Database DatabaseConfig

	// Email configuration (Resend only)
	Email EmailConfig

	// JWT configuration
	JWT JWTConfig

	// App configuration
	App AppConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxConns     int
	MinConns     int
	MaxConnLife  string
	MaxConnIdle  string
}

type EmailConfig struct {
	// Resend configuration
	ResendAPIKey string
	FromEmail    string
	FromName     string
}

type JWTConfig struct {
	Secret string
	TTL    string
}

type AppConfig struct {
	Name         string
	Environment  string
	BaseURL      string
	MagicLinkTTL string
}

// Load loads configuration from environment variables using Viper
func Load() *Config {
	// Initialize Viper
	v := viper.New()

	// Set default environment if not specified
	env := getEnv("APP_ENV", "development")

	// Configure Viper
	configureViper(v, env)

	return &Config{
		Server: ServerConfig{
			Host:         v.GetString("SERVER_HOST"),
			Port:         v.GetString("SERVER_PORT"),
			ReadTimeout:  v.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("SERVER_WRITE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:         v.GetString("DB_HOST"),
			Port:         v.GetString("DB_PORT"),
			User:         v.GetString("DB_USER"),
			Password:     v.GetString("DB_PASSWORD"),
			Name:         v.GetString("DB_NAME"),
			SSLMode:      v.GetString("DB_SSL_MODE"),
			MaxConns:     v.GetInt("DB_MAX_CONNS"),
			MinConns:     v.GetInt("DB_MIN_CONNS"),
			MaxConnLife:  v.GetString("DB_MAX_CONN_LIFE"),
			MaxConnIdle:  v.GetString("DB_MAX_CONN_IDLE"),
		},
		Email: EmailConfig{
			ResendAPIKey: v.GetString("RESEND_API_KEY"),
			FromEmail:    v.GetString("FROM_EMAIL"),
			FromName:     v.GetString("FROM_NAME"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("JWT_SECRET"),
			TTL:    v.GetString("JWT_TTL"),
		},
		App: AppConfig{
			Name:         v.GetString("APP_NAME"),
			Environment:  v.GetString("APP_ENV"),
			BaseURL:      v.GetString("APP_BASE_URL"),
			MagicLinkTTL: v.GetString("MAGIC_LINK_TTL"),
		},
	}
}

// configureViper sets up Viper configuration
func configureViper(v *viper.Viper, env string) {
	// Set the config name and type
	v.SetConfigType("env")

	// Add config paths
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")

	// Try to read environment-specific config file first
	configFile := fmt.Sprintf("env.%s", env)
	v.SetConfigName(configFile)

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file %s: %v", configFile, err)
		log.Printf("Falling back to environment variables only")
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	// Automatically read environment variables
	v.AutomaticEnv()

	// Set default values
	setDefaults(v)
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("SERVER_HOST", "localhost")
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("SERVER_READ_TIMEOUT", "10s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "10s")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5433")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "trading_alchemist_db")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_CONNS", 25)
	v.SetDefault("DB_MIN_CONNS", 5)
	v.SetDefault("DB_MAX_CONN_LIFE", "1h")
	v.SetDefault("DB_MAX_CONN_IDLE", "30m")

	// Email defaults
	v.SetDefault("RESEND_API_KEY", "")
	v.SetDefault("FROM_EMAIL", "noreply@example.com")
	v.SetDefault("FROM_NAME", "Trading Alchemist")

	// JWT defaults
	v.SetDefault("JWT_SECRET", "your-secret-key")
	v.SetDefault("JWT_TTL", "24h")

	// App defaults
	v.SetDefault("APP_NAME", "Trading Alchemist")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_BASE_URL", "http://localhost:8080")
	v.SetDefault("MAGIC_LINK_TTL", "15m")
}

// LoadForEnvironment loads configuration for a specific environment
func LoadForEnvironment(env string) *Config {
	// Set the environment variable
	os.Setenv("APP_ENV", env)
	return Load()
}

// ValidateConfig validates the loaded configuration
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key" && c.App.Environment == "production" {
		return fmt.Errorf("JWT secret must be changed for production environment")
	}

	if c.Email.ResendAPIKey == "" {
		log.Printf("Warning: RESEND_API_KEY is not set, email functionality will not work")
	}

	if c.Database.Password == "postgres" && c.App.Environment == "production" {
		return fmt.Errorf("database password must be changed for production environment")
	}

	return nil
}

// GetConfigPath returns the path to the config file for the given environment
func GetConfigPath(env string) string {
	return filepath.Join("configs", fmt.Sprintf("env.%s", env))
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Legacy helper functions (kept for compatibility but deprecated)
func getIntEnv(key string, defaultValue int) int {
	log.Printf("Warning: getIntEnv is deprecated, use Viper configuration instead")
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	log.Printf("Warning: getDurationEnv is deprecated, use Viper configuration instead")
	return defaultValue
} 