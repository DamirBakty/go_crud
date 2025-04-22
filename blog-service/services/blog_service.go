package services

import (
	"blog-service/blogpb"
	"blog-service/config"
	"blog-service/repos"
	"context"
	"log"
)

type BlogServer struct {
	blogpb.UnimplementedBlogServiceServer
	Config      config.AppConfig
	BlogRepo    *repos.BlogRepository
	CommentRepo *repos.CommentRepository
	LikeRepo    *repos.LikeRepository
}

func (s *BlogServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.BlogResponse, error) {
	log.Printf("CreateBlog: user_id=%s, title=%s", req.UserId, req.Title)

	blog, err := s.BlogRepo.CreateBlog(req.UserId, req.Title, req.Text)
	if err != nil {
		log.Printf("Failed to create blog: %v", err)
		return &blogpb.BlogResponse{
			Error:   true,
			Message: err.Error(),
		}, nil
	}

	return &blogpb.BlogResponse{
		Error:   false,
		Message: "Blog created successfully",
		Blog: &blogpb.Blog{
			Id:        blog.ID,
			UserId:    blog.UserID,
			Title:     blog.Title,
			Text:      blog.Text,
			CreatedAt: blog.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: blog.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *BlogServer) GetBlog(ctx context.Context, req *blogpb.GetBlogRequest) (*blogpb.BlogResponse, error) {
	log.Printf("GetBlog: id=%s", req.Id)

	blog, err := s.BlogRepo.GetBlog(req.Id)
	if err != nil {
		log.Printf("Failed to get blog: %v", err)
		return &blogpb.BlogResponse{
			Error:   true,
			Message: err.Error(),
		}, nil
	}

	return &blogpb.BlogResponse{
		Error:   false,
		Message: "Blog retrieved successfully",
		Blog: &blogpb.Blog{
			Id:        blog.ID,
			UserId:    blog.UserID,
			Title:     blog.Title,
			Text:      blog.Text,
			CreatedAt: blog.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: blog.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *BlogServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.BlogResponse, error) {
	log.Printf("UpdateBlog: id=%s", req.Id)

	blog, err := s.BlogRepo.UpdateBlog(req.Id, req.Title, req.Text)
	if err != nil {
		log.Printf("Failed to update blog: %v", err)
		return &blogpb.BlogResponse{
			Error:   true,
			Message: err.Error(),
		}, nil
	}

	return &blogpb.BlogResponse{
		Error:   false,
		Message: "Blog updated successfully",
		Blog: &blogpb.Blog{
			Id:        blog.ID,
			UserId:    blog.UserID,
			Title:     blog.Title,
			Text:      blog.Text,
			CreatedAt: blog.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: blog.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *BlogServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	log.Printf("DeleteBlog: id=%s", req.Id)

	err := s.BlogRepo.DeleteBlog(req.Id)
	if err != nil {
		log.Printf("Failed to delete blog: %v", err)
		return &blogpb.DeleteBlogResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &blogpb.DeleteBlogResponse{
		Success: true,
		Message: "Blog deleted successfully",
	}, nil
}

func (s *BlogServer) ListBlogs(ctx context.Context, req *blogpb.ListBlogsRequest) (*blogpb.ListBlogsResponse, error) {
	log.Printf("ListBlogs: user_id=%s, page=%d, limit=%d", req.UserId, req.Page, req.Limit)

	blogs, total, err := s.BlogRepo.ListBlogs(req.UserId, int(req.Page), int(req.Limit))
	if err != nil {
		log.Printf("Failed to list blogs: %v", err)
		return &blogpb.ListBlogsResponse{
			Error:   true,
			Message: err.Error(),
		}, nil
	}

	var pbBlogs []*blogpb.Blog
	for _, blog := range blogs {
		pbBlogs = append(pbBlogs, &blogpb.Blog{
			Id:        blog.ID,
			UserId:    blog.UserID,
			Title:     blog.Title,
			Text:      blog.Text,
			CreatedAt: blog.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: blog.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &blogpb.ListBlogsResponse{
		Error:   false,
		Message: "Blogs listed successfully",
		Blogs:   pbBlogs,
		Total:   int32(total),
	}, nil
}
