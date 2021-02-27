package identity

import (
	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
)

// Repository describes the persistence interface for the identity service
type Repository interface {
	CreateUser(tx database.Tx, user *cirrus.User) error
}
