package files

import "errors"

var (
	ErrNotImplemented          = errors.New("Not implemented")
	ErrRecordNotFound          = errors.New("Record not found")
	ErrInternalServerError     = errors.New("Internal server error")
	ErrDirectoryNameEmpty      = errors.New("Directory name can not be empty")
	ErrParentDirectoryNotFound = errors.New("Parent directory not found")
	ErrInvalidID               = errors.New("Invalid ID")
	ErrParentNotDirectory      = errors.New("Parent is not a directory")
)
