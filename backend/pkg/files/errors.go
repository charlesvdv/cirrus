package files

import "errors"

var (
	ErrNotImplemented      = errors.New("Not implemented")
	ErrRecordNotFound      = errors.New("Record not found")
	ErrInternalServerError = errors.New("Internal server error")
)
