package sqlutil

type Row interface {
	Scan(dest ...interface{}) error
}
