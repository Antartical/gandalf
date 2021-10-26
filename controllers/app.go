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

func RegisterAppRoutes(
	router *gin.Engine,
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	appService services.IAppService,
) {

	controller := AppController{
		appService:     appService,
		authMiddleware: authBearerMiddleware,
	}

	writeRoutes := router.Group("/apps")
	{
		scopes := []string{security.ScopeAppWrite}
		writeRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		writeRoutes.POST("", controller.CreateApp)
	}
}

// Controller for /app endpoints
type AppController struct {
	appService     services.IAppService
	authMiddleware middlewares.IAuthBearerMiddleware
}

// @Summary Creates a new app
// @Description creates an app
// @ID app-create
// @Tags Apps
// @Accept json
// @Produce json
// @Param app body validators.AppCreateData true "Creates an app"
// @Success 201 {object} serializers.AppSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /apps [post]
func (controller AppController) CreateApp(c *gin.Context) {
	var input validators.AppCreateData
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	user := controller.authMiddleware.GetAuthorizedUser(c)
	app, err := controller.appService.Create(input, *user)
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, serializers.NewAppSerializer(*app))
}
