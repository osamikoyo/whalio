package config

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server settings
	Port         string        `json:"port"`
	Host         string        `json:"host"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`

	// Sqlite db file path
	DatabasePath string `json:"database_path"`
	// Development settings
	Debug       bool   `json:"debug"`
	Environment string `json:"environment"`

	//Storage for images
	ImageDir string `json:"image_dir"`
	// Storage for songs
	UploadDir string `json:"upload_dir"`
	// Static files
	StaticDir string `json:"static_dir"`

	// CORS settings
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`

	// Rate limiting
	RateLimitEnabled bool `json:"rate_limit_enabled"`
	RateLimit        int  `json:"rate_limit"`

	// Logging
	LogLevel  string `json:"log_level"`
	LogFormat string `json:"log_format"` // "json" or "console"
}

// Default configuration values
const (
	DefaultPort         = "8080"
	DefaultHost         = "localhost"
	DefaultReadTimeout  = 15 * time.Second
	DefaultWriteTimeout = 15 * time.Second
	DefaultIdleTimeout  = 60 * time.Second
	DefaultEnvironment  = "development"
	DefaultStaticDir    = "./static"
	DefaultImageDir     = "./images"
	DefaultUploadDir    = "./files"
	DefaultLogLevel     = "info"
	DefaultLogFormat    = "console"
	DefaultRateLimit    = 100
	DefaultDatabasePath = "database.db"
)

// Load loads configuration from environment variables and command line flags
func Load() *Config {
	cfg := &Config{
		Port:             getEnv("PORT", DefaultPort),
		Host:             getEnv("HOST", DefaultHost),
		ReadTimeout:      getDurationEnv("READ_TIMEOUT", DefaultReadTimeout),
		WriteTimeout:     getDurationEnv("WRITE_TIMEOUT", DefaultWriteTimeout),
		IdleTimeout:      getDurationEnv("IDLE_TIMEOUT", DefaultIdleTimeout),
		Debug:            getBoolEnv("DEBUG", false),
		UploadDir:        getEnv("UPLOAD_DIR", DefaultUploadDir),
		Environment:      getEnv("ENVIRONMENT", DefaultEnvironment),
		DatabasePath:     getEnv("DATABASE_PATH", DefaultDatabasePath),
		StaticDir:        getEnv("STATIC_DIR", DefaultStaticDir),
		ImageDir:         getEnv("IMAGE_DIR", DefaultImageDir),
		LogLevel:         getEnv("LOG_LEVEL", DefaultLogLevel),
		LogFormat:        getEnv("LOG_FORMAT", DefaultLogFormat),
		RateLimitEnabled: getBoolEnv("RATE_LIMIT_ENABLED", true),
		RateLimit:        getIntEnv("RATE_LIMIT", DefaultRateLimit),
	}

	// Set default CORS settings
	cfg.AllowedOrigins = getSliceEnv("ALLOWED_ORIGINS", []string{"*"})
	cfg.AllowedMethods = getSliceEnv("ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	cfg.AllowedHeaders = getSliceEnv("ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "HX-Request", "HX-Trigger", "HX-Target"})

	// Parse command line flags
	parseFlags(cfg)

	return cfg
}

// parseFlags parses command line flags and overrides config values
func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.Port, "port", cfg.Port, "Server port")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "Server host")
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "Enable debug mode")
	flag.StringVar(&cfg.Environment, "env", cfg.Environment, "Environment (development, staging, production)")
	flag.StringVar(&cfg.StaticDir, "static-dir", cfg.StaticDir, "Static files directory")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Log level (debug, info, warn, error)")
	flag.StringVar(&cfg.LogFormat, "log-format", cfg.LogFormat, "Log format (json, console)")

	flag.Parse()
}

// Address returns the server address string
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if c.Environment == "" {
		return fmt.Errorf("environment cannot be empty")
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !slices.Contains(validLogLevels, c.LogLevel) {
		return fmt.Errorf("invalid log level: %s (valid: %v)", c.LogLevel, validLogLevels)
	}

	// Validate log format
	if c.LogFormat != "json" && c.LogFormat != "console" {
		return fmt.Errorf("invalid log format: %s (valid: json, console)", c.LogFormat)
	}

	return nil
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple split by comma - in production you might want more sophisticated parsing
		result := []string{}
		for _, v := range []string{value} { // Simplified for now
			result = append(result, v)
		}
		return result
	}
	return defaultValue
}
