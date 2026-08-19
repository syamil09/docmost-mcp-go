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
