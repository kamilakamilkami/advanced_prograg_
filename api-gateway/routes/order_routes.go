package routes

import (
	"api-gateway/services"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.Engine) {
	r.POST("/orders", services.CreateOrder)
    r.GET("/orders/:id", services.GetOrderByID)
    r.PATCH("/orders/:id", services.UpdateOrderStatus)
    r.GET("/orders", services.GetUserOrders)
}
