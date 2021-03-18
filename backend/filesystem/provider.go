package filesystem

import (
	"context"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
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
	// Do nothing for now...
	return nil
}

// GetUserFilesystem gets a filesystem implemented based on the view of a given user.
func (s ServiceProvider) GetUserFilesystem(ctx context.Context, user cirrus.User) (FilesystemService, error) {
	return FilesystemService{
		user:       user,
		repository: s.repository,
		db:         s.db,
	}, nil
}
