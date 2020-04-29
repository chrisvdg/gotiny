package business_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/chrisvdg/gotiny/backend"
	"github.com/chrisvdg/gotiny/business"
	"github.com/chrisvdg/gotiny/utils"
	"github.com/stretchr/testify/assert"
)

var testDir string

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

func Test_List(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)

	result, err := l.List()
	assert.NoError(err)
	assert.Equal([]byte("[]"), result)

	entries := make(map[string]string)
	entries["foo"] = "https://foo.bar"
	entries["ping"] = "https://ping.pong"

	for id, url := range entries {
		_, err := l.Create(id, url)
		assert.NoError(err)
	}

	result, err = l.List()
	assert.NoError(err)
	var tResult []backend.TinyURL
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Len(tResult, 2)

	// add an entry and check len
	_, err = l.Create("lorem", "https://lorem.ipsum")
	assert.NoError(err)
	result, err = l.List()
	assert.NoError(err)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Len(tResult, 3)

	// remove an entry and check len
	err = l.Delete("lorem")
	assert.NoError(err)
	assert.NoError(err)
	result, err = l.List()
	assert.NoError(err)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Len(tResult, 2)

	rEntries := make(map[string]string)
	for _, r := range tResult {
		rEntries[r.ID] = r.URL
	}

	for id, url := range entries {
		rURL, ok := rEntries[id]
		assert.True(ok, "Expected key %s not found", id)
		assert.Equal(url, rURL)
	}
}

func Test_Create(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)
	var tResult backend.TinyURL

	result, err := l.Create("foo", "http://foo.bar")
	assert.NoError(err)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Equal("foo", tResult.ID)
	assert.Equal("http://foo.bar", tResult.URL)

	result, err = l.Create("hello", "https://hello.world")
	assert.NoError(err)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Equal("hello", tResult.ID)
	assert.Equal("https://hello.world", tResult.URL)

	result, err = l.Create("ping", "ping.ping")
	assert.NoError(err)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Equal("ping", tResult.ID)
	assert.Equal("http://ping.ping", tResult.URL)

}

func Test_CreateNoID(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)
	var tResult backend.TinyURL

	result, err := l.Create("", "http://foo.bar")
	assert.NoError(err)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Len(tResult.ID, 5)
	assert.Equal("http://foo.bar", tResult.URL)
}

func Test_CreateInvalidID(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)

	result, err := l.Create("foo bar", "http://foo.bar")
	assert.Error(err)
	assert.Nil(result)
}

func Test_CreateInvalidURL(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)

	result, err := l.Create("foo", "http://foo bar")
	assert.Error(err)
	assert.Nil(result)
}

func Test_Get(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)
	var tResult backend.TinyURL

	id := "foo"
	url := "http://foo.bar"

	_, err = l.Create(id, url)
	assert.NoError(err)
	result, err := l.Get(id)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Equal(id, tResult.ID)
	assert.Equal(url, tResult.URL)

}

func Test_Update(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)
	var tResult backend.TinyURL

	id := "foo"
	url := "http://foo.bar"
	url2 := "http://hello.world"

	_, err = l.Create(id, url)
	assert.NoError(err)

	err = l.Update(id, url2)
	assert.NoError(err)

	result, err := l.Get(id)
	assert.NoError(err)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Equal(id, tResult.ID)
	assert.Equal(url2, tResult.URL)
	assert.NotEqual(url, tResult.URL)
}

func Test_Delete(t *testing.T) {
	assert := assert.New(t)
	l, err := business.NewFileBackedLogic(getFilePath(), false, 5)
	assert.NoError(err)
	var tResult backend.TinyURL

	id := "foo"
	url := "http://foo.bar"

	_, err = l.Create(id, url)
	assert.NoError(err)

	result, err := l.Get(id)
	err = json.Unmarshal(result, &tResult)
	assert.NoError(err)
	assert.Equal(id, tResult.ID)
	assert.Equal(url, tResult.URL)

	err = l.Delete(id)
	assert.NoError(err)
	result, err = l.Get(id)
	assert.EqualError(err, backend.ErrNotFound.Error())
	assert.Nil(result)

	err = l.Delete(id)
	assert.NoError(err)
}

func getFilePath() string {
	return path.Join(testDir, fmt.Sprintf("backend%s.json", utils.GenerateID(5)))
}
