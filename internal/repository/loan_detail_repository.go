package repository

import (
	"context"
	"electronic-library/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanDetailRepository interface {
	CreateLoanDetail(ctx context.Context, loan *model.LoanDetail) (*model.LoanDetail, error)
}

type DbLoanDetailRepository struct {
	Pool *pgxpool.Pool
}

func NewLoanDetailRepository(pool *pgxpool.Pool) *DbLoanDetailRepository {
	return &DbLoanDetailRepository{Pool: pool}
}

func (r *DbLoanDetailRepository) CreateLoanDetail(ctx context.Context, loan *model.LoanDetail) (*model.LoanDetail, error) {
	query := `
		INSERT INTO
			LOAN_DETAILS (NAME_OF_BORROWER, BOOK_ID, LOAN_DATE, RETURN_DATE)
		VALUES
			($1, $2, $3, $4)
		RETURNING ID, NAME_OF_BORROWER, BOOK_ID, LOAN_DATE, RETURN_DATE
	`
	err := r.Pool.QueryRow(ctx, query, loan.NameOfBorrower, loan.BookID, loan.LoanDate, loan.ReturnDate).
		Scan(&loan.ID, &loan.NameOfBorrower, &loan.BookID, &loan.LoanDate, &loan.ReturnDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create loan detail: %w", err)
	}

	return loan, err
}
