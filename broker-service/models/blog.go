package models

type Blog struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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

type DeleteBlogResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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

type Comment struct {
	ID        string `json:"id"`
	BlogID    string `json:"blog_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
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

type DeleteCommentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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

type Like struct {
	ID        string `json:"id"`
	BlogID    string `json:"blog_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
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
