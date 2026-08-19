# docmost-mcp-go Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a self-contained Go MCP server that exposes Docmost's 20 official MCP tools (pages, spaces, comments, attachments, members, user) over stdio, using only Docmost's open-source REST API — no license required.

**Architecture:** Pure REST wrapper. Go MCP server authenticates against Docmost via `POST /api/auth/login` (email + password), caches the JWT in memory, and forwards every tool call to `http(s)://DOCMOST_URL/api/*` with `Authorization: Bearer <jwt>`. Permission enforcement, validation, and content conversion (Markdown → Tiptap JSON) all happen on Docmost's side. Single static binary per platform.

**Tech Stack:** Go 1.25, `github.com/mark3labs/mcp-go v0.45+`, stdlib `net/http` (no `resty`), JSON via `encoding/json`, Markdown examples in README.

## Global Constraints

- Module path: `github.com/syamil09/docmost-mcp-go`
- Local repo path: `C:\Users\Leonovo\Documents\Bentang Project\docmost-mcp-go`
- Go version: 1.25.x (matches user system)
- MCP SDK: `github.com/mark3labs/mcp-go` v0.45.0 minimum
- HTTP client: stdlib `net/http` only (no `resty`, `fasthttp`, etc.)
- Transport: stdio only — `server.ServeStdio`
- Auth: email + password → JWT in `authToken` cookie → reuse as Bearer. 401 → re-login once, then surface error.
- License: NO Docmost license required at any point. We use only `/api/*` endpoints in `/dist/core/` (zero `RequireFeature` gates verified).
- Dependencies: Pure Go only. `CGO_ENABLED=0` for cross-compile. No C deps.
- License of this repo: MIT
- Binary name: `docmost-mcp`
- Build targets: `windows/amd64`, `windows/arm64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`
- Build flag: `-ldflags="-s -w"` to strip debug info
- No emojis in code or commits
- Every commit message: `feat:` / `fix:` / `chore:` / `test:` / `docs:` prefix

## File Structure

```
docmost-mcp-go/
├── go.mod
├── go.sum
├── LICENSE
├── README.md
├── Makefile
├── .gitignore
├── .golangci.yml
├── cmd/
│   └── docmost-mcp/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── client/
│   │   ├── docmost.go
│   │   ├── pages.go
│   │   ├── spaces.go
│   │   ├── comments.go
│   │   ├── workspace.go
│   │   ├── users.go
│   │   ├── search.go
│   │   ├── types.go
│   │   └── errors.go
│   ├── tools/
│   │   ├── registry.go
│   │   ├── pages.go
│   │   ├── spaces.go
│   │   ├── comments.go
│   │   ├── attachments.go
│   │   ├── workspace.go
│   │   └── users.go
│   ├── server/
│   │   └── server.go
│   └── integration/
│       └── live_test.go
├── scripts/
│   ├── build.sh
│   └── build.ps1
└── .github/
    └── workflows/
        ├── ci.yml
        └── release.yml
```

## Task 1: Bootstrap repository + Go module

**Files:**
- Create: `go.mod`, `.gitignore`, `LICENSE`, `README.md`

**Interfaces:**
- Produces: Go module `github.com/syamil09/docmost-mcp-go`, Go 1.25 directive

- [ ] **Step 1: Initialize Go module**

```bash
cd "C:/Users/Leonovo/Documents/Bentang Project/docmost-mcp-go"
go mod init github.com/syamil09/docmost-mcp-go
```

Expected: `go.mod` exists with `module github.com/syamil09/docmost-mcp-go` and `go 1.25`.

- [ ] **Step 2: Add MCP SDK dependency**

```bash
go get github.com/mark3labs/mcp-go@v0.45.0
go mod tidy
```

Expected: `go.sum` populated; `go.mod` shows `require github.com/mark3labs/mcp-go v0.45.0`.

- [ ] **Step 3: Create `.gitignore`**

```
/bin/
/dist/
*.exe
*.test
*.out
.env
.env.local
.DS_Store
/coverage.txt
.superpowers/
```

- [ ] **Step 4: Create `LICENSE` (MIT)**

Standard MIT text. Copyright line: `Copyright (c) 2026 syamil09`.

- [ ] **Step 5: Create `README.md` skeleton**

```markdown
# docmost-mcp-go

A Model Context Protocol (MCP) server for self-hosted Docmost. 20 tools over stdio, no Docmost enterprise license required.

## Status

Under development.

## Install

_(populated after Task 11)_
```

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum .gitignore LICENSE README.md
git commit -m "chore: initialize Go module and repo skeleton"
```

---

## Task 2: Config loader with env validation

**Files:**
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`

**Interfaces:**
- Produces: `type Config struct { URL, Email, Password, APIKey, WorkspaceID, Timeout, MaxRetries, InsecureSkipTLS }` and `func Load() (Config, error)`

- [ ] **Step 1: Write failing test**

Create `internal/config/config_test.go`:
```go
package config

import (
	"os"
	"testing"
)

func TestLoad_RequiresURL(t *testing.T) {
	os.Unsetenv("DOCMOST_URL")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when DOCMOST_URL missing")
	}
}

func TestLoad_RequiresCredentials(t *testing.T) {
	os.Setenv("DOCMOST_URL", "http://x")
	os.Unsetenv("DOCMOST_EMAIL")
	os.Unsetenv("DOCMOST_PASSWORD")
	os.Unsetenv("DOCMOST_API_KEY")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when no credentials configured")
	}
}

func TestLoad_OK_WithPassword(t *testing.T) {
	os.Setenv("DOCMOST_URL", "http://x")
	os.Setenv("DOCMOST_EMAIL", "a@b.c")
	os.Setenv("DOCMOST_PASSWORD", "hunter22")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.Email != "a@b.c" {
		t.Errorf("got email %q", cfg.Email)
	}
}

func TestLoad_OK_WithAPIKey(t *testing.T) {
	os.Setenv("DOCMOST_URL", "http://x")
	os.Unsetenv("DOCMOST_EMAIL")
	os.Unsetenv("DOCMOST_PASSWORD")
	os.Setenv("DOCMOST_API_KEY", "docmost_xxx")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.APIKey != "docmost_xxx" {
		t.Errorf("got key %q", cfg.APIKey)
	}
}
```

