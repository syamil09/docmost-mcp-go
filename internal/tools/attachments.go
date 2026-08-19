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
