package postgres

import (
	"database/sql"

	fs "github.com/charlesvdv/cirrus/backend/pkg/filesystem"
)

type FilesystemStore struct {
	db *sql.DB
}

func (s *FilesystemStore) GetDirectory(id fs.InodeID) (fs.Directory, error) {
	// TODO
	return fs.Directory{}, nil
}

func (s *FilesystemStore) SaveDirectory(dir fs.Directory) error {
	// TODO
	return nil
}
