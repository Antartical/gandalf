package controllers

import (
	"fmt"
	"gandalf/middlewares"
	"gandalf/security"
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
		publicRoutes.POST("/reset/resend", controller.ResendVerificationEmail)
	}

	verifyRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserVerify}
		verifyRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		verifyRoutes.PATCH("/verify", controller.VerificateUser)
	}

	readRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserRead}
		readRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		readRoutes.GET("/me", controller.Me)
	}

	changePasswordRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserChangePassword}
		changePasswordRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		changePasswordRoutes.PATCH("/reset", controller.ResetUserPassword)
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
		*user, []string{security.ScopeUserVerify},
	).AccessToken
	emailData := validators.PelipperUserVerifyEmail{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		VerificationLink: fmt.Sprintf(
			"%s?code=%s", input.VerificationURL, verifyToken,
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
		*user, []string{security.ScopeUserVerify},
	).AccessToken
	emailData := validators.PelipperUserVerifyEmail{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		VerificationLink: fmt.Sprintf(
			"%s?code=%s", input.VerificationURL, verifyToken,
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

/*
ResendResetPasswordEmail -> resend the user reset password email to the given
one
*/
func (controller UserController) ResendResetPasswordEmail(c *gin.Context) {
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

	changePasswordToken := controller.authService.GenerateTokens(
		*user, []string{security.ScopeUserChangePassword},
	).AccessToken
	emailData := validators.PelipperUserChangePassword{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		ChangePasswordLink: fmt.Sprintf(
			"%s?code=%s", input.VerificationURL, changePasswordToken,
		),
	}
	go controller.pelipperService.SendUserChangePasswordEmail(emailData)
	c.JSON(http.StatusCreated, nil)
}

/*
ResetUserPassword -> reset the password ftom the user who perform the request
*/
func (controller UserController) ResetUserPassword(c *gin.Context) {
	var input validators.UserResetPasswordData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.userService.ResetPassword(controller.authMiddleware.GetAuthorizedUser(c), input.Password)
	c.JSON(http.StatusOK, nil)
}

/*
Me -> return the user data who perform the request
*/
func (controller UserController) Me(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	c.JSON(http.StatusOK, serializers.NewUserSerializer(*user))
}
