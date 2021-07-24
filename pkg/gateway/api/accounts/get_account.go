package accounts

import (
	"net/http"
	"time"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/api/responses"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetAccount gets an account
// @Summary Gets an account
// @Description Retrieve an account by its ID
// @Tags Accounts
// @Param account_id path string true "Account ID"
// @Accept json
// @Produce json
// @Success 200 {object} GetAccountResponse
// @Failure 400 "Could not parse request"
// @Failure 404 "Account not found"
// @Failure 422 "Could not create account"
// @Failure 500 "Internal server error"
// @Router /accounts/{account_id} [get]
func (h Handler) GetAccount(r *http.Request) responses.Response {
	operation := "accounts.Handler.GetAccount"

	ctx := r.Context()
	accID, err := uuid.Parse(mux.Vars(r)["account_id"])
	if err != nil {
		return responses.BadRequest(domain.Error(operation, err), responses.ErrInvalidAccID)
	}

	acc, err := h.Usecase.GetAccountByID(ctx, vos.AccountID(accID.String()))
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	return responses.OK(GetAccountResponse{
		ID:              acc.ID,
		Document:        acc.Document,
		Balance:         acc.Balance,
		AvailableCredit: acc.AvailableCredit,
		CreatedAt:       acc.CreatedAt,
		UpdateAt:        acc.UpdateAt,
	})
}

// GetAccountResponse payload
type GetAccountResponse struct {
	ID              vos.AccountID `json:"account_id"`
	Document        vos.Document  `json:"document_number"`
	Balance         vos.Money     `json:"balance"`
	AvailableCredit vos.Money     `json:"available_credit_limit"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdateAt        time.Time     `json:"updated_at"`
}
