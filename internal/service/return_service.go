package service

import (
	"context"
	"electronic-library/internal/model"
	"electronic-library/internal/repository"
	"electronic-library/pkg/errors"
	"time"
)

type ReturnService struct {
	BookDetailRepo repository.BookDetailRepository
	LoanDetailRepo repository.LoanDetailRepository
}

type ReturnedAlreadyError struct {
	Message string
}

func (e *ReturnedAlreadyError) Error() string {
	return e.Message
}

func NewReturnService(bdr repository.BookDetailRepository, ldr repository.LoanDetailRepository) *ReturnService {
	return &ReturnService{
		BookDetailRepo: bdr,
		LoanDetailRepo: ldr,
	}
}

func (s *ReturnService) Call(ctx context.Context, ld *model.LoanDetail) (*model.LoanDetail, *errors.APIError) {
	tx, err := s.LoanDetailRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, errors.New(500, "Transaction error")
	}
	defer tx.Rollback(ctx)

	ld, err = s.LoanDetailRepo.GetByID(ctx, ld.ID)
	if err != nil {
		return nil, errors.New(404, "Loan detail not found")
	}

	if ld.ReturnedOn != nil {
		return nil, errors.New(422, "Book is already returned")
	}

	currentTime := time.Now().UTC()
	ld.ReturnedOn = &currentTime

	updatedLoanDetail, err := s.LoanDetailRepo.UpdateLoanDetail(ctx, ld)
	if err != nil {
		return nil, errors.New(500, "Transaction: failure to update loan detail")
	}

	book, err := s.BookDetailRepo.GetByID(ctx, ld.BookID)
	if err != nil {
		return nil, errors.New(500, "Transaction: failure to find the book")
	}

	err = s.BookDetailRepo.UpdateAvailableCopies(ctx, book.ID, book.AvailableCopies+1)
	if err != nil {
		return nil, errors.New(500, "Transaction: failure to update the book")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.New(500, "Transaction: failure to update the book")
	}

	return updatedLoanDetail, nil
}
