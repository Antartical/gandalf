package controllers

import (
	"fmt"
	"gandalf/middlewares"
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

/*
RegisterUserRoutes -> register user endpoints to the given router
*/
func RegisterUserRoutes(
	router *gin.Engine,
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	authService services.IAuthService,
	userService services.IUserService,
	pelipperService services.IPelipperService,
) {
	controller := UserController{
		authService:     authService,
		userService:     userService,
		pelipperService: pelipperService,
		authMiddleware:  authBearerMiddleware,
	}

	publicRoutes := router.Group("/users")
	{
		publicRoutes.POST("", controller.CreateUser)
		publicRoutes.POST("/verify/resend", controller.ResendVerificationEmail)
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
	authService     services.IAuthService
	userService     services.IUserService
	pelipperService services.IPelipperService
	authMiddleware  middlewares.IAuthBearerMiddleware
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

	verifyToken := controller.authService.GenerateTokens(
		*user, []string{services.ScopeUserVerify},
	).AccessToken
	emailData := validators.PelipperUserVerifyEmail{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		VerificationLink: fmt.Sprintf(
			"%s?code=%s", os.Getenv("EMAIL_VERIFICATION_URL"), verifyToken,
		),
	}

	go controller.pelipperService.SendUserVerifyEmail(emailData)
	c.JSON(http.StatusCreated, serializers.NewUserSerializer(*user))
}

/*
ResendVerificationEmail -> resend the user verification email to the given
one
*/
func (controller UserController) ResendVerificationEmail(c *gin.Context) {
	var input validators.UserResendEmail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := controller.userService.ReadByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusCreated, nil)
		return
	}

	verifyToken := controller.authService.GenerateTokens(
		*user, []string{services.ScopeUserVerify},
	).AccessToken
	emailData := validators.PelipperUserVerifyEmail{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		VerificationLink: fmt.Sprintf(
			"%s?code=%s", os.Getenv("EMAIL_VERIFICATION_URL"), verifyToken,
		),
	}
	go controller.pelipperService.SendUserVerifyEmail(emailData)
	c.JSON(http.StatusCreated, nil)
}

/*
VerificateUser -> verificates the user who perform the request
*/
func (controller UserController) VerificateUser(c *gin.Context) {
	controller.userService.Verificate(controller.authMiddleware.GetAuthorizedUser(c))
	c.JSON(http.StatusOK, nil)
}
