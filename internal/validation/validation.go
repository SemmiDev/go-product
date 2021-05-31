package validation

import (
	"github.com/SemmiDev/go-product/internal/constant"
	"github.com/go-playground/validator"
)

func Struct(s interface{}) error {
	err := validator.New().Struct(s)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return constant.NewErrFieldValidation(e)
		}
	}
	return nil
}
