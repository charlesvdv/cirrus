package sqlite_test

import (
	"context"
	"testing"

	"github.com/charlesvdv/cirrus/backend/database/sqlite"
	"github.com/stretchr/testify/require"
)

func TestMigration(t *testing.T) {
	sqlite.NewTestDatabase()
}

func TestTxCommit(t *testing.T) {
	db := sqlite.NewTestDatabase()
	defer db.Close()

	tx, err := db.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)
}

func TestTxRollbackAfterCommit(t *testing.T) {
	db := sqlite.NewTestDatabase()
	defer db.Close()

	tx, err := db.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)
}

func TestTxCommitAfterRollback(t *testing.T) {
	db := sqlite.NewTestDatabase()
	defer db.Close()

	tx, err := db.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)

	err = tx.Commit()
	require.Error(t, err)
}

func TestTxRollback(t *testing.T) {
	db := sqlite.NewTestDatabase()
	defer db.Close()

	tx, err := db.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)
}
