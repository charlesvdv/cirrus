package http

import (
	"net/http"

	"github.com/charlesvdv/cirrus/backend/pkg/files"
)

type DirectoryHandler struct {
	handler    files.DirectoryHandler
}

type CreateDirectoryRequest struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

type DirectoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ParentID    string `json:"parent-id"`
	CreatedTime string `json:"created-time"`
}

func (h *DirectoryHandler) CreateDirectory(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateDirectoryRequest
	if err := decodeRequest(r, requestBody); err != nil {
		encodeError(w, err)
	}

	directory, err := h.handler.Create(files.CreateDirectoryInfo{
		Name:     requestBody.Name,
		ParentID: requestBody.Parent,
	})
	if err != nil {
		encodeError(w, err)
	}

	encodeResponse(w, encodeDirResponse(directory))
}

func encodeDirResponse(dir files.DirectoryMetadata) DirectoryResponse {
	return DirectoryResponse{
		ID:          dir.ID().String(),
		Name:        dir.Name(),
		ParentID:    dir.ParentID().String(),
		CreatedTime: formatTime(dir.CreatedTime()),
	}
}
