package object_store

import (
	"github.com/charlesvdv/cirrus/pkg/storage"
)

type StorageProviderFactory struct{}

func NewStorageProviderFactory() StorageProviderFactory {
	return StorageProviderFactory{}
}

func (factory *StorageProviderFactory) Create() storage.Storer {
	return nil
}
