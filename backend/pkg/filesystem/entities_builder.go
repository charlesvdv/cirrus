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

type DirectoryBuilder struct {
	inodeMetadata inodeMetadata
}

func NewDirectoryBuilder() DirectoryBuilder {
	return DirectoryBuilder{
		inodeMetadata: defaultInodeMetadata(),
	}
}

func (b DirectoryBuilder) WithID(id InodeID) DirectoryBuilder {
	b.inodeMetadata.id = id
	return b
}

func (b DirectoryBuilder) WithName(name string) DirectoryBuilder {
	b.inodeMetadata.name = name
	return b
}

func (b DirectoryBuilder) WithParentID(id InodeID) DirectoryBuilder {
	b.inodeMetadata.id = id
	return b
}

func (b DirectoryBuilder) WithCreatedTime(time time.Time) DirectoryBuilder {
	b.inodeMetadata.createdTime = time
	return b
}

func (b DirectoryBuilder) Build() Directory {
	return Directory{
		inodeMetadata: b.inodeMetadata,
	}
}
