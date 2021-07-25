package accounts

import (
	"context"
	"strconv"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/entities"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"github.com/sirupsen/logrus"
)

// CreateAccount creates an account
func (u Usecase) CreateAccount(ctx context.Context, doc vos.Document, creditLimit vos.Money) (vos.AccountID, error) {
	const operation = "accounts.Usecase.CreateAccount"

	log := logger.FromCtx(ctx).WithFields(logrus.Fields{
		"doc": doc,
	})

	log.Infoln("creating account")

	_, err := strconv.Atoi(doc.String())
	if err != nil {
		return "", ErrInvalidDocument
	}

	if creditLimit < 0 {
		return "", ErrInvalidCreditLimit
	}

	acc := entities.NewAccount(doc, 0, creditLimit)
	accID, err := u.accRepo.CreateAccount(ctx, acc)
	if err != nil {
		return "", domain.Error(operation, err)
	}

	log.WithField("id", accID).Infoln("account successfully created")

	return accID, nil
}
