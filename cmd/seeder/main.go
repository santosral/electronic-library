package main

import (
	"context"
	"electronic-library/config"
	"electronic-library/internal/db"
	"log"
)

func main() {
	ctx := context.Background()

	cfg_path := "config/config.json"
	cfg, err := config.LoadConfig(cfg_path)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dbConnection, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)
	}
	defer dbConnection.Close()

	sqlFilePath := "internal/db/seeds/initial_data.sql"

	err = db.RunSeeder(ctx, sqlFilePath, dbConnection.Pool)
	if err != nil {
		log.Fatalf("Failed to run seeder: %v\n", err)
	}

	log.Print("SQL seeder executed successfully")
}
