package main

import (
	"log"
	"net"

	grpc "github.com/fernandodr19/mybank-acc/pkg/gateway/gRPC"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.BuildHandler(nil)

	log.Println("Starting server...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}

}
