package filesystem

import "time"

func defaultInodeMetadata() inodeMetadata {
	return inodeMetadata{
		id:          generateInodeID(),
		createdTime: time.Now().UTC(),
		parentID:    emptyInodeID(),
		name:        "",
	}
}

type directoryBuilder struct {
	inodeMetadata inodeMetadata
}

func newDirectoryBuilder() directoryBuilder {
	return directoryBuilder{
		inodeMetadata: defaultInodeMetadata(),
	}
}

func (b directoryBuilder) withName(name string) directoryBuilder {
	b.inodeMetadata.name = name
	return b
}

func (b directoryBuilder) withParentID(id InodeID) directoryBuilder {
	b.inodeMetadata.id = id
	return b
}

func (b directoryBuilder) build() Directory {
	return Directory{
		inodeMetadata: b.inodeMetadata,
	}
}
