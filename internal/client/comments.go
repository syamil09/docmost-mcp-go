package client

import "context"

type ListCommentsResult struct {
	Items      []Comment `json:"items"`
	NextCursor *string   `json:"nextCursor,omitempty"`
}

type CreateCommentInput struct {
	PageID          string `json:"pageId"`
	Content         string `json:"content"`
	Selection       string `json:"selection,omitempty"`
	ParentCommentID string `json:"parentCommentId,omitempty"`
}

type UpdateCommentInput struct {
	CommentID string `json:"commentId"`
	Content   string `json:"content"`
}

func (c *Client) ListComments(ctx context.Context, pageID string, limit int, cursor *string) (*ListCommentsResult, error) {
	var r ListCommentsResult
	err := c.Do(ctx, "/api/comments/", map[string]any{"pageId": pageID, "limit": limit, "cursor": cursor}, &r)
	return &r, err
}
func (c *Client) CreateComment(ctx context.Context, in CreateCommentInput) (*Comment, error) {
	var c2 Comment
	err := c.Do(ctx, "/api/comments/create", in, &c2)
	return &c2, err
}
func (c *Client) UpdateComment(ctx context.Context, in UpdateCommentInput) (*Comment, error) {
	var c2 Comment
	err := c.Do(ctx, "/api/comments/update", in, &c2)
	return &c2, err
}
