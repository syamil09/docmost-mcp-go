package client

import "context"

type SearchInput struct {
	Query   string `json:"query"`
	SpaceID string `json:"spaceId,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Type    string `json:"type,omitempty"`
}

type SearchHit struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Snippet string `json:"snippet,omitempty"`
	SpaceID string `json:"spaceId,omitempty"`
	Type    string `json:"type,omitempty"`
}

type SearchResult struct {
	Items []SearchHit `json:"items"`
	Total int         `json:"total,omitempty"`
}

func (c *Client) SearchPages(ctx context.Context, in SearchInput) (*SearchResult, error) {
	in.Type = "page"
	var r SearchResult
	err := c.Do(ctx, "/api/search", in, &r)
	return &r, err
}

func (c *Client) SearchAttachments(ctx context.Context, in SearchInput) (*SearchResult, error) {
	in.Type = "attachment"
	var r SearchResult
	err := c.Do(ctx, "/api/search", in, &r)
	return &r, err
}
