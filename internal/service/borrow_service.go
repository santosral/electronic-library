package service

import (
	"context"
	"electronic-library/internal/model"
	"electronic-library/internal/repository"
	"electronic-library/pkg/errors"
	"time"
)

const loanDuration = 4 * 7 * 24 * time.Hour

type BorrowService struct {
	BookDetailRepo repository.BookDetailRepository
	LoanDetailRepo repository.LoanDetailRepository
}

func NewBorrowService(bdr repository.BookDetailRepository, ldr repository.LoanDetailRepository) *BorrowService {
	return &BorrowService{
		BookDetailRepo: bdr,
		LoanDetailRepo: ldr,
	}
}

func (s *BorrowService) Call(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, *errors.APIError) {
	tx, err := s.BookDetailRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, errors.New(500, "Transaction error")
	}
	defer tx.Rollback(ctx)

	book, err := s.BookDetailRepo.GetByID(ctx, ld.BookID)
	if err != nil {
		return nil, errors.New(404, "Book detail not found")
	}

	if book.AvailableCopies == 0 {
		return nil, errors.New(422, "Book is unavailable")
	}

	currentTime := time.Now().UTC()
	ld.LoanDate = currentTime
	ld.ReturnDate = currentTime.Add(loanDuration)

	createdLoanDetail, err := s.LoanDetailRepo.CreateLoanDetail(ctx, ld)
	if err != nil {
		return nil, errors.New(500, "Transaction: failed to create loan detail")
	}

	err = s.BookDetailRepo.UpdateAvailableCopies(ctx, book.ID, book.AvailableCopies-1)
	if err != nil {
		return nil, errors.New(500, "Transaction: failed to update book detail")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.New(500, "Transaction: failed to commit transaction")
	}

	return createdLoanDetail, nil
}
