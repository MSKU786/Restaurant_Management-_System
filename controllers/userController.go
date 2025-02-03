package controller

import (
	"context"
	"net/http"
	"restaurant-managment-system/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user");


func GetUsers() gin.HandlerFunc{
		return func(c *gin.Context) {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);


				recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage"));
				if err != nil || recordsPerPage < 1 {
					recordsPerPage = 10;
				}

				page, err1 := strconv.Atoi(c.Query("page"));
				if err1 != nil || page < 1{
					page = 1;
				}

				startIndex = (page-1) * recordsPerPage;
				startIndex, err = strconv.Atoi(c.Query("startIndex"));

				matchStage := bson.D{{"$match", bson.D{}}};
				
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
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);

			//convert the JSON data coming from postman to what golang can undertand
			var newUser User;

			if err := c.BindJSON(&newUser); err !=- nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});
				return;
			}

			// Validate the data based on struct defined
			if newUser.Email == "" {

			}

			//First check if the email already exists
			emailCount, err := userCollection.CountDocuments(context.Background(), bson.M{"email": newUser.Email});
			

			// Hash the password


			// need to check if phone number is already registerd by other users
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