package identity

import (
	"github.com/charlesvdv/cirrus/backend/db"
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
		return err
	}

	tx, err := s.txProvider.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = s.repository.Create(ctx, tx, user)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
