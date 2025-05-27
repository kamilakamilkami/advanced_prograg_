package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"api-gateway/routes"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)

	log.Println("API Gateway running on :8080")
	r.Run(":8080")
}
