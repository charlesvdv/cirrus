package filemetadata

import (
	"database/sql"
	"errors"
	"time"

	"github.com/charlesvdv/cirrus/backend/pkg/files"
)

type Store struct {
	db *sql.DB
}

func (s *Store) Get(id *files.ID) (files.Metadata, error) {
	row := s.db.QueryRow("SELECT parent_id, name, created_time, type FROM cirrus.metadata WHERE id = $1", id.String())

	var parentId, name, metadataType string
	var createdTime time.Time
	err := row.Scan(&id, &parentId, &name, &createdTime, &metadataType)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return files.Metadata{}, files.ErrRecordNotFound
	} else {
		return files.Metadata{}, files.ErrInternalServerError
	}

	metadata := files.NewMetadataBuilder().
		WithID(id).
		WithParentID(files.IDFromString(parentId)).
		WithName(name).
		WithCreatedTime(createdTime).
		WithType(metadataType).
		Build()

	return metadata, nil
}
