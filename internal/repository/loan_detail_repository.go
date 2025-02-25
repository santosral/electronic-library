package repository

import (
	"context"
	"electronic-library/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanDetailRepository interface {
	GetByID(ctx context.Context, id string) (*model.LoanDetail, error)
	CreateLoanDetail(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, error)
	UpdateLoanDetail(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, error)
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
}

type DbLoanDetailRepository struct {
	Pool *pgxpool.Pool
}

func NewLoanDetailRepository(pool *pgxpool.Pool) *DbLoanDetailRepository {
	return &DbLoanDetailRepository{Pool: pool}
}

func (repo *DbLoanDetailRepository) GetByID(ctx context.Context, id string) (*model.LoanDetail, error) {
	query := `
		SELECT
			ID,
			NAME_OF_BORROWER,
			LOAN_DATE,
			RETURN_DATE,
			RETURNED_ON
		FROM
			LOAN_DETAILS
		WHERE
			ID = $1;
	`

	row := repo.Pool.QueryRow(ctx, query, id)

	var loanDetail model.LoanDetail
	err := row.Scan(&loanDetail.ID, &loanDetail.NameOfBorrower, &loanDetail.LoanDate, &loanDetail.ReturnDate, &loanDetail.ReturnedOn)
	if err != nil {
		return nil, err
	}

	return &loanDetail, nil
}

func (repo *DbLoanDetailRepository) CreateLoanDetail(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, error) {
	query := `
		INSERT INTO
			LOAN_DETAILS (NAME_OF_BORROWER, BOOK_ID, LOAN_DATE, RETURN_DATE)
		VALUES
			($1, $2, $3, $4)
		RETURNING ID, NAME_OF_BORROWER, BOOK_ID, LOAN_DATE, RETURN_DATE
	`
	err := repo.Pool.QueryRow(ctx, query, ld.NameOfBorrower, ld.BookID, ld.LoanDate, ld.ReturnDate).
		Scan(&ld.ID, &ld.NameOfBorrower, &ld.BookID, &ld.LoanDate, &ld.ReturnDate)
	if err != nil {
		return nil, err
	}

	return ld, err
}

func (repo *DbLoanDetailRepository) UpdateLoanDetail(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, error) {
	query := `
		UPDATE LOAN_DETAILS
		SET
			RETURN_DATE = $2,
			RETURNED_ON = $3
		WHERE
			ID = $1
		RETURNING ID, NAME_OF_BORROWER, BOOK_ID, LOAN_DATE, RETURN_DATE, RETURNED_ON;
	`

	err := repo.Pool.QueryRow(ctx, query, ld.ID, ld.ReturnDate, ld.ReturnedOn).
		Scan(&ld.ID, &ld.NameOfBorrower, &ld.BookID, &ld.LoanDate, &ld.ReturnDate, &ld.ReturnedOn)
	if err != nil {
		return nil, fmt.Errorf("failed to update loan detail: %w", err)
	}

	return ld, err
}

func (repo *DbLoanDetailRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return repo.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
}
