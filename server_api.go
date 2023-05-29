package main

import (
	"trust-bank/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/cliente", routes.GetClientInfo)

	router.POST("/api/deposito", routes.DepositToWallet)

	router.POST("/api/giro", routes.WithdrawFromWallet)
	
	router.POST("/api/inicio_sesion", routes.Login)
	
	router.POST("/api/transferencia", routes.TransferFunds)

	router.Run() // listen and serve on 0.0.0.0:8080
}
