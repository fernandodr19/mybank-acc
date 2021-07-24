package accounts

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"github.com/sirupsen/logrus"
)

// CreateAccount creates an account
func (u Usecase) CreateAccount(ctx context.Context, doc vos.Document) (vos.AccountID, error) {
	const operation = "accounts.Usecase.CreateAccount"

	log := logger.FromCtx(ctx).WithFields(logrus.Fields{
		"doc": doc,
	})

	log.Infoln("creating account")

	if doc == "" {
		return "", ErrInvalidDocument
	}

	accID, err := u.accRepo.CreateAccount(ctx, doc)
	if err != nil {
		return "", domain.Error(operation, err)
	}

	log.WithField("id", accID).Infoln("account successfully created")

	return accID, nil
}
