package main

import (
	"auth-service/authpb"
	"auth-service/config"
	userrepos "auth-service/repos"
	"auth-service/services"
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
	dbUrl := cfg.DbUrl
	fmt.Println(dbUrl)

	db, err := userrepos.NewPostgresDB(dbUrl)

	appConfig := config.AppConfig{
		DB: db,
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &services.AuthServer{})

	log.Println("Auth service is running on port 50051...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
