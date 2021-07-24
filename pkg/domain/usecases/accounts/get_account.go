package accounts

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/entities"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"github.com/sirupsen/logrus"
)

// GetAccountByID retrieves an account based a given ID
func (u Usecase) GetAccountByID(ctx context.Context, accID vos.AccountID) (entities.Account, error) {
	const operation = "accounts.Usecase.GetAccountByID"

	log := logger.FromCtx(ctx).WithFields(logrus.Fields{
		"accID": accID,
	})

	log.Infoln("getting account")

	acc, err := u.accRepo.GetAccountByID(ctx, accID)
	if err != nil {
		return entities.Account{}, domain.Error(operation, err)
	}

	log.WithField("doc", acc.Document).Infoln("account successfully retrieved")

	return acc, nil
}
