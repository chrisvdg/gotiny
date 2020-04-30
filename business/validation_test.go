package business_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/chrisvdg/gotiny/business"
	"github.com/chrisvdg/gotiny/utils"
	"github.com/stretchr/testify/assert"
)

func Test_IsValidationError(t *testing.T) {
	assert := assert.New(t)
	for _, err := range business.ValidationErrors {
		assert.True(business.IsValidationError(err), "Expected '%s' to be a validation error", err)
	}

	for _, err := range utils.ValidationErrors {
		assert.True(business.IsValidationError(err), "Expected '%s' to be a validation error", err)
	}
}

func Test_IsValidationErrorInvalid(t *testing.T) {
	assert := assert.New(t)

	cases := []error{
		errors.New("Hello"),
		fmt.Errorf("foobar"),
		io.EOF,
		io.ErrUnexpectedEOF,
		io.ErrClosedPipe,
		io.ErrNoProgress,
		io.ErrShortBuffer,
		io.ErrShortWrite,
		http.ErrAbortHandler,
		http.ErrHijacked,
		http.ErrLineTooLong,
		http.ErrMissingFile,
		os.ErrClosed,
		os.ErrExist,
		os.ErrInvalid,
		os.ErrNoDeadline,
		os.ErrNotExist,
		os.ErrPermission,
	}

	for _, tc := range cases {
		assert.False(business.IsValidationError(tc), "Expected '%s' not to be a validation error", tc)
	}
}
