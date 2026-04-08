package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	configureDockerHost()
	os.Exit(m.Run())
}

// configureDockerHost sets up Docker socket path for Colima on macOS.
// This is required for testcontainers to work with Colima.
func configureDockerHost() {
	if os.Getenv("DOCKER_HOST") != "" {
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	colimaSocket := fmt.Sprintf("%s/.colima/default/docker.sock", homeDir)
	if _, err := os.Stat(colimaSocket); err == nil {
		os.Setenv("DOCKER_HOST", fmt.Sprintf("unix://%s", colimaSocket))
		os.Setenv("TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE", "/var/run/docker.sock")
	}
}

type LostUpdateSuite struct {
	suite.Suite
	pgContainer *tcpostgres.PostgresContainer
	db          *sql.DB
}

func TestLostUpdateSuite(t *testing.T) {
	suite.Run(t, new(LostUpdateSuite))
}

func (s *LostUpdateSuite) SetupSuite() {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("demo"),
		tcpostgres.WithUsername("acid"),
		tcpostgres.WithPassword("strong_password_123"),
		tcpostgres.BasicWaitStrategies(),
	)
	s.Require().NoError(err)
	s.pgContainer = pgContainer

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err)

	db, err := sql.Open("postgres", connStr)
	s.Require().NoError(err)
	s.Require().NoError(db.Ping())
	s.db = db
}

func (s *LostUpdateSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.pgContainer != nil {
		s.pgContainer.Terminate(context.Background()) //nolint:errcheck
	}
}

