package tests

import (
	"context"
	"testing"

	"github.com/fernandodr19/mybank-acc/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Deposit(t *testing.T) {
	ctx := context.Background()
	testTable := []struct {
		Name            string
		AccID           vos.AccountID
		Amount          vos.Money
		Setup           func() (vos.AccountID, error)
		ExpectedError   error
		ExpectedBalance vos.Money
	}{
		{
			Name:          "expected invalid amount",
			AccID:         "e031a99d-6191-4d02-8616-b5e3530caccb",
			Amount:        -10,
			ExpectedError: accounts.ErrInvalidAmount,
		},
		{
			Name:          "expected invalid acc id",
			AccID:         "24dde2d4-5763-419d-9a93-3365ef55255c",
			Amount:        10,
			ExpectedError: accounts.ErrAccountNotFound,
		},
		{
			Name: "deposit happy path",
			Setup: func() (vos.AccountID, error) {
				return testEnv.App.Accounts.CreateAccount(ctx, "123", 0)
			},
			Amount:          10,
			ExpectedBalance: 10,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			defer truncatePostgresTables()

			// prepare
			if tt.Setup != nil {
				accID, err := tt.Setup()
				require.NoError(t, err)
				tt.AccID = accID
			}

			// test
			err := testEnv.GrpcFakeClient.Deposit(ctx, tt.AccID, tt.Amount)

			// assert
			assert.ErrorIs(t, tt.ExpectedError, err)

			if err != nil {
				return
			}

			acc, err := testEnv.App.Accounts.GetAccountByID(ctx, tt.AccID)
			require.NoError(t, err)
			assert.Equal(t, tt.ExpectedBalance, acc.Balance)
		})
	}
}

func Test_Withdraw(t *testing.T) {
	ctx := context.Background()
	testTable := []struct {
		Name            string
		AccID           vos.AccountID
		Amount          vos.Money
		Setup           func(t *testing.T) vos.AccountID
		ExpectedError   error
		ExpectedBalance vos.Money
	}{
		{
			Name:          "expected invalid amount",
			AccID:         "e031a99d-6191-4d02-8616-b5e3530caccb",
			Amount:        -10,
			ExpectedError: accounts.ErrInvalidAmount,
		},
		{
			Name:          "expected invalid acc id",
			AccID:         "24dde2d4-5763-419d-9a93-3365ef55255c",
			Amount:        10,
			ExpectedError: accounts.ErrAccountNotFound,
		},
		{
			Name: "insufficient balance",
			Setup: func(t *testing.T) vos.AccountID {
				accID, err := testEnv.App.Accounts.CreateAccount(ctx, "123", 0)
				require.NoError(t, err)

				return accID
			},
			Amount:        10,
			ExpectedError: accounts.ErrInsufficientBalance,
		},
		{
			Name: "withdraw happy path",
			Setup: func(t *testing.T) vos.AccountID {
				accID, err := testEnv.App.Accounts.CreateAccount(ctx, "123", 0)
				require.NoError(t, err)

				err = testEnv.App.Accounts.Deposit(ctx, accID, 10)
				require.NoError(t, err)

				return accID
			},
			Amount:          10,
			ExpectedBalance: 0,
		},
		{
			Name: "withdraw happy path rich",
			Setup: func(t *testing.T) vos.AccountID {
				accID, err := testEnv.App.Accounts.CreateAccount(ctx, "123", 0)
				require.NoError(t, err)

				err = testEnv.App.Accounts.Deposit(ctx, accID, 200000000)
				require.NoError(t, err)

				return accID
			},
			Amount:          10,
			ExpectedBalance: 199999990,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			defer truncatePostgresTables()

			// prepare
			if tt.Setup != nil {
				tt.AccID = tt.Setup(t)
			}

			// test
			err := testEnv.GrpcFakeClient.Withdrawal(ctx, tt.AccID, tt.Amount)

			// assert
			assert.ErrorIs(t, tt.ExpectedError, err)

			if err != nil {
				return
			}

			acc, err := testEnv.App.Accounts.GetAccountByID(ctx, tt.AccID)
			require.NoError(t, err)
			assert.Equal(t, tt.ExpectedBalance, acc.Balance)
		})
	}
}

func Test_Credit(t *testing.T) {
	ctx := context.Background()
	testTable := []struct {
		Name                    string
		AccID                   vos.AccountID
		Amount                  vos.Money
		Setup                   func() (vos.AccountID, error)
		ExpectedError           error
		ExpectedAvailableCredit vos.Money
	}{
		{
			Name:          "expected invalid amount",
			AccID:         "e031a99d-6191-4d02-8616-b5e3530caccb",
			Amount:        -10,
			ExpectedError: accounts.ErrInvalidAmount,
		},
		{
			Name:          "expected invalid acc id",
			AccID:         "24dde2d4-5763-419d-9a93-3365ef55255c",
			Amount:        10,
			ExpectedError: accounts.ErrAccountNotFound,
		},
		{
			Name: "insufficient credit limit",
			Setup: func() (vos.AccountID, error) {
				return testEnv.App.Accounts.CreateAccount(ctx, "123", 10)
			},
			Amount:                  11,
			ExpectedError:           accounts.ErrInsufficientCredit,
			ExpectedAvailableCredit: 10,
		},
		{
			Name: "credit happy path",
			Setup: func() (vos.AccountID, error) {
				return testEnv.App.Accounts.CreateAccount(ctx, "123", 10)
			},
			Amount: 10,
		},
		{
			Name: "credit happy path rich",
			Setup: func() (vos.AccountID, error) {
				return testEnv.App.Accounts.CreateAccount(ctx, "123", 999999990)
			},
			Amount:                  70,
			ExpectedAvailableCredit: 999999920,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			defer truncatePostgresTables()

			// prepare
			if tt.Setup != nil {
				accID, err := tt.Setup()
				require.NoError(t, err)
				tt.AccID = accID
			}

			// test
			err := testEnv.GrpcFakeClient.ReserveCreditLimit(ctx, tt.AccID, tt.Amount)

			// assert
			assert.ErrorIs(t, tt.ExpectedError, err)

			if err != nil {
				return
			}

			acc, err := testEnv.App.Accounts.GetAccountByID(ctx, tt.AccID)
			require.NoError(t, err)
			assert.Equal(t, tt.ExpectedAvailableCredit, acc.AvailableCredit)
		})
	}
}
