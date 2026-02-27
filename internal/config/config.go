package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Supabase    SupabaseConfig
	JWT         JWTConfig
	GoogleOAuth GoogleOAuthConfig
	CORS        CORSConfig
	Storage     StorageConfig
	Logging     LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port       string
	Env        string
	APIVersion string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI      string
	Database string
}

// SupabaseConfig holds Supabase configuration
type SupabaseConfig struct {
	URL            string
	AnonKey        string
	ServiceRoleKey string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

// GoogleOAuthConfig holds Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	MaxUploadSize    int64
	AllowedFileTypes []string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:       getEnv("PORT", "8080"),
			Env:        getEnv("ENV", "development"),
			APIVersion: getEnv("API_VERSION", "v1"),
		},
		Database: DatabaseConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "fmp_app"),
		},
		Supabase: SupabaseConfig{
			URL:            getEnv("SUPABASE_URL", ""),
			AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "change-this-secret"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		GoogleOAuth: GoogleOAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("GOOGLE_REDIRECT_URI", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		},
		Storage: StorageConfig{
			MaxUploadSize:    getEnvAsInt64("MAX_UPLOAD_SIZE", 5242880), // 5MB
			AllowedFileTypes: getEnvAsSlice("ALLOWED_FILE_TYPES", []string{"image/jpeg", "image/png"}),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.URI == "" {
		return fmt.Errorf("MONGODB_URI is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("MONGODB_DATABASE is required")
	}
	if c.Supabase.URL == "" {
		return fmt.Errorf("SUPABASE_URL is required")
	}
	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
