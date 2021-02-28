package filesystem

import (
	"context"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
	"github.com/rs/zerolog/log"
)

// ServiceProvider implements the high level service that will be used to create a filesystem.
type ServiceProvider struct {
	db         database.TxUtils
	repository Repository
}

// NewServiceProvider creates a `ServiceProvider` with all of its dependencies
func NewServiceProvider(db database.TxProvider, repository Repository) ServiceProvider {
	return ServiceProvider{
		db:         database.WrapTxProvider(db),
		repository: repository,
	}
}

// InitUserFilesystem is a callback function used when a user is created.
func (s ServiceProvider) InitUserFilesystem(ctx context.Context, tx database.Tx, user cirrus.User) error {
	err := s.repository.CreateFilesystem(tx, user.ID)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("init filesystem failed")
		return err
	}
	return nil
}