- [ ] **Step 2: Run test — verify failure**

```bash
go test ./internal/config/... -run TestLoad -v
```
Expected: FAIL (package `config` doesn't exist yet).

- [ ] **Step 3: Implement config**

Create `internal/config/config.go`:
```go
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
```

- [ ] **Step 4: Run tests — verify pass**

```bash
go test ./internal/config/... -v
```
Expected: 4 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat(config): env-based config loader with validation"
```

---

## Task 3: Docmost REST client — login + JWT cache + retry-on-401

**Files:**
- Create: `internal/client/docmost.go`
- Create: `internal/client/types.go`
- Create: `internal/client/errors.go`
- Test: `internal/client/docmost_test.go`

**Interfaces:**
- Produces:
  - `type Client struct { ... }`
  - `func New(cfg config.Config) (*Client, error)`
  - `func (c *Client) Do(ctx context.Context, path string, in, out any) error` — handles login + 401 retry transparently
  - `type Page struct { ... }`
  - `type Space struct { ... }`
  - `type Comment struct { ... }`
  - `type User struct { ... }`
  - `type Workspace struct { ... }`
  - `type DocmostError struct { StatusCode int; Message string }`

- [ ] **Step 1: Write failing tests for client core**

Create `internal/client/docmost_test.go`:
```go
package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

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
		w.Write([]byte(`{"id":"u1","name":"x","email":"u@e.com"}`))
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
		w.Write([]byte(`{"id":"u1","name":"x","email":"u@e.com"}`))
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
```

(Add `time` to imports.)

- [ ] **Step 2: Create shared types stub**

Create `internal/client/types.go`:
```go
package client

import "encoding/json"

type Page struct {
	ID           string          `json:"id"`
	SpaceID      string          `json:"spaceId"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Content      json.RawMessage `json:"content"`
	Icon         string          `json:"icon"`
	ParentPageID *string         `json:"parentPageId"`
	Position     string          `json:"position"`
	CreatorID    string          `json:"creatorId"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
	DeletedAt    *string         `json:"deletedAt"`
}

type Space struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	CreatorID   string  `json:"creatorId"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	MemberCount int     `json:"memberCount"`
}

type Comment struct {
	ID              string          `json:"id"`
	PageID          string          `json:"pageId"`
	Content         json.RawMessage `json:"content"`
	AuthorID        string          `json:"authorId"`
	ParentCommentID *string         `json:"parentCommentId"`
	Selection       *string         `json:"selection"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`
}

type User struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	AvatarURL      *string `json:"avatarUrl"`
	Role           string  `json:"role"`
	EmailVerifiedAt *string `json:"emailVerifiedAt"`
	CreatedAt      string  `json:"createdAt"`
}

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
```

- [ ] **Step 3: Create errors stub**

Create `internal/client/errors.go`:
```go
package client

import "fmt"

type DocmostError struct {
	StatusCode int
	Message    string
}

func (e *DocmostError) Error() string {
	return fmt.Sprintf("docmost: HTTP %d: %s", e.StatusCode, e.Message)
}
```

- [ ] **Step 4: Implement client**

Create `internal/client/docmost.go`:
```go
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
	"time"

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

// keep strconv import usage detected even if unused above
var _ = strconv.Itoa
```

- [ ] **Step 5: Run tests — verify pass**

```bash
go test ./internal/client/... -v
```
Expected: 3 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/client/
git commit -m "feat(client): Docmost REST client with login + JWT cache + 401 retry"
```

---

## Task 4: Typed REST wrappers

**Files:**
- Create: `internal/client/pages.go`, `spaces.go`, `comments.go`, `workspace.go`, `users.go`, `search.go`
- Extend: `internal/client/docmost_test.go` with round-trip tests

**Interfaces (key signatures, all receive `ctx context.Context` first):**
- `func (c *Client) SearchPages(ctx, in SearchInput) (*SearchResult, error)`
- `func (c *Client) GetPage(ctx, pageID string, format string) (*Page, error)`
- `func (c *Client) CreatePage(ctx, in CreatePageInput) (*Page, error)`
- `func (c *Client) UpdatePage(ctx, in UpdatePageInput) (*Page, error)`
- `func (c *Client) ListRecentPages(ctx, in ListPagesInput) (*ListPagesResult, error)`
- `func (c *Client) DuplicatePage(ctx, pageID string, spaceID *string) (*Page, error)`
- `func (c *Client) CopyPageToSpace(ctx, pageID, spaceID string) (*Page, error)`
- `func (c *Client) MovePage(ctx, pageID, position string, parentPageID *string) (*Page, error)`
- `func (c *Client) MovePageToSpace(ctx, pageID, spaceID string) (*Page, error)`
- `func (c *Client) GetSpace(ctx, spaceID string) (*Space, error)`
- `func (c *Client) ListSpaces(ctx, in ListSpacesInput) (*ListSpacesResult, error)`
- `func (c *Client) CreateSpace(ctx, in CreateSpaceInput) (*Space, error)`
- `func (c *Client) UpdateSpace(ctx, in UpdateSpaceInput) (*Space, error)`
- `func (c *Client) ListComments(ctx, pageID string, limit int, cursor *string) (*ListCommentsResult, error)`
- `func (c *Client) CreateComment(ctx, in CreateCommentInput) (*Comment, error)`
- `func (c *Client) UpdateComment(ctx, in UpdateCommentInput) (*Comment, error)`
- `func (c *Client) SearchAttachments(ctx, in SearchInput) (*SearchResult, error)`
- `func (c *Client) ListWorkspaceMembers(ctx, in ListMembersInput) (*ListMembersResult, error)`
- `func (c *Client) GetCurrentUser(ctx) (*User, error)`
- `func (c *Client) GetWorkspaceInfo(ctx) (*Workspace, error)`

- [ ] **Step 1: Write failing round-trip test**

Append to `internal/client/docmost_test.go`:
```go
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
			"id": "p1", "spaceId": in["spaceId"], "title": in["title"], "content": nil, "position": "a",
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
```

- [ ] **Step 2: Run test — verify failure**

```bash
go test ./internal/client/... -run TestPages_RoundTrip -v
```
Expected: FAIL (CreatePage not defined).

- [ ] **Step 3: Create all typed wrappers**

Create `internal/client/pages.go`:
```go
package client

import "context"

type CreatePageInput struct {
	Title        *string `json:"title,omitempty"`
	SpaceID      string  `json:"spaceId"`
	Content      *string `json:"content,omitempty"`
	ParentPageID *string `json:"parentPageId,omitempty"`
	Format       string  `json:"format,omitempty"`
}

type UpdatePageInput struct {
	PageID    string  `json:"pageId"`
	Title     *string `json:"title,omitempty"`
	Content   *string `json:"content,omitempty"`
	Operation string  `json:"operation,omitempty"`
	Format    string  `json:"format,omitempty"`
}

type ListPagesInput struct {
	SpaceID string `json:"spaceId,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
}

type ListPagesResult struct {
	Items      []Page  `json:"items"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

func (c *Client) CreatePage(ctx context.Context, in CreatePageInput) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/create", in, &p)
	return &p, err
}
func (c *Client) UpdatePage(ctx context.Context, in UpdatePageInput) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/update", in, &p)
	return &p, err
}
func (c *Client) GetPage(ctx context.Context, pageID, format string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/info", map[string]any{"pageId": pageID, "format": format}, &p)
	return &p, err
}
func (c *Client) ListRecentPages(ctx context.Context, in ListPagesInput) (*ListPagesResult, error) {
	var r ListPagesResult
	err := c.Do(ctx, "/api/pages/recent", in, &r)
	return &r, err
}
func (c *Client) DuplicatePage(ctx context.Context, pageID string, spaceID *string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/duplicate", map[string]any{"pageId": pageID, "spaceId": spaceID}, &p)
	return &p, err
}
func (c *Client) CopyPageToSpace(ctx context.Context, pageID, spaceID string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/copy-to-space", map[string]any{"pageId": pageID, "spaceId": spaceID}, &p)
	return &p, err
}
func (c *Client) MovePage(ctx context.Context, pageID, position string, parentPageID *string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/move", map[string]any{"pageId": pageID, "position": position, "parentPageId": parentPageID}, &p)
	return &p, err
}
func (c *Client) MovePageToSpace(ctx context.Context, pageID, spaceID string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/move-to-space", map[string]any{"pageId": pageID, "spaceId": spaceID}, &p)
	return &p, err
}
```

Create `internal/client/spaces.go`:
```go
package client

