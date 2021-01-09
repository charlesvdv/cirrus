package identity

import (
	"errors"

	"github.com/charlesvdv/cirrus/backend/db"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type UserRepository interface {
	Create(ctx context.Context, tx db.Tx, user User) (User, error)
}

func NewUserService(txProvider db.TxProvider, repository UserRepository) UserService {
	return UserService{
		txProvider: txProvider,
		repository: repository,
	}
}

type UserService struct {
	txProvider db.TxProvider
	repository UserRepository
}

type SignupInfo struct {
	Email    string
	Password string
}

func (s *UserService) Signup(ctx context.Context, info SignupInfo) error {
	user, err := newSignupUser(info.Email, info.Password)
	if err != nil {
		log.Ctx(ctx).Debug().Err(err).Msg("Invalid new user")
		return ErrInvalidUsernameOrPassword
	}

	err = s.txProvider.WithTransaction(ctx, func(tx db.Tx) error {
		_, err = s.repository.Create(ctx, tx, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		var dupErr db.DuplicateError
		if errors.As(err, &dupErr) {
			return ErrUserAlreadyExists
		}
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to create user")
		return ErrInternal
	}

	return nil
}
