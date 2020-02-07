package files

import (
	"errors"
	"strings"

	"github.com/google/logger"
)

const (
	rootParentID = "00000000-0000-0000-0000-000000000000"
)

func NewDirectoryHandler(storer MetadataStorer) DirectoryHandler {
	return directoryHandlerImpl{
		storer: storer,
	}
}

type DirectoryHandler interface {
	Create(metadata CreateDirectoryInfo) (DirectoryMetadata, error)
	List(parentID string) ([]Metadata, error)
	Delete(id string) error
}

type CreateDirectoryInfo struct {
	ParentID string
	Name     string
}

type directoryHandlerImpl struct {
	storer MetadataStorer
}

func (impl directoryHandlerImpl) Create(info CreateDirectoryInfo) (DirectoryMetadata, error) {
	if info.Name == "" {
		return Metadata{}, ErrDirectoryNameEmpty
	}

	parentID, err := impl.getParentID(info.ParentID)
	if err != nil {
		return Metadata{}, err
	}

	directory := NewDirectoryBuilder().
		WithID(NewID()).
		WithParentID(parentID).
		WithName(info.Name).
		Build()

	err = impl.storer.Create(directory)
	if err != nil {
		return Metadata{}, ErrInternalServerError
	}

	return directory, nil
}

func (impl directoryHandlerImpl) List(rawParentID string) ([]Metadata, error) {
	parentID, err := impl.getParentID(rawParentID)
	if err != nil {
		return []Metadata{}, err
	}

	return impl.storer.List(parentID)
}

func (impl directoryHandlerImpl) Delete(rawID string) error {
	return ErrNotImplemented
}

func (impl directoryHandlerImpl) get(rawID string) (Metadata, error) {
	rawID = strings.TrimSpace(rawID)
	if rawID == "" {
		rawID = rootParentID
	}

	id, err := ParseID(rawID)
	if err != nil {
		logger.Error(err)
		return Metadata{}, err
	}

	metadata, err := impl.storer.Get(id)
	if err != nil {
		return Metadata{}, err
	}

	return metadata, nil
}

func (impl directoryHandlerImpl) getParentID(rawParentID string) (ID, error) {
	rawParentID = strings.TrimSpace(rawParentID)

	if rawParentID == "" {
		return MustParseID(rootParentID), nil
	}

	parent, err := impl.getParent(rawParentID)
	if err != nil {
		return ID{}, err
	}

	return parent.ID(), nil
}

func (impl directoryHandlerImpl) getParent(rawParentID string) (Metadata, error) {
	parent, err := impl.get(rawParentID)
	if errors.Is(err, ErrRecordNotFound) {
		return Metadata{}, ErrParentDirectoryNotFound
	} else if err != nil {
		return Metadata{}, ErrInternalServerError
	}

	if parent.Type() != MetadataTypeDirectory {
		return Metadata{}, ErrParentNotDirectory
	}

	return parent, nil
}
