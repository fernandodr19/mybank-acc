package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/fernandodr19/mybank-acc/pkg/domain"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/api/responses"
)

// CreateAccount creates an account
// @Summary Creates an account
// @Description Creates an bank account
// @Tags Accounts
// @Param Body body CreateAccountRequest true "Body"
// @Accept json
// @Produce json
// @Success 201 {object} CreateAccountResponse
// @Failure 400 "Could not parse request"
// @Failure 409 "Account already registered"
// @Failure 422 "Could not create account"
// @Failure 500 "Internal server error"
// @Router /accounts [post]
func (h Handler) CreateAccount(r *http.Request) responses.Response {
	operation := "accounts.Handler.CreateAccount"

	ctx := r.Context()
	var body CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return responses.BadRequest(domain.Error(operation, err), responses.ErrInvalidBody)
	}

	accID, err := h.Usecase.CreateAccount(ctx, body.Document)
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	return responses.Created(CreateAccountResponse{
		AccountID: accID,
	})
}

// CreateAccountRequest payload
type CreateAccountRequest struct {
	Document    vos.Document `json:"document_number"`
	CreditLimit vos.Money    `json:"credit_limit"`
}

// CreateAccountResponse payload
type CreateAccountResponse struct {
	AccountID vos.AccountID `json:"account_id"`
}
