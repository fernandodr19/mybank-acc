package accounts

import (
	"context"
	"net/http"

	"github.com/fernandodr19/mybank-acc/pkg/domain/entities"
	"github.com/fernandodr19/mybank-acc/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/api/middleware"
	"github.com/gorilla/mux"
)

//go:generate moq -skip-ensure -stub -out mocks.gen.go . Usecase:AccountsMockUsecase

var _ Usecase = accounts.Usecase{}

// Usecase of accoutns
type Usecase interface {
	CreateAccount(ctx context.Context, doc vos.Document, creditLimit vos.Money) (vos.AccountID, error)
	GetAccountByID(ctx context.Context, accID vos.AccountID) (entities.Account, error)
}

// Handler handles account relared REST requests
type Handler struct {
	Usecase
}

// NewHandler builds accounts handler
func NewHandler(public *mux.Router, admin *mux.Router, usecase Usecase) *Handler {
	h := &Handler{
		Usecase: usecase,
	}

	public.Handle("/accounts",
		middleware.Handle(h.CreateAccount)).
		Methods(http.MethodPost)

	public.Handle("/accounts/{account_id}",
		middleware.Handle(h.GetAccount)).
		Methods(http.MethodGet)

	return h
}
