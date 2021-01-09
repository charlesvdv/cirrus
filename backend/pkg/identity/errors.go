package identity

import "errors"

var (
	ErrInternal                  = errors.New("Internal error")
	ErrInvalidUsernameOrPassword = errors.New("Invalid user name or password")
	ErrUserAlreadyExists         = errors.New("User already exists")
)
