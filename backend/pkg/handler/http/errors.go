package http

type BadRequestError struct {
	msg string
}

func NewBadRequestError(msg string) error {
	return BadRequestError{
		msg: msg,
	}
}

func (e BadRequestError) Error() string {
	return e.msg
}
