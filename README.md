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

Download the binary for your platform from [Releases](https://github.com/syamil09/docmost-mcp-go/releases):

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
