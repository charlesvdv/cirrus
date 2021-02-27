package identity

import (
	"context"
	"time"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
	"github.com/rs/zerolog/log"
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

// UserCreatedCallback defines a callback that is called when the user has been created.
// The transaction passed should be used to create any other persistent changes. Don't commit nor
// rollback the transaction.
// If an error is returned, the user creation will be rollbacked as well.
type UserCreatedCallback func(ctx context.Context, tx database.Tx, user cirrus.User) error

// CreateUser creates a user
func (s UserService) CreateUser(ctx context.Context, user *cirrus.User, callback UserCreatedCallback) (err error) {
	user.CreatedAt = time.Now().UTC()

	err = s.db.WithTransaction(ctx, func(tx database.Tx) error {
		err = s.repository.CreateUser(tx, user)
		if err != nil {
			return err
		}

		return callback(ctx, tx, *user)
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create user")
		return cirrus.Errorf(cirrus.ErrCodeInternal, "create user")
	}

	return nil
}
