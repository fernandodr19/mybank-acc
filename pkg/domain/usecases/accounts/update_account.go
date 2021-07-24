package accounts

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
)

// Deposit deposits money on an account
func (u Usecase) Deposit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Usecase.Deposit"

	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := u.accRepo.Deposit(ctx, accID, amount)

	if err != nil {
		return domain.Error(operation, err)
	}

	return nil
}

// Withdraw Withdraws money from an account
func (u Usecase) Withdraw(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Usecase.Withdraw"

	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := u.accRepo.Withdraw(ctx, accID, amount)

	if err != nil {
		return domain.Error(operation, err)
	}

	return nil
}

// ReserveCreditLimit decrease the account's credit limit
func (u Usecase) ReserveCreditLimit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Usecase.ReserveCreditLimit"

	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := u.accRepo.ReserveCreditLimit(ctx, accID, amount)

	if err != nil {
		return domain.Error(operation, err)
	}

	return nil
}
