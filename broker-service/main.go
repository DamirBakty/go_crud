package main

import (
	"broker-service/clients"
	"broker-service/config"
	"broker-service/handlers"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	webPort = "8080"
)

type Config struct {
	AuthGRPCAddress string
	AuthClient      *clients.AuthClient
	BlogGRPCAddress string
	BlogClient      *clients.BlogClient
	AppConfig       *config.AppConfig
	Handlers        *handlers.Config
}

func main() {
	appConfig := &config.AppConfig{}

	app := Config{
		AuthGRPCAddress: "localhost:50051",
		BlogGRPCAddress: "localhost:50052",
		AppConfig:       appConfig,
	}

	authClient, err := clients.NewAuthClient(app.AuthGRPCAddress)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}
	defer authClient.Close()
	app.AuthClient = authClient

	blogClient, err := clients.NewBlogClient(app.BlogGRPCAddress)
	if err != nil {
		log.Fatalf("Failed to create blog client: %v", err)
	}
	defer blogClient.Close()
	app.BlogClient = blogClient

	app.Handlers = &handlers.Config{
		AppConfig:       appConfig,
		AuthGRPCAddress: app.AuthGRPCAddress,
		AuthClient:      app.AuthClient,
		BlogGRPCAddress: app.BlogGRPCAddress,
		BlogClient:      app.BlogClient,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", webPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting broker service on port %s\n", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.Handlers.HandleRoot)
	mux.HandleFunc("/auth", app.Handlers.HandleAuth)
	mux.HandleFunc("/blog", app.Handlers.HandleBlog)
	mux.HandleFunc("/blogs", app.Handlers.HandleListBlogs)

	return mux
}
