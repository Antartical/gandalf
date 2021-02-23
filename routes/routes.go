package routes

import (
	"gandalf/connections"
	"gandalf/controllers"
	"gandalf/middlewares"
	"gandalf/services"

	"github.com/gin-gonic/gin"
)

/*
Routes -> resgister backend routes in the given router
*/
func Routes(router *gin.Engine) {
	db := connections.NewGormPostgresConnection().Connect()

	// Services
	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)
	pelipperService := services.NewPelipperService()

	// Middlewares
	authBearerMiddleware := middlewares.NewAuthBearerMiddleware(authService)

	// Routes
	controllers.RegisterPingRoutes(router)
	controllers.RegisterUserRoutes(
		router, authBearerMiddleware,
		authService, userService, pelipperService,
	)
}
