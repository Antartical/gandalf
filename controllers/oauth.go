package controllers

import (
	"fmt"
	"gandalf/helpers"
	"gandalf/middlewares"
	"gandalf/security"
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// Register oauth2 endpoints to the given router
func RegisterOauth2Routes(
	router *gin.Engine,
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	authService services.IAuthService,
	userService services.IUserService,
	appService services.IAppService,
) {
	controller := Oauth2Controller{
		authService:    authService,
		userService:    userService,
		appService:     appService,
		authMiddleware: authBearerMiddleware,
	}

	publicRoutes := router.Group("/oauth")
	{
		publicRoutes.POST("/login", controller.Oauth2Login)
	}

	authorizeRoutes := router.Group("/oauth")
	{
		scopes := []string{security.ScopeUserAuthorizeApp}
		authorizeRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		authorizeRoutes.POST("/authorize")
	}
}

// Controller for /oauth2 endpoints
type Oauth2Controller struct {
	authService    services.IAuthService
	appService     services.IAppService
	userService    services.IUserService
	authMiddleware middlewares.IAuthBearerMiddleware
}

func (controller Oauth2Controller) Oauth2Login(c *gin.Context) {
	var input validators.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	user, err := controller.authService.Authenticate(input, false)
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusForbidden, err)
		return
	}

	tokens := controller.authService.GenerateTokens(*user, security.GroupUserOauth2Request)
	c.JSON(http.StatusOK, serializers.NewTokensSerializer(tokens))
}

func (controller Oauth2Controller) Oauth2Authorize(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	var input validators.OauthAuthorizeData
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	clientID, _ := uuid.FromString(input.ClientID)
	app, err := controller.appService.ReadByClientID(clientID)
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	controller.authService.Authorize(app, user, input)

	redirectUrl := fmt.Sprintf("%s?code=%s&state=%s", input.RedirectURI, "", input.State)
	c.Redirect(http.StatusFound, redirectUrl)
}
