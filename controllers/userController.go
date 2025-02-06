package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-managment-system/database"
	"restaurant-managment-system/helpers"
	"restaurant-managment-system/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

			// Generate token and refresh token
			token, refreshToken, _ := helpers.GenerateAllTokens(*newUser.Email, *newUser.First_name, *newUser.Last_name, *newUser.User_id)
			newUser.Token = &token;
			newUser.Refresh_token = &refreshToken;

			//If all okay then insert the data into the database
			result, insertErr := userCollection.InsertOne(ctx, newUser);

			if insertErr != nil {
				msg := fmt.Sprintf("User is not created")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg});
				return;
			}

			defer cancel();

			//return status result and return okay
			c.JSON(http.StatusOK, result);
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second);

			var user models.User;
			var foundUser models.User;

			if err := c.Bind(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()});
				return;
			}

			// Find the user with the email
			err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser);

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email"});
				return;
			}

			// Verify the password
			passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password);
			if (!passwordIsValid) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg});
				return;
			}

			// Generate Token and refresh token
			token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *&foundUser.First_name, *&foundUser.Last_name, *foundUser.User_id);

			// update the token and refersh token
			helpers.UpdateAllToken(token, refreshToken, foundUser.User_id);

			defer cancel();

			c.JSON(http.StatusOK, gin.H{"message": "User Logged in successfully"});
	}
}


func HashPassword(password string) string{
		bytes, err := bcrypt.GenerateFromPassword([]byte(password, 14));
		if err != nil {
			log.Panic(err);
		}

		return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {
		err := bcrypt.CompareHashAndPassword([[] byte(providePassword), []byte(userPassword)]);
		check := true
		msg := ""

		if err != nil {
			msg = fmt.Sprintf("Invalid password");
			check := false;
		}

		return check, msg;
}