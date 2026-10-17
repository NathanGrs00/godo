package config

import (
	"strings"
	"testing"
	"time"
)

func lookup(env map[string]string) func(string) string {
	return func(key string) string { return env[key] }
}

func TestDefaults(t *testing.T) {
	cfg, err := FromLookup(lookup(map[string]string{
		"MONGO_URI":  "mongodb://db",
		"JWT_SECRET": strings.Repeat("s", 32),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "8080" || cfg.MongoDatabase != "godo" || cfg.TokenTTL != 168*time.Hour {
		t.Errorf("defaults = %+v", cfg)
	}
	if len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "http://localhost:5173" {
		t.Errorf("origins = %v", cfg.AllowedOrigins)
	}
}

func TestOverrides(t *testing.T) {
	cfg, err := FromLookup(lookup(map[string]string{
		"MONGO_URI":    "mongodb://db",
		"JWT_SECRET":   strings.Repeat("s", 40),
		"PORT":         "9000",
		"MONGO_DB":     "other",
		"TOKEN_TTL":    "1h",
		"CORS_ORIGINS": "http://a.test, http://b.test,",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "9000" || cfg.MongoDatabase != "other" || cfg.TokenTTL != time.Hour {
		t.Errorf("cfg = %+v", cfg)
	}
	if len(cfg.AllowedOrigins) != 2 || cfg.AllowedOrigins[1] != "http://b.test" {
		t.Errorf("origins = %v", cfg.AllowedOrigins)
	}
}

func TestErrors(t *testing.T) {
	_, err := FromLookup(lookup(map[string]string{"JWT_SECRET": "short", "TOKEN_TTL": "soon"}))
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{"MONGO_URI", "JWT_SECRET", "TOKEN_TTL"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should mention %s", err, want)
		}
	}
}
