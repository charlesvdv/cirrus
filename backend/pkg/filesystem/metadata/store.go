package metadata

import (
	fs "github.com/charlesvdv/cirrus/backend/pkg/filesystem"
)

type Store struct {
}

func (s *Store) GetDirectory(id fs.InodeID) (fs.Directory, error) {
	// TODO
	return fs.Directory{}, nil
}

func (s *Store) SaveDirectory(dir fs.Directory) error {
	// TODO
	return nil
}
