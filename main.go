package main

import (
	"fmt"
	"log"
	"os"

	routes "gandalf/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.Routes(router)
	log.Fatal(router.Run(fmt.Sprintf(":%s", os.Getenv("GANDALF_PORT"))))
}
