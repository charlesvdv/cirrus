package postgres

import (
	"database/sql"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
)

const dbPassword = "password"

func TestOpenAndSetupDB(t *testing.T) {
	pool, err := dockertest.NewPool("")
	assert.NoError(t, err)

	connConfig := Config{
		DBName:   "postgres",
		User:     "postgres",
		Password: "password",
		Host:     "localhost",
		sslMode:  "disable",
	}

	containerEnvs := []string{"POSTGRES_PASSWORD=" + connConfig.Password, "POSTGRES_DB=" + connConfig.DBName}
	resource, err := pool.Run("postgres", "alpine", containerEnvs)
	assert.NoError(t, err)
	defer resource.Close()

	publishedPort, err := strconv.Atoi(resource.GetPort("5432/tcp"))
	assert.NoError(t, err)
	connConfig.Port = uint16(publishedPort)

	err = pool.Retry(func() error {
		db, err := sql.Open("postgres", formatConnectionString(connConfig))
		if err != nil {
			return err
		}
		return db.Ping()
	})
	assert.NoError(t, err)

	db, err := Open(connConfig)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
