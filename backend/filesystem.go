package cirrus

// File describes a user file
type File struct {
	ID string
}

// Directory describes a user directory
type Directory struct {
	ID   string
	Name string
}

// Filesystem describes the operation that can be made on filesystem.
type Filesystem interface {
}
