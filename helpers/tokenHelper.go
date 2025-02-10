package helpers

import (
	"os"
	"restaurant-managment-system/database"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type SingerDetails struct {
	Email string
	First_name string
	Last_name string
	Uid string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")
func GenerateAllTokens() {

}


func UpdateAllToken() {

}


func ValidateToken() {

}


