package model

import (
	"time"
)

type LoanDetail struct {
	ID             string     `json:"id"`
	NameOfBorrower string     `json:"name_of_borrower"`
	LoanDate       time.Time  `json:"loan_date"`
	ReturnDate     time.Time  `json:"return_date"`
	BookID         string     `json:"book_id,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}
