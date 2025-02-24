package main

import (
	"context"
	"electronic-library/config"
	"electronic-library/internal/db"
	"electronic-library/internal/handler"
	"electronic-library/internal/repository"
	"electronic-library/internal/service"
	"log"
	"net/http"
	"strconv"
)

func main() {
	ctx := context.Background()

	cfg_path := "config/config.json"
	cfg, err := config.LoadConfig(cfg_path)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Successfully loaded config from file: %s", cfg_path)

	log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)

	dbConnection, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer dbConnection.Close()
	log.Printf("Successfully connected to the database")

	bookRepo := repository.NewBookDetailRepository(dbConnection.Pool)
	loanDetailRepo := repository.NewLoanDetailRepository(dbConnection.Pool)

	borrowService := service.NewBorrowService(bookRepo, loanDetailRepo)
	extendService := service.NewExtendService(bookRepo, loanDetailRepo)
	returnService := service.NewReturnService(bookRepo, loanDetailRepo)

	bookDetailHandler := handler.NewBookDetailHandler(bookRepo, *borrowService, *returnService)
	loanDetailHandler := handler.NewLoanDetailHandler(*extendService)

	mux := http.NewServeMux()

	mux.HandleFunc("/book-details/search", bookDetailHandler.SearchBooks)
	mux.HandleFunc("/borrow", bookDetailHandler.BorrowBook)
	mux.HandleFunc("/extend", loanDetailHandler.Extend)
	mux.HandleFunc("/return", bookDetailHandler.ReturnBook)

	portStr := strconv.Itoa(cfg.Server.Port)
	log.Printf("Server is running on port %s...\n", portStr)
	log.Fatal(http.ListenAndServe(":"+portStr, mux))
}
