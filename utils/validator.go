package utils

import (
	"github.com/go-playground/validator"
)

// Validator instance (singleton-like, initialized once)
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct based on the 'validate' tags
func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}
