package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/charlesvdv/cirrus/backend/pkg/files"
)

type filesHandler struct {
	filesService files.FilesFacade
}

type contentType = string

func (h *filesHandler) createFileHandler(w http.ResponseWriter, r *http.Request) {
	uploadTypeQueryString := r.URL.Query()["type"]
	if uploadTypeQueryString == nil || len(uploadTypeQueryString) == 0 || len(uploadTypeQueryString) > 1 {
		badRequestError(w, errors.New("Upload type required"))
		return
	}
	uploadType := uploadTypeQueryString[0]

	if uploadType == "multipart" {
		h.createFileMultipartHandler(w, r)
	} else if uploadType == "resumable" {
		// TODO: create resumable upload
		badRequestError(w, errors.New("Resumable upload not yet supported"))
	} else {
		badRequestError(w, errors.New("Upload type unknown"))
		return
	}
}

func (h *filesHandler) createFileMultipartHandler(w http.ResponseWriter, r *http.Request) {
	multipart, err := r.MultipartReader()
	if err != nil {
		badRequestError(w, errors.New("Multipart body expected"))
		return
	}

	part, err := multipart.NextPart()
	if err != nil {
		badRequestError(w, errors.New("Invalid multipart"))
		return
	}
	fileMetadata, err := decodeMetadataPart(part)
	if err != nil {
		badRequestError(w, err)
		return
	}

	part, err = multipart.NextPart()
	if err != nil {
		badRequestError(w, errors.New("Invalid multipart"))
		return
	}
	fileContent, err := decodeFilePart(part)
	if err != nil {
		badRequestError(w, errors.New("Unable to read files"))
		return
	}

	file, err := h.filesService.Create(files.CreateFileRequest{
		Name:        fileMetadata.Name,
		ParentID:    fileMetadata.ParentID,
		ContentType: fileMetadata.ContentType,
		FileContent: fileContent,
	})
	if err != nil {
		// TODO: handle functional error
		return
	}

	createdSuccess(w, encodeFile(file))
}

func decodeMetadataPart(part *multipart.Part) (createFileRequest, error) {
	var request createFileRequest

	if part.Header.Get("Content-Type") != "application/json" {
		return request, errors.New("Metadata should be json encoded")
	}

	rawContent, err := ioutil.ReadAll(part)
	if err != nil {
		return request, errors.New("Unable to read metadata")
	}

	err = json.Unmarshal(rawContent, &request)
	if err != nil {
		return request, errors.New("Unable to parse json")
	}

	return request, nil
}

func decodeFilePart(part *multipart.Part) ([]byte, error) {
	return ioutil.ReadAll(part)
}
