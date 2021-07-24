package accounts

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
)

// CreateAccount creates an account
func (u Usecase) CreateAccount(ctx context.Context, doc vos.Document) (vos.AccountID, error) {
	const operation = "accounts.Usecase.CreateAccount"

	if doc == "" {
		return "", ErrInvalidDocument
	}

	accID, err := u.accRepo.CreateAccount(ctx, doc)
	if err != nil {
		return "", domain.Error(operation, err)
	}

	return accID, nil
}
