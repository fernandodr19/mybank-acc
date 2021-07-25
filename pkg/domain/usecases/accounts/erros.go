package accounts

import "errors"

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrAccountConflict     = errors.New("account alreagy regiteres")
	ErrInvalidAccID        = errors.New("invalid account id")
	ErrInvalidCreditLimit  = errors.New("invalid credit limit")
	ErrInvalidDocument     = errors.New("invalid document")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInsufficientCredit  = errors.New("insufficient credit")
)
