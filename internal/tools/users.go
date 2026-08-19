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
