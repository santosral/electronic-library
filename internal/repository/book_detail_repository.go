package repository

import (
	"context"
	"electronic-library/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookDetailRepository interface {
	GetByID(ctx context.Context, id string) (*model.BookDetail, error)
	SearchByTitle(ctx context.Context, title string, limit int, offset int) ([]model.BookDetail, error)
	UpdateAvailableCopies(ctx context.Context, id string, newAvailableCopies int) error
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
}

type BookNotFound struct {
	Message string
}

func (e *BookNotFound) Error() string {
	return e.Message
}

type DbBookDetailRepository struct {
	Pool *pgxpool.Pool
}

func NewBookDetailRepository(pool *pgxpool.Pool) *DbBookDetailRepository {
	return &DbBookDetailRepository{Pool: pool}
}

func (repo *DbBookDetailRepository) GetByID(ctx context.Context, id string) (*model.BookDetail, error) {
	query := `
		SELECT
			ID,
			TITLE,
			AVAILABLE_COPIES
		FROM
			BOOK_DETAILS
		WHERE
			ID = $1;
	`

	row := repo.Pool.QueryRow(ctx, query, id)

	var bookDetail model.BookDetail
	err := row.Scan(&bookDetail.ID, &bookDetail.Title, &bookDetail.AvailableCopies)
	if err != nil {
		return nil, &BookNotFound{Message: "book not found"}
	}

	return &bookDetail, nil
}

func (repo *DbBookDetailRepository) SearchByTitle(ctx context.Context, title string, limit int, offset int) ([]model.BookDetail, error) {
	query := `
		SELECT
			ID,
			TITLE,
			AVAILABLE_COPIES
		FROM
			BOOK_DETAILS
		WHERE
			TO_TSVECTOR('english', TITLE) @@ WEBSEARCH_TO_TSQUERY('english', $1)
			AND AVAILABLE_COPIES > 0
		ORDER BY
			TS_RANK(
				TO_TSVECTOR('english', TITLE),
				WEBSEARCH_TO_TSQUERY('english', $1)
			) DESC
		LIMIT
			$2
		OFFSET
			$3;
	`

	rows, err := repo.Pool.Query(ctx, query, title, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookDetails, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.BookDetail])
	if err != nil {
		return nil, err
	}

	return bookDetails, nil
}

func (repo *DbBookDetailRepository) UpdateAvailableCopies(ctx context.Context, id string, newAvailableCopies int) error {
	query := `
		UPDATE BOOK_DETAILS
		SET AVAILABLE_COPIES = $1
		WHERE ID = $2;
	`
	_, err := repo.Pool.Exec(ctx, query, newAvailableCopies, id)
	return err
}

func (repo *DbBookDetailRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return repo.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
}
