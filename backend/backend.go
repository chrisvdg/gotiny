package backend

import (
	"errors"
	"strconv"
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

// TinyURL represents a tiny url entry
type TinyURL struct {
	ID      string   `json:"id"`
	URL     string   `json:"url"`
	Created JSONTime `json:"created"`
}

// JSONTime is a time.Time wrapper that JSON (un)marshals into a unix timestamp
type JSONTime time.Time

// MarshalJSON is used to convert the timestamp to JSON
func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *JSONTime) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

// Unix returns the unix time stamp of the underlaying time object
func (t JSONTime) Unix() int64 {
	return time.Time(t).Unix()
}

// Time returns the JSON time as a time.Time instance
func (t JSONTime) Time() time.Time {
	return time.Time(t)
}

// String returns time as a formatted string
func (t JSONTime) String() string {
	return t.Time().String()
}
