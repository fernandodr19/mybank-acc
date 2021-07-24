package grpc

import (
	"context"
	"errors"

	usecase "github.com/fernandodr19/mybank-acc/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/mybank-acc/pkg/domain/vos"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/grpc/accounts"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Usecase interface for accoutns usecases
type Usecase interface {
	Deposit(ctx context.Context, accID vos.AccountID, amount vos.Money) error
	Withdraw(ctx context.Context, accID vos.AccountID, amount vos.Money) error
	ReserveCreditLimit(ctx context.Context, accID vos.AccountID, amount vos.Money) error
}

// Server grpc
type Server struct {
	Usecase
	accounts.UnimplementedAccountsServiceServer
}

// BuildHandler builds grpc handler
func BuildHandler(accUsecase Usecase) *grpc.Server {
	s := Server{
		Usecase: accUsecase,
	}
	grpcServer := grpc.NewServer()
	accounts.RegisterAccountsServiceServer(grpcServer, &s)
	return grpcServer
}

// Deposit handles deposit requests
func (s *Server) Deposit(ctx context.Context, req *accounts.Request) (*accounts.Response, error) {
	err := s.Usecase.Deposit(ctx, vos.AccountID(req.AccountID), vos.Money(req.Amount))
	return &accounts.Response{}, errorResponse(ctx, err)
}

// Withdrawal handles withdrawals requests
func (s *Server) Withdrawal(ctx context.Context, req *accounts.Request) (*accounts.Response, error) {
	err := s.Usecase.Withdraw(ctx, vos.AccountID(req.AccountID), vos.Money(req.Amount))
	return &accounts.Response{}, errorResponse(ctx, err)
}

// ReserveCreditLimit handles reserve credit limit requests
func (s *Server) ReserveCreditLimit(ctx context.Context, req *accounts.Request) (*accounts.Response, error) {
	err := s.Usecase.ReserveCreditLimit(ctx, vos.AccountID(req.AccountID), vos.Money(req.Amount))
	return &accounts.Response{}, errorResponse(ctx, err)
}

var (
	ErrAcountNotFound      = status.New(codes.NotFound, "err::account_not_found").Err()
	ErrInvalidAmount       = status.New(codes.InvalidArgument, "err::invalid_amount").Err()
	ErrInsufficientBalance = status.New(codes.InvalidArgument, "err::insufficient_balance").Err()
	ErrInsufficientCredit  = status.New(codes.InvalidArgument, "err::insufficient_credit").Err()
	ErrUnknown             = status.New(codes.Unknown, "err::unknown").Err()
)

// errorResponse maps response error
func errorResponse(ctx context.Context, err error) error {
	logger.FromCtx(ctx).Errorln(err)
	switch {
	case errors.Is(err, usecase.ErrAccountNotFound):
		return ErrAcountNotFound
	case errors.Is(err, usecase.ErrInvalidAmount):
		return ErrInvalidAmount
	case errors.Is(err, usecase.ErrInsufficientBalance):
		return ErrInsufficientBalance
	case errors.Is(err, usecase.ErrInsufficientCredit):
		return ErrInsufficientCredit
	default:
		return ErrUnknown
	}
}
