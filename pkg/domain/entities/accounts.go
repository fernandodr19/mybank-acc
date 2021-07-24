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
