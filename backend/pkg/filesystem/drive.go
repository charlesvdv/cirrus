package filesystem

type DriveResolver interface {
	Resolve(id InodeID) (Drive, error)
}

type Drive interface {
	CreateDirectory(dir Directory) (Directory, error)
}
