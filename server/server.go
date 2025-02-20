package server

import (
	"crud/config"
	markethandlers "crud/market-service/handlers"
	marketrepos "crud/market-service/repos"
	marketservices "crud/market-service/services"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbUrl := cfg.DbUrl
	fmt.Println(dbUrl)
	db, err := marketrepos.NewPostgresDB(dbUrl)

	marketRepo := marketrepos.NewMarketRepo(db)
	marketService := marketservices.NewMarketService(marketRepo)
	marketHandler := markethandlers.NewMarketHandler(marketService)
	itemRepo := marketrepos.NewItemRepo(db)
	itemService := marketservices.NewItemService(itemRepo)
	itemHandler := markethandlers.NewItemHandler(itemService)
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
