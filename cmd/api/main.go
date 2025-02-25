package main

import (
	"context"
	"electronic-library/config"
	"electronic-library/internal/db"
	"electronic-library/internal/handler"
	"electronic-library/internal/repository"
	"electronic-library/internal/routes"
	"electronic-library/internal/service"
	"electronic-library/pkg"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	logger := pkg.NewLogger()
	ctx := context.Background()
	ctx = pkg.SetLoggerToContext(ctx, logger)

	cfg_path := "config/config.json"
	cfg, err := config.LoadConfig(cfg_path)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to load config", "path", cfg_path, "error", err)
		return
	}
	logger.InfoContext(ctx, "Successfully loaded config", "configPath", cfg_path)

	logger.InfoContext(ctx, "Starting database connection")
	dbConnection, err := db.Connect(ctx, cfg)
	if err != nil {
		logger.ErrorContext(ctx, "Error connecting to the database", "trace", err)
		return
	}
	defer dbConnection.Close()
	logger.InfoContext(ctx, "Successfully connected to the database")

	bookRepo := repository.NewBookDetailRepository(dbConnection.Pool)
	loanDetailRepo := repository.NewLoanDetailRepository(dbConnection.Pool)

	borrowService := service.NewBorrowService(bookRepo, loanDetailRepo)
	extendService := service.NewExtendService(bookRepo, loanDetailRepo)
	returnService := service.NewReturnService(bookRepo, loanDetailRepo)

	bookDetailHandler := handler.NewBookDetailHandler(bookRepo, borrowService, returnService)
	loanDetailHandler := handler.NewLoanDetailHandler(extendService)

	mux := http.NewServeMux()

	routes.SetupRoutes(mux, bookDetailHandler, loanDetailHandler)

	portStr := strconv.Itoa(cfg.Server.Port)
	server := &http.Server{
		Addr:    ":" + portStr,
		Handler: mux,
	}

	logger.InfoContext(ctx, "Server is starting", "port", portStr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(ctx, "Error starting server", "trace", err)
		}
	}()

	<-stop
	logger.InfoContext(ctx, "Received shutdown signal, shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.ErrorContext(ctx, "Error during server shutdown", "trace", err)
		return
	}

	logger.InfoContext(ctx, "Server gracefully stopped")
}
