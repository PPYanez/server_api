package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"trust-bank/api/db"
	"trust-bank/api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetClients(c *gin.Context) {
	client := db.GetClient()
	coll := client.Database("trustBank").Collection("Clientes")
	cur, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var clients []models.Client

	for cur.Next(context.TODO()) {
		var client models.Client
		err := cur.Decode(&client)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(client)
		clients = append(clients, client)
	}

	c.JSON(http.StatusOK, gin.H{"clients": clients})
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
