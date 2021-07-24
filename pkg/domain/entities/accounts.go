package entities

import (
	"time"

	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
)

// Account entity
type Account struct {
	ID              vos.AccountID
	Document        vos.Document
	Balance         vos.Money
	AvailableCredit vos.Money
	CreatedAt       time.Time
	UpdateAt        time.Time
}

func NewAccount(doc vos.Document, balance vos.Money, AvailableCredit vos.Money) Account {
	return Account{
		Document:        doc,
		Balance:         balance,
		AvailableCredit: AvailableCredit,
	}
}
