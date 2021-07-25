package tests

import (
	"context"
	"fmt"
	"net"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/fernandodr19/mybank-acc/pkg/config"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/api"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/db/postgres"
	acc_grpc "github.com/fernandodr19/mybank-acc/pkg/gateway/gRPC"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"github.com/fernandodr19/mybank-acc/pkg/tests/clients"
	"google.golang.org/grpc"

	app "github.com/fernandodr19/mybank-acc/pkg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/testcontainers/testcontainers-go"
)

type _testEnv struct {
	// Server
	Server *httptest.Server

	// 3rd party fake Clients
	GrpcFakeClient *clients.FakeClient

	// DB
	Conn    *pgx.Conn
	AccRepo *postgres.AccountsRepository

	// App
	App *app.App
}

var testEnv _testEnv

func TestMain(m *testing.M) {
	teardown := setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func setup() func() {
	log := logger.Default()
	log.Info("setting up integration tests env")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("failed loading config")
	}

	err = setupDockerTest()
	if err != nil {
		log.WithError(err).Fatal("failed setting up docker")
	}

	// Setup postgres
	cfg.Postgres.DBName = "test"
	cfg.Postgres.Port = "5435"
	dbConn, err := setupPostgresTest(cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("failed setting up postgres")
	}

	testEnv.Conn = dbConn
	testEnv.AccRepo = postgres.NewAccountsRepository(dbConn)

	app := app.BuildApp(dbConn)

	testEnv.App = app

	apiHandler := api.BuildHandler(app)

	lis, err := net.Listen(cfg.GRPC.Protocol, cfg.GRPC.Address())
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := acc_grpc.BuildHandler(app)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.WithError(err).Fatalf("Failed to serve gRPC server")
		}
	}()

	// Fake gRPC client
	clintGrpcConn, err := grpc.Dial(cfg.GRPC.Address(), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalln("failed connecting grpc")
	}
	testEnv.GrpcFakeClient = clients.NewFakeAccountslient(clintGrpcConn)

	testEnv.Server = httptest.NewServer(apiHandler)

	return func() {
		clintGrpcConn.Close()
	}
}

func setupDockerTest() error {
	running, err := isDockerRunning([]string{
		"pg-test",
	})
	if err != nil {
		return err
	}

	if running {
		logger.Default().Infoln("necessary containers already running...")
		return nil
	}

	compose := testcontainers.NewLocalDockerCompose(
		[]string{"./docker-compose.yml"},
		strings.ToLower(uuid.New().String()),
	)
	execErr := compose.WithCommand([]string{"up", "-d"}).Invoke()
	if execErr.Error != nil {
		return execErr.Error
	}
	return nil
}

func isDockerRunning(expectedImages []string) (bool, error) {
	stdout, err := exec.Command("docker", "ps").Output()
	if err != nil {
		return false, err
	}

	ps := string(stdout)
	if err != nil {
		return false, err
	}

	running := true
	for _, image := range expectedImages {
		if !strings.Contains(ps, image) {
			running = false
			break
		}
	}
	return running, nil
}

func setupPostgresTest(cfg config.Postgres) (*pgx.Conn, error) {
	done := make(chan bool, 1)
	var dbConn *pgx.Conn
	var err error

	// tries to connect within 5 seconds timeout
	go func() {
		for {
			dbConn, err = postgres.NewConnection(context.Background(), cfg)
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			break
		}
		close(done)
	}()

	select {
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timed out trying to set up postgres: %w", err)
	case <-done:
	}

	return dbConn, nil
}

func truncatePostgresTables() {
	testEnv.Conn.Exec(context.Background(),
		`TRUNCATE TABLE 
			accounts
		CASCADE`,
	)
}
