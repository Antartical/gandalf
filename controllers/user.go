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
	//pelipper services.IPelipperService,
) {
	controller := UserController{
		userService: userService,
		//pelipper:          pelipper,
		getAuthorizedUser: authBearerMiddleware.GetAuthorizedUser,
	}

	publicRoutes := router.Group("/users")
	{
		publicRoutes.POST("", controller.CreateUser)
	}

	verifyRoutes := router.Group("/users")
	{
		scopes := []string{services.ScopeUserVerify}
		verifyRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		verifyRoutes.PATCH("/verify", controller.VerificateUser)
	}
}

/*
UserController -> controller fot /users endpoints
*/
type UserController struct {
	userService       services.IUserService
	pelipper          services.IPelipperService
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

	//verifyToken :=
	emailData := validators.PelipperUserVerifyEmail{
		Email:            user.Email,
		Name:             user.Name,
		Subject:          "Welcome",
		VerificationLink: "",
	}

	go controller.pelipper.SendUserVerifyEmail(emailData)
	c.JSON(http.StatusCreated, serializers.NewUserSerializer(*user))
}

/*
VerificateUser -> verificates the user who perform the request
*/
func (controller UserController) VerificateUser(c *gin.Context) {
	controller.userService.Verificate(controller.getAuthorizedUser(c))
	c.JSON(http.StatusOK, nil)
}
