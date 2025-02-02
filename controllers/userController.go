package controller

import (
	"context"
	"net/http"
	"restaurant-managment-system/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user");


func GetUsers() gin.HandlerFunc{
		return func(c *gin.Context) {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);

				var allUsers []bson.M;

				result ,err := userCollection.Find(context.TODO(), bson.M{})

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error":"Error occured while listing users"})
					return;
				}


				if err = result.All(ctx, &allUsers); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding users"});
					return;
				}

				c.JSON(http.StatusOK, allUsers);
		}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func LogIn() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}


func HashPassword(password string) string{

}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {

}