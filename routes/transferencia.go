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

func TransferFunds(c *gin.Context) {
	// Parse request body
	transferObject := new (models.Transfer)
	if err := c.BindJSON(&transferObject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	// Instantiate mongodb
	mongoClient := db.GetClient()

	// Check if provided clients exist
	clientColl := mongoClient.Database("trustBank").Collection("Clientes")

	foundClient := new (models.Client)
	filter := bson.D{{"numero_identificacion", transferObject.NumeroClienteOrigen}}
	err := clientColl.FindOne(c, filter).Decode(&foundClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_origen_no_encontrado"})
		return
	}

	foundClient = new (models.Client)
	filter = bson.D{{"numero_identificacion", transferObject.NumeroClienteDestino}}
	err = clientColl.FindOne(c, filter).Decode(&foundClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_destino_no_encontrado"})
		return
	}

	// Check if provided wallets exist
	walletColl := mongoClient.Database("trustBank").Collection("Billeteras")

	foundWallet := new (models.Wallet)
	filter = bson.D{{"nro_cliente", transferObject.NumeroClienteDestino}}
	err = walletColl.FindOne(c, filter).Decode(&foundWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_destino_no_encontrada"})
		return
	}

	foundWallet = new (models.Wallet)
	filter = bson.D{{"nro_cliente", transferObject.NumeroClienteOrigen}}
	err = walletColl.FindOne(c, filter).Decode(&foundWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_origen_no_encontrada"})
		return
	}

	// Check if origin wallet has sufficient funds
	montoFloat, _ := strconv.ParseFloat(transferObject.Monto, 32)
	saldoFloat, _ := strconv.ParseFloat(foundWallet.Saldo, 32)

	if(montoFloat > saldoFloat) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"estado": "billetera_origen_sin_fondos_suficientes"})
		return
	}

	// Send transfer message
	messages.RunProducer(fmt.Sprintf(
		"transfer %s %s %s", // ammount, from, to
		transferObject.Monto, 
		transferObject.NumeroClienteOrigen,
		transferObject.NumeroClienteDestino,
	))

	c.JSON(http.StatusOK, gin.H{"estado": "transferencia_enviada"})
}