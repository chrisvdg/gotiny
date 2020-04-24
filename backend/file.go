package backend

import (
	"strconv"
	"time"
)

// NewFile returns a new file backend
func NewFile() (*File, error) {
	return &File{}, nil
}

// File represents a file backend implementation
type File struct {
}

// FileConfig represents a file backend config
type FileConfig struct{}

// Entry represents a tiny URL entry in the backend
type Entry struct {
	ID      string
	URL     string
	Created Time
}

// List represents a list of backend entries
type List []Entry

// Time represents a unix time stamp
type Time time.Time

// MarshalJSON is used to convert the timestamp to JSON
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

// Unix returns the unix time stamp of the underlaying time object
func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

// Time returns the JSON time as a time.Time instance
func (t Time) Time() time.Time {
	return time.Time(t)
}

// String returns time as a formatted string
func (t Time) String() string {
	return t.Time().String()
}
