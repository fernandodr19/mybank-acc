package accounts

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/entities"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
)

// GetAccountByID retrieves an account based a given ID
func (u Usecase) GetAccountByID(ctx context.Context, accID vos.AccountID) (entities.Account, error) {
	const operation = "accounts.Usecase.GetAccountByID"

	acc, err := u.accRepo.GetAccountByID(ctx, accID)
	if err != nil {
		return entities.Account{}, domain.Error(operation, err)
	}

	return acc, nil
}
