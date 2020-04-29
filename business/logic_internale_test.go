package business

import (
	"testing"

	"github.com/chrisvdg/gotiny/backend"
	"github.com/stretchr/testify/assert"
)

func Test_FormatList(t *testing.T) {
	assert := assert.New(t)
	input := []backend.TinyURL{
		{
			ID:  "FirstEntry",
			URL: "foo.bar",
		},
		{
			ID:  "SecondEntry",
			URL: "lorem.ipsum",
		},
	}
	expectedOutput := []byte(`[{"id":"FirstEntry","url":"foo.bar","created":-62135596800},{"id":"SecondEntry","url":"lorem.ipsum","created":-62135596800}]`)
	expectedPrettyOutput := []byte(`[
	{
		"id": "FirstEntry",
		"url": "foo.bar",
		"created": -62135596800
	},
	{
		"id": "SecondEntry",
		"url": "lorem.ipsum",
		"created": -62135596800
	}
]`)

	out, err := formatList(input, false)
	assert.NoError(err)
	assert.Equal(expectedOutput, out)

	prettyOut, err := formatList(input, true)
	assert.NoError(err)
	assert.Equal(expectedPrettyOutput, prettyOut)
}

func Test_FormatEntry(t *testing.T) {
	assert := assert.New(t)
	input := backend.TinyURL{
		ID:  "AnotherEntry",
		URL: "ping.pong",
	}
	expectedOutput := []byte(`{"id":"AnotherEntry","url":"ping.pong","created":-62135596800}`)
	expectedPrettyOutput := []byte(`{
	"id": "AnotherEntry",
	"url": "ping.pong",
	"created": -62135596800
}`)

	out, err := formatEntry(input, false)
	assert.NoError(err)
	assert.Equal(expectedOutput, out)
	prettyOut, err := formatEntry(input, true)
	assert.NoError(err)
	assert.Equal(expectedPrettyOutput, prettyOut)
}
