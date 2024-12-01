package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-managment-system/database"
	"restaurant-managment-system/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)



var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")


func GetOrders() gin.HandlerFunc{
		return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

			result , err := orderCollection.Find(context.TODO(), bson.M{});
			defer cancel();

			if (err != nil) {
				c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing orders"})
			}

			var allOrders []bson.M;

			if err = result.All(ctx, &allOrders) ; err != nil {
				log.Fatal(err);
			}

			c.JSON(http.StatusOK, allOrders);

		}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode( &order);
		defer cancel()

		if err!=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the food"})
		}

		c.JSON( http.StatusOK, order);
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}
