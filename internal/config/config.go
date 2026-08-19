package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	URL             string
	Email           string
	Password        string
	APIKey          string
	WorkspaceID     string
	Timeout         time.Duration
	MaxRetries      int
	InsecureSkipTLS bool
}

func Load() (Config, error) {
	c := Config{
		URL:             os.Getenv("DOCMOST_URL"),
		Email:           os.Getenv("DOCMOST_EMAIL"),
		Password:        os.Getenv("DOCMOST_PASSWORD"),
		APIKey:          os.Getenv("DOCMOST_API_KEY"),
		WorkspaceID:     os.Getenv("DOCMOST_WORKSPACE_ID"),
		Timeout:         30 * time.Second,
		MaxRetries:      2,
		InsecureSkipTLS: os.Getenv("DOCMOST_INSECURE_SKIP_TLS") == "true",
	}
	if t := os.Getenv("DOCMOST_TIMEOUT"); t != "" {
		if d, err := time.ParseDuration(t); err == nil {
			c.Timeout = d
		}
	}
	if r := os.Getenv("DOCMOST_MAX_RETRIES"); r != "" {
		if n, err := strconv.Atoi(r); err == nil {
			c.MaxRetries = n
		}
	}
	if c.URL == "" {
		return c, errors.New("DOCMOST_URL is required")
	}
	if c.APIKey == "" && (c.Email == "" || c.Password == "") {
		return c, errors.New("either DOCMOST_API_KEY or DOCMOST_EMAIL + DOCMOST_PASSWORD must be set")
	}
	if !isHTTPSorHTTP(c.URL) {
		return c, fmt.Errorf("DOCMOST_URL must start with http:// or https://, got %q", c.URL)
	}
	return c, nil
}

func isHTTPSorHTTP(u string) bool {
	return len(u) > 7 && (u[:7] == "http://" || (len(u) > 8 && u[:8] == "https://"))
}
