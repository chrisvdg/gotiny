package utils

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

const base64URLCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateID returns a random base64URL string of provided length
// Not guaranteed to be unique
func GenerateID(length int) string {
	r := make([]byte, length)
	for i := range r {
		r[i] = base64URLCharset[rand.Intn(len(base64URLCharset))]
	}
	return string(r)
}

// ValidateID validates a tiny URL ID
func ValidateID(id string) error {
	urlSafe := url.QueryEscape(id)
	if urlSafe != id {
		return fmt.Errorf("ID contains illegal characters")
	}

	return nil
}

// ValidateURL validates provided URL
func ValidateURL(urlStr string) error {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("Invalid URL")
	}

	return err
}
