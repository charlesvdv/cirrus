package http

import (
	"github.com/charlesvdv/cirrus/backend/pkg/files"
)

type createFileRequest struct {
	Name        string `json:"name"`
	ParentID    string `json:"parent-id"`
	ContentType string `json:"content-type"`
}

type fileResponse struct {
	ID       string `json:"id"`
	ParentID string `json:"parent-id"`
}

func encodeFile(file files.File) fileResponse {
	return fileResponse{
		ID: file.ID(),
	}
}
