package postgres

import (
	"context"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/entities"
	"github.com/fernandodr19/mybank-acc/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/db/postgres/sqlc"
	"github.com/jackc/pgx/v4"
	pgx_errors "github.com/jackc/pgx/v4"
)

var _ accounts.Repository = &AccountsRepository{}

// AccountsRepository is the repository of accounts
type AccountsRepository struct {
	conn *pgx.Conn
	q    *sqlc.Queries
}

// NewAccountsRepository returns an acc repository
func NewAccountsRepository(conn *pgx.Conn) *AccountsRepository {
	return &AccountsRepository{
		conn: conn,
		q:    sqlc.New(conn),
	}
}

// CreateAccount inserts an account on DB returning its ID
func (r AccountsRepository) CreateAccount(ctx context.Context, acc entities.Account) (vos.AccountID, error) {
	const operation = "postgres.AccountsRepository.CreateAccount"

	accID, err := r.q.CreateAccount(ctx, sqlc.CreateAccountParams{
		Document:        acc.Document.String(),
		Balance:         acc.Balance.Int64(),
		AvailableCredit: acc.AvailableCredit.Int64(),
	})
	if err != nil {
		return "", domain.Error(operation, err)
	}

	return vos.AccountID(accID), nil
}

// GetAccountByID retrieves an account by ID
func (r AccountsRepository) GetAccountByID(ctx context.Context, accID vos.AccountID) (entities.Account, error) {
	const operation = "postgres.AccountsRepository.GetAccountByID"

	rawAcc, err := r.q.GetAccountByID(ctx, accID.String())
	if err != nil {
		if err == pgx_errors.ErrNoRows {
			return entities.Account{}, accounts.ErrAccountNotFound
		}
		return entities.Account{}, domain.Error(operation, err)
	}

	return mapRawAccount(rawAcc), nil
}

// Deposit increments account balance
func (r AccountsRepository) Deposit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "postgres.AccountsRepository.Deposit"

	rows, err := r.q.Deposit(ctx, sqlc.DepositParams{
		ID:     accID.String(),
		Amount: amount.Int64(),
	})
	if err != nil {
		return domain.Error(operation, err)
	}
	if rows == 0 {
		return accounts.ErrAccountNotFound
	}

	return nil
}

// Withdraw decreases account balance
func (r AccountsRepository) Withdraw(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "postgres.AccountsRepository.Withdraw"

	rows, err := r.q.Withdraw(ctx, sqlc.WithdrawParams{
		ID:     accID.String(),
		Amount: amount.Int64(),
	})
	if err != nil {
		return domain.Error(operation, err)
	}
	if rows == 0 {
		return accounts.ErrAccountNotFound
	}

	return nil
}

// DecreaseAvailableCredit decreases account available credit
func (r AccountsRepository) DecreaseAvailableCredit(ctx context.Context, accID vos.AccountID, amount vos.Money) error {
	const operation = "postgres.AccountsRepository.DecreaseAvailableCredit"

	rows, err := r.q.DecreaseAvailableCredit(ctx, sqlc.DecreaseAvailableCreditParams{
		ID:     accID.String(),
		Amount: amount.Int64(),
	})
	if err != nil {
		return domain.Error(operation, err)
	}
	if rows == 0 {
		return accounts.ErrAccountNotFound
	}

	return nil
}

func mapRawAccount(rawAcc sqlc.Account) entities.Account {
	return entities.Account{
		ID:              vos.AccountID(rawAcc.ID),
		Document:        vos.Document(rawAcc.Document),
		Balance:         vos.Money(rawAcc.Balance),
		AvailableCredit: vos.Money(rawAcc.AvailableCredit),
		CreatedAt:       rawAcc.CreatedAt,
		UpdateAt:        rawAcc.UpdatedAt,
	}
}
