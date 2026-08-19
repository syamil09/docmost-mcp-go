package client

import "context"

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var u User
	err := c.Do(ctx, "/api/users/me", nil, &u)
	return &u, err
}
