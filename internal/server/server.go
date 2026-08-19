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
