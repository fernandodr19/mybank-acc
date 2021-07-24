package accounts

import "errors"

var (
	ErrInvalidDocument = errors.New("invalid document")
	ErrInvalidAmount   = errors.New("invalid amount")
)
