package controllers

import (
	"fmt"
	"gandalf/helpers"
	"gandalf/security"
	"gandalf/services"
	"gandalf/validators"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Register notifications endpoints to the given router
func RegisterNotificationRoutes(
	router *gin.Engine,
	authService services.IAuthService,
	userService services.IUserService,
	pelipperService services.IPelipperService,
) {
	controller := NotificationController{
		authService:     authService,
		userService:     userService,
		pelipperService: pelipperService,
	}

	publicRoutes := router.Group("/notifications/emails")
	{
		publicRoutes.POST("/verify", controller.UserVerificationEmail)
		publicRoutes.POST("/reset-password", controller.UserResetPasswordEmail)
	}

}

// Controller for /notifications endpoints
type NotificationController struct {
	authService     services.IAuthService
	userService     services.IUserService
	pelipperService services.IPelipperService
}

// @Summary Sends verification email
// @Description Sends verification email
// @ID notifications-emails-verification
// @Tags Notification
// @Accept json
// @Produce json
// @Param data body validators.UserResendEmail true "sends the verification email"
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /notifications/emails/verify [post]
func (controller NotificationController) UserVerificationEmail(c *gin.Context) {
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

// @Summary Sends reset password email
// @Description Sends reset password email
// @ID notifications-emails-reset-password
// @Tags Notification
// @Accept json
// @Produce json
// @Param data body validators.UserResendEmail true "resend the reset password email"
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /notifications/emails/reset-password [post]
func (controller NotificationController) UserResetPasswordEmail(c *gin.Context) {
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
