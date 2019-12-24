package filesystem

import (
	"time"

	"github.com/google/uuid"
)

type InodeID struct {
	id []byte
}

func (i InodeID) String() string {
	return string(i.id)
}

func generateInodeID() InodeID {
	id, err := uuid.New().MarshalText()
	if err != nil {
		// Should not happen because the implementation of `MarshalText` never return an error
		// TODO: log something
	}
	return InodeID {
		id: id,
	}
}

func emptyInodeID() InodeID {
	return InodeID {
		id: []byte{},
	}
}

func inodeIDFromString(rawID string) InodeID {
	if rawID == "" {
		return emptyInodeID()
	}
	return InodeID {
		id: []byte(rawID),
	}
}

type inode struct {
	id          InodeID
	name        string
	parentID    InodeID
	createdTime time.Time
}

type Directory struct {
	inode
}

func (i *inode) ID() InodeID {
	return i.id
}

func (i *inode) Name() string {
	return i.name
}

func (i *inode) ParentID() InodeID {
	return i.parentID
}

func (i *inode) CreatedTime() time.Time {
	return i.createdTime
}
