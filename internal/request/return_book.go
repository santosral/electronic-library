package request

import (
	"electronic-library/pkg/errors"

	"github.com/go-playground/validator/v10"
)

type ReturnBook struct {
	ID string `json:"id" validate:"uuid4_rfc4122"`
}

func (s *ReturnBook) Validate(v *validator.Validate) error {
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
