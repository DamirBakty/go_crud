package repos

import (
	"blog-service/models"
	"database/sql"
	"errors"
	"time"
)

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) CreateBlog(userID, title, text string) (*models.Blog, error) {
	if userID == "" || title == "" || text == "" {
		return nil, ErrInvalidInput
	}

	var blog models.Blog
	err := r.db.QueryRow(
		`INSERT INTO blogs (user_id, title, text) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, title, text, created_at, updated_at`,
		userID, title, text,
	).Scan(&blog.ID, &blog.UserID, &blog.Title, &blog.Text, &blog.CreatedAt, &blog.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (r *BlogRepository) GetBlog(id string) (*models.Blog, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	var blog models.Blog
	err := r.db.QueryRow(
		`SELECT id, user_id, title, text, created_at, updated_at 
		FROM blogs 
		WHERE id = $1`,
		id,
	).Scan(&blog.ID, &blog.UserID, &blog.Title, &blog.Text, &blog.CreatedAt, &blog.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBlogNotFound
		}
		return nil, err
	}

	return &blog, nil
}

func (r *BlogRepository) UpdateBlog(id, title, text string) (*models.Blog, error) {
	if id == "" || (title == "" && text == "") {
		return nil, ErrInvalidInput
	}

	blog, err := r.GetBlog(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		blog.Title = title
	}
	if text != "" {
		blog.Text = text
	}

	_, err = r.db.Exec(
		`UPDATE blogs 
		SET title = $1, text = $2, updated_at = $3 
		WHERE id = $4`,
		blog.Title, blog.Text, time.Now(), id,
	)

	if err != nil {
		return nil, err
	}

	return r.GetBlog(id)
}

func (r *BlogRepository) DeleteBlog(id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	result, err := r.db.Exec("DELETE FROM blogs WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBlogNotFound
	}

	return nil
}

func (r *BlogRepository) ListBlogs(userID string, page, limit int) ([]*models.Blog, int, error) {
	if userID == "" {
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
	err := r.db.QueryRow("SELECT COUNT(*) FROM blogs WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT id, user_id, title, text, created_at, updated_at 
		FROM blogs 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var blogs []*models.Blog
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.ID, &blog.UserID, &blog.Title, &blog.Text, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		blogs = append(blogs, &blog)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}
