package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupPingRouter() *gin.Engine {
	router := gin.Default()
	RegisterPingRoutes(router)
	return router
}

func TestPing(t *testing.T) {
	assert := require.New(t)
	router := setupPingRouter()
	var response gin.H

	t.Run("Test ping successfully", func(t *testing.T) {
		expectedResponse := gin.H{
			"data": "pong",
		}

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(recorder, request)

		err := json.Unmarshal(recorder.Body.Bytes(), &response)

		if err != nil {
			assert.Fail("Payload does not match with the expected one")
		}

		assert.Equal(recorder.Result().StatusCode, http.StatusOK)
		assert.Equal(expectedResponse, response)
	})

}
