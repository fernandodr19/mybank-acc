package accounts

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain/entities"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
)

// Repository of transactions
type Repository interface {
	CreateAccount(ctx context.Context, acc entities.Account) (vos.AccountID, error)
	GetAccountByID(ctx context.Context, accID vos.AccountID) (entities.Account, error)
	Deposit(ctx context.Context, accID vos.AccountID, amount vos.Money) error
	Withdraw(ctx context.Context, accID vos.AccountID, amount vos.Money) error
	DecreaseAvailableCredit(ctx context.Context, accID vos.AccountID, amount vos.Money) error
}

// Usecase of accounts
type Usecase struct {
	accRepo Repository
}

// NewUsecase builds an acc usecase
func NewUsecase(accRepo Repository) *Usecase {
	return &Usecase{
		accRepo: accRepo,
	}
}
