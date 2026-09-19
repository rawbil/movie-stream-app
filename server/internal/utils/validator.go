package utils

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = NewValidator()

func NewValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("password_format", PasswordValidator)
	return validate
}

var (
	PassUpper = regexp.MustCompile(`[A-Z]`)
	PassLower = regexp.MustCompile(`[a-z]`)
	PassNum   = regexp.MustCompile(`[0-9]`)
	PassChar  = regexp.MustCompile(`[^A-Za-z0-9]`)
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	return PassUpper.MatchString(password) && PassLower.MatchString(password) && PassNum.MatchString(password) && PassChar.MatchString(password) && len(password) >= 6 && len(password) <= 20
}

func ValidateUpdateGenre(arg UpdateGenreParams) error {
	return validate.Struct(arg)
}

func ValidationErrors(tag string, err error) bool {
	var validationErrors validator.ValidationErrors

	if !errors.As(err, &validationErrors) {
		return false
	}

	for _, error := range validationErrors {

		if error.Tag() == tag {
			return true
		}
	}
	return false
}
