package controller

import "github.com/gin-gonic/gin"

func GetUsers() gin.HandleFunc{
		return func(c *gin.Context) {

		}
}

func GetUser() gin.HandleFunc{
	return func(c *gin.Context) {

	}
}

func SignUp() gin.HandleFunc{
	return func(c *gin.Context) {

	}
}

func LogIn() gin.HandleFunc{
	return func(c *gin.Context) {

	}
}


func HashPassword(password string) string{

}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {
	
}