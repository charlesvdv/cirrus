package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"

	"github.com/charlesvdv/cirrus/backend/pkg/db/postgres"
	"github.com/charlesvdv/cirrus/backend/pkg/files"
)

func getDirectoryHandler() files.DirectoryHandler {
	metadatastorer := postgres.NewMetadataStore(testDB)
	return files.NewDirectoryHandler(metadatastorer)
}

func TestCreateAndRetrieveRootDirectory(t *testing.T) {
	dirHandler := getDirectoryHandler()

	dirName := gofakeit.BeerName()
	_, err := dirHandler.Create(files.CreateDirectoryInfo{
		Name: dirName,
	})
	assert.NoError(t, err)

	rootContent, err := dirHandler.List("")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rootContent))
	assert.Equal(t, dirName, rootContent[0].Name())

	// TODO: cleanup
}
