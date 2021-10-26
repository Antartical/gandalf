package controllers

import (
	"gandalf/helpers"
	"gandalf/middlewares"
	"gandalf/security"
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register me endpoints to the given router
func RegisterMeRoutes(
	router *gin.Engine,
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	authService services.IAuthService,
	userService services.IUserService,
	appService services.IAppService,
	pelipperService services.IPelipperService,
) {
	controller := MeController{
		authService:     authService,
		userService:     userService,
		pelipperService: pelipperService,
		authMiddleware:  authBearerMiddleware,
		appService:      appService,
	}

	verifyRoutes := router.Group("/me")
	{
		scopes := []string{security.ScopeUserVerify}
		verifyRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		verifyRoutes.POST("/verify", controller.VerificateMe)
	}

	readRoutes := router.Group("/me")
	{
		scopes := []string{security.ScopeUserRead}
		readRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		readRoutes.GET("", controller.ReadMe)
	}

	updateRoutes := router.Group("/me")
	{
		scopes := []string{security.ScopeUserWrite}
		updateRoutes.Use(authBearerMiddleware.HasScopes(scopes))
		updateRoutes.PATCH("", controller.UpdateMe)
	}

	deleteRoutes := router.Group("/me")
	{
		scopes := []string{security.ScopeUserDelete}
		deleteRoutes.Use(authBearerMiddleware.HasScopes(scopes))
		deleteRoutes.DELETE("", controller.DeleteMe)
	}

	changePasswordRoutes := router.Group("/me")
	{
		scopes := []string{security.ScopeUserChangePassword}
		changePasswordRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		changePasswordRoutes.POST("/reset-password", controller.ResetMyPassword)
	}

	readAppsRoutes := router.Group("/me")
	{
		scopes := []string{security.ScopeAppRead}
		readAppsRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		readAppsRoutes.GET("/apps", controller.GetMyApps)
		readAppsRoutes.GET("/connected-apps", controller.GetMyConnectedApps)
	}
}

// Controller for /me endpoints
type MeController struct {
	authService     services.IAuthService
	userService     services.IUserService
	pelipperService services.IPelipperService
	appService      services.IAppService
	authMiddleware  middlewares.IAuthBearerMiddleware
}

// @Summary Get me
// @Description get the user who performs the request
// @ID me-read
// @Tags Me
// @Accept json
// @Produce json
// @Success 200 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /me [get]
func (controller MeController) ReadMe(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	c.JSON(http.StatusOK, serializers.NewUserSerializer(*user))
}

// @Summary Update me
// @Description update me
// @ID me-update
// @Tags Me
// @Accept json
// @Produce json
// @Param user body validators.UserUpdateData true "Updates the user with the given data"
// @Success 200 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /me [patch]
func (controller MeController) UpdateMe(c *gin.Context) {
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

// @Summary Delete me
// @Description deletes the user who perform the request
// @ID me-delete
// @Tags Me
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /me [delete]
func (controller MeController) DeleteMe(c *gin.Context) {
	user := controller.authMiddleware.GetAuthorizedUser(c)
	if err := controller.userService.Delete(user.UUID); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

// @Summary Verify me
// @Description Verify me
// @ID me-verify
// @Tags Me
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /me/verify [post]
func (controller MeController) VerificateMe(c *gin.Context) {
	controller.userService.Verificate(controller.authMiddleware.GetAuthorizedUser(c))
	c.JSON(http.StatusNoContent, nil)
}

// @Summary Reset my password
// @Description Reset my password
// @ID me-reset-password
// @Tags Me
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /me/reset-password [post]
func (controller MeController) ResetMyPassword(c *gin.Context) {
	var input validators.UserResetPasswordData
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	controller.userService.ResetPassword(controller.authMiddleware.GetAuthorizedUser(c), input.Password)
	c.JSON(http.StatusNoContent, nil)
}

// @Summary Get user's created apps
// @Description Get user's created apps
// @ID me-apps
// @Tags Me
// @Accept json
// @Produce json
// @Param page query int false "cursor's page"
// @Param limit query int false "cursor's limit"
// @Success 200 {object} serializers.PaginatedAppsSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /me/apps [get]
func (controller MeController) GetMyApps(c *gin.Context) {
	var input validators.PaginationQuery
	if err := c.ShouldBindQuery(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	user := controller.authMiddleware.GetAuthorizedUser(c)
	cursor := helpers.NewCursor(input.Page, input.PageSize)
	apps := controller.appService.ListApps(*user, &cursor)

	c.JSON(http.StatusOK, serializers.NewPaginatedAppsSerializer(apps, cursor))
}

// @Summary Get user's connected apps
// @Description Get user's connected apps
// @ID me-connected-apps
// @Tags Me
// @Accept json
// @Produce json
// @Param page query int false "cursor's page"
// @Param limit query int false "cursor's limit"
// @Success 200 {object} serializers.PaginatedAppsPublicSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /me/connected-apps [get]
func (controller MeController) GetMyConnectedApps(c *gin.Context) {
	var input validators.PaginationQuery
	if err := c.ShouldBindQuery(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	user := controller.authMiddleware.GetAuthorizedUser(c)
	cursor := helpers.NewCursor(input.Page, input.PageSize)
	apps := controller.appService.ListConnectedApps(*user, &cursor)

	c.JSON(http.StatusOK, serializers.NewPaginatedAppsPublicSerializer(apps, cursor))
}
