package models

import (
	"time"
)

type Blog struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBlogRequest struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type UpdateBlogRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type BlogResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Blog    *Blog  `json:"blog,omitempty"`
}

type ListBlogsRequest struct {
	UserID string `json:"user_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type ListBlogsResponse struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Blogs   []*Blog `json:"blogs,omitempty"`
	Total   int     `json:"total"`
}
