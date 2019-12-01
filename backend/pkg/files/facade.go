package files

type FilesFacade interface {
	Create(req CreateFileRequest) (File, error)
}

type CreateFileRequest struct {
	Name        string
	ParentID    InodeID
	ContentType string
	FileContent []byte
}
