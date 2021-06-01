package constant

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
)

var (
	ErrServer = errors.New("Something went wrong")

	ErrUrlPathParameter  = errors.New("Invalid url path parameter")
	ErrUrlQueryParameter = errors.New("Invalid url query parameter")
	ErrRequestBody       = errors.New("Invalid request body")
	ErrUnauthorized      = errors.New("You are not authorized to perform this action")
	ErrFieldValidation   = errors.New("Field is not valid")

	ErrMerchantNotFound   = errors.New("Merchant not found")
	ErrEmailRegistered    = errors.New("Email already in use")
	ErrEmailNotRegistered = errors.New("Email not registered")
	ErrWrongPassword      = errors.New("Password incorrect")

	ErrProductNotFound = errors.New("Product not found")
)

func NewErrFieldValidation(err validator.FieldError) error {
	return fmt.Errorf("%s: %w; format must be (%s=%s)", err.Field(), ErrFieldValidation, err.ActualTag(), err.Param())
}
