package http

import "errors"

var (
	ErrInternalServerError = errors.New("Internal server error")
	ErrRequestBodyInvalid  = errors.New("Request body is invalid")
	ErrInvalidContentType  = errors.New("Invalid content type")
)
