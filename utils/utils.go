package utils

import (
	"math/rand"
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
