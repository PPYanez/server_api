package routes

import (
	"context"
	"log"
	"net/http"
	"trust-bank/api/db"
	"trust-bank/api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClientInfo(c *gin.Context) {
	client := db.GetClient()

	numero_identificacion := c.Query("numero_identificacion")

	foundClient := new(models.Client)
	opts := options.FindOne().SetProjection(bson.D{{"contrasena", 0}})
	coll := client.Database("trustBank").Collection("Clientes")
	err := coll.FindOne(c, bson.D{{"numero_identificacion", numero_identificacion}}, opts).Decode(&foundClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_no_encontrado"})
	}

	c.JSON(http.StatusOK, foundClient)
}

func Login(c *gin.Context) {
	client := db.GetClient()
	usuarioLog := new(models.Client)

	if err := c.BindJSON(&usuarioLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	coll := client.Database("trustBank").Collection("Clientes")
	var clientModel models.Client
	c.BindJSON(&clientModel)
	var result models.Client
	err := coll.FindOne(context.TODO(), bson.M{
		"numero_identificacion": usuarioLog.NumeroIdentificacion,
		"contrasena":            usuarioLog.Contrasena,
	}).Decode(&result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"estado": " no_exitoso"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"estado": "exitoso"})
}
