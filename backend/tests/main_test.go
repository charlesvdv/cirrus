package tests

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/charlesvdv/cirrus/backend/pkg/db/postgres"
	"github.com/ory/dockertest"
)

var (
	testDB *sql.DB
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}

	postgresResource, err := setupPostgres(pool)
	if err != nil {
		log.Fatalf("Unable to start postgres: %v", err)
	}

	code := m.Run()

	if err := pool.Purge(postgresResource); err != nil {
		log.Fatalf("Unable cleanup postgres container: %v", err)
	}

	os.Exit(code)
}

func setupPostgres(pool *dockertest.Pool) (*dockertest.Resource, error) {
	connConfig := postgres.Config{
		DBName:   "postgres",
		User:     "postgres",
		Password: "password",
		Host:     "localhost",
		SslMode:  "disable",
	}

	containerEnvs := []string{"POSTGRES_PASSWORD=" + connConfig.Password, "POSTGRES_DB=" + connConfig.DBName}
	resource, err := pool.Run("postgres", "alpine", containerEnvs)
	if err != nil {
		return nil, err
	}

	publishedPort, err := strconv.Atoi(resource.GetPort("5432/tcp"))
	if err != nil {
		return nil, err
	}
	connConfig.Port = uint16(publishedPort)

	err = pool.Retry(func() error {
		db, err := sql.Open("postgres", postgres.FormatConnectionString(connConfig))
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		return nil, err
	}

	db, err := postgres.Open(connConfig)
	if err != nil {
		return nil, err
	}
	testDB = db

	return resource, nil
}
