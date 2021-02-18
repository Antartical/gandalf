package routes

import (
	"gandalf/connections"
	"gandalf/controllers"
	"gandalf/services"

	"github.com/gin-gonic/gin"
)

/*
Routes -> resgister backend routes in the given router
*/
func Routes(router *gin.Engine) {
	db := connections.NewGormPostgresConnection().Connect()

	controllers.RegisterPingRoutes(router)
	controllers.RegisterUserRoutes(router, services.NewUserService(db))
}
