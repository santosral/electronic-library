package errors

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func FormatError(err validator.FieldError) error {
	return errors.New("Field validation for '" + err.Field() + "' failed on the '" + err.Tag() + "' tag")
}
