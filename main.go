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
// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://localhost:9100/oauth/token
// @authorizationurl https://localhost:3000/oauth
// @scope.user:me:verify Grants access to verify created user
// @scope.user:me:change-password Grants access to change self password
// @scope.user:me:read Grants access to read self user
// @scope.user:me:write Grants access to write self user
// @scope.user:me:delete Grants access to delete self user
// @scope.user:me:authorized-app Grants access an app to get information about the user
// @scope.app:me:write Grants access to write self created apps
// @scope.app:me:read Grants access to read self created apps
func main() {
	docs.SwaggerInfo.Title = "Gandalf API"
	router := gin.Default()

	// Cors configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	config.AllowCredentials = true
	router.Use(cors.New(config))

	routes.Routes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(router.Run(fmt.Sprintf(":%s", os.Getenv("GANDALF_PORT"))))
}
