package main

import (
	"fmt"
	"log"
	"net"

	"user-service/internal/repository"
	
	grpc_hand "user-service/delivery/grpc"

	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"

	"user-service/internal/usecase"
	"proto/userpb"
	
	"google.golang.org/grpc"
	"user-service/internal/redis"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ Failed to ping MongoDB:", err)
	} else {
		log.Println("✅ Successfully connected to MongoDB")
	}

	db := client.Database("user_db")
	
	productCollection := db.Collection("users")
	redisClient := redis.NewClient()
	userRepo := repository.NewUserRepository(productCollection.Database())

	userUsecase := usecase.NewUserUsecase(userRepo, redisClient)

	grpcServer := grpc.NewServer()

	userServiceServer := grpc_hand.NewUserServer(userUsecase)

	userpb.RegisterUserServiceServer(grpcServer, userServiceServer)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("UserService running on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
