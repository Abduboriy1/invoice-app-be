package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal error")
)

type AppError struct {
	Err     error
	Message string
	Code    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Err:     ErrNotFound,
		Message: message,
		Code:    "NOT_FOUND",
	}
}

func NewInvalidInput(message string) *AppError {
	return &AppError{
		Err:     ErrInvalidInput,
		Message: message,
		Code:    "INVALID_INPUT",
	}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{
		Err:     ErrUnauthorized,
		Message: message,
		Code:    "UNAUTHORIZED",
	}
}

func NewInternal(message string) *AppError {
	return &AppError{
		Err:     ErrInternal,
		Message: message,
		Code:    "INTERNAL_ERROR",
	}
}
