package files

import "time"

type DirectoryBuilder interface {
	WithID(id ID) MetadataBuilder
	WithParentID(id ID) MetadataBuilder
	WithName(name string) MetadataBuilder
	WithCreatedTime(createdTime time.Time) MetadataBuilder
	BuildDirectory() DirectoryMetadata
}

func NewDirectoryBuilder() DirectoryBuilder {
	return NewMetadataBuilder().
		WithType(MetadataTypeDirectory)
}

type FileBuilder interface {
	WithID(id ID) MetadataBuilder
	WithParentID(id ID) MetadataBuilder
	WithName(name string) MetadataBuilder
	WithCreatedTime(createdTime time.Time) MetadataBuilder
	WithSize(size uint64) MetadataBuilder
	BuildFile() FileMetadata
}

func NewFileBuilder() FileBuilder {
	return NewMetadataBuilder().
		WithType(MetadataTypeFile)
}

func defaultMetadata() Metadata {
	return Metadata{
		id:           MustParseID(rootParentID),
		parentID:     MustParseID(rootParentID),
		name:         "",
		createdTime:  time.Now().UTC(),
		size:         0,
		metadataType: MetadataTypeUndefined,
	}
}

type MetadataBuilder struct {
	partialMetadata Metadata
}

func NewMetadataBuilder() MetadataBuilder {
	return MetadataBuilder{
		partialMetadata: defaultMetadata(),
	}
}

func (b MetadataBuilder) WithID(id ID) MetadataBuilder {
	b.partialMetadata.id = id
	return b
}

func (b MetadataBuilder) WithParentID(id ID) MetadataBuilder {
	b.partialMetadata.parentID = id
	return b
}

func (b MetadataBuilder) WithName(name string) MetadataBuilder {
	b.partialMetadata.name = name
	return b
}

func (b MetadataBuilder) WithCreatedTime(createdTime time.Time) MetadataBuilder {
	b.partialMetadata.createdTime = createdTime
	return b
}

func (b MetadataBuilder) WithSize(size uint64) MetadataBuilder {
	b.partialMetadata.size = size
	return b
}

func (b MetadataBuilder) WithType(metadataType string) MetadataBuilder {
	b.partialMetadata.metadataType = metadataType
	return b
}

func (b MetadataBuilder) Build() Metadata {
	return b.partialMetadata
}

func (b MetadataBuilder) BuildDirectory() DirectoryMetadata {
	return b.Build()
}

func (b MetadataBuilder) BuildFile() FileMetadata {
	return b.Build()
}
