package db

import "context"

type TxProvider interface {
	BeginTx(ctx context.Context) (Tx, error)
	WithTransaction(ctx context.Context, txFn func(Tx) error) error
}

type Tx interface {
	Commit() error
	Rollback() error
}

type TxError struct {
	err error
}

func (e TxError) Error() string {
	return e.err.Error()
}
