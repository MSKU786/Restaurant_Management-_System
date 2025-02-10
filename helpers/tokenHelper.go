package helpers

import (
	"log"
	"os"
	"restaurant-managment-system/database"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
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


func UpdateAllToken() {

}


func ValidateToken() {

}


