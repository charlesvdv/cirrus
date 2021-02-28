package sqlite

import (
	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
	"github.com/rs/zerolog/log"
)

// FilesystemRepository contains the sqlite persistence layer for filesystem related data
type FilesystemRepository struct {
}

// CreateFilesystem impl
func (fr FilesystemRepository) CreateFilesystem(_tx database.Tx, userID cirrus.UserID) error {
	tx := getTx(_tx)

	dirStmt := tx.Prep("INSERT INTO directories (parent_id, name) VALUES (NULL, 'xxx-root-directory-xxx');")
	fsStmt := tx.Prep("INSERT INTO user_filesystems(user_id, directory_root_id) VALUES($user_id, last_insert_rowid());")
	fsStmt.SetInt64("$user_id", formatUserID(userID))

	_, err := dirStmt.Step()
	if err != nil {
		log.Err(err).Msg("failed to create fs root directory")
		return formatError(err)
	}

	_, err = fsStmt.Step()
	if err != nil {
		log.Err(err).Msg("failed to create filesystem")
		return formatError(err)
	}
	return nil
}
