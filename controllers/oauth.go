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
		publicRoutes.POST("/token", controller.Oauth2Token)
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

// @Summary Login an user and retrieve auth token
// @Description logs an user
// @ID oauth-login
// @Tags Oauth
// @Accept json
// @Produce json
// @Param user body validators.Credentials true "Logs an user"
// @Success 201 {object} serializers.TokensSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /oauth/login [post]
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

// @Summary Authorize an app to get the user data
// @Description authorize app
// @ID oauth-authorize
// @Tags Oauth
// @Accept json
// @Produce json
// @Param user body validators.OauthAuthorizeData true "Authorize app to get user's data"
// @Success 302
// @Router /oauth/authorize [post]
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

// @Summary Retrieves access token form the authorization one
// @Description Retrieves access token form the authorization one
// @ID oauth-token
// @Tags Oauth
// @Accept application/x-www-form-urlencoded
// @Accept json
// @Produce json
// @Param user body validators.OauthExchangeToken true "Token exchange data"
// @Success 201 {object} serializers.TokensSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /oauth/token [post]
func (controller Oauth2Controller) Oauth2Token(c *gin.Context) {
	var input validators.OauthExchangeToken
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	app, err := controller.appService.ReadByClientID(uuid.FromStringOrNil(input.ClientID))
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
	}

	tokens, err := controller.authService.ExchangeOauthToken(*app, input)
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusUnauthorized, err)
	}
	c.JSON(http.StatusOK, serializers.NewTokensSerializer(*tokens))
}
