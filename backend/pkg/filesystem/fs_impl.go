package filesystem

import "github.com/charlesvdv/cirrus/backend/pkg/errors"

type userFilesystemImpl struct {
	driveResolver DriveResolver
}

func (fs *userFilesystemImpl) CreateDirectory(req CreateDirectoryRequest) (Directory, error) {
	if req.Name == "" {
		return Directory{}, errors.NewFunctionalError("Directory name can not be empty")
	}

	parentID := inodeIDFromString(req.ParentID)
	drive, err := fs.driveResolver.Resolve(parentID)
	if err != nil {
		return Directory{}, err
	}

	dir := newDirectoryBuilder().
		withName(req.Name).
		withParentID(parentID).
		build()

	dir, err = drive.CreateDirectory(dir)
	if err != nil {
		return dir, err
	}

	return dir, nil
}
