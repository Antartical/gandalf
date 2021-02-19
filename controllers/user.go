package controllers

import (
	"gandalf/middlewares"
	"gandalf/models"
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
RegisterUserRoutes -> register user endpoints to the given router
*/
func RegisterUserRoutes(
	router *gin.Engine,
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	userService services.IUserService,
) {
	routes := router.Group("/users")
	controller := UserController{
		userService:       userService,
		getAuthorizedUser: authBearerMiddleware.GetAuthorizedUser,
	}

	routes.POST("", controller.CreateUser)
}

/*
UserController -> controller fot /users endpoints
*/
type UserController struct {
	userService       services.IUserService
	getAuthorizedUser func(c *gin.Context) *models.User
}

/*
CreateUser -> creates a new user
*/
func (controller UserController) CreateUser(c *gin.Context) {

	var input validators.UserCreateData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := controller.userService.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serializers.NewUserSerializer(*user))
}
