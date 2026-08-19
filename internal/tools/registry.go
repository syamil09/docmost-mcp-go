package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func RegisterAll(s *server.MCPServer, c *client.Client) {
	registerUserTools(s, c)
	registerSpaceTools(s, c)
}
