package business

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisvdg/gotiny/backend"
	"github.com/chrisvdg/gotiny/utils"
	log "github.com/sirupsen/logrus"
)

// NewFileBackedLogic creates a new Logic instance
func NewFileBackedLogic(backendPath string, prettyJSON bool, defaultIDLen int) (*Logic, error) {
	b, err := backend.NewFile(backendPath)
	if err != nil {
		return nil, err
	}

	if defaultIDLen <= 0 {
		defaultIDLen = 5
	}
	l := &Logic{
		backend:      b,
		prettyJSON:   prettyJSON,
		defaultIDLen: defaultIDLen,
	}

	return l, nil
}

// Logic contains a stateful set of business logic
type Logic struct {
	backend backend.Backend
	// Prettifies the json respresentation
	prettyJSON   bool
	defaultIDLen int
}

// List retrieves a list of entries from the backend and returns a json encoding of that list
func (l *Logic) List() ([]byte, error) {
	bData, err := l.backend.List()
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("Failed to list tiny URL entries")
	}

	result, err := formatList(bData, l.prettyJSON)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("Failed to list tiny URL entries")
	}

	return result, nil
}

// Create creates a new entry in the backend
func (l *Logic) Create(id string, url string) ([]byte, error) {
	if id == "" {
		id = utils.GenerateID(l.defaultIDLen)
	}

	err := utils.ValidateID(id)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	err = utils.ValidateURL(url)
	if err != nil {
		return nil, err
	}

	res, err := l.backend.Create(id, url)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("Failed to create new entry")
	}

	data, err := formatEntry(res, l.prettyJSON)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("Failed to create new entry")
	}

	return data, nil
}

// Get returns a json endcoded
func (l *Logic) Get(id string) ([]byte, error) {
	entry, err := l.backend.Get(id)
	if err != nil {
		return nil, err
	}
	data, err := formatEntry(entry, l.prettyJSON)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("Failed to get tiny URL entry")
	}
	return data, nil
}

// Update updates an entry in the backend
func (l *Logic) Update(id string, url string) error {
	original, err := l.backend.Get(id)
	if err != nil {
		log.Error(err)
		return err
	}
	if url == original.URL {
		return nil
	}

	err = utils.ValidateURL(url)
	if err != nil {
		return err
	}

	entry := backend.TinyURL{
		ID:  id,
		URL: url,
	}
	err = l.backend.Update(entry)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("Failed to update entry")
	}

	return nil
}

// Delete deletes an entry from the backend
func (l *Logic) Delete(id string) error {
	_, err := l.backend.Get(id)
	if err != nil {
		if err == backend.ErrNotFound {
			return nil
		}
		log.Error(err)
		return fmt.Errorf("Failed to delete entry")
	}

	err = l.backend.Remove(id)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("Failed to delete entry")
	}

	return nil
}

func formatList(entries []backend.TinyURL, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(entries, "", "\t")
	}

	return json.Marshal(entries)
}

func formatEntry(entry backend.TinyURL, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(entry, "", "\t")
	}

	return json.Marshal(entry)
}
