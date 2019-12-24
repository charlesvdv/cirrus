package directories

import (
	"net/http"

	fs "github.com/charlesvdv/cirrus/backend/pkg/filesystem"
	httputils "github.com/charlesvdv/cirrus/backend/pkg/handler/http"
)

type DirectoriesHandler struct {
	filesystem fs.UserFilesystem
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

func (h *DirectoriesHandler) CreateDirectory(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateDirectoryRequest
	if err := httputils.DecodeRequest(r, requestBody); err != nil {
		httputils.EncodeError(w, err)
	}

	directory, err := h.filesystem.CreateDirectory(fs.CreateDirectoryRequest{
		Name:     requestBody.Name,
		ParentID: requestBody.Parent,
	})
	if err != nil {
		httputils.EncodeError(w, err)
	}

	httputils.EncodeResponse(w, encodeDirResponse(directory))
}

func encodeDirResponse(dir fs.Directory) DirectoryResponse {
	return DirectoryResponse{
		ID:          dir.ID().String(),
		Name:        dir.Name(),
		ParentID:    dir.ParentID().String(),
		CreatedTime: httputils.FormatTime(dir.CreatedTime()),
	}
}
