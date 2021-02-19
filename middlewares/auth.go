package middlewares

import (
	"errors"
	"gandalf/models"
	"gandalf/services"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
IAuthBearerMiddleware -> interface for auth based on bearer token middleware
*/
type IAuthBearerMiddleware interface {
	Authorize() gin.HandlerFunc
	GetAuthorizedUser(c *gin.Context) *models.User
}

/*
AuthBearerMiddleware -> auth middleware for authenticate users with
Bearer tokens
*/
type AuthBearerMiddleware struct {
	authService services.IAuthService
}

/*
NewAuthBearerMiddleware -> creates a new auth middleware.
*/
func NewAuthBearerMiddleware(authService services.IAuthService) AuthBearerMiddleware {
	return AuthBearerMiddleware{authService: authService}
}

/*
Authorize -> identifies the user who perform the request and
return it if him has permissions.
*/
func (middleware AuthBearerMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := strings.Split(c.GetHeader("Authorization"), "Bearer")
		if len(bearer) < 2 {
			c.AbortWithError(400, errors.New("Invalid authorization header"))
		}
		user, err := middleware.authService.Authorize(bearer[1])
		if err != nil {
			c.AbortWithError(403, err)
		}
		c.Set("authorizedUser", user)
	}
}

/*
GetAuthorizedUser -> return the authorized user from the given gin context
*/
func (middleware AuthBearerMiddleware) GetAuthorizedUser(c *gin.Context) *models.User {
	user, exists := c.Get("authorizedUser")
	if !exists {
		panic(AuthBearerMiddlewareNotCalledError{})
	}
	return user.(*models.User)
}
