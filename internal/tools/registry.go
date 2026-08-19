package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func RegisterAll(s *server.MCPServer, c *client.Client) {
	registerPageTools(s, c)
	registerSpaceTools(s, c)
	registerCommentTools(s, c)
	registerAttachmentTools(s, c)
	registerWorkspaceTools(s, c)
	registerUserTools(s, c)
}
