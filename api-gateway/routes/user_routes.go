package routes

import (
	"api-gateway/services"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	r.POST("/register", services.RegisterUser)
	r.POST("/authenticate", services.AuthenticateUser)
	r.GET("/profile/:user_id", services.GetUserProfile)
}
