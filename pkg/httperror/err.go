package httperror

import "fmt"

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func New(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code int, err error, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: fmt.Sprintf("%s: %v", message, err),
	}
}
