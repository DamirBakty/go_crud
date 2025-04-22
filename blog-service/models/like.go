package models

import (
	"time"
)

type Like struct {
	ID        string    `json:"id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLikeRequest struct {
	BlogID string `json:"blog_id"`
	UserID string `json:"user_id"`
}

type LikeResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Like    *Like  `json:"like,omitempty"`
}

type DeleteLikeRequest struct {
	BlogID string `json:"blog_id"`
	UserID string `json:"user_id"`
}

type DeleteLikeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ListLikesRequest struct {
	BlogID string `json:"blog_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type ListLikesResponse struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Likes   []*Like `json:"likes,omitempty"`
	Total   int     `json:"total"`
}

type CheckLikeRequest struct {
	BlogID string `json:"blog_id"`
	UserID string `json:"user_id"`
}

type CheckLikeResponse struct {
	Liked bool `json:"liked"`
}
