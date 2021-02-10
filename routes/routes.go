package routes

import (
	"gandalf/controllers"

	"github.com/gin-gonic/gin"
)

/*
Routes -> resgister backend routes in the given router
*/
func Routes(router *gin.Engine) {
	controllers.RegisterPingRoutes(router)
}
