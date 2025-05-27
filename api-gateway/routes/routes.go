package routes

import (
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	RegisterUserRoutes(r)
	RegisterInventoryRoutes(r)
    RegisterOrderRoutes(r)
}


// func RegisterRoutes(r *gin.Engine) {
//     // Inventory

//     // Orders
    // r.POST("/orders", services.CreateOrder)
    // r.GET("/orders/:id", services.GetOrderByID)
    // r.PATCH("/orders/:id", services.UpdateOrderStatus)
    // r.GET("/orders", services.GetUserOrders)

//     // new route
//     r.GET("/register", services.Register)

//     // Auth
//     r.POST("/login", services.Login)
// }
