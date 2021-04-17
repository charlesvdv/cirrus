package filesystem_test

import (
	"context"
	"os"
	"testing"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database/sqlite"
	"github.com/charlesvdv/cirrus/backend/filesystem"
	"github.com/charlesvdv/cirrus/backend/identity"
	"github.com/stretchr/testify/require"
)

var db sqlite.Database
var userService identity.UserService
var fsProvider filesystem.ServiceProvider

func getNewFs(t *testing.T) filesystem.FilesystemService {
	var user cirrus.User
	err := userService.CreateUser(context.Background(), &user, fsProvider.InitUserFilesystem)
	require.NoError(t, err)

	fs, err := fsProvider.GetUserFilesystem(context.Background(), user)
	require.NoError(t, err)

	return fs

}

func TestMain(m *testing.M) {
	db = sqlite.NewTestDatabase()
	userService = identity.NewUserService(db, sqlite.IdentityRepository{})
	fsProvider = filesystem.NewServiceProvider(db, sqlite.FilesystemRepository{})

	status := m.Run()

	db.Close()

	os.Exit(status)
}
