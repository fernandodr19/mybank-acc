package app

import (
	"github.com/fernandodr19/mybank-acc/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/db/postgres"

	"github.com/jackc/pgx/v4"
)

// App contains application's usecases
type App struct {
	Accounts *accounts.Usecase
}

// BuildApp builds application struct with its necessary usecases
func BuildApp(dbConn *pgx.Conn) *App {
	accRepo := postgres.NewAccountsRepository(dbConn)
	return &App{
		Accounts: accounts.NewUsecase(accRepo),
	}
}
