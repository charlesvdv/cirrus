package files

type MetadataStorer interface {
	Get(id ID) (Metadata, error)
	List(parent ID) ([]Metadata, error)
	Create(metadata Metadata) error
}
