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
}
