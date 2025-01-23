package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-managment-system/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table");

func GetTables() gin.HandlerFunc{
		return func(c *gin.Context) {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second );

				result, err := tableCollection.Find(context.TODO(), bson.M{});
				defer cancel()();

				if (err != nil) {
					c.JSON(https.StatusInternalServerError, gin.H{"error": "error occured while fetching the tables"})
				}

				var allTables []bson.M;

				if err = result.All(ctx, &allTables) : err != nil {
					log.Fatal(err);
				}
				c.JSON(http.StatusOK, allTables);
		}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}
