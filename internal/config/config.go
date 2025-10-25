package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Application
	AppEnv  string
	AppPort int
	AppName string
	AppURL  string

	// Database
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBMaxOpenConns int
	DBMaxIdleConns int

	// JWT
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Stripe
	StripeSecretKey     string
	StripeWebhookSecret string

	// Email
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	// S3
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3Region    string
	S3UseSSL    bool

	// DGII
	DGIIEnabled  bool
	DGIIAPI      string
	DGIICertPath string
	DGIIKeyPath  string

	// Rate Limiting
	RateLimitEnabled     bool
	RateLimitMaxRequests int
	RateLimitWindow      time.Duration

	// Logging
	LogLevel  string
	LogFormat string

	// Metrics
	MetricsEnabled bool
	MetricsPort    int
}

func Load() *Config {
	return &Config{
		// Application
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnvAsInt("APP_PORT", 8080),
		AppName: getEnv("APP_NAME", "Vendix"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),

		// Database
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "vendix"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),

		// JWT
		JWTSecret:        getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTAccessExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),

		// Redis
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// Stripe
		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),

		// Email
		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@billing.example.com"),

		// S3
		S3Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    getEnv("S3_BUCKET", "billing-documents"),
		S3Region:    getEnv("S3_REGION", "us-east-1"),
		S3UseSSL:    getEnvAsBool("S3_USE_SSL", false),

		// DGII
		DGIIEnabled:  getEnvAsBool("DGII_ENABLED", false),
		DGIIAPI:      getEnv("DGII_API_URL", "https://api.dgii.gov.do"),
		DGIICertPath: getEnv("DGII_CERT_PATH", "./certs/dgii.pem"),
		DGIIKeyPath:  getEnv("DGII_KEY_PATH", "./certs/dgii.key"),

		// Rate Limiting
		RateLimitEnabled:     getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitMaxRequests: getEnvAsInt("RATE_LIMIT_MAX_REQUESTS", 100),
		RateLimitWindow:      getEnvAsDuration("RATE_LIMIT_WINDOW", 1*time.Minute),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "debug"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// Metrics
		MetricsEnabled: getEnvAsBool("METRICS_ENABLED", true),
		MetricsPort:    getEnvAsInt("METRICS_PORT", 9090),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
