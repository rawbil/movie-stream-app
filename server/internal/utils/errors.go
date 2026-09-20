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
)

func ErrorResponse(c *gin.Context, status_code int, error_msg string, err error) {
	c.JSON(status_code, gin.H{
		"error": error_msg,
	})

	Log.Error(strings.ToUpper(error_msg), "error", err)
}
