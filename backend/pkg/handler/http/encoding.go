package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func decodeRequest(r *http.Request, container interface{}) error {
	if err := validateContentType(r); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("%w: unable to read request body", ErrInternalServerError)
	}

	if err := json.Unmarshal(body, container); err != nil {
		return ErrRequestBodyInvalid
	}

	return nil
}

func encodeResponse(w http.ResponseWriter, container interface{}) {
	if container != nil {
		if err := encodeAndWriteBody(w, container); err != nil {
			encodeError(w, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func encodeError(w http.ResponseWriter, receivedErr error) {
	statusCode := http.StatusBadRequest
	if isInternalServerError(receivedErr) {
		statusCode = http.StatusInternalServerError
	}

	err := encodeAndWriteBody(w, errorResponse{Message: receivedErr.Error()})
	if err != nil {
		// TODO: log something here
		statusCode = 501
		_ = encodeAndWriteBody(w, errorResponse{Message: "internal error"})
	}

	w.WriteHeader(statusCode)
}

func encodeAndWriteBody(w http.ResponseWriter, container interface{}) error {
	encodedBody, err := json.Marshal(container)
	if err != nil {
		return ErrInternalServerError
	}

	w.Write(encodedBody)
	return nil
}

func validateContentType(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("%w with value: '%v'", ErrInvalidContentType, contentType)
	}
	return nil
}

func isInternalServerError(err error) bool {
	return errors.Is(err, ErrInternalServerError)
}
