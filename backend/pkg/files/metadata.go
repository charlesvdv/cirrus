package files

import (
	"time"

	"github.com/google/uuid"
)

const (
	MetadataTypeUndefined = "undefined"
	MetadataTypeFile      = "file"
	MetadataTypeDirectory = "directory"
)

type DirectoryMetadata interface {
	ID() *ID
	Name() string
	ParentID() *ID
	CreatedTime() time.Time
}

type FileMetadata interface {
	ID() *ID
	Name() string
	ParentID() *ID
	CreatedTime() time.Time
	Size() uint64
}

type Metadata struct {
	id           *ID
	name         string
	parentID     *ID
	createdTime  time.Time
	size         uint64
	metadataType string
}

func (m Metadata) ID() *ID {
	return m.id
}

func (m Metadata) Name() string {
	return m.name
}

func (m Metadata) ParentID() *ID {
	return m.parentID
}

func (m Metadata) CreatedTime() time.Time {
	return m.createdTime
}

func (m Metadata) Size() uint64 {
	return m.size
}

func (m Metadata) Type() string {
	return m.metadataType
}

type ID struct {
	val uuid.UUID
}

func (id ID) String() string {
	return id.val.String()
}

func IDFromString(rawID string) *ID {
	if rawID == "" {
		return nil
	} else {
		return &ID{
			val: uuid.New(),
		}
	}
}
