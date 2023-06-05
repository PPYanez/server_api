package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"trust-bank/api/db"
	"trust-bank/api/messages"
	"trust-bank/api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func WithdrawFromWallet(c *gin.Context) {
	// Parse request body
	withdrawObject := new(models.Deposit)
	if err := c.BindJSON(&withdrawObject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	// Instantiate mongodb
	mongoClient := db.GetClient()

	// Check if the provided client exists
	clientColl := mongoClient.Database("trustBank").Collection("Clientes")
	foundClient := new(models.Client)
	filter := bson.D{{"numero_identificacion", withdrawObject.NumeroCliente}}
	err := clientColl.FindOne(c, filter).Decode(&foundClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_no_encontrado"})
		return
	}

	// Check if the provided wallet exists
	walletColl := mongoClient.Database("trustBank").Collection("Billeteras")
	foundWallet := new(models.Wallet)
	filter = bson.D{{"nro_cliente", withdrawObject.NumeroCliente}}
	err = walletColl.FindOne(c, filter).Decode(&foundWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_no_encontrada"})
		return
	}
	montoFloat, _ := strconv.ParseFloat(withdrawObject.Monto, 32)
	saldoFloat, _ := strconv.ParseFloat(foundWallet.Saldo, 32)

	// Check if wallet has sufficient funds
	if montoFloat > saldoFloat {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"estado": "billetera_origen_sin_fondos_suficientes"})
		return
	}

	// Send transfer message
	messages.RunProducer(fmt.Sprintf(
		"withdraw %s %s", // ammount, from
		withdrawObject.Monto,
		withdrawObject.NumeroCliente,
	))

	c.JSON(http.StatusOK, gin.H{"estado": "giro_enviado"})
	return
}
