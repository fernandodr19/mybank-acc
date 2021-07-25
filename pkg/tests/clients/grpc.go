package clients

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	usecase "github.com/fernandodr19/mybank-acc/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/grpc/accounts"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FakeClient gRPC of accounts
type FakeClient struct {
	client accounts.AccountsServiceClient
}

// NewFakeAccountslient returns a gRPC client
func NewFakeAccountslient(conn *grpc.ClientConn) *FakeClient {
	return &FakeClient{
		client: accounts.NewAccountsServiceClient(conn),
	}

}

// Deposit requests a deposit to the accounts server
func (c FakeClient) Deposit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Client.Deposit"
	_, err := c.client.Deposit(ctx, &accounts.Request{
		AccountID: accID.String(),
		Amount:    amount.Int64(),
	})
	if err != nil {
		return parseServerErr(operation, err)
	}
	return nil
}

// Withdrawal requests a withdrawal to the accounts server
func (c FakeClient) Withdrawal(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Client.Withdrawal"
	_, err := c.client.Withdrawal(ctx, &accounts.Request{
		AccountID: accID.String(),
		Amount:    amount.Int64(),
	})
	if err != nil {
		return parseServerErr(operation, err)
	}
	return nil
}

// ReserveCreditLimit requests a credit limit reserval to the accounts server
func (c FakeClient) ReserveCreditLimit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "accounts.Client.ReserveCreditLimit"
	_, err := c.client.ReserveCreditLimit(ctx, &accounts.Request{
		AccountID: accID.String(),
		Amount:    amount.Int64(),
	})
	if err != nil {
		return parseServerErr(operation, err)
	}
	return nil
}

func parseServerErr(operation string, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return domain.Error(operation, err)
	}
	//nolint
	switch st.Code() {
	case codes.NotFound:
		return usecase.ErrAccountNotFound
	case codes.InvalidArgument:
		switch st.Message() {
		case "err::insufficient_balance":
			return usecase.ErrInsufficientBalance
		case "err::insufficient_credit":
			return usecase.ErrInsufficientCredit
		case "err::invalid_amount":
			return usecase.ErrInvalidAmount
		}
	}

	return domain.Error(operation, err)
}
