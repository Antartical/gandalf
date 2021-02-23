package controllers

import (
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
RegisterAuthRoutes -> register auth endpoints to the given router
*/
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

/*
AuthController -> controller fot /auth endpoints
*/
type AuthController struct {
	authService services.IAuthService
}

/*
Login -> logs the given user into the system
*/
func (controller AuthController) Login(c *gin.Context) {
	var input validators.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := controller.authService.Authenticate(input)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	tokens := controller.authService.GenerateTokens(*user, input.Scopes)
	c.JSON(http.StatusOK, serializers.NewTokensSerializer(tokens))
}

/*
Refresh -> refresh the accessing token
*/
func (controller AuthController) Refresh(c *gin.Context) {
	var input validators.AuthTokens

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTokens, err := controller.authService.RefreshToken(input.AcessToken, input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, serializers.NewTokensSerializer(*newTokens))
}
