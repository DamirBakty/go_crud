package clients

import (
	"auth-service/authpb"
	"broker-service/models"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		conn:   conn,
		client: authpb.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) Register(username, email, password string) (*models.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &authpb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	res, err := c.client.Register(ctx, req)
	if err != nil {
		return &models.AuthResponse{
			Error:   true,
			Message: fmt.Sprintf("Error calling Register: %v", err),
		}, nil
	}

	return &models.AuthResponse{
		Error:        false,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}, nil
}

func (c *AuthClient) Login(email, password string) (*models.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &authpb.LoginRequest{
		Email:    email,
		Password: password,
	}

	res, err := c.client.Login(ctx, req)
	if err != nil {
		return &models.AuthResponse{
			Error:   true,
			Message: fmt.Sprintf("Error calling Login: %v", err),
		}, nil
	}

	return &models.AuthResponse{
		Error:        false,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}, nil
}

func (c *AuthClient) VerifyToken(token string) (*models.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &authpb.VerifyTokenRequest{
		AccessToken: token,
	}

	res, err := c.client.VerifyToken(ctx, req)
	if err != nil {
		return &models.AuthResponse{
			Error:   true,
			Message: fmt.Sprintf("Error calling VerifyToken: %v", err),
		}, nil
	}

	return &models.AuthResponse{
		Error:    false,
		Valid:    res.Valid,
		UserID:   res.UserId,
		Username: res.Username,
	}, nil
}
