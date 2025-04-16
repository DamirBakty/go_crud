package services

import (
	"auth-service/authpb"
	"auth-service/server/AppConfig"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
	//"time"
)

var secretKey = []byte("secret-key")

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	Config AppConfig
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	log.Printf("Register: %s, %s", req.Username, req.Email)
	return &authpb.AuthResponse{
		AccessToken:  "fake-access-token",
		RefreshToken: "fake-refresh-token",
		ExpiresIn:    "3600",
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	log.Printf("Login: %s", req.Email)
	return &authpb.AuthResponse{
		AccessToken:  "access-token-login",
		RefreshToken: "refresh-token-login",
		ExpiresIn:    "3600",
	}, nil
}

func (s *AuthServer) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	valid := req.AccessToken == "access-token-login"
	return &authpb.VerifyTokenResponse{
		Valid:    valid,
		UserId:   "user123",
		Username: "damir",
	}, nil
}
