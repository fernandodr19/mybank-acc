package accounts

import "errors"

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrInvalidAccID        = errors.New("invalid account id")
	ErrInvalidDocument     = errors.New("invalid document")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInsufficientCredit  = errors.New("insufficient credit")
)
