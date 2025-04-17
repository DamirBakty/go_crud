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
	Email    string
	Password string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, password string) (*User, error) {
	// Check if user with this email already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailAlreadyExists
	}

	// Hash the password
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// Insert the new user
	var id string
	err = r.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		email, string(hashedPassword),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Email:    email,
		Password: string(hashedPassword),
	}, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
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

	// Hash the provided password and compare with stored hash
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	if user.Password != hashedPassword {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
