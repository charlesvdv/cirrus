package sqlite_test

import (
	"context"
	"testing"

	"github.com/charlesvdv/cirrus/backend/database/sqlite"
	"github.com/stretchr/testify/require"
)

func TestSqliteMigration(t *testing.T) {
	sqlite.NewTestDatabase()
}

func TestSqliteTxCommit(t *testing.T) {
	db := sqlite.NewTestDatabase()

	tx, err := db.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)
}

func TestSqliteTxRollback(t *testing.T) {
	db := sqlite.NewTestDatabase()

	tx, err := db.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)
}
