package service

import (
	"context"
	"electronic-library/internal/model"
	"electronic-library/internal/repository"
	"fmt"
	"time"
)

// 3 weeks
const ExtendDuration = 3 * 7 * 24 * time.Hour

type ExtendService struct {
	BookDetailRepo repository.BookDetailRepository
	LoanDetailRepo repository.LoanDetailRepository
}

func NewExtendService(bdr repository.BookDetailRepository, ldr repository.LoanDetailRepository) *ExtendService {
	return &ExtendService{
		BookDetailRepo: bdr,
		LoanDetailRepo: ldr,
	}
}

func (s *ExtendService) Call(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, error) {
	tx, err := s.LoanDetailRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	ld, err = s.LoanDetailRepo.GetByID(ctx, ld.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching loan detail: %w", err)
	}

	newReturnDate := ld.ReturnDate.Add(ExtendDuration)
	ld.ReturnDate = newReturnDate

	updatedLoanDetail, err := s.LoanDetailRepo.UpdateLoanDetail(ctx, ld)
	if err != nil {
		return nil, fmt.Errorf("failed to update loan: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedLoanDetail, nil
}
