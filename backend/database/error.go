package database

import "errors"

var (
	// ErrInternal describes an internal error.
	ErrInternal = errors.New("internal")
	// ErrDuplicate is returned when a unique constraint is throwned.
	ErrDuplicate = errors.New("duplicate")
)
