package filesystem

import "github.com/charlesvdv/cirrus/backend/pkg/errors"

type userFilesystemImpl struct {
	metadataStore MetadataStorer
}

func (fs *userFilesystemImpl) CreateDirectory(req CreateDirectoryRequest) (Directory, error) {
	if req.Name == "" {
		return Directory{}, errors.NewFunctionalError("Directory name can not be empty")
	}

	var parentID InodeID
	if req.ParentID == "" {
		parentID = emptyInodeID()
	} else {
		parentID := InodeIDFromString(req.ParentID)
		_, err := fs.metadataStore.GetDirectory(parentID)
		if err != nil {
			return Directory{}, errors.NewFunctionalErrorWithCause("Unknow parent directory", err)
		}
	}

	dir := NewDirectoryBuilder().
		WithName(req.Name).
		WithParentID(parentID).
		Build()

	err := fs.metadataStore.SaveDirectory(dir)
	if err != nil {
		return dir, err
	}

	return dir, nil
}
