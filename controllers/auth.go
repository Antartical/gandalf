package controllers

import (
	"gandalf/helpers"
	"gandalf/security"
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register auth endpoints to the given router
func RegisterAuthRoutes(router *gin.Engine, authService services.IAuthService) {
	controller := AuthController{
		authService: authService,
	}

	publicRoutes := router.Group("/auth")
	{
		publicRoutes.POST("/login", controller.Login)
		publicRoutes.POST("/refresh", controller.Refresh)
	}
}

// Controller fot /auth endpoints
type AuthController struct {
	authService services.IAuthService
}

// @Summary Login admin
// @Description Logs an user into the system
// @ID auth-login
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body validators.Credentials true "Logs into the system with the given credentials"
// @Success 200 {object} serializers.TokensSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /auth/login [post]
func (controller AuthController) Login(c *gin.Context) {
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

	tokens := controller.authService.GenerateTokens(*user, security.GroupUserSelf)
	c.JSON(http.StatusOK, serializers.NewTokensSerializer(tokens))
}

// @Summary Refresh
// @Description Refresh the given access token
// @ID auth-refresh
// @Tags Auth
// @Accept json
// @Produce json
// @Param tokens body validators.AuthTokens true "Refresh the given access token with the refresh one"
// @Success 200 {object} serializers.TokensSerializer
// @Failure 400 {object} helpers.HTTPError
// @Router /auth/refresh [post]
func (controller AuthController) Refresh(c *gin.Context) {
	var input validators.AuthTokens

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	newTokens, err := controller.authService.RefreshToken(input.AcessToken, input.RefreshToken)
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, serializers.NewTokensSerializer(*newTokens))
}
