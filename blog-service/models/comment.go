package models

import (
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCommentRequest struct {
	BlogID  string `json:"blog_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type UpdateCommentRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type CommentResponse struct {
	Error   bool     `json:"error"`
	Message string   `json:"message"`
	Comment *Comment `json:"comment,omitempty"`
}

type ListCommentsRequest struct {
	BlogID string `json:"blog_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type ListCommentsResponse struct {
	Error    bool       `json:"error"`
	Message  string     `json:"message"`
	Comments []*Comment `json:"comments,omitempty"`
	Total    int        `json:"total"`
}

type DeleteCommentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
