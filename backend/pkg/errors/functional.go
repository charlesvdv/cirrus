package errors

func NewFunctionalError(msg string) FunctionalError {
	return NewFunctionalErrorWithCause(msg, nil)
}

func NewFunctionalErrorWithCause(msg string, err error) FunctionalError {
	return FunctionalError{
		cause: err,
		msg:   msg,
	}
}

type FunctionalError struct {
	cause error
	msg   string
}

func (e FunctionalError) Cause() error {
	return e.cause
}

func (e FunctionalError) Error() string {
	return e.msg
}
