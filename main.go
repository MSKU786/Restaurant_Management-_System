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
		router.Use(gin.Logger());
		routes.UserRoutes(router)
		router.Use(middleware.Authentication())

		routes.MenuRoutes(router);
		routes.FoodRoutes(router);
		routes.TableRoutes(router);
		routes.OrderRoutes(router);
		routes.OrderItemsRoutes(router);
		routes.InvoiceRoutes(router);
		 
		router.Run(":" + port);
}