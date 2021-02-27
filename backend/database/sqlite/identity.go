package sqlite

import (
	"strconv"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
	"github.com/rs/zerolog/log"
)

// IdentityRepository implements the identity.Repository interface for sqlite
type IdentityRepository struct {
}

// CreateUser impl
func (ir IdentityRepository) CreateUser(_tx database.Tx, user *cirrus.User) error {
	tx := getTx(_tx)

	stmt := tx.Prep("INSERT INTO users (created_at) VALUES ($created_at);")
	stmt.SetText("$created_at", formatTime(user.CreatedAt))
	_, err := stmt.Step()
	if err != nil {
		log.Err(err).Msg("failed to create user")
		return formatError(err)
	}

	user.ID = cirrus.UserID(strconv.FormatInt(tx.LastInsertRowID(), 10))

	return nil
}
