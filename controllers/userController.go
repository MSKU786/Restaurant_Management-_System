package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-managment-system/database"
	"restaurant-managment-system/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

				startIndex := (page-1) * recordsPerPage;
				startIndex, err = strconv.Atoi(c.Query("startIndex"));

				matchStage := bson.D{{Key: "$match", Value: bson.D{}}};
				projectStage := bson.D{
					{Key: "$project", Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "total_count", Value: 1},
						{Key: "user_items", Value: bson.D{
							{Key: "$slice", Value: []interface{}{"$data", startIndex, recordsPerPage},
						}}},
					}}}

				
				result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage});
				defer cancel();

				var allUsers []bson.M;

				// result ,err := userCollection.Find(context.TODO(), bson.M{})

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
			var newUser models.User;

			if err := c.BindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});
				return;
			}

			// Validate the data based on struct defined
			validateErr := validate.Struct(newUser);
			if validateErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()});
				return;
			}

			// Yow will the check if the email has already been used by another user
			//First check if the email already exists
			emailCount, err := userCollection.CountDocuments(context.Background(), bson.M{"email": newUser.Email});
			defer cancel();
			if err != nil {
				log.Panic(err);
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"});
				return;
			}

			// Hash Password
			password := HashPassword(*newUser.Password);
			newUser.Password = &password;

			if emailCount > 0 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Email alread exists"});
				return;
			}

			// need to check if phone number is already registerd by other users
			phoneCount, err := userCollection.CountDocuments(context.Background(), bson.M{"phone": newUser.Phone});
			defer cancel();
			if err != nil {
				log.Panic(err);
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"});
				return;
			}

			if phoneCount > 0 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Phone number already exists"});
				return;
			}
			

			// Create some extra fields
			newUser.Created_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123));
			newUser.Updated_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123));
			newUser.ID = primitive.NewObjectID();
			newUser.User_id = newUser.ID.Hex();
			

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}


func HashPassword(password string) string{

}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {

}