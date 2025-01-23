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
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
			tableId := c.Param("table_id");

			var table model.Table;

			err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table);
			defer cancel()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching the table"})
			}

			c.JSON(http.StatusOK, table);
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
