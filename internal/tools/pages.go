package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
)

func registerPageTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("search_pages",
		mcp.WithDescription("Full-text search across all pages the user has access to."),
		mcp.WithString("query", mcp.Required()),
		mcp.WithString("spaceId", mcp.Description("Filter to one space (UUID)")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Min(1), mcp.Max(25)),
		mcp.WithNumber("offset", mcp.DefaultNumber(0), mcp.Min(0)),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		r, err := c.SearchPages(ctx, client.SearchInput{
			Query:   req.GetString("query", ""),
			SpaceID: req.GetString("spaceId", ""),
			Limit:   int(req.GetFloat("limit", 10)),
			Offset:  int(req.GetFloat("offset", 0)),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})

	s.AddTool(mcp.NewTool("get_page",
		mcp.WithDescription("Get a page by ID. Returns metadata and content in the requested format."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("format", mcp.Enum("markdown", "html", "json"), mcp.DefaultString("markdown")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p, err := c.GetPage(ctx, req.GetString("pageId", ""), req.GetString("format", "markdown"))
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	s.AddTool(mcp.NewTool("create_page",
		mcp.WithDescription("Create a new page. Content can be markdown, html, or json."),
		mcp.WithString("title"),
		mcp.WithString("spaceId", mcp.Required()),
		mcp.WithString("content"),
		mcp.WithString("parentPageId"),
		mcp.WithString("format", mcp.Enum("markdown", "html", "json"), mcp.DefaultString("markdown")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		in := client.CreatePageInput{
			SpaceID: req.GetString("spaceId", ""),
			Format:  req.GetString("format", "markdown"),
		}
		if v := req.GetString("title", ""); v != "" {
			in.Title = &v
		}
		if v := req.GetString("content", ""); v != "" {
			in.Content = &v
		}
		if v := req.GetString("parentPageId", ""); v != "" {
			in.ParentPageID = &v
		}
		p, err := c.CreatePage(ctx, in)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	s.AddTool(mcp.NewTool("update_page",
		mcp.WithDescription("Update a page title and/or content. operation controls how content is applied."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("title"),
		mcp.WithString("content"),
		mcp.WithString("operation", mcp.Enum("append", "prepend", "replace"), mcp.DefaultString("append")),
		mcp.WithString("format", mcp.Enum("markdown", "html", "json"), mcp.DefaultString("markdown")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		in := client.UpdatePageInput{
			PageID:    req.GetString("pageId", ""),
			Operation: req.GetString("operation", "append"),
			Format:    req.GetString("format", "markdown"),
		}
		if v := req.GetString("title", ""); v != "" {
			in.Title = &v
		}
		if v := req.GetString("content", ""); v != "" {
			in.Content = &v
		}
		p, err := c.UpdatePage(ctx, in)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	s.AddTool(mcp.NewTool("list_pages",
		mcp.WithDescription("List recent pages in a space."),
		mcp.WithString("spaceId", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(20), mcp.Min(1), mcp.Max(100)),
		mcp.WithString("cursor"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		r, err := c.ListRecentPages(ctx, client.ListPagesInput{
			SpaceID: req.GetString("spaceId", ""),
			Limit:   int(req.GetFloat("limit", 20)),
			Cursor:  req.GetString("cursor", ""),
		})
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(r), nil
	})

	s.AddTool(mcp.NewTool("list_child_pages",
		mcp.WithDescription("List child pages of a parent page."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Min(1), mcp.Max(200)),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := int(req.GetFloat("limit", 50))
		r, err := c.ListRecentPages(ctx, client.ListPagesInput{
			ParentPageID: req.GetString("pageId", ""),
			Limit:        limit,
		})
		if err != nil {
			return errorResult(err), nil
		}
		children := make([]client.Page, 0, len(r.Items))
		for _, p := range r.Items {
			if p.ParentPageID != nil && *p.ParentPageID == req.GetString("pageId", "") {
				children = append(children, p)
			}
		}
		r.Items = children
		r.NextCursor = nil
		return jsonResult(r), nil
	})

	s.AddTool(mcp.NewTool("duplicate_page",
		mcp.WithDescription("Duplicate a page within its space (or to a different space)."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("spaceId"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var sid *string
		if v := req.GetString("spaceId", ""); v != "" {
			sid = &v
		}
		p, err := c.DuplicatePage(ctx, req.GetString("pageId", ""), sid)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	s.AddTool(mcp.NewTool("copy_page_to_space",
		mcp.WithDescription("Copy a page to a different space."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("spaceId", mcp.Required()),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p, err := c.CopyPageToSpace(ctx, req.GetString("pageId", ""), req.GetString("spaceId", ""))
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	s.AddTool(mcp.NewTool("move_page",
		mcp.WithDescription("Move a page to a different position or parent within the same space."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("position", mcp.Required()),
		mcp.WithString("parentPageId"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var pid *string
		if v := req.GetString("parentPageId", ""); v != "" {
			pid = &v
		}
		p, err := c.MovePage(ctx, req.GetString("pageId", ""), req.GetString("position", ""), pid)
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})

	s.AddTool(mcp.NewTool("move_page_to_space",
		mcp.WithDescription("Move a page to a different space."),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("spaceId", mcp.Required()),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p, err := c.MovePageToSpace(ctx, req.GetString("pageId", ""), req.GetString("spaceId", ""))
		if err != nil {
			return errorResult(err), nil
		}
		return jsonResult(p), nil
	})
}
