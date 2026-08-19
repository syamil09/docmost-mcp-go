package client

import "context"

type ListMembersInput struct {
	Query  string `json:"query,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type WorkspaceMember struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

type ListMembersResult struct {
	Items      []WorkspaceMember `json:"items"`
	NextCursor *string           `json:"nextCursor,omitempty"`
}

func (c *Client) ListWorkspaceMembers(ctx context.Context, in ListMembersInput) (*ListMembersResult, error) {
	var r ListMembersResult
	err := c.Do(ctx, "/api/workspace/members", in, &r)
	return &r, err
}
func (c *Client) GetWorkspaceInfo(ctx context.Context) (*Workspace, error) {
	var w Workspace
	err := c.Do(ctx, "/api/workspace/public", nil, &w)
	return &w, err
}
