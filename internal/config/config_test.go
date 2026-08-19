package config

import (
	"testing"
)

func TestLoad_RequiresURL(t *testing.T) {
	t.Setenv("DOCMOST_URL", "")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when DOCMOST_URL missing")
	}
}

func TestLoad_RequiresCredentials(t *testing.T) {
	t.Setenv("DOCMOST_URL", "http://x")
	t.Setenv("DOCMOST_EMAIL", "")
	t.Setenv("DOCMOST_PASSWORD", "")
	t.Setenv("DOCMOST_API_KEY", "")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when no credentials configured")
	}
}

func TestLoad_OK_WithPassword(t *testing.T) {
	t.Setenv("DOCMOST_URL", "http://x")
	t.Setenv("DOCMOST_EMAIL", "a@b.c")
	t.Setenv("DOCMOST_PASSWORD", "hunter22")
	t.Setenv("DOCMOST_API_KEY", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.Email != "a@b.c" {
		t.Errorf("got email %q", cfg.Email)
	}
}

func TestLoad_OK_WithAPIKey(t *testing.T) {
	t.Setenv("DOCMOST_URL", "http://x")
	t.Setenv("DOCMOST_EMAIL", "")
	t.Setenv("DOCMOST_PASSWORD", "")
	t.Setenv("DOCMOST_API_KEY", "docmost_xxx")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.APIKey != "docmost_xxx" {
		t.Errorf("got key %q", cfg.APIKey)
	}
}

