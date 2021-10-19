package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register ping endpoints to the given router
func RegisterPingRoutes(router *gin.Engine) {
	router.GET("/ping", Ping)
}

type pong struct {
	Data string `json:"data" example:"pong"`
}

// @Summary Ping
// @Description  ping the system to healthcare purposes
// @ID ping
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} pong
// @Failure 500 {object} helpers.HTTPError
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, pong{Data: "pong"})
	return
}
