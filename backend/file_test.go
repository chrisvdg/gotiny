package backend_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/chrisvdg/gotiny/backend"
	"github.com/chrisvdg/gotiny/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testDir = ""

func TestMain(m *testing.M) {
	var err error
	testDir, err = ioutil.TempDir("", "filebackend_test")
	if err != nil {
		log.Fatalf("Failed to create test dir: %s", err)
	}
	exitCode := m.Run()
	os.RemoveAll(testDir)
	os.Exit(exitCode)
}

func Test_BackendCreation(t *testing.T) {
	assert := assert.New(t)

	// Empty file location should return an error
	_, err := backend.NewFile("")
	assert.Error(err)

	backendFile := path.Join(testDir, generateBackendfilename())
	_, err = backend.NewFile(backendFile)

	// Make sure the backend file is created by the constructor
	assert.FileExists(backendFile)
}

func Test_Create(t *testing.T) {
	assert := assert.New(t)
	backendFile := path.Join(testDir, generateBackendfilename())
	b, err := backend.NewFile(backendFile)
	assert.NoError(err)

	id := utils.GenerateID(5)
	url := "http://foo.bar"

	res, err := b.Create(id, url)
	assert.NoError(err)
	assert.Equal(res.ID, id)
	assert.Equal(res.URL, url)
	assert.NotNil(res.Created)
}

func Test_CreateIDInUser(t *testing.T) {
	assert := assert.New(t)
	backendFile := path.Join(testDir, generateBackendfilename())
	b, err := backend.NewFile(backendFile)
	assert.NoError(err)

	id := utils.GenerateID(5)
	url := "http://foo.bar"
	url2 := "http://lorem.ipsum"

	res, err := b.Create(id, url)
	assert.NoError(err)
	assert.Equal(res.ID, id)
	assert.Equal(res.URL, url)
	assert.NotNil(res.Created)

	res, err = b.Create(id, url)
	assert.NoError(err)
	assert.Equal(res.ID, id)
	assert.Equal(res.URL, url)
	assert.NotNil(res.Created)

	_, err = b.Create(id, url2)
	assert.EqualError(err, backend.ErrIDInUse.Error())
}

func Test_Get(t *testing.T) {
	assert := assert.New(t)
	_, b := createFilebackend(t)
	id, url := addEntry(t, b)

	res, err := b.Get(id)
	assert.NoError(err)
	assert.Equal(res.ID, id)
	assert.Equal(res.URL, url)
	assert.NotNil(res.Created)
}

func Test_List(t *testing.T) {
	assert := assert.New(t)
	_, b := createFilebackend(t)
	id1, url1 := addEntry(t, b)
	id2, url2 := addEntry(t, b)

	res, err := b.List()
	assert.NoError(err)
	assert.Len(res, 2)
	exp_entries := make(map[string]string)
	exp_entries[id1] = url1
	exp_entries[id2] = url2
	ids_res := []string{}
	urls_res := []string{}
	for _, entry := range res {
		ids_res = append(ids_res, entry.ID)
		urls_res = append(urls_res, entry.URL)
	}
	for id, url := range exp_entries {
		assert.Contains(ids_res, id)
		assert.Contains(urls_res, url)
	}
}

func Test_Update(t *testing.T) {
	assert := assert.New(t)
	_, b := createFilebackend(t)
	id, url1 := addEntry(t, b)
	url2 := generateURL()

	res1, err := b.Get(id)
	assert.NoError(err)
	assert.Equal(res1.URL, url1)

	req := backend.TinyURL{
		ID:  id,
		URL: url2,
	}
	err = b.Update(req)
	assert.NoError(err)

	res2, err := b.Get(id)
	assert.NoError(err)
	assert.Equal(res2.URL, url2)
}

func Test_Remove(t *testing.T) {
	assert := assert.New(t)
	_, b := createFilebackend(t)
	id1, _ := addEntry(t, b)
	id2, _ := addEntry(t, b)

	resList, err := b.List()
	assert.NoError(err)
	assert.Len(resList, 2)

	err = b.Remove(id1)
	assert.NoError(err)
	resList, err = b.List()
	assert.NoError(err)
	assert.Len(resList, 1)
	_, err = b.Get(id1)
	assert.EqualError(err, backend.ErrNotFound.Error())
	resGet, err := b.Get(id2)
	assert.NoError(err)
	assert.Equal(id2, resGet.ID)

	err = b.Remove(id2)
	assert.NoError(err)
	resList, err = b.List()
	assert.NoError(err)
	assert.Len(resList, 0)
	_, err = b.Get(id2)
	assert.EqualError(err, backend.ErrNotFound.Error())
}

func Test_ReuseBackendFile(t *testing.T) {
	assert := assert.New(t)
	backendFilePath, b := createFilebackend(t)
	id1, url1 := addEntry(t, b)
	id2, url2 := addEntry(t, b)

	res, err := b.Get(id1)
	assert.NoError(err)
	assert.Equal(url1, res.URL)
	res, err = b.Get(id2)
	assert.NoError(err)
	assert.Equal(url2, res.URL)

	b2, err := backend.NewFile(backendFilePath)
	assert.NoError(err)

	list, err := b2.List()
	fmt.Println(list)

	res, err = b2.Get(id1)
	assert.NoError(err)
	assert.Equal(url1, res.URL)
	res, err = b2.Get(id2)
	assert.NoError(err)
	assert.Equal(url2, res.URL)
}

// Test_UseOfExistingEmptyFile tests initiation of a backend with JSON file content
// The content in the file can represent different versions of empty JSON structures
func Test_UseOfExistingEmptyFile(t *testing.T) {
	assert := assert.New(t)
	file_data := []byte("")
	backendFilePath := path.Join(testDir, generateBackendfilename())

	err := ioutil.WriteFile(backendFilePath, file_data, 0666)
	assert.NoError(err)
	_, err = backend.NewFile(backendFilePath)
	assert.NoError(err)

	file_data = []byte("{}")
	err = ioutil.WriteFile(backendFilePath, file_data, 0666)
	assert.NoError(err)
	_, err = backend.NewFile(backendFilePath)
	assert.NoError(err)

	file_data = []byte("[]")
	err = ioutil.WriteFile(backendFilePath, file_data, 0666)
	assert.NoError(err)
	_, err = backend.NewFile(backendFilePath)
	assert.NoError(err)
}

// create_backend_single_entry is a convenience function to create a backend
// returns backendfile, backend object
func createFilebackend(t *testing.T) (string, backend.Backend) {
	assert := assert.New(t)
	backendFile := path.Join(testDir, generateBackendfilename())
	b, err := backend.NewFile(backendFile)
	assert.NoError(err)

	return backendFile, b
}

// Adds an entry with random ID and url to the backend
// returns id and url of the created entry
func addEntry(t *testing.T, b backend.Backend) (string, string) {
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
