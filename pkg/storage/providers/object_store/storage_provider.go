package object_store

// Implement the `Storer` interface for object storage.
type objectStoreStorageProvider struct {
	objectStore ObjectStorer
}

func newObjectStoreStorageProvider(objectStore ObjectStorer) objectStoreStorageProvider {
	return objectStoreStorageProvider{
		objectStore: objectStore,
	}
}
