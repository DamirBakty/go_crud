package clients

import (
	"blog-service/blogpb"
	"broker-service/models"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type BlogClient struct {
	conn   *grpc.ClientConn
	client blogpb.BlogServiceClient
}

func NewBlogClient(addr string) (*BlogClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blog service: %w", err)
	}

	return &BlogClient{
		conn:   conn,
		client: blogpb.NewBlogServiceClient(conn),
	}, nil
}

func (c *BlogClient) Close() error {
	return c.conn.Close()
}

func (c *BlogClient) CreateBlog(userID, title, text string) (*models.BlogResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &blogpb.CreateBlogRequest{
		UserId: userID,
		Title:  title,
		Text:   text,
	}

	res, err := c.client.CreateBlog(ctx, req)
	if err != nil {
		return &models.BlogResponse{
			Error:   true,
			Message: fmt.Sprintf("Error calling CreateBlog: %v", err),
		}, nil
	}

	return &models.BlogResponse{
		Error:   res.Error,
		Message: res.Message,
		Blog: &models.Blog{
			ID:        res.Blog.Id,
			UserID:    res.Blog.UserId,
			Title:     res.Blog.Title,
			Text:      res.Blog.Text,
			CreatedAt: res.Blog.CreatedAt,
			UpdatedAt: res.Blog.UpdatedAt,
		},
	}, nil
}

func (c *BlogClient) GetBlog(id string) (*models.BlogResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &blogpb.GetBlogRequest{
		Id: id,
	}

	res, err := c.client.GetBlog(ctx, req)
	if err != nil {
		return &models.BlogResponse{
			Error:   true,
			Message: fmt.Sprintf("Error calling GetBlog: %v", err),
		}, nil
	}

	return &models.BlogResponse{
		Error:   res.Error,
		Message: res.Message,
		Blog: &models.Blog{
			ID:        res.Blog.Id,
			UserID:    res.Blog.UserId,
			Title:     res.Blog.Title,
			Text:      res.Blog.Text,
			CreatedAt: res.Blog.CreatedAt,
			UpdatedAt: res.Blog.UpdatedAt,
		},
	}, nil
}

func (c *BlogClient) UpdateBlog(id, title, text string) (*models.BlogResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &blogpb.UpdateBlogRequest{
		Id:    id,
		Title: title,
		Text:  text,
	}

	res, err := c.client.UpdateBlog(ctx, req)
	if err != nil {
		return &models.BlogResponse{
			Error:   true,
			Message: fmt.Sprintf("Error calling UpdateBlog: %v", err),
		}, nil
	}

	return &models.BlogResponse{
		Error:   res.Error,
		Message: res.Message,
		Blog: &models.Blog{
			ID:        res.Blog.Id,
			UserID:    res.Blog.UserId,
			Title:     res.Blog.Title,
			Text:      res.Blog.Text,
			CreatedAt: res.Blog.CreatedAt,
			UpdatedAt: res.Blog.UpdatedAt,
		},
	}, nil
}

func (c *BlogClient) DeleteBlog(id string) (*models.DeleteBlogResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &blogpb.DeleteBlogRequest{
		Id: id,
	}

	res, err := c.client.DeleteBlog(ctx, req)
	if err != nil {
		return &models.DeleteBlogResponse{
			Success: false,
			Message: fmt.Sprintf("Error calling DeleteBlog: %v", err),
		}, nil
	}

	return &models.DeleteBlogResponse{
		Success: res.Success,
		Message: res.Message,
	}, nil
}

func (c *BlogClient) ListBlogs(userID string, page, limit int) (*models.ListBlogsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &blogpb.ListBlogsRequest{
		UserId: userID,
		Page:   int32(page),
		Limit:  int32(limit),
	}

	res, err := c.client.ListBlogs(ctx, req)
	if err != nil {
		return &models.ListBlogsResponse{
			Error:   true,
			Message: fmt.Sprintf("Error calling ListBlogs: %v", err),
		}, nil
	}

	var blogs []*models.Blog
	for _, blog := range res.Blogs {
		blogs = append(blogs, &models.Blog{
			ID:        blog.Id,
			UserID:    blog.UserId,
			Title:     blog.Title,
			Text:      blog.Text,
			CreatedAt: blog.CreatedAt,
			UpdatedAt: blog.UpdatedAt,
		})
	}

	return &models.ListBlogsResponse{
		Error:   res.Error,
		Message: res.Message,
		Blogs:   blogs,
		Total:   int(res.Total),
	}, nil
}