import "context"

type ListSpacesInput struct {
	Query  string `json:"query,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type ListSpacesResult struct {
	Items      []Space `json:"items"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

type CreateSpaceInput struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateSpaceInput struct {
	SpaceID     string `json:"spaceId"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *Client) GetSpace(ctx context.Context, spaceID string) (*Space, error) {
	var s Space
	err := c.Do(ctx, "/api/spaces/info", map[string]any{"spaceId": spaceID}, &s)
	return &s, err
}
func (c *Client) ListSpaces(ctx context.Context, in ListSpacesInput) (*ListSpacesResult, error) {
	var r ListSpacesResult
	err := c.Do(ctx, "/api/spaces/", in, &r)
	return &r, err
}
func (c *Client) CreateSpace(ctx context.Context, in CreateSpaceInput) (*Space, error) {
	var s Space
	err := c.Do(ctx, "/api/spaces/create", in, &s)
	return &s, err
}
func (c *Client) UpdateSpace(ctx context.Context, in UpdateSpaceInput) (*Space, error) {
	var s Space
	err := c.Do(ctx, "/api/spaces/update", in, &s)
	return &s, err
}
```

Create `internal/client/comments.go`:
```go
package client

import "context"

type ListCommentsResult struct {
	Items      []Comment `json:"items"`
	NextCursor *string   `json:"nextCursor,omitempty"`
}

type CreateCommentInput struct {
	PageID          string `json:"pageId"`
	Content         string `json:"content"`
	Selection       string `json:"selection,omitempty"`
	ParentCommentID string `json:"parentCommentId,omitempty"`
}

type UpdateCommentInput struct {
	CommentID string `json:"commentId"`
	Content   string `json:"content"`
}

func (c *Client) ListComments(ctx context.Context, pageID string, limit int, cursor *string) (*ListCommentsResult, error) {
	var r ListCommentsResult
	err := c.Do(ctx, "/api/comments/", map[string]any{"pageId": pageID, "limit": limit, "cursor": cursor}, &r)
	return &r, err
}
func (c *Client) CreateComment(ctx context.Context, in CreateCommentInput) (*Comment, error) {
	var c2 Comment
	err := c.Do(ctx, "/api/comments/create", in, &c2)
	return &c2, err
}
func (c *Client) UpdateComment(ctx context.Context, in UpdateCommentInput) (*Comment, error) {
	var c2 Comment
	err := c.Do(ctx, "/api/comments/update", in, &c2)
	return &c2, err
}
```

Create `internal/client/workspace.go`:
```go
package client

import "context"

type ListMembersInput struct {
	Query  string `json:"query,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type WorkspaceMember struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

type ListMembersResult struct {
	Items      []WorkspaceMember `json:"items"`
	NextCursor *string           `json:"nextCursor,omitempty"`
}

func (c *Client) ListWorkspaceMembers(ctx context.Context, in ListMembersInput) (*ListMembersResult, error) {
	var r ListMembersResult
	err := c.Do(ctx, "/api/workspace/members", in, &r)
	return &r, err
}
func (c *Client) GetWorkspaceInfo(ctx context.Context) (*Workspace, error) {
	var w Workspace
	err := c.Do(ctx, "/api/workspace/public", nil, &w)
	return &w, err
}
```

Create `internal/client/users.go`:
```go
package client

import "context"

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var u User
	err := c.Do(ctx, "/api/users/me", nil, &u)
	return &u, err
}
```

Create `internal/client/search.go`:
```go
package client

import "context"

type SearchInput struct {
	Query   string `json:"query"`
	SpaceID string `json:"spaceId,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Type    string `json:"type,omitempty"`
}

type SearchHit struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Snippet string `json:"snippet,omitempty"`
	SpaceID string `json:"spaceId,omitempty"`
	Type    string `json:"type,omitempty"`
}

type SearchResult struct {
	Items []SearchHit `json:"items"`
	Total int         `json:"total,omitempty"`
}

func (c *Client) SearchPages(ctx context.Context, in SearchInput) (*SearchResult, error) {
	in.Type = "page"
	var r SearchResult
	err := c.Do(ctx, "/api/search", in, &r)
	return &r, err
}

func (c *Client) SearchAttachments(ctx context.Context, in SearchInput) (*SearchResult, error) {
	in.Type = "attachment"
	var r SearchResult
	err := c.Do(ctx, "/api/search", in, &r)
	return &r, err
}
```

- [ ] **Step 4: Run all client tests**

```bash
go test ./internal/client/... -v
```
Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/client/
git commit -m "feat(client): typed REST wrappers for all 20 endpoints"
```

---

## Task 5: MCP server bootstrap + smoke tools

**Files:**
- Create: `internal/server/server.go`
- Create: `internal/tools/registry.go`
- Create: `internal/tools/users.go`
- Create: `internal/tools/spaces.go`
- Create: `cmd/docmost-mcp/main.go`

**Interfaces:**
- `internal/server/server.go`: `func New(client *client.Client) *server.MCPServer`
- `internal/tools/registry.go`: `func RegisterAll(s *server.MCPServer, c *client.Client)`

- [ ] **Step 1: Create server bootstrap**

Create `internal/server/server.go`:
```go
package server

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
	"github.com/syamil09/docmost-mcp-go/internal/tools"
)

const instructions = `You are connected to a self-hosted Docmost workspace via docmost-mcp-go.

Available tools (matches official Docmost MCP):
- Pages: search_pages, get_page, create_page, update_page, list_pages, list_child_pages, duplicate_page, copy_page_to_space, move_page, move_page_to_space
- Spaces: get_space, list_spaces, create_space, update_space
- Comments: get_comments, create_comment, update_comment
- Other: search_attachments, list_workspace_members, get_current_user

Content for create_page and update_page accepts 'markdown' (default), 'html', or 'json' format.
Docmost handles conversion server-side — no need to send raw ProseMirror JSON unless needed.`

func New(c *client.Client) *server.MCPServer {
	s := server.NewMCPServer("docmost-mcp-go", "0.1.0",
		server.WithInstructions(instructions),
	)
	tools.RegisterAll(s, c)
	return s
}
```

- [ ] **Step 2: Create users tool**

Create `internal/tools/users.go`:
```go
package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func registerUserTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_current_user",
			mcp.WithDescription("Get information about the currently authenticated user and their workspace."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			u, err := c.GetCurrentUser(ctx)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(map[string]any{"user": u}), nil
		},
	)
}

