package repos

import (
	"blog-service/models"
	"database/sql"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) CreateLike(blogID, userID string) (*models.Like, error) {
	if blogID == "" || userID == "" {
		return nil, ErrInvalidInput
	}

	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE blog_id = $1 AND user_id = $2", blogID, userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrAlreadyLiked
	}

	var like models.Like
	err = r.db.QueryRow(
		`INSERT INTO likes (blog_id, user_id) 
		VALUES ($1, $2) 
		RETURNING id, blog_id, user_id, created_at`,
		blogID, userID,
	).Scan(&like.ID, &like.BlogID, &like.UserID, &like.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &like, nil
}

func (r *LikeRepository) DeleteLike(blogID, userID string) error {
	if blogID == "" || userID == "" {
		return ErrInvalidInput
	}

	result, err := r.db.Exec("DELETE FROM likes WHERE blog_id = $1 AND user_id = $2", blogID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrLikeNotFound
	}

	return nil
}

func (r *LikeRepository) ListLikes(blogID string, page, limit int) ([]*models.Like, int, error) {
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
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE blog_id = $1", blogID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT id, blog_id, user_id, created_at 
		FROM likes 
		WHERE blog_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		blogID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var likes []*models.Like
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.BlogID, &like.UserID, &like.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		likes = append(likes, &like)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return likes, total, nil
}

func (r *LikeRepository) CheckLike(blogID, userID string) (bool, error) {
	if blogID == "" || userID == "" {
		return false, ErrInvalidInput
	}

	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE blog_id = $1 AND user_id = $2", blogID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
