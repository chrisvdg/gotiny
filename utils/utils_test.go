package utils_test

import (
	"testing"

	"github.com/chrisvdg/gotiny/utils"
	"github.com/stretchr/testify/assert"
)

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
		"aws.amazon.com",
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
	}

	for _, tc := range cases {
		err := utils.ValidateURL(tc)
		assert.Error(err, "%s should fail validation", tc)
	}
}
