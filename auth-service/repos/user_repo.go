package repos

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, email, password string) (*User, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailAlreadyExists
	}

	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var id string
	err = r.db.QueryRow(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		username, email, string(hashedPassword),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		"SELECT id, username, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ValidateCredentials(email, password string) (*User, error) {
	user, err := r.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	if user.Password != hashedPassword {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
