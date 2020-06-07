package business

import (
	"github.com/chrisvdg/gotiny/backend"
	"github.com/chrisvdg/gotiny/utils"
)

var (
	//ValidationErrors contains a list of possible validation error
	ValidationErrors = []error{
		backend.ErrIDInUse,
	}
	// ErrTinyURLNotFound represents an error where a Tiny URL could not be found in the backend
	ErrTinyURLNotFound = backend.ErrNotFound
)

func init() {
	ValidationErrors = append(ValidationErrors, utils.ValidationErrors...)
}

// IsValidationError is a convenience error to check if an error was a validation error
func IsValidationError(err error) bool {
	for _, e := range ValidationErrors {
		if e == err {
			return true
		}
	}

	return false
}
