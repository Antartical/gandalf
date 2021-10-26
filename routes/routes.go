package routes

import (
	"gandalf/connections"
	"gandalf/controllers"
	"gandalf/middlewares"
	"gandalf/services"

	"github.com/gin-gonic/gin"
)

// Register backend routes and controllers into the
// given router
func Routes(router *gin.Engine) {
	db := connections.NewGormPostgresConnection().Connect()

	// Services
	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)
	appService := services.NewAppService(db)
	pelipperService := services.NewPelipperService()

	// Middlewares
	authBearerMiddleware := middlewares.NewAuthBearerMiddleware(authService)

	// Routes
	controllers.RegisterAuthRoutes(router, authService)
	controllers.RegisterNotificationRoutes(
		router, authService,
		userService, pelipperService,
	)
	controllers.RegisterPingRoutes(router)
	controllers.RegisterUserRoutes(
		router, authBearerMiddleware,
		authService, userService, pelipperService,
	)
	controllers.RegisterOauth2Routes(
		router, authBearerMiddleware,
		authService, userService, appService,
	)
}
