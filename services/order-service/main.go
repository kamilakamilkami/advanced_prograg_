package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"log"

	// "github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"proto/orderpb"
	grpc_hand "order-service/delivery/grpc"
	"net"
	"fmt"
	// "order-service/delivery/http"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	"order-service/nats"
	"order-service/internal/redis"
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
	db := client.Database("orders_db")
	orderCollection := db.Collection("orders")
	redisClient := redis.NewClient()
    repo := repository.NewOrderRepository(orderCollection.Database())
	publisher := nats.NewPublisher()
    usecase := usecase.NewOrderUsecase(repo, publisher, redisClient)

	grpcServer := grpc.NewServer()

	orderServiceServer := grpc_hand.NewOrderServer(usecase)

	orderpb.RegisterOrderServiceServer(grpcServer, orderServiceServer)

	lis, err := net.Listen("tcp", ":5002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Inventary Service running on :5002")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}