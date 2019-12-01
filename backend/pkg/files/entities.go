package files

import (
	"time"
)

type UserID = string
type InodeID = string

type File struct {
	inode
	size uint64
}

type Directory struct {
	inode
}

type inode struct {
	id          InodeID
	name        string
	parentID    InodeID
	createdTime time.Time
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

func (f *File) Size() uint64 {
	return f.size
}
