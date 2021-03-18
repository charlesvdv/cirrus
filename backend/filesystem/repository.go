package filesystem

import (
	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
)

// Repository describes the interface to the persistence layer.
type Repository interface {
	CreateDirectory(tx database.Tx, directory *cirrus.Directory) error
	ListDirectoryContent(tx database.Tx, ownerID cirrus.UserID, parent *cirrus.ObjectID) ([]cirrus.FilesystemObject, error)
	ResolvePath(tx database.Tx, ownerID cirrus.UserID, path cirrus.Path) ([]cirrus.FilesystemObject, error)
}
