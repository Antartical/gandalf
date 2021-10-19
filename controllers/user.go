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

// Register user endpoints to the given router
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
		publicRoutes.POST("/email/verify/resend", controller.ResendVerificationEmail)
		publicRoutes.POST("/email/reset-password/resend", controller.ResendResetPasswordEmail)
	}

	verifyRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserVerify}
		verifyRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		verifyRoutes.POST("/me/verify", controller.VerificateUser)
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

		changePasswordRoutes.POST("/me/reset-password", controller.ResetUserPassword)
	}
}

// Controller for /users endpoints
type UserController struct {
	authService     services.IAuthService
	userService     services.IUserService
	pelipperService services.IPelipperService
	authMiddleware  middlewares.IAuthBearerMiddleware
}

// @Summary Create User
// @Description Creates a new user
// @ID user-create
// @Tags User
// @Accept json
// @Produce json
// @Param user body validators.UserCreateData true "Creates a new user"
// @Success 201 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
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

// @Summary Update user
// @Description updates an user
// @ID user-update
// @Tags User
// @Accept json
// @Produce json
// @Param user body validators.UserUpdateData true "Updates the user with the given data"
// @Success 200 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/me [patch]
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

// @Summary Get user
// @Description get an user by his uuid
// @ID user-read-uuid
// @Tags User
// @Accept json
// @Produce json
// @Param uuid path string true "User uuid"
// @Success 200 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/{uuid} [get]
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

// @Summary Get me
// @Description get the user who performs the request
// @ID user-read-me
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/me [get]
func (controller UserController) Me(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	c.JSON(http.StatusOK, serializers.NewUserSerializer(*user))
}

// @Summary Delete me
// @Description deletes the user who perform the request
// @ID user-delete-me
// @Tags User
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/me [delete]
func (controller UserController) DeleteUser(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	if err := controller.userService.Delete(user.UUID); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

// @Summary Verify user
// @Description Verify an user
// @ID user-verify
// @Tags User
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/me/verify [post]
func (controller UserController) VerificateUser(c *gin.Context) {
	controller.userService.Verificate(controller.authMiddleware.GetAuthorizedUser(c))
	c.JSON(http.StatusNoContent, nil)
}

// @Summary Reset user password
// @Description Reset user password
// @ID user-reset-password
// @Tags User
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/me/reset-password [post]
func (controller UserController) ResetUserPassword(c *gin.Context) {
	var input validators.UserResetPasswordData
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	controller.userService.ResetPassword(controller.authMiddleware.GetAuthorizedUser(c), input.Password)
	c.JSON(http.StatusNoContent, nil)
}

// @Summary Resend verification email
// @Description Resend verification email
// @ID user-resend-verification-email
// @Tags Notification
// @Accept json
// @Produce json
// @Param data body validators.UserResendEmail true "resen the verification email"
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/email/verify/resend [post]
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
	c.JSON(http.StatusNoContent, nil)
}

// @Summary Resend reset password email
// @Description Resend reset password email
// @ID user-resend-reset-password-email
// @Tags Notification
// @Accept json
// @Produce json
// @Param data body validators.UserResendEmail true "resend the reset password email"
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/email/reset-password/resend [post]
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
	c.JSON(http.StatusNoContent, nil)
}
