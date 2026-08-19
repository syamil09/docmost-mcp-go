package client

import "encoding/json"

type Page struct {
	ID           string          `json:"id"`
	SpaceID      string          `json:"spaceId"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Content      json.RawMessage `json:"content"`
	Icon         string          `json:"icon"`
	ParentPageID *string         `json:"parentPageId"`
	Position     string          `json:"position"`
	CreatorID    string          `json:"creatorId"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
	DeletedAt    *string         `json:"deletedAt"`
}

type Space struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	CreatorID   string  `json:"creatorId"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	MemberCount int     `json:"memberCount"`
}

type Comment struct {
	ID              string          `json:"id"`
	PageID          string          `json:"pageId"`
	Content         json.RawMessage `json:"content"`
	AuthorID        string          `json:"authorId"`
	ParentCommentID *string         `json:"parentCommentId"`
	Selection       *string         `json:"selection"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`
}

type User struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	AvatarURL       *string `json:"avatarUrl"`
	Role            string  `json:"role"`
	EmailVerifiedAt *string `json:"emailVerifiedAt"`
	CreatedAt       string  `json:"createdAt"`
}

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}