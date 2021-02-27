package database

import (
	"context"

	"github.com/rs/zerolog/log"
)

// TxFn is function executed in the scope of a transaction.
// If error is nil, the transaction is commited. Otherwise,
// the transaction is rollbacked.
type TxFn func(Tx) error

// WrapTxProvider wraps the TxProvider with other utilities function.
func WrapTxProvider(provider TxProvider) TxUtils {
	return TxUtils{
		TxProvider: provider,
	}
}

// TxUtils provides safe utility function to work with transaction.
type TxUtils struct {
	TxProvider
}

// WithTransaction gives a safe way for the user to execute database code inside
// a transaction. If fn returns nil, the transaction is commited. Otherwise,
// the transaction is rollbacked.
func (utils TxUtils) WithTransaction(ctx context.Context, fn TxFn) (err error) {
	tx, err := utils.BeginTx(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("WithTransaction begin tx failed")
		return ErrInternal
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("WithTransaction commit failed")
		return ErrInternal
	}

	return nil
}

// TxProvider provides a database-independant way to create a transaction.
//
// The context provided is used to commit or rollback the transaction.
// If the context is cancelled, the transaction will rollback.
type TxProvider interface {
	BeginTx(ctx context.Context) (Tx, error)
}

// Tx exposes a database-independant transaction interface.
//
// Database-dependent code must cast this transaction interface
// to the database implementation specific object.
type Tx interface {
	// Commit the transaction.
	Commit() error

	// Rollback aborts the transaction. If a transaction already has been
	// commited, do nothing.
	// This behavior defers from `database/sql` package.
	Rollback() error
}
