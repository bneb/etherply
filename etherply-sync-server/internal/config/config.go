// Package config provides centralized configuration for the EtherPly sync server.
// Configuration is loaded from environment variables with sensible defaults.
package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/sync"
)

// Config holds all server configuration.
type Config struct {
	// Server settings
	Port            string
	ShutdownTimeout time.Duration

	// Authentication
	JWTSecret string

	// Storage
	BadgerPath string

	// Sync strategy (automerge, lww, server-auth)
	SyncStrategy sync.StrategyType

	// Replication
	NATSURLs []string
	Region   string
	ServerID string

	// Webhooks
	WebhookURL string

	// Logging
	LogFormat string // "json" or "text"
	LogLevel  slog.Level
}

// Load reads configuration from environment variables.
func Load() *Config {
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT_SECONDS", 30*time.Second),
		JWTSecret:       os.Getenv("ETHERPLY_JWT_SECRET"),
		BadgerPath:      getEnv("BADGER_PATH", "./badger.db"),
		SyncStrategy:    sync.StrategyType(getEnv("SYNC_STRATEGY", string(sync.StrategyAutomerge))),
		Region:          getEnv("REGION", "default"),
		ServerID:        os.Getenv("SERVER_ID"),
		WebhookURL:      os.Getenv("WEBHOOK_URL"),
		LogFormat:       getEnv("LOG_FORMAT", "json"),
		LogLevel:        parseLogLevel(getEnv("LOG_LEVEL", "info")),
	}

	// Parse NATS_URL as comma-separated list
	if natsURL := os.Getenv("NATS_URL"); natsURL != "" {
		cfg.NATSURLs = splitTrim(natsURL, ",")
	}

	// Generate server ID if not set
	if cfg.ServerID == "" && len(cfg.NATSURLs) > 0 {
		cfg.ServerID = "sync-server-" + cfg.Port
	}

	return cfg
}

// NewLogger creates a logger based on configuration.
func (c *Config) NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{Level: c.LogLevel}

	var handler slog.Handler
	if c.LogFormat == "text" {
		handler = slog.NewTextHandler(os.Stderr, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

// Validate checks required configuration.
func (c *Config) Validate() error {
	// JWT secret is required in production
	if c.JWTSecret == "" {
		return &ConfigError{Field: "ETHERPLY_JWT_SECRET", Message: "required for secure operation"}
	}

	// Validate strategy
	switch c.SyncStrategy {
	case sync.StrategyAutomerge, sync.StrategyLWW, sync.StrategyServerAuthoritative:
		// Valid
	default:
		return &ConfigError{
			Field:   "SYNC_STRATEGY",
			Message: "must be one of: automerge, lww, server-auth",
		}
	}

	return nil
}

// ConfigError represents a configuration validation error.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if seconds, err := strconv.Atoi(v); err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}

func parseLogLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func splitTrim(s, sep string) []string {
	if s == "" {
		return nil
	}
	parts := make([]string, 0)
	for _, p := range splitString(s, sep) {
		if trimmed := trimSpace(p); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	result := make([]string, 0)
	idx := 0
	for {
		i := indexOf(s[idx:], sep)
		if i == -1 {
			result = append(result, s[idx:])
			break
		}
		result = append(result, s[idx:idx+i])
		idx += i + len(sep)
	}
	return result
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
