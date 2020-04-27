package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/chrisvdg/gotiny/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_TestSave(t *testing.T) {
	assert := assert.New(t)
	filePath, b := createFilebackend(t)
	data, err := ioutil.ReadFile(filePath)
	assert.NoError(err)
	assert.Len(data, 0)

	id, url := addEntry(t, b)
	data, err = ioutil.ReadFile(filePath)
	assert.NoError(err)
	assert.NotZero(len(data))

	result := make(fileData)
	err = json.Unmarshal(data, &result)
	assert.NoError(err)
	_, found := result[id]
	assert.True(found)
	assert.Equal(result[id].URL, url)
}

func createFilebackend(t *testing.T) (string, Backend) {
	assert := assert.New(t)
	testDir, err := ioutil.TempDir("", "filebackend_internal_test")
	if err != nil {
		log.Fatalf("Failed to create test dir: %s", err)
	}
	backendFile := path.Join(testDir, generateBackendfilename())
	b, err := NewFile(backendFile)
	assert.NoError(err)

	return backendFile, b
}

// Adds an entry with random ID and url to the backend
// returns id and url of the created entry
func addEntry(t *testing.T, b Backend) (string, string) {
	assert := assert.New(t)
	id := utils.GenerateID(5)
	url := generateURL()
	_, err := b.Create(id, url)
	assert.NoError(err)

	return id, url
}

func generateBackendfilename() string {
	return fmt.Sprintf("backend%s.json", utils.GenerateID(5))
}

func generateURL() string {
	return fmt.Sprintf("%s.%s", utils.GenerateID(5), utils.GenerateID(3))
}
