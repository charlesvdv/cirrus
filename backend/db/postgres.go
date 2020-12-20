package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

func (c *PostgresConfig) format() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", c.User, c.Password, c.Host, c.Port, c.Database)
}

func NewPostgresDatabase(conf PostgresConfig) (PostgresDatabase, error) {
	log.Debug().Str("user", conf.User).Str("host", conf.Host).
		Uint16("port", conf.Port).Str("database", conf.Database).Msg("Init postgres database")

	var pool *pgxpool.Pool
	var err error
	for retryCount := 0; retryCount < 5; retryCount++ {
		pool, err = pgxpool.Connect(context.Background(), conf.format())
		if err != nil {
			log.Debug().Int("retry", retryCount).Err(err).Msg("Retry to connect to database failed")
			time.Sleep(time.Second)
			continue
		}
	}

	if err != nil {
		return PostgresDatabase{}, err
	}

	database := PostgresDatabase{
		pool: pool,
	}

	return database, nil
}

type PostgresDatabase struct {
	pool *pgxpool.Pool
}

func (db *PostgresDatabase) UpdateSchemas() error {
	// TODO: do proper migration script
	_, err := db.pool.Exec(context.Background(), postgresSchema)
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	log.Debug().Msg("Schema updated successfully")
	return nil
}

func (db *PostgresDatabase) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return PostgresTx{}, err
	}

	return PostgresTx{
		context: ctx,
		Tx:      tx,
	}, nil
}

type PostgresTx struct {
	pgx.Tx
	context context.Context
}

func (tx PostgresTx) Commit() error {
	return tx.Tx.Commit(tx.context)
}

func (tx PostgresTx) Rollback() error {
	return tx.Tx.Rollback(tx.context)
}
