package filesystem

type UserFilesystem interface {
	CreateDirectory(req CreateDirectoryRequest) (Directory, error)
}

type CreateDirectoryRequest struct {
	Name   string
	ParentID string
}
