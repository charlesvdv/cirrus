package filesystem

type MetadataStorer interface {
	GetDirectory(id InodeID) (Directory, error)
	SaveDirectory(dir Directory) error
}
