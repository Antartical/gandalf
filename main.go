package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	routes "gandalf/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Cors configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	router.Use(cors.New(config))

	routes.Routes(router)
	log.Fatal(router.Run(fmt.Sprintf(":%s", os.Getenv("GANDALF_PORT"))))
}
