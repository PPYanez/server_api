package routes

import (
	"fmt"
	"log"
	"net/http"
	"trust-bank/api/db"
	"trust-bank/api/messages"
	"trust-bank/api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DepositToWallet(c *gin.Context) {
	// Parse request body
	depositObject := new (models.Deposit)
	if err := c.BindJSON(&depositObject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	// Instantiate mongodb
	mongoClient := db.GetClient()

	// Check if the provided client exists
	clientColl := mongoClient.Database("trustBank").Collection("Clientes")
	foundClient := new (models.Client)
	filter := bson.D{{"numero_identificacion", depositObject.NumeroCliente}}
	err := clientColl.FindOne(c, filter).Decode(&foundClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_no_encontrado"})
		return
	}

	// Check if the provided wallet exists
	walletColl := mongoClient.Database("trustBank").Collection("Billeteras")
	foundWallet := new (models.Wallet)
	filter = bson.D{{"nro_cliente", depositObject.NumeroCliente}}
	err = walletColl.FindOne(c, filter).Decode(&foundWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_no_encontrada"})
		return
	}

	// Send deposit message
	messages.RunProducer(fmt.Sprintf("deposit %s %s", depositObject.Monto, depositObject.NumeroCliente))

	c.JSON(http.StatusOK, gin.H{"estado": "deposito_enviado"})
}