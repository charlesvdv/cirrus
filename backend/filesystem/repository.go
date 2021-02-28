package filesystem

import (
	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
)

// Repository describes the interface to the persistence layer.
type Repository interface {
	CreateFilesystem(tx database.Tx, userID cirrus.UserID) error
}
