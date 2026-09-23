package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AllFieldsRequiredError = errors.New("All fields are required")
	NoRecordError          = errors.New("No Record Found")
	DuplicateRecordError   = errors.New("Record already exists")
	MinTitleError = errors.New("title should have a minimum of 3 characters")
	MaxTitleError = errors.New("title should have a maximum of 100 characters")
	InvalidUrlError = errors.New("poster_path should be a valid url")
	GenreMissing = errors.New("provide at least one genre")
	MovieExistsError = errors.New("movie already exists")
	InvalidEmailFormat = errors.New("invalid email format")
	InvalidPasswordFormat = errors.New("password should be at between 6-20 characters long, have at least alphanumerical and have at least one special character")
	NoEmptyGenre = errors.New("All genre fields should be populated")
)

func ErrorResponse(c *gin.Context, status_code int, error_msg string, err error) {
	c.JSON(status_code, gin.H{
		"error": error_msg,
	})

	Log.Error(strings.ToUpper(error_msg), "error", err)
}
