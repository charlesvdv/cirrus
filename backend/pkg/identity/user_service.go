package identity

import (
	"github.com/charlesvdv/cirrus/backend/db"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type UserRepository interface {
	Create(ctx context.Context, tx db.Tx, user User) (User, error)
	GetUserWithID(ctx context.Context, tx db.Tx, userID UserID) (User, error)
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
		if db.IsErrDuplicate(err) {
			return ErrUserAlreadyExists
		}
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to create user")
		return ErrInternal
	}

	return nil
}

func (s *UserService) GetUser(ctx context.Context, userID UserID) (User, error) {
	var user User
	err := s.txProvider.WithTransaction(ctx, func(tx db.Tx) error {
		var err error
		user, err = s.repository.GetUserWithID(ctx, tx, userID)
		if err != nil {
			if db.IsErrNoRows(err) {
				return ErrUserNoLongerExists
			}
			log.Ctx(ctx).Warn().Uint64("user-id", userID).Err(err).Msg("Failed to retrieve user")
			return ErrInternal
		}
		return nil
	})

	return user, err
}
