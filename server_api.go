package main

import (
	"trust-bank/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/cliente", routes.GetClients)
	router.POST("/api/inicio_sesion", routes.Login)
	router.Run() // listen and serve on 0.0.0.0:8080
}
