package responses

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fernandodr19/mybank-acc/pkg/domain/usecases/accounts"
)

// Response represents an API response
type Response struct {
	Status  int
	Error   error
	Payload interface{}
	headers map[string]string
}

// Headers list response headers
func (r *Response) Headers() map[string]string {
	return r.headers
}

// SetHeader set response header
func (r *Response) SetHeader(key, value string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}

	r.headers[key] = value
}

// Error code & description
type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ErrorPayload represents error response payload
type ErrorPayload struct {
	Error `json:"errors"`
}

// shared
var (
	ErrInternalServerError = ErrorPayload{Error: Error{Code: "error:internal_server_error", Description: "Internal Server Error"}}
	ErrInvalidBody         = ErrorPayload{Error: Error{Code: "error:invalid_body", Description: "Invalid body"}}
	ErrInvalidParams       = ErrorPayload{Error: Error{Code: "error:invalid_parameters", Description: "Invalid query parameters"}}
	ErrNotImplemented      = ErrorPayload{Error: Error{Code: "error:not_implemented", Description: "Not implemented"}}
)

// accounts
var (
	ErrAccountNotFound     = ErrorPayload{Error: Error{Code: "error:account_not_found", Description: "Account not found"}}
	ErrAccountConflict     = ErrorPayload{Error: Error{Code: "error:account_already_registered", Description: "Account already registered"}}
	ErrInsufficientBalance = ErrorPayload{Error: Error{Code: "error:insufficient_balance", Description: "Insufficient balance"}}
	ErrInsufficientCredit  = ErrorPayload{Error: Error{Code: "error:insufficient_credit", Description: "Insufficient credit"}}
	ErrInvalidAccID        = ErrorPayload{Error: Error{Code: "error:invalid_account_id", Description: "Account id must be a UUIDv4"}}
	ErrInvalidAmount       = ErrorPayload{Error: Error{Code: "error:invalid_amount", Description: "Amount must be greater than 0"}}
	ErrInvalidCreditLimit  = ErrorPayload{Error: Error{Code: "error:invalid_credit_limit", Description: "Credit limit must be greater than 0"}}
	ErrInvalidDocument     = ErrorPayload{Error: Error{Code: "error:invalid_document", Description: "Invalid document"}}
)

// ErrorResponse maps response error
func ErrorResponse(err error) Response {
	switch {
	case errors.Is(err, accounts.ErrInvalidDocument):
		return UnprocessableEntity(err, ErrInvalidDocument)
	case errors.Is(err, accounts.ErrInvalidCreditLimit):
		return UnprocessableEntity(err, ErrInvalidCreditLimit)
	case errors.Is(err, accounts.ErrAccountConflict):
		return Conflict(err, ErrAccountConflict)
	case errors.Is(err, accounts.ErrInvalidAmount):
		return UnprocessableEntity(err, ErrInvalidAmount)
	case errors.Is(err, accounts.ErrInsufficientBalance):
		return UnprocessableEntity(err, ErrInsufficientBalance)
	case errors.Is(err, accounts.ErrInsufficientCredit):
		return UnprocessableEntity(err, ErrInsufficientCredit)
	default:
		return InternalServerError(err)
	}
}

// InternalServerError 500
func InternalServerError(err error) Response {
	return Response{
		Status:  http.StatusInternalServerError,
		Error:   err,
		Payload: ErrInternalServerError,
	}
}

// NotImplemented 501
func NotImplemented(err error) Response {
	return Response{
		Status:  http.StatusNotImplemented,
		Error:   err,
		Payload: ErrNotImplemented,
	}
}

// BadRequest 400
func BadRequest(err error, payload ErrorPayload) Response {
	return genericError(http.StatusBadRequest, err, payload)
}

// Unauthorized 401
func Unauthorized(err error, payload ErrorPayload) Response {
	return genericError(http.StatusUnauthorized, err, payload)
}

// NotFound 404
func NotFound(err error, payload ErrorPayload) Response {
	return genericError(http.StatusNotFound, err, payload)
}

// Conflict 409
func Conflict(err error, payload ErrorPayload) Response {
	return genericError(http.StatusConflict, err, payload)
}

// UnprocessableEntity 422
func UnprocessableEntity(err error, payload ErrorPayload) Response {
	return genericError(http.StatusUnprocessableEntity, err, payload)
}

func genericError(status int, err error, payload ErrorPayload) Response {
	return Response{
		Status:  status,
		Error:   err,
		Payload: payload,
	}
}

// NoContent 204
func NoContent() Response {
	return Response{
		Status: http.StatusNoContent,
	}
}

// OK 200
func OK(payload interface{}) Response {
	return Response{
		Status:  http.StatusOK,
		Payload: payload,
	}
}

// Created 201
func Created(payload interface{}) Response {
	return Response{
		Status:  http.StatusCreated,
		Payload: payload,
	}
}

// Accepted 202
func Accepted(payload interface{}) Response {
	return Response{
		Status:  http.StatusAccepted,
		Payload: payload,
	}
}

// SendJSON responds requests based on
func SendJSON(w http.ResponseWriter, payload interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil { // Blank body is not valid JSON.
		return nil
	}

	switch p := payload.(type) {
	case string:
		if p == "" {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(payload)
}
