package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/syamil09/docmost-mcp-go/internal/config"
)

func newTestServer(t *testing.T) (*httptest.Server, *Client) {
	t.Helper()
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	cfg := config.Config{
		URL:             ts.URL,
		Email:           "u@e.com",
		Password:        "hunter22",
		Timeout:         time.Duration(5) * time.Second,
		MaxRetries:      0,
		InsecureSkipTLS: false,
	}
	c, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return ts, c
}

func TestNew_Success(t *testing.T) {
	_, c := newTestServer(t)
	if c == nil {
		t.Fatal("expected client")
	}
}

func TestClient_Login_OnFirstCall(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "authToken", Value: "jwt-abc", Path: "/", HttpOnly: true})
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer jwt-abc" {
			t.Errorf("missing/wrong bearer: %q", got)
		}
		w.Write([]byte(`{"data":{"id":"u1","name":"x","email":"u@e.com"}}`))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	c, _ := New(config.Config{URL: ts.URL, Email: "u@e.com", Password: "hunter22", Timeout: 5 * time.Second})
	var u User
	if err := c.Do(context.Background(), "/api/users/me", nil, &u); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if u.ID != "u1" {
		t.Errorf("got user %+v", u)
	}
}

func TestClient_ReLogin_On401(t *testing.T) {
	var loginCount int
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		loginCount++
		http.SetCookie(w, &http.Cookie{Name: "authToken", Value: "jwt-fresh", Path: "/"})
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer jwt-stale" {
			w.WriteHeader(401)
			w.Write([]byte(`{"message":"unauthorized"}`))
			return
		}
		w.Write([]byte(`{"data":{"id":"u1","name":"x","email":"u@e.com"}}`))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	c, _ := New(config.Config{URL: ts.URL, Email: "u@e.com", Password: "hunter22", Timeout: 5 * time.Second})
	c.token = "jwt-stale"
	var u User
	if err := c.Do(context.Background(), "/api/users/me", nil, &u); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if loginCount != 1 {
		t.Errorf("expected 1 re-login, got %d", loginCount)
	}
}

func TestPages_RoundTrip_CreateGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "authToken", Value: "t", Path: "/"})
		w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/api/pages/create", func(w http.ResponseWriter, r *http.Request) {
		var in map[string]any
		json.NewDecoder(r.Body).Decode(&in)
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"id": "p1", "spaceId": in["spaceId"], "title": in["title"], "content": nil, "position": "a",
			},
		})
		w.Write(body)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	c, _ := New(config.Config{URL: ts.URL, Email: "a@b.c", Password: "hunter22", Timeout: 5 * time.Second})
	p, err := c.CreatePage(context.Background(), CreatePageInput{SpaceID: "s1", Title: strPtr("T")})
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "p1" {
		t.Errorf("got %+v", p)
	}
}

func strPtr(s string) *string { return &s }
