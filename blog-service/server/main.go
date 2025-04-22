package main

import (
	"blog-service/blogpb"
	"blog-service/config"
	"blog-service/repos"
	"blog-service/services"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("Database URL:", cfg.DbUrl)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	blogRepo := repos.NewBlogRepository(cfg.DB)
	commentRepo := repos.NewCommentRepository(cfg.DB)
	likeRepo := repos.NewLikeRepository(cfg.DB)

	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &services.BlogServer{
		Config:      *cfg,
		BlogRepo:    blogRepo,
		CommentRepo: commentRepo,
		LikeRepo:    likeRepo,
	})

	log.Println("Blog service is running on port 50052...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
