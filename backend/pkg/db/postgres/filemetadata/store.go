package filemetadata

import (
	"database/sql"
	"fmt"
	"time"

	fs "github.com/charlesvdv/cirrus/backend/pkg/filesystem"
)

const (
	inodeTypeDirectory = "directory"
	inodeTypeFile      = "file"
)

type Store struct {
	db *sql.DB
}

func (s *Store) GetDirectory(id fs.InodeID) (fs.Directory, error) {
	row := s.db.QueryRow("SELECT parent_id, name, created_time, type FROM cirrus.inode WHERE id = $1", id.String())

	var parentId, name, inodeType string
	var createdTime time.Time
	err := row.Scan(&id, &parentId, &name, &createdTime, &inodeType)
	if err != nil {
		return fs.Directory{}, err
	}

	if inodeType != inodeTypeDirectory {
		return fs.Directory{}, fmt.Errorf("Inode '%s' is not a directory", id.String())
	}

	directory := fs.NewDirectoryBuilder().
		WithID(id).
		WithParentID(fs.InodeIDFromString(parentId)).
		WithName(name).
		WithCreatedTime(createdTime).
		Build()

	return directory, nil
}

func (s *Store) SaveDirectory(dir fs.Directory) error {
	// TODO
	return nil
}
