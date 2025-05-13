package repository

import "errors"

// ErrTicketNotFound is returned when a ticket is not found in the repository
var ErrTicketNotFound = errors.New("ticket not found")

// ErrMessageNotFound is returned when a message is not found in the repository
var ErrMessageNotFound = errors.New("message not found")
