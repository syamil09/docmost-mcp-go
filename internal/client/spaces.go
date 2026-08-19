package client

import "context"

type ListSpacesInput struct {
	Query  string `json:"query,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type ListSpacesResult struct {
	Items      []Space `json:"items"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

type CreateSpaceInput struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateSpaceInput struct {
	SpaceID     string `json:"spaceId"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *Client) GetSpace(ctx context.Context, spaceID string) (*Space, error) {
	var s Space
	err := c.Do(ctx, "/api/spaces/info", map[string]any{"spaceId": spaceID}, &s)
	return &s, err
}
func (c *Client) ListSpaces(ctx context.Context, in ListSpacesInput) (*ListSpacesResult, error) {
	var r ListSpacesResult
	err := c.Do(ctx, "/api/spaces/", in, &r)
	return &r, err
}
func (c *Client) CreateSpace(ctx context.Context, in CreateSpaceInput) (*Space, error) {
	var s Space
	err := c.Do(ctx, "/api/spaces/create", in, &s)
	return &s, err
}
func (c *Client) UpdateSpace(ctx context.Context, in UpdateSpaceInput) (*Space, error) {
	var s Space
	err := c.Do(ctx, "/api/spaces/update", in, &s)
	return &s, err
}
