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
