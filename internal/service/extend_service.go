package service

import (
	"context"
	"electronic-library/internal/model"
	"electronic-library/internal/repository"
	"electronic-library/pkg/errors"
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

func (s *ExtendService) Call(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, *errors.APIError) {
	tx, err := s.LoanDetailRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, errors.New(500, "Transaction failed for extend service")
	}
	defer tx.Rollback(ctx)

	ld, err = s.LoanDetailRepo.GetByID(ctx, ld.ID)
	if err != nil {
		return nil, errors.New(401, "loan detail not found")
	}

	newReturnDate := ld.ReturnDate.Add(ExtendDuration)
	ld.ReturnDate = newReturnDate

	updatedLoanDetail, err := s.LoanDetailRepo.UpdateLoanDetail(ctx, ld)
	if err != nil {
		return nil, errors.New(500, "Transaction: failed to update loan detail")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.New(500, "Transaction: failed to commit")
	}

	return updatedLoanDetail, nil
}
