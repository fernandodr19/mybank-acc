package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/fernandodr19/mybank-acc/pkg/gateway/grpc/accounts"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	accounts.RegisterAccountsServiceServer(grpcServer, &s)

	log.Println("Starting server...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}

}

type Server struct {
	accounts.UnimplementedAccountsServiceServer
}

func (s *Server) Deposit(ctx context.Context, req *accounts.Request) (*accounts.Response, error) {
	fmt.Println("DEPOSIT !!!")
	return &accounts.Response{Success: true}, nil
}

func (s *Server) Withdrawal(ctx context.Context, req *accounts.Request) (*accounts.Response, error) {
	fmt.Println("WITH !!!")
	return &accounts.Response{Success: true}, nil
}

func (s *Server) ReserveCreditLimit(ctx context.Context, req *accounts.Request) (*accounts.Response, error) {
	fmt.Println("RESERVE !!!")
	return &accounts.Response{Success: true}, nil
}
