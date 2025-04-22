package repos

import (
	"errors"
)

var (
	ErrBlogNotFound    = errors.New("blog not found")
	ErrCommentNotFound = errors.New("comment not found")
	ErrLikeNotFound    = errors.New("like not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrAlreadyLiked    = errors.New("user has already liked this blog")
)
