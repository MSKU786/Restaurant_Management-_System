package routes

import (
	controller "restaurant-managment-system/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemsRoutes(incomingRoutes *gin.Engine) {
	 incomingRoutes.GET("/orderItems", controller.GetorderItems())
	 incomingRoutes.GET("/orderItems/:orderItem_id", controller.GetOrderItem())
	 incomingRoutes.POST("/orderItems", controller.CreateOrderItems())
	 incomingRoutes.PATCH("/orderItems/:orderItem_id", controller.UpdateOrderItems())
}