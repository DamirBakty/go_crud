package repos

import (
	"blog-service/models"
	"database/sql"
	"errors"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(blogID, userID, content string) (*models.Comment, error) {
	if blogID == "" || userID == "" || content == "" {
		return nil, ErrInvalidInput
	}

	var comment models.Comment
	err := r.db.QueryRow(
		`INSERT INTO comments (blog_id, user_id, content) 
		VALUES ($1, $2, $3) 
		RETURNING id, blog_id, user_id, content, created_at`,
		blogID, userID, content,
	).Scan(&comment.ID, &comment.BlogID, &comment.UserID, &comment.Content, &comment.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) GetComment(id string) (*models.Comment, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	var comment models.Comment
	err := r.db.QueryRow(
		`SELECT id, blog_id, user_id, content, created_at 
		FROM comments 
		WHERE id = $1`,
		id,
	).Scan(&comment.ID, &comment.BlogID, &comment.UserID, &comment.Content, &comment.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) UpdateComment(id, content string) (*models.Comment, error) {
	if id == "" || content == "" {
		return nil, ErrInvalidInput
	}

	comment, err := r.GetComment(id)
	if err != nil {
		return nil, err
	}

	comment.Content = content

	_, err = r.db.Exec(
		`UPDATE comments 
		SET content = $1 
		WHERE id = $2`,
		content, id,
	)

	if err != nil {
		return nil, err
	}

	return r.GetComment(id)
}

func (r *CommentRepository) DeleteComment(id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	result, err := r.db.Exec("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrCommentNotFound
	}

	return nil
}

func (r *CommentRepository) ListComments(blogID string, page, limit int) ([]*models.Comment, int, error) {
	if blogID == "" {
		return nil, 0, ErrInvalidInput
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM comments WHERE blog_id = $1", blogID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT id, blog_id, user_id, content, created_at 
		FROM comments 
		WHERE blog_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		blogID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.BlogID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
