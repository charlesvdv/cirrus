package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	
	"github.com/charlesvdv/cirrus/backend/pkg/errors"
)

type errorResponse struct {
	Message string `json:"message"`
}

func DecodeRequest(r *http.Request, container interface{}) error {
	if err := validateContentType(r); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.NewBadRequestError("Unable to read the body content")
	}

	if err := json.Unmarshal(body, container); err != nil {
		return errors.NewBadRequestError("Invalid body content")
	}

	return nil
}

func EncodeResponse(w http.ResponseWriter, container interface{}) {
	if container != nil {
		if err := encodeAndWriteBody(w, container); err != nil {
			EncodeError(w, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func EncodeError(w http.ResponseWriter, receivedErr error) {
	statusCode := http.StatusInternalServerError
	message := "internal error"
	if _, ok := receivedErr.(*errors.BadRequestError); ok {
		statusCode = http.StatusBadRequest
		message = receivedErr.Error()
	}

	err := encodeAndWriteBody(w, errorResponse{Message: message})
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
		return err
	}

	w.Write(encodedBody)
	return nil
}

func validateContentType(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return errors.NewBadRequestError("Invalid Content-Type")
	}
	return nil
}
