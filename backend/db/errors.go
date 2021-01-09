package db

type DuplicateError struct {
	err error
}

func (e DuplicateError) Error() string {
	return e.err.Error()
}
