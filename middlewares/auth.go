package middlewares

import (
	"errors"
	"gandalf/models"
	auth "gandalf/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Interface for bearer authentication middleware
type IAuthBearerMiddleware interface {
	HasScopes(scopes []string) gin.HandlerFunc
	GetAuthorizedUser(c *gin.Context) *models.User
}

// Auth middleware for authenticate users with Bearer tokens
type AuthBearerMiddleware struct {
	authService auth.IAuthService
}

// Creates a new auth middleware
func NewAuthBearerMiddleware(authService auth.IAuthService) AuthBearerMiddleware {
	return AuthBearerMiddleware{authService: authService}
}

// Check if the user who perform the request has the given scopes
func (middleware AuthBearerMiddleware) HasScopes(scopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := strings.Split(c.GetHeader("Authorization"), "Bearer ")
		if len(bearer) < 2 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Invalid authorization header"))
			return
		}
		user, err := middleware.authService.GetAuthorizedUser(bearer[1], scopes)
		if err != nil {
			c.AbortWithError(http.StatusForbidden, err)
			return
		}

		c.Set("authorizedUser", user)
	}
}

// Return the authorized user from the given gin context
func (middleware AuthBearerMiddleware) GetAuthorizedUser(c *gin.Context) *models.User {
	user, exists := c.Get("authorizedUser")
	if !exists {
		panic(AuthBearerMiddlewareNotCalledError{})
	}
	return user.(*models.User)
}
