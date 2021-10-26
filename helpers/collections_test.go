package helpers

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPqStringArrayContains(t *testing.T) {
	assert := require.New(t)

	t.Run("Test contains", func(t *testing.T) {
		element := "Me"
		pqArray := pq.StringArray{"wowowo", element}

		assert.True(PqStringArrayContains(pqArray, element))
	})

	t.Run("Test not contains", func(t *testing.T) {
		element := "Me"
		pqArray := pq.StringArray{"wowowo"}
		assert.False(PqStringArrayContains(pqArray, element))
	})
}
