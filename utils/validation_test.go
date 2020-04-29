package utils_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/chrisvdg/gotiny/utils"
	"github.com/stretchr/testify/assert"
)

func Test_IsValidationError(t *testing.T) {
	assert := assert.New(t)
	for _, err := range utils.ValidationErrors {
		assert.True(utils.IsValidationError(err), "Expected %s to pass validation", err)
	}
}

func Test_IsValidationErrorInvalid(t *testing.T) {
	assert := assert.New(t)

	cases := []error{
		fmt.Errorf("test"),
		errors.New("foo"),
		io.EOF,
	}

	for _, tc := range cases {
		assert.False(utils.IsValidationError(tc), "Expected %s to fail validation", tc)
	}
}

func Test_ValidateID(t *testing.T) {
	assert := assert.New(t)
	cases := []string{
		"hello",
		"foobar",
		"foo.bar",
		"H3110",
	}

	for _, tc := range cases {
		err := utils.ValidateID(tc)
		assert.NoError(err, "%s should not fail validation", tc)
	}
}

func Test_ValidateIDInvalid(t *testing.T) {
	assert := assert.New(t)
	cases := []string{
		"foo/bar",
		"foo>bar",
		"foo{bar}",
		"foo[bar]",
		"foo?bar",
		"foo&bar",
		"foo^bar",
		"foo%bar",
		"#foobar",
		"foo|bar",
		"foo bar",
	}

	for _, tc := range cases {
		err := utils.ValidateID(tc)
		assert.Error(err, "%s should fail validation", tc)
	}
}

func Test_ValidateURL(t *testing.T) {
	assert := assert.New(t)
	cases := []string{
		"http://google.com",
		"https://google.com",
		"http://google.com/search",
		"http://google.com?q=searching%20for%20something",
		"http://foo.bar",
	}

	for _, tc := range cases {
		err := utils.ValidateURL(tc)
		assert.NoError(err, "%s should not fail validation", tc)
	}
}

func Test_ValidateURLInvalid(t *testing.T) {
	assert := assert.New(t)
	cases := []string{
		"derp",
		"eggqahrjaqj,..P;FEO534564",
		"http://foo.^.bar",
		"https://derp.{}",
		"https://foo bar",
	}

	for _, tc := range cases {
		err := utils.ValidateURL(tc)
		assert.Error(err, "%s should fail validation", tc)
	}
}
