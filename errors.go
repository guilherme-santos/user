package user

import "fmt"

type ValidationError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewMissingFieldError(field string) *ValidationError {
	return &ValidationError{
		Code:    "missing_field",
		Field:   field,
		Message: "Field was not provided",
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("code=%s field=%s message=%s", e.Code, e.Field, e.Message)
}
