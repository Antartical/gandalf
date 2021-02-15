package controllers

import (
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
RegisterUserRoutes -> register user endpoints to the given router
*/
func RegisterUserRoutes(router *gin.Engine, userService services.IUserService) {
	users := router.Group("/users")
	users.Use(func(c *gin.Context) {
		c.Set("userService", userService)
	})
	users.POST("", CreateUser)
}

/*
CreateUser -> creates a new user
*/
func CreateUser(c *gin.Context) {
	userService := c.MustGet("userService").(services.IUserService)

	var input validators.UserCreateData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userService.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serializers.NewUserSerializer(user))
}
