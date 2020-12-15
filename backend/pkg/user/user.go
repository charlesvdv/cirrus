package user

import (
	"errors"

	"golang.org/x/net/context"
)

type ID = uint64

type User struct {
}

func NewUserManager() UserManager {
	return UserManager{}
}

type UserManager struct {
}

type SignupInfo struct {
	Email    string
	Password string
}

func (m *UserManager) Signup(ctx context.Context, info SignupInfo) error {
	return errors.New("unimplemented")
}
