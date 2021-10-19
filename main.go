package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	docs "gandalf/docs"
	routes "gandalf/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Gandalf API
// @version 1.0
// @description Oauth2 server.
// @host localhost:9100
// @x-extension-openapi {"example": "value on a json format"}
// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://localhost:9100/auth/login
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information
// @query.collection.format multi
func main() {
	docs.SwaggerInfo.Title = "Gandalf API"
	router := gin.Default()

	// Cors configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	router.Use(cors.New(config))

	routes.Routes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(router.Run(fmt.Sprintf(":%s", os.Getenv("GANDALF_PORT"))))
}
