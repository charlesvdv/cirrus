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

func newServiceUtils(db sqlite.Database) serviceUtils {
	return serviceUtils{
		db:          db,
		userService: identity.NewUserService(db, sqlite.IdentityRepository{}),
		fsProvider:  filesystem.NewServiceProvider(db, sqlite.FilesystemRepository{}),
	}
}

type serviceUtils struct {
	db          sqlite.Database
	userService identity.UserService
	fsProvider  filesystem.ServiceProvider
}

func (sw serviceUtils) init(t *testing.T) filesystem.FilesystemService {
	var user cirrus.User
	err := sw.userService.CreateUser(context.Background(), &user, sw.fsProvider.InitUserFilesystem)
	require.NoError(t, err)

	fs, err := sw.fsProvider.GetUserFilesystem(context.Background(), user)
	require.NoError(t, err)

	return fs
}

func (sw serviceUtils) cleanUser(t *testing.T, fs filesystem.FilesystemService) {
	// TODO: delete user when possible
}

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

func TestFilesystem_MakeDirectory(t *testing.T) {
	db := sqlite.NewTestDatabase()
	defer db.Close()

	utils := newServiceUtils(db)

	ctx := context.Background()

	t.Run("duplicate directory", func(t *testing.T) {
		fs := utils.init(t)

		dirCreate := cirrus.DirectoryCreate{
			Name: "test",
		}

		_, err := fs.MakeDirectory(ctx, cirrus.MustParsePath(""), dirCreate)
		require.NoError(t, err)

		_, err = fs.MakeDirectory(ctx, cirrus.MustParsePath(""), dirCreate)
		require.Error(t, err)
		cerr, ok := err.(*cirrus.Error)
		require.True(t, ok)
		require.Equal(t, cirrus.ErrCodeAlreadyExist, cerr.Code)
	})

	t.Run("path does not exist", func(t *testing.T) {
		fs := utils.init(t)

		dirCreate := cirrus.DirectoryCreate{
			Name: "test",
		}

		_, err := fs.MakeDirectory(ctx, cirrus.MustParsePath("/random/path"), dirCreate)
		require.Error(t, err)
		cerr, ok := err.(*cirrus.Error)
		require.True(t, ok)
		require.Equal(t, cirrus.ErrCodeInvalidInput, cerr.Code)
	})
}

func TestFilesystem_List(t *testing.T) {
	db := sqlite.NewTestDatabase()
	defer db.Close()

	utils := newServiceUtils(db)

	ctx := context.Background()

	t.Run("empty", func(t *testing.T) {
		fs := utils.init(t)
		content, err := fs.List(ctx, cirrus.MustParsePath(""))
		require.NoError(t, err)
		require.Len(t, content, 0)
	})

	t.Run("with directories", func(t *testing.T) {
		fs := utils.init(t)
		directoryCreate := cirrus.DirectoryCreate{
			Name: "test",
		}
		_, err := fs.MakeDirectory(ctx, cirrus.MustParsePath(""), directoryCreate)
		require.NoError(t, err)

		content, err := fs.List(ctx, cirrus.MustParsePath(""))
		require.NoError(t, err)
		require.Len(t, content, 1)

		otherDirectoryCreate := cirrus.DirectoryCreate{
			Name: "test1",
		}
		_, err = fs.MakeDirectory(ctx, cirrus.MustParsePath(""), otherDirectoryCreate)
		require.NoError(t, err)

		content, err = fs.List(ctx, cirrus.MustParsePath(""))
		require.NoError(t, err)
		require.Len(t, content, 2)
	})

	t.Run("with nested directories", func(t *testing.T) {
		fs := utils.init(t)
		directoryCreate := cirrus.DirectoryCreate{
			Name: "parent",
		}
		_, err := fs.MakeDirectory(ctx, cirrus.MustParsePath(""), directoryCreate)
		require.NoError(t, err)

		content, err := fs.List(ctx, cirrus.MustParsePath(""))
		require.NoError(t, err)
		require.Len(t, content, 1)

		directory, ok := content[0].(cirrus.Directory)
		require.True(t, ok)
		require.Equal(t, "parent", directory.Name)

		nestedDirectoryCreate := cirrus.DirectoryCreate{
			Name: "nested",
		}
		_, err = fs.MakeDirectory(ctx, cirrus.MustParsePath("parent"), nestedDirectoryCreate)
		require.NoError(t, err)

		content, err = fs.List(ctx, cirrus.MustParsePath(""))
		require.NoError(t, err)
		require.Len(t, content, 1)

		content, err = fs.List(ctx, cirrus.MustParsePath("parent"))
		require.NoError(t, err)
		require.Len(t, content, 1)

		nestedDirectory, ok := content[0].(cirrus.Directory)
		require.True(t, ok)
		require.Equal(t, "nested", nestedDirectory.Name)
	})
}
