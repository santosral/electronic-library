package service

import (
	"context"
	"electronic-library/internal/model"
	"electronic-library/internal/repository"
	"fmt"
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

func (s *BorrowService) Call(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, error) {
	tx, err := s.BookDetailRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	book, err := s.BookDetailRepo.GetByID(ctx, ld.BookID)
	if err != nil {
		return nil, fmt.Errorf("error fetching book details: %w", err)
	}

	if book.AvailableCopies == 0 {
		return nil, fmt.Errorf("no copies available for loan")
	}

	currentTime := time.Now().UTC()
	ld.LoanDate = currentTime
	ld.ReturnDate = currentTime.Add(loanDuration)

	createdLoanDetail, err := s.LoanDetailRepo.CreateLoanDetail(ctx, ld)
	if err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	err = s.BookDetailRepo.UpdateAvailableCopies(ctx, book.ID, book.AvailableCopies-1)
	if err != nil {
		return nil, fmt.Errorf("failed to update book availability: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdLoanDetail, nil
}
