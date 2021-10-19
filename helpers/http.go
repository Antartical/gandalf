package helpers

import (
	"github.com/gin-gonic/gin"
)

// Struct for http error
type HTTPError struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"status bad request"`
}

// Creates a new http error
func NewHTTPError(status int, err error) HTTPError {
	return HTTPError{
		Code:  status,
		Error: err.Error(),
	}
}

// Write the status and a the given error serialized to a JSON
// into the gin context
func AbortWithStatus(c *gin.Context, status int, err error) {
	c.JSON(status, NewHTTPError(status, err))
}
