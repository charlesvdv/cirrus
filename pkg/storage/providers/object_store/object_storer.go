package object_store

import (
	"io"
	"time"
)

// Store defines an object store interface providing the low-level
// semantic to build a filesystem on top of.
//
// The bucket name is supposed to be passed when the interface is constructed.
// The bucket is expected to be present or will be created during the
// object initialization.
type ObjectStorer interface {
	Get(object string) (io.ReadCloser, error)
	Put(object string, data io.Reader) error
	Stat(object string) (ObjectInfo, error)
	Remove(object string) error
	List(prefix string, handler FileListingHandler) error
}

type ObjectInfo struct {
	Size         uint64
	LastModified time.Time
}

type FileListingHandler = func(name string)
