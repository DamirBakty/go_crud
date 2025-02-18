package server

import (
	"crud/handlers"
	"crud/repos"
	"crud/services"
	"log"
	"net/http"
)

func Run() {
	db, err := repos.NewPostgresDB("postgres://test_user:test_password@localhost:5432/test_market?sslmode=disable")
	marketRepo := repos.NewMarketRepo(db)
	marketService := services.NewMarketService(marketRepo)
	marketHandler := handlers.NewMarketHandler(marketService)
	itemRepo := repos.NewItemRepo(db)
	itemService := services.NewItemService(itemRepo)
	itemHandler := handlers.NewItemHandler(itemService)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		itemHandler.HandleItem(w, r)
	})

	mux.HandleFunc("/markets", func(w http.ResponseWriter, r *http.Request) {
		marketHandler.HandleMarket(w, r)
	})

	log.Println("server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
