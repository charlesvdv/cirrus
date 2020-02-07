package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/logger"

	"github.com/charlesvdv/cirrus/backend/pkg/db/sqlutil"
	"github.com/charlesvdv/cirrus/backend/pkg/files"
)

const (
	metadataFieldsSQLList = "id, parent_id, name, created_time, type"
)

func NewMetadataStore(db *sql.DB) MetadataStore {
	return MetadataStore{
		db: db,
	}
}

type MetadataStore struct {
	db *sql.DB
}

func (s MetadataStore) Get(id files.ID) (files.Metadata, error) {
	row := s.db.QueryRow("SELECT "+metadataFieldsSQLList+" FROM cirrus.metadata WHERE id = $1", id.String())

	metadata, err := scanMetadataRow(row)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}

func (s MetadataStore) Create(metadata files.Metadata) error {
	_, err := s.db.Exec("INSERT INTO cirrus.metadata ("+metadataFieldsSQLList+") VALUES($1, $2, $3, $4, $5)",
		metadata.ID().String(), metadata.ParentID().String(), metadata.Name(), metadata.CreatedTime(), metadata.Type())
	if err != nil {
		logger.Error(err)
		return files.ErrInternalServerError
	}

	return nil
}

func (s MetadataStore) List(parent files.ID) ([]files.Metadata, error) {
	metadatas := []files.Metadata{}

	rows, err := s.db.Query("SELECT "+metadataFieldsSQLList+" FROM cirrus.metadata WHERE parent_id = $1", parent.String())
	if err != nil {
		logger.Error(err)
		return metadatas, files.ErrInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		newEntry, err := scanMetadataRow(rows)
		if err != nil {
			return metadatas, err
		}
		metadatas = append(metadatas, newEntry)
	}

	if err = rows.Err(); err != nil {
		logger.Error(err)
		return metadatas, files.ErrInternalServerError
	}

	return metadatas, nil
}

func scanMetadataRow(row sqlutil.Row) (files.Metadata, error) {
	var id, parentID, name, metadataType string
	var createdTime time.Time
	err := row.Scan(&id, &parentID, &name, &createdTime, &metadataType)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return files.Metadata{}, files.ErrRecordNotFound
	} else if err != nil {
		logger.Error(err)
		return files.Metadata{}, files.ErrInternalServerError
	}

	metadata := files.NewMetadataBuilder().
		WithID(files.MustParseID(id)).
		WithParentID(files.MustParseID(parentID)).
		WithName(name).
		WithCreatedTime(createdTime).
		WithType(metadataType).
		Build()

	return metadata, nil
}