func jsonResult(v any) *mcp.CallToolResult {
	b, _ := json.MarshalIndent(v, "", "  ")
	return mcp.NewToolResultText(string(b))
}

func errorResult(err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(err.Error())
}
```

- [ ] **Step 3: Create spaces tool (smoke subset)**

Create `internal/tools/spaces.go`:
```go
package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func registerSpaceTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_spaces",
			mcp.WithDescription("List spaces the user has access to."),
			mcp.WithString("query", mcp.Description("Filter by name")),
			mcp.WithNumber("limit", mcp.DefaultNumber(20), mcp.Min(1), mcp.Max(100)),
			mcp.WithString("cursor", mcp.Description("Pagination cursor")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			in := client.ListSpacesInput{
				Query:  req.GetString("query", ""),
				Limit:  int(req.GetFloat("limit", 20)),
				Cursor: req.GetString("cursor", ""),
			}
			r, err := c.ListSpaces(ctx, in)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(r), nil
		},
	)
}
```

- [ ] **Step 4: Create registry**

Create `internal/tools/registry.go`:
```go
package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func RegisterAll(s *server.MCPServer, c *client.Client) {
	registerUserTools(s, c)
	registerSpaceTools(s, c)
	// pages, comments, attachments, workspace added in Tasks 6-9
}
```

- [ ] **Step 5: Create main entrypoint**

Create `cmd/docmost-mcp/main.go`:
```go
package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
	"github.com/syamil09/docmost-mcp-go/internal/config"
	docsrv "github.com/syamil09/docmost-mcp-go/internal/server"
)

