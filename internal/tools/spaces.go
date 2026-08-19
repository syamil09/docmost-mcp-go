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
				Limit:  req.GetInt("limit", 20),
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
