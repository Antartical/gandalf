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
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// RegisterUserRoutes -> register user endpoints to the given router
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
		publicRoutes.POST("/reset/resend", controller.ResendResetPasswordEmail)
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
		readRoutes.GET(":uuid", controller.ReadUser)
	}

	updateRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserWrite}
		updateRoutes.Use(authBearerMiddleware.HasScopes(scopes))
		updateRoutes.PATCH("/me", controller.UpdateUser)
	}

	deleteRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserDelete}
		deleteRoutes.Use(authBearerMiddleware.HasScopes(scopes))
		deleteRoutes.DELETE("/me", controller.DeleteUser)
	}

	changePasswordRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserChangePassword}
		changePasswordRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		changePasswordRoutes.PATCH("/reset", controller.ResetUserPassword)
	}
}

// UserController -> controller fot /users endpoints
type UserController struct {
	authService     services.IAuthService
	userService     services.IUserService
	pelipperService services.IPelipperService
	authMiddleware  middlewares.IAuthBearerMiddleware
}

// @Summary Creates a new user
// @Description Creates a new user
// @ID create-user
// @Accept json
// @Produce json
// @Param user body validators.UserCreateData true "the user who will be created in the database"
// @Success 200 {object} serializers.UserSerializer
// @Failure 400 {object} serializers.HTTPErrorSerializer
// @Router /users [post]
func (controller UserController) CreateUser(c *gin.Context) {
	var input validators.UserCreateData
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userService.Create(input)
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	verifyToken := controller.authService.GenerateTokens(
		*user, []string{security.ScopeUserVerify},
	).AccessToken

	url := os.Getenv("EMAIL_VERIFICATION_URL")
	emailData := validators.PelipperUserVerifyEmail{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		VerificationLink: fmt.Sprintf(
			"%s?code=%s", url, verifyToken,
		),
	}

	go controller.pelipperService.SendUserVerifyEmail(emailData)
	c.JSON(http.StatusCreated, serializers.NewUserSerializer(*user))
}

/*
UpdateUser -> updates the user who perform the request
*/
func (controller UserController) UpdateUser(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)

	var input validators.UserUpdateData
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userService.Update(user.UUID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, serializers.NewUserSerializer(*user))
}

/*
ReadUser -> read an user by his UUID
*/
func (controller UserController) ReadUser(c *gin.Context) {
	var input validators.UserReadData
	if err := c.ShouldBindUri(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	uuid, _ := uuid.FromString(input.UUID)
	user, err := controller.userService.Read(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, serializers.NewUserSerializer(*user))
}

/*
Me -> return the user data who perform the request
*/
func (controller UserController) Me(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	c.JSON(http.StatusOK, serializers.NewUserSerializer(*user))
}

/*
Delente -> deletes the user who performs the request
*/
func (controller UserController) DeleteUser(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	if err := controller.userService.Delete(user.UUID); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

/*
ResendVerificationEmail -> resend the user verification email to the given
one
*/
func (controller UserController) ResendVerificationEmail(c *gin.Context) {
	var input validators.UserResendEmail
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
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
	url := os.Getenv("EMAIL_VERIFICATION_URL")
	emailData := validators.PelipperUserVerifyEmail{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		VerificationLink: fmt.Sprintf(
			"%s?code=%s", url, verifyToken,
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
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
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
	url := os.Getenv("PASSWORD_CHANGE_URL")
	emailData := validators.PelipperUserChangePassword{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		ChangePasswordLink: fmt.Sprintf(
			"%s?code=%s", url, changePasswordToken,
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
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	controller.userService.ResetPassword(controller.authMiddleware.GetAuthorizedUser(c), input.Password)
	c.JSON(http.StatusOK, nil)
}