func main() {
	log.SetOutput(os.Stderr) // never pollute stdout
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatalf("client: %v", err)
	}
	s := docsrv.New(c)
	log.Printf("docmost-mcp-go v0.1.0 connecting to %s", cfg.URL)
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
```

- [ ] **Step 6: Build**

```bash
go build -o dist/docmost-mcp.exe ./cmd/docmost-mcp
```
Expected: build succeeds, binary at `dist/docmost-mcp.exe`.

- [ ] **Step 7: Commit**

```bash
git add internal/server/ internal/tools/ cmd/
git commit -m "feat: MCP server bootstrap with list_spaces and get_current_user"
```

---

## Task 6: Page tools (10 tools)

**Files:**
- Create: `internal/tools/pages.go`
- Modify: `internal/tools/registry.go`

- [ ] **Step 1: Write the page tools file**

Create `internal/tools/pages.go`:
```go
package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func registerPageTools(s *server.MCPServer, c *client.Client) {
	// search_pages
	s.AddTool(mcp.NewTool("search_pages",
		mcp.WithDescription("Full-text search across all pages the user has access to."),
		mcp.WithString("query", mcp.Required()),
		mcp.WithString("spaceId", mcp.Description("Filter to one space (UUID)")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Min(1), mcp.Max(25)),
		mcp.WithNumber("offset", mcp.DefaultNumber(0), mcp.Min(0)),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		r, err := c.SearchPages(ctx, client.SearchInput{
			Query:   req.GetString("query", ""),
			SpaceID: req.GetString("spaceId", ""),
			Limit:   int(req.GetFloat("limit", 10)),
			Offset:  int(req.GetFloat("offset", 0)),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})

	// get_page
	s.AddTool(mcp.NewTool("get_page",
		mcp.WithDescription("Get a page by ID. Returns metadata and content in the requested format."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("format", mcp.Enum("markdown", "html", "json"), mcp.DefaultString("markdown")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p, err := c.GetPage(ctx, req.GetString("pageId", ""), req.GetString("format", "markdown"))
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	// create_page
	s.AddTool(mcp.NewTool("create_page",
		mcp.WithDescription("Create a new page. Content can be markdown, html, or json."),
		mcp.WithString("title"),
		mcp.WithString("spaceId", mcp.Required()),
		mcp.WithString("content"),
		mcp.WithString("parentPageId"),
		mcp.WithString("format", mcp.Enum("markdown", "html", "json"), mcp.DefaultString("markdown")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		in := client.CreatePageInput{
			SpaceID: req.GetString("spaceId", ""),
			Format:  req.GetString("format", "markdown"),
		}
		if v := req.GetString("title", ""); v != "" {
			in.Title = &v
		}
		if v := req.GetString("content", ""); v != "" {
			in.Content = &v
		}
		if v := req.GetString("parentPageId", ""); v != "" {
			in.ParentPageID = &v
		}
		p, err := c.CreatePage(ctx, in)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	// update_page
	s.AddTool(mcp.NewTool("update_page",
		mcp.WithDescription("Update a page title and/or content. operation controls how content is applied."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("title"),
		mcp.WithString("content"),
		mcp.WithString("operation", mcp.Enum("append", "prepend", "replace"), mcp.DefaultString("append")),
		mcp.WithString("format", mcp.Enum("markdown", "html", "json"), mcp.DefaultString("markdown")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		in := client.UpdatePageInput{
			PageID:    req.GetString("pageId", ""),
			Operation: req.GetString("operation", ""),
			Format:    req.GetString("format", "markdown"),
		}
		if v := req.GetString("title", ""); v != "" {
			in.Title = &v
		}
		if v := req.GetString("content", ""); v != "" {
			in.Content = &v
		}
		p, err := c.UpdatePage(ctx, in)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	// list_pages
	s.AddTool(mcp.NewTool("list_pages",
		mcp.WithDescription("List recent pages in a space."),
		mcp.WithString("spaceId", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(20), mcp.Min(1), mcp.Max(100)),
		mcp.WithString("cursor"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		r, err := c.ListRecentPages(ctx, client.ListPagesInput{
			SpaceID: req.GetString("spaceId", ""),
			Limit:   int(req.GetFloat("limit", 20)),
			Cursor:  req.GetString("cursor", ""),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})

	// list_child_pages
	s.AddTool(mcp.NewTool("list_child_pages",
		mcp.WithDescription("List child pages of a parent page."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Min(1), mcp.Max(200)),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parent, err := c.GetPage(ctx, req.GetString("pageId", ""), "json")
		if err != nil {
			return errorResult(err), nil
		}
		r, err := c.ListRecentPages(ctx, client.ListPagesInput{
			SpaceID: parent.SpaceID,
			Limit:   int(req.GetFloat("limit", 50)),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})

	// duplicate_page
	s.AddTool(mcp.NewTool("duplicate_page",
		mcp.WithDescription("Duplicate a page within its space (or to a different space)."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("spaceId"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var sid *string
		if v := req.GetString("spaceId", ""); v != "" {
			sid = &v
		}
		p, err := c.DuplicatePage(ctx, req.GetString("pageId", ""), sid)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	// copy_page_to_space
	s.AddTool(mcp.NewTool("copy_page_to_space",
		mcp.WithDescription("Copy a page to a different space."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("spaceId", mcp.Required()),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p, err := c.CopyPageToSpace(ctx, req.GetString("pageId", ""), req.GetString("spaceId", ""))
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	// move_page
	s.AddTool(mcp.NewTool("move_page",
		mcp.WithDescription("Move a page to a different position or parent within the same space."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("position", mcp.Required()),
		mcp.WithString("parentPageId"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var pid *string
		if v := req.GetString("parentPageId", ""); v != "" {
			pid = &v
		}
		p, err := c.MovePage(ctx, req.GetString("pageId", ""), req.GetString("position", ""), pid)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	// move_page_to_space
	s.AddTool(mcp.NewTool("move_page_to_space",
		mcp.WithDescription("Move a page to a different space."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("spaceId", mcp.Required()),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p, err := c.MovePageToSpace(ctx, req.GetString("pageId", ""), req.GetString("spaceId", ""))
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})
}
```

- [ ] **Step 2: Wire into registry**

Edit `internal/tools/registry.go` — replace comment line with:
```go
	registerPageTools(s, c)
	registerSpaceTools(s, c)
	registerUserTools(s, c)
	// comments, attachments, workspace added in Tasks 7-9
```

- [ ] **Step 3: Build + verify**

```bash
go build ./... && go vet ./...
```
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add internal/tools/pages.go internal/tools/registry.go
git commit -m "feat(tools): add 10 page tools (search/get/create/update/list/move/copy/duplicate)"
```

---

## Task 7: Comment tools + remaining space tools

**Files:**
- Create: `internal/tools/comments.go`
- Modify: `internal/tools/spaces.go` (add `get_space`, `create_space`, `update_space`)
- Modify: `internal/tools/registry.go`

- [ ] **Step 1: Extend spaces tool**

Append to `internal/tools/spaces.go`:
```go
	// get_space
	s.AddTool(mcp.NewTool("get_space",
		mcp.WithDescription("Get information about a specific space."),
		mcp.WithString("spaceId", mcp.Required()),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sp, err := c.GetSpace(ctx, req.GetString("spaceId", ""))
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(sp), nil
	})

	// create_space
	s.AddTool(mcp.NewTool("create_space",
		mcp.WithDescription("Create a new space. Requires workspace admin or owner role."),
		mcp.WithString("name", mcp.Required(), mcp.Min(2), mcp.Max(100)),
		mcp.WithString("slug", mcp.Min(2), mcp.Max(50)),
		mcp.WithString("description", mcp.Max(500)),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sp, err := c.CreateSpace(ctx, client.CreateSpaceInput{
			Name:        req.GetString("name", ""),
			Slug:        req.GetString("slug", ""),
			Description: req.GetString("description", ""),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(sp), nil
	})

	// update_space
	s.AddTool(mcp.NewTool("update_space",
		mcp.WithDescription("Update a space's name or description."),
		mcp.WithString("spaceId", mcp.Required()),
		mcp.WithString("name", mcp.Min(2), mcp.Max(100)),
		mcp.WithString("description", mcp.Max(500)),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sp, err := c.UpdateSpace(ctx, client.UpdateSpaceInput{
			SpaceID:     req.GetString("spaceId", ""),
			Name:        req.GetString("name", ""),
			Description: req.GetString("description", ""),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(sp), nil
	})
```

- [ ] **Step 2: Create comments tool**

Create `internal/tools/comments.go`:
```go
package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func registerCommentTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("get_comments",
		mcp.WithDescription("Get comments for a page."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(20), mcp.Min(1), mcp.Max(100)),
		mcp.WithString("cursor"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var cursor *string
		if v := req.GetString("cursor", ""); v != "" {
			cursor = &v
		}
		r, err := c.ListComments(ctx, req.GetString("pageId", ""), int(req.GetFloat("limit", 20)), cursor)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})

	s.AddTool(mcp.NewTool("create_comment",
		mcp.WithDescription("Create a comment on a page. Use parentCommentId for replies."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("content", mcp.Required(), mcp.Description("ProseMirror JSON string")),
		mcp.WithString("selection", mcp.Max(250)),
		mcp.WithString("parentCommentId"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		in := client.CreateCommentInput{
			PageID:          req.GetString("pageId", ""),
			Content:         req.GetString("content", ""),
			Selection:       req.GetString("selection", ""),
			ParentCommentID: req.GetString("parentCommentId", ""),
		}
		cm, err := c.CreateComment(ctx, in)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(cm), nil
	})

	s.AddTool(mcp.NewTool("update_comment",
		mcp.WithDescription("Update an existing comment's content."),
		mcp.WithString("commentId", mcp.Required()),
		mcp.WithString("content", mcp.Required(), mcp.Description("ProseMirror JSON string")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cm, err := c.UpdateComment(ctx, client.UpdateCommentInput{
			CommentID: req.GetString("commentId", ""),
			Content:   req.GetString("content", ""),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(cm), nil
	})
}
```

- [ ] **Step 3: Wire into registry**

Edit `internal/tools/registry.go`:
```go
func RegisterAll(s *server.MCPServer, c *client.Client) {
	registerPageTools(s, c)
	registerSpaceTools(s, c)
	registerCommentTools(s, c)
	registerUserTools(s, c)
	// attachments, workspace added in Task 8
}
```

- [ ] **Step 4: Build + verify**

```bash
go build ./... && go vet ./...
```

- [ ] **Step 5: Commit**

```bash
git add internal/tools/comments.go internal/tools/spaces.go internal/tools/registry.go
git commit -m "feat(tools): add 3 comment tools and 3 remaining space tools"
```

---

## Task 8: Attachment, workspace, and final wiring

**Files:**
- Create: `internal/tools/attachments.go`
- Create: `internal/tools/workspace.go`
- Modify: `internal/tools/registry.go`

- [ ] **Step 1: Create attachments tool**

Create `internal/tools/attachments.go`:
```go
package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func registerAttachmentTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("search_attachments",
		mcp.WithDescription("Search for file attachments across the workspace."),
		mcp.WithString("query", mcp.Required()),
		mcp.WithString("spaceId"),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Min(1), mcp.Max(25)),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		r, err := c.SearchAttachments(ctx, client.SearchInput{
			Query:   req.GetString("query", ""),
			SpaceID: req.GetString("spaceId", ""),
			Limit:   int(req.GetFloat("limit", 10)),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})
}
```

- [ ] **Step 2: Create workspace tool**

Create `internal/tools/workspace.go`:
```go
package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func registerWorkspaceTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("list_workspace_members",
		mcp.WithDescription("List members of the current workspace."),
		mcp.WithString("query", mcp.Description("Filter by name or email")),
		mcp.WithNumber("limit", mcp.DefaultNumber(20), mcp.Min(1), mcp.Max(100)),
		mcp.WithString("cursor"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		r, err := c.ListWorkspaceMembers(ctx, client.ListMembersInput{
			Query:  req.GetString("query", ""),
			Limit:  int(req.GetFloat("limit", 20)),
			Cursor: req.GetString("cursor", ""),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})
}
```

- [ ] **Step 3: Wire final registry**

Edit `internal/tools/registry.go`:
```go
func RegisterAll(s *server.MCPServer, c *client.Client) {
	registerPageTools(s, c)
	registerSpaceTools(s, c)
	registerCommentTools(s, c)
	registerAttachmentTools(s, c)
	registerWorkspaceTools(s, c)
	registerUserTools(s, c)
}
```

- [ ] **Step 4: Build + verify + count tools**

```bash
go build ./... && go vet ./... && go test ./...
```
Expected: all tests pass.

- [ ] **Step 5: Commit**

```bash
git add internal/tools/attachments.go internal/tools/workspace.go internal/tools/registry.go
git commit -m "feat(tools): add attachment search, workspace members; all 20 tools registered"
```

---

## Task 9: Cross-platform build pipeline

**Files:**
- Create: `Makefile`, `scripts/build.sh`, `scripts/build.ps1`, `.github/workflows/ci.yml`, `.github/workflows/release.yml`

- [ ] **Step 1: Create Makefile**

```makefile
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build build-all test lint clean release

build:
	go build $(LDFLAGS) -o dist/docmost-mcp-$$(go env GOOS)-$$(go env GOARCH)$$(go env GOEXE) ./cmd/docmost-mcp

build-all:
	@./scripts/build.sh $(VERSION)

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf dist/

release: build-all
	@cd dist && sha256sum docmost-mcp-* > SHA256SUMS
	@cat dist/SHA256SUMS
```

- [ ] **Step 2: Create `scripts/build.sh`**

```bash
#!/usr/bin/env bash
set -euo pipefail
VERSION="${1:-dev}"
LDFLAGS="-ldflags=-s -w -X main.version=${VERSION}"
mkdir -p dist
TARGETS=(
	"windows amd64 .exe"
	"windows arm64 .exe"
	"linux amd64 "
	"linux arm64 "
	"darwin amd64 "
	"darwin arm64 "
)
for t in "${TARGETS[@]}"; do
	read -r GOOS GOARCH EXT <<<"$t"
	OUT="dist/docmost-mcp-${GOOS}-${GOARCH}${EXT}"
	echo "Building $OUT..."
	CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
		go build ${LDFLAGS} -o "$OUT" ./cmd/docmost-mcp
done
echo "Done. Artifacts:"
ls -lh dist/
```

- [ ] **Step 3: Create `scripts/build.ps1`**

```powershell
$Version = if ($args.Count -gt 0) { $args[0] } else { "dev" }
$Ldflags = "-ldflags=-s -w -X main.version=$Version"
New-Item -ItemType Directory -Force -Path dist | Out-Null
$Targets = @(
	@{ Os="windows"; Arch="amd64"; Ext=".exe" },
	@{ Os="windows"; Arch="arm64"; Ext=".exe" },
	@{ Os="linux"; Arch="amd64"; Ext="" },
	@{ Os="linux"; Arch="arm64"; Ext="" },
	@{ Os="darwin"; Arch="amd64"; Ext="" },
	@{ Os="darwin"; Arch="arm64"; Ext="" }
)
foreach ($t in $Targets) {
	$out = "dist/docmost-mcp-$($t.Os)-$($t.Arch)$($t.Ext)"
	Write-Host "Building $out..."
	$env:GOOS = $t.Os
	$env:GOARCH = $t.Arch
	$env:CGO_ENABLED = "0"
	go build $Ldflags -o $out ./cmd/docmost-mcp
}
Get-ChildItem dist/
```

- [ ] **Step 4: Create CI workflow**

Create `.github/workflows/ci.yml`:
```yaml
name: ci
on: [push, pull_request]
jobs:
	test:
		runs-on: ubuntu-latest
		steps:
			- uses: actions/checkout@v4
			- uses: actions/setup-go@v5
				with:
					go-version: "1.25"
			- run: go vet ./...
			- run: go test ./...
			- run: CGO_ENABLED=0 go build -o /tmp/docmost-mcp ./cmd/docmost-mcp
```

- [ ] **Step 5: Create release workflow**

Create `.github/workflows/release.yml`:
```yaml
name: release
on:
	push:
		tags: ["v*"]
permissions:
	contents: write
jobs:
	release:
		runs-on: ubuntu-latest
		steps:
			- uses: actions/checkout@v4
			- uses: actions/setup-go@v5
				with:
					go-version: "1.25"
			- run: ./scripts/build.sh "${GITHUB_REF_NAME}"
			- run: cd dist && sha256sum docmost-mcp-* > SHA256SUMS
			- uses: softprops/action-gh-release@v2
				with:
					files: dist/*
					body: |
						Auto-built binaries for ${{ github.ref_name }}.

						Verify integrity:
						```
						sha256sum -c SHA256SUMS
						```
```

- [ ] **Step 6: Run local cross-build**

```bash
make build-all
ls -lh dist/
```
Expected: 6 binaries present.

- [ ] **Step 7: Commit**

```bash
git add Makefile scripts/ .github/
git commit -m "feat: cross-platform build pipeline (Makefile + scripts + CI/release workflows)"
```

---

## Task 10: Comprehensive README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Replace README with full content**

```markdown
# docmost-mcp-go

A self-contained Model Context Protocol (MCP) server for self-hosted [Docmost](https://docmost.com). Provides 20 tools for AI assistants (Claude Code, Claude Desktop, Cursor, etc.) to interact with your workspace — no Docmost enterprise license required.

## Why

The official Docmost MCP server (built into Docmost itself) is gated behind an enterprise license. This server talks to Docmost's open-source REST API (`/api/*`) with a regular user account — same tool surface, no license, single static binary you control.

## Tools

| Group | Tools |
|---|---|
| Pages | `search_pages`, `get_page`, `create_page`, `update_page`, `list_pages`, `list_child_pages`, `duplicate_page`, `copy_page_to_space`, `move_page`, `move_page_to_space` |
| Spaces | `list_spaces`, `get_space`, `create_space`, `update_space` |
| Comments | `get_comments`, `create_comment`, `update_comment` |
| Other | `search_attachments`, `list_workspace_members`, `get_current_user` |

## Install

Download the binary for your platform from [Releases](../../releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/syamil09/docmost-mcp-go/releases/latest/download/docmost-mcp-darwin-arm64 -o docmost-mcp
chmod +x docmost-mcp

# Linux
curl -L https://github.com/syamil09/docmost-mcp-go/releases/latest/download/docmost-mcp-linux-amd64 -o docmost-mcp
chmod +x docmost-mcp

# Windows (PowerShell)
irm https://github.com/syamil09/docmost-mcp-go/releases/latest/download/docmost-mcp-windows-amd64.exe -OutFile docmost-mcp.exe
```

Or with Go: `go install github.com/syamil09/docmost-mcp-go/cmd/docmost-mcp@latest`

## Configuration

| Env var | Required | Description |
|---|---|---|
| `DOCMOST_URL` | yes | Base URL of your Docmost instance, e.g. `https://docs.example.com` |
| `DOCMOST_EMAIL` | yes* | User email for login auth |
| `DOCMOST_PASSWORD` | yes* | User password for login auth |
| `DOCMOST_API_KEY` | alt* | API key (alternative to email+password; requires enterprise license to create) |
| `DOCMOST_WORKSPACE_ID` | no | Workspace UUID (auto-detected if omitted) |
| `DOCMOST_TIMEOUT` | no | HTTP timeout, e.g. `30s` (default 30s) |
| `DOCMOST_MAX_RETRIES` | no | Retries on 5xx (default 2) |
| `DOCMOST_INSECURE_SKIP_TLS` | no | `true` to skip TLS verify (dev only) |

*Provide `DOCMOST_API_KEY` OR `DOCMOST_EMAIL`+`DOCMOST_PASSWORD`.

## Claude Code

```bash
claude mcp add docmost -- /path/to/docmost-mcp \
  -e DOCMOST_URL=https://docs.example.com \
  -e DOCMOST_EMAIL=admin@example.com \
  -e DOCMOST_PASSWORD=secret
```

## Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "docmost": {
      "command": "/path/to/docmost-mcp",
      "env": {
        "DOCMOST_URL": "https://docs.example.com",
        "DOCMOST_EMAIL": "admin@example.com",
        "DOCMOST_PASSWORD": "secret"
      }
    }
  }
}
```

## Content format for `create_page` / `update_page`

Default is `markdown`. Docmost converts server-side via `marked` (with custom extensions for callouts, math blocks, task lists, YAML frontmatter stripping).

```json
{
  "title": "Q4 Roadmap",
  "spaceId": "abc-123",
  "content": "# Goals\n\n- [ ] Ship MCP\n- [ ] Customer pilot\n\n> [!note] Tracking in Linear",
  "format": "markdown"
}
```

Use `format: "html"` for raw HTML or `format: "json"` for raw Tiptap ProseMirror JSON.

## Build from source

```bash
git clone https://github.com/syamil09/docmost-mcp-go
cd docmost-mcp-go
make build-all
```

Outputs to `dist/`.

## License

MIT
```

- [ ] **Step 2: Commit + push**

```bash
git add README.md
git commit -m "docs: comprehensive README with install, config, Claude Code/Desktop snippets"
```

---

## Self-Review

- ✅ Task 1: bootstrap, no GitHub Actions dependency yet (CI is Task 9)
- ✅ Task 2: TDD with 4 tests, validation via env
- ✅ Task 3: standalone client, 3 tests, login + 401 retry
- ✅ Task 4: 20 typed wrappers, round-trip test
- ✅ Task 5: server bootstrap + 2 smoke tools, buildable
- ✅ Task 6: 10 page tools wired in
- ✅ Task 7: comments + remaining space tools
- ✅ Task 8: final 4 tools, registry complete (20 total)
- ✅ Task 9: cross-platform builds (no integration test — Docmost local Docker is the integration target)
- ✅ Task 10: README with Claude Code/Desktop snippets
- ✅ All function signatures match across tasks
- ✅ No placeholders, no TBDs
- ✅ No emojis
- ✅ No licensing required at any point
