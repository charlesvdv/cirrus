package filesystem

import (
	"context"
	"errors"
	"time"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
	"github.com/rs/zerolog/log"
)

// FilesystemService implements the methods accessing and mutating
// a user filesystem.
// A filesystem is completely isolated from other users' filesystem.
type FilesystemService struct {
	user       cirrus.User
	repository Repository
	db         database.TxUtils
}

// Owner returns the owner (user) of this filesystem view.
func (fs FilesystemService) Owner() cirrus.User {
	return fs.user
}

// List lists all of the filesystem objects at a given path.
func (fs FilesystemService) List(ctx context.Context, path cirrus.Path) ([]cirrus.FilesystemObject, error) {
	var dirContent []cirrus.FilesystemObject

	err := fs.db.WithTransaction(ctx, func(tx database.Tx) error {
		directoryID, err := fs.resolveDirectoryIDFromPath(tx, fs.user.ID, path)
		if err != nil {
			return err
		}

		dirContentTemp, err := fs.repository.ListDirectoryContent(tx, fs.user.ID, directoryID)
		if err != nil {
			log.Err(err).Msg("failed to list directory content")
			return cirrus.Errorf(cirrus.ErrCodeInternal, "list")
		}
		dirContent = dirContentTemp
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirContent, nil
}

// MakeDirectory creates a directory in a given path.
func (fs FilesystemService) MakeDirectory(ctx context.Context, parent cirrus.Path, directoryCreate cirrus.DirectoryCreate) (cirrus.Directory, error) {
	err := directoryCreate.Validate()
	if err != nil {
		return cirrus.Directory{}, err
	}

	directory := cirrus.Directory{
		Name:      directoryCreate.Name,
		Owner:     fs.user.ID,
		CreatedAt: time.Now().UTC(),
	}

	err = fs.db.WithTransaction(ctx, func(tx database.Tx) error {
		parentID, err := fs.resolveDirectoryIDFromPath(tx, fs.user.ID, parent)
		if err != nil {
			return err
		}

		directory.Parent = parentID
		err = fs.repository.CreateDirectory(tx, &directory)
		if err != nil {
			log.Err(err).Msg("failed to create directory")
			if errors.Is(err, database.ErrDuplicate) {
				return cirrus.Errorf(cirrus.ErrCodeAlreadyExist, "path %s already exist", parent.AppendChild(directoryCreate.Name))
			}
			return cirrus.Errorf(cirrus.ErrCodeInternal, "create directory")
		}
		return nil
	})
	if err != nil {
		return directory, err
	}
	return directory, nil
}

func (fs FilesystemService) resolveDirectoryIDFromPath(tx database.Tx, ownerID cirrus.UserID, path cirrus.Path) (*cirrus.ObjectID, error) {
	var directoryID *cirrus.ObjectID

	if path.IsRoot() {
		directoryID = nil
	} else {
		resolvedPaths, err := fs.resolvePath(tx, fs.user.ID, path)
		if err != nil {
			return directoryID, err
		}
		lastElemInPath := resolvedPaths[len(resolvedPaths)-1]
		if !cirrus.IsDirectory(lastElemInPath) {
			return directoryID, cirrus.Errorf(cirrus.ErrCodeInvalidInput, "path is not a directory")
		}
		directoryIDTemp := lastElemInPath.(cirrus.Directory).ID
		directoryID = &directoryIDTemp
	}

	return directoryID, nil
}

func (fs FilesystemService) resolvePath(tx database.Tx, ownerID cirrus.UserID, path cirrus.Path) ([]cirrus.FilesystemObject, error) {
	resolvedPath, err := fs.repository.ResolvePath(tx, ownerID, path)
	if err != nil {
		log.Err(err).Str("path", path.String()).Msg("failed to resolve path")
		return resolvedPath, cirrus.Errorf(cirrus.ErrCodeInternal, "failed path resolution")
	}

	if len(resolvedPath) != path.ElementCount() {
		return resolvedPath, cirrus.Errorf(cirrus.ErrCodeInvalidInput, "path does not exist")
	}

	return resolvedPath, nil
}
