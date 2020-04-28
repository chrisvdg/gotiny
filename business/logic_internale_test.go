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
			ID:  "First entry",
			URL: "foo.bar",
		},
		{
			ID:  "Second entry",
			URL: "lorem.ipsum",
		},
	}
	expectedOutput := []byte(`[{"id":"First entry","url":"foo.bar","created":-62135596800},{"id":"Second entry","url":"lorem.ipsum","created":-62135596800}]`)
	expectedPrettyOutput := []byte(`[
	{
		"id": "First entry",
		"url": "foo.bar",
		"created": -62135596800
	},
	{
		"id": "Second entry",
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
