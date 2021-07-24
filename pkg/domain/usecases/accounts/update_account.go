package accounts

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"github.com/sirupsen/logrus"
)

// Deposit deposits money on an account
func (u Usecase) Deposit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Usecase.Deposit"

	log := logger.FromCtx(ctx).WithFields(logrus.Fields{
		"accID":  accID,
		"amount": amount.Int(),
	})

	log.Infoln("processing a deposit")

	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := u.accRepo.Deposit(ctx, accID, amount)

	if err != nil {
		return domain.Error(operation, err)
	}

	log.Infoln("deposit successfully processed")

	return nil
}

// Withdraw Withdraws money from an account
func (u Usecase) Withdraw(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Usecase.Withdraw"

	log := logger.FromCtx(ctx).WithFields(logrus.Fields{
		"accID":  accID,
		"amount": amount.Int(),
	})

	log.Infoln("processing a withdrawal")

	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := u.accRepo.Withdraw(ctx, accID, amount)

	if err != nil {
		return domain.Error(operation, err)
	}

	log.Infoln("withdrawal successfully processed")

	return nil
}

// ReserveCreditLimit decrease the account's credit limit
func (u Usecase) ReserveCreditLimit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Usecase.ReserveCreditLimit"

	log := logger.FromCtx(ctx).WithFields(logrus.Fields{
		"accID":  accID,
		"amount": amount.Int(),
	})

	log.Infoln("processing a credit reservation")

	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := u.accRepo.ReserveCreditLimit(ctx, accID, amount)

	if err != nil {
		return domain.Error(operation, err)
	}

	log.Infoln("credit limit successfully reserved")

	return nil
}
