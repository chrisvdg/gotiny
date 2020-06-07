package utils

import (
	"errors"
	"net/url"
)

var (
	// ValidationErrors is a list that contains all validation errors
	ValidationErrors = []error{ErrInvalidID, ErrInvalidURL}
	// ErrInvalidID respresents an invalid tiny URL id error
	ErrInvalidID = errors.New("ID contains illegal characters")
	// ErrInvalidURL respresents an invalid tiny URL url error
	ErrInvalidURL = errors.New("invalid URL")
)

// IsValidationError returns true if provided error is a validation error
func IsValidationError(err error) bool {
	for _, e := range ValidationErrors {
		if err == e {
			return true
		}
	}

	return false
}

// ValidateID validates a tiny URL ID
func ValidateID(id string) error {
	urlSafe := url.QueryEscape(id)
	if urlSafe != id {
		return ErrInvalidID
	}

	return nil
}

// ValidateURL validates provided URL
func ValidateURL(urlStr string) error {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return ErrInvalidURL
	}

	return err
}
