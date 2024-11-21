package routes

import (
	controller "restaurant-managment-system/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(incomingRoutes *gin.Engine) {
	 incomingRoutes.GET("/Tables", controller.GetTables())
	 incomingRoutes.GET("/Tables/:table_id", controller.GetTable())
	 incomingRoutes.POST("/Tables", controller.CreateTable())
	 incomingRoutes.PATCH("/Tables/:table_id", controller.UpdateTable())
}