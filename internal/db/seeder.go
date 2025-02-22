package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunSeeder(ctx context.Context, sqlFilePath string, dbConnection *pgxpool.Pool) error {
	sql, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	_, err = dbConnection.Exec(ctx, string(sql))
	if err != nil {
		return fmt.Errorf("failed to execute SQL script: %w", err)
	}

	return nil
}
