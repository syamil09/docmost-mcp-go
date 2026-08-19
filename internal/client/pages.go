package client

import "context"

type CreatePageInput struct {
	Title        *string `json:"title,omitempty"`
	SpaceID      string  `json:"spaceId"`
	Content      *string `json:"content,omitempty"`
	ParentPageID *string `json:"parentPageId,omitempty"`
	Format       string  `json:"format,omitempty"`
}

type UpdatePageInput struct {
	PageID    string  `json:"pageId"`
	Title     *string `json:"title,omitempty"`
	Content   *string `json:"content,omitempty"`
	Operation string  `json:"operation,omitempty"`
	Format    string  `json:"format,omitempty"`
}

type ListPagesInput struct {
	SpaceID string `json:"spaceId,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
}

type ListPagesResult struct {
	Items      []Page  `json:"items"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

func (c *Client) CreatePage(ctx context.Context, in CreatePageInput) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/create", in, &p)
	return &p, err
}
func (c *Client) UpdatePage(ctx context.Context, in UpdatePageInput) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/update", in, &p)
	return &p, err
}
func (c *Client) GetPage(ctx context.Context, pageID, format string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/info", map[string]any{"pageId": pageID, "format": format}, &p)
	return &p, err
}
func (c *Client) ListRecentPages(ctx context.Context, in ListPagesInput) (*ListPagesResult, error) {
	var r ListPagesResult
	err := c.Do(ctx, "/api/pages/recent", in, &r)
	return &r, err
}
func (c *Client) DuplicatePage(ctx context.Context, pageID string, spaceID *string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/duplicate", map[string]any{"pageId": pageID, "spaceId": spaceID}, &p)
	return &p, err
}
func (c *Client) CopyPageToSpace(ctx context.Context, pageID, spaceID string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/copy-to-space", map[string]any{"pageId": pageID, "spaceId": spaceID}, &p)
	return &p, err
}
func (c *Client) MovePage(ctx context.Context, pageID, position string, parentPageID *string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/move", map[string]any{"pageId": pageID, "position": position, "parentPageId": parentPageID}, &p)
	return &p, err
}
func (c *Client) MovePageToSpace(ctx context.Context, pageID, spaceID string) (*Page, error) {
	var p Page
	err := c.Do(ctx, "/api/pages/move-to-space", map[string]any{"pageId": pageID, "spaceId": spaceID}, &p)
	return &p, err
}
