package main

import (
	"log"

	"github.com/gin-gonic/gin"

	ports "pay-service/internal/ports"
	pay "pay-service/internal/services"
)

func main() {
	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	service := pay.NewService()
	handler := ports.NewHandler(service)

	api := server.Group("/api/v1")
	ports.Routers(api, handler)

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
