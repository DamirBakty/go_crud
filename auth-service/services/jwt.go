package services

import (
	"auth-service/authpb"
	"auth-service/config"
	"auth-service/repos"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

var (
	accessSecretKey  = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	refreshSecretKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	Config   config.AppConfig
	UserRepo *repos.UserRepository
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	log.Printf("Register: %s, %s", req.Username, req.Email)

	user, err := s.UserRepo.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}

	accessToken, err := generateToken(user.ID, user.Email, AccessToken, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(user.ID, user.Email, RefreshToken, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &authpb.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    "86400",
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	log.Printf("Login: %s", req.Email)

	user, err := s.UserRepo.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to validate credentials: %v", err)
		return nil, err
	}

	accessToken, err := generateToken(user.ID, user.Email, AccessToken, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(user.ID, user.Email, RefreshToken, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &authpb.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    "86400",
	}, nil
}

func generateToken(userID string, email string, tokenType TokenType, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	var secretKey []byte
	if tokenType == AccessToken {
		secretKey = accessSecretKey
	} else {
		secretKey = refreshSecretKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (s *AuthServer) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	token, err := jwt.ParseWithClaims(req.AccessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return accessSecretKey, nil
	})

	if err != nil {
		return &authpb.VerifyTokenResponse{
			Valid: false,
		}, nil
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &authpb.VerifyTokenResponse{
			Valid:    true,
			UserId:   claims.UserID,
			Username: claims.Email,
		}, nil
	}

	return &authpb.VerifyTokenResponse{
		Valid: false,
	}, nil
}
