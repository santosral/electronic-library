package service

import (
	"context"
	"electronic-library/internal/model"
	"electronic-library/internal/repository"
	"fmt"
	"time"
)

type ReturnService struct {
	BookDetailRepo repository.BookDetailRepository
	LoanDetailRepo repository.LoanDetailRepository
}

func NewReturnService(bdr repository.BookDetailRepository, ldr repository.LoanDetailRepository) *ReturnService {
	return &ReturnService{
		BookDetailRepo: bdr,
		LoanDetailRepo: ldr,
	}
}

func (s *ReturnService) Call(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, error) {
	tx, err := s.LoanDetailRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	ld, err = s.LoanDetailRepo.GetByID(ctx, ld.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching loan detail: %w", err)
	}

	returnedOn := time.Now().UTC()
	ld.ReturnedOn = &returnedOn

	updatedLoanDetail, err := s.LoanDetailRepo.UpdateLoanDetail(ctx, ld)
	if err != nil {
		return nil, fmt.Errorf("failed to update loan: %w", err)
	}

	book, err := s.BookDetailRepo.GetByID(ctx, ld.BookID)
	if err != nil {
		return nil, fmt.Errorf("error fetching book details: %w", err)
	}

	err = s.BookDetailRepo.UpdateAvailableCopies(ctx, book.ID, book.AvailableCopies+1)
	if err != nil {
		return nil, fmt.Errorf("failed to update book availability: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedLoanDetail, nil
}
