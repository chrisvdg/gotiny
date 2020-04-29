package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const filePerm os.FileMode = 0666

// NewFile returns a new file backend
func NewFile(filePath string) (*File, error) {
	if filePath == "" {
		return nil, fmt.Errorf("no backend file location provided")
	}
	backend := &File{
		filePath: filePath,
		data:     fileData{},
	}
	err := backend.ensureFile()
	if err != nil {
		return nil, fmt.Errorf("failed to ensure backend file exists: %s", err)
	}

	err = backend.read()
	if err != nil {
		return nil, fmt.Errorf("Failed to read data from backendfile: %s", err)
	}

	return backend, nil
}

// File represents a file backend implementation
type File struct {
	filePath string
	data     fileData
}

// List implements backend.List
func (f *File) List() ([]TinyURL, error) {
	result := []TinyURL{}
	for k, v := range f.data {
		result = append(result, TinyURL{
			ID:      k,
			URL:     v.URL,
			Created: v.Created,
		})
	}

	return result, nil
}

// Create implements backend.Create
func (f *File) Create(id string, url string) (TinyURL, error) {
	if res, ok := f.data[id]; ok {
		if res.URL == url {
			return f.Get(id)
		}
		return TinyURL{}, ErrIDInUse
	}
	t := TinyURL{
		ID:      id,
		URL:     url,
		Created: JSONTime(time.Now()),
	}
	f.data[id] = fileEntry{
		URL:     url,
		Created: JSONTime(t.Created),
	}

	err := f.save()
	if err != nil {
		// Remove new entry on error
		delete(f.data, id)
		return TinyURL{}, fmt.Errorf("failed to save to backend: %s", err)
	}

	return t, nil
}

// Get implements backend.Get
func (f *File) Get(id string) (TinyURL, error) {
	val, ok := f.data[id]
	if !ok {
		return TinyURL{}, ErrNotFound
	}

	result := TinyURL{
		ID:      id,
		URL:     val.URL,
		Created: val.Created,
	}

	return result, nil
}

// Update implements backend.Update
func (f *File) Update(entry TinyURL) error {
	val, ok := f.data[entry.ID]
	if !ok {
		return ErrNotFound
	}

	f.data[entry.ID] = fileEntry{
		URL:     entry.URL,
		Created: val.Created, // Created time stamp should not be updated, maybe add updated timestamp in later release
	}

	err := f.save()
	if err != nil {
		// Undo update when saving failed
		f.data[entry.ID] = val
		return fmt.Errorf("failed to save update to file backend: %err", err)
	}

	return nil
}

// Remove implents backend.Remove
func (f *File) Remove(id string) error {
	val, ok := f.data[id]
	if !ok {
		return nil
	}
	delete(f.data, id)
	err := f.save()
	if err != nil {
		f.data[id] = val
		return fmt.Errorf("failed to save delete to file backend: %err", err)
	}

	return nil
}

// Close implements backend.Close
func (f *File) Close() error {
	return f.save()
}

// save writes the current file backend data to the backend file
func (f *File) save() error {
	data, err := json.Marshal(f.data)
	if err != nil {
		return fmt.Errorf("failed to marshal backend data to json: %s", err)
	}
	err = ioutil.WriteFile(f.filePath, data, filePerm)
	if err != nil {
		return fmt.Errorf("failed to openfile to write to: %s", err)
	}
	return nil
}

// read reads the backend file to in memory objects for the file backend
func (f *File) read() error {
	data, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to read backend file: %s", err)
	}

	if string(data) == "" || string(data) == "[]" {
		return nil
	}
	err = json.Unmarshal(data, &f.data)
	if err != nil {
		return fmt.Errorf("failed to parse data from backend file: %s", err)
	}

	return nil
}

// ensureFile ensures that the backend file exists
func (f *File) ensureFile() error {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, filePerm)
	if err != nil {
		return fmt.Errorf("something went wrong creating/reading backend file: %s", err)
	}

	return file.Close()
}

// List represents a list of backend entries
type fileData map[string]fileEntry

// Entry represents a tiny URL entry in the backend
type fileEntry struct {
	URL     string   `json:"url"`
	Created JSONTime `json:"created"`
}
