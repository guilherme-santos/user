package user

import "fmt"

type Error struct {
	Type    Type   `json:"-"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%s message=%s", e.Code, e.Message)
}

// uerr is used internally because other more specific errors wanted to
// embed Error, but we also have a .Error() method.
type Err = Error

type ValidationError struct {
	Err
	Field string `json:"field"`
}

func NewMissingFieldError(field string) *ValidationError {
	return &ValidationError{
		Err: Error{
			Type:    InvalidArgument,
			Code:    "missing_field",
			Message: "Field was not provided",
		},
		Field: field,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("code=%s field=%s message=%s", e.Code, e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return &e.Err
}

type Type uint16

const (
	Unknown Type = iota
	InvalidArgument
	NotFound
)
