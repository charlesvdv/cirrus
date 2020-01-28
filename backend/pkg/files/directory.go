package files

type DirectoryHandler interface {
	Create(metadata CreateDirectoryInfo) (DirectoryMetadata, error)
	Delete(id string) error
}

type CreateDirectoryInfo struct {
	ParentID string
	Name     string
}

type directoryHandlerImpl struct {
}

func (impl directoryHandlerImpl) Create(metadata CreateDirectoryInfo) (DirectoryMetadata, error) {
	return Metadata{}, ErrNotImplemented
}

func (impl directoryHandlerImpl) Delete(id string) error {
	return ErrNotImplemented
}
