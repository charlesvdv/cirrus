package filesystem

import "time"

func defaultInode() inode {
	return inode{
		id:          generateInodeID(),
		createdTime: time.Now().UTC(),
		parentID: emptyInodeID(),
		name: "",
	}
}

type directoryBuilder struct {
	inode inode
}

func newDirectoryBuilder() directoryBuilder {
	return directoryBuilder{
		inode: defaultInode(),
	}
}

func (b directoryBuilder) withName(name string) directoryBuilder {
	b.inode.name = name
	return b
}

func (b directoryBuilder) withParentID(id InodeID) directoryBuilder {
	b.inode.id = id
	return b
}

func (b directoryBuilder) build() Directory {
	return Directory {
		inode: b.inode,
	}
}