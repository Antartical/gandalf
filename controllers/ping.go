package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
RegisterPingRoutes -> register ping endpoints to the given router
*/
func RegisterPingRoutes(router *gin.Engine) {
	router.GET("/ping", Ping)
}

/*
Ping -> handler for /ping route
*/
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "pong",
	})
	return
}
