package cirrus

import "fmt"

const (
	// ErrCodeInternal describes an internal error
	ErrCodeInternal = "internal"
)

// Error describes a business logic error
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Errorf creates an error
func Errorf(code, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
