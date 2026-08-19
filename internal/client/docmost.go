package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/syamil09/docmost-mcp-go/internal/config"
)

type Client struct {
	cfg   config.Config
	http  *http.Client
	base  *url.URL
	mu    sync.Mutex
	token string
}

func New(cfg config.Config) (*Client, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipTLS}, //nolint:gosec
	}
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout, Transport: tr},
		base: u,
	}, nil
}

func (c *Client) Do(ctx context.Context, path string, in, out any) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}
	if err := c.doOnce(ctx, path, in, out); err != nil {
		if de, ok := err.(*DocmostError); ok && de.StatusCode == http.StatusUnauthorized && c.cfg.APIKey == "" {
			c.mu.Lock()
			c.token = ""
			c.mu.Unlock()
			if err := c.ensureToken(ctx); err != nil {
				return err
			}
			return c.doOnce(ctx, path, in, out)
		}
		return err
	}
	return nil
}

func (c *Client) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	hasToken := c.token != ""
	c.mu.Unlock()
	if hasToken || c.cfg.APIKey != "" {
		return nil
	}
	return c.login(ctx)
}

func (c *Client) login(ctx context.Context) error {
	body, _ := json.Marshal(map[string]string{
		"email":    c.cfg.Email,
		"password": c.cfg.Password,
	})
	req, _ := http.NewRequestWithContext(ctx, "POST", c.base.String()+"/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return &DocmostError{StatusCode: resp.StatusCode, Message: string(b)}
	}
	for _, ck := range resp.Cookies() {
		if ck.Name == "authToken" {
			c.mu.Lock()
			c.token = ck.Value
			c.mu.Unlock()
			return nil
		}
	}
	return fmt.Errorf("login: no authToken cookie in response")
}

func (c *Client) doOnce(ctx context.Context, path string, in, out any) error {
	var rdr io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.base.String()+path, rdr)
	if err != nil {
		return err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	} else {
		c.mu.Lock()
		tk := c.token
		c.mu.Unlock()
		req.Header.Set("Authorization", "Bearer "+tk)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := extractMessage(b)
		if resp.StatusCode == http.StatusUnauthorized {
			msg = "unauthorized"
		}
		return &DocmostError{StatusCode: resp.StatusCode, Message: msg}
	}
	if out != nil && len(b) > 0 {
		if err := json.Unmarshal(b, out); err != nil {
			return fmt.Errorf("decode %s: %w", path, err)
		}
	}
	return nil
}

func extractMessage(b []byte) string {
	var p struct {
		Message    string `json:"message"`
		StatusCode int    `json:"statusCode"`
	}
	if json.Unmarshal(b, &p) == nil && p.Message != "" {
		return p.Message
	}
	return strings.TrimSpace(string(b))
}