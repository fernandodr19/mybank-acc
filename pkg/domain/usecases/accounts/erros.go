package accounts

import "errors"

var (
	ErrInvalidDocument     = errors.New("invalid document")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInsufficientCredit  = errors.New("insufficient credit")
)
