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