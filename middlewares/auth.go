package middlewares

import (
	"errors"
	"gandalf/models"
	auth "gandalf/services"
	"strings"

	set "github.com/deckarep/golang-set"

	"github.com/gin-gonic/gin"
)

/*
IAuthBearerMiddleware -> interface for auth based on bearer token middleware
*/
type IAuthBearerMiddleware interface {
	HasScopes(scopes set.Set) gin.HandlerFunc
	GetAuthorizedUser(c *gin.Context) *models.User
}

/*
AuthBearerMiddleware -> auth middleware for authenticate users with
Bearer tokens
*/
type AuthBearerMiddleware struct {
	authService auth.IAuthService
}

/*
NewAuthBearerMiddleware -> creates a new auth middleware.
*/
func NewAuthBearerMiddleware(authService auth.IAuthService) AuthBearerMiddleware {
	return AuthBearerMiddleware{authService: authService}
}

/*
HasScopes -> check if the user who perform the request has the given scopes
*/
func (middleware AuthBearerMiddleware) HasScopes(scopes set.Set) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := strings.Split(c.GetHeader("Authorization"), "Bearer")
		if len(bearer) < 2 {
			c.AbortWithError(400, errors.New("Invalid authorization header"))
		}
		user, err := middleware.authService.GetAuthorizedUser(bearer[1], scopes)
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
