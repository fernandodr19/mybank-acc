package main

import (
	"context"
	"net"
	"net/http"
	"time"

	app "github.com/fernandodr19/mybank-acc/pkg"
	"github.com/fernandodr19/mybank-acc/pkg/config"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/api"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/db/postgres"
	grpc_acc "github.com/fernandodr19/mybank-acc/pkg/gateway/gRPC"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	log := logger.Default()
	log.Infoln("=== My Bank ACC ===")

	ctx := context.Background()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("failed loading config")
	}

	// Setup postgres
	dbConn, err := postgres.NewConnection(ctx, cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("failed setting up postgres")
	}

	// Build app
	app, err := app.BuildApp(dbConn)
	if err != nil {
		log.WithError(err).Fatalln("failed to build app")
	}

	// Build gRPC handler
	grpcHandler := grpc_acc.BuildHandler(app)

	// Build API handler
	apiHandler := api.BuildHandler(app)

	// Server up application
	serveApp(apiHandler, grpcHandler, cfg, log)
}

func serveApp(apiHandler http.Handler, grpcHandler *grpc.Server, cfg *config.Config, log *logrus.Entry) {
	// gRPC server
	go func() {
		protocol, address := cfg.GRPC.Protocol, cfg.GRPC.Address()
		l, err := net.Listen(protocol, address)
		if err != nil {
			log.WithFields(logrus.Fields{
				"protocol": protocol,
				"address":  address,
			}).WithError(err).Fatalln("failed to listen gRPC")
		}

		log.WithField("address", address).Infoln("gRPC server starting...")
		log.Fatal(grpcHandler.Serve(l))
	}()

	// REST server
	server := &http.Server{
		Handler:      apiHandler,
		Addr:         cfg.API.Address(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("address", cfg.API.Address()).Info("rest server starting...")
	log.Fatal(server.ListenAndServe())
}
