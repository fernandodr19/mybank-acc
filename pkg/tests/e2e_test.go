package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/fernandodr19/mybank-acc/pkg/gateway/api/accounts"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccount(t *testing.T) {
	testTable := []struct {
		Name               string
		Req                accounts.CreateAccountRequest
		Setup              func(t *testing.T)
		ExpectedStatusCode int
	}{
		{
			Name:               "bad request: invalid body",
			Req:                accounts.CreateAccountRequest{},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "unprocessable entity: invalid document",
			Req:                accounts.CreateAccountRequest{Document: "invalid"},
			ExpectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:               "unprocessable entity: invalid credit limit",
			Req:                accounts.CreateAccountRequest{Document: "123", CreditLimit: -1},
			ExpectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name: "conflict: account alread registered",
			Req:  accounts.CreateAccountRequest{Document: "999", CreditLimit: 10},
			Setup: func(t *testing.T) {
				// create account to simulate duplication
				_, err := testEnv.App.Accounts.CreateAccount(context.Background(), "999", 0)
				require.NoError(t, err)

			},
			ExpectedStatusCode: http.StatusConflict,
		},
		{
			Name:               "create account happy path",
			Req:                accounts.CreateAccountRequest{Document: "123", CreditLimit: 5000},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			Name:               "create account happy path with default credit limit (zero)",
			Req:                accounts.CreateAccountRequest{Document: "123"},
			ExpectedStatusCode: http.StatusCreated,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			defer truncatePostgresTables()

			// prepare
			if tt.Setup != nil {
				tt.Setup(t)
			}

			target := testEnv.ApiServer.URL + "/api/v1/accounts"
			body, err := json.Marshal(tt.Req)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
			require.NoError(t, err)

			// test
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// assert
			require.Equal(t, tt.ExpectedStatusCode, resp.StatusCode)
		})
	}
}
