package cirrus

import (
	"context"
	"strings"
	"time"
)

// ObjectID is a unique identifier given to any fileobject.
type ObjectID string

func (id ObjectID) String() string {
	return string(id)
}

// File describes a user file
type File struct {
	ID     ObjectID
	Name   string
	Parent *ObjectID
	Owner  UserID
}

// Directory describes a user directory
type Directory struct {
	ID        ObjectID
	Name      string
	Parent    *ObjectID
	CreatedAt time.Time
	Owner     UserID
}

// DirectoryCreate describes the field needed to be inputed by the user
// to create a directory.
type DirectoryCreate struct {
	Name string
}

// Validate checks if a DirectoryCreate is valid
func (d *DirectoryCreate) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	if d.Name == "" {
		return Errorf(ErrCodeInvalidInput, "directory name is empty")
	}

	return nil
}

// FilesystemObject describes any object representing a filesystem.
// For example: directories, files, ...
type FilesystemObject interface {
	fsObject()
}

func (f File) fsObject()      {}
func (d Directory) fsObject() {}

// IsDirectory checks if a FilesystemObject is a directory.
func IsDirectory(obj FilesystemObject) bool {
	_, ok := obj.(Directory)
	return ok
}

// Filesystem describes the operation that can be made on filesystem.
type Filesystem interface {
	List(ctx context.Context, path Path) ([]FilesystemObject, error)
	MakeDirectory(ctx context.Context, directory *Directory) error
}
