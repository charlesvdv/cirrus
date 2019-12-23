package filesystem

type CreateDirectoryRequest struct {
	Name   string
	Parent string
}

type DirectoriesService interface {
	CreateDirectory(req CreateDirectoryRequest) (Directory, error)
}
