package filesystem

import "time"

type InodeID = string

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
