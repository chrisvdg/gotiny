package utils_test

import (
	"testing"

	"github.com/chrisvdg/gotiny/utils"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateIDLen(t *testing.T) {
	assert := assert.New(t)
	for i := 0; i <= 1000; i++ {
		id := utils.GenerateID(i)
		assert.Len(id, i)
	}
}

func Test_GenerateUniqueness(t *testing.T) {
	assert := assert.New(t)
	generated := []string{}

	// Chose 750 because this unit test then takes about 1 second
	for i := 0; i <= 750; i++ {
		generated = append(generated, utils.GenerateID(5))
	}

	for i, toCheck := range generated {
		for j, check := range generated {
			if i == j {
				continue
			}
			assert.NotEqual(toCheck, check)
		}
	}
}
