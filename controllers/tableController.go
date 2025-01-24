package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-managment-system/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);	
			var table model.Table;

			if err := c.BindJSON(&table); err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return;
			}

			validationError := validate.Struct(table);

			if (validationError != nil) {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
				return;
			}

			table.created_at , _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
			table.updated.at , _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
			table.ID = primitive.NewObjectID();
			table.Table_id = table.ID.Hex();

			result, insertErr := tableCollection.InsertOne(ctx, table);

			if insertErr != nil {
				msg: fmt.Sprintf("Table item was not created")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return;
			}

			defer cancel();
			c.JSON(http.StatusOK, table);


	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}
