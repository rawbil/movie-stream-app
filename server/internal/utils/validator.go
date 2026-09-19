package utils

import "github.com/go-playground/validator/v10"

var validate = validator.New(validator.WithRequiredStructEnabled())

// func ValidateGenre(genre CreateGenre) error {
// 	return validate.Struct()
// }
