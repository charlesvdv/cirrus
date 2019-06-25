// storage provides an storage agnostic interface for storing and manipulating remote files.
//
// The storage interface is build with remote storage in mindlike S3.
//
// Below, you will find a non-exhaustive list of features supported by this package:
//      - File versionning
//      - Basic POSIX filesystem manipulation (move, remove, copy, ...)
//      - Download & upload from a presigned remote url
//
// The durability of files is only garanteed by the storage providers.
package storage

type Storer interface {
	FilesystemOperator
	FileUploader
	FileDownloader
}

type FilesystemOperator interface {
	CreateFile(fileInfo FileInfo) error
	Remove(path Path, recursive bool) error
	Copy(src VersionedFile, dest Path) error
	Move(src Path, dest Path) error
	MakeDirectory(path Path) error
	GetFileInfo(file VersionedFile) (FileInfo, error)
	UpdateFileInfo(fileInfo UpdatableFileInfo) error
	ListDirectory(path Path) ([]Path, error)
}

type FileUploader interface {
	GetUploadUrls(req UploadRequest) (UploadResponse, error)
}

type FileDownloader interface {
	GetDownloadUrls(req DownloadRequest) (DownloadResponse, error)
}
