package main

import (
	"os"
	"restaurant-managment-system/database"
	"restaurant-managment-system/middleware"
	"restaurant-managment-system/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


var foodCollection *mongo.Collection  = database.OpenCollection(database.Client, "food")
func main () {
		port := os.Getenv("PORT");

		if port == "" {
			port = "8800"
		}
 
		router := gin.New();
		routes.Use(gin.Logger());
		routes.UseRoutes(router)
		router.Use(middleware.Authentication())


		routes.FoodRoutes(router);
		routes.MenuRoutes(router);
		routes.TableRoutes(router);
		routes.OrderRoutes(router);
		routes.OrderItemRoutes(router);
		routes.UserRoutes(router);
		routes.InvoiceRoutes(router);
		 
		router.Run(":" + port);
}