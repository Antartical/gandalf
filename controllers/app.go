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
	"github.com/gofrs/uuid"
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

	readRoutes := router.Group("/apps")
	{
		scopes := []string{security.ScopeAppReadAll}
		readRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		readRoutes.GET(":uuid", controller.ReadApp)
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
// @Tags App
// @Accept json
// @Produce json
// @Param app body validators.AppCreateData true "Creates an app"
// @Success 201 {object} serializers.AppSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Security OAuth2AccessCode[app:me:write]
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

// @Summary Get an app
// @Description get an app by his uuid
// @ID app-read
// @Tags App
// @Accept json
// @Produce json
// @Param uuid path string true "App uuid"
// @Success 200 {object} serializers.AppSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Security OAuth2AccessCode[app:all:read]
// @Router /apps/{uuid} [get]
func (controller AppController) ReadApp(c *gin.Context) {
	var input validators.AppReadData
	if err := c.ShouldBindUri(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	uuid, _ := uuid.FromString(input.UUID)
	app, err := controller.appService.Read(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, serializers.NewAppSerializer(*app))
}
