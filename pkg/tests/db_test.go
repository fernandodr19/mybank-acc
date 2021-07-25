package tests

import (
	"context"
	"testing"

	"github.com/fernandodr19/mybank-acc/pkg/domain/entities"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CreateAccount_DB(t *testing.T) {
	testTable := []struct {
		Name            string
		Document        vos.Document
		Balance         vos.Money
		AvailableCredit vos.Money
		ExpectedError   bool
	}{
		{
			Name:            "create acc happy path",
			Document:        "123",
			Balance:         0,
			AvailableCredit: 5000,
		},
	}

	ctx := context.Background()

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			defer truncatePostgresTables()

			// test
			accID, err := testEnv.AccRepo.CreateAccount(ctx, entities.NewAccount(tt.Document, tt.Balance, tt.AvailableCredit))
			if tt.ExpectedError {
				assert.Error(t, err)
				return
			}

			// assert
			assert.NoError(t, err)
			_, err = uuid.Parse(accID.String())
			assert.NoError(t, err)
		})
	}
}
