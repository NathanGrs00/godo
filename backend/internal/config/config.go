// Package config loads the server's settings from environment variables.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds every setting the server needs.
type Config struct {
	// Port is the HTTP port to listen on (PORT, default "8080").
	Port string
	// MongoURI is the MongoDB connection string (MONGO_URI, required).
	MongoURI string
	// MongoDatabase is the database name (MONGO_DB, default "godo").
	MongoDatabase string
	// JWTSecret signs login tokens (JWT_SECRET, required, at least 32 characters).
	JWTSecret string
	// TokenTTL is how long a login stays valid (TOKEN_TTL, default "168h").
	TokenTTL time.Duration
	// AllowedOrigins are the frontend URLs allowed by CORS
	// (CORS_ORIGINS, comma separated, default "http://localhost:5173").
	AllowedOrigins []string
}

// minSecretLength is the shortest JWT secret accepted.
const minSecretLength = 32

// Load reads the configuration using os.Getenv.
func Load() (Config, error) {
	return FromLookup(os.Getenv)
}

// FromLookup reads the configuration using getenv, which makes it easy to test.
func FromLookup(getenv func(string) string) (Config, error) {
	cfg := Config{
		Port:          withDefault(getenv("PORT"), "8080"),
		MongoURI:      getenv("MONGO_URI"),
		MongoDatabase: withDefault(getenv("MONGO_DB"), "godo"),
		JWTSecret:     getenv("JWT_SECRET"),
	}

	var errs []error
	if cfg.MongoURI == "" {
		errs = append(errs, errors.New("MONGO_URI is required"))
	}
	if len(cfg.JWTSecret) < minSecretLength {
		errs = append(errs, fmt.Errorf("JWT_SECRET must be at least %d characters", minSecretLength))
	}

	ttl, err := time.ParseDuration(withDefault(getenv("TOKEN_TTL"), "168h"))
	if err != nil || ttl <= 0 {
		errs = append(errs, errors.New("TOKEN_TTL must be a positive duration such as 168h"))
	}
	cfg.TokenTTL = ttl

	for _, origin := range strings.Split(withDefault(getenv("CORS_ORIGINS"), "http://localhost:5173"), ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, origin)
		}
	}

	return cfg, errors.Join(errs...)
}

// withDefault returns value, or fallback when value is empty.
func withDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
