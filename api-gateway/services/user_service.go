package services

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"proto/userpb"

	"google.golang.org/grpc"
)

var userClient userpb.UserServiceClient

func init() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}
	userClient = userpb.NewUserServiceClient(conn)
}

func PingUserService(c *gin.Context) {
	_, err := userClient.GetUserProfile(context.Background(), &userpb.UserID{UserId: ""})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user service unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "user service is running"})
}

func RegisterUser(c *gin.Context) {
	var input userpb.UserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	res, err := userClient.RegisterUser(context.Background(), &input)
	if err != nil || !res.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": res.UserId})
}

func AuthenticateUser(c *gin.Context) {
	var input userpb.AuthRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	res, err := userClient.AuthenticateUser(context.Background(), &input)
	if err != nil || !res.Success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": res.UserId})
}

func GetUserProfile(c *gin.Context) {
	userId := c.Param("user_id")
	res, err := userClient.GetUserProfile(context.Background(), &userpb.UserID{UserId: userId})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": res.UserId, "username": res.Username})
}
