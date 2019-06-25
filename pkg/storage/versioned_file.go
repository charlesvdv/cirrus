package storage

import (
	"math"
)

type VersionedFile struct {
	version uint16
	path    Path
}

func NewFile(path Path) VersionedFile {
	return VersionedFile{
		path:    path,
		version: math.MaxUint16,
	}
}

func NewFileWithVersion(path Path, version uint16) VersionedFile {
	return VersionedFile{
		path:    path,
		version: version,
	}
}

func (f *VersionedFile) Path() Path {
	return f.path
}

func (f *VersionedFile) Version() uint16 {
	// TODO: what should be done when f.IsLastestVersion()?
	return f.version
}

func (f *VersionedFile) IsLastestVersion() bool {
	return f.version == math.MaxUint16
}
