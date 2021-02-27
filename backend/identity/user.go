package identity

import (
	"context"
	"time"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
)

// UserService implements the user management with persistence
type UserService struct {
	db         database.TxUtils
	repository Repository
}

// NewUserService creates a UserService with its explicit dependency
func NewUserService(db database.TxProvider, repository Repository) UserService {
	return UserService{
		db:         database.WrapTxProvider(db),
		repository: repository,
	}
}

// CreateUser creates a user
func (s UserService) CreateUser(ctx context.Context, user *cirrus.User) (err error) {
	user.CreatedAt = time.Now().UTC()

	err = s.db.WithTransaction(ctx, func(tx database.Tx) error {
		return s.repository.CreateUser(tx, user)
	})
	if err != nil {
		return cirrus.Errorf(cirrus.ErrCodeInternal, "create user")
	}

	return nil
}
