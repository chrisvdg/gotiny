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
func NewFileBackedLogic() (*Logic, error) {

	return nil, nil
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
		log.Debug(err)
		return nil, fmt.Errorf("Failed to list tiny URL entries from the backend")
	}

	result, err := formatList(bData, l.prettyJSON)
	if err != nil {
		log.Debug(err)
		return nil, fmt.Errorf("Failed to list tiny URL entries from the backend")
	}

	return result, nil
}

func formatList(entries []backend.TinyURL, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(entries, "", "\t")
	}

	return json.Marshal(entries)
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

	l.backend.Create(id, url)

	return nil, nil
}

// Get returns a json endcoded
func (l *Logic) Get(id string) []byte {

	return nil
}

// Update updates an entry in the backend
func (l *Logic) Update(id string, url string) error {

	return nil
}

// Delete deletes an entry from the backend
func (l *Logic) Delete(id string) error {

	return nil
}
