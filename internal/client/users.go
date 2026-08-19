package client

import "context"

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var wrapper struct {
		User User `json:"user"`
	}
	err := c.Do(ctx, "/api/users/me", nil, &wrapper)
	return &wrapper.User, err
}
