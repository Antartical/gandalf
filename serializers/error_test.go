package serializers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorSerializer(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		status := http.StatusBadRequest
		err := errors.New("Whoops!")

		httpErrorSerializer := NewHTTPErrorSerializer(status, err)

		assert.Equal(status, httpErrorSerializer.Code)
		assert.Equal(err.Error(), httpErrorSerializer.Error)
	})
}
