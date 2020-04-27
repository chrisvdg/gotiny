package backend

import (
	"errors"
	"time"
)

// ErrNotFound represents an error where a tiny url entry could not be found
var ErrNotFound error = errors.New("tiny URL entry not found")

// Backend defines the interface to the backend
type Backend interface {
	// List returns a list of the current tiny URL entries
	List() ([]TinyURL, error)
	// Create creates a new tiny URL entry
	Create(id string, url string) (TinyURL, error)
	// Get returns information of a tiny URL matching provided ID
	Get(id string) (TinyURL, error)
	// Update updates the tiny URL of the provided ID in the entry with the provided values
	Update(entry TinyURL) error
	// Remove removes an entry from the backend
	Remove(id string) error
}

// TinyURL represents a tiny url entry in the backend
type TinyURL struct {
	ID      string
	URL     string
	Created time.Time
}
