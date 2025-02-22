package repository

import (
	"context"
	"electronic-library/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookDetailRepository struct {
	Pool *pgxpool.Pool
}

func NewBookDetailRepository(pool *pgxpool.Pool) *BookDetailRepository {
	return &BookDetailRepository{Pool: pool}
}

func (repo *BookDetailRepository) SearchByTitle(ctx context.Context, title string, limit int, offset int) ([]model.BookDetail, int, error) {
	tx, err := repo.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("unable to begin transaction: %w", err)
	}

	query := `
		SELECT id, title, available_copies
		FROM book_details
		WHERE to_tsvector('english', title) @@ plainto_tsquery('english', $1)
		AND available_copies > 0
		ORDER BY ts_rank(to_tsvector('english', title), plainto_tsquery('english', $1)) DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, title, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error performing text search on title: %w", err)
	}
	defer rows.Close()

	bookDetails, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.BookDetail])
	if err != nil {
		return nil, 0, fmt.Errorf("error scanning title: %w", err)
	}

	countQuery := `
		SELECT COUNT(*)
		FROM book_details
		WHERE to_tsvector('english', title) @@ plainto_tsquery('english', $1)
		AND available_copies > 0
	`
	var totalCount int
	err = tx.QueryRow(ctx, countQuery, title).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching total count: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return bookDetails, totalCount, nil
}
