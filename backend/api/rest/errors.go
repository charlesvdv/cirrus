package rest

import (
	"fmt"
	"net/http"
)

func errBusinessLogic(err error) *restError {
	return &restError{
		err:     err,
		Message: err.Error(),
		status:  http.StatusBadRequest,
	}
}

func errInternalError(format string, v ...interface{}) *restError {
	return errInternalErrorErr(nil, format, v...)
}

func errInternalErrorErr(err error, format string, v ...interface{}) *restError {
	return newRestError(err, http.StatusInternalServerError, format, v...)
}

func errBadRequest(format string, v ...interface{}) *restError {
	return errBadRequestErr(nil, format, v...)
}

func errBadRequestErr(err error, format string, v ...interface{}) *restError {
	return newRestError(err, http.StatusBadRequest, format, v...)
}

func newRestError(err error, status int, format string, v ...interface{}) *restError {
	return &restError{
		err:     err,
		status:  status,
		Message: fmt.Sprintf(format, v...),
	}
}

type restError struct {
	err     error
	Message string `json:"message"`
	status  int
}

func (e *restError) Error() string {
	if e.err == nil || e.err.Error() == e.Message {
		return fmt.Sprintf("%s: %s", http.StatusText(e.status), e.Message)
	}
	return fmt.Sprintf("%s: %s: %s", http.StatusText(e.status), e.Message, e.err.Error())
}

func (e *restError) Unwrap() error {
	return e.err
}
