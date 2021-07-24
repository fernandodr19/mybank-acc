package tests

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/fernandodr19/mybank-acc/pkg/config"
	"github.com/fernandodr19/mybank-acc/pkg/gateway/db/postgres"
	"github.com/fernandodr19/mybank-acc/pkg/instrumentation/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/testcontainers/testcontainers-go"
)

type _testEnv struct {
	Conn *pgx.Conn
}

var testEnv _testEnv

func TestMain(m *testing.M) {
	teardown := setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func setup() func() {

	return func() {}
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
