package http

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	ErrorType string `json:"type"`
	Message   string `json:"message"`
}

func badRequestError(w http.ResponseWriter, err error) {
	body, err := json.Marshal(errorResponse{
		Message: err.Error(),
	})
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
	w.WriteHeader(http.StatusBadRequest)
}

func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func createdSuccess(w http.ResponseWriter, resp interface{}) {
	genericSuccess(w, http.StatusCreated, resp)
}

func okSuccess(w http.ResponseWriter, resp interface{}) {
	genericSuccess(w, http.StatusOK, resp)
}

func genericSuccess(w http.ResponseWriter, statusCode int, resp interface{}) {
	if resp != nil {
		body, err := json.Marshal(resp)
		if err != nil {
			internalServerError(w)
			return
		}
		w.Write(body)
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(statusCode)
}
