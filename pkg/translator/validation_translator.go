package validatorhelper

import (
	"fmt"
	"github.com/go-playground/validator"
	"strings"
)

// ErrorMessages is a map that holds the user-friendly messages for validation errors
var ErrorMessages = map[string]string{
	"required": "{field} is required.",
	"min":      "{field} must be at least {param} characters long.",
	"max":      "{field} must be no more than {param} characters long.",
	"email":    "{field} must be a valid email address.",
	"eqfield":  "{field} must match {param}.",
}

// Validator is the instance of the validator
var Validator *validator.Validate

func init() {
	Validator = validator.New()
}

// ValidateStruct validates the struct and returns error messages
func ValidateStruct(s interface{}) error {
	if err := Validator.Struct(s); err != nil {
		return translateValidationErrors(err.(validator.ValidationErrors))
	}
	return nil
}

// translateValidationErrors translates validation errors into user-friendly messages
func translateValidationErrors(errs validator.ValidationErrors) error {
	var errMessages []string
	for _, fieldErr := range errs {
		errMessage := formatErrorMessage(fieldErr)
		errMessages = append(errMessages, errMessage)
	}
	return fmt.Errorf(strings.Join(errMessages, "; "))
}

// formatErrorMessage formats an error message based on the validation tag and field name
func formatErrorMessage(fieldErr validator.FieldError) string {
	message, ok := ErrorMessages[fieldErr.Tag()]
	if !ok {
		// Default message if the tag is not defined in the map
		message = fmt.Sprintf("{field} is invalid.")
	}

	// Replace placeholders in the message with actual values
	message = strings.ReplaceAll(message, "{field}", fieldErr.Field())
	message = strings.ReplaceAll(message, "{param}", fieldErr.Param())

	return message
}
