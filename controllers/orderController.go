package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-managment-system/database"
	"restaurant-managment-system/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		validationErr := 	validate.Struct(order)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return 
		}

		if order.Table_id != nil  {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel();
			if err != nil {
				msg := fmt.Sprintf("message: table id not found");
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return;
			}
		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
		order.ID = primitive.NewObjectID();
		order.Order_id = order.ID.Hex()
		result, insertErr := orderCollection.InsertOne(ctx, order);
		
		if insertErr !=nil {
			msg := fmt.Sprintf("order is not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return;
		}

		defer cancel()
		c.JSON( http.StatusOK, result);
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
		var order models.Order
		var table models.Table

		var updatedObj primitive.D
		order_id := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		}



		if order.Table_id != nil  {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel();
			if err != nil {
				msg := fmt.Sprintf("message: table id not found");
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return;
			}

			updatedObj = append(updatedObj, bson.E{Key: "table_id", Value: order.Table_id});
		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: order.Updated_at})

		upsert := true;

		filter := bson.M{"order_id": order_id}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(
			ctx, 
			filter,
			bson.D{
				{Key: "$set", Value: updatedObj},
			},
			&opt,
		)

		if err !=nil {
			msg := fmt.Sprintf("Food item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return;
		}

		defer cancel()
		c.JSON( http.StatusOK, result);
	}
}


func OrderItemOrderCreator(order models.Order) string {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
		
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
	order.ID = primitive.NewObjectID();
	order.Order_id = order.ID.Hex()
	orderCollection.InsertOne(ctx, order);

	defer cancel()
	return order.Order_id;
}