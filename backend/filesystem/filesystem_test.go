package filesystem_test

import (
	"context"
	"testing"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database/sqlite"
	"github.com/charlesvdv/cirrus/backend/filesystem"
	"github.com/charlesvdv/cirrus/backend/identity"
	"github.com/stretchr/testify/require"
)

func TestFilesystem_InitUserFilesystem(t *testing.T) {
	identityRepository := sqlite.IdentityRepository{}
	repository := sqlite.FilesystemRepository{}
	db := sqlite.NewTestDatabase()
	defer db.Close()

	userService := identity.NewUserService(db, identityRepository)
	service := filesystem.NewServiceProvider(db, repository)

	var user cirrus.User
	err := userService.CreateUser(context.Background(), &user, service.InitUserFilesystem)
	require.NoError(t, err)
}
