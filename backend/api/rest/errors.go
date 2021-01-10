package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/charlesvdv/cirrus/backend/pkg/identity"
)

func convertIdentityErr(err error) *restError {
	if errors.Is(err, identity.ErrInternal) {
		return &restError{
			err:     err,
			Message: err.Error(),
			status:  http.StatusInternalServerError,
		}
	}
	if errors.Is(err, identity.ErrUnauthorized) {
		return errUnauthorized(err)
	}

	return errBusinessLogic(err)
}

func errBusinessLogic(err error) *restError {
	return &restError{
		err:     err,
		Message: err.Error(),
		status:  http.StatusBadRequest,
	}
}

func errUnauthorized(err error) *restError {
	return &restError{
		err:     err,
		Message: "Unauthorized",
		status:  http.StatusUnauthorized,
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
