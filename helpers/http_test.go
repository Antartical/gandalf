package helpers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAbortWithStatus(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(recorder)
		status := http.StatusBadRequest
		err := errors.New("Whoops!")

		AbortWithStatus(context, status, err)

		assert.Equal(status, recorder.Code)
	})
}

func TestHTTPErrror(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		status := http.StatusBadRequest
		err := errors.New("Whoops!")

		httpErrorSerializer := NewHTTPError(status, err)

		assert.Equal(status, httpErrorSerializer.Code)
		assert.Equal(err.Error(), httpErrorSerializer.Error)
	})
}
