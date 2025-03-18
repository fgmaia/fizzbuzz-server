package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig
	Telemetry  TelemetryConfig
	Prometheus PrometheusConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// TelemetryConfig holds OpenTelemetry configuration
type TelemetryConfig struct {
	OTLPEndpoint      string
	ServiceName       string
	ResourceAttributes string
}

// PrometheusConfig holds Prometheus configuration
type PrometheusConfig struct {
	Enabled      bool
	Endpoint     string
	PushGateway  string
	PushInterval time.Duration
}

// Global configuration instance
var config *Config

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config = &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Telemetry: TelemetryConfig{
			OTLPEndpoint:       getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317"),
			ServiceName:        getEnv("OTEL_SERVICE_NAME", "fizzbuzz-server"),
			ResourceAttributes: getEnv("OTEL_RESOURCE_ATTRIBUTES", "service.version=1.0.0,deployment.environment=development"),
		},
		Prometheus: PrometheusConfig{
			Enabled:      getBoolEnv("PROMETHEUS_ENABLED", true),
			Endpoint:     getEnv("PROMETHEUS_ENDPOINT", "/metrics"),
			PushGateway:  getEnv("PROMETHEUS_PUSH_GATEWAY", "http://localhost:9091"),
			PushInterval: getDurationEnv("PROMETHEUS_PUSH_INTERVAL", 10*time.Second),
		},
	}

	return config, nil
}

// Get returns the current configuration
func Get() *Config {
	if config == nil {
		_, _ = Load()
	}
	return config
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}