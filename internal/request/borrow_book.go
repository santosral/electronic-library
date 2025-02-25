package request

import (
	"electronic-library/pkg/errors"

	"github.com/go-playground/validator/v10"
)

type BorrowBook struct {
	BookID         string `json:"book_id" validate:"uuid4_rfc4122"`
	NameOfBorrower string `json:"name_of_borrower" validate:"required"`
}

func (s *BorrowBook) Validate(v *validator.Validate) error {
	err := v.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		if len(validationErrors) > 0 {
			firstErr := validationErrors[0]
			return errors.FormatError(firstErr)
		}
	}
	return nil
}
