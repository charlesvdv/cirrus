package sqlite

import (
	"context"

	"github.com/charlesvdv/cirrus/backend/database"
	"github.com/rs/zerolog/log"

	driver "crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

const schema = `
CREATE TABLE users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at TEXT NOT NULL
);
`

// NewTestDatabase creates a new test database. Panics if something goes wrong
// while initializing the database. No persistance is garanteed.
func NewTestDatabase() Database {
	dbpool, err := sqlitex.Open("file:memory?mode=memory&cache=shared", 0, 10)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open test sqlite db")
	}

	db, err := newDatabaseFromPool(dbpool)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create db")
	}

	return db
}

func newDatabaseFromPool(dbpool *sqlitex.Pool) (Database, error) {
	db := Database{
		pool: dbpool,
	}

	err := db.migrate()
	if err != nil {
		log.Err(err).Msg("migration failed")
		return db, err
	}

	return db, nil
}

// Database provides a sqlite database wrapper.
type Database struct {
	pool *sqlitex.Pool
}

func (db Database) migrate() error {
	// TODO: use proper migration script...
	conn := db.pool.Get(context.Background())
	if conn == nil {
		log.Error().Msg("failed to get connection")
		return database.ErrInternal
	}
	defer db.pool.Put(conn)

	err := sqlitex.ExecScript(conn, schema)
	if err != nil {
		log.Err(err).Msg("failed to migrate schema")
		return database.ErrInternal
	}

	return nil
}

// BeginTx starts a transaction on the database.
func (db Database) BeginTx(ctx context.Context) (database.Tx, error) {
	conn := db.pool.Get(ctx)
	if conn == nil {
		log.Ctx(ctx).Error().Msg("failed to get database connection")
		return Tx{}, database.ErrInternal
	}

	stmt := conn.Prep("BEGIN TRANSACTION;")
	_, err := stmt.Step()
	if err != nil {
		log.Ctx(ctx).Error().Msg("failed to start transaction")
		return Tx{}, database.ErrInternal
	}

	return Tx{
		Conn: conn,
		ctx:  ctx,
	}, nil
}

// Tx implements a transaction for the sqlite package.
type Tx struct {
	*driver.Conn
	ctx context.Context
}

// Commit implements commit for sqlite package
func (tx Tx) Commit() error {
	log.Ctx(tx.ctx).Debug().Msg("commit tx")

	stmt := tx.Prep("COMMIT;")
	_, err := stmt.Step()
	if err != nil {
		log.Ctx(tx.ctx).Err(err).Msg("failed to commit tx")
		return database.ErrInternal
	}

	return nil
}

// Rollback implements rollback for sqlite package
func (tx Tx) Rollback() error {
	log.Ctx(tx.ctx).Debug().Msg("rollback tx")

	stmt := tx.Prep("ROLLBACK;")
	_, err := stmt.Step()
	if err != nil {
		log.Ctx(tx.ctx).Err(err).Msg("failed to rollback tx")
		return database.ErrInternal
	}

	return nil
}

func getTx(tx database.Tx) Tx {
	return tx.(Tx)
}

func formatError(err error) error {
	if err == nil {
		return nil
	}

	return database.ErrInternal
}
