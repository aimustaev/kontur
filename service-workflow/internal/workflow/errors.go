package workflow

import "errors"

var (
	// ErrEmptyMessage is returned when the message is empty
	ErrEmptyMessage = errors.New("message cannot be empty")
	// ErrInvalidInput is returned when the input type is invalid
	ErrInvalidInput = errors.New("invalid input type")
)
