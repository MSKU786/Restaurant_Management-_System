package helpers

import (
	"context"
	"log"
	"os"
	"restaurant-managment-system/database"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email string
	First_name string
	Last_name string
	Uid string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, first_name string, last_name string, uid string) (signedToken string, signedrRefreshToken string, err error) {
		claims := &SignedDetails{
			Email: email, 
			First_name: first_name, 
			Last_name: last_name,
			Uid: uid , 
			StandardClaims: jwt.StandardClaims{
					ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			},
		}

		refreshClaims := &SignedDetails{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			},
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY));

		refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY));

		if (err != nil) {
			log.Panic(err);
			return;
		}

		return token, refreshToken, err;
}


func UpdateAllToken(signedToken string, signedRefreshToken string, userId string) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);

		var updateObj primitive.D;

		updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken});
		updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken});
		
		Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at});

		upsert := true
		filter := bson.M{"user_id":userId};
		opt := options.UpdateOptions{
		 Upsert: &upsert,
		}

		_, err := userCollection.UpdateOne(
			ctx, 
			filter, 
			bson.D{
				{Key: "$set", Value: updateObj},

			},
			&opt,
		)

		defer cancel();

		if err != nil {
			log.Panic(err);
			return;
		}

		return;
}


func ValidateToken() {

}


