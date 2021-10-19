package helpers

import (
	"gandalf/serializers"

	"github.com/gin-gonic/gin"
)

// Write the status and a the given error serialized to a JSON
// into the gin context
func AbortWithStatus(c *gin.Context, status int, err error) {
	c.JSON(status, serializers.NewHTTPErrorSerializer(status, err))
}
